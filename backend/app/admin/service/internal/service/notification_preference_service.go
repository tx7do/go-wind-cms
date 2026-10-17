package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	notifyV1 "go-wind-cms/api/gen/go/notify/service/v1"
)

// NotificationPreferenceService 用户通知偏好（admin 门面，透传 core；自助，读写的都是操作人自己的偏好）。
type NotificationPreferenceService struct {
	adminV1.NotificationPreferenceServiceHTTPServer

	log          *log.Helper
	preferenceCli notifyV1.NotificationPreferenceServiceClient
}

func NewNotificationPreferenceService(ctx *bootstrap.Context, preferenceClient notifyV1.NotificationPreferenceServiceClient) *NotificationPreferenceService {
	return &NotificationPreferenceService{
		log:           ctx.NewLoggerHelper("notification-preference/service/admin-service"),
		preferenceCli: preferenceClient,
	}
}

func (s *NotificationPreferenceService) GetMyNotificationPreference(ctx context.Context, req *emptypb.Empty) (*notifyV1.NotificationPreference, error) {
	return s.preferenceCli.GetMyNotificationPreference(ctx, req)
}

func (s *NotificationPreferenceService) UpdateMyNotificationPreference(ctx context.Context, req *notifyV1.UpdateNotificationPreferenceRequest) (*notifyV1.NotificationPreference, error) {
	return s.preferenceCli.UpdateMyNotificationPreference(ctx, req)
}

func (s *NotificationPreferenceService) ListMyNotificationCategories(ctx context.Context, req *emptypb.Empty) (*notifyV1.ListMyNotificationCategoriesResponse, error) {
	return s.preferenceCli.ListMyNotificationCategories(ctx, req)
}
