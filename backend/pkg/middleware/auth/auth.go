package auth

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"

	"github.com/tx7do/go-crud/viewer"
	"go.opentelemetry.io/otel/trace"

	authenticationV1 "go-wind-cms/api/gen/go/authentication/service/v1"

	authnEngine "github.com/tx7do/kratos-authn/engine"
	authzEngine "github.com/tx7do/kratos-authz/engine"

	appViewer "go-wind-cms/pkg/entgo/viewer"
	"go-wind-cms/pkg/metadata"
)

var defaultAction = authzEngine.Action("ANY")

// Server 衔接认证和鉴权
func Server(opts ...Option) middleware.Middleware {
	op := options{
		log: log.NewHelper(log.With(log.DefaultLogger, "module", "auth/middleware")),

		injectOperatorId: false,
		injectTenantId:   false,
		enableAuthz:      true,
		injectEnt:        true,
		injectMetadata:   true,
	}
	for _, o := range opts {
		o(&op)
	}

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {
			//op.log.Debugf("auth middleware: processing request")

			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				op.log.Errorf("auth middleware: missing transport in context")
				return nil, ErrWrongContext
			}

			token, err := authnEngine.AuthFromMD(ctx, authnEngine.BearerWord, authnEngine.ContextTypeKratosMetaData)
			if err != nil {
				return nil, ErrMissingBearerToken
			}

			if op.accessTokenChecker == nil {
				op.log.Errorf("auth middleware: access token checker is not configured")
				return nil, ErrAccessTokenCheckerNotConfigured
			}

			var tokenPayload *authenticationV1.UserTokenPayload
			var valid bool
			if valid, tokenPayload = op.accessTokenChecker.IsValidAccessToken(ctx, token, false); !valid {
				op.log.Errorf("auth middleware: invalid access token")
				return nil, ErrAccessTokenExpired
			}

			ctx = NewContext(ctx, tokenPayload)

			if op.injectOperatorId {
				if err = setRequestOperationId(req, tokenPayload); err != nil {
					op.log.Errorf("auth middleware: invalid token payload in context [%s]", err.Error())
					return nil, err
				}
			}
			if op.injectTenantId {
				if err = setRequestTenantId(req, tokenPayload); err != nil {
					op.log.Errorf("auth middleware: invalid token payload in context [%s]", err.Error())
					return nil, err
				}
			}

			if op.injectEnt {
				var traceID string
				spanContext := trace.SpanContextFromContext(ctx)
				if spanContext.HasTraceID() {
					traceID = spanContext.TraceID().String()
				}

				userViewer := appViewer.NewUserViewer(
					uint64(tokenPayload.GetUserId()),
					uint64(tokenPayload.GetTenantId()),
					uint64(tokenPayload.GetOrgUnitId()),
					traceID,
					tokenPayload.GetDataScope(),
				)
				ctx = viewer.WithContext(ctx, userViewer)
			}

			if op.injectMetadata {
				//op.log.Debugf("auth middleware: injecting operator metadata into context")
				operator := metadata.NewUserOperator(
					uint64(tokenPayload.GetUserId()),
					uint64(tokenPayload.GetTenantId()),
					uint64(tokenPayload.GetOrgUnitId()),
					tokenPayload.GetDataScope(),
				)
				ctx, err = metadata.NewContext(ctx, operator)
				if err != nil {
					op.log.Errorf("auth middleware: failed to inject operator metadata into context [%s]", err.Error())
					return nil, err
				}
				//op.log.Debugf("auth middleware: successfully injected operator metadata into context")
				if htr, ok := tr.(*http.Transport); ok {
					if err = metadata.SetOperatorToRequestHeader(htr.RequestHeader(), operator); err != nil {
						op.log.Errorf("auth middleware: failed to set operator metadata to request header [%s]", err.Error())
					}
				}
			}

			if op.enableAuthz {
				ctx, err = processAuthz(ctx, tr, tokenPayload)
				if err != nil {
					op.log.Errorf("auth middleware: authorization failed [%s]", err.Error())
					return nil, err
				}
			}

			return handler(ctx, req)
		}
	}
}
