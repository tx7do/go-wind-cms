package logging

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	authenticationV1 "go-wind-cms/api/gen/go/authentication/service/v1"
)

// ── getClientRealIP ──

func TestGetClientRealIP(t *testing.T) {
	tests := []struct {
		name       string
		xff        string
		xri        string
		remoteAddr string
		want       string
	}{
		{"xff first valid wins", "203.0.113.9, 70.41.3.1", "", "192.168.1.1:1234", "203.0.113.9"},
		{"xff skips invalid entries", "garbage, 70.41.3.1", "", "", "70.41.3.1"},
		{"xff trimmed", "  70.41.3.1  ", "", "", "70.41.3.1"},
		{"xff all invalid falls through", "garbage, still-bad", "198.51.100.7", "", "198.51.100.7"},
		{"xff empty falls to xri", "", "198.51.100.7", "10.0.0.1:99", "198.51.100.7"},
		{"xri invalid falls to remote", "", "not-an-ip", "203.0.113.5:80", "203.0.113.5"},
		{"no headers uses remote with port", "", "", "203.0.113.5:80", "203.0.113.5"},
		{"ipv6 remote", "", "", "[2001:db8::1]:443", "2001:db8::1"},
		{"remote without port", "", "", "203.0.113.5", "203.0.113.5"},
		{"nothing resolvable", "", "", "not-an-addr", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.xff != "" {
				req.Header.Set("X-Forwarded-For", tt.xff)
			}
			if tt.xri != "" {
				req.Header.Set("X-Real-IP", tt.xri)
			}
			req.RemoteAddr = tt.remoteAddr
			assert.Equal(t, tt.want, getClientRealIP(req))
		})
	}
}

func TestClientRealIPEportedAlias(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "203.0.113.5:80"
	assert.Equal(t, getClientRealIP(req), ClientRealIP(req))
}

func TestGetClientRealIPNilRequest(t *testing.T) {
	assert.Equal(t, "", getClientRealIP(nil))
}

// ── getRequestId ──

func TestGetRequestId(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Request-Id", "rid-1")
	assert.Equal(t, "rid-1", getRequestId(req))

	req.Header.Set("X-Correlation-Id", "cid-1")
	req.Header.Del("X-Request-Id")
	assert.Equal(t, "cid-1", getRequestId(req))

	req.Header.Set("X-Fc-Request-Id", "fcid-1")
	req.Header.Del("X-Correlation-Id")
	assert.Equal(t, "fcid-1", getRequestId(req))

	// 无头：生成新 GUIDv4（32 位 hex，无连字符）。语义钉测：每次调用都
	// 生成新 ID——中间件对无 ID 请求必须各自独立取值。
	req.Header.Del("X-Fc-Request-Id")
	generated := getRequestId(req)
	assert.NotEmpty(t, generated)
	assert.Len(t, generated, 32)
	assert.NotEqual(t, generated, getRequestId(req), "每次调用生成全新 ID")
}

func TestGetRequestIdNilRequest(t *testing.T) {
	assert.Equal(t, "", getRequestId(nil))
}

// ── getClientID ──

func TestGetClientID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	// 无头无 token 返回空串。
	assert.Equal(t, "", getClientID(req, nil))

	// token clientId 作回退。
	tok := &authenticationV1.UserTokenPayload{ClientId: stringPtr("from-token")}
	assert.Equal(t, "from-token", getClientID(req, tok))

	// 头优先于 token。
	req.Header.Set("X-Client-IP", "client-header")
	assert.Equal(t, "client-header", getClientID(req, tok))
}

func stringPtr(s string) *string { return &s }

// ── getStatusCode ──

func TestGetStatusCode(t *testing.T) {
	code, reason, ok := getStatusCode(nil)
	assert.Equal(t, uint32(200), code)
	assert.Equal(t, "", reason)
	assert.True(t, ok)

	forbidden := errors.Forbidden("TEST_FORBIDDEN", "no")
	code, reason, ok = getStatusCode(forbidden)
	assert.Equal(t, uint32(403), code)
	assert.Equal(t, "TEST_FORBIDDEN", reason)
	assert.False(t, ok, "4xx 是失败请求")

	serverErr := errors.InternalServer("TEST_INTERNAL", "boom")
	code, _, ok = getStatusCode(serverErr)
	assert.Equal(t, uint32(500), code)
	assert.False(t, ok)

	unauthorized := errors.Unauthorized("UNAUTHORIZED", "no token")
	code, _, ok = getStatusCode(unauthorized)
	assert.Equal(t, uint32(401), code)
	assert.False(t, ok)
}

// ── parseUsernameFromBytes / extractUsernameFromRequest ──

func TestParseUsernameFromBytes(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		want    string
		wantErr bool
	}{
		{"json", `{"username":"alice","password":"x"}`, "alice", false},
		{"json spaces", `{"username" :  "bob" }`, "bob", false},
		{"json empty username no match", `{"username":""}`, "", true},
		{"form", "username=carol&password=x", "carol", false},
		{"form encoded", "username=dave%21&x=1", "dave!", false},
		{"neither", `{"user":"eve"}`, "", true},
		{"empty body", "", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseUsernameFromBytes([]byte(tt.body))
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExtractUsernameFromRequestRestoresBody(t *testing.T) {
	body := `{"username":"alice"}`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))

	got, err := extractUsernameFromRequest(req)
	require.NoError(t, err)
	assert.Equal(t, "alice", got)

	// Body 必须被恢复：后续（proto 绑定）还能读到。
	rest := make([]byte, len(body))
	n, _ := req.Body.Read(rest)
	assert.Equal(t, body, string(rest[:n]), "body must be restored for downstream binding")
}

