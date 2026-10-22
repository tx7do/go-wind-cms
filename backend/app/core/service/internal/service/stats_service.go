package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	"go-wind-cms/app/core/service/internal/data"

	statsV1 "go-wind-cms/api/gen/go/stats/service/v1"
)

// 趋势天数与 TOP N 的 clamp 边界。
const (
	defaultTrendDays = 30
	maxTrendDays     = 90
	defaultTopN      = 10
	maxTopN          = 50
)

// StatsService 是后台分析页的只读统计服务。
//
// 全部查询无副作用，不做 platform-only gate：admin BFF 已完成 token 鉴权，
// 路由授权含 tenant_manager（租户管理员看本租户），租户范围由 viewer
// 上下文在 repo 层过滤决定。
type StatsService struct {
	statsV1.UnimplementedStatsServiceServer

	statsRepo *data.StatsRepo
	log       *log.Helper
}

func NewStatsService(ctx *bootstrap.Context, statsRepo *data.StatsRepo) *StatsService {
	return &StatsService{
		statsRepo: statsRepo,
		log:       ctx.NewLoggerHelper("stats/service/core-service"),
	}
}

// GetDashboardOverview 分析页概览卡片：各实体总量 + 近 7 天新增量。
func (s *StatsService) GetDashboardOverview(ctx context.Context, _ *statsV1.GetDashboardOverviewRequest) (*statsV1.GetDashboardOverviewResponse, error) {
	counts, err := s.statsRepo.GetOverview(ctx)
	if err != nil {
		return nil, err
	}
	return &statsV1.GetDashboardOverviewResponse{
		UserCount:                counts.UserCount,
		PostCount:                counts.PostCount,
		CommentCount:             counts.CommentCount,
		InteractionCount:         counts.InteractionCount,
		NewUserCountWeek:         counts.NewUserCountWeek,
		NewPostCountWeek:         counts.NewPostCountWeek,
		NewCommentCountWeek:      counts.NewCommentCount,
		NewLikeCountWeek:         counts.NewLikeCountWeek,
		MediaAssetCount:          counts.MediaAssetCount,
		MediaStorageBytes:        counts.MediaStorageBytes,
		NewMediaCountWeek:        counts.NewMediaCountWeek,
		NewMediaStorageBytesWeek: counts.NewMediaStorageWeek,
		PendingCommentCount:      counts.PendingCommentCount,
		DraftPostCount:           counts.DraftPostCount,
		ScheduledPostCount:       counts.ScheduledPostCount,
	}, nil
}

// GetContentTrend 内容增长趋势：近 N 天按天聚合的新增用户/帖子/评论序列。
func (s *StatsService) GetContentTrend(ctx context.Context, req *statsV1.GetContentTrendRequest) (*statsV1.GetContentTrendResponse, error) {
	days := clampDays(req.GetDays())

	users, posts, comments, err := s.statsRepo.GetDailyTrend(ctx, days)
	if err != nil {
		return nil, err
	}
	return &statsV1.GetContentTrendResponse{
		Users:    toDailyCounts(users),
		Posts:    toDailyCounts(posts),
		Comments: toDailyCounts(comments),
	}, nil
}

// GetInteractionStats 互动统计：点赞 TOP N 内容 + 各 metric 累计总量分布。
func (s *StatsService) GetInteractionStats(ctx context.Context, req *statsV1.GetInteractionStatsRequest) (*statsV1.GetInteractionStatsResponse, error) {
	topN := clampTopN(req.GetTopN())

	top, err := s.statsRepo.GetTopLikedPosts(ctx, topN)
	if err != nil {
		return nil, err
	}
	topViewed, err := s.statsRepo.GetTopViewedPosts(ctx, topN)
	if err != nil {
		return nil, err
	}
	likes, watches, views, err := s.statsRepo.GetInteractionTotals(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]*statsV1.InteractionTopItem, 0, len(top))
	for _, item := range top {
		items = append(items, &statsV1.InteractionTopItem{
			TargetId:  item.TargetID,
			Title:     item.Title,
			LikeCount: item.LikeCount,
		})
	}
	viewedItems := make([]*statsV1.InteractionTopItem, 0, len(topViewed))
	for _, item := range topViewed {
		viewedItems = append(viewedItems, &statsV1.InteractionTopItem{
			TargetId:  item.TargetID,
			Title:     item.Title,
			LikeCount: item.LikeCount, // 复用字段：此处语义为浏览计数
		})
	}
	return &statsV1.GetInteractionStatsResponse{
		TopLikedPosts:  items,
		TotalLikes:     likes,
		TotalWatches:   watches,
		TopViewedPosts: viewedItems,
		TotalViews:     views,
	}, nil
}

