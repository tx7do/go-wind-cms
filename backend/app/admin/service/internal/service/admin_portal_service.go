package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/go-utils/sliceutil"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	identityV1 "go-wind-cms/api/gen/go/identity/service/v1"
	permissionV1 "go-wind-cms/api/gen/go/permission/service/v1"
	resourceV1 "go-wind-cms/api/gen/go/resource/service/v1"

	"go-wind-cms/pkg/middleware/auth"
)

type AdminPortalService struct {
	adminV1.AdminPortalServiceHTTPServer

	log *log.Helper

	menuServiceClient       resourceV1.MenuServiceClient
	permissionServiceClient permissionV1.PermissionServiceClient

	roleServiceClient permissionV1.RoleServiceClient
	userServiceClient identityV1.UserServiceClient
}

func NewRouterService(
	ctx *bootstrap.Context,
	menuServiceClient resourceV1.MenuServiceClient,
	permissionServiceClient permissionV1.PermissionServiceClient,
	roleServiceClient permissionV1.RoleServiceClient,
	userServiceClient identityV1.UserServiceClient,
) *AdminPortalService {
	return &AdminPortalService{
		log:                     ctx.NewLoggerHelper("admin-portal/service/admin-service"),
		menuServiceClient:       menuServiceClient,
		permissionServiceClient: permissionServiceClient,
		roleServiceClient:       roleServiceClient,
		userServiceClient:       userServiceClient,
	}
}

func (s *AdminPortalService) menuListToQueryString(menus []uint32, onlyButton bool) string {
	var ids []string
	for _, menu := range menus {
		ids = append(ids, fmt.Sprintf("\"%d\"", menu))
	}
	idsStr := fmt.Sprintf("[%s]", strings.Join(ids, ", "))
	query := map[string]string{"id__in": idsStr}

	if onlyButton {
		query["type"] = resourceV1.Menu_BUTTON.String()
	} else {
		query["type__not"] = resourceV1.Menu_BUTTON.String()
	}

	query["status"] = "ON"

	queryStr, err := json.Marshal(query)
	if err != nil {
		return ""
	}

	return string(queryStr)
}

// queryMultipleRolesMenusByRoleIds 使用RoleIDs查询菜单，即多个角色的菜单
func (s *AdminPortalService) queryMultipleRolesMenusByRoleIds(ctx context.Context, roleIDs []uint32) ([]uint32, error) {
	if len(roleIDs) == 0 {
		s.log.Warnf("user has no roles assigned")
		return nil, nil
	}

	permissionIDs, err := s.roleServiceClient.ListPermissionIds(ctx, &permissionV1.ListPermissionIdsRequest{
		RoleIds: roleIDs,
	})
	if err != nil {
		return nil, adminV1.ErrorInternalServerError("query roles permissions failed")
	}

	resourceResp, err := s.permissionServiceClient.ListPermissionResources(ctx, &permissionV1.ListPermissionResourcesRequest{
		PermissionIds: permissionIDs.PermissionIds,
		ResourceTypes: []permissionV1.ListPermissionResourcesRequest_ResourceType{
			permissionV1.ListPermissionResourcesRequest_MENU,
		},
	})
	if err != nil {
		s.log.Errorf("list permission resources failed [%s]", err.Error())
		return nil, adminV1.ErrorInternalServerError("list permission resources failed")
	}

	idSet, ok := resourceResp.Resources[int32(permissionV1.ListPermissionResourcesRequest_MENU.Number())]
	if !ok {
		s.log.Warnf("roles [%v] has no menu permission", roleIDs)
		return nil, nil
	}
	if len(idSet.Ids) == 0 {
		s.log.Warnf("roles [%v] has no menu permission", roleIDs)
		return nil, nil
	}

	idSet.Ids = sliceutil.Unique(idSet.Ids)

	return idSet.Ids, nil
}

func (s *AdminPortalService) GetMyPermissionCode(ctx context.Context, _ *emptypb.Empty) (*adminV1.ListPermissionCodeResponse, error) {
	// 获取操作人信息
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	permissionIdsResp, err := s.roleServiceClient.ListPermissionIds(ctx, &permissionV1.ListPermissionIdsRequest{
		QueryBy: &permissionV1.ListPermissionIdsRequest_UserId{
			UserId: operator.GetUserId(),
		},
	})
	if err != nil {
		s.log.Errorf("list user role ids failed [%s]", err.Error())
		return nil, adminV1.ErrorInternalServerError("list user role ids failed")
	}

	permissionCodeResp, err := s.permissionServiceClient.ListPermissionCodesByIds(ctx, &permissionV1.ListPermissionCodesByIdsRequest{
		PermissionIds: permissionIdsResp.PermissionIds,
	})
	if err != nil {
		s.log.Errorf("list permission codes by ids failed [%s]", err.Error())
		return nil, adminV1.ErrorInternalServerError("list permission codes by ids failed")
	}

	return &adminV1.ListPermissionCodeResponse{
		Codes: permissionCodeResp.PermissionCodes,
	}, nil
}

func (s *AdminPortalService) fillRouteItem(menus []*resourceV1.Menu) []*resourceV1.MenuRouteItem {
	if len(menus) == 0 {
		return nil
	}

	var routers []*resourceV1.MenuRouteItem

	for _, v := range menus {
		if v.GetStatus() != resourceV1.Menu_ON {
			continue
		}
		if v.GetType() == resourceV1.Menu_BUTTON {
			continue
		}

		item := &resourceV1.MenuRouteItem{
			Path:      v.Path,
			Component: v.Component,
			Name:      v.Name,
			Redirect:  v.Redirect,
			Alias:     v.Alias,
			Meta:      v.Meta,
		}

		if len(v.Children) > 0 {
			item.Children = s.fillRouteItem(v.Children)
		}

		routers = append(routers, item)
	}

	return routers
}

func (s *AdminPortalService) GetNavigation(ctx context.Context, _ *emptypb.Empty) (*adminV1.ListRouteResponse, error) {
	// 获取操作人信息
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	user, err := s.userServiceClient.Get(ctx, &identityV1.GetUserRequest{
		QueryBy: &identityV1.GetUserRequest_Id{
			Id: operator.UserId,
		},
	})
	if err != nil {
		s.log.Errorf("query user failed[%s]", err.Error())
		return nil, adminV1.ErrorInternalServerError("query user failed")
	}

	// 多角色的菜单
	roleMenus, err := s.queryMultipleRolesMenusByRoleIds(ctx, user.GetRoleIds())
	if err != nil {
		return nil, err
	}

	menuList, err := s.menuServiceClient.List(ctx, &paginationV1.PagingRequest{
		NoPaging: trans.Ptr(true),
		FilteringType: &paginationV1.PagingRequest_Query{
			Query: s.menuListToQueryString(roleMenus, false),
		},
	})
	if err != nil {
		s.log.Errorf("list route failed [%s]", err.Error())
		return nil, adminV1.ErrorInternalServerError("list route failed")
	}

	return &adminV1.ListRouteResponse{Items: s.fillRouteItem(menuList.Items)}, nil
}
