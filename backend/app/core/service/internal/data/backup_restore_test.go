package data

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/go-kratos/kratos/v2/log"
	entCrud "github.com/tx7do/go-crud/entgo"
	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	appViewer "go-wind-cms/pkg/entgo/viewer"

	"go-wind-cms/app/core/service/internal/data/client"
	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/enttest"
	"go-wind-cms/app/core/service/internal/data/ent/menu"
	"go-wind-cms/app/core/service/internal/data/ent/tenant"

	_ "github.com/xiaoqidun/entps"
)

// newTestBackupRepo 构造一套连内存 sqlite 的备份仓储。
// 必须经 enttest 建客户端：enttest 真实 import ent/runtime，触发其 init() 填充各实体
// Hooks[0]，否则任何 mutation 都会报 "ent: uninitialized hook (forgotten import ent/runtime?)"。
// enttest.NewClient 同时自动跑 schema migration。
func newTestBackupRepo(t *testing.T) (*ent.Client, *BackupRepo) {
	t.Helper()

	// 每个测试用独立命名的内存库：cache=shared 下同名 DSN 在同进程内共享，
	// 数据会跨测试撞车。
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

	return db, NewBackupRepo(bctx, wrapped)
}

// newBackupTestSystemContext 系统旁路上下文：导出/恢复要跨租户全量读写。
func newBackupTestSystemContext() context.Context {
	return appViewer.NewSystemViewerContext(context.Background())
}

// TestBackupRestoreRoundTrip 恢复工具链的往返实证：
// 导出核心表 → 清空全部目标表 → RestoreCoreTables 恢复 → 行数与内容一致。
// 非空表前置校验与单事务语义由其余用例钉住。
func TestBackupRestoreRoundTrip(t *testing.T) {
	entClient, repo := newTestBackupRepo(t)
	client_ := entClient
	ctx := newBackupTestSystemContext()

	// 1. 播种：两个租户（带 Logo/时间字段）+ 一条菜单（父子结构）
	t1, err := client_.Tenant.Create().
		SetName("往返租户甲").
		SetCode("RT_A").
		SetStatus(tenant.StatusOn).
		Save(ctx)
	require.NoError(t, err)
	t2, err := client_.Tenant.Create().
		SetName("往返租户乙").
		SetCode("RT_B").
		SetStatus(tenant.StatusOn).
		Save(ctx)
	require.NoError(t, err)

	menu1, err := client_.Menu.Create().
		SetName("往返菜单").
		SetPath("/rt").
		SetStatus(menu.StatusOn).
		Save(ctx)
	require.NoError(t, err)

	// 2. 导出：data 段即 map[表名]行数组（RestoreCoreTables 的入参形状）
	exported, err := repo.ExportCoreTables(ctx)
	require.NoError(t, err)
	rawTables, err := json.Marshal(exported)
	require.NoError(t, err)

	// 3. 清空全部目标表（模拟"全新库"）
	_, err = client_.Menu.Delete().Exec(ctx)
	require.NoError(t, err)
	_, err = client_.Tenant.Delete().Exec(ctx)
	require.NoError(t, err)

	// 中途抽检：确认真空了
	count, err := client_.Tenant.Query().Count(ctx)
	require.NoError(t, err)
	require.Zero(t, count)

	// 4. 恢复
	counts, err := repo.RestoreCoreTables(ctx, rawTables)
	require.NoError(t, err)
	require.EqualValues(t, 2, counts["tenants"], "租户应恢复 2 行")
	require.EqualValues(t, 1, counts["menus"], "菜单应恢复 1 行")

	// 5. 内容一致性抽检：字段值与 ID 保真
	restored1, err := client_.Tenant.Query().
		Where(tenant.IDEQ(t1.ID)).
		Only(ctx)
	require.NoError(t, err)
	require.Equal(t, "往返租户甲", *restored1.Name)
	require.Equal(t, "RT_A", *restored1.Code)

	restoredMenu, err := client_.Menu.Query().
		Where(menu.IDEQ(menu1.ID)).
		Only(ctx)
	require.NoError(t, err)
	require.Equal(t, "往返菜单", *restoredMenu.Name)
	require.Equal(t, menu1.ID, restoredMenu.ID)
	_ = t2
}

// TestRestoreRefusesNonEmptyTables 非空表前置校验：任一目标表非空即中止且不写入。
func TestRestoreRefusesNonEmptyTables(t *testing.T) {
	entClient, repo := newTestBackupRepo(t)
	client_ := entClient
	ctx := newBackupTestSystemContext()

	// 造一个非空的 menus（目标表之一），其余全空
	_, err := client_.Menu.Create().
		SetName("既有菜单").
		SetPath("/existing").
		SetStatus(menu.StatusOn).
		Save(ctx)
	require.NoError(t, err)

	t1, err := client_.Tenant.Create().
		SetName("存活租户").
		SetCode("KEEP").
		SetStatus(tenant.StatusOn).
		Save(ctx)
	require.NoError(t, err)

	// 造一份合法备份 payload（tenants 有 1 行）
	exported, err := repo.ExportCoreTables(ctx)
	require.NoError(t, err)
	// 从导出里删掉 menus 键，模拟"备份里 menus 为空、但库里 menus 非空"的最小非空场景
	delete(exported, "menus")
	rawTables, err := json.Marshal(exported)
	require.NoError(t, err)

	// 恢复应被拒绝：menus 非空
	_, err = repo.RestoreCoreTables(ctx, rawTables)
	require.Error(t, err)
	require.Contains(t, err.Error(), "not empty")

	// 中止语义：零写入——tenants 里应只有先前那行，备份中的租户没进来
	count, err := client_.Tenant.Query().Count(ctx)
	require.NoError(t, err)
	require.Equal(t, 1, count)

	got, err := client_.Tenant.Query().Where(tenant.IDEQ(t1.ID)).Only(ctx)
	require.NoError(t, err)
	require.Equal(t, "KEEP", *got.Code)
}

// TestRestoreUnknownTableIgnored 备份 JSON 里的未知键（如 _meta 误入 data 段）被忽略。
func TestRestoreUnknownTableIgnored(t *testing.T) {
	_, repo := newTestBackupRepo(t)
	ctx := newBackupTestSystemContext()

	payload := map[string]json.RawMessage{
		"tenants": json.RawMessage(`[]`),
		"bogus":   json.RawMessage(`[{"id":1}]`),
	}
	raw, err := json.Marshal(payload)
	require.NoError(t, err)

	counts, err := repo.RestoreCoreTables(ctx, raw)
	require.NoError(t, err)
	require.EqualValues(t, 0, counts["tenants"])
	_, restored := counts["bogus"]
	require.False(t, restored, "未知表不应出现在计数里")
}
