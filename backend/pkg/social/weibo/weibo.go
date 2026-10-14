// Package weibo 微博发布适配器（国内）。
//
// 认证模型：微博开放平台 OAuth2 access_token（需注册开发者 + 创建应用过审，
// 发博接口权限按应用审批——外部审核是使用前提，非代码问题）。
//
// 用到的 Open API（https://api.weibo.com/2）：
//   - GET  /account/rate_limit_status.json  校验 token
//   - POST /statuses/update.json            发博
//   - GET  /statuses/show.json              回拉互动计数（转评赞）
//
// 微博不回传帖子 URL，externalURI 留空。
package weibo

import (
	"bytes"
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
	Platform = "PLATFORM_WEIBO"

	apiBase         = "https://api.weibo.com/2"
	maxStatusRunes  = 2000 // 长文本上限；普通应用未开权限时服务端按 140 拒绝
	defaultTimeout  = 30 * time.Second
	maxResponseSize = 10 << 20 // 10MB
)

// Client 微博客户端（并发安全）。
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

type rateLimitView struct {
	UserID int64 `json:"uid"`
}

type statusView struct {
	ID             int64  `json:"id"`
	CommentsCount  int64  `json:"comments_count"`
	RepostsCount   int64  `json:"reposts_count"`
	AttitudesCount int64  `json:"attitudes_count"`
	User           struct {
		ID int64 `json:"id"`
	} `json:"user"`
}

// ── social.Provider 实现 ──

// ValidateCredentials 校验 token：rate_limit_status 通过即有效。
func (c *Client) ValidateCredentials(ctx context.Context, creds social.Credentials) error {
	var view rateLimitView
	return c.do(ctx, creds, http.MethodGet, "/account/rate_limit_status.json", nil, &view)
}

// Publish 发博。带配图时走 statuses/upload（发文+传图一体，取第一张配图）。
// 微博 API 不回传帖子 URL，externalURI 留空。
func (c *Client) Publish(ctx context.Context, creds social.Credentials, input *social.PublishInput) (*social.PublishResult, error) {
	if input == nil {
		return nil, fmt.Errorf("weibo: empty input")
	}
	if len(input.Media) == 0 && strings.TrimSpace(input.Text) == "" {
		return nil, fmt.Errorf("weibo: empty text")
	}
	if utf8.RuneCountInString(input.Text) > maxStatusRunes {
		return nil, social.ErrContentTooLong
	}

	if len(input.Media) != 0 {
		return c.publishWithPic(ctx, creds, input)
	}

	form := url.Values{}
	form.Set("status", input.Text)

	var status statusView
	if err := c.do(ctx, creds, http.MethodPost, "/statuses/update.json", form, &status); err != nil {
		return nil, err
	}
	if status.ID == 0 {
		return nil, fmt.Errorf("weibo: unexpected publish response")
	}

	return &social.PublishResult{
		ExternalID: fmt.Sprintf("%d", status.ID),
	}, nil
}

// publishWithPic 带图发博（multipart statuses/upload，取第一张配图）。
func (c *Client) publishWithPic(ctx context.Context, creds social.Credentials, input *social.PublishInput) (*social.PublishResult, error) {
	token := creds.Get("access_token")
	if token == "" {
		return nil, social.ErrInvalidCredentials
	}

	media := input.Media[0]
	buf, fw := social.NewMultipartWriter()
	_ = fw.WriteField("status", input.Text)
	if perr := social.WriteFilePart(fw, "pic", media.Filename, media.Data); perr != nil {
		return nil, fmt.Errorf("weibo: build pic part: %w", perr)
	}
	if ferr := fw.Close(); ferr != nil {
		return nil, ferr
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.endpoint+"/statuses/upload.json?access_token="+url.QueryEscape(token), bytes.NewReader(buf.Bytes()))
	if err != nil {
		return nil, fmt.Errorf("weibo: build request: %w", err)
	}
	req.Header.Set("User-Agent", social.UserAgent)
	req.Header.Set("Content-Type", fw.FormDataContentType())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("weibo: upload: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseSize))
	if err != nil {
		return nil, fmt.Errorf("weibo: read response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, classifyError(resp.StatusCode, raw)
	}

	var status statusView
	if err := json.Unmarshal(raw, &status); err != nil {
		return nil, fmt.Errorf("weibo: decode response: %w", err)
	}
	if status.ID == 0 {
		return nil, fmt.Errorf("weibo: unexpected publish response")
	}

	return &social.PublishResult{
		ExternalID: fmt.Sprintf("%d", status.ID),
	}, nil
}

// FetchStats 回拉转评赞计数（读接口权限取决于应用审批结果，运行时不足则报错）。
func (c *Client) FetchStats(ctx context.Context, creds social.Credentials, externalID, _ string) (*social.Stats, error) {
	if externalID == "" {
		return nil, fmt.Errorf("weibo: empty external id")
	}

	var status statusView
	if err := c.do(ctx, creds, http.MethodGet, "/statuses/show.json?id="+url.QueryEscape(externalID), nil, &status); err != nil {
		return nil, err
	}

	return &social.Stats{
		Likes:   clampUint32(status.AttitudesCount),
		Replies: clampUint32(status.CommentsCount),
		Reposts: clampUint32(status.RepostsCount),
	}, nil
}

// ── 内部实现 ──

// do 执行一次 Open API 调用：token 走 query 参数（官方惯例）、form/JSON 编码、
// 响应限长、错误归类。
func (c *Client) do(ctx context.Context, creds social.Credentials, method, path string, body any, out any) error {
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
			return fmt.Errorf("weibo: marshal request: %w", err)
		}
		reqBody = strings.NewReader(string(buf))
		contentType = "application/json"
	}

	sep := "?"
	if strings.Contains(path, "?") {
		sep = "&"
	}
	req, err := http.NewRequestWithContext(ctx, method,
		c.endpoint+path+sep+"access_token="+url.QueryEscape(token), reqBody)
	if err != nil {
		return fmt.Errorf("weibo: build request: %w", err)
	}
	req.Header.Set("User-Agent", "go-wind-cms-social/1.0")
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("weibo: request %s: %w", path, err)
	}
	defer func() { _ = resp.Body.Close() }()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseSize))
	if err != nil {
		return fmt.Errorf("weibo: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return classifyError(resp.StatusCode, raw)
	}

	if out != nil && len(raw) != 0 {
		if err = json.Unmarshal(raw, out); err != nil {
			return fmt.Errorf("weibo: decode response: %w", err)
		}
	}
	return nil
}

// classifyError 把非 2xx 响应归类为 pkg/social 的语义化错误。
func classifyError(status int, raw []byte) error {
	var apiErr struct {
		ErrorCode int    `json:"error_code"`
		Error     string `json:"error"`
	}
	_ = json.Unmarshal(raw, &apiErr)

	detail := strings.TrimSpace(apiErr.Error)
	msg := fmt.Sprintf("weibo %d (%d): %s", status, apiErr.ErrorCode, detail)

	switch {
	case status == http.StatusUnauthorized || status == http.StatusForbidden ||
		apiErr.ErrorCode == 21332 || apiErr.ErrorCode == 21327:
		// token 无效/过期/应用审核未过
		return social.ErrInvalidCredentials
	case status == http.StatusTooManyRequests || apiErr.ErrorCode == 130:
		return social.ErrRateLimited
	case status == http.StatusNotFound:
		return social.ErrNotFound
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
