package task

// ============================================================================
// 内容健康度巡检任务
//
// 每日凌晨对已发布文章做一轮健康检查：
//   1. 翻译缺摘要（SEO/前台卡片展示受影响）
//   2. 正文外链可达性（HEAD 探测，超时/非 2xx 记为疑似死链）
// 问题按作者聚合成一条站内信（走 internalMessages 缝）。检查是旁路观测，
// 不改任何内容数据；上限参数防单轮扫爆。
// ============================================================================

// ContentHealthSweepTaskType 内容健康度巡检任务类型。
const ContentHealthSweepTaskType = "content.health.sweep"

// ContentHealthSweepCronSpec 每日 03:17（错开整点的备份 03:00 与高峰）。
const ContentHealthSweepCronSpec = "17 3 * * *"

// ContentHealthSweepTaskData 巡检 payload（无参数；边界在服务内钳位）。
type ContentHealthSweepTaskData struct{}
