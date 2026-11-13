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

	identityV1 "go-wind-cms/api/gen/go/identity/service/v1"
	permissionV1 "go-wind-cms/api/gen/go/permission/service/v1"

	converter "github.com/tx7do/go-utils/converter"
	"go-wind-cms/pkg/constants"
	appViewer "go-wind-cms/pkg/entgo/viewer"
	menuconverter "go-wind-cms/pkg/utils/converter"
)

type PermissionService struct {
	permissionV1.UnimplementedPermissionServiceServer

	log *log.Helper

	permissionRepo      *data.PermissionRepo
	permissionGroupRepo *data.PermissionGroupRepo

	menuRepo *data.MenuRepo
	apiRepo  *data.ApiRepo

	roleRepo *data.RoleRepo

	menuPermissionConverter *menuconverter.MenuPermissionConverter
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
		menuPermissionConverter: menuconverter.NewMenuPermissionConverter(),
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
	groupSet aggregator.ResourceMap[int64, *permissionV1.PermissionGroup],
) {
	for _, p := range permissions {
		if p.GetGroupId() > 0 {
			groupSet[p.GetGroupId()] = nil
		}
	}
}

func (s *PermissionService) fetchRelationInfo(
	ctx context.Context,
	groupSet aggregator.ResourceMap[int64, *permissionV1.PermissionGroup],
) error {
	if len(groupSet) > 0 {
		groupIds := make([]int64, 0, len(groupSet))
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
	groupSet aggregator.ResourceMap[int64, *permissionV1.PermissionGroup],
) {
	aggregator.Populate(
		permissions,
		groupSet,
		func(ou *permissionV1.Permission) int64 { return ou.GetGroupId() },
		func(ou *permissionV1.Permission, g *permissionV1.PermissionGroup) {
			ou.GroupName = g.Name
		},
	)
}

func (s *PermissionService) enrichRelations(ctx context.Context, permissions []*permissionV1.Permission) error {
	var groupSet = make(aggregator.ResourceMap[int64, *permissionV1.PermissionGroup])
	s.extractRelationIDs(permissions, groupSet)
	if err := s.fetchRelationInfo(ctx, groupSet); err != nil {
		return err
	}
	s.bindRelations(permissions, groupSet)
	return nil
}

func (s *PermissionService) List(ctx context.Context, req *paginationV1.PagingRequest) (*permissionV1.ListPermissionResponse, error) {
	resp, err := s.permissionRepo.List(ctx, req, nil)
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
	resp, err := s.permissionRepo.Get(ctx, req)
	if err != nil {
		return nil, err
	}

	fakeItems := []*permissionV1.Permission{resp}
	_ = s.enrichRelations(ctx, fakeItems)

	return resp, nil
}

func (s *PermissionService) Create(ctx context.Context, req *permissionV1.CreatePermissionRequest) (*emptypb.Empty, error) {
	if constants.IsReservedPermissionCode(req.GetData().GetCode()) {
		return nil, identityV1.ErrorBadRequest("permission code is reserved")
	}

	if req == nil || req.Data == nil {
		return nil, permissionV1.ErrorBadRequest("invalid parameter")
	}

	if err := s.permissionRepo.Create(ctx, req); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *PermissionService) Update(ctx context.Context, req *permissionV1.UpdatePermissionRequest) (*emptypb.Empty, error) {
	if constants.IsReservedPermissionCode(req.GetData().GetCode()) {
		return nil, identityV1.ErrorBadRequest("permission code is reserved")
	}

	if req == nil || req.Data == nil {
		return nil, permissionV1.ErrorBadRequest("invalid parameter")
	}

	if err := s.permissionRepo.Update(ctx, req); err != nil {
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
	operatorUserId int64,
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
		apis   []int64
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
				apis:   []int64{api.GetId()},
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

// SyncPermissions 同步权限点。
//
// 按自然键 code 的 **MERGE 增量对齐**（与接口同步对账式 upsert 同一教义，仓内
// 不再存在删除重建型同步）：派生集（BFF 菜单派生 + 本服务 appendAPis 补 API
// 权限）里 code 已存在的权限行原位更新（**保 id**——sys_role_permissions 角色
// 授权、接口/菜单绑定因此免于漂移，此前 Truncate+BatchCreate 全量重建会让角色
// 授权整体悬空），缺失才插入；库里有、派生集没有的业务权限删除并级联清理其
// 角色授权与接口/菜单绑定。管理端手工停用的权限 status 不被同步触碰。
// 分组行无授权绑定引用，保持原重建语义（TruncateBizGroup + 重插）。
func (s *PermissionService) SyncPermissions(ctx context.Context, req *permissionV1.SyncPermissionsRequest) (*emptypb.Empty, error) {
	if req == nil {
		return nil, permissionV1.ErrorBadRequest("invalid request")
	}

	// 同步是系统级全量操作，数据读写统一走系统上下文：管理端用户 viewer 会把
	// 基线查询（ListBizPermissions 等）按租户过滤——基线恒空会让派生集全判
	// "新增"而撞 code 唯一约束。
	ctx = appViewer.NewSystemViewerContext(ctx)

	// 为权限追加对应的 API 资源 ID 列表
	// 操作人身份从 viewer context 推导，忽略客户端 operator_id，防止伪造审计归属
	var operatorUserID int64
	if opID, hasOp := viewerUserIDFromContext(ctx); hasOp {
		operatorUserID = opID
	}
	var mapPermissions = make(map[string][]*permissionV1.Permission)
	if err := s.appendAPis(ctx, &req.Permissions, &mapPermissions, operatorUserID); err != nil {
		s.log.Errorf("append apis to permissions failed: %s", err.Error())
		return nil, err
	}

	var err error

	// 分组行无授权绑定引用：保持重建语义（先重建分组，派生权限挂新组 id）
	if err = s.permissionGroupRepo.TruncateBizGroup(ctx); err != nil {
		s.log.Errorf("truncate biz permission groups failed: %s", err.Error())
		return nil, err
	}
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

	// —— MERGE 分段：按 code 对齐派生集与库存 ——
	derivedByCode := make(map[string]*permissionV1.Permission, len(req.Permissions))
	for _, p := range req.Permissions {
		derivedByCode[p.GetCode()] = p
	}

	var existing []*permissionV1.Permission
	if existing, err = s.permissionRepo.ListBizPermissions(ctx); err != nil {
		return nil, err
	}
	existingByCode := make(map[string]*permissionV1.Permission, len(existing))
	for _, e := range existing {
		existingByCode[e.GetCode()] = e
	}

	// 更新段：code 已存在 → 原位更新（保 id）+ 按派生集整体替换接口/菜单绑定
	// mergeHit 记录已消化码：插入段与删除段都以它判定
	updated, added, removed := 0, 0, 0
	mergeHit := make(map[string]struct{}, len(derivedByCode))
	for code, d := range derivedByCode {
		exist, ok := existingByCode[code]
		if !ok {
			continue
		}
		if err = s.permissionRepo.UpdateForSync(ctx, exist.GetId(), d); err != nil {
			s.log.Errorf("sync permissions update [%s] failed: %s", code, err.Error())
			return nil, err
		}
		if err = s.permissionRepo.AssignSyncBindings(ctx, exist.GetId(), d.GetApiIds(), d.GetMenuIds()); err != nil {
			return nil, err
		}
		updated++
		mergeHit[code] = struct{}{}
	}

	// 插入段：派生集有、库没有
	toInsert := make([]*permissionV1.Permission, 0, len(derivedByCode))
	for _, d := range derivedByCode {
		if _, hit := mergeHit[d.GetCode()]; hit {
			continue
		}
		toInsert = append(toInsert, d)
		added++
	}
	if len(toInsert) > 0 {
		if err = s.permissionRepo.BatchCreate(ctx, toInsert); err != nil {
			s.log.Errorf("batch create permissions failed: %s", err.Error())
			return nil, err
		}
	}

	// 删除段：库有、派生集没有 → 删行并级联清角色授权/接口/菜单绑定。放最后：
	// 中途失败残留的是"本该删的旧码"，重跑同步即自愈，好过删多了。
	goneIDs := make([]int64, 0)
	for code, e := range existingByCode {
		if _, hit := mergeHit[code]; hit {
			continue
		}
		goneIDs = append(goneIDs, e.GetId())
		removed++
	}
	if len(goneIDs) > 0 {
		if err = s.permissionRepo.DeleteByIDs(ctx, goneIDs); err != nil {
			return nil, err
		}
	}
	s.log.Infof("权限同步完成：新增 %d、更新 %d、移除 %d", added, updated, removed)

	return &emptypb.Empty{}, nil
}

func (s *PermissionService) ListPermissionResources(ctx context.Context, req *permissionV1.ListPermissionResourcesRequest) (*permissionV1.ListPermissionResourcesResponse, error) {
	if req == nil {
		return nil, permissionV1.ErrorBadRequest("invalid request")
	}
	if len(req.PermissionIds) == 0 && len(req.RoleIds) == 0 {
		return nil, permissionV1.ErrorBadRequest("permission_ids and role_ids cannot be both empty")
	}

	if len(req.RoleIds) > 0 {
		var limitPermissionIDs []int64
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
