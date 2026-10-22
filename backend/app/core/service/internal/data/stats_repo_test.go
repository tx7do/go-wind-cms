package data

import (
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	entCrud "github.com/tx7do/go-crud/entgo"

	interactionV1 "go-wind-cms/api/gen/go/interaction/service/v1"

	"go-wind-cms/app/core/service/internal/data/client"
	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/comment"
	"go-wind-cms/app/core/service/internal/data/ent/enttest"
	"go-wind-cms/app/core/service/internal/data/ent/loginauditlog"
	"go-wind-cms/app/core/service/internal/data/ent/mediaasset"
	"go-wind-cms/app/core/service/internal/data/ent/post"
)

// newTestStatsRepo 构造一个连内存 sqlite 的 StatsRepo。
// 同 newTestInteractionRepo:必须经 enttest 建客户端以触发 ent/runtime 的
// hook 注册,并自动跑 schema migration;同时注册与生产一致的时间兜底 hook。
func newTestStatsRepo(t *testing.T) (*StatsRepo, *ent.Client, func()) {
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

	repo := &StatsRepo{
		entClient: wrapped,
		log:       log.NewHelper(log.DefaultLogger),
	}
	cleanup := func() {
		_ = db.Close()
	}
	return repo, db, cleanup
}

func createTestUser(t *testing.T, db *ent.Client, tid uint32, name string) {
	t.Helper()
	_, err := db.User.Create().
		SetTenantID(tid).
		SetUsername(name).
		Save(viewerCtx(tid, 1))
	require.NoError(t, err, "create test user failed")
}

func createTestComment(t *testing.T, db *ent.Client, tid uint32, spam bool) {
	t.Helper()
	c := db.Comment.Create().
		SetTenantID(tid).
		SetContent("test comment")
	if spam {
		c = c.SetStatus(comment.StatusStatusSpam)
	}
	_, err := c.Save(viewerCtx(tid, 1))
	require.NoError(t, err, "create test comment failed")
}

func createTestLoginLog(t *testing.T, db *ent.Client, tid uint32, action loginauditlog.ActionType, status loginauditlog.Status, at time.Time) {
	t.Helper()
	_, err := db.LoginAuditLog.Create().
		SetTenantID(tid).
		SetActionType(action).
		SetStatus(status).
		SetCreatedAt(at).
		Save(viewerCtx(tid, 1))
	require.NoError(t, err, "create test login log failed")
}

func createTestCounter(t *testing.T, db *ent.Client, tid, targetID uint32, metric interactionV1.CounterMetric, count int64) {
	t.Helper()
	_, err := db.InteractionCounter.Create().
		SetTenantID(tid).
		SetTargetType(uint8(interactionV1.TargetType_TARGET_TYPE_POST)).
		SetTargetID(targetID).
		SetMetric(uint8(metric)).
		SetCount(count).
		Save(viewerCtx(tid, 1))
	require.NoError(t, err, "create test counter failed")
}

func createTestMediaAsset(t *testing.T, db *ent.Client, tid uint32, size int64) {
	t.Helper()
	_, err := db.MediaAsset.Create().
		SetTenantID(tid).
		SetSize(uint64(size)).
		Save(viewerCtx(tid, 1))
	require.NoError(t, err, "create test media asset failed")
}

