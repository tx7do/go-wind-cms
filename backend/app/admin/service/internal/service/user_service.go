package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/dtm-labs/dtm/client/dtmcli"
	"github.com/dtm-labs/dtm/client/dtmgrpc/dtmgimp"
	"github.com/dtm-labs/dtm/client/workflow"

	"github.com/tx7do/go-utils/id"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/log"

	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	paginationFilter "github.com/tx7do/go-crud/pagination/filter"
	"github.com/tx7do/go-utils/aggregator"
	"github.com/tx7do/go-utils/sliceutil"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	authenticationV1 "go-wind-cms/api/gen/go/authentication/service/v1"
	identityV1 "go-wind-cms/api/gen/go/identity/service/v1"
	permissionV1 "go-wind-cms/api/gen/go/permission/service/v1"

	"go-wind-cms/pkg/constants"
	"go-wind-cms/pkg/middleware/auth"
	"go-wind-cms/pkg/utils"
)

const (
	WorkflowUserServiceCreateUserSAGA = "workflow_create_user_saga"
)

type UserService struct {
	adminV1.UserServiceHTTPServer

	log *log.Helper

	codec encoding.Codec

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
		codec:                       encoding.GetCodec("proto"),
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
	var err error

	// Register SAGA workflow
	if err = s.registerWorkflowCreateUser(); err != nil {
		s.log.Fatalf("注册工作流失败: %v", err)
		return
	}
}

func (s *UserService) extractRelationIDs(
	users []*identityV1.User,
	roleSet aggregator.ResourceMap[uint32, *permissionV1.Role],
	tenantSet aggregator.ResourceMap[uint32, *identityV1.Tenant],
	orgUnitSet aggregator.ResourceMap[uint32, *identityV1.OrgUnit],
	posSet aggregator.ResourceMap[uint32, *identityV1.Position],
) {
	for _, v := range users {
		if v == nil {
			continue
		}

		if tenantId := v.GetTenantId(); tenantId > 0 {
			tenantSet[tenantId] = nil
		}

		for _, roleId := range v.RoleIds {
			if roleId > 0 {
				roleSet[roleId] = nil
			}
		}

		if v.GetOrgUnitId() > 0 {
			orgUnitSet[v.GetOrgUnitId()] = nil
		}
		if len(v.OrgUnitIds) > 0 {
			for _, orgID := range v.OrgUnitIds {
				if orgID > 0 {
					orgUnitSet[orgID] = nil
				}
			}
		}

		if v.GetPositionId() > 0 {
			posSet[v.GetPositionId()] = nil
		}
		if len(v.PositionIds) > 0 {
			for _, posID := range v.PositionIds {
				if posID > 0 {
					posSet[posID] = nil
				}
			}
		}

	}
}

