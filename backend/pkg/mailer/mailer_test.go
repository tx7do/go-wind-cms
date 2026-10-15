package mailer

import (
	"strings"
	"testing"
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
	if cfg.addr() != "smtp.example.com:465" {
		t.Fatalf("addr = %q", cfg.addr())
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

func TestBuildMessage(t *testing.T) {
	msg := string(buildMessage("noreply@example.com", []string{"u@example.com"}, "通知 标题", "<p>正文</p>"))
	if !strings.Contains(msg, "From: noreply@example.com") {
		t.Fatalf("missing From header")
	}
	if !strings.Contains(msg, "To: u@example.com") {
		t.Fatalf("missing To header")
	}
	if !strings.Contains(msg, "=?UTF-8?B?") {
		t.Fatalf("subject should be RFC 2047 encoded")
	}
	if !strings.Contains(msg, "text/html; charset=UTF-8") {
		t.Fatalf("missing html content type")
	}
	if !strings.HasSuffix(msg, "<p>正文</p>") {
		t.Fatalf("body mismatch")
	}
}
