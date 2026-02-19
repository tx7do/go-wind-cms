package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	"go-wind-cms/app/core/service/internal/data"

	internalMessageV1 "go-wind-cms/api/gen/go/internal_message/service/v1"
)

type InternalMessageCategoryService struct {
	internalMessageV1.UnimplementedInternalMessageCategoryServiceServer

	log *log.Helper

	repo *data.InternalMessageCategoryRepo
}

func NewInternalMessageCategoryService(ctx *bootstrap.Context, repo *data.InternalMessageCategoryRepo) *InternalMessageCategoryService {
	return &InternalMessageCategoryService{
		log:  ctx.NewLoggerHelper("internal-message-category/service/core-service"),
		repo: repo,
	}
}

func (s *InternalMessageCategoryService) List(ctx context.Context, req *paginationV1.PagingRequest) (*internalMessageV1.ListInternalMessageCategoryResponse, error) {
	return s.repo.List(ctx, req)
}

func (s *InternalMessageCategoryService) Get(ctx context.Context, req *internalMessageV1.GetInternalMessageCategoryRequest) (*internalMessageV1.InternalMessageCategory, error) {
	return s.repo.Get(ctx, req)
}

func (s *InternalMessageCategoryService) Create(ctx context.Context, req *internalMessageV1.CreateInternalMessageCategoryRequest) (*emptypb.Empty, error) {
	if req.Data == nil {
		return nil, internalMessageV1.ErrorBadRequest("invalid parameter")
	}

	if err := s.repo.Create(ctx, req); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *InternalMessageCategoryService) Update(ctx context.Context, req *internalMessageV1.UpdateInternalMessageCategoryRequest) (*emptypb.Empty, error) {
	if req.Data == nil {
		return nil, internalMessageV1.ErrorBadRequest("invalid parameter")
	}

	if err := s.repo.Update(ctx, req); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *InternalMessageCategoryService) Delete(ctx context.Context, req *internalMessageV1.DeleteInternalMessageCategoryRequest) (*emptypb.Empty, error) {
	if err := s.repo.Delete(ctx, req); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
