package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	resourceV1 "go-wind-cms/api/gen/go/resource/service/v1"

	"go-wind-cms/pkg/middleware/auth"
)

type MenuService struct {
	adminV1.MenuServiceHTTPServer

	log *log.Helper

	menuServiceClient resourceV1.MenuServiceClient
}

func NewMenuService(ctx *bootstrap.Context, menuServiceClient resourceV1.MenuServiceClient) *MenuService {
	return &MenuService{
		log:               ctx.NewLoggerHelper("menu/service/admin-service"),
		menuServiceClient: menuServiceClient,
	}
}

func (s *MenuService) List(ctx context.Context, req *paginationV1.PagingRequest) (*resourceV1.ListMenuResponse, error) {
	ret, err := s.menuServiceClient.List(ctx, req)
	if err != nil {

		return nil, err
	}

	return ret, nil
}

func (s *MenuService) Get(ctx context.Context, req *resourceV1.GetMenuRequest) (*resourceV1.Menu, error) {
	ret, err := s.menuServiceClient.Get(ctx, req)
	if err != nil {

		return nil, err
	}

	return ret, nil
}

func (s *MenuService) Create(ctx context.Context, req *resourceV1.CreateMenuRequest) (*emptypb.Empty, error) {
	if req.Data == nil {
		return nil, adminV1.ErrorBadRequest("invalid parameter")
	}

	// 获取操作人信息
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.Data.CreatedBy = trans.Ptr(operator.UserId)

	return s.menuServiceClient.Create(ctx, req)
}

func (s *MenuService) Update(ctx context.Context, req *resourceV1.UpdateMenuRequest) (*emptypb.Empty, error) {
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

	return s.menuServiceClient.Update(ctx, req)
}

func (s *MenuService) Delete(ctx context.Context, req *resourceV1.DeleteMenuRequest) (*emptypb.Empty, error) {
	// 获取操作人信息
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.OperatorId = trans.Ptr(operator.UserId)

	return s.menuServiceClient.Delete(ctx, req)
}
