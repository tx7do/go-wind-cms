package channel

// 扩展渠道声明：IM 组（15 个）。字段名（settings JSON 键）与各渠道包 Config 一一对应，
// 敏感字段（token/secret）由 repo 按 Sensitive 标记加密落库、回显脱敏。

import (
	"github.com/tx7do/go-utils/notify"

	"github.com/tx7do/go-utils/notify/dingtalk"
	"github.com/tx7do/go-utils/notify/discord"
	"github.com/tx7do/go-utils/notify/feishu"
	"github.com/tx7do/go-utils/notify/gchat"
	"github.com/tx7do/go-utils/notify/line"
	"github.com/tx7do/go-utils/notify/matrix"
	"github.com/tx7do/go-utils/notify/mattermost"
	"github.com/tx7do/go-utils/notify/rocketchat"
	"github.com/tx7do/go-utils/notify/slack"
	"github.com/tx7do/go-utils/notify/teams"
	"github.com/tx7do/go-utils/notify/telegram"
	"github.com/tx7do/go-utils/notify/webex"
	"github.com/tx7do/go-utils/notify/wechat"
	"github.com/tx7do/go-utils/notify/wecom"
	"github.com/tx7do/go-utils/notify/zulip"
)

func init() {
	registerSettingsChannel(SettingsSchema{
		Type: "WECOM_APP", Group: "IM", Label: "企业微信应用消息",
		TargetLabel: "成员 UserID（多个用 | 分隔）", NeedsTarget: true,
		Fields: []SettingsField{
			{Name: "corpId", Label: "企业 ID", Type: "text", Required: true},
			{Name: "corpSecret", Label: "应用 Secret", Type: "password", Required: true, Sensitive: true},
			{Name: "agentId", Label: "应用 AgentID", Type: "number", Required: true},
			{Name: "apiBase", Label: "API 地址", Type: "text", Hint: "留空用官方地址"},
		},
	}, func(m settingsMap) notify.Service {
		return wecom.New(settingsHTTPClient, wecom.Config{
			Ref: m.str("ref"), CorpID: m.str("corpId"), CorpSecret: m.str("corpSecret"),
			AgentID: m.int64("agentId"), APIBase: m.str("apiBase"),
		})
	})

	registerSettingsChannel(SettingsSchema{
		Type: "DINGTALK_APP", Group: "IM", Label: "钉钉工作通知",
		TargetLabel: "成员 UserID（多个用 , 分隔）", NeedsTarget: true,
		Fields: []SettingsField{
			{Name: "appKey", Label: "AppKey", Type: "text", Required: true},
			{Name: "appSecret", Label: "AppSecret", Type: "password", Required: true, Sensitive: true},
			{Name: "agentId", Label: "应用 AgentID", Type: "number", Required: true},
			{Name: "apiBase", Label: "API 地址", Type: "text", Hint: "留空用官方地址"},
		},
	}, func(m settingsMap) notify.Service {
		return dingtalk.New(settingsHTTPClient, dingtalk.Config{
			Ref: m.str("ref"), AppKey: m.str("appKey"), AppSecret: m.str("appSecret"),
			AgentID: m.int64("agentId"), APIBase: m.str("apiBase"),
		})
	})

	registerSettingsChannel(SettingsSchema{
		Type: "FEISHU_APP", Group: "IM", Label: "飞书应用消息",
		TargetLabel: "接收者 ID（按 ID 类型）", NeedsTarget: true,
		Fields: []SettingsField{
			{Name: "appId", Label: "App ID", Type: "text", Required: true},
			{Name: "appSecret", Label: "App Secret", Type: "password", Required: true, Sensitive: true},
			{Name: "receiveIdType", Label: "ID 类型", Type: "select", Options: []string{"open_id", "user_id", "union_id", "email", "chat_id"}},
			{Name: "apiBase", Label: "API 地址", Type: "text", Hint: "留空用官方地址"},
		},
	}, func(m settingsMap) notify.Service {
		return feishu.New(settingsHTTPClient, feishu.Config{
			Ref: m.str("ref"), AppID: m.str("appId"), AppSecret: m.str("appSecret"),
			ReceiveIDType: m.str("receiveIdType"), APIBase: m.str("apiBase"),
		})
	})

	registerSettingsChannel(SettingsSchema{
		Type: "WECHAT", Group: "IM", Label: "微信公众号模板消息",
		TargetLabel: "接收者 OpenID", NeedsTarget: true,
		Fields: []SettingsField{
			{Name: "appId", Label: "AppID", Type: "text", Required: true},
			{Name: "appSecret", Label: "AppSecret", Type: "password", Required: true, Sensitive: true},
			{Name: "templateId", Label: "模板 ID", Type: "text", Required: true},
			{Name: "apiBase", Label: "API 地址", Type: "text", Hint: "留空用官方地址"},
		},
	}, func(m settingsMap) notify.Service {
		return wechat.New(settingsHTTPClient, wechat.Config{
			Ref: m.str("ref"), AppID: m.str("appId"), AppSecret: m.str("appSecret"),
			TemplateID: m.str("templateId"), APIBase: m.str("apiBase"),
		})
	})

	registerSettingsChannel(SettingsSchema{
		Type: "TELEGRAM", Group: "IM", Label: "Telegram Bot",
		TargetLabel: "Chat ID", NeedsTarget: true,
		Fields: []SettingsField{
			{Name: "botToken", Label: "Bot Token", Type: "password", Required: true, Sensitive: true},
			{Name: "apiBase", Label: "API 地址", Type: "text", Hint: "留空用官方 API；自建网关可覆盖"},
			{Name: "parseMode", Label: "解析模式", Type: "select", Options: []string{"", "MarkdownV2", "HTML"}},
		},
	}, func(m settingsMap) notify.Service {
		return telegram.New(settingsHTTPClient, telegram.Config{
			Ref: m.str("ref"), BotToken: m.str("botToken"),
			APIBase: m.str("apiBase"), ParseMode: m.str("parseMode"),
		})
	})

	registerSettingsChannel(SettingsSchema{
		Type: "SLACK", Group: "IM", Label: "Slack（Incoming Webhook）",
		Fields: []SettingsField{
			{Name: "webhookUrl", Label: "Webhook 地址", Type: "password", Required: true, Sensitive: true, Placeholder: "https://hooks.slack.com/services/..."},
		},
	}, func(m settingsMap) notify.Service {
		return slack.New(settingsHTTPClient, slack.Config{Ref: m.str("ref"), WebhookURL: m.str("webhookUrl")})
	})

	registerSettingsChannel(SettingsSchema{
		Type: "DISCORD", Group: "IM", Label: "Discord（Webhook）",
		Fields: []SettingsField{
			{Name: "webhookUrl", Label: "Webhook 地址", Type: "password", Required: true, Sensitive: true},
		},
	}, func(m settingsMap) notify.Service {
		return discord.New(settingsHTTPClient, discord.Config{Ref: m.str("ref"), WebhookURL: m.str("webhookUrl")})
	})

	registerSettingsChannel(SettingsSchema{
		Type: "TEAMS", Group: "IM", Label: "Microsoft Teams（Workflow/Webhook）",
		Fields: []SettingsField{
			{Name: "webhookUrl", Label: "Webhook 地址", Type: "password", Required: true, Sensitive: true},
		},
	}, func(m settingsMap) notify.Service {
		return teams.New(settingsHTTPClient, teams.Config{Ref: m.str("ref"), WebhookURL: m.str("webhookUrl")})
	})

	registerSettingsChannel(SettingsSchema{
		Type: "GCHAT", Group: "IM", Label: "Google Chat（Webhook）",
		Fields: []SettingsField{
			{Name: "webhookUrl", Label: "Webhook 地址", Type: "password", Required: true, Sensitive: true},
		},
	}, func(m settingsMap) notify.Service {
		return gchat.New(settingsHTTPClient, gchat.Config{Ref: m.str("ref"), WebhookURL: m.str("webhookUrl")})
	})

	registerSettingsChannel(SettingsSchema{
		Type: "WEBEX", Group: "IM", Label: "Webex Bot",
		TargetLabel: "收件人邮箱或 Room ID", NeedsTarget: true,
		Fields: []SettingsField{
			{Name: "botToken", Label: "Bot Token", Type: "password", Required: true, Sensitive: true},
			{Name: "apiBase", Label: "API 地址", Type: "text", Hint: "留空用官方地址"},
		},
	}, func(m settingsMap) notify.Service {
		return webex.New(settingsHTTPClient, webex.Config{Ref: m.str("ref"), BotToken: m.str("botToken"), APIBase: m.str("apiBase")})
	})

	registerSettingsChannel(SettingsSchema{
		Type: "LINE", Group: "IM", Label: "LINE Messaging API",
		TargetLabel: "用户 ID", NeedsTarget: true,
		Fields: []SettingsField{
			{Name: "channelAccessToken", Label: "Channel Access Token", Type: "password", Required: true, Sensitive: true},
			{Name: "apiBase", Label: "API 地址", Type: "text", Hint: "留空用官方地址"},
		},
	}, func(m settingsMap) notify.Service {
		return line.New(settingsHTTPClient, line.Config{Ref: m.str("ref"), ChannelAccessToken: m.str("channelAccessToken"), APIBase: m.str("apiBase")})
	})

	registerSettingsChannel(SettingsSchema{
		Type: "MATRIX", Group: "IM", Label: "Matrix",
		TargetLabel: "房间 ID（!room:server）", NeedsTarget: true,
		Fields: []SettingsField{
			{Name: "serverUrl", Label: "服务地址", Type: "text", Required: true},
			{Name: "accessToken", Label: "Access Token", Type: "password", Required: true, Sensitive: true},
		},
	}, func(m settingsMap) notify.Service {
		return matrix.New(settingsHTTPClient, matrix.Config{Ref: m.str("ref"), ServerURL: m.str("serverUrl"), AccessToken: m.str("accessToken")})
	})

	registerSettingsChannel(SettingsSchema{
		Type: "MATTERMOST", Group: "IM", Label: "Mattermost（Incoming Webhook）",
		Fields: []SettingsField{
			{Name: "webhookUrl", Label: "Webhook 地址", Type: "password", Required: true, Sensitive: true},
			{Name: "channel", Label: "频道（可选）", Type: "text"},
		},
	}, func(m settingsMap) notify.Service {
		return mattermost.New(settingsHTTPClient, mattermost.Config{Ref: m.str("ref"), WebhookURL: m.str("webhookUrl"), Channel: m.str("channel")})
	})

	registerSettingsChannel(SettingsSchema{
		Type: "ROCKETCHAT", Group: "IM", Label: "Rocket.Chat",
		TargetLabel: "频道名（#开头）或 @用户", NeedsTarget: true,
		Fields: []SettingsField{
			{Name: "serverUrl", Label: "服务地址", Type: "text", Required: true},
			{Name: "userId", Label: "User ID", Type: "text", Required: true},
			{Name: "authToken", Label: "Auth Token", Type: "password", Required: true, Sensitive: true},
		},
	}, func(m settingsMap) notify.Service {
		return rocketchat.New(settingsHTTPClient, rocketchat.Config{Ref: m.str("ref"), ServerURL: m.str("serverUrl"), UserID: m.str("userId"), AuthToken: m.str("authToken")})
	})

	registerSettingsChannel(SettingsSchema{
		Type: "ZULIP", Group: "IM", Label: "Zulip",
		TargetLabel: "接收者（用户邮箱或频道名）", NeedsTarget: true,
		Fields: []SettingsField{
			{Name: "serverUrl", Label: "服务地址", Type: "text", Required: true},
			{Name: "botEmail", Label: "Bot 邮箱", Type: "text", Required: true},
			{Name: "botApiKey", Label: "Bot API Key", Type: "password", Required: true, Sensitive: true},
			{Name: "topic", Label: "话题（可选）", Type: "text"},
		},
	}, func(m settingsMap) notify.Service {
		return zulip.New(settingsHTTPClient, zulip.Config{Ref: m.str("ref"), ServerURL: m.str("serverUrl"), BotEmail: m.str("botEmail"), BotAPIKey: m.str("botApiKey"), Topic: m.str("topic")})
	})
}
