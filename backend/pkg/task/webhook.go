package task

// ============================================================================
// Webhook 内容事件外发任务类型定义
//
// 两段式投递（与 asynq 前缀路由约定兼容，两类任务名互不为前缀）：
//   1. webhook.fanout  —— 内容写路径入队一次；worker 按事件列出启用的订阅，
//      逐个再入队 webhook.deliver（单订阅粒度重试，互不牵连）
//   2. webhook.deliver —— 单订阅投递：落 delivery 记录 → HMAC 签名 POST →
//      回写结果；失败返回 err 由 asynq MaxRetry 指数退避重试
//
// 安全（与 search.reindex 同款模型）：payload 只含 id 与事件名，不含正文；
// 投递 body 为最小 envelope（事件名 + 实体 id + 状态），不携带全文内容。
// ============================================================================

const (
	// WebhookFanoutTaskType 事件扇出任务类型。
	WebhookFanoutTaskType = "webhook.fanout"
	// WebhookDeliverTaskType 单订阅投递任务类型。
	WebhookDeliverTaskType = "webhook.deliver"
)

// Webhook 内容事件名。
const (
	WebhookEventPostCreated   = "post.created"
	WebhookEventPostUpdated   = "post.updated"
	WebhookEventPostPublished = "post.published"
	WebhookEventPostDeleted   = "post.deleted"
	WebhookEventPageCreated   = "page.created"
	WebhookEventPageUpdated   = "page.updated"
	WebhookEventPagePublished = "page.published"
	WebhookEventPageDeleted   = "page.deleted"

	// v2 扩展事件族（仅 Flow 触发；webhook 订阅暂不放开）：
	// 负载带字段快照（status/title 等），条件判定优先吃快照。
	WebhookEventCommentCreated = "comment.created"
	WebhookEventMediaCreated   = "media.created"
	WebhookEventUserCreated    = "user.created"
)

// WebhookFanoutPayload 事件扇出任务 payload。
type WebhookFanoutPayload struct {
	Event     string `json:"event"`      // 事件名（post.published 等）
	EntityID  uint32 `json:"entity_id"`  // 事件主体 ID
	TenantID  uint32 `json:"tenant_id"`  // 租户 ID（扇出查询按租户过滤）
}

// WebhookDeliverPayload 单订阅投递任务 payload。
type WebhookDeliverPayload struct {
	WebhookID uint32 `json:"webhook_id"` // 订阅 ID
	Event     string `json:"event"`      // 事件名
	EntityID  uint32 `json:"entity_id"`  // 事件主体 ID
	TenantID  uint32 `json:"tenant_id"`  // 租户 ID（仅日志/排查）
}
