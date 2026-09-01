package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	"go-wind-cms/app/app/service/internal/data"
	"go-wind-cms/app/app/service/internal/server"
	"go-wind-cms/app/app/service/internal/service"
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
// 本服务是 C 端 BFF 网关:不直连数据库,仓储层由对 core-service 的 gRPC 服务客户端替代。
// 约定 / Conventions:
//   - 新增模块:先加服务客户端行,再加服务构造行,并传给 NewRestServer;漏接由编译器在调用处报错。
//     To add a module: append a service-client line, then a service constructor, then pass it to
//     NewRestServer; a missing connection is a compile error at the call site.
//   - 本文件只做构造与传参,不写业务逻辑。
//     Construction and parameter passing only; no business logic in this file.
func initApp(ctx *bootstrap.Context) (*kratos.App, func(), error) {
	// cleanup 注册表:本服务无持有 cleanup 的资源,rollback 为空实现,保持与
	// admin/core 两服务同构的返回签名。
	// No cleanup-owning resources in this service; rollback is a no-op kept for
	// signature parity with the admin/core wiring.
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
	minIOClient := data.NewMinIoClient(ctx)

	// 认证:认证服务客户端 → 访问令牌校验器;租户解析器供中间件解析 Host→租户
	authenticationServiceClient := data.NewAuthenticationServiceClient(ctx, discovery)
	accessTokenChecker := auth.NewTokenChecker(ctx, authenticationServiceClient, clientType)
	tenantServiceClient := data.NewTenantServiceClient(ctx, discovery)
	tenantResolver := data.NewTenantResolver(tenantServiceClient)

	// ═══════════════════════ 二、服务客户端(internal/data) ═══════════════════════

	// 身份
	userServiceClient := data.NewUserServiceClient(ctx, discovery)
	roleServiceClient := data.NewRoleServiceClient(ctx, discovery)
	orgUnitServiceClient := data.NewOrgUnitServiceClient(ctx, discovery)
	positionServiceClient := data.NewPositionServiceClient(ctx, discovery)
	userCredentialServiceClient := data.NewUserCredentialServiceClient(ctx, discovery)

	// 文件
	fileServiceClient := data.NewFileServiceClient(ctx, discovery)

	// 内容
	postServiceClient := data.NewPostServiceClient(ctx, discovery)
	categoryServiceClient := data.NewCategoryServiceClient(ctx, discovery)
	commentServiceClient := data.NewCommentServiceClient(ctx, discovery)
	interactionServiceClient := data.NewInteractionServiceClient(ctx, discovery)
	tagServiceClient := data.NewTagServiceClient(ctx, discovery)
	pageServiceClient := data.NewPageServiceClient(ctx, discovery)

	// 站点配置
	siteServiceClient := data.NewSiteServiceClient(ctx, discovery)
	navigationServiceClient := data.NewNavigationServiceClient(ctx, discovery)

	// ── register:client ── 新模块服务客户端在此行后注册(make register 工具锚点,勿删)

	// ═══════════════════════ 三、服务层(internal/service) ═══════════════════════

	// 认证
	authenticationService := service.NewAuthenticationService(ctx, authenticationServiceClient)

	// 文件
	fileTransferService := service.NewFileTransferService(ctx, minIOClient, fileServiceClient)

	// 身份
	userProfileService := service.NewUserProfileService(ctx, userServiceClient, tenantServiceClient, orgUnitServiceClient, positionServiceClient, roleServiceClient, userCredentialServiceClient)

	// 内容
	postService := service.NewPostService(ctx, postServiceClient)
	categoryService := service.NewCategoryService(ctx, categoryServiceClient)
	commentService := service.NewCommentService(ctx, commentServiceClient)
	interactionService := service.NewInteractionService(ctx, interactionServiceClient)
	tagService := service.NewTagService(ctx, tagServiceClient)
	pageService := service.NewPageService(ctx, pageServiceClient)

	// 站点配置
	navigationService := service.NewNavigationService(ctx, navigationServiceClient)
	siteService := service.NewSiteService(ctx, siteServiceClient)

	// ── register:service ── 新模块服务在此行后注册(make register 工具锚点,勿删)

	// ═══════════════════════ 四、传输层(internal/server) ═══════════════════════

	restMiddlewares := server.NewRestMiddleware(ctx, accessTokenChecker, engine, tenantResolver)

	httpServer := server.NewRestServer(ctx, restMiddlewares,
		authenticationService,
		fileTransferService,
		userProfileService,
		postService, categoryService, commentService, interactionService, tagService, pageService,
		navigationService, siteService,
		// register:rest-arg ── 新模块服务实参在此行后追加(make register 工具锚点,勿删)
	)

	grpcMiddlewares := server.NewGrpcMiddleware(ctx)
	grpcServer, err := server.NewGrpcServer(ctx, grpcMiddlewares)
	if err != nil {
		rollback()
		return nil, nil, err
	}

	sseServer := server.NewSseServer(ctx, authenticationServiceClient)

	return newApp(ctx, httpServer, grpcServer, sseServer), rollback, nil
}
