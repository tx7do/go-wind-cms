package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	statsV1 "go-wind-cms/api/gen/go/stats/service/v1"
)

// StatsService 是 admin BFF 对 core StatsService 的纯透传转发器，
// 为后台分析页提供只读统计数据。鉴权由 REST 中间件（token + authz）完成，
// 租户范围经 viewer 上下文透传至 core 侧过滤。
type StatsService struct {
	adminV1.StatsServiceHTTPServer

	statsServiceClient statsV1.StatsServiceClient
	log                *log.Helper
}

func NewStatsService(
	ctx *bootstrap.Context,
	statsServiceClient statsV1.StatsServiceClient,
) *StatsService {
	return &StatsService{
		statsServiceClient: statsServiceClient,
		log:                ctx.NewLoggerHelper("stats/service/admin-service"),
	}
}

func (s *StatsService) GetDashboardOverview(ctx context.Context, req *statsV1.GetDashboardOverviewRequest) (*statsV1.GetDashboardOverviewResponse, error) {
	return s.statsServiceClient.GetDashboardOverview(ctx, req)
}

func (s *StatsService) GetContentTrend(ctx context.Context, req *statsV1.GetContentTrendRequest) (*statsV1.GetContentTrendResponse, error) {
	return s.statsServiceClient.GetContentTrend(ctx, req)
}

func (s *StatsService) GetInteractionStats(ctx context.Context, req *statsV1.GetInteractionStatsRequest) (*statsV1.GetInteractionStatsResponse, error) {
	return s.statsServiceClient.GetInteractionStats(ctx, req)
}

func (s *StatsService) GetLoginActivity(ctx context.Context, req *statsV1.GetLoginActivityRequest) (*statsV1.GetLoginActivityResponse, error) {
	return s.statsServiceClient.GetLoginActivity(ctx, req)
}
