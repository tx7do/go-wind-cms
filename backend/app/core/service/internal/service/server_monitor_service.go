package service

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/timestamppb"

	entCrud "github.com/tx7do/go-crud/entgo"

	serverMonitorV1 "go-wind-cms/api/gen/go/server_monitor/service/v1"

	entgo "go-wind-cms/app/core/service/internal/data/ent"
)

// startedAt 进程启动时间（uptime 计算锚；包级初始化即进程启动近似值）。
var startedAt = time.Now()

// ServerMonitorService 服务监控查询（只读聚合：Go 运行时 / 数据库连接池 / 主机信息）。
// 采集全部在 core 进程内完成，与监控告警的阈值指标同源。
type ServerMonitorService struct {
	serverMonitorV1.UnimplementedServerMonitorServiceServer

	entClient *entCrud.EntClient[*entgo.Client]
	log       *log.Helper
}

func NewServerMonitorService(
	ctx *bootstrap.Context,
	entClient *entCrud.EntClient[*entgo.Client],
) *ServerMonitorService {
	return &ServerMonitorService{
		entClient: entClient,
		log:       ctx.NewLoggerHelper("server-monitor/service/core-service"),
	}
}

// Get 聚合一次服务监控快照。任一来源失败不阻断其余视图（fail-soft）。
func (s *ServerMonitorService) Get(ctx context.Context, _ *serverMonitorV1.GetServerMonitorRequest) (*serverMonitorV1.ServerMonitorInfo, error) {
	if err := requireNotificationPlatformContext(ctx, "server monitor"); err != nil {
		return nil, err
	}

	info := &serverMonitorV1.ServerMonitorInfo{
		Go:          s.goRuntimeInfo(),
		Host:        s.hostInfo(),
		CollectedAt: timestamppb.Now(),
	}
	if db, err := s.databaseInfo(ctx); err != nil {
		s.log.Errorf("collect database info failed: %s", err.Error())
	} else {
		info.Database = db
	}

	return info, nil
}

func (s *ServerMonitorService) goRuntimeInfo() *serverMonitorV1.GoRuntimeInfo {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	return &serverMonitorV1.GoRuntimeInfo{
		Version:       trans.Ptr(runtime.Version()),
		NumGoroutine:  trans.Ptr(uint32(runtime.NumGoroutine())),
		MemAllocBytes: trans.Ptr(ms.Alloc),
		MemSysBytes:   trans.Ptr(ms.Sys),
		GcCycles:      trans.Ptr(uint32(ms.NumGC)),
		UptimeSeconds: trans.Ptr(uint64(time.Since(startedAt).Seconds())),
		StartedAt:     timestamppb.New(startedAt),
	}
}

// databaseInfo 数据库连接池统计与连通性（经 ent 底层 *sql.DB）。
func (s *ServerMonitorService) databaseInfo(ctx context.Context) (*serverMonitorV1.DatabaseInfo, error) {
	if s.entClient == nil {
		return nil, fmt.Errorf("ent client is nil")
	}
	db := s.entClient.DB()
	if db == nil {
		return nil, fmt.Errorf("sql db is nil")
	}

	info := &serverMonitorV1.DatabaseInfo{
		Driver: trans.Ptr("postgres"),
	}
	stats := db.Stats()
	info.MaxOpenConnections = trans.Ptr(uint32(stats.MaxOpenConnections))
	info.OpenConnections = trans.Ptr(uint32(stats.OpenConnections))
	info.InUseConnections = trans.Ptr(uint32(stats.InUse))
	info.IdleConnections = trans.Ptr(uint32(stats.Idle))

	if err := db.PingContext(ctx); err != nil {
		info.PingOk = transPtr(false)
		info.PingError = transPtr(err.Error())
	} else {
		info.PingOk = transPtr(true)
	}
	return info, nil
}

func (s *ServerMonitorService) hostInfo() *serverMonitorV1.HostInfo {
	hostname, _ := os.Hostname()
	return &serverMonitorV1.HostInfo{
		Os:       trans.Ptr(runtime.GOOS),
		Arch:     trans.Ptr(runtime.GOARCH),
		NumCpu:   trans.Ptr(uint32(runtime.NumCPU())),
		Hostname: trans.Ptr(hostname),
	}
}

func transPtr[T any](v T) *T { return &v }
