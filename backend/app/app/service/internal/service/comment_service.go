package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"github.com/tx7do/go-utils/trans"
	"google.golang.org/protobuf/types/known/emptypb"

	appV1 "go-wind-cms/api/gen/go/app/service/v1"
	commentV1 "go-wind-cms/api/gen/go/comment/service/v1"

	"go-wind-cms/pkg/middleware/auth"
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
	if req == nil || req.Data == nil {
		return nil, commentV1.ErrorBadRequest("invalid parameter")
	}

	// 获取操作人信息，强制以服务端身份覆盖 CreatedBy，防止客户端伪造评论作者
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.Data.CreatedBy = trans.Ptr(operator.UserId)

	return s.commentClient.Create(ctx, req)
}

// ensureCommentOwner 校验调用者是否为目标评论的作者，防止任意登录用户改/删他人评论（IDOR）。
func (s *CommentService) ensureCommentOwner(ctx context.Context, commentID uint32) error {
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return commentV1.ErrorUnauthorized("authentication required")
	}

	comment, err := s.commentClient.Get(ctx, &commentV1.GetCommentRequest{
		Id: commentID,
	})
	if err != nil {
		return err
	}

	if comment.GetCreatedBy() != operator.GetUserId() {
		return commentV1.ErrorForbidden("you can only modify your own comments")
	}
	return nil
}

func (s *CommentService) Update(ctx context.Context, req *commentV1.UpdateCommentRequest) (*commentV1.Comment, error) {
	if err := s.ensureCommentOwner(ctx, req.GetId()); err != nil {
		return nil, err
	}

	// 获取操作人信息，强制以服务端身份覆盖 UpdatedBy，审计归属真实
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	req.Data.UpdatedBy = trans.Ptr(operator.GetUserId())

	return s.commentClient.Update(ctx, req)
}

func (s *CommentService) Delete(ctx context.Context, req *commentV1.DeleteCommentRequest) (*emptypb.Empty, error) {
	if err := s.ensureCommentOwner(ctx, req.GetId()); err != nil {
		return nil, err
	}
	return s.commentClient.Delete(ctx, req)
}
