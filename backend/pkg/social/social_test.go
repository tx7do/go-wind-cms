package social

import (
	"context"
	"errors"
	"mime"
	"mime/multipart"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── Credentials ──

func TestCredentialsGet(t *testing.T) {
	var nilCreds Credentials
	assert.Equal(t, "", nilCreds.Get("handle"), "nil map must not panic")

	creds := Credentials{"handle": "alice.bsky.social"}
	assert.Equal(t, "alice.bsky.social", creds.Get("handle"))
	assert.Equal(t, "", creds.Get("missing"))
}

func TestParseCredentials(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		wantErr bool
		want    map[string]string
	}{
		{"empty string", "", true, nil},
		{"empty json object", "{}", true, nil},
		{"invalid json", "not-json", true, nil},
		{"json array", "[]", true, nil},
		{
			"valid",
			`{"handle":"alice.bsky.social","app_password":"secret"}`,
			false,
			map[string]string{"handle": "alice.bsky.social", "app_password": "secret"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			creds, err := ParseCredentials(tt.raw)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, map[string]string(creds))
		})
	}
}

func TestMarshalCredentialsRoundTrip(t *testing.T) {
	in := Credentials{"server": "https://mastodon.social", "access_token": "tok"}
	raw, err := MarshalCredentials(in)
	require.NoError(t, err)

	out, err := ParseCredentials(raw)
	require.NoError(t, err)
	assert.Equal(t, in, out)
}

func TestMarshalCredentialsNil(t *testing.T) {
	raw, err := MarshalCredentials(nil)
	require.NoError(t, err)
	// json.Marshal(nil map) 产 "null"；ParseCredentials 会以"empty credentials"拒绝，
	// 空 map 无法经 JSON 往返成立——落库侧不应允许空凭据。
	assert.Equal(t, "null", raw)
	_, perr := ParseCredentials(raw)
	assert.Error(t, perr)
}

// ── 平台注册表 ──

// fakeProvider 注册表测试桩。
type fakeProvider struct{ inbox bool }

func (f *fakeProvider) ValidateCredentials(ctx context.Context, creds Credentials) error {
	return nil
}
func (f *fakeProvider) Publish(ctx context.Context, creds Credentials, input *PublishInput) (*PublishResult, error) {
	return &PublishResult{ExternalID: "ext-1", ExternalURI: "https://example.com/ext-1"}, nil
}
func (f *fakeProvider) FetchStats(ctx context.Context, creds Credentials, externalID, externalURI string) (*Stats, error) {
	return &Stats{Likes: 1, Replies: 2, Reposts: 3}, nil
}

// inboxFakeProvider 同时实现 InboxFetcher 的桩。
type inboxFakeProvider struct{ fakeProvider }

func (f *inboxFakeProvider) FetchInbox(ctx context.Context, creds Credentials, since time.Time) ([]InboxItem, error) {
	return []InboxItem{{ExternalID: "it-1", Type: "REPLY"}}, nil
}

func TestProviderRegistry(t *testing.T) {
	const platform = "PLATFORM_UNIT_TEST_FAKE"

	_, err := GetProvider(platform)
	require.Error(t, err, "unregistered platform must error")

	p := &fakeProvider{}
	Register(platform, p)

	got, err := GetProvider(platform)
	require.NoError(t, err)
	assert.Same(t, p, got)
}

func TestAsInboxFetcher(t *testing.T) {
	p := &fakeProvider{}
	_, ok := AsInboxFetcher(p)
	assert.False(t, ok, "provider without inbox capability must not assert")

	q := &inboxFakeProvider{}
	f, ok := AsInboxFetcher(q)
	require.True(t, ok)
	items, err := f.FetchInbox(context.Background(), nil, time.Time{})
	require.NoError(t, err)
	assert.Len(t, items, 1)
}

// ── multipart 辅助 ──

