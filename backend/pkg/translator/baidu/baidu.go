// Package baidu 提供百度通用翻译 API 的 Translator 实现，
// 与 github.com/tx7do/go-utils/translator.Translator 接口兼容。
// 百度端点在大陆网络可达（google 系端点不可达），凭据申请见
// https://fanyi-api.baidu.com/ ，标准版有免费额度。
package baidu

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const DefaultEndpoint = "https://fanyi-api.baidu.com/api/trans/vip/translate"

// bcp47ToBaidu 将 BCP47 语言标签（zh-CN/en-US/ar，编辑器使用的格式）
// 映射为百度语种码（zh/en/ara）。未收录的语种去掉地区后缀后小写透传。
var bcp47ToBaidu = map[string]string{
	"auto":   "auto",
	"zh":     "zh",
	"zh-cn":  "zh",
	"zh-tw":  "cht",
	"zh-hk":  "cht",
	"en":     "en",
	"en-us":  "en",
	"en-gb":  "en",
	"ar":     "ara",
	"ja":     "jp",
	"ja-jp":  "jp",
	"ko":     "kor",
	"ko-kr":  "kor",
	"fr":     "fra",
	"fr-fr":  "fra",
	"es":     "spa",
	"es-es":  "spa",
	"de":     "de",
	"de-de":  "de",
	"ru":     "ru",
	"ru-ru":  "ru",
	"pt":     "pt",
	"pt-br":  "pt",
	"it":     "it",
	"it-it":  "it",
	"th":     "th",
	"th-th":  "th",
	"vi":     "vie",
	"vi-vn":  "vie",
}

func NormalizeLang(lang string) string {
	l := strings.ToLower(strings.TrimSpace(lang))
	if v, ok := bcp47ToBaidu[l]; ok {
		return v
	}
	if i := strings.IndexByte(l, '-'); i > 0 {
		if v, ok := bcp47ToBaidu[l[:i]]; ok {
			return v
		}
		return l[:i]
	}
	return l
}

type Translator struct {
	appID      string
	secret     string
	endpoint   string
	httpClient *http.Client
}

type Option func(*Translator)

func WithAppID(appID string) Option {
	return func(t *Translator) { t.appID = appID }
}

func WithSecret(secret string) Option {
	return func(t *Translator) { t.secret = secret }
}

// WithEndpoint 覆盖请求端点，测试注入 mock 服务用
func WithEndpoint(endpoint string) Option {
	return func(t *Translator) { t.endpoint = endpoint }
}

func WithHTTPClient(c *http.Client) Option {
	return func(t *Translator) { t.httpClient = c }
}

func NewTranslator(opts ...Option) *Translator {
	t := &Translator{
		endpoint:   DefaultEndpoint,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
	for _, o := range opts {
		o(t)
	}
	return t
}

// Translate 翻译文本。sourceLang 传 "auto" 时由百度自动检测源语言。
func (t *Translator) Translate(source, sourceLang, targetLang string) (string, error) {
	if t.appID == "" || t.secret == "" {
		return "", errors.New("baidu translator: appid/secret is not configured")
	}
	if strings.TrimSpace(source) == "" {
		return "", nil
	}

	from := NormalizeLang(sourceLang)
	if from == "" {
		from = "auto"
	}
	to := NormalizeLang(targetLang)

	salt := strconv.FormatUint(rand.Uint64(), 10)

	form := url.Values{}
	form.Set("q", source)
	form.Set("from", from)
	form.Set("to", to)
	form.Set("appid", t.appID)
	form.Set("salt", salt)
	form.Set("sign", t.sign(source, salt))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, t.endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("baidu translator: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("baidu translator: request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 10<<20))
	if err != nil {
		return "", fmt.Errorf("baidu translator: read response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("baidu translator: unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var out struct {
		ErrorCode   string `json:"error_code"`
		ErrorMsg    string `json:"error_msg"`
		TransResult []struct {
			Src string `json:"src"`
			Dst string `json:"dst"`
		} `json:"trans_result"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return "", fmt.Errorf("baidu translator: decode response: %w", err)
	}
	if out.ErrorCode != "" {
		return "", fmt.Errorf("baidu translator: api error %s: %s", out.ErrorCode, out.ErrorMsg)
	}

	var sb strings.Builder
	for i, r := range out.TransResult {
		if i > 0 {
			sb.WriteByte('\n')
		}
		sb.WriteString(r.Dst)
	}
	return sb.String(), nil
}

// sign 生成百度签名：MD5(appid + q + salt + 密钥)
func (t *Translator) sign(q, salt string) string {
	sum := md5.Sum([]byte(t.appID + q + salt + t.secret))
	return hex.EncodeToString(sum[:])
}
