// Package password 提供等保要求的口令策略：复杂度、有效期、历史口令。
// 阈值统一自 sys_config 平台参数读取（参数管理页可调），由调用方传入本包——
// 本包只做纯校验，不依赖存储与配置源（单测无需桩）。
package password

import (
	"errors"
	"strconv"
	"unicode"
)

// ErrWeakPassword 复杂度不达标（调用方转成 4xx 响应）。
var ErrWeakPassword = errors.New("password does not meet complexity requirements: at least 3 of (lowercase, uppercase, digit, symbol)")

// ValidateComplexity 校验明文口令复杂度：长度达标且至少包含
// 小写/大写/数字/符号 四类中的三类。minLen 由调用方自 sys_config 读取。
func ValidateComplexity(plain string, minLen int) error {
	if minLen <= 0 {
		minLen = 8
	}
	if len(plain) < minLen {
		return errors.New("password too short: minimum length is " + strconv.Itoa(minLen))
	}
	var hasLower, hasUpper, hasDigit, hasSymbol bool
	for _, r := range plain {
		switch {
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsDigit(r):
			hasDigit = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			hasSymbol = true
		}
	}
	classes := 0
	for _, b := range []bool{hasLower, hasUpper, hasDigit, hasSymbol} {
		if b {
			classes++
		}
	}
	if classes < 3 {
		return ErrWeakPassword
	}
	return nil
}
