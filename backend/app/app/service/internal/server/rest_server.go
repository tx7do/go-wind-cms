package server

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport/http"

	authzEngine "github.com/tx7do/kratos-authz/engine"
	authz "github.com/tx7do/kratos-authz/middleware"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"github.com/tx7do/kratos-bootstrap/rpc"

	swaggerUI "github.com/tx7do/kratos-swagger-ui"

	"go-wind-cms/app/app/service/cmd/server/assets"
	"go-wind-cms/app/app/service/internal/service"

	appV1 "go-wind-cms/api/gen/go/app/service/v1"
	auditV1 "go-wind-cms/api/gen/go/audit/service/v1"

	"go-wind-cms/pkg/middleware/auth"
	applogging "go-wind-cms/pkg/middleware/logging"
)

// NewRestMiddleware 创建中间件
func NewRestMiddleware(
	ctx *bootstrap.Context,
	accessTokenChecker auth.AccessTokenChecker,
	authorizer authzEngine.Engine,
) []middleware.Middleware {
	var ms []middleware.Middleware
	ms = append(ms, logging.Server(ctx.GetLogger()))

	// add white list for authentication.
	rpc.AddWhiteList(
		appV1.OperationAuthenticationServiceLogin,
	)

	ms = append(ms, applogging.Server(
		applogging.WithWriteApiLogFunc(func(ctx context.Context, data *auditV1.ApiAuditLog) error {
			return nil
		}),
		applogging.WithWriteLoginLogFunc(func(ctx context.Context, data *auditV1.LoginAuditLog) error {
			return nil
		}),
	))

	ms = append(ms, selector.Server(
		auth.Server(
			auth.WithAccessTokenChecker(accessTokenChecker),
			auth.WithInjectMetadata(true),
			auth.WithInjectEnt(true),
		),
		authz.Server(authorizer),
	).Match(rpc.NewRestWhiteListMatcher()).Build())

	return ms
}

// NewRestServer new an REST server.
func NewRestServer(
	ctx *bootstrap.Context,

	middlewares []middleware.Middleware,

	authenticationService *service.AuthenticationService,

	postService *service.PostService,
	categoryService *service.CategoryService,
	commentService *service.CommentService,
	tagService *service.TagService,
	pageService *service.PageService,
) *http.Server {
	cfg := ctx.GetConfig()

	if cfg == nil || cfg.Server == nil || cfg.Server.Rest == nil {
		return nil
	}

	srv, err := rpc.CreateRestServer(cfg, middlewares...)
	if err != nil {
		panic(err)
	}

	appV1.RegisterAuthenticationServiceHTTPServer(srv, authenticationService)

	appV1.RegisterPostServiceHTTPServer(srv, postService)
	appV1.RegisterCategoryServiceHTTPServer(srv, categoryService)
	appV1.RegisterTagServiceHTTPServer(srv, tagService)
	appV1.RegisterCommentServiceHTTPServer(srv, commentService)
	appV1.RegisterPageServiceHTTPServer(srv, pageService)

	if cfg.GetServer().GetRest().GetEnableSwagger() {
		swaggerUI.RegisterSwaggerUIServerWithOption(
			srv,
			swaggerUI.WithTitle("CMS App Service"),
			swaggerUI.WithMemoryData(assets.OpenApiData, "yaml"),
		)
	}

	return srv
}
