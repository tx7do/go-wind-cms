package server

import (
	"github.com/go-kratos/kratos/v2/log"

	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"github.com/tx7do/kratos-transport/transport/asynq"

	"go-wind-cms/app/core/service/internal/service"

	appViewer "go-wind-cms/pkg/entgo/viewer"
	"go-wind-cms/pkg/task"
)

// NewAsynqServer creates a new asynq server.
func NewAsynqServer(ctx *bootstrap.Context, taskService *service.TaskService, searchService *service.SearchService, aiKnowledgeService *service.AiKnowledgeService, aiQueryService *service.AiQueryService) *asynq.Server {
	cfg := ctx.GetConfig()

	if cfg == nil || cfg.Server == nil || cfg.Server.Asynq == nil {
		return nil
	}

	srv := asynq.NewServer(
		asynq.WithCodec(cfg.Server.Asynq.GetCodec()),
		asynq.WithRedisURI(cfg.Server.Asynq.GetUri()),
		asynq.WithLocation(cfg.Server.Asynq.GetLocation()),
		asynq.WithGracefullyShutdown(cfg.Server.Asynq.GetEnableGracefullyShutdown()),
		asynq.WithShutdownTimeout(cfg.Server.Asynq.GetShutdownTimeout().AsDuration()),
	)

	taskService.RegisterTaskScheduler(srv)

	var err error

	// 注册任务
	if err = asynq.RegisterSubscriber(srv, task.BackupTaskType, taskService.AsyncBackup); err != nil {
		log.Error(err)
	}

	// 注册搜索重索引任务订阅者。
	// worker 收到 search.reindex 任务后，调 SearchService.ReindexPost 从 DB 取
	// 最新数据写入/删除 ES。详见 search_service.go / search_repo.go 安全模型。
	if err = asynq.RegisterSubscriber(srv, task.SearchReindexTaskType, searchService.ReindexPost); err != nil {
		log.Error(err)
	}

	// 注册知识库切片重索引任务：更换 embedding 模型/provider 后全量重算向量。
	if err = asynq.RegisterSubscriber(srv, task.AiDocReindexTaskType, aiKnowledgeService.ReindexDocsWorker); err != nil {
		log.Error(err)
	}

	// 注册定时问数任务：执行问题 → 站内信推送结果给平台侧用户。
	if err = asynq.RegisterSubscriber(srv, task.AiQueryRunTaskType, aiQueryService.AsyncAiQueryRun); err != nil {
		log.Error(err)
	}

	// 启动所有的任务
	_, _ = taskService.StartAllTask(appViewer.NewSystemViewerContext(ctx.Context()), nil)

	return srv
}
