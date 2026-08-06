package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	appV1 "go-wind-cms/api/gen/go/app/service/v1"
	siteV1 "go-wind-cms/api/gen/go/site/service/v1"
)

type NavigationService struct {
	appV1.NavigationServiceHTTPServer

	categoryClient siteV1.NavigationServiceClient
	log            *log.Helper
}

func NewNavigationService(ctx *bootstrap.Context, categoryClient siteV1.NavigationServiceClient) *NavigationService {
	return &NavigationService{
		log:            ctx.NewLoggerHelper("Navigation/service/app-service"),
		categoryClient: categoryClient,
	}
}

func (s *NavigationService) List(ctx context.Context, req *paginationV1.PagingRequest) (*siteV1.ListNavigationResponse, error) {
	return s.categoryClient.List(ctx, req)
}

func (s *NavigationService) Get(ctx context.Context, req *siteV1.GetNavigationRequest) (*siteV1.Navigation, error) {
	return s.categoryClient.Get(ctx, req)
}

// Create/Update/Delete 在 app（公开站点）服务上禁用：站点导航的写操作应经由 admin 服务。
func (s *NavigationService) Create(_ context.Context, _ *siteV1.CreateNavigationRequest) (*siteV1.Navigation, error) {
	return nil, siteV1.ErrorForbidden("content mutation is not allowed on the public app service")
}

func (s *NavigationService) Update(_ context.Context, _ *siteV1.UpdateNavigationRequest) (*siteV1.Navigation, error) {
	return nil, siteV1.ErrorForbidden("content mutation is not allowed on the public app service")
}

func (s *NavigationService) Delete(_ context.Context, _ *siteV1.DeleteNavigationRequest) (*emptypb.Empty, error) {
	return nil, siteV1.ErrorForbidden("content mutation is not allowed on the public app service")
}
