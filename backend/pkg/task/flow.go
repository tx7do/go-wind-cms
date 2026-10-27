package task

// ============================================================================
// Flow 规则执行任务类型定义
//
// 单段式（与 webhook 两段式不同：规则数远少于订阅数，worker 内直接执行动作，
// 失败整规则重试即可；动作数少时无重试粒度隔离的收益）：
//   flow.execute      —— 内容事件/定时触发；worker 加载规则 → 条件判定 → 逐动作执行
//   flow.resume       —— 延时动作续跑（DELAY 任务，ProcessIn 投递）
//   flow.schedule.tick —— 每分钟 tick，命中 cron 的定时规则入队 flow.execute
//
// 与 asynq 前缀路由约定兼容：三个类型互不为前缀
//（flow.execute.resume 之类的前缀形态会被 flow.execute 吞掉，禁用）。
// ============================================================================

// FlowExecuteTaskType 规则执行任务类型。
const FlowExecuteTaskType = "flow.execute"

// FlowExecutePayload 规则执行任务 payload。
type FlowExecutePayload struct {
	RuleID   uint32 `json:"rule_id"`   // 规则 ID
	Event    string `json:"event"`     // 事件名（post.published 等，同 Webhook 事件族；定时触发固定 "schedule"）
	EntityID uint32 `json:"entity_id"` // 事件主体 ID（定时触发为 0）
	TenantID uint32 `json:"tenant_id"` // 租户 ID（条件查询按租户过滤）
	// ActionIndex 延时续跑起始动作下标（0=完整执行：含条件判定；>0=跳过条件，
	// 从该下标继续动作链——命中判定已在触发时完成）。
	ActionIndex int `json:"action_index,omitempty"`
	// Chain 延时续跑所属动作链："main"（缺省）/ "else"（条件不满足分支）。
	Chain string `json:"chain,omitempty"`
	// Snapshot 触发时刻的实体字段快照（值统一字符串；status/prev_status/title/
	// slug/author_id/category_id/comment 系 content 等）。条件判定优先吃快照，
	// 兼修删除事件回查 404 与"worker 回查到执行时而非触发时状态"两类时序问题；
	// 旧任务无快照时回退 DB 回查（仅 status 可回退）。
	Snapshot map[string]string `json:"snapshot,omitempty"`
}

// FlowResumeTaskType 延时续跑任务类型。
// 命名约束：不得以 flow.execute 为前缀（asynq 精确匹配失败后按前缀兜底路由，
// flow.execute.resume 会被 flow.execute 的 handler 吞掉导致 payload 解码失败）。
const FlowResumeTaskType = "flow.resume"

// FlowScheduleTickTaskType 定时规则调度 tick（每分钟）。
const FlowScheduleTickTaskType = "flow.schedule.tick"

// FlowScheduleTickPayload 定时调度 tick payload（无参数，规则查询在 worker 内）。
type FlowScheduleTickPayload struct{}