func (s *UserService) fetchRelationInfo(
	ctx context.Context,
	roleSet aggregator.ResourceMap[uint32, *permissionV1.Role],
	tenantSet aggregator.ResourceMap[uint32, *identityV1.Tenant],
	orgUnitSet aggregator.ResourceMap[uint32, *identityV1.OrgUnit],
	posSet aggregator.ResourceMap[uint32, *identityV1.Position],
) error {
	if len(roleSet) > 0 {
		roleIds := make([]uint32, 0, len(roleSet))
		for roleId := range roleSet {
			roleIds = append(roleIds, roleId)
		}

		roles, err := s.roleServiceClient.List(ctx, &paginationV1.PagingRequest{
			NoPaging: trans.Ptr(true),
			FilteringType: &paginationV1.PagingRequest_Query{
				Query: fmt.Sprintf(`{"id__in": "[%s]"}`, strings.Join(
					sliceutil.Map(roleIds, func(value uint32, _ int, _ []uint32) string {
						return strconv.FormatUint(uint64(value), 10)
					}),
					","),
				),
			},
		})
		if err != nil {
			s.log.Errorf("query roles err: %v", err)
			return err
		}

		for _, role := range roles.Items {
			roleSet[role.GetId()] = role
		}
	}

	if len(tenantSet) > 0 {
		tenantIds := make([]uint32, 0, len(tenantSet))
		for tenantId := range tenantSet {
			tenantIds = append(tenantIds, tenantId)
		}

		tenants, err := s.tenantServiceClient.List(ctx, &paginationV1.PagingRequest{
			NoPaging: trans.Ptr(true),
			FilteringType: &paginationV1.PagingRequest_Query{
				Query: fmt.Sprintf(`{"id__in": "[%s]"}`, strings.Join(
					sliceutil.Map(tenantIds, func(value uint32, _ int, _ []uint32) string {
						return strconv.FormatUint(uint64(value), 10)
					}),
					","),
				),
			},
		})
		if err != nil {
			s.log.Errorf("query tenants err: %v", err)
			return err
		}

		for _, tenant := range tenants.Items {
			tenantSet[tenant.GetId()] = tenant
		}
	}

	if len(orgUnitSet) > 0 {
		orgUnitIds := make([]uint32, 0, len(orgUnitSet))
		for i := range orgUnitSet {
			orgUnitIds = append(orgUnitIds, i)
		}

		orgUnits, err := s.orgUnitServiceClient.List(ctx, &paginationV1.PagingRequest{
			NoPaging: trans.Ptr(true),
			FilteringType: &paginationV1.PagingRequest_Query{
				Query: fmt.Sprintf(`{"id__in": "[%s]"}`, strings.Join(
					sliceutil.Map(orgUnitIds, func(value uint32, _ int, _ []uint32) string {
						return strconv.FormatUint(uint64(value), 10)
					}),
					","),
				),
			},
		})
		if err != nil {
			s.log.Errorf("query orgUnits err: %v", err)
			return err
		}

		for _, orgUnit := range orgUnits.Items {
			orgUnitSet[orgUnit.GetId()] = orgUnit
		}
	}

	if len(posSet) > 0 {
		posIds := make([]uint32, 0, len(posSet))
		for i := range posSet {
			posIds = append(posIds, i)
		}

		positions, err := s.positionServiceClient.List(ctx, &paginationV1.PagingRequest{
			NoPaging: trans.Ptr(true),
			FilteringType: &paginationV1.PagingRequest_Query{
				Query: fmt.Sprintf(`{"id__in": "[%s]"}`, strings.Join(
					sliceutil.Map(posIds, func(value uint32, _ int, _ []uint32) string {
						return strconv.FormatUint(uint64(value), 10)
					}),
					","),
				),
			},
		})
		if err != nil {
			s.log.Errorf("query positions err: %v", err)
			return err
		}

		for _, position := range positions.Items {
			posSet[position.GetId()] = position
		}
	}

	return nil
}

func (s *UserService) bindRelations(
	users []*identityV1.User,
	roleSet aggregator.ResourceMap[uint32, *permissionV1.Role],
	tenantSet aggregator.ResourceMap[uint32, *identityV1.Tenant],
	orgUnitSet aggregator.ResourceMap[uint32, *identityV1.OrgUnit],
	posSet aggregator.ResourceMap[uint32, *identityV1.Position],
) {
	aggregator.PopulateMulti(
		users,
		roleSet,
		func(ou *identityV1.User) []uint32 { return ou.GetRoleIds() },
		func(ou *identityV1.User, r []*permissionV1.Role) {
			for _, role := range r {
				ou.RoleNames = append(ou.RoleNames, role.GetName())
				ou.Roles = append(ou.Roles, role.GetCode())
			}
		},
	)
	aggregator.Populate(
		users,
		roleSet,
		func(ou *identityV1.User) uint32 { return ou.GetRoleId() },
		func(ou *identityV1.User, r *permissionV1.Role) {
			ou.RoleNames = append(ou.RoleNames, r.GetName())
			ou.Roles = append(ou.Roles, r.GetCode())
		},
	)

	aggregator.Populate(
		users,
		tenantSet,
		func(ou *identityV1.User) uint32 { return ou.GetTenantId() },
		func(ou *identityV1.User, r *identityV1.Tenant) {
			ou.TenantName = r.Name
		},
	)

	aggregator.PopulateMulti(
		users,
		posSet,
		func(ou *identityV1.User) []uint32 { return ou.GetPositionIds() },
		func(ou *identityV1.User, r []*identityV1.Position) {
			for _, pos := range r {
				ou.PositionNames = append(ou.PositionNames, pos.GetName())
			}
		},
	)
	aggregator.Populate(
		users,
		posSet,
		func(ou *identityV1.User) uint32 { return ou.GetPositionId() },
		func(ou *identityV1.User, r *identityV1.Position) {
			ou.PositionName = r.Name
		},
	)

	aggregator.PopulateMulti(
		users,
		orgUnitSet,
		func(ou *identityV1.User) []uint32 { return ou.GetOrgUnitIds() },
		func(ou *identityV1.User, orgs []*identityV1.OrgUnit) {
			for _, org := range orgs {
				ou.OrgUnitNames = append(ou.OrgUnitNames, org.GetName())
			}
		},
	)
	aggregator.Populate(
		users,
		orgUnitSet,
		func(ou *identityV1.User) uint32 { return ou.GetOrgUnitId() },
		func(ou *identityV1.User, org *identityV1.OrgUnit) {
			ou.OrgUnitName = org.Name
		},
	)
}

