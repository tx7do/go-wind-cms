package data

import (
	"context"
	"testing"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	entCrud "github.com/tx7do/go-crud/entgo"

	"github.com/tx7do/go-utils/mapper"
	"github.com/tx7do/go-utils/trans"

	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/enttest"
	"go-wind-cms/app/core/service/internal/data/ent/posttranslation"

	_ "github.com/xiaoqidun/entps"

	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
)

// newTestPostTranslationRepo 构造一个连内存 sqlite 的 PostTranslationRepo。
// 必须经 enttest 建客户端（enttest 真实 import ent/runtime，填充实体 Hooks，
// 见 interaction_repo_test.go 的注释），事务内读写也统一走 tx 客户端。
// 各测试共用同一块共享内存库，postID 取不同值避免互相干扰。
func newTestPostTranslationRepo(t *testing.T) (*PostTranslationRepo, *ent.Client, func()) {
	t.Helper()

	drv, err := entCrud.CreateDriver(
		"sqlite3",
		"file:ent?mode=memory&cache=shared&_fk=1",
		false, false,
	)
	require.NoError(t, err, "创建 sqlite driver 失败")

	db := enttest.NewClient(t, enttest.WithOptions(
		ent.Driver(drv),
		ent.Log(func(a ...any) { t.Log(a...) }),
	))

	wrapped := entCrud.NewEntClient(db, drv)

	repo := &PostTranslationRepo{
		entClient: wrapped,
		log:       log.NewHelper(log.DefaultLogger),
		mapper:    mapper.NewCopierMapper[contentV1.PostTranslation, ent.PostTranslation](),
	}
	repo.init()

	cleanup := func() {
		_ = db.Close()
	}
	return repo, db, cleanup
}

func upsertPostTranslations(t *testing.T, repo *PostTranslationRepo, db *ent.Client, ctx context.Context, postID uint32, items ...*contentV1.PostTranslation) {
	t.Helper()

	tx, err := db.Tx(ctx)
	require.NoError(t, err, "开启事务失败")
	require.NoError(t, repo.UpsertTranslations(ctx, tx, postID, items), "upsert 译文失败")
	require.NoError(t, tx.Commit(), "提交事务失败")
}

func countPostTranslations(t *testing.T, db *ent.Client, ctx context.Context, postID uint32, languageCode string) int {
	t.Helper()

	builder := db.PostTranslation.Query().Where(posttranslation.PostIDEQ(postID))
	if languageCode != "" {
		builder = builder.Where(posttranslation.LanguageCodeEQ(languageCode))
	}
	n, err := builder.Count(ctx)
	require.NoError(t, err, "统计译文行数失败")
	return n
}

// 新语言应插入，且自动生成 slug/摘要/字数
func TestPostUpsertTranslationsInsertsNewLanguage(t *testing.T) {
	repo, db, cleanup := newTestPostTranslationRepo(t)
	defer cleanup()

	ctx := viewerCtx(1, 1)
	upsertPostTranslations(t, repo, db, ctx, 100, &contentV1.PostTranslation{
		Title:        trans.Ptr("Hello World"),
		Content:      trans.Ptr("some content"),
		LanguageCode: trans.Ptr("en-US"),
	})

	assert.Equal(t, 1, countPostTranslations(t, db, ctx, 100, "en-US"))

	rows, err := db.PostTranslation.Query().Where(posttranslation.PostIDEQ(100)).All(ctx)
	require.NoError(t, err)
	require.Len(t, rows, 1)
	require.NotNil(t, rows[0].Slug)
	assert.NotEmpty(t, *rows[0].Slug, "应自动生成 slug")
	require.NotNil(t, rows[0].Summary)
	assert.NotEmpty(t, *rows[0].Summary, "应自动生成摘要")
}

