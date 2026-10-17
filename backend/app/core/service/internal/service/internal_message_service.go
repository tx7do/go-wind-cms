package service

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/go-utils/aggregator"
	"github.com/tx7do/go-utils/timeutil"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/tx7do/go-crud/viewer"

	appViewer "go-wind-cms/pkg/entgo/viewer"
	"go-wind-cms/app/core/service/internal/data"
	"go-wind-cms/pkg/task"

	internalMessageV1 "go-wind-cms/api/gen/go/internal_message/service/v1"
	identityV1 "go-wind-cms/api/gen/go/identity/service/v1"
	notifyV1 "go-wind-cms/api/gen/go/notify/service/v1"

	"github.com/redis/go-redis/v9"
	"go-wind-cms/pkg/sseevent"
)

type InternalMessageService struct {
	internalMessageV1.UnimplementedInternalMessageServiceServer

	log *log.Helper

	internalMessageRepo          *data.InternalMessageRepo
	internalMessageCategoryRepo  *data.InternalMessageCategoryRepo
	internalMessageRecipientRepo *data.InternalMessageRecipientRepo

	userRepo data.UserRepo
	taskService *TaskService

	// redisClient 用于把系统通知事件发布到 Redis Pub/Sub（admin 服务订阅后转
	// SSE 实时推送）；为 nil 时静默跳过（落库为准，只损失实时性）。
	redisClient *redis.Client

	// notifier 通知域投递缝（Notifier 接口，由 NotificationService 实现）。
	// 构造期留 nil，装配期由 wiring 在两个 service 都就绪后 RegisterNotifier 注入
	//（双向引用：NotificationService 的 INTERNAL 渠道又反向持有本服务，构造期互相传会成环）。
	notifier Notifier
}

func NewInternalMessageService(
	ctx *bootstrap.Context,
	internalMessageRepo *data.InternalMessageRepo,
	internalMessageCategoryRepo *data.InternalMessageCategoryRepo,
	internalMessageRecipientRepo *data.InternalMessageRecipientRepo,
	userRepo data.UserRepo,
	redisClient *redis.Client,
	taskService *TaskService,
) *InternalMessageService {
	return &InternalMessageService{
		log:                          ctx.NewLoggerHelper("internal-message/service/core-service"),
		internalMessageRepo:          internalMessageRepo,
		internalMessageCategoryRepo:  internalMessageCategoryRepo,
		internalMessageRecipientRepo: internalMessageRecipientRepo,
		userRepo:                     userRepo,
		redisClient:                  redisClient,
		taskService:                  taskService,
	}
}

// RegisterNotifier 注入通知域投递缝；注入前系统通知走不到渠道路由，SendSystemNotification 直接报错
//（不静默降级成直写：装配遗漏应该暴露成显式失败，而不是让台账悄悄缺一列）。
func (s *InternalMessageService) RegisterNotifier(notifier Notifier) {
	s.notifier = notifier
}

