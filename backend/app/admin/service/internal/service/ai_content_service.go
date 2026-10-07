package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	aiV1 "go-wind-cms/api/gen/go/ai/service/v1"
	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
)

// AiContentService AI 内容生成（表单助手）BFF 直通转发。
type AiContentService struct {
	adminV1.AiContentServiceHTTPServer

	aiContentServiceClient aiV1.AiContentServiceClient
	log                    *log.Helper
}

func NewAiContentService(ctx *bootstrap.Context, aiContentServiceClient aiV1.AiContentServiceClient) *AiContentService {
	return &AiContentService{
		log:                    ctx.NewLoggerHelper("ai_content/service/admin-service"),
		aiContentServiceClient: aiContentServiceClient,
	}
}

func (s *AiContentService) GenerateContent(ctx context.Context, req *aiV1.GenerateContentRequest) (*aiV1.GenerateContentResponse, error) {
	return s.aiContentServiceClient.GenerateContent(ctx, req)
}

func (s *AiContentService) SemanticSearch(ctx context.Context, req *aiV1.SemanticSearchRequest) (*aiV1.SemanticSearchResponse, error) {
	return s.aiContentServiceClient.SemanticSearch(ctx, req)
}

func (s *AiContentService) RebuildSearchIndex(ctx context.Context, _ *emptypb.Empty) (*aiV1.RebuildSearchIndexResponse, error) {
	return s.aiContentServiceClient.RebuildSearchIndex(ctx, &emptypb.Empty{})
}

func (s *AiContentService) GetSafetyInsights(ctx context.Context, req *aiV1.SafetyInsightsRequest) (*aiV1.SafetyInsightsResponse, error) {
	return s.aiContentServiceClient.GetSafetyInsights(ctx, req)
}
