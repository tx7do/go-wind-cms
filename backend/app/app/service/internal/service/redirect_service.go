package service

import (
	"context"
	"strconv"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	kratosMetadata "github.com/go-kratos/kratos/v2/metadata"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	appV1 "go-wind-cms/api/gen/go/app/service/v1"
	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
	entmiddleware "go-wind-cms/pkg/middleware/ent"
)

// RedirectService 公开站点的重定向查询（匿名可调用，C 端路由中间件消费）。
// 租户归属与浏览/表单链路同款：按请求 Host 解析后经内网可信元数据注入。
type RedirectService struct {
	appV1.RedirectServiceHTTPServer

	redirectServiceClient contentV1.RedirectServiceClient
	tenantResolve         entmiddleware.TenantResolver
	log                   *log.Helper
}

func NewRedirectService(ctx *bootstrap.Context, redirectServiceClient contentV1.RedirectServiceClient, tenantResolve entmiddleware.TenantResolver) *RedirectService {
	return &RedirectService{
		log:                   ctx.NewLoggerHelper("redirect/service/app-service"),
		redirectServiceClient: redirectServiceClient,
		tenantResolve:         tenantResolve,
	}
}

// ListEnabled 全量启用项（列表量级小，直接全量；前端按 from_path 前缀/精确匹配）。
func (s *RedirectService) ListEnabled(ctx context.Context, req *contentV1.ListRedirectsRequest) (*contentV1.ListRedirectsResponse, error) {
	if s.tenantResolve != nil {
		if hr, ok := khttp.RequestFromServerContext(ctx); ok && hr != nil && hr.Host != "" {
			if tid, err := s.tenantResolve.ResolveTenantIDByDomain(ctx, hr.Host); err == nil && tid > 0 {
				ctx = kratosMetadata.AppendToClientContext(ctx, "x-md-global-tenant-id", strconv.FormatUint(uint64(tid), 10))
			}
		}
	}
	// core 无 ListEnabled 专用 RPC：全量拉取后客户端过滤启用项（量级小）
	resp, err := s.redirectServiceClient.List(ctx, &paginationV1.PagingRequest{
		Page:     trans.Ptr(uint32(1)),
		PageSize: trans.Ptr(uint32(500)),
	})
	if err != nil {
		return nil, err
	}
	enabled := make([]*contentV1.Redirect, 0, len(resp.GetItems()))
	for _, item := range resp.GetItems() {
		if item.GetEnabled() {
			enabled = append(enabled, item)
		}
	}
	return &contentV1.ListRedirectsResponse{Items: enabled, Total: uint64(len(enabled))}, nil
}
