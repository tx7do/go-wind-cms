package service

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/hibiken/asynq"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"

	entCrud "github.com/tx7do/go-crud/entgo"

	notifyV1 "go-wind-cms/api/gen/go/notify/service/v1"

	appViewer "go-wind-cms/pkg/entgo/viewer"
	"go-wind-cms/pkg/task"

	"go-wind-cms/app/core/service/internal/data"
	"go-wind-cms/app/core/service/internal/data/channel"
	"go-wind-cms/app/core/service/internal/data/client"
	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/enttest"

	_ "github.com/xiaoqidun/entps"
)

// fakeSender 可编程的渠道实现：Send 的返回由字段决定，用于驱动分发状态机的各条分支。
type fakeSender struct {
	ch      notifyV1.Channel
	sendErr error
	receipt *channel.SendReceipt
	sent    []*channel.SendRequest
}

func (f *fakeSender) Channel() notifyV1.Channel { return f.ch }

func (f *fakeSender) Send(_ context.Context, req *channel.SendRequest) (*channel.SendReceipt, error) {
	f.sent = append(f.sent, req)
	if f.receipt != nil {
		return f.receipt, f.sendErr
	}
	return nil, f.sendErr
}

// notificationTestStack 服务层测试的依赖栈：内存 sqlite + 通知域四个仓储。
type notificationTestStack struct {
	ent         *ent.Client
	delivery    *data.NotificationDeliveryRepo
	rule        *data.NotificationRuleRepo
	template    *data.NotificationTemplateRepo
	preference  *data.NotificationPreferenceRepo
	channelRepo *data.NotificationChannelRepo
}

// newNotificationTestStack 构造内存 sqlite 栈。必须经 enttest 建客户端：
// enttest 真实 import ent/runtime，触发其 init() 填充各实体 Hooks[0]。
func newNotificationTestStack(t *testing.T) *notificationTestStack {
	t.Helper()

	// 每个测试用独立命名的内存库：cache=shared 下同名 DSN 在同进程内共享，
	// 规则表等唯一约束会跨测试撞车（第一个测试播的种毒害第二个测试）。
	drv, err := entCrud.CreateDriver(
		"sqlite3",
		"file:ent_"+strings.ReplaceAll(t.Name(), "/", "_")+"?mode=memory&cache=shared&_fk=1&_time_format=sqlite",
		false, false,
	)
	require.NoError(t, err, "创建 sqlite driver 失败")

	db := enttest.NewClient(t, enttest.WithOptions(
		ent.Driver(drv),
		ent.Log(func(a ...any) { t.Log(a...) }),
	))
	client.ApplyTimeDefaultHooks(db)

	wrapped := entCrud.NewEntClient(db, drv)
	bctx := bootstrap.NewContextWithParam(context.Background(), &conf.AppInfo{}, &conf.Bootstrap{}, log.DefaultLogger)

	return &notificationTestStack{
		ent:         db,
		delivery:    data.NewNotificationDeliveryRepo(bctx, wrapped),
		rule:        data.NewNotificationRuleRepo(bctx, wrapped),
		template:    data.NewNotificationTemplateRepo(bctx, wrapped),
		preference:  data.NewNotificationPreferenceRepo(bctx, wrapped),
		channelRepo: data.NewNotificationChannelRepo(bctx, wrapped),
	}
}

func newTestNotificationService(t *testing.T) (*NotificationService, *data.NotificationDeliveryRepo, *fakeSender, *channel.Registry) {
	t.Helper()

	stack := newNotificationTestStack(t)
	bctx := bootstrap.NewContextWithParam(context.Background(), &conf.AppInfo{}, &conf.Bootstrap{}, log.DefaultLogger)

	registry := channel.NewRegistry()
	sender := &fakeSender{ch: notifyV1.Channel_INTERNAL}
	registry.Register(sender)

	return NewNotificationService(bctx, stack.delivery, stack.rule, stack.template, registry),
		stack.delivery, sender, registry
}

