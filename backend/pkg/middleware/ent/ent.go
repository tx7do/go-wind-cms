package ent

import (
	"context"
	"reflect"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"

	"github.com/tx7do/go-crud/viewer"
	"go.opentelemetry.io/otel/trace"

	authenticationV1 "go-wind-cms/api/gen/go/authentication/service/v1"

	appViewer "go-wind-cms/pkg/entgo/viewer"
	"go-wind-cms/pkg/metadata"
)

// TenantResolver 按 Host 解析 tenant_id。
//
// 由 app BFF 注入（实现为对 core TenantService.ResolveTenantByDomain 的转发）。
// 仅用于匿名（白名单）请求：无 OperatorMetadata 时，按请求 Host 解析出真实
// tenant_id，注入只读 tenant-scoped viewer，使公开读按 tenant 隔离而非跨租户。
// admin/core 不注入，保持原 SystemViewer 兜底。
type TenantResolver interface {
	ResolveTenantIDByDomain(ctx context.Context, domain string) (uint32, error)
}

type entMiddlewareOptions struct {
	tenantResolver TenantResolver
}

// Option 配置 ent 中间件。
type Option func(*entMiddlewareOptions)

// WithTenantResolver 注入 tenant 解析器。注入后，匿名（白名单）请求按 Host
// 解析 tenant_id 并注入 AnonymousTenantViewer；解析失败 fail-closed 注入
// noopContext（拒绝查询），不再回退 SystemViewer。未注入时保持原 SystemViewer
// 兜底行为（admin/core 场景）。
func WithTenantResolver(r TenantResolver) Option {
	return func(o *entMiddlewareOptions) {
		o.tenantResolver = r
	}
}

// Server 设置 Ent Viewer 到 Context 中的中间件。
//
// 已认证请求（携带 OperatorMetadata）构建 UserViewer。
// 匿名请求（无 OperatorMetadata）：
//   - 注入了 TenantResolver（app BFF）：按 Host 解析 tenant_id，成功注入
//     AnonymousTenantViewer；失败注入 noopContext fail-closed。
//   - 未注入（admin/core）：注入 SystemViewer 兜底，保持原行为。
func Server(opts ...Option) middleware.Middleware {
	o := entMiddlewareOptions{}
	for _, opt := range opts {
		opt(&o)
	}

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {
			md, err := metadata.FromServerContext(ctx)
			if err != nil {
				reqType := "<nil>"
				if req != nil {
					reqType = reflect.TypeOf(req).String()
				}
				t, _ := transport.FromServerContext(ctx)
				kind, operation, endpoint := "", "", ""
				if t != nil {
					kind = string(t.Kind())
					operation = t.Operation()
					endpoint = t.Endpoint()
				}
				log.Errorf("ent middleware: failed to get metadata from context: %v; req_type=%s; transport=%s; operation=%s; endpoint=%s",
					err, reqType, kind, operation, endpoint)
			}

			var traceID string
			spanContext := trace.SpanContextFromContext(ctx)
			if spanContext.HasTraceID() {
				traceID = spanContext.TraceID().String()
			}

			if md == nil {
				ctx = viewer.WithContext(ctx, resolveAnonymousViewer(ctx, o.tenantResolver, traceID))
				return handler(ctx, req)
			}

			ctx = viewer.WithContext(ctx, metaDataToUserViewerContext(md, traceID))

			return handler(ctx, req)
		}
	}
}

// resolveAnonymousViewer 决定匿名（无 OperatorMetadata）请求注入的 viewer。
//
// 注入了 resolver：按请求 Host 解析 tenant_id。成功 → AnonymousTenantViewer
// （按 tenant 隔离、只读）；失败/未配 domain → noopContext fail-closed（拒绝查询，
// 不回退 SystemViewer 避免跨租户泄漏）。
// 未注入 resolver（admin/core）：保持 SystemViewer 兜底。
func resolveAnonymousViewer(ctx context.Context, resolver TenantResolver, traceID string) viewer.Context {
	if resolver == nil {
		return appViewer.NewSystemViewer()
	}
	host := hostFromContext(ctx)
	if host == "" {
		log.Warnf("ent middleware: anonymous request without resolvable host, fail-closed")
		return viewer.NewNoopContext()
	}
	tid, err := resolver.ResolveTenantIDByDomain(ctx, host)
	if err != nil || tid == 0 {
		log.Warnf("ent middleware: tenant resolve failed for host %q, fail-closed: %v", host, err)
		return viewer.NewNoopContext()
	}
	return appViewer.NewAnonymousTenantViewer(uint64(tid), traceID)
}

// hostFromContext 从 server transport 提取请求 Host（仅 HTTP transport 有意义）。
func hostFromContext(ctx context.Context) string {
	hr, ok := http.RequestFromServerContext(ctx)
	if !ok || hr == nil {
		return ""
	}
	return hr.Host
}

func metaDataToUserViewerContext(md *authenticationV1.OperatorMetadata, traceID string) viewer.Context {
	if md == nil {
		return nil
	}

	userViewer := appViewer.NewUserViewer(
		md.GetUserId(),
		md.GetTenantId(),
		md.GetOrgUnitId(),
		traceID,
		md.GetDataScope(),
	)
	return userViewer
}
