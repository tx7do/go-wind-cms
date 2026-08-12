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
	entmiddleware "go-wind-cms/pkg/middleware/ent"
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

		appV1.OperationNavigationServiceList,
		appV1.OperationPageServiceList,
		appV1.OperationPostServiceList,
		appV1.OperationCategoryServiceList,
		appV1.OperationCommentServiceList,
		appV1.OperationTagServiceList,

		appV1.OperationPageServiceGet,
		appV1.OperationSectionServiceList,
		appV1.OperationSectionServiceGet,
		appV1.OperationSectionServiceGetTranslation,
		appV1.OperationPostServiceGet,
		appV1.OperationCategoryServiceGet,
		appV1.OperationCommentServiceGet,
		appV1.OperationTagServiceGet,
	)

	ms = append(ms, applogging.Server(
		applogging.WithWriteApiLogFunc(func(ctx context.Context, data *auditV1.ApiAuditLog) error {
			return nil
		}),
		applogging.WithWriteLoginLogFunc(func(ctx context.Context, data *auditV1.LoginAuditLog) error {
			return nil
		}),
	))

	// 鉴权必须在 ent.Server() 之前执行：auth.Server 对非白名单请求注入
	// OperatorMetadata，随后 ent.Server() 才能据此构建带租户作用域的 UserViewer。
	// 若顺序颠倒，ent.Server() 总以 md==nil 兜底为 SystemViewer，导致租户隔离失效。
	ms = append(ms, selector.Server(
		auth.Server(
			auth.WithAccessTokenChecker(accessTokenChecker),
			auth.WithInjectMetadata(true),
			auth.WithInjectEnt(true),
		),
		authz.Server(authorizer),
	).Match(rpc.NewRestWhiteListMatcher()).Build())

	// ent.Server() 必须在 auth.Server 之后：此时非白名单请求已注入 OperatorMetadata，
	// 可构建 UserViewer；白名单请求（登录/公开内容）md==nil 但在白名单内，兜底 SystemViewer。
	ms = append(ms, entmiddleware.Server())

	return ms
}

// NewRestServer new an REST server.
func NewRestServer(
	ctx *bootstrap.Context,

	middlewares []middleware.Middleware,

	authenticationService *service.AuthenticationService,
	fileTransferService *service.FileTransferService,
	userProfileService *service.UserProfileService,

	postService *service.PostService,
	categoryService *service.CategoryService,
	commentService *service.CommentService,
	interactionService *service.InteractionService,
	tagService *service.TagService,
	pageService *service.PageService,
	sectionService *service.SectionService,
	navigationService *service.NavigationService,
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
	appV1.RegisterFileTransferServiceHTTPServer(srv, fileTransferService)
	appV1.RegisterUserProfileServiceHTTPServer(srv, userProfileService)

	appV1.RegisterNavigationServiceHTTPServer(srv, navigationService)

	appV1.RegisterPostServiceHTTPServer(srv, postService)
	appV1.RegisterCategoryServiceHTTPServer(srv, categoryService)
	appV1.RegisterTagServiceHTTPServer(srv, tagService)
	appV1.RegisterPageServiceHTTPServer(srv, pageService)
	appV1.RegisterSectionServiceHTTPServer(srv, sectionService)

	appV1.RegisterCommentServiceHTTPServer(srv, commentService)

	appV1.RegisterInteractionServiceHTTPServer(srv, interactionService)

	if cfg.GetServer().GetRest().GetEnableSwagger() {
		swaggerUI.RegisterSwaggerUIServerWithOption(
			srv,
			swaggerUI.WithTitle("GoWind Content Hub App API"),
			swaggerUI.WithMemoryData(assets.OpenApiData, "yaml"),
		)
	}

	return srv
}
