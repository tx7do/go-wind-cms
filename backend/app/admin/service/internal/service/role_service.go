package service

import (
	"context"
	"go-wind-cms/pkg/utils"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	identityV1 "go-wind-cms/api/gen/go/identity/service/v1"
	permissionV1 "go-wind-cms/api/gen/go/permission/service/v1"

	"go-wind-cms/pkg/middleware/auth"
)

type RoleService struct {
	adminV1.RoleServiceHTTPServer

	log *log.Helper

	roleServiceClient   permissionV1.RoleServiceClient
	tenantServiceClient identityV1.TenantServiceClient
}

func NewRoleService(
	ctx *bootstrap.Context,
	roleServiceClient permissionV1.RoleServiceClient,
	tenantServiceClient identityV1.TenantServiceClient,
) *RoleService {
	l := log.NewHelper(log.With(ctx.GetLogger(), "module", "role/service/admin-service"))
	return &RoleService{
		log:                 l,
		roleServiceClient:   roleServiceClient,
		tenantServiceClient: tenantServiceClient,
	}
}

func (s *RoleService) List(ctx context.Context, req *paginationV1.PagingRequest) (*permissionV1.ListRoleResponse, error) {
	resp, err := s.roleServiceClient.List(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *RoleService) Count(ctx context.Context, req *paginationV1.PagingRequest) (*permissionV1.CountRoleResponse, error) {
	return s.roleServiceClient.Count(ctx, req)
}

func (s *RoleService) Get(ctx context.Context, req *permissionV1.GetRoleRequest) (*permissionV1.Role, error) {
	resp, err := s.roleServiceClient.Get(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *RoleService) Create(ctx context.Context, req *permissionV1.CreateRoleRequest) (*emptypb.Empty, error) {
	if req == nil || req.Data == nil {
		return nil, adminV1.ErrorBadRequest("invalid request")
	}

	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.Data.CreatedBy = trans.Ptr(operator.GetUserId())

	if operator.GetTenantId() > 0 && req.Data.GetType() != permissionV1.Role_TENANT {
		req.Data.Type = trans.Ptr(permissionV1.Role_TENANT)
	}

	_, err = s.roleServiceClient.Create(ctx, req)

	return &emptypb.Empty{}, err
}

func (s *RoleService) Update(ctx context.Context, req *permissionV1.UpdateRoleRequest) (*emptypb.Empty, error) {
	if req == nil || req.Data == nil {
		return nil, adminV1.ErrorBadRequest("invalid request")
	}

	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.Data.UpdatedBy = trans.Ptr(operator.GetUserId())
	if req.UpdateMask != nil {
		req.UpdateMask.Paths = append(req.UpdateMask.Paths, "updated_by")
	}

	if operator.GetTenantId() > 0 && req.Data.GetType() != permissionV1.Role_TENANT {
		req.Data.Type = trans.Ptr(permissionV1.Role_TENANT)
	}

	r, err := s.roleServiceClient.Get(ctx, &permissionV1.GetRoleRequest{
		QueryBy: &permissionV1.GetRoleRequest_Id{
			Id: req.Data.GetId(),
		},
	})
	if err != nil {
		return nil, err
	}

	// 非系统管理员禁止修改系统角色
	if r.GetType() == permissionV1.Role_SYSTEM && operator.GetTenantId() > 0 {
		return nil, adminV1.ErrorForbidden("no permission to update system role")
	}

	// 保护角色字段不可修改
	if r.GetIsProtected() {
		if len(req.GetUpdateMask().Paths) > 0 {
			req.GetUpdateMask().Paths = utils.FilterBlacklist(req.GetUpdateMask().Paths, []string{
				"is_protected",
				"type",
				"status",
				"code",
			})
		} else {
			req.Data.IsProtected = nil
			req.Data.Type = nil
			req.Data.Status = nil
			req.Data.Code = nil
		}
	}

	return s.roleServiceClient.Update(ctx, req)
}

func (s *RoleService) Delete(ctx context.Context, req *permissionV1.DeleteRoleRequest) (*emptypb.Empty, error) {
	if req == nil {
		return nil, adminV1.ErrorBadRequest("invalid request")
	}

	return s.roleServiceClient.Delete(ctx, req)
}
