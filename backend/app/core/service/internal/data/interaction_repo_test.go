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

	"go-wind-cms/app/core/service/internal/data/client"
	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/enttest"
	"go-wind-cms/app/core/service/internal/data/ent/interactioncounter"
	"go-wind-cms/app/core/service/internal/data/ent/post"
	"go-wind-cms/app/core/service/internal/data/ent/postlike"
	appViewer "go-wind-cms/pkg/entgo/viewer"

	_ "github.com/xiaoqidun/entps"

	interactionV1 "go-wind-cms/api/gen/go/interaction/service/v1"
)

// newTestInteractionRepo 构造一个连内存 sqlite 的 InteractionRepo。
// 必须经 enttest 建客户端：enttest 真实 import ent/runtime（见 enttest.go），
// 触发其 init() 填充各实体 Hooks[0]，否则任何 mutation 都会报
// "ent: uninitialized hook (forgotten import ent/runtime?)"。
// enttest.NewClient 同时自动跑 schema migration，无需手动 Schema.Create。
func newTestInteractionRepo(t *testing.T) (*InteractionRepo, *ent.Client, func()) {
	t.Helper()

	drv, err := entCrud.CreateDriver(
		"sqlite3",
		"file:ent?mode=memory&cache=shared&_fk=1&_time_format=sqlite",
		false, false,
	)
	require.NoError(t, err, "创建 sqlite driver 失败")

	db := enttest.NewClient(t, enttest.WithOptions(
		ent.Driver(drv),
		ent.Log(func(a ...any) { t.Log(a...) }),
	))

	client.ApplyTimeDefaultHooks(db)

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

// readCounterRow 直接查 interaction_counter 表中 (tenant, target_type, target_id, metric=LIKE) 的行，
// 返回当前计数。不存在行返回 (0, false)。
func readCounterRow(t *testing.T, db *ent.Client, tid, targetID uint32, targetType interactionV1.TargetType, metric interactionV1.CounterMetric) (int64, bool) {
	t.Helper()
	row, err := db.InteractionCounter.Query().
		Where(
			interactioncounter.TenantIDEQ(tid),
			interactioncounter.TargetTypeEQ(uint8(targetType)),
			interactioncounter.TargetIDEQ(targetID),
			interactioncounter.MetricEQ(uint8(metric)),
		).
		Only(viewerCtx(tid, 1))
	if err != nil {
		if ent.IsNotFound(err) {
			return 0, false
		}
		require.NoError(t, err, "query interaction_counter row failed")
	}
	if row.Count == nil {
		return 0, true
	}
	return *row.Count, true
}

// TestLike_IncrementsCounter 验证点赞后 interaction_counter 表出现 count=1 的行，
// 且 ledger 行存在，GetCounts 返回 1。
func TestLike_IncrementsCounter(t *testing.T) {
	repo, db, cleanup := newTestInteractionRepo(t)
	defer cleanup()

	pid := createTestPost(t, db, 1)
	ctx := viewerCtx(1, 100)

	liked, likeCount, err := repo.Like(ctx, 100, interactionV1.TargetType_TARGET_TYPE_POST, pid)
	require.NoError(t, err)
	assert.True(t, liked, "点赞后应返回 liked=true")
	assert.Equal(t, int32(1), likeCount, "点赞后计数应为 1")

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

	// interaction_counter 表出现 count=1 的行
	cnt, exists := readCounterRow(t, db, 1, pid, interactionV1.TargetType_TARGET_TYPE_POST, interactionV1.CounterMetric_COUNTER_METRIC_LIKE)
	assert.True(t, exists, "点赞后 interaction_counter 应有行")
	assert.Equal(t, int64(1), cnt, "interaction_counter 行 count 应为 1")

	// GetCounts 返回 1
	counts, err := repo.GetCounts(viewerCtx(1, 100), interactionV1.TargetType_TARGET_TYPE_POST, []uint32{pid}, []interactionV1.CounterMetric{interactionV1.CounterMetric_COUNTER_METRIC_LIKE})
	require.NoError(t, err)
	cm, ok := counts[pid]
	if assert.True(t, ok, "GetCounts 应返回该 pid 的条目") {
		assert.Equal(t, int64(1), cm.Counts[0].Count, "GetCounts 返回的计数应为 1")
		assert.Equal(t, interactionV1.CounterMetric_COUNTER_METRIC_LIKE, cm.Counts[0].Metric)
	}
}

// TestUnlike_DecrementsCounter 验证取消点赞后 counter 行删除（count→0→删行），ledger 行删除。
func TestUnlike_DecrementsCounter(t *testing.T) {
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
	assert.Equal(t, int32(0), likeCount, "取消后计数应为 0")

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

	// interaction_counter 行应已删（count 归 0 触发删行）
	cnt, exists := readCounterRow(t, db, 1, pid, interactionV1.TargetType_TARGET_TYPE_POST, interactionV1.CounterMetric_COUNTER_METRIC_LIKE)
	assert.False(t, exists, "取消后 interaction_counter 行应删除")
	assert.Equal(t, int64(0), cnt, "无行时计数为 0")

	// GetCounts 不返回该 pid 的条目
	counts, err := repo.GetCounts(viewerCtx(1, 100), interactionV1.TargetType_TARGET_TYPE_POST, []uint32{pid}, []interactionV1.CounterMetric{interactionV1.CounterMetric_COUNTER_METRIC_LIKE})
	require.NoError(t, err)
	_, ok := counts[pid]
	assert.False(t, ok, "取消后 GetCounts 不应返回该 pid 的条目")
}

// TestLike_Idempotent 验证重复点赞幂等：counter 不变 2，ledger 仍单行。
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
	assert.Equal(t, int32(1), likeCount2, "重复点赞计数不变 2")

	// counter 仍 count=1
	cnt, _ := readCounterRow(t, db, 1, pid, interactionV1.TargetType_TARGET_TYPE_POST, interactionV1.CounterMetric_COUNTER_METRIC_LIKE)
	assert.Equal(t, int64(1), cnt, "重复点赞后 counter 仍应为 1")

	// ledger 仍单行
	cntLedger, err := db.PostLike.Query().
		Where(
			postlike.TenantIDEQ(1),
			postlike.UserIDEQ(100),
			postlike.PostIDEQ(pid),
		).
		Count(viewerCtx(1, 100))
	require.NoError(t, err)
	assert.Equal(t, 1, cntLedger, "ledger 应保持单行（幂等不新增）")
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

	// 跨租户 GetCounts 也不可见
	counts, err := repo.GetCounts(ctxB, interactionV1.TargetType_TARGET_TYPE_POST, []uint32{pid}, []interactionV1.CounterMetric{interactionV1.CounterMetric_COUNTER_METRIC_LIKE})
	require.NoError(t, err)
	_, ok := counts[pid]
	assert.False(t, ok, "跨租户 GetCounts 不应返回该 pid 的条目")
}

