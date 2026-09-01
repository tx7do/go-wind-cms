package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	"go-wind-cms/app/admin/service/internal/data"
	"go-wind-cms/app/admin/service/internal/server"
	"go-wind-cms/app/admin/service/internal/service"
	"go-wind-cms/pkg/middleware/auth"
)

// initApp 手写装配整个应用,是本服务的依赖注入点(替代原 google/wire 生成)。
// initApp assembles the whole application by hand — the dependency injection point.
//
// 装配严格单向分层,自上而下的阅读顺序即依赖方向:
// The wiring is strictly layered; reading top-down follows the dependency direction:
//
//	基础设施 → 服务客户端(data) → 服务层(service) → 传输层(server)
//
// 本服务是 BFF 网关:不直连数据库,仓储层由对 core-service 的 gRPC 服务客户端替代。
// 约定 / Conventions:
//   - 新增模块:先加服务客户端行,再加服务构造行,并传给 NewRestServer;漏接由编译器在调用处报错。
//     To add a module: append a service-client line, then a service constructor, then pass it to
//     NewRestServer; a missing connection is a compile error at the call site.
//   - 持有 cleanup 的资源创建成功后立即注册进 cleanups;任何一步失败,rollback 逆序执行已注册
//     的清理(shutdown 与中途失败共用同一条 LIFO 路径)。
//     Resources owning a cleanup register it immediately; rollback runs them LIFO both on
//     mid-way failure and on shutdown.
//   - 本文件只做构造与传参,不写业务逻辑。
//     Construction and parameter passing only; no business logic in this file.
func initApp(ctx *bootstrap.Context) (*kratos.App, func(), error) {
	// cleanup 注册表:rollback 时逆序执行。
	// Cleanup registry; rollback runs entries in reverse order.
	var cleanups []func()
	rollback := func() {
		for i := len(cleanups) - 1; i >= 0; i-- {
			cleanups[i]()
		}
	}

	// ═══════════════════════ 一、基础设施 ═══════════════════════

	// 服务发现与鉴权基建
	discovery := data.NewDiscovery(ctx)
	clientType := data.NewClientType()
	engine := data.NewAuthorizer()

	redisClient, cleanupRedis, err := data.NewRedisClient(ctx)
	if err != nil {
		return nil, nil, err
	}
	cleanups = append(cleanups, cleanupRedis)

	captcha := data.NewCaptcha(redisClient)
	minIOClient := data.NewMinIoClient(ctx)
	translator := data.NewTranslator(ctx)

	// 认证:认证服务客户端 → 访问令牌校验器
	authenticationServiceClient := data.NewAuthenticationServiceClient(ctx, discovery)
	accessTokenChecker := auth.NewTokenChecker(ctx, authenticationServiceClient, clientType)

	// ═══════════════════════ 二、服务客户端(internal/data) ═══════════════════════

	// 身份与组织
	userServiceClient := data.NewUserServiceClient(ctx, discovery)
	roleServiceClient := data.NewRoleServiceClient(ctx, discovery)
	tenantServiceClient := data.NewTenantServiceClient(ctx, discovery)
	orgUnitServiceClient := data.NewOrgUnitServiceClient(ctx, discovery)
	positionServiceClient := data.NewPositionServiceClient(ctx, discovery)
	userCredentialServiceClient := data.NewUserCredentialServiceClient(ctx, discovery)
	loginPolicyServiceClient := data.NewLoginPolicyServiceClient(ctx, discovery)

	// RBAC
	menuServiceClient := data.NewMenuServiceClient(ctx, discovery)
	apiServiceClient := data.NewApiServiceClient(ctx, discovery)
	permissionServiceClient := data.NewPermissionServiceClient(ctx, discovery)
	permissionGroupServiceClient := data.NewPermissionGroupServiceClient(ctx, discovery)

	// 审计日志
	apiAuditLogServiceClient := data.NewApiAuditLogServiceClient(ctx, discovery)
	loginAuditLogServiceClient := data.NewLoginAuditLogServiceClient(ctx, discovery)
	operationAuditLogServiceClient := data.NewOperationAuditLogServiceClient(ctx, discovery)
	permissionAuditLogServiceClient := data.NewPermissionAuditLogServiceClient(ctx, discovery)
	dataAccessAuditLogServiceClient := data.NewDataAccessAuditLogServiceClient(ctx, discovery)
	policyEvaluationLogServiceClient := data.NewPolicyEvaluationLogServiceClient(ctx, discovery)

	// 字典与多语言
	dictTypeServiceClient := data.NewDictTypeServiceClient(ctx, discovery)
	dictEntryServiceClient := data.NewDictEntryServiceClient(ctx, discovery)
	languageServiceClient := data.NewLanguageServiceClient(ctx, discovery)

	// 文件 / 任务
	fileServiceClient := data.NewFileServiceClient(ctx, discovery)
	taskServiceClient := data.NewTaskServiceClient(ctx, discovery)
	mediaAssetServiceClient := data.NewMediaAssetServiceClient(ctx, discovery)

	// 站内信
	internalMessageServiceClient := data.NewInternalMessageServiceClient(ctx, discovery)
	internalMessageCategoryServiceClient := data.NewInternalMessageCategoryServiceClient(ctx, discovery)
	internalMessageRecipientServiceClient := data.NewInternalMessageRecipientServiceClient(ctx, discovery)

	// 内容
	postServiceClient := data.NewPostServiceClient(ctx, discovery)
	categoryServiceClient := data.NewCategoryServiceClient(ctx, discovery)
	tagServiceClient := data.NewTagServiceClient(ctx, discovery)
	contentModelServiceClient := data.NewContentModelServiceClient(ctx, discovery)
	pageServiceClient := data.NewPageServiceClient(ctx, discovery)
	commentServiceClient := data.NewCommentServiceClient(ctx, discovery)
	interactionAdminServiceClient := data.NewInteractionAdminServiceClient(ctx, discovery)
	statsServiceClient := data.NewStatsServiceClient(ctx, discovery)

	// 站点配置
	siteServiceClient := data.NewSiteServiceClient(ctx, discovery)
	siteSettingServiceClient := data.NewSiteSettingServiceClient(ctx, discovery)
	navigationServiceClient := data.NewNavigationServiceClient(ctx, discovery)
	navigationItemServiceClient := data.NewNavigationItemServiceClient(ctx, discovery)

	// ── register:client ── 新模块服务客户端在此行后注册(make register 工具锚点,勿删)

	// ═══════════════════════ 三、服务层(internal/service) ═══════════════════════

	// 认证与登录策略
	authenticationService := service.NewAuthenticationService(ctx, authenticationServiceClient, captcha)
	loginPolicyService := service.NewLoginPolicyService(ctx, loginPolicyServiceClient)

	// 身份与组织
	userService := service.NewUserService(ctx, userServiceClient, tenantServiceClient, orgUnitServiceClient, positionServiceClient, roleServiceClient, userCredentialServiceClient)
	userProfileService := service.NewUserProfileService(ctx, userServiceClient, tenantServiceClient, orgUnitServiceClient, positionServiceClient, roleServiceClient, userCredentialServiceClient)
	roleService := service.NewRoleService(ctx, roleServiceClient, tenantServiceClient)
	tenantService := service.NewTenantService(ctx, userServiceClient, userCredentialServiceClient, tenantServiceClient, roleServiceClient)
	orgUnitService := service.NewOrgUnitService(ctx, orgUnitServiceClient, userServiceClient)
	positionService := service.NewPositionService(ctx, positionServiceClient, orgUnitServiceClient)

	// RBAC 与管理门户
	menuService := service.NewMenuService(ctx, menuServiceClient)
	apiService := service.NewApiService(ctx, apiServiceClient)
	permissionGroupService := service.NewPermissionGroupService(ctx, permissionServiceClient, permissionGroupServiceClient)
	permissionService := service.NewPermissionService(ctx, permissionServiceClient, permissionGroupServiceClient, roleServiceClient, apiServiceClient, menuServiceClient)
	adminPortalService := service.NewRouterService(ctx, menuServiceClient, permissionServiceClient, roleServiceClient, userServiceClient)

	// 任务 / 文件 / 翻译
	taskService := service.NewTaskService(ctx, taskServiceClient)
	fileService := service.NewFileService(ctx, fileServiceClient)
	fileTransferService := service.NewFileTransferService(ctx, minIOClient, fileServiceClient, mediaAssetServiceClient)
	translatorService := service.NewTranslatorService(ctx, translator)

	// 字典与多语言
	dictTypeService := service.NewDictTypeService(ctx, dictTypeServiceClient)
	dictEntryService := service.NewDictEntryService(ctx, dictEntryServiceClient)
	languageService := service.NewLanguageService(ctx, languageServiceClient)

	// 审计日志
	apiAuditLogService := service.NewApiAuditLogService(ctx, apiAuditLogServiceClient, apiServiceClient)
	dataAccessAuditLogService := service.NewDataAccessAuditLogService(ctx, dataAccessAuditLogServiceClient)
	loginAuditLogService := service.NewLoginAuditLogService(ctx, loginAuditLogServiceClient)
	policyEvaluationLogService := service.NewPolicyEvaluationLogService(ctx, policyEvaluationLogServiceClient)
	operationAuditLogService := service.NewOperationAuditLogService(ctx, operationAuditLogServiceClient)
	permissionAuditLogService := service.NewPermissionAuditLogService(ctx, permissionAuditLogServiceClient)

	// 站内信
	internalMessageService := service.NewInternalMessageService(ctx, internalMessageServiceClient, internalMessageCategoryServiceClient, internalMessageRecipientServiceClient, authenticationServiceClient, userServiceClient, clientType)
	internalMessageCategoryService := service.NewInternalMessageCategoryService(ctx, internalMessageCategoryServiceClient)
	internalMessageRecipientService := service.NewInternalMessageRecipientService(ctx, internalMessageServiceClient, internalMessageRecipientServiceClient)

	// 内容
	commentService := service.NewCommentService(ctx, commentServiceClient)
	interactionAdminService := service.NewInteractionAdminService(ctx, interactionAdminServiceClient)
	statsService := service.NewStatsService(ctx, statsServiceClient)
	postService := service.NewPostService(ctx, postServiceClient)
	categoryService := service.NewCategoryService(ctx, categoryServiceClient)
	tagService := service.NewTagService(ctx, tagServiceClient)
	contentModelService := service.NewContentModelService(ctx, contentModelServiceClient)
	pageService := service.NewPageService(ctx, pageServiceClient)

	// 站点配置与媒体
	siteService := service.NewSiteService(ctx, siteServiceClient)
	siteSettingService := service.NewSiteSettingService(ctx, siteSettingServiceClient)
	navigationService := service.NewNavigationService(ctx, navigationServiceClient)
	navigationItemService := service.NewNavigationItemService(ctx, navigationItemServiceClient)
	mediaAssetService := service.NewMediaAssetService(ctx, mediaAssetServiceClient)

	// ── register:service ── 新模块服务在此行后注册(make register 工具锚点,勿删)

	// ═══════════════════════ 四、传输层(internal/server) ═══════════════════════

	restMiddlewares := server.NewRestMiddleware(ctx, accessTokenChecker, engine, apiAuditLogServiceClient, loginAuditLogServiceClient)

	httpServer := server.NewRestServer(ctx, restMiddlewares,
		userService, userProfileService, roleService, tenantService, orgUnitService, positionService,
		menuService, apiService, permissionGroupService, permissionService,
		adminPortalService, taskService,
		authenticationService, loginPolicyService,
		dictTypeService, dictEntryService, languageService,
		fileService, fileTransferService,
		translatorService,
		internalMessageService, internalMessageCategoryService, internalMessageRecipientService,
		apiAuditLogService, dataAccessAuditLogService, loginAuditLogService, policyEvaluationLogService, operationAuditLogService, permissionAuditLogService,
		commentService, interactionAdminService, statsService,
		postService, categoryService, tagService, contentModelService, pageService,
		siteService, siteSettingService, navigationService, navigationItemService,
		mediaAssetService,
		// register:rest-arg ── 新模块服务实参在此行后追加(make register 工具锚点,勿删)
	)

	grpcMiddlewares := server.NewGrpcMiddleware(ctx)
	grpcServer, err := server.NewGrpcServer(ctx, grpcMiddlewares)
	if err != nil {
		rollback()
		return nil, nil, err
	}

	sseServer := server.NewSseServer(ctx, internalMessageService)

	return newApp(ctx, httpServer, grpcServer, sseServer), rollback, nil
}
