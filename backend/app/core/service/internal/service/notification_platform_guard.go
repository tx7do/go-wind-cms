package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/tx7do/go-crud/viewer"
)

// requireNotificationPlatformContext 把通知域的配置与台账挡在平台视图之内。
//
// 渠道、路由规则、投递台账三张表都没有 tenant_id 列（形状与理由见
// ent/schema/notification_*.go 注释）：租户侧能读就读到全平台的 SMTP 主机名与
// 投递记录，能写就改得动全站的事件路由。一道护栏不该靠路由/权限漏配偶然存在，
// 所以判定显式放在服务层（与 interaction_admin_service 的 requireAdminOperator 同缝）。
func requireNotificationPlatformContext(ctx context.Context, resource string) error {
	vc, exist := viewer.FromContext(ctx)
	if !exist || vc == nil {
		return errors.Forbidden("NOTIFICATION", "operator identity required")
	}

	if vc.IsPlatformContext() || vc.IsSystemContext() {
		return nil
	}

	return errors.Forbidden("NOTIFICATION", resource+" is platform-level configuration, only platform administrators may access it")
}