// TestWatch_DoesNotTouchCounter 验证收藏不在 counter 表中产生行。
// TestWatch_IncrementsWatchCounter 验证收藏后 interaction_counter 表出现 WATCH metric 行 count=1，
// 返回 watchCount=1，且不影响 LIKE metric 行（LIKE 行不存在）。
func TestWatch_IncrementsWatchCounter(t *testing.T) {
	repo, db, cleanup := newTestInteractionRepo(t)
	defer cleanup()

	pid := createTestPost(t, db, 1)
	ctx := viewerCtx(1, 100)

	watched, watchCount, err := repo.Watch(ctx, 100, pid)
	require.NoError(t, err)
	assert.True(t, watched, "收藏后应返回 watched=true")
	assert.Equal(t, int32(1), watchCount, "收藏后 watchCount 应为 1")

	// WATCH metric 行存在，count=1
	watchCnt, watchExists := readCounterRow(t, db, 1, pid, interactionV1.TargetType_TARGET_TYPE_POST, interactionV1.CounterMetric_COUNTER_METRIC_WATCH)
	assert.True(t, watchExists, "收藏后 interaction_counter 应有 WATCH 行")
	assert.Equal(t, int64(1), watchCnt, "WATCH 行 count 应为 1")

	// LIKE metric 行不应存在（收藏不影响点赞计数）
	likeCnt, likeExists := readCounterRow(t, db, 1, pid, interactionV1.TargetType_TARGET_TYPE_POST, interactionV1.CounterMetric_COUNTER_METRIC_LIKE)
	assert.False(t, likeExists, "收藏不应产生 LIKE 行")
	assert.Equal(t, int64(0), likeCnt, "LIKE 行不存在时计数为 0")
}

