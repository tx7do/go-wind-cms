package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	siteV1 "go-wind-cms/api/gen/go/site/service/v1"

	"go-wind-cms/pkg/middleware/auth"
)

type SiteService struct {
	adminV1.SiteServiceHTTPServer

	siteServiceClient siteV1.SiteServiceClient
	log               *log.Helper
}

func NewSiteService(ctx *bootstrap.Context, siteServiceClient siteV1.SiteServiceClient) *SiteService {
	return &SiteService{
		log:               ctx.NewLoggerHelper("site/service/admin-service"),
		siteServiceClient: siteServiceClient,
	}
}

func (s *SiteService) List(ctx context.Context, req *paginationV1.PagingRequest) (*siteV1.ListSiteResponse, error) {
	return s.siteServiceClient.List(ctx, req)
}

func (s *SiteService) Get(ctx context.Context, req *siteV1.GetSiteRequest) (*siteV1.Site, error) {
	return s.siteServiceClient.Get(ctx, req)
}

func (s *SiteService) Create(ctx context.Context, req *siteV1.CreateSiteRequest) (*siteV1.Site, error) {
	if req.Data == nil {
		return nil, adminV1.ErrorBadRequest("invalid parameter")
	}

	// 获取操作人信息
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.Data.CreatedBy = trans.Ptr(operator.UserId)

	return s.siteServiceClient.Create(ctx, req)
}

func (s *SiteService) Update(ctx context.Context, req *siteV1.UpdateSiteRequest) (*siteV1.Site, error) {
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

	return s.siteServiceClient.Update(ctx, req)
}

func (s *SiteService) Delete(ctx context.Context, req *siteV1.DeleteSiteRequest) (*emptypb.Empty, error) {
	return s.siteServiceClient.Delete(ctx, req)
}
