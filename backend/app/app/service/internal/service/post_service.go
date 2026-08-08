package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	appV1 "go-wind-cms/api/gen/go/app/service/v1"
	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
)

type PostService struct {
	appV1.PostServiceHTTPServer

	postClient contentV1.PostServiceClient
	log        *log.Helper
}

func NewPostService(ctx *bootstrap.Context, postClient contentV1.PostServiceClient) *PostService {
	return &PostService{
		log:        ctx.NewLoggerHelper("post/service/app-service"),
		postClient: postClient,
	}
}

func (s *PostService) List(ctx context.Context, req *paginationV1.PagingRequest) (*contentV1.ListPostResponse, error) {
	resp, err := s.postClient.List(ctx, req)
	if err != nil {
		return nil, err
	}
	// 公开端点仅返回已发布文章，过滤草稿/归档/私有等状态
	if resp != nil {
		filtered := make([]*contentV1.Post, 0, len(resp.GetItems()))
		for _, p := range resp.GetItems() {
			if p != nil && p.GetStatus() == contentV1.Post_POST_STATUS_PUBLISHED {
				filtered = append(filtered, p)
			}
		}
		resp.Items = filtered
		resp.Total = uint64(len(filtered))
	}
	return resp, nil
}

func (s *PostService) Get(ctx context.Context, req *contentV1.GetPostRequest) (*contentV1.Post, error) {
	resp, err := s.postClient.Get(ctx, req)
	if err != nil {
		return nil, err
	}
	// 公开端点仅返回已发布文章，草稿/归档/私有等状态按未找到处理
	if resp == nil || resp.GetStatus() != contentV1.Post_POST_STATUS_PUBLISHED {
		return nil, contentV1.ErrorNotFound("post not found")
	}
	return resp, nil
}

// Create/Update/Delete 在 app（公开站点）服务上禁用：CMS 内容的写操作应经由 admin 服务，
// 公开站点登录用户不应直接创建/修改/删除文章。RBAC 为故意的 noop，故在此显式拒绝。
func (s *PostService) Create(_ context.Context, _ *contentV1.CreatePostRequest) (*contentV1.Post, error) {
	return nil, contentV1.ErrorForbidden("content mutation is not allowed on the public app service")
}

func (s *PostService) Update(_ context.Context, _ *contentV1.UpdatePostRequest) (*contentV1.Post, error) {
	return nil, contentV1.ErrorForbidden("content mutation is not allowed on the public app service")
}

func (s *PostService) Delete(_ context.Context, _ *contentV1.DeletePostRequest) (*emptypb.Empty, error) {
	return nil, contentV1.ErrorForbidden("content mutation is not allowed on the public app service")
}

func (s *PostService) GetTranslation(ctx context.Context, req *contentV1.GetPostRequest) (*contentV1.PostTranslation, error) {
	return s.postClient.GetTranslation(ctx, req)
}
