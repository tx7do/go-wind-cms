package service

import (
	"context"
	"sort"
	"strconv"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/go-crud/pagination"
	"github.com/tx7do/go-crud/pagination/filter"

	"github.com/tx7do/go-utils/sliceutil"
	"github.com/tx7do/go-utils/trans"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	permissionV1 "go-wind-cms/api/gen/go/permission/service/v1"
	resourceV1 "go-wind-cms/api/gen/go/resource/service/v1"

	"go-wind-cms/pkg/constants"
	appViewer "go-wind-cms/pkg/entgo/viewer"
	"go-wind-cms/pkg/middleware/auth"
	"go-wind-cms/pkg/utils/converter"
)

type PermissionService struct {
	adminV1.PermissionServiceHTTPServer

	log *log.Helper

	permissionServiceClient      permissionV1.PermissionServiceClient
	permissionGroupServiceClient permissionV1.PermissionGroupServiceClient

	roleServiceClient permissionV1.RoleServiceClient

	apiServiceClient  resourceV1.ApiServiceClient
	menuServiceClient resourceV1.MenuServiceClient

	menuPermissionConverter *converter.MenuPermissionConverter
	apiPermissionConverter  *converter.ApiPermissionConverter
}

func NewPermissionService(
	ctx *bootstrap.Context,
	permissionServiceClient permissionV1.PermissionServiceClient,
	permissionGroupServiceClient permissionV1.PermissionGroupServiceClient,
	roleServiceClient permissionV1.RoleServiceClient,
	apiServiceClient resourceV1.ApiServiceClient,
	menuServiceClient resourceV1.MenuServiceClient,
) *PermissionService {
	svc := &PermissionService{
		log: ctx.NewLoggerHelper("permission/service/admin-service"),

		permissionServiceClient:      permissionServiceClient,
		permissionGroupServiceClient: permissionGroupServiceClient,
		roleServiceClient:            roleServiceClient,
		apiServiceClient:             apiServiceClient,
		menuServiceClient:            menuServiceClient,

		menuPermissionConverter: converter.NewMenuPermissionConverter(),
		apiPermissionConverter:  converter.NewApiPermissionConverter(),
	}

	svc.init()

	return svc
}

func (s *PermissionService) init() {
	ctx := appViewer.NewSystemViewerContext(context.Background())
	if resp, _ := s.permissionServiceClient.Count(ctx, nil); resp.Count == 0 || resp.Count == (uint64)(len(constants.DefaultPermissions)) {
		apiCount, _ := s.apiServiceClient.Count(ctx, nil)

		menusCount, _ := s.menuServiceClient.Count(ctx, nil)

		if apiCount != nil && apiCount.Count > 0 && menusCount != nil && menusCount.Count > 0 {
			_, _ = s.SyncPermissions(ctx, &emptypb.Empty{})
		}
	}
}

