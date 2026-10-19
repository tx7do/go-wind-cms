package channel

// 扩展渠道声明：个人推送（8）+ 移动推送服务（2）。

import (
	"github.com/tx7do/go-utils/notify"

	"github.com/tx7do/go-utils/notify/bark"
	"github.com/tx7do/go-utils/notify/getui"
	"github.com/tx7do/go-utils/notify/gotify"
	"github.com/tx7do/go-utils/notify/jpush"
	"github.com/tx7do/go-utils/notify/ntfy"
	"github.com/tx7do/go-utils/notify/pushbullet"
	"github.com/tx7do/go-utils/notify/pushdeer"
	"github.com/tx7do/go-utils/notify/pushover"
	"github.com/tx7do/go-utils/notify/serverchan"
	"github.com/tx7do/go-utils/notify/wxpusher"
)

func init() {
	registerSettingsChannel(SettingsSchema{
		Type: "BARK", Group: "PUSH", Label: "Bark（iOS）",
		TargetLabel: "设备 Key", NeedsTarget: true,
		Fields: []SettingsField{
			{Name: "serverUrl", Label: "服务地址", Type: "text", Hint: "留空用官方 api.day.app；自建服务可覆盖"},
		},
	}, func(m settingsMap) notify.Service {
		return bark.New(settingsHTTPClient, bark.Config{Ref: m.str("ref"), ServerURL: m.str("serverUrl")})
	})

	registerSettingsChannel(SettingsSchema{
		Type: "NTFY", Group: "PUSH", Label: "ntfy",
		TargetLabel: "Topic", NeedsTarget: true,
		Fields: []SettingsField{
			{Name: "serverUrl", Label: "服务地址", Type: "text", Hint: "留空用官方 ntfy.sh"},
			{Name: "priority", Label: "默认优先级（1~5，0 默认）", Type: "number"},
			{Name: "user", Label: "用户名", Type: "text"},
			{Name: "password", Label: "密码", Type: "password", Sensitive: true},
			{Name: "accessToken", Label: "Access Token", Type: "password", Sensitive: true, Hint: "tk_ 开头；优先于用户名/密码"},
			{Name: "tags", Label: "默认标签", Type: "text", Hint: "逗号分隔（emoji shortcode 或文本）"},
		},
	}, func(m settingsMap) notify.Service {
		return ntfy.New(settingsHTTPClient, ntfy.Config{
			Ref: m.str("ref"), ServerURL: m.str("serverUrl"), Priority: m.int("priority"),
			User: m.str("user"), Password: m.str("password"),
			AccessToken: m.str("accessToken"), Tags: m.strSlice("tags"),
		})
	})

	registerSettingsChannel(SettingsSchema{
		Type: "GOTIFY", Group: "PUSH", Label: "Gotify（自托管）",
		Fields: []SettingsField{
			{Name: "serverUrl", Label: "服务地址", Type: "text", Required: true},
			{Name: "appToken", Label: "App Token", Type: "password", Required: true, Sensitive: true},
			{Name: "priority", Label: "默认优先级", Type: "number"},
		},
	}, func(m settingsMap) notify.Service {
		return gotify.New(settingsHTTPClient, gotify.Config{Ref: m.str("ref"), ServerURL: m.str("serverUrl"), AppToken: m.str("appToken"), Priority: m.int("priority")})
	})

	registerSettingsChannel(SettingsSchema{
		Type: "SERVERCHAN", Group: "PUSH", Label: "Server酱",
		Fields: []SettingsField{
			{Name: "sendKey", Label: "SendKey", Type: "password", Required: true, Sensitive: true},
			{Name: "apiBase", Label: "API 地址", Type: "text", Hint: "留空用官方地址"},
		},
	}, func(m settingsMap) notify.Service {
		return serverchan.New(settingsHTTPClient, serverchan.Config{Ref: m.str("ref"), SendKey: m.str("sendKey"), APIBase: m.str("apiBase")})
	})

	registerSettingsChannel(SettingsSchema{
		Type: "WXPUSHER", Group: "PUSH", Label: "WxPusher（微信推送）",
		TargetLabel: "用户 UID", NeedsTarget: true,
		Fields: []SettingsField{
			{Name: "appToken", Label: "App Token", Type: "password", Required: true, Sensitive: true},
			{Name: "apiBase", Label: "API 地址", Type: "text", Hint: "留空用官方地址"},
		},
	}, func(m settingsMap) notify.Service {
		return wxpusher.New(settingsHTTPClient, wxpusher.Config{Ref: m.str("ref"), AppToken: m.str("appToken"), APIBase: m.str("apiBase")})
	})

	registerSettingsChannel(SettingsSchema{
		Type: "PUSHOVER", Group: "PUSH", Label: "Pushover",
		TargetLabel: "User Key", NeedsTarget: true,
		Fields: []SettingsField{
			{Name: "appToken", Label: "App Token", Type: "password", Required: true, Sensitive: true},
			{Name: "apiBase", Label: "API 地址", Type: "text", Hint: "留空用官方地址"},
		},
	}, func(m settingsMap) notify.Service {
		return pushover.New(settingsHTTPClient, pushover.Config{Ref: m.str("ref"), AppToken: m.str("appToken"), APIBase: m.str("apiBase")})
	})

	registerSettingsChannel(SettingsSchema{
		Type: "PUSHBULLET", Group: "PUSH", Label: "Pushbullet",
		TargetLabel: "设备 ID（可选，留空发全部设备）", NeedsTarget: true,
		Fields: []SettingsField{
			{Name: "apiKey", Label: "API Key", Type: "password", Required: true, Sensitive: true},
			{Name: "apiBase", Label: "API 地址", Type: "text", Hint: "留空用官方地址"},
		},
	}, func(m settingsMap) notify.Service {
		return pushbullet.New(settingsHTTPClient, pushbullet.Config{Ref: m.str("ref"), APIKey: m.str("apiKey"), APIBase: m.str("apiBase")})
	})

	registerSettingsChannel(SettingsSchema{
		Type: "PUSHDEER", Group: "PUSH", Label: "PushDeer",
		Fields: []SettingsField{
			{Name: "pushKey", Label: "Push Key", Type: "password", Required: true, Sensitive: true},
			{Name: "serverUrl", Label: "服务地址", Type: "text", Hint: "留空用官方地址；自建可覆盖"},
		},
	}, func(m settingsMap) notify.Service {
		return pushdeer.New(settingsHTTPClient, pushdeer.Config{Ref: m.str("ref"), PushKey: m.str("pushKey"), ServerURL: m.str("serverUrl")})
	})

	registerSettingsChannel(SettingsSchema{
		Type: "JPUSH", Group: "MOBILE", Label: "极光推送（JPush）",
		TargetLabel: "Registration ID", NeedsTarget: true,
		Fields: []SettingsField{
			{Name: "appKey", Label: "AppKey", Type: "text", Required: true},
			{Name: "masterSecret", Label: "Master Secret", Type: "password", Required: true, Sensitive: true},
			{Name: "apiBase", Label: "API 地址", Type: "text", Hint: "留空用官方地址"},
		},
	}, func(m settingsMap) notify.Service {
		return jpush.New(settingsHTTPClient, jpush.Config{Ref: m.str("ref"), AppKey: m.str("appKey"), MasterSecret: m.str("masterSecret"), APIBase: m.str("apiBase")})
	})

	registerSettingsChannel(SettingsSchema{
		Type: "GETUI", Group: "MOBILE", Label: "个推（GeTui）",
		TargetLabel: "CID（设备标识）", NeedsTarget: true,
		Fields: []SettingsField{
			{Name: "appId", Label: "App ID", Type: "text", Required: true},
			{Name: "appKey", Label: "AppKey", Type: "text", Required: true},
			{Name: "masterSecret", Label: "Master Secret", Type: "password", Required: true, Sensitive: true},
			{Name: "apiBase", Label: "API 地址", Type: "text", Hint: "留空用官方地址"},
		},
	}, func(m settingsMap) notify.Service {
		return getui.New(settingsHTTPClient, getui.Config{Ref: m.str("ref"), AppID: m.str("appId"), AppKey: m.str("appKey"), MasterSecret: m.str("masterSecret"), APIBase: m.str("apiBase")})
	})
}