// GetLoginActivity 登录活跃趋势：近 N 天每日登录成功/失败次数 + 活跃用户数。
func (s *StatsService) GetLoginActivity(ctx context.Context, req *statsV1.GetLoginActivityRequest) (*statsV1.GetLoginActivityResponse, error) {
	days := clampDays(req.GetDays())

	success, failed, activeUsers, err := s.statsRepo.GetLoginActivity(ctx, days)
	if err != nil {
		return nil, err
	}
	return &statsV1.GetLoginActivityResponse{
		Success:     toDailyCounts(success),
		Failed:      toDailyCounts(failed),
		ActiveUsers: toDailyCounts(activeUsers),
	}, nil
}

// GetAiUsageTrend AI 用量趋势：近 N 天每日模型调用次数与 token 消耗。
func (s *StatsService) GetAiUsageTrend(ctx context.Context, req *statsV1.TrendRequest) (*statsV1.GetAiUsageTrendResponse, error) {
	days := clampDays(req.GetDays())

	calls, tokens, err := s.statsRepo.GetAiUsageTrend(ctx, days)
	if err != nil {
		return nil, err
	}
	return &statsV1.GetAiUsageTrendResponse{
		Calls:  toDailyCounts(calls),
		Tokens: toDailyCounts(tokens),
	}, nil
}

// GetApiCallTrend API 调用量趋势：近 N 天每日网关请求次数。
func (s *StatsService) GetApiCallTrend(ctx context.Context, req *statsV1.TrendRequest) (*statsV1.GetApiCallTrendResponse, error) {
	days := clampDays(req.GetDays())

	calls, err := s.statsRepo.GetApiCallTrend(ctx, days)
	if err != nil {
		return nil, err
	}
	return &statsV1.GetApiCallTrendResponse{
		Calls: toDailyCounts(calls),
	}, nil
}

// GetPostStatusDistribution 帖子状态分布（排除回收站）。
func (s *StatsService) GetPostStatusDistribution(ctx context.Context, _ *statsV1.GetPostStatusDistributionRequest) (*statsV1.GetPostStatusDistributionResponse, error) {
	rows, err := s.statsRepo.GetPostStatusDistribution(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]*statsV1.StatusCount, 0, len(rows))
	for _, row := range rows {
		items = append(items, &statsV1.StatusCount{Status: row.Status, Count: row.Count})
	}
	return &statsV1.GetPostStatusDistributionResponse{Items: items}, nil
}

// GetCommentStatusDistribution 评论状态分布（NULL 状态按已通过计）。
func (s *StatsService) GetCommentStatusDistribution(ctx context.Context, _ *statsV1.GetCommentStatusDistributionRequest) (*statsV1.GetCommentStatusDistributionResponse, error) {
	rows, err := s.statsRepo.GetCommentStatusDistribution(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]*statsV1.StatusCount, 0, len(rows))
	for _, row := range rows {
		items = append(items, &statsV1.StatusCount{Status: row.Status, Count: row.Count})
	}
	return &statsV1.GetCommentStatusDistributionResponse{Items: items}, nil
}

// GetMediaAssetTypeDistribution 媒体类型分布（NULL 类型按其他计）。
func (s *StatsService) GetMediaAssetTypeDistribution(ctx context.Context, _ *statsV1.GetMediaAssetTypeDistributionRequest) (*statsV1.GetMediaAssetTypeDistributionResponse, error) {
	rows, err := s.statsRepo.GetMediaAssetTypeDistribution(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]*statsV1.StatusCount, 0, len(rows))
	for _, row := range rows {
		items = append(items, &statsV1.StatusCount{Status: row.Status, Count: row.Count})
	}
	return &statsV1.GetMediaAssetTypeDistributionResponse{Items: items}, nil
}

// clampDays 趋势天数：缺省/非正数取 30，上限 90。
func clampDays(days int32) int32 {
	if days <= 0 {
		return defaultTrendDays
	}
	if days > maxTrendDays {
		return maxTrendDays
	}
	return days
}

// clampTopN TOP N：缺省/非正数取 10，上限 50。
func clampTopN(topN int32) int32 {
	if topN <= 0 {
		return defaultTopN
	}
	if topN > maxTopN {
		return maxTopN
	}
	return topN
}

// toDailyCounts repo 行 → proto 序列。
func toDailyCounts(rows []data.DailyCountRow) []*statsV1.DailyCount {
	out := make([]*statsV1.DailyCount, 0, len(rows))
	for _, row := range rows {
		out = append(out, &statsV1.DailyCount{Date: row.Date, Value: row.Value})
	}
	return out
}
