package task

// ============================================================================
// 社交分发发布任务类型定义
//
// 场景：social_distributions 创建为 SCHEDULED（或用户点"立即发布"）时入队一个
// social.distribution.publish 任务，asynq worker 到点后从 DB 取分发记录与账号
// 凭据，调用平台适配器发布并回写状态。
//
// 安全（与 search.reindex 同款模型）：
//   - payload 只含 id，不含文案与凭据，体积小、无敏感字段
//   - 实际内容由 worker 从 DB 取（带 SystemViewer 跨租户读）
//   - 凭据在 worker 内解密，只传入平台适配器，不进日志
//   - 发布失败回写 FAILED + error_message 后返回 err，由 asynq MaxRetry 重试
// ============================================================================

// SocialDistributionPublishTaskType 社交分发发布任务的 asynq 任务类型。
const SocialDistributionPublishTaskType = "social.distribution.publish"

// SocialDistributionPublishPayload 社交分发发布任务 payload。
//
// TenantID 仅供日志/排查展示，worker 以 DB 记录的 tenant_id 为准（不可用
// payload 值覆盖，与 SearchReindexPayload 的既有口径一致）。
type SocialDistributionPublishPayload struct {
	DistributionID uint32 `json:"distribution_id"`
	TenantID       uint32 `json:"tenant_id"`
}
