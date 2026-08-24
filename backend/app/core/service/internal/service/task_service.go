package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/hibiken/asynq"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"go-wind-cms/app/core/service/internal/data"

	taskV1 "go-wind-cms/api/gen/go/task/service/v1"

	"go-wind-cms/pkg/task"
)

// TaskScheduler 任务调度接口
type TaskScheduler interface {
	TaskTypeExists(taskType string) bool
	GetRegisteredTaskTypes() []string

	NewTask(typeName string, msg any, opts ...asynq.Option) error
	NewWaitResultTask(typeName string, msg any, opts ...asynq.Option) error
	NewPeriodicTask(cronSpec, typeName string, msg any, opts ...asynq.Option) (string, error)

	RemovePeriodicTask(id string) error
	RemoveAllPeriodicTask()
}

// TaskService 任务服务
type TaskService struct {
	taskV1.UnimplementedTaskServiceServer

	log *log.Helper

	taskScheduler TaskScheduler

	userRepo data.UserRepo
	taskRepo *data.TaskRepo

	// asynq Redis URI（与 asynq server 同源），用于懒建只读 Inspector 查执行记录
	asynqRedisURI string
	asynqInspector *asynq.Inspector
}

func NewTaskService(
	ctx *bootstrap.Context,
	taskRepo *data.TaskRepo,
	userRepo data.UserRepo,
) *TaskService {
	svc := &TaskService{
		log:      ctx.NewLoggerHelper("task/service/core-service"),
		taskRepo: taskRepo,
		userRepo: userRepo,
	}

	if cfg := ctx.GetConfig(); cfg != nil && cfg.Server != nil && cfg.Server.Asynq != nil {
		svc.asynqRedisURI = cfg.Server.Asynq.GetUri()
	}

	return svc
}

func (s *TaskService) RegisterTaskScheduler(taskScheduler TaskScheduler) {
	s.taskScheduler = taskScheduler
}

// getAsynqInspector 懒建只读 Inspector（与 worker 同一 Redis，来自 server.yaml 同源配置）。
func (s *TaskService) getAsynqInspector() (*asynq.Inspector, error) {
	if s.asynqInspector != nil {
		return s.asynqInspector, nil
	}
	if s.asynqRedisURI == "" {
		return nil, taskV1.ErrorInternalServerError("asynq redis uri is not configured")
	}
	connOpt, err := asynq.ParseRedisURI(s.asynqRedisURI)
	if err != nil {
		s.log.Errorf("parse asynq redis uri failed: %s", err.Error())
		return nil, taskV1.ErrorInternalServerError("parse asynq redis uri failed")
	}
	s.asynqInspector = asynq.NewInspector(connOpt)
	return s.asynqInspector, nil
}