// 同语言重复 upsert 应就地更新而非追加行；标题未变时 slug 保持不变
func TestPostUpsertTranslationsUpdatesInPlace(t *testing.T) {
	repo, db, cleanup := newTestPostTranslationRepo(t)
	defer cleanup()

	ctx := viewerCtx(1, 1)

	upsertPostTranslations(t, repo, db, ctx, 200, &contentV1.PostTranslation{
		Title:        trans.Ptr("Hello"),
		Content:      trans.Ptr("content v1"),
		LanguageCode: trans.Ptr("en-US"),
	})
	firstRows, err := db.PostTranslation.Query().Where(posttranslation.PostIDEQ(200)).All(ctx)
	require.NoError(t, err)
	require.Len(t, firstRows, 1)
	firstSlug := *firstRows[0].Slug

	upsertPostTranslations(t, repo, db, ctx, 200, &contentV1.PostTranslation{
		Title:        trans.Ptr("Hello"),
		Content:      trans.Ptr("content v2"),
		LanguageCode: trans.Ptr("en-US"),
	})

	assert.Equal(t, 1, countPostTranslations(t, db, ctx, 200, "en-US"), "同语言重复发布不得累积重复行")

	rows, err := db.PostTranslation.Query().Where(posttranslation.PostIDEQ(200)).All(ctx)
	require.NoError(t, err)
	require.Len(t, rows, 1)
	require.NotNil(t, rows[0].Content)
	assert.Equal(t, "content v2", *rows[0].Content, "内容应为最新一次提交")
	assert.Equal(t, firstSlug, *rows[0].Slug, "标题未变时 slug 应保持不变")
	require.NotNil(t, rows[0].Summary)
	assert.Contains(t, *rows[0].Summary, "v2", "客户端未提交摘要时应按最新 content 重算")
}

// 标题变化时应重新生成 slug
func TestPostUpsertTranslationsRegeneratesSlugOnTitleChange(t *testing.T) {
	repo, db, cleanup := newTestPostTranslationRepo(t)
	defer cleanup()

	ctx := viewerCtx(1, 1)

	upsertPostTranslations(t, repo, db, ctx, 250, &contentV1.PostTranslation{
		Title:        trans.Ptr("First Title"),
		Content:      trans.Ptr("c"),
		LanguageCode: trans.Ptr("en-US"),
	})
	upsertPostTranslations(t, repo, db, ctx, 250, &contentV1.PostTranslation{
		Title:        trans.Ptr("Second Title"),
		Content:      trans.Ptr("c"),
		LanguageCode: trans.Ptr("en-US"),
	})

	rows, err := db.PostTranslation.Query().Where(posttranslation.PostIDEQ(250)).All(ctx)
	require.NoError(t, err)
	require.Len(t, rows, 1)
	require.NotNil(t, rows[0].Slug)
	assert.Contains(t, *rows[0].Slug, "second", "标题变更后 slug 应基于新标题生成")
}

// upsert 某一语言不得影响其他语言（回归：page/category/tag 旧实现整表清空重建）
func TestPostUpsertTranslationsPreservesOtherLanguages(t *testing.T) {
	repo, db, cleanup := newTestPostTranslationRepo(t)
	defer cleanup()

	ctx := viewerCtx(1, 1)

	upsertPostTranslations(t, repo, db, ctx, 300, &contentV1.PostTranslation{
		Title:        trans.Ptr("中文标题"),
		Content:      trans.Ptr("中文内容 v1"),
		LanguageCode: trans.Ptr("zh-CN"),
	})
	upsertPostTranslations(t, repo, db, ctx, 300, &contentV1.PostTranslation{
		Title:        trans.Ptr("English Title"),
		Content:      trans.Ptr("english content"),
		LanguageCode: trans.Ptr("en-US"),
	})
	upsertPostTranslations(t, repo, db, ctx, 300, &contentV1.PostTranslation{
		Title:        trans.Ptr("中文标题"),
		Content:      trans.Ptr("中文内容 v2"),
		LanguageCode: trans.Ptr("zh-CN"),
	})

	assert.Equal(t, 1, countPostTranslations(t, db, ctx, 300, "zh-CN"))
	assert.Equal(t, 1, countPostTranslations(t, db, ctx, 300, "en-US"), "更新一语言不得删除其他语言")

	enRows, err := db.PostTranslation.Query().
		Where(posttranslation.PostIDEQ(300), posttranslation.LanguageCodeEQ("en-US")).All(ctx)
	require.NoError(t, err)
	require.Len(t, enRows, 1)
	require.NotNil(t, enRows[0].Content)
	assert.Equal(t, "english content", *enRows[0].Content)
}

