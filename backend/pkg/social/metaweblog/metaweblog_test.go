package metaweblog

import (
	"context"
	"encoding/xml"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-wind-cms/pkg/social"
)

func creds(server string) social.Credentials {
	return social.Credentials{
		"rpc_url":  server,
		"username": "alice",
		"password": "secret",
	}
}

// xmlRPCResponse 包一层标准的 methodResponse 信封。
func xmlRPCResponse(inner string) string {
	return `<?xml version="1.0"?><methodResponse><params>` + inner + `</params></methodResponse>`
}

func xmlRPCFault(code int, msg string) string {
	return `<?xml version="1.0"?><methodResponse><fault><value><struct>` +
		`<member><name>faultCode</name><value><int>` + itoa(code) + `</int></value></member>` +
		`<member><name>faultString</name><value><string>` + msg + `</string></value></member>` +
		`</struct></value></fault></methodResponse>`
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	digits := ""
	for i > 0 {
		digits = string(rune('0'+i%10)) + digits
		i /= 10
	}
	return digits
}

// newTestServer 起一个假 MetaWeblog 端点，按方法名分发并断言 XML-RPC 载荷。
func newTestServer(t *testing.T, passwordOK bool, postID string) (*Client, *httptest.Server, *[]string) {
	var methods []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var method string
		var paramsBlock string
		// 极简解析：取 methodName 与 payload 关键词
		raw := readAll(t, r)
		for _, m := range []string{"blogger.getUsersBlogs", "metaWeblog.newPost"} {
			if strings.Contains(raw, m) {
				method = m
				break
			}
		}
		methods = append(methods, method)
		paramsBlock = raw

		w.Header().Set("Content-Type", "text/xml")
		switch method {
		case "blogger.getUsersBlogs":
			if !paramsBlockContains(paramsBlock, "alice", "secret") || !passwordOK {
				_, _ = w.Write([]byte(xmlRPCFault(403, "用户名或密码错误")))
				return
			}
			_, _ = w.Write([]byte(xmlRPCResponse(`<param><value><array><data><value><struct>` +
				`<member><name>blogid</name><value><string>blog-1</string></value></member>` +
				`<member><name>url</name><value><string>https://alice.cnblogs.com</string></value></member>` +
				`<member><name>blogName</name><value><string>Alice 的博客</string></value></member>` +
				`</struct></value></data></array></value></param>`)))
		case "metaWeblog.newPost":
			if !paramsBlockContains(paramsBlock, "alice", "secret") || !passwordOK {
				_, _ = w.Write([]byte(xmlRPCFault(403, "用户名或密码错误")))
				return
			}
			_, _ = w.Write([]byte(xmlRPCResponse(`<param><value><string>` + postID + `</string></value></param>`)))
		default:
			t.Errorf("unexpected method: %q (payload: %s)", method, raw[:min(len(raw), 300)])
		}
	}))
	t.Cleanup(srv.Close)
	return NewClient(), srv, &methods
}

func paramsBlockContains(payload, needle1, needle2 string) bool {
	return strings.Contains(payload, "<string>"+needle1+"</string>") &&
		strings.Contains(payload, "<string>"+needle2+"</string>")
}

func readAll(t *testing.T, r *http.Request) string {
	t.Helper()
	buf := make([]byte, 0, r.ContentLength)
	tmp := make([]byte, 4096)
	for {
		n, err := r.Body.Read(tmp)
		buf = append(buf, tmp[:n]...)
		if err != nil {
			break
		}
	}
	return string(buf)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

var _ = xml.Unmarshal // 保留 xml 包引用位

func TestValidateCredentials_OK(t *testing.T) {
	c, srv, _ := newTestServer(t, true, "")
	if err := c.ValidateCredentials(context.Background(), creds(srv.URL)); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestValidateCredentials_BadPassword(t *testing.T) {
	c, srv, _ := newTestServer(t, false, "")
	err := c.ValidateCredentials(context.Background(), creds(srv.URL))
	if !errors.Is(err, social.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestValidateCredentials_MissingFields(t *testing.T) {
	if err := NewClient().ValidateCredentials(context.Background(), social.Credentials{"rpc_url": "https://x"}); !errors.Is(err, social.ErrInvalidCredentials) {
		t.Fatal("expected ErrInvalidCredentials for missing username/password")
	}
	if err := NewClient().ValidateCredentials(context.Background(), social.Credentials{"username": "a"}); !errors.Is(err, social.ErrInvalidCredentials) {
		t.Fatal("expected ErrInvalidCredentials for missing rpc_url")
	}
}

func TestPublish_OK_AutoBlogID(t *testing.T) {
	c, srv, methods := newTestServer(t, true, "10086")

	res, err := c.Publish(context.Background(), creds(srv.URL),
		&social.PublishInput{Title: "标题", Text: "正文第一行\n\n第二段"})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if res.ExternalID != "10086" {
		t.Errorf("unexpected external id: %s", res.ExternalID)
	}
	if res.ExternalURI != "" {
		t.Errorf("metaweblog has no permalink, got %s", res.ExternalURI)
	}
	// blog_id 未填 → 先 getUsersBlogs 再 newPost
	if len(*methods) < 2 || (*methods)[0] != "blogger.getUsersBlogs" || (*methods)[1] != "metaWeblog.newPost" {
		t.Errorf("unexpected call sequence: %v", *methods)
	}
}

func TestPublish_ExplicitBlogID_SkipsDiscovery(t *testing.T) {
	c, srv, methods := newTestServer(t, true, "20001")
	credsWithBlog := creds(srv.URL)
	credsWithBlog["blog_id"] = "blog-9"

	res, err := c.Publish(context.Background(), credsWithBlog,
		&social.PublishInput{Title: "T", Text: "B"})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if res.ExternalID != "20001" {
		t.Errorf("unexpected external id: %s", res.ExternalID)
	}
	for _, m := range *methods {
		if m == "blogger.getUsersBlogs" {
			t.Errorf("blog_id provided, discovery should be skipped: %v", *methods)
		}
	}
}

func TestFetchStats_Unsupported(t *testing.T) {
	if _, err := NewClient().FetchStats(context.Background(), creds("x"), "1", ""); !errors.Is(err, social.ErrStatsUnsupported) {
		t.Fatalf("expected ErrStatsUnsupported, got %v", err)
	}
}

func TestRegistryRegistration(t *testing.T) {
	p, err := social.GetProvider(Platform)
	if err != nil {
		t.Fatalf("metaweblog not registered via init(): %v", err)
	}
	if _, ok := p.(*Client); !ok {
		t.Fatalf("registered provider is not *Client")
	}
}
