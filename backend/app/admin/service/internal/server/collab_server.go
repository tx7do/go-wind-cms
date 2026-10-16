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
	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
	"go-wind-cms/app/admin/service/internal/collab"
)

// 默认 ws 监听地址；COLLAB_ADDR=off/disabled 停用（与 GRAPHQL_ADDR 同约定）。
const defaultCollabAddr = ":6602"

// NewCollab 构造协同服务与 ws 服务器。
//
// 环境变量：
//   - COLLAB_ADDR            监听地址（默认 ":6602"；off/disabled 停用，
//                            停用时 server 为 nil、service 仍返回供删除踢房 no-op）
//   - COLLAB_ALLOWED_ORIGINS 允许的 ws 来源，逗号分隔
//     （默认 "http://localhost:15999"；生产必须改成实际后台来源，
//     留空回退同源校验）
func NewCollab(
	ctx *bootstrap.Context,
	authenticationServiceClient authenticationV1.AuthenticationServiceClient,
	collabDocStateServiceClient contentV1.CollabDocStateServiceClient,
	postServiceClient contentV1.PostServiceClient,
) (*collab.Service, *CollabServer) {
	addr := strings.TrimSpace(os.Getenv("COLLAB_ADDR"))
	switch strings.ToLower(addr) {
	case "off", "disabled", "false", "0":
		return collab.New(ctx, authenticationServiceClient, collabDocStateServiceClient, postServiceClient, nil), nil
	case "":
		addr = defaultCollabAddr
	}

	origins := collab.ParseAllowedOrigins(os.Getenv("COLLAB_ALLOWED_ORIGINS"))
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
