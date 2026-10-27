package service

import (
	"context"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
	"github.com/tx7do/go-crud/viewer"

	"go-wind-cms/app/core/service/internal/data"
)

// ReviewPolicy 发布审批多级策略：级数与各级角色来自系统参数
//（review.post.levels / review.post.levelRoles，configservice 懒加载+失效广播）。
// roleCodes 解析器由 wiring 注入（userRepo/roleRepo 在身份域，避免循环依赖）。
type ReviewPolicy struct {
	configRepo *data.ConfigRepo
	roleCodes  func(ctx context.Context, userID uint32) []string
	log        *log.Helper
}

func NewReviewPolicy(ctx *bootstrap.Context, configRepo *data.ConfigRepo, roleCodes func(ctx context.Context, userID uint32) []string) *ReviewPolicy {
	return &ReviewPolicy{
		configRepo: configRepo,
		roleCodes:  roleCodes,
		log:        ctx.NewLoggerHelper("review-policy/service/core-service"),
	}
}

// Levels 发布审批级数（钳位 1..8；配置缺失/非法回落 1=单级，历史行为）。
func (p *ReviewPolicy) Levels(ctx context.Context, tenantID uint32) int {
	if p == nil || p.configRepo == nil {
		return 1
	}
	n := p.configRepo.GetConfigInt(ctx, "review.post.levels", 1)
	if n < 1 {
		n = 1
	}
	if n > 8 {
		n = 8
	}
	return n
}

// levelRoles 第 level 级（1 起）允许的审批角色 code 列表（空切片=不校验角色）。
func (p *ReviewPolicy) levelRoles(ctx context.Context, tenantID uint32, level int) []string {
	if p == nil || p.configRepo == nil {
		return nil
	}
	raw := p.configRepo.GetConfigString(ctx, "review.post.levelRoles", "")
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	// 逗号列表按级对应：超出配置长度的级别沿用最后一级的角色
	idx := level - 1
	if idx < 0 {
		idx = 0
	}
	if idx >= len(parts) {
		idx = len(parts) - 1
	}
	roles := make([]string, 0, 2)
	for _, r := range strings.Split(parts[idx], "|") {
		if t := strings.TrimSpace(r); t != "" {
			roles = append(roles, t)
		}
	}
	return roles
}

// CheckApprover 校验操作者可审第 level 级（配置了角色才校验；无匹配角色 403）。
func (p *ReviewPolicy) CheckApprover(ctx context.Context, tenantID uint32, userID uint32, level int) error {
	required := p.levelRoles(ctx, tenantID, level)
	if len(required) == 0 || p.roleCodes == nil {
		return nil
	}
	codes := p.roleCodes(ctx, userID)
	for _, have := range codes {
		for _, want := range required {
			if strings.EqualFold(have, want) {
				return nil
			}
		}
	}
	p.log.Warnf("review approver [%d] lacks required role for level %d (have=%v want=%v)",
		userID, level, codes, required)
	return contentV1.ErrorForbidden("insufficient review role for this level")
}

// operatorTenantUserID 从 viewer 取操作者与租户（审批留痕与校验用）。
func operatorTenantUserID(ctx context.Context) (userID, tenantID uint32, ok bool) {
	vc, exist := viewer.FromContext(ctx)
	if !exist || vc == nil {
		return 0, 0, false
	}
	return uint32(vc.UserID()), uint32(vc.TenantID()), true
}
