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

type CategoryService struct {
	contentV1.UnimplementedCategoryServiceServer

	categoryRepo *data.CategoryRepo
	log          *log.Helper
}

func NewCategoryService(ctx *bootstrap.Context, uc *data.CategoryRepo) *CategoryService {
	return &CategoryService{
		log:          ctx.NewLoggerHelper("category/service/core-service"),
		categoryRepo: uc,
	}
}

func (s *CategoryService) List(ctx context.Context, req *paginationV1.PagingRequest) (*contentV1.ListCategoryResponse, error) {
	return s.categoryRepo.List(ctx, req)
}

func (s *CategoryService) Get(ctx context.Context, req *contentV1.GetCategoryRequest) (*contentV1.Category, error) {
	return s.categoryRepo.Get(ctx, req)
}

func (s *CategoryService) Create(ctx context.Context, req *contentV1.CreateCategoryRequest) (*contentV1.Category, error) {
	return s.categoryRepo.Create(ctx, req)
}

func (s *CategoryService) Update(ctx context.Context, req *contentV1.UpdateCategoryRequest) (*contentV1.Category, error) {
	return s.categoryRepo.Update(ctx, req)
}

func (s *CategoryService) Delete(ctx context.Context, req *contentV1.DeleteCategoryRequest) (*emptypb.Empty, error) {
	err := s.categoryRepo.Delete(ctx, req)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
