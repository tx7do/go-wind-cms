package data

import (
	"os"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	authzEngine "github.com/tx7do/kratos-authz/engine"
	"github.com/tx7do/kratos-authz/engine/noop"

	"github.com/go-kratos/kratos/v2/registry"

	"github.com/tx7do/kratos-bootstrap/bootstrap"
	redisClient "github.com/tx7do/kratos-bootstrap/cache/redis"
	bRegistry "github.com/tx7do/kratos-bootstrap/registry"
	"github.com/tx7do/kratos-bootstrap/rpc"

	"github.com/tx7do/go-utils/captcha"
	"github.com/tx7do/go-utils/translator"
	"github.com/tx7do/go-utils/translator/google"

	"go-wind-cms/pkg/translator/baidu"

	auditV1 "go-wind-cms/api/gen/go/audit/service/v1"
	authenticationV1 "go-wind-cms/api/gen/go/authentication/service/v1"
	commentV1 "go-wind-cms/api/gen/go/comment/service/v1"
	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
	dictV1 "go-wind-cms/api/gen/go/dict/service/v1"
	identityV1 "go-wind-cms/api/gen/go/identity/service/v1"
	internalMessageV1 "go-wind-cms/api/gen/go/internal_message/service/v1"
	interactionV1 "go-wind-cms/api/gen/go/interaction/service/v1"
	mediaV1 "go-wind-cms/api/gen/go/media/service/v1"
	permissionV1 "go-wind-cms/api/gen/go/permission/service/v1"
	siteV1 "go-wind-cms/api/gen/go/site/service/v1"
	storageV1 "go-wind-cms/api/gen/go/storage/service/v1"
	taskV1 "go-wind-cms/api/gen/go/task/service/v1"

	"go-wind-cms/pkg/oss"
	"go-wind-cms/pkg/serviceid"
)

func NewClientType() authenticationV1.ClientType {
	return authenticationV1.ClientType_admin
}

// NewRedisClient 创建Redis客户端
func NewRedisClient(ctx *bootstrap.Context) (*redis.Client, func(), error) {
	cfg := ctx.GetConfig()
	if cfg == nil {
		return nil, func() {}, nil
	}

	l := ctx.NewLoggerHelper("redis/data/admin-service")

	cli := redisClient.NewClient(cfg.Data, l)

	return cli, func() {
		if err := cli.Close(); err != nil {
			l.Error(err)
		}
	}, nil
}

func NewCaptcha(rdb *redis.Client) *captcha.Captcha {
	captchaInstance := captcha.NewCaptcha(rdb,
		captcha.WithDriverType(captcha.DriverString),
		captcha.WithExpire(10*time.Minute),
		captcha.WithKeyPrefix(serviceid.ProjectName+":captcha"),
		captcha.WithStringCount(6),
		captcha.WithStringSource("ABCDEFGHJKLMNPQRSTUVWXYZ23456789"),
	)
	return captchaInstance
}

// NewDiscovery 创建服务发现客户端
func NewDiscovery(ctx *bootstrap.Context) registry.Discovery {
	cfg := ctx.GetConfig()
	if cfg == nil {
		return nil
	}

	discovery, err := bRegistry.NewDiscovery(cfg.Registry)
	if err != nil {
		return nil
	}

	NewDtmDriver(discovery)

	return discovery
}

func NewMinIoClient(ctx *bootstrap.Context) *oss.MinIOClient {
	return oss.NewMinIoClient(ctx.GetConfig(), ctx.GetLogger())
}

