package main

import (
	"os"

	"github.com/go-kratos/kratos/v2"
	kratosTransport "github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	pkgcrypto "go-wind-cms/pkg/crypto"

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

	// 全局存储加密器(AES-256-GCM, pkg/crypto):与 core 服务必须配置同一个
	// ENCRYPTION_KEY,否则一边加密一边解不开。未配置时降级明文并告警(生产必须配置)。
	if encKey := os.Getenv("ENCRYPTION_KEY"); encKey != "" {
		if err := pkgcrypto.InitGlobalEncryptor(encKey, true); err != nil {
			return nil, nil, err
		}
	} else {
		log.Warnf("ENCRYPTION_KEY not set: credential encryption DISABLED (secrets stored in plaintext)")
	}

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
	aiProviderServiceClient := data.NewAiProviderServiceClient(ctx, discovery)
	aiContentServiceClient := data.NewAiContentServiceClient(ctx, discovery)
	aiChatServiceClient := data.NewAiChatServiceClient(ctx, discovery)
	aiConversationServiceClient := data.NewAiConversationServiceClient(ctx, discovery)
	aiMessageServiceClient := data.NewAiMessageServiceClient(ctx, discovery)
	aiKnowledgeBaseServiceClient := data.NewAiKnowledgeBaseServiceClient(ctx, discovery)
	aiUsageLogServiceClient := data.NewAiUsageLogServiceClient(ctx, discovery)
	aiQueryServiceClient := data.NewAiQueryServiceClient(ctx, discovery)
	auditExportService := service.NewAuditExportService(ctx, loginAuditLogServiceClient, apiAuditLogServiceClient, operationAuditLogServiceClient, dataAccessAuditLogServiceClient, permissionAuditLogServiceClient, policyEvaluationLogServiceClient, aiUsageLogServiceClient, accessTokenChecker)

	// 站点配置
	siteServiceClient := data.NewSiteServiceClient(ctx, discovery)
	siteSettingServiceClient := data.NewSiteSettingServiceClient(ctx, discovery)
	navigationServiceClient := data.NewNavigationServiceClient(ctx, discovery)
	navigationItemServiceClient := data.NewNavigationItemServiceClient(ctx, discovery)

	// 社交分发
	socialServiceClient := data.NewSocialServiceClient(ctx, discovery)
	webhookServiceClient := data.NewWebhookServiceClient(ctx, discovery)
	flowRuleServiceClient := data.NewFlowRuleServiceClient(ctx, discovery)
	formSubmissionServiceClient := data.NewFormSubmissionServiceClient(ctx, discovery)
	redirectServiceClient := data.NewRedirectServiceClient(ctx, discovery)
	experimentServiceClient := data.NewExperimentServiceClient(ctx, discovery)
	editLockServiceClient := data.NewEditLockServiceClient(ctx, discovery)
	collabDocStateServiceClient := data.NewCollabDocStateServiceClient(ctx, discovery)

	// 协同服务（COLLAB_ADDR，默认 :6602；off 时仅 service 供删除踢房 no-op）
	collabSvc, collabServer := server.NewCollab(ctx, authenticationServiceClient, collabDocStateServiceClient, postServiceClient)

	// 通知域（渠道花名册 / 投递台账 / 路由规则 / 模板 / 用户偏好）
	notificationServiceClient := data.NewNotificationServiceClient(ctx, discovery)
	monitorAlertServiceClient := data.NewMonitorAlertServiceClient(ctx, discovery)
	onlineSessionServiceClient := data.NewOnlineSessionServiceClient(ctx, discovery)
	accessKeyServiceClient := data.NewAccessKeyServiceClient(ctx, discovery)
	tenantQuotaServiceClient := data.NewTenantQuotaServiceClient(ctx, discovery)
	configServiceClient := data.NewConfigServiceClient(ctx, discovery)
	aiAgentServiceClient := data.NewAiAgentServiceClient(ctx, discovery)
	serverMonitorServiceClient := data.NewServerMonitorServiceClient(ctx, discovery)
	redisCacheMonitorServiceClient := data.NewRedisCacheMonitorServiceClient(ctx, discovery)
	notificationChannelServiceClient := data.NewNotificationChannelServiceClient(ctx, discovery)
	notificationRuleServiceClient := data.NewNotificationRuleServiceClient(ctx, discovery)
	notificationTemplateServiceClient := data.NewNotificationTemplateServiceClient(ctx, discovery)
	notificationPreferenceServiceClient := data.NewNotificationPreferenceServiceClient(ctx, discovery)

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
	previewService := service.NewPreviewService(ctx)

	// 任务 / 文件 / 翻译
	taskService := service.NewTaskService(ctx, taskServiceClient)
	fileService := service.NewFileService(ctx, fileServiceClient)
	fileTransferService := service.NewFileTransferService(ctx, minIOClient, fileServiceClient, mediaAssetServiceClient, tenantQuotaServiceClient)
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

	// 订阅 core 服务经 Redis Pub/Sub 发出的系统通知事件，转 SSE 实时推送
	//（落库在 core 侧完成；订阅协程随进程退出停止，注册进 cleanups）
	internalMessageService.StartNotifySubscriber(redisClient)
	cleanups = append(cleanups, internalMessageService.StopNotifySubscriber)
	internalMessageCategoryService := service.NewInternalMessageCategoryService(ctx, internalMessageCategoryServiceClient)
	internalMessageRecipientService := service.NewInternalMessageRecipientService(ctx, internalMessageServiceClient, internalMessageRecipientServiceClient)

	// 内容
	commentService := service.NewCommentService(ctx, commentServiceClient)
	interactionAdminService := service.NewInteractionAdminService(ctx, interactionAdminServiceClient)
	statsService := service.NewStatsService(ctx, statsServiceClient)
	aiProviderService := service.NewAiProviderService(ctx, aiProviderServiceClient)
	aiContentService := service.NewAiContentService(ctx, aiContentServiceClient)
	aiChatService := service.NewAiChatService(ctx, aiChatServiceClient, authenticationServiceClient)
	aiConversationService := service.NewAiConversationService(ctx, aiConversationServiceClient)
	aiMessageService := service.NewAiMessageService(ctx, aiMessageServiceClient)
	aiKnowledgeService := service.NewAiKnowledgeService(ctx, aiKnowledgeBaseServiceClient)
	aiUsageLogService := service.NewAiUsageLogService(ctx, aiUsageLogServiceClient)
	aiQueryService := service.NewAiQueryService(ctx, aiQueryServiceClient)
	postService := service.NewPostService(ctx, postServiceClient, collabDocStateServiceClient, collabSvc)
	categoryService := service.NewCategoryService(ctx, categoryServiceClient)
	tagService := service.NewTagService(ctx, tagServiceClient)
	contentModelService := service.NewContentModelService(ctx, contentModelServiceClient)
	componentDefinitionService := service.NewComponentDefinitionService(ctx, contentModelServiceClient)
	pageService := service.NewPageService(ctx, pageServiceClient, collabDocStateServiceClient, collabSvc)

	// 站点配置与媒体
	siteService := service.NewSiteService(ctx, siteServiceClient)
	siteSettingService := service.NewSiteSettingService(ctx, siteSettingServiceClient)
	navigationService := service.NewNavigationService(ctx, navigationServiceClient)
	navigationItemService := service.NewNavigationItemService(ctx, navigationItemServiceClient)
	mediaAssetService := service.NewMediaAssetService(ctx, mediaAssetServiceClient)

	// 社交分发
	socialService := service.NewSocialService(ctx, socialServiceClient)
	webhookService := service.NewWebhookService(ctx, webhookServiceClient)
	formSubmissionService := service.NewFormSubmissionService(ctx, formSubmissionServiceClient)
	redirectService := service.NewRedirectService(ctx, redirectServiceClient)
	experimentService := service.NewExperimentService(ctx, experimentServiceClient)
	editLockService := service.NewEditLockService(ctx, editLockServiceClient)

	flowRuleService := service.NewFlowRuleService(ctx, flowRuleServiceClient)

	notificationService := service.NewNotificationService(ctx, notificationServiceClient)
	notificationChannelService := service.NewNotificationChannelService(ctx, notificationChannelServiceClient)
	notificationRuleService := service.NewNotificationRuleService(ctx, notificationRuleServiceClient)
	notificationTemplateService := service.NewNotificationTemplateService(ctx, notificationTemplateServiceClient)
	notificationPreferenceService := service.NewNotificationPreferenceService(ctx, notificationPreferenceServiceClient)
	monitorAlertService := service.NewMonitorAlertService(ctx, monitorAlertServiceClient)
	onlineSessionService := service.NewOnlineSessionService(ctx, onlineSessionServiceClient)
	accessKeyService := service.NewAccessKeyService(ctx, accessKeyServiceClient)
	tenantQuotaService := service.NewTenantQuotaService(ctx, tenantQuotaServiceClient)
	configService := service.NewConfigService(ctx, configServiceClient)
	aiAgentService := service.NewAiAgentService(ctx, aiAgentServiceClient)
	serverMonitorService := service.NewServerMonitorService(ctx, serverMonitorServiceClient)
	redisCacheMonitorService := service.NewRedisCacheMonitorService(ctx, redisCacheMonitorServiceClient)
	// ── register:service ── 新模块服务在此行后注册(make register 工具锚点,勿删)

	// ═══════════════════════ 四、传输层(internal/server) ═══════════════════════

	restMiddlewares := server.NewRestMiddleware(ctx, accessTokenChecker, engine, tenantQuotaServiceClient, apiAuditLogServiceClient, loginAuditLogServiceClient, operationAuditLogServiceClient)

	httpServer := server.NewRestServer(ctx, restMiddlewares,
		userService, userProfileService, roleService, tenantService, orgUnitService, positionService,
		menuService, apiService, permissionGroupService, permissionService,
		adminPortalService, previewService, taskService,
		authenticationService, loginPolicyService,
		dictTypeService, dictEntryService, languageService,
		fileService, fileTransferService,
		translatorService,
		internalMessageService, internalMessageCategoryService, internalMessageRecipientService,
		apiAuditLogService, dataAccessAuditLogService, loginAuditLogService, policyEvaluationLogService, operationAuditLogService, permissionAuditLogService,
		commentService, interactionAdminService, statsService,
		aiProviderService, aiContentService,
		aiChatService, aiConversationService, aiMessageService,
		aiKnowledgeService,
		aiUsageLogService,
		aiQueryService,
		auditExportService,
		postService, categoryService, tagService, contentModelService, componentDefinitionService, pageService,
		siteService, siteSettingService, navigationService, navigationItemService,
		mediaAssetService,
		socialService,
		webhookService,
		formSubmissionService,
		redirectService,
		experimentService,
		editLockService,
		flowRuleService,
		notificationService,
		monitorAlertService,
		onlineSessionService,
		accessKeyService,
		tenantQuotaService,
		configService,
		aiAgentService,
		serverMonitorService,
		redisCacheMonitorService,
		notificationChannelService,
		notificationRuleService,
		notificationTemplateService,
		notificationPreferenceService,
		// register:rest-arg ── 新模块服务实参在此行后追加(make register 工具锚点,勿删)
	)

	grpcMiddlewares := server.NewGrpcMiddleware(ctx)
	grpcServer, err := server.NewGrpcServer(ctx, grpcMiddlewares)
	if err != nil {
		rollback()
		return nil, nil, err
	}

	aiAgentService.StartAgentRunSubscriber(redisClient, authenticationServiceClient)
	sseServer := server.NewSseServer(ctx, internalMessageService, aiChatService, aiAgentService)

	// 协同编辑 ws（COLLAB_ADDR，默认 :6602；off 停用）。
	// off 时 NewCollab 返回 nil CollabServer——nil 不能进 servers：
	// kratos buildInstance 会对每个 Server 调 Endpoint()，typed-nil 上
	// 解引用 s.addr 直接 panic（COLLAB_ADDR=off 停用路径曾必崩）。
	servers := []kratosTransport.Server{httpServer, grpcServer, sseServer}
	if collabServer != nil {
		servers = append(servers, collabServer)
	}
	// MCP 服务器（可选）：server.mcp 配置段或 MCP_ADDR 设置时把 CMS 查询/发布能力暴露给 AI 代理
	if mcpServer := server.NewMcpServer(ctx, postService, socialService); mcpServer != nil {
		servers = append(servers, mcpServer)
	}

	return newApp(ctx, servers...), rollback, nil
}
