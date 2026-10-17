package service

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"

	notifyV1 "go-wind-cms/api/gen/go/notify/service/v1"

	"go-wind-cms/app/core/service/internal/data"
)

// NotificationChannelService 通知渠道管理（平台级配置）。
// 实现 EMAIL（SMTP）与 WEBHOOK 渠道：CRUD 留在本服务，"发"这件事已经交还给
// NotificationService（渠道选择策略与 SMTP/webhook 调用收在 data/channel/ 一处）。
type NotificationChannelService struct {
	notifyV1.UnimplementedNotificationChannelServiceServer

	log  *log.Helper
	repo *data.NotificationChannelRepo

	notifier Notifier
}

func NewNotificationChannelService(
	ctx *bootstrap.Context,
	repo *data.NotificationChannelRepo,
	notifier Notifier,
) *NotificationChannelService {
	return &NotificationChannelService{
		log:      ctx.NewLoggerHelper("notification-channel/service/core-service"),
		repo:     repo,
		notifier: notifier,
	}
}

func (s *NotificationChannelService) List(ctx context.Context, req *paginationV1.PagingRequest) (*notifyV1.ListNotificationChannelResponse, error) {
	if err := requireNotificationPlatformContext(ctx, "notification channels"); err != nil {
		return nil, err
	}

	return s.repo.List(ctx, req)
}

func (s *NotificationChannelService) Get(ctx context.Context, req *notifyV1.GetNotificationChannelRequest) (*notifyV1.NotificationChannel, error) {
	if req == nil || req.GetId() == 0 {
		return nil, errors.BadRequest("NOTIFICATION_CHANNEL", "id is required")
	}

	if err := requireNotificationPlatformContext(ctx, "notification channels"); err != nil {
		return nil, err
	}

	return s.repo.Get(ctx, req.GetId())
}

func (s *NotificationChannelService) Create(ctx context.Context, req *notifyV1.CreateNotificationChannelRequest) (*notifyV1.NotificationChannel, error) {
	if req == nil || req.Data == nil {
		return nil, errors.BadRequest("NOTIFICATION_CHANNEL", "invalid parameter")
	}

	if err := requireNotificationPlatformContext(ctx, "notification channels"); err != nil {
		return nil, err
	}

	operatorID, _ := viewerUserIDFromContext(ctx)
	req.Data.CreatedBy = trans.Ptr(operatorID)

	created, err := s.repo.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *NotificationChannelService) Update(ctx context.Context, req *notifyV1.UpdateNotificationChannelRequest) (*notifyV1.NotificationChannel, error) {
	if req == nil || req.Data == nil {
		return nil, errors.BadRequest("NOTIFICATION_CHANNEL", "invalid parameter")
	}

	if err := requireNotificationPlatformContext(ctx, "notification channels"); err != nil {
		return nil, err
	}

	operatorID, _ := viewerUserIDFromContext(ctx)
	req.Data.UpdatedBy = trans.Ptr(operatorID)

	return s.repo.Update(ctx, req)
}

func (s *NotificationChannelService) Delete(ctx context.Context, req *notifyV1.DeleteNotificationChannelRequest) (*emptypb.Empty, error) {
	if err := requireNotificationPlatformContext(ctx, "notification channels"); err != nil {
		return nil, err
	}

	if err := s.repo.Delete(ctx, req.GetId()); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// SendTestEmail 向指定收件人发送测试邮件，验证渠道配置是否可用。
//
// 类型/启用状态/凭据一律不在这层判断：SendDirect 显式带 channel_id 时，
// EmailSender 会做同样的校验并把原因包进 error（渠道没配/没启用/类型不对），
// 因此这里能直接把原始错误回给配置页——正是这个功能唯一有用的输出。
func (s *NotificationChannelService) SendTestEmail(ctx context.Context, req *notifyV1.SendTestEmailRequest) (*notifyV1.SendTestEmailResponse, error) {
	if req == nil || req.GetId() == 0 {
		return nil, errors.BadRequest("NOTIFICATION_CHANNEL", "id is required")
	}
	if req.GetRecipient() == "" {
		return nil, errors.BadRequest("NOTIFICATION_CHANNEL", "recipient is required")
	}

	if err := requireNotificationPlatformContext(ctx, "notification channels"); err != nil {
		return nil, err
	}

	operatorID, _ := viewerUserIDFromContext(ctx)
	id := req.GetId()
	recipient := req.GetRecipient()

	dto, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	sendReq := &notifyV1.SendDirectNotificationRequest{
		ChannelId:      trans.Ptr(id),
		Target:         recipient,
		OperatorUserId: trans.Ptr(operatorID),
	}

	switch dto.GetType() {
	case notifyV1.NotificationChannel_EMAIL:
		sendReq.EventType = notifyV1.EventType_CHANNEL_TEST_EMAIL
		sendReq.Title = fmt.Sprintf("通知渠道测试邮件（渠道 #%d）", id)
		sendReq.Content = fmt.Sprintf("这是一封测试邮件（由操作员 #%d 触发）。收到本邮件说明该通知渠道的 SMTP 配置工作正常。", operatorID)

	case notifyV1.NotificationChannel_WEBHOOK:
		return nil, errors.BadRequest("NOTIFICATION_CHANNEL", "webhook channels are tested from the notification-rules page")

	default:
		return nil, errors.BadRequest("NOTIFICATION_CHANNEL", "unsupported channel type "+dto.GetType().String())
	}

	if _, err = s.notifier.SendDirect(ctx, sendReq); err != nil {
		s.log.Errorf("send test via channel [%d] to [%s] failed: %v", id, recipient, err)
		return nil, errors.BadRequest("NOTIFICATION_CHANNEL", "send test message failed: "+err.Error())
	}

	s.log.Infof("test message sent via channel [%d] to [%s] by operator [%d]", id, recipient, operatorID)
	return &notifyV1.SendTestEmailResponse{Success: true}, nil
}
