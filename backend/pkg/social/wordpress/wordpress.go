// Package wordpress WordPress（REST API）发布适配器。
//
// 认证模型：WP 5.6+ 内置「应用程序密码」（用户 → 个人资料 → 应用密码），
// Basic Auth（用户名:应用密码）即开即用，零审核；覆盖自建站与 WordPress.com
// （后者需在后台创建应用密码并确保 REST 可访问）。
//
// 用到的 REST：
//   - GET /wp-json/wp/v2/users/me?context=edit  校验凭据
//   - POST /wp-json/wp/v2/posts                 发表文章
//
// 统计：自建站无平台侧互动统计，FetchStats 返回 ErrStatsUnsupported。
package wordpress

import (
	"context"
	"encoding/base64"
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
	Platform = "PLATFORM_WORDPRESS"

	defaultTimeout  = 30 * time.Second
	maxResponseSize = 10 << 20 // 10MB
)

// Client WordPress REST 客户端（并发安全）。
type Client struct {
	httpClient *http.Client
}

// Option 客户端选项。
type Option func(*Client)

// WithHTTPClient 覆盖底层 HTTP 客户端。
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		if httpClient != nil {
			c.httpClient = httpClient
		}
	}
}

// NewClient 构造客户端（站点地址来自凭据 server 键）。
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
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type postView struct {
	ID   int64  `json:"id"`
	Link string `json:"link"`
}

type mediaView struct {
	ID        int64  `json:"id"`
	SourceURL string `json:"source_url"`
}

type createPostRequest struct {
	Title        string `json:"title"`
	Content      string `json:"content"`
	Status       string `json:"status"`
	FeaturedMedia int64  `json:"featured_media,omitempty"`
}

// ── social.Provider 实现 ──

// ValidateCredentials 校验凭据：users/me 通过即有效。
func (c *Client) ValidateCredentials(ctx context.Context, creds social.Credentials) error {
	var user userView
	if err := c.do(ctx, creds, http.MethodGet, "/users/me?context=edit", nil, &user); err != nil {
		return err
	}
	return nil
}

// Publish 发表文章：标题取 input.Title（缺省取正文首行），正文转基础 HTML 段落。
// 带配图时第一张上传为特色图（featured_media），其余以 <img> 追加进正文。
func (c *Client) Publish(ctx context.Context, creds social.Credentials, input *social.PublishInput) (*social.PublishResult, error) {
	if input == nil || strings.TrimSpace(input.Text) == "" {
		return nil, fmt.Errorf("wordpress: empty content")
	}
	if utf8.RuneCountInString(input.Text) > 1<<20 {
		return nil, social.ErrContentTooLong
	}

	title := deriveTitle(input.Title, input.Text)
	if title == "" {
		return nil, fmt.Errorf("wordpress: empty title")
	}

	contentHTML := social.PlainTextToHTML(input.Text)

	body := createPostRequest{
		Title:   title,
		Content: contentHTML,
		Status:  "publish",
	}

	// 配图：第一张为特色图，其余追加进正文
	if len(input.Media) != 0 {
		var imgTags []string
		for i, m := range input.Media {
			mediaURL, uerr := c.uploadMedia(ctx, creds, m)
			if uerr != nil {
				return nil, uerr
			}
			if i == 0 {
				body.FeaturedMedia = mediaURL.ID
				continue
			}
			imgTags = append(imgTags, "<p><img src=\""+mediaURL.SourceURL+"\" alt=\"\" /></p>")
		}
		if len(imgTags) != 0 {
			body.Content += strings.Join(imgTags, "")
		}
	}

	var post postView
	if err := c.do(ctx, creds, http.MethodPost, "/posts", body, &post); err != nil {
		return nil, err
	}
	if post.ID == 0 {
		return nil, fmt.Errorf("wordpress: unexpected publish response")
	}

	return &social.PublishResult{
		ExternalID:  fmt.Sprintf("%d", post.ID),
		ExternalURI: post.Link,
	}, nil
}

