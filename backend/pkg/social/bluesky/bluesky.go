// Package bluesky Bluesky（AT Protocol）发布适配器。
//
// 认证模型：App Password + createSession 换临时 access JWT，会话不持久化——
// 每次发布/回拉现签现用，免除 token 刷新与过期态管理（App Password 可随时在
// bsky.app 撤销重生成）。
//
// 用到的 XRPC：
//   - com.atproto.server.createSession  登录换 JWT
//   - com.atproto.repo.createRecord     发帖（app.bsky.feed.post）
//   - app.bsky.feed.getPosts            回拉互动计数
package bluesky

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
	Platform = "PLATFORM_BLUESKY"

	// DefaultPDS 公共聚合实例；自建 PDS 可通过 WithEndpoint 覆盖。
	DefaultPDS = "https://bsky.social"

	postCollection  = "app.bsky.feed.post"
	maxPostRunes    = 300  // Bluesky 单帖上限（服务端按 grapheme 计数，这里按 rune 预检，接近即可）
	defaultTimeout  = 30 * time.Second
	maxResponseSize = 10 << 20 // 10MB
)

// Client Bluesky 客户端（并发安全）。
type Client struct {
	httpClient *http.Client
	endpoint   string
}

// Option 客户端选项。
type Option func(*Client)

