package task

// ============================================================================
// Flow 规则执行任务类型定义
//
// 单段式（与 webhook 两段式不同：规则数远少于订阅数，worker 内直接执行动作，
// 失败整规则重试即可；动作数少时无重试粒度隔离的收益）：
//   flow.execute —— 内容事件触发；worker 加载规则 → 条件判定 → 逐动作执行
//
// 与 asynq 前缀路由约定兼容：flow.execute 不与其他任务类型互为前缀。
// ============================================================================

// FlowExecuteTaskType 规则执行任务类型。
const FlowExecuteTaskType = "flow.execute"

// FlowExecutePayload 规则执行任务 payload。
type FlowExecutePayload struct {
	RuleID   uint32 `json:"rule_id"`   // 规则 ID
	Event    string `json:"event"`     // 事件名（post.published 等，同 Webhook 事件族）
	EntityID uint32 `json:"entity_id"` // 事件主体 ID
	TenantID uint32 `json:"tenant_id"` // 租户 ID（条件查询按租户过滤）
}
