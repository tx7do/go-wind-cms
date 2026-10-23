package data

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/go-kratos/kratos/v2/log"
	entCrud "github.com/tx7do/go-crud/entgo"
	"github.com/tx7do/go-crud/viewer"
	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	aiV1 "go-wind-cms/api/gen/go/ai/service/v1"
	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
	identityV1 "go-wind-cms/api/gen/go/identity/service/v1"
	appViewer "go-wind-cms/pkg/entgo/viewer"

	"go-wind-cms/app/core/service/internal/data/client"
	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/aiknowledgebase"
	"go-wind-cms/app/core/service/internal/data/ent/contentmodel"
	"go-wind-cms/app/core/service/internal/data/ent/enttest"
	"go-wind-cms/app/core/service/internal/data/ent/fielddefinition"

	_ "github.com/xiaoqidun/entps"
)

// newTenantIsolationTestEnv 构造一套连内存 sqlite 的 ent 客户端与被测仓储。
// 必须经 enttest 建客户端：真实 import ent/runtime 以挂载 TenantPrivacy 策略
// （mixin.TenantID.Policy），使查询路径的租户注入与生产行为一致。
func newTenantIsolationTestEnv(t *testing.T) (*ent.Client, *AiKnowledgeRepo, *ContentModelRepo) {
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
	return db, NewAiKnowledgeRepo(bctx, wrapped), NewContentModelRepo(bctx, wrapped)
}

// tenantViewerContext 模拟租户用户（tid=tenantID）的请求上下文。
func tenantViewerContext(tenantID uint32) context.Context {
	return viewer.WithContext(context.Background(),
		appViewer.NewUserViewer(uint64(tenantID)+10, uint64(tenantID), 0, "", identityV1.DataScope_DATA_SCOPE_UNSPECIFIED))
}

func systemViewerContext() context.Context {
	return appViewer.NewSystemViewerContext(context.Background())
}

// TestTenantWriteIsolationAiKnowledge 钉住知识库写路径的跨租户隔离：
// Update/Delete/DeleteDoc 是 ent 变更路径，TenantPrivacy 不覆盖，repo 必须显式
// 加租户谓词——否则 WHERE 仅含 id，任意租户可改/删他人知识库（CWE-639）。
func TestTenantWriteIsolationAiKnowledge(t *testing.T) {
	db, kbRepo, _ := newTenantIsolationTestEnv(t)
	sysCtx := systemViewerContext()

	// 播种：租户 1 / 租户 2 各一个知识库，租户 2 的库下挂一个文档
	kb1, err := db.AiKnowledgeBase.Create().SetName("kb-t1").SetTenantID(1).Save(sysCtx)
	require.NoError(t, err)
	kb2, err := db.AiKnowledgeBase.Create().SetName("kb-t2").SetTenantID(2).Save(sysCtx)
	require.NoError(t, err)
	doc2, err := db.AiDoc.Create().SetBaseID(kb2.ID).SetName("doc-t2").SetTenantID(2).Save(sysCtx)
	require.NoError(t, err)

	ctx := tenantViewerContext(1)

	// 跨租户 Update：repo 层谓词保证他人数据不被改动
	updErr := kbRepo.Update(ctx, &aiV1.UpdateAiKnowledgeBaseRequest{
		Id:   kb2.ID,
		Data: &aiV1.AiKnowledgeBase{Name: ptrString("PWNED")},
	})
	require.NoError(t, updErr, "repo 层静默 no-op（归属语义由 service 预检兜住），但数据必须不被改动")
	after, err := db.AiKnowledgeBase.Query().Where(aiknowledgebase.IDEQ(kb2.ID)).Only(sysCtx)
	require.NoError(t, err)
	require.Equal(t, "kb-t2", *after.Name, "跨租户 Update 不得改动他人知识库名称")

	// 跨租户 Delete / DeleteDoc：必须 NotFound 且行原样保留
	require.Error(t, kbRepo.Delete(ctx, kb2.ID), "跨租户 Delete 必须 NotFound")
	require.Error(t, kbRepo.DeleteDoc(ctx, doc2.ID), "跨租户 DeleteDoc 必须 NotFound")
	got2, err := db.AiKnowledgeBase.Query().Where(aiknowledgebase.IDEQ(kb2.ID)).Only(sysCtx)
	require.NoError(t, err, "跨租户 Delete 后他人知识库必须仍存在")
	require.Equal(t, "kb-t2", *got2.Name)
	gotDoc, err := db.AiDoc.Get(sysCtx, doc2.ID)
	require.NoError(t, err, "跨租户 DeleteDoc 后他人文档必须仍存在")
	require.Equal(t, "doc-t2", *gotDoc.Name)

	// 归属读侧（service 预检依赖）：租户 1 视角查他人库必须 NotFound
	_, err = kbRepo.GetEntityByID(ctx, kb2.ID)
	require.Error(t, err, "跨租户 GetEntityByID 必须 NotFound（privacy 查询注入）")
	_, err = kbRepo.GetDocByID(ctx, doc2.ID)
	require.Error(t, err, "跨租户 GetDocByID 必须 NotFound（privacy 查询注入）")

	// 正向对照：本租户的更新/删除不受影响
	require.NoError(t, kbRepo.Update(ctx, &aiV1.UpdateAiKnowledgeBaseRequest{
		Id:   kb1.ID,
		Data: &aiV1.AiKnowledgeBase{Name: ptrString("kb-t1-renamed")},
	}))
	after1, err := db.AiKnowledgeBase.Query().Where(aiknowledgebase.IDEQ(kb1.ID)).Only(sysCtx)
	require.NoError(t, err)
	require.Equal(t, "kb-t1-renamed", *after1.Name)

	doc1, err := db.AiDoc.Create().SetBaseID(kb1.ID).SetName("doc-t1").SetTenantID(1).Save(sysCtx)
	require.NoError(t, err)
	require.NoError(t, kbRepo.DeleteDoc(ctx, doc1.ID), "本租户 DeleteDoc 必须成功")
	_, err = db.AiDoc.Get(sysCtx, doc1.ID)
	require.Error(t, err, "本租户文档应已被删除")
}

