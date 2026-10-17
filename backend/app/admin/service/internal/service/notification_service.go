package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	notifyV1 "go-wind-cms/api/gen/go/notify/service/v1"
)

// NotificationService 通知投递台账（admin 门面，透传 core）。
// 只读视图：发送能力不开 HTTP 口（SendDirect 只由 core 进程内业务经 Notifier 调用）。
type NotificationService struct {
	adminV1.NotificationServiceHTTPServer

	log             *log.Helper
	notificationCli notifyV1.NotificationServiceClient
}

func NewNotificationService(ctx *bootstrap.Context, notificationClient notifyV1.NotificationServiceClient) *NotificationService {
	return &NotificationService{
		log:             ctx.NewLoggerHelper("notification/service/admin-service"),
		notificationCli: notificationClient,
	}
}

func (s *NotificationService) ListDelivery(ctx context.Context, req *paginationV1.PagingRequest) (*notifyV1.ListNotificationDeliveryResponse, error) {
	return s.notificationCli.ListDelivery(ctx, req)
}

func (s *NotificationService) GetDelivery(ctx context.Context, req *notifyV1.GetNotificationDeliveryRequest) (*notifyV1.NotificationDelivery, error) {
	return s.notificationCli.GetDelivery(ctx, req)
}
