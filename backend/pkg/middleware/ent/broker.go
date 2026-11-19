package ent

import (
	"context"
	"strconv"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/go-crud/viewer"
	"github.com/tx7do/kratos-transport/broker"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/proto"

	"github.com/go-kratos/kratos/v2/encoding"
	_ "github.com/go-kratos/kratos/v2/encoding/json"

	appViewer "go-wind-cms/pkg/entgo/viewer"
)

var codec = encoding.GetCodec("proto")

const (
	MetadataHeaderKey        = "x-metadata-bin"
	TraceIDHeaderKey         = "x-trace-id"
	UserIDHeaderKey          = "x-user-id"
	TenantIDHeaderKey        = "x-tenant-id"
	OrgUnitIDHeaderKey       = "x-orgunit-id"
	DataScopeHeaderKey       = "x-data-scope"
	PermissionsHeaderKey     = "x-perms"
	RolesHeaderKey           = "x-roles"
	SystemContextHeaderKey   = "x-system"
	PlatformContextHeaderKey = "x-platform"
	TenantContextHeaderKey   = "x-tenant"
	ShouldAuditHeaderKey     = "x-audit"
)

func Publish() broker.PublishMiddleware {
	return func(handler broker.PublishHandler) broker.PublishHandler {
		return func(ctx context.Context, topic string, msg *broker.Message, opts ...broker.PublishOption) error {
			vCtx, ok := viewer.FromContext(ctx)
			if !ok {
				log.Warnf("ent broker middleware: no viewer found in context; topic=%s", topic)
				return handler(ctx, topic, msg, opts...)
			}

			if vCtx == nil {
				return handler(ctx, topic, msg, opts...)
			}

			var traceID string
			spanContext := trace.SpanContextFromContext(ctx)
			if spanContext.HasTraceID() {
				traceID = spanContext.TraceID().String()
				msg.SetHeader(TraceIDHeaderKey, traceID)
			}

			msg.SetHeader(UserIDHeaderKey, strconv.FormatInt(vCtx.UserID(), 10))
			msg.SetHeader(TenantIDHeaderKey, strconv.FormatInt(vCtx.TenantID(), 10))
			msg.SetHeader(OrgUnitIDHeaderKey, strconv.FormatInt(vCtx.OrgUnitID(), 10))
			if len(vCtx.Permissions()) > 0 {
				msg.SetHeader(PermissionsHeaderKey, strings.Join(vCtx.Permissions(), ","))
			}
			if len(vCtx.Roles()) > 0 {
				msg.SetHeader(RolesHeaderKey, strings.Join(vCtx.Roles(), ","))
			}
			if vCtx.IsSystemContext() {
				msg.SetHeader(SystemContextHeaderKey, strconv.FormatBool(vCtx.IsSystemContext()))
			}
			//msg.SetHeader("x-platform", strconv.FormatBool(vCtx.IsPlatformContext()))
			//msg.SetHeader("x-tenant", strconv.FormatBool(vCtx.IsTenantContext()))
			//msg.SetHeader("x-audit", strconv.FormatBool(vCtx.ShouldAudit()))

			if len(vCtx.DataScope()) > 0 {
				// 防护钉子：proto codec 的 Marshal 对非 proto.Message 直接类型断言
				// （panic），而 viewer.DataScope 是普通 struct——UserViewer 默认自带
				// 一个 datascope，走到这里必炸。先校验可编码性，不可编码 warn 跳过
				// （trace/uid 头照常携带，消费端仅丢 datascope）。
				if _, ok := any(vCtx.DataScope()).(proto.Message); !ok {
					log.Warnf("ent broker middleware: data scope is not proto-encodable (%T); header skipped", vCtx.DataScope())
				} else {
					buf, err := codec.Marshal(vCtx.DataScope())
					if err != nil {
						log.Warnf("ent broker middleware: failed to marshal data scope: %v", err)
					}
					msg.SetHeader(DataScopeHeaderKey, string(buf))
				}
			}

			return handler(ctx, topic, msg, opts...)
		}
	}
}

func Subscriber() broker.SubscriberMiddleware {
	return func(handler broker.Handler) broker.Handler {
		return func(ctx context.Context, evt broker.Event) error {
			if len(evt.Message().Headers) == 0 {
				return handler(ctx, evt)
			}

			traceID := evt.Message().GetHeader(TraceIDHeaderKey)
			if traceID != "" {
				// Publish 侧写的是 TraceID().String()（32 位 hex）。此前用
				// trace.TraceID([]byte(...)) 直接转 16 字节数组——长度不是 16
				// 就 panic（32 位 hex 必炸，消费端整条消息处理中断）。
				// 改走 TraceIDFromHex：hex 合法才回填，否则 warn 跳过。
				if scID, err := trace.TraceIDFromHex(traceID); err == nil {
					ctx = trace.ContextWithSpanContext(ctx, trace.SpanContextFromContext(ctx).WithTraceID(scID))
				} else {
					log.Warnf("ent broker middleware: invalid trace id header %q: %v", traceID, err)
				}
			}

			userIDStr := evt.Message().GetHeader(UserIDHeaderKey)
			tenantIDStr := evt.Message().GetHeader(TenantIDHeaderKey)
			orgUnitIDStr := evt.Message().GetHeader(OrgUnitIDHeaderKey)
			//permsStr := evt.Message().GetHeader(PermissionsHeaderKey)
			//rolesStr := evt.Message().GetHeader(RolesHeaderKey)
			systemContextStr := evt.Message().GetHeader(SystemContextHeaderKey)
			dataScopeStr := evt.Message().GetHeader(DataScopeHeaderKey)

			if systemContextStr != "" {
				isSystemContext, err := strconv.ParseBool(systemContextStr)
				if err != nil {
					log.Warnf("ent broker middleware: invalid system context header: %v", err)
				}

				if isSystemContext {
					ctx = appViewer.NewSystemViewerContext(ctx)
					return handler(ctx, evt)
				}
			}

			var userID, tenantID, orgUnitID int64
			var err error

			if userIDStr != "" {
				userID, err = strconv.ParseInt(userIDStr, 10, 64)
				if err != nil {
					log.Warnf("ent broker middleware: invalid user id header: %v", err)
				}
			}

			if tenantIDStr != "" {
				tenantID, err = strconv.ParseInt(tenantIDStr, 10, 64)
				if err != nil {
					log.Warnf("ent broker middleware: invalid tenant id header: %v", err)
				}
			}

			if orgUnitIDStr != "" {
				orgUnitID, err = strconv.ParseInt(orgUnitIDStr, 10, 64)
				if err != nil {
					log.Warnf("ent broker middleware: invalid org unit id header: %v", err)
				}
			}

			var dataScopes []viewer.DataScope
			if dataScopeStr != "" {
				err = codec.Unmarshal([]byte(dataScopeStr), &dataScopes)
				if err != nil {
					log.Warnf("ent broker middleware: failed to unmarshal data scopes: %v", err)
				}
			}

			userViewer := appViewer.NewUserViewerWithDataScopes(
				userID,
				tenantID,
				orgUnitID,
				traceID,
				dataScopes,
			)
			ctx = viewer.WithContext(ctx, userViewer)

			return handler(ctx, evt)
		}
	}
}
