package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"google.golang.org/protobuf/types/known/emptypb"

	identityV1 "go-wind-cms/api/gen/go/identity/service/v1"
	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/privacy"
	"go-wind-cms/app/core/service/internal/data/ent/tenantquota"
	appViewer "go-wind-cms/pkg/entgo/viewer"
)

// ─── 配额套餐模板（平台级命名模板，一键套用到租户）───
//
// 模板管理走平台管理员视角（sys_quota_plans 行 tenant_id=0）；Apply 用系统
// 权限跨租户写目标租户的 sys_tenant_quotas（逐项 upsert，模板未含的类型不动——
// 已有独立配置保留）。

// ListQuotaPlans 列出模板。
func (s *TenantQuotaService) ListQuotaPlans(ctx context.Context, req *identityV1.ListQuotaPlansRequest) (*identityV1.ListQuotaPlansResponse, error) {
	if s.quotaPlanRepo == nil {
		return nil, errors.InternalServer("QUOTA_PLAN", "quota plan repo unavailable")
	}
	items, err := s.quotaPlanRepo.List(appViewer.NewSystemViewerContext(ctx))
	if err != nil {
		return nil, err
	}
	return &identityV1.ListQuotaPlansResponse{Items: items}, nil
}

// CreateQuotaPlan 创建模板。
func (s *TenantQuotaService) CreateQuotaPlan(ctx context.Context, req *identityV1.CreateQuotaPlanRequest) (*identityV1.QuotaPlan, error) {
	if s.quotaPlanRepo == nil {
		return nil, errors.InternalServer("QUOTA_PLAN", "quota plan repo unavailable")
	}
	if req == nil || req.Data == nil || req.GetData().GetName() == "" {
		return nil, errors.BadRequest("QUOTA_PLAN", "name is required")
	}
	if len(req.GetData().GetItems()) == 0 {
		return nil, errors.BadRequest("QUOTA_PLAN", "at least one quota item is required")
	}
	for _, it := range req.GetData().GetItems() {
		if it.GetQuotaType() == identityV1.TenantQuota_QUOTA_TYPE_UNSPECIFIED {
			return nil, errors.BadRequest("QUOTA_PLAN", "quota item type is required")
		}
		if it.GetLimit() < 0 {
			return nil, errors.BadRequest("QUOTA_PLAN", "quota item limit must be >= 0")
		}
	}
	operatorID, _ := viewerUserIDFromContext(ctx)
	id, err := s.quotaPlanRepo.Create(appViewer.NewSystemViewerContext(ctx), req.GetData(), operatorID)
	if err != nil {
		return nil, err
	}
	return s.quotaPlanRepo.Get(appViewer.NewSystemViewerContext(ctx), id)
}

// DeleteQuotaPlan 删除模板（已套用的租户配额是落库快照，不受影响）。
func (s *TenantQuotaService) DeleteQuotaPlan(ctx context.Context, req *identityV1.DeleteQuotaPlanRequest) (*emptypb.Empty, error) {
	if s.quotaPlanRepo == nil {
		return nil, errors.InternalServer("QUOTA_PLAN", "quota plan repo unavailable")
	}
	if req.GetId() == 0 {
		return nil, errors.BadRequest("QUOTA_PLAN", "id is required")
	}
	if err := s.quotaPlanRepo.Delete(appViewer.NewSystemViewerContext(ctx), req.GetId()); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// ApplyQuotaPlan 套用模板：逐项 upsert 目标租户配额，返回写入项数。
func (s *TenantQuotaService) ApplyQuotaPlan(ctx context.Context, req *identityV1.ApplyQuotaPlanRequest) (*identityV1.ApplyQuotaPlanResponse, error) {
	if s.quotaPlanRepo == nil {
		return nil, errors.InternalServer("QUOTA_PLAN", "quota plan repo unavailable")
	}
	if req.GetPlanId() == 0 || req.GetTenantId() == 0 {
		return nil, errors.BadRequest("QUOTA_PLAN", "plan id and tenant id are required")
	}
	plan, err := s.quotaPlanRepo.Get(appViewer.NewSystemViewerContext(ctx), req.GetPlanId())
	if err != nil {
		return nil, err
	}

	// 现有配额按类型建索引（模板未含的类型保留不动）；ListByTenant 返回 ent 行
	existing, err := s.repo.ListByTenant(appViewer.NewSystemViewerContext(ctx), req.GetTenantId())
	if err != nil {
		return nil, err
	}
	quotaTypeOf := func(e *ent.TenantQuota) identityV1.TenantQuota_QuotaType {
		switch e.QuotaType {
		case tenantquota.QuotaTypeUserLimit:
			return identityV1.TenantQuota_USER_LIMIT
		case tenantquota.QuotaTypeStorage:
			return identityV1.TenantQuota_STORAGE
		case tenantquota.QuotaTypeApiCall:
			return identityV1.TenantQuota_API_CALL
		case tenantquota.QuotaTypeAiTokens:
			return identityV1.TenantQuota_AI_TOKENS
		default:
			return identityV1.TenantQuota_QUOTA_TYPE_UNSPECIFIED
		}
	}
	byType := make(map[identityV1.TenantQuota_QuotaType]*ent.TenantQuota, len(existing))
	for _, e := range existing {
		byType[quotaTypeOf(e)] = e
	}

	sysCtx := privacy.DecisionContext(appViewer.NewSystemViewerContext(ctx), privacy.Allow)
	applied := uint32(0)
	for _, item := range plan.GetItems() {
		limit := uint64(item.GetLimit())
		if cur, ok := byType[item.GetQuotaType()]; ok {
			if cur.QuotaValue == limit {
				continue // 同值跳过，不刷 updated_at
			}
			if uerr := s.repo.Update(sysCtx, cur.ID, &identityV1.TenantQuota{
				QuotaValue: &limit,
			}, 0); uerr != nil {
				return nil, uerr
			}
		} else {
			if _, cerr := s.repo.Create(sysCtx, &identityV1.TenantQuota{
				TenantId:   &req.TenantId,
				QuotaType:  item.QuotaType.Enum(),
				QuotaValue: &limit,
			}, 0); cerr != nil {
				return nil, cerr
			}
		}
		applied++
	}
	s.log.Infof("quota plan [%d] applied to tenant [%d]: %d item(s)", req.GetPlanId(), req.GetTenantId(), applied)
	return &identityV1.ApplyQuotaPlanResponse{Applied: applied}, nil
}