func (s *UserService) enrichRelations(ctx context.Context, users []*identityV1.User) error {
	var roleSet = make(aggregator.ResourceMap[uint32, *permissionV1.Role])
	var tenantSet = make(aggregator.ResourceMap[uint32, *identityV1.Tenant])
	var orgUnitSet = make(aggregator.ResourceMap[uint32, *identityV1.OrgUnit])
	var posSet = make(aggregator.ResourceMap[uint32, *identityV1.Position])

	s.extractRelationIDs(users, roleSet, tenantSet, orgUnitSet, posSet)
	if err := s.fetchRelationInfo(ctx, roleSet, tenantSet, orgUnitSet, posSet); err != nil {
		return err
	}
	s.bindRelations(users, roleSet, tenantSet, orgUnitSet, posSet)
	return nil
}

// FilterFields 过滤掉不需要的字段条件
func FilterFields(filterExpr *paginationV1.FilterExpr, excludeFields []string) []*paginationV1.FilterCondition {
	if filterExpr == nil || len(filterExpr.Conditions) == 0 {
		return []*paginationV1.FilterCondition{}
	}

	exclude := make(map[string]struct{}, len(excludeFields))
	for _, f := range excludeFields {
		if f == "" {
			continue
		}
		exclude[f] = struct{}{}
	}

	includeConditions := make([]*paginationV1.FilterCondition, 0, len(filterExpr.Conditions))
	excludeConditions := make([]*paginationV1.FilterCondition, 0, len(filterExpr.Conditions))
	for _, cond := range filterExpr.Conditions {
		if cond == nil || cond.Field == "" {
			continue
		}
		if _, skip := exclude[cond.Field]; skip {
			excludeConditions = append(excludeConditions, cond)
			continue
		}
		includeConditions = append(includeConditions, cond)
	}

	filterExpr.Conditions = includeConditions

	return excludeConditions
}

func (s *UserService) queryUserIDsByRelationIDs(ctx context.Context, roleIDs []uint32, orgUnitIDs []uint32, positionIDs []uint32) ([]uint32, error) {
	if len(roleIDs) == 0 && len(orgUnitIDs) == 0 && len(positionIDs) == 0 {
		return nil, nil
	}

	switch constants.DefaultUserTenantRelationType {
	default:
		fallthrough
	case constants.UserTenantRelationOneToOne:
		return s.queryUserIDsByRelationIDsUserTenantRelationOneToOne(ctx, roleIDs, orgUnitIDs, positionIDs)
	case constants.UserTenantRelationOneToMany:
		return s.queryUserIDsByRelationIDsUserTenantRelationOneToMany(ctx, roleIDs, orgUnitIDs, positionIDs)
	}
}