// TestStatsRepo_GetOverview 验证总量计数、口径排除(TRASHED/SPAM)与近 7 天增量。
func TestStatsRepo_GetOverview(t *testing.T) {
	repo, db, cleanup := newTestStatsRepo(t)
	defer cleanup()

	ctx := viewerCtx(1, 1)

	createTestUser(t, db, 1, "u1")
	createTestUser(t, db, 1, "u2")

	createTestPost(t, db, 1)
	createTestPost(t, db, 1)
	// 回收站帖子不计入
	_, err := db.Post.Create().
		SetTenantID(1).
		SetStatus(post.StatusPostStatusTrashed).
		Save(ctx)
	require.NoError(t, err)

	createTestComment(t, db, 1, false)
	// SPAM 评论不计入
	createTestComment(t, db, 1, true)
	// 待审评论计入待办数
	_, err = db.Comment.Create().
		SetTenantID(1).
		SetContent("pending comment").
		SetStatus(comment.StatusStatusPending).
		Save(ctx)
	require.NoError(t, err)

	pid := createTestPost(t, db, 1)
	_, err = db.PostLike.Create().
		SetTenantID(1).
		SetUserID(100).
		SetPostID(pid).
		Save(ctx)
	require.NoError(t, err)

	createTestCounter(t, db, 1, pid, interactionV1.CounterMetric_COUNTER_METRIC_LIKE, 5)
	createTestCounter(t, db, 1, pid, interactionV1.CounterMetric_COUNTER_METRIC_WATCH, 3)

	// 媒体库：2 个资产共 300 字节；草稿/定时各 1 篇计入待办数
	createTestMediaAsset(t, db, 1, 100)
	createTestMediaAsset(t, db, 1, 200)
	_, err = db.Post.Create().
		SetTenantID(1).
		SetStatus(post.StatusPostStatusDraft).
		Save(ctx)
	require.NoError(t, err)
	_, err = db.Post.Create().
		SetTenantID(1).
		SetStatus(post.StatusPostStatusScheduled).
		Save(ctx)
	require.NoError(t, err)

	// 其他租户数据不应出现
	createTestUser(t, db, 2, "other-tenant-user")
	createTestMediaAsset(t, db, 2, 999)

	counts, err := repo.GetOverview(ctx)
	require.NoError(t, err)

	assert.Equal(t, int64(2), counts.UserCount, "用户总量应为 2(不含其他租户)")
	assert.Equal(t, int64(2), counts.NewUserCountWeek, "近 7 天新增用户应为 2")
	assert.Equal(t, int64(5), counts.PostCount, "帖子总量应为 5(排除回收站 1 篇)")
	assert.Equal(t, int64(5), counts.NewPostCountWeek, "近 7 天新增帖子应为 5")
	assert.Equal(t, int64(2), counts.CommentCount, "评论总量应为 2(排除 SPAM)")
	assert.Equal(t, int64(2), counts.NewCommentCount, "近 7 天新增评论应为 2")
	assert.Equal(t, int64(8), counts.InteractionCount, "互动总量应为 5+3")
	assert.Equal(t, int64(1), counts.NewLikeCountWeek, "近 7 天新增点赞应为 1")
	assert.Equal(t, int64(2), counts.MediaAssetCount, "媒体资产应为 2(不含其他租户)")
	assert.Equal(t, int64(300), counts.MediaStorageBytes, "媒体存储应为 100+200 字节")
	assert.Equal(t, int64(2), counts.NewMediaCountWeek, "近 7 天新增媒体应为 2")
	assert.Equal(t, int64(300), counts.NewMediaStorageWeek, "近 7 天新增媒体存储应为 300 字节")
	assert.Equal(t, int64(1), counts.PendingCommentCount, "待审评论应为 1")
	assert.Equal(t, int64(1), counts.DraftPostCount, "草稿帖子应为 1")
	assert.Equal(t, int64(1), counts.ScheduledPostCount, "定时发布帖子应为 1")
}

// TestStatsRepo_GetDailyTrend 验证按天序列长度、日期连续与当日计数。
func TestStatsRepo_GetDailyTrend(t *testing.T) {
	repo, db, cleanup := newTestStatsRepo(t)
	defer cleanup()

	ctx := viewerCtx(1, 1)

	createTestUser(t, db, 1, "u1")
	createTestUser(t, db, 1, "u2")

	users, posts, comments, err := repo.GetDailyTrend(ctx, 7)
	require.NoError(t, err)

	assert.Len(t, users, 7, "用户序列应为 7 天")
	assert.Len(t, posts, 7, "帖子序列应为 7 天")
	assert.Len(t, comments, 7, "评论序列应为 7 天")

	today := time.Now().Format("2006-01-02")
	assert.Equal(t, today, users[6].Date, "序列最后一天应为今天(UTC)")
	assert.Equal(t, int64(2), users[6].Value, "今日新增用户应为 2")
	for i := 0; i < 6; i++ {
		assert.Equal(t, int64(0), users[i].Value, "无数据日期应补 0")
	}
	assert.Equal(t, int64(0), posts[6].Value, "今日无新增帖子")
}

