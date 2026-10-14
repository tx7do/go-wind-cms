// Package mastodon Mastodon 发布适配器。
//
// 认证模型：个人访问令牌（Preferences → Development → New Application），
// Bearer 长期有效，无需刷新逻辑。
//
// 用到的 API：
//   - GET  /api/v1/accounts/verify_credentials  校验令牌
//   - POST /api/v1/statuses                     发帖（form 编码）
//   - GET  /api/v1/statuses/:id                 回拉互动计数
package mastodon

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode/utf8"

	"go-wind-cms/pkg/social"
)

const (
	// Platform 本适配器注册进 pkg/social 注册表的键（proto 枚举名）。
	Platform = "PLATFORM_MASTODON"

	maxStatusRunes  = 500 // Mastodon 单帖上限（多数实例 500；实例自定义上限以服务端报错为准）
	defaultTimeout  = 30 * time.Second
	maxResponseSize = 10 << 20 // 10MB
)

// Client Mastodon 客户端（并发安全）。
type Client struct {
	httpClient *http.Client
}

// Option 客户端选项。
type Option func(*Client)

// WithHTTPClient 覆盖底层 HTTP 客户端（测试注入 / 全局代理）。
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		if httpClient != nil {
			c.httpClient = httpClient
		}
	}
}

// NewClient 构造客户端（服务器地址来自凭据的 server 键，逐调用解析）。
func NewClient(opts ...Option) *Client {
	c := &Client{httpClient: &http.Client{Timeout: defaultTimeout}}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// init 注册进 pkg/social 注册表（key 与 ent/proto 枚举值一致）。
func init() {
	social.Register(Platform, NewClient())
}

// ── 数据结构 ──

type accountView struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type statusView struct {
	ID            string `json:"id"`
	URL           string `json:"url"`
	RepliesCount  int64  `json:"replies_count"`
	ReblogsCount  int64  `json:"reblogs_count"`
	FavouritesCount int64 `json:"favourites_count"`
}

// ── social.Provider 实现 ──

// ValidateCredentials 校验令牌：verify_credentials 通过即有效。
func (c *Client) ValidateCredentials(ctx context.Context, creds social.Credentials) error {
	var account accountView
	if err := c.do(ctx, creds, http.MethodGet, "/api/v1/accounts/verify_credentials", nil, &account); err != nil {
		return err
	}
	return nil
}

// Publish 发帖。externalID 为平台 status id（统计回拉定位键），externalURI 为帖子链接。
func (c *Client) Publish(ctx context.Context, creds social.Credentials, input *social.PublishInput) (*social.PublishResult, error) {
	if input == nil || strings.TrimSpace(input.Text) == "" {
		return nil, fmt.Errorf("mastodon: empty text")
	}
	if utf8.RuneCountInString(input.Text) > maxStatusRunes {
		return nil, social.ErrContentTooLong
	}

	form := url.Values{}
	form.Set("status", input.Text)
	form.Set("visibility", "public")

	var status statusView
	if err := c.do(ctx, creds, http.MethodPost, "/api/v1/statuses", form, &status); err != nil {
		return nil, err
	}
	if status.ID == "" {
		return nil, fmt.Errorf("mastodon: unexpected publish response")
	}

	return &social.PublishResult{
		ExternalID:  status.ID,
		ExternalURI: status.URL,
	}, nil
}

// FetchStats 回拉帖子互动计数。
func (c *Client) FetchStats(ctx context.Context, creds social.Credentials, externalID, _ string) (*social.Stats, error) {
	if externalID == "" {
		return nil, fmt.Errorf("mastodon: empty external id")
	}

	var status statusView
	if err := c.do(ctx, creds, http.MethodGet, "/api/v1/statuses/"+url.PathEscape(externalID), nil, &status); err != nil {
		return nil, err
	}

	return &social.Stats{
		Likes:   clampUint32(status.FavouritesCount),
		Replies: clampUint32(status.RepliesCount),
		Reposts: clampUint32(status.ReblogsCount),
	}, nil
}

// ── 内部实现 ──

// server 从凭据解析实例地址（https://mastodon.social 等自选实例）。
func server(creds social.Credentials) (string, error) {
	s := strings.TrimRight(strings.TrimSpace(creds.Get("server")), "/")
	if s == "" {
		return "", social.ErrInvalidCredentials
	}
	if !strings.HasPrefix(s, "http://") && !strings.HasPrefix(s, "https://") {
		s = "https://" + s
	}
	return s, nil
}

// do 执行一次 API 调用：Bearer 注入、form/JSON 编码、响应限长、错误归类。
func (c *Client) do(ctx context.Context, creds social.Credentials, method, path string, body any, out any) error {
	base, err := server(creds)
	if err != nil {
		return err
	}
	token := creds.Get("access_token")
	if token == "" {
		return social.ErrInvalidCredentials
	}

	var reqBody io.Reader
	contentType := ""
	if form, ok := body.(url.Values); ok {
		reqBody = strings.NewReader(form.Encode())
		contentType = "application/x-www-form-urlencoded"
	} else if body != nil {
		buf, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("mastodon: marshal request: %w", err)
		}
		reqBody = strings.NewReader(string(buf))
		contentType = "application/json"
	}

	req, err := http.NewRequestWithContext(ctx, method, base+path, reqBody)
	if err != nil {
		return fmt.Errorf("mastodon: build request: %w", err)
	}
	req.Header.Set("User-Agent", "go-wind-cms-social/1.0")
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("mastodon: request %s: %w", path, err)
	}
	defer func() { _ = resp.Body.Close() }()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseSize))
	if err != nil {
		return fmt.Errorf("mastodon: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return classifyError(resp.StatusCode, raw)
	}

	if out != nil && len(raw) != 0 {
		if err = json.Unmarshal(raw, out); err != nil {
			return fmt.Errorf("mastodon: decode response: %w", err)
		}
	}
	return nil
}

// classifyError 把非 2xx 响应归类为 pkg/social 的语义化错误。
func classifyError(status int, raw []byte) error {
	var apiErr struct {
		Error string `json:"error"`
	}
	_ = json.Unmarshal(raw, &apiErr)

	detail := strings.TrimSpace(apiErr.Error)
	msg := fmt.Sprintf("mastodon: http %d: %s", status, detail)

	switch {
	case status == http.StatusUnauthorized || status == http.StatusForbidden:
		return social.ErrInvalidCredentials
	case status == http.StatusTooManyRequests:
		return social.ErrRateLimited
	case status == http.StatusNotFound:
		return social.ErrNotFound
	case status == http.StatusUnprocessableEntity && strings.Contains(strings.ToLower(detail), "character"):
		return social.ErrContentTooLong
	default:
		return fmt.Errorf("social platform error: %s", msg)
	}
}

func clampUint32(v int64) uint32 {
	if v < 0 {
		return 0
	}
	if v > int64(^uint32(0)) {
		return ^uint32(0)
	}
	return uint32(v)
}
