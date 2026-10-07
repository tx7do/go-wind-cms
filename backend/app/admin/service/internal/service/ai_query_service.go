package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	aiV1 "go-wind-cms/api/gen/go/ai/service/v1"
	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
)

// AiQueryService 智能问数（BFF 直通转发；
// 白名单/只读护栏/租户谓词校验全部在 core 侧执行）。
type AiQueryService struct {
	adminV1.AiQueryServiceHTTPServer

	aiQueryServiceClient aiV1.AiQueryServiceClient
	log                  *log.Helper
}

func NewAiQueryService(ctx *bootstrap.Context, aiQueryServiceClient aiV1.AiQueryServiceClient) *AiQueryService {
	return &AiQueryService{
		log:                  ctx.NewLoggerHelper("ai_query/service/admin-service"),
		aiQueryServiceClient: aiQueryServiceClient,
	}
}

func (s *AiQueryService) Ask(ctx context.Context, req *aiV1.AskAiQueryRequest) (*aiV1.AskAiQueryResponse, error) {
	return s.aiQueryServiceClient.Ask(ctx, req)
}