func (s *PermissionService) List(ctx context.Context, req *paginationV1.PagingRequest) (*permissionV1.ListPermissionResponse, error) {
	// 获取操作人信息
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	var limitPermissionIDs []uint32
	if operator.GetTenantId() > 0 {
		var limitPermissionResp *permissionV1.ListPermissionIdsResponse
		limitPermissionResp, err = s.roleServiceClient.ListPermissionIds(ctx, &permissionV1.ListPermissionIdsRequest{
			RoleCodes: operator.GetRoles(),
		})
		if err != nil {
			return nil, err
		}

		limitPermissionIDs = limitPermissionResp.PermissionIds

		// 没有任何 permission 可访问，直接返回空列表
		if len(limitPermissionIDs) == 0 {
			return &permissionV1.ListPermissionResponse{
				Items: []*permissionV1.Permission{},
				Total: 0,
			}, nil
		}
	}

	if len(limitPermissionIDs) > 0 {
		filterExpr, err := filter.ConvertFilterByPagingRequest(req)
		if err != nil {
			return nil, err
		}
		if filterExpr != nil {
			pagination.ClearFilterExprByFieldNames(filterExpr, "id")
			req.FilteringType = &paginationV1.PagingRequest_FilterExpr{FilterExpr: filterExpr}
		}

		condition := &paginationV1.FilterCondition{
			Field: "id",
			Op:    paginationV1.Operator_IN,
			Values: sliceutil.Map(limitPermissionIDs, func(value uint32, _ int, _ []uint32) string {
				return strconv.FormatUint(uint64(value), 10)
			}),
		}

		if req.FilteringType == nil {
			req.FilteringType = &paginationV1.PagingRequest_FilterExpr{
				FilterExpr: &paginationV1.FilterExpr{
					Type: paginationV1.ExprType_AND,
					Conditions: []*paginationV1.FilterCondition{
						condition,
					},
				},
			}
		} else {
			req.GetFilterExpr().Conditions = append(req.GetFilterExpr().GetConditions(), condition)
		}
	}

	resp, err := s.permissionServiceClient.List(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *PermissionService) Get(ctx context.Context, req *permissionV1.GetPermissionRequest) (*permissionV1.Permission, error) {
	// 获取操作人信息
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := s.permissionServiceClient.Get(ctx, req)
	if err != nil {
		return nil, err
	}

	if operator.GetTenantId() > 0 {
		var limitPermissionResp *permissionV1.ListPermissionIdsResponse
		limitPermissionResp, err = s.roleServiceClient.ListPermissionIds(ctx, &permissionV1.ListPermissionIdsRequest{
			RoleCodes: operator.GetRoles(),
		})
		if err != nil {
			return nil, err
		}

		found := false
		for _, pid := range limitPermissionResp.PermissionIds {
			if pid == resp.GetId() {
				found = true
				break
			}
		}
		if !found {
			return nil, adminV1.ErrorForbidden("no access to the permission")
		}
	}

	return resp, nil
}

func (s *PermissionService) Create(ctx context.Context, req *permissionV1.CreatePermissionRequest) (*emptypb.Empty, error) {
	if req.Data == nil {
		return nil, adminV1.ErrorBadRequest("invalid parameter")
	}

	// 获取操作人信息
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.Data.CreatedBy = trans.Ptr(operator.UserId)

	return s.permissionServiceClient.Create(ctx, req)
}

func (s *PermissionService) Update(ctx context.Context, req *permissionV1.UpdatePermissionRequest) (*emptypb.Empty, error) {
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

	return s.permissionServiceClient.Update(ctx, req)
}

func (s *PermissionService) Delete(ctx context.Context, req *permissionV1.DeletePermissionRequest) (*emptypb.Empty, error) {
	return s.permissionServiceClient.Delete(ctx, req)
}

// SyncPermissions 同步权限点
func (s *PermissionService) SyncPermissions(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	// 获取操作人信息
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	// 查询所有启用的菜单
	menus, err := s.menuServiceClient.List(ctx, &paginationV1.PagingRequest{
		NoPaging: trans.Ptr(true),
		FilteringType: &paginationV1.PagingRequest_Query{
			Query: `{"status":"ON"}`,
		},
		OrderBy: trans.Ptr("id desc"),
	})
	if err != nil {
		return nil, err
	}

	s.menuPermissionConverter.ComposeMenuPaths(menus.Items)

	sort.SliceStable(menus.Items, func(i, j int) bool {
		return menus.Items[i].GetParentId() < menus.Items[j].GetParentId()
	})

	var permissionGroups []*permissionV1.PermissionGroup
	var permissions []*permissionV1.Permission
	var mapPermissions = make(map[string][]*permissionV1.Permission)

	permissionGroups = append(permissionGroups, &permissionV1.PermissionGroup{
		Name:      trans.Ptr("未分类"),
		Module:    trans.Ptr(constants.UncategorizedPermissionGroup),
		Status:    trans.Ptr(permissionV1.PermissionGroup_ON),
		SortOrder: trans.Ptr(uint32(len(permissionGroups) + 1)),
		CreatedBy: trans.Ptr(operator.UserId),
		UpdatedBy: trans.Ptr(operator.UserId),
	})

	for _, menu := range menus.Items {
		var title string
		//if menu.GetMeta() != nil && menu.GetMeta().GetTitle() != "" {
		//	title = menu.GetMeta().GetTitle()
		//} else {
		//	title = menu.GetName()
		//}
		title = menu.GetName()

		var permissionCode string
		permissionCode = s.menuPermissionConverter.ConvertCode(menu.GetPath(), title, menu.GetType())
		if permissionCode == "" {
			continue
		}

		module := s.menuPermissionConverter.MenuPathToModuleName(menu.GetPath())

		// 以目录类型的菜单作为权限组
		if menu.GetType() == resourceV1.Menu_CATALOG {
			permissionGroups = append(permissionGroups, &permissionV1.PermissionGroup{
				Name:      trans.Ptr(title),
				Module:    trans.Ptr(module),
				Status:    trans.Ptr(permissionV1.PermissionGroup_ON),
				SortOrder: trans.Ptr(uint32(len(permissionGroups) + 1)),
				CreatedBy: trans.Ptr(operator.UserId),
				UpdatedBy: trans.Ptr(operator.UserId),
			})
			//s.log.Debugf("SyncPermissions: created permission group for menu %s - %s", menu.GetName(), permissionCode)
		}

		perm := &permissionV1.Permission{
			Name:      trans.Ptr(title),
			Code:      trans.Ptr(permissionCode),
			Status:    trans.Ptr(permissionV1.Permission_ON),
			MenuIds:   []uint32{menu.GetId()},
			CreatedBy: trans.Ptr(operator.UserId),
			UpdatedBy: trans.Ptr(operator.UserId),
		}

		permissions = append(permissions, perm)

		mapPermissions[module] = append(mapPermissions[module], perm)
	}

	return s.permissionServiceClient.SyncPermissions(ctx, &permissionV1.SyncPermissionsRequest{
		PermissionGroups: permissionGroups,
		Permissions:      permissions,
		OperatorId:       trans.Ptr(operator.UserId),
	})
}