// WithEndpoint 覆盖 PDS 地址（自建 PDS / 测试桩）。
func WithEndpoint(endpoint string) Option {
	return func(c *Client) {
		if endpoint != "" {
			c.endpoint = strings.TrimRight(endpoint, "/")
		}
	}
}

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
	c := &Client{
		endpoint:   DefaultPDS,
		httpClient: &http.Client{Timeout: defaultTimeout},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// init 注册进 pkg/social 注册表（key 与 ent/proto 枚举值一致）。
func init() {
	social.Register(Platform, NewClient())
}

// ── XRPC 数据结构 ──

type xrpcError struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type createSessionRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type createSessionResponse struct {
	Did       string `json:"did"`
	Handle    string `json:"handle"`
	AccessJwt string `json:"accessJwt"`
}

type createRecordRequest struct {
	Repo       string        `json:"repo"`
	Collection string        `json:"collection"`
	Record     postRecordDoc `json:"record"`
}

// maxEmbedImages Bluesky 单帖配图上限。
const maxEmbedImages = 4

// postRecordDoc app.bsky.feed.post 记录体（纯文本/配图 embed + 创建时间）。
type postRecordDoc struct {
	Type      string          `json:"$type"`
	Text      string          `json:"text"`
	CreatedAt string          `json:"createdAt"`
	Embed     json.RawMessage `json:"embed,omitempty"`
}

type uploadBlobResponse struct {
	Blob json.RawMessage `json:"blob"`
}

type createRecordResponse struct {
	Uri string `json:"uri"`
	Cid string `json:"cid"`
}

type getPostsResponse struct {
	Posts []postView `json:"posts"`
}

type postView struct {
	Uri         string `json:"uri"`
	Cid         string `json:"cid"`
	LikeCount    int64  `json:"likeCount"`
	ReplyCount   int64  `json:"replyCount"`
	RepostCount  int64  `json:"repostCount"`
}

// ── social.Provider 实现 ──

// ValidateCredentials 校验凭据：能换取会话即有效。
func (c *Client) ValidateCredentials(ctx context.Context, creds social.Credentials) error {
	_, err := c.createSession(ctx, creds)
	return err
}

// Publish 发帖。成功返回 at:// URI（统计回拉定位键，存 external_id）
// 与 bsky.app 人类可读链接（存 external_uri）。
func (c *Client) Publish(ctx context.Context, creds social.Credentials, input *social.PublishInput) (*social.PublishResult, error) {
	if input == nil || strings.TrimSpace(input.Text) == "" {
		return nil, fmt.Errorf("bluesky: empty text")
	}
	if utf8.RuneCountInString(input.Text) > maxPostRunes {
		return nil, social.ErrContentTooLong
	}

	sess, err := c.createSession(ctx, creds)
	if err != nil {
		return nil, err
	}

	// 配图：uploadBlob 逐张上传（至多 4 张），组装 app.bsky.embed.images
	embed := json.RawMessage("")
	images := input.Media
	if len(images) > maxEmbedImages {
		images = images[:maxEmbedImages]
	}
	if len(images) != 0 {
		uploaded := make([]map[string]any, 0, len(images))
		for _, m := range images {
			blob, uerr := c.uploadBlob(ctx, sess.AccessJwt, m.ContentType, m.Data)
			if uerr != nil {
				return nil, uerr
			}
			uploaded = append(uploaded, map[string]any{"alt": "", "image": blob})
		}
		if embedRaw, merr := json.Marshal(map[string]any{
			"$type":  "app.bsky.embed.images",
			"images": uploaded,
		}); merr == nil {
			embed = embedRaw
		}
	}

	body := createRecordRequest{
		Repo:       sess.Did,
		Collection: postCollection,
		Record: postRecordDoc{
			Type:      postCollection,
			Text:      input.Text,
			CreatedAt: time.Now().UTC().Format(time.RFC3339),
			Embed:     embed,
		},
	}

	var resp createRecordResponse
	if err = c.do(ctx, http.MethodPost, "/xrpc/com.atproto.repo.createRecord", body, sess.AccessJwt, &resp); err != nil {
		return nil, err
	}

	return &social.PublishResult{
		ExternalID:  resp.Uri, // at://did/app.bsky.feed.post/rkey
		ExternalURI: webURL(creds.Get("handle"), resp.Uri),
	}, nil
}

// inboxNotification app.bsky.feed.notifications 的通知项。
type inboxNotification struct {
	URI       string `json:"uri"`
	Cid       string `json:"cid"`
	Reason    string `json:"reason"` // reply / mention / like / repost / quote / follow
	IndexedAt string `json:"indexedAt"`
	Author    struct {
		Did         string `json:"did"`
		Handle      string `json:"handle"`
		DisplayName string `json:"displayName"`
	} `json:"author"`
	Record struct {
		Text string `json:"text"`
	} `json:"record"`
}

type notificationsResponse struct {
	Notifications []inboxNotification `json:"notifications"`
}

// FetchInbox 回拉回复/提及（like/repost 仅计数无内容，跳过）。
func (c *Client) FetchInbox(ctx context.Context, creds social.Credentials, since time.Time) ([]social.InboxItem, error) {
	sess, err := c.createSession(ctx, creds)
	if err != nil {
		return nil, err
	}

	var resp notificationsResponse
	if err = c.do(ctx, http.MethodGet, "/xrpc/app.bsky.feed.notifications?limit=50", nil, sess.AccessJwt, &resp); err != nil {
		return nil, err
	}

	items := make([]social.InboxItem, 0, len(resp.Notifications))
	for _, n := range resp.Notifications {
		if n.Reason != "reply" && n.Reason != "mention" {
			continue
		}
		at := parseATTime(n.IndexedAt)
		if at.Before(since) {
			continue
		}
		typ := "REPLY"
		if n.Reason == "mention" {
			typ = "MENTION"
		}
		items = append(items, social.InboxItem{
			ExternalID:   n.URI,
			Type:         typ,
			AuthorName:   n.Author.DisplayName,
			AuthorHandle: n.Author.Handle,
			Content:      n.Record.Text,
			URL:          fmt.Sprintf("https://bsky.app/profile/%s/post/%s", n.Author.Handle, rkeyOf(n.URI)),
			OccurredAt:   at,
		})
	}
	return items, nil
}

// rkeyOf 从 at:// URI 取最后一段 rkey。
func rkeyOf(atURI string) string {
	parts := strings.Split(atURI, "/")
	if len(parts) == 0 {
		return atURI
	}
	return parts[len(parts)-1]
}

// parseATTime 解析 AT Protocol 的 RFC3339 时间戳。
func parseATTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}
	}
	return t
}

// FetchStats 回拉帖子互动计数。externalID 为发帖时记录的 at:// URI。
func (c *Client) FetchStats(ctx context.Context, creds social.Credentials, externalID, _ string) (*social.Stats, error) {
	if externalID == "" {
		return nil, fmt.Errorf("bluesky: empty external id")
	}

	sess, err := c.createSession(ctx, creds)
	if err != nil {
		return nil, err
	}

	var resp getPostsResponse
	if err = c.do(ctx, http.MethodGet, "/xrpc/app.bsky.feed.getPosts?uris="+url.QueryEscape(externalID), nil, sess.AccessJwt, &resp); err != nil {
		return nil, err
	}
	if len(resp.Posts) == 0 {
		return nil, social.ErrNotFound
	}

	p := resp.Posts[0]
	return &social.Stats{
		Likes:   clampUint32(p.LikeCount),
		Replies: clampUint32(p.ReplyCount),
		Reposts: clampUint32(p.RepostCount),
	}, nil
}

// ── 内部实现 ──

