package server

import (
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/transport/grpc"

	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"github.com/tx7do/kratos-bootstrap/rpc"

	"go-wind-cms/app/core/service/internal/service"

	auditV1 "go-wind-cms/api/gen/go/audit/service/v1"
	authenticationV1 "go-wind-cms/api/gen/go/authentication/service/v1"
	commentV1 "go-wind-cms/api/gen/go/comment/service/v1"
	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
	dictV1 "go-wind-cms/api/gen/go/dict/service/v1"
	identityV1 "go-wind-cms/api/gen/go/identity/service/v1"
	internalMessageV1 "go-wind-cms/api/gen/go/internal_message/service/v1"
	mediaV1 "go-wind-cms/api/gen/go/media/service/v1"
	permissionV1 "go-wind-cms/api/gen/go/permission/service/v1"
	resourceV1 "go-wind-cms/api/gen/go/resource/service/v1"
	siteV1 "go-wind-cms/api/gen/go/site/service/v1"
	storageV1 "go-wind-cms/api/gen/go/storage/service/v1"
	taskV1 "go-wind-cms/api/gen/go/task/service/v1"

	"go-wind-cms/pkg/middleware/ent"
)

func NewGrpcMiddleware(ctx *bootstrap.Context) []middleware.Middleware {
	var ms []middleware.Middleware
	ms = append(ms, logging.Server(ctx.GetLogger()))
	ms = append(ms, ent.Server())
	return ms
}

// NewGrpcServer new a gRPC server.
func NewGrpcServer(
	ctx *bootstrap.Context,
	middlewares []middleware.Middleware,

	authenticationService *service.AuthenticationService,
	loginPolicyService *service.LoginPolicyService,

	taskService *service.TaskService,

	uEditorService *service.UEditorService,
	fileService *service.FileService,

	dictTypeService *service.DictTypeService,
	dictEntryService *service.DictEntryService,
	languageService *service.LanguageService,

	tenantService *service.TenantService,
	userService *service.UserService,
	roleService *service.RoleService,
	positionService *service.PositionService,
	orgUnitService *service.OrgUnitService,

	menuService *service.MenuService,
	apiService *service.ApiService,
	permissionService *service.PermissionService,
	permissionGroupService *service.PermissionGroupService,
	permissionAuditLogService *service.PermissionAuditLogService,
	policyEvaluationLogService *service.PolicyEvaluationLogService,

	loginAuditLogService *service.LoginAuditLogService,
	apiAuditLogService *service.ApiAuditLogService,
	operationAuditLogService *service.OperationAuditLogService,
	dataAccessAuditLogService *service.DataAccessAuditLogService,

	internalMessageService *service.InternalMessageService,
	internalMessageCategoryService *service.InternalMessageCategoryService,
	internalMessageRecipientService *service.InternalMessageRecipientService,

	commentService *service.CommentService,

	postService *service.PostService,
	categoryService *service.CategoryService,
	tagService *service.TagService,
	pageService *service.PageService,

	siteSettingService *service.SiteSettingService,
	navigationService *service.NavigationService,

	mediaAssetService *service.MediaAssetService,
) (*grpc.Server, error) {
	cfg := ctx.GetConfig()

	if cfg == nil || cfg.Server == nil || cfg.Server.Grpc == nil {
		return nil, nil
	}

	srv, err := rpc.CreateGrpcServer(cfg, middlewares...)
	if err != nil {
		return nil, err
	}

	taskV1.RegisterTaskServiceServer(srv, taskService)

	authenticationV1.RegisterLoginPolicyServiceServer(srv, loginPolicyService)
	authenticationV1.RegisterAuthenticationServiceServer(srv, authenticationService)

	dictV1.RegisterDictTypeServiceServer(srv, dictTypeService)
	dictV1.RegisterDictEntryServiceServer(srv, dictEntryService)
	dictV1.RegisterLanguageServiceServer(srv, languageService)

	resourceV1.RegisterApiServiceServer(srv, apiService)
	resourceV1.RegisterMenuServiceServer(srv, menuService)

	permissionV1.RegisterPermissionServiceServer(srv, permissionService)
	permissionV1.RegisterPermissionGroupServiceServer(srv, permissionGroupService)
	permissionV1.RegisterPolicyEvaluationLogServiceServer(srv, policyEvaluationLogService)
	permissionV1.RegisterRoleServiceServer(srv, roleService)

	identityV1.RegisterUserServiceServer(srv, userService)
	identityV1.RegisterOrgUnitServiceServer(srv, orgUnitService)
	identityV1.RegisterPositionServiceServer(srv, positionService)
	identityV1.RegisterTenantServiceServer(srv, tenantService)

	auditV1.RegisterLoginAuditLogServiceServer(srv, loginAuditLogService)
	auditV1.RegisterApiAuditLogServiceServer(srv, apiAuditLogService)
	auditV1.RegisterOperationAuditLogServiceServer(srv, operationAuditLogService)
	auditV1.RegisterDataAccessAuditLogServiceServer(srv, dataAccessAuditLogService)
	auditV1.RegisterPermissionAuditLogServiceServer(srv, permissionAuditLogService)

	storageV1.RegisterFileServiceServer(srv, fileService)
	storageV1.RegisterUEditorServiceServer(srv, uEditorService)

	internalMessageV1.RegisterInternalMessageServiceServer(srv, internalMessageService)
	internalMessageV1.RegisterInternalMessageCategoryServiceServer(srv, internalMessageCategoryService)
	internalMessageV1.RegisterInternalMessageRecipientServiceServer(srv, internalMessageRecipientService)

	commentV1.RegisterCommentServiceServer(srv, commentService)

	contentV1.RegisterPostServiceServer(srv, postService)
	contentV1.RegisterCategoryServiceServer(srv, categoryService)
	contentV1.RegisterTagServiceServer(srv, tagService)
	contentV1.RegisterPageServiceServer(srv, pageService)

	siteV1.RegisterSiteSettingServiceServer(srv, siteSettingService)
	siteV1.RegisterNavigationServiceServer(srv, navigationService)

	mediaV1.RegisterMediaAssetServiceServer(srv, mediaAssetService)

	return srv, nil
}
