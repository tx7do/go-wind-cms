package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
)

// ReviewLogService 发布审批留痕 BFF 面：代理 core ReviewLogService。
type ReviewLogService struct {
	adminV1.ReviewLogServiceHTTPServer

	log *log.Helper

	reviewLogServiceClient contentV1.ReviewLogServiceClient
}

func NewReviewLogService(
	ctx *bootstrap.Context,
	reviewLogServiceClient contentV1.ReviewLogServiceClient,
) *ReviewLogService {
	return &ReviewLogService{
		log:                    ctx.NewLoggerHelper("review-log/service/admin-service"),
		reviewLogServiceClient: reviewLogServiceClient,
	}
}

func (s *ReviewLogService) List(ctx context.Context, req *contentV1.ListReviewLogsRequest) (*contentV1.ListReviewLogsResponse, error) {
	return s.reviewLogServiceClient.List(ctx, req)
}
