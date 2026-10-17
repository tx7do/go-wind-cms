package service

import (
	"context"
	"errors"
	"time"

	"github.com/go-kratos/kratos/v2/log"

	notifyV1 "go-wind-cms/api/gen/go/notify/service/v1"

	"go-wind-cms/app/core/service/internal/data"
	"go-wind-cms/app/core/service/internal/data/channel"
)

// InternalMessageSender 站内信（INTERNAL）渠道。
//
// 为什么它住在 internal/service 而不是 data/channel/：站内信的实现就是同包里的
// InternalMessageService（收件行落库 + SSE 推送），而 data 包 import service 包会成环。
// 于是本适配器由 wiring.go 在装配期注册进 data/channel 的 Registry——
// 这是全仓唯一一条"渠道实现反向依赖服务层"的边。
//
// 它持有的是 *InternalMessageService 本身而不是收件仓储 + redis client：
// 两者都可能滞后注入（redis 在 store 起来后才有，publisher 缝在两个 service
// 都就绪后互相注册），构造期快照会把它永久钉在 noop 上——站内信落库成功
// 但没人收到实时通知，且不报错。
type InternalMessageSender struct {
	messages *InternalMessageService
	// prefRepo 用户通知偏好（静音时段判定）；nil 视作"无人配置偏好"（fail-open 不静音）。
	prefRepo *data.NotificationPreferenceRepo
	log      *log.Helper
}

var _ channel.Sender = (*InternalMessageSender)(nil)

func NewInternalMessageSender(messages *InternalMessageService, prefRepo *data.NotificationPreferenceRepo, log *log.Helper) *InternalMessageSender {
	return &InternalMessageSender{messages: messages, prefRepo: prefRepo, log: log}
}

func (s *InternalMessageSender) Channel() notifyV1.Channel {
	return notifyV1.Channel_INTERNAL
}

// Send 把一个收件人挂到一条已存在的消息本体上：落收件行 + 尽力实时推送。
//
// 回执不带 ChannelID：站内信不吃 notification_channels 花名册（没有"用哪个账号发"的
// 问题），台账该列留空即"这条投递不来自任何渠道配置"。
//
// 实时推送失败不算投递失败——站内信以落库为准，用户离线时重连后仍能从收件箱补取；
// 偏好启用静音时段时只抑制实时推送（桌面角标即时性），收件行照常落库，
// 用户打开收件箱仍能看到全部消息（与上游通知域同一条语义）。
func (s *InternalMessageSender) Send(ctx context.Context, req *channel.SendRequest) (*channel.SendReceipt, error) {
	if req.RelatedID == 0 {
		return nil, errors.New("internal message delivery requires related_id: the internal_messages row to attach this recipient to")
	}
	if req.RecipientUserID == 0 {
		return nil, errors.New("internal message delivery requires recipient_user_id")
	}

	now := time.Now()
	if err := s.messages.sendNotification(ctx, req.RelatedID, req.RecipientUserID, req.OperatorUserID, &now, req.Title, req.Content); err != nil {
		return nil, err
	}

	if !userInQuietWindow(ctx, s.prefRepo, s.log, req.RecipientUserID, now) {
		s.messages.publishNotifyEvent(req.RecipientUserID, req.Title, req.Content)
	}

	return &channel.SendReceipt{}, nil
}

// userInQuietWindow 判定收件人此刻是否处在实时推送静音时段。
// 全链 fail-open：偏好缺失/读失败/配置非法一律不静音（推送是即时性的加分项，不该成为故障面）。
// 跨零点窗口（start > end，如 22:00→08:00）合法；start == end 视为无效配置。
func userInQuietWindow(ctx context.Context, prefRepo *data.NotificationPreferenceRepo, log *log.Helper, userID uint32, now time.Time) bool {
	if prefRepo == nil {
		return false
	}
	pref, err := prefRepo.GetByUserID(ctx, userID)
	if err != nil {
		log.Warnf("read notification preference for user [%d] failed, treat as not in quiet window: %s", userID, err.Error())
		return false
	}
	if pref == nil || !pref.GetQuietEnabled() {
		return false
	}

	start, end := int(pref.GetQuietStartMinute()), int(pref.GetQuietEndMinute())
	if start < 0 || end < 0 || start > 1439 || end > 1439 || start == end {
		return false
	}

	cur := now.Hour()*60 + now.Minute()
	if start < end {
		return cur >= start && cur < end
	}
	return cur >= start || cur < end
}
