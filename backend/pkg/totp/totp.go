// Package totp 实现 RFC 6238 TOTP（SHA1/6 位/30s 步长，主流验证器 App 默认口径）。
// 仅依赖标准库：secret 为 Base32（RFC 4648 无填充）编码，otpauth URI 手工拼装。
package totp

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"crypto/subtle"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// 生成 secret 的字节数（160-bit，与主流验证器兼容）
const secretBytes = 20

// 允许的时钟偏移步数（±1 = 接受上一步/当前步/下一步，共 90s 容忍窗口）
const verifyWindow = 1

var base32Encoding = base32.StdEncoding.WithPadding(base32.NoPadding)

// GenerateSecret 生成新的 TOTP secret，返回 Base32 无填充编码字符串。
func GenerateSecret() (string, error) {
	buf := make([]byte, secretBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("read random bytes: %w", err)
	}
	return base32Encoding.EncodeToString(buf), nil
}

// ValidateSecret 校验 Base32 secret 形态合法（规范化为大写无填充形态返回）。
func ValidateSecret(secret string) (string, error) {
	normalized := strings.ToUpper(strings.ReplaceAll(strings.TrimRight(secret, "="), " ", ""))
	if _, err := base32Encoding.DecodeString(normalized); err != nil {
		return "", fmt.Errorf("invalid base32 secret: %w", err)
	}
	return normalized, nil
}

// codeAt 计算 time.Time 所在步长的 6 位验证码。
func codeAt(secret string, t time.Time) (string, error) {
	normalized, err := ValidateSecret(secret)
	if err != nil {
		return "", err
	}
	key, err := base32Encoding.DecodeString(normalized)
	if err != nil {
		return "", err
	}
	counter := uint64(t.Unix()) / 30
	var msg [8]byte
	binary.BigEndian.PutUint64(msg[:], counter)

	mac := hmac.New(sha1.New, key)
	mac.Write(msg[:])
	sum := mac.Sum(nil)

	offset := sum[len(sum)-1] & 0x0f
	bin := (uint32(sum[offset])&0x7f)<<24 |
		uint32(sum[offset+1])<<16 |
		uint32(sum[offset+2])<<8 |
		uint32(sum[offset+3])
	return fmt.Sprintf("%06d", bin%1_000_000), nil
}

// Verify 校验 6 位验证码（±1 步时钟偏移窗口；常数时间比较防时序攻击）。
func Verify(secret, code string, at time.Time) bool {
	normalized, err := ValidateSecret(secret)
	if err != nil {
		return false
	}
	code = strings.TrimSpace(code)
	if len(code) != 6 {
		return false
	}
	for i := 0; i <= verifyWindow; i++ {
		expected, err := codeAt(normalized, at.Add(time.Duration(-i)*30*time.Second))
		if err == nil && subtle.ConstantTimeCompare([]byte(expected), []byte(code)) == 1 {
			return true
		}
		expected, err = codeAt(normalized, at.Add(time.Duration(i)*30*time.Second))
		if err == nil && subtle.ConstantTimeCompare([]byte(expected), []byte(code)) == 1 {
			return true
		}
	}
	return false
}

// OTPAuthURI 生成 otpauth:// 绑定 URI（验证器 App 扫码/手动录入用）。
func OTPAuthURI(issuer, account, secret string) string {
	label := url.PathEscape(issuer) + ":" + url.PathEscape(account)
	q := url.Values{}
	q.Set("secret", secret)
	if issuer != "" {
		q.Set("issuer", issuer)
	}
	q.Set("algorithm", "SHA1")
	q.Set("digits", "6")
	q.Set("period", "30")
	return "otpauth://totp/" + label + "?" + q.Encode()
}
