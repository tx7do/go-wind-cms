package baidu

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeLang(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"zh-CN", "zh"},
		{"en-US", "en"},
		{"ar", "ara"},
		{"zh-TW", "cht"},
		{"auto", "auto"},
		{"", ""},
		{"fr-FR", "fra"},
		{"ja-JP", "jp"},
		{"ko-KR", "kor"},
		{"pt-BR", "pt"},
		{"xxx-YY", "xxx"},
	}
	for _, c := range cases {
		assert.Equal(t, c.want, NormalizeLang(c.in), "input: %q", c.in)
	}
}

func TestTranslate(t *testing.T) {
	var gotForm url.Values

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		gotForm = r.PostForm
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"from": "en",
			"to":   "zh",
			"trans_result": []map[string]string{
				{"src": "hello", "dst": "你好"},
				{"src": "world", "dst": "世界"},
			},
		})
	}))
	defer srv.Close()

	tr := NewTranslator(
		WithAppID("appid123"),
		WithSecret("secret456"),
		WithEndpoint(srv.URL),
	)

	out, err := tr.Translate("hello\nworld", "auto", "zh-CN")
	require.NoError(t, err)
	assert.Equal(t, "你好\n世界", out)

	require.NotNil(t, gotForm, "服务端未收到请求")
	assert.Equal(t, "hello\nworld", gotForm.Get("q"))
	assert.Equal(t, "auto", gotForm.Get("from"))
	assert.Equal(t, "zh", gotForm.Get("to"))
	assert.Equal(t, "appid123", gotForm.Get("appid"))

	salt := gotForm.Get("salt")
	require.NotEmpty(t, salt, "salt 缺失")
	sum := md5.Sum([]byte("appid123" + "hello\nworld" + salt + "secret456"))
	assert.Equal(t, hex.EncodeToString(sum[:]), gotForm.Get("sign"), "签名应为 MD5(appid+q+salt+secret)")
}

func TestTranslateLanguageMapping(t *testing.T) {
	var gotForm url.Values

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		gotForm = r.PostForm
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"trans_result": []map[string]string{{"src": "x", "dst": "y"}},
		})
	}))
	defer srv.Close()

	tr := NewTranslator(WithAppID("a"), WithSecret("b"), WithEndpoint(srv.URL))

	_, err := tr.Translate("x", "en-US", "ar")
	require.NoError(t, err)
	assert.Equal(t, "en", gotForm.Get("from"))
	assert.Equal(t, "ara", gotForm.Get("to"))
}

func TestTranslateAPIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"error_code": "54003", "error_msg": "Too fast"})
	}))
	defer srv.Close()

	tr := NewTranslator(WithAppID("a"), WithSecret("b"), WithEndpoint(srv.URL))

	_, err := tr.Translate("x", "auto", "zh-CN")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "54003")
}

func TestTranslateNoCredentials(t *testing.T) {
	tr := NewTranslator()
	_, err := tr.Translate("x", "auto", "zh-CN")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not configured")
}

func TestTranslateEmptyInput(t *testing.T) {
	tr := NewTranslator(WithAppID("a"), WithSecret("b"))
	out, err := tr.Translate("", "auto", "zh-CN")
	require.NoError(t, err)
	assert.Empty(t, out)
}