// TestUnwatch_DecrementsWatchCounter 验证取消收藏后 WATCH metric 行删除，watchCount=0。
func TestUnwatch_DecrementsWatchCounter(t *testing.T) {
	repo, db, cleanup := newTestInteractionRepo(t)
	defer cleanup()

	pid := createTestPost(t, db, 1)
	ctx := viewerCtx(1, 100)

	// 先收藏
	_, _, err := repo.Watch(ctx, 100, pid)
	require.NoError(t, err)

	// 取消收藏
	watched, watchCount, err := repo.Unwatch(ctx, 100, pid)
	require.NoError(t, err)
	assert.False(t, watched, "取消后应返回 watched=false")
	assert.Equal(t, int32(0), watchCount, "取消后 watchCount 应为 0")

	// WATCH 行应已删
	watchCnt, watchExists := readCounterRow(t, db, 1, pid, interactionV1.TargetType_TARGET_TYPE_POST, interactionV1.CounterMetric_COUNTER_METRIC_WATCH)
	assert.False(t, watchExists, "取消后 WATCH 行应删除")
	assert.Equal(t, int64(0), watchCnt, "无行时计数为 0")
}

// TestWatch_Idempotent 验证重复收藏幂等：WATCH 行 count 不变 2，ledger 仍单行。
func TestWatch_Idempotent(t *testing.T) {
	repo, db, cleanup := newTestInteractionRepo(t)
	defer cleanup()

	pid := createTestPost(t, db, 1)
	ctx := viewerCtx(1, 100)

	// 第一次收藏
	_, watchCount1, err := repo.Watch(ctx, 100, pid)
	require.NoError(t, err)
	assert.Equal(t, int32(1), watchCount1)

	// 第二次收藏（幂等）
	watched2, watchCount2, err := repo.Watch(ctx, 100, pid)
	require.NoError(t, err)
	assert.True(t, watched2, "已收藏应返回 watched=true")
	assert.Equal(t, int32(1), watchCount2, "重复收藏 watchCount 不变 2")

	// WATCH 行仍 count=1
	watchCnt, _ := readCounterRow(t, db, 1, pid, interactionV1.TargetType_TARGET_TYPE_POST, interactionV1.CounterMetric_COUNTER_METRIC_WATCH)
	assert.Equal(t, int64(1), watchCnt, "重复收藏后 WATCH 行仍应为 1")
}

// TestGetCounts_Batch 验证多目标批量查询返回正确的 map。
func TestGetCounts_Batch(t *testing.T) {
	repo, db, cleanup := newTestInteractionRepo(t)
	defer cleanup()

	pid1 := createTestPost(t, db, 1)
	pid2 := createTestPost(t, db, 1)
	ctx := viewerCtx(1, 100)

	// 给 pid1 点赞 2 次（幂等，仍 count=1），pid2 点赞 1 次
	_, _, err := repo.Like(ctx, 100, interactionV1.TargetType_TARGET_TYPE_POST, pid1)
	require.NoError(t, err)
	_, _, err = repo.Like(ctx, 100, interactionV1.TargetType_TARGET_TYPE_POST, pid2)
	require.NoError(t, err)

	// 批量查 [pid1, pid2]：应都返回 count=1
	counts, err := repo.GetCounts(ctx, interactionV1.TargetType_TARGET_TYPE_POST, []uint32{pid1, pid2}, []interactionV1.CounterMetric{interactionV1.CounterMetric_COUNTER_METRIC_LIKE})
	require.NoError(t, err)
	assert.Len(t, counts, 2, "应返回两个目标的条目")
	for _, tid := range []uint32{pid1, pid2} {
		cm, ok := counts[tid]
		if !assert.True(t, ok, "目标 %d 应有条目", tid) {
			continue
		}
		assert.Len(t, cm.Counts, 1, "目标 %d 应有 1 个 metric 条目", tid)
		if assert.Len(t, cm.Counts, 1, "目标 %d 应有 1 个 metric 条目", tid) {
			assert.Equal(t, interactionV1.CounterMetric_COUNTER_METRIC_LIKE, cm.Counts[0].Metric)
			assert.Equal(t, int64(1), cm.Counts[0].Count, "目标 %d 计数应为 1", tid)
		}
	}
}

