package service

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/go-kratos/kratos/v2/log"

	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/go-utils/sliceutil"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	authenticationV1 "go-wind-cms/api/gen/go/authentication/service/v1"
	identityV1 "go-wind-cms/api/gen/go/identity/service/v1"
	permissionV1 "go-wind-cms/api/gen/go/permission/service/v1"

	"go-wind-cms/pkg/middleware/auth"
	"go-wind-cms/pkg/utils"
)

type UserService struct {
	adminV1.UserServiceHTTPServer

	log *log.Helper

	userServiceClient     identityV1.UserServiceClient
	tenantServiceClient   identityV1.TenantServiceClient
	orgUnitServiceClient  identityV1.OrgUnitServiceClient
	positionServiceClient identityV1.PositionServiceClient

	userCredentialServiceClient authenticationV1.UserCredentialServiceClient
	roleServiceClient           permissionV1.RoleServiceClient
}

func NewUserService(
	ctx *bootstrap.Context,
	userServiceClient identityV1.UserServiceClient,
	tenantServiceClient identityV1.TenantServiceClient,
	orgUnitServiceClient identityV1.OrgUnitServiceClient,
	positionServiceClient identityV1.PositionServiceClient,
	roleServiceClient permissionV1.RoleServiceClient,
	userCredentialServiceClient authenticationV1.UserCredentialServiceClient,
) *UserService {
	svc := &UserService{
		log:                         ctx.NewLoggerHelper("user/service/admin-service"),
		userServiceClient:           userServiceClient,
		tenantServiceClient:         tenantServiceClient,
		orgUnitServiceClient:        orgUnitServiceClient,
		positionServiceClient:       positionServiceClient,
		roleServiceClient:           roleServiceClient,
		userCredentialServiceClient: userCredentialServiceClient,
	}
	svc.init()
	return svc
}

func (s *UserService) init() {
}

func (s *UserService) List(ctx context.Context, req *paginationV1.PagingRequest) (*identityV1.ListUserResponse, error) {
	return s.userServiceClient.List(ctx, req)
}

func (s *UserService) Count(ctx context.Context, req *paginationV1.PagingRequest) (*identityV1.CountUserResponse, error) {
	return s.userServiceClient.Count(ctx, req)
}