// NewTranslator 创建翻译器（供 admin 一键翻译使用）。
// 引擎与凭据由环境变量控制，全部缺省时保持原有行为（google 免费接口 v1，
// 无需凭据，但端点在大陆网络不可达且为非官方接口）：
//
//	TRANSLATOR_PROVIDER    = google | baidu（默认 google）
//	TRANSLATOR_VERSION     = v1 | v2 | v3（google 专用；v2/v3 为官方 API）
//	TRANSLATOR_API_KEY     = google 官方 API key（v2/v3 必填）
//	BAIDU_TRANSLATE_APP_ID = 百度通用翻译 appid（provider=baidu 必填）
//	BAIDU_TRANSLATE_SECRET = 百度通用翻译密钥（provider=baidu 必填）
func NewTranslator(_ *bootstrap.Context) translator.Translator {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("TRANSLATOR_PROVIDER"))) {
	case "baidu":
		return baidu.NewTranslator(
			baidu.WithAppID(strings.TrimSpace(os.Getenv("BAIDU_TRANSLATE_APP_ID"))),
			baidu.WithSecret(strings.TrimSpace(os.Getenv("BAIDU_TRANSLATE_SECRET"))),
		)
	}

	opts := []google.Option{google.WithVersion("v1")}
	if v := strings.TrimSpace(os.Getenv("TRANSLATOR_VERSION")); v != "" {
		opts[0] = google.WithVersion(v)
	}
	if k := strings.TrimSpace(os.Getenv("TRANSLATOR_API_KEY")); k != "" {
		opts = append(opts, google.WithApiKey(k))
	}
	return google.NewTranslator(opts...)
}

// NewAuthorizer 创建权鉴器
func NewAuthorizer() authzEngine.Engine {
	return noop.State{}
}

func NewAuthenticationServiceClient(ctx *bootstrap.Context, r registry.Discovery) authenticationV1.AuthenticationServiceClient {
	cli, err := rpc.CreateGrpcClient(ctx.Context(), r, serviceid.NewDiscoveryName(serviceid.CoreService), ctx.GetConfig())
	if err != nil {
		return nil
	}

	return authenticationV1.NewAuthenticationServiceClient(cli)
}

func NewUserCredentialServiceClient(ctx *bootstrap.Context, r registry.Discovery) authenticationV1.UserCredentialServiceClient {
	cli, err := rpc.CreateGrpcClient(ctx.Context(), r, serviceid.NewDiscoveryName(serviceid.CoreService), ctx.GetConfig())
	if err != nil {
		return nil
	}

	return authenticationV1.NewUserCredentialServiceClient(cli)
}

func NewLoginPolicyServiceClient(ctx *bootstrap.Context, r registry.Discovery) authenticationV1.LoginPolicyServiceClient {
	cli, err := rpc.CreateGrpcClient(ctx.Context(), r, serviceid.NewDiscoveryName(serviceid.CoreService), ctx.GetConfig())
	if err != nil {
		return nil
	}

	return authenticationV1.NewLoginPolicyServiceClient(cli)
}

func NewUserServiceClient(ctx *bootstrap.Context, r registry.Discovery) identityV1.UserServiceClient {
	cli, err := rpc.CreateGrpcClient(ctx.Context(), r, serviceid.NewDiscoveryName(serviceid.CoreService), ctx.GetConfig())
	if err != nil {
		return nil
	}

	return identityV1.NewUserServiceClient(cli)
}

func NewTenantServiceClient(ctx *bootstrap.Context, r registry.Discovery) identityV1.TenantServiceClient {
	cli, err := rpc.CreateGrpcClient(ctx.Context(), r, serviceid.NewDiscoveryName(serviceid.CoreService), ctx.GetConfig())
	if err != nil {
		return nil
	}

	return identityV1.NewTenantServiceClient(cli)
}

func NewRoleServiceClient(ctx *bootstrap.Context, r registry.Discovery) permissionV1.RoleServiceClient {
	cli, err := rpc.CreateGrpcClient(ctx.Context(), r, serviceid.NewDiscoveryName(serviceid.CoreService), ctx.GetConfig())
	if err != nil {
		return nil
	}

	return permissionV1.NewRoleServiceClient(cli)
}

func NewOrgUnitServiceClient(ctx *bootstrap.Context, r registry.Discovery) identityV1.OrgUnitServiceClient {
	cli, err := rpc.CreateGrpcClient(ctx.Context(), r, serviceid.NewDiscoveryName(serviceid.CoreService), ctx.GetConfig())
	if err != nil {
		return nil
	}

	return identityV1.NewOrgUnitServiceClient(cli)
}

func NewPositionServiceClient(ctx *bootstrap.Context, r registry.Discovery) identityV1.PositionServiceClient {
	cli, err := rpc.CreateGrpcClient(ctx.Context(), r, serviceid.NewDiscoveryName(serviceid.CoreService), ctx.GetConfig())
	if err != nil {
		return nil
	}

	return identityV1.NewPositionServiceClient(cli)
}

