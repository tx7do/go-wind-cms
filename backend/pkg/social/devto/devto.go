// Package devto DEV.to（Forem）发布适配器。
//
// 认证模型：账户 API Key（Settings → Extensions → DEV Community API Keys），
// 自定义请求头 api-key 携带，长期有效。
//
// 用到的 API：
//   - GET  /api/users/me     校验 API Key
//   - POST /api/articles     发表文章（title + body_markdown）
//   - GET  /api/articles/:id 回拉互动计数
//
// 支持自建 Forem 实例：凭据 server 键可覆盖缺省 https://dev.to。
package devto

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
	Platform = "PLATFORM_DEVTO"

	// DefaultServer 官方社区；自建 Forem 可通过凭据 server 键覆盖。
	DefaultServer = "https://dev.to"

	maxBodyRunes    = 400000 // Forem 文章正文上限的宽松预检（实际以服务端为准）
	maxTitleRunes   = 128    // 标题从正文推导时的截断长度
	defaultTimeout  = 30 * time.Second
	maxResponseSize = 10 << 20 // 10MB
)

// Client DEV.to 客户端（并发安全）。
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

// NewClient 构造客户端。
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

type userView struct {
	Username string `json:"username"`
	ID       int    `json:"id"`
}

type articleView struct {
	ID                    int    `json:"id"`
	URL                   string `json:"url"`
	Title                 string `json:"title"`
	PositiveReactionsCount int64 `json:"positive_reactions_count"`
	CommentsCount         int64  `json:"comments_count"`
}

type createArticleRequest struct {
	Article articlePayload `json:"article"`
}

type articlePayload struct {
	Title        string `json:"title"`
	BodyMarkdown string `json:"body_markdown"`
	Published    bool   `json:"published"`
}

// ── social.Provider 实现 ──

// ValidateCredentials 校验 API Key：/api/users/me 通过即有效。
func (c *Client) ValidateCredentials(ctx context.Context, creds social.Credentials) error {
	var user userView
	if err := c.do(ctx, creds, http.MethodGet, "/api/users/me", nil, &user); err != nil {
		return err
	}
	return nil
}

// Publish 发表文章。externalID 为文章 ID（统计回拉定位键），externalURI 为文章链接。
// 标题取 input.Title；为空时从正文首行推导（去除 Markdown 标题符后截断）。
func (c *Client) Publish(ctx context.Context, creds social.Credentials, input *social.PublishInput) (*social.PublishResult, error) {
	if input == nil || strings.TrimSpace(input.Text) == "" {
		return nil, fmt.Errorf("devto: empty body")
	}
	if utf8.RuneCountInString(input.Text) > maxBodyRunes {
		return nil, social.ErrContentTooLong
	}
	title := deriveTitle(input.Title, input.Text)
	if title == "" {
		return nil, fmt.Errorf("devto: empty title")
	}

	body := createArticleRequest{Article: articlePayload{
		Title:        title,
		BodyMarkdown: input.Text,
		Published:    true,
	}}

	var article articleView
	if err := c.do(ctx, creds, http.MethodPost, "/api/articles", body, &article); err != nil {
		return nil, err
	}
	if article.ID == 0 {
		return nil, fmt.Errorf("devto: unexpected publish response")
	}

	return &social.PublishResult{
		ExternalID:  fmt.Sprintf("%d", article.ID),
		ExternalURI: article.URL,
	}, nil
}

// FetchStats 回拉文章互动计数（点赞 → positive_reactions，回复 → comments）。
func (c *Client) FetchStats(ctx context.Context, creds social.Credentials, externalID, _ string) (*social.Stats, error) {
	if externalID == "" {
		return nil, fmt.Errorf("devto: empty external id")
	}

	var article articleView
	if err := c.do(ctx, creds, http.MethodGet, "/api/articles/"+url.PathEscape(externalID), nil, &article); err != nil {
		return nil, err
	}

	return &social.Stats{
		Likes:   clampUint32(article.PositiveReactionsCount),
		Replies: clampUint32(article.CommentsCount),
	}, nil
}

// ── 内部实现 ──

// deriveTitle 解析标题：显式指定优先；否则取正文首个非空行，
// 去除 Markdown 标题符与首尾空白后截断（Forem 标题也可由 front matter 提供，
// 这里统一在请求体 title 字段给值，避免依赖用户写 front matter）。
func deriveTitle(title, body string) string {
	t := strings.TrimSpace(title)
	if t != "" {
		return truncate(t)
	}
	for _, line := range strings.Split(body, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		line = strings.TrimSpace(strings.TrimLeft(line, "#"))
		if line == "" {
			continue
		}
		return truncate(line)
	}
	return ""
}

func truncate(s string) string {
	s = strings.TrimSpace(s)
	if utf8.RuneCountInString(s) <= maxTitleRunes {
		return s
	}
	runes := []rune(s)
	return string(runes[:maxTitleRunes])
}

// server 从凭据解析实例地址。
func server(creds social.Credentials) (string, error) {
	s := strings.TrimRight(strings.TrimSpace(creds.Get("server")), "/")
	if s == "" {
		s = DefaultServer
	}
	if !strings.HasPrefix(s, "http://") && !strings.HasPrefix(s, "https://") {
		s = "https://" + s
	}
	return s, nil
}

// do 执行一次 API 调用：api-key 头注入、JSON 编解码、响应限长、错误归类。
func (c *Client) do(ctx context.Context, creds social.Credentials, method, path string, body any, out any) error {
	base, err := server(creds)
	if err != nil {
		return err
	}
	apiKey := creds.Get("api_key")
	if apiKey == "" {
		return social.ErrInvalidCredentials
	}

	var reqBody io.Reader
	if body != nil {
		buf, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("devto: marshal request: %w", err)
		}
		reqBody = strings.NewReader(string(buf))
	}

	req, err := http.NewRequestWithContext(ctx, method, base+path, reqBody)
	if err != nil {
		return fmt.Errorf("devto: build request: %w", err)
	}
	req.Header.Set("User-Agent", "go-wind-cms-social/1.0")
	req.Header.Set("api-key", apiKey)
	req.Header.Set("Accept", "application/vnd.forem.api-v1+json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("devto: request %s: %w", path, err)
	}
	defer func() { _ = resp.Body.Close() }()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseSize))
	if err != nil {
		return fmt.Errorf("devto: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return classifyError(resp.StatusCode, raw)
	}

	if out != nil && len(raw) != 0 {
		if err = json.Unmarshal(raw, out); err != nil {
			return fmt.Errorf("devto: decode response: %w", err)
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
	msg := fmt.Sprintf("devto: http %d: %s", status, detail)

	switch {
	case status == http.StatusUnauthorized || status == http.StatusForbidden:
		return social.ErrInvalidCredentials
	case status == http.StatusTooManyRequests:
		return social.ErrRateLimited
	case status == http.StatusNotFound:
		return social.ErrNotFound
	case status == http.StatusUnprocessableEntity && strings.Contains(strings.ToLower(detail), "title"):
		return fmt.Errorf("social platform error: %s", msg)
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
