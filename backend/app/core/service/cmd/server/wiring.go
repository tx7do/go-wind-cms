package main

import (
	"os"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	pkgcrypto "go-wind-cms/pkg/crypto"

	"github.com/tx7do/go-utils/geoip/geolite"

	"go-wind-cms/app/core/service/internal/data"
	"go-wind-cms/app/core/service/internal/data/channel"
	"go-wind-cms/app/core/service/internal/data/client"
	"go-wind-cms/app/core/service/internal/server"
	"go-wind-cms/app/core/service/internal/service"
)

// initApp 手写装配整个应用,是本服务的依赖注入点(替代原 google/wire 生成)。
// initApp assembles the whole application by hand — the dependency injection point.
//
// 装配严格单向分层,自上而下的阅读顺序即依赖方向:
// The wiring is strictly layered; reading top-down follows the dependency direction:
//
//	基础设施 → 仓储层(data) → 服务层(service) → 传输层(server)
//
// 约定 / Conventions:
//   - 新增仓储/服务:在对应分层小节追加一行构造,并传给下游消费者;漏接由编译器在调用处报错。
//     To add a repo/service: append one constructor line in its layer section, then pass it to
//     downstream consumers; a missing connection is a compile error at the call site.
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

	// 全局存储加密器(AES-256-GCM, pkg/crypto):凭据类字段落库加密的总开关。
	// key 从 ENCRYPTION_KEY 环境变量读取;未配置时降级明文并告警(生产必须配置)。
	if encKey := os.Getenv("ENCRYPTION_KEY"); encKey != "" {
		if err := pkgcrypto.InitGlobalEncryptor(encKey, true); err != nil {
			return nil, nil, err
		}
	} else {
		log.Warnf("ENCRYPTION_KEY not set: credential encryption DISABLED (secrets stored in plaintext)")
	}

	// 认证基建:令牌缓存 → 认证器
	authenticatorOption := data.NewAuthenticatorConfig(ctx)

	redisClient, cleanupRedis, err := client.NewRedisClient(ctx)
	if err != nil {
		return nil, nil, err
	}
	cleanups = append(cleanups, cleanupRedis)

	userTokenCache := data.NewUserTokenCache(ctx, redisClient)
	vcodeCache := data.NewVCodeCache(redisClient)
	authenticator := data.NewAuthenticator(ctx, authenticatorOption, userTokenCache)

	entClient, cleanupEnt, err := client.NewEntClient(ctx)
	if err != nil {
		rollback()
		return nil, nil, err
	}
	cleanups = append(cleanups, cleanupEnt)

	crypto := data.NewPasswordCrypto()
	minIOClient := client.NewMinIoClient(ctx)

	opensearchClient, cleanupSearch, err := client.NewElasticSearchClient(ctx)
	if err != nil {
		rollback()
		return nil, nil, err
	}
	cleanups = append(cleanups, cleanupSearch)

	// ═══════════════════════ 二、仓储层(internal/data) ═══════════════════════

	// 身份与账号
	configRepo := data.NewConfigRepo(ctx, entClient, redisClient)
	userCredentialRepo := data.NewUserCredentialRepo(ctx, entClient, crypto, configRepo)
	loginPolicyRepo := data.NewLoginPolicyRepo(ctx, entClient)
	userRoleRepo := data.NewUserRoleRepo(ctx, entClient)
	userOrgUnitRepo := data.NewUserOrgUnitRepo(ctx, entClient)
	userPositionRepo := data.NewUserPositionRepo(ctx, entClient)
	membershipRoleRepo := data.NewMembershipRoleRepo(ctx, entClient)
	membershipPositionRepo := data.NewMembershipPositionRepo(ctx, entClient)
	membershipOrgUnitRepo := data.NewMembershipOrgUnitRepo(ctx, entClient)
	membershipRepo := data.NewMembershipRepo(ctx, entClient, membershipRoleRepo, membershipPositionRepo, membershipOrgUnitRepo)
	userRepo := data.NewUserRepo(ctx, entClient, userRoleRepo, userOrgUnitRepo, userPositionRepo, membershipRepo)

	// 租户
	tenantRepo := data.NewTenantRepo(ctx, entClient)

	// RBAC:权限(聚合 Api/Menu 挂载点) → 角色(聚合权限)
	rolePermissionRepo := data.NewRolePermissionRepo(ctx, entClient)
	permissionApiRepo := data.NewPermissionApiRepo(ctx, entClient)
	permissionMenuRepo := data.NewPermissionMenuRepo(ctx, entClient)
	permissionRepo := data.NewPermissionRepo(ctx, entClient, permissionApiRepo, permissionMenuRepo, rolePermissionRepo)
	roleMetadataRepo := data.NewRoleMetadataRepo(ctx, entClient)
	roleRepo := data.NewRoleRepo(ctx, entClient, rolePermissionRepo, permissionRepo, roleMetadataRepo, userRoleRepo)

	// RBAC 挂载点:菜单 / Api / 权限组
	menuRepo := data.NewMenuRepo(ctx, entClient)
	apiRepo := data.NewApiRepo(ctx, entClient)
	permissionGroupRepo := data.NewPermissionGroupRepo(ctx, entClient)

	// 组织架构
	positionRepo := data.NewPositionRepo(ctx, entClient)
	orgUnitRepo := data.NewOrgUnitRepo(ctx, entClient)

	// 审计日志
	permissionAuditLogRepo := data.NewPermissionAuditLogRepo(ctx, entClient)
	policyEvaluationLogRepo := data.NewPolicyEvaluationLogRepo(ctx, entClient)
	loginAuditLogRepo := data.NewLoginAuditLogRepo(ctx, entClient)
	apiAuditLogRepo := data.NewApiAuditLogRepo(ctx, entClient)
	operationAuditLogRepo := data.NewOperationAuditLogRepo(ctx, entClient)
	dataAccessAuditLogRepo := data.NewDataAccessAuditLogRepo(ctx, entClient)

	// 站内信
	internalMessageRecipientRepo := data.NewInternalMessageRecipientRepo(ctx, entClient)
	internalMessageRepo := data.NewInternalMessageRepo(ctx, entClient, internalMessageRecipientRepo)
	internalMessageCategoryRepo := data.NewInternalMessageCategoryRepo(ctx, entClient)

	// 任务 / 文件
	taskRepo := data.NewTaskRepo(ctx, entClient)
	backupRepo := data.NewBackupRepo(ctx, entClient)
	fileRepo := data.NewFileRepo(ctx, entClient)

	// 字典与多语言
	dictEntryI18nRepo := data.NewDictEntryI18nRepo(ctx, entClient)
	dictEntryRepo := data.NewDictEntryRepo(ctx, entClient, dictEntryI18nRepo)
	dictTypeRepo := data.NewDictTypeRepo(ctx, entClient, dictEntryRepo)
	languageRepo := data.NewLanguageRepo(ctx, entClient)

	// 内容:评论 / 互动 / 搜索 / 文章(聚合翻译/分类/标签) / 分类 / 标签 / 页面(聚合 Section)
	commentRepo := data.NewCommentRepo(ctx, entClient)
	interactionRepo := data.NewInteractionRepo(ctx, entClient)
	searchRepo := data.NewSearchRepo(ctx, opensearchClient)
	postTranslationRepo := data.NewPostTranslationRepo(ctx, entClient)
	postCategoryRepo := data.NewPostCategoryRepo(ctx, entClient)
	postTagRepo := data.NewPostTagRepo(ctx, entClient)
	postRepo := data.NewPostRepo(ctx, entClient, postTranslationRepo, postCategoryRepo, postTagRepo)
	statsRepo := data.NewStatsRepo(ctx, entClient)
	aiProviderRepo := data.NewAiProviderRepo(ctx, entClient)
	aiConversationRepo := data.NewAiConversationRepo(ctx, entClient)
	aiMessageRepo := data.NewAiMessageRepo(ctx, entClient)
	aiUsageLogRepo := data.NewAiUsageLogRepo(ctx, entClient)
	aiKnowledgeRepo := data.NewAiKnowledgeRepo(ctx, entClient)
	contentModelRepo := data.NewContentModelRepo(ctx, entClient)
	categoryTranslationRepo := data.NewCategoryTranslationRepo(ctx, entClient)
	categoryRepo := data.NewCategoryRepo(ctx, entClient, categoryTranslationRepo)
	tagTranslationRepo := data.NewTagTranslationRepo(ctx, entClient)
	tagRepo := data.NewTagRepo(ctx, entClient, tagTranslationRepo)
	pageTranslationRepo := data.NewPageTranslationRepo(ctx, entClient)
	sectionTranslationRepo := data.NewSectionTranslationRepo(ctx, entClient)
	sectionRepo := data.NewSectionRepo(ctx, entClient, sectionTranslationRepo)
	pageRepo := data.NewPageRepo(ctx, entClient, pageTranslationRepo, sectionRepo)

	// 内容版本快照（保存自动落版本，恢复走常规 Update 链路）
	postVersionRepo := data.NewPostVersionRepo(ctx, entClient)
	pageVersionRepo := data.NewPageVersionRepo(ctx, entClient)

	// Webhook 订阅与投递记录
	webhookRepo := data.NewWebhookRepo(ctx, entClient)

	// 表单提交（页面 FORM 区块）
	formSubmissionRepo := data.NewFormSubmissionRepo(ctx, entClient)

	// URL 重定向
	redirectRepo := data.NewRedirectRepo(ctx, entClient)
	experimentRepo := data.NewExperimentRepo(ctx, entClient)
	collabDocStateRepo := data.NewCollabDocStateRepo(ctx, entClient)

	// 站点配置
	siteRepo := data.NewSiteRepo(ctx, entClient)
	siteSettingRepo := data.NewSiteSettingRepo(ctx, entClient)
	navigationItemRepo := data.NewNavigationItemRepo(ctx, entClient)
	navigationRepo := data.NewNavigationRepo(ctx, entClient, navigationItemRepo)

	// 媒体
	mediaVariantRepo := data.NewMediaVariantRepo(ctx, entClient)
	mediaAssetRepo := data.NewMediaAssetRepo(ctx, entClient, mediaVariantRepo)
	contentReferenceRepo := data.NewContentReferenceRepo(ctx, entClient)
	flowRepo := data.NewFlowRepo(ctx, entClient)

	// 社交分发
	socialAccountRepo := data.NewSocialAccountRepo(ctx, entClient)
	socialDistributionRepo := data.NewSocialDistributionRepo(ctx, entClient)
	socialStatsHistoryRepo := data.NewSocialStatsHistoryRepo(ctx, entClient)
	socialInboxRepo := data.NewSocialInboxRepo(ctx, entClient)

	// 通知域：渠道花名册 / 投递台账 / 路由规则 / 模板 / 用户偏好（全部平台级，无租户列）
	notificationChannelRepo := data.NewNotificationChannelRepo(ctx, entClient)
	monitorAlertRuleRepo := data.NewMonitorAlertRuleRepo(ctx, entClient)
	notificationDeliveryRepo := data.NewNotificationDeliveryRepo(ctx, entClient)
	notificationRuleRepo := data.NewNotificationRuleRepo(ctx, entClient)
	notificationTemplateRepo := data.NewNotificationTemplateRepo(ctx, entClient)
	notificationPreferenceRepo := data.NewNotificationPreferenceRepo(ctx, entClient)

	// ── register:repo ── 新模块仓储在此行后注册(make register 工具锚点,勿删)

	// 访问凭证（AK/SK）：OpenAPI 机器凭证
	accessKeyRepo := data.NewAccessKeyRepo(ctx, entClient)
	accessKeyRateLimiter := data.NewAccessKeyRateLimiter(ctx, redisClient)
	tenantQuotaRepo := data.NewTenantQuotaRepo(ctx, entClient)
	quotaPlanRepo := data.NewQuotaPlanRepo(ctx, entClient)
	contentHealthService := service.NewContentHealthService(ctx, postRepo)
	tenantQuotaService := service.NewTenantQuotaService(ctx, tenantQuotaRepo, quotaPlanRepo)
	configService := service.NewConfigService(ctx, configRepo)
	aiAgentAppRepo := data.NewAiAgentAppRepo(ctx, entClient)
	aiAgentRunRepo := data.NewAiAgentRunRepo(ctx, entClient)
	aiAgentService := service.NewAiAgentService(ctx, aiAgentAppRepo, aiAgentRunRepo, aiProviderRepo, aiKnowledgeRepo, aiUsageLogRepo, redisClient)

	// ═══════════════════════ 三、服务层(internal/service) ═══════════════════════

	// 认证与登录策略
	geoClient, _ := geolite.NewClient()
	authenticationService := service.NewAuthenticationService(ctx, authenticator, userCredentialRepo, userRepo, roleRepo, tenantRepo, permissionRepo, redisClient, vcodeCache, loginPolicyRepo, geoClient)
	mfaService := service.NewMFAService(ctx, authenticationService)
	loginPolicyService := service.NewLoginPolicyService(ctx, loginPolicyRepo)
	userCredentialService := service.NewUserCredentialService(ctx, userCredentialRepo)

	// 任务 / 文件
	taskService := service.NewTaskService(ctx, taskRepo, userRepo, backupRepo, minIOClient)
	backupArchiveService := service.NewBackupArchiveService(ctx, minIOClient)

	// Webhook 事件外发（订阅 CRUD + 扇出/投递 worker；依赖 taskService）
	webhookService := service.NewWebhookService(ctx, webhookRepo, taskService)

	// 表单提交（页面 FORM 区块；校验依赖 pageRepo）
	formSubmissionService := service.NewFormSubmissionService(ctx, formSubmissionRepo, pageRepo)

	// URL 重定向
	redirectService := service.NewRedirectService(ctx, redirectRepo)

	// 内容实验（A/B 分流 + 受众定向）
	experimentService := service.NewExperimentService(ctx, experimentRepo)

	// 编辑锁（协同编辑 v1，Redis 软锁）
	editLockService := service.NewEditLockService(ctx, redisClient, userRepo)

	// 协同文档状态（collab ws 引擎的持久化后端）
	collabDocStateService := service.NewCollabDocStateService(ctx, collabDocStateRepo)

	// 邮件模板（系统通知出站可视化模板）

	// 邮件出站 SMTP 配置
	fileService := service.NewFileService(ctx, fileRepo, minIOClient)

	// 字典与多语言
	dictTypeService := service.NewDictTypeService(ctx, dictTypeRepo)
	dictEntryService := service.NewDictEntryService(ctx, dictEntryRepo)
	languageService := service.NewLanguageService(ctx, languageRepo)

	// 租户与身份
	tenantService := service.NewTenantService(ctx, tenantRepo, userRepo, userCredentialRepo, roleRepo)
	internalMessageService := service.NewInternalMessageService(ctx, internalMessageRepo, internalMessageCategoryRepo, internalMessageRecipientRepo, userRepo, redisClient, taskService)
	flowService := service.NewFlowService(ctx, flowRepo, postRepo, pageRepo, taskService, internalMessageService)
	userService := service.NewUserService(ctx, userRepo, roleRepo, userCredentialRepo, positionRepo, orgUnitRepo, tenantRepo, membershipRepo, tenantQuotaService, flowService)
	roleFieldPermissionRepo := data.NewRoleFieldPermissionRepo(ctx, entClient)
	roleService := service.NewRoleService(ctx, roleRepo, tenantRepo, userRoleRepo, userRepo, roleFieldPermissionRepo)
	positionService := service.NewPositionService(ctx, positionRepo, orgUnitRepo)
	orgUnitService := service.NewOrgUnitService(ctx, orgUnitRepo, userRepo)

	// RBAC
	menuService := service.NewMenuService(ctx, menuRepo)
	apiService := service.NewApiService(ctx, apiRepo)
	permissionService := service.NewPermissionService(ctx, permissionRepo, permissionGroupRepo, menuRepo, apiRepo, roleRepo)
	permissionGroupService := service.NewPermissionGroupService(ctx, permissionGroupRepo, permissionRepo)

	// 审计日志
	permissionAuditLogService := service.NewPermissionAuditLogService(ctx, permissionAuditLogRepo)
	policyEvaluationLogService := service.NewPolicyEvaluationLogService(ctx, policyEvaluationLogRepo)
	loginAuditLogService := service.NewLoginAuditLogService(ctx, loginAuditLogRepo)
	apiAuditLogService := service.NewApiAuditLogService(ctx, apiAuditLogRepo, apiRepo)
	operationAuditLogService := service.NewOperationAuditLogService(ctx, operationAuditLogRepo)
	dataAccessAuditLogService := service.NewDataAccessAuditLogService(ctx, dataAccessAuditLogRepo)

	// 站内信
	internalMessageCategoryService := service.NewInternalMessageCategoryService(ctx, internalMessageCategoryRepo)
	internalMessageRecipientService := service.NewInternalMessageRecipientService(ctx, internalMessageRepo, internalMessageRecipientRepo)

	// 通知域（投递台账 + 渠道花名册 + 路由规则 + 模板 + 用户偏好）
	// 注册表先注册自包含的渠道 sender；INTERNAL 渠道的 sender 依赖两个 service 互相就绪，
	// 在下方"双向引用缝"处补注册（构造期互相传会成环）。
	channelRegistry := channel.NewRegistry()
	channelRegistry.Register(channel.NewEmailSender(notificationChannelRepo))
	channelRegistry.Register(channel.NewWebhookSender(notificationChannelRepo))
	// 扩展渠道（EXT_* 39 种）：settings 列配置 → go-utils/notify 统一发送
	channelRegistry.Register(channel.NewSettingsSender(notificationChannelRepo))
	// 短信渠道（阿里云/腾讯云）：仅承载验证码类事件
	channelRegistry.Register(channel.NewSmsSender(notificationChannelRepo))
	notificationService := service.NewNotificationService(ctx, notificationDeliveryRepo, notificationRuleRepo, notificationTemplateRepo, channelRegistry)
	notificationChannelService := service.NewNotificationChannelService(ctx, notificationChannelRepo, notificationService)
	notificationRuleService := service.NewNotificationRuleService(ctx, notificationRuleRepo, notificationChannelRepo, notificationService)
	notificationTemplateService := service.NewNotificationTemplateService(ctx, notificationTemplateRepo)
	notificationPreferenceService := service.NewNotificationPreferenceService(ctx, notificationPreferenceRepo, internalMessageCategoryRepo)
	monitorAlertService := service.NewMonitorAlertService(ctx, monitorAlertRuleRepo, entClient, redisClient)
	onlineSessionService := service.NewOnlineSessionService(ctx, authenticator)
	accessKeyService := service.NewAccessKeyService(ctx, accessKeyRepo, authenticator, accessKeyRateLimiter)
	serverMonitorService := service.NewServerMonitorService(ctx, entClient)
	redisCacheMonitorService := service.NewRedisCacheMonitorService(ctx, redisClient)
	monitorAlertService.RegisterNotifier(notificationService)

	// 双向引用缝：站内信的系统通知经投递缝（SendDirect）走渠道路由，
	// INTERNAL 渠道的实现又反向持有站内信 service——两个 service 构造完之后再互相注册。
	channelRegistry.Register(service.NewInternalMessageSender(internalMessageService, notificationPreferenceRepo, ctx.NewLoggerHelper("internal-message-sender/channel")))
	internalMessageService.RegisterNotifier(notificationService)
	authenticationService.RegisterNotifier(notificationService)

	// 站内信投递缝：配额水位告警经系统通知落租户管理员收件箱
	tenantQuotaService.RegisterInternalMessages(internalMessageService)
	contentHealthService.RegisterInternalMessages(internalMessageService)

	// 内容
	commentAIDeps := &service.CommentAIModerationDeps{ProviderRepo: aiProviderRepo, UsageLogRepo: aiUsageLogRepo}
	commentService := service.NewCommentService(ctx, commentRepo, flowService, configRepo, commentAIDeps)
	interactionService := service.NewInteractionService(ctx, interactionRepo, postRepo)
	interactionAdminService := service.NewInteractionAdminService(ctx, interactionRepo, operationAuditLogRepo)
	statsService := service.NewStatsService(ctx, statsRepo)
	aiProviderService := service.NewAiProviderService(ctx, aiProviderRepo)
	aiContentService := service.NewAiContentService(ctx, aiProviderRepo, aiUsageLogRepo, aiKnowledgeRepo, entClient, menuRepo, roleRepo, userRepo)
	aiChatService := service.NewAiChatService(ctx, aiConversationRepo, aiMessageRepo, aiProviderRepo, aiUsageLogRepo, aiKnowledgeRepo, entClient)
	aiConversationService := service.NewAiConversationService(ctx, aiConversationRepo)
	aiMessageService := service.NewAiMessageService(ctx, aiMessageRepo)
	aiKnowledgeService := service.NewAiKnowledgeService(ctx, aiKnowledgeRepo, aiProviderRepo, aiUsageLogRepo, taskService)
	aiUsageLogService := service.NewAiUsageLogService(ctx, aiUsageLogRepo, aiProviderRepo)
	aiQueryService := service.NewAiQueryService(ctx, aiProviderRepo, aiUsageLogRepo, entClient, internalMessageRepo, internalMessageService)
	searchService := service.NewSearchService(ctx, searchRepo, postRepo)
	referenceTrackingService := service.NewReferenceTrackingService(ctx, mediaAssetRepo, contentReferenceRepo)
	postService := service.NewPostService(ctx, postRepo, searchService, taskService, contentModelRepo, postVersionRepo, webhookService, internalMessageService, referenceTrackingService, flowService)
	categoryService := service.NewCategoryService(ctx, categoryRepo, contentModelRepo)
	tagService := service.NewTagService(ctx, tagRepo)
	pageService := service.NewPageService(ctx, pageRepo, contentModelRepo, pageVersionRepo, webhookService, internalMessageService, referenceTrackingService, flowService)
	contentModelService := service.NewContentModelService(ctx, contentModelRepo, postRepo, pageRepo)

	// 站点配置
	siteService := service.NewSiteService(ctx, siteRepo)
	siteSettingService := service.NewSiteSettingService(ctx, siteSettingRepo)
	navigationService := service.NewNavigationService(ctx, navigationRepo)
	navigationItemService := service.NewNavigationItemService(ctx, navigationItemRepo)

	// 媒体
	mediaAssetService := service.NewMediaAssetService(ctx, mediaAssetRepo, contentReferenceRepo, postRepo, pageRepo, flowService)

	// 社交分发
	socialService := service.NewSocialService(ctx, socialAccountRepo, socialDistributionRepo, socialStatsHistoryRepo, mediaAssetRepo, fileRepo, minIOClient, postRepo, taskService, internalMessageService, socialInboxRepo)

	// ── register:service ── 新模块服务在此行后注册(make register 工具锚点,勿删)

	// ═══════════════════════ 四、传输层(internal/server) ═══════════════════════

	grpcMiddlewares := server.NewGrpcMiddleware(ctx)

	grpcServer, err := server.NewGrpcServer(ctx, grpcMiddlewares,
		authenticationService, mfaService, backupArchiveService, loginPolicyService, userCredentialService,
		taskService,
		monitorAlertService,
		onlineSessionService,
		accessKeyService,
		tenantQuotaService,
		configService,
		aiAgentService,
		serverMonitorService,
		redisCacheMonitorService,
		fileService,
		dictTypeService, dictEntryService, languageService,
		tenantService, userService, roleService, positionService, orgUnitService,
		menuService, apiService, permissionService, permissionGroupService, permissionAuditLogService, policyEvaluationLogService,
		loginAuditLogService, apiAuditLogService, operationAuditLogService, dataAccessAuditLogService,
		internalMessageService, internalMessageCategoryService, internalMessageRecipientService,
		commentService,
		interactionService, interactionAdminService,
		statsService,
		aiProviderService, aiContentService,
		aiChatService, aiConversationService, aiMessageService,
		aiKnowledgeService,
		aiUsageLogService,
		aiQueryService,
		flowService,
		postService, categoryService, tagService, pageService,
		contentModelService,
		siteService, siteSettingService, navigationService, navigationItemService,
		mediaAssetService,
		socialService,
		webhookService,
		formSubmissionService,
		redirectService,
		experimentService,
		editLockService,
		collabDocStateService,
		notificationService,
		notificationChannelService,
		notificationRuleService,
		notificationTemplateService,
		notificationPreferenceService,
		// register:grpc-arg ── 新模块服务实参在此行后追加(make register 工具锚点,勿删)
	)
	if err != nil {
		rollback()
		return nil, nil, err
	}

	asynqServer := server.NewAsynqServer(ctx, taskService, searchService, aiKnowledgeService, aiQueryService, aiContentService, socialService, postService, webhookService, internalMessageService, flowService, notificationService, monitorAlertService, tenantQuotaService, contentHealthService)

	return newApp(ctx, grpcServer, asynqServer), rollback, nil
}
