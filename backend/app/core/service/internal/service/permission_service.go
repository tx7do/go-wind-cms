package service

import (
	"context"
	"sort"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/tx7do/go-utils/aggregator"
	"github.com/tx7do/go-utils/sliceutil"
	"github.com/tx7do/go-utils/stringcase"
	"github.com/tx7do/go-utils/trans"

	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	"go-wind-cms/app/core/service/internal/data"

	permissionV1 "go-wind-cms/api/gen/go/permission/service/v1"

	"go-wind-cms/pkg/constants"
	appViewer "go-wind-cms/pkg/entgo/viewer"
	"go-wind-cms/pkg/middleware/auth"
	"go-wind-cms/pkg/utils/converter"
)

type PermissionService struct {
	permissionV1.UnimplementedPermissionServiceServer

	log *log.Helper

	permissionRepo      *data.PermissionRepo
	permissionGroupRepo *data.PermissionGroupRepo

	menuRepo *data.MenuRepo
	apiRepo  *data.ApiRepo

	roleRepo *data.RoleRepo

	menuPermissionConverter *converter.MenuPermissionConverter
	apiPermissionConverter  *converter.ApiPermissionConverter
}

func NewPermissionService(
	ctx *bootstrap.Context,
	permissionRepo *data.PermissionRepo,
	permissionGroupRepo *data.PermissionGroupRepo,
	menuRepo *data.MenuRepo,
	apiRepo *data.ApiRepo,
	roleRepo *data.RoleRepo,
) *PermissionService {
	svc := &PermissionService{
		log:                     ctx.NewLoggerHelper("permission/service/core-service"),
		permissionRepo:          permissionRepo,
		permissionGroupRepo:     permissionGroupRepo,
		menuRepo:                menuRepo,
		apiRepo:                 apiRepo,
		roleRepo:                roleRepo,
		menuPermissionConverter: converter.NewMenuPermissionConverter(),
		apiPermissionConverter:  converter.NewApiPermissionConverter(),
	}

	svc.init()

	return svc
}

func (s *PermissionService) init() {
	ctx := appViewer.NewSystemViewerContext(context.Background())
	if resp, _ := s.permissionRepo.Count(ctx, nil); resp != nil && resp.Count == 0 {
		_ = s.createDefaultPermissions(ctx)
	}
}

func (s *PermissionService) extractRelationIDs(
	permissions []*permissionV1.Permission,
	groupSet aggregator.ResourceMap[uint32, *permissionV1.PermissionGroup],
) {
	for _, p := range permissions {
		if p.GetGroupId() > 0 {
			groupSet[p.GetGroupId()] = nil
		}
	}
}

func (s *PermissionService) fetchRelationInfo(
	ctx context.Context,
	groupSet aggregator.ResourceMap[uint32, *permissionV1.PermissionGroup],
) error {
	if len(groupSet) > 0 {
		groupIds := make([]uint32, 0, len(groupSet))
		for id := range groupSet {
			groupIds = append(groupIds, id)
		}

		groups, err := s.permissionGroupRepo.ListByIDs(ctx, groupIds)
		if err != nil {
			s.log.Errorf("query permission group err: %v", err)
			return err
		}

		for _, g := range groups {
			groupSet[g.GetId()] = g
		}
	}

	return nil
}

func (s *PermissionService) bindRelations(
	permissions []*permissionV1.Permission,
	groupSet aggregator.ResourceMap[uint32, *permissionV1.PermissionGroup],
) {
	aggregator.Populate(
		permissions,
		groupSet,
		func(ou *permissionV1.Permission) uint32 { return ou.GetGroupId() },
		func(ou *permissionV1.Permission, g *permissionV1.PermissionGroup) {
			ou.GroupName = g.Name
		},
	)
}

func (s *PermissionService) enrichRelations(ctx context.Context, permissions []*permissionV1.Permission) error {
	var groupSet = make(aggregator.ResourceMap[uint32, *permissionV1.PermissionGroup])
	s.extractRelationIDs(permissions, groupSet)
	if err := s.fetchRelationInfo(ctx, groupSet); err != nil {
		return err
	}
	s.bindRelations(permissions, groupSet)
	return nil
}

