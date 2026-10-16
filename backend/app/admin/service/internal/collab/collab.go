// Package collab 协同编辑（ygo 内嵌）：admin 服务的 ws 端点。
//
// v1 设计：房间纯内存（不落 CRDT 增量），core 内容仍是唯一事实源——
// 持久化保持编辑页既有的显式保存链路（保存按钮 → Update 全链路，
// 版本化/Webhook/索引不变），协同层只做"连接期间的无冲突合并 + 在线感知"。
// 服务重启后房间按需从客户端本地（y-indexeddb）或 core 内容重建。
package collab

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/reearth/ygo/crdt"
	"github.com/reearth/ygo/provider/websocket"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	authenticationV1 "go-wind-cms/api/gen/go/authentication/service/v1"
	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
)

// Service 协同编辑服务。
type Service struct {
	authClient authenticationV1.AuthenticationServiceClient
	log        *log.Helper

	ws *websocket.Server
}

// New 构造协同服务：upgrade 前经 ValidateToken 准入；传入 docStateClient 时
// 房间状态持久化到 core（重启可恢复，消除全员断线/重启的丢失窗口），
// 未传入则纯内存。
func New(ctx *bootstrap.Context, authClient authenticationV1.AuthenticationServiceClient, docStateClient contentV1.CollabDocStateServiceClient, allowedOrigins []string) *Service {
	s := &Service{
		authClient: authClient,
		log:        ctx.NewLoggerHelper("collab/admin-service"),
	}

	var srv *websocket.Server
	if docStateClient != nil {
		srv = websocket.NewServerWithPersistence(newGRPCPersistence(docStateClient))
	} else {
		srv = websocket.NewServer()
	}
	// 原生 y-websocket 客户端帧（非 Hocuspocus docName 前缀帧）；
	// 房间名取自 URL 路径末段（如 /post:12:zh-CN）。
	srv.AuthFunc = s.authenticate
	srv.AllowedOrigins = allowedOrigins
	srv.OnLoadDocument = func(_ context.Context, room string, _ *crdt.Doc) error {
		s.log.Infof("collab room loaded: %s", room)
		return nil
	}
	s.ws = srv
	return s
}

// Handler 返回 ws 挂载点（http.Handler）。
func (s *Service) Handler() http.Handler {
	return s.ws
}

// Shutdown 优雅停机（等待在途持久化写，当前无适配器即等待连接收尾）。
func (s *Service) Shutdown(ctx context.Context) error {
	return s.ws.Shutdown(ctx)
}

// authenticate ws upgrade 准入：query token → core ValidateToken（admin 端）。
// 浏览器 WebSocket 无法带 Authorization 头，token 走 query 参数；
// 内容级细粒度权限由保存链路（core Update）兜底，准入门槛与 SSE 同级。
func (s *Service) authenticate(r *http.Request) bool {
	token := r.URL.Query().Get("token")
	if token == "" {
		return false
	}
	resp, err := s.authClient.ValidateToken(context.Background(), &authenticationV1.ValidateTokenRequest{
		ClientType:    authenticationV1.ClientType_admin,
		Token:         token,
		TokenCategory: authenticationV1.TokenCategory_ACCESS,
	})
	if err != nil {
		s.log.Errorf("collab token authentication failed: %v", err)
		return false
	}
	if resp.GetIsBlocked() {
		s.log.Warnf("collab token is blocked")
		return false
	}
	return resp.GetIsValid()
}

// ParseAllowedOrigins 解析逗号分隔的允许来源列表（环境变量 COLLAB_ALLOWED_ORIGINS）。
func ParseAllowedOrigins(raw string) []string {
	var origins []string
	for _, o := range strings.Split(raw, ",") {
		if o = strings.TrimSpace(o); o != "" {
			origins = append(origins, o)
		}
	}
	return origins
}
