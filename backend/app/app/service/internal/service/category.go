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

type CategoryService struct {
	appV1.CategoryServiceHTTPServer

	categoryClient contentV1.CategoryServiceClient
	log            *log.Helper
}

func NewCategoryService(ctx *bootstrap.Context, categoryClient contentV1.CategoryServiceClient) *CategoryService {
	return &CategoryService{
		log:            ctx.NewLoggerHelper("category/service/app-service"),
		categoryClient: categoryClient,
	}
}

func (s *CategoryService) List(ctx context.Context, req *paginationV1.PagingRequest) (*contentV1.ListCategoryResponse, error) {
	return s.categoryClient.List(ctx, req)
}

func (s *CategoryService) Get(ctx context.Context, req *contentV1.GetCategoryRequest) (*contentV1.Category, error) {
	return s.categoryClient.Get(ctx, req)
}

func (s *CategoryService) Create(ctx context.Context, req *contentV1.CreateCategoryRequest) (*contentV1.Category, error) {
	return s.categoryClient.Create(ctx, req)
}

func (s *CategoryService) Update(ctx context.Context, req *contentV1.UpdateCategoryRequest) (*contentV1.Category, error) {
	return s.categoryClient.Update(ctx, req)
}

func (s *CategoryService) Delete(ctx context.Context, req *contentV1.DeleteCategoryRequest) (*emptypb.Empty, error) {
	return s.categoryClient.Delete(ctx, req)
}
