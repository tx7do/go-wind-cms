package task

// ============================================================================
// 搜索重索引任务类型定义
//
// 用于 PostgreSQL → OpenSearch 的双写同步：post/post_translation 的
// Create/Update/Delete 在 DB 事务提交后，入队一个 search.reindex 任务，
// asynq worker 收到后从 DB 取最新数据写入/删除 ES 文档。
//
// 安全：
//   - payload 只含 id，不含文档内容，体积小、无敏感字段
//   - 实际文档内容由 worker 从 DB 取（带 SystemViewer 跨租户读），写入 ES 的
//     tenant_id 取自 DB 记录字段，非 payload
//   - 失败由 asynq 自动重试；崩溃窗口内漏掉的文档由周期 ReindexAll 修复
// ============================================================================

const (
	// SearchReindexTaskType 单条重索引任务的 asynq 任务类型。
	SearchReindexTaskType = "search.reindex"

	// SearchReindexAllTaskType 周期全量重索引任务类型（自愈路径）。
	SearchReindexAllTaskType = "search.reindex.all"
)

// SearchReindexPayload 单条重索引任务的 payload。
//
// Op 取值：
//   - "index"  ：post 或其翻译被创建/更新，worker 从 DB 取最新数据 upsert 到 ES
//   - "delete" ：post 被删除，worker 按post_id+tenant_id 删除 ES 中所有语言文档
//
// Entity 取值："post"（v1 仅此一种）。
//
// TenantID 来自 DB 记录的 tenant_id（reindex 路径），用于 delete 操作时
// 双重定位删除；index 操作时 worker 仍从 DB 取真实 tenant_id 写入 ES。
type SearchReindexPayload struct {
	Entity   string `json:"entity"`
	ID       uint32 `json:"id"`
	TenantID uint32 `json:"tenant_id"`
	Op       string `json:"op"`
}