// 历史重复行（旧 Update 直插遗留）应被清理：仅保留 id 最小的一行并更新
func TestPostUpsertTranslationsCleansLegacyDuplicates(t *testing.T) {
	repo, db, cleanup := newTestPostTranslationRepo(t)
	defer cleanup()

	ctx := viewerCtx(1, 1)

	// 模拟旧 bug 产生的同语言重复行
	r1, err := db.PostTranslation.Create().
		SetPostID(400).
		SetLanguageCode("zh-CN").
		SetTitle("旧标题A").
		SetSlug("old-a").
		Save(ctx)
	require.NoError(t, err)
	_, err = db.PostTranslation.Create().
		SetPostID(400).
		SetLanguageCode("zh-CN").
		SetTitle("旧标题B").
		SetSlug("old-b").
		Save(ctx)
	require.NoError(t, err)

	upsertPostTranslations(t, repo, db, ctx, 400, &contentV1.PostTranslation{
		Title:        trans.Ptr("新标题"),
		Content:      trans.Ptr("新内容"),
		LanguageCode: trans.Ptr("zh-CN"),
	})

	assert.Equal(t, 1, countPostTranslations(t, db, ctx, 400, "zh-CN"), "历史重复行应被清理")

	rows, err := db.PostTranslation.Query().Where(posttranslation.PostIDEQ(400)).All(ctx)
	require.NoError(t, err)
	require.Len(t, rows, 1)
	assert.Equal(t, r1.ID, rows[0].ID, "应保留 id 最小的原行")
	require.NotNil(t, rows[0].Title)
	assert.Equal(t, "新标题", *rows[0].Title)
}

// 缺译文不是错误：返回 (nil, nil) 交由调用方处理（回归：旧实现 500）
func TestPostGetTranslationMissingReturnsNil(t *testing.T) {
	repo, _, cleanup := newTestPostTranslationRepo(t)
	defer cleanup()

	dto, err := repo.GetTranslation(viewerCtx(1, 1), 999, "en-US")
	require.NoError(t, err)
	assert.Nil(t, dto)
}

// GetTranslation 在存在重复行时取 id 最小的一行，结果确定
func TestPostGetTranslationPicksMinIDRow(t *testing.T) {
	repo, db, cleanup := newTestPostTranslationRepo(t)
	defer cleanup()

	ctx := viewerCtx(1, 1)

	r1, err := db.PostTranslation.Create().
		SetPostID(500).
		SetLanguageCode("ja").
		SetTitle("一行目").
		SetSlug("row-1").
		Save(ctx)
	require.NoError(t, err)
	_, err = db.PostTranslation.Create().
		SetPostID(500).
		SetLanguageCode("ja").
		SetTitle("二行目").
		SetSlug("row-2").
		Save(ctx)
	require.NoError(t, err)

	dto, err := repo.GetTranslation(ctx, 500, "ja")
	require.NoError(t, err)
	require.NotNil(t, dto)
	assert.Equal(t, r1.ID, dto.GetId(), "应取 id 最小的行")
	assert.Equal(t, "一行目", dto.GetTitle())
}

// 指定语言缺译时 ListTranslationsWithFallback 回落全量译文，
// 命中时只回传该语言（C 端 List/Get 嵌入依赖此语义做前端兜底）
func TestPostListTranslationsWithFallback(t *testing.T) {
	repo, db, cleanup := newTestPostTranslationRepo(t)
	defer cleanup()

	ctx := viewerCtx(1, 1)

	upsertPostTranslations(t, repo, db, ctx, 600, &contentV1.PostTranslation{
		Title:        trans.Ptr("中文标题"),
		Content:      trans.Ptr("c"),
		LanguageCode: trans.Ptr("zh-CN"),
	})
	upsertPostTranslations(t, repo, db, ctx, 600, &contentV1.PostTranslation{
		Title:        trans.Ptr("English Title"),
		Content:      trans.Ptr("c"),
		LanguageCode: trans.Ptr("en-US"),
	})

	// 语言命中：仅回传该语言
	hit, err := repo.ListTranslationsWithFallback(ctx, 600, "en-US", nil)
	require.NoError(t, err)
	require.Len(t, hit, 1)
	assert.Equal(t, "en-US", hit[0].GetLanguageCode())

	// 语言缺译：回落全量，前端可按"匹配→回退第一条"兜底
	miss, err := repo.ListTranslationsWithFallback(ctx, 600, "ja", nil)
	require.NoError(t, err)
	assert.Len(t, miss, 2, "缺译语言应回落全量译文")

	// 无 locale：行为与原 ListTranslations 一致
	all, err := repo.ListTranslationsWithFallback(ctx, 600, "", nil)
	require.NoError(t, err)
	assert.Len(t, all, 2)
}
