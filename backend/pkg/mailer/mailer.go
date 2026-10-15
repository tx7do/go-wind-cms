package mailer

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"
)

// 基于 stdlib net/smtp 的最小 SMTP 出站客户端。
//
// 配置全部来自环境变量（生产典型姿势；管理页可视化配置留待后续引入
// provider 表时迁移）：
//
//	SMTP_HOST / SMTP_PORT（默认 587）/ SMTP_USERNAME / SMTP_PASSWORD / SMTP_FROM
//	SMTP_TLS: "starttls"（默认）| "ssl"（465 直连）| "none"（仅内网调试）
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

// addr host:port。
func (c Config) addr() string {
	return net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
}

// Send 发送一封 HTML 邮件给多个收件人。
// 返回 (已接受收件人, error)：SMTP 层拒绝部分收件人时不视为整体失败。
func (c Config) Send(to []string, subject, htmlBody string) error {
	if !c.Enabled() {
		return fmt.Errorf("mailer: smtp not configured")
	}
	if len(to) == 0 {
		return fmt.Errorf("mailer: no recipients")
	}

	msg := buildMessage(c.From, to, subject, htmlBody)

	var auth smtp.Auth
	if c.Username != "" {
		auth = smtp.PlainAuth("", c.Username, c.Password, c.Host)
	}

	switch c.TLSMode {
	case "ssl":
		return sendOverSSL(c, auth, to, msg)
	case "none":
		return smtp.SendMail(c.addr(), auth, c.From, to, msg)
	default: // starttls
		return smtp.SendMail(c.addr(), auth, c.From, to, msg)
	}
}

// sendOverSSL 走隐式 TLS（465 端口姿势）。stdlib smtp.SendMail 不支持，
// 手动拨 TLS 连接后复用 smtp.NewClient。
func sendOverSSL(c Config, auth smtp.Auth, to []string, msg []byte) error {
	conn, err := tls.Dial("tcp", c.addr(), &tls.Config{ServerName: c.Host})
	if err != nil {
		return fmt.Errorf("mailer: tls dial failed: %w", err)
	}
	defer func() { _ = conn.Close() }()

	client, err := smtp.NewClient(conn, c.Host)
	if err != nil {
		return fmt.Errorf("mailer: smtp client failed: %w", err)
	}
	defer func() { _ = client.Close() }()

	if ok, _ := client.Extension("AUTH"); ok && auth != nil {
		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("mailer: smtp auth failed: %w", err)
		}
	}
	if err = client.Mail(c.From); err != nil {
		return fmt.Errorf("mailer: MAIL FROM failed: %w", err)
	}
	for _, rcpt := range to {
		if err = client.Rcpt(rcpt); err != nil {
			continue // 单收件人被拒不阻断其余
		}
	}
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("mailer: DATA failed: %w", err)
	}
	if _, err = w.Write(msg); err != nil {
		return fmt.Errorf("mailer: write body failed: %w", err)
	}
	if err = w.Close(); err != nil {
		return fmt.Errorf("mailer: close data failed: %w", err)
	}
	return client.Quit()
}

// buildMessage 拼 RFC 5322 报文（UTF-8 + HTML）。
func buildMessage(from string, to []string, subject, htmlBody string) []byte {
	var b strings.Builder
	b.WriteString("From: " + from + "\r\n")
	b.WriteString("To: " + strings.Join(to, ", ") + "\r\n")
	b.WriteString("Subject: =?UTF-8?B?" + base64Subject(subject) + "?=\r\n")
	b.WriteString("MIME-Version: 1.0\r\n")
	b.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	b.WriteString("Date: " + time.Now().Format(time.RFC1123Z) + "\r\n")
	b.WriteString("\r\n")
	b.WriteString(htmlBody)
	return []byte(b.String())
}

// base64Subject 按 RFC 2047 B-encoding 编码主题（支持中文等非 ASCII）。
func base64Subject(subject string) string {
	return base64.StdEncoding.EncodeToString([]byte(subject))
}