// TestStatsRepo_GetDailyTrend 跨天窗口:昨天的数据归入倒数第二天。
func TestStatsRepo_GetDailyTrend_Yesterday(t *testing.T) {
	repo, db, cleanup := newTestStatsRepo(t)
	defer cleanup()

	ctx := viewerCtx(1, 1)
	// 本地昨天 01:00,稳定落在前一个统计日(避免跨零点边界)
	now := time.Now()
	yesterday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, -1).Add(time.Hour)
	_, err := db.User.Create().
		SetTenantID(1).
		SetUsername("yesterday-user").
		SetCreatedAt(yesterday).
		Save(ctx)
	require.NoError(t, err)

	users, _, _, err := repo.GetDailyTrend(ctx, 7)
	require.NoError(t, err)
	assert.Equal(t, int64(1), users[5].Value, "昨天创建的用户应归入倒数第二天")
	assert.Equal(t, int64(0), users[6].Value, "今天应无新增")
}

// TestStatsRepo_GetTopLikedPosts 验证 TOP 排序、标题关联与回收站剔除。
func TestStatsRepo_GetTopLikedPosts(t *testing.T) {
	repo, db, cleanup := newTestStatsRepo(t)
	defer cleanup()

	ctx := viewerCtx(1, 1)

	pid1 := createTestPost(t, db, 1)
	pid2 := createTestPost(t, db, 1)

	_, err := db.PostTranslation.Create().
		SetTenantID(1).
		SetPostID(pid2).
		SetTitle("Post Two").
		Save(ctx)
	require.NoError(t, err)

	createTestCounter(t, db, 1, pid1, interactionV1.CounterMetric_COUNTER_METRIC_LIKE, 5)
	createTestCounter(t, db, 1, pid2, interactionV1.CounterMetric_COUNTER_METRIC_LIKE, 9)

	top, err := repo.GetTopLikedPosts(ctx, 10)
	require.NoError(t, err)
	require.Len(t, top, 2, "应有 2 条 TOP 记录")
	assert.Equal(t, pid2, top[0].TargetID, "点赞 9 的帖子应排第一")
	assert.Equal(t, int64(9), top[0].LikeCount)
	assert.Equal(t, "Post Two", top[0].Title, "应关联翻译标题")
	assert.Equal(t, pid1, top[1].TargetID)
	assert.Empty(t, top[1].Title, "无翻译的帖子标题为空串")

	// 回收 pid2 后应被剔除
	_, err = db.Post.UpdateOneID(pid2).
		SetStatus(post.StatusPostStatusTrashed).
		Save(ctx)
	require.NoError(t, err)

	top, err = repo.GetTopLikedPosts(ctx, 10)
	require.NoError(t, err)
	require.Len(t, top, 1, "回收站帖子应从榜单剔除")
	assert.Equal(t, pid1, top[0].TargetID)
}

// TestStatsRepo_GetLoginActivity 验证登录成功/失败/活跃用户按天聚合,以及非 LOGIN 动作与 PARTIAL 状态不计入。
func TestStatsRepo_GetLoginActivity(t *testing.T) {
	repo, db, cleanup := newTestStatsRepo(t)
	defer cleanup()

	now := time.Now()
	createTestLoginLog(t, db, 1, loginauditlog.ActionTypeLogin, loginauditlog.StatusSuccess, now)
	createTestLoginLog(t, db, 1, loginauditlog.ActionTypeLogin, loginauditlog.StatusSuccess, now)
	createTestLoginLog(t, db, 1, loginauditlog.ActionTypeLogin, loginauditlog.StatusFailed, now)
	// LOGOUT 与 PARTIAL 不应计入
	createTestLoginLog(t, db, 1, loginauditlog.ActionTypeLogout, loginauditlog.StatusSuccess, now)
	createTestLoginLog(t, db, 1, loginauditlog.ActionTypeLogin, loginauditlog.StatusPartial, now)
	// user_id=7 成功登录两次(另一次独立创建),活跃用户去重后应为 1
	for i := 0; i < 2; i++ {
		_, err := db.LoginAuditLog.Create().
			SetTenantID(1).
			SetUserID(7).
			SetActionType(loginauditlog.ActionTypeLogin).
			SetStatus(loginauditlog.StatusSuccess).
			SetCreatedAt(now).
			Save(viewerCtx(1, 1))
		require.NoError(t, err)
	}

	success, failed, activeUsers, err := repo.GetLoginActivity(viewerCtx(1, 1), 7)
	require.NoError(t, err)

	assert.Len(t, success, 7)
	assert.Len(t, failed, 7)
	assert.Len(t, activeUsers, 7)
	assert.Equal(t, int64(4), success[6].Value, "今日登录成功应为 4")
	assert.Equal(t, int64(1), failed[6].Value, "今日登录失败应为 1")
	assert.Equal(t, int64(1), activeUsers[6].Value, "今日活跃用户去重后应为 1(user_id=0 的行不计)")
}

