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
func NewAsynqServer(ctx *bootstrap.Context, taskService *service.TaskService, searchService *service.SearchService, aiKnowledgeService *service.AiKnowledgeService, aiQueryService *service.AiQueryService, aiContentService *service.AiContentService, socialService *service.SocialService, postService *service.PostService, webhookService *service.WebhookService, internalMessageService *service.InternalMessageService, flowService *service.FlowService) *asynq.Server {
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

	// 注册帖子语义索引任务：帖子增删改后同步 AI 全局搜索索引。
	if err = asynq.RegisterSubscriber(srv, task.AiPostIndexTaskType, aiContentService.IndexPostWorker); err != nil {
		log.Error(err)
	}

	// 注册社交分发发布任务：到点发布到社交平台并回写状态（详见 social_service.go 状态机）。
	if err = asynq.RegisterSubscriber(srv, task.SocialDistributionPublishTaskType, socialService.PublishDistributionWorker); err != nil {
		log.Error(err)
	}

	// 注册社交收件箱回拉任务：回复/提及落库（详见 social_service.go InboxSyncWorker）。
	if err = asynq.RegisterSubscriber(srv, task.SocialInboxSyncTaskType, socialService.InboxSyncWorker); err != nil {
		log.Error(err)
	}

	// 注册文章定时发布任务：到点把 SCHEDULED 文章转正为 PUBLISHED 并刷新索引
	//（修复前端承诺"到点自动发布"但后端无执行器的既有缺口）。
	if err = asynq.RegisterSubscriber(srv, task.PostPublishTaskType, postService.PublishScheduledPostWorker); err != nil {
		log.Error(err)
	}

	// 周期兜底扫描（每分钟）：DELAY 任务在 Redis 丢失/崩溃窗口内可能静默消失，
	// 到期未处理的 SCHEDULED 记录由此重新入队；worker 幂等，重复入队无害。
	if _, err = srv.NewPeriodicTask("*/1 * * * *", task.PostPublishSweepTaskType, &task.PublishSweepPayload{}); err != nil {
		log.Error(err)
	}
	if _, err = srv.NewPeriodicTask("*/1 * * * *", task.SocialPublishSweepTaskType, &task.PublishSweepPayload{}); err != nil {
		log.Error(err)
	}
	if _, err = srv.NewPeriodicTask("*/15 * * * *", task.SocialInboxSyncTaskType, &task.PublishSweepPayload{}); err != nil {
		log.Error(err)
	}

	// 注册邮件通知投递任务：系统通知的附加邮件触达（SMTP 未配置时 worker 内跳过）。
	if err = asynq.RegisterSubscriber(srv, task.NotifyEmailTaskType, internalMessageService.SendEmailWorker); err != nil {
		log.Error(err)
	}

	// 注册 Webhook 投递任务：fanout 按事件扇出到单订阅 deliver（各自由 asynq 退避重试）。
	if err = asynq.RegisterSubscriber(srv, task.WebhookFanoutTaskType, webhookService.FanoutWorker); err != nil {
		log.Error(err)
	}
	if err = asynq.RegisterSubscriber(srv, task.WebhookDeliverTaskType, webhookService.DeliverWorker); err != nil {
		log.Error(err)
	}
	if err = asynq.RegisterSubscriber(srv, task.FlowExecuteTaskType, flowService.ExecuteWorker); err != nil {
		log.Error(err)
	}

	// 注册社交统计周期回拉（每 30 分钟）：批量刷新近 30 天内已发布分发的互动快照。
	if err = asynq.RegisterSubscriber(srv, task.SocialStatsSyncTaskType, socialService.StatsSyncWorker); err != nil {
		log.Error(err)
	}
	if _, err = srv.NewPeriodicTask("*/30 * * * *", task.SocialStatsSyncTaskType, &task.StatsSyncPayload{}); err != nil {
		log.Error(err)
	}

	// 启动所有的任务
	_, _ = taskService.StartAllTask(appViewer.NewSystemViewerContext(ctx.Context()), nil)

	return srv
}