func NewInternalMessageCategoryServiceClient(ctx *bootstrap.Context, r registry.Discovery) internalMessageV1.InternalMessageCategoryServiceClient {
	cli, err := rpc.CreateGrpcClient(ctx.Context(), r, serviceid.NewDiscoveryName(serviceid.CoreService), ctx.GetConfig())
	if err != nil {
		return nil
	}

	return internalMessageV1.NewInternalMessageCategoryServiceClient(cli)
}

func NewInternalMessageServiceClient(ctx *bootstrap.Context, r registry.Discovery) internalMessageV1.InternalMessageServiceClient {
	cli, err := rpc.CreateGrpcClient(ctx.Context(), r, serviceid.NewDiscoveryName(serviceid.CoreService), ctx.GetConfig())
	if err != nil {
		return nil
	}

	return internalMessageV1.NewInternalMessageServiceClient(cli)
}

func NewInternalMessageRecipientServiceClient(ctx *bootstrap.Context, r registry.Discovery) internalMessageV1.InternalMessageRecipientServiceClient {
	cli, err := rpc.CreateGrpcClient(ctx.Context(), r, serviceid.NewDiscoveryName(serviceid.CoreService), ctx.GetConfig())
	if err != nil {
		return nil
	}

	return internalMessageV1.NewInternalMessageRecipientServiceClient(cli)
}

func NewFileServiceClient(ctx *bootstrap.Context, r registry.Discovery) storageV1.FileServiceClient {
	cli, err := rpc.CreateGrpcClient(ctx.Context(), r, serviceid.NewDiscoveryName(serviceid.CoreService), ctx.GetConfig())
	if err != nil {
		return nil
	}

	return storageV1.NewFileServiceClient(cli)
}

func NewPermissionGroupServiceClient(ctx *bootstrap.Context, r registry.Discovery) permissionV1.PermissionGroupServiceClient {
	cli, err := rpc.CreateGrpcClient(ctx.Context(), r, serviceid.NewDiscoveryName(serviceid.CoreService), ctx.GetConfig())
	if err != nil {
		return nil
	}

	return permissionV1.NewPermissionGroupServiceClient(cli)
}

func NewPermissionServiceClient(ctx *bootstrap.Context, r registry.Discovery) permissionV1.PermissionServiceClient {
	cli, err := rpc.CreateGrpcClient(ctx.Context(), r, serviceid.NewDiscoveryName(serviceid.CoreService), ctx.GetConfig())
	if err != nil {
		return nil
	}

	return permissionV1.NewPermissionServiceClient(cli)
}

func NewApiServiceClient(ctx *bootstrap.Context, r registry.Discovery) permissionV1.ApiServiceClient {
	cli, err := rpc.CreateGrpcClient(ctx.Context(), r, serviceid.NewDiscoveryName(serviceid.CoreService), ctx.GetConfig())
	if err != nil {
		return nil
	}

	return permissionV1.NewApiServiceClient(cli)
}

func NewMenuServiceClient(ctx *bootstrap.Context, r registry.Discovery) permissionV1.MenuServiceClient {
	cli, err := rpc.CreateGrpcClient(ctx.Context(), r, serviceid.NewDiscoveryName(serviceid.CoreService), ctx.GetConfig())
	if err != nil {
		return nil
	}

	return permissionV1.NewMenuServiceClient(cli)
}

func NewPermissionAuditLogServiceClient(ctx *bootstrap.Context, r registry.Discovery) auditV1.PermissionAuditLogServiceClient {
	cli, err := rpc.CreateGrpcClient(ctx.Context(), r, serviceid.NewDiscoveryName(serviceid.CoreService), ctx.GetConfig())
	if err != nil {
		return nil
	}

	return auditV1.NewPermissionAuditLogServiceClient(cli)
}

func NewPolicyEvaluationLogServiceClient(ctx *bootstrap.Context, r registry.Discovery) permissionV1.PolicyEvaluationLogServiceClient {
	cli, err := rpc.CreateGrpcClient(ctx.Context(), r, serviceid.NewDiscoveryName(serviceid.CoreService), ctx.GetConfig())
	if err != nil {
		return nil
	}

	return permissionV1.NewPolicyEvaluationLogServiceClient(cli)
}

