package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	notifyV1 "go-wind-cms/api/gen/go/notify/service/v1"
)

// NotificationRuleService 通知路由规则管理（admin 门面，透传 core）。
type NotificationRuleService struct {
	adminV1.NotificationRuleServiceHTTPServer

	log      *log.Helper
	ruleCli  notifyV1.NotificationRuleServiceClient
}

func NewNotificationRuleService(ctx *bootstrap.Context, ruleClient notifyV1.NotificationRuleServiceClient) *NotificationRuleService {
	return &NotificationRuleService{
		log:     ctx.NewLoggerHelper("notification-rule/service/admin-service"),
		ruleCli: ruleClient,
	}
}

func (s *NotificationRuleService) List(ctx context.Context, req *paginationV1.PagingRequest) (*notifyV1.ListNotificationRuleResponse, error) {
	return s.ruleCli.List(ctx, req)
}

func (s *NotificationRuleService) Get(ctx context.Context, req *notifyV1.GetNotificationRuleRequest) (*notifyV1.NotificationRule, error) {
	return s.ruleCli.Get(ctx, req)
}

func (s *NotificationRuleService) Create(ctx context.Context, req *notifyV1.CreateNotificationRuleRequest) (*notifyV1.NotificationRule, error) {
	return s.ruleCli.Create(ctx, req)
}

func (s *NotificationRuleService) Update(ctx context.Context, req *notifyV1.UpdateNotificationRuleRequest) (*notifyV1.NotificationRule, error) {
	return s.ruleCli.Update(ctx, req)
}

func (s *NotificationRuleService) Delete(ctx context.Context, req *notifyV1.DeleteNotificationRuleRequest) (*emptypb.Empty, error) {
	return s.ruleCli.Delete(ctx, req)
}

func (s *NotificationRuleService) TestDispatch(ctx context.Context, req *notifyV1.TestDispatchNotificationRequest) (*notifyV1.TestDispatchNotificationResponse, error) {
	return s.ruleCli.TestDispatch(ctx, req)
}
