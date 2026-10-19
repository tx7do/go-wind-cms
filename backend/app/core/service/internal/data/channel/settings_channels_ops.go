package channel

// 扩展渠道声明：短信/消息（5）+ 邮件 API（3）+ 效率自动化（3）+ 运维告警（3）。

import (
	"github.com/tx7do/go-utils/notify"

	"github.com/tx7do/go-utils/notify/datadog"
	"github.com/tx7do/go-utils/notify/homeassistant"
	"github.com/tx7do/go-utils/notify/huaweisms"
	"github.com/tx7do/go-utils/notify/ifttt"
	"github.com/tx7do/go-utils/notify/join"
	"github.com/tx7do/go-utils/notify/mailgun"
	"github.com/tx7do/go-utils/notify/opsgenie"
	"github.com/tx7do/go-utils/notify/pagerduty"
	"github.com/tx7do/go-utils/notify/postmark"
	"github.com/tx7do/go-utils/notify/sendgrid"
	"github.com/tx7do/go-utils/notify/twilio"
	"github.com/tx7do/go-utils/notify/vonage"
	"github.com/tx7do/go-utils/notify/whatsapp"
	"github.com/tx7do/go-utils/notify/yunpian"
)

func init() {
	// ---- 短信/消息（5）----

	registerSettingsChannel(SettingsSchema{
		Type: "HUAWEISMS", Group: "SMS", Label: "华为云短信",
		TargetLabel: "手机号（E.164）", NeedsTarget: true,
		Fields: []SettingsField{
			{Name: "appKey", Label: "App Key", Type: "text", Required: true},
			{Name: "appSecret", Label: "App Secret", Type: "password", Required: true, Sensitive: true},
			{Name: "sender", Label: "签名通道号", Type: "text", Required: true},
			{Name: "templateId", Label: "模板 ID", Type: "text", Required: true},
			{Name: "statusCallback", Label: "状态回调地址（可选）", Type: "text"},
			{Name: "apiBase", Label: "API 地址", Type: "text", Hint: "留空用官方地址"},
		},
	}, func(m settingsMap) notify.Service {
		return huaweisms.New(settingsHTTPClient, huaweisms.Config{
			Ref: m.str("ref"), AppKey: m.str("appKey"), AppSecret: m.str("appSecret"),
			Sender: m.str("sender"), TemplateID: m.str("templateId"),
			StatusCallback: m.str("statusCallback"), APIBase: m.str("apiBase"),
		})
	})

	registerSettingsChannel(SettingsSchema{
		Type: "YUNPIAN", Group: "SMS", Label: "云片短信",
		TargetLabel: "手机号（大陆）", NeedsTarget: true,
		Fields: []SettingsField{
			{Name: "apiKey", Label: "API Key", Type: "password", Required: true, Sensitive: true},
			{Name: "apiBase", Label: "API 地址", Type: "text", Hint: "留空用官方地址"},
		},
	}, func(m settingsMap) notify.Service {
		return yunpian.New(settingsHTTPClient, yunpian.Config{Ref: m.str("ref"), ApiKey: m.str("apiKey"), APIBase: m.str("apiBase")})
	})

	registerSettingsChannel(SettingsSchema{
		Type: "TWILIO", Group: "SMS", Label: "Twilio 短信",
		TargetLabel: "手机号（E.164）", NeedsTarget: true,
		Fields: []SettingsField{
			{Name: "accountSid", Label: "Account SID", Type: "text", Required: true},
			{Name: "authToken", Label: "Auth Token", Type: "password", Required: true, Sensitive: true},
			{Name: "from", Label: "发信号码", Type: "text", Required: true},
			{Name: "apiBase", Label: "API 地址", Type: "text", Hint: "留空用官方地址"},
		},
	}, func(m settingsMap) notify.Service {
		return twilio.New(settingsHTTPClient, twilio.Config{Ref: m.str("ref"), AccountSID: m.str("accountSid"), AuthToken: m.str("authToken"), From: m.str("from"), APIBase: m.str("apiBase")})
	})

	registerSettingsChannel(SettingsSchema{
		Type: "VONAGE", Group: "SMS", Label: "Vonage 短信",
		TargetLabel: "手机号（E.164）", NeedsTarget: true,
		Fields: []SettingsField{
			{Name: "apiKey", Label: "API Key", Type: "text", Required: true},
			{Name: "apiSecret", Label: "API Secret", Type: "password", Required: true, Sensitive: true},
			{Name: "from", Label: "发信号码/品牌名", Type: "text", Required: true},
			{Name: "apiBase", Label: "API 地址", Type: "text", Hint: "留空用官方地址"},
		},
	}, func(m settingsMap) notify.Service {
		return vonage.New(settingsHTTPClient, vonage.Config{Ref: m.str("ref"), APIKey: m.str("apiKey"), APISecret: m.str("apiSecret"), From: m.str("from"), APIBase: m.str("apiBase")})
	})

	registerSettingsChannel(SettingsSchema{
		Type: "WHATSAPP", Group: "SMS", Label: "WhatsApp Business",
		TargetLabel: "手机号（E.164）", NeedsTarget: true,
		Fields: []SettingsField{
			{Name: "phoneNumberId", Label: "Phone Number ID", Type: "text", Required: true},
			{Name: "accessToken", Label: "Access Token", Type: "password", Required: true, Sensitive: true},
			{Name: "apiVersion", Label: "API 版本", Type: "text", Hint: "留空用库内默认版本"},
			{Name: "apiBase", Label: "API 地址", Type: "text", Hint: "留空用官方地址"},
		},
	}, func(m settingsMap) notify.Service {
		return whatsapp.New(settingsHTTPClient, whatsapp.Config{
			Ref: m.str("ref"), PhoneNumberID: m.str("phoneNumberId"), AccessToken: m.str("accessToken"),
			APIVersion: m.str("apiVersion"), APIBase: m.str("apiBase"),
		})
	})

	// ---- 邮件 API（3）----

	registerSettingsChannel(SettingsSchema{
		Type: "MAILGUN", Group: "MAIL", Label: "Mailgun",
		TargetLabel: "收件人邮箱", NeedsTarget: true,
		Fields: []SettingsField{
			{Name: "domain", Label: "域名", Type: "text", Required: true},
			{Name: "from", Label: "发件人", Type: "text", Required: true},
			{Name: "apiKey", Label: "API Key", Type: "password", Required: true, Sensitive: true},
			{Name: "apiBase", Label: "API 地址", Type: "text", Hint: "留空用官方地址（EU/US 区域不同）"},
		},
	}, func(m settingsMap) notify.Service {
		return mailgun.New(settingsHTTPClient, mailgun.Config{Ref: m.str("ref"), Domain: m.str("domain"), From: m.str("from"), APIKey: m.str("apiKey"), APIBase: m.str("apiBase")})
	})

	registerSettingsChannel(SettingsSchema{
		Type: "SENDGRID", Group: "MAIL", Label: "SendGrid",
		TargetLabel: "收件人邮箱", NeedsTarget: true,
		Fields: []SettingsField{
			{Name: "apiKey", Label: "API Key", Type: "password", Required: true, Sensitive: true},
			{Name: "from", Label: "发件人", Type: "text", Required: true},
			{Name: "apiBase", Label: "API 地址", Type: "text", Hint: "留空用官方地址"},
		},
	}, func(m settingsMap) notify.Service {
		return sendgrid.New(settingsHTTPClient, sendgrid.Config{Ref: m.str("ref"), APIKey: m.str("apiKey"), From: m.str("from"), APIBase: m.str("apiBase")})
	})

	registerSettingsChannel(SettingsSchema{
		Type: "POSTMARK", Group: "MAIL", Label: "Postmark",
		TargetLabel: "收件人邮箱", NeedsTarget: true,
		Fields: []SettingsField{
			{Name: "serverToken", Label: "Server Token", Type: "password", Required: true, Sensitive: true},
			{Name: "from", Label: "发件人", Type: "text", Required: true},
			{Name: "messageStream", Label: "Message Stream", Type: "text", Hint: "留空用默认流"},
			{Name: "apiBase", Label: "API 地址", Type: "text", Hint: "留空用官方地址"},
		},
	}, func(m settingsMap) notify.Service {
		return postmark.New(settingsHTTPClient, postmark.Config{Ref: m.str("ref"), ServerToken: m.str("serverToken"), From: m.str("from"), MessageStream: m.str("messageStream"), APIBase: m.str("apiBase")})
	})

	// ---- 效率自动化（3）----

	registerSettingsChannel(SettingsSchema{
		Type: "IFTTT", Group: "AUTO", Label: "IFTTT",
		TargetLabel: "Event 名称", NeedsTarget: true,
		Fields: []SettingsField{
			{Name: "apiKey", Label: "API Key", Type: "password", Required: true, Sensitive: true},
			{Name: "apiBase", Label: "API 地址", Type: "text", Hint: "留空用官方地址"},
		},
	}, func(m settingsMap) notify.Service {
		return ifttt.New(settingsHTTPClient, ifttt.Config{Ref: m.str("ref"), APIKey: m.str("apiKey"), APIBase: m.str("apiBase")})
	})

	registerSettingsChannel(SettingsSchema{
		Type: "JOIN", Group: "AUTO", Label: "Join（Join by joaoapps）",
		TargetLabel: "设备 ID（可选，留空发全部设备）", NeedsTarget: true,
		Fields: []SettingsField{
			{Name: "apiKey", Label: "API Key", Type: "password", Required: true, Sensitive: true},
			{Name: "apiBase", Label: "API 地址", Type: "text", Hint: "留空用官方地址"},
		},
	}, func(m settingsMap) notify.Service {
		return join.New(settingsHTTPClient, join.Config{Ref: m.str("ref"), APIKey: m.str("apiKey"), APIBase: m.str("apiBase")})
	})

	registerSettingsChannel(SettingsSchema{
		Type: "HOMEASSISTANT", Group: "AUTO", Label: "Home Assistant 通知",
		TargetLabel: "目标实体（可选）", NeedsTarget: true,
		Fields: []SettingsField{
			{Name: "serverUrl", Label: "服务地址", Type: "text", Required: true},
			{Name: "token", Label: "长期访问令牌", Type: "password", Required: true, Sensitive: true},
			{Name: "service", Label: "通知服务", Type: "text", Hint: "如 notify.mobile_app_xxx；留空用默认"},
			{Name: "apiBase", Label: "API 地址", Type: "text", Hint: "留空跟随服务地址"},
		},
	}, func(m settingsMap) notify.Service {
		return homeassistant.New(settingsHTTPClient, homeassistant.Config{Ref: m.str("ref"), ServerURL: m.str("serverUrl"), Token: m.str("token"), Service: m.str("service"), APIBase: m.str("apiBase")})
	})

	// ---- 运维告警（3）----

	registerSettingsChannel(SettingsSchema{
		Type: "DATADOG", Group: "OPS", Label: "Datadog 事件",
		TargetLabel: "Host（可选）", NeedsTarget: true,
		Fields: []SettingsField{
			{Name: "apiKey", Label: "API Key", Type: "password", Required: true, Sensitive: true},
			{Name: "priority", Label: "优先级", Type: "select", Options: []string{"normal", "low"}},
			{Name: "tags", Label: "标签", Type: "text", Hint: "逗号分隔"},
			{Name: "apiBase", Label: "API 地址", Type: "text", Hint: "留空用官方地址（EU 站不同）"},
		},
	}, func(m settingsMap) notify.Service {
		return datadog.New(settingsHTTPClient, datadog.Config{Ref: m.str("ref"), APIKey: m.str("apiKey"), Priority: m.str("priority"), Tags: m.strSlice("tags"), APIBase: m.str("apiBase")})
	})

	registerSettingsChannel(SettingsSchema{
		Type: "OPSGENIE", Group: "OPS", Label: "Opsgenie 告警",
		TargetLabel: "别名（可选，用于去重）", NeedsTarget: true,
		Fields: []SettingsField{
			{Name: "apiKey", Label: "API Key", Type: "password", Required: true, Sensitive: true},
			{Name: "apiBase", Label: "API 地址", Type: "text", Hint: "留空用官方地址（EU 站不同）"},
		},
	}, func(m settingsMap) notify.Service {
		return opsgenie.New(settingsHTTPClient, opsgenie.Config{Ref: m.str("ref"), APIKey: m.str("apiKey"), APIBase: m.str("apiBase")})
	})

	registerSettingsChannel(SettingsSchema{
		Type: "PAGERDUTY", Group: "OPS", Label: "PagerDuty 事件",
		TargetLabel: "Dedup Key（可选，用于去重）", NeedsTarget: true,
		Fields: []SettingsField{
			{Name: "routingKey", Label: "Routing Key", Type: "password", Required: true, Sensitive: true},
			{Name: "severity", Label: "级别", Type: "select", Options: []string{"critical", "error", "warning", "info"}},
			{Name: "source", Label: "来源标识", Type: "text", Hint: "默认 go-wind-admin"},
			{Name: "apiBase", Label: "API 地址", Type: "text", Hint: "留空用官方地址"},
		},
	}, func(m settingsMap) notify.Service {
		return pagerduty.New(settingsHTTPClient, pagerduty.Config{Ref: m.str("ref"), RoutingKey: m.str("routingKey"), Severity: m.str("severity"), Source: m.str("source"), APIBase: m.str("apiBase")})
	})
}
