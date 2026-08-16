package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	interactionV1 "go-wind-cms/api/gen/go/interaction/service/v1"
)

// InteractionAdminService 是 admin BFF 对 core InteractionAdminService 的纯透传转发器。
// 鉴权与审计均在 core 侧执行（运营上下文 guard + OperationAuditLog）。
type InteractionAdminService struct {
	adminV1.InteractionAdminServiceHTTPServer

	interactionAdminServiceClient interactionV1.InteractionAdminServiceClient
	log                           *log.Helper
}

func NewInteractionAdminService(
	ctx *bootstrap.Context,
	interactionAdminServiceClient interactionV1.InteractionAdminServiceClient,
) *InteractionAdminService {
	return &InteractionAdminService{
		interactionAdminServiceClient: interactionAdminServiceClient,
		log:                           ctx.NewLoggerHelper("interaction-admin/service/admin-service"),
	}
}

func (s *InteractionAdminService) PurgeTargetInteractions(ctx context.Context, req *interactionV1.PurgeTargetInteractionsRequest) (*interactionV1.PurgeTargetInteractionsResponse, error) {
	return s.interactionAdminServiceClient.PurgeTargetInteractions(ctx, req)
}

func (s *InteractionAdminService) PurgeUserInteractions(ctx context.Context, req *interactionV1.PurgeUserInteractionsRequest) (*interactionV1.PurgeUserInteractionsResponse, error) {
	return s.interactionAdminServiceClient.PurgeUserInteractions(ctx, req)
}

func (s *InteractionAdminService) ResetCounter(ctx context.Context, req *interactionV1.ResetCounterRequest) (*interactionV1.ResetCounterResponse, error) {
	return s.interactionAdminServiceClient.ResetCounter(ctx, req)
}
