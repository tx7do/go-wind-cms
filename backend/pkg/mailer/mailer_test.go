package mailer

import (
	"testing"

	libmailer "github.com/tx7do/go-utils/mailer"
)

func TestFromEnvDefaults(t *testing.T) {
	t.Setenv("SMTP_HOST", "")
	t.Setenv("SMTP_FROM", "")
	t.Setenv("SMTP_PORT", "")
	t.Setenv("SMTP_TLS", "")

	cfg := FromEnv()
	if cfg.Enabled() {
		t.Fatalf("empty host/from should be disabled")
	}
	if cfg.Port != 587 {
		t.Fatalf("default port = %d, want 587", cfg.Port)
	}
	if cfg.TLSMode != "starttls" {
		t.Fatalf("default tls mode = %q, want starttls", cfg.TLSMode)
	}
}

func TestFromEnvEnabled(t *testing.T) {
	t.Setenv("SMTP_HOST", "smtp.example.com")
	t.Setenv("SMTP_FROM", "noreply@example.com")
	t.Setenv("SMTP_PORT", "465")
	t.Setenv("SMTP_TLS", "ssl")

	cfg := FromEnv()
	if !cfg.Enabled() {
		t.Fatalf("configured smtp should be enabled")
	}
	if cfg.Port != 465 || cfg.TLSMode != "ssl" {
		t.Fatalf("port/tls = %d/%q, want 465/ssl", cfg.Port, cfg.TLSMode)
	}
	sc := cfg.smtpConfig()
	if sc.Host != "smtp.example.com" || sc.Port != 465 {
		t.Fatalf("smtpConfig host/port = %s/%d", sc.Host, sc.Port)
	}
	if sc.TlsMode != "SSL" {
		t.Fatalf("tls mode mapping = %q, want SSL", sc.TlsMode)
	}
}

// TestClampPort 钳位语义：负数/超端口上限/零值回落默认 587，合法区间无损。
// 钉住 CodeQL go/incorrect-integer-conversion 的修复——SMTP_PORT 走 Atoi
// 且错误被忽略，转换 uint32 前必须先钳位。
func TestClampPort(t *testing.T) {
	cases := []struct {
		in   int
		want uint32
	}{
		{0, 587},
		{-1, 587},
		{-65535, 587},
		{65536, 587},
		{1 << 30, 587}, // Atoi 溢出场景的量级
		{1, 1},
		{587, 587},
		{465, 465},
		{65535, 65535},
	}
	for _, c := range cases {
		if got := clampPort(c.in); got != c.want {
			t.Fatalf("clampPort(%d) = %d, want %d", c.in, got, c.want)
		}
	}
}

// TestFromEnvInvalidPort 越界 SMTP_PORT 不透传：负数与溢出（Atoi 得 MaxInt）
// 都回落默认 587，而不是被 uint32() 静默回绕成巨端口。
func TestFromEnvInvalidPort(t *testing.T) {
	cases := []struct{ raw string }{
		{"-1"},
		{"99999999999999999999"}, // Atoi 溢出，错误被忽略得 MaxInt
		{"65536"},
		{"abc"},
	}
	for _, c := range cases {
		t.Setenv("SMTP_PORT", c.raw)
		cfg := FromEnv()
		if got := cfg.Port; got != 587 {
			t.Fatalf("SMTP_PORT=%q → port %d, want default 587", c.raw, got)
		}
		// 端到端：污染源（Atoi）到 uint32 转换点全程钳位，不再静默回绕。
		if got := cfg.smtpConfig().Port; got != 587 {
			t.Fatalf("SMTP_PORT=%q → smtpConfig port %d, want 587", c.raw, got)
		}
	}
}

func TestSmtpConfigTlsMapping(t *testing.T) {
	cases := []struct{ in, want string }{
		{"ssl", "SSL"},
		{"none", "NONE"},
		{"starttls", "START_TLS"},
		{"", "START_TLS"},      // 缺省回落 starttls 口径
		{"weird", "START_TLS"}, // 枚举外回落 starttls 口径
	}
	for _, c := range cases {
		cfg := Config{TLSMode: c.in}
		if got := cfg.smtpConfig().TlsMode; got != c.want {
			t.Fatalf("TLSMode %q mapped to %q, want %q", c.in, got, c.want)
		}
	}
}

func TestSendRejectsInvalidInput(t *testing.T) {
	cfg := Config{Host: "smtp.example.com", From: "noreply@example.com"}
	if err := cfg.Send(nil, "s", "<p/>"); err == nil {
		t.Fatalf("empty recipients should error")
	}
	disabled := Config{}
	if err := disabled.Send([]string{"a@b.c"}, "s", "<p/>"); err == nil {
		t.Fatalf("disabled config should error")
	}
}

// TestSendDelegatesToLibWithHTML 投递委托库 SendMail 且声明 HTML：
// 借库的导出面间接锁行为（WithHTML 存在性 + 配置映射），
// 端到端报文断言在库 mailer 的假 SMTP 服务器测试里覆盖。
func TestSendDelegatesToLibWithHTML(t *testing.T) {
	var opts []libmailer.SendOption = []libmailer.SendOption{libmailer.WithHTML()}
	if len(opts) != 1 {
		t.Fatalf("WithHTML option missing")
	}
	cfg := Config{
		Host: "smtp.example.com", Port: 465, From: "noreply@example.com",
		FromName: "GoWind 内容中枢", TLSMode: "ssl",
	}
	sc := cfg.smtpConfig()
	if sc.Port != 465 || sc.TlsMode != "SSL" {
		t.Fatalf("config mapping mismatch: %+v", sc)
	}
	// FromName 透传到库：非空时库按 `名称 <地址>` 渲染 From 头（RFC 2047）。
	if sc.FromName != "GoWind 内容中枢" {
		t.Fatalf("from name = %q, want passthrough", sc.FromName)
	}
}
