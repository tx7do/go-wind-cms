package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/tx7do/go-crud/entgo/mixin"
)

// InteractionCounter holds the schema definition for the InteractionCounter entity.
//
// 交互计数内聚表：以 (tenant, target_type, target_id, metric) 为唯一键存储累计计数。
// 取代散落在 post.likes / comment.like_count 上的缓存列，由 InteractionService
// 作为唯一写入方在 ledger 操作的同一事务内 upsert。复合 unique 索引保证单目标
// 单 metric 仅一行，兼作查询索引。
type InteractionCounter struct {
	ent.Schema
}

func (InteractionCounter) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "interaction_counters",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
		entsql.WithComments(true),
		schema.Comment("交互计数内聚表"),
	}
}

func (InteractionCounter) Fields() []ent.Field {
	return []ent.Field{
		field.Uint8("target_type").
			Comment("目标类型，对应 interaction.TargetType 枚举").
			Optional().
			Nillable(),
		field.Uint32("target_id").
			Comment("目标ID").
			Optional().
			Nillable(),
		field.Uint8("metric").
			Comment("计数指标，对应 interaction.CounterMetric 枚举").
			Optional().
			Nillable(),
		field.Int64("count").
			Comment("累计计数").
			Default(0).
			Optional().
			Nillable(),
	}
}

func (InteractionCounter) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.AutoIncrementId{},
		mixin.TimeAt{},
		mixin.TenantID[uint32]{},
	}
}

func (InteractionCounter) Indexes() []ent.Index {
	return []ent.Index{
		// 复合 unique 索引：单目标单 metric 仅一行，兼作查询索引
		index.Fields("tenant_id", "target_type", "target_id", "metric").
			Unique(),
	}
}
