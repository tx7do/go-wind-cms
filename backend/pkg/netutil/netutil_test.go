package netutil

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/go-kratos/kratos/v2/transport"
)

// fakeTransporter 测试用 transport.Transporter（含 HTTP 扩展面），
// kratos 自带的 http.Transport 字段全未导出、无法从包外构造。
type fakeTransporter struct {
	headers      http.Header
	request      *http.Request
	operation    string
	pathTemplate string
}

func (f *fakeTransporter) Kind() transport.Kind      { return transport.KindHTTP }
func (f *fakeTransporter) Endpoint() string          { return "http://localhost" }
func (f *fakeTransporter) Operation() string         { return f.operation }
func (f *fakeTransporter) RequestHeader() transport.Header {
	return headerCarrier(f.headers)
}
func (f *fakeTransporter) ReplyHeader() transport.Header { return headerCarrier(http.Header{}) }
func (f *fakeTransporter) Request() *http.Request        { return f.request }
func (f *fakeTransporter) PathTemplate() string          { return f.pathTemplate }

type headerCarrier http.Header

func (h headerCarrier) Get(key string) string        { return http.Header(h).Get(key) }
func (h headerCarrier) Set(key, value string)        { http.Header(h).Set(key, value) }
func (h headerCarrier) Add(key, value string)        { http.Header(h).Add(key, value) }
func (h headerCarrier) Keys() []string               { return []string{} }
func (h headerCarrier) Values(key string) []string   { return http.Header(h).Values(key) }

// newTransportContext 把 fake transporter 挂到 server context。
func newTransportContext(t *testing.T, req *http.Request) context.Context {
	t.Helper()
	tr := &fakeTransporter{
		headers:      req.Header,
		request:      req,
		operation:    "/test/Operation",
		pathTemplate: "/test/{id}",
	}
	return transport.NewServerContext(context.Background(), tr)
}

// ── ValidateURL ──

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		rawURL  string
		wantErr bool
	}{
		{"empty", "", true},
		{"invalid parse", "://bad", true},
		{"http ok", "http://example.com/page", false},
		{"https ok", "https://example.com/page?q=1", false},
		{"ftp scheme", "ftp://example.com/file", true},
		{"file scheme", "file:///etc/passwd", true},
		{"javascript scheme", "javascript:alert(1)", true},
		{"missing host", "http:///path", true},
		{"userinfo denied", "http://user:pass@example.com/", true},
		{"userinfo user only", "http://user@example.com/", true},
		{"port ok", "https://example.com:8443/", false},
		{"ipv4 ok", "http://93.184.216.34/", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := ValidateURL(tt.rawURL)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("ValidateURL(%q) error = nil, want error", tt.rawURL)
				}
				return
			}
			if err != nil {
				t.Fatalf("ValidateURL(%q) unexpected error: %v", tt.rawURL, err)
			}
			if u == nil {
				t.Fatal("parsed url must not be nil on success")
			}
		})
	}
}

// ── LimitReader ──

func TestLimitReader(t *testing.T) {
	all, err := io.ReadAll(LimitReader(strings.NewReader("0123456789")))
	if err != nil {
		t.Fatal(err)
	}
	if string(all) != "0123456789" {
		t.Fatalf("got %q, want full payload", all)
	}

	all2, _ := io.ReadAll(LimitReader(strings.NewReader(strings.Repeat("x", 20))))
	if len(all2) != 20 {
		t.Fatalf("short payload should not be truncated, got %d", len(all2))
	}
}

// ── isBlockedIP（SSRF 保留段黑名单）──

func TestIsBlockedIP(t *testing.T) {
	cases := []struct {
		ip      string
		blocked bool
	}{
		{"127.0.0.1", true},       // loopback
		{"10.1.2.3", true},        // private
		{"192.168.1.1", true},     // private
		{"172.16.0.9", true},      // private
		{"169.254.169.254", true}, // cloud metadata / link-local
		{"0.0.0.0", true},         // this-network
		{"100.64.1.2", true},      // CGNAT
		{"::1", true},             // ipv6 loopback
		{"fe80::1", true},         // ipv6 link-local
		{"fc00::1", true},         // ipv6 ULA
		{"192.0.2.1", true},       // TEST-NET-1
		{"198.51.100.7", true},    // TEST-NET-2
		{"203.0.113.9", true},     // TEST-NET-3
		{"8.8.8.8", false},        // public
		{"1.1.1.1", false},        // public
		{"93.184.216.34", false},  // public
		{"2606:4700::1111", false},
	}
	for _, c := range cases {
		ip := net.ParseIP(c.ip)
		if ip == nil {
			t.Fatalf("bad test ip %q", c.ip)
		}
		if got := isBlockedIP(ip); got != c.blocked {
			t.Errorf("isBlockedIP(%s) = %v, want %v", c.ip, got, c.blocked)
		}
	}

	if !isBlockedIP(nil) {
		t.Error("nil IP must be treated as blocked")
	}
}

