package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	contentV1 "go-wind-cms/api/gen/go/content/service/v1"

	"go-wind-cms/app/core/service/internal/data"
)

// ReviewLogService 发布审批留痕查询（admin 面）。
type ReviewLogService struct {
	contentV1.UnimplementedReviewLogServiceServer

	repo *data.ReviewLogRepo
	log  *log.Helper
}

func NewReviewLogService(ctx *bootstrap.Context, repo *data.ReviewLogRepo) *ReviewLogService {
	return &ReviewLogService{
		log:  ctx.NewLoggerHelper("review-log/service/core-service"),
		repo: repo,
	}
}

// List 按实体查询审批记录（时间正序）。
func (s *ReviewLogService) List(ctx context.Context, req *contentV1.ListReviewLogsRequest) (*contentV1.ListReviewLogsResponse, error) {
	if req == nil || req.GetEntityType() == "" || req.GetEntityId() == 0 {
		return nil, contentV1.ErrorBadRequest("entity type and id are required")
	}
	items, err := s.repo.ListByEntity(ctx, req.GetEntityType(), req.GetEntityId())
	if err != nil {
		return nil, err
	}
	return &contentV1.ListReviewLogsResponse{Items: items}, nil
}

// Count 留痕计数。
func (s *ReviewLogService) Count(ctx context.Context, req *contentV1.CountReviewLogsRequest) (*contentV1.CountReviewLogsResponse, error) {
	if req == nil || req.GetEntityType() == "" || req.GetEntityId() == 0 {
		return nil, contentV1.ErrorBadRequest("entity type and id are required")
	}
	items, err := s.repo.ListByEntity(ctx, req.GetEntityType(), req.GetEntityId())
	if err != nil {
		return nil, err
	}
	return &contentV1.CountReviewLogsResponse{Count: uint64(len(items))}, nil
}
