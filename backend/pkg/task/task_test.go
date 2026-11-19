package task

import (
	"encoding/json"
	"testing"

	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── 任务类型命名约束 ──
//
// asynq ServeMux 在精确匹配失败后按【前缀】兜底路由：互为前缀的任务类型会让
// payload 被路由到错误的 handler 解码失败被吞（历史教训见 AGENTS.md）。
// 本测试对全部已注册任务类型做两两前缀检查，新增任务类型自动纳入。

// allTaskTypes 汇总所有 asynq 任务类型常量。
func allTaskTypes() []string {
	return []string{
		BackupTaskType,
		ContentHealthSweepTaskType,
		FlowExecuteTaskType,
		FlowResumeTaskType,
		FlowScheduleTickTaskType,
		MonitorAlertScanTaskType,
		NotificationDispatchTaskType,
		NotificationDeliverySweepTaskType,
		NotifyEmailTaskType,
		PostPublishTaskType,
		PostPublishSweepTaskType,
		SocialPublishSweepTaskType,
		SocialInboxSyncTaskType,
		SearchReindexTaskType,
		SearchReindexAllTaskType,
		AiDocReindexTaskType,
		AiQueryRunTaskType,
		AiPostIndexTaskType,
		SocialDistributionPublishTaskType,
		SocialStatsSyncTaskType,
		TenantQuotaWatermarkTaskType,
		WebhookFanoutTaskType,
		WebhookDeliverTaskType,
	}
}

func TestTaskTypesNonEmpty(t *testing.T) {
	for _, tt := range allTaskTypes() {
		if tt == "" {
			t.Fatal("task type constant must not be empty")
		}
	}
}

// knownPrefixPairs 已知的前缀组合豁免。search.reindex 与 search.reindex.all
// 互为前缀：asynq ServeMux 精确匹配优先，两者当前均注册在同一 mux
// （asynq_server.go），精确匹配命中、行为正确。但若 search.reindex.all 的
// 注册行丢失（AGENTS.md 记载过 sweep 任务"handler 在、注册行丢"的前科），
// 周期任务会被前缀兜底路由进单条重索引 handler，payload `{}` 解码后 no-op，
// 任务静默丢失——与 post.publish.sweep 事故同款。新增任务类型必须走严格
// 互不为前缀的命名，不得扩大本豁免表。
var knownPrefixPairs = [][2]string{
	{"search.reindex", "search.reindex.all"},
}

func isKnownPrefixPair(a, b string) bool {
	for _, p := range knownPrefixPairs {
		if (p[0] == a && p[1] == b) || (p[0] == b && p[1] == a) {
			return true
		}
	}
	return false
}

func TestTaskTypesNoPrefixCollision(t *testing.T) {
	types := allTaskTypes()
	for i, a := range types {
		for j, b := range types {
			if i == j {
				continue
			}
			if len(a) < len(b) && b[:len(a)] == a && !isKnownPrefixPair(a, b) {
				t.Errorf("task type %q is a prefix of %q — asynq 前缀兜底路由会在 b 未注册时把任务吞给 a 的 handler", a, b)
			}
			if len(b) < len(a) && a[:len(b)] == b && !isKnownPrefixPair(b, a) {
				t.Errorf("task type %q is a prefix of %q — asynq 前缀兜底路由会在 a 未注册时把任务吞给 b 的 handler", b, a)
			}
		}
	}
}

func TestTaskTypesUnique(t *testing.T) {
	seen := map[string]bool{}
	for _, tt := range allTaskTypes() {
		if seen[tt] {
			t.Errorf("duplicated task type %q", tt)
		}
		seen[tt] = true
	}
}

// ── cron 规格有效性 ──

func TestCronSpecsValid(t *testing.T) {
	specs := map[string]string{
		"content health sweep":  ContentHealthSweepCronSpec,
		"monitor alert scan":    MonitorAlertScanCronSpec,
		"notification sweep":    NotificationDeliverySweepCronSpec,
		"search reindex all":    SearchReindexAllCronSpec,
		"tenant quota scan":     TenantQuotaWatermarkCronSpec,
	}
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	for name, spec := range specs {
		if _, err := parser.Parse(spec); err != nil {
			t.Errorf("cron spec for %s (%q) is invalid: %v", name, spec, err)
		}
	}
}

// 兜底扫描任务与发布任务必须错开整点高峰（文档口径：03:00 备份 / 03:17 健康巡检 / 04:44 重索引）。
func TestNightlyCronSpecsPinned(t *testing.T) {
	assert.Equal(t, "17 3 * * *", ContentHealthSweepCronSpec)
	assert.Equal(t, "44 4 * * *", SearchReindexAllCronSpec)
	assert.Equal(t, "0 9 * * *", TenantQuotaWatermarkCronSpec)
}

// ── payload JSON 往返 ──
//
// asynq 以 JSON 编解码 payload，往返不丢字段是 worker 正确解码的前提。

func jsonRoundTrip(t *testing.T, in, out any) {
	t.Helper()
	b, err := json.Marshal(in)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(b, out))
}

