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
	siteV1 "go-wind-cms/api/gen/go/site/service/v1"
)

type SiteSettingService struct {
	adminV1.SiteSettingServiceHTTPServer

	siteSettingServiceClient siteV1.SiteSettingServiceClient
	log                      *log.Helper
}

func NewSiteSettingService(ctx *bootstrap.Context, siteSettingServiceClient siteV1.SiteSettingServiceClient) *SiteSettingService {
	return &SiteSettingService{
		log:                      ctx.NewLoggerHelper("site-setting/service/admin-service"),
		siteSettingServiceClient: siteSettingServiceClient,
	}
}

func (s *SiteSettingService) List(ctx context.Context, req *paginationV1.PagingRequest) (*siteV1.ListSiteSettingResponse, error) {
	return s.siteSettingServiceClient.List(ctx, req)
}

func (s *SiteSettingService) Get(ctx context.Context, req *siteV1.GetSiteSettingRequest) (*siteV1.SiteSetting, error) {
	return s.siteSettingServiceClient.Get(ctx, req)
}

func (s *SiteSettingService) Create(ctx context.Context, req *siteV1.CreateSiteSettingRequest) (*siteV1.SiteSetting, error) {
	if req.Data == nil {
		return nil, adminV1.ErrorBadRequest("invalid parameter")
	}

	// 获取操作人信息
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.Data.CreatedBy = trans.Ptr(operator.UserId)

	return s.siteSettingServiceClient.Create(ctx, req)
}

func (s *SiteSettingService) Update(ctx context.Context, req *siteV1.UpdateSiteSettingRequest) (*siteV1.SiteSetting, error) {
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

	return s.siteSettingServiceClient.Update(ctx, req)
}

func (s *SiteSettingService) Delete(ctx context.Context, req *siteV1.DeleteSiteSettingRequest) (*emptypb.Empty, error) {
	return s.siteSettingServiceClient.Delete(ctx, req)
}
