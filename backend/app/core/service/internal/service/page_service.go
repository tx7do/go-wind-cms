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

type PageService struct {
	contentV1.UnimplementedPageServiceServer

	pageRepo *data.PageRepo
	log      *log.Helper
}

func NewPageService(ctx *bootstrap.Context, uc *data.PageRepo) *PageService {
	return &PageService{
		log:      ctx.NewLoggerHelper("page/service/core-service"),
		pageRepo: uc,
	}
}

func (s *PageService) List(ctx context.Context, req *paginationV1.PagingRequest) (*contentV1.ListPageResponse, error) {
	return s.pageRepo.List(ctx, req)
}

func (s *PageService) Get(ctx context.Context, req *contentV1.GetPageRequest) (*contentV1.Page, error) {
	return s.pageRepo.Get(ctx, req)
}

func (s *PageService) Create(ctx context.Context, req *contentV1.CreatePageRequest) (*contentV1.Page, error) {
	return s.pageRepo.Create(ctx, req)
}

func (s *PageService) Update(ctx context.Context, req *contentV1.UpdatePageRequest) (*contentV1.Page, error) {
	return s.pageRepo.Update(ctx, req)
}

func (s *PageService) Delete(ctx context.Context, req *contentV1.DeletePageRequest) (*emptypb.Empty, error) {
	err := s.pageRepo.Delete(ctx, req)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
