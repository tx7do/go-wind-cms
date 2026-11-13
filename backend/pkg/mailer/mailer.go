package mailer

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	libmailer "github.com/tx7do/go-utils/mailer"
)

// 基于库 mailer（github.com/tx7do/go-utils/mailer）的薄封装：
// 本包只负责 env 配置装载与启用判定，投递统一走库 SendMail（HTML + RFC2047 主题 + Date 头）。
//
// 环境变量（生产典型姿势；管理页可视化配置优先，本兜底仅在其落空时使用）：
//
//	SMTP_HOST / SMTP_PORT（默认 587）/ SMTP_USERNAME / SMTP_PASSWORD / SMTP_FROM
//	SMTP_TLS: "starttls"（默认）| "ssl"（465 直连）| "none"
//
// 任一关键项缺失即视为未启用（Enabled()=false），调用方跳过发送并告警一次。

// Config SMTP 连接配置。
type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	TLSMode  string // starttls | ssl | none
}

// FromEnv 从环境变量读取配置。
func FromEnv() Config {
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if port == 0 {
		port = 587
	}
	mode := strings.ToLower(os.Getenv("SMTP_TLS"))
	if mode == "" {
		mode = "starttls"
	}
	return Config{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     port,
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
		From:     os.Getenv("SMTP_FROM"),
		TLSMode:  mode,
	}
}

// Enabled 判断 SMTP 是否已配置可用。
func (c Config) Enabled() bool {
	return c.Host != "" && c.From != ""
}

// smtpConfig 映射为库 SmtpConfig：本包小写 TLS 口径 → 库大写枚举。
func (c Config) smtpConfig() libmailer.SmtpConfig {
	mode := "START_TLS"
	switch strings.ToLower(c.TLSMode) {
	case "ssl":
		mode = "SSL"
	case "none":
		mode = "NONE"
	}
	return libmailer.SmtpConfig{
		Host:     c.Host,
		Port:     uint32(c.Port),
		Username: c.Username,
		Password: c.Password,
		From:     c.From,
		TlsMode:  mode,
	}
}

// Send 发送一封 HTML 邮件给多个收件人（库版 fail-fast：任一 RCPT 被拒即整体失败）。
func (c Config) Send(to []string, subject, htmlBody string) error {
	if !c.Enabled() {
		return fmt.Errorf("mailer: smtp not configured")
	}
	if len(to) == 0 {
		return fmt.Errorf("mailer: no recipients")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return libmailer.SendMail(ctx, c.smtpConfig(), to, subject, htmlBody, libmailer.WithHTML())
}
