package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	identityV1 "go-wind-cms/api/gen/go/identity/service/v1"

	"go-wind-cms/pkg/middleware/auth"
)

// TenantQuotaService 租户配额管理（BFF 薄代理：配额数据/用量聚合/闸门判定都在 core）。
type TenantQuotaService struct {
	adminV1.TenantQuotaServiceHTTPServer

	log               *log.Helper
	tenantQuotaClient identityV1.TenantQuotaServiceClient
}

func NewTenantQuotaService(
	ctx *bootstrap.Context,
	tenantQuotaClient identityV1.TenantQuotaServiceClient,
) *TenantQuotaService {
	l := log.NewHelper(log.With(ctx.GetLogger(), "module", "tenant-quota/service/admin-service"))
	return &TenantQuotaService{
		log:               l,
		tenantQuotaClient: tenantQuotaClient,
	}
}

// List 租户配额列表（透传 core）。
func (s *TenantQuotaService) List(ctx context.Context, req *paginationV1.PagingRequest) (*identityV1.ListTenantQuotaResponse, error) {
	return s.tenantQuotaClient.List(ctx, req)
}

// Get 租户配额详情（透传 core）。
func (s *TenantQuotaService) Get(ctx context.Context, req *identityV1.GetTenantQuotaRequest) (*identityV1.TenantQuota, error) {
	return s.tenantQuotaClient.Get(ctx, req)
}

// Create 创建配额（透传 core）。
func (s *TenantQuotaService) Create(ctx context.Context, req *identityV1.CreateTenantQuotaRequest) (*emptypb.Empty, error) {
	return s.tenantQuotaClient.Create(ctx, req)
}

// Update 更新配额（透传 core）。
func (s *TenantQuotaService) Update(ctx context.Context, req *identityV1.UpdateTenantQuotaRequest) (*emptypb.Empty, error) {
	return s.tenantQuotaClient.Update(ctx, req)
}

// Delete 删除配额（透传 core）。
func (s *TenantQuotaService) Delete(ctx context.Context, req *identityV1.DeleteTenantQuotaRequest) (*emptypb.Empty, error) {
	return s.tenantQuotaClient.Delete(ctx, req)
}

// GetUsage 指定租户用量与配额（平台管理员视图，透传 core）。
func (s *TenantQuotaService) GetUsage(ctx context.Context, req *identityV1.GetTenantUsageRequest) (*identityV1.TenantUsage, error) {
	return s.tenantQuotaClient.GetUsage(ctx, req)
}

// GetMyTenantUsage 当前租户用量与配额（租户自助视图）。
// 租户身份不跨服务透传：从认证上下文取出 tenant_id 显式传给 core。
func (s *TenantQuotaService) GetMyTenantUsage(ctx context.Context, req *identityV1.GetMyTenantUsageRequest) (*identityV1.TenantUsage, error) {
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	req.TenantId = operator.GetTenantId()
	return s.tenantQuotaClient.GetMyTenantUsage(ctx, req)
}

// ── 配额套餐模板（透传 core）──

func (s *TenantQuotaService) ListQuotaPlans(ctx context.Context, req *identityV1.ListQuotaPlansRequest) (*identityV1.ListQuotaPlansResponse, error) {
	return s.tenantQuotaClient.ListQuotaPlans(ctx, req)
}

func (s *TenantQuotaService) CreateQuotaPlan(ctx context.Context, req *identityV1.CreateQuotaPlanRequest) (*identityV1.QuotaPlan, error) {
	return s.tenantQuotaClient.CreateQuotaPlan(ctx, req)
}

func (s *TenantQuotaService) DeleteQuotaPlan(ctx context.Context, req *identityV1.DeleteQuotaPlanRequest) (*emptypb.Empty, error) {
	return s.tenantQuotaClient.DeleteQuotaPlan(ctx, req)
}

func (s *TenantQuotaService) ApplyQuotaPlan(ctx context.Context, req *identityV1.ApplyQuotaPlanRequest) (*identityV1.ApplyQuotaPlanResponse, error) {
	return s.tenantQuotaClient.ApplyQuotaPlan(ctx, req)
}
