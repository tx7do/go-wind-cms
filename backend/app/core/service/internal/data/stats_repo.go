package data

import (
	"context"
	"fmt"
	"time"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	entCrud "github.com/tx7do/go-crud/entgo"

	interactionV1 "go-wind-cms/api/gen/go/interaction/service/v1"
	statsV1 "go-wind-cms/api/gen/go/stats/service/v1"

	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/comment"
	"go-wind-cms/app/core/service/internal/data/ent/interactioncounter"
	"go-wind-cms/app/core/service/internal/data/ent/loginauditlog"
	"go-wind-cms/app/core/service/internal/data/ent/post"
	"go-wind-cms/app/core/service/internal/data/ent/postlike"
	"go-wind-cms/app/core/service/internal/data/ent/posttranslation"
	"go-wind-cms/app/core/service/internal/data/ent/user"
)

// localTZOffsetMinutes 返回服务器本地时区相对 UTC 的偏移(分钟)。
// 统计日按本地时区对齐:DATE() 对 timestamptz 折算 UTC 取日,若不平移,
// 本地凌晨 0-N 点(N=UTC+时差)的数据会被归入前一天。
func localTZOffsetMinutes() int {
	_, offsetSec := time.Now().Zone()
	return offsetSec / 60
}

// dailyDateExpr 生成按天截断的日期表达式(YYYY-MM-DD 文本),按方言平移本地时区。
// CAST 成 TEXT 保证 postgres/sqlite 双端 Scan 到 string 的行为一致。
func dailyDateExpr(d string) string {
	mins := localTZOffsetMinutes()
	sign := "+"
	if mins < 0 {
		sign = "-"
		mins = -mins
	}
	mod := fmt.Sprintf("%s%d minutes", sign, mins)
	if d == dialect.Postgres {
		return fmt.Sprintf("CAST(DATE(created_at + INTERVAL '%s') AS TEXT)", mod)
	}
	return fmt.Sprintf("CAST(DATE(created_at, '%s') AS TEXT)", mod)
}

// counter 是各 ent 查询共有的计数能力（方法级 duck typing，避免泛型样板）。
type counter interface {
	Count(context.Context) (int, error)
}

// StatsRepo 是后台分析页的跨域只读聚合仓储。
//
// 设计要点：
//   - 纯只读，无事务、无副作用；不触碰 ledger/counter 写入路径。
//   - 租户范围由 viewer 上下文决定：maybeTenantFromViewer 有租户则过滤本租户，
//     平台/系统上下文（无租户）查全局——与 interaction_repo 的租户语义一致。
//   - 按天聚合走 ent Modify + 原生 SQL（仓库首个 GroupBy 用法），避免全表
//     拉回内存；Go 侧仅做缺日补零。GROUP BY 引用 SELECT 别名 day
//     （sqlite/postgres 均支持），因为 Selector.GroupBy 会对参数加引号，
//     无法直接传表达式。
//   - 统计口径：帖子排除回收站（TRASHED），评论排除垃圾（SPAM），用户全量；
//     登录活跃仅统计 action_type=LOGIN 的 SUCCESS/FAILED 行。
type StatsRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper
}

func NewStatsRepo(ctx *bootstrap.Context, entClient *entCrud.EntClient[*ent.Client]) *StatsRepo {
	return &StatsRepo{
		entClient: entClient,
		log:       ctx.NewLoggerHelper("stats/repo/core-service"),
	}
}

// OverviewCounts 概览卡片的总量 + 近 7 天增量。
type OverviewCounts struct {
	UserCount        int64
	PostCount        int64
	CommentCount     int64
	InteractionCount int64 // 点赞 + 收藏计数累计之和
	NewUserCountWeek int64
	NewPostCountWeek int64
	NewCommentCount  int64
	NewLikeCountWeek int64 // 近 7 天 post_like ledger 行数
}

// DailyCountRow 单日计数点。
type DailyCountRow struct {
	Date  string
	Value int64
}

// TopLikedPost 点赞 TOP 榜条目。
type TopLikedPost struct {
	TargetID  uint32
	Title     string
	LikeCount int64
}

// dailyRow 按天聚合 SQL 的 Scan 目标（json tag 与列别名对应）。
type dailyRow struct {
	Day   string `json:"day"`
	Value int64  `json:"value"`
}

