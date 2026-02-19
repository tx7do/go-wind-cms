package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	permissionV1 "go-wind-cms/api/gen/go/permission/service/v1"

	"go-wind-cms/pkg/middleware/auth"
)

type PermissionGroupService struct {
	adminV1.PermissionGroupServiceHTTPServer

	log *log.Helper

	permissionServiceClient      permissionV1.PermissionServiceClient
	permissionGroupServiceClient permissionV1.PermissionGroupServiceClient
}

func NewPermissionGroupService(
	ctx *bootstrap.Context,
	permissionServiceClient permissionV1.PermissionServiceClient,
	permissionGroupServiceClient permissionV1.PermissionGroupServiceClient,
) *PermissionGroupService {
	svc := &PermissionGroupService{
		log:                          ctx.NewLoggerHelper("permission-group/service/admin-service"),
		permissionServiceClient:      permissionServiceClient,
		permissionGroupServiceClient: permissionGroupServiceClient,
	}

	svc.init()

	return svc
}

func (s *PermissionGroupService) init() {

}

func (s *PermissionGroupService) List(ctx context.Context, req *paginationV1.PagingRequest) (*permissionV1.ListPermissionGroupResponse, error) {
	return s.permissionGroupServiceClient.List(ctx, req)
}

func (s *PermissionGroupService) Get(ctx context.Context, req *permissionV1.GetPermissionGroupRequest) (*permissionV1.PermissionGroup, error) {
	return s.permissionGroupServiceClient.Get(ctx, req)
}

func (s *PermissionGroupService) Create(ctx context.Context, req *permissionV1.CreatePermissionGroupRequest) (*emptypb.Empty, error) {
	if req.Data == nil {
		return nil, adminV1.ErrorBadRequest("invalid parameter")
	}

	// 获取操作人信息
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.Data.CreatedBy = trans.Ptr(operator.UserId)

	if _, err = s.permissionGroupServiceClient.Create(ctx, req); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *PermissionGroupService) Update(ctx context.Context, req *permissionV1.UpdatePermissionGroupRequest) (*emptypb.Empty, error) {
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

	return s.permissionGroupServiceClient.Update(ctx, req)
}

func (s *PermissionGroupService) Delete(ctx context.Context, req *permissionV1.DeletePermissionGroupRequest) (*emptypb.Empty, error) {
	var err error

	if _, err = s.permissionGroupServiceClient.Delete(ctx, req); err != nil {
		return nil, err
	}

	if _, err = s.permissionServiceClient.Delete(ctx, &permissionV1.DeletePermissionRequest{
		DeleteBy: &permissionV1.DeletePermissionRequest_GroupId{GroupId: req.GetId()},
	}); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
