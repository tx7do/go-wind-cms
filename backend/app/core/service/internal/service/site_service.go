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

type SiteService struct {
	siteV1.UnimplementedSiteServiceServer

	siteSettingRepo *data.SiteRepo
	log             *log.Helper
}

func NewSiteService(ctx *bootstrap.Context, uc *data.SiteRepo) *SiteService {
	return &SiteService{
		log:             ctx.NewLoggerHelper("site/service/core-service"),
		siteSettingRepo: uc,
	}
}

func (s *SiteService) List(ctx context.Context, req *paginationV1.PagingRequest) (*siteV1.ListSiteResponse, error) {
	return s.siteSettingRepo.List(ctx, req)
}

func (s *SiteService) Get(ctx context.Context, req *siteV1.GetSiteRequest) (*siteV1.Site, error) {
	return s.siteSettingRepo.Get(ctx, req)
}

func (s *SiteService) Create(ctx context.Context, req *siteV1.CreateSiteRequest) (*siteV1.Site, error) {
	return s.siteSettingRepo.Create(ctx, req)
}

func (s *SiteService) Update(ctx context.Context, req *siteV1.UpdateSiteRequest) (*siteV1.Site, error) {
	return s.siteSettingRepo.Update(ctx, req)
}

func (s *SiteService) Delete(ctx context.Context, req *siteV1.DeleteSiteRequest) (*emptypb.Empty, error) {
	_, err := s.siteSettingRepo.Delete(ctx, req)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
