package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/tx7do/go-crud/entgo/mixin"
)

// SocialDistribution 内容分发记录：CMS post 推送到社交平台的一次发布任务及其结果。
// 内容事实源在 CMS（post_id 指向 posts），本表只存平台覆盖文案、发布状态与回拉的互动统计快照。
type SocialDistribution struct{ ent.Schema }

func (SocialDistribution) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "social_distributions", Charset: "utf8mb4", Collation: "utf8mb4_bin"},
		entsql.WithComments(true),
		schema.Comment("社交内容分发记录"),
	}
}

func (SocialDistribution) Fields() []ent.Field {
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
			Comment("使用的社交账号ID").
			Optional().
			Nillable(),

		field.Uint32("post_id").
			Comment("关联的 CMS 文章ID").
			Optional().
			Nillable(),

		field.Text("override_text").
			Comment("按平台覆盖的发布文案").
			Optional().
			Nillable(),

		field.Enum("status").
			Comment("分发状态").
			NamedValues(
				"StatusDraft", "STATUS_DRAFT",
				"StatusScheduled", "STATUS_SCHEDULED",
				"StatusPublishing", "STATUS_PUBLISHING",
				"StatusPublished", "STATUS_PUBLISHED",
				"StatusFailed", "STATUS_FAILED",
			).
			Optional().
			Nillable(),

		field.Time("publish_time").
			Comment("计划发布时间").
			Optional().
			Nillable(),

		field.Time("published_at").
			Comment("实际发布时间").
			Optional().
			Nillable(),

		field.String("external_uri").
			Comment("平台侧帖子 URI").
			Optional().
			Nillable(),

		field.String("external_id").
			Comment("平台侧帖子 ID").
			Optional().
			Nillable(),

		field.Text("error_message").
			Comment("最近一次发布失败的原因").
			Optional().
			Nillable(),

		field.String("title").
			Comment("平台文章标题（DEV.to 等需要标题的平台）").
			Optional().
			Nillable(),

		field.JSON("media_ids", []uint32{}).
			Comment("配图媒体资产ID列表（支持图片的平台）").
			Optional(),

		field.Uint32("likes").
			Comment("点赞数（平台回拉快照）").
			Default(0).
			Optional().
			Nillable(),

		field.Uint32("replies").
			Comment("回复数（平台回拉快照）").
			Default(0).
			Optional().
			Nillable(),

		field.Uint32("reposts").
			Comment("转发数（平台回拉快照）").
			Default(0).
			Optional().
			Nillable(),

		field.Time("stats_synced_at").
			Comment("统计回拉时间").
			Optional().
			Nillable(),

		field.Uint32("publish_attempts").
			Comment("发布投递尝试次数（含首次与 asynq 重试），用于判定最终失败").
			Default(0).
			Optional().
			Nillable(),
	}
}

func (SocialDistribution) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.AutoIncrementId{},
		mixin.TimeAt{},
		mixin.OperatorID{},
		mixin.TenantID[uint32]{},
	}
}

func (SocialDistribution) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("tenant_id", "status").StorageKey("idx_social_distributions_tenant_status"),
		index.Fields("account_id").StorageKey("idx_social_distributions_account"),
		index.Fields("post_id").StorageKey("idx_social_distributions_post"),
		index.Fields("platform").StorageKey("idx_social_distributions_platform"),
		index.Fields("publish_time").StorageKey("idx_social_distributions_publish_time"),
	}
}
