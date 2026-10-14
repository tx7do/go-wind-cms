package task

// ============================================================================
// 社交分发统计周期回拉任务类型定义
//
// 场景：已发布分发的互动统计平时靠管理页"刷新统计"手动回拉；本任务每 30 分钟
// 批量回拉近 30 天内发布、且距上次回拉超过 2 小时的分发记录，让列表页的
// 统计快照保持新鲜（急用仍走手动 RefreshSocialDistributionStats，即时生效）。
//
// 安全：payload 为空、SystemViewer 扫描；单次限量（防长事务），失败逐条跳过
// 不中断批次。
// ============================================================================

// SocialStatsSyncTaskType 社交分发统计周期回拉的 asynq 周期任务类型。
// 命名注意：asynq ServeMux 按前缀兜底路由，类型名不得与其他已注册类型互为前缀。
const SocialStatsSyncTaskType = "social.distribution.stats.sync"

// StatsSyncPayload 统计回拉任务 payload（无字段，扫描范围由 DB 决定）。
type StatsSyncPayload struct{}
