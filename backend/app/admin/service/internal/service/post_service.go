package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	contentV1 "go-wind-cms/api/gen/go/content/service/v1"

	"go-wind-cms/pkg/middleware/auth"
)

type PostService struct {
	adminV1.PostServiceHTTPServer

	postServiceClient contentV1.PostServiceClient
	log               *log.Helper
}

func NewPostService(ctx *bootstrap.Context, postServiceClient contentV1.PostServiceClient) *PostService {
	return &PostService{
		log:               ctx.NewLoggerHelper("post/service/admin-service"),
		postServiceClient: postServiceClient,
	}
}

func (s *PostService) List(ctx context.Context, req *paginationV1.PagingRequest) (*contentV1.ListPostResponse, error) {
	return s.postServiceClient.List(ctx, req)
}

func (s *PostService) Get(ctx context.Context, req *contentV1.GetPostRequest) (*contentV1.Post, error) {
	return s.postServiceClient.Get(ctx, req)
}

func (s *PostService) Create(ctx context.Context, req *contentV1.CreatePostRequest) (*contentV1.Post, error) {
	if req.Data == nil {
		return nil, adminV1.ErrorBadRequest("invalid parameter")
	}

	// 获取操作人信息
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.Data.CreatedBy = trans.Ptr(operator.UserId)
	req.Data.AuthorId = trans.Ptr(operator.UserId)

	for i := range req.Data.Translations {
		req.Data.Translations[i].CreatedBy = trans.Ptr(operator.UserId)
	}

	return s.postServiceClient.Create(ctx, req)
}

func (s *PostService) Update(ctx context.Context, req *contentV1.UpdatePostRequest) (*contentV1.Post, error) {
	if req.Data == nil {
		return nil, adminV1.ErrorBadRequest("invalid parameter")
	}

	// 获取操作人信息
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.Data.Id = trans.Ptr(req.GetId())

	req.Data.UpdatedBy = trans.Ptr(operator.GetUserId())
	if req.UpdateMask != nil {
		req.UpdateMask.Paths = append(req.UpdateMask.Paths, "updated_by")
	}

	for i := range req.Data.Translations {
		req.Data.Translations[i].UpdatedBy = trans.Ptr(operator.UserId)
	}

	return s.postServiceClient.Update(ctx, req)
}

func (s *PostService) Delete(ctx context.Context, req *contentV1.DeletePostRequest) (*emptypb.Empty, error) {
	return s.postServiceClient.Delete(ctx, req)
}

func (s *PostService) IsExistTranslation(ctx context.Context, req *contentV1.IsExistTranslationRequest) (*contentV1.IsExistTranslationResponse, error) {
	return s.postServiceClient.IsExistTranslation(ctx, req)
}