// uploadMedia 上传一张图片到媒体库，返回媒体 ID 与公开地址。
func (c *Client) uploadMedia(ctx context.Context, creds social.Credentials, m social.MediaAttachment) (*mediaView, error) {
	base, err := siteURL(creds)
	if err != nil {
		return nil, err
	}
	username := creds.Get("username")
	appPassword := creds.Get("app_password")
	if username == "" || appPassword == "" {
		return nil, social.ErrInvalidCredentials
	}

	buf, fw := social.NewMultipartWriter()
	if werr := social.WriteFilePart(fw, "file", m.Filename, m.Data); werr != nil {
		return nil, fmt.Errorf("wordpress: build media part: %w", werr)
	}
	if ferr := fw.Close(); ferr != nil {
		return nil, ferr
	}

	req, rerr := http.NewRequestWithContext(ctx, http.MethodPost, base+"/wp-json/wp/v2/media", buf)
	if rerr != nil {
		return nil, fmt.Errorf("wordpress: build media request: %w", rerr)
	}
	req.Header.Set("User-Agent", social.UserAgent)
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(username+":"+appPassword)))
	req.Header.Set("Content-Type", fw.FormDataContentType())

	resp, derr := c.httpClient.Do(req)
	if derr != nil {
		return nil, fmt.Errorf("wordpress: upload media: %w", derr)
	}
	defer func() { _ = resp.Body.Close() }()

	raw, rerr := io.ReadAll(io.LimitReader(resp.Body, maxResponseSize))
	if rerr != nil {
		return nil, fmt.Errorf("wordpress: read media response: %w", rerr)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, classifyError(resp.StatusCode, raw)
	}

	var media mediaView
	if err := json.Unmarshal(raw, &media); err != nil {
		return nil, fmt.Errorf("wordpress: decode media response: %w", err)
	}
	if media.ID == 0 {
		return nil, fmt.Errorf("wordpress: unexpected media response")
	}
	return &media, nil
}

// FetchStats 自建站无平台侧互动统计。
func (c *Client) FetchStats(context.Context, social.Credentials, string, string) (*social.Stats, error) {
	return nil, social.ErrStatsUnsupported
}

// ── 内部实现 ──

// deriveTitle 标题：显式指定优先；否则取正文首个非空行去 Markdown 标题符。
func deriveTitle(title, body string) string {
	t := strings.TrimSpace(title)
	if t != "" {
		return t
	}
	for _, line := range strings.Split(body, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		return strings.TrimSpace(strings.TrimLeft(line, "#"))
	}
	return ""
}

// siteURL 从凭据解析站点地址（须含协议；如 https://example.com）。
func siteURL(creds social.Credentials) (string, error) {
	s := strings.TrimRight(strings.TrimSpace(creds.Get("server")), "/")
	if s == "" {
		return "", social.ErrInvalidCredentials
	}
	if !strings.HasPrefix(s, "http://") && !strings.HasPrefix(s, "https://") {
		s = "https://" + s
	}
	return s, nil
}

// do 执行一次 REST 调用：Basic Auth（用户名:应用密码）、JSON 编解码、错误归类。
func (c *Client) do(ctx context.Context, creds social.Credentials, method, path string, body any, out any) error {
	base, err := siteURL(creds)
	if err != nil {
		return err
	}
	username := creds.Get("username")
	appPassword := creds.Get("app_password")
	if username == "" || appPassword == "" {
		return social.ErrInvalidCredentials
	}

	var reqBody io.Reader
	if body != nil {
		buf, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("wordpress: marshal request: %w", err)
		}
		reqBody = strings.NewReader(string(buf))
	}

	req, err := http.NewRequestWithContext(ctx, method, base+"/wp-json/wp/v2"+path, reqBody)
	if err != nil {
		return fmt.Errorf("wordpress: build request: %w", err)
	}
	req.Header.Set("User-Agent", "go-wind-cms-social/1.0")
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(username+":"+appPassword)))
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("wordpress: request %s: %w", path, err)
	}
	defer func() { _ = resp.Body.Close() }()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseSize))
	if err != nil {
		return fmt.Errorf("wordpress: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return classifyError(resp.StatusCode, raw)
	}

	if out != nil && len(raw) != 0 {
		if err = json.Unmarshal(raw, out); err != nil {
			return fmt.Errorf("wordpress: decode response: %w", err)
		}
	}
	return nil
}

// classifyError 把非 2xx 响应归类为 pkg/social 的语义化错误。
func classifyError(status int, raw []byte) error {
	var apiErr struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	_ = json.Unmarshal(raw, &apiErr)

	detail := strings.TrimSpace(apiErr.Message)
	if detail == "" {
		detail = apiErr.Code
	}
	msg := fmt.Sprintf("wordpress %d: %s", status, detail)

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