func TestExtractUsernameFromRequestNil(t *testing.T) {
	_, err := extractUsernameFromRequest(nil)
	assert.Error(t, err)
}

// ── tokenHash ──

func TestTokenHash(t *testing.T) {
	h := tokenHash("secret-token")
	// SHA-256 hex 摘要：64 个十六进制字符。
	assert.Len(t, h, 64)
	assert.NotContains(t, h, "secret-token", "不得泄漏明文")

	sum := sha256.Sum256([]byte("secret-token"))
	assert.Equal(t, hex.EncodeToString(sum[:]), h)

	assert.Equal(t, h, tokenHash("secret-token"), "同输入稳定")
	assert.NotEqual(t, h, tokenHash("other"), "异输入不同")
}

// ── generateECDSAKeyPair + encodeDER ──

func TestEncodeDERSignatureVerifies(t *testing.T) {
	priv, pub, err := generateECDSAKeyPair()
	require.NoError(t, err)

	// encodeDER 的产物必须能被标准库 VerifyASN1 接受（属性测试，多轮随机签名）。
	for i := 0; i < 20; i++ {
		digest := sha256.Sum256([]byte("payload-" + string(rune('a'+i))))
		r, s, err := ecdsa.Sign(rand.Reader, priv, digest[:])
		require.NoError(t, err)

		der, err := encodeDER(r, s)
		require.NoError(t, err)
		assert.True(t, ecdsa.VerifyASN1(pub, digest[:], der),
			"DER-encoded signature must verify (round %d)", i)
	}
}

func TestEncodeDERHighBits(t *testing.T) {
	// r/s 首字节带 0x80 时需补 0x00 前缀（DER 正整数规则），产出仍可验签。
	priv, pub, err := generateECDSAKeyPair()
	require.NoError(t, err)

	digest := sha256.Sum256([]byte("high-bit"))
	for i := 0; i < 50; i++ {
		r, s, err := ecdsa.Sign(rand.Reader, priv, digest[:])
		require.NoError(t, err)
		der, err := encodeDER(r, s)
		require.NoError(t, err)
		assert.True(t, ecdsa.VerifyASN1(pub, digest[:], der))
	}
}

// ── isPrivateIP ──

func TestIsPrivateIP(t *testing.T) {
	private := []string{
		"10.0.0.1", "10.255.255.255",
		"172.16.0.1", "172.31.255.255",
		"192.168.1.1",
		"127.0.0.1",
		"169.254.1.1",
		"::1",
		"fc00::1", "fd12::1",
		" 192.168.1.1 ",
	}
	for _, ip := range private {
		assert.True(t, isPrivateIP(ip), "%s should be private", ip)
	}

	public := []string{
		"8.8.8.8",
		"172.32.0.1", // 172.32 已出私网段
		"172.15.255.255",
		"192.169.0.1",
		"2001:db8::1",
		"",
		"garbage",
	}
	for _, ip := range public {
		assert.False(t, isPrivateIP(ip), "%s should not be private", ip)
	}
}

// ── detectPlatformFromUA ──

func TestDetectPlatformFromUA(t *testing.T) {
	tests := []struct {
		name string
		ua   string
		want string
	}{
		{"empty", "", PlatformOther},
		{"okhttp android app", "okhttp/4.9.0", PlatformAndroidApp},
		{"android webview", "Mozilla/5.0 (Linux; Android 13; SM-G991B; wv) AppleWebKit/537.36", PlatformAndroidApp},
		{"android package", "MyApp/1.0 (Android; com.example.app)", PlatformAndroidApp},
		{"ios native", "MyApp/2.1 CFNetwork/1404.0.5 Darwin/22.2.0", PlatformiOSApp},
		{"iphone", "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X)", PlatformiOSApp},
		{"ipad", "Mozilla/5.0 (iPad; CPU OS 16_0)", PlatformiOSApp},
		{"electron windows", "Mozilla/5.0 (Windows NT 10.0) Electron/22.0", PlatformDesktopWindows},
		{"desktop mac", "MyDesktopApp/1.0 (Macintosh; macOS 13)", PlatformDesktopMac},
		{"desktop linux", "DesktopApp/1.0 (X11; Ubuntu Linux)", PlatformDesktopLinux},
		{"web chrome windows", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/120.0", PlatformWeb},
		{"web mac", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) Safari/605.1.15", PlatformWeb},
		{"web android chrome", "Mozilla/5.0 (Linux; Android 13) Chrome/120 Mobile", PlatformWeb},
		{"unknown bot", "python-requests/2.31", PlatformOther},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, detectPlatformFromUA(tt.ua))
		})
	}
}

// ── Benchmarks ──

func BenchmarkGetClientRealIP(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Forwarded-For", "203.0.113.9, 70.41.3.1")
	req.RemoteAddr = "192.168.1.1:1234"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = getClientRealIP(req)
	}
}

func BenchmarkDetectPlatformFromUA(b *testing.B) {
	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = detectPlatformFromUA(ua)
	}
}

func BenchmarkParseUsernameFromBytesJSON(b *testing.B) {
	body := []byte(`{"username":"alice","password":"x","captchaId":"id","captchaValue":"v"}`)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = parseUsernameFromBytes(body)
	}
}

func BenchmarkTokenHash(b *testing.B) {
	token := strings.Repeat("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.", 6)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = tokenHash(token)
	}
}

func BenchmarkEncodeDER(b *testing.B) {
	priv, _, err := generateECDSAKeyPair()
	if err != nil {
		b.Fatal(err)
	}
	digest := sha256.Sum256([]byte("bench"))
	r, s, err := ecdsa.Sign(rand.Reader, priv, digest[:])
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = encodeDER(r, s)
	}
}