func NewApiAuditLogServiceClient(ctx *bootstrap.Context, r registry.Discovery) auditV1.ApiAuditLogServiceClient {
	cli, err := rpc.CreateGrpcClient(ctx.Context(), r, serviceid.NewDiscoveryName(serviceid.CoreService), ctx.GetConfig())
	if err != nil {
		return nil
	}

	return auditV1.NewApiAuditLogServiceClient(cli)
}

func NewDataAccessAuditLogServiceClient(ctx *bootstrap.Context, r registry.Discovery) auditV1.DataAccessAuditLogServiceClient {
	cli, err := rpc.CreateGrpcClient(ctx.Context(), r, serviceid.NewDiscoveryName(serviceid.CoreService), ctx.GetConfig())
	if err != nil {
		return nil
	}

	return auditV1.NewDataAccessAuditLogServiceClient(cli)
}

func NewLoginAuditLogServiceClient(ctx *bootstrap.Context, r registry.Discovery) auditV1.LoginAuditLogServiceClient {
	cli, err := rpc.CreateGrpcClient(ctx.Context(), r, serviceid.NewDiscoveryName(serviceid.CoreService), ctx.GetConfig())
	if err != nil {
		return nil
	}

	return auditV1.NewLoginAuditLogServiceClient(cli)
}

func NewOperationAuditLogServiceClient(ctx *bootstrap.Context, r registry.Discovery) auditV1.OperationAuditLogServiceClient {
	cli, err := rpc.CreateGrpcClient(ctx.Context(), r, serviceid.NewDiscoveryName(serviceid.CoreService), ctx.GetConfig())
	if err != nil {
		return nil
	}

	return auditV1.NewOperationAuditLogServiceClient(cli)
}

func NewDictEntryServiceClient(ctx *bootstrap.Context, r registry.Discovery) dictV1.DictEntryServiceClient {
	cli, err := rpc.CreateGrpcClient(ctx.Context(), r, serviceid.NewDiscoveryName(serviceid.CoreService), ctx.GetConfig())
	if err != nil {
		return nil
	}

	return dictV1.NewDictEntryServiceClient(cli)
}

func NewDictTypeServiceClient(ctx *bootstrap.Context, r registry.Discovery) dictV1.DictTypeServiceClient {
	cli, err := rpc.CreateGrpcClient(ctx.Context(), r, serviceid.NewDiscoveryName(serviceid.CoreService), ctx.GetConfig())
	if err != nil {
		return nil
	}

	return dictV1.NewDictTypeServiceClient(cli)
}

func NewLanguageServiceClient(ctx *bootstrap.Context, r registry.Discovery) dictV1.LanguageServiceClient {
	cli, err := rpc.CreateGrpcClient(ctx.Context(), r, serviceid.NewDiscoveryName(serviceid.CoreService), ctx.GetConfig())
	if err != nil {
		return nil
	}

	return dictV1.NewLanguageServiceClient(cli)
}

func NewTaskServiceClient(ctx *bootstrap.Context, r registry.Discovery) taskV1.TaskServiceClient {
	cli, err := rpc.CreateGrpcClient(ctx.Context(), r, serviceid.NewDiscoveryName(serviceid.CoreService), ctx.GetConfig())
	if err != nil {
		return nil
	}

	return taskV1.NewTaskServiceClient(cli)
}

func NewCommentServiceClient(ctx *bootstrap.Context, r registry.Discovery) commentV1.CommentServiceClient {
	cli, err := rpc.CreateGrpcClient(ctx.Context(), r, serviceid.NewDiscoveryName(serviceid.CoreService), ctx.GetConfig())
	if err != nil {
		return nil
	}

	return commentV1.NewCommentServiceClient(cli)
}

func NewInteractionAdminServiceClient(ctx *bootstrap.Context, r registry.Discovery) interactionV1.InteractionAdminServiceClient {
	cli, err := rpc.CreateGrpcClient(ctx.Context(), r, serviceid.NewDiscoveryName(serviceid.CoreService), ctx.GetConfig())
	if err != nil {
		return nil
	}

	return interactionV1.NewInteractionAdminServiceClient(cli)
}

