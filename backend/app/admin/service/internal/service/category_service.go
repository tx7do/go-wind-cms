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
	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
)

type CategoryService struct {
	adminV1.CategoryServiceHTTPServer

	categoryServiceClient contentV1.CategoryServiceClient
	log                   *log.Helper
}

func NewCategoryService(ctx *bootstrap.Context, categoryServiceClient contentV1.CategoryServiceClient) *CategoryService {
	return &CategoryService{
		log:                   ctx.NewLoggerHelper("category/service/admin-service"),
		categoryServiceClient: categoryServiceClient,
	}
}

func (s *CategoryService) List(ctx context.Context, req *paginationV1.PagingRequest) (*contentV1.ListCategoryResponse, error) {
	return s.categoryServiceClient.List(ctx, req)
}

func (s *CategoryService) Get(ctx context.Context, req *contentV1.GetCategoryRequest) (*contentV1.Category, error) {
	return s.categoryServiceClient.Get(ctx, req)
}

func (s *CategoryService) Create(ctx context.Context, req *contentV1.CreateCategoryRequest) (*contentV1.Category, error) {
	if req.Data == nil {
		return nil, adminV1.ErrorBadRequest("invalid parameter")
	}

	// 获取操作人信息
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.Data.CreatedBy = trans.Ptr(operator.UserId)

	return s.categoryServiceClient.Create(ctx, req)
}

func (s *CategoryService) Update(ctx context.Context, req *contentV1.UpdateCategoryRequest) (*contentV1.Category, error) {
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

	return s.categoryServiceClient.Update(ctx, req)
}

func (s *CategoryService) Delete(ctx context.Context, req *contentV1.DeleteCategoryRequest) (*emptypb.Empty, error) {
	return s.categoryServiceClient.Delete(ctx, req)
}
