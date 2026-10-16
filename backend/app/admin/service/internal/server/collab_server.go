package server

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	kratosTransport "github.com/go-kratos/kratos/v2/transport"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	authenticationV1 "go-wind-cms/api/gen/go/authentication/service/v1"
	confv1 "go-wind-cms/api/gen/go/conf/v1"
	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
	"go-wind-cms/app/admin/service/internal/collab"
)

// 默认 ws 监听地址；server.collab.addr / COLLAB_ADDR=off/disabled 停用（与 GRAPHQL_ADDR 同约定）。
const defaultCollabAddr = ":6602"

// collabConfigExtKey RegisterCollabConfigExt 登记扩展配置所用的键。
const collabConfigExtKey = "collab-config-ext"

// RegisterCollabConfigExt 把 server.collab 配置段（项目自定义 proto 扩展
// confv1.BootstrapExtension）登记进 bootstrap 的配置扫描链。
//
// 必须在 bootstrap 读配置（LoadBootstrapConfig）之前调用——扫描列表在那一刻
// 定格，晚了 yaml 段不会落进消息且无任何报错。当前在 cmd/server/main.go 的
// runApp 中调用。
func RegisterCollabConfigExt(ctx *bootstrap.Context) {
	ctx.RegisterCustomConfig(collabConfigExtKey, &confv1.BootstrapExtension{})
}

// collabExtConfig 取已登记的 server.collab 扩展配置；未登记/类型不符时
// 返回零值消息（等价"未配置"，走 env 兜底）。
func collabExtConfig(ctx *bootstrap.Context) *confv1.Collab {
	if v, ok := ctx.GetCustomConfig(collabConfigExtKey); ok {
		if ext, ok := v.(*confv1.BootstrapExtension); ok {
			return ext.GetServer().GetCollab()
		}
	}
	return &confv1.Collab{}
}

// NewCollab 构造协同服务与 ws 服务器。
//
// 配置来源：config 的 server.collab 段（confv1.Collab，项目自定义扩展消息，
// 见 api/protos/conf/v1/bootstrap_extension.proto）优先；未配置回退环境变量
// （兼容既有部署）：
//   - server.collab.addr / COLLAB_ADDR            监听地址（默认 ":6602"；
//     off/disabled 停用，停用时 server 为 nil、service 仍返回供删除踢房 no-op）
//   - server.collab.allowed_origins /
//     COLLAB_ALLOWED_ORIGINS                      允许的 ws 来源，逗号分隔
//     （默认 "http://localhost:15999"；生产必须改成实际后台来源，
//     留空回退同源校验）
func NewCollab(
	ctx *bootstrap.Context,
	authenticationServiceClient authenticationV1.AuthenticationServiceClient,
	collabDocStateServiceClient contentV1.CollabDocStateServiceClient,
	postServiceClient contentV1.PostServiceClient,
) (*collab.Service, *CollabServer) {
	collabCfg := collabExtConfig(ctx)

	addr := strings.TrimSpace(collabCfg.GetAddr())
	if addr == "" {
		addr = strings.TrimSpace(os.Getenv("COLLAB_ADDR"))
	}
	switch strings.ToLower(addr) {
	case "off", "disabled", "false", "0":
		return collab.New(ctx, authenticationServiceClient, collabDocStateServiceClient, postServiceClient, nil), nil
	case "":
		addr = defaultCollabAddr
	}

	// yaml 与 env 的 Origin 白名单归一都走 ParseAllowedOrigins（trim + 丢弃空项）
	origins := collab.ParseAllowedOrigins(strings.Join(collabCfg.GetAllowedOrigins(), ","))
	if len(origins) == 0 {
		origins = collab.ParseAllowedOrigins(os.Getenv("COLLAB_ALLOWED_ORIGINS"))
	}
	if len(origins) == 0 {
		origins = []string{"http://localhost:15999"}
	}

	svc := collab.New(ctx, authenticationServiceClient, collabDocStateServiceClient, postServiceClient, origins)

	l := ctx.NewLoggerHelper("collab-server/admin-service")
	return svc, &CollabServer{
		log:  l,
		addr: addr,
		svc:  svc,
		http: &http.Server{ReadHeaderTimeout: 10 * time.Second},
	}
}

// CollabServer kratos transport.Server 适配：包装 ygo ws handler 的独立端口监听。
type CollabServer struct {
	log  *log.Helper
	addr string
	svc  *collab.Service
	http *http.Server
	lis  net.Listener
}

func (s *CollabServer) Start(ctx context.Context) error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	s.lis = lis
	s.http.Handler = s.svc.Handler()
	s.log.Infof("collab ws server listening on: %s (rooms: /post:{id}:{lang})", lis.Addr().String())

	go func() {
		<-ctx.Done()
		_ = s.http.Close()
	}()

	if err := s.http.Serve(lis); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *CollabServer) Stop(ctx context.Context) error {
	s.log.Info("collab ws server stopping...")
	if err := s.http.Shutdown(ctx); err != nil {
		return err
	}
	if err := s.svc.Shutdown(ctx); err != nil {
		return err
	}
	s.log.Info("collab ws server stopped.")
	return nil
}

func (s *CollabServer) Endpoint() (*url.URL, error) {
	// bootstrap 在 Start 前就会取端点，返回配置地址而非要求监听已建立
	return url.Parse("ws://localhost" + s.addr)
}

var _ kratosTransport.Server = (*CollabServer)(nil)