func (s *PermissionService) List(ctx context.Context, req *paginationV1.PagingRequest) (*permissionV1.ListPermissionResponse, error) {
	// 获取操作人信息
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	var limitPermissionIDs []uint32
	if operator.GetTenantId() > 0 {
		limitPermissionIDs, err = s.roleRepo.ListPermissionIDsByRoleCodes(ctx, operator.GetRoles())
		if err != nil {
			return nil, err
		}

		// 没有任何 permission 可访问，直接返回空列表
		if len(limitPermissionIDs) == 0 {
			return &permissionV1.ListPermissionResponse{
				Items: []*permissionV1.Permission{},
				Total: 0,
			}, nil
		}
	}

	resp, err := s.permissionRepo.List(ctx, req, limitPermissionIDs)
	if err != nil {
		return nil, err
	}

	_ = s.enrichRelations(ctx, resp.Items)

	return resp, nil
}

func (s *PermissionService) Count(ctx context.Context, req *paginationV1.PagingRequest) (*permissionV1.CountPermissionResponse, error) {
	return s.permissionRepo.Count(ctx, req)
}

func (s *PermissionService) Get(ctx context.Context, req *permissionV1.GetPermissionRequest) (*permissionV1.Permission, error) {
	// 获取操作人信息
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := s.permissionRepo.Get(ctx, req)
	if err != nil {
		return nil, err
	}

	if operator.GetTenantId() > 0 {
		permissionIDs, err := s.roleRepo.ListPermissionIDsByRoleCodes(ctx, operator.GetRoles())
		if err != nil {
			return nil, err
		}

		found := false
		for _, pid := range permissionIDs {
			if pid == resp.GetId() {
				found = true
				break
			}
		}
		if !found {
			return nil, permissionV1.ErrorForbidden("no access to the permission")
		}
	}

	fakeItems := []*permissionV1.Permission{resp}
	_ = s.enrichRelations(ctx, fakeItems)

	return resp, nil
}

