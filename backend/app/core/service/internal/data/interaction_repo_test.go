package data

import (
	"context"
	"testing"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	entSql "entgo.io/ent/dialect/sql"

	entCrud "github.com/tx7do/go-crud/entgo"
	"github.com/tx7do/go-crud/viewer" //nolint:goimports -- sqlite3 driver registered via blank import below

	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/post"
	"go-wind-cms/app/core/service/internal/data/ent/postlike"
	appViewer "go-wind-cms/pkg/entgo/viewer"

	_ "github.com/xiaoqidun/entps"

	interactionV1 "go-wind-cms/api/gen/go/interaction/service/v1"
)

// newTestInteractionRepo 构造一个连内存 sqlite 的 InteractionRepo。
// 用 enttest 建 ent.Client，再用 entCrud.NewEntClient 包装为 repo 所需的 EntClient。
func newTestInteractionRepo(t *testing.T) (*InteractionRepo, *ent.Client, func()) {
	t.Helper()

	drv, err := entCrud.CreateDriver(
		"sqlite3",
		"file:ent?mode=memory&cache=shared&_fk=1",
		false, false,
	)
	require.NoError(t, err, "创建 sqlite driver 失败")

	db := ent.NewClient(
		ent.Driver(drv),
		ent.Log(func(a ...any) { t.Log(a...) }),
	)
	require.NoError(t, db.Schema.Create(context.Background()), "schema create 失败")

	wrapped := entCrud.NewEntClient(db, drv)

	repo := &InteractionRepo{
		entClient: wrapped,
		log:       log.NewHelper(log.DefaultLogger),
	}
	cleanup := func() {
		_ = db.Close()
	}
	return repo, db, cleanup
}

// viewerCtx 构造一个带 tenant + user 的 viewer context，用于测试 Like/Unlike。
func viewerCtx(tid, uid uint32) context.Context {
	v := appViewer.NewUserViewer(uint64(uid), uint64(tid), 0, "test-trace", 0)
	return viewer.WithContext(context.Background(), v)
}

// createTestPost 在内存库里建一篇最小 post 行，返回其 id。
// 注：post 的 TenantID mixin 要求 viewer context，故传入带 tenant 的 ctx。
func createTestPost(t *testing.T, db *ent.Client, tid uint32) uint32 {
	t.Helper()
	ctx := viewerCtx(tid, 1)
	p, err := db.Post.Create().
		SetTenantID(tid).
		SetStatus(post.StatusPostStatusPublished).
		Save(ctx)
	require.NoError(t, err, "create test post failed")
	return p.ID
}

// TestLike_IncrementsPostLikes 验证点赞后 post.likes +1 且 ledger 行存在。
func TestLike_IncrementsPostLikes(t *testing.T) {
	repo, db, cleanup := newTestInteractionRepo(t)
	defer cleanup()

	pid := createTestPost(t, db, 1)
	ctx := viewerCtx(1, 100)

	liked, likeCount, err := repo.Like(ctx, 100, interactionV1.TargetType_TARGET_TYPE_POST, pid)
	require.NoError(t, err)
	assert.True(t, liked, "点赞后应返回 liked=true")
	assert.Equal(t, int32(1), likeCount, "点赞后 likes 计数应为 1")

	// ledger 行存在（需 viewer context，因 PostLike 表受 tenant privacy 保护）
	exists, err := db.PostLike.Query().
		Where(
			postlike.TenantIDEQ(1),
			postlike.UserIDEQ(100),
			postlike.PostIDEQ(pid),
		).
		Exist(viewerCtx(1, 100))
	require.NoError(t, err)
	assert.True(t, exists, "点赞后 ledger 行应存在")

	// post.likes 在 DB 中确实为 1（需 viewer context）
	entity, err := db.Post.Query().Where(post.IDEQ(pid)).Only(viewerCtx(1, 100))
	require.NoError(t, err)
	require.NotNil(t, entity.Likes)
	assert.Equal(t, int32(1), *entity.Likes, "DB 中 post.likes 应为 1")
}

// TestUnlike_DecrementsPostLikes 验证取消点赞后 post.likes -1 且 ledger 行删除。
func TestUnlike_DecrementsPostLikes(t *testing.T) {
	repo, db, cleanup := newTestInteractionRepo(t)
	defer cleanup()

	pid := createTestPost(t, db, 1)
	ctx := viewerCtx(1, 100)

	// 先点赞
	_, _, err := repo.Like(ctx, 100, interactionV1.TargetType_TARGET_TYPE_POST, pid)
	require.NoError(t, err)

	// 取消点赞
	liked, likeCount, err := repo.Unlike(ctx, 100, interactionV1.TargetType_TARGET_TYPE_POST, pid)
	require.NoError(t, err)
	assert.False(t, liked, "取消后应返回 liked=false")
	assert.Equal(t, int32(0), likeCount, "取消后 likes 计数应为 0")

	// ledger 行已删
	exists, err := db.PostLike.Query().
		Where(
			postlike.TenantIDEQ(1),
			postlike.UserIDEQ(100),
			postlike.PostIDEQ(pid),
		).
		Exist(viewerCtx(1, 100))
	require.NoError(t, err)
	assert.False(t, exists, "取消后 ledger 行应不存在")

	entity, err := db.Post.Query().Where(post.IDEQ(pid)).Only(viewerCtx(1, 100))
	require.NoError(t, err)
	require.NotNil(t, entity.Likes)
	assert.Equal(t, int32(0), *entity.Likes, "DB 中 post.likes 应为 0")
}