func (s *UserService) queryUserIDsByRelationIDsUserTenantRelationOneToMany(_ context.Context, _, _, _ []uint32) ([]uint32, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *UserService) queryUserIDsByRelationIDsUserTenantRelationOneToOne(ctx context.Context, roleIDs, orgUnitIDs, positionIDs []uint32) ([]uint32, error) {
	if len(roleIDs) == 0 && len(orgUnitIDs) == 0 && len(positionIDs) == 0 {
		return nil, nil
	}

	var orgUnitUserIDs []uint32
	var positionUserIDs []uint32
	var roleUserIDs []uint32
	if len(orgUnitIDs) > 0 {
		resp, err := s.userServiceClient.ListUserIDsByOrgUnitIDs(ctx, &identityV1.ListUserIDsByOrgUnitIDsRequest{
			OrgUnitIds: orgUnitIDs,
		})
		if err != nil {
			return nil, err
		}
		if resp != nil {
			orgUnitUserIDs = resp.GetUserIds()
		}
	}
	if len(positionIDs) > 0 {
		resp, err := s.userServiceClient.ListUserIDsByPositionIDs(ctx, &identityV1.ListUserIDsByPositionIDsRequest{
			PositionIds: positionIDs,
		})
		if err != nil {
			return nil, err
		}
		if resp != nil {
			positionUserIDs = resp.GetUserIds()
		}
	}
	if len(roleIDs) > 0 {
		resp, err := s.roleServiceClient.ListUserIDsByRoleIDs(ctx, &permissionV1.ListUserIDsByRoleIDsRequest{
			RoleIds: roleIDs,
		})
		if err != nil {
			return nil, err
		}

		if resp != nil {
			roleUserIDs = resp.GetUserIds()
		}
	}

	// 收集所有非空列表用于求交集
	lists := make([][]uint32, 0, 3)
	if orgUnitUserIDs != nil {
		lists = append(lists, orgUnitUserIDs)
	}
	if positionUserIDs != nil {
		lists = append(lists, positionUserIDs)
	}
	if roleUserIDs != nil {
		lists = append(lists, roleUserIDs)
	}

	// 如果没有任何实际列表（例如对应 ids 为空导致查询未执行），返回空
	if len(lists) == 0 {
		return []uint32{}, nil
	}

	// 逐步求交集
	result := lists[0]
	for i := 1; i < len(lists); i++ {
		result = sliceutil.Intersection(result, lists[i])
		if len(result) == 0 {
			break
		}
	}

	return result, nil
}

func (s *UserService) List(ctx context.Context, req *paginationV1.PagingRequest) (*identityV1.ListUserResponse, error) {
	filterExpr, err := paginationFilter.ConvertFilterByPagingRequest(req)
	if err != nil {
		s.log.Errorf("convert filter by paging request failed: %s", err.Error())
		return nil, err
	}

	excludeConditions := FilterFields(filterExpr, []string{
		"org_unit_id", "org_unit_ids",
		"position_id", "position_ids",
		"role_id", "role_ids",
	})

	var orgUnitIDs []uint32
	var positionIDs []uint32
	var roleIDs []uint32
	for _, cond := range excludeConditions {
		//r.log.Debugf("excluding filter condition: field=%s operator=%s value=%v", cond.GetField(), cond.GetOp(), cond.GetValue())

		var val uint64
		switch cond.GetField() {
		case "org_unit_id":
			if val, err = strconv.ParseUint(cond.GetValue(), 10, 64); err == nil {
				orgUnitIDs = append(orgUnitIDs, uint32(val))
			} else {
				s.log.Errorf("parse org_unit_id value failed: %s", err.Error())
			}
		case "org_unit_ids":
			for _, v := range cond.GetValues() {
				if val, err = strconv.ParseUint(v, 10, 64); err == nil {
					orgUnitIDs = append(orgUnitIDs, uint32(val))
				} else {
					s.log.Errorf("parse org_unit_ids value failed: %s", err.Error())
				}
			}

		case "position_id":
			if val, err = strconv.ParseUint(cond.GetValue(), 10, 64); err == nil {
				positionIDs = append(positionIDs, uint32(val))
			} else {
				s.log.Errorf("parse position_id value failed: %s", err.Error())
			}
		case "position_ids":
			for _, v := range cond.GetValues() {
				if val, err = strconv.ParseUint(v, 10, 64); err == nil {
					positionIDs = append(positionIDs, uint32(val))
				} else {
					s.log.Errorf("parse position_ids value failed: %s", err.Error())
				}
			}

		case "role_id":
			if val, err = strconv.ParseUint(cond.GetValue(), 10, 64); err == nil {
				roleIDs = append(roleIDs, uint32(val))
			} else {
				s.log.Errorf("parse role_id value failed: %s", err.Error())
			}
		case "role_ids":
			for _, v := range cond.GetValues() {
				if val, err = strconv.ParseUint(v, 10, 64); err == nil {
					roleIDs = append(roleIDs, uint32(val))
				} else {
					s.log.Errorf("parse role_ids value failed: %s", err.Error())
				}
			}
		}
	}

	var mergedUserIDs []uint32
	mergedUserIDs, err = s.queryUserIDsByRelationIDs(ctx, roleIDs, orgUnitIDs, positionIDs)
	if err != nil {
		s.log.Errorf("query user ids by relation ids failed: %s", err.Error())
		return nil, err
	}

	hasRelationFilter := len(roleIDs) > 0 || len(orgUnitIDs) > 0 || len(positionIDs) > 0
	if hasRelationFilter && len(mergedUserIDs) == 0 {
		// 如果有关系过滤条件但没有匹配的用户ID，直接返回空结果
		return &identityV1.ListUserResponse{Total: 0, Items: nil}, nil
	}

	if len(mergedUserIDs) > 0 {
		filterExpr.Conditions = append(filterExpr.Conditions, &paginationV1.FilterCondition{
			Field: "id",
			Op:    paginationV1.Operator_IN,
			Values: func() []string {
				values := make([]string, 0, len(mergedUserIDs))
				for _, d := range mergedUserIDs {
					values = append(values, strconv.FormatUint(uint64(d), 10))
				}
				return values
			}(),
		})
	}

	req.FilteringType = &paginationV1.PagingRequest_FilterExpr{FilterExpr: filterExpr}

	resp, err := s.userServiceClient.List(ctx, req)
	if err != nil {
		return nil, err
	}

	for _, user := range resp.Items {
		roleIdsResp, err := s.roleServiceClient.ListUserRoleIDs(ctx, &permissionV1.ListUserRoleIDsRequest{
			QueryBy: &permissionV1.ListUserRoleIDsRequest_UserId{
				UserId: user.GetId(),
			},
		})
		if err != nil {
			s.log.Errorf("list user role ids failed [%s]", err.Error())
			return nil, adminV1.ErrorInternalServerError("list user role ids failed")
		}

		user.RoleIds = roleIdsResp.GetRoleIds()
	}

	_ = s.enrichRelations(ctx, resp.Items)

	return resp, nil
}

