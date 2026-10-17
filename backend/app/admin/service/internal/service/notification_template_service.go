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

// NotificationTemplateService 通知模板管理（admin 门面，透传 core）。
type NotificationTemplateService struct {
	adminV1.NotificationTemplateServiceHTTPServer

	log        *log.Helper
	templateCli notifyV1.NotificationTemplateServiceClient
}

func NewNotificationTemplateService(ctx *bootstrap.Context, templateClient notifyV1.NotificationTemplateServiceClient) *NotificationTemplateService {
	return &NotificationTemplateService{
		log:         ctx.NewLoggerHelper("notification-template/service/admin-service"),
		templateCli: templateClient,
	}
}

func (s *NotificationTemplateService) List(ctx context.Context, req *paginationV1.PagingRequest) (*notifyV1.ListNotificationTemplateResponse, error) {
	return s.templateCli.List(ctx, req)
}

func (s *NotificationTemplateService) Get(ctx context.Context, req *notifyV1.GetNotificationTemplateRequest) (*notifyV1.NotificationTemplate, error) {
	return s.templateCli.Get(ctx, req)
}

func (s *NotificationTemplateService) Create(ctx context.Context, req *notifyV1.CreateNotificationTemplateRequest) (*notifyV1.NotificationTemplate, error) {
	return s.templateCli.Create(ctx, req)
}

func (s *NotificationTemplateService) Update(ctx context.Context, req *notifyV1.UpdateNotificationTemplateRequest) (*notifyV1.NotificationTemplate, error) {
	return s.templateCli.Update(ctx, req)
}

func (s *NotificationTemplateService) Delete(ctx context.Context, req *notifyV1.DeleteNotificationTemplateRequest) (*emptypb.Empty, error) {
	return s.templateCli.Delete(ctx, req)
}

func (s *NotificationTemplateService) Render(ctx context.Context, req *notifyV1.RenderNotificationTemplateRequest) (*notifyV1.RenderNotificationTemplateResponse, error) {
	return s.templateCli.Render(ctx, req)
}
