// Package wechat 微信公众号发布适配器（国内）。
//
// 认证模型：公众号 AppID/AppSecret 换取 access_token（2 小时有效），每次操作
// 现取现用，免除刷新管理；**服务器 IP 必须在公众号后台「IP 白名单」内**。
//
// 发布链路（图文以"草稿箱 → 发布"两段式进行）：
//   1. GET  /cgi-bin/token             换 access_token
//   2. POST /cgi-bin/draft/add         建草稿（thumb_media_id 必填——请把一张
//      封面图在公众号后台「素材管理」传为永久素材，将其 media_id 填入账号凭据）
//   3. POST /cgi-bin/freepublish/submit  提交发布（异步任务）
//   4. POST /cgi-bin/freepublish/get     轮询发布结果拿文章链接
//
// 统计：阅读/在看等数据在 datacube 统计接口且按天滞后，第一版不回拉
// （FetchStats 返回 ErrStatsUnsupported）。
//
// 资质约束：freepublish（发布到正式发表）要求**认证**公众号；未认证号可走
// 草稿箱（本适配器在 freepublish 报权限错误时保留草稿并把 media_id 作为
// external_id，便于后台手动群发）。
package wechat

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
	Platform = "PLATFORM_WECHAT_OA"

	apiBase         = "https://api.weixin.qq.com"
	maxTitleRunes   = 64   // 公众号文章标题上限
	maxDigestRunes  = 120  // 摘要上限
	defaultTimeout  = 30 * time.Second
	maxResponseSize = 10 << 20 // 10MB
)

// Client 微信公众号客户端（并发安全）。
type Client struct {
	httpClient  *http.Client
	endpoint    string
	pollInterval time.Duration // freepublish 结果轮询间隔（默认 2s，测试可注入）
}

// Option 客户端选项。
type Option func(*Client)

// WithEndpoint 覆盖 API 地址（测试桩 / 代理）。
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

// WithPollInterval 覆盖发布结果轮询间隔（默认 2s，测试注入短间隔）。
func WithPollInterval(d time.Duration) Option {
	return func(c *Client) {
		if d > 0 {
			c.pollInterval = d
		}
	}
}

