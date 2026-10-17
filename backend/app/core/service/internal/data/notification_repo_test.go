package data

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tx7do/go-utils/trans"

	entSql "entgo.io/ent/dialect/sql"

	entCrud "github.com/tx7do/go-crud/entgo"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"

	"go-wind-cms/app/core/service/internal/data/client"
	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/enttest"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"

	notifyV1 "go-wind-cms/api/gen/go/notify/service/v1"

	_ "github.com/xiaoqidun/entps"
)

// newTestNotificationRepos 构造一套连内存 sqlite 的通知域仓储。
// 必须经 enttest 建客户端：enttest 真实 import ent/runtime，触发其 init() 填充各实体
// Hooks[0]，否则任何 mutation 都会报 "ent: uninitialized hook (forgotten import ent/runtime?)"。
// enttest.NewClient 同时自动跑 schema migration。
// updateDeliveryCreatedAtForTest 直接改台账行的 created_at（TimeAt mixin 把该列标成不可变，
// 生成的 builder 没有 SetCreatedAt；清扫测试要造陈旧行只能绕过 ORM 走裸 SQL）。
func updateDeliveryCreatedAtForTest(ctx context.Context, drv *entSql.Driver, id uint32, at time.Time) error {
	_, err := drv.ExecContext(ctx, "UPDATE notification_deliveries SET created_at = ? WHERE id = ?", at, id)
	return err
}

func newTestNotificationRepos(t *testing.T) (*ent.Client, *entSql.Driver, *NotificationChannelRepo, *NotificationDeliveryRepo, *NotificationRuleRepo, *NotificationTemplateRepo, *NotificationPreferenceRepo) {
	t.Helper()

	// 每个测试用独立命名的内存库：cache=shared 下同名 DSN 在同进程内共享，
	// 唯一约束（request_id+channel 等）会跨测试撞车。
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

	return db, drv,
		NewNotificationChannelRepo(bctx, wrapped),
		NewNotificationDeliveryRepo(bctx, wrapped),
		NewNotificationRuleRepo(bctx, wrapped),
		NewNotificationTemplateRepo(bctx, wrapped),
		NewNotificationPreferenceRepo(bctx, wrapped)
}

func TestNotificationDeliveryStateMachine(t *testing.T) {
	_, _, _, deliveryRepo, _, _, _ := newTestNotificationRepos(t)
	ctx := context.Background()

	// 落一条 SENDING 台账
	created, err := deliveryRepo.Create(ctx, &notifyV1.NotificationDelivery{
		EventType: trans.Ptr(notifyV1.EventType_INTERNAL_MESSAGE),
		Channel:   trans.Ptr(notifyV1.Channel_INTERNAL),
		Target:    trans.Ptr("42"),
		RequestId: trans.Ptr("req-1"),
	})
	require.NoError(t, err)
	require.NotNil(t, created.GetId())
	assert.Equal(t, notifyV1.DeliveryStatus_SENDING, created.GetStatus(), "新建台账默认 SENDING")
	assert.Equal(t, uint32(0), created.GetAttempts(), "attempts 是结果列，Create 不收请求值")

	// 结果回写：SENT + attempts + sent_at
	sentAt := time.Now()
	err = deliveryRepo.MarkResult(ctx, created.GetId(), &DeliveryOutcome{
		Status:    notifyV1.DeliveryStatus_SENT,
		SentAt:    &sentAt,
		Attempts:  trans.Ptr(uint32(1)),
		ChannelID: trans.Ptr(uint32(7)),
	})
	require.NoError(t, err)

	row, err := deliveryRepo.Get(ctx, created.GetId())
	require.NoError(t, err)
	assert.Equal(t, notifyV1.DeliveryStatus_SENT, row.GetStatus())
	assert.Equal(t, uint32(1), row.GetAttempts())
	assert.Equal(t, uint32(7), row.GetChannelId())
	assert.NotNil(t, row.GetSentAt())
}