// TestNotificationServiceSendDirectStateMachine 驱动 SendDirect 的同步投递路径：
// 规则命中 → 落 SENDING 台账 → sender 拨号 → 结论回写（成功 / 配置类失败两种终态）。
func TestNotificationServiceSendDirectStateMachine(t *testing.T) {
	svc, deliveryRepo, sender, _ := newTestNotificationService(t)
	ctx := appViewer.NewSystemViewerContext(context.Background())

	// 播种路由：INTERNAL_MESSAGE → INTERNAL
	_, err := svc.ruleRepo.Create(ctx, &notifyV1.NotificationRule{
		EventType: trans.Ptr(notifyV1.EventType_INTERNAL_MESSAGE),
		Channel:   trans.Ptr(notifyV1.Channel_INTERNAL),
		IsEnabled: trans.Ptr(true),
	}, 0)
	require.NoError(t, err)

	// 成功路径：SENT + attempts=1 + relatedId 透传
	resp, err := svc.SendDirect(ctx, &notifyV1.SendDirectNotificationRequest{
		EventType:       notifyV1.EventType_INTERNAL_MESSAGE,
		Target:          "42",
		Title:           "标题",
		Content:         "正文",
		RecipientUserId: trans.Ptr(uint32(42)),
		RelatedId:       trans.Ptr(uint32(7)),
	})
	require.NoError(t, err)
	require.NotZero(t, resp.GetDeliveryId())
	assert.Equal(t, notifyV1.DeliveryStatus_SENT, resp.GetStatus())
	require.Len(t, sender.sent, 1)
	assert.Equal(t, uint32(7), sender.sent[0].RelatedID, "relatedId 透传给渠道（站内信挂收件行靠它）")
	assert.Equal(t, "42", sender.sent[0].Target)

	row, err := deliveryRepo.Get(ctx, resp.GetDeliveryId())
	require.NoError(t, err)
	assert.Equal(t, notifyV1.DeliveryStatus_SENT, row.GetStatus())
	assert.Equal(t, uint32(1), row.GetAttempts())
	assert.Equal(t, "42", row.GetTarget(), "INTERNAL 渠道 target 原样存（收件用户 ID）")

	// 渠道报错路径：FAILED + last_error，response 与 error 同时非 nil（Notifier 契约）
	sender.sendErr = errors.New("smtp dial timeout")
	resp2, err2 := svc.SendDirect(ctx, &notifyV1.SendDirectNotificationRequest{
		EventType: notifyV1.EventType_INTERNAL_MESSAGE,
		Target:    "43",
		Title:     "标题",
		Content:   "正文",
	})
	require.Error(t, err2)
	assert.Equal(t, notifyV1.DeliveryStatus_FAILED, resp2.GetStatus())
	row2, err := deliveryRepo.Get(ctx, resp2.GetDeliveryId())
	require.NoError(t, err)
	assert.Equal(t, notifyV1.DeliveryStatus_FAILED, row2.GetStatus())
	assert.Contains(t, row2.GetLastError(), "smtp dial timeout")
}

// TestNotificationServiceSkipsUnconfiguredChannel 渠道未配置（SKIPPED 语义）：
// fakeSender 返回包了 ErrChannelNotConfigured 的错误，台账记 SKIPPED 而非 FAILED。
func TestNotificationServiceSkipsUnconfiguredChannel(t *testing.T) {
	svc, deliveryRepo, sender, _ := newTestNotificationService(t)
	ctx := appViewer.NewSystemViewerContext(context.Background())

	_, err := svc.ruleRepo.Create(ctx, &notifyV1.NotificationRule{
		EventType: trans.Ptr(notifyV1.EventType_INTERNAL_MESSAGE),
		Channel:   trans.Ptr(notifyV1.Channel_INTERNAL),
		IsEnabled: trans.Ptr(true),
	}, 0)
	require.NoError(t, err)

	sender.sendErr = errors.Join(channel.ErrChannelNotConfigured, errors.New("no enabled SMTP account"))

	resp, err := svc.SendDirect(ctx, &notifyV1.SendDirectNotificationRequest{
		EventType: notifyV1.EventType_INTERNAL_MESSAGE,
		Target:    "44",
		Title:     "标题",
		Content:   "正文",
	})
	require.Error(t, err)
	assert.Equal(t, notifyV1.DeliveryStatus_SKIPPED, resp.GetStatus(), "配置类错误 = SKIPPED（没拨过号）")

	row, err := deliveryRepo.Get(ctx, resp.GetDeliveryId())
	require.NoError(t, err)
	assert.Equal(t, notifyV1.DeliveryStatus_SKIPPED, row.GetStatus())
}