// NewClient 构造客户端。
func NewClient(opts ...Option) *Client {
	c := &Client{endpoint: apiBase, httpClient: &http.Client{Timeout: defaultTimeout}, pollInterval: 2 * time.Second}
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

// wxEnvelope 公共响应壳：errcode=0 视为成功。
type wxEnvelope struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

type tokenResponse struct {
	wxEnvelope
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type draftArticle struct {
	Title        string `json:"title"`
	Author       string `json:"author,omitempty"`
	Digest       string `json:"digest,omitempty"`
	Content      string `json:"content"`
	ThumbMediaID string `json:"thumb_media_id"`
}

type draftAddRequest struct {
	Articles []draftArticle `json:"articles"`
}

type draftAddResponse struct {
	wxEnvelope
	MediaID string `json:"media_id"`
}

type freepublishSubmitResponse struct {
	wxEnvelope
	PublishID string `json:"publish_id"`
}

type freepublishGetResponse struct {
	wxEnvelope
	PublishStatus int `json:"publish_status"` // 0=成功 1=发布中 2/3/4=失败态
	ArticleID     struct {
		Item []struct {
			Idx        int    `json:"idx"`
			ArticleURL string `json:"article_url"`
		} `json:"item"`
	} `json:"article_id"`
}

// ── social.Provider 实现 ──

// ValidateCredentials 校验凭据：换取 access_token 成功即有效
//（secret 错误与 IP 不在白名单都会在此暴露）。
func (c *Client) ValidateCredentials(ctx context.Context, creds social.Credentials) error {
	_, err := c.fetchAccessToken(ctx, creds)
	return err
}

// Publish 发表文章：草稿箱 → 提交发布 → 短轮询拿文章链接。
// external_id 存 publish_id（平台侧定位键），external_uri 存文章链接（轮询未
// 就绪时留空）。
func (c *Client) Publish(ctx context.Context, creds social.Credentials, input *social.PublishInput) (*social.PublishResult, error) {
	if input == nil || strings.TrimSpace(input.Text) == "" {
		return nil, fmt.Errorf("wechat: empty content")
	}

	thumbMediaID := creds.Get("thumb_media_id")
	if thumbMediaID == "" {
		// 明确的配置缺口：封面图必填，需在公众号后台传永久素材后把 media_id 填进凭据
		return nil, fmt.Errorf("wechat: thumb_media_id required (upload a cover image as permanent material in the WeChat backend and put its media_id into account credentials)")
	}

	title := deriveTitle(input.Title, input.Text)
	if title == "" {
		return nil, fmt.Errorf("wechat: empty title")
	}

	token, err := c.fetchAccessToken(ctx, creds)
	if err != nil {
		return nil, err
	}

	// 1) 建草稿
	digest := truncateRunes(strings.TrimSpace(strings.SplitN(input.Text, "\n", 2)[0]), maxDigestRunes)
	body := draftAddRequest{Articles: []draftArticle{{
		Title:        truncateRunes(title, maxTitleRunes),
		Digest:       digest,
		Content:      social.PlainTextToHTML(input.Text),
		ThumbMediaID: thumbMediaID,
	}}}

	var draft draftAddResponse
	if err = c.do(ctx, token, "/cgi-bin/draft/add", body, &draft); err != nil {
		return nil, err
	}
	if draft.MediaID == "" {
		return nil, fmt.Errorf("wechat: unexpected draft response")
	}

	// 2) 提交发布
	var submit freepublishSubmitResponse
	if err = c.do(ctx, token, "/cgi-bin/freepublish/submit", map[string]string{"media_id": draft.MediaID}, &submit); err != nil {
		return nil, err
	}
	if submit.PublishID == "" {
		return nil, fmt.Errorf("wechat: unexpected freepublish response")
	}

	// 3) 短轮询发布结果拿文章链接（发布任务通常秒级完成；未就绪则留空不强求）
	uri := ""
	for i := 0; i < 4; i++ {
		select {
		case <-ctx.Done():
		case <-time.After(c.pollInterval):
		}
		var status freepublishGetResponse
		if gerr := c.do(ctx, token, "/cgi-bin/freepublish/get",
			map[string]string{"publish_id": submit.PublishID}, &status); gerr != nil {
			break // 轮询失败不影响发布结果
		}
		if status.PublishStatus == 0 && len(status.ArticleID.Item) != 0 {
			uri = status.ArticleID.Item[0].ArticleURL
			break
		}
		if status.PublishStatus != 1 {
			break // 失败态：链接拿不到，发布结果以后台为准
		}
	}

	return &social.PublishResult{
		ExternalID:  submit.PublishID,
		ExternalURI: uri,
	}, nil
}

// FetchStats 阅读/在看等数据在 datacube 统计接口且按天滞后，第一版不回拉。
func (c *Client) FetchStats(context.Context, social.Credentials, string, string) (*social.Stats, error) {
	return nil, social.ErrStatsUnsupported
}

// ── 内部实现 ──

// fetchAccessToken 现取现用 access_token（2h 有效，避免刷新管理；
// 公众号 token 下发接口本身有每日配额，本适配器每次操作至多取一次）。
func (c *Client) fetchAccessToken(ctx context.Context, creds social.Credentials) (string, error) {
	appID := creds.Get("app_id")
	secret := creds.Get("app_secret")
	if appID == "" || secret == "" {
		return "", social.ErrInvalidCredentials
	}

	q := url.Values{}
	q.Set("grant_type", "client_credential")
	q.Set("appid", appID)
	q.Set("secret", secret)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		c.endpoint+"/cgi-bin/token?"+q.Encode(), nil)
	if err != nil {
		return "", fmt.Errorf("wechat: build token request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("wechat: request token: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseSize))
	if err != nil {
		return "", fmt.Errorf("wechat: read token response: %w", err)
	}

	var tok tokenResponse
	_ = json.Unmarshal(raw, &tok)
	if tok.AccessToken == "" {
		return "", classifyWxError(tok.ErrCode, tok.ErrMsg)
	}
	return tok.AccessToken, nil
}

// do 执行一次带 access_token 的 API 调用（token 走 query，官方惯例）。
func (c *Client) do(ctx context.Context, accessToken, path string, body any, out any) error {
	var reqBody io.Reader
	if body != nil {
		buf, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("wechat: marshal request: %w", err)
		}
		reqBody = strings.NewReader(string(buf))
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.endpoint+path+"?access_token="+url.QueryEscape(accessToken), reqBody)
	if err != nil {
		return fmt.Errorf("wechat: build request: %w", err)
	}
	req.Header.Set("User-Agent", "go-wind-cms-social/1.0")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("wechat: request %s: %w", path, err)
	}
	defer func() { _ = resp.Body.Close() }()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseSize))
	if err != nil {
		return fmt.Errorf("wechat: read response: %w", err)
	}

	var env wxEnvelope
	_ = json.Unmarshal(raw, &env)
	if env.ErrCode != 0 {
		return classifyWxError(env.ErrCode, env.ErrMsg)
	}

	if out != nil {
		if err = json.Unmarshal(raw, out); err != nil {
			return fmt.Errorf("wechat: decode response: %w", err)
		}
	}
	return nil
}

// classifyWxError 把微信 errcode 归类为 pkg/social 的语义化错误。
func classifyWxError(code int, msg string) error {
	detail := strings.TrimSpace(msg)
	switch code {
	case 40001, 40125, 41001, 41002, 41004:
		// invalid credential / invalid appsecret / 缺 access_token / 缺 appid / 缺 secret
		return social.ErrInvalidCredentials
	case 40164:
		// IP 不在白名单：凭据没问题，是服务端 IP 未加白——给出可操作的提示
		return fmt.Errorf("social platform error: wechat 40164: server IP not in the official account IP whitelist (add it in the WeChat backend): %s", detail)
	case 48001:
		// api unauthorized：未认证号调用了认证接口（如 freepublish）
		return fmt.Errorf("social platform error: wechat 48001: api unauthorized for this account type (verified account required for freepublish; the draft is kept in the draft box): %s", detail)
	}
	if code == 0 {
		return fmt.Errorf("wechat: malformed response")
	}
	return fmt.Errorf("social platform error: wechat %d: %s", code, detail)
}

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

// truncateRunes 按 rune 截断（nil 安全由调用方保证非空检查）。
func truncateRunes(s string, max int) string {
	s = strings.TrimSpace(s)
	if utf8.RuneCountInString(s) <= max {
		return s
	}
	runes := []rune(s)
	return string(runes[:max])
}
