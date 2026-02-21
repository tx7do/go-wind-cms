package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	internalMessageV1 "go-wind-cms/api/gen/go/internal_message/service/v1"
)

type InternalMessageRecipientService struct {
	adminV1.InternalMessageRecipientServiceHTTPServer

	log *log.Helper

	messageServiceClient          internalMessageV1.InternalMessageServiceClient
	messageRecipientServiceClient internalMessageV1.InternalMessageRecipientServiceClient
}

func NewInternalMessageRecipientService(
	ctx *bootstrap.Context,
	messageServiceClient internalMessageV1.InternalMessageServiceClient,
	messageRecipientServiceClient internalMessageV1.InternalMessageRecipientServiceClient,
) *InternalMessageRecipientService {
	l := log.NewHelper(log.With(ctx.GetLogger(), "module", "internal-message-recipient/service/admin-service"))
	return &InternalMessageRecipientService{
		log:                           l,
		messageServiceClient:          messageServiceClient,
		messageRecipientServiceClient: messageRecipientServiceClient,
	}
}

// ListUserInbox 获取用户的收件箱列表 (通知类)
func (s *InternalMessageRecipientService) ListUserInbox(ctx context.Context, req *paginationV1.PagingRequest) (*internalMessageV1.ListUserInboxResponse, error) {
	resp, err := s.messageRecipientServiceClient.ListUserInbox(ctx, req)
	if err != nil {
		return nil, err
	}

	for _, d := range resp.Items {
		if d.MessageId == nil {
			continue
		}

		msg, err := s.messageServiceClient.GetMessage(ctx, &internalMessageV1.GetInternalMessageRequest{
			QueryBy: &internalMessageV1.GetInternalMessageRequest_Id{
				Id: d.GetMessageId(),
			},
		})
		if err != nil {
			s.log.Errorf("list user inbox failed, get message failed: %s", err)
			continue
		}

		d.Title = msg.Title
		d.Content = msg.Content
	}

	return resp, nil
}

func (s *InternalMessageRecipientService) DeleteNotificationFromInbox(ctx context.Context, req *internalMessageV1.DeleteNotificationFromInboxRequest) (*emptypb.Empty, error) {
	return s.messageRecipientServiceClient.DeleteNotificationFromInbox(ctx, req)
}

// MarkNotificationAsRead 将通知标记为已读
func (s *InternalMessageRecipientService) MarkNotificationAsRead(ctx context.Context, req *internalMessageV1.MarkNotificationAsReadRequest) (*emptypb.Empty, error) {
	return s.messageRecipientServiceClient.MarkNotificationAsRead(ctx, req)
}

// MarkNotificationsStatus 标记特定用户的某些或所有通知的状态
func (s *InternalMessageRecipientService) MarkNotificationsStatus(ctx context.Context, req *internalMessageV1.MarkNotificationsStatusRequest) (*emptypb.Empty, error) {
	return s.messageRecipientServiceClient.MarkNotificationsStatus(ctx, req)
}
