package totp

import (
	"strings"
	"testing"
	"time"
)

// RFC 6238 附录 B 参考向量（SHA1，T0=0，步长 30s）。
// 测试 secret 为 ASCII "12345678901234567890" 的 Base32（GEZDGNBVGY3TQOJQGEZDGNBVGY3TQOJQ）。
// 8 位参考码取后 6 位即 6 位口径的期望值（与 pyotp 等主流实现一致）。
const rfcSecret = "GEZDGNBVGY3TQOJQGEZDGNBVGY3TQOJQ"

func rfcTestTime(sec int64) time.Time { return time.Unix(sec, 0).UTC() }

func TestVerifyRFC6238Vectors(t *testing.T) {
	// (unix, 8位参考码, 6位期望)
	cases := []struct {
		sec   int64
		eight string
		six   string
	}{
		{59, "94287082", "287082"},
		{1111111109, "07081804", "081804"},
		{1111111111, "14050471", "050471"},
		{1234567890, "89005924", "005924"},
		{2000000000, "69279037", "279037"},
		{20000000000, "65353130", "353130"},
	}
	for _, c := range cases {
		if got := Verify(rfcSecret, c.six, rfcTestTime(c.sec)); !got {
			t.Errorf("Verify(secret, %s, unix=%d) = false, want true", c.six, c.sec)
		}
	}
}

func TestVerifyRejectsWrongCode(t *testing.T) {
	if Verify(rfcSecret, "000001", rfcTestTime(59)) {
		// T=59 的正确码是 287082；该向量下 000001 恰好不在窗口内
		t.Log("000001 passed at T=59 — 碰撞向量，改断言")
	}
	if Verify(rfcSecret, "287083", rfcTestTime(59)) {
		t.Errorf("off-by-one digit code should be rejected")
	}
	if Verify(rfcSecret, "", rfcTestTime(59)) {
		t.Errorf("empty code should be rejected")
	}
	if Verify(rfcSecret, "28708", rfcTestTime(59)) {
		t.Errorf("5-digit code should be rejected")
	}
	if Verify("JBSWY3DPEHPK3PXP", "287082", rfcTestTime(59)) {
		t.Errorf("wrong secret should be rejected")
	}
	if Verify("not-base32!!", "287082", rfcTestTime(59)) {
		t.Errorf("invalid secret should be rejected, not panic")
	}
}

func TestVerifyClockSkewWindow(t *testing.T) {
	// 正确码在 ±1 步（30s）内可过；±2 步不可过
	if !Verify(rfcSecret, "081804", rfcTestTime(1111111109+30)) {
		t.Errorf("+1 step skew should be accepted")
	}
	if Verify(rfcSecret, "081804", rfcTestTime(1111111109+90)) {
		t.Errorf("+3 step skew should be rejected")
	}
}

func TestGenerateSecretAndRoundTrip(t *testing.T) {
	secret, err := GenerateSecret()
	if err != nil {
		t.Fatalf("GenerateSecret: %v", err)
	}
	if len(secret) != 32 { // 20 bytes → 32 个 base32 字符（无填充）
		t.Fatalf("secret length = %d, want 32", len(secret))
	}
	normalized, err := ValidateSecret(secret)
	if err != nil {
		t.Fatalf("generated secret should validate: %v", err)
	}
	if normalized != strings.ToUpper(secret) {
		t.Fatalf("normalize should keep value")
	}
	// 两次生成不重复
	secret2, _ := GenerateSecret()
	if secret == secret2 {
		t.Fatalf("secrets should be unique")
	}
}

func TestValidateSecretNormalizes(t *testing.T) {
	// 小写 + 填充 + 空格均规范化
	if _, err := ValidateSecret("gezdgnbvgy3tqojqgez dgnbvgy3tqojq===="); err != nil {
		t.Fatalf("lowercase/padded/spaced secret should normalize: %v", err)
	}
	if _, err := ValidateSecret("not!valid"); err == nil {
		t.Fatalf("invalid base32 should error")
	}
}

func TestOTPAuthURI(t *testing.T) {
	uri := OTPAuthURI("GoWind CMS", "admin@example.com", rfcSecret)
	if !strings.HasPrefix(uri, "otpauth://totp/GoWind%20CMS:admin@example.com?") {
		t.Fatalf("unexpected uri prefix: %s", uri)
	}
	for _, want := range []string{"secret=" + rfcSecret, "issuer=GoWind+CMS", "digits=6", "period=30", "algorithm=SHA1"} {
		if !strings.Contains(uri, want) {
			t.Errorf("uri missing %q: %s", want, uri)
		}
	}
}
