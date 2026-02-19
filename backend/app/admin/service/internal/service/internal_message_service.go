package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"

	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"github.com/tx7do/kratos-transport/transport/sse"

	"github.com/tx7do/go-utils/id"
	"github.com/tx7do/go-utils/timeutil"
	"github.com/tx7do/go-utils/trans"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	authenticationV1 "go-wind-cms/api/gen/go/authentication/service/v1"
	identityV1 "go-wind-cms/api/gen/go/identity/service/v1"
	internalMessageV1 "go-wind-cms/api/gen/go/internal_message/service/v1"

	"go-wind-cms/pkg/middleware/auth"
)

type InternalMessageService struct {
	adminV1.InternalMessageServiceHTTPServer

	log *log.Helper

	internalMessageServiceClient          internalMessageV1.InternalMessageServiceClient
	internalMessageCategoryServiceClient  internalMessageV1.InternalMessageCategoryServiceClient
	internalMessageRecipientServiceClient internalMessageV1.InternalMessageRecipientServiceClient

	userServiceClient           identityV1.UserServiceClient
	authenticationServiceClient authenticationV1.AuthenticationServiceClient

	sseServer *sse.Server
}

func NewInternalMessageService(
	ctx *bootstrap.Context,
	internalMessageRepo internalMessageV1.InternalMessageServiceClient,
	internalMessageCategoryRepo internalMessageV1.InternalMessageCategoryServiceClient,
	internalMessageRecipientRepo internalMessageV1.InternalMessageRecipientServiceClient,
	authenticationRepo authenticationV1.AuthenticationServiceClient,
	userRepo identityV1.UserServiceClient,
	sseServer *sse.Server,
) *InternalMessageService {
	l := log.NewHelper(log.With(ctx.GetLogger(), "module", "internal-message/service/admin-service"))
	return &InternalMessageService{
		log:                                   l,
		internalMessageServiceClient:          internalMessageRepo,
		internalMessageCategoryServiceClient:  internalMessageCategoryRepo,
		internalMessageRecipientServiceClient: internalMessageRecipientRepo,
		authenticationServiceClient:           authenticationRepo,
		userServiceClient:                     userRepo,
		sseServer:                             sseServer,
	}
}

