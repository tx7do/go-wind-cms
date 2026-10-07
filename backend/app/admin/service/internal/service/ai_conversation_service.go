package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	aiV1 "go-wind-cms/api/gen/go/ai/service/v1"
	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
)

// AiConversationService AI 对话会话管理（BFF 直通转发；
// 归属校验由 core 侧基于 viewer 上下文完成）。
type AiConversationService struct {
	adminV1.AiConversationServiceHTTPServer

	aiConversationServiceClient aiV1.AiConversationServiceClient
	log                         *log.Helper
}

func NewAiConversationService(ctx *bootstrap.Context, aiConversationServiceClient aiV1.AiConversationServiceClient) *AiConversationService {
	return &AiConversationService{
		log:                         ctx.NewLoggerHelper("ai_conversation/service/admin-service"),
		aiConversationServiceClient: aiConversationServiceClient,
	}
}

func (s *AiConversationService) List(ctx context.Context, req *paginationV1.PagingRequest) (*aiV1.ListAiConversationResponse, error) {
	return s.aiConversationServiceClient.List(ctx, req)
}

func (s *AiConversationService) Get(ctx context.Context, req *aiV1.GetAiConversationRequest) (*aiV1.AiConversation, error) {
	return s.aiConversationServiceClient.Get(ctx, req)
}

func (s *AiConversationService) Update(ctx context.Context, req *aiV1.UpdateAiConversationRequest) (*emptypb.Empty, error) {
	return s.aiConversationServiceClient.Update(ctx, req)
}

func (s *AiConversationService) Delete(ctx context.Context, req *aiV1.DeleteAiConversationRequest) (*emptypb.Empty, error) {
	return s.aiConversationServiceClient.Delete(ctx, req)
}

func (s *AiConversationService) Count(ctx context.Context, req *paginationV1.PagingRequest) (*aiV1.CountAiConversationResponse, error) {
	return s.aiConversationServiceClient.Count(ctx, req)
}
