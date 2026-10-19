package data

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/v9"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
)

// 访问凭证（AK/SK）令牌交换失败计数与锁定常量。
// 与上游登录限流器同语义（5 次失败锁 15 分钟），但独立成键——CMS 登录链路
// 未接限流器，本计数器只服务交换端点的防爆破。
const (
	// accessKeyFailThreshold 触发锁定的失败次数阈值。
	accessKeyFailThreshold = 5
	// accessKeyLockoutDuration 触发锁定后的窗口时长。
	accessKeyLockoutDuration = 15 * time.Minute
)

// 失败计数键：按 IP 与 AK 双维度（同上游 IP+用户名双维度的防爆破思路）。
const (
	accessKeyFailKeyByIPFmt = "gwc:akx:fail:ip:%s"
	accessKeyFailKeyByAKFmt = "gwc:akx:fail:ak:%s"
)

// incrIfNotLockedScript 原子「判定是否已锁定 → 自增失败计数 → 设置 TTL」Lua 脚本。
// 语义：若当前计数已 >= 阈值，则判定为锁定状态、不再自增；否则自增并在首次失败时设置 TTL。
// 返回值：{已锁定 flag(0/1), 当前失败次数}
var incrIfNotLockedScript = redis.NewScript(`
	local key = KEYS[1]
	local threshold = tonumber(ARGV[1])
	local ttl = tonumber(ARGV[2])

	local current = tonumber(redis.call('GET', key) or "0")
	if current >= threshold then
		return {1, current}
	end

	current = redis.call('INCR', key)
	if current == 1 then
		redis.call('EXPIRE', key, ttl)
	end
	return {0, current}
`)

// AccessKeyRateLimiter 基于 Redis 的令牌交换失败计数器，按 IP 与 AK 双维度限流。
// Redis 不可用时 fail-open（限流属防御性增强，不应让 Redis 故障阻断合法机器调用）。
type AccessKeyRateLimiter struct {
	log *log.Helper
	rdb *redis.Client
}

func NewAccessKeyRateLimiter(ctx *bootstrap.Context, rdb *redis.Client) *AccessKeyRateLimiter {
	return &AccessKeyRateLimiter{
		rdb: rdb,
		log: ctx.NewLoggerHelper("access-key/rate-limiter/core-service"),
	}
}

// CheckAndIncr 在交换失败时自增失败计数，并返回是否已触发锁定。
// 任一维度超阈值即视为锁定；retryAfter 当前固定为锁定窗口。
func (l *AccessKeyRateLimiter) CheckAndIncr(ctx context.Context, ip, accessKey string) (locked bool, attempts int, retryAfter time.Duration, err error) {
	if l == nil || l.rdb == nil {
		return false, 0, 0, nil
	}

	keys := l.failKeys(ip, accessKey)
	threshold := int64(accessKeyFailThreshold)
	ttl := int64(accessKeyLockoutDuration.Seconds())

	maxAttempts := 0
	for _, key := range keys {
		res, rerr := incrIfNotLockedScript.Run(ctx, l.rdb, []string{key}, threshold, ttl).Result()
		if rerr != nil {
			l.log.Errorf("incr access key fail counter failed (key=%s): %s", key, rerr.Error())
			continue
		}
		// res 形如 [1 "5"] 或 [0 "3"]
		vals, ok := res.([]interface{})
		if !ok || len(vals) < 2 {
			continue
		}
		isLocked, _ := vals[0].(int64)
		count, _ := vals[1].(int64)
		if int(count) > maxAttempts {
			maxAttempts = int(count)
		}
		if isLocked == 1 {
			l.log.Warnf("access key exchange locked: key=%s, attempts=%d", key, count)
			return true, int(count), accessKeyLockoutDuration, nil
		}
	}
	return false, maxAttempts, 0, nil
}

// IsLocked 查询当前 IP/AK 是否处于锁定状态（不自增计数）。交换前置拦截用。
func (l *AccessKeyRateLimiter) IsLocked(ctx context.Context, ip, accessKey string) (bool, error) {
	if l == nil || l.rdb == nil {
		return false, nil
	}
	threshold := int64(accessKeyFailThreshold)
	for _, key := range l.failKeys(ip, accessKey) {
		cnt, rerr := l.rdb.Get(ctx, key).Int64()
		if rerr != nil && rerr != redis.Nil {
			l.log.Errorf("get access key fail counter failed (key=%s): %s", key, rerr.Error())
			continue
		}
		if cnt >= threshold {
			return true, nil
		}
	}
	return false, nil
}

// Reset 交换成功后清零失败计数。
func (l *AccessKeyRateLimiter) Reset(ctx context.Context, ip, accessKey string) {
	if l == nil || l.rdb == nil {
		return
	}
	for _, key := range l.failKeys(ip, accessKey) {
		if err := l.rdb.Del(ctx, key).Err(); err != nil {
			l.log.Errorf("reset access key fail counter failed (key=%s): %s", key, err.Error())
		}
	}
}

// failKeys 根据 IP 与 AK 生成两个维度的失败计数键；空值维度跳过。
func (l *AccessKeyRateLimiter) failKeys(ip, accessKey string) []string {
	keys := make([]string, 0, 2)
	if ip != "" {
		keys = append(keys, fmt.Sprintf(accessKeyFailKeyByIPFmt, ip))
	}
	if accessKey != "" {
		keys = append(keys, fmt.Sprintf(accessKeyFailKeyByAKFmt, accessKey))
	}
	return keys
}
