package task

// ============================================================================
// 文章定时发布任务类型定义
//
// 场景：文章保存为 SCHEDULED（定时发布）时入队一个 post.publish 任务，
// asynq worker 到点后校验状态与排期，把 SCHEDULED 流转为 PUBLISHED。
// 修复点：前端编辑页一直承诺"到达时间后由后端调度器自动置为已发布"，
// 此前该执行器不存在，SCHEDULED 文章会被无声搁置。
//
// 安全（与 search.reindex / social.distribution.publish 同款模型）：
//   - payload 只含 id，worker 从 DB 取最新状态（SystemViewer 跨租户读）
//   - worker 幂等：状态已非 SCHEDULED 直接跳过，重复投递/双击无副作用
//   - 改期后旧任务到点自检 publish_time 仍在未来则自弃，由改期时新入队任务执行
//   - 发布成功后由 service 层钩子入队 ES/AI 索引重写（复用既有双写链）
// ============================================================================

// PostPublishTaskType 文章定时发布任务的 asynq 任务类型。
const PostPublishTaskType = "post.publish"

// PostPublishPayload 文章定时发布任务 payload。
//
// TenantID 仅供日志/排查展示，worker 以 DB 记录为准。
type PostPublishPayload struct {
	PostID   uint32 `json:"post_id"`
	TenantID uint32 `json:"tenant_id"`
}
