package service

import (
	"context"
	"fmt"
	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	"strconv"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	paginationFilter "github.com/tx7do/go-crud/pagination/filter"
	"github.com/tx7do/go-utils/aggregator"
	"github.com/tx7do/go-utils/sliceutil"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	"go-wind-cms/app/core/service/internal/data"

	authenticationV1 "go-wind-cms/api/gen/go/authentication/service/v1"
	identityV1 "go-wind-cms/api/gen/go/identity/service/v1"
	permissionV1 "go-wind-cms/api/gen/go/permission/service/v1"

	"go-wind-cms/pkg/constants"
	appViewer "go-wind-cms/pkg/entgo/viewer"
)

type UserService struct {
	identityV1.UnimplementedUserServiceServer

	log *log.Helper

	userRepo           data.UserRepo
	userCredentialRepo *data.UserCredentialRepo

	roleRepo     *data.RoleRepo
	positionRepo *data.PositionRepo
	orgUnitRepo  *data.OrgUnitRepo
	tenantRepo   *data.TenantRepo

	membershipRepo *data.MembershipRepo
}

func NewUserService(
	ctx *bootstrap.Context,
	userRepo data.UserRepo,
	roleRepo *data.RoleRepo,
	userCredentialRepo *data.UserCredentialRepo,
	positionRepo *data.PositionRepo,
	orgUnitRepo *data.OrgUnitRepo,
	tenantRepo *data.TenantRepo,
	membershipRepo *data.MembershipRepo,
) *UserService {
	svc := &UserService{
		log:                ctx.NewLoggerHelper("user/service/core-service"),
		userRepo:           userRepo,
		roleRepo:           roleRepo,
		userCredentialRepo: userCredentialRepo,
		positionRepo:       positionRepo,
		orgUnitRepo:        orgUnitRepo,
		tenantRepo:         tenantRepo,
		membershipRepo:     membershipRepo,
	}

	svc.init()

	return svc
}

