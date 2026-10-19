package service

import (
	"context"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/v9"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/timestamppb"

	redisCacheV1 "go-wind-cms/api/gen/go/redis_cache/service/v1"
)

// slowLogFetchLimit SLOWLOG GET 抓取的最近条目数上限。
const slowLogFetchLimit = 10

// RedisCacheMonitorService Redis 缓存监控查询（只读聚合：INFO/DBSIZE/SLOWLOG）。
// 数据来源为共享 Redis 的只读命令，不做任何写操作。
type RedisCacheMonitorService struct {
	redisCacheV1.UnimplementedRedisCacheMonitorServiceServer

	rdb *redis.Client
	log *log.Helper
}

func NewRedisCacheMonitorService(
	ctx *bootstrap.Context,
	rdb *redis.Client,
) *RedisCacheMonitorService {
	return &RedisCacheMonitorService{
		rdb: rdb,
		log: ctx.NewLoggerHelper("redis-cache-monitor/service/core-service"),
	}
}

// Get 聚合返回 Redis 的 INFO（按 section 拆分）、DBSIZE、SLOWLOG 三类只读运维指标。
// 任一命令失败时仅记录告警并将对应字段置空，不阻断其余字段（fail-soft）。
func (s *RedisCacheMonitorService) Get(ctx context.Context, _ *redisCacheV1.GetRedisCacheMonitorRequest) (*redisCacheV1.RedisCacheMonitorInfo, error) {
	if err := requireNotificationPlatformContext(ctx, "redis cache monitor"); err != nil {
		return nil, err
	}

	info := &redisCacheV1.RedisCacheMonitorInfo{}

	if s.rdb == nil {
		// Redis 不可用时返回空视图，而非错误。
		return info, nil
	}

	if raw, err := s.rdb.Info(ctx).Result(); err != nil {
		s.log.Errorf("redis INFO failed: %s", err.Error())
	} else {
		info.Sections = parseInfoSections(raw)
	}

	if n, err := s.rdb.DBSize(ctx).Result(); err != nil {
		s.log.Errorf("redis DBSIZE failed: %s", err.Error())
	} else if n >= 0 {
		info.DbSize = uint64(n)
	}

	if entries, err := s.rdb.SlowLogGet(ctx, slowLogFetchLimit).Result(); err != nil {
		s.log.Errorf("redis SLOWLOG GET failed: %s", err.Error())
	} else {
		info.Slowlog = mapSlowLogEntries(entries)
	}

	return info, nil
}

// mapSlowLogEntries 将 redis.SlowLog 列表映射为 proto 消息列表。
func mapSlowLogEntries(in []redis.SlowLog) []*redisCacheV1.SlowLogEntry {
	if len(in) == 0 {
		return nil
	}
	out := make([]*redisCacheV1.SlowLogEntry, 0, len(in))
	for i := range in {
		slow := &in[i]
		out = append(out, &redisCacheV1.SlowLogEntry{
			Id:           slow.ID,
			CreatedAt:    timestamppb.New(slow.Time),
			DurationUsec: slow.Duration.Microseconds(),
			Args:         sanitizeSlowLogArgs(slow.Args),
			ClientAddr:   slow.ClientAddr,
			ClientName:   strings.ToValidUTF8(slow.ClientName, "\uFFFD"),
		})
	}
	return out
}

// sanitizeSlowLogArgs 将 slowlog 命令参数净化为合法 UTF-8。
// Redis 命令参数二进制安全，而 proto string 强制 UTF-8——未净化直接透传会让
// 响应在序列化阶段整体失败，表现为与本服务 fail-soft 设计相悖的 500。
func sanitizeSlowLogArgs(args []string) []string {
	if len(args) == 0 {
		return nil
	}
	out := make([]string, 0, len(args))
	for _, a := range args {
		out = append(out, strings.ToValidUTF8(a, "\uFFFD"))
	}
	return out
}

// parseInfoSections 解析 Redis INFO 原始串为 section/entry 树。
// INFO 输出格式：以 `# <SectionName>` 行标记新 section，其后 `key:value` 行为
// 该 section 的条目，空行与无法识别的行被忽略。Keyspace 的 `db0:keys=...`
// 同样按 key/value 落入此结构，前端按通用 kv 表渲染。
func parseInfoSections(raw string) []*redisCacheV1.InfoSection {
	var sections []*redisCacheV1.InfoSection
	var current *redisCacheV1.InfoSection

	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimRight(line, "\r")
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "# ") {
			name := strings.TrimSpace(line[2:])
			if name == "" {
				// 畸形 section 头：停止向其挂条目，避免误归属。
				current = nil
				continue
			}
			current = &redisCacheV1.InfoSection{Name: name}
			sections = append(sections, current)
			continue
		}

		idx := strings.IndexByte(line, ':')
		if idx < 0 {
			continue
		}
		if current == nil {
			// 无归属 section 的孤立条目，忽略。
			continue
		}
		current.Entries = append(current.Entries, &redisCacheV1.InfoEntry{
			Key:   line[:idx],
			Value: line[idx+1:],
		})
	}

	return sections
}
