// Package metaweblog MetaWeblog（XML-RPC）博客发布通用适配器。
//
// 覆盖对象：博客园（https://rpc.cnblogs.com/metaweblog/{博客名}）、Typecho 及
// 一切声明支持 MetaWeblog API 的博客系统——RPC 端点由用户凭据给出，本适配器
// 不绑定具体站点。博客园需在其后台确认 MetaWeblog/API 访问已开启。
//
// 协议（XML-RPC，本适配器用到的三个方法）：
//   - blogger.getUsersBlogs(appKey, username, password)  校验凭据 + 取博客 ID
//   - metaWeblog.newPost(blogid, username, password, struct, publish)  发文
//     struct: {title, description(HTML 正文), categories[], dateCreated?}
//   - （回拉互动计数不在 MetaWeblog 规范内，FetchStats 返回 ErrStatsUnsupported）
//
// 凭据：
//   - rpc_url:  XML-RPC 端点（必填）
//   - username / password: 博客账号
//   - blog_id:  博客 ID（可选；留空时自动取 getUsersBlogs 的第一个博客）
package metaweblog

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/kolo/xmlrpc"

	"go-wind-cms/pkg/social"
)

// Platform 本适配器注册进 pkg/social 注册表的键（proto 枚举名）。
const Platform = "PLATFORM_METAWEBLOG"

// blogContent metaWeblog.newPost 的 content 结构体字段。
type blogContent struct {
	Title       string   `xmlrpc:"title"`
	Description string   `xmlrpc:"description"`
	Categories  []string `xmlrpc:"categories,omitempty"`
}

type userBlog struct {
	BlogID   string `xmlrpc:"blogid"`
	URL      string `xmlrpc:"url"`
	BlogName string `xmlrpc:"blogName"`
}

// Client MetaWeblog 客户端（并发安全；底层每次调用新建 XML-RPC 连接）。
type Client struct{}

// NewClient 构造客户端（端点来自凭据 rpc_url 键，逐调用解析）。
func NewClient() *Client {
	return &Client{}
}

// init 注册进 pkg/social 注册表（key 与 ent/proto 枚举值一致）。
func init() {
	social.Register(Platform, NewClient())
}

// ── social.Provider 实现 ──

// ValidateCredentials 校验凭据：getUsersBlogs 能取到博客即有效。
func (c *Client) ValidateCredentials(ctx context.Context, creds social.Credentials) error {
	_, err := c.listUserBlogs(ctx, creds)
	return err
}

// Publish 发博文：标题取 input.Title（缺省取正文首行），正文转基础 HTML 段落。
// external_id 为平台返回的 postid；MetaWeblog 不回传 permalink，externalURI 留空。
func (c *Client) Publish(ctx context.Context, creds social.Credentials, input *social.PublishInput) (*social.PublishResult, error) {
	if input == nil || strings.TrimSpace(input.Text) == "" {
		return nil, fmt.Errorf("metaweblog: empty content")
	}

	title := deriveTitle(input.Title, input.Text)
	if title == "" {
		return nil, fmt.Errorf("metaweblog: empty title")
	}

	blogID := creds.Get("blog_id")
	if blogID == "" {
		blogs, err := c.listUserBlogs(ctx, creds)
		if err != nil {
			return nil, err
		}
		if len(blogs) == 0 {
			return nil, fmt.Errorf("metaweblog: no blog found for this account")
		}
		blogID = blogs[0].BlogID
	}

	content := blogContent{
		Title:       title,
		Description: social.PlainTextToHTML(input.Text),
	}

	postID, err := c.callPostID(ctx, creds, blogID, content)
	if err != nil {
		return nil, err
	}

	return &social.PublishResult{
		ExternalID: postID,
	}, nil
}

// FetchStats 互动计数不在 MetaWeblog 规范内。
func (c *Client) FetchStats(context.Context, social.Credentials, string, string) (*social.Stats, error) {
	return nil, social.ErrStatsUnsupported
}

// ── 内部实现 ──

// callXMLRPC 建立一次 XML-RPC 连接并调用指定方法（端点来自凭据 rpc_url）。
func (c *Client) callXMLRPC(ctx context.Context, creds social.Credentials, method string, args []any, out any) error {
	rpcURL := strings.TrimSpace(creds.Get("rpc_url"))
	if rpcURL == "" {
		return social.ErrInvalidCredentials
	}
	if !strings.HasPrefix(rpcURL, "http://") && !strings.HasPrefix(rpcURL, "https://") {
		rpcURL = "https://" + rpcURL
	}

	client, err := xmlrpc.NewClient(rpcURL, &ctxTransport{ctx: ctx, base: http.DefaultTransport})
	if err != nil {
		return fmt.Errorf("metaweblog: connect %s: %w", rpcURL, err)
	}

	if err = client.Call(method, args, out); err != nil {
		return classifyXMLError(err)
	}
	return nil
}

// ctxTransport 把调用级 ctx 透传给 XML-RPC 库的 HTTP 层
//（kolo/xmlrpc 的 NewClient 收 RoundTripper 而非 ctx）。
type ctxTransport struct {
	ctx  context.Context
	base http.RoundTripper
}

func (t *ctxTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.base.RoundTrip(req.WithContext(t.ctx))
}

func (c *Client) listUserBlogs(ctx context.Context, creds social.Credentials) ([]userBlog, error) {
	username := creds.Get("username")
	password := creds.Get("password")
	if username == "" || password == "" {
		return nil, social.ErrInvalidCredentials
	}

	var blogs []userBlog
	// appKey 参数在博客园等实现中不校验，传空串
	if err := c.callXMLRPC(ctx, creds, "blogger.getUsersBlogs",
		[]any{"", username, password}, &blogs); err != nil {
		return nil, err
	}
	return blogs, nil
}

func (c *Client) callPostID(ctx context.Context, creds social.Credentials, blogID string, content blogContent) (string, error) {
	username := creds.Get("username")
	password := creds.Get("password")
	if username == "" || password == "" {
		return "", social.ErrInvalidCredentials
	}

	var postID string
	if err := c.callXMLRPC(ctx, creds, "metaWeblog.newPost",
		[]any{blogID, username, password, content, true}, &postID); err != nil {
		return "", err
	}
	if strings.TrimSpace(postID) == "" {
		return "", fmt.Errorf("metaweblog: unexpected empty post id")
	}
	return postID, nil
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

// classifyXMLError 把 XML-RPC 错误归类为 pkg/social 的语义化错误。
// 各实现 faultString 不统一，按常见认证关键词归类，其余归平台错误。
func classifyXMLError(err error) error {
	msg := strings.ToLower(err.Error())
	authMarkers := []string{"password", "username", "用户名", "密码", "认证", "auth", "403", "login"}
	for _, marker := range authMarkers {
		if strings.Contains(msg, marker) {
			return social.ErrInvalidCredentials
		}
	}
	return fmt.Errorf("social platform error: metaweblog: %w", err)
}
