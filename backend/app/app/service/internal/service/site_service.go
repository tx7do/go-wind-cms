package service

import (
	"context"
	"net"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"github.com/go-kratos/kratos/v2/transport/http"
	"google.golang.org/protobuf/types/known/emptypb"

	appV1 "go-wind-cms/api/gen/go/app/service/v1"
	siteV1 "go-wind-cms/api/gen/go/site/service/v1"
)

type SiteService struct {
	appV1.SiteServiceHTTPServer

	siteClient siteV1.SiteServiceClient
	log        *log.Helper
}

func NewSiteService(ctx *bootstrap.Context, siteClient siteV1.SiteServiceClient) *SiteService {
	return &SiteService{
		log:        ctx.NewLoggerHelper("site/service/app-service"),
		siteClient: siteClient,
	}
}

// GetSiteByDomain 按当前请求 Host 解析域名，向 core 查询当前站点的渲染必需配置。
//
// domain 不接受调用方传入，始终由请求 Host 头提取（与 ent middleware 的 tenant 解析
// 一致），保证前台只能拿到"自己所在域名"的站点配置。core 端返回前已裁剪非渲染字段。
func (s *SiteService) GetSiteByDomain(ctx context.Context, req *siteV1.GetSiteByDomainRequest) (*siteV1.Site, error) {
	if req == nil {
		return nil, siteV1.ErrorBadRequest("invalid parameter")
	}

	host := hostFromContext(ctx)
	if host == "" {
		return nil, siteV1.ErrorBadRequest("resolvable host is required")
	}

	// 去端口：HTTP Host 形如 example.com:8080，只取域名部分。
	domain := host
	if h, _, err := net.SplitHostPort(host); err == nil && h != "" {
		domain = h
	}

	req.Domain = domain

	return s.siteClient.GetSiteByDomain(ctx, req)
}

// hostFromContext 从 server transport 提取请求 Host（仅 HTTP transport 有意义）。
func hostFromContext(ctx context.Context) string {
	hr, ok := http.RequestFromServerContext(ctx)
	if !ok || hr == nil {
		return ""
	}
	return hr.Host
}

// List/Create/Update/Delete 在 app（公开站点）服务上禁用：站点配置的写操作应经由 admin 服务。
func (s *SiteService) List(_ context.Context, _ *paginationV1.PagingRequest) (*siteV1.ListSiteResponse, error) {
	return nil, siteV1.ErrorForbidden("content mutation is not allowed on the public app service")
}

func (s *SiteService) Create(_ context.Context, _ *siteV1.CreateSiteRequest) (*siteV1.Site, error) {
	return nil, siteV1.ErrorForbidden("content mutation is not allowed on the public app service")
}

func (s *SiteService) Update(_ context.Context, _ *siteV1.UpdateSiteRequest) (*siteV1.Site, error) {
	return nil, siteV1.ErrorForbidden("content mutation is not allowed on the public app service")
}

func (s *SiteService) Delete(_ context.Context, _ *siteV1.DeleteSiteRequest) (*emptypb.Empty, error) {
	return nil, siteV1.ErrorForbidden("content mutation is not allowed on the public app service")
}