func (s *PermissionService) Create(ctx context.Context, req *permissionV1.CreatePermissionRequest) (*emptypb.Empty, error) {
	if req.Data == nil {
		return nil, permissionV1.ErrorBadRequest("invalid parameter")
	}

	// 获取操作人信息
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.Data.CreatedBy = trans.Ptr(operator.UserId)

	if err = s.permissionRepo.Create(ctx, req); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *PermissionService) Update(ctx context.Context, req *permissionV1.UpdatePermissionRequest) (*emptypb.Empty, error) {
	if req.Data == nil {
		return nil, permissionV1.ErrorBadRequest("invalid parameter")
	}

	// 获取操作人信息
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.Data.UpdatedBy = trans.Ptr(operator.UserId)
	if req.UpdateMask != nil {
		req.UpdateMask.Paths = append(req.UpdateMask.Paths, "updated_by")
	}

	if err = s.permissionRepo.Update(ctx, req); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *PermissionService) Delete(ctx context.Context, req *permissionV1.DeletePermissionRequest) (*emptypb.Empty, error) {
	if err := s.permissionRepo.Delete(ctx, req); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// appendAPis 为权限追加对应的 API 资源 ID 列表
func (s *PermissionService) appendAPis(
	ctx context.Context,
	permissions *[]*permissionV1.Permission,
	mapPermissions *map[string][]*permissionV1.Permission,
	operatorUserId uint32,
) error {
	// 查询所有启用的 API 资源
	apis, err := s.apiRepo.List(ctx, &paginationV1.PagingRequest{
		NoPaging: trans.Ptr(true),
		FilteringType: &paginationV1.PagingRequest_Query{
			Query: `{"status":"ON"}`,
		},
		OrderBy: trans.Ptr("operation"),
	})
	if err != nil {
		return err
	}

	sort.SliceStable(apis.Items, func(i, j int) bool {
		a, b := apis.Items[i], apis.Items[j]
		if a.GetModule() != b.GetModule() {
			return a.GetModule() < b.GetModule()
		}
		if a.GetPath() != b.GetPath() {
			return a.GetPath() > b.GetPath()
		}
		return a.GetOperation() > b.GetOperation()
	})

	type moduleApis struct {
		module string
		apis   []uint32
	}

	codes := make(map[string]*moduleApis)
	for _, api := range apis.Items {
		//code := s.apiPermissionConverter.ConvertCodeByOperationID(api.GetOperation())
		code := s.apiPermissionConverter.ConvertCodeByPath(api.GetMethod(), api.GetPath())
		if code == "" {
			continue
		}

		s.log.Debugf("appendAPis: processing api [%s] [%s] with code [%s]", api.GetMethod(), api.GetPath(), code)

		if curCode, exist := codes[code]; !exist {
			var module string
			for k, perms := range *mapPermissions {
				if len(perms) == 0 {
					continue
				}

				for _, perm := range perms {
					code1Prefix := strings.Split(perm.GetCode(), ":")[0]
					code2Prefix := strings.Split(code, ":")[0]
					if strings.HasPrefix(code2Prefix, code1Prefix) {
						module = k
						break
					}
				}
			}

			if module == "" {
				module = constants.UncategorizedPermissionGroup
			}

			codes[code] = &moduleApis{
				module: module,
				apis:   []uint32{api.GetId()},
			}
		} else {
			curCode.apis = append(curCode.apis, api.GetId())
		}
	}

	for _, perm := range *permissions {
		permIds, exist := codes[perm.GetCode()]
		if exist {
			perm.ApiIds = append(perm.ApiIds, permIds.apis...)
			delete(codes, perm.GetCode())
		}
	}

	//s.log.Debugf("appendAPis: unmatched permission codes: %v", codes)

	for code, apiIDs := range codes {
		name := strings.ReplaceAll(code, ":", "_")
		name = stringcase.ToPascalCase(name)

		perm := &permissionV1.Permission{
			Name:      trans.Ptr(name),
			Code:      trans.Ptr(code),
			Status:    trans.Ptr(permissionV1.Permission_ON),
			ApiIds:    apiIDs.apis,
			CreatedBy: trans.Ptr(operatorUserId),
			UpdatedBy: trans.Ptr(operatorUserId),
		}

		*permissions = append(*permissions, perm)

		(*mapPermissions)[apiIDs.module] = append((*mapPermissions)[apiIDs.module], perm)

		//s.log.Debugf("appendAPis: create permission for api code: [%s][%s]", apiIDs.module, code)
	}

	return nil
}

// SyncPermissions 同步权限点
func (s *PermissionService) SyncPermissions(ctx context.Context, req *permissionV1.SyncPermissionsRequest) (*emptypb.Empty, error) {

	// 清理菜单相关权限
	_ = s.permissionRepo.TruncateBizPermissions(ctx)
	_ = s.permissionGroupRepo.TruncateBizGroup(ctx)

	// 为权限追加对应的 API 资源 ID 列表
	var mapPermissions = make(map[string][]*permissionV1.Permission)
	_ = s.appendAPis(ctx, &req.Permissions, &mapPermissions, req.GetOperatorId())

	var err error

	var finalPermissionGroups []*permissionV1.PermissionGroup
	if finalPermissionGroups, err = s.permissionGroupRepo.BatchCreate(ctx, req.PermissionGroups); err != nil {
		s.log.Errorf("batch create permission groups failed: %s", err.Error())
		return nil, err
	}

	// 为权限分配权限组 ID
	for _, pg := range finalPermissionGroups {
		curPers := mapPermissions[pg.GetModule()]
		for _, p := range curPers {
			p.GroupId = pg.Id
		}
	}

	if err = s.permissionRepo.BatchCreate(ctx, req.Permissions); err != nil {
		s.log.Errorf("batch create permissions failed: %s", err.Error())
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *PermissionService) ListPermissionResources(ctx context.Context, req *permissionV1.ListPermissionResourcesRequest) (*permissionV1.ListPermissionResourcesResponse, error) {
	if len(req.PermissionIds) == 0 && len(req.RoleIds) == 0 {
		return nil, permissionV1.ErrorBadRequest("permission_ids and role_ids cannot be both empty")
	}

	if len(req.RoleIds) > 0 {
		var limitPermissionIDs []uint32
		limitPermissionIDs, err := s.roleRepo.ListPermissionIDsByRoleIDs(ctx, req.GetRoleIds())
		if err != nil {
			return nil, err
		}

		req.PermissionIds = append(req.PermissionIds, limitPermissionIDs...)
		req.PermissionIds = sliceutil.Unique(req.PermissionIds)
	}

	return s.permissionRepo.ListPermissionResources(ctx, req)
}

func (s *PermissionService) ListPermissionCodesByIds(ctx context.Context, req *permissionV1.ListPermissionCodesByIdsRequest) (*permissionV1.ListPermissionCodesByIdsResponse, error) {
	permissionCodes, err := s.permissionRepo.ListPermissionCodesByIds(ctx, req.GetPermissionIds())
	if err != nil {
		return nil, err
	}

	return &permissionV1.ListPermissionCodesByIdsResponse{
		PermissionCodes: permissionCodes,
	}, nil
}

func (s *PermissionService) createDefaultPermissions(ctx context.Context) error {
	var err error

	for _, d := range constants.DefaultPermissions {
		if err = s.permissionRepo.Create(ctx, &permissionV1.CreatePermissionRequest{
			Data: d,
		}); err != nil {
			s.log.Errorf("create default permission %s failed: %v", d.GetCode(), err)
			return err
		}
	}

	return nil
}
