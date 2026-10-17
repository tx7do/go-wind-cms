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

// NotificationChannelService 通知渠道管理（admin 门面，透传 core）。
type NotificationChannelService struct {
	adminV1.NotificationChannelServiceHTTPServer

	log         *log.Helper
	channelCli  notifyV1.NotificationChannelServiceClient
}

func NewNotificationChannelService(ctx *bootstrap.Context, channelClient notifyV1.NotificationChannelServiceClient) *NotificationChannelService {
	return &NotificationChannelService{
		log:        ctx.NewLoggerHelper("notification-channel/service/admin-service"),
		channelCli: channelClient,
	}
}

func (s *NotificationChannelService) List(ctx context.Context, req *paginationV1.PagingRequest) (*notifyV1.ListNotificationChannelResponse, error) {
	return s.channelCli.List(ctx, req)
}

func (s *NotificationChannelService) Get(ctx context.Context, req *notifyV1.GetNotificationChannelRequest) (*notifyV1.NotificationChannel, error) {
	return s.channelCli.Get(ctx, req)
}

func (s *NotificationChannelService) Create(ctx context.Context, req *notifyV1.CreateNotificationChannelRequest) (*notifyV1.NotificationChannel, error) {
	return s.channelCli.Create(ctx, req)
}

func (s *NotificationChannelService) Update(ctx context.Context, req *notifyV1.UpdateNotificationChannelRequest) (*notifyV1.NotificationChannel, error) {
	return s.channelCli.Update(ctx, req)
}

func (s *NotificationChannelService) Delete(ctx context.Context, req *notifyV1.DeleteNotificationChannelRequest) (*emptypb.Empty, error) {
	return s.channelCli.Delete(ctx, req)
}

func (s *NotificationChannelService) SendTestEmail(ctx context.Context, req *notifyV1.SendTestEmailRequest) (*notifyV1.SendTestEmailResponse, error) {
	return s.channelCli.SendTestEmail(ctx, req)
}
