package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"

	aiV1 "go-wind-cms/api/gen/go/ai/service/v1"
	"go-wind-cms/app/core/service/internal/data"
)

// AiMessageService AI 对话消息查询/删除。
// List 仅限本人会话的消息（经会话归属校验）。
type AiMessageService struct {
	aiV1.UnimplementedAiMessageServiceServer
	log  *log.Helper
	repo *data.AiMessageRepo
}

func NewAiMessageService(ctx *bootstrap.Context, repo *data.AiMessageRepo) *AiMessageService {
	return &AiMessageService{
		log:  ctx.NewLoggerHelper("ai_message/service/core-service"),
		repo: repo,
	}
}

// List 仅返回当前用户自己的会话消息（repo 层经 ownerUserID 过滤）。
func (s *AiMessageService) List(ctx context.Context, req *paginationV1.PagingRequest) (*aiV1.ListAiMessageResponse, error) {
	userId, _ := viewerUserIDFromContext(ctx)
	return s.repo.List(ctx, req, userId)
}

func (s *AiMessageService) Count(ctx context.Context, req *paginationV1.PagingRequest) (*aiV1.CountAiMessageResponse, error) {
	return s.repo.Count(ctx, req)
}

func (s *AiMessageService) Delete(ctx context.Context, req *aiV1.DeleteAiMessageRequest) (*emptypb.Empty, error) {
	if err := s.repo.Delete(ctx, req); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
