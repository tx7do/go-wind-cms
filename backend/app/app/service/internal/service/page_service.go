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

type PageService struct {
	appV1.PageServiceHTTPServer

	pageServiceClient contentV1.PageServiceClient
	log               *log.Helper
}

func NewPageService(ctx *bootstrap.Context, pageServiceClient contentV1.PageServiceClient) *PageService {
	return &PageService{
		log:               ctx.NewLoggerHelper("page/service/app-service"),
		pageServiceClient: pageServiceClient,
	}
}

func (s *PageService) List(ctx context.Context, req *paginationV1.PagingRequest) (*contentV1.ListPageResponse, error) {
	resp, err := s.pageServiceClient.List(ctx, req)
	if err != nil {
		return nil, err
	}
	// 公开端点仅返回已发布页面，过滤草稿/归档/私有等状态
	if resp != nil {
		filtered := make([]*contentV1.Page, 0, len(resp.GetItems()))
		for _, p := range resp.GetItems() {
			if p != nil && p.GetStatus() == contentV1.Page_PAGE_STATUS_PUBLISHED {
				filtered = append(filtered, p)
			}
		}
		resp.Items = filtered
		resp.Total = uint64(len(filtered))
	}
	return resp, nil
}

func (s *PageService) Get(ctx context.Context, req *contentV1.GetPageRequest) (*contentV1.Page, error) {
	resp, err := s.pageServiceClient.Get(ctx, req)
	if err != nil {
		return nil, err
	}
	// 公开端点仅返回已发布页面，草稿/归档/私有等状态按未找到处理
	if resp == nil || resp.GetStatus() != contentV1.Page_PAGE_STATUS_PUBLISHED {
		return nil, contentV1.ErrorNotFound("page not found")
	}
	return resp, nil
}

// Create/Update/Delete 在 app（公开站点）服务上禁用：CMS 内容的写操作应经由 admin 服务。
func (s *PageService) Create(_ context.Context, _ *contentV1.CreatePageRequest) (*contentV1.Page, error) {
	return nil, contentV1.ErrorForbidden("content mutation is not allowed on the public app service")
}

func (s *PageService) Update(_ context.Context, _ *contentV1.UpdatePageRequest) (*contentV1.Page, error) {
	return nil, contentV1.ErrorForbidden("content mutation is not allowed on the public app service")
}

func (s *PageService) Delete(_ context.Context, _ *contentV1.DeletePageRequest) (*emptypb.Empty, error) {
	return nil, contentV1.ErrorForbidden("content mutation is not allowed on the public app service")
}

func (s *PageService) GetTranslation(ctx context.Context, req *contentV1.GetPageRequest) (*contentV1.PageTranslation, error) {
	return s.pageServiceClient.GetTranslation(ctx, req)
}