// ── SafeHTTPClient：出站拨号必须拒绝内网/保留地址 ──

func TestSafeHTTPClientBlocksLoopback(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	client := SafeHTTPClient()
	// httptest 服务必然监听在 127.0.0.1，SSRF 防护必须拦截。
	resp, err := client.Get(srv.URL)
	if err == nil {
		resp.Body.Close()
		t.Fatal("fetch to loopback must be blocked, got success")
	}
	if !strings.Contains(err.Error(), "blocked") {
		t.Fatalf("error should mention blocked, got: %v", err)
	}
}

func TestSafeHTTPClientTimeout(t *testing.T) {
	client := SafeHTTPClient()
	if client.Timeout != 30*time.Second {
		t.Fatalf("client timeout = %v, want 30s", client.Timeout)
	}
	if client.Transport == nil {
		t.Fatal("client transport must be set")
	}
}

// ── HeaderFromContext ──

func TestHeaderFromContext(t *testing.T) {
	if h := HeaderFromContext(context.Background()); h != nil {
		t.Errorf("plain context should yield nil header, got %v", h)
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Captcha-Id", "cap-1")
	ctx := newTransportContext(t, req)

	h := HeaderFromContext(ctx)
	if h == nil {
		t.Fatal("http transport context should yield request header")
	}
	if got := h.Get("X-Captcha-Id"); got != "cap-1" {
		t.Errorf("header value = %q, want %q", got, "cap-1")
	}
}

// ── ParseUnverifiedBearerJWTClaims ──

// signTestToken 签一个仅供解析（不验签）用的 HS256 token。
func signTestToken(t *testing.T, claims jwt.MapClaims) string {
	t.Helper()
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := tok.SignedString([]byte("unit-test-secret"))
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	return s
}

func newBearerContext(t *testing.T, authorization string) context.Context {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	if authorization != "" {
		req.Header.Set("Authorization", authorization)
	}
	return newTransportContext(t, req)
}

func TestParseUnverifiedBearerJWTClaims(t *testing.T) {
	tests := []struct {
		name    string
		authz   string
		wantUID int64
		wantJTI string
		wantErr bool
	}{
		{
			name:    "no header",
			authz:   "",
			wantErr: true,
		},
		{
			name:    "not bearer",
			authz:   "Basic dXNlcjpwYXNz",
			wantErr: true,
		},
		{
			name:    "bare bearer word",
			authz:   "Bearer ",
			wantErr: true,
		},
		{
			name:    "garbage token",
			authz:   "Bearer not-a-jwt",
			wantErr: true,
		},
		{
			name:    "missing jti",
			authz:   "Bearer " + signTestToken(t, jwt.MapClaims{"uid": 42}),
			wantErr: true,
		},
		{
			name:    "missing uid",
			authz:   "Bearer " + signTestToken(t, jwt.MapClaims{"jti": "abc"}),
			wantErr: true,
		},
		{
			name:    "uid zero",
			authz:   "Bearer " + signTestToken(t, jwt.MapClaims{"uid": 0, "jti": "abc"}),
			wantErr: true,
		},
		{
			name:    "valid",
			authz:   "Bearer " + signTestToken(t, jwt.MapClaims{"uid": 42, "jti": "jti-xyz"}),
			wantUID: 42,
			wantJTI: "jti-xyz",
		},
		{
			name:    "valid expired token is still parseable",
			authz:   "Bearer " + signTestToken(t, jwt.MapClaims{"uid": 7, "jti": "jti-1", "exp": time.Now().Add(-time.Hour).Unix()}),
			wantUID: 7,
			wantJTI: "jti-1",
		},
		{
			name:    "uid as string is rejected (numeric expected)",
			authz:   "Bearer " + signTestToken(t, jwt.MapClaims{"uid": "42", "jti": "abc"}),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := newBearerContext(t, tt.authz)
			uid, jti, err := ParseUnverifiedBearerJWTClaims(ctx)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got uid=%d jti=%q", uid, jti)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if uid != tt.wantUID || jti != tt.wantJTI {
				t.Fatalf("got uid=%d jti=%q, want uid=%d jti=%q", uid, jti, tt.wantUID, tt.wantJTI)
			}
		})
	}
}

func TestParseUnverifiedBearerJWTClaimsNoTransport(t *testing.T) {
	if _, _, err := ParseUnverifiedBearerJWTClaims(context.Background()); err == nil {
		t.Fatal("expected error without transport in context")
	}
}

// ── Benchmarks ──

func BenchmarkValidateURL(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = ValidateURL("https://example.com/content/post/hello-world?q=1&lang=zh-CN")
	}
}

func BenchmarkIsBlockedIP(b *testing.B) {
	ip := net.ParseIP("93.184.216.34")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = isBlockedIP(ip)
	}
}
