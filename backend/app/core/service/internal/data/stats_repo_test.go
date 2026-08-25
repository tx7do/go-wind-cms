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

	pid := createTestPost(t, db, 1)
	_, err = db.PostLike.Create().
		SetTenantID(1).
		SetUserID(100).
		SetPostID(pid).
		Save(ctx)
	require.NoError(t, err)

	createTestCounter(t, db, 1, pid, interactionV1.CounterMetric_COUNTER_METRIC_LIKE, 5)
	createTestCounter(t, db, 1, pid, interactionV1.CounterMetric_COUNTER_METRIC_WATCH, 3)

	// 其他租户数据不应出现
	createTestUser(t, db, 2, "other-tenant-user")

	counts, err := repo.GetOverview(ctx)
	require.NoError(t, err)

	assert.Equal(t, int64(2), counts.UserCount, "用户总量应为 2(不含其他租户)")
	assert.Equal(t, int64(2), counts.NewUserCountWeek, "近 7 天新增用户应为 2")
	assert.Equal(t, int64(3), counts.PostCount, "帖子总量应为 3(排除回收站 1 篇)")
	assert.Equal(t, int64(3), counts.NewPostCountWeek, "近 7 天新增帖子应为 3")
	assert.Equal(t, int64(1), counts.CommentCount, "评论总量应为 1(排除 SPAM)")
	assert.Equal(t, int64(1), counts.NewCommentCount, "近 7 天新增评论应为 1")
	assert.Equal(t, int64(8), counts.InteractionCount, "互动总量应为 5+3")
	assert.Equal(t, int64(1), counts.NewLikeCountWeek, "近 7 天新增点赞应为 1")
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

// TestStatsRepo_GetLoginActivity 验证登录成功/失败按天聚合,以及非 LOGIN 动作与 PARTIAL 状态不计入。
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

	success, failed, err := repo.GetLoginActivity(viewerCtx(1, 1), 7)
	require.NoError(t, err)

	assert.Len(t, success, 7)
	assert.Len(t, failed, 7)
	assert.Equal(t, int64(2), success[6].Value, "今日登录成功应为 2")
	assert.Equal(t, int64(1), failed[6].Value, "今日登录失败应为 1")
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

	likes, watches, err := repo.GetInteractionTotals(ctx2)
	require.NoError(t, err)
	assert.Equal(t, int64(0), likes)
	assert.Equal(t, int64(0), watches)
}