// TestTenantWriteIsolationContentModel 钉住内容模型写路径的跨租户隔离：
// Update/Delete（级联字段定义与翻译）必须先过 id+租户作用域检查。
func TestTenantWriteIsolationContentModel(t *testing.T) {
	db, _, cmRepo := newTenantIsolationTestEnv(t)
	sysCtx := systemViewerContext()

	m1, err := db.ContentModel.Create().SetName("model-t1").SetCode("t1_model").SetTenantID(1).Save(sysCtx)
	require.NoError(t, err)
	m2, err := db.ContentModel.Create().SetName("model-t2").SetCode("t2_model").SetTenantID(2).Save(sysCtx)
	require.NoError(t, err)
	_, err = db.FieldDefinition.Create().SetContentModelID(m2.ID).SetName("f2").SetTenantID(2).Save(sysCtx)
	require.NoError(t, err)

	ctx := tenantViewerContext(1)

	// 跨租户 Update：必须 NotFound 且他人模型不被改动
	_, err = cmRepo.Update(ctx, &contentV1.UpdateContentModelRequest{
		Id:   m2.ID,
		Data: &contentV1.ContentModel{Name: ptrString("PWNED")},
	})
	require.Error(t, err, "跨租户 Update 必须 NotFound")
	after, err := db.ContentModel.Query().Where(contentmodel.IDEQ(m2.ID)).Only(sysCtx)
	require.NoError(t, err)
	require.Equal(t, "model-t2", *after.Name, "跨租户 Update 不得改动他人模型")

	// 跨租户 Delete：必须 NotFound，模型与级联字段定义原样保留
	require.Error(t, cmRepo.Delete(ctx, &contentV1.DeleteContentModelRequest{
		QueryBy: &contentV1.DeleteContentModelRequest_Id{Id: m2.ID},
	}), "跨租户 Delete 必须 NotFound")
	_, err = db.ContentModel.Query().Where(contentmodel.IDEQ(m2.ID)).Only(sysCtx)
	require.NoError(t, err, "跨租户 Delete 后他人模型必须仍存在")
	fdCount, err := db.FieldDefinition.Query().Where(fielddefinition.ContentModelIDEQ(m2.ID)).Count(sysCtx)
	require.NoError(t, err)
	require.Equal(t, 1, fdCount, "级联字段定义必须原样保留")

	// 正向对照：本租户的更新/删除（含级联）不受影响
	_, err = cmRepo.Update(ctx, &contentV1.UpdateContentModelRequest{
		Id:   m1.ID,
		Data: &contentV1.ContentModel{Name: ptrString("model-t1-renamed")},
	})
	require.NoError(t, err)
	after1, err := db.ContentModel.Query().Where(contentmodel.IDEQ(m1.ID)).Only(sysCtx)
	require.NoError(t, err)
	require.Equal(t, "model-t1-renamed", *after1.Name)

	_, err = db.FieldDefinition.Create().SetContentModelID(m1.ID).SetName("f1").SetTenantID(1).Save(sysCtx)
	require.NoError(t, err)
	require.NoError(t, cmRepo.Delete(ctx, &contentV1.DeleteContentModelRequest{
		QueryBy: &contentV1.DeleteContentModelRequest_Id{Id: m1.ID},
	}), "本租户 Delete 必须成功")
	_, err = db.ContentModel.Query().Where(contentmodel.IDEQ(m1.ID)).Only(sysCtx)
	require.Error(t, err, "本租户模型应已被删除")
	fdCount, err = db.FieldDefinition.Query().Where(fielddefinition.ContentModelIDEQ(m1.ID)).Count(sysCtx)
	require.NoError(t, err)
	require.Equal(t, 0, fdCount, "本租户级联字段定义应已被删除")
}

func ptrString(s string) *string { return &s }
