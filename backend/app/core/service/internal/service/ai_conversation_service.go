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

// AiConversationService AI 对话会话管理。
// 归属校验基于 viewer 上下文（BFF 透传操作者身份）：
// List 只回当前用户自己的会话，Get/Update/Delete 非本人一律 NotFound
// （不泄露他人会话的存在性）。
type AiConversationService struct {
	aiV1.UnimplementedAiConversationServiceServer
	log  *log.Helper
	repo *data.AiConversationRepo
}

func NewAiConversationService(ctx *bootstrap.Context, repo *data.AiConversationRepo) *AiConversationService {
	return &AiConversationService{
		log:  ctx.NewLoggerHelper("ai_conversation/service/core-service"),
		repo: repo,
	}
}

// List 只返回当前用户自己的会话。
func (s *AiConversationService) List(ctx context.Context, req *paginationV1.PagingRequest) (*aiV1.ListAiConversationResponse, error) {
	userId, _ := viewerUserIDFromContext(ctx)
	return s.repo.List(ctx, req, userId)
}

// Get 会话详情；归属校验失败一律 NotFound。
func (s *AiConversationService) Get(ctx context.Context, req *aiV1.GetAiConversationRequest) (*aiV1.AiConversation, error) {
	userId, _ := viewerUserIDFromContext(ctx)

	dto, err := s.repo.Get(ctx, req)
	if err != nil {
		return nil, err
	}
	if dto.GetUserId() != userId {
		return nil, aiV1.ErrorNotFound("ai conversation not found")
	}
	return dto, nil
}

func (s *AiConversationService) Update(ctx context.Context, req *aiV1.UpdateAiConversationRequest) (*emptypb.Empty, error) {
	if req.Data == nil {
		return nil, aiV1.ErrorBadRequest("invalid parameter")
	}
	userId, _ := viewerUserIDFromContext(ctx)

	// 归属校验：非本人按 NotFound 拒绝
	existing, err := s.repo.GetEntityByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	if existing.UserID == nil || *existing.UserID != userId {
		return nil, aiV1.ErrorNotFound("ai conversation not found")
	}

	if err = s.repo.Update(ctx, req); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *AiConversationService) Delete(ctx context.Context, req *aiV1.DeleteAiConversationRequest) (*emptypb.Empty, error) {
	userId, _ := viewerUserIDFromContext(ctx)

	existing, err := s.repo.GetEntityByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	if existing.UserID == nil || *existing.UserID != userId {
		return nil, aiV1.ErrorNotFound("ai conversation not found")
	}

	if err = s.repo.Delete(ctx, req); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *AiConversationService) Count(ctx context.Context, req *paginationV1.PagingRequest) (*aiV1.CountAiConversationResponse, error) {
	return s.repo.Count(ctx, req)
}
