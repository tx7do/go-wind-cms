package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	"go-wind-cms/app/core/service/internal/data"

	siteV1 "go-wind-cms/api/gen/go/site/service/v1"
)

type NavigationService struct {
	siteV1.UnimplementedNavigationServiceServer

	navigationRepo *data.NavigationRepo
	log            *log.Helper
}

func NewNavigationService(ctx *bootstrap.Context, uc *data.NavigationRepo) *NavigationService {
	return &NavigationService{
		log:            ctx.NewLoggerHelper("navigation/service/core-service"),
		navigationRepo: uc,
	}
}

func (s *NavigationService) List(ctx context.Context, req *paginationV1.PagingRequest) (*siteV1.ListNavigationResponse, error) {
	return s.navigationRepo.List(ctx, req)
}

func (s *NavigationService) Get(ctx context.Context, req *siteV1.GetNavigationRequest) (*siteV1.Navigation, error) {
	return s.navigationRepo.Get(ctx, req)
}

func (s *NavigationService) Create(ctx context.Context, req *siteV1.CreateNavigationRequest) (*siteV1.Navigation, error) {
	return s.navigationRepo.Create(ctx, req)
}

func (s *NavigationService) Update(ctx context.Context, req *siteV1.UpdateNavigationRequest) (*siteV1.Navigation, error) {
	return s.navigationRepo.Update(ctx, req)
}

func (s *NavigationService) Delete(ctx context.Context, req *siteV1.DeleteNavigationRequest) (*emptypb.Empty, error) {
	err := s.navigationRepo.Delete(ctx, req)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
