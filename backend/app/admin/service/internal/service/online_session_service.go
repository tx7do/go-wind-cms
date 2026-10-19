package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	onlineSessionV1 "go-wind-cms/api/gen/go/online_session/service/v1"

	"go-wind-cms/pkg/middleware/auth"
)

// OnlineSessionService 在线会话管理（BFF 薄代理：会话数据与吊销都在 core）。
type OnlineSessionService struct {
	adminV1.OnlineSessionServiceHTTPServer

	log             *log.Helper
	onlineSessionCli onlineSessionV1.OnlineSessionServiceClient
}

func NewOnlineSessionService(
	ctx *bootstrap.Context,
	onlineSessionCli onlineSessionV1.OnlineSessionServiceClient,
) *OnlineSessionService {
	l := log.NewHelper(log.With(ctx.GetLogger(), "module", "online-session/service/admin-service"))
	return &OnlineSessionService{
		log:              l,
		onlineSessionCli: onlineSessionCli,
	}
}

// ListOnlineSession 全量在线会话列表（透传 core）。
func (s *OnlineSessionService) ListOnlineSession(ctx context.Context, req *onlineSessionV1.ListOnlineSessionRequest) (*onlineSessionV1.ListOnlineSessionResponse, error) {
	return s.onlineSessionCli.List(ctx, req)
}

// ForceLogoutSession 强制下线指定会话（透传 core）。
func (s *OnlineSessionService) ForceLogoutSession(ctx context.Context, req *onlineSessionV1.ForceLogoutSessionRequest) (*onlineSessionV1.ForceLogoutSessionResponse, error) {
	return s.onlineSessionCli.ForceLogout(ctx, req)
}

// ListMyOnlineSession 当前登录用户的在线会话列表。
// 用户身份不跨服务透传：从认证上下文取出 userId 与当前 jti 显式传给 core。
func (s *OnlineSessionService) ListMyOnlineSession(ctx context.Context, req *onlineSessionV1.ListMyOnlineSessionRequest) (*onlineSessionV1.ListOnlineSessionResponse, error) {
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.UserId = trans.Ptr(operator.UserId)
	req.CurrentJti = operator.Jti

	return s.onlineSessionCli.ListMy(ctx, req)
}

// RevokeMyOnlineSession 当前用户强制下线自己的指定会话。
func (s *OnlineSessionService) RevokeMyOnlineSession(ctx context.Context, req *onlineSessionV1.RevokeMyOnlineSessionRequest) (*onlineSessionV1.RevokeMyOnlineSessionResponse, error) {
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	// 自撤只允许操作自己的会话：userId 一律以认证上下文为准，不信任请求体
	req.UserId = trans.Ptr(operator.UserId)

	return s.onlineSessionCli.RevokeMy(ctx, req)
}
