package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/dtm-labs/dtm/client/dtmcli"
	"github.com/dtm-labs/dtm/client/dtmgrpc/dtmgimp"
	"github.com/dtm-labs/dtm/client/workflow"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/log"

	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/go-utils/aggregator"
	"github.com/tx7do/go-utils/id"
	"github.com/tx7do/go-utils/sliceutil"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	authenticationV1 "go-wind-cms/api/gen/go/authentication/service/v1"
	identityV1 "go-wind-cms/api/gen/go/identity/service/v1"
	permissionV1 "go-wind-cms/api/gen/go/permission/service/v1"

	"go-wind-cms/pkg/middleware/auth"
)

const (
	WorkflowTenantServiceCreateTenantWithAdminSAGA = "workflow_tenant_create_with_admin_saga"
)

type TenantService struct {
	adminV1.TenantServiceHTTPServer

	log *log.Helper

	codec encoding.Codec

	userServiceClient           identityV1.UserServiceClient
	userCredentialServiceClient authenticationV1.UserCredentialServiceClient
	tenantServiceClient         identityV1.TenantServiceClient
	roleServiceClient           permissionV1.RoleServiceClient
}

func NewTenantService(
	ctx *bootstrap.Context,
	userServiceClient identityV1.UserServiceClient,
	userCredentialServiceClient authenticationV1.UserCredentialServiceClient,
	tenantServiceClient identityV1.TenantServiceClient,
	roleServiceClient permissionV1.RoleServiceClient,
) *TenantService {
	svc := &TenantService{
		log:                         ctx.NewLoggerHelper("tenant/service/admin-service"),
		userServiceClient:           userServiceClient,
		userCredentialServiceClient: userCredentialServiceClient,
		tenantServiceClient:         tenantServiceClient,
		roleServiceClient:           roleServiceClient,
		codec:                       encoding.GetCodec("proto"),
	}

	svc.init()

	return svc
}

func (s *TenantService) init() {
	var err error

	// Register SAGA workflow
	if err = s.registerWorkflowCreateTenantWithAdminUser(); err != nil {
		s.log.Fatalf("注册工作流失败: %v", err)
		return
	}
}

func (s *TenantService) extractRelationIDs(
	tenants []*identityV1.Tenant,
	userSet aggregator.ResourceMap[uint32, *identityV1.User],
) {
	for _, t := range tenants {
		if t.GetAdminUserId() > 0 {
			userSet[t.GetAdminUserId()] = nil
		}
	}
}

func (s *TenantService) fetchRelationInfo(
	ctx context.Context,
	userSet aggregator.ResourceMap[uint32, *identityV1.User],
) error {
	if len(userSet) > 0 {
		userIds := make([]uint32, 0, len(userSet))
		for userId := range userSet {
			userIds = append(userIds, userId)
		}

		users, err := s.userServiceClient.List(ctx, &paginationV1.PagingRequest{
			NoPaging: trans.Ptr(true),
			FilteringType: &paginationV1.PagingRequest_Query{
				Query: fmt.Sprintf(`{"id__in": "[%s]"}`, strings.Join(
					sliceutil.Map(userIds, func(value uint32, _ int, _ []uint32) string {
						return strconv.FormatUint(uint64(value), 10)
					}),
					","),
				),
			},
		})
		if err != nil {
			s.log.Errorf("query users err: %v", err)
			return err
		}

		for _, u := range users.Items {
			userSet[u.GetId()] = u
		}
	}

	return nil
}

func (s *TenantService) bindRelations(
	tenants []*identityV1.Tenant,
	userSet aggregator.ResourceMap[uint32, *identityV1.User],
) {
	aggregator.Populate(
		tenants,
		userSet,
		func(ou *identityV1.Tenant) uint32 { return ou.GetAdminUserId() },
		func(ou *identityV1.Tenant, r *identityV1.User) {
			ou.AdminUserName = r.Username
		},
	)
}

func (s *TenantService) enrichRelations(ctx context.Context, tenants []*identityV1.Tenant) error {
	var userSet = make(aggregator.ResourceMap[uint32, *identityV1.User])
	s.extractRelationIDs(tenants, userSet)
	if err := s.fetchRelationInfo(ctx, userSet); err != nil {
		return err
	}
	s.bindRelations(tenants, userSet)
	return nil
}

func (s *TenantService) List(ctx context.Context, req *paginationV1.PagingRequest) (*identityV1.ListTenantResponse, error) {
	resp, err := s.tenantServiceClient.List(ctx, req)
	if err != nil {
		return nil, err
	}

	_ = s.enrichRelations(ctx, resp.Items)

	return resp, nil
}

func (s *TenantService) Count(ctx context.Context, req *paginationV1.PagingRequest) (*identityV1.CountTenantResponse, error) {
	return s.tenantServiceClient.Count(ctx, req)
}

