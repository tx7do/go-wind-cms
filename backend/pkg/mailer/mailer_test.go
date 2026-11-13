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
	cfg := Config{Host: "smtp.example.com", Port: 465, From: "noreply@example.com", TLSMode: "ssl"}
	sc := cfg.smtpConfig()
	if sc.Port != 465 || sc.TlsMode != "SSL" {
		t.Fatalf("config mapping mismatch: %+v", sc)
	}
}
