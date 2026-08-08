package service

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/go-utils/aggregator"
	"github.com/tx7do/go-utils/timeutil"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"github.com/tx7do/go-crud/viewer"
	"google.golang.org/protobuf/types/known/emptypb"

	"go-wind-cms/app/core/service/internal/data"

	identityV1 "go-wind-cms/api/gen/go/identity/service/v1"
	permissionV1 "go-wind-cms/api/gen/go/permission/service/v1"

	"go-wind-cms/pkg/constants"
	appViewer "go-wind-cms/pkg/entgo/viewer"
	"go-wind-cms/pkg/utils"
)

type RoleService struct {
	permissionV1.UnimplementedRoleServiceServer

	log *log.Helper

	roleRepo     *data.RoleRepo
	tenantRepo   *data.TenantRepo
	userRoleRepo *data.UserRoleRepo
	userRepo     data.UserRepo
}

func NewRoleService(
	ctx *bootstrap.Context,
	roleRepo *data.RoleRepo,
	tenantRepo *data.TenantRepo,
	userRoleRepo *data.UserRoleRepo,
	userRepo data.UserRepo,
) *RoleService {
	svc := &RoleService{
		log:          ctx.NewLoggerHelper("role/service/core-service"),
		roleRepo:     roleRepo,
		tenantRepo:   tenantRepo,
		userRoleRepo: userRoleRepo,
		userRepo:     userRepo,
	}

	svc.init()

	return svc
}

func (s *RoleService) init() {
	ctx := appViewer.NewSystemViewerContext(context.Background())
	if count, _ := s.roleRepo.Count(ctx, nil); count == 0 {
		_ = s.createDefaultRoles(ctx)
	}
}

func (s *RoleService) extractRelationIDs(
	roles []*permissionV1.Role,
	tenantSet aggregator.ResourceMap[uint32, *identityV1.Tenant],
) {
	for _, p := range roles {
		if p.GetTenantId() > 0 {
			tenantSet[p.GetTenantId()] = nil
		}
	}
}

func (s *RoleService) fetchRelationInfo(
	ctx context.Context,
	tenantSet aggregator.ResourceMap[uint32, *identityV1.Tenant],
) error {
	if len(tenantSet) > 0 {
		tenantIds := make([]uint32, 0, len(tenantSet))
		for id := range tenantSet {
			tenantIds = append(tenantIds, id)
		}

		tenants, err := s.tenantRepo.ListTenantsByIds(ctx, tenantIds)
		if err != nil {
			s.log.Errorf("query tenants err: %v", err)
			return err
		}

		for _, tenant := range tenants {
			tenantSet[tenant.GetId()] = tenant
		}
	}

	return nil
}

func (s *RoleService) bindRelations(
	roles []*permissionV1.Role,
	tenantSet aggregator.ResourceMap[uint32, *identityV1.Tenant],
) {
	aggregator.Populate(
		roles,
		tenantSet,
		func(ou *permissionV1.Role) uint32 { return ou.GetTenantId() },
		func(ou *permissionV1.Role, r *identityV1.Tenant) {
			ou.TenantName = r.Name
		},
	)
}

func (s *RoleService) enrichRelations(ctx context.Context, roles []*permissionV1.Role) error {
	var tenantSet = make(aggregator.ResourceMap[uint32, *identityV1.Tenant])
	s.extractRelationIDs(roles, tenantSet)
	if err := s.fetchRelationInfo(ctx, tenantSet); err != nil {
		return err
	}
	s.bindRelations(roles, tenantSet)
	return nil
}

func (s *RoleService) List(ctx context.Context, req *paginationV1.PagingRequest) (*permissionV1.ListRoleResponse, error) {
	resp, err := s.roleRepo.List(ctx, req)
	if err != nil {
		return nil, err
	}

	_ = s.enrichRelations(ctx, resp.Items)

	return resp, nil
}

