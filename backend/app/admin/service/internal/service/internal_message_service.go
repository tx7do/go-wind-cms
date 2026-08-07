package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
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

type InternalMessagePublisher interface {
	Publish(ctx context.Context, streamId sse.StreamID, event *sse.Event)
}

type InternalMessageService struct {
	adminV1.InternalMessageServiceHTTPServer

	log *log.Helper

	internalMessageServiceClient          internalMessageV1.InternalMessageServiceClient
	internalMessageCategoryServiceClient  internalMessageV1.InternalMessageCategoryServiceClient
	internalMessageRecipientServiceClient internalMessageV1.InternalMessageRecipientServiceClient

	userServiceClient           identityV1.UserServiceClient
	authenticationServiceClient authenticationV1.AuthenticationServiceClient

	internalMessagePublisher InternalMessagePublisher

	clientType authenticationV1.ClientType
}

func NewInternalMessageService(
	ctx *bootstrap.Context,
	internalMessageRepo internalMessageV1.InternalMessageServiceClient,
	internalMessageCategoryRepo internalMessageV1.InternalMessageCategoryServiceClient,
	internalMessageRecipientRepo internalMessageV1.InternalMessageRecipientServiceClient,
	authenticationRepo authenticationV1.AuthenticationServiceClient,
	userRepo identityV1.UserServiceClient,
	clientType authenticationV1.ClientType,
) *InternalMessageService {
	l := log.NewHelper(log.With(ctx.GetLogger(), "module", "internal-message/service/admin-service"))
	return &InternalMessageService{
		log:                                   l,
		internalMessageServiceClient:          internalMessageRepo,
		internalMessageCategoryServiceClient:  internalMessageCategoryRepo,
		internalMessageRecipientServiceClient: internalMessageRecipientRepo,
		authenticationServiceClient:           authenticationRepo,
		userServiceClient:                     userRepo,
		clientType:                            clientType,
	}
}

func (s *InternalMessageService) RegisterInternalMessagePublisher(internalMessagePublisher InternalMessagePublisher) {
	s.internalMessagePublisher = internalMessagePublisher
}

func (s *InternalMessageService) HandleSubscribe(streamID sse.StreamID, _ *sse.Subscriber) {
	// streamID 即访问令牌，记录其指纹而非原始值，避免日志泄露可冒充的凭证
	s.log.Infof("subscriber [%s] connected", hashToken(string(streamID)))
}

func (s *InternalMessageService) HandleAuthorize(r *http.Request, token string) error {
	resp, err := s.authenticationServiceClient.ValidateToken(context.Background(), &authenticationV1.ValidateTokenRequest{
		ClientType:    s.clientType,
		Token:         token,
		TokenCategory: authenticationV1.TokenCategory_ACCESS,
	})
	if err != nil {
		s.log.Errorf("token authentication failed: %s", err)
		return err
	}

	tokenHash := hashToken(token)
	if resp.GetIsBlocked() {
		s.log.Warnf("token is blocked: %s", tokenHash)
		return authenticationV1.ErrorForbidden("token is blocked")
	}
	if !resp.GetIsValid() {
		s.log.Warnf("token is invalid: %s", tokenHash)
		return authenticationV1.ErrorUnauthorized("invalid token")
	}

	// 绑定 stream 到已验证令牌：通知流的 ID 即为接收方访问令牌，
	// 客户端只能订阅与其自身令牌同名的流，防止用合法令牌订阅他人通知流。
	// 框架从 Authorization/X-Token/?token= 提取令牌做校验，但 ?stream= 是独立
	// 提取的，若不在此处绑定，攻击者可携带 A 的令牌订阅 B 的流。
	if r == nil || r.URL == nil {
		return authenticationV1.ErrorForbidden("invalid request")
	}
	streamID := r.URL.Query().Get("stream")
	if streamID == "" || streamID != token {
		s.log.Warnf("stream/token mismatch, token: %s", tokenHash)
		return authenticationV1.ErrorForbidden("stream does not match token")
	}

	s.log.Debugf("token authenticated successfully, userId: [%d]", resp.GetPayload().GetUserId())

	return nil
}

// hashToken returns a short SHA-256 fingerprint of a token suitable for
// log correlation without exposing the raw credential.
func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])[:16]
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
	if req == nil || req.Data == nil {
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
	if req == nil || req.Data == nil {
		return nil, adminV1.ErrorBadRequest("invalid parameter")
	}

	// 获取操作人信息
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.Data.Id = trans.Ptr(req.GetId())

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

	// 平台上下文(tenant==0)可向任意租户用户发送；普通租户用户只能向本租户用户发送
	senderTenantID := operator.GetTenantId()
	if req.GetTargetAll() {
		users, err := s.userServiceClient.List(ctx, &paginationV1.PagingRequest{NoPaging: trans.Ptr(true)})
		if err != nil {
			s.log.Errorf("send message failed, list users failed, %s", err)
		} else {
			for _, user := range users.Items {
				// 非平台上下文时，跳过其他租户的用户
				if senderTenantID != 0 && user.GetTenantId() != senderTenantID {
					continue
				}
				_ = s.sendNotification(ctx, msg.GetId(), user.GetId(), operator.GetUserId(), &now, msg.GetTitle(), msg.GetContent())
			}
		}
	} else {
		if req.RecipientUserId != nil {
			if s.isRecipientAllowed(ctx, req.GetRecipientUserId(), senderTenantID) {
				_ = s.sendNotification(ctx, msg.GetId(), req.GetRecipientUserId(), operator.GetUserId(), &now, msg.GetTitle(), msg.GetContent())
			}
		} else {
			if len(req.TargetUserIds) != 0 {
				for _, uid := range req.TargetUserIds {
					if s.isRecipientAllowed(ctx, uid, senderTenantID) {
						_ = s.sendNotification(ctx, msg.GetId(), uid, operator.GetUserId(), &now, msg.GetTitle(), msg.GetContent())
					}
				}
			}
		}
	}

	return &internalMessageV1.SendMessageResponse{
		MessageId: msg.GetId(),
	}, nil
}

// isRecipientAllowed 校验收件人是否在发送者的可发送范围内。
// 平台上下文(senderTenantID==0)可向任意租户用户发送；普通租户用户只能向本租户用户发送。
// 收件人不存在或跨租户时拒绝。与 core InternalMessageService.SendMessage 的实现保持一致。
func (s *InternalMessageService) isRecipientAllowed(ctx context.Context, recipientUserID uint32, senderTenantID uint32) bool {
	recipient, err := s.userServiceClient.Get(ctx, &identityV1.GetUserRequest{
		QueryBy: &identityV1.GetUserRequest_Id{Id: recipientUserID},
	})
	if err != nil || recipient == nil {
		s.log.Errorf("send message failed, recipient not found [%d]: %v", recipientUserID, err)
		return false
	}
	if senderTenantID != 0 && recipient.GetTenantId() != senderTenantID {
		s.log.Errorf("send message forbidden, tenant mismatch: sender tenant [%d], recipient [%d] tenant [%d]", senderTenantID, recipientUserID, recipient.GetTenantId())
		return false
	}
	return true
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
		s.internalMessagePublisher.Publish(ctx, sse.StreamID(streamId), &sse.Event{
			ID:    []byte(id.NewGUIDv7(false)),
			Data:  recipientJson,
			Event: []byte("notification"),
		})
	}

	return nil
}
