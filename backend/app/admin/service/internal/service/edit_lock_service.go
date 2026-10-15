package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
)

// EditLockService 编辑锁门面（透传 core）。
type EditLockService struct {
	adminV1.EditLockServiceHTTPServer

	editLockServiceClient contentV1.EditLockServiceClient
	log                   *log.Helper
}

func NewEditLockService(ctx *bootstrap.Context, editLockServiceClient contentV1.EditLockServiceClient) *EditLockService {
	return &EditLockService{
		log:                   ctx.NewLoggerHelper("edit-lock/service/admin-service"),
		editLockServiceClient: editLockServiceClient,
	}
}

func (s *EditLockService) Acquire(ctx context.Context, req *contentV1.AcquireEditLockRequest) (*contentV1.EditLock, error) {
	return s.editLockServiceClient.Acquire(ctx, req)
}

func (s *EditLockService) Release(ctx context.Context, req *contentV1.ReleaseEditLockRequest) (*emptypb.Empty, error) {
	return s.editLockServiceClient.Release(ctx, req)
}
