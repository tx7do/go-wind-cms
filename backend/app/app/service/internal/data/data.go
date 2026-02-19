package data

import (
	"github.com/redis/go-redis/v9"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	authnEngine "github.com/tx7do/kratos-authn/engine"
	"github.com/tx7do/kratos-authn/engine/jwt"

	authzEngine "github.com/tx7do/kratos-authz/engine"
	"github.com/tx7do/kratos-authz/engine/noop"

	"github.com/go-kratos/kratos/v2/registry"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	redisClient "github.com/tx7do/kratos-bootstrap/cache/redis"
	bRegistry "github.com/tx7do/kratos-bootstrap/registry"
	"github.com/tx7do/kratos-bootstrap/rpc"

	authenticationV1 "go-wind-cms/api/gen/go/authentication/service/v1"
	commentV1 "go-wind-cms/api/gen/go/comment/service/v1"
	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
	mediaV1 "go-wind-cms/api/gen/go/media/service/v1"
	siteV1 "go-wind-cms/api/gen/go/site/service/v1"

	"go-wind-cms/pkg/serviceid"
)

func NewClientType() authenticationV1.ClientType {
	return authenticationV1.ClientType_app
}

// NewRedisClient 创建Redis客户端
func NewRedisClient(ctx *bootstrap.Context) (*redis.Client, func(), error) {
	cfg := ctx.GetConfig()
	if cfg == nil {
		return nil, func() {}, nil
	}

	l := ctx.NewLoggerHelper("redis/data/app-service")

	cli := redisClient.NewClient(cfg.Data, l)

	return cli, func() {
		if err := cli.Close(); err != nil {
			l.Error(err)
		}
	}, nil
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

// NewAuthenticator 创建认证器
func NewAuthenticator(cfg *conf.Bootstrap) authnEngine.Authenticator {
	authenticator, _ := jwt.NewAuthenticator(
		jwt.WithKey([]byte(cfg.Server.Rest.Middleware.Auth.Key)),
		jwt.WithSigningMethod(cfg.Server.Rest.Middleware.Auth.Method),
	)
	return authenticator
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

func NewCommentServiceClient(ctx *bootstrap.Context, r registry.Discovery) commentV1.CommentServiceClient {
	cli, err := rpc.CreateGrpcClient(ctx.Context(), r, serviceid.NewDiscoveryName(serviceid.CoreService), ctx.GetConfig())
	if err != nil {
		return nil
	}

	return commentV1.NewCommentServiceClient(cli)
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

func NewPostServiceClient(ctx *bootstrap.Context, r registry.Discovery) contentV1.PostServiceClient {
	cli, err := rpc.CreateGrpcClient(ctx.Context(), r, serviceid.NewDiscoveryName(serviceid.CoreService), ctx.GetConfig())
	if err != nil {
		return nil
	}

	return contentV1.NewPostServiceClient(cli)
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

func NewSiteSettingServiceClient(ctx *bootstrap.Context, r registry.Discovery) siteV1.SiteSettingServiceClient {
	cli, err := rpc.CreateGrpcClient(ctx.Context(), r, serviceid.NewDiscoveryName(serviceid.CoreService), ctx.GetConfig())
	if err != nil {
		return nil
	}

	return siteV1.NewSiteSettingServiceClient(cli)
}

func NewMediaAssetServiceClient(ctx *bootstrap.Context, r registry.Discovery) mediaV1.MediaAssetServiceClient {
	cli, err := rpc.CreateGrpcClient(ctx.Context(), r, serviceid.NewDiscoveryName(serviceid.CoreService), ctx.GetConfig())
	if err != nil {
		return nil
	}

	return mediaV1.NewMediaAssetServiceClient(cli)
}
