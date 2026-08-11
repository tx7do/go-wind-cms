package data

import (
	"errors"

	"github.com/go-kratos/kratos/v2/registry"

	"github.com/tx7do/go-utils/password"

	"github.com/tx7do/kratos-bootstrap/bootstrap"
	redisClient "github.com/tx7do/kratos-bootstrap/cache/redis"
	"github.com/tx7do/kratos-bootstrap/database/elasticsearch"
	bRegistry "github.com/tx7do/kratos-bootstrap/registry"

	elasticsearchCrud "github.com/tx7do/go-crud/elasticsearch"

	"github.com/redis/go-redis/v9"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"

	"go-wind-cms/pkg/oss"
)

// NewRedisClient 创建Redis客户端
func NewRedisClient(ctx *bootstrap.Context) (*redis.Client, func(), error) {
	cfg := ctx.GetConfig()
	if cfg == nil {
		return nil, func() {}, nil
	}

	l := ctx.NewLoggerHelper("redis/data/core-service")

	cli := redisClient.NewClient(cfg.Data, l)

	return cli, func() {
		if err := cli.Close(); err != nil {
			l.Error(err)
		}
	}, nil
}

func NewElasticSearchClient(ctx *bootstrap.Context) (*elasticsearchCrud.Client, func(), error) {
	cfg := ctx.GetConfig()
	if cfg == nil {
		return nil, func() {}, nil
	}

	// AUD9-M7: 拒绝使用仓库内硬编码的默认 fallback 密码启动。
	// 生产环境必须通过 OPENSEARCH_PASSWORD 环境变量注入强密码；
	// 若落到 fallback "@Abcd#123456" 说明环境变量未设置，启动即失败，
	// 避免 ES admin 凭据以弱默认值暴露。
	if cfg.Data != nil && cfg.Data.Elasticsearch != nil {
		pw := cfg.Data.Elasticsearch.GetPassword()
		if pw == "" || pw == "@Abcd#123456" || pw == "dev_only_change_me" {
			return nil, func() {}, errors.New("elasticsearch password must be set via OPENSEARCH_PASSWORD environment variable; refusing to start with default/empty password")
		}
	}

	cli, err := elasticsearch.NewClient(ctx.GetLogger(), cfg)
	if err != nil {
		return nil, func() {}, err
	}

	return cli, func() {
	}, nil
}

// NewDiscovery 创建服务发现客户端
func NewDiscovery(ctx *bootstrap.Context) registry.Discovery {
	cfg := ctx.GetConfig()
	if cfg == nil {
		return nil
	}

	ret, err := bRegistry.NewDiscovery(cfg.Registry)
	if err != nil {
		return nil
	}
	return ret
}

func NewMinIoClient(ctx *bootstrap.Context) *oss.MinIOClient {
	return oss.NewMinIoClient(ctx.GetConfig(), ctx.GetLogger())
}

func NewPasswordCrypto() password.Crypto {
	crypto, err := password.CreateCrypto("bcrypt")
	if err != nil {
		panic(err)
	}
	return crypto
}
