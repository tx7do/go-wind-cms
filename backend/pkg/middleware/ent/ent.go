package ent

import (
	"context"
	"reflect"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"

	"github.com/tx7do/go-crud/viewer"
	"go.opentelemetry.io/otel/trace"

	authenticationV1 "go-wind-cms/api/gen/go/authentication/service/v1"

	appViewer "go-wind-cms/pkg/entgo/viewer"
	"go-wind-cms/pkg/metadata"
)

// Server 设置 Ent Viewer 到 Context 中的中间件
func Server() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {
			md, err := metadata.FromServerContext(ctx)
			if err != nil {
				reqType := "<nil>"
				if req != nil {
					reqType = reflect.TypeOf(req).String()
				}
				t, _ := transport.FromServerContext(ctx)
				kind, operation, endpoint := "", "", ""
				if t != nil {
					kind = string(t.Kind())
					operation = t.Operation()
					endpoint = t.Endpoint()
				}
				log.Errorf("ent middleware: failed to get metadata from context: %v; req_type=%s; transport=%s; operation=%s; endpoint=%s",
					err, reqType, kind, operation, endpoint)
			}
			if md == nil {
				return handler(ctx, req)
			}

			var traceID string
			spanContext := trace.SpanContextFromContext(ctx)
			if spanContext.HasTraceID() {
				traceID = spanContext.TraceID().String()
			}

			ctx = viewer.WithContext(ctx, metaDataToUserViewerContext(md, traceID))

			return handler(ctx, req)
		}
	}
}

func metaDataToUserViewerContext(md *authenticationV1.OperatorMetadata, traceID string) viewer.Context {
	if md == nil {
		return nil
	}

	userViewer := appViewer.NewUserViewer(
		md.GetUserId(),
		md.GetTenantId(),
		md.GetOrgUnitId(),
		traceID,
		md.GetDataScope(),
	)
	return userViewer
}
