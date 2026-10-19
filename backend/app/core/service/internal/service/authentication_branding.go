package service

import (
	"context"
	"strings"

	authenticationV1 "go-wind-cms/api/gen/go/authentication/service/v1"
	identityV1 "go-wind-cms/api/gen/go/identity/service/v1"
)

// GetTenantBranding 租户白标查询（免鉴权，登录前用）。
//
// 只暴露展示字段（name/logo_url）：免鉴权端点的暴露面必须最小，
// 运营/配额/管理员等字段一律不给。编号未命中返回 found=false（而非 404）——
// 登录页对"不存在"与"未设置白标"同等对待：保持默认品牌即可，不区别提示，
// 避免给枚举租户编号的探测者提供命中/未命中的区分信号。
func (s *AuthenticationService) GetTenantBranding(ctx context.Context, req *authenticationV1.GetTenantBrandingRequest) (*authenticationV1.GetTenantBrandingResponse, error) {
	if req == nil || strings.TrimSpace(req.GetCode()) == "" {
		return &authenticationV1.GetTenantBrandingResponse{Found: false}, nil
	}

	tenant, err := s.tenantRepo.Get(ctx, &identityV1.GetTenantRequest{
		QueryBy: &identityV1.GetTenantRequest_Code{Code: req.GetCode()},
	})
	if err != nil || tenant == nil {
		// 未命中不是错误（防枚举），只有基础设施故障才 500——但这里的 err
		// 来自 repo（已带日志），归一为 found=false 即可
		return &authenticationV1.GetTenantBrandingResponse{Found: false}, nil
	}

	return &authenticationV1.GetTenantBrandingResponse{
		Found:   true,
		Name:    tenant.GetName(),
		LogoUrl: tenant.GetLogoUrl(),
	}, nil
}
