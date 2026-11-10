package data

import (
	"context"
	"strings"
	"testing"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/stretchr/testify/require"

	"github.com/tx7do/go-utils/trans"

	entCrud "github.com/tx7do/go-crud/entgo"
	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	"go-wind-cms/app/core/service/internal/data/client"
	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/enttest"
	"go-wind-cms/app/core/service/internal/data/ent/permission"
	"go-wind-cms/app/core/service/internal/data/ent/permissionapi"
	"go-wind-cms/app/core/service/internal/data/ent/permissionmenu"
	"go-wind-cms/app/core/service/internal/data/ent/rolepermission"

	permissionV1 "go-wind-cms/api/gen/go/permission/service/v1"
	"go-wind-cms/pkg/constants"
	appViewer "go-wind-cms/pkg/entgo/viewer"

	_ "github.com/xiaoqidun/entps"
)

// newTestPermissionSyncRepo 构造一套连内存 sqlite 的权限仓储簇（权限同步 MERGE 用）。
// 必须经 enttest 建客户端：真实 import ent/runtime 触发 Hooks 注册，同时自动跑
// schema migration。每个测试用独立命名内存库，避免 cache=shared 跨测试撞车。
func newTestPermissionSyncRepo(t *testing.T) (*ent.Client, *PermissionRepo) {
	t.Helper()

	drv, err := entCrud.CreateDriver(
		"sqlite3",
		"file:ent_"+strings.ReplaceAll(t.Name(), "/", "_")+"?mode=memory&cache=shared&_fk=1&_time_format=sqlite",
		false, false,
	)
	require.NoError(t, err, "创建 sqlite driver 失败")

	db := enttest.NewClient(t, enttest.WithOptions(
		ent.Driver(drv),
		ent.Log(func(a ...any) { t.Log(a...) }),
	))

	client.ApplyTimeDefaultHooks(db)

	wrapped := entCrud.NewEntClient(db, drv)
	bctx := bootstrap.NewContextWithParam(context.Background(), &conf.AppInfo{}, &conf.Bootstrap{}, log.DefaultLogger)

	return db, NewPermissionRepo(bctx, wrapped,
		NewPermissionApiRepo(bctx, wrapped),
		NewPermissionMenuRepo(bctx, wrapped),
		NewRolePermissionRepo(bctx, wrapped),
	)
}

