package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	"go-wind-cms/app/core/service/internal/data"

	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
	interactionV1 "go-wind-cms/api/gen/go/interaction/service/v1"
)

// InteractionService 是点赞/收藏计数内聚子系统的服务端入口。
//
// 安全契约：
//   - viewer 用户身份强制从鉴权上下文（viewer context）提取，绝不接受请求体
//     中的 user_id。未登录调用一律 401。
//   - 跨租户隔离、防重复、幂等性由 InteractionRepo 保证。
//
// 依赖注入：同时持有 InteractionRepo（ledger/计数写入）与 PostRepo（ListWatchedPosts
// 复用其附带查询），遵循 user_profile_service 的多依赖聚合先例。
type InteractionService struct {
	interactionV1.UnimplementedInteractionServiceServer

	interactionRepo *data.InteractionRepo
	postRepo        *data.PostRepo
	log             *log.Helper
}

func NewInteractionService(
	ctx *bootstrap.Context,
	interactionRepo *data.InteractionRepo,
	postRepo *data.PostRepo,
) *InteractionService {
	return &InteractionService{
		log:             ctx.NewLoggerHelper("interaction/service/core-service"),
		interactionRepo: interactionRepo,
		postRepo:        postRepo,
	}
}

// Like 点赞。viewer 身份从鉴权上下文提取，请求体的 user_id 被忽略。
func (s *InteractionService) Like(ctx context.Context, req *interactionV1.LikeRequest) (*interactionV1.LikeResponse, error) {
	viewerUserID, ok := viewerUserIDFromContext(ctx)
	if !ok {
		return nil, interactionV1.ErrorUnauthorized("login required")
	}
	liked, likeCount, err := s.interactionRepo.Like(ctx, viewerUserID, req.GetTargetType(), req.GetTargetId())
	if err != nil {
		return nil, err
	}
	return &interactionV1.LikeResponse{Liked: liked, LikeCount: likeCount}, nil
}

// Unlike 取消点赞。
func (s *InteractionService) Unlike(ctx context.Context, req *interactionV1.LikeRequest) (*interactionV1.LikeResponse, error) {
	viewerUserID, ok := viewerUserIDFromContext(ctx)
	if !ok {
		return nil, interactionV1.ErrorUnauthorized("login required")
	}
	liked, likeCount, err := s.interactionRepo.Unlike(ctx, viewerUserID, req.GetTargetType(), req.GetTargetId())
	if err != nil {
		return nil, err
	}
	return &interactionV1.LikeResponse{Liked: liked, LikeCount: likeCount}, nil
}

// Watch 收藏 post。
func (s *InteractionService) Watch(ctx context.Context, req *interactionV1.WatchRequest) (*interactionV1.WatchResponse, error) {
	viewerUserID, ok := viewerUserIDFromContext(ctx)
	if !ok {
		return nil, interactionV1.ErrorUnauthorized("login required")
	}
	watched, err := s.interactionRepo.Watch(ctx, viewerUserID, req.GetPostId())
	if err != nil {
		return nil, err
	}
	return &interactionV1.WatchResponse{Watched: watched}, nil
}

// Unwatch 取消收藏 post。
func (s *InteractionService) Unwatch(ctx context.Context, req *interactionV1.WatchRequest) (*interactionV1.WatchResponse, error) {
	viewerUserID, ok := viewerUserIDFromContext(ctx)
	if !ok {
		return nil, interactionV1.ErrorUnauthorized("login required")
	}
	watched, err := s.interactionRepo.Unwatch(ctx, viewerUserID, req.GetPostId())
	if err != nil {
		return nil, err
	}
	return &interactionV1.WatchResponse{Watched: watched}, nil
}

// GetInteractionStatus 批量查询当前 viewer 对指定目标的交互状态。
func (s *InteractionService) GetInteractionStatus(ctx context.Context, req *interactionV1.GetInteractionStatusRequest) (*interactionV1.GetInteractionStatusResponse, error) {
	viewerUserID, ok := viewerUserIDFromContext(ctx)
	if !ok {
		return nil, interactionV1.ErrorUnauthorized("login required")
	}
	statuses, err := s.interactionRepo.GetInteractionStatus(ctx, viewerUserID, req.GetTargetType(), req.GetTargetIds())
	if err != nil {
		return nil, err
	}
	return &interactionV1.GetInteractionStatusResponse{Statuses: statuses}, nil
}

// ListWatchedPosts 列出当前 viewer 收藏的 post。
// 复用 PostRepo.Get 的附带查询与 view_mask 处理。
func (s *InteractionService) ListWatchedPosts(ctx context.Context, req *paginationV1.PagingRequest) (*contentV1.ListPostResponse, error) {
	viewerUserID, ok := viewerUserIDFromContext(ctx)
	if !ok {
		return nil, interactionV1.ErrorUnauthorized("login required")
	}
	return s.interactionRepo.ListWatchedPosts(ctx, viewerUserID, s.postRepo, req)
}
