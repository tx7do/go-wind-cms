package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	"go-wind-cms/app/core/service/internal/data"

	commentV1 "go-wind-cms/api/gen/go/comment/service/v1"
)

type CommentService struct {
	commentV1.UnimplementedCommentServiceServer

	commentRepo *data.CommentRepo
	log         *log.Helper
}

func NewCommentService(ctx *bootstrap.Context, uc *data.CommentRepo) *CommentService {
	return &CommentService{
		log:         ctx.NewLoggerHelper("comment/service/core-service"),
		commentRepo: uc,
	}
}

func (s *CommentService) List(ctx context.Context, req *paginationV1.PagingRequest) (*commentV1.ListCommentResponse, error) {
	return s.commentRepo.List(ctx, req)
}

func (s *CommentService) Get(ctx context.Context, req *commentV1.GetCommentRequest) (*commentV1.Comment, error) {
	return s.commentRepo.Get(ctx, req)
}

func (s *CommentService) Create(ctx context.Context, req *commentV1.CreateCommentRequest) (*commentV1.Comment, error) {
	return s.commentRepo.Create(ctx, req)
}

func (s *CommentService) Update(ctx context.Context, req *commentV1.UpdateCommentRequest) (*commentV1.Comment, error) {
	return s.commentRepo.Update(ctx, req)
}

func (s *CommentService) Delete(ctx context.Context, req *commentV1.DeleteCommentRequest) (*emptypb.Empty, error) {
	err := s.commentRepo.Delete(ctx, req)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