// TestPermissionSyncRepoMergePrimitives 权限同步 MERGE 三原语的语义钉子：
// ListBizPermissions 的业务集过滤、UpdateForSync 的保 id 与 status 不触碰、
// DeleteByIDs 的从属绑定级联清理（接口/菜单/角色授权）。服务层的分段逻辑
// （更新/插入/删除三段）由本原语组合而成，原语语义不漂移则分段不漂移。
func TestPermissionSyncRepoMergePrimitives(t *testing.T) {
	entClient, repo := newTestPermissionSyncRepo(t)
	ctx := appViewer.NewSystemViewerContext(context.Background())

	// —— 播种：业务权限两条（一条稍后更新、一条稍后删除）+ 系统种子权限一条 ——
	pKeep, err := entClient.Permission.Create().
		SetName("旧名称").
		SetCode("opm:user:list").
		Save(ctx)
	require.NoError(t, err)

	pGone, err := entClient.Permission.Create().
		SetName("待删除").
		SetCode("content:post:purge").
		Save(ctx)
	require.NoError(t, err)

	_, err = entClient.Permission.Create().
		SetName("系统种子").
		SetCode(constants.SystemPermissionCodePrefix + "login").
		Save(ctx)
	require.NoError(t, err)

	// 待删除权限的三类从属绑定 + 角色授权行
	require.NoError(t, entClient.PermissionApi.Create().
		SetPermissionID(pGone.ID).SetAPIID(9001).Exec(ctx))
	require.NoError(t, entClient.PermissionMenu.Create().
		SetPermissionID(pGone.ID).SetMenuID(9101).Exec(ctx))
	_, err = entClient.RolePermission.Create().
		SetRoleID(88).SetPermissionID(pGone.ID).Save(ctx)
	require.NoError(t, err)

	// 待更新权限也挂一条绑定，稍后验证 AssignSyncBindings 的整体替换语义
	require.NoError(t, entClient.PermissionApi.Create().
		SetPermissionID(pKeep.ID).SetAPIID(9002).Exec(ctx))

	// —— ListBizPermissions：业务集含两条业务权限、不含 sys: 前缀 ——
	biz, err := repo.ListBizPermissions(ctx)
	require.NoError(t, err)
	codes := make(map[string]uint64, len(biz))
	for _, p := range biz {
		codes[p.GetCode()] = p.GetId()
	}
	require.Len(t, biz, 2, "应只含两条业务权限（sys: 前缀排除）")
	require.Contains(t, codes, "opm:user:list")
	require.Contains(t, codes, "content:post:purge")
	require.NotContains(t, codes, constants.SystemPermissionCodePrefix+"login")

	// —— UpdateForSync：保 id、更新名称/分组、不触碰手工停用的 status ——
	require.NoError(t, entClient.Permission.UpdateOneID(pKeep.ID).SetStatus(permission.StatusOff).Exec(ctx))
	newGroup, err := entClient.PermissionGroup.Create().SetName("用户管理").SetModule("opm").Save(ctx)
	require.NoError(t, err)
	derived := &permissionV1.Permission{
		Name:    trans.Ptr("新名称"),
		GroupId: trans.Ptr(newGroup.ID),
	}
	require.NoError(t, repo.UpdateForSync(ctx, pKeep.ID, derived))

	after, err := entClient.Permission.Get(ctx, pKeep.ID)
	require.NoError(t, err)
	require.NotNil(t, after.Status)
	require.Equal(t, permission.StatusOff, *after.Status, "同步更新不得重新启用手工停用的权限")
	require.NotNil(t, after.Name)
	require.Equal(t, "新名称", *after.Name)
	require.NotNil(t, after.GroupID)
	require.Equal(t, newGroup.ID, *after.GroupID)
	require.Equal(t, pKeep.ID, after.ID, "原位更新必须保 id")

	// —— AssignSyncBindings：绑定按派生集整体替换（旧 9002 → 新 11,12） ——
	require.NoError(t, repo.AssignSyncBindings(ctx, pKeep.ID, []uint64{11, 12}, []uint64{21}))
	apiIDs, err := repo.ListApiIDsByPermissionIDs(ctx, []uint64{pKeep.ID})
	require.NoError(t, err)
	require.ElementsMatch(t, []uint64{11, 12}, apiIDs)
	menuIDs, err := repo.ListMenuIDsByPermissionIDs(ctx, []uint64{pKeep.ID})
	require.NoError(t, err)
	require.ElementsMatch(t, []uint64{21}, menuIDs)
	// 空派生 = 清空该权限全部绑定
	require.NoError(t, repo.AssignSyncBindings(ctx, pKeep.ID, nil, nil))
	apiIDs, err = repo.ListApiIDsByPermissionIDs(ctx, []uint64{pKeep.ID})
	require.NoError(t, err)
	require.Empty(t, apiIDs)

	// —— DeleteByIDs：删行并级联清理接口/菜单/角色授权绑定 ——
	require.NoError(t, repo.DeleteByIDs(ctx, []uint64{pGone.ID}))
	_, err = entClient.Permission.Get(ctx, pGone.ID)
	require.Error(t, err, "权限行应已删除")

	apiCnt, err := entClient.PermissionApi.Query().Where(permissionapi.PermissionIDIn(pGone.ID)).Count(ctx)
	require.NoError(t, err)
	require.Zero(t, apiCnt, "接口绑定应级联清理")

	menuCnt, err := entClient.PermissionMenu.Query().Where(permissionmenu.PermissionIDIn(pGone.ID)).Count(ctx)
	require.NoError(t, err)
	require.Zero(t, menuCnt, "菜单绑定应级联清理")

	rpCnt, err := entClient.RolePermission.Query().Where(rolepermission.PermissionIDIn(pGone.ID)).Count(ctx)
	require.NoError(t, err)
	require.Zero(t, rpCnt, "角色授权应级联清理（此前重建换 id 的悬空授权根源）")

	// 其他权限的行与绑定不受波及
	_, err = entClient.Permission.Get(ctx, pKeep.ID)
	require.NoError(t, err)
}

// TestPermissionSyncRepoDeleteByIDsEmpty 空集直返，不启事务不报错。
func TestPermissionSyncRepoDeleteByIDsEmpty(t *testing.T) {
	_, repo := newTestPermissionSyncRepo(t)
	ctx := appViewer.NewSystemViewerContext(context.Background())
	require.NoError(t, repo.DeleteByIDs(ctx, nil))
}
