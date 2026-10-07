package service

import (
	"context"

	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"

	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"

	aiV1 "go-wind-cms/api/gen/go/ai/service/v1"
	"go-wind-cms/app/core/service/internal/data"
)

// AiProviderService AI 模型提供商管理：core 侧纯持久化透传，
// 操作人（created_by/updated_by）与密钥加密由 admin BFF 层负责。
type AiProviderService struct {
	aiV1.UnimplementedAiProviderServiceServer
	log  *log.Helper
	repo *data.AiProviderRepo
}

func NewAiProviderService(ctx *bootstrap.Context, repo *data.AiProviderRepo) *AiProviderService {
	return &AiProviderService{
		log:  ctx.NewLoggerHelper("ai_provider/service/core-service"),
		repo: repo,
	}
}

func (s *AiProviderService) Count(ctx context.Context, req *paginationV1.PagingRequest) (*aiV1.CountAiProviderResponse, error) {
	return s.repo.Count(ctx, req)
}

func (s *AiProviderService) List(ctx context.Context, req *paginationV1.PagingRequest) (*aiV1.ListAiProviderResponse, error) {
	return s.repo.List(ctx, req)
}

func (s *AiProviderService) Get(ctx context.Context, req *aiV1.GetAiProviderRequest) (*aiV1.AiProvider, error) {
	return s.repo.Get(ctx, req)
}

func (s *AiProviderService) Create(ctx context.Context, req *aiV1.CreateAiProviderRequest) (*emptypb.Empty, error) {
	if err := s.repo.Create(ctx, req); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *AiProviderService) Update(ctx context.Context, req *aiV1.UpdateAiProviderRequest) (*emptypb.Empty, error) {
	if err := s.repo.Update(ctx, req); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *AiProviderService) Delete(ctx context.Context, req *aiV1.DeleteAiProviderRequest) (*emptypb.Empty, error) {
	if err := s.repo.Delete(ctx, req); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
