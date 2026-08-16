package data

import (
	"context"

	identityV1 "go-wind-cms/api/gen/go/identity/service/v1"
	entmiddleware "go-wind-cms/pkg/middleware/ent"
)

// tenantResolver 把 core TenantServiceClient 适配为 ent 中间件的 TenantResolver。
//
// 供 app BFF 在匿名（白名单）请求中按 Host 解析 tenant_id。仅转发
// ResolveTenantByDomain，不暴露其他 tenant 查询面。
type tenantResolver struct {
	client identityV1.TenantServiceClient
}

// NewTenantResolver 构造注入 ent 中间件的 tenant 解析器。
func NewTenantResolver(client identityV1.TenantServiceClient) entmiddleware.TenantResolver {
	if client == nil {
		return nil
	}
	return tenantResolver{client: client}
}

func (r tenantResolver) ResolveTenantIDByDomain(ctx context.Context, domain string) (uint32, error) {
	resp, err := r.client.ResolveTenantByDomain(ctx, &identityV1.ResolveTenantByDomainRequest{Domain: domain})
	if err != nil {
		return 0, err
	}
	if resp == nil {
		return 0, nil
	}
	return resp.GetTenantId(), nil
}