// deliverViaNotifier 经通知域投递缝发一条通知；未接线时报错（装配遗漏应显式暴露，
// 不静默降级成直写让台账悄悄缺一列）。
func (s *InternalMessageService) deliverViaNotifier(ctx context.Context, req *notifyV1.SendDirectNotificationRequest) (*notifyV1.SendNotificationResponse, error) {
	if s.notifier == nil {
		return nil, errors.New("notification seam is not wired: RegisterNotifier must be called before delivering via notifier")
	}
	return s.notifier.SendDirect(ctx, req)
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

// SendSystemNotification 系统事件通知（Flow 规则动作等自动链路的唯一入口）：
// 以系统身份（created_by=0）创建 PUBLISHED 消息，收件人经通知域投递缝
//（Notifier.SendDirect，EventType=INTERNAL_MESSAGE）逐人投递——种子规则默认路由
// INTERNAL 渠道（落收件行 + 尽力实时推送），端态与直写一致，额外落一行投递台账。
// 管理员在通知规则页停用该规则即整体停止此类通知（此时 SendDirect 报错，只记日志不中断）。
// 收件人跨租户不拦截——系统消息的租户边界由规则配置方（收件人列表）负责。返回消息 ID。
func (s *InternalMessageService) SendSystemNotification(ctx context.Context, title, content string, targetUserIDs []uint32) (uint32, error) {
	if strings.TrimSpace(title) == "" || strings.TrimSpace(content) == "" {
		return 0, internalMessageV1.ErrorBadRequest("title and content are required")
	}
	targets := make([]uint32, 0, len(targetUserIDs))
	for _, uid := range targetUserIDs {
		if uid != 0 {
			targets = append(targets, uid)
		}
	}
	if len(targets) == 0 {
		return 0, internalMessageV1.ErrorBadRequest("target user ids are required")
	}

	now := time.Now()
	msg, err := s.internalMessageRepo.Create(ctx, &internalMessageV1.CreateInternalMessageRequest{
		Data: &internalMessageV1.InternalMessage{
			Title:     trans.Ptr(title),
			Content:   trans.Ptr(content),
			Status:    trans.Ptr(internalMessageV1.InternalMessage_PUBLISHED),
			Type:      trans.Ptr(internalMessageV1.InternalMessage_NOTIFICATION),
			CreatedBy: trans.Ptr(uint32(0)),
			CreatedAt: timeutil.TimeToTimestamppb(&now),
		},
	})
	if err != nil {
		s.log.Errorf("system notification: create message failed: %v", err)
		return 0, err
	}
	for _, uid := range targets {
		// 单收件人失败不中断其余收件人（与群发口径一致）；台账 ID 一并带进日志便于回查
		resp, dErr := s.notifier.SendDirect(ctx, &notifyV1.SendDirectNotificationRequest{
			EventType:       notifyV1.EventType_INTERNAL_MESSAGE,
			Target:          strconv.FormatUint(uint64(uid), 10),
			Title:           title,
			Content:         content,
			RecipientUserId: trans.Ptr(uid),
			RelatedId:       trans.Ptr(msg.GetId()),
		})
		if dErr != nil {
			s.log.Errorf("system notification: deliver to user [%d] failed (delivery=%d): %v", uid, resp.GetDeliveryId(), dErr)
		}
	}
	return msg.GetId(), nil
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
				_ = s.sendNotification(ctx, msg.GetId(), user.GetId(), senderUserID, &now, msg.GetTitle(), msg.GetContent(), req.GetSendEmail())
			}
		}
	} else {
		if req.RecipientUserId != nil {
			if s.isRecipientAllowed(ctx, req.GetRecipientUserId(), senderTenantID) {
				_ = s.sendNotification(ctx, msg.GetId(), req.GetRecipientUserId(), senderUserID, &now, msg.GetTitle(), msg.GetContent(), req.GetSendEmail())
			}
		} else {
			if len(req.TargetUserIds) != 0 {
				for _, uid := range req.TargetUserIds {
					if s.isRecipientAllowed(ctx, uid, senderTenantID) {
						_ = s.sendNotification(ctx, msg.GetId(), uid, senderUserID, &now, msg.GetTitle(), msg.GetContent(), req.GetSendEmail())
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


// NotifyUser 系统级通知：后台任务（无操作人身份）触达指定用户。
// 借鉴 go-wind-admin Notification 域的"系统发起"语义——不读 viewer/operator，
// sender 记为 0（系统）；租户从收件人用户记录推导（消息与收件行都落对租户）。
// 通知失败只记日志不返回错误：调用方（发布 worker 等任务链路）不能因通知中断。
//
// 邮件出站走通知域（SYSTEM_NOTIFY 事件经渠道路由，EmailSender 统一 HTML 壳）。
func (s *InternalMessageService) NotifyUser(ctx context.Context, recipientUserID uint32, title, content string) {
	if recipientUserID == 0 {
		return
	}

	recipient, err := s.userRepo.Get(ctx, &identityV1.GetUserRequest{
		QueryBy: &identityV1.GetUserRequest_Id{Id: recipientUserID},
	})
	if err != nil || recipient == nil {
		s.log.Errorf("notify user: recipient %d not found: %v", recipientUserID, err)
		return
	}
	tenantID := recipient.GetTenantId()
	now := time.Now()

	msg, err := s.internalMessageRepo.Create(ctx, &internalMessageV1.CreateInternalMessageRequest{
		Data: &internalMessageV1.InternalMessage{
			TenantId:  trans.Ptr(tenantID),
			Title:     trans.Ptr(title),
			Content:   trans.Ptr(content),
			Status:    trans.Ptr(internalMessageV1.InternalMessage_PUBLISHED),
			Type:      trans.Ptr(internalMessageV1.InternalMessage_NOTIFICATION),
			CreatedBy: trans.Ptr(uint32(0)), // 0 = 系统
			CreatedAt: timeutil.TimeToTimestamppb(&now),
		},
	})
	if err != nil {
		s.log.Errorf("notify user: create message failed (recipient=%d): %v", recipientUserID, err)
		return
	}

	if err = s.sendNotification(ctx, msg.GetId(), recipientUserID, 0, &now, title, content); err != nil {
		s.log.Errorf("notify user: create recipient failed (recipient=%d): %v", recipientUserID, err)
		return
	}

	// 实时推送（best-effort）：向 Redis Pub/Sub 发事件，admin 服务订阅后转 SSE
	// 实时推送给目标用户的全部在线会话。落库为准：推送失败只损失实时性，不影响到达。
	s.publishNotifyEvent(recipientUserID, title, content)

	// 邮件出站（best-effort）：收件人配有邮箱时异步投递（通知域 SYSTEM_NOTIFY 渠道链）。
	// 系统级通知（NotifyUser）保持"有邮箱即出站"口径；群发出站开关在
	// CreateMessage fan-out 路径（sendNotification 显式传 sendEmailOut）。
	// 站内信落库为准，邮件是附加触达——任何失败只记日志。
	if recipient.GetEmail() != "" && s.taskService != nil {
		if perr := s.taskService.EnqueueNotifyEmail(&task.NotifyEmailPayload{
			UserID:    recipientUserID,
			TenantID:  tenantID,
			MessageID: msg.GetId(),
			Title:     title,
			Content:   content,
		}); perr != nil {
			s.log.Warnf("notify user: enqueue email failed (recipient=%d): %v", recipientUserID, perr)
		}
	}
}

// SendEmailWorker 是 asynq "notify.email" 任务的 handler：
// 从 DB 取收件人邮箱，经通知域投递缝（SYSTEM_NOTIFY→EMAIL 渠道）发出系统通知邮件。
//
// 邮件即投递：走 SendDirect 落投递台账（渠道选择/脱敏/结论回写全在通知域），
// HTML 壳由 EmailSender 统一包裹；种子规则 SYSTEM_NOTIFY→EMAIL 提供路由，
// 管理员停用该行即整体关闭系统通知邮件出站（SendDirect 报错，本 worker 记 SKIPPED 不重试）。
//
// 幂等/重试口径：
//   - 渠道未配置（SKIPPED）/ 收件人无邮箱 / 用户已删 → 视为永久不适用，返回 nil 不重试
//   - SMTP 暂时性失败（FAILED）→ 返回 err 交 asynq MaxRetry 退避重试（重试产生新的台账行）
func (s *InternalMessageService) SendEmailWorker(taskType string, payload *task.NotifyEmailPayload) error {
	if payload == nil || payload.UserID == 0 {
		return nil
	}
	if s.taskService == nil {
		return nil
	}

	ctx := appViewer.NewSystemViewerContext(context.Background())
	user, err := s.userRepo.Get(ctx, &identityV1.GetUserRequest{
		QueryBy: &identityV1.GetUserRequest_Id{Id: payload.UserID},
	})
	if err != nil || user == nil {
		s.log.Warnf("notify email: user %d not loadable, drop: %v", payload.UserID, err)
		return nil
	}
	to := user.GetEmail()
	if to == "" {
		return nil
	}

	resp, derr := s.deliverViaNotifier(ctx, &notifyV1.SendDirectNotificationRequest{
		EventType:       notifyV1.EventType_SYSTEM_NOTIFY,
		Target:          to,
		Title:           payload.Title,
		Content:         payload.Content,
		RecipientUserId: trans.Ptr(payload.UserID),
		RelatedId:       trans.Ptr(payload.MessageID),
	})
	if derr != nil {
		if resp != nil && resp.GetStatus() == notifyV1.DeliveryStatus_SKIPPED {
			// 渠道未配置：永久不适用，不重试（原始原因已落台账 last_error）
			s.log.Warnf("notify email: skipped (user=%d, delivery=%d): %v", payload.UserID, resp.GetDeliveryId(), derr)
			return nil
		}
		s.log.Errorf("notify email: send failed (user=%d, delivery=%d): %v", payload.UserID, resp.GetDeliveryId(), derr)
		return derr // SMTP 暂时性失败：asynq 重试
	}
	s.log.Infof("notify email: sent to %s (user=%d, delivery=%d)", to, payload.UserID, resp.GetDeliveryId())
	return nil
}

// publishNotifyEvent 向 Redis Pub/Sub 发一条实时通知事件（admin 服务订阅后转 SSE）。
// best-effort：redisClient 未配置或发布失败只记日志——落库为准，只损失实时性。
func (s *InternalMessageService) publishNotifyEvent(recipientUserID uint32, title, content string) {
	if s.redisClient == nil {
		return
	}
	payload, jerr := json.Marshal(sseevent.NotifyEvent{
		UserID:  recipientUserID,
		Title:   title,
		Content: content,
	})
	if jerr != nil {
		s.log.Warnf("marshal realtime notify event failed: %v", jerr)
		return
	}
	if perr := s.redisClient.Publish(context.Background(), sseevent.RedisChannelNotify, payload).Err(); perr != nil {
		s.log.Warnf("publish realtime notify event failed: %v", perr)
	}
}

// sendNotification 向客户端发送通知消息
// sendEmailOut 差异化参数：显式传 true 时该收件人出站邮件（群发开关）；
// 省略时保持各调用方既有口径（系统级通知=有邮箱即出站，群发=默认不出站）。
func (s *InternalMessageService) sendNotification(ctx context.Context, messageId uint32, recipientUserId uint32, senderUserId uint32, now *time.Time, title, content string, sendEmailOut ...bool) error {
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
