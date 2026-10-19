package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/tx7do/go-crud/entgo/mixin"
)

// NotificationChannel holds the schema definition for the NotificationChannel entity.
// 通知渠道为平台级全局配置（无租户隔离）：由平台管理员维护，供站内信外发/
// webhook 通知等能力使用；读侧只开放平台视图（service 层 IsPlatformContext 守卫）。
type NotificationChannel struct {
	ent.Schema
}

func (NotificationChannel) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "notification_channels",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
		entsql.WithComments(true),
		schema.Comment("通知渠道表（EMAIL/WEBHOOK 花名册）"),
	}
}

// Fields of the NotificationChannel.
func (NotificationChannel) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Comment("渠道名称").
			NotEmpty(),
		field.Enum("type").
			Comment("渠道类型（成员名与 proto Type 枚举逐字相同；Nillable 与全域枚举列同规：copier 枚举转换是指针↔指针对）。EXT_* 扩展渠道的配置整包存 settings 列").
			NamedValues(
				"Email", "EMAIL",
				"Webhook", "WEBHOOK",
				// ---- 扩展渠道（EXT：settings 列承载配置，proto Type 段位 100+）----
				"WecomApp", "WECOM_APP", "DingtalkApp", "DINGTALK_APP", "FeishuApp", "FEISHU_APP", "Wechat", "WECHAT",
				"Telegram", "TELEGRAM", "Slack", "SLACK", "Discord", "DISCORD", "Teams", "TEAMS",
				"Gchat", "GCHAT", "Webex", "WEBEX", "Line", "LINE", "Matrix", "MATRIX",
				"Mattermost", "MATTERMOST", "Rocketchat", "ROCKETCHAT", "Zulip", "ZULIP",
				"Bark", "BARK", "Ntfy", "NTFY", "Gotify", "GOTIFY", "Serverchan", "SERVERCHAN",
				"Wxpusher", "WXPUSHER", "Pushover", "PUSHOVER", "Pushbullet", "PUSHBULLET", "Pushdeer", "PUSHDEER",
				"Jpush", "JPUSH", "Getui", "GETUI",
				"Huaweisms", "HUAWEISMS", "Yunpian", "YUNPIAN", "Twilio", "TWILIO", "Vonage", "VONAGE",
				"Whatsapp", "WHATSAPP",
				"Mailgun", "MAILGUN", "Sendgrid", "SENDGRID", "Postmark", "POSTMARK",
				"Ifttt", "IFTTT", "Join", "JOIN", "Homeassistant", "HOMEASSISTANT",
				"Datadog", "DATADOG", "Opsgenie", "OPSGENIE", "Pagerduty", "PAGERDUTY",
			).
			Default("EMAIL").
			Optional().
			Nillable(),
		field.String("smtp_host").
			Comment("SMTP 服务器地址").
			Optional().
			Nillable(),
		field.Uint32("smtp_port").
			Comment("SMTP 端口（25/465/587）").
			Optional().
			Nillable(),
		field.String("smtp_username").
			Comment("SMTP 用户名").
			Optional().
			Nillable(),
		field.String("smtp_password").
			Comment("SMTP 密码/授权码（AES-GCM 加密存储）").
			Sensitive().
			Optional().
			Nillable(),
		field.String("smtp_from").
			Comment("发件人地址").
			Optional().
			Nillable(),
		field.String("smtp_from_name").
			Comment("发件人显示名（收件客户端直接显示该名）").
			Optional().
			Nillable(),
		field.Enum("smtp_tls").
			Comment("加密方式").
			NamedValues(
				"None", "NONE",
				"StartTls", "START_TLS",
				"Ssl", "SSL",
			).
			Default("START_TLS").
			Optional().
			Nillable(),

		// WEBHOOK 渠道的四列：SMTP 那几列对它没有意义。按渠道加列而不上 settings JSON
		// ——同一行只允许一种真相源。
		field.String("webhook_url").
			Comment("Webhook 回调地址（仅 WEBHOOK 渠道）").
			Optional().
			Nillable(),
		field.String("webhook_secret").
			Comment("Webhook 签名密钥（AES-GCM 加密存储，仅 WEBHOOK 渠道）").
			Sensitive().
			Optional().
			Nillable(),
		// 出站风格与载荷模板：同属"这一条 webhook 到底发成什么形状"，跟着渠道行走，
		// 不进规则表（规则表答的是"这个事件走哪个渠道"）。
		field.Enum("webhook_sign_style").
			Comment("Webhook 出站风格（签名位置/算法与应答判据，仅 WEBHOOK 渠道）").
			NamedValues(
				"Custom", "CUSTOM",
				"None", "NONE",
				"Dingtalk", "DINGTALK",
				"Feishu", "FEISHU",
				"Wecom", "WECOM",
			).
			Default("CUSTOM").
			Optional().
			Nillable(),
		field.String("webhook_payload_template").
			Comment("Webhook 载荷模板（{{占位符}} 渲染；留空则用该风格的内置默认 JSON）").
			Optional().
			Nillable(),
		field.String("settings").
			Comment("扩展渠道（EXT_* 类型）配置整包：JSON，sensitive 字段加密落库；EMAIL/WEBHOOK 不使用").
			Optional().
			Nillable(),
	}
}

// Mixin of the NotificationChannel.
func (NotificationChannel) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.AutoIncrementId{},
		mixin.TimeAt{},
		mixin.OperatorID{},
		mixin.Remark{},
		mixin.SwitchStatus{},
	}
}

// Indexes of the NotificationChannel.
func (NotificationChannel) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").Unique().
			StorageKey("uidx_notification_channel_name"),
	}
}
