// Package linkedin LinkedIn 个人动态发布适配器。
//
// 认证模型：LinkedIn 开发者应用 + OAuth2 access_token（self-serve 应用可即时
// 创建，个人动态 w_member_social 无需额外审核；组织号发文需过 Community
// Management 审核，本适配器仅支持个人动态）。
//
// 用到的 API：
//   - GET  /v2/userinfo   校验 token 并取成员 sub（OIDC）
//   - POST /rest/posts    发个人动态（Posts API）
//
// 统计：互动计数需额外的 socialActions 权限与 URN 拼接，第一版不回拉
// （FetchStats 返回 ErrStatsUnsupported）。
package linkedin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"go-wind-cms/pkg/social"
)

const (
	// Platform 本适配器注册进 pkg/social 注册表的键（proto 枚举名）。
	Platform = "PLATFORM_LINKEDIN"

	apiBase             = "https://api.linkedin.com"
	linkedinAPIVersion  = "202405" // Posts API 要求带版本头
	maxCommentaryRunes  = 3000     // Posts API commentary 上限
	defaultTimeout      = 30 * time.Second
	maxResponseSize     = 10 << 20 // 10MB
)

// Client LinkedIn 客户端（并发安全）。
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

type userInfoView struct {
	Sub         string `json:"sub"` // OIDC 成员标识 → urn:li:person:{sub}
	Name        string `json:"name"`
}

type createPostRequest struct {
	Author      string            `json:"author"`
	Commentary  string            `json:"commentary"`
	Visibility  string            `json:"visibility"`
	Distribution map[string]any   `json:"distribution"`
}

// ── social.Provider 实现 ──

// ValidateCredentials 校验 token：userinfo 通过即有效。
func (c *Client) ValidateCredentials(ctx context.Context, creds social.Credentials) error {
	var user userInfoView
	if err := c.do(ctx, creds, http.MethodGet, "/v2/userinfo", nil, &user); err != nil {
		return err
	}
	if user.Sub == "" {
		return fmt.Errorf("linkedin: unexpected userinfo response")
	}
	return nil
}

// Publish 发个人动态（PUBLIC 可见性）。external_id 为 RestLi 帖子 ID，
// external_uri 构造为 feed 永久链接。
func (c *Client) Publish(ctx context.Context, creds social.Credentials, input *social.PublishInput) (*social.PublishResult, error) {
	if input == nil || strings.TrimSpace(input.Text) == "" {
		return nil, fmt.Errorf("linkedin: empty text")
	}
	if utf8.RuneCountInString(input.Text) > maxCommentaryRunes {
		return nil, social.ErrContentTooLong
	}

	var user userInfoView
	if err := c.do(ctx, creds, http.MethodGet, "/v2/userinfo", nil, &user); err != nil {
		return nil, err
	}
	if user.Sub == "" {
		return nil, fmt.Errorf("linkedin: cannot resolve member urn")
	}

	body := createPostRequest{
		Author:      "urn:li:person:" + user.Sub,
		Commentary:  input.Text,
		Visibility:  "PUBLIC",
		Distribution: map[string]any{"linkedInDistributionSource": struct{}{}},
	}

	restLiID, err := c.doPost(ctx, creds, "/rest/posts", body)
	if err != nil {
		return nil, err
	}

	return &social.PublishResult{
		ExternalID:  restLiID,
		ExternalURI: "https://www.linkedin.com/feed/update/urn:li:activity:" + restLiID,
	}, nil
}

// FetchStats 互动计数需 socialActions 权限与 URN 拼接，第一版不回拉。
func (c *Client) FetchStats(context.Context, social.Credentials, string, string) (*social.Stats, error) {
	return nil, social.ErrStatsUnsupported
}

// ── 内部实现 ──

// do 执行一次 GET 调用（Bearer）。
func (c *Client) do(ctx context.Context, creds social.Credentials, method, path string, body any, out any) error {
	respBody, status, err := c.request(ctx, creds, method, path, body)
	if err != nil {
		return err
	}
	if status < 200 || status >= 300 {
		return classifyError(status, respBody)
	}
	if out != nil && len(respBody) != 0 {
		if err = json.Unmarshal(respBody, out); err != nil {
			return fmt.Errorf("linkedin: decode response: %w", err)
		}
	}
	return nil
}

// doPost 执行一次 POST 并返回 X-RestLi-Id 响应头（创建类接口的定位键）。
func (c *Client) doPost(ctx context.Context, creds social.Credentials, path string, body any) (string, error) {
	respBody, status, hdrs, err := c.requestWithHeaders(ctx, creds, path, body)
	if err != nil {
		return "", err
	}
	if status < 200 || status >= 300 {
		return "", classifyError(status, respBody)
	}
	id := strings.TrimSpace(hdrs.Get("X-RestLi-Id"))
	if id == "" {
		return "", fmt.Errorf("linkedin: unexpected publish response (missing X-RestLi-Id)")
	}
	return id, nil
}

func (c *Client) request(ctx context.Context, creds social.Credentials, method, path string, body any) ([]byte, int, error) {
	respBody, status, _, err := c.requestWithHeadersFull(ctx, creds, method, path, body)
	return respBody, status, err
}

func (c *Client) requestWithHeaders(ctx context.Context, creds social.Credentials, path string, body any) ([]byte, int, http.Header, error) {
	return c.requestWithHeadersFull(ctx, creds, http.MethodPost, path, body)
}

func (c *Client) requestWithHeadersFull(ctx context.Context, creds social.Credentials, method, path string, body any) ([]byte, int, http.Header, error) {
	token := creds.Get("access_token")
	if token == "" {
		return nil, 0, nil, social.ErrInvalidCredentials
	}

	var reqBody io.Reader
	if body != nil {
		buf, err := json.Marshal(body)
		if err != nil {
			return nil, 0, nil, fmt.Errorf("linkedin: marshal request: %w", err)
		}
		reqBody = strings.NewReader(string(buf))
	}

	req, err := http.NewRequestWithContext(ctx, method, c.endpoint+path, reqBody)
	if err != nil {
		return nil, 0, nil, fmt.Errorf("linkedin: build request: %w", err)
	}
	req.Header.Set("User-Agent", "go-wind-cms-social/1.0")
	req.Header.Set("Authorization", "Bearer "+token)
	// Posts API 必需头：版本号 + RestLi 协议声明
	req.Header.Set("LinkedIn-Version", linkedinAPIVersion)
	req.Header.Set("X-Restli-Protocol", "2.0.0")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, nil, fmt.Errorf("linkedin: request %s: %w", path, err)
	}
	defer func() { _ = resp.Body.Close() }()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseSize))
	if err != nil {
		return nil, 0, nil, fmt.Errorf("linkedin: read response: %w", err)
	}
	return raw, resp.StatusCode, resp.Header, nil
}

// classifyError 把非 2xx 响应归类为 pkg/social 的语义化错误。
func classifyError(status int, raw []byte) error {
	var apiErr struct {
		Message string `json:"message"`
	}
	_ = json.Unmarshal(raw, &apiErr)

	detail := strings.TrimSpace(apiErr.Message)
	msg := fmt.Sprintf("linkedin %d: %s", status, detail)

	switch status {
	case http.StatusUnauthorized, http.StatusForbidden:
		return social.ErrInvalidCredentials
	case http.StatusTooManyRequests:
		return social.ErrRateLimited
	case http.StatusNotFound:
		return social.ErrNotFound
	default:
		return fmt.Errorf("social platform error: %s", msg)
	}
}