func (s *RoleService) Count(ctx context.Context, req *paginationV1.PagingRequest) (*permissionV1.CountRoleResponse, error) {
	count, err := s.roleRepo.Count(ctx, req)
	if err != nil {
		return nil, err
	}

	return &permissionV1.CountRoleResponse{
		Count: uint64(count),
	}, nil
}

func (s *RoleService) Get(ctx context.Context, req *permissionV1.GetRoleRequest) (*permissionV1.Role, error) {
	resp, err := s.roleRepo.Get(ctx, req)
	if err != nil {
		return nil, err
	}

	fakeItems := []*permissionV1.Role{resp}
	_ = s.enrichRelations(ctx, fakeItems)

	return resp, nil
}

func (s *RoleService) Create(ctx context.Context, req *permissionV1.CreateRoleRequest) (*emptypb.Empty, error) {
	if req == nil || req.Data == nil {
		return nil, identityV1.ErrorBadRequest("invalid parameter")
	}

	if err := s.roleRepo.Create(ctx, req); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *RoleService) Update(ctx context.Context, req *permissionV1.UpdateRoleRequest) (*emptypb.Empty, error) {
	if req == nil || req.Data == nil {
		return nil, identityV1.ErrorBadRequest("invalid parameter")
	}

	r, err := s.roleRepo.Get(ctx, &permissionV1.GetRoleRequest{
		QueryBy: &permissionV1.GetRoleRequest_Id{
			Id: req.Data.GetId(),
		},
	})
	if err != nil {
		return nil, err
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

	if err = s.roleRepo.Update(ctx, req); err != nil {
		s.log.Errorf("update role error: %v", err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *RoleService) Delete(ctx context.Context, req *permissionV1.DeleteRoleRequest) (*emptypb.Empty, error) {
	var err error

	if err = s.roleRepo.Delete(ctx, req); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *RoleService) ListRoleCodesByIds(ctx context.Context, req *permissionV1.ListRoleCodesByIdsRequest) (*permissionV1.ListRoleCodesByIdsResponse, error) {
	ids, err := s.roleRepo.ListRoleCodesByRoleIds(ctx, req.GetRoleIds())
	if err != nil {
		return nil, err
	}

	return &permissionV1.ListRoleCodesByIdsResponse{
		RoleCodes: ids,
	}, nil
}

func (s *RoleService) ListRoleIdsByCodes(ctx context.Context, req *permissionV1.ListRoleIdsByCodesRequest) (*permissionV1.ListRoleIdsByCodesResponse, error) {
	ids, err := s.roleRepo.ListRoleIDsByRoleCodes(ctx, req.GetRoleCodes())
	if err != nil {
		return nil, err
	}

	return &permissionV1.ListRoleIdsByCodesResponse{
		RoleIds: ids,
	}, nil
}

func (s *RoleService) ListPermissionIds(ctx context.Context, req *permissionV1.ListPermissionIdsRequest) (*permissionV1.ListPermissionIdsResponse, error) {
	if req == nil {
		return nil, permissionV1.ErrorBadRequest("invalid request")
	}
	var permissionIDs []uint32
	var err error

	switch req.QueryBy.(type) {
	case *permissionV1.ListPermissionIdsRequest_RoleId:
		permissionIDs, err = s.roleRepo.ListPermissionIDsByRoleIDs(ctx, []uint32{req.GetRoleId()})
		if err != nil {
			return nil, err
		}

	case *permissionV1.ListPermissionIdsRequest_RoleCode:
		permissionIDs, err = s.roleRepo.ListPermissionIDsByRoleCodes(ctx, []string{req.GetRoleCode()})
		if err != nil {
			return nil, err
		}

	case *permissionV1.ListPermissionIdsRequest_UserId:
		// 校验目标用户归属当前调用者租户，避免跨租户枚举他人权限 ID
		if err := s.validateTargetUserTenant(ctx, req.GetUserId()); err != nil {
			return nil, err
		}
		permissionIDs, err = s.roleRepo.ListPermissionIDsByUserID(ctx, req.GetUserId())
		if err != nil {
			return nil, err
		}

	default:
		if len(req.RoleIds) > 0 {
			permissionIDs, err = s.roleRepo.ListPermissionIDsByRoleIDs(ctx, req.GetRoleIds())
			if err != nil {
				return nil, err
			}
		}

		if len(req.RoleCodes) > 0 {
			permissionIDs, err = s.roleRepo.ListPermissionIDsByRoleCodes(ctx, req.GetRoleCodes())
			if err != nil {
				return nil, err
			}
		}
	}

	return &permissionV1.ListPermissionIdsResponse{
		PermissionIds: permissionIDs,
	}, nil
}

func (s *RoleService) ListUserRoleIDs(ctx context.Context, req *permissionV1.ListUserRoleIDsRequest) (*permissionV1.ListUserRoleIDsResponse, error) {
	// 校验目标用户归属当前调用者租户，避免跨租户枚举他人角色 ID
	if err := s.validateTargetUserTenant(ctx, req.GetUserId()); err != nil {
		return nil, err
	}
	roleIDs, err := s.userRoleRepo.ListRoleIDs(ctx, req.GetUserId(), false)
	if err != nil {
		return nil, err
	}

	return &permissionV1.ListUserRoleIDsResponse{
		RoleIds: roleIDs,
	}, nil
}

// AssignRolesToUser 为用户分配角色（替换旧关联）
func (s *RoleService) AssignRolesToUser(ctx context.Context, req *permissionV1.AssignRolesToUserRequest) (*emptypb.Empty, error) {
	if req == nil || req.GetUserId() == 0 || len(req.GetRoleIds()) == 0 {
		return nil, permissionV1.ErrorBadRequest("invalid parameter")
	}

	// 操作人身份与租户必须从 viewer context 推导，忽略客户端传入的 operator_id/tenant_id，
	// 防止越权将角色绑定写入他租户或伪造审计归属
	callerUserID, hasUser := viewerUserIDFromContext(ctx)
	if !hasUser {
		return nil, permissionV1.ErrorBadRequest("operator identity is required")
	}
	var callerTenantID uint32
	isPlatform := false
	if vc, exist := viewer.FromContext(ctx); exist && vc != nil {
		callerTenantID = uint32(vc.TenantID())
		isPlatform = vc.IsPlatformContext() || vc.IsSystemContext()
	}

	// 校验目标用户归属当前调用者租户：平台上下文放行；租户用户仅可向本租户用户授权。
	// (userRepo.Get 经 EvalQuery 注入 tenant 谓词，租户用户取不到他租户用户即自动拒绝)
	targetUser, err := s.userRepo.Get(ctx, &identityV1.GetUserRequest{
		QueryBy: &identityV1.GetUserRequest_Id{Id: req.GetUserId()},
	})
	if err != nil || targetUser == nil {
		return nil, permissionV1.ErrorBadRequest("target user not found")
	}
	if !isPlatform && targetUser.GetTenantId() != callerTenantID {
		return nil, permissionV1.ErrorForbidden("cannot assign roles to a user in another tenant")
	}

	// 逐个校验角色可分配性（启用/非模板/非阻断/非平台管理员模板），CanAssignRole 内部 Get 经
	// EvalQuery 注入 tenant 谓词，租户用户取不到他租户角色即自动拒绝
	for _, roleID := range req.GetRoleIds() {
		ok, err := s.roleRepo.CanAssignRole(ctx, roleID)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, permissionV1.ErrorForbidden("role %d is not assignable", roleID)
		}
	}

	now := time.Now()

	userRoles := make([]*permissionV1.UserRole, 0, len(req.GetRoleIds()))
	for _, roleID := range req.GetRoleIds() {
		userRoles = append(userRoles, &permissionV1.UserRole{
			TenantId:   trans.Ptr(callerTenantID),
			UserId:     trans.Ptr(req.GetUserId()),
			RoleId:     trans.Ptr(roleID),
			Status:     permissionV1.UserRole_ACTIVE.Enum(),
			AssignedBy: trans.Ptr(callerUserID),
			AssignedAt: timeutil.TimeToTimestamppb(&now),
			StartAt:    timeutil.TimeToTimestamppb(&now),
			CreatedBy:  trans.Ptr(callerUserID),
		})
	}

	if err := s.userRoleRepo.AssignRolesToUser(ctx, req.GetUserId(), userRoles); err != nil {
		s.log.Errorf("assign roles to user failed: %v", err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// GetUserRoles 获取用户的角色绑定列表
func (s *RoleService) GetUserRoles(ctx context.Context, req *permissionV1.GetUserRolesRequest) (*permissionV1.GetUserRolesResponse, error) {
	if req == nil || req.GetUserId() == 0 {
		return nil, permissionV1.ErrorBadRequest("invalid parameter")
	}

	// 校验目标用户归属当前调用者租户，避免跨租户枚举他人角色绑定
	if err := s.validateTargetUserTenant(ctx, req.GetUserId()); err != nil {
		return nil, err
	}

	bindings, err := s.userRoleRepo.ListByUserID(ctx, req.GetUserId(), req.GetIncludeExpired())
	if err != nil {
		return nil, err
	}

	return &permissionV1.GetUserRolesResponse{
		Bindings: bindings,
	}, nil
}

// UnassignRolesFromUser 从用户移除指定角色
func (s *RoleService) UnassignRolesFromUser(ctx context.Context, req *permissionV1.UnassignRolesFromUserRequest) (*permissionV1.UnassignRolesFromUserResponse, error) {
	if req == nil || req.GetUserId() == 0 || len(req.GetRoleIds()) == 0 {
		return nil, permissionV1.ErrorBadRequest("invalid parameter")
	}

	if err := s.userRoleRepo.RemoveRolesFromUser(ctx, req.GetUserId(), req.GetRoleIds()); err != nil {
		s.log.Errorf("unassign roles from user failed: %v", err)
		return nil, err
	}

	return &permissionV1.UnassignRolesFromUserResponse{
		RemovedRoleIds: req.GetRoleIds(),
	}, nil
}

// validateTargetUserTenant 校验目标用户归属当前调用者租户：
// 平台/系统上下文放行；租户用户仅可查询本租户用户的数据。
// userRepo.Get 经 EvalQuery 注入 tenant 谓词，租户用户取不到他租户用户即自动拒绝。
// 返回 nil 表示放行，非 nil 表示拒绝（调用方应直接返回该错误）。
func (s *RoleService) validateTargetUserTenant(ctx context.Context, targetUserID uint32) error {
	var callerTenantID uint32
	isPlatform := false
	if vc, exist := viewer.FromContext(ctx); exist && vc != nil {
		callerTenantID = uint32(vc.TenantID())
		isPlatform = vc.IsPlatformContext() || vc.IsSystemContext()
	}
	targetUser, err := s.userRepo.Get(ctx, &identityV1.GetUserRequest{
		QueryBy: &identityV1.GetUserRequest_Id{Id: targetUserID},
	})
	if err != nil || targetUser == nil {
		return permissionV1.ErrorBadRequest("target user not found")
	}
	if !isPlatform && targetUser.GetTenantId() != callerTenantID {
		return permissionV1.ErrorForbidden("cannot access a user in another tenant")
	}
	return nil
}

// createDefaultRoles 创建默认角色(包括超级管理员)
func (s *RoleService) createDefaultRoles(ctx context.Context) error {
	var err error

	for _, d := range constants.DefaultRoles {
		err = s.roleRepo.Create(ctx, &permissionV1.CreateRoleRequest{
			Data: d,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
