package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/tx7do/go-crud/entgo/mixin"
)

// SocialInboxItem 社交收件箱条目：从平台回拉的回复/提及等互动内容。
// (platform, external_id) 为去重键——重复回拉不产生新行。
type SocialInboxItem struct{ ent.Schema }

func (SocialInboxItem) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "social_inbox_items", Charset: "utf8mb4", Collation: "utf8mb4_bin"},
		entsql.WithComments(true),
		schema.Comment("社交收件箱条目"),
	}
}

func (SocialInboxItem) Fields() []ent.Field {
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

		field.Uint32("account_id").
			Comment("关联的社交账号ID").
			Optional().
			Nillable(),

		field.String("type").
			Comment("条目类型（REPLY/MENTION/LIKE/OTHER）").
			Optional().
			Nillable(),

		field.String("external_id").
			Comment("平台侧条目唯一标识").
			Optional().
			Nillable(),

		field.String("author_name").
			Comment("作者显示名").
			Optional().
			Nillable(),

		field.String("author_handle").
			Comment("作者账号标识").
			Optional().
			Nillable(),

		field.Text("content").
			Comment("内容摘要").
			Optional().
			Nillable(),

		field.String("url").
			Comment("平台侧链接").
			Optional().
			Nillable(),

		field.Time("occurred_at").
			Comment("平台侧发生时间").
			Optional().
			Nillable(),

		field.Time("read_at").
			Comment("标记已读时间").
			Optional().
			Nillable(),
	}
}

func (SocialInboxItem) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.AutoIncrementId{},
		mixin.TimeAt{},
		mixin.OperatorID{},
		mixin.TenantID[uint32]{},
	}
}

func (SocialInboxItem) Indexes() []ent.Index {
	return []ent.Index{
		// 去重键：平台 + 平台侧条目 ID
		index.Fields("platform", "external_id").StorageKey("uidx_social_inbox_platform_external"),
		index.Fields("account_id", "read_at").StorageKey("idx_social_inbox_account_read"),
		index.Fields("tenant_id", "created_at").StorageKey("idx_social_inbox_tenant_created"),
	}
}
