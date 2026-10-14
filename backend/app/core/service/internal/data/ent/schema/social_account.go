package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/tx7do/go-crud/entgo/mixin"
)

// SocialAccount 社交账号：一行 = 一个租户在某个社交平台上的发布身份。
// 凭据（App Password / API Key 等）以 JSON 整体存于 credentials，
// AES-GCM 加密落库（enc: 前缀），读取路径由 service 层清空、只回脱敏 hint。
type SocialAccount struct{ ent.Schema }

func (SocialAccount) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "social_accounts", Charset: "utf8mb4", Collation: "utf8mb4_bin"},
		entsql.WithComments(true),
		schema.Comment("社交账号"),
	}
}

func (SocialAccount) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("platform").
			Comment("社交平台").
			NamedValues(
				"PlatformBluesky", "PLATFORM_BLUESKY",
				"PlatformMastodon", "PLATFORM_MASTODON",
				"PlatformDevto", "PLATFORM_DEVTO",
				"PlatformTelegram", "PLATFORM_TELEGRAM",
				"PlatformWordpress", "PLATFORM_WORDPRESS",
				"PlatformWechatOa", "PLATFORM_WECHAT_OA",
				"PlatformNostr", "PLATFORM_NOSTR",
				"PlatformWeibo", "PLATFORM_WEIBO",
				"PlatformLinkedin", "PLATFORM_LINKEDIN",
				"PlatformMetaweblog", "PLATFORM_METAWEBLOG",
			).
			Optional().
			Nillable(),
		field.String("name").Comment("账号显示名称").Optional().Nillable(),
		field.String("handle").Comment("平台内账号标识（如 user.bsky.social）").Optional().Nillable(),
		field.Text("credentials").Comment("平台凭据 JSON（AES-GCM 加密落库，enc: 前缀；仅写入，读取由服务层清空）").Optional().Nillable().Sensitive(),
		field.String("credential_hint").Comment("凭据脱敏提示（如 app-password ***wxyz）").Optional().Nillable(),
		field.String("remark").Comment("备注").Optional().Nillable(),
	}
}

func (SocialAccount) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.AutoIncrementId{},
		mixin.TimeAt{},
		mixin.OperatorID{},
		mixin.IsEnabled{},
		mixin.TenantID[uint32]{},
	}
}

func (SocialAccount) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("tenant_id", "platform").StorageKey("idx_social_accounts_tenant_platform"),
		index.Fields("platform").StorageKey("idx_social_accounts_platform"),
		index.Fields("created_at").StorageKey("idx_social_accounts_created_at"),
	}
}