func TestNotificationDeliverySweepOnlySettlesStaleSending(t *testing.T) {
	_, drv, _, deliveryRepo, _, _, _ := newTestNotificationRepos(t)
	ctx := context.Background()

	// 新鲜 SENDING（刚落）+ 陈旧 SENDING（created_at 拨回 1 小时前）+ 已定案 SENT
	fresh, err := deliveryRepo.Create(ctx, &notifyV1.NotificationDelivery{
		EventType: trans.Ptr(notifyV1.EventType_INTERNAL_MESSAGE),
		Channel:   trans.Ptr(notifyV1.Channel_INTERNAL),
		Target:    trans.Ptr("1"),
		RequestId: trans.Ptr("req-fresh"),
	})
	require.NoError(t, err)

	stale, err := deliveryRepo.Create(ctx, &notifyV1.NotificationDelivery{
		EventType: trans.Ptr(notifyV1.EventType_CHANNEL_TEST_EMAIL),
		Channel:   trans.Ptr(notifyV1.Channel_EMAIL),
		Target:    trans.Ptr("a***@example.com"),
		RequestId: trans.Ptr("req-stale"),
	})
	require.NoError(t, err)
	// 直接经 ent 改 created_at（repo 不开放写 created_at）
	past := time.Now().Add(-time.Hour)
	err = updateDeliveryCreatedAtForTest(ctx, drv, stale.GetId(), past)
	require.NoError(t, err)

	sent, err := deliveryRepo.Create(ctx, &notifyV1.NotificationDelivery{
		EventType: trans.Ptr(notifyV1.EventType_INTERNAL_MESSAGE),
		Channel:   trans.Ptr(notifyV1.Channel_INTERNAL),
		Target:    trans.Ptr("2"),
		RequestId: trans.Ptr("req-sent"),
	})
	require.NoError(t, err)
	err = deliveryRepo.MarkResult(ctx, sent.GetId(), &DeliveryOutcome{Status: notifyV1.DeliveryStatus_SENT})
	require.NoError(t, err)

	swept, err := deliveryRepo.SweepStaleSending(ctx, time.Now().Add(-15*time.Minute), "swept by test", 500)
	require.NoError(t, err)
	assert.Equal(t, 1, swept, "只清扫超期 SENDING")

	freshRow, err := deliveryRepo.Get(ctx, fresh.GetId())
	require.NoError(t, err)
	assert.Equal(t, notifyV1.DeliveryStatus_SENDING, freshRow.GetStatus(), "新鲜 SENDING 不动")

	staleRow, err := deliveryRepo.Get(ctx, stale.GetId())
	require.NoError(t, err)
	assert.Equal(t, notifyV1.DeliveryStatus_FAILED, staleRow.GetStatus(), "超期 SENDING 定案 FAILED")
	assert.Contains(t, staleRow.GetLastError(), "swept by test")

	sentRow, err := deliveryRepo.Get(ctx, sent.GetId())
	require.NoError(t, err)
	assert.Equal(t, notifyV1.DeliveryStatus_SENT, sentRow.GetStatus(), "已定案行不被清扫改写")
}

func TestNotificationRuleCreateAndRoute(t *testing.T) {
	_, _, _, _, ruleRepo, _, _ := newTestNotificationRepos(t)
	ctx := context.Background()

	// GetByEventType 空表：配置事实返回 (nil, nil)，不报错
	rule, err := ruleRepo.GetByEventType(ctx, notifyV1.EventType_INTERNAL_MESSAGE)
	require.NoError(t, err)
	assert.Nil(t, rule, "空表返回 nil 而非 NotFound")

	id, err := ruleRepo.Create(ctx, &notifyV1.NotificationRule{
		EventType: trans.Ptr(notifyV1.EventType_INTERNAL_MESSAGE),
		Channel:   trans.Ptr(notifyV1.Channel_INTERNAL),
		IsEnabled: trans.Ptr(true),
	}, 0)
	require.NoError(t, err)
	require.NotZero(t, id)

	rule, err = ruleRepo.GetByEventType(ctx, notifyV1.EventType_INTERNAL_MESSAGE)
	require.NoError(t, err)
	require.NotNil(t, rule)
	assert.Equal(t, notifyV1.Channel_INTERNAL, rule.GetChannel())
	assert.True(t, rule.GetIsEnabled())

	// 一事件一行：重复 event_type 显式拒绝
	_, err = ruleRepo.Create(ctx, &notifyV1.NotificationRule{
		EventType: trans.Ptr(notifyV1.EventType_INTERNAL_MESSAGE),
		Channel:   trans.Ptr(notifyV1.Channel_EMAIL),
	}, 0)
	require.Error(t, err, "event_type 唯一")

	// Update 不改 event_type（唯一锚不随编辑漂移）
	_, err = ruleRepo.Update(ctx, &notifyV1.UpdateNotificationRuleRequest{
		Id:   id,
		Data: &notifyV1.NotificationRule{Channel: trans.Ptr(notifyV1.Channel_EMAIL), IsEnabled: trans.Ptr(false)},
	}, 0)
	require.NoError(t, err)
	rule, err = ruleRepo.GetByEventType(ctx, notifyV1.EventType_INTERNAL_MESSAGE)
	require.NoError(t, err)
	require.NotNil(t, rule)
	assert.Equal(t, notifyV1.Channel_EMAIL, rule.GetChannel())
	assert.False(t, rule.GetIsEnabled())
}

