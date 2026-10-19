package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
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
	asynqRedisURI  string
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

// NewTask 通用入队口：调用方自带任务类型/载荷/选项（通知域异步派发等）。
// 调度器未注册（asynq 未配置）时返回错误，调用方据此退回同步路径。
func (s *TaskService) NewTask(typeName string, msg any, opts ...asynq.Option) error {
	if s.taskScheduler == nil {
		return fmt.Errorf("task scheduler is not registered")
	}
	return s.taskScheduler.NewTask(typeName, msg, opts...)
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

// systemTaskInspectMaxFailures 单任务类型保留的最近失败条数上限（防一张表被历史失败刷满）。
const systemTaskInspectMaxFailures = 10

// systemTaskQueue asynq 队列名：与 NewPeriodicTask/ListTaskExecutions 同口径（未配置多队列路由）。
const systemTaskQueue = "default"

// InspectSystemTasks 一站式巡检系统级常驻任务（只读 asynq Inspector，无副作用）。
//
// 方案 C'：零新表。调度计划与队列状态全部来自 asynq 自身——
//   - SchedulerEntries：cron 表达式 + Next/Prev 入队时间（Prev 零值 = 从未跑过）；
//   - 队列的 active/pending/retry/archived 计数与失败明细（LastErr/LastFailedAt）。
//
// asynq 给不了的：成功执行的持久历史（成功即出队消失，那在 ListTaskExecutions，
// 深度受 retention 约束）与领域级结果（通知投递表/审计台账）。
func (s *TaskService) InspectSystemTasks(ctx context.Context, _ *taskV1.InspectSystemTasksRequest) (*taskV1.InspectSystemTasksResponse, error) {
	inspector, err := s.getAsynqInspector()
	if err != nil {
		return nil, err
	}

	resp := &taskV1.InspectSystemTasksResponse{Queue: systemTaskQueue}

	entries, err := inspector.SchedulerEntries()
	if err != nil {
		s.log.Errorf("task monitor inspect: list scheduler entries failed: %s", err.Error())
		return nil, taskV1.ErrorInternalServerError("list scheduler entries failed")
	}
	for _, e := range entries {
		resp.Schedules = append(resp.Schedules, scheduleEntryToDTO(e))
	}

	// 只巡检调度条目覆盖到的任务类型：队列里可能躺着别的项目/别的类型的任务
	// （队列无命名空间，共库时以此过滤），未知类型的遗留任务在失败明细里仍可见。
	monitored := make(map[string]bool, len(entries))
	for _, e := range entries {
		monitored[e.Task.Type()] = true
	}

	summary, err := s.summarizeQueues(inspector, monitored)
	if err != nil {
		s.log.Errorf("task monitor inspect: summarize queues failed: %s", err.Error())
		return nil, taskV1.ErrorInternalServerError("inspect queues failed")
	}
	resp.Summaries = summary

	return resp, nil
}

// summarizeQueues 按任务类型聚合 active/pending/retry/archived 计数与最近失败。
func (s *TaskService) summarizeQueues(inspector *asynq.Inspector, monitored map[string]bool) ([]*taskV1.SystemTaskStateSummary, error) {
	byType := map[string]*taskV1.SystemTaskStateSummary{}
	get := func(t string) *taskV1.SystemTaskStateSummary {
		if summary, ok := byType[t]; ok {
			return summary
		}
		summary := &taskV1.SystemTaskStateSummary{TaskType: t}
		byType[t] = summary
		return summary
	}

	// 队列计数总量按 TaskInfo 逐条归到类型上：把四种状态的每条任务各读一遍太重，
	// 这里用 ListTasks 首页（前 N 条）组合——计数只做"有没有存货"的量级提示，
	// 失败明细才是运营要看的主体。
	listOpts := []asynq.ListOption{asynq.Page(1), asynq.PageSize(systemTaskInspectMaxFailures * 10)}

	active, err := inspector.ListActiveTasks(systemTaskQueue, listOpts...)
	if err != nil {
		return nil, err
	}
	for _, ti := range active {
		get(ti.Type).Active++
	}
	pending, err := inspector.ListPendingTasks(systemTaskQueue, listOpts...)
	if err != nil {
		return nil, err
	}
	for _, ti := range pending {
		get(ti.Type).Pending++
	}
	retry, err := inspector.ListRetryTasks(systemTaskQueue, listOpts...)
	if err != nil {
		return nil, err
	}
	for _, ti := range retry {
		summary := get(ti.Type)
		summary.Retry++
		if monitored[ti.Type] || len(monitored) == 0 {
			summary.RecentFailures = append(summary.RecentFailures, failureToDTO(ti, "retry"))
		}
	}
	archived, err := inspector.ListArchivedTasks(systemTaskQueue, listOpts...)
	if err != nil {
		return nil, err
	}
	for _, ti := range archived {
		summary := get(ti.Type)
		summary.Archived++
		if monitored[ti.Type] || len(monitored) == 0 {
			summary.RecentFailures = append(summary.RecentFailures, failureToDTO(ti, "archived"))
		}
	}

	for _, summary := range byType {
		sortFailures(summary.RecentFailures)
		if len(summary.RecentFailures) > systemTaskInspectMaxFailures {
			summary.RecentFailures = summary.RecentFailures[:systemTaskInspectMaxFailures]
		}
	}

	// 稳定输出：按任务类型名排序
	result := make([]*taskV1.SystemTaskStateSummary, 0, len(byType))
	for _, summary := range byType {
		result = append(result, summary)
	}
	sort.SliceStable(result, func(i, j int) bool {
		return result[i].GetTaskType() < result[j].GetTaskType()
	})
	return result, nil
}

// ==== 巡检 DTO 映射/排序 helpers ====

func scheduleEntryToDTO(e *asynq.SchedulerEntry) *taskV1.SystemTaskScheduleEntry {
	entry := &taskV1.SystemTaskScheduleEntry{
		TaskType: e.Task.Type(),
		CronSpec: e.Spec,
	}
	if !e.Next.IsZero() {
		entry.NextEnqueueAt = timestamppb.New(e.Next)
	}
	// Prev 零值时间（从未入队）不转：proto 里保持空，前端按"从未运行"渲染
	if !e.Prev.IsZero() {
		entry.PrevEnqueueAt = timestamppb.New(e.Prev)
	}
	return entry
}

func failureToDTO(ti *asynq.TaskInfo, state string) *taskV1.SystemTaskFailure {
	failure := &taskV1.SystemTaskFailure{
		TaskType:  ti.Type,
		Queue:     ti.Queue,
		State:     state,
		LastError: ti.LastErr,
		Retried:   uint32(ti.Retried),
		MaxRetry:  uint32(ti.MaxRetry),
	}
	if !ti.LastFailedAt.IsZero() {
		failure.LastFailedAt = timestamppb.New(ti.LastFailedAt)
	}
	if !ti.NextProcessAt.IsZero() {
		failure.NextProcessAt = timestamppb.New(ti.NextProcessAt)
	}
	return failure
}

func failureTime(f *taskV1.SystemTaskFailure) *time.Time {
	ts := f.GetLastFailedAt()
	if ts == nil {
		return nil
	}
	t := ts.AsTime()
	return &t
}

func sortFailures(failures []*taskV1.SystemTaskFailure) {
	// 按失败时间倒序（零值排最后）
	sort.SliceStable(failures, func(i, j int) bool {
		a := failureTime(failures[i])
		b := failureTime(failures[j])
		if a == nil {
			return false
		}
		if b == nil {
			return true
		}
		return a.After(*b)
	})
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

// EnqueueAiDocReindex 入队一个知识库切片重索引任务（baseId=0 全量）。
// 入队 best-effort：调度器未注入（开发环境未起 asynq）时静默跳过。
// EnqueueAiPostIndex 入队一个帖子语义索引任务（op=index/delete）。
// 入队 best-effort：调度器未注入（开发环境未起 asynq）时静默跳过。
func (s *TaskService) EnqueueAiPostIndex(payload *task.AiPostIndexPayload) error {
	if payload == nil || payload.PostID == 0 {
		return nil
	}
	if s.taskScheduler == nil {
		s.log.Warnf("ai post index skipped: task scheduler not available (post=%d op=%s)",
			payload.PostID, payload.Op)
		return nil
	}
	if err := s.taskScheduler.NewTask(task.AiPostIndexTaskType, payload); err != nil {
		s.log.Errorf("enqueue ai post index failed (post=%d op=%s): %v",
			payload.PostID, payload.Op, err)
		return err
	}
	s.log.Infof("enqueued ai post index (post=%d op=%s)", payload.PostID, payload.Op)
	return nil
}

// EnqueueNotifyEmail 入队一个邮件通知投递任务（best-effort：调度器未注入或
// SMTP 未配置时静默跳过——邮件是站内信之外的附加触达，不阻断通知主链路）。
func (s *TaskService) EnqueueNotifyEmail(payload *task.NotifyEmailPayload) error {
	if payload == nil || payload.UserID == 0 {
		return nil
	}
	if s.taskScheduler == nil {
		s.log.Warnf("notify email skipped: task scheduler not available (user=%d)", payload.UserID)
		return nil
	}
	if err := s.taskScheduler.NewTask(task.NotifyEmailTaskType, payload, asynq.MaxRetry(3)); err != nil {
		s.log.Errorf("enqueue notify email failed (user=%d): %v", payload.UserID, err)
		return err
	}
	return nil
}

// EnqueueFlowExecute 入队一个 Flow 规则执行任务（best-effort：调度器未注入时静默跳过）。
func (s *TaskService) EnqueueFlowExecute(payload *task.FlowExecutePayload) error {
	if payload == nil || payload.Event == "" {
		return nil
	}
	if s.taskScheduler == nil {
		s.log.Warnf("flow execute skipped: task scheduler not available (event=%s entity=%d)",
			payload.Event, payload.EntityID)
		return nil
	}
	if err := s.taskScheduler.NewTask(task.FlowExecuteTaskType, payload, asynq.MaxRetry(3)); err != nil {
		s.log.Errorf("enqueue flow execute failed (event=%s entity=%d): %v", payload.Event, payload.EntityID, err)
		return err
	}
	return nil
}

// EnqueueFlowResume 入队延时续跑任务（asynq DELAY 语义，ProcessIn 延迟投递）。
// 动作链遇 delay 动作时由 FlowService 调用：剩余动作由续跑任务接续。
func (s *TaskService) EnqueueFlowResume(payload *task.FlowExecutePayload, delay time.Duration) error {
	if payload == nil || payload.RuleID == 0 {
		return nil
	}
	if s.taskScheduler == nil {
		s.log.Warnf("flow resume skipped: task scheduler not available (rule=%d)", payload.RuleID)
		return nil
	}
	if err := s.taskScheduler.NewTask(task.FlowResumeTaskType, payload,
		asynq.MaxRetry(3), asynq.ProcessIn(delay)); err != nil {
		s.log.Errorf("enqueue flow resume failed (rule=%d): %v", payload.RuleID, err)
		return err
	}
	return nil
}

// EnqueueWebhookFanout 入队一个 Webhook 事件扇出任务（best-effort：调度器未注入时静默跳过，
// 内容保存主链路不因外发失败而受阻）。
func (s *TaskService) EnqueueWebhookFanout(payload *task.WebhookFanoutPayload) error {
	if payload == nil || payload.Event == "" {
		return nil
	}
	if s.taskScheduler == nil {
		s.log.Warnf("webhook fanout skipped: task scheduler not available (event=%s entity=%d)",
			payload.Event, payload.EntityID)
		return nil
	}
	if err := s.taskScheduler.NewTask(task.WebhookFanoutTaskType, payload, asynq.MaxRetry(3)); err != nil {
		s.log.Errorf("enqueue webhook fanout failed (event=%s entity=%d): %v", payload.Event, payload.EntityID, err)
		return err
	}
	return nil
}

// EnqueueWebhookDeliver 入队一个单订阅投递任务（MaxRetry 5 指数退避，见 DeliverWorker）。
func (s *TaskService) EnqueueWebhookDeliver(payload *task.WebhookDeliverPayload) error {
	if payload == nil || payload.WebhookID == 0 {
		return nil
	}
	if s.taskScheduler == nil {
		s.log.Warnf("webhook deliver skipped: task scheduler not available (webhook=%d event=%s)",
			payload.WebhookID, payload.Event)
		return nil
	}
	if err := s.taskScheduler.NewTask(task.WebhookDeliverTaskType, payload, asynq.MaxRetry(5)); err != nil {
		s.log.Errorf("enqueue webhook deliver failed (webhook=%d event=%s): %v", payload.WebhookID, payload.Event, err)
		return err
	}
	return nil
}

func (s *TaskService) EnqueueAiDocReindex(baseID uint32) error {
	if s.taskScheduler == nil {
		s.log.Warnf("ai doc reindex skipped: task scheduler not available (base=%d)", baseID)
		return nil
	}
	if err := s.taskScheduler.NewTask(task.AiDocReindexTaskType, &task.AiDocReindexPayload{BaseID: baseID}); err != nil {
		s.log.Errorf("enqueue ai doc reindex failed (base=%d): %v", baseID, err)
		return err
	}
	s.log.Infof("enqueued ai doc reindex (base=%d)", baseID)
	return nil
}

// EnqueueSocialDistributionPublish 入队一个社交分发发布任务。
//
// 与重索引类任务的 best-effort 不同：发布是用户显式动作，调度器未注入时
// 返回错误（否则 SCHEDULED 记录会被无声搁置）。
//
// processAt 零值表示立即执行；否则精确到该时刻（asynq DELAY 语义）。
// 重试由 asynq 承担（MaxRetry 3），每次重试前 worker 会把状态回写 FAILED。
func (s *TaskService) EnqueueSocialDistributionPublish(payload *task.SocialDistributionPublishPayload, processAt time.Time) error {
	if payload == nil || payload.DistributionID == 0 {
		return errors.New("nil social distribution publish payload")
	}
	if s.taskScheduler == nil {
		s.log.Errorf("enqueue social publish rejected: task scheduler not available (distribution=%d)",
			payload.DistributionID)
		return errors.New("task scheduler not available")
	}

	opts := []asynq.Option{asynq.MaxRetry(3)}
	if !processAt.IsZero() {
		opts = append(opts, asynq.ProcessAt(processAt))
	}

	if err := s.taskScheduler.NewTask(task.SocialDistributionPublishTaskType, payload, opts...); err != nil {
		s.log.Errorf("enqueue social publish failed (distribution=%d): %v", payload.DistributionID, err)
		return err
	}
	s.log.Infof("enqueued social publish (distribution=%d processAt=%s)",
		payload.DistributionID, processAt)
	return nil
}

// EnqueuePostPublish 入队一个文章定时发布任务（到点把 SCHEDULED 流转为 PUBLISHED）。
//
// 与 social 分发同语义：发布是用户显式动作，调度器未注入时返回错误
// （否则 SCHEDULED 文章会被无声搁置，前端承诺的"到点自动发布"落空）。
// processAt 零值表示立即执行。
func (s *TaskService) EnqueuePostPublish(payload *task.PostPublishPayload, processAt time.Time) error {
	if payload == nil || payload.PostID == 0 {
		return errors.New("nil post publish payload")
	}
	if s.taskScheduler == nil {
		s.log.Errorf("enqueue post publish rejected: task scheduler not available (post=%d)",
			payload.PostID)
		return errors.New("task scheduler not available")
	}

	opts := []asynq.Option{asynq.MaxRetry(3)}
	if !processAt.IsZero() {
		opts = append(opts, asynq.ProcessAt(processAt))
	}

	if err := s.taskScheduler.NewTask(task.PostPublishTaskType, payload, opts...); err != nil {
		s.log.Errorf("enqueue post publish failed (post=%d): %v", payload.PostID, err)
		return err
	}
	s.log.Infof("enqueued post publish (post=%d processAt=%s)", payload.PostID, processAt)
	return nil
}
