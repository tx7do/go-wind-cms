package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	aiV1 "go-wind-cms/api/gen/go/ai/service/v1"
	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"

	"go-wind-cms/pkg/middleware/auth"
)

// AiUsageLogService 用量流水查询（BFF 薄包装；租户口径由 core 侧裁定）。
type AiUsageLogService struct {
	adminV1.AiUsageLogServiceHTTPServer

	aiUsageLogServiceClient aiV1.AiUsageLogServiceClient
	log                     *log.Helper
}

func NewAiUsageLogService(ctx *bootstrap.Context, aiUsageLogServiceClient aiV1.AiUsageLogServiceClient) *AiUsageLogService {
	return &AiUsageLogService{
		log:                     ctx.NewLoggerHelper("ai_usage_log/service/admin-service"),
		aiUsageLogServiceClient: aiUsageLogServiceClient,
	}
}

func (s *AiUsageLogService) List(ctx context.Context, req *paginationV1.PagingRequest) (*aiV1.ListAiUsageLogResponse, error) {
	return s.aiUsageLogServiceClient.List(ctx, req)
}

func (s *AiUsageLogService) Count(ctx context.Context, req *paginationV1.PagingRequest) (*aiV1.CountAiUsageLogResponse, error) {
	return s.aiUsageLogServiceClient.Count(ctx, req)
}

// GetUsageSummary 当月用量汇总。本仓无套餐配额：quota 字段恒 false/0。
func (s *AiUsageLogService) GetUsageSummary(ctx context.Context, _ *emptypb.Empty) (*aiV1.UsageSummaryResponse, error) {
	if _, err := auth.FromContext(ctx); err != nil {
		return nil, err
	}
	return s.aiUsageLogServiceClient.GetUsageSummary(ctx, &emptypb.Empty{})
}
