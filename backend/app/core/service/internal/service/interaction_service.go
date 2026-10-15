package service

import (
	"context"
	"strconv"

	"github.com/go-kratos/kratos/v2/log"
	kratosMetadata "github.com/go-kratos/kratos/v2/metadata"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/go-crud/viewer"
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
	watched, watchCount, err := s.interactionRepo.Watch(ctx, viewerUserID, req.GetPostId())
	if err != nil {
		return nil, err
	}
	return &interactionV1.WatchResponse{Watched: watched, WatchCount: watchCount}, nil
}

// Unwatch 取消收藏 post。
func (s *InteractionService) Unwatch(ctx context.Context, req *interactionV1.WatchRequest) (*interactionV1.WatchResponse, error) {
	viewerUserID, ok := viewerUserIDFromContext(ctx)
	if !ok {
		return nil, interactionV1.ErrorUnauthorized("login required")
	}
	watched, watchCount, err := s.interactionRepo.Unwatch(ctx, viewerUserID, req.GetPostId())
	if err != nil {
		return nil, err
	}
	return &interactionV1.WatchResponse{Watched: watched, WatchCount: watchCount}, nil
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

// GetCounts 批量查询指定目标的计数（如点赞数）。
//
// 计数查询本身不依赖 viewer 用户身份，仅按 tenant 隔离。登录用户由 repo 层
// 从 viewer context 提取 tenant_id 过滤；未登录（SystemViewer，tenant_id==0）
// 由 repo 层返回空结果，避免跨租户泄漏的同时不抛 401 打断公开页渲染。
// Like/Unlike/Watch 等写操作仍强制登录。
func (s *InteractionService) GetCounts(ctx context.Context, req *interactionV1.GetCountsRequest) (*interactionV1.GetCountsResponse, error) {
	counts, err := s.interactionRepo.GetCounts(ctx, req.GetTargetType(), req.GetTargetIds(), req.GetMetrics())
	if err != nil {
		return nil, err
	}
	return &interactionV1.GetCountsResponse{Counts: counts}, nil
}

// RecordView 浏览量上报（匿名可调用）。
//
// 租户归属解析顺序与评论链路一致：viewer 携带租户（登录态）优先；
// 匿名链路 app BFF 经内网可信元数据 x-md-global-tenant-id 注入（不可被外部伪造）；
// 两者皆缺省（本地直连 core 等）按 0（平台）落库。
func (s *InteractionService) RecordView(ctx context.Context, req *interactionV1.RecordViewRequest) (*interactionV1.RecordViewResponse, error) {
	if req == nil || req.GetTargetId() == 0 || req.GetTargetType() == interactionV1.TargetType_TARGET_TYPE_UNSPECIFIED {
		return nil, interactionV1.ErrorBadRequest("invalid parameter")
	}

	var tid uint32
	if vc, exist := viewer.FromContext(ctx); exist && vc != nil && vc.TenantID() > 0 {
		tid = uint32(vc.TenantID())
	} else {
		tid = tenantIDFromInternalMetadata(ctx)
	}

	count, err := s.interactionRepo.RecordView(ctx, tid, req.GetTargetId(), req.GetTargetType())
	if err != nil {
		return nil, err
	}
	return &interactionV1.RecordViewResponse{ViewCount: count}, nil
}

// tenantIDFromInternalMetadata 从内网可信元数据 x-md-global-tenant-id 读租户
//（app BFF 按请求 Host 解析后注入；该 header 仅内网 gRPC 可达，外部不可伪造）。
func tenantIDFromInternalMetadata(ctx context.Context) uint32 {
	if md, ok := kratosMetadata.FromServerContext(ctx); ok {
		if v := md.Get("x-md-global-tenant-id"); v != "" {
			if tid, err := strconv.ParseUint(v, 10, 32); err == nil && tid > 0 {
				return uint32(tid)
			}
		}
	}
	return 0
}