func (s *UserService) init() {
	ctx := appViewer.NewSystemViewerContext(context.Background())
	if count, _ := s.userRepo.Count(ctx, nil); count == 0 {
		_ = s.createDefaultUser(ctx)
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

		if id := v.GetTenantId(); id > 0 {
			tenantSet[id] = nil
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
		for id := range roleSet {
			roleIds = append(roleIds, id)
		}

		roles, err := s.roleRepo.ListRolesByRoleIds(ctx, roleIds)
		if err != nil {
			s.log.Errorf("query roles err: %v", err)
			return err
		}

		for _, role := range roles {
			roleSet[role.GetId()] = role
		}
	}

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

	if len(orgUnitSet) > 0 {
		orgUnitIds := make([]uint32, 0, len(orgUnitSet))
		for id := range orgUnitSet {
			orgUnitIds = append(orgUnitIds, id)
		}

		orgUnits, err := s.orgUnitRepo.ListOrgUnitsByIds(ctx, orgUnitIds)
		if err != nil {
			s.log.Errorf("query orgUnits err: %v", err)
			return err
		}

		for _, orgUnit := range orgUnits {
			orgUnitSet[orgUnit.GetId()] = orgUnit
		}
	}

	if len(posSet) > 0 {
		posIds := make([]uint32, 0, len(posSet))
		for id := range posSet {
			posIds = append(posIds, id)
		}

		positions, err := s.positionRepo.ListPositionByIds(ctx, posIds)
		if err != nil {
			s.log.Errorf("query positions err: %v", err)
			return err
		}

		for _, position := range positions {
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

	var err error

	var orgUnitUserIDs []uint32
	var positionUserIDs []uint32
	var roleUserIDs []uint32
	if len(orgUnitIDs) > 0 {
		orgUnitUserIDs, err = s.userRepo.ListUserIDsByOrgUnitIDs(ctx, orgUnitIDs, true)
		if err != nil {
			return nil, err
		}
	}
	if len(positionIDs) > 0 {
		positionUserIDs, err = s.userRepo.ListUserIDsByPositionIDs(ctx, positionIDs, true)
		if err != nil {
			return nil, err
		}
	}
	if len(roleIDs) > 0 {
		roleUserIDs, err = s.userRepo.ListUserIDsByRoleIDs(ctx, roleIDs, true)
		if err != nil {
			return nil, err
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

	resp, err := s.userRepo.List(ctx, req)
	if err != nil {
		return nil, err
	}

	for _, user := range resp.Items {
		roleIds, err := s.userRepo.ListRoleIDsByUserID(ctx, user.GetId())
		if err != nil {
			s.log.Errorf("list user role ids failed [%s]", err.Error())
			return nil, adminV1.ErrorInternalServerError("list user role ids failed")
		}

		user.RoleIds = roleIds
	}

	_ = s.enrichRelations(ctx, resp.Items)

	return resp, nil
}

func (s *UserService) Count(ctx context.Context, req *paginationV1.PagingRequest) (*identityV1.CountUserResponse, error) {
	count, err := s.userRepo.Count(ctx, req)
	if err != nil {
		return nil, err
	}

	return &identityV1.CountUserResponse{
		Count: uint64(count),
	}, nil
}

func (s *UserService) Get(ctx context.Context, req *identityV1.GetUserRequest) (*identityV1.User, error) {
	resp, err := s.userRepo.Get(ctx, req)
	if err != nil {
		return nil, err
	}

	fakeItems := []*identityV1.User{resp}
	_ = s.enrichRelations(ctx, fakeItems)

	return resp, nil
}

func (s *UserService) Create(ctx context.Context, req *identityV1.CreateUserRequest) (*identityV1.User, error) {
	if req.Data == nil {
		return nil, identityV1.ErrorBadRequest("invalid parameter")
	}

	// 创建用户
	var user *identityV1.User
	var err error
	if user, err = s.userRepo.Create(ctx, req); err != nil {
		return nil, err
	}

	if len(req.GetPassword()) == 0 {
		// 如果没有设置密码，则设置为默认密码。
		req.Password = trans.Ptr(constants.DefaultUserPassword)
	}

	if len(req.GetPassword()) > 0 {
		if err = s.userCredentialRepo.Create(ctx, &authenticationV1.CreateUserCredentialRequest{
			Data: &authenticationV1.UserCredential{
				UserId:   user.Id,
				TenantId: user.TenantId,

				IdentityType: authenticationV1.UserCredential_USERNAME.Enum(),
				Identifier:   req.Data.Username,

				CredentialType: authenticationV1.UserCredential_PASSWORD_HASH.Enum(),
				Credential:     req.Password,

				IsPrimary: trans.Ptr(true),
				Status:    authenticationV1.UserCredential_ENABLED.Enum(),
			},
		}); err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (s *UserService) Update(ctx context.Context, req *identityV1.UpdateUserRequest) (*emptypb.Empty, error) {
	if req.Data == nil {
		return nil, identityV1.ErrorBadRequest("invalid parameter")
	}

	// 更新用户
	if err := s.userRepo.Update(ctx, req); err != nil {
		s.log.Error(err)
		return nil, err
	}

	if len(req.GetPassword()) > 0 {
		if err := s.userCredentialRepo.ResetCredential(ctx, &authenticationV1.ResetCredentialRequest{
			IdentityType:  authenticationV1.UserCredential_USERNAME,
			Identifier:    req.Data.GetUsername(),
			NewCredential: req.GetPassword(),
		}); err != nil {
			return nil, err
		}
	}

	return &emptypb.Empty{}, nil
}

func (s *UserService) Delete(ctx context.Context, req *identityV1.DeleteUserRequest) (*emptypb.Empty, error) {
	err := s.userRepo.Delete(ctx, req)
	return &emptypb.Empty{}, err
}

func (s *UserService) UserExists(ctx context.Context, req *identityV1.UserExistsRequest) (*identityV1.UserExistsResponse, error) {
	return s.userRepo.UserExists(ctx, req)
}

// EditUserPassword 修改用户密码
func (s *UserService) EditUserPassword(ctx context.Context, req *identityV1.EditUserPasswordRequest) (*emptypb.Empty, error) {
	// 获取操作者的用户信息
	u, err := s.userRepo.Get(ctx, &identityV1.GetUserRequest{
		QueryBy: &identityV1.GetUserRequest_Id{
			Id: req.GetUserId(),
		},
	})
	if err != nil {
		return nil, err
	}

	if err = s.userCredentialRepo.ResetCredential(ctx, &authenticationV1.ResetCredentialRequest{
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

// createDefaultUser 创建默认用户，即超级用户
func (s *UserService) createDefaultUser(ctx context.Context) error {
	var err error

	// 创建默认用户
	for _, user := range constants.DefaultUsers {
		if _, err = s.userRepo.Create(ctx, &identityV1.CreateUserRequest{
			Data: user,
		}); err != nil {
			s.log.Errorf("create default user err: %v", err)
			return err
		}
	}

	// 创建默认用户凭证
	for _, userCredential := range constants.DefaultUserCredentials {
		if err = s.userCredentialRepo.Create(ctx, &authenticationV1.CreateUserCredentialRequest{
			Data: userCredential,
		}); err != nil {
			s.log.Errorf("create default user credential err: %v", err)
			return err
		}
	}

	switch constants.DefaultUserTenantRelationType {
	default:
		fallthrough
	case constants.UserTenantRelationOneToOne:
		// 创建默认用户角色关联关系
		for _, userRole := range constants.DefaultUserRoles {
			if err = s.userRepo.AssignUserRole(ctx, userRole); err != nil {
				s.log.Errorf("create default user role relation err: %v", err)
				return err
			}
		}

	case constants.UserTenantRelationOneToMany:
		// 创建默认用户租户关联关系
		for _, membership := range constants.DefaultMemberships {
			if err = s.membershipRepo.AssignTenantMembershipWith(ctx, membership); err != nil {
				s.log.Errorf("create default user membership err: %v", err)
				return err
			}
		}
	}

	return err
}