// TestStatsRepo_GetAiUsageTrend 验证调用次数/token 按天聚合与租户隔离。
func TestStatsRepo_GetAiUsageTrend(t *testing.T) {
	repo, db, cleanup := newTestStatsRepo(t)
	defer cleanup()

	ctx := viewerCtx(1, 1)
	for i := 0; i < 2; i++ {
		_, err := db.AiUsageLog.Create().
			SetTenantID(1).
			SetTotalTokens(100).
			Save(ctx)
		require.NoError(t, err)
	}
	// 其他租户不计入
	_, err := db.AiUsageLog.Create().
		SetTenantID(2).
		SetTotalTokens(999).
		Save(viewerCtx(2, 1))
	require.NoError(t, err)

	calls, tokens, err := repo.GetAiUsageTrend(ctx, 7)
	require.NoError(t, err)

	assert.Len(t, calls, 7)
	assert.Len(t, tokens, 7)
	assert.Equal(t, int64(2), calls[6].Value, "今日 AI 调用应为 2")
	assert.Equal(t, int64(200), tokens[6].Value, "今日 token 消耗应为 200")
}

// TestStatsRepo_GetApiCallTrend 验证 API 请求按天聚合与租户隔离。
func TestStatsRepo_GetApiCallTrend(t *testing.T) {
	repo, db, cleanup := newTestStatsRepo(t)
	defer cleanup()

	ctx := viewerCtx(1, 1)
	for i := 0; i < 3; i++ {
		_, err := db.ApiAuditLog.Create().
			SetTenantID(1).
			Save(ctx)
		require.NoError(t, err)
	}
	_, err := db.ApiAuditLog.Create().
		SetTenantID(2).
		Save(viewerCtx(2, 1))
	require.NoError(t, err)

	calls, err := repo.GetApiCallTrend(ctx, 7)
	require.NoError(t, err)

	assert.Len(t, calls, 7)
	assert.Equal(t, int64(3), calls[6].Value, "今日 API 请求应为 3(不含其他租户)")
}

// TestStatsRepo_GetPostStatusDistribution 验证状态分组与回收站排除。
func TestStatsRepo_GetPostStatusDistribution(t *testing.T) {
	repo, db, cleanup := newTestStatsRepo(t)
	defer cleanup()

	ctx := viewerCtx(1, 1)
	createTestPost(t, db, 1) // PUBLISHED
	_, err := db.Post.Create().SetTenantID(1).SetStatus(post.StatusPostStatusDraft).Save(ctx)
	require.NoError(t, err)
	_, err = db.Post.Create().SetTenantID(1).SetStatus(post.StatusPostStatusDraft).Save(ctx)
	require.NoError(t, err)
	_, err = db.Post.Create().SetTenantID(1).SetStatus(post.StatusPostStatusTrashed).Save(ctx)
	require.NoError(t, err)

	rows, err := repo.GetPostStatusDistribution(ctx)
	require.NoError(t, err)

	byStatus := make(map[string]int64, len(rows))
	for _, row := range rows {
		byStatus[row.Status] = row.Count
	}
	assert.Equal(t, int64(2), byStatus[string(post.StatusPostStatusDraft)], "草稿应为 2")
	assert.Equal(t, int64(1), byStatus[string(post.StatusPostStatusPublished)], "已发布应为 1")
	_, hasTrashed := byStatus[string(post.StatusPostStatusTrashed)]
	assert.False(t, hasTrashed, "回收站状态不应出现在分布里")
}