func (s *TenantService) Get(ctx context.Context, req *identityV1.GetTenantRequest) (*identityV1.Tenant, error) {
	resp, err := s.tenantServiceClient.Get(ctx, req)
	if err != nil {
		return nil, err
	}

	fakeItems := []*identityV1.Tenant{resp}
	_ = s.enrichRelations(ctx, fakeItems)

	return resp, nil
}

func (s *TenantService) Create(ctx context.Context, req *identityV1.CreateTenantRequest) (*emptypb.Empty, error) {
	if req.Data == nil {
		return nil, adminV1.ErrorBadRequest("invalid parameter")
	}

	// 获取操作人信息
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.Data.CreatedBy = trans.Ptr(operator.UserId)

	if _, err = s.tenantServiceClient.Create(ctx, req); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *TenantService) Update(ctx context.Context, req *identityV1.UpdateTenantRequest) (*emptypb.Empty, error) {
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

	return s.tenantServiceClient.Update(ctx, req)
}

func (s *TenantService) Delete(ctx context.Context, req *identityV1.DeleteTenantRequest) (*emptypb.Empty, error) {
	return s.tenantServiceClient.Delete(ctx, req)
}

func (s *TenantService) TenantExists(ctx context.Context, req *identityV1.TenantExistsRequest) (*identityV1.TenantExistsResponse, error) {
	return s.tenantServiceClient.TenantExists(ctx, req)
}

// registerWorkflowCreateTenantWithAdminUser 注册创建租户及管理员用户的 SAGA 工作流
func (s *TenantService) registerWorkflowCreateTenantWithAdminUser() (err error) {
	err = workflow.Register(WorkflowTenantServiceCreateTenantWithAdminSAGA, func(wf *workflow.Workflow, data []byte) error {
		var req identityV1.CreateTenantWithAdminUserRequest
		if len(data) > 0 {
			if err = s.codec.Unmarshal(data, &req); err != nil {
				s.log.Errorf("工作流数据反序列化失败: %v", err)
				return identityV1.ErrorInternalServerError("工作流数据反序列化失败")
			}
		}

		// Create tenant
		wf.NewBranch().OnRollback(func(bb *dtmcli.BranchBarrier) error {
			if _, err = s.tenantServiceClient.Delete(wf.Context, &identityV1.DeleteTenantRequest{
				DeleteBy: &identityV1.DeleteTenantRequest_Code{Code: req.Tenant.GetCode()},
			}); err != nil {
				s.log.Errorf("工作流回滚，创建租户失败: %v", err)
				return identityV1.ErrorInternalServerError("工作流回滚，创建租户失败")
			}
			return nil
		})
		var tenant *identityV1.Tenant
		if tenant, err = s.tenantServiceClient.Create(wf.Context, &identityV1.CreateTenantRequest{Data: req.Tenant}); err != nil {
			s.log.Errorf("create tenant err: %v", err)
			return err
		}
		req.User.TenantId = tenant.Id

		// copy tenant manager role to tenant
		wf.NewBranch().OnRollback(func(bb *dtmcli.BranchBarrier) error {
			if _, err = s.roleServiceClient.Delete(wf.Context, &permissionV1.DeleteRoleRequest{
				DeleteBy: &permissionV1.DeleteRoleRequest_TenantId{TenantId: tenant.GetId()},
			}); err != nil {
				s.log.Errorf("工作流回滚，创建租户角色失败: %v", err)
				return identityV1.ErrorInternalServerError("工作流回滚，创建租户角色失败")
			}
			return nil
		})
		var role *permissionV1.Role
		if role, err = s.roleServiceClient.CreateTenantRoleFromTemplate(wf.Context, &permissionV1.CreateTenantRoleFromTemplateRequest{
			TenantId:   tenant.GetId(),
			OperatorId: req.User.GetCreatedBy(),
		}); err != nil {
			s.log.Errorf("copy tenant admin role template to tenant err: %v", err)
			return err
		}

		// Create tenant admin user
		wf.NewBranch().OnRollback(func(bb *dtmcli.BranchBarrier) error {
			if _, err = s.userServiceClient.Delete(wf.Context, &identityV1.DeleteUserRequest{
				DeleteBy: &identityV1.DeleteUserRequest_Username{Username: req.User.GetUsername()},
			}); err != nil {
				s.log.Errorf("工作流回滚，创建租户管理员用户失败: %v", err)
				return identityV1.ErrorInternalServerError("工作流回滚，创建租户管理员用户失败")
			}
			return nil
		})
		var adminUser *identityV1.User
		req.User.RoleId = role.Id
		if adminUser, err = s.userServiceClient.Create(wf.Context, &identityV1.CreateUserRequest{Data: req.User}); err != nil {
			s.log.Errorf("create tenant admin user err: %v", err)
			return err
		}

		// Create user credential
		wf.NewBranch().OnRollback(func(bb *dtmcli.BranchBarrier) error {
			if _, err = s.userCredentialServiceClient.Delete(wf.Context, &authenticationV1.DeleteUserCredentialRequest{
				DeleteBy: &authenticationV1.DeleteUserCredentialRequest_UserId{UserId: adminUser.GetId()},
			}); err != nil {
				s.log.Errorf("工作流回滚，创建租户管理员用户认证信息失败: %v", err)
				return identityV1.ErrorInternalServerError("工作流回滚，创建租户管理员用户认证信息失败")
			}
			return nil
		})
		if _, err = s.userCredentialServiceClient.Create(wf.Context, &authenticationV1.CreateUserCredentialRequest{
			Data: &authenticationV1.UserCredential{
				UserId:         adminUser.Id,
				TenantId:       tenant.Id,
				IdentityType:   authenticationV1.UserCredential_USERNAME.Enum(),
				Identifier:     req.User.Username,
				CredentialType: authenticationV1.UserCredential_PASSWORD_HASH.Enum(),
				Credential:     trans.Ptr(req.GetPassword()),
				IsPrimary:      trans.Ptr(true),
				Status:         authenticationV1.UserCredential_ENABLED.Enum(),
				CreatedBy:      req.User.CreatedBy,
			},
		}); err != nil {
			s.log.Errorf("create tenant admin user credential err: %v", err)
			return err
		}

		// assign role to admin user
		wf.NewBranch().OnRollback(func(bb *dtmcli.BranchBarrier) error {
			if _, err = s.roleServiceClient.UnassignRolesFromUser(wf.Context, &permissionV1.UnassignRolesFromUserRequest{
				UserId:   adminUser.GetId(),
				TenantId: tenant.GetId(),
				RoleIds: []uint32{
					role.GetId(),
				},
				OperatorId: req.User.GetCreatedBy(),
				Reason:     trans.Ptr("工作流回滚，撤销租户管理员用户角色分配"),
			}); err != nil {
				s.log.Errorf("工作流回滚，分配租户管理员用户角色失败: %v", err)
				return identityV1.ErrorInternalServerError("工作流回滚，分配租户管理员用户角色失败")
			}
			return nil
		})
		if _, err = s.roleServiceClient.AssignRolesToUser(wf.Context, &permissionV1.AssignRolesToUserRequest{
			UserId:   adminUser.GetId(),
			TenantId: tenant.GetId(),
			RoleIds: []uint32{
				role.GetId(),
			},
			OperatorId: req.User.GetCreatedBy(),
			Reason:     "初始化租户管理员用户角色分配",
		}); err != nil {
			s.log.Errorf("assign role to admin user err: %v", err)
			return err
		}

		// assign admin user id to tenant
		if _, err = s.tenantServiceClient.AssignTenantAdmin(wf.Context, &identityV1.AssignTenantAdminRequest{
			TenantId: tenant.GetId(),
			UserId:   adminUser.GetId(),
		}); err != nil {
			s.log.Errorf("assign admin user id to tenant err: %v", err)
			return err
		}

		return nil
	})
	if err != nil {
		s.log.Errorf("工作流[%s] 注册失败: %v", WorkflowTenantServiceCreateTenantWithAdminSAGA, err)
	}

	return err
}

