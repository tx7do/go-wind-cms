package channel

import (
	"context"
	"errors"
	"fmt"
	"html"
	"os"
	"strings"

	"go-wind-cms/pkg/mailer"

	notifyV1 "go-wind-cms/api/gen/go/notify/service/v1"

	"go-wind-cms/app/core/service/internal/data"
)

// EmailSender 邮件（SMTP）渠道。
//
// 出站走 CMS pkg/mailer（HTML 正文）：通知模板/调用方给的是纯文本标题+正文，
// 本 sender 统一包一层品牌 HTML 壳（标题/正文转义、站点链接、页脚），与
// 旧邮件模板链（notify/system + 内置 HTML）的观感对齐——通知域接手邮件出站后
// 系统通知邮件不因迁移降级成纯文本。
//
// 账号解析策略收在这一处：显式指定 channel_id 则用该条（配置页的测试邮件），
// 否则取第一个启用的 EMAIL 渠道；花名册无可用行时回退 SMTP_* 环境变量
//（与旧 resolveMailerConfig 同一条兜底口径，空环境起步行为不变），
// 两者皆无才 ErrChannelNotConfigured。
//
// 不持有 logger：发送失败由服务层统一记账并写日志（脱敏后的目标 + 台账 ID），
// 在这里再记一遍只会漏掉台账关联。
type EmailSender struct {
	channel *data.NotificationChannelRepo
}

var _ Prechecker = (*EmailSender)(nil)

func NewEmailSender(channelRepo *data.NotificationChannelRepo) *EmailSender {
	return &EmailSender{channel: channelRepo}
}

func (s *EmailSender) Channel() notifyV1.Channel {
	return notifyV1.Channel_EMAIL
}

// Send 投递一封邮件。返回的 error 一律是"投递失败"，包括渠道不可用——
// 由服务层用 ErrChannelNotConfigured 判别是否记 SKIPPED。
func (s *EmailSender) Send(ctx context.Context, req *SendRequest) (*SendReceipt, error) {
	if req == nil || req.Target == "" {
		return nil, errors.New("email target is empty")
	}

	cfg, accountID, err := s.pickAccount(ctx, req.ChannelID)
	if err != nil {
		return nil, err
	}

	if err = cfg.Send([]string{req.Target}, req.Title, buildEmailHTML(req.Title, req.Content)); err != nil {
		// 账号已选出：失败也必须把"是哪条渠道没发出去"带回给台账（见 Sender 的契约）；
		// env 兜底账号没有渠道行，accountID=0 → 台账 channel_id 留空。
		if accountID > 0 {
			return &SendReceipt{ChannelID: accountID}, fmt.Errorf("send mail via channel [%d] failed: %w", accountID, err)
		}
		return nil, fmt.Errorf("send mail via smtp env fallback failed: %w", err)
	}

	return &SendReceipt{ChannelID: accountID}, nil
}

// Precheck 只解析 SMTP 配置并自检可用性，不拨号。见 channel.Prechecker。
func (s *EmailSender) Precheck(ctx context.Context, channelID uint32) error {
	_, _, err := s.pickAccount(ctx, channelID)
	return err
}

// pickAccount 解析实际使用的 SMTP 配置：显式 ID 优先，并校验类型与启用状态；
// 花名册无可用行时回退 SMTP_* 环境变量（accountID=0 标记兜底账号）。
//
// 一切"渠道用不了"（没配、没启用、类型不对、查不到、env 也没配）都包成
// ErrChannelNotConfigured 返回——服务层据此把台账记 SKIPPED（没尝试过投递）
// 而非 FAILED（投递报错），原始原因留在 error 文本里。
func (s *EmailSender) pickAccount(ctx context.Context, channelID uint32) (mailer.Config, uint32, error) {
	cfg, account, err := s.resolveAccount(ctx, channelID)
	if err != nil {
		// 花名册自选落空 → env 兜底；显式点名渠道失败不兜底（管理员点名的就是那一行）
		if channelID == 0 && errors.Is(err, ErrChannelNotConfigured) {
			envCfg := mailer.FromEnv()
			if envCfg.Enabled() {
				return envCfg, 0, nil
			}
		}
		return mailer.Config{}, 0, err
	}
	if err = checkUsable(cfg); err != nil {
		return mailer.Config{}, 0, err
	}
	return cfg, account, nil
}

