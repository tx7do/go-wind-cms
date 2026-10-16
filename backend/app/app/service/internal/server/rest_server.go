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
	authenticationV1 "go-wind-cms/api/gen/go/authentication/service/v1"

	"go-wind-cms/pkg/middleware/auth"
	applogging "go-wind-cms/pkg/middleware/logging"
	entmiddleware "go-wind-cms/pkg/middleware/ent"
)

// OperationAuthenticationServiceRegister 手工挂载的 C 端注册路由操作名。
// core 的 RegisterUser RPC 已实现,但 app proto 尚未声明该接口,为规避 proto
// 再生成,路由按生成代码的等价方式手工挂载;须与 AddWhiteList 保持一致。
const OperationAuthenticationServiceRegister = "/app.v1.AuthenticationService/Register"

// SEO 域手工挂载端点的操作名（sitemap.xml / robots.txt，非 proto 路由）。
const (
	OperationSeoServiceSitemap = "SeoService_Sitemap"
	OperationSeoServiceRobots  = "SeoService_Robots"
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

	// AuthenticationService.GetOidcAuthorizeUrl：OIDC 登录的前置步骤（登录页
	// 尚无 token），与 Login 同列免鉴权；state 防 CSRF 由 core 生成校验。
	appV1.OperationAuthenticationServiceGetOidcAuthorizeUrl,

		appV1.OperationNavigationServiceList,

		// AuthenticationService.RefreshToken：access token 过期后凭 refresh_token 换新，
		// 调用时 bearer 往往已过期，若要求有效 token 刷新将永远 401。
		appV1.OperationAuthenticationServiceRefreshToken,

		// C 端注册（手工挂载路由，见 NewRestServer；app proto 暂未声明该接口）。
		OperationAuthenticationServiceRegister,

		// CommentService.Create：游客评论策略在此操作内自行执行
		// （可选认证 + enable_comments/allow_guest_comments 开关），
		// 认证白名单仅为放行游客请求。
		appV1.OperationCommentServiceCreate,

		// SiteService.GetSiteByDomain：公开站点配置（template/theme/default_locale 等
		// 渲染必需字段）。domain 由 BFF 按请求 Host 填入，调用方不可指定；core 端返回前
		// 已裁剪租户/状态/备用域名等运维字段。其余 SiteService 操作需登录，不在此登记。
		appV1.OperationSiteServiceGetSiteByDomain,

		appV1.OperationPageServiceList,
		appV1.OperationPostServiceList,
		appV1.OperationCategoryServiceList,
		appV1.OperationCommentServiceList,
		appV1.OperationTagServiceList,

		appV1.OperationPageServiceGet,
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

		// InteractionService.RecordView：浏览量上报，匿名可调用（PV 语义，
		// 无 viewer 个人状态；租户归属由服务层按 Host 解析内网注入）。
		appV1.OperationInteractionServiceRecordView,

		// FormSubmissionService.Submit：页面 FORM 区块的访客提交，匿名可调用
		//（校验在 core：页面须已发布、区块须为该页面的 FORM 区块、字段集受限）。
		appV1.OperationFormSubmissionServiceSubmit,

		// RedirectService.ListEnabled：公开重定向查询（C 端路由中间件 3xx 匹配）。
		appV1.OperationRedirectServiceListEnabled,

		// PostService/PageService.GetBySlug：公开 SEO 路由（/post/{slug}、/page/{slug}
		// 的数据源），仅返回已发布内容，匿名可访问。
		appV1.OperationPostServiceGetBySlug,
		appV1.OperationPageServiceGetBySlug,

		// SeoService（手工挂载非 proto 路由）：sitemap.xml / robots.txt，
		// 搜索引擎抓取必须匿名可达；操作名为约定字符串非 proto 生成。
		OperationSeoServiceSitemap,
		OperationSeoServiceRobots,
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
	seoService *service.SeoService,
	formSubmissionService *service.FormSubmissionService,
	redirectService *service.RedirectService,
	navigationService *service.NavigationService,
	siteService *service.SiteService,

	// register:param ── 新模块服务形参在此行后注册(make register 工具锚点,勿删)
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

	// ── 手工挂载:C 端注册 POST /app/v1/register ──
	// 与生成代码等价的注册方式;操作名须与 AddWhiteList 一致。
	{
		registerRoute := srv.Route("/")
		registerRoute.POST("/app/v1/register", func(ctx http.Context) error {
			http.SetOperation(ctx, OperationAuthenticationServiceRegister)
			var in authenticationV1.RegisterUserRequest
			if err := ctx.Bind(&in); err != nil {
				return err
			}
			out, err := authenticationService.Register(ctx, &in)
			if err != nil {
				return err
			}
			return ctx.Result(200, out)
		})
	}

	// ── 手工挂载:SEO 域 sitemap.xml / robots.txt ──
	// 非 proto 定义的部署域端点：操作名为约定字符串（与 AddWhiteList 一致），
	// 直接写响应体（XML/TXT），不走 kratos JSON 编码。
	{
		seoRoute := srv.Route("/")
		seoRoute.GET("/app/v1/sitemap.xml", func(ctx http.Context) error {
			http.SetOperation(ctx, OperationSeoServiceSitemap)
			body, err := seoService.Sitemap(ctx)
			if err != nil {
				return err
			}
			ctx.Response().Header().Set("Content-Type", "application/xml; charset=utf-8")
			return ctx.Blob(200, "application/xml; charset=utf-8", []byte(body))
		})
		seoRoute.GET("/app/v1/robots.txt", func(ctx http.Context) error {
			http.SetOperation(ctx, OperationSeoServiceRobots)
			body, err := seoService.Robots(ctx)
			if err != nil {
				return err
			}
			return ctx.Blob(200, "text/plain; charset=utf-8", []byte(body))
		})
	}

	appV1.RegisterFileTransferServiceHTTPServer(srv, fileTransferService)
	appV1.RegisterUserProfileServiceHTTPServer(srv, userProfileService)

	appV1.RegisterNavigationServiceHTTPServer(srv, navigationService)

	appV1.RegisterSiteServiceHTTPServer(srv, siteService)

	appV1.RegisterPostServiceHTTPServer(srv, postService)
	appV1.RegisterCategoryServiceHTTPServer(srv, categoryService)
	appV1.RegisterTagServiceHTTPServer(srv, tagService)
	appV1.RegisterPageServiceHTTPServer(srv, pageService)
	appV1.RegisterFormSubmissionServiceHTTPServer(srv, formSubmissionService)
	appV1.RegisterRedirectServiceHTTPServer(srv, redirectService)

	appV1.RegisterCommentServiceHTTPServer(srv, commentService)

	appV1.RegisterInteractionServiceHTTPServer(srv, interactionService)

	// register:route ── 新模块路由在此行后注册(make register 工具锚点,勿删)

	if cfg.GetServer().GetRest().GetEnableSwagger() {
		swaggerUI.RegisterSwaggerUIServerWithOption(
			srv,
			swaggerUI.WithTitle("GoWind Content Hub App API"),
			swaggerUI.WithMemoryData(assets.OpenApiData, "yaml"),
		)
	}

	return srv
}
