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

	swaggerUI "github.com/tx7do/kratos-swagger-ui"

	"github.com/tx7do/kratos-bootstrap/rpc"

	"go-wind-cms/app/admin/service/cmd/server/assets"
	"go-wind-cms/app/admin/service/internal/service"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	auditV1 "go-wind-cms/api/gen/go/audit/service/v1"

	"go-wind-cms/pkg/middleware/auth"
	applogging "go-wind-cms/pkg/middleware/logging"
)

// NewRestMiddleware 创建中间件
func NewRestMiddleware(
	ctx *bootstrap.Context,
	accessTokenChecker auth.AccessTokenChecker,
	authorizer authzEngine.Engine,
	apiAuditLogRepo auditV1.ApiAuditLogServiceClient,
	loginLogRepo auditV1.LoginAuditLogServiceClient,
) []middleware.Middleware {
	var ms []middleware.Middleware
	ms = append(ms, logging.Server(ctx.GetLogger()))

	// add white list for authentication.
	rpc.AddWhiteList(
		adminV1.OperationAuthenticationServiceLogin,
	)

	ms = append(ms, applogging.Server(
		applogging.WithWriteApiLogFunc(func(ctx context.Context, data *auditV1.ApiAuditLog) error {
			// TODO 如果系统的负载比较小，可以同步写入数据库，否则，建议使用异步方式，即投递进队列。
			_, err := apiAuditLogRepo.Create(ctx, &auditV1.CreateApiAuditLogRequest{Data: data})
			return err
		}),
		applogging.WithWriteLoginLogFunc(func(ctx context.Context, data *auditV1.LoginAuditLog) error {
			// TODO 如果系统的负载比较小，可以同步写入数据库，否则，建议使用异步方式，即投递进队列。
			_, err := loginLogRepo.Create(ctx, &auditV1.CreateLoginAuditLogRequest{Data: data})
			return err
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

	fileSvc *service.FileService,
	fileTransferService *service.FileTransferService,
	ueditorSvc *service.UEditorService,

	internalMessageService *service.InternalMessageService,
	internalMessageCategoryService *service.InternalMessageCategoryService,
	internalMessageRecipientService *service.InternalMessageRecipientService,

	apiAuditLogService *service.ApiAuditLogService,
	dataAccessAuditLogService *service.DataAccessAuditLogService,
	loginAuditLogService *service.LoginAuditLogService,
	policyEvaluationLogService *service.PolicyEvaluationLogService,
	operationAuditLogService *service.OperationAuditLogService,
	permissionAuditLogService *service.PermissionAuditLogService,

	postSvc *service.PostService,
	cateSvc *service.CategoryService,
	commentSvc *service.CommentService,
	tagSvc *service.TagService,
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

	adminV1.RegisterFileServiceHTTPServer(srv, fileSvc)
	adminV1.RegisterUEditorServiceHTTPServer(srv, ueditorSvc)

	// 注册文件传输服务，用于处理文件上传下载等功能
	// TODO 它不能够使用代码生成器生成的Handler，需要手动注册。代码生成器生成的Handler无法处理文件上传下载的请求。
	// 但，代码生成器生成代码可以提供给OpenAPI使用。
	registerFileTransferServiceHandler(srv, fileTransferService)

	adminV1.RegisterPostServiceHTTPServer(srv, postSvc)
	adminV1.RegisterCategoryServiceHTTPServer(srv, cateSvc)
	adminV1.RegisterTagServiceHTTPServer(srv, tagSvc)
	adminV1.RegisterCommentServiceHTTPServer(srv, commentSvc)

	if cfg.GetServer().GetRest().GetEnableSwagger() {
		swaggerUI.RegisterSwaggerUIServerWithOption(
			srv,
			swaggerUI.WithTitle("CMS Admin Service"),
			swaggerUI.WithMemoryData(assets.OpenApiData, "yaml"),
		)
	}

	return srv
}
