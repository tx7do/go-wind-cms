package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"

	"go-wind-cms/pkg/middleware/auth"
)

// requirePlatformAdmin 把平台专用配置与台账挡在平台管理员之内。
//
// 租户与套餐同登记在 MODULE_TENANT、菜单/API/权限登记在 MODULE_PERMISSION，
// 这些模块因租户自助与角色抽屉读权限树而常态在租户套餐白名单内，叠加
// authz.type 默认 noop 全放行，意味着默认部署下租户操作者即可打到平台专用
// 变更端点。一道护栏不该靠漏配存在，判定显式放这里（移植自上游 da2b3a5f）。
func requirePlatformAdmin(ctx context.Context, log *log.Helper, resource string) error {
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return err
	}

	if operator.GetIsPlatformAdmin() {
		return nil
	}

	log.Errorf("operator [%d] (tenant %d) is not a platform administrator, %s access denied",
		operator.GetUserId(), operator.GetTenantId(), resource)

	return adminV1.ErrorForbidden("%s is platform-level configuration, only platform administrators may access it", resource)
}
