package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/go-utils/trans"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	taskV1 "go-wind-cms/api/gen/go/task/service/v1"

	"go-wind-cms/pkg/middleware/auth"
)

type TaskService struct {
	adminV1.TaskServiceHTTPServer

	log *log.Helper

	taskServiceClient taskV1.TaskServiceClient
}

func NewTaskService(
	ctx *bootstrap.Context,
	taskServiceClient taskV1.TaskServiceClient,
) *TaskService {
	l := log.NewHelper(log.With(ctx.GetLogger(), "module", "task/service/admin-service"))
	return &TaskService{
		log:               l,
		taskServiceClient: taskServiceClient,
	}
}

func (s *TaskService) List(ctx context.Context, req *paginationV1.PagingRequest) (*taskV1.ListTaskResponse, error) {
	return s.taskServiceClient.List(ctx, req)
}

func (s *TaskService) Get(ctx context.Context, req *taskV1.GetTaskRequest) (*taskV1.Task, error) {
	return s.taskServiceClient.Get(ctx, req)
}

func (s *TaskService) ListTaskTypeName(ctx context.Context, req *emptypb.Empty) (*taskV1.ListTaskTypeNameResponse, error) {
	return s.taskServiceClient.ListTaskTypeName(ctx, req)
}

func (s *TaskService) Create(ctx context.Context, req *taskV1.CreateTaskRequest) (*emptypb.Empty, error) {
	if req.Data == nil {
		return nil, adminV1.ErrorBadRequest("invalid parameter")
	}

	// 获取操作人信息
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.Data.CreatedBy = trans.Ptr(operator.UserId)

	_, err = s.taskServiceClient.Create(ctx, req)
	return &emptypb.Empty{}, err
}

func (s *TaskService) Update(ctx context.Context, req *taskV1.UpdateTaskRequest) (*emptypb.Empty, error) {
	if req.Data == nil {
		return nil, adminV1.ErrorBadRequest("invalid parameter")
	}

	// 获取操作人信息
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.Data.UpdatedBy = trans.Ptr(operator.GetUserId())
	if req.UpdateMask != nil {
		req.UpdateMask.Paths = append(req.UpdateMask.Paths, "updated_by")
	}

	_, err = s.taskServiceClient.Update(ctx, req)
	return &emptypb.Empty{}, err
}

func (s *TaskService) Delete(ctx context.Context, req *taskV1.DeleteTaskRequest) (*emptypb.Empty, error) {
	return s.taskServiceClient.Delete(ctx, req)
}

// ControlTask 控制调度任务
func (s *TaskService) ControlTask(ctx context.Context, req *taskV1.ControlTaskRequest) (*emptypb.Empty, error) {
	return s.taskServiceClient.ControlTask(ctx, req)
}

// StopAllTask 停止所有的调度任务
func (s *TaskService) StopAllTask(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	return s.taskServiceClient.StopAllTask(ctx, req)
}

// StartAllTask 启动所有的调度任务
func (s *TaskService) StartAllTask(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	return s.taskServiceClient.StartAllTask(ctx, req)
}

// RestartAllTask 重启所有的调度任务
func (s *TaskService) RestartAllTask(ctx context.Context, req *emptypb.Empty) (*taskV1.RestartAllTaskResponse, error) {
	return s.taskServiceClient.RestartAllTask(ctx, req)
}
