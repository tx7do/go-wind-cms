// Package nostr Nostr（开放协议）发布适配器。
//
// 认证模型：无平台、无审核——账户即密钥。凭据为 NIP-01 私钥（64 位 hex 或
// nsec bech32），发布 = 构造 kind=1 事件并 Schnorr 签名后向一组 relay 广播，
// 任一 relay 接受即成功。
//
// 凭据：
//   - private_key: 私钥（hex 或 nsec）
//   - relays:      relay 地址列表（逗号分隔，ws:// 或 wss://）
//
// 统计：互动（NIP-25 reaction）需跨 relay 查询聚合，第一版不回拉
// （FetchStats 返回 ErrStatsUnsupported）。
package nostr

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	stdnostr "github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"

	"go-wind-cms/pkg/social"
)

const (
	// Platform 本适配器注册进 pkg/social 注册表的键（proto 枚举名）。
	Platform = "PLATFORM_NOSTR"

	publishTimeout = 15 * time.Second
)

// publishOutcome 单个 relay 的投递结果（url → error，nil 表示被 relay 接受）。
type publishOutcome struct {
	URL string
	Err error
}

// publisher relay 广播的实现缝：默认走 go-nostr SimplePool，测试注入假实现。
type publisher func(ctx context.Context, urls []string, evt stdnostr.Event) []publishOutcome

// Client Nostr 客户端（并发安全）。
type Client struct {
	pool    *stdnostr.SimplePool
	publish publisher
}

// Option 客户端选项。
type Option func(*Client)

// WithPublisher 覆盖 relay 广播实现（测试注入）。
func WithPublisher(p publisher) Option {
	return func(c *Client) {
		if p != nil {
			c.publish = p
		}
	}
}

// NewClient 构造客户端。
func NewClient(opts ...Option) *Client {
	c := &Client{}
	for _, opt := range opts {
		opt(c)
	}
	if c.publish == nil {
		c.pool = stdnostr.NewSimplePool(context.Background())
		c.publish = c.publishViaPool
	}
	return c
}

// init 注册进 pkg/social 注册表（key 与 ent/proto 枚举值一致）。
func init() {
	social.Register(Platform, NewClient())
}

// ── social.Provider 实现 ──

// ValidateCredentials 校验私钥：能解码并推导出公钥即有效。
func (c *Client) ValidateCredentials(_ context.Context, creds social.Credentials) error {
	_, err := decodePrivateKey(creds.Get("private_key"))
	return err
}

// Publish 构造并签名 kind=1 事件，向凭据内的全部 relay 广播；
// 任一 relay 接受即成功，全部失败才报错（附各 relay 结果）。
// external_id 为事件 ID（hex），external_uri 留空（无规范全局链接）。
func (c *Client) Publish(_ context.Context, creds social.Credentials, input *social.PublishInput) (*social.PublishResult, error) {
	if input == nil || strings.TrimSpace(input.Text) == "" {
		return nil, fmt.Errorf("nostr: empty text")
	}

	sk, err := decodePrivateKey(creds.Get("private_key"))
	if err != nil {
		return nil, err
	}
	relays := parseRelays(creds.Get("relays"))
	if len(relays) == 0 {
		return nil, fmt.Errorf("nostr: no relay configured")
	}

	pk, err := stdnostr.GetPublicKey(sk)
	if err != nil {
		return nil, fmt.Errorf("nostr: derive pubkey: %w", err)
	}

	evt := stdnostr.Event{
		PubKey:    pk,
		CreatedAt: stdnostr.Now(),
		Kind:      1,
		Content:   input.Text,
	}
	if err = evt.Sign(sk); err != nil {
		return nil, fmt.Errorf("nostr: sign event: %w", err)
	}
	evt.ID = evt.GetID()

	ctx, cancel := context.WithTimeout(context.Background(), publishTimeout)
	defer cancel()

	outcomes := c.publish(ctx, relays, evt)
	var failures []string
	for _, oc := range outcomes {
		if oc.Err != nil {
			failures = append(failures, fmt.Sprintf("%s: %v", oc.URL, oc.Err))
		}
	}
	if len(failures) == len(relays) {
		return nil, fmt.Errorf("social platform error: nostr all relays rejected: %s",
			strings.Join(failures, "; "))
	}

	return &social.PublishResult{
		ExternalID: evt.ID,
	}, nil
}

// FetchStats 互动聚合需跨 relay 查询，第一版不回拉。
func (c *Client) FetchStats(context.Context, social.Credentials, string, string) (*social.Stats, error) {
	return nil, social.ErrStatsUnsupported
}

// ── 内部实现 ──

// publishViaPool 默认广播实现：SimplePool 并发投递，收集全部 relay 结果。
func (c *Client) publishViaPool(ctx context.Context, urls []string, evt stdnostr.Event) []publishOutcome {
	ch := c.pool.PublishMany(ctx, urls, evt)

	outcomes := make([]publishOutcome, 0, len(urls))
	for range urls {
		select {
		case <-ctx.Done():
			return append(outcomes, publishOutcome{Err: fmt.Errorf("publish timeout")})
		case r := <-ch:
			outcomes = append(outcomes, publishOutcome{URL: r.RelayURL, Err: r.Error})
		}
	}
	return outcomes
}

// decodePrivateKey 私钥解析：64 位 hex 或 nsec bech32（统一转 hex）。
func decodePrivateKey(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", social.ErrInvalidCredentials
	}

	if strings.HasPrefix(raw, "nsec") {
		prefix, value, err := nip19.Decode(raw)
		if err != nil || prefix != "nsec" {
			return "", social.ErrInvalidCredentials
		}
		hexKey, ok := value.(string)
		if !ok || hexKey == "" {
			return "", social.ErrInvalidCredentials
		}
		raw = hexKey
	}

	b, err := hex.DecodeString(raw)
	if err != nil || len(b) != 32 {
		return "", social.ErrInvalidCredentials
	}
	return raw, nil
}

// parseRelays 解析 relay 列表（逗号/空白分隔，保留 ws:// 或 wss:// 原样）。
func parseRelays(raw string) []string {
	fields := strings.FieldsFunc(raw, func(r rune) bool {
		return r == ',' || r == ' ' || r == '\n' || r == '\t' || r == ';'
	})
	urls := make([]string, 0, len(fields))
	for _, f := range fields {
		f = strings.TrimSpace(f)
		if f == "" {
			continue
		}
		if !strings.HasPrefix(f, "ws://") && !strings.HasPrefix(f, "wss://") {
			f = "wss://" + f
		}
		urls = append(urls, f)
	}
	return urls
}
