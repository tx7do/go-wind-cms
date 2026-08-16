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
	tenantResolver entmiddleware.TenantResolver,
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

		// PostService.SearchPosts：公开全文搜索，与文章列表/详情的匿名可见性一致。
		// tenant_id 由 core 端从 viewer（匿名经路线2 注入的 AnonymousTenantViewer，
		// 登录为 UserViewer）提取，按 tenant 隔离，仅返回 PUBLISHED。调用方无法
		// 指定或绕过 tenant。
		appV1.OperationPostServiceSearchPosts,

		// InteractionService.GetCounts：公开计数（如点赞数）随文章列表展示，
		// 仅按 tenant 隔离、不依赖 viewer 身份。Like/Unlike/Watch 等写操作
		// 及 GetInteractionStatus（含 viewer 个人状态）仍需登录，故不在此登记。
		appV1.OperationInteractionServiceGetCounts,
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
	// 可构建 UserViewer；白名单请求（公开内容）md==nil，由注入的 TenantResolver 按
	// Host 解析 tenant_id 并注入只读 AnonymousTenantViewer（按 tenant 隔离）；解析失败
	// fail-closed 注入 noopContext（拒绝），不再回退 SystemViewer 避免跨租户泄漏。
	ms = append(ms, entmiddleware.Server(entmiddleware.WithTenantResolver(tenantResolver)))

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