// TestGetCounts_Unauthenticated 验证无 viewer context 时返回 Unauthorized。
func TestGetCounts_Unauthenticated(t *testing.T) {
	repo, _, cleanup := newTestInteractionRepo(t)
	defer cleanup()

	// 无 viewer context（context.Background 无 tenant）
	_, err := repo.GetCounts(context.Background(), interactionV1.TargetType_TARGET_TYPE_POST, []uint32{1}, []interactionV1.CounterMetric{interactionV1.CounterMetric_COUNTER_METRIC_LIKE})
	assert.Error(t, err, "无 viewer context 应返回错误")
}

// TestGetCounts_EmptyRequest 验证空 target_ids 或 metrics 返回空 map。
func TestGetCounts_EmptyRequest(t *testing.T) {
	repo, _, cleanup := newTestInteractionRepo(t)
	defer cleanup()
	ctx := viewerCtx(1, 100)

	// 空 target_ids
	counts, err := repo.GetCounts(ctx, interactionV1.TargetType_TARGET_TYPE_POST, []uint32{}, []interactionV1.CounterMetric{interactionV1.CounterMetric_COUNTER_METRIC_LIKE})
	require.NoError(t, err)
	assert.Empty(t, counts, "空 target_ids 应返回空 map")

	// 空 metrics
	counts2, err := repo.GetCounts(ctx, interactionV1.TargetType_TARGET_TYPE_POST, []uint32{1}, []interactionV1.CounterMetric{})
	require.NoError(t, err)
	assert.Empty(t, counts2, "空 metrics 应返回空 map")
}

// 引入 entSql 以防 go vet 报未使用（driver 实际由 entCrud 管理）
var _ = entSql.Dialect

// ==============================
// 清数据测试（PurgeTargetInteractions / PurgeUserInteractions / ResetCounter）
// ==============================

// TestPurgeTargetInteractions_DeletesLedgerAndCounter 验证 purge 单条 target：
// 两个用户给同 post 点赞后，purge 该 post → post_like 行全删 + LIKE counter 行删。
func TestPurgeTargetInteractions_DeletesLedgerAndCounter(t *testing.T) {
	repo, db, cleanup := newTestInteractionRepo(t)
	defer cleanup()

	pid := createTestPost(t, db, 1)
	ctxA := viewerCtx(1, 100)
	ctxB := viewerCtx(1, 200)

	// 两用户各点赞
	_, _, err := repo.Like(ctxA, 100, interactionV1.TargetType_TARGET_TYPE_POST, pid)
	require.NoError(t, err)
	_, _, err = repo.Like(ctxB, 200, interactionV1.TargetType_TARGET_TYPE_POST, pid)
	require.NoError(t, err)

	// counter 行存在 count=2
	cnt, exists := readCounterRow(t, db, 1, pid, interactionV1.TargetType_TARGET_TYPE_POST, interactionV1.CounterMetric_COUNTER_METRIC_LIKE)
	assert.True(t, exists, "两用户点赞后应有 LIKE counter 行")
	assert.Equal(t, int64(2), cnt, "counter 应为 2")

	// purge 该 target（用 tenant 1 的 viewer context）
	affected, err := repo.PurgeTargetInteractions(ctxA, interactionV1.TargetType_TARGET_TYPE_POST, pid)
	require.NoError(t, err)
	assert.Equal(t, uint32(2), affected, "应删 2 行 ledger")

	// ledger 行已删
	ledgerCnt, err := db.PostLike.Query().
		Where(
			postlike.TenantIDEQ(1),
			postlike.PostIDEQ(pid),
		).
		Count(viewerCtx(1, 100))
	require.NoError(t, err)
	assert.Equal(t, 0, ledgerCnt, "purge 后 post_like 应无该 post 的行")

	// counter 行已删
	cnt, exists = readCounterRow(t, db, 1, pid, interactionV1.TargetType_TARGET_TYPE_POST, interactionV1.CounterMetric_COUNTER_METRIC_LIKE)
	_ = cnt
	assert.False(t, exists, "purge 后 LIKE counter 行应删除")
}

