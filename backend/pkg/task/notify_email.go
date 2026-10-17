package task

// ============================================================================
// 邮件通知投递任务类型定义
//
// 场景：core InternalMessageService.NotifyUser（系统级通知）落库后，若收件人
// 配有邮箱且 SMTP 已配置，入队 notify.email 任务异步投递邮件。
//
// 安全/可靠性：
//   - payload 只含 id 与文案，不含凭据；收件邮箱由 worker 从 DB 取
//   - best-effort：入队失败只损失邮件不阻断通知落库
//   - SMTP 暂时性失败返回 err 交 asynq 退避重试（MaxRetry 3）
// ============================================================================

// NotifyEmailTaskType 邮件通知投递任务的 asynq 任务类型。
const NotifyEmailTaskType = "notify.email"

// NotifyEmailPayload 邮件通知投递任务 payload。
type NotifyEmailPayload struct {
	UserID    uint32 `json:"user_id"`    // 收件人用户 ID
	TenantID  uint32 `json:"tenant_id"`  // 租户 ID（仅日志/排查）
	MessageID uint32 `json:"message_id"` // 关联的站内信消息 ID（台账 related_id 用；0=无关联）
	Title     string `json:"title"`      // 通知标题（邮件主题）
	Content   string `json:"content"`    // 通知内容（邮件正文，纯文本）
}