func (s *UserService) Count(ctx context.Context, req *paginationV1.PagingRequest) (*identityV1.CountUserResponse, error) {
	return s.userServiceClient.Count(ctx, req)
}

func (s *UserService) Get(ctx context.Context, req *identityV1.GetUserRequest) (*identityV1.User, error) {
	resp, err := s.userServiceClient.Get(ctx, req)
	if err != nil {
		return nil, err
	}

	roleIdsResp, err := s.roleServiceClient.ListUserRoleIDs(ctx, &permissionV1.ListUserRoleIDsRequest{
		QueryBy: &permissionV1.ListUserRoleIDsRequest_UserId{
			UserId: resp.GetId(),
		},
	})
	if err != nil {
		s.log.Errorf("list user role ids failed [%s]", err.Error())
		return nil, adminV1.ErrorInternalServerError("list user role ids failed")
	}

	resp.RoleIds = roleIdsResp.GetRoleIds()

	fakeItems := []*identityV1.User{resp}
	_ = s.enrichRelations(ctx, fakeItems)

	return resp, nil
}

func (s *UserService) registerWorkflowCreateUser() (err error) {
	err = workflow.Register2(WorkflowUserServiceCreateUserSAGA, func(wf *workflow.Workflow, data []byte) ([]byte, error) {
		var req identityV1.CreateUserRequest
		if len(data) > 0 {
			if err = s.codec.Unmarshal(data, &req); err != nil {
				s.log.Errorf("工作流数据反序列化失败: %v", err)
				return nil, identityV1.ErrorInternalServerError("工作流数据反序列化失败")
			}
		}

		// Create user
		wf.NewBranch().OnRollback(func(bb *dtmcli.BranchBarrier) error {
			if _, err = s.userServiceClient.Delete(wf.Context, &identityV1.DeleteUserRequest{
				DeleteBy: &identityV1.DeleteUserRequest_Username{Username: req.Data.GetUsername()},
			}); err != nil {
				s.log.Errorf("工作流回滚，创建用户失败: %v", err)
				return identityV1.ErrorInternalServerError("工作流回滚，创建用户失败")
			}
			return nil
		})
		var outUser *identityV1.User
		if outUser, err = s.userServiceClient.Create(wf.Context, &identityV1.CreateUserRequest{
			Data: req.Data,
		}); err != nil {
			log.Errorf("create user err: %v", err)
			return nil, err
		}

		// Create user credential
		if len(req.GetPassword()) > 0 {
			wf.NewBranch().OnRollback(func(bb *dtmcli.BranchBarrier) error {
				if _, err = s.userCredentialServiceClient.Delete(wf.Context, &authenticationV1.DeleteUserCredentialRequest{
					DeleteBy: &authenticationV1.DeleteUserCredentialRequest_UserId{UserId: outUser.GetId()},
				}); err != nil {
					s.log.Errorf("工作流回滚，创建用户认证信息失败: %v", err)
					return identityV1.ErrorInternalServerError("工作流回滚，创建用户认证信息失败")
				}
				return nil
			})
			if _, err = s.userCredentialServiceClient.Create(wf.Context, &authenticationV1.CreateUserCredentialRequest{
				Data: &authenticationV1.UserCredential{
					UserId:         outUser.Id,
					IdentityType:   authenticationV1.UserCredential_USERNAME.Enum(),
					Identifier:     outUser.Username,
					CredentialType: authenticationV1.UserCredential_PASSWORD_HASH.Enum(),
					Credential:     trans.Ptr(req.GetPassword()),
					IsPrimary:      trans.Ptr(true),
					Status:         authenticationV1.UserCredential_ENABLED.Enum(),
				},
			}); err != nil {
				log.Errorf("create user credential err: %v", err)
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

		// assign role to user
		wf.NewBranch().OnRollback(func(bb *dtmcli.BranchBarrier) error {
			if _, err = s.roleServiceClient.UnassignRolesFromUser(wf.Context, &permissionV1.UnassignRolesFromUserRequest{
				UserId:     outUser.GetId(),
				TenantId:   outUser.GetTenantId(),
				RoleIds:    roleIds,
				OperatorId: req.Data.GetCreatedBy(),
				Reason:     trans.Ptr("工作流回滚，撤销用户角色分配"),
			}); err != nil {
				s.log.Errorf("工作流回滚，分配用户角色失败: %v", err)
				return identityV1.ErrorInternalServerError("工作流回滚，分配用户角色失败")
			}
			return nil
		})
		if _, err = s.roleServiceClient.AssignRolesToUser(wf.Context, &permissionV1.AssignRolesToUserRequest{
			UserId:     outUser.GetId(),
			TenantId:   outUser.GetTenantId(),
			RoleIds:    roleIds,
			OperatorId: req.Data.GetCreatedBy(),
			Reason:     "初始创建用户，分配角色",
		}); err != nil {
			s.log.Errorf("assign roles to user err: %v", err)
			return nil, err
		}

		return dtmgimp.MustProtoMarshal(outUser), nil
	})
	if err != nil {
		s.log.Errorf("工作流[%s] 注册失败: %v", WorkflowUserServiceCreateUserSAGA, err)
	}

	return err
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

	if len(req.GetPassword()) == 0 {
		// 如果没有设置密码，则设置为默认密码。
		req.Password = trans.Ptr(constants.DefaultUserPassword)
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

	req.Data.RoleId = nil
	req.Data.RoleIds = roleIds

	// 生成全局唯一事务 ID (GID)
	gid := id.NewGUIDv7(false)

	s.log.Infof("开始SAGA工作流事务，GID: %s", gid)

	// 提交工作流
	var outData []byte
	if outData, err = workflow.ExecuteCtx(ctx, WorkflowUserServiceCreateUserSAGA, gid, dtmgimp.MustProtoMarshal(req)); err != nil {
		s.log.Errorf("SAGA工作流事务提交失败: %v", err)
		return nil, identityV1.ErrorInternalServerError("SAGA工作流事务提交失败")
	}

	s.log.Infof("SAGA工作流事务提交成功，GID: %s", gid)

	var outUser *identityV1.User
	outUser = &identityV1.User{}
	if err = s.codec.Unmarshal(outData, outUser); err != nil {
		s.log.Errorf("工作流结果反序列化失败: %v", err)
		return nil, identityV1.ErrorInternalServerError("工作流结果反序列化失败")
	}

	return outUser, nil
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
