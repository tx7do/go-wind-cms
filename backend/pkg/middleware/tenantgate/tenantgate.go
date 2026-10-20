package tenantgate

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	ktransport "github.com/go-kratos/kratos/v2/transport"
	"google.golang.org/grpc"

	identityV1 "go-wind-cms/api/gen/go/identity/service/v1"
	"go-wind-cms/pkg/metadata"
)

// Checker 租户访问闸门的判定接口（对齐 core TenantQuotaService 的 gRPC 客户端面）。
type Checker interface {
	CheckTenantAccess(ctx context.Context, in *identityV1.CheckTenantAccessRequest, opts ...grpc.CallOption) (*identityV1.CheckTenantAccessResponse, error)
}

// Server 租户闸门中间件：每个已认证请求（白名单外的）经此做租户状态与
// API_CALL 配额检查。判定在 core（CheckTenantAccess），全链 fail-open——
// 配额缺失 = 未限售，基础设施抖动不打死租户业务；拒绝只在"配额明确存在
// 且已用满"时发生（core 侧返回 403）。
//
// 必须挂在 auth.Server 之后：依赖其注入 context 的操作者元数据（tenantId）。
// tenant==0（平台身份）跳过闸门；机器令牌继承凭证归属租户的闸门（与上游一致）。
func Server(checker Checker) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {
			if checker == nil {
				return handler(ctx, req)
			}

			if op, err := metadata.FromContext(ctx, true); err == nil && op != nil {
				if tenantID := uint32(op.GetTenantId()); tenantID != 0 {
					method, path := requestInfo(ctx)
					if _, err := checker.CheckTenantAccess(ctx, &identityV1.CheckTenantAccessRequest{
						TenantId: tenantID,
						Method:   method,
						Path:     path,
					}); err != nil {
						return nil, err
					}
				}
			}

			return handler(ctx, req)
		}
	}
}

// requestInfo 取当前请求的 HTTP 方法与路径（core 侧拒绝日志用；取不到不影响判定）。
func requestInfo(ctx context.Context) (method, path string) {
	tr, ok := ktransport.FromServerContext(ctx)
	if !ok {
		return "", ""
	}
	path = tr.Operation()
	if htr, ok := tr.(*khttp.Transport); ok && htr.Request() != nil {
		method = htr.Request().Method
	}
	return method, path
}
