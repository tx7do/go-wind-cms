package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	appV1 "go-wind-cms/api/gen/go/app/service/v1"
	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
)

type CategoryService struct {
	appV1.CategoryServiceHTTPServer

	categoryClient contentV1.CategoryServiceClient
	log            *log.Helper
}

func NewCategoryService(ctx *bootstrap.Context, categoryClient contentV1.CategoryServiceClient) *CategoryService {
	return &CategoryService{
		log:            ctx.NewLoggerHelper("category/service/app-service"),
		categoryClient: categoryClient,
	}
}

func (s *CategoryService) List(ctx context.Context, req *paginationV1.PagingRequest) (*contentV1.ListCategoryResponse, error) {
	resp, err := s.categoryClient.List(ctx, req)
	if err != nil {
		return nil, err
	}
	// 公开端点仅返回启用（前台可见）分类，过滤隐藏/归档分类
	if resp != nil {
		filtered := make([]*contentV1.Category, 0, len(resp.GetItems()))
		for _, c := range resp.GetItems() {
			if c != nil && c.GetStatus() == contentV1.Category_CATEGORY_STATUS_ACTIVE {
				filtered = append(filtered, c)
			}
		}
		resp.Items = filtered
		resp.Total = uint64(len(filtered))
	}
	return resp, nil
}

func (s *CategoryService) Get(ctx context.Context, req *contentV1.GetCategoryRequest) (*contentV1.Category, error) {
	resp, err := s.categoryClient.Get(ctx, req)
	if err != nil {
		return nil, err
	}
	// 公开端点仅返回启用（前台可见）分类
	if resp == nil || resp.GetStatus() != contentV1.Category_CATEGORY_STATUS_ACTIVE {
		return nil, contentV1.ErrorNotFound("category not found")
	}
	return resp, nil
}

// Create/Update/Delete 在 app（公开站点）服务上禁用：CMS 内容的写操作应经由 admin 服务。
func (s *CategoryService) Create(_ context.Context, _ *contentV1.CreateCategoryRequest) (*contentV1.Category, error) {
	return nil, contentV1.ErrorForbidden("content mutation is not allowed on the public app service")
}

func (s *CategoryService) Update(_ context.Context, _ *contentV1.UpdateCategoryRequest) (*contentV1.Category, error) {
	return nil, contentV1.ErrorForbidden("content mutation is not allowed on the public app service")
}

func (s *CategoryService) Delete(_ context.Context, _ *contentV1.DeleteCategoryRequest) (*emptypb.Empty, error) {
	return nil, contentV1.ErrorForbidden("content mutation is not allowed on the public app service")
}

func (s *CategoryService) GetTranslation(ctx context.Context, req *contentV1.GetCategoryRequest) (*contentV1.CategoryTranslation, error) {
	return s.categoryClient.GetTranslation(ctx, req)
}