func TestNotificationTemplateCode(t *testing.T) {
	_, _, _, _, _, templateRepo, _ := newTestNotificationRepos(t)
	ctx := context.Background()

	tpl, err := templateRepo.GetByCode(ctx, "notify/system")
	require.NoError(t, err)
	assert.Nil(t, tpl, "未配置模板返回 nil 而非 NotFound")

	id, err := templateRepo.Create(ctx, &notifyV1.NotificationTemplate{
		Name:            trans.Ptr("系统通知"),
		Code:            trans.Ptr("notify/system"),
		TitleTemplate:   trans.Ptr("{{title}}"),
		ContentTemplate: trans.Ptr("{{content}}"),
		IsEnabled:       trans.Ptr(true),
	}, 0)
	require.NoError(t, err)
	require.NotZero(t, id)

	_, err = templateRepo.Create(ctx, &notifyV1.NotificationTemplate{
		Code:            trans.Ptr("notify/system"),
		TitleTemplate:   trans.Ptr("x"),
		ContentTemplate: trans.Ptr("y"),
	}, 0)
	require.Error(t, err, "code 全局唯一")

	tpl, err = templateRepo.GetByCode(ctx, "notify/system")
	require.NoError(t, err)
	require.NotNil(t, tpl)
	assert.Equal(t, "系统通知", tpl.GetName())
}

func TestNotificationPreferenceUpsert(t *testing.T) {
	_, _, _, _, _, _, prefRepo := newTestNotificationRepos(t)
	ctx := context.Background()

	pref, err := prefRepo.GetByUserID(ctx, 42)
	require.NoError(t, err)
	assert.Nil(t, pref, "未配置偏好返回 nil（fail-open）")

	saved, err := prefRepo.Upsert(ctx, 42, &notifyV1.UpdateNotificationPreferenceRequest{
		QuietEnabled:     trans.Ptr(true),
		QuietStartMinute: trans.Ptr(int32(1320)),
		QuietEndMinute:   trans.Ptr(int32(480)),
	}, 42)
	require.NoError(t, err)
	assert.True(t, saved.GetQuietEnabled())

	// 二次 Upsert 走更新而非新建
	saved, err = prefRepo.Upsert(ctx, 42, &notifyV1.UpdateNotificationPreferenceRequest{
		QuietEnabled:     trans.Ptr(false),
		QuietStartMinute: trans.Ptr(int32(0)),
		QuietEndMinute:   trans.Ptr(int32(0)),
	}, 42)
	require.NoError(t, err)
	assert.False(t, saved.GetQuietEnabled())

	list, err := prefRepo.ListByUserIDs(ctx, []uint32{42, 43})
	require.NoError(t, err)
	assert.Len(t, list, 1, "缺偏好的用户不在返回值里")
}

func TestNotificationChannelListFillsFlagsAndEnabled(t *testing.T) {
	_, _, channelRepo, _, _, _, _ := newTestNotificationRepos(t)
	ctx := context.Background()

	created, err := channelRepo.Create(ctx, &notifyV1.CreateNotificationChannelRequest{
		Data: &notifyV1.NotificationChannel{
			Name:     trans.Ptr("主 SMTP"),
			Type:     trans.Ptr(notifyV1.NotificationChannel_EMAIL),
			Enabled:  trans.Ptr(true),
			SmtpTls:  trans.Ptr(notifyV1.NotificationChannel_START_TLS),
		},
		Password: trans.Ptr("secret-password"),
	})
	require.NoError(t, err)
	assert.True(t, created.GetEnabled(), "读视图回填 enabled（status 枚举 → 布尔）")
	assert.True(t, created.GetHasPassword(), "读视图回填凭据有无标识")
	// DTO 类型上根本没有 password 字段（只有 Create/Update 请求级明文入参）——
	// 凭据永不回显是编译期保证，不需要运行期断言。

	resp, err := channelRepo.List(ctx, &paginationV1.PagingRequest{NoPaging: trans.Ptr(true)})
	require.NoError(t, err)
	require.Len(t, resp.GetItems(), 1)
	item := resp.GetItems()[0]
	assert.Equal(t, notifyV1.NotificationChannel_EMAIL, item.GetType())
	assert.True(t, item.GetEnabled())
	assert.True(t, item.GetHasPassword())
	assert.False(t, item.GetHasWebhookSecret())

	// 自选策略：取第一个启用的 EMAIL 渠道并解密密码
	account, err := channelRepo.GetFirstEnabledEmailChannel(ctx)
	require.NoError(t, err)
	assert.Equal(t, "secret-password", account.Password, "sender 侧拿到解密明文")
	assert.NotEmpty(t, account.ID)
}
