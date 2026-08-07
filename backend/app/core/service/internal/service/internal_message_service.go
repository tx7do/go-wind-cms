package service

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/go-utils/aggregator"
	"github.com/tx7do/go-utils/timeutil"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/tx7do/go-crud/viewer"

	"go-wind-cms/app/core/service/internal/data"

	internalMessageV1 "go-wind-cms/api/gen/go/internal_message/service/v1"
	identityV1 "go-wind-cms/api/gen/go/identity/service/v1"
)

type InternalMessageService struct {
	internalMessageV1.UnimplementedInternalMessageServiceServer

	log *log.Helper

	internalMessageRepo          *data.InternalMessageRepo
	internalMessageCategoryRepo  *data.InternalMessageCategoryRepo
	internalMessageRecipientRepo *data.InternalMessageRecipientRepo

	userRepo data.UserRepo
}

func NewInternalMessageService(
	ctx *bootstrap.Context,
	internalMessageRepo *data.InternalMessageRepo,
	internalMessageCategoryRepo *data.InternalMessageCategoryRepo,
	internalMessageRecipientRepo *data.InternalMessageRecipientRepo,
	userRepo data.UserRepo,
) *InternalMessageService {
	return &InternalMessageService{
		log:                          ctx.NewLoggerHelper("internal-message/service/core-service"),
		internalMessageRepo:          internalMessageRepo,
		internalMessageCategoryRepo:  internalMessageCategoryRepo,
		internalMessageRecipientRepo: internalMessageRecipientRepo,
		userRepo:                     userRepo,
	}
}

func (s *InternalMessageService) extractRelationIDs(
	messages []*internalMessageV1.InternalMessage,
	categorySet aggregator.ResourceMap[uint32, *internalMessageV1.InternalMessageCategory],
) {
	for _, p := range messages {
		if p.GetCategoryId() > 0 {
			categorySet[p.GetCategoryId()] = nil
		}
	}
}

func (s *InternalMessageService) fetchRelationInfo(
	ctx context.Context,
	categorySet aggregator.ResourceMap[uint32, *internalMessageV1.InternalMessageCategory],
) error {
	if len(categorySet) > 0 {
		categoryIds := make([]uint32, 0, len(categorySet))
		for id := range categorySet {
			categoryIds = append(categoryIds, id)
		}

		categories, err := s.internalMessageCategoryRepo.ListCategoriesByIds(ctx, categoryIds)
		if err != nil {
			s.log.Errorf("query internal message category err: %v", err)
			return err
		}

		for _, g := range categories {
			categorySet[g.GetId()] = g
		}
	}

	return nil
}

func (s *InternalMessageService) bindRelations(
	messages []*internalMessageV1.InternalMessage,
	categorySet aggregator.ResourceMap[uint32, *internalMessageV1.InternalMessageCategory],
) {
	aggregator.Populate(
		messages,
		categorySet,
		func(ou *internalMessageV1.InternalMessage) uint32 { return ou.GetCategoryId() },
		func(ou *internalMessageV1.InternalMessage, c *internalMessageV1.InternalMessageCategory) {
			ou.CategoryName = c.Name
		},
	)
}

func (s *InternalMessageService) enrichRelations(ctx context.Context, messages []*internalMessageV1.InternalMessage) error {
	var categorySet = make(aggregator.ResourceMap[uint32, *internalMessageV1.InternalMessageCategory])
	s.extractRelationIDs(messages, categorySet)
	if err := s.fetchRelationInfo(ctx, categorySet); err != nil {
		return err
	}
	s.bindRelations(messages, categorySet)
	return nil
}

func (s *InternalMessageService) ListMessage(ctx context.Context, req *paginationV1.PagingRequest) (*internalMessageV1.ListInternalMessageResponse, error) {
	resp, err := s.internalMessageRepo.List(ctx, req)
	if err != nil {
		return nil, err
	}

	_ = s.enrichRelations(ctx, resp.Items)

	return resp, nil
}

func (s *InternalMessageService) GetMessage(ctx context.Context, req *internalMessageV1.GetInternalMessageRequest) (*internalMessageV1.InternalMessage, error) {
	resp, err := s.internalMessageRepo.Get(ctx, req)
	if err != nil {
		return nil, err
	}

	fakeItems := []*internalMessageV1.InternalMessage{resp}
	_ = s.enrichRelations(ctx, fakeItems)

	return resp, nil
}

func (s *InternalMessageService) CreateMessage(ctx context.Context, req *internalMessageV1.CreateInternalMessageRequest) (*internalMessageV1.InternalMessage, error) {
	if req == nil || req.Data == nil {
		return nil, internalMessageV1.ErrorBadRequest("invalid parameter")
	}

	var created *internalMessageV1.InternalMessage
	var err error
	if created, err = s.internalMessageRepo.Create(ctx, req); err != nil {
		return nil, err
	}

	return created, nil
}