// trendStart 返回覆盖最近 days 天(含今天)的起始本地零点时刻。
// 统计日按服务器本地时区对齐,与 dailyDateExpr 的 SQL 平移口径一致。
func trendStart(days int32) time.Time {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return start.AddDate(0, 0, -(int(days) - 1))
}

// dailyModify 生成按天聚合计数的 Modify 片段(SELECT day, COUNT(*) GROUP BY day)。
func dailyModify() func(s *sql.Selector) {
	return func(s *sql.Selector) {
		s.Select(
			sql.As(dailyDateExpr(s.Dialect()), "day"),
			sql.As(sql.Count("*"), "value"),
		)
		s.GroupBy("day")
	}
}

// fillDailySeries 把聚合行补齐为从起始日到今天、每日一点的连续序列（无数据补 0）。
func fillDailySeries(rows []dailyRow, days int32) []DailyCountRow {
	byDate := make(map[string]int64, len(rows))
	for _, row := range rows {
		byDate[row.Day] = row.Value
	}
	start := trendStart(days)
	out := make([]DailyCountRow, 0, days)
	for i := int32(0); i < days; i++ {
		d := start.AddDate(0, 0, int(i)).Format("2006-01-02")
		out = append(out, DailyCountRow{Date: d, Value: byDate[d]})
	}
	return out
}

// countOf 执行计数并转为 int64。
func countOf(ctx context.Context, q counter) (int64, error) {
	n, err := q.Count(ctx)
	return int64(n), err
}

// GetOverview 概览计数。7 天增量按 created_at >= 7 天前零点计算。
func (r *StatsRepo) GetOverview(ctx context.Context) (*OverviewCounts, error) {
	tid, hasTenant := maybeTenantFromViewer(ctx)
	client := r.entClient.Client()
	since := trendStart(7)
	counts := &OverviewCounts{}

	// 用户全量
	userQ := client.User.Query()
	user7Q := client.User.Query().Where(user.CreatedAtGTE(since))
	if hasTenant {
		userQ = userQ.Where(user.TenantIDEQ(tid))
		user7Q = user7Q.Where(user.TenantIDEQ(tid))
	}

	postQ := client.Post.Query().Where(post.StatusNEQ(post.StatusPostStatusTrashed))
	post7Q := client.Post.Query().Where(post.StatusNEQ(post.StatusPostStatusTrashed), post.CreatedAtGTE(since))
	if hasTenant {
		postQ = postQ.Where(post.TenantIDEQ(tid))
		post7Q = post7Q.Where(post.TenantIDEQ(tid))
	}

	commentQ := client.Comment.Query().Where(comment.Or(comment.StatusIsNil(), comment.StatusNEQ(comment.StatusStatusSpam)))
	comment7Q := client.Comment.Query().Where(comment.Or(comment.StatusIsNil(), comment.StatusNEQ(comment.StatusStatusSpam)), comment.CreatedAtGTE(since))
	if hasTenant {
		commentQ = commentQ.Where(comment.TenantIDEQ(tid))
		comment7Q = comment7Q.Where(comment.TenantIDEQ(tid))
	}

	like7Q := client.PostLike.Query().Where(postlike.CreatedAtGTE(since))
	if hasTenant {
		like7Q = like7Q.Where(postlike.TenantIDEQ(tid))
	}

	var err error
	if counts.UserCount, err = countOf(ctx, userQ); err != nil {
		r.log.Errorf("count users failed: %s", err.Error())
		return nil, statsV1.ErrorInternalServerError("count users failed")
	}
	if counts.NewUserCountWeek, err = countOf(ctx, user7Q); err != nil {
		r.log.Errorf("count new users failed: %s", err.Error())
		return nil, statsV1.ErrorInternalServerError("count new users failed")
	}
	if counts.PostCount, err = countOf(ctx, postQ); err != nil {
		r.log.Errorf("count posts failed: %s", err.Error())
		return nil, statsV1.ErrorInternalServerError("count posts failed")
	}
	if counts.NewPostCountWeek, err = countOf(ctx, post7Q); err != nil {
		r.log.Errorf("count new posts failed: %s", err.Error())
		return nil, statsV1.ErrorInternalServerError("count new posts failed")
	}
	if counts.CommentCount, err = countOf(ctx, commentQ); err != nil {
		r.log.Errorf("count comments failed: %s", err.Error())
		return nil, statsV1.ErrorInternalServerError("count comments failed")
	}
	if counts.NewCommentCount, err = countOf(ctx, comment7Q); err != nil {
		r.log.Errorf("count new comments failed: %s", err.Error())
		return nil, statsV1.ErrorInternalServerError("count new comments failed")
	}
	if counts.NewLikeCountWeek, err = countOf(ctx, like7Q); err != nil {
		r.log.Errorf("count new likes failed: %s", err.Error())
		return nil, statsV1.ErrorInternalServerError("count new likes failed")
	}

	likes, watches, err := r.GetInteractionTotals(ctx)
	if err != nil {
		return nil, err
	}
	counts.InteractionCount = likes + watches

	return counts, nil
}