// ListTaskExecutions 查询近期的任务执行实例（完成/归档失败/进行中/待处理），
// 数据来自 asynq Inspector。完成态深度受 retention 约束（默认 24h，见
// convertTaskOption），归档态由 asynq janitor 保留（约 7 天）。按 type_name
// 过滤后在内存分页；Inspector 不支持按类型过滤，故每状态拉取分页窗口量再合并。
func (s *TaskService) ListTaskExecutions(ctx context.Context, req *taskV1.ListTaskExecutionsRequest) (*taskV1.ListTaskExecutionsResponse, error) {
	inspector, err := s.getAsynqInspector()
	if err != nil {
		return nil, err
	}

	const queue = "default"

	fetchSize := int(req.GetPageSize()) * int(req.GetPage())
	if fetchSize <= 0 {
		fetchSize = 30
	}
	if fetchSize > 200 {
		fetchSize = 200
	}
	listOpts := []asynq.ListOption{asynq.Page(1), asynq.PageSize(fetchSize)}

	var infos []*asynq.TaskInfo
	for _, list := range []func(string, ...asynq.ListOption) ([]*asynq.TaskInfo, error){
		inspector.ListCompletedTasks,
		inspector.ListArchivedTasks,
		inspector.ListActiveTasks,
		inspector.ListPendingTasks,
	} {
		items, listErr := list(queue, listOpts...)
		if listErr != nil {
			// 单状态查询失败不致命（如队列键尚未创建），跳过继续
			s.log.Warnf("list asynq tasks (state query) failed: %s", listErr.Error())
			continue
		}
		infos = append(infos, items...)
	}

	// 按任务类型过滤
	if typeName := req.GetTypeName(); typeName != "" {
		filtered := infos[:0]
		for _, info := range infos {
			if info.Type == typeName {
				filtered = append(filtered, info)
			}
		}
		infos = filtered
	}

	total := uint64(len(infos))

	// 内存分页
	page := int(req.GetPage())
	if page <= 0 {
		page = 1
	}
	pageSize := int(req.GetPageSize())
	if pageSize <= 0 {
		pageSize = 10
	}
	start := (page - 1) * pageSize
	if start > len(infos) {
		start = len(infos)
	}
	end := start + pageSize
	if end > len(infos) {
		end = len(infos)
	}

	items := make([]*taskV1.TaskExecution, 0, end-start)
	for _, info := range infos[start:end] {
		item := &taskV1.TaskExecution{
			Id:        trans.Ptr(info.ID),
			TypeName:  trans.Ptr(info.Type),
			State:     trans.Ptr(info.State.String()),
			Payload:   trans.Ptr(string(info.Payload)),
			Retried:   trans.Ptr(uint32(info.Retried)),
			LastError: trans.Ptr(info.LastErr),
		}
		if !info.LastFailedAt.IsZero() {
			item.LastFailedAt = timestamppb.New(info.LastFailedAt)
		}
		if !info.CompletedAt.IsZero() {
			item.CompletedAt = timestamppb.New(info.CompletedAt)
		}
		if !info.NextProcessAt.IsZero() {
			item.NextProcessAt = timestamppb.New(info.NextProcessAt)
		}
		items = append(items, item)
	}

	return &taskV1.ListTaskExecutionsResponse{
		Items: items,
		Total: total,
	}, nil
}

func (s *TaskService) List(ctx context.Context, req *paginationV1.PagingRequest) (*taskV1.ListTaskResponse, error) {
	return s.taskRepo.List(ctx, req)
}

func (s *TaskService) Get(ctx context.Context, req *taskV1.GetTaskRequest) (*taskV1.Task, error) {
	return s.taskRepo.Get(ctx, req)
}

func (s *TaskService) ListTaskTypeName(_ context.Context, _ *emptypb.Empty) (*taskV1.ListTaskTypeNameResponse, error) {
	// taskScheduler 未注入时返回不可用，避免 nil 解引用 panic
	if s.taskScheduler == nil {
		return nil, taskV1.ErrorServiceUnavailable("task scheduler is not available")
	}
	typeNames := s.taskScheduler.GetRegisteredTaskTypes()
	return &taskV1.ListTaskTypeNameResponse{
		TypeNames: typeNames,
	}, nil
}

func (s *TaskService) Create(ctx context.Context, req *taskV1.CreateTaskRequest) (*emptypb.Empty, error) {
	if req == nil || req.Data == nil {
		return nil, taskV1.ErrorBadRequest("invalid parameter")
	}

	// 任务类型与命名必须由服务端校验，避免空 type_name/cron_spec 进入调度器导致
	// 注册失败或以空名注册互相覆盖
	if err := validateTaskFields(req.Data); err != nil {
		return nil, err
	}

	var t *taskV1.Task
	var err error
	if t, err = s.taskRepo.Create(ctx, req); err != nil {
		return nil, err
	}

	if err = s.startTask(t); err != nil {
		s.log.Error(err)
	}

	return &emptypb.Empty{}, nil
}

