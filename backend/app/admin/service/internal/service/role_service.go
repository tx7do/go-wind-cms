package service

import (
	"context"
	"fmt"
	"go-wind-cms/pkg/utils"
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
	resp, err := s.roleServiceClient.List(ctx, req)
	if err != nil {
		return nil, err
	}

	_ = s.enrichRelations(ctx, resp.Items)

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

	fakeItems := []*permissionV1.Role{resp}
	_ = s.enrichRelations(ctx, fakeItems)

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