func TestPayloadJSONRoundTrip(t *testing.T) {
	tenant := int64(7)

	var backupOut BackupTaskData
	jsonRoundTrip(t, BackupTaskData{Name: "daily"}, &backupOut)
	assert.Equal(t, "daily", backupOut.Name)

	var flowOut FlowExecutePayload
	flowIn := FlowExecutePayload{
		RuleID:      11,
		Event:       "post.published",
		EntityID:    22,
		TenantID:    tenant,
		ActionIndex: 3,
		Chain:       "else",
		Snapshot:    map[string]string{"status": "PUBLISHED", "title": "标题"},
	}
	jsonRoundTrip(t, flowIn, &flowOut)
	assert.Equal(t, flowIn, flowOut)

	var flowTickOut FlowScheduleTickPayload
	jsonRoundTrip(t, FlowScheduleTickPayload{}, &flowTickOut)

	var sweepOut PublishSweepPayload
	jsonRoundTrip(t, PublishSweepPayload{}, &sweepOut)

	var statsSyncOut StatsSyncPayload
	jsonRoundTrip(t, StatsSyncPayload{}, &statsSyncOut)

	var healthOut ContentHealthSweepTaskData
	jsonRoundTrip(t, ContentHealthSweepTaskData{}, &healthOut)

	var monitorOut MonitorAlertScanTaskData
	jsonRoundTrip(t, MonitorAlertScanTaskData{}, &monitorOut)

	var quotaOut TenantQuotaWatermarkTaskData
	jsonRoundTrip(t, TenantQuotaWatermarkTaskData{}, &quotaOut)

	var notifSweepOut NotificationDeliverySweepTaskData
	jsonRoundTrip(t, NotificationDeliverySweepTaskData{}, &notifSweepOut)

	var notifOut NotificationDispatchTaskData
	notifIn := NotificationDispatchTaskData{
		DeliveryId: 33,
		Target:     "user@example.com",
		Title:      "标题",
		Content:    "内容",
		Params:     map[string]string{"k": "v"},
	}
	jsonRoundTrip(t, notifIn, &notifOut)
	assert.Equal(t, notifIn, notifOut)

	var emailOut NotifyEmailPayload
	emailIn := NotifyEmailPayload{UserID: 1, TenantID: tenant, MessageID: 9, Title: "t", Content: "c"}
	jsonRoundTrip(t, emailIn, &emailOut)
	assert.Equal(t, emailIn, emailOut)

	var postOut PostPublishPayload
	jsonRoundTrip(t, PostPublishPayload{PostID: 5, TenantID: tenant}, &postOut)
	assert.Equal(t, PostPublishPayload{PostID: 5, TenantID: tenant}, postOut)

	var reindexOut SearchReindexPayload
	reindexIn := SearchReindexPayload{Entity: "post", ID: 5, TenantID: tenant, Op: "index"}
	jsonRoundTrip(t, reindexIn, &reindexOut)
	assert.Equal(t, reindexIn, reindexOut)

	var reindexAllOut SearchReindexAllTaskData
	jsonRoundTrip(t, SearchReindexAllTaskData{}, &reindexAllOut)

	var aiDocOut AiDocReindexPayload
	jsonRoundTrip(t, AiDocReindexPayload{BaseID: 3}, &aiDocOut)
	assert.Equal(t, AiDocReindexPayload{BaseID: 3}, aiDocOut)

	var aiQueryOut AiQueryRunPayload
	jsonRoundTrip(t, AiQueryRunPayload{Question: "q", Lang: "zh-CN"}, &aiQueryOut)
	assert.Equal(t, AiQueryRunPayload{Question: "q", Lang: "zh-CN"}, aiQueryOut)

	var aiPostOut AiPostIndexPayload
	jsonRoundTrip(t, AiPostIndexPayload{PostID: 5, TenantID: tenant, Op: "delete"}, &aiPostOut)
	assert.Equal(t, AiPostIndexPayload{PostID: 5, TenantID: tenant, Op: "delete"}, aiPostOut)

	var socialOut SocialDistributionPublishPayload
	jsonRoundTrip(t, SocialDistributionPublishPayload{DistributionID: 8, TenantID: tenant}, &socialOut)
	assert.Equal(t, SocialDistributionPublishPayload{DistributionID: 8, TenantID: tenant}, socialOut)

	var fanoutOut WebhookFanoutPayload
	jsonRoundTrip(t, WebhookFanoutPayload{Event: "post.published", EntityID: 5, TenantID: tenant}, &fanoutOut)
	assert.Equal(t, WebhookFanoutPayload{Event: "post.published", EntityID: 5, TenantID: tenant}, fanoutOut)

	var deliverOut WebhookDeliverPayload
	jsonRoundTrip(t, WebhookDeliverPayload{WebhookID: 2, Event: "post.created", EntityID: 5, TenantID: tenant}, &deliverOut)
	assert.Equal(t, WebhookDeliverPayload{WebhookID: 2, Event: "post.created", EntityID: 5, TenantID: tenant}, deliverOut)
}