// TestNotificationServiceDispatchIdempotency 异步消费的幂等门与重试语义：
// 已定案的行直接跳过；有重试额度的失败留在 SENDING；额度用尽才落 FAILED。
func TestNotificationServiceDispatchIdempotency(t *testing.T) {
	svc, deliveryRepo, sender, _ := newTestNotificationService(t)
	ctx := appViewer.NewSystemViewerContext(context.Background())

	// 1. 已定案（SENT）的行：handler 直接入幂等门返回 nil，不再拨号
	settled, err := deliveryRepo.Create(ctx, &notifyV1.NotificationDelivery{
		EventType: trans.Ptr(notifyV1.EventType_INTERNAL_MESSAGE),
		Channel:   trans.Ptr(notifyV1.Channel_INTERNAL),
		Target:    trans.Ptr("42"),
		Status:    trans.Ptr(notifyV1.DeliveryStatus_SENT),
	})
	require.NoError(t, err)
	sender.sendErr = errors.New("should not dial")
	err = svc.AsyncNotificationDispatch(ctx, "notification.dispatch", &task.NotificationDispatchTaskData{DeliveryId: settled.GetId()})
	require.NoError(t, err, "幂等门：非 SENDING 行静默跳过")
	assert.Len(t, sender.sent, 0, "已定案行不拨号")

	// 2. SENDING 行 + 失败：第一次失败留 SENDING（重试额度还在）
	sending, err := deliveryRepo.Create(ctx, &notifyV1.NotificationDelivery{
		EventType: trans.Ptr(notifyV1.EventType_INTERNAL_MESSAGE),
		Channel:   trans.Ptr(notifyV1.Channel_INTERNAL),
		Target:    trans.Ptr("43"),
	})
	require.NoError(t, err)
	require.Len(t, sender.sent, 0)
	sender.sendErr = errors.New("smtp busy 450")

	// attempts=0 → 第一次尝试后 attempts=1 < MaxAttempts(4)：留 SENDING
	err = svc.AsyncNotificationDispatch(ctx, "notification.dispatch", &task.NotificationDispatchTaskData{DeliveryId: sending.GetId()})
	require.Error(t, err, "投递失败上抛给 asynq 触发重试")
	row, err := deliveryRepo.Get(ctx, sending.GetId())
	require.NoError(t, err)
	assert.Equal(t, notifyV1.DeliveryStatus_SENDING, row.GetStatus(), "有重试额度时留在 SENDING（否则幂等门挡掉重投）")
	assert.Equal(t, uint32(1), row.GetAttempts())
	assert.Contains(t, row.GetLastError(), "smtp busy")

	// 3. attempts 已到 MaxAttempts：这次失败落 FAILED 终态
	require.NoError(t, deliveryRepo.MarkAttempted(ctx, sending.GetId(), notificationDispatchMaxAttempts))
	err = svc.AsyncNotificationDispatch(ctx, "notification.dispatch", &task.NotificationDispatchTaskData{DeliveryId: sending.GetId()})
	require.Error(t, err)
	row, err = deliveryRepo.Get(ctx, sending.GetId())
	require.NoError(t, err)
	assert.Equal(t, notifyV1.DeliveryStatus_FAILED, row.GetStatus(), "额度用尽定案 FAILED")

	// 4. SKIPPED（配置类）：立即定案并 SkipRetry
	skipping, err := deliveryRepo.Create(ctx, &notifyV1.NotificationDelivery{
		EventType: trans.Ptr(notifyV1.EventType_INTERNAL_MESSAGE),
		Channel:   trans.Ptr(notifyV1.Channel_INTERNAL),
		Target:    trans.Ptr("44"),
	})
	require.NoError(t, err)
	sender.sendErr = errors.Join(channel.ErrChannelNotConfigured, errors.New("no SMTP"))
	err = svc.AsyncNotificationDispatch(ctx, "notification.dispatch", &task.NotificationDispatchTaskData{DeliveryId: skipping.GetId()})
	require.Error(t, err)
	assert.True(t, errors.Is(err, asynq.SkipRetry), "配置类失败跳过重试（载荷正文不该躺在归档队列里过期）")
	row, err = deliveryRepo.Get(ctx, skipping.GetId())
	require.NoError(t, err)
	assert.Equal(t, notifyV1.DeliveryStatus_SKIPPED, row.GetStatus())
}

