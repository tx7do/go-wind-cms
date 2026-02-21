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

type NavigationService struct {
	adminV1.NavigationServiceHTTPServer

	navigationServiceClient siteV1.NavigationServiceClient
	log                     *log.Helper
}

func NewNavigationService(ctx *bootstrap.Context, navigationServiceClient siteV1.NavigationServiceClient) *NavigationService {
	return &NavigationService{
		log:                     ctx.NewLoggerHelper("navigation/service/admin-service"),
		navigationServiceClient: navigationServiceClient,
	}
}

func (s *NavigationService) List(ctx context.Context, req *paginationV1.PagingRequest) (*siteV1.ListNavigationResponse, error) {
	return s.navigationServiceClient.List(ctx, req)
}

func (s *NavigationService) Get(ctx context.Context, req *siteV1.GetNavigationRequest) (*siteV1.Navigation, error) {
	return s.navigationServiceClient.Get(ctx, req)
}

func (s *NavigationService) Create(ctx context.Context, req *siteV1.CreateNavigationRequest) (*siteV1.Navigation, error) {
	if req.Data == nil {
		return nil, adminV1.ErrorBadRequest("invalid parameter")
	}

	// 获取操作人信息
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.Data.CreatedBy = trans.Ptr(operator.UserId)

	return s.navigationServiceClient.Create(ctx, req)
}

func (s *NavigationService) Update(ctx context.Context, req *siteV1.UpdateNavigationRequest) (*siteV1.Navigation, error) {
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

	return s.navigationServiceClient.Update(ctx, req)
}

func (s *NavigationService) Delete(ctx context.Context, req *siteV1.DeleteNavigationRequest) (*emptypb.Empty, error) {
	return s.navigationServiceClient.Delete(ctx, req)
}
