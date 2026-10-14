package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/tx7do/go-crud/entgo/mixin"
)

// SocialStatsHistory 社交分发互动统计历史快照：统计回拉成功时逐次追加一行，
// 供将来绘制互动趋势（7/30/90 天）使用。行只通过 distribution_id 访问
// （租户过滤由父表分发记录承担），故不挂 TenantID mixin——与
// go-wind-admin 通知域的全局表同一设计取舍。
type SocialStatsHistory struct{ ent.Schema }

func (SocialStatsHistory) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "social_stats_history", Charset: "utf8mb4", Collation: "utf8mb4_bin"},
		entsql.WithComments(true),
		schema.Comment("社交分发统计历史快照"),
	}
}

func (SocialStatsHistory) Fields() []ent.Field {
	return []ent.Field{
		field.Uint32("distribution_id").
			Comment("分发记录ID").
			Optional().
			Nillable(),

		field.Uint32("likes").
			Comment("点赞数快照").
			Default(0).
			Optional().
			Nillable(),

		field.Uint32("replies").
			Comment("回复数快照").
			Default(0).
			Optional().
			Nillable(),

		field.Uint32("reposts").
			Comment("转发数快照").
			Default(0).
			Optional().
			Nillable(),
	}
}

func (SocialStatsHistory) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.AutoIncrementId{},
		mixin.TimeAt{}, // created_at 即快照时间
	}
}

func (SocialStatsHistory) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("distribution_id", "created_at").StorageKey("idx_social_stats_history_dist_time"),
		index.Fields("created_at").StorageKey("idx_social_stats_history_created_at"),
	}
}
