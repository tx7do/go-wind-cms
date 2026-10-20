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
func NewAsynqServer(ctx *bootstrap.Context, taskService *service.TaskService, searchService *service.SearchService, aiKnowledgeService *service.AiKnowledgeService, aiQueryService *service.AiQueryService, aiContentService *service.AiContentService, socialService *service.SocialService, postService *service.PostService, webhookService *service.WebhookService, internalMessageService *service.InternalMessageService, flowService *service.FlowService, notificationService *service.NotificationService, monitorAlertService *service.MonitorAlertService, tenantQuotaService *service.TenantQuotaService) *asynq.Server {
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

	// 注册定时发布/社交分发兜底扫描订阅：周期任务每分钟入队（下方 NewPeriodicTask），
	// 订阅曾整体缺失——任务以 handler not found 无限重试堆积（retry 积压数千条），
	// 兜底扫描实际从未执行过。
	if err = asynq.RegisterSubscriber(srv, task.PostPublishSweepTaskType, postService.PublishSweepWorker); err != nil {
		log.Error(err)
	}
	if err = asynq.RegisterSubscriber(srv, task.SocialPublishSweepTaskType, socialService.PublishSweepWorker); err != nil {
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

	// 注册 Flow 延时续跑任务：动作链遇 delay 动作后由该 DELAY 任务接续剩余动作
	//（同一 handler，按 payload.ActionIndex 分流；任务类型不能以 flow.execute 为
	// 前缀，见 pkg/task/flow.go 前缀路由约定）。
	if err = asynq.RegisterSubscriber(srv, task.FlowResumeTaskType, flowService.ExecuteWorker); err != nil {
		log.Error(err)
	}

	// Flow 定时规则调度 tick（每分钟）：解析启用规则的 cron，命中当前分钟的
	// 按其租户入队 flow.execute（Event=schedule、EntityID=0、不走扇出）。
	if err = asynq.RegisterSubscriber(srv, task.FlowScheduleTickTaskType, flowService.ScheduleTickWorker); err != nil {
		log.Error(err)
	}
	// 监控告警周期扫描：每 5 分钟评估一轮启用规则（告警/恢复/冷却重发内核见 monitor_alert_service）
	if _, err = srv.NewPeriodicTask(task.MonitorAlertScanCronSpec, task.MonitorAlertScanTaskType, &task.MonitorAlertScanTaskData{}); err != nil {
		return nil
	}
	if err = asynq.RegisterSubscriber(srv, task.MonitorAlertScanTaskType, monitorAlertService.AsyncMonitorAlertScan); err != nil {
		return nil
	}

	if _, err = srv.NewPeriodicTask("*/1 * * * *", task.FlowScheduleTickTaskType, &task.FlowScheduleTickPayload{}); err != nil {
		log.Error(err)
	}

	// 注册社交统计周期回拉（每 30 分钟）：批量刷新近 30 天内已发布分发的互动快照。
	if err = asynq.RegisterSubscriber(srv, task.SocialStatsSyncTaskType, socialService.StatsSyncWorker); err != nil {
		log.Error(err)
	}
	if _, err = srv.NewPeriodicTask("*/30 * * * *", task.SocialStatsSyncTaskType, &task.StatsSyncPayload{}); err != nil {
		log.Error(err)
	}

	// 通知域异步派发：规则行标 is_async 的事件由该任务投递。带 ctx 版注册——
	// asynq 的截止时间要透传进 mailer 的 DialContext 才能真掐断 SMTP 握手。
	notificationService.RegisterTaskEnqueuer(taskService)
	if err = asynq.RegisterSubscriberWithCtx(srv, task.NotificationDispatchTaskType, notificationService.AsyncNotificationDispatch); err != nil {
		log.Error(err)
	}

	// 通知台账清扫（每 5 分钟）：超期停留在 SENDING 的行定案 FAILED（只定案不补投）。
	// 系统级常驻任务，不进 sys_tasks（任务管理页不可见不可停）。
	if err = asynq.RegisterSubscriber(srv, task.NotificationDeliverySweepTaskType, notificationService.AsyncDeliverySweep); err != nil {
		log.Error(err)
	}
	if _, err = srv.NewPeriodicTask(task.NotificationDeliverySweepCronSpec, task.NotificationDeliverySweepTaskType, &task.NotificationDeliverySweepTaskData{}); err != nil {
		log.Error(err)
	}

	// 租户配额水位扫描（每日 09:00）：只读扫描 + 提前告警（站内信投给租户管理员），
	// 与闸门硬执行互补。系统级常驻任务，不进 sys_tasks。
	if err = asynq.RegisterSubscriber(srv, task.TenantQuotaWatermarkTaskType, tenantQuotaService.AsyncTenantQuotaWatermarkScan); err != nil {
		log.Error(err)
	}
	if _, err = srv.NewPeriodicTask(task.TenantQuotaWatermarkCronSpec, task.TenantQuotaWatermarkTaskType, &task.TenantQuotaWatermarkTaskData{}); err != nil {
		log.Error(err)
	}

	// 启动所有的任务
	_, _ = taskService.StartAllTask(appViewer.NewSystemViewerContext(ctx.Context()), nil)

	return srv
}
