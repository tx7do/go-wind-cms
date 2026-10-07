package service

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/go-crud/viewer"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"

	aiV1 "go-wind-cms/api/gen/go/ai/service/v1"
	"go-wind-cms/app/core/service/internal/data"
)

// AiUsageLogService 用量流水查询（配额记账事实源的只读视图）。
// 口径：租户用户按本租户，平台管理员（tenantId=0）见全量。
// 本仓无套餐配额体系，Summary 的 quota 字段恒 false/0。
type AiUsageLogService struct {
	aiV1.UnimplementedAiUsageLogServiceServer
	log  *log.Helper
	repo *data.AiUsageLogRepo
}

func NewAiUsageLogService(ctx *bootstrap.Context, repo *data.AiUsageLogRepo) *AiUsageLogService {
	return &AiUsageLogService{
		log:  ctx.NewLoggerHelper("ai_usage_log/service/core-service"),
		repo: repo,
	}
}

func (s *AiUsageLogService) List(ctx context.Context, req *paginationV1.PagingRequest) (*aiV1.ListAiUsageLogResponse, error) {
	tenantId, _ := s.operatorTenant(ctx)
	return s.repo.List(ctx, req, tenantId)
}

func (s *AiUsageLogService) Count(ctx context.Context, req *paginationV1.PagingRequest) (*aiV1.CountAiUsageLogResponse, error) {
	tenantId, _ := s.operatorTenant(ctx)
	return s.repo.Count(ctx, req, tenantId)
}

// GetUsageSummary 当月用量汇总：tokens/调用次数。
func (s *AiUsageLogService) GetUsageSummary(ctx context.Context, _ *emptypb.Empty) (*aiV1.UsageSummaryResponse, error) {
	tenantId, _ := s.operatorTenant(ctx)
	monthStart := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Local)
	tokens, calls, err := s.repo.MonthStats(ctx, tenantId, monthStart)
	if err != nil {
		return nil, err
	}

	return &aiV1.UsageSummaryResponse{
		MonthTokens: uint32(tokens),
		MonthCalls:  uint32(calls),
		// 本仓无套餐配额：quota_configured 恒 false（前端据此不渲染配额进度条）
	}, nil
}

// operatorTenant 取调用者租户（viewer 上下文；平台用户 tenantId=0 见全量）。
func (s *AiUsageLogService) operatorTenant(ctx context.Context) (uint32, bool) {
	vc, exist := viewer.FromContext(ctx)
	if !exist || vc == nil {
		return 0, false
	}
	return uint32(vc.TenantID()), true
}
