package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/go-utils/aggregator"
	"github.com/tx7do/go-utils/sliceutil"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	authenticationV1 "go-wind-cms/api/gen/go/authentication/service/v1"
	identityV1 "go-wind-cms/api/gen/go/identity/service/v1"
	permissionV1 "go-wind-cms/api/gen/go/permission/service/v1"

	"go-wind-cms/pkg/middleware/auth"
)

type UserProfileService struct {
	adminV1.UserProfileServiceHTTPServer

	log *log.Helper

	userServiceClient     identityV1.UserServiceClient
	tenantServiceClient   identityV1.TenantServiceClient
	orgUnitServiceClient  identityV1.OrgUnitServiceClient
	positionServiceClient identityV1.PositionServiceClient

	roleServiceClient permissionV1.RoleServiceClient

	userCredentialServiceClient authenticationV1.UserCredentialServiceClient
}

func NewUserProfileService(
	ctx *bootstrap.Context,
	userServiceClient identityV1.UserServiceClient,
	tenantServiceClient identityV1.TenantServiceClient,
	orgUnitServiceClient identityV1.OrgUnitServiceClient,
	positionServiceClient identityV1.PositionServiceClient,
	roleServiceClient permissionV1.RoleServiceClient,
	userCredentialServiceClient authenticationV1.UserCredentialServiceClient,
) *UserProfileService {
	return &UserProfileService{
		log:                         ctx.NewLoggerHelper("user-profile/service/admin-service"),
		userServiceClient:           userServiceClient,
		tenantServiceClient:         tenantServiceClient,
		orgUnitServiceClient:        orgUnitServiceClient,
		positionServiceClient:       positionServiceClient,
		roleServiceClient:           roleServiceClient,
		userCredentialServiceClient: userCredentialServiceClient,
	}
}

func (s *UserProfileService) extractRelationIDs(
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

func (s *UserProfileService) fetchRelationInfo(
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
		for id := range tenantSet {
			tenantIds = append(tenantIds, id)
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
		for id := range orgUnitSet {
			orgUnitIds = append(orgUnitIds, id)
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
		for id := range posSet {
			posIds = append(posIds, id)
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

func (s *UserProfileService) bindRelations(
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

func (s *UserProfileService) enrichRelations(ctx context.Context, users []*identityV1.User) error {
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

func (s *UserProfileService) GetUser(ctx context.Context, _ *emptypb.Empty) (*identityV1.User, error) {
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := s.userServiceClient.Get(ctx, &identityV1.GetUserRequest{
		QueryBy: &identityV1.GetUserRequest_Id{
			Id: operator.UserId,
		},
	})
	if err != nil {
		return nil, err
	}

	roleIdsResp, err := s.roleServiceClient.ListUserRoleIDs(ctx, &permissionV1.ListUserRoleIDsRequest{
		QueryBy: &permissionV1.ListUserRoleIDsRequest_UserId{
			UserId: operator.GetUserId(),
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

func (s *UserProfileService) UpdateUser(ctx context.Context, req *identityV1.UpdateUserRequest) (*emptypb.Empty, error) {
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.Data.Id = trans.Ptr(operator.UserId)
	req.Id = operator.UserId

	return s.userServiceClient.Update(ctx, req)
}

func (s *UserProfileService) ChangePassword(ctx context.Context, req *identityV1.ChangePasswordRequest) (*emptypb.Empty, error) {
	// 获取操作人信息
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	return s.userCredentialServiceClient.ChangeCredential(ctx, &authenticationV1.ChangeCredentialRequest{
		IdentityType:  authenticationV1.UserCredential_USERNAME,
		Identifier:    operator.GetUsername(),
		OldCredential: req.GetOldPassword(),
		NewCredential: req.GetNewPassword(),
	})
}