func (s *InternalMessageService) UpdateMessage(ctx context.Context, req *internalMessageV1.UpdateInternalMessageRequest) (*emptypb.Empty, error) {
	if req == nil || req.Data == nil {
		return nil, internalMessageV1.ErrorBadRequest("invalid parameter")
	}

	if err := s.internalMessageRepo.Update(ctx, req); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *InternalMessageService) DeleteMessage(ctx context.Context, req *internalMessageV1.DeleteInternalMessageRequest) (*emptypb.Empty, error) {
	if err := s.internalMessageRepo.Delete(ctx, req.GetId()); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// RevokeMessage 撤销某条消息
func (s *InternalMessageService) RevokeMessage(ctx context.Context, req *internalMessageV1.RevokeMessageRequest) (*emptypb.Empty, error) {
	// 仅消息发送者本人（或平台/系统上下文）可撤销消息体，防止跨租户/越权删除他人消息
	callerUserID, hasUser := viewerUserIDFromContext(ctx)
	isPlatform := false
	if vc, exist := viewer.FromContext(ctx); exist && vc != nil && (vc.IsPlatformContext() || vc.IsSystemContext()) {
		isPlatform = true
	}
	if !hasUser && !isPlatform {
		return nil, internalMessageV1.ErrorBadRequest("missing viewer context")
	}

	msg, err := s.internalMessageRepo.Get(ctx, &internalMessageV1.GetInternalMessageRequest{
		QueryBy: &internalMessageV1.GetInternalMessageRequest_Id{Id: req.GetMessageId()},
	})
	if err != nil || msg == nil {
		return nil, internalMessageV1.ErrorNotFound("internal message not found")
	}
	if !isPlatform {
		sender := msg.GetCreatedBy()
		if !hasUser || sender == 0 || sender != callerUserID {
			return nil, internalMessageV1.ErrorForbidden("only the sender can revoke this message")
		}
	}

	// 消息体删除失败立即返回，不继续清理收件人，避免不一致状态
	if err = s.internalMessageRepo.Delete(ctx, req.GetMessageId()); err != nil {
		s.log.Errorf("delete internal message failed: [%d] %s", req.GetMessageId(), err.Error())
		return nil, err
	}

	if err = s.internalMessageRecipientRepo.RevokeMessage(ctx, req); err != nil {
		s.log.Errorf("delete internal message inbox failed: [%d][%d] %s", req.GetMessageId(), req.GetUserId(), err.Error())
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// SendMessage 发送消息
func (s *InternalMessageService) SendMessage(ctx context.Context, req *internalMessageV1.SendMessageRequest) (*internalMessageV1.SendMessageResponse, error) {
	if req == nil {
		return nil, internalMessageV1.ErrorBadRequest("invalid request")
	}

	// 发送者身份必须从 viewer context 推导，忽略客户端传入的 send_user_id，防止越权伪造
	senderUserID, hasUser := viewerUserIDFromContext(ctx)
	if !hasUser {
		return nil, internalMessageV1.ErrorBadRequest("sender identity is required")
	}

	// 从 viewer context 获取发送者的租户 ID，用于限定消息发送范围
	senderTenantID := uint32(0)
	if vc, exist := viewer.FromContext(ctx); exist && vc != nil {
		senderTenantID = uint32(vc.TenantID())
	}

	now := time.Now()

	var err error
	var msg *internalMessageV1.InternalMessage
	if msg, err = s.internalMessageRepo.Create(ctx, &internalMessageV1.CreateInternalMessageRequest{
		Data: &internalMessageV1.InternalMessage{
			Title:      req.Title,
			Content:    trans.Ptr(req.GetContent()),
			Status:     trans.Ptr(internalMessageV1.InternalMessage_PUBLISHED),
			Type:       trans.Ptr(req.GetType()),
			CategoryId: req.CategoryId,
			CreatedBy:  trans.Ptr(senderUserID),
			CreatedAt:  timeutil.TimeToTimestamppb(&now),
		},
	}); err != nil {
		s.log.Errorf("create internal message failed: %s", err)
		return nil, err
	}

	// 平台管理员（tenant_id=0）可向所有租户用户发送；普通租户用户只能向本租户用户发送
	if req.GetTargetAll() {
		users, err := s.userRepo.List(ctx, &paginationV1.PagingRequest{NoPaging: trans.Ptr(true)})
		if err != nil {
			s.log.Errorf("send message failed, list users failed, %s", err)
		} else {
			for _, user := range users.Items {
				// 非平台上下文时，跳过其他租户的用户
				if senderTenantID != 0 && user.GetTenantId() != senderTenantID {
					continue
				}
				_ = s.sendNotification(ctx, msg.GetId(), user.GetId(), senderUserID, &now, msg.GetTitle(), msg.GetContent())
			}
		}
	} else {
		if req.RecipientUserId != nil {
			if s.isRecipientAllowed(ctx, req.GetRecipientUserId(), senderTenantID) {
				_ = s.sendNotification(ctx, msg.GetId(), req.GetRecipientUserId(), senderUserID, &now, msg.GetTitle(), msg.GetContent())
			}
		} else {
			if len(req.TargetUserIds) != 0 {
				for _, uid := range req.TargetUserIds {
					if s.isRecipientAllowed(ctx, uid, senderTenantID) {
						_ = s.sendNotification(ctx, msg.GetId(), uid, senderUserID, &now, msg.GetTitle(), msg.GetContent())
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
// 平台管理员（senderTenantID==0）可向任意租户用户发送；
// 普通租户用户只能向本租户用户发送。收件人不存在或跨租户时拒绝。
func (s *InternalMessageService) isRecipientAllowed(ctx context.Context, recipientUserID uint32, senderTenantID uint32) bool {
	recipient, err := s.userRepo.Get(ctx, &identityV1.GetUserRequest{
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
	if entity, err = s.internalMessageRecipientRepo.Create(ctx, recipient); err != nil {
		s.log.Errorf("send message failed, send to user failed, %s", err)
		return err
	}
	recipient.Id = entity.Id

	return nil
}
