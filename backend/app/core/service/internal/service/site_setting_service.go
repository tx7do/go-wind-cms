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

type SiteSettingService struct {
	siteV1.UnimplementedSiteSettingServiceServer

	siteSettingRepo *data.SiteSettingRepo
	log             *log.Helper
}

func NewSiteSettingService(ctx *bootstrap.Context, uc *data.SiteSettingRepo) *SiteSettingService {
	return &SiteSettingService{
		log:             ctx.NewLoggerHelper("site-setting/service/core-service"),
		siteSettingRepo: uc,
	}
}

func (s *SiteSettingService) List(ctx context.Context, req *paginationV1.PagingRequest) (*siteV1.ListSiteSettingResponse, error) {
	return s.siteSettingRepo.List(ctx, req)
}

func (s *SiteSettingService) Get(ctx context.Context, req *siteV1.GetSiteSettingRequest) (*siteV1.SiteSetting, error) {
	return s.siteSettingRepo.Get(ctx, req)
}

func (s *SiteSettingService) Create(ctx context.Context, req *siteV1.CreateSiteSettingRequest) (*siteV1.SiteSetting, error) {
	return s.siteSettingRepo.Create(ctx, req)
}

func (s *SiteSettingService) Update(ctx context.Context, req *siteV1.UpdateSiteSettingRequest) (*siteV1.SiteSetting, error) {
	return s.siteSettingRepo.Update(ctx, req)
}

func (s *SiteSettingService) Delete(ctx context.Context, req *siteV1.DeleteSiteSettingRequest) (*emptypb.Empty, error) {
	_, err := s.siteSettingRepo.Delete(ctx, req)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