func NewCategoryServiceClient(ctx *bootstrap.Context, r registry.Discovery) contentV1.CategoryServiceClient {
	cli, err := rpc.CreateGrpcClient(ctx.Context(), r, serviceid.NewDiscoveryName(serviceid.CoreService), ctx.GetConfig())
	if err != nil {
		return nil
	}

	return contentV1.NewCategoryServiceClient(cli)
}

func NewPageServiceClient(ctx *bootstrap.Context, r registry.Discovery) contentV1.PageServiceClient {
	cli, err := rpc.CreateGrpcClient(ctx.Context(), r, serviceid.NewDiscoveryName(serviceid.CoreService), ctx.GetConfig())
	if err != nil {
		return nil
	}

	return contentV1.NewPageServiceClient(cli)
}

func NewSectionServiceClient(ctx *bootstrap.Context, r registry.Discovery) contentV1.SectionServiceClient {
	cli, err := rpc.CreateGrpcClient(ctx.Context(), r, serviceid.NewDiscoveryName(serviceid.CoreService), ctx.GetConfig())
	if err != nil {
		return nil
	}

	return contentV1.NewSectionServiceClient(cli)
}

func NewPostServiceClient(ctx *bootstrap.Context, r registry.Discovery) contentV1.PostServiceClient {
	cli, err := rpc.CreateGrpcClient(ctx.Context(), r, serviceid.NewDiscoveryName(serviceid.CoreService), ctx.GetConfig())
	if err != nil {
		return nil
	}

	return contentV1.NewPostServiceClient(cli)
}

func NewContentModelServiceClient(ctx *bootstrap.Context, r registry.Discovery) contentV1.ContentModelServiceClient {
	cli, err := rpc.CreateGrpcClient(ctx.Context(), r, serviceid.NewDiscoveryName(serviceid.CoreService), ctx.GetConfig())
	if err != nil {
		return nil
	}

	return contentV1.NewContentModelServiceClient(cli)
}

func NewTagServiceClient(ctx *bootstrap.Context, r registry.Discovery) contentV1.TagServiceClient {
	cli, err := rpc.CreateGrpcClient(ctx.Context(), r, serviceid.NewDiscoveryName(serviceid.CoreService), ctx.GetConfig())
	if err != nil {
		return nil
	}

	return contentV1.NewTagServiceClient(cli)
}

func NewNavigationServiceClient(ctx *bootstrap.Context, r registry.Discovery) siteV1.NavigationServiceClient {
	cli, err := rpc.CreateGrpcClient(ctx.Context(), r, serviceid.NewDiscoveryName(serviceid.CoreService), ctx.GetConfig())
	if err != nil {
		return nil
	}

	return siteV1.NewNavigationServiceClient(cli)
}

func NewNavigationItemServiceClient(ctx *bootstrap.Context, r registry.Discovery) siteV1.NavigationItemServiceClient {
	cli, err := rpc.CreateGrpcClient(ctx.Context(), r, serviceid.NewDiscoveryName(serviceid.CoreService), ctx.GetConfig())
	if err != nil {
		return nil
	}

	return siteV1.NewNavigationItemServiceClient(cli)
}

func NewSiteSettingServiceClient(ctx *bootstrap.Context, r registry.Discovery) siteV1.SiteSettingServiceClient {
	cli, err := rpc.CreateGrpcClient(ctx.Context(), r, serviceid.NewDiscoveryName(serviceid.CoreService), ctx.GetConfig())
	if err != nil {
		return nil
	}

	return siteV1.NewSiteSettingServiceClient(cli)
}

func NewSiteServiceClient(ctx *bootstrap.Context, r registry.Discovery) siteV1.SiteServiceClient {
	cli, err := rpc.CreateGrpcClient(ctx.Context(), r, serviceid.NewDiscoveryName(serviceid.CoreService), ctx.GetConfig())
	if err != nil {
		return nil
	}

	return siteV1.NewSiteServiceClient(cli)
}

func NewMediaAssetServiceClient(ctx *bootstrap.Context, r registry.Discovery) mediaV1.MediaAssetServiceClient {
	cli, err := rpc.CreateGrpcClient(ctx.Context(), r, serviceid.NewDiscoveryName(serviceid.CoreService), ctx.GetConfig())
	if err != nil {
		return nil
	}

	return mediaV1.NewMediaAssetServiceClient(cli)
}
