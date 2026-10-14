// Package social 社交平台适配层：定义统一的 Provider 接口与凭据载体，
// 各平台的 HTTP 集成放在子包（如 bluesky）并在其 init() 中注册。
//
// 凭据只在 core service 层解密后以 Credentials 传入，适配器不接触存储；
// 与 pkg/ai 的"密钥运行时注入"同款安全模型。
package social

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
)

// 语义化哨兵错误：core service 据此映射 proto 错误码，适配器负责返回这些值。
var (
	// ErrInvalidCredentials 平台认证失败（凭据错误/过期）。
	ErrInvalidCredentials = errors.New("social: invalid credentials")
	// ErrRateLimited 平台限流。
	ErrRateLimited = errors.New("social: rate limited")
	// ErrContentTooLong 文案超出平台长度限制。
	ErrContentTooLong = errors.New("social: content too long")
	// ErrNotFound 平台侧内容不存在（已删除或不可见）。
	ErrNotFound = errors.New("social: content not found")
	// ErrStatsUnsupported 平台不提供互动统计能力（如 Telegram Bot API 无浏览数、
	// 自建 WordPress 无平台侧统计）。统计回拉链路据此静默跳过而非报错。
	ErrStatsUnsupported = errors.New("social: stats not supported by platform")
)

// Credentials 平台凭据的解密后载体。各平台自行约定键名：
//   - bluesky: handle / app_password / pds_url（可选，缺省 https://bsky.social）
//   - mastodon（预留）: server / access_token
//   - devto（预留）: api_key
//
// 存储形态即此结构的 JSON，经 AES-GCM 加密落库（enc: 前缀）。
type Credentials map[string]string

// Get 取凭据字段，缺失返回空串。
func (c Credentials) Get(key string) string {
	if c == nil {
		return ""
	}
	return c[key]
}

// ParseCredentials 解密后的凭据 JSON → Credentials。
func ParseCredentials(raw string) (Credentials, error) {
	if raw == "" {
		return nil, fmt.Errorf("social: empty credentials")
	}
	creds := Credentials{}
	if err := json.Unmarshal([]byte(raw), &creds); err != nil {
		return nil, fmt.Errorf("social: parse credentials: %w", err)
	}
	if len(creds) == 0 {
		return nil, fmt.Errorf("social: empty credentials")
	}
	return creds, nil
}

// MarshalCredentials Credentials → 存储 JSON（由调用方加密后落库）。
func MarshalCredentials(creds Credentials) (string, error) {
	b, err := json.Marshal(creds)
	if err != nil {
		return "", fmt.Errorf("social: marshal credentials: %w", err)
	}
	return string(b), nil
}

// MediaAttachment 一张随内容发布的图片（已从媒体库解析出的字节）。
type MediaAttachment struct {
	Filename    string
	ContentType string
	Data        []byte
}

// PublishInput 一次发布的内容输入。
type PublishInput struct {
	// Text 发布文案（来自分发记录的平台覆盖文案）。
	Text string
	// Title 平台文章标题（DEV.to 等以"文章"为单位的平台使用；可为空）。
	Title string
	// Media 配图（仅支持图片的平台消费；按平台自身的张数上限取舍）。
	Media []MediaAttachment
}

// PublishResult 发布成功后的平台侧定位信息。
type PublishResult struct {
	// ExternalID 平台侧帖子标识（统计回拉的定位键；格式由平台自定）。
	ExternalID string
	// ExternalURI 平台侧帖子的人类可读链接（展示用）。
	ExternalURI string
}

// Stats 平台互动统计快照。
type Stats struct {
	Likes   uint32
	Replies uint32
	Reposts uint32
}

// Provider 社交平台适配器接口。凭据在调用时传入（解密后的值不越过本层）。
type Provider interface {
	// ValidateCredentials 校验凭据可用性（账号保存时同步验证）。
	ValidateCredentials(ctx context.Context, creds Credentials) error
	// Publish 发布内容。
	Publish(ctx context.Context, creds Credentials, input *PublishInput) (*PublishResult, error)
	// FetchStats 回拉指定帖子的互动统计。
	FetchStats(ctx context.Context, creds Credentials, externalID, externalURI string) (*Stats, error)
}

// ── 平台注册表 ──
// key 为 proto 枚举名字符串（如 "PLATFORM_BLUESKY"），与 ent 枚举值一致。

var (
	registryMu sync.RWMutex
	registry   = map[string]Provider{}
)

// Register 注册平台适配器。各平台子包在 init() 中调用。
func Register(platform string, p Provider) {
	registryMu.Lock()
	defer registryMu.Unlock()
	registry[platform] = p
}

// GetProvider 取平台适配器；未注册返回错误。
func GetProvider(platform string) (Provider, error) {
	registryMu.RLock()
	defer registryMu.RUnlock()
	p, ok := registry[platform]
	if !ok {
		return nil, fmt.Errorf("social provider %q not registered", platform)
	}
	return p, nil
}