// TestNotificationRuleSeedingOnEmptyTable 规则服务构造期空表播种（INTERNAL_MESSAGE→INTERNAL，同步启用）。
func TestNotificationRuleSeedingOnEmptyTable(t *testing.T) {
	stack := newNotificationTestStack(t)
	bctx := bootstrap.NewContextWithParam(context.Background(), &conf.AppInfo{}, &conf.Bootstrap{}, log.DefaultLogger)

	svc := NewNotificationRuleService(bctx, stack.rule, nil, nil)

	rule, err := stack.rule.GetByEventType(appViewer.NewSystemViewerContext(context.Background()), notifyV1.EventType_INTERNAL_MESSAGE)
	require.NoError(t, err)
	require.NotNil(t, rule, "空表构造后应有播种行")
	assert.Equal(t, notifyV1.Channel_INTERNAL, rule.GetChannel())
	assert.True(t, rule.GetIsEnabled())
	assert.False(t, rule.GetIsAsync())
	assert.NotNil(t, svc)
}

// TestRenderTemplate 模板渲染：命中/未识别占位符报错带名/未闭合原样保留。
func TestRenderTemplate(t *testing.T) {
	out, err := renderTemplate("Hi {{name}}, order {{orderNo}}", map[string]string{"name": "Bob", "orderNo": "A1"})
	require.NoError(t, err)
	assert.Equal(t, "Hi Bob, order A1", out)

	_, err = renderTemplate("{{missing}}", map[string]string{"other": "x"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "{{missing}}", "未识别占位符报错带名字")

	out, err = renderTemplate("{{unclosed", nil)
	require.NoError(t, err)
	assert.Equal(t, "{{unclosed", out, "未闭合 {{ 原样保留")

	out, err = renderTemplate("", nil)
	require.NoError(t, err)
	assert.Empty(t, out)
}

// TestMaskTarget 台账脱敏：邮箱留首字符+域名，其余留末 4 位，INTERNAL 原样。
func TestMaskTarget(t *testing.T) {
	assert.Equal(t, "b***@example.com", maskTarget(notifyV1.Channel_EMAIL, "bob@example.com"))
	assert.Equal(t, "****abcd", maskTarget(notifyV1.Channel_WEBHOOK, "https://hook.example.com/abcd"))
	assert.Equal(t, "****", maskTarget(notifyV1.Channel_EMAIL, "abcd"))
	assert.Equal(t, "42", maskTarget(notifyV1.Channel_INTERNAL, "42"), "站内信目标是用户 ID，原样存")
}

// TestUserInQuietWindow 静音窗口判定：fail-open、跨零点、零宽不静音。
func TestUserInQuietWindow(t *testing.T) {
	repos := newNotificationTestStack(t)
	ctx := appViewer.NewSystemViewerContext(context.Background())

	// 无偏好仓库 / 无偏好行：fail-open 不静音
	assert.False(t, userInQuietWindow(ctx, nil, log.NewHelper(log.DefaultLogger), 1, time.Now()))
	assert.False(t, userInQuietWindow(ctx, repos.preference, log.NewHelper(log.DefaultLogger), 1, time.Now()))

	// 跨零点窗口 22:00→08:00：23:00 命中，12:00 不命中
	_, err := repos.preference.Upsert(ctx, 1, &notifyV1.UpdateNotificationPreferenceRequest{
		QuietEnabled:     trans.Ptr(true),
		QuietStartMinute: trans.Ptr(int32(1320)),
		QuietEndMinute:   trans.Ptr(int32(480)),
	}, 1)
	require.NoError(t, err)

	at23 := time.Date(2026, 10, 17, 23, 0, 0, 0, time.Local)
	at12 := time.Date(2026, 10, 17, 12, 0, 0, 0, time.Local)
	at7 := time.Date(2026, 10, 17, 7, 0, 0, 0, time.Local)
	assert.True(t, userInQuietWindow(ctx, repos.preference, log.NewHelper(log.DefaultLogger), 1, at23))
	assert.True(t, userInQuietWindow(ctx, repos.preference, log.NewHelper(log.DefaultLogger), 1, at7), "跨零点窗口的后半段")
	assert.False(t, userInQuietWindow(ctx, repos.preference, log.NewHelper(log.DefaultLogger), 1, at12))

	// 零宽窗口：按无效配置处理，不静音
	_, err = repos.preference.Upsert(ctx, 2, &notifyV1.UpdateNotificationPreferenceRequest{
		QuietEnabled:     trans.Ptr(true),
		QuietStartMinute: trans.Ptr(int32(600)),
		QuietEndMinute:   trans.Ptr(int32(600)),
	}, 2)
	require.NoError(t, err)
	assert.False(t, userInQuietWindow(ctx, repos.preference, log.NewHelper(log.DefaultLogger), 2, at12), "零宽窗口不静音")
}