// GetDailyTrend 近 N 天每日新增用户/帖子/评论序列（缺日补 0）。
func (r *StatsRepo) GetDailyTrend(ctx context.Context, days int32) (users, posts, comments []DailyCountRow, err error) {
	tid, hasTenant := maybeTenantFromViewer(ctx)
	client := r.entClient.Client()
	start := trendStart(days)

	userQ := client.User.Query().Where(user.CreatedAtGTE(start))
	postQ := client.Post.Query().Where(post.StatusNEQ(post.StatusPostStatusTrashed), post.CreatedAtGTE(start))
	commentQ := client.Comment.Query().Where(comment.Or(comment.StatusIsNil(), comment.StatusNEQ(comment.StatusStatusSpam)), comment.CreatedAtGTE(start))
	if hasTenant {
		userQ = userQ.Where(user.TenantIDEQ(tid))
		postQ = postQ.Where(post.TenantIDEQ(tid))
		commentQ = commentQ.Where(comment.TenantIDEQ(tid))
	}

	var userRows, postRows, commentRows []dailyRow

	if err = userQ.Modify(dailyModify()).Scan(ctx, &userRows); err != nil {
		r.log.Errorf("aggregate daily users failed: %s", err.Error())
		return nil, nil, nil, statsV1.ErrorInternalServerError("aggregate daily users failed")
	}
	if err = postQ.Modify(dailyModify()).Scan(ctx, &postRows); err != nil {
		r.log.Errorf("aggregate daily posts failed: %s", err.Error())
		return nil, nil, nil, statsV1.ErrorInternalServerError("aggregate daily posts failed")
	}
	if err = commentQ.Modify(dailyModify()).Scan(ctx, &commentRows); err != nil {
		r.log.Errorf("aggregate daily comments failed: %s", err.Error())
		return nil, nil, nil, statsV1.ErrorInternalServerError("aggregate daily comments failed")
	}

	return fillDailySeries(userRows, days), fillDailySeries(postRows, days), fillDailySeries(commentRows, days), nil
}

// GetInteractionTotals 各 metric 的计数累计总量。
func (r *StatsRepo) GetInteractionTotals(ctx context.Context) (likes, watches int64, err error) {
	tid, hasTenant := maybeTenantFromViewer(ctx)

	query := r.entClient.Client().InteractionCounter.Query()
	if hasTenant {
		query = query.Where(interactioncounter.TenantIDEQ(tid))
	}

	var rows []struct {
		Metric uint8 `json:"metric"`
		Total  int64 `json:"total"`
	}
	err = query.
		Modify(func(s *sql.Selector) {
			s.Select(sql.As("metric", "metric"), sql.As(sql.Sum("count"), "total"))
			s.GroupBy("metric")
		}).
		Scan(ctx, &rows)
	if err != nil {
		r.log.Errorf("aggregate interaction totals failed: %s", err.Error())
		return 0, 0, statsV1.ErrorInternalServerError("aggregate interaction totals failed")
	}

	for _, row := range rows {
		switch row.Metric {
		case uint8(interactionV1.CounterMetric_COUNTER_METRIC_LIKE):
			likes = row.Total
		case uint8(interactionV1.CounterMetric_COUNTER_METRIC_WATCH):
			watches = row.Total
		}
	}
	return likes, watches, nil
}

