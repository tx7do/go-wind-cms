package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	appV1 "go-wind-cms/api/gen/go/app/service/v1"
	interactionV1 "go-wind-cms/api/gen/go/interaction/service/v1"

	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
)

// InteractionService 是 app 网关的点赞/收藏转发器。
//
// 作为纯透传层，把 HTTP 请求转发给 core 服务的 InteractionService。所有鉴权、
// viewer 身份提取、幂等、计数一致性均在 core 侧完成。本层不做业务过滤。
//
// 安全说明：Interaction 的写操作（Like/Unlike/Watch/Unwatch）及含 viewer 个人
// 状态的查询（GetInteractionStatus/ListWatchedPosts）均需登录，在 rest_server 的
// rpc.AddWhiteList 中刻意不登记，强制走鉴权中间件。仅 GetCounts（公开计数，按
// tenant 隔离、不依赖 viewer 身份）登记为白名单，供文章列表/详情匿名展示点赞数。
type InteractionService struct {
	appV1.InteractionServiceHTTPServer

	interactionClient interactionV1.InteractionServiceClient
	log               *log.Helper
}

func NewInteractionService(ctx *bootstrap.Context, interactionClient interactionV1.InteractionServiceClient) *InteractionService {
	return &InteractionService{
		log:               ctx.NewLoggerHelper("interaction/service/app-service"),
		interactionClient: interactionClient,
	}
}

func (s *InteractionService) Like(ctx context.Context, req *interactionV1.LikeRequest) (*interactionV1.LikeResponse, error) {
	return s.interactionClient.Like(ctx, req)
}

func (s *InteractionService) Unlike(ctx context.Context, req *interactionV1.LikeRequest) (*interactionV1.LikeResponse, error) {
	return s.interactionClient.Unlike(ctx, req)
}

func (s *InteractionService) Watch(ctx context.Context, req *interactionV1.WatchRequest) (*interactionV1.WatchResponse, error) {
	return s.interactionClient.Watch(ctx, req)
}

func (s *InteractionService) Unwatch(ctx context.Context, req *interactionV1.WatchRequest) (*interactionV1.WatchResponse, error) {
	return s.interactionClient.Unwatch(ctx, req)
}

func (s *InteractionService) GetInteractionStatus(ctx context.Context, req *interactionV1.GetInteractionStatusRequest) (*interactionV1.GetInteractionStatusResponse, error) {
	return s.interactionClient.GetInteractionStatus(ctx, req)
}

func (s *InteractionService) ListWatchedPosts(ctx context.Context, req *paginationV1.PagingRequest) (*contentV1.ListPostResponse, error) {
	return s.interactionClient.ListWatchedPosts(ctx, req)
}

func (s *InteractionService) GetCounts(ctx context.Context, req *interactionV1.GetCountsRequest) (*interactionV1.GetCountsResponse, error) {
	return s.interactionClient.GetCounts(ctx, req)
}