func (s *TenantService) CreateTenantWithAdminUser(ctx context.Context, req *identityV1.CreateTenantWithAdminUserRequest) (*emptypb.Empty, error) {
	if req.Tenant == nil || req.User == nil {
		s.log.Error("invalid parameter: tenant or user is nil", req)
		return nil, adminV1.ErrorBadRequest("invalid parameter")
	}

	// 获取操作人信息
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.Tenant.CreatedBy = trans.Ptr(operator.UserId)
	req.User.CreatedBy = trans.Ptr(operator.UserId)

	// Check if tenant code or tenant name already exists
	if _, err = s.tenantServiceClient.TenantExists(ctx, &identityV1.TenantExistsRequest{
		Code: req.GetTenant().GetCode(),
		Name: req.GetTenant().GetName(),
	}); err != nil {
		s.log.Errorf("check tenant code exists err: %v", err)
		return nil, err
	}

	// Check if admin user exists
	if _, err = s.userServiceClient.UserExists(ctx, &identityV1.UserExistsRequest{
		QueryBy: &identityV1.UserExistsRequest_Username{
			Username: req.GetUser().GetUsername(),
		},
	}); err != nil {
		s.log.Errorf("check admin user exists err: %v", err)
		return nil, err
	}

	// 生成全局唯一事务 ID (GID)
	gid := id.NewGUIDv7(false)

	s.log.Infof("开始SAGA工作流事务，GID: %s", gid)

	// 提交工作流
	if _, err = workflow.ExecuteCtx(ctx, WorkflowTenantServiceCreateTenantWithAdminSAGA, gid, dtmgimp.MustProtoMarshal(req)); err != nil {
		s.log.Errorf("SAGA工作流事务提交失败: %v", err)
		return nil, identityV1.ErrorInternalServerError("SAGA工作流事务提交失败")
	}

	s.log.Infof("SAGA工作流事务提交成功，GID: %s", gid)

	return &emptypb.Empty{}, nil
}
