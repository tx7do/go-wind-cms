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
	maxMediaCount   = 4   // 单帖配图上限
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

// Publish 发帖。带配图时先逐张上传媒体（至多 4 张），再以 media_ids 发帖。
// externalID 为平台 status id（统计回拉定位键），externalURI 为帖子链接。
func (c *Client) Publish(ctx context.Context, creds social.Credentials, input *social.PublishInput) (*social.PublishResult, error) {
	if input == nil {
		return nil, fmt.Errorf("mastodon: empty input")
	}
	if len(input.Media) == 0 && strings.TrimSpace(input.Text) == "" {
		return nil, fmt.Errorf("mastodon: empty text")
	}
	if utf8.RuneCountInString(input.Text) > maxStatusRunes {
		return nil, social.ErrContentTooLong
	}

	form := url.Values{}
	form.Set("status", input.Text)
	form.Set("visibility", "public")

	media := input.Media
	if len(media) > maxMediaCount {
		media = media[:maxMediaCount]
	}
	for _, m := range media {
		mediaID, merr := c.uploadMedia(ctx, creds, m)
		if merr != nil {
			return nil, merr
		}
		form.Add("media_ids[]", mediaID)
	}

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

// uploadMedia 上传一张图片，返回媒体 ID。
func (c *Client) uploadMedia(ctx context.Context, creds social.Credentials, m social.MediaAttachment) (string, error) {
	base, err := server(creds)
	if err != nil {
		return "", err
	}
	token := creds.Get("access_token")
	if token == "" {
		return "", social.ErrInvalidCredentials
	}

	buf, fw := social.NewMultipartWriter()
	if werr := social.WriteFilePart(fw, "file", m.Filename, m.Data); werr != nil {
		return "", fmt.Errorf("mastodon: build media part: %w", werr)
	}
	if ferr := fw.Close(); ferr != nil {
		return "", ferr
	}

	req, rerr := social.NewMultipartRequest(ctx, http.MethodPost, base+"/api/v1/media", buf, fw)
	if rerr != nil {
		return "", fmt.Errorf("mastodon: build media request: %w", rerr)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, derr := c.httpClient.Do(req)
	if derr != nil {
		return "", fmt.Errorf("mastodon: upload media: %w", derr)
	}
	defer func() { _ = resp.Body.Close() }()

	raw, rerr := io.ReadAll(io.LimitReader(resp.Body, maxResponseSize))
	if rerr != nil {
		return "", fmt.Errorf("mastodon: read media response: %w", rerr)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", classifyError(resp.StatusCode, raw)
	}

	var media struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(raw, &media); err != nil {
		return "", fmt.Errorf("mastodon: decode media response: %w", err)
	}
	if media.ID == "" {
		return "", fmt.Errorf("mastodon: unexpected media response")
	}
	return media.ID, nil
}

// inboxNotification /api/v1/notifications 的通知项。
type inboxNotification struct {
	ID     string `json:"id"`
	Type   string `json:"type"` // mention / favourite / reblog / follow ...
	Status *struct {
		ID      string `json:"id"`
		URL     string `json:"url"`
		Content string `json:"content"` // HTML
		Account struct {
			Username    string `json:"username"`
			DisplayName string `json:"display_name"`
		} `json:"account"`
	} `json:"status"`
	Account struct {
		Username    string `json:"username"`
		DisplayName string `json:"display_name"`
	} `json:"account"`
	CreatedAt string `json:"created_at"`
}

// stripHTMLTags 去除 HTML 标签为纯文本（收件箱摘要用，够用就好）。
func stripHTMLTags(s string) string {
	var b strings.Builder
	inTag := false
	for _, r := range s {
		switch r {
		case '<':
			inTag = true
		case '>':
			inTag = false
			b.WriteRune(' ')
		default:
			if !inTag {
				b.WriteRune(r)
			}
		}
	}
	return strings.Join(strings.Fields(b.String()), " ")
}

// parseMastoTime 解析 Mastodon 的 RFC3339 时间戳。
func parseMastoTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}
	}
	return t
}

// FetchInbox 回拉提及（mention 通知，含对他人帖子的回复）。
func (c *Client) FetchInbox(ctx context.Context, creds social.Credentials, since time.Time) ([]social.InboxItem, error) {
	base, err := server(creds)
	if err != nil {
		return nil, err
	}
	token := creds.Get("access_token")
	if token == "" {
		return nil, social.ErrInvalidCredentials
	}

	req, rerr := http.NewRequestWithContext(ctx, http.MethodGet,
		base+"/api/v1/notifications?types[]=mention&limit=40", nil)
	if rerr != nil {
		return nil, fmt.Errorf("mastodon: build notifications request: %w", rerr)
	}
	req.Header.Set("User-Agent", social.UserAgent)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, derr := c.httpClient.Do(req)
	if derr != nil {
		return nil, fmt.Errorf("mastodon: notifications: %w", derr)
	}
	defer func() { _ = resp.Body.Close() }()

	raw, rerr := io.ReadAll(io.LimitReader(resp.Body, maxResponseSize))
	if rerr != nil {
		return nil, fmt.Errorf("mastodon: read notifications: %w", rerr)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, classifyError(resp.StatusCode, raw)
	}

	var notifications []inboxNotification
	if err = json.Unmarshal(raw, &notifications); err != nil {
		return nil, fmt.Errorf("mastodon: decode notifications: %w", err)
	}

	items := make([]social.InboxItem, 0, len(notifications))
	for _, n := range notifications {
		occurred := parseMastoTime(n.CreatedAt)
		if occurred.Before(since) {
			continue
		}
		item := social.InboxItem{
			Type:         "MENTION",
			AuthorName:   n.Account.DisplayName,
			AuthorHandle: n.Account.Username,
			OccurredAt:   occurred,
		}
		if n.Status != nil {
			item.ExternalID = n.Status.ID
			item.URL = n.Status.URL
			item.Content = stripHTMLTags(n.Status.Content)
			item.AuthorHandle = n.Status.Account.Username
		} else {
			item.ExternalID = n.ID
		}
		if item.ExternalID == "" {
			continue
		}
		items = append(items, item)
	}
	return items, nil
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