// checkUsable 已解析配置的可用性自检。
//
// 加密方式是字符串列，库里存着枚举外的值时原本要一路走到拨号才报错，而"没拨过号"
// 和"拨号被拒"在台账上是两种结论（SKIPPED / FAILED）：在这里拦掉。
func checkUsable(cfg mailer.Config) error {
	if !cfg.Enabled() {
		return fmt.Errorf("%w: smtp host/from is empty", ErrChannelNotConfigured)
	}
	return nil
}

func (s *EmailSender) resolveAccount(ctx context.Context, channelID uint32) (mailer.Config, uint32, error) {
	if channelID == 0 {
		account, err := s.channel.GetFirstEnabledEmailChannel(ctx)
		if err != nil {
			return mailer.Config{}, 0, fmt.Errorf("%w: %v", ErrChannelNotConfigured, err)
		}
		return smtpConfigOf(account), account.ID, nil
	}

	// 先走脱敏的 Get 校验类型/启用：GetDecryptedSmtpAccount 不区分渠道类型，
	// 直接取会把一条 WEBHOOK 配置当 SMTP 用。
	dto, err := s.channel.Get(ctx, channelID)
	if err != nil {
		return mailer.Config{}, 0, fmt.Errorf("%w: channel [%d] not found: %v", ErrChannelNotConfigured, channelID, err)
	}
	if dto.GetType() != notifyV1.NotificationChannel_EMAIL {
		return mailer.Config{}, 0, fmt.Errorf("%w: channel [%d] is not an EMAIL channel", ErrChannelNotConfigured, channelID)
	}
	if !dto.GetEnabled() {
		return mailer.Config{}, 0, fmt.Errorf("%w: channel [%d] is disabled", ErrChannelNotConfigured, channelID)
	}

	account, err := s.channel.GetDecryptedSmtpAccount(ctx, channelID)
	if err != nil {
		return mailer.Config{}, 0, fmt.Errorf("%w: channel [%d] credentials: %v", ErrChannelNotConfigured, channelID, err)
	}

	return smtpConfigOf(account), account.ID, nil
}

// smtpConfigOf 花名册账号 → pkg/mailer 配置。TLS 枚举列存大写（NONE/START_TLS/SSL），
// mailer 认小写（starttls/ssl/none），映射而非透传；认不出的值回落 starttls。
func smtpConfigOf(account *data.SmtpAccount) mailer.Config {
	tlsMode := strings.ToLower(account.TlsMode)
	switch tlsMode {
	case "ssl", "none":
	default:
		tlsMode = "starttls"
	}
	port := account.Port
	if port == 0 {
		port = 587
	}
	return mailer.Config{
		Host:     account.Host,
		Port:     int(port),
		Username: account.Username,
		Password: account.Password,
		From:     account.From,
		TLSMode:  tlsMode,
	}
}

// buildEmailHTML 通知域邮件的统一 HTML 壳：标题/正文转义、正文换行转 <br/>、
// 站点链接（SITE_URL，未配置不附）、自动发送页脚。与旧邮件模板链的内置模板同形——
// 系统通知邮件从旧链迁到通知域后观感不降级。壳在本 sender 收口，通知模板
// 只负责标题与正文文字，不做 HTML 拼接。
func buildEmailHTML(title, content string) string {
	siteURL := strings.TrimRight(os.Getenv("SITE_URL"), "/")
	body := strings.ReplaceAll(html.EscapeString(content), "\n", "<br/>")
	var b strings.Builder
	b.WriteString(`<div style="max-width:560px;margin:0 auto;font-family:-apple-system,'Segoe UI',Roboto,'Helvetica Neue',Arial,sans-serif;color:#1f2329;">`)
	b.WriteString(`<h2 style="font-size:18px;margin:24px 0 12px;">` + html.EscapeString(title) + `</h2>`)
	b.WriteString(`<div style="font-size:14px;line-height:1.7;">` + body + `</div>`)
	if siteURL != "" {
		b.WriteString(`<p style="margin:24px 0 8px;"><a href="` + siteURL + `" style="color:#006be6;">` + html.EscapeString(siteURL) + `</a></p>`)
	}
	b.WriteString(`<hr style="border:none;border-top:1px solid #e5e6eb;margin:16px 0;"/><p style="font-size:12px;color:#86909c;">` +
		html.EscapeString("本邮件由系统自动发送，请勿直接回复。") + `</p></div>`)
	return b.String()
}
