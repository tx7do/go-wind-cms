package service

import (
	"context"
	"go-wind-cms/pkg/middleware/auth"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	commentV1 "go-wind-cms/api/gen/go/comment/service/v1"
)

type CommentService struct {
	adminV1.CommentServiceHTTPServer

	commentServiceClient commentV1.CommentServiceClient
	log                  *log.Helper
}

func NewCommentService(ctx *bootstrap.Context, commentServiceClient commentV1.CommentServiceClient) *CommentService {
	return &CommentService{
		log:                  ctx.NewLoggerHelper("comment/service/admin-service"),
		commentServiceClient: commentServiceClient,
	}
}

func (s *CommentService) List(ctx context.Context, req *paginationV1.PagingRequest) (*commentV1.ListCommentResponse, error) {
	return s.commentServiceClient.List(ctx, req)
}

func (s *CommentService) Get(ctx context.Context, req *commentV1.GetCommentRequest) (*commentV1.Comment, error) {
	return s.commentServiceClient.Get(ctx, req)
}

func (s *CommentService) Create(ctx context.Context, req *commentV1.CreateCommentRequest) (*commentV1.Comment, error) {
	if req.Data == nil {
		return nil, adminV1.ErrorBadRequest("invalid parameter")
	}

	// 获取操作人信息
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.Data.CreatedBy = trans.Ptr(operator.UserId)

	return s.commentServiceClient.Create(ctx, req)
}

func (s *CommentService) Update(ctx context.Context, req *commentV1.UpdateCommentRequest) (*commentV1.Comment, error) {
	if req.Data == nil {
		return nil, adminV1.ErrorBadRequest("invalid parameter")
	}

	// 获取操作人信息
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.Data.UpdatedBy = trans.Ptr(operator.GetUserId())
	if req.UpdateMask != nil {
		req.UpdateMask.Paths = append(req.UpdateMask.Paths, "updated_by")
	}

	return s.commentServiceClient.Update(ctx, req)
}

func (s *CommentService) Delete(ctx context.Context, req *commentV1.DeleteCommentRequest) (*emptypb.Empty, error) {
	return s.commentServiceClient.Delete(ctx, req)
}