func (s *InternalMessageService) ListMessage(ctx context.Context, req *paginationV1.PagingRequest) (*internalMessageV1.ListInternalMessageResponse, error) {
	resp, err := s.internalMessageServiceClient.ListMessage(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *InternalMessageService) GetMessage(ctx context.Context, req *internalMessageV1.GetInternalMessageRequest) (*internalMessageV1.InternalMessage, error) {
	resp, err := s.internalMessageServiceClient.GetMessage(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *InternalMessageService) CreateMessage(ctx context.Context, req *internalMessageV1.CreateInternalMessageRequest) (*emptypb.Empty, error) {
	if req.Data == nil {
		return nil, adminV1.ErrorBadRequest("invalid parameter")
	}

	// 获取操作人信息
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.Data.CreatedBy = trans.Ptr(operator.UserId)

	if _, err = s.internalMessageServiceClient.CreateMessage(ctx, req); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *InternalMessageService) UpdateMessage(ctx context.Context, req *internalMessageV1.UpdateInternalMessageRequest) (*emptypb.Empty, error) {
	if req.Data == nil {
		return nil, adminV1.ErrorBadRequest("invalid parameter")
	}

	// 获取操作人信息
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.Data.UpdatedBy = trans.Ptr(operator.GetUserId())
	if req.UpdateMask != nil {
		req.UpdateMask.Paths = append(req.UpdateMask.Paths, "updated_by")
	}

	return s.internalMessageServiceClient.UpdateMessage(ctx, req)
}

func (s *InternalMessageService) DeleteMessage(ctx context.Context, req *internalMessageV1.DeleteInternalMessageRequest) (*emptypb.Empty, error) {
	return s.internalMessageServiceClient.DeleteMessage(ctx, req)
}

// RevokeMessage 撤销某条消息
func (s *InternalMessageService) RevokeMessage(ctx context.Context, req *internalMessageV1.RevokeMessageRequest) (*emptypb.Empty, error) {
	return s.internalMessageServiceClient.RevokeMessage(ctx, req)
}

// SendMessage 发送消息
func (s *InternalMessageService) SendMessage(ctx context.Context, req *internalMessageV1.SendMessageRequest) (*internalMessageV1.SendMessageResponse, error) {
	// 获取操作人信息
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	var msg *internalMessageV1.InternalMessage
	if msg, err = s.internalMessageServiceClient.CreateMessage(ctx, &internalMessageV1.CreateInternalMessageRequest{
		Data: &internalMessageV1.InternalMessage{
			Title:      req.Title,
			Content:    trans.Ptr(req.GetContent()),
			Status:     trans.Ptr(internalMessageV1.InternalMessage_PUBLISHED),
			Type:       trans.Ptr(req.GetType()),
			CategoryId: req.CategoryId,
			CreatedBy:  trans.Ptr(operator.GetUserId()),
			CreatedAt:  timeutil.TimeToTimestamppb(&now),
		},
	}); err != nil {
		s.log.Errorf("create internal message failed: %s", err)
		return nil, err
	}

	if req.GetTargetAll() {
		users, err := s.userServiceClient.List(ctx, &paginationV1.PagingRequest{NoPaging: trans.Ptr(true)})
		if err != nil {
			s.log.Errorf("send message failed, list users failed, %s", err)
		} else {
			for _, user := range users.Items {
				_ = s.sendNotification(ctx, msg.GetId(), user.GetId(), operator.GetUserId(), &now, msg.GetTitle(), msg.GetContent())
			}
		}
	} else {
		if req.RecipientUserId != nil {
			_ = s.sendNotification(ctx, msg.GetId(), req.GetRecipientUserId(), operator.GetUserId(), &now, msg.GetTitle(), msg.GetContent())
		} else {
			if len(req.TargetUserIds) != 0 {
				for _, uid := range req.TargetUserIds {
					_ = s.sendNotification(ctx, msg.GetId(), uid, operator.GetUserId(), &now, msg.GetTitle(), msg.GetContent())
				}
			}
		}
	}

	return &internalMessageV1.SendMessageResponse{
		MessageId: msg.GetId(),
	}, nil
}

// sendNotification 向客户端发送通知消息
func (s *InternalMessageService) sendNotification(ctx context.Context, messageId uint32, recipientUserId uint32, senderUserId uint32, now *time.Time, title, content string) error {
	recipient := &internalMessageV1.InternalMessageRecipient{
		MessageId:       trans.Ptr(messageId),
		RecipientUserId: trans.Ptr(recipientUserId),
		Status:          trans.Ptr(internalMessageV1.InternalMessageRecipient_SENT),
		CreatedBy:       trans.Ptr(senderUserId),
		CreatedAt:       timeutil.TimeToTimestamppb(now),
		Title:           trans.Ptr(title),
		Content:         trans.Ptr(content),
	}

	var err error
	var entity *internalMessageV1.InternalMessageRecipient
	if entity, err = s.internalMessageRecipientServiceClient.Create(ctx, &internalMessageV1.CreateInternalMessageRecipientRequest{
		Data: recipient,
	}); err != nil {
		s.log.Errorf("send message failed, send to user failed, %s", err)
		return err
	}
	recipient.Id = entity.Id

	recipientJson, _ := json.Marshal(recipient)

	recipientStreamIds, err := s.authenticationServiceClient.GetAccessTokens(ctx, &authenticationV1.GetAccessTokensRequest{
		UserId:     recipientUserId,
		ClientType: authenticationV1.ClientType_admin,
	})
	if err != nil {
		s.log.Errorf("send message failed, get user access tokens failed, %s", err)
		return err
	}
	for _, streamId := range recipientStreamIds.AccessTokens {
		s.sseServer.Publish(ctx, sse.StreamID(streamId), &sse.Event{
			ID:    []byte(id.NewGUIDv7(false)),
			Data:  recipientJson,
			Event: []byte("notification"),
		})
	}

	return nil
}
