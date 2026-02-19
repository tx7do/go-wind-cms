package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	"go-wind-cms/app/core/service/internal/data"

	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
)

type PostService struct {
	contentV1.UnimplementedPostServiceServer

	postRepo *data.PostRepo
	log      *log.Helper
}

func NewPostService(ctx *bootstrap.Context, uc *data.PostRepo) *PostService {
	return &PostService{
		log:      ctx.NewLoggerHelper("post/service/core-service"),
		postRepo: uc,
	}
}

func (s *PostService) List(ctx context.Context, req *paginationV1.PagingRequest) (*contentV1.ListPostResponse, error) {
	return s.postRepo.List(ctx, req)
}

func (s *PostService) Get(ctx context.Context, req *contentV1.GetPostRequest) (*contentV1.Post, error) {
	return s.postRepo.Get(ctx, req)
}

func (s *PostService) Create(ctx context.Context, req *contentV1.CreatePostRequest) (*contentV1.Post, error) {
	return s.postRepo.Create(ctx, req)
}

func (s *PostService) Update(ctx context.Context, req *contentV1.UpdatePostRequest) (*contentV1.Post, error) {
	return s.postRepo.Update(ctx, req)
}

func (s *PostService) Delete(ctx context.Context, req *contentV1.DeletePostRequest) (*emptypb.Empty, error) {
	err := s.postRepo.Delete(ctx, req)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
