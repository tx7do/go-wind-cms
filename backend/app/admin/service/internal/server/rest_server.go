package server

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport/http"

	authzEngine "github.com/tx7do/kratos-authz/engine"
	authz "github.com/tx7do/kratos-authz/middleware"

	swaggerUI "github.com/tx7do/kratos-swagger-ui"

	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"github.com/tx7do/kratos-bootstrap/rpc"

	"go-wind-cms/app/admin/service/cmd/server/assets"
	"go-wind-cms/app/admin/service/internal/service"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	auditV1 "go-wind-cms/api/gen/go/audit/service/v1"
	identityV1 "go-wind-cms/api/gen/go/identity/service/v1"

	"go-wind-cms/pkg/metadata"
	"go-wind-cms/pkg/middleware/auth"
	applogging "go-wind-cms/pkg/middleware/logging"
	entmiddleware "go-wind-cms/pkg/middleware/ent"
)

// NewRestMiddleware 创建中间件
func NewRestMiddleware(
	ctx *bootstrap.Context,
	accessTokenChecker auth.AccessTokenChecker,
	authorizer authzEngine.Engine,
	apiAuditLogServiceClient auditV1.ApiAuditLogServiceClient,
	loginAuditLogServiceClient auditV1.LoginAuditLogServiceClient,
) []middleware.Middleware {
	var ms []middleware.Middleware
	ms = append(ms, logging.Server(ctx.GetLogger()))

	// add white list for authentication.
	// GenerateCaptcha 和 VerifyCaptcha 必须在白名单中，因为登录需要验证码，
	// 而获取验证码时用户尚无 token（鸡生蛋问题）。
	rpc.AddWhiteList(
		adminV1.OperationAuthenticationServiceLogin,
		adminV1.OperationAuthenticationServiceGenerateCaptcha,
		adminV1.OperationAuthenticationServiceVerifyCaptcha,
	)

	ms = append(ms, applogging.Server(
		applogging.WithWriteApiLogFunc(func(ctx context.Context, data *auditV1.ApiAuditLog) error {
			ctx, _ = metadata.NewContext(ctx, metadata.NewUserOperator(0, 0, 0, identityV1.DataScope_ALL))
			// TODO 如果系统的负载比较小，可以同步写入数据库，否则，建议使用异步方式，即投递进队列。
			_, err := apiAuditLogServiceClient.Create(ctx, &auditV1.CreateApiAuditLogRequest{Data: data})
			return err
		}),
		applogging.WithWriteLoginLogFunc(func(ctx context.Context, data *auditV1.LoginAuditLog) error {
			// TODO 如果系统的负载比较小，可以同步写入数据库，否则，建议使用异步方式，即投递进队列。
			ctx, _ = metadata.NewContext(ctx, metadata.NewUserOperator(0, 0, 0, identityV1.DataScope_ALL))
			_, err := loginAuditLogServiceClient.Create(ctx, &auditV1.CreateLoginAuditLogRequest{Data: data})
			return err
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
	// 可构建 UserViewer；白名单请求（登录/验证码）md==nil 但在白名单内，兜底 SystemViewer。
	ms = append(ms, entmiddleware.Server())

	return ms
}

// NewRestServer new an REST server.
func NewRestServer(
	ctx *bootstrap.Context,

	middlewares []middleware.Middleware,

	userService *service.UserService,
	userProfileService *service.UserProfileService,
	roleService *service.RoleService,
	tenantService *service.TenantService,
	orgUnitService *service.OrgUnitService,
	positionService *service.PositionService,

	menuSvc *service.MenuService,
	apiService *service.ApiService,
	permissionGroupService *service.PermissionGroupService,
	permissionService *service.PermissionService,

	adminPortalService *service.AdminPortalService,
	taskService *service.TaskService,

	authenticationService *service.AuthenticationService,
	loginPolicyService *service.LoginPolicyService,

	dictTypeService *service.DictTypeService,
	dictEntryService *service.DictEntryService,
	languageService *service.LanguageService,

	fileSvc *service.FileService,
	fileTransferService *service.FileTransferService,

	translatorService *service.TranslatorService,

	internalMessageService *service.InternalMessageService,
	internalMessageCategoryService *service.InternalMessageCategoryService,
	internalMessageRecipientService *service.InternalMessageRecipientService,

	apiAuditLogService *service.ApiAuditLogService,
	dataAccessAuditLogService *service.DataAccessAuditLogService,
	loginAuditLogService *service.LoginAuditLogService,
	policyEvaluationLogService *service.PolicyEvaluationLogService,
	operationAuditLogService *service.OperationAuditLogService,
	permissionAuditLogService *service.PermissionAuditLogService,

	commentService *service.CommentService,
	interactionAdminService *service.InteractionAdminService,

	postService *service.PostService,
	categoryService *service.CategoryService,
	tagService *service.TagService,
	contentModelService *service.ContentModelService,
	pageService *service.PageService,
	sectionService *service.SectionService,

	siteService *service.SiteService,
	siteSettingService *service.SiteSettingService,
	navigationService *service.NavigationService,
	navigationItemService *service.NavigationItemService,

	mediaAssetService *service.MediaAssetService,
) *http.Server {
	cfg := ctx.GetConfig()

	if cfg == nil || cfg.Server == nil || cfg.Server.Rest == nil {
		return nil
	}

	srv, err := rpc.CreateRestServer(cfg, middlewares...)
	if err != nil {
		panic(err)
	}

	adminV1.RegisterAuthenticationServiceHTTPServer(srv, authenticationService)
	adminV1.RegisterLoginPolicyServiceHTTPServer(srv, loginPolicyService)

	adminV1.RegisterUserProfileServiceHTTPServer(srv, userProfileService)
	adminV1.RegisterUserServiceHTTPServer(srv, userService)
	adminV1.RegisterRoleServiceHTTPServer(srv, roleService)
	adminV1.RegisterTenantServiceHTTPServer(srv, tenantService)
	adminV1.RegisterOrgUnitServiceHTTPServer(srv, orgUnitService)
	adminV1.RegisterPositionServiceHTTPServer(srv, positionService)

	adminV1.RegisterAdminPortalServiceHTTPServer(srv, adminPortalService)
	adminV1.RegisterTaskServiceHTTPServer(srv, taskService)

	adminV1.RegisterDictTypeServiceHTTPServer(srv, dictTypeService)
	adminV1.RegisterDictEntryServiceHTTPServer(srv, dictEntryService)
	adminV1.RegisterLanguageServiceHTTPServer(srv, languageService)

	adminV1.RegisterApiServiceHTTPServer(srv, apiService)
	adminV1.RegisterMenuServiceHTTPServer(srv, menuSvc)
	adminV1.RegisterPermissionGroupServiceHTTPServer(srv, permissionGroupService)
	adminV1.RegisterPermissionServiceHTTPServer(srv, permissionService)

	adminV1.RegisterApiAuditLogServiceHTTPServer(srv, apiAuditLogService)
	adminV1.RegisterDataAccessAuditLogServiceHTTPServer(srv, dataAccessAuditLogService)
	adminV1.RegisterLoginAuditLogServiceHTTPServer(srv, loginAuditLogService)
	adminV1.RegisterOperationAuditLogServiceHTTPServer(srv, operationAuditLogService)
	adminV1.RegisterPermissionAuditLogServiceHTTPServer(srv, permissionAuditLogService)
	adminV1.RegisterPolicyEvaluationLogServiceHTTPServer(srv, policyEvaluationLogService)

	adminV1.RegisterInternalMessageServiceHTTPServer(srv, internalMessageService)
	adminV1.RegisterInternalMessageCategoryServiceHTTPServer(srv, internalMessageCategoryService)
	adminV1.RegisterInternalMessageRecipientServiceHTTPServer(srv, internalMessageRecipientService)

	adminV1.RegisterTranslatorServiceHTTPServer(srv, translatorService)

	// 注册文件传输服务，用于处理文件上传下载等功能
	// TODO 它不能够使用代码生成器生成的Handler，需要手动注册。代码生成器生成的Handler无法处理文件上传下载的请求。
	// 但，代码生成器生成代码可以提供给OpenAPI使用。
	registerFileTransferServiceHandler(srv, fileTransferService)
	adminV1.RegisterFileServiceHTTPServer(srv, fileSvc)

	adminV1.RegisterPostServiceHTTPServer(srv, postService)

	// 内容模型管理 HTTP 路由
	adminV1.RegisterContentModelServiceHTTPServer(srv, contentModelService)
	adminV1.RegisterCategoryServiceHTTPServer(srv, categoryService)
	adminV1.RegisterTagServiceHTTPServer(srv, tagService)
	adminV1.RegisterCommentServiceHTTPServer(srv, commentService)
	adminV1.RegisterInteractionAdminServiceHTTPServer(srv, interactionAdminService)
	adminV1.RegisterPageServiceHTTPServer(srv, pageService)
	adminV1.RegisterSectionServiceHTTPServer(srv, sectionService)

	adminV1.RegisterSiteSettingServiceHTTPServer(srv, siteSettingService)
	adminV1.RegisterSiteServiceHTTPServer(srv, siteService)
	adminV1.RegisterNavigationServiceHTTPServer(srv, navigationService)
	adminV1.RegisterNavigationItemServiceHTTPServer(srv, navigationItemService)

	adminV1.RegisterMediaAssetServiceHTTPServer(srv, mediaAssetService)

	if cfg.GetServer().GetRest().GetEnableSwagger() {
		swaggerUI.RegisterSwaggerUIServerWithOption(
			srv,
			swaggerUI.WithTitle("GoWind Content Hub Admin API"),
			swaggerUI.WithMemoryData(assets.OpenApiData, "yaml"),
		)
	}

	return srv
}
