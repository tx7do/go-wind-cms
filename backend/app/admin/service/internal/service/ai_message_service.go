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

// AiMessageService AI 对话消息查询/删除（BFF 直通转发；
// 范围过滤由 core 侧基于 viewer 上下文完成）。
type AiMessageService struct {
	adminV1.AiMessageServiceHTTPServer

	aiMessageServiceClient aiV1.AiMessageServiceClient
	log                    *log.Helper
}

func NewAiMessageService(ctx *bootstrap.Context, aiMessageServiceClient aiV1.AiMessageServiceClient) *AiMessageService {
	return &AiMessageService{
		log:                    ctx.NewLoggerHelper("ai_message/service/admin-service"),
		aiMessageServiceClient: aiMessageServiceClient,
	}
}

func (s *AiMessageService) List(ctx context.Context, req *paginationV1.PagingRequest) (*aiV1.ListAiMessageResponse, error) {
	return s.aiMessageServiceClient.List(ctx, req)
}

func (s *AiMessageService) Count(ctx context.Context, req *paginationV1.PagingRequest) (*aiV1.CountAiMessageResponse, error) {
	return s.aiMessageServiceClient.Count(ctx, req)
}

func (s *AiMessageService) Delete(ctx context.Context, req *aiV1.DeleteAiMessageRequest) (*emptypb.Empty, error) {
	return s.aiMessageServiceClient.Delete(ctx, req)
}
