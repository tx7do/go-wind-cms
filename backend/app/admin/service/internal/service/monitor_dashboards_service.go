package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	redisCacheV1 "go-wind-cms/api/gen/go/redis_cache/service/v1"
	serverMonitorV1 "go-wind-cms/api/gen/go/server_monitor/service/v1"
)

// ServerMonitorService 服务监控查询（BFF 薄代理：采集在 core 进程内）。
type ServerMonitorService struct {
	adminV1.ServerMonitorServiceHTTPServer

	log              *log.Helper
	serverMonitorCli serverMonitorV1.ServerMonitorServiceClient
}

func NewServerMonitorService(
	ctx *bootstrap.Context,
	serverMonitorCli serverMonitorV1.ServerMonitorServiceClient,
) *ServerMonitorService {
	l := log.NewHelper(log.With(ctx.GetLogger(), "module", "server-monitor/service/admin-service"))
	return &ServerMonitorService{
		log:              l,
		serverMonitorCli: serverMonitorCli,
	}
}

// Get 聚合一次服务监控快照（透传 core）。
func (s *ServerMonitorService) Get(ctx context.Context, req *serverMonitorV1.GetServerMonitorRequest) (*serverMonitorV1.ServerMonitorInfo, error) {
	return s.serverMonitorCli.Get(ctx, req)
}

// RedisCacheMonitorService Redis 缓存监控查询（BFF 薄代理）。
type RedisCacheMonitorService struct {
	adminV1.RedisCacheMonitorServiceHTTPServer

	log           *log.Helper
	redisCacheCli redisCacheV1.RedisCacheMonitorServiceClient
}

func NewRedisCacheMonitorService(
	ctx *bootstrap.Context,
	redisCacheCli redisCacheV1.RedisCacheMonitorServiceClient,
) *RedisCacheMonitorService {
	l := log.NewHelper(log.With(ctx.GetLogger(), "module", "redis-cache-monitor/service/admin-service"))
	return &RedisCacheMonitorService{
		log:           l,
		redisCacheCli: redisCacheCli,
	}
}

// Get 聚合一次 Redis 缓存监控快照（透传 core）。
func (s *RedisCacheMonitorService) Get(ctx context.Context, req *redisCacheV1.GetRedisCacheMonitorRequest) (*redisCacheV1.RedisCacheMonitorInfo, error) {
	return s.redisCacheCli.Get(ctx, req)
}