// TestStatsRepo_GetCommentStatusDistribution 验证状态分组与 NULL→已通过归并。
func TestStatsRepo_GetCommentStatusDistribution(t *testing.T) {
	repo, db, cleanup := newTestStatsRepo(t)
	defer cleanup()

	ctx := viewerCtx(1, 1)
	createTestComment(t, db, 1, false) // NULL 状态 → 归并为已通过
	_, err := db.Comment.Create().
		SetTenantID(1).
		SetContent("approved").
		SetStatus(comment.StatusStatusApproved).
		Save(ctx)
	require.NoError(t, err)
	_, err = db.Comment.Create().
		SetTenantID(1).
		SetContent("pending").
		SetStatus(comment.StatusStatusPending).
		Save(ctx)
	require.NoError(t, err)
	createTestComment(t, db, 1, true) // SPAM 保留在分布里

	rows, err := repo.GetCommentStatusDistribution(ctx)
	require.NoError(t, err)

	byStatus := make(map[string]int64, len(rows))
	for _, row := range rows {
		byStatus[row.Status] = row.Count
	}
	assert.Equal(t, int64(2), byStatus[string(comment.StatusStatusApproved)], "已通过应为 2(NULL 归并)")
	assert.Equal(t, int64(1), byStatus[string(comment.StatusStatusPending)], "待审应为 1")
	assert.Equal(t, int64(1), byStatus[string(comment.StatusStatusSpam)], "垃圾应为 1(分布保留)")
}

// TestStatsRepo_GetMediaAssetTypeDistribution 验证类型分组与 NULL→其他归并。
func TestStatsRepo_GetMediaAssetTypeDistribution(t *testing.T) {
	repo, db, cleanup := newTestStatsRepo(t)
	defer cleanup()

	ctx := viewerCtx(1, 1)
	_, err := db.MediaAsset.Create().
		SetTenantID(1).
		SetType(mediaasset.TypeAssetTypeImage).
		SetSize(100).
		Save(ctx)
	require.NoError(t, err)
	_, err = db.MediaAsset.Create().
		SetTenantID(1).
		SetType(mediaasset.TypeAssetTypeImage).
		SetSize(200).
		Save(ctx)
	require.NoError(t, err)
	// NULL 类型 → 归并为其他
	createTestMediaAsset(t, db, 1, 50)
	createTestMediaAsset(t, db, 2, 999)

	rows, err := repo.GetMediaAssetTypeDistribution(ctx)
	require.NoError(t, err)

	byType := make(map[string]int64, len(rows))
	for _, row := range rows {
		byType[row.Status] = row.Count
	}
	assert.Equal(t, int64(2), byType[string(mediaasset.TypeAssetTypeImage)], "图片应为 2")
	assert.Equal(t, int64(1), byType[string(mediaasset.TypeAssetTypeOther)], "其他应为 1(NULL 归并)")
	assert.Len(t, rows, 2, "租户 2 的数据不应出现")
}

// TestStatsRepo_TenantIsolation 验证租户 2 的 viewer 看不到租户 1 的数据。
func TestStatsRepo_TenantIsolation(t *testing.T) {
	repo, db, cleanup := newTestStatsRepo(t)
	defer cleanup()

	createTestUser(t, db, 1, "tenant1-user")
	pid := createTestPost(t, db, 1)
	createTestCounter(t, db, 1, pid, interactionV1.CounterMetric_COUNTER_METRIC_LIKE, 5)

	ctx2 := viewerCtx(2, 1)
	counts, err := repo.GetOverview(ctx2)
	require.NoError(t, err)
	assert.Equal(t, int64(0), counts.UserCount, "租户 2 应看不到租户 1 的用户")
	assert.Equal(t, int64(0), counts.PostCount, "租户 2 应看不到租户 1 的帖子")
	assert.Equal(t, int64(0), counts.InteractionCount, "租户 2 应看不到租户 1 的互动计数")

	likes, watches, _, err := repo.GetInteractionTotals(ctx2)
	require.NoError(t, err)
	assert.Equal(t, int64(0), likes)
	assert.Equal(t, int64(0), watches)
}
