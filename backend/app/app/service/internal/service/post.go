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
	return s.postClient.List(ctx, req)
}

func (s *PostService) Get(ctx context.Context, req *contentV1.GetPostRequest) (*contentV1.Post, error) {
	return s.postClient.Get(ctx, req)
}

func (s *PostService) Create(ctx context.Context, req *contentV1.CreatePostRequest) (*contentV1.Post, error) {
	return s.postClient.Create(ctx, req)
}

func (s *PostService) Update(ctx context.Context, req *contentV1.UpdatePostRequest) (*contentV1.Post, error) {
	return s.postClient.Update(ctx, req)
}

func (s *PostService) Delete(ctx context.Context, req *contentV1.DeletePostRequest) (*emptypb.Empty, error) {
	return s.postClient.Delete(ctx, req)
}
