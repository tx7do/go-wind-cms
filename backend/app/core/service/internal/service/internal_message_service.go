package service

import (
	"context"
	"encoding/json"
	"html"
	"os"
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
	"go-wind-cms/pkg/mailer"
	"go-wind-cms/pkg/task"

	internalMessageV1 "go-wind-cms/api/gen/go/internal_message/service/v1"
	identityV1 "go-wind-cms/api/gen/go/identity/service/v1"

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
	emailProviderRepo *data.EmailProviderRepo
	taskService *TaskService

	// redisClient 用于把系统通知事件发布到 Redis Pub/Sub（admin 服务订阅后转
	// SSE 实时推送）；为 nil 时静默跳过（落库为准，只损失实时性）。
	redisClient *redis.Client
}

func NewInternalMessageService(
	ctx *bootstrap.Context,
	internalMessageRepo *data.InternalMessageRepo,
	internalMessageCategoryRepo *data.InternalMessageCategoryRepo,
	internalMessageRecipientRepo *data.InternalMessageRecipientRepo,
	userRepo data.UserRepo,
	redisClient *redis.Client,
	taskService *TaskService,
	emailProviderRepo *data.EmailProviderRepo,
) *InternalMessageService {
	return &InternalMessageService{
		log:                          ctx.NewLoggerHelper("internal-message/service/core-service"),
		internalMessageRepo:          internalMessageRepo,
		internalMessageCategoryRepo:  internalMessageCategoryRepo,
		internalMessageRecipientRepo: internalMessageRecipientRepo,
		userRepo:                     userRepo,
		redisClient:                  redisClient,
		taskService:                  taskService,
		emailProviderRepo:            emailProviderRepo,
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

// NotifyUser 系统级通知：后台任务（无操作人身份）触达指定用户。
// 借鉴 go-wind-admin Notification 域的"系统发起"语义——不读 viewer/operator，
// sender 记为 0（系统）；租户从收件人用户记录推导（消息与收件行都落对租户）。
// 通知失败只记日志不返回错误：调用方（发布 worker 等任务链路）不能因通知中断。
//
// 设计取舍（对照 go-wind-admin 的完整 Notification 域）：站内信落库即达
// （收件中心可见），实时 SSE 推送延后；渠道抽象/路由规则/模板表等有第二类
// 投递渠道需求时再引入。
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
	if s.redisClient != nil {
		if payload, jerr := json.Marshal(sseevent.NotifyEvent{
			UserID:  recipientUserID,
			Title:   title,
			Content: content,
		}); jerr == nil {
			if perr := s.redisClient.Publish(context.Background(), sseevent.RedisChannelNotify, payload).Err(); perr != nil {
				s.log.Warnf("notify user: publish realtime event failed: %v", perr)
			}
		}
	}

	// 邮件出站（best-effort）：收件人配有邮箱且 SMTP 已配置时异步投递。
	// 站内信落库为准，邮件是附加触达——任何失败只记日志。
	if recipient.GetEmail() != "" && s.taskService != nil {
		if perr := s.taskService.EnqueueNotifyEmail(&task.NotifyEmailPayload{
			UserID:   recipientUserID,
			TenantID: tenantID,
			Title:    title,
			Content:  content,
		}); perr != nil {
			s.log.Warnf("notify user: enqueue email failed (recipient=%d): %v", recipientUserID, perr)
		}
	}
}

// SendEmailWorker 是 asynq "notify.email" 任务的 handler：
// 从 DB 取收件人邮箱，以最小 HTML 模板投递系统通知邮件。
//
// 幂等/重试口径：
//   - SMTP 未配置 / 收件人无邮箱 / 用户已删 → 视为永久不适用，返回 nil 不重试
//   - SMTP 暂时性失败 → 返回 err 交 asynq MaxRetry 退避重试
func (s *InternalMessageService) SendEmailWorker(taskType string, payload *task.NotifyEmailPayload) error {
	if payload == nil || payload.UserID == 0 {
		return nil
	}
	if s.taskService == nil {
		return nil
	}

	ctx := appViewer.NewSystemViewerContext(context.Background())
	cfg := s.resolveMailerConfig(ctx)
	if !cfg.Enabled() {
		s.log.Warnf("notify email: smtp not configured (DB/env), drop (user=%d)", payload.UserID)
		return nil
	}
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

	if err = cfg.Send([]string{to}, payload.Title, buildNotifyEmailHTML(payload.Title, payload.Content)); err != nil {
		s.log.Errorf("notify email: send failed (user=%d): %v", payload.UserID, err)
		return err // SMTP 暂时性失败：asynq 重试
	}
	s.log.Infof("notify email: sent to %s (user=%d)", to, payload.UserID)
	return nil
}

// resolveMailerConfig 邮件配置解析：DB 启用项（管理页可视化配置）优先，
// 未配置时回退 SMTP_* 环境变量（兼容未启用可视化配置的部署）。
func (s *InternalMessageService) resolveMailerConfig(ctx context.Context) mailer.Config {
	if s.emailProviderRepo != nil {
		if cfg, err := s.emailProviderRepo.GetActive(ctx); err == nil && cfg != nil {
			return mailer.Config{
				Host:     cfg.Host,
				Port:     int(cfg.Port),
				Username: cfg.Username,
				Password: cfg.Password,
				From:     cfg.From,
				TLSMode:  cfg.TLSMode,
			}
		} else if err != nil {
			s.log.Errorf("notify email: load db provider config failed, fallback to env: %v", err)
		}
	}
	return mailer.FromEnv()
}

// buildNotifyEmailHTML 系统通知邮件的最小 HTML 模板。
// 站点链接取 SITE_URL 环境变量（未配置则不附链接）。
func buildNotifyEmailHTML(title, content string) string {
	siteURL := strings.TrimRight(os.Getenv("SITE_URL"), "/")
	body := strings.ReplaceAll(html.EscapeString(content), "\n", "<br/>")
	var b strings.Builder
	b.WriteString(`<div style="max-width:560px;margin:0 auto;font-family:-apple-system,'Segoe UI',Roboto,'Helvetica Neue',Arial,sans-serif;color:#1f2329;">`)
	b.WriteString(`<h2 style="font-size:18px;margin:24px 0 12px;">` + html.EscapeString(title) + `</h2>`)
	b.WriteString(`<div style="font-size:14px;line-height:1.7;">` + body + `</div>`)
	if siteURL != "" {
		b.WriteString(`<p style="margin:24px 0 8px;"><a href="` + siteURL + `" style="color:#006be6;">` + html.EscapeString(siteURL) + `</a></p>`)
	}
	b.WriteString(`<hr style="border:none;border-top:1px solid #e5e6eb;margin:16px 0;"/><p style="font-size:12px;color:#86909c;">` +
		html.EscapeString("本邮件由系统自动发送，请勿直接回复。") + `</p></div>`)
	return b.String()
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