// ── 事件名与关键常量钉测 ──

func TestWebhookEventNamesPinned(t *testing.T) {
	// 事件名是 webhook 订阅方与 Flow 规则的双边约定，不可改名。
	assert.Equal(t, "post.created", WebhookEventPostCreated)
	assert.Equal(t, "post.updated", WebhookEventPostUpdated)
	assert.Equal(t, "post.published", WebhookEventPostPublished)
	assert.Equal(t, "post.deleted", WebhookEventPostDeleted)
	assert.Equal(t, "page.created", WebhookEventPageCreated)
	assert.Equal(t, "page.updated", WebhookEventPageUpdated)
	assert.Equal(t, "page.published", WebhookEventPagePublished)
	assert.Equal(t, "page.deleted", WebhookEventPageDeleted)
	assert.Equal(t, "comment.created", WebhookEventCommentCreated)
	assert.Equal(t, "media.created", WebhookEventMediaCreated)
	assert.Equal(t, "user.created", WebhookEventUserCreated)
}

func TestTaskTypeValuesPinned(t *testing.T) {
	// 曾用名 post.publish.sweep 命中 post.publish 的前缀兜底路由，改名钉死。
	assert.Equal(t, "post.scheduled.sweep", PostPublishSweepTaskType)
	assert.Equal(t, "social.distribution.sweep", SocialPublishSweepTaskType)
	assert.Equal(t, "flow.execute", FlowExecuteTaskType)
	assert.Equal(t, "flow.resume", FlowResumeTaskType)
	assert.Equal(t, "flow.schedule.tick", FlowScheduleTickTaskType)
	assert.Equal(t, "backup", BackupTaskType)
}

// ── Flow 快照载荷语义 ──

func TestFlowExecutePayloadChainDefaults(t *testing.T) {
	// Chain 空串 = main 链（缺省语义），不应序列化出冗余字段。
	p := FlowExecutePayload{RuleID: 1}
	b, err := json.Marshal(p)
	require.NoError(t, err)
	assert.NotContains(t, string(b), "chain")
	assert.NotContains(t, string(b), "action_index")

	p.Chain = "else"
	b2, err := json.Marshal(p)
	require.NoError(t, err)
	assert.Contains(t, string(b2), `"chain":"else"`)
}
