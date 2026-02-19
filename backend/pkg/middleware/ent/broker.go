package ent

import (
	"context"
	"strconv"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/go-crud/viewer"
	"github.com/tx7do/kratos-transport/broker"
	"go.opentelemetry.io/otel/trace"

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

			msg.SetHeader(UserIDHeaderKey, strconv.FormatUint(vCtx.UserID(), 10))
			msg.SetHeader(TenantIDHeaderKey, strconv.FormatUint(vCtx.TenantID(), 10))
			msg.SetHeader(OrgUnitIDHeaderKey, strconv.FormatUint(vCtx.OrgUnitID(), 10))
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
				buf, err := codec.Marshal(vCtx.DataScope())
				if err != nil {
					log.Warnf("ent broker middleware: failed to marshal data scope: %v", err)
				}
				msg.SetHeader(DataScopeHeaderKey, string(buf))
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
				ctx = trace.ContextWithSpanContext(ctx, trace.SpanContextFromContext(ctx).WithTraceID(
					trace.TraceID([]byte(traceID))),
				)
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

			var userID, tenantID, orgUnitID uint64
			var err error

			if userIDStr != "" {
				userID, err = strconv.ParseUint(userIDStr, 10, 64)
				if err != nil {
					log.Warnf("ent broker middleware: invalid user id header: %v", err)
				}
			}

			if tenantIDStr != "" {
				tenantID, err = strconv.ParseUint(tenantIDStr, 10, 64)
				if err != nil {
					log.Warnf("ent broker middleware: invalid tenant id header: %v", err)
				}
			}

			if orgUnitIDStr != "" {
				orgUnitID, err = strconv.ParseUint(orgUnitIDStr, 10, 64)
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