func (s *UserService) Get(ctx context.Context, req *identityV1.GetUserRequest) (*identityV1.User, error) {
	resp, err := s.userServiceClient.Get(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *UserService) Create(ctx context.Context, req *identityV1.CreateUserRequest) (*identityV1.User, error) {
	if req == nil || req.Data == nil {
		return nil, adminV1.ErrorBadRequest("invalid request")
	}

	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.Data.CreatedBy = trans.Ptr(operator.GetUserId())
	if operator.GetTenantId() > 0 {
		req.Data.TenantId = operator.TenantId
	}

	// 获取操作者的用户信息
	_, err = s.userServiceClient.Get(ctx, &identityV1.GetUserRequest{
		QueryBy: &identityV1.GetUserRequest_Id{
			Id: req.Data.GetCreatedBy(),
		},
	})
	if err != nil {
		return nil, err
	}

	var roleIds []uint32
	if len(req.Data.GetRoleIds()) > 0 {
		roleIds = req.Data.GetRoleIds()
	}
	if req.Data.RoleId != nil && *req.Data.RoleId > 0 {
		roleIds = append(roleIds, *req.Data.RoleId)
	}
	roleIds = sliceutil.Unique(roleIds)
	if len(roleIds) == 0 {
		s.log.Errorf("role_ids is required")
		return nil, identityV1.ErrorBadRequest("role_ids is required")
	}

	var queryString string
	if operator.GetTenantId() > 0 || req.Data.GetTenantId() > 0 {
		queryString = fmt.Sprintf(`{"id__in": "[%s]", "type": "TENANT", "tenant_id": %d}`,
			utils.NumberSliceToString(roleIds),
			req.Data.GetTenantId(),
		)
	} else {
		queryString = fmt.Sprintf(`{"id__in": "[%s]", "type": "SYSTEM"}`,
			utils.NumberSliceToString(roleIds),
		)
	}
	roles, err := s.roleServiceClient.List(ctx, &paginationV1.PagingRequest{
		NoPaging: trans.Ptr(true),
		FilteringType: &paginationV1.PagingRequest_Query{
			Query: queryString,
		},
	})
	if err != nil {
		s.log.Errorf("query roles err: %v", err)
		return nil, err
	}

	if len(roles.Items) != len(roleIds) {
		s.log.Errorf("some roles not found, requested role ids: %v", roleIds)
		return nil, identityV1.ErrorBadRequest("some roles not found")
	}
	if len(roles.Items) == 0 {
		s.log.Errorf("at least one role is required")
		return nil, identityV1.ErrorBadRequest("at least one role is required")
	}

	req.Data.RoleId = nil
	req.Data.RoleIds = roleIds

	return s.userServiceClient.Create(ctx, req)
}

func (s *UserService) Update(ctx context.Context, req *identityV1.UpdateUserRequest) (*emptypb.Empty, error) {
	if req == nil || req.Data == nil {
		return nil, adminV1.ErrorBadRequest("invalid request")
	}

	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	// 获取操作者的用户信息
	_, err = s.userServiceClient.Get(ctx, &identityV1.GetUserRequest{
		QueryBy: &identityV1.GetUserRequest_Id{
			Id: operator.GetUserId(),
		},
	})
	if err != nil {
		return nil, err
	}

	req.Data.UpdatedBy = trans.Ptr(operator.GetUserId())
	if req.UpdateMask != nil {
		req.UpdateMask.Paths = append(req.UpdateMask.Paths, "updated_by")
	}

	if operator.GetTenantId() > 0 {
		req.Data.TenantId = operator.TenantId
	}

	req.Data.Id = trans.Ptr(req.GetId())

	if req.GetPassword() != "" {
		if _, err = s.userCredentialServiceClient.ResetCredential(ctx, &authenticationV1.ResetCredentialRequest{
			IdentityType:  authenticationV1.UserCredential_USERNAME,
			Identifier:    req.Data.GetUsername(),
			NewCredential: req.GetPassword(),
			NeedDecrypt:   false,
		}); err != nil {
			s.log.Errorf("reset user password err: %v", err)
			return nil, err
		}
	}

	var roleIds []uint32
	if len(req.Data.GetRoleIds()) > 0 {
		roleIds = req.Data.GetRoleIds()
	}
	if req.Data.RoleId != nil && *req.Data.RoleId > 0 {
		roleIds = append(roleIds, *req.Data.RoleId)
	}
	roleIds = sliceutil.Unique(roleIds)
	if len(roleIds) == 0 {
		s.log.Errorf("role_ids is required")
		return nil, adminV1.ErrorBadRequest("role_ids is required")
	}

	var queryString string
	if operator.GetTenantId() > 0 || req.Data.GetTenantId() > 0 {
		queryString = fmt.Sprintf(`{"id__in": "[%s]", "type": "TENANT", "tenant_id": %d}`,
			utils.NumberSliceToString(roleIds),
			req.Data.GetTenantId(),
		)
	} else {
		queryString = fmt.Sprintf(`{"id__in": "[%s]", "type": "SYSTEM"}`,
			utils.NumberSliceToString(roleIds),
		)
	}
	roles, err := s.roleServiceClient.List(ctx, &paginationV1.PagingRequest{
		NoPaging: trans.Ptr(true),
		FilteringType: &paginationV1.PagingRequest_Query{
			Query: queryString,
		},
	})
	if err != nil {
		s.log.Errorf("query roles err: %v", err)
		return nil, err
	}

	if len(roles.Items) != len(roleIds) {
		s.log.Errorf("some roles not found, requested role ids: %v", roleIds)
		return nil, adminV1.ErrorBadRequest("some roles not found")
	}
	if len(roles.Items) == 0 {
		s.log.Errorf("at least one role is required")
		return nil, adminV1.ErrorBadRequest("at least one role is required")
	}

	if len(roleIds) > 0 {
		if _, err = s.roleServiceClient.AssignRolesToUser(ctx, &permissionV1.AssignRolesToUserRequest{
			UserId:     req.Data.GetId(),
			TenantId:   req.Data.GetTenantId(),
			RoleIds:    roleIds,
			OperatorId: req.Data.GetUpdatedBy(),
			Reason:     "更新用户信息，分配角色",
		}); err != nil {
			s.log.Errorf("assign roles to user err: %v", err)
			return nil, err
		}
	}

	return s.userServiceClient.Update(ctx, req)
}

func (s *UserService) Delete(ctx context.Context, req *identityV1.DeleteUserRequest) (*emptypb.Empty, error) {
	if req == nil {
		return nil, adminV1.ErrorBadRequest("invalid request")
	}

	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.DeletedBy = trans.Ptr(operator.UserId)

	// 获取操作者的用户信息
	_, err = s.userServiceClient.Get(ctx, &identityV1.GetUserRequest{
		QueryBy: &identityV1.GetUserRequest_Id{
			Id: operator.UserId,
		},
	})
	if err != nil {
		return nil, err
	}

	getRequest := &identityV1.GetUserRequest{
		QueryBy: &identityV1.GetUserRequest_Id{
			Id: operator.UserId,
		},
	}
	switch req.DeleteBy.(type) {
	case *identityV1.DeleteUserRequest_Username:
		getRequest.QueryBy = &identityV1.GetUserRequest_Id{
			Id: req.GetId(),
		}
	case *identityV1.DeleteUserRequest_Id:
		getRequest.QueryBy = &identityV1.GetUserRequest_Username{
			Username: req.GetUsername(),
		}
	default:
		return nil, adminV1.ErrorBadRequest("invalid request delete_by")
	}

	var deleteUser *identityV1.User
	deleteUser, err = s.userServiceClient.Get(ctx, getRequest)
	if err != nil {
		return nil, err
	}

	userRoles, err := s.roleServiceClient.GetUserRoles(ctx, &permissionV1.GetUserRolesRequest{
		UserId:         deleteUser.GetId(),
		TenantId:       deleteUser.GetTenantId(),
		IncludeExpired: true,
	})
	if err != nil {
		return nil, err
	}
	var roleIds []uint32
	for _, role := range userRoles.GetBindings() {
		roleIds = append(roleIds, role.GetRoleId())
	}

	if _, err = s.roleServiceClient.UnassignRolesFromUser(ctx, &permissionV1.UnassignRolesFromUserRequest{
		UserId:     deleteUser.GetId(),
		TenantId:   deleteUser.GetTenantId(),
		RoleIds:    roleIds,
		OperatorId: req.GetDeletedBy(),
		Reason:     trans.Ptr("删除用户，撤销角色分配"),
	}); err != nil {
		s.log.Errorf("unassign roles from user err: %v", err)
		return nil, err
	}

	return s.userServiceClient.Delete(ctx, req)
}

func (s *UserService) UserExists(ctx context.Context, req *identityV1.UserExistsRequest) (*identityV1.UserExistsResponse, error) {
	return s.userServiceClient.UserExists(ctx, req)
}

// EditUserPassword 修改用户密码
func (s *UserService) EditUserPassword(ctx context.Context, req *identityV1.EditUserPasswordRequest) (*emptypb.Empty, error) {
	// 获取操作者的用户信息
	u, err := s.userServiceClient.Get(ctx, &identityV1.GetUserRequest{
		QueryBy: &identityV1.GetUserRequest_Id{
			Id: req.GetUserId(),
		},
	})
	if err != nil {
		return nil, err
	}

	if _, err = s.userCredentialServiceClient.ResetCredential(ctx, &authenticationV1.ResetCredentialRequest{
		IdentityType:  authenticationV1.UserCredential_USERNAME,
		Identifier:    u.GetUsername(),
		NewCredential: req.GetNewPassword(),
		NeedDecrypt:   false,
	}); err != nil {
		s.log.Errorf("reset user password err: %v", err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
