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

type NavigationItemService struct {
	adminV1.NavigationItemServiceHTTPServer

	navigationServiceClient siteV1.NavigationItemServiceClient
	log                     *log.Helper
}

func NewNavigationItemService(ctx *bootstrap.Context, navigationServiceClient siteV1.NavigationItemServiceClient) *NavigationItemService {
	return &NavigationItemService{
		log:                     ctx.NewLoggerHelper("navigation-item/service/admin-service"),
		navigationServiceClient: navigationServiceClient,
	}
}

func (s *NavigationItemService) List(ctx context.Context, req *paginationV1.PagingRequest) (*siteV1.ListNavigationItemResponse, error) {
	return s.navigationServiceClient.List(ctx, req)
}

func (s *NavigationItemService) Get(ctx context.Context, req *siteV1.GetNavigationItemRequest) (*siteV1.NavigationItem, error) {
	return s.navigationServiceClient.Get(ctx, req)
}

func (s *NavigationItemService) Create(ctx context.Context, req *siteV1.CreateNavigationItemRequest) (*siteV1.NavigationItem, error) {
	if req == nil || req.Data == nil {
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

func (s *NavigationItemService) Update(ctx context.Context, req *siteV1.UpdateNavigationItemRequest) (*siteV1.NavigationItem, error) {
	if req == nil || req.Data == nil {
		return nil, adminV1.ErrorBadRequest("invalid parameter")
	}

	// 获取操作人信息
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.Data.Id = trans.Ptr(req.GetId())

	req.Data.UpdatedBy = trans.Ptr(operator.GetUserId())
	if req.UpdateMask != nil {
		req.UpdateMask.Paths = append(req.UpdateMask.Paths, "updated_by")
	}

	return s.navigationServiceClient.Update(ctx, req)
}

func (s *NavigationItemService) Delete(ctx context.Context, req *siteV1.DeleteNavigationItemRequest) (*emptypb.Empty, error) {
	return s.navigationServiceClient.Delete(ctx, req)
}
