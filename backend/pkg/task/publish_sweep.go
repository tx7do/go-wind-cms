package task

// ============================================================================
// 定时发布兜底扫描任务类型定义
//
// 场景：post.publish / social.distribution.publish 都是 asynq DELAY 任务，
// 在 Redis 数据丢失、服务崩溃窗口内可能静默丢失——SCHEDULED 记录会被无声搁置。
// 两个 sweep 任务每分钟扫描"已到期但仍是 SCHEDULED"的记录并重新入队，
// 由各自的幂等 worker 执行（重复入队无害：worker 会跳过已发布/改期记录）。
//
// 安全：与对应发布 worker 同款模型——payload 为空、扫描用 SystemViewer、
// 重新入队的任务仍走 DB 状态自检。
// ============================================================================

// 任务类型命名注意：asynq ServeMux 在精确匹配失败后按【前缀】兜底路由
// （strings.HasPrefix），类型名不得是其他已注册类型的前缀——
// 曾用名 "post.publish.sweep" 会命中 "post.publish" 的 handler，payload
// 解码失败被吞。现改为互不为前缀的 "post.scheduled.sweep" / "social.distribution.sweep"。
const (
	// PostPublishSweepTaskType 文章定时发布兜底扫描的 asynq 周期任务类型。
	PostPublishSweepTaskType = "post.scheduled.sweep"

	// SocialPublishSweepTaskType 社交分发发布兜底扫描的 asynq 周期任务类型。
	SocialPublishSweepTaskType = "social.distribution.sweep"
)

// PublishSweepPayload 兜底扫描任务 payload（无字段，扫描范围由 DB 决定）。
type PublishSweepPayload struct{}