func TestMultipartRequest(t *testing.T) {
	buf, fw := NewMultipartWriter()
	require.NoError(t, WriteFilePart(fw, "photo", "cat.jpg", []byte("jpeg-bytes")))
	require.NoError(t, fw.Close())

	assert.Equal(t, UserAgent, "go-wind-cms-social/1.0")

	req, err := NewMultipartRequest(context.Background(), "POST", "https://example.com/upload", buf, fw)
	require.NoError(t, err)
	assert.Equal(t, UserAgent, req.Header.Get("User-Agent"))
	assert.Equal(t, "POST", req.Method)
	assert.Equal(t, fw.FormDataContentType(), req.Header.Get("Content-Type"))

	// 服务端按 multipart/form-data 解析，字段与内容应完整还原。
	ct := req.Header.Get("Content-Type")
	mediaType, params, err := mime.ParseMediaType(ct)
	require.NoError(t, err)
	assert.Equal(t, "multipart/form-data", mediaType)
	boundary := params["boundary"]
	require.NotEmpty(t, boundary)

	parser := multipart.NewReader(req.Body, boundary)
	form, err := parser.ReadForm(1 << 20)
	require.NoError(t, err)
	defer form.RemoveAll()

	file := form.File["photo"]
	require.Len(t, file, 1)
	fh := file[0]
	assert.Equal(t, "cat.jpg", fh.Filename)
	content := make([]byte, fh.Size)
	f, err := fh.Open()
	require.NoError(t, err)
	_, err = f.Read(content)
	require.NoError(t, err)
	assert.Equal(t, "jpeg-bytes", string(content))
}

func TestWriteFilePartEmptyData(t *testing.T) {
	buf, fw := NewMultipartWriter()
	require.NoError(t, WriteFilePart(fw, "media", "empty.bin", []byte{}))
	require.NoError(t, fw.Close())
	assert.NotEmpty(t, buf.Bytes())
}

// ── PlainTextToHTML ──

func TestPlainTextToHTML(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"empty", "", ""},
		{"whitespace only", "  \n\t ", ""},
		{"single paragraph", "hello world", "<p>hello world</p>"},
		{"two paragraphs", "first\n\nsecond", "<p>first</p><p>second</p>"},
		{"line break inside paragraph", "line1\nline2", "<p>line1<br/>line2</p>"},
		{
			"mixed paragraphs and line breaks",
			"para1-a\npara1-b\n\npara2",
			"<p>para1-a<br/>para1-b</p><p>para2</p>",
		},
		{"blank block skipped", "a\n\n\n\nb", "<p>a</p><p>b</p>"},
		{"html escaped", "<script>alert(1)</script>", "<p>&lt;script&gt;alert(1)&lt;/script&gt;</p>"},
		{"entity escaped", "a & b", "<p>a &amp; b</p>"},
		{"quote escaped", `say "hi"`, "<p>say &#34;hi&#34;</p>"},
		{"leading trailing trimmed", "  spaced  ", "<p>spaced</p>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, PlainTextToHTML(tt.in))
		})
	}
}

// ── 哨兵错误 ──

func TestSentinelErrors(t *testing.T) {
	// core service 按哨兵错误映射 proto 错误码，errors.Is 语义必须成立。
	err := errors.Join(ErrRateLimited)
	assert.True(t, errors.Is(err, ErrRateLimited))
	assert.ErrorIs(t, ErrStatsUnsupported, ErrStatsUnsupported)
	assert.NotErrorIs(t, ErrNotFound, ErrInvalidCredentials)
}

// ── Benchmarks ──

func BenchmarkPlainTextToHTML(b *testing.B) {
	long := strings.Repeat("这是第一段内容，包含中文。\n\n", 20)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = PlainTextToHTML(long)
	}
}

func BenchmarkMarshalCredentials(b *testing.B) {
	creds := Credentials{"handle": "alice.bsky.social", "app_password": "secret", "pds_url": "https://bsky.social"}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = MarshalCredentials(creds)
	}
}

func BenchmarkParseCredentials(b *testing.B) {
	raw := `{"handle":"alice.bsky.social","app_password":"secret","pds_url":"https://bsky.social"}`
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = ParseCredentials(raw)
	}
}
