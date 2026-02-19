package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	appV1 "go-wind-cms/api/gen/go/app/service/v1"
	commentV1 "go-wind-cms/api/gen/go/comment/service/v1"
)

type CommentService struct {
	appV1.CommentServiceHTTPServer

	commentClient commentV1.CommentServiceClient
	log           *log.Helper
}

func NewCommentService(ctx *bootstrap.Context, commentClient commentV1.CommentServiceClient) *CommentService {
	return &CommentService{
		log:           ctx.NewLoggerHelper("comment/service/app-service"),
		commentClient: commentClient,
	}
}

func (s *CommentService) List(ctx context.Context, req *paginationV1.PagingRequest) (*commentV1.ListCommentResponse, error) {
	return s.commentClient.List(ctx, req)
}

func (s *CommentService) Get(ctx context.Context, req *commentV1.GetCommentRequest) (*commentV1.Comment, error) {
	return s.commentClient.Get(ctx, req)
}

func (s *CommentService) Create(ctx context.Context, req *commentV1.CreateCommentRequest) (*commentV1.Comment, error) {
	return s.commentClient.Create(ctx, req)
}

func (s *CommentService) Update(ctx context.Context, req *commentV1.UpdateCommentRequest) (*commentV1.Comment, error) {
	return s.commentClient.Update(ctx, req)
}

func (s *CommentService) Delete(ctx context.Context, req *commentV1.DeleteCommentRequest) (*emptypb.Empty, error) {
	return s.commentClient.Delete(ctx, req)
}