// uploadBlob 上传一张图片到用户的 blob 存储，返回可直接嵌入 record 的 blob JSON。
func (c *Client) uploadBlob(ctx context.Context, accessToken, contentType string, data []byte) (json.RawMessage, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.endpoint+"/xrpc/com.atproto.repo.uploadBlob", bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("bluesky: build upload request: %w", err)
	}
	req.Header.Set("User-Agent", social.UserAgent)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", contentType)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("bluesky: upload blob: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseSize))
	if err != nil {
		return nil, fmt.Errorf("bluesky: read upload response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, classifyError(resp.StatusCode, raw)
	}

	var up uploadBlobResponse
	if err = json.Unmarshal(raw, &up); err != nil {
		return nil, fmt.Errorf("bluesky: decode upload response: %w", err)
	}
	if len(up.Blob) == 0 {
		return nil, fmt.Errorf("bluesky: empty blob in upload response")
	}
	return up.Blob, nil
}

// createSession 登录换取临时会话。
func (c *Client) createSession(ctx context.Context, creds social.Credentials) (*createSessionResponse, error) {
	handle := creds.Get("handle")
	password := creds.Get("app_password")
	if handle == "" || password == "" {
		return nil, social.ErrInvalidCredentials
	}

	var resp createSessionResponse
	if err := c.do(ctx, http.MethodPost, "/xrpc/com.atproto.server.createSession",
		createSessionRequest{Identifier: handle, Password: password}, "", &resp); err != nil {
		return nil, err
	}
	if resp.AccessJwt == "" || resp.Did == "" {
		return nil, fmt.Errorf("bluesky: unexpected session response")
	}
	return &resp, nil
}

// do 执行一次 XRPC 调用：JSON 编解码、Bearer 注入、响应限长、错误归类。
func (c *Client) do(ctx context.Context, method, pathAndQuery string, body any, accessToken string, out any) error {
	var reader io.Reader
	if body != nil {
		buf, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("bluesky: marshal request: %w", err)
		}
		reader = strings.NewReader(string(buf))
	}

	req, err := http.NewRequestWithContext(ctx, method, c.endpoint+pathAndQuery, reader)
	if err != nil {
		return fmt.Errorf("bluesky: build request: %w", err)
	}
	req.Header.Set("User-Agent", "go-wind-cms-social/1.0")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("bluesky: request %s: %w", pathAndQuery, err)
	}
	defer func() { _ = resp.Body.Close() }()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseSize))
	if err != nil {
		return fmt.Errorf("bluesky: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return classifyError(resp.StatusCode, raw)
	}

	if out != nil && len(raw) != 0 {
		if err = json.Unmarshal(raw, out); err != nil {
			return fmt.Errorf("bluesky: decode response: %w", err)
		}
	}
	return nil
}

// classifyError 把非 2xx 响应归类为 pkg/social 的语义化错误。
func classifyError(status int, raw []byte) error {
	var xe xrpcError
	_ = json.Unmarshal(raw, &xe) // 错误体解析失败时保留状态码信息

	detail := strings.TrimSpace(xe.Message)
	if detail == "" {
		detail = xe.Error
	}
	msg := fmt.Sprintf("bluesky: http %d: %s", status, detail)

	switch {
	case status == http.StatusUnauthorized || status == http.StatusForbidden:
		return social.ErrInvalidCredentials
	case status == http.StatusTooManyRequests:
		return social.ErrRateLimited
	case status == http.StatusNotFound:
		return social.ErrNotFound
	case strings.Contains(strings.ToLower(xe.Message), "long") ||
		strings.Contains(strings.ToLower(xe.Message), "grapheme") ||
		strings.Contains(xe.Message, "300"):
		return social.ErrContentTooLong
	default:
		return fmt.Errorf("social platform error: %s", msg)
	}
}

// webURL 把 at:// URI 换成 bsky.app 人类可读链接。
func webURL(handle, atURI string) string {
	// at://did:plc:xxx/app.bsky.feed.post/rkey → https://bsky.app/profile/{handle}/post/{rkey}
	parts := strings.Split(atURI, "/")
	// ["at:", "", "did:plc:xxx", "app.bsky.feed.post", "rkey"]
	if len(parts) < 5 || parts[0] != "at:" {
		return atURI
	}
	rkey := parts[len(parts)-1]
	if handle == "" {
		handle = parts[2] // 无 handle 时退回 did 形式，链接仍可用
	}
	return fmt.Sprintf("https://bsky.app/profile/%s/post/%s", handle, rkey)
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