// GetTopLikedPosts 帖子点赞 TOP N（降序）。回收站/已不存在的帖子被剔除。
func (r *StatsRepo) GetTopLikedPosts(ctx context.Context, topN int32) ([]TopLikedPost, error) {
	tid, hasTenant := maybeTenantFromViewer(ctx)
	client := r.entClient.Client()

	query := client.InteractionCounter.Query().Where(
		interactioncounter.TargetTypeEQ(uint8(interactionV1.TargetType_TARGET_TYPE_POST)),
		interactioncounter.MetricEQ(uint8(interactionV1.CounterMetric_COUNTER_METRIC_LIKE)),
	)
	if hasTenant {
		query = query.Where(interactioncounter.TenantIDEQ(tid))
	}

	var rows []struct {
		TargetID uint32 `json:"target_id"`
		Total    int64  `json:"total"`
	}
	err := query.
		Modify(func(s *sql.Selector) {
			s.Select(sql.As("target_id", "target_id"), sql.As(sql.Sum("count"), "total"))
			s.GroupBy("target_id")
			s.OrderBy(sql.Desc("total"))
			s.Limit(int(topN))
		}).
		Scan(ctx, &rows)
	if err != nil {
		r.log.Errorf("aggregate top liked posts failed: %s", err.Error())
		return nil, statsV1.ErrorInternalServerError("aggregate top liked posts failed")
	}
	if len(rows) == 0 {
		return nil, nil
	}

	// 剔除回收站/已不存在的帖子
	ids := make([]uint32, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.TargetID)
	}
	existQ := client.Post.Query().Where(post.IDIn(ids...), post.StatusNEQ(post.StatusPostStatusTrashed))
	if hasTenant {
		existQ = existQ.Where(post.TenantIDEQ(tid))
	}
	existing, err := existQ.IDs(ctx)
	if err != nil {
		r.log.Errorf("query existing posts failed: %s", err.Error())
		return nil, statsV1.ErrorInternalServerError("query existing posts failed")
	}
	existingSet := make(map[uint32]struct{}, len(existing))
	for _, id := range existing {
		existingSet[id] = struct{}{}
	}

	// 标题取任一语言翻译的首条非空值（榜单展示用途，不做语言精确匹配）
	titleQ := client.PostTranslation.Query().Where(posttranslation.PostIDIn(existing...))
	if hasTenant {
		titleQ = titleQ.Where(posttranslation.TenantIDEQ(tid))
	}
	titleRows, err := titleQ.All(ctx)
	if err != nil {
		r.log.Errorf("query post titles failed: %s", err.Error())
		return nil, statsV1.ErrorInternalServerError("query post titles failed")
	}
	titles := make(map[uint32]string, len(titleRows))
	for _, t := range titleRows {
		if t.PostID == nil || t.Title == nil || *t.Title == "" {
			continue
		}
		if _, ok := titles[*t.PostID]; !ok {
			titles[*t.PostID] = *t.Title
		}
	}

	out := make([]TopLikedPost, 0, len(rows))
	for _, row := range rows {
		if _, ok := existingSet[row.TargetID]; !ok {
			continue
		}
		out = append(out, TopLikedPost{
			TargetID:  row.TargetID,
			Title:     titles[row.TargetID],
			LikeCount: row.Total,
		})
	}
	return out, nil
}

// GetLoginActivity 近 N 天每日登录成功/失败次数。
func (r *StatsRepo) GetLoginActivity(ctx context.Context, days int32) (success, failed []DailyCountRow, err error) {
	tid, hasTenant := maybeTenantFromViewer(ctx)
	client := r.entClient.Client()
	start := trendStart(days)

	loginQuery := func() *ent.LoginAuditLogQuery {
		q := client.LoginAuditLog.Query().Where(
			loginauditlog.ActionTypeEQ(loginauditlog.ActionTypeLogin),
			loginauditlog.CreatedAtGTE(start),
		)
		if hasTenant {
			q = q.Where(loginauditlog.TenantIDEQ(tid))
		}
		return q
	}

	var successRows, failedRows []dailyRow
	if err = loginQuery().
		Where(loginauditlog.StatusEQ(loginauditlog.StatusSuccess)).
		Modify(dailyModify()).
		Scan(ctx, &successRows); err != nil {
		r.log.Errorf("aggregate daily login success failed: %s", err.Error())
		return nil, nil, statsV1.ErrorInternalServerError("aggregate daily login success failed")
	}
	if err = loginQuery().
		Where(loginauditlog.StatusEQ(loginauditlog.StatusFailed)).
		Modify(dailyModify()).
		Scan(ctx, &failedRows); err != nil {
		r.log.Errorf("aggregate daily login failed failed: %s", err.Error())
		return nil, nil, statsV1.ErrorInternalServerError("aggregate daily login failed failed")
	}

	return fillDailySeries(successRows, days), fillDailySeries(failedRows, days), nil
}
