package data

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// VCodeCache 验证码缓存（Redis）：一次性消费 + TTL。
//
// 与图形验证码（captcha 包自管）分开的原因：验证码按业务 scope 隔离
//（"reset_password"/"phone_bind"…），键形 vcode:{scope}:{identifier}，
// Verify 命中即删——单次有效由删除保证，不做重放窗口。
type VCodeCache struct {
	rdb *redis.Client
}

func NewVCodeCache(rdb *redis.Client) *VCodeCache {
	return &VCodeCache{rdb: rdb}
}

func vcodeKey(scope, identifier string) string {
	return fmt.Sprintf("vcode:%s:%s", scope, identifier)
}

// Save 写入验证码（覆盖旧值并重置 TTL）。
func (c *VCodeCache) Save(ctx context.Context, scope, identifier, code string, ttl time.Duration) error {
	if c.rdb == nil {
		return fmt.Errorf("redis client is nil")
	}
	return c.rdb.Set(ctx, vcodeKey(scope, identifier), code, ttl).Err()
}

// Verify 校验并消费（命中即删，单次有效；未命中/已过期/已消费均返回 false）。
func (c *VCodeCache) Verify(ctx context.Context, scope, identifier, code string) bool {
	if c.rdb == nil {
		return false
	}
	key := vcodeKey(scope, identifier)
	stored, err := c.rdb.Get(ctx, key).Result()
	if err != nil || stored == "" {
		return false
	}
	// 恒删除：错误尝试也消费掉，防爆破（重发验证码会覆盖新值）
	_ = c.rdb.Del(ctx, key).Err()
	return stored == code
}