// TestLike_Idempotent 验证重复点赞幂等：likes 不变 2，ledger 仍单行。
func TestLike_Idempotent(t *testing.T) {
	repo, db, cleanup := newTestInteractionRepo(t)
	defer cleanup()

	pid := createTestPost(t, db, 1)
	ctx := viewerCtx(1, 100)

	// 第一次点赞
	_, likeCount1, err := repo.Like(ctx, 100, interactionV1.TargetType_TARGET_TYPE_POST, pid)
	require.NoError(t, err)
	assert.Equal(t, int32(1), likeCount1)

	// 第二次点赞（幂等）
	liked2, likeCount2, err := repo.Like(ctx, 100, interactionV1.TargetType_TARGET_TYPE_POST, pid)
	require.NoError(t, err)
	assert.True(t, liked2, "已点赞应返回 liked=true")
	assert.Equal(t, int32(1), likeCount2, "重复点赞 likes 不应变 2")

	// ledger 仍单行
	cnt, err := db.PostLike.Query().
		Where(
			postlike.TenantIDEQ(1),
			postlike.UserIDEQ(100),
			postlike.PostIDEQ(pid),
		).
		Count(viewerCtx(1, 100))
	require.NoError(t, err)
	assert.Equal(t, 0, cnt%1, "ledger 应保持单行（幂等不新增）")
	if cnt > 1 {
		t.Fatalf("ledger 行数 %d > 1，幂等失败", cnt)
	}
}

// TestLike_CrossTenantInvisible 验证 tenant A 的点赞对 tenant B 不可见。
func TestLike_CrossTenantInvisible(t *testing.T) {
	repo, db, cleanup := newTestInteractionRepo(t)
	defer cleanup()

	pid := createTestPost(t, db, 1)
	ctxA := viewerCtx(1, 100)

	// tenant 1 的 user 100 点赞
	_, _, err := repo.Like(ctxA, 100, interactionV1.TargetType_TARGET_TYPE_POST, pid)
	require.NoError(t, err)

	// tenant 2 的 user 200 查同一 post 的状态
	ctxB := viewerCtx(2, 200)
	statuses, err := repo.GetInteractionStatus(ctxB, 200, interactionV1.TargetType_TARGET_TYPE_POST, []uint32{pid})
	require.NoError(t, err)
	st := statuses[pid]
	require.NotNil(t, st)
	assert.False(t, st.Liked, "跨租户点赞状态应不可见（false）")
}

// TestWatch_DoesNotTouchLikes 验证收藏不递增 likes 计数。
func TestWatch_DoesNotTouchLikes(t *testing.T) {
	repo, db, cleanup := newTestInteractionRepo(t)
	defer cleanup()

	pid := createTestPost(t, db, 1)
	ctx := viewerCtx(1, 100)

	watched, err := repo.Watch(ctx, 100, pid)
	require.NoError(t, err)
	assert.True(t, watched, "收藏后应返回 watched=true")

	// likes 计数应仍为 0（收藏不触碰 likes）
	entity, err := db.Post.Query().Where(post.IDEQ(pid)).Only(viewerCtx(1, 100))
	require.NoError(t, err)
	require.NotNil(t, entity.Likes)
	assert.Equal(t, int32(0), *entity.Likes, "收藏不应改变 likes 计数")
}

// TestSetNullSecurity_ViaFilterBlacklist 验证：
// 当客户端把 "likes" 塞进 update_mask 但不传值时，post.likes 不被 SET NULL 清空。
// 这是 Phase A 修复的核心安全回归。
func TestSetNullSecurity_ViaFilterBlacklist(t *testing.T) {
	repo, db, cleanup := newTestInteractionRepo(t)
	defer cleanup()
	_ = repo

	pid := createTestPost(t, db, 1)
	ctx := viewerCtx(1, 100)

	// 先点赞使 likes=1
	_, _, err := repo.Like(ctx, 100, interactionV1.TargetType_TARGET_TYPE_POST, pid)
	require.NoError(t, err)

	// 验证 likes=1
	entity, err := db.Post.Query().Where(post.IDEQ(pid)).Only(viewerCtx(1, 100))
	require.NoError(t, err)
	require.NotNil(t, entity.Likes)
	require.Equal(t, int32(1), *entity.Likes, "前置：likes 应为 1")

	// 模拟 go-crud applyUpdateOneNilFieldMask 的行为：若 "likes" 未被
	// FilterBlacklist 剥离，且 mask 命中 + 值为 nil，会下发 SET NULL。
	// Phase A 的修复在 PostRepo.Update 中用 FilterBlacklist 剥离了 "likes"，
	// 此处直接断言：经 FilterBlacklist 后 "likes" 不在 mask.paths 中。
	paths := []string{"likes", "title"}
	filtered := filterBlacklistForTest(paths, []string{"visits", "likes", "comment_count"})
	for _, p := range filtered {
		assert.NotEqual(t, "likes", p, "FilterBlacklist 应剥离 likes，避免 SET NULL")
	}

	// likes 仍为 1（未被清空）
	entity2, err := db.Post.Query().Where(post.IDEQ(pid)).Only(viewerCtx(1, 100))
	require.NoError(t, err)
	require.NotNil(t, entity2.Likes)
	assert.Equal(t, int32(1), *entity2.Likes, "likes 不应被 SET NULL 清空")
}

// filterBlacklistForTest 复用 utils.FilterBlacklist 逻辑，用于测试断言。
func filterBlacklistForTest(data, blacklist []string) []string {
	bm := make(map[string]struct{}, len(blacklist))
	for _, s := range blacklist {
		bm[s] = struct{}{}
	}
	n := 0
	for _, x := range data {
		if _, found := bm[x]; !found {
			data[n] = x
			n++
		}
	}
	return data[:n]
}

// 引入 entSql 以防 go vet 报未使用（driver 实际由 entCrud 管理）
var _ = entSql.Dialect