func (s *TaskService) Update(ctx context.Context, req *taskV1.UpdateTaskRequest) (*emptypb.Empty, error) {
	if req == nil || req.Data == nil {
		return nil, taskV1.ErrorBadRequest("invalid parameter")
	}

	// 任务类型与命名必须由服务端校验，避免空 type_name/cron_spec 进入调度器导致
	// 注册失败或以空名注册互相覆盖
	if err := validateTaskFields(req.Data); err != nil {
		return nil, err
	}

	// 获取更新前的任务，用于判断调度器中是否有正在运行的注册项
	oldTask, _ := s.taskRepo.Get(ctx, &taskV1.GetTaskRequest{QueryBy: &taskV1.GetTaskRequest_Id{Id: req.GetId()}})

	var t *taskV1.Task
	var err error
	if t, err = s.taskRepo.Update(ctx, req); err != nil {

		return nil, err
	}

	// 更新结果为空时不可继续访问其字段，避免 nil 解引用
	if t == nil {
		return &emptypb.Empty{}, nil
	}

	// 先移除调度器中旧的注册项（若存在），避免停用后仍运行、或启用时重复注册。
	// 直接调用 RemovePeriodicTask 绕过 stopTask 内部的 enable 保护，
	// 因为 oldTask 可能已是禁用状态，但其注册项可能仍残留在调度器中。
	if oldTask != nil && oldTask.GetType() == taskV1.Task_PERIODIC && oldTask.GetTypeName() != "" {
		if removeErr := s.taskScheduler.RemovePeriodicTask(oldTask.GetTypeName()); removeErr != nil {
			s.log.Warnf("移除旧定时任务注册项失败[%s]: %v", oldTask.GetTypeName(), removeErr)
		}
	}

	// 根据更新后的 enable 状态决定是否重新启动
	if err = s.startTask(t); err != nil {
		s.log.Error(err)
	}

	return &emptypb.Empty{}, nil
}

// validateTaskFields 校验任务必填字段：type_name 始终必填；定时任务另需 cron_spec。
// 避免空值进入调度器导致注册失败或以空名互相覆盖。
func validateTaskFields(t *taskV1.Task) error {
	if t == nil {
		return taskV1.ErrorBadRequest("invalid parameter")
	}
	if t.GetTypeName() == "" {
		return taskV1.ErrorBadRequest("type_name is required")
	}
	if t.GetType() == taskV1.Task_PERIODIC && t.GetCronSpec() == "" {
		return taskV1.ErrorBadRequest("cron_spec is required for periodic task")
	}
	return nil
}

func (s *TaskService) Delete(ctx context.Context, req *taskV1.DeleteTaskRequest) (*emptypb.Empty, error) {
	if req == nil {
		return nil, taskV1.ErrorBadRequest("invalid request")
	}
	var err error
	var t *taskV1.Task
	if t, err = s.taskRepo.Get(ctx, &taskV1.GetTaskRequest{QueryBy: &taskV1.GetTaskRequest_Id{Id: req.GetId()}}); err != nil {
		return nil, err
	}

	if err = s.taskRepo.Delete(ctx, req); err != nil {
		return nil, err
	}

	if t != nil {
		_ = s.stopTask(t)
	}

	return &emptypb.Empty{}, nil
}

// ControlTask 控制调度任务
func (s *TaskService) ControlTask(ctx context.Context, req *taskV1.ControlTaskRequest) (*emptypb.Empty, error) {
	t, err := s.taskRepo.Get(ctx, &taskV1.GetTaskRequest{QueryBy: &taskV1.GetTaskRequest_TypeName{TypeName: req.GetTypeName()}})
	if err != nil {
		s.log.Errorf("获取任务失败[%s]", err.Error())
		return nil, err
	}

	switch req.GetControlType() {
	case taskV1.ControlTaskRequest_Restart:
		if err = s.stopTask(t); err != nil {
			return nil, err
		}

		if err = s.startTask(t); err != nil {
			return nil, err
		}

	case taskV1.ControlTaskRequest_Stop:
		err = s.stopTask(t)
		return nil, err

	case taskV1.ControlTaskRequest_Start:
		err = s.startTask(t)
		return nil, err

	default:
		return nil, taskV1.ErrorBadRequest("unknown control type: %v", req.GetControlType())
	}

	return &emptypb.Empty{}, nil
}

// StopAllTask 停止所有的调度任务
func (s *TaskService) StopAllTask(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	if s.taskScheduler == nil {
		return nil, taskV1.ErrorServiceUnavailable("task scheduler is not available")
	}
	s.stopAllTask(ctx)
	return &emptypb.Empty{}, nil
}