// TestPurgeTargetInteractions_Idempotent 验证对无交互的 target purge 是 no-op，返回 0。
func TestPurgeTargetInteractions_Idempotent(t *testing.T) {
	repo, db, cleanup := newTestInteractionRepo(t)
	defer cleanup()

	pid := createTestPost(t, db, 1)
	ctx := viewerCtx(1, 100)

	// 无任何交互，直接 purge
	affected, err := repo.PurgeTargetInteractions(ctx, interactionV1.TargetType_TARGET_TYPE_POST, pid)
	require.NoError(t, err)
	assert.Equal(t, uint32(0), affected, "无交互的 target purge 应返回 0")

	// counter 行不存在
	_, exists := readCounterRow(t, db, 1, pid, interactionV1.TargetType_TARGET_TYPE_POST, interactionV1.CounterMetric_COUNTER_METRIC_LIKE)
	assert.False(t, exists, "无交互时不应有 counter 行")
}

// TestPurgeUserInteractions_DeletesAcrossTargets 验证 purge 某用户在全站的交互：
// 1 用户给 3 个 post 点赞 → purge 该 user → 三 target 的 ledger + counter 全清。
func TestPurgeUserInteractions_DeletesAcrossTargets(t *testing.T) {
	repo, db, cleanup := newTestInteractionRepo(t)
	defer cleanup()

	pid1 := createTestPost(t, db, 1)
	pid2 := createTestPost(t, db, 1)
	pid3 := createTestPost(t, db, 1)
	ctx := viewerCtx(1, 100)

	// 用户 100 给 3 个 post 各点赞
	for _, pid := range []uint32{pid1, pid2, pid3} {
		_, _, err := repo.Like(ctx, 100, interactionV1.TargetType_TARGET_TYPE_POST, pid)
		require.NoError(t, err)
	}

	// 三 target 各有 counter 行 count=1
	for _, pid := range []uint32{pid1, pid2, pid3} {
		cnt, exists := readCounterRow(t, db, 1, pid, interactionV1.TargetType_TARGET_TYPE_POST, interactionV1.CounterMetric_COUNTER_METRIC_LIKE)
		assert.True(t, exists, "post %d 应有 counter 行", pid)
		assert.Equal(t, int64(1), cnt, "post %d counter 应为 1", pid)
	}

	// purge 该用户
	affected, err := repo.PurgeUserInteractions(ctx, 100)
	require.NoError(t, err)
	assert.Equal(t, uint32(3), affected, "应删 3 行 ledger")

	// 三 target 的 ledger + counter 全清
	for _, pid := range []uint32{pid1, pid2, pid3} {
		ledgerCnt, err := db.PostLike.Query().
			Where(
				postlike.TenantIDEQ(1),
				postlike.PostIDEQ(pid),
			).
			Count(viewerCtx(1, 100))
		require.NoError(t, err)
		assert.Equal(t, 0, ledgerCnt, "post %d 的 post_like 应清空", pid)

		_, exists := readCounterRow(t, db, 1, pid, interactionV1.TargetType_TARGET_TYPE_POST, interactionV1.CounterMetric_COUNTER_METRIC_LIKE)
		assert.False(t, exists, "post %d 的 LIKE counter 行应删除", pid)
	}
}

// TestResetCounter_RecountsFromLedger 验证 reset 把 counter 校正为 ledger 真实计数：
// 点赞 1 次（counter=1）→ 手动把 counter 改高到 5 → reset → recount=1，counter 行 count=1。
func TestResetCounter_RecountsFromLedger(t *testing.T) {
	repo, db, cleanup := newTestInteractionRepo(t)
	defer cleanup()

	pid := createTestPost(t, db, 1)
	ctx := viewerCtx(1, 100)

	// 点赞 1 次，counter=1
	_, _, err := repo.Like(ctx, 100, interactionV1.TargetType_TARGET_TYPE_POST, pid)
	require.NoError(t, err)

	// 手动把 counter 行改高到 5（模拟漂移）
	row, err := db.InteractionCounter.Query().
		Where(
			interactioncounter.TenantIDEQ(1),
			interactioncounter.TargetIDEQ(pid),
			interactioncounter.TargetTypeEQ(uint8(interactionV1.TargetType_TARGET_TYPE_POST)),
			interactioncounter.MetricEQ(uint8(interactionV1.CounterMetric_COUNTER_METRIC_LIKE)),
		).
		Only(viewerCtx(1, 100))
	require.NoError(t, err)
	_, err = db.InteractionCounter.UpdateOneID(row.ID).SetCount(5).Save(viewerCtx(1, 100))
	require.NoError(t, err)

	// reset：应重算为 1（ledger 真实计数）
	recount, err := repo.ResetCounter(ctx, interactionV1.TargetType_TARGET_TYPE_POST, pid, interactionV1.CounterMetric_COUNTER_METRIC_LIKE)
	require.NoError(t, err)
	assert.Equal(t, int64(1), recount, "reset 后 recount 应为 ledger 真实计数 1")

	// counter 行被校正回 1
	cnt, exists := readCounterRow(t, db, 1, pid, interactionV1.TargetType_TARGET_TYPE_POST, interactionV1.CounterMetric_COUNTER_METRIC_LIKE)
	assert.True(t, exists, "reset 后 counter 行应存在（count=1）")
	assert.Equal(t, int64(1), cnt, "counter 应被校正回 1")
}

