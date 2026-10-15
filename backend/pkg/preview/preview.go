package preview

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// 无状态草稿预览令牌：admin BFF 签发，app（公开站点）BFF 校验后放行非已发布内容。
//
// 格式：base64url(kind|id|exp_unix) + "." + base64url(HMAC-SHA256(payload, key))
// 无 Redis 依赖（app 服务未接 Redis），过期由 exp 承担；泄露窗口=剩余有效期，
// 生产靠短 TTL（默认 15 分钟）+ HTTPS 收敛风险。

const (
	KindPost = "post"
	KindPage = "page"

	// DefaultTTL 预览链接有效期。
	DefaultTTL = 15 * time.Minute
)

// keyFromEnv 读签名密钥：优先 ENCRYPTION_KEY（与凭据加密同源，双服务同值约定）；
// 未配置时降级为项目名派生密钥并提示——仅限本地开发，生产必须配置。
func keyFromEnv() []byte {
	if k := os.Getenv("ENCRYPTION_KEY"); k != "" {
		return []byte(k)
	}
	return []byte("gowind-preview-dev:" + os.Getenv("PROJECT_NAME"))
}

func sign(payload string, key []byte) string {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

// Mint 签发预览令牌。
func Mint(kind string, id uint32, ttl time.Duration, now time.Time) (string, time.Time, error) {
	if id == 0 {
		return "", time.Time{}, fmt.Errorf("preview: target id required")
	}
	if ttl <= 0 {
		ttl = DefaultTTL
	}
	exp := now.Add(ttl)
	payload := fmt.Sprintf("%s|%d|%d", kind, id, exp.Unix())
	token := base64.RawURLEncoding.EncodeToString([]byte(payload)) + "." + sign(payload, keyFromEnv())
	return token, exp, nil
}

// Verify 校验令牌，返回目标 id。kind 不匹配 / 签名不符 / 已过期均返回错误。
func Verify(token, kind string, now time.Time) (uint32, error) {
	dot := strings.LastIndex(token, ".")
	if dot <= 0 {
		return 0, fmt.Errorf("preview: malformed token")
	}
	rawPayload, sig := token[:dot], token[dot+1:]
	payload, err := base64.RawURLEncoding.DecodeString(rawPayload)
	if err != nil {
		return 0, fmt.Errorf("preview: malformed payload")
	}
	if !hmac.Equal([]byte(sign(string(payload), keyFromEnv())), []byte(sig)) {
		return 0, fmt.Errorf("preview: signature mismatch")
	}
	parts := strings.Split(string(payload), "|")
	if len(parts) != 3 || parts[0] != kind {
		return 0, fmt.Errorf("preview: kind mismatch")
	}
	id64, err := strconv.ParseUint(parts[1], 10, 32)
	if err != nil {
		return 0, fmt.Errorf("preview: malformed id")
	}
	expUnix, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil || now.Unix() > expUnix {
		return 0, fmt.Errorf("preview: token expired")
	}
	return uint32(id64), nil
}