// StartAllTask 启动所有的调度任务
func (s *TaskService) StartAllTask(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	_, err := s.startAllTask(ctx)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// RestartAllTask 重启所有的调度任务
func (s *TaskService) RestartAllTask(ctx context.Context, _ *emptypb.Empty) (*taskV1.RestartAllTaskResponse, error) {
	if s.taskScheduler == nil {
		return nil, taskV1.ErrorServiceUnavailable("task scheduler is not available")
	}
	// 停止当前上下文（租户）范围内的任务
	s.stopAllTask(ctx)

	// 重新启动当前上下文（租户）范围内的任务
	count, err := s.startAllTask(ctx)

	return &taskV1.RestartAllTaskResponse{
		Count: count,
	}, err
}

// StartAllTask 启动所有的任务
func (s *TaskService) startAllTask(ctx context.Context) (int32, error) {
	//_, _ = s.asynqServer.NewPeriodicTask("*/1 * * * ?", task.BackupTaskType, task.BackupTaskData{Name: "test"})

	resp, err := s.List(ctx, &paginationV1.PagingRequest{
		NoPaging: trans.Ptr(true),
	})
	if err != nil {
		s.log.Errorf("获取任务列表失败[%s]", err.Error())
		return 0, err
	}

	s.log.Infof("开始开启定时任务，总计[%d]个", resp.GetTotal())

	// 重新启动任务
	var count int32
	for _, t := range resp.GetItems() {
		if s.startTask(t) != nil {
			continue
		} else {
			count++
		}
	}

	s.log.Infof("总共成功开启定时任务[%d]个", count)

	return count, nil
}

// stopAllTask 停止当前上下文（租户）范围内的任务。
// 仅遍历经 List（EvalQuery 已按调用者租户过滤）返回的任务，按 type_name 逐个停止，
// 避免调用 RemoveAllPeriodicTask() 清空其它租户的任务。
func (s *TaskService) stopAllTask(ctx context.Context) {
	// taskScheduler 未注入时跳过，避免 nil 解引用 panic
	if s.taskScheduler == nil {
		s.log.Warnf("task scheduler is not available, skip stopAllTask")
		return
	}

	resp, err := s.List(ctx, &paginationV1.PagingRequest{
		NoPaging: trans.Ptr(true),
	})
	if err != nil {
		s.log.Errorf("获取任务列表失败[%s]", err.Error())
		return
	}

	s.log.Infof("开始停止当前租户定时任务，总计[%d]个", resp.GetTotal())

	var count int32
	for _, t := range resp.GetItems() {
		if s.stopTask(t) != nil {
			continue
		} else {
			count++
		}
	}

	s.log.Infof("总共成功停止定时任务[%d]个", count)
}

// stopTask 停止一个任务
func (s *TaskService) stopTask(t *taskV1.Task) error {
	if t == nil {
		return errors.New("task is nil")
	}

	// taskScheduler 未注入时返回不可用，避免 nil 解引用 panic
	if s.taskScheduler == nil {
		return taskV1.ErrorServiceUnavailable("task scheduler is not available")
	}

	if t.GetEnable() == false {
		return errors.New("task is not enable")
	}

	switch t.GetType() {
	case taskV1.Task_PERIODIC:
		return s.taskScheduler.RemovePeriodicTask(t.GetTypeName())

	case taskV1.Task_DELAY:

	case taskV1.Task_WAIT_RESULT:
	}

	return nil
}

// convertTaskOption 转换任务选项
func (s *TaskService) convertTaskOption(t *taskV1.Task) (opts []asynq.Option, payload any) {
	if t == nil {
		return
	}

	if len(t.GetTaskPayload()) > 0 {
		_ = json.Unmarshal([]byte(t.GetTaskPayload()), &payload)
	}

	if t.TaskOptions != nil {
		if t.GetTaskOptions().GetMaxRetry() > 0 {
			opts = append(opts, asynq.MaxRetry(int(t.GetTaskOptions().GetMaxRetry())))
		}
		if t.GetTaskOptions().Timeout != nil {
			opts = append(opts, asynq.Timeout(t.GetTaskOptions().GetTimeout().AsDuration()))
		}
		if t.GetTaskOptions().Deadline != nil {
			opts = append(opts, asynq.Deadline(t.GetTaskOptions().GetDeadline().AsTime()))
		}
		if t.GetTaskOptions().ProcessIn != nil {
			opts = append(opts, asynq.ProcessIn(t.GetTaskOptions().GetProcessIn().AsDuration()))
		}
		if t.GetTaskOptions().ProcessAt != nil {
			opts = append(opts, asynq.ProcessAt(t.GetTaskOptions().GetProcessAt().AsTime()))
		}
		if t.GetTaskOptions().UniqueTtl != nil {
			opts = append(opts, asynq.Unique(t.GetTaskOptions().GetUniqueTtl().AsDuration()))
		}
		if t.GetTaskOptions().Retention != nil {
			opts = append(opts, asynq.Retention(t.GetTaskOptions().GetRetention().AsDuration()))
		} else {
			// 默认保留 24h：无 retention 的完成任务会被 asynq 立即删除，
			// ListTaskExecutions 的完成态执行记录将无数据可读
			opts = append(opts, asynq.Retention(24*time.Hour))
		}
		if t.GetTaskOptions().Group != nil {
			opts = append(opts, asynq.Group(t.GetTaskOptions().GetGroup()))
		}
		if t.GetTaskOptions().TaskId != nil {
			opts = append(opts, asynq.TaskID(t.GetTaskOptions().GetTaskId()))
		}
	}

	return
}

// startTask 启动一个任务
func (s *TaskService) startTask(t *taskV1.Task) error {
	if t == nil {
		return errors.New("task is nil")
	}

	// taskScheduler 未注入时返回不可用，避免 nil 解引用 panic
	if s.taskScheduler == nil {
		return taskV1.ErrorServiceUnavailable("task scheduler is not available")
	}

	if t.GetEnable() == false {
		return errors.New("task is not enable")
	}

	var opts []asynq.Option
	var payload any
	var err error

	switch t.GetType() {
	case taskV1.Task_PERIODIC:
		opts, payload = s.convertTaskOption(t)
		if _, err = s.taskScheduler.NewPeriodicTask(t.GetCronSpec(), t.GetTypeName(), payload, opts...); err != nil {
			s.log.Errorf("[%s] 创建定时任务失败[%s]", t.GetTypeName(), err.Error())
			return err
		}

	case taskV1.Task_DELAY:
		opts, payload = s.convertTaskOption(t)
		if err = s.taskScheduler.NewTask(t.GetTypeName(), payload, opts...); err != nil {
			s.log.Errorf("[%s] 创建延迟任务失败[%s]", t.GetTypeName(), err.Error())
			return err
		}

	case taskV1.Task_WAIT_RESULT:
		opts, payload = s.convertTaskOption(t)
		if err = s.taskScheduler.NewWaitResultTask(t.GetTypeName(), payload, opts...); err != nil {
			s.log.Errorf("[%s] 创建等待结果任务失败[%s]", t.GetTypeName(), err.Error())
			return err
		}
	}

	return nil
}

// AsyncBackup 异步备份
func (s *TaskService) AsyncBackup(taskType string, taskData *task.BackupTaskData) error {
	s.log.Infof("AsyncBackup [%s] [%+v] [%s]", taskType, taskData, taskData.Name)
	return nil
}

// EnqueueSearchReindex 入队一个 OpenSearch 单条重索引任务。
//
// 由 PostService / PostTranslationRepo 在 DB 事务提交成功后调用，
// 用于把 post / post_translation 的变更同步到 ES。
//
// 入队是 best-effort：失败仅记日志、不回滚 DB。崩溃窗口内漏掉的文档由
// SearchService.ReindexAll 周期任务修复（最长 1 小时滞后）。
//
// 安全：payload 只含 id，不含文档内容；ES 文档 tenant_id 由 worker 从 DB
// 记录取（GetReindexDocuments 内部从 ent.Post.TenantID 取），非 payload。
func (s *TaskService) EnqueueSearchReindex(payload *task.SearchReindexPayload) error {
	if payload == nil {
		return errors.New("nil search reindex payload")
	}
	if s.taskScheduler == nil {
		// 调度器未注入（如开发环境未起 asynq）→ 静默跳过，不阻断主业务
		s.log.Warnf("search reindex skipped: task scheduler not available (entity=%s id=%d op=%s)",
			payload.Entity, payload.ID, payload.Op)
		return nil
	}
	if err := s.taskScheduler.NewTask(task.SearchReindexTaskType, payload); err != nil {
		s.log.Errorf("enqueue search reindex failed (entity=%s id=%d op=%s): %v",
			payload.Entity, payload.ID, payload.Op, err)
		return err
	}
	s.log.Infof("enqueued search reindex (entity=%s id=%d op=%s)",
		payload.Entity, payload.ID, payload.Op)
	return nil
}
