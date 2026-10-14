// Package telegram Telegram 频道发布适配器。
//
// 认证模型：@BotFather 创建 bot 即得 token（零审核）；bot 需被设为目标频道的
// 管理员才能以频道身份发帖。
//
// 用到的 Bot API：
//   - GET /bot{token}/getMe      校验 token
//   - POST /bot{token}/sendMessage  发帖
//
// 统计：Bot API 不暴露频道帖浏览数，FetchStats 返回 ErrStatsUnsupported。
package telegram

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
	Platform = "PLATFORM_TELEGRAM"

	maxMessageRunes = 4096 // Telegram 单条消息上限
	apiBase         = "https://api.telegram.org"
	defaultTimeout  = 30 * time.Second
	maxResponseSize = 10 << 20 // 10MB
)

// Client Telegram 客户端（并发安全）。
type Client struct {
	httpClient *http.Client
	endpoint   string
}

// Option 客户端选项。
type Option func(*Client)

// WithEndpoint 覆盖 API 地址（测试桩）。
func WithEndpoint(endpoint string) Option {
	return func(c *Client) {
		if endpoint != "" {
			c.endpoint = strings.TrimRight(endpoint, "/")
		}
	}
}

// WithHTTPClient 覆盖底层 HTTP 客户端。
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		if httpClient != nil {
			c.httpClient = httpClient
		}
	}
}

// NewClient 构造客户端。
func NewClient(opts ...Option) *Client {
	c := &Client{endpoint: apiBase, httpClient: &http.Client{Timeout: defaultTimeout}}
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

type apiResponse struct {
	OK          bool            `json:"ok"`
	ErrorCode   int             `json:"error_code"`
	Description string          `json:"description"`
	Result      json.RawMessage `json:"result"`
}

type userView struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

type messageView struct {
	MessageID int64 `json:"message_id"`
	Chat      struct {
		ID       int64  `json:"id"`
		Username string `json:"username"`
	} `json:"chat"`
}

// ── social.Provider 实现 ──

// ValidateCredentials 校验 token：getMe 通过即有效。
func (c *Client) ValidateCredentials(ctx context.Context, creds social.Credentials) error {
	var user userView
	return c.do(ctx, creds, http.MethodGet, "/getMe", nil, &user)
}

// Publish 发帖到 chat_id 指定的频道/群。公开频道（chat_id 为 @name）可构造
// t.me 链接；私有频道无法构造，externalURI 留空。
func (c *Client) Publish(ctx context.Context, creds social.Credentials, input *social.PublishInput) (*social.PublishResult, error) {
	if input == nil || strings.TrimSpace(input.Text) == "" {
		return nil, fmt.Errorf("telegram: empty text")
	}
	if utf8.RuneCountInString(input.Text) > maxMessageRunes {
		return nil, social.ErrContentTooLong
	}

	chatID := creds.Get("chat_id")
	if chatID == "" {
		return nil, social.ErrInvalidCredentials
	}

	form := url.Values{}
	form.Set("chat_id", chatID)
	form.Set("text", input.Text)

	var msg messageView
	if err := c.do(ctx, creds, http.MethodPost, "/sendMessage", form, &msg); err != nil {
		return nil, err
	}
	if msg.MessageID == 0 {
		return nil, fmt.Errorf("telegram: unexpected publish response")
	}

	uri := ""
	if strings.HasPrefix(chatID, "@") {
		uri = fmt.Sprintf("https://t.me/%s/%d", strings.TrimPrefix(chatID, "@"), msg.MessageID)
	}

	return &social.PublishResult{
		ExternalID:  fmt.Sprintf("%d", msg.MessageID),
		ExternalURI: uri,
	}, nil
}

// FetchStats Bot API 不提供频道帖浏览/互动计数。
func (c *Client) FetchStats(context.Context, social.Credentials, string, string) (*social.Stats, error) {
	return nil, social.ErrStatsUnsupported
}

// ── 内部实现 ──

// do 执行一次 Bot API 调用：token 注入路径、form/JSON 编码、响应限长、错误归类。
func (c *Client) do(ctx context.Context, creds social.Credentials, method, path string, body any, out any) error {
	token := creds.Get("bot_token")
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
			return fmt.Errorf("telegram: marshal request: %w", err)
		}
		reqBody = strings.NewReader(string(buf))
		contentType = "application/json"
	}

	req, err := http.NewRequestWithContext(ctx, method, c.endpoint+"/bot"+url.PathEscape(token)+path, reqBody)
	if err != nil {
		return fmt.Errorf("telegram: build request: %w", err)
	}
	req.Header.Set("User-Agent", "go-wind-cms-social/1.0")
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("telegram: request %s: %w", path, err)
	}
	defer func() { _ = resp.Body.Close() }()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseSize))
	if err != nil {
		return fmt.Errorf("telegram: read response: %w", err)
	}

	var api apiResponse
	_ = json.Unmarshal(raw, &api)
	if !api.OK {
		return classifyError(api.ErrorCode, api.Description)
	}

	if out != nil && len(api.Result) != 0 {
		if err = json.Unmarshal(api.Result, out); err != nil {
			return fmt.Errorf("telegram: decode response: %w", err)
		}
	}
	return nil
}

// classifyError 把 Bot API 错误归类为 pkg/social 的语义化错误。
func classifyError(code int, description string) error {
	msg := strings.TrimSpace(description)
	switch code {
	case http.StatusUnauthorized, http.StatusForbidden:
		return social.ErrInvalidCredentials
	case http.StatusTooManyRequests:
		return social.ErrRateLimited
	case http.StatusNotFound:
		return social.ErrNotFound
	}
	if code == 0 && msg == "" {
		return fmt.Errorf("telegram: malformed response")
	}
	return fmt.Errorf("social platform error: telegram %d: %s", code, msg)
}