// TestResetCounter_ZeroRecountDeletesRow 验证 ledger 已空时 reset 删除 counter 行：
// 点赞后 unlike（ledger 空，counter 已删）→ reset → recount=0，无 counter 行。
func TestResetCounter_ZeroRecountDeletesRow(t *testing.T) {
	repo, db, cleanup := newTestInteractionRepo(t)
	defer cleanup()

	pid := createTestPost(t, db, 1)
	ctx := viewerCtx(1, 100)

	// 点赞后取消（counter 行已删）
	_, _, err := repo.Like(ctx, 100, interactionV1.TargetType_TARGET_TYPE_POST, pid)
	require.NoError(t, err)
	_, _, err = repo.Unlike(ctx, 100, interactionV1.TargetType_TARGET_TYPE_POST, pid)
	require.NoError(t, err)

	// 此时 counter 行已不存在（unlike 归 0 删行）
	_, exists := readCounterRow(t, db, 1, pid, interactionV1.TargetType_TARGET_TYPE_POST, interactionV1.CounterMetric_COUNTER_METRIC_LIKE)
	assert.False(t, exists, "unlike 后 counter 行应已删")

	// reset：recount=0，无行
	recount, err := repo.ResetCounter(ctx, interactionV1.TargetType_TARGET_TYPE_POST, pid, interactionV1.CounterMetric_COUNTER_METRIC_LIKE)
	require.NoError(t, err)
	assert.Equal(t, int64(0), recount, "ledger 空时 recount 应为 0")

	// 仍无 counter 行
	_, exists = readCounterRow(t, db, 1, pid, interactionV1.TargetType_TARGET_TYPE_POST, interactionV1.CounterMetric_COUNTER_METRIC_LIKE)
	assert.False(t, exists, "reset 后不应有 counter 行")
}

// TestPurgeTargetInteractions_CrossTenantIsolated 验证 purge tenant A 的 target
// 不影响 tenant B 的 counter。
func TestPurgeTargetInteractions_CrossTenantIsolated(t *testing.T) {
	repo, db, cleanup := newTestInteractionRepo(t)
	defer cleanup()

	// tenant 1 与 tenant 2 各建一个 post
	pidA := createTestPost(t, db, 1)
	pidB := createTestPost(t, db, 2)
	ctxA := viewerCtx(1, 100)
	ctxB := viewerCtx(2, 200)

	// 两租户各点赞自己的 post
	_, _, err := repo.Like(ctxA, 100, interactionV1.TargetType_TARGET_TYPE_POST, pidA)
	require.NoError(t, err)
	_, _, err = repo.Like(ctxB, 200, interactionV1.TargetType_TARGET_TYPE_POST, pidB)
	require.NoError(t, err)

	// purge tenant 1 的 target
	affected, err := repo.PurgeTargetInteractions(ctxA, interactionV1.TargetType_TARGET_TYPE_POST, pidA)
	require.NoError(t, err)
	assert.Equal(t, uint32(1), affected, "tenant 1 purge 应删 1 行")

	// tenant 1 的 counter 行已删
	_, existsA := readCounterRow(t, db, 1, pidA, interactionV1.TargetType_TARGET_TYPE_POST, interactionV1.CounterMetric_COUNTER_METRIC_LIKE)
	assert.False(t, existsA, "tenant 1 的 counter 行应删")

	// tenant 2 的 counter 行仍存在 count=1，不受影响
	cntB, existsB := readCounterRow(t, db, 2, pidB, interactionV1.TargetType_TARGET_TYPE_POST, interactionV1.CounterMetric_COUNTER_METRIC_LIKE)
	_ = cntB
	assert.True(t, existsB, "tenant 2 的 counter 行应仍存在（跨租户隔离）")
}
