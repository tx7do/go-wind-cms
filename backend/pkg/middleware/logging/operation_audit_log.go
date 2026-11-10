package logging

import (
	"context"
	"strings"

	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/tx7do/go-utils/trans"

	auditV1 "go-wind-cms/api/gen/go/audit/service/v1"

	appViewer "go-wind-cms/pkg/entgo/viewer"
)

type OperationAuditLogMiddleware struct {
	op *options
}

func NewOperationAuditLogMiddleware(op *options) *OperationAuditLogMiddleware {
	return &OperationAuditLogMiddleware{
		op: op,
	}
}

func (o *OperationAuditLogMiddleware) Name() string {
	return "OperationAuditLogMiddleware"
}

// Handle 记录写操作（POST/PUT/DELETE/PATCH）的操作审计日志。
// 读操作不记录（避免噪声）；登录/登出由 LoginAuditLogMiddleware 处理。
func (o *OperationAuditLogMiddleware) Handle(ctx context.Context, htr *http.Transport, middleErr error) {
	if o.op == nil || o.op.writeOperationLogFunc == nil {
		return
	}
	if htr == nil || htr.Request() == nil {
		return
	}

	method := strings.ToUpper(htr.Request().Method)
	action, ok := operationActionByMethod(method)
	if !ok {
		return
	}

	if htr.Operation() == o.op.loginOperation || htr.Operation() == o.op.logoutOperation {
		return
	}

	_, reason, success := getStatusCode(middleErr)

	operationLog := &auditV1.OperationAuditLog{}
	operationLog.Action = trans.Ptr(action)
	resourceType, resourceId := extractResource(htr.PathTemplate(), htr.Request().URL.Path)
	operationLog.ResourceType = trans.Ptr(resourceType)
	operationLog.ResourceId = trans.Ptr(resourceId)

	clientIp := getClientRealIP(htr.Request())
	operationLog.IpAddress = trans.Ptr(clientIp)
	operationLog.RequestId = trans.Ptr(getRequestId(htr.Request()))

	ut := extractAuthToken(htr)
	if ut != nil {
		operationLog.UserId = trans.Ptr(ut.UserId)
		operationLog.Username = ut.Username
	}

	operationLog.GeoLocation = fillGeoLocation(clientIp)

	operationLog.Success = trans.Ptr(success)
	if !success {
		operationLog.FailureReason = trans.Ptr(reason)
	}

	ctx = appViewer.NewSystemViewerContext(ctx)
	_ = o.op.writeOperationLogFunc(ctx, operationLog)
}

// operationActionByMethod 将 HTTP 方法映射为操作审计动作类型
func operationActionByMethod(method string) (auditV1.OperationAuditLog_ActionType, bool) {
	switch method {
	case "POST":
		return auditV1.OperationAuditLog_CREATE, true
	case "PUT", "PATCH":
		return auditV1.OperationAuditLog_UPDATE, true
	case "DELETE":
		return auditV1.OperationAuditLog_DELETE, true
	default:
		return auditV1.OperationAuditLog_ACTION_TYPE_UNSPECIFIED, false
	}
}

// extractResource 从路径模板与实际请求路径提取资源类型与资源ID。
// 模板末段为 {id} 占位符时，用实际请求路径同位置段作为资源ID。
// 例：模板 /admin/v1/users/{id} + 请求 /admin/v1/users/42 → ("users", "42")；
//
//	模板 /admin/v1/users + 请求 /admin/v1/users → ("users", "")
func extractResource(pathTemplate string, requestPath string) (resourceType string, resourceId string) {
	if pathTemplate == "" {
		return "", ""
	}

	tplSegments := strings.Split(strings.Trim(pathTemplate, "/"), "/")
	reqSegments := strings.Split(strings.Trim(requestPath, "/"), "/")

	// 对齐前缀：跳过模块段(admin/app)与版本段(v1),定位资源起始下标
	offset := 0
	for i, seg := range tplSegments {
		if seg == "v1" {
			offset = i + 1
			break
		}
		_ = i
	}
	if offset >= len(tplSegments) || offset >= len(reqSegments) {
		return "", ""
	}

	tplRes := tplSegments[offset:]
	reqRes := reqSegments[offset:]

	if len(tplRes) == 0 {
		return "", ""
	}

	resourceType = tplRes[0]
	// 自定义动词后缀(如 tasks:control)只保留资源段
	if idx := strings.Index(resourceType, ":"); idx > 0 {
		resourceType = resourceType[:idx]
	}
	// 末段为占位符时,取实际路径同位置段作为资源ID
	if last := tplRes[len(tplRes)-1]; strings.HasPrefix(last, "{") && strings.HasSuffix(last, "}") {
		if len(reqRes) == len(tplRes) {
			resourceId = reqRes[len(reqRes)-1]
		}
	}
	return resourceType, resourceId
}
