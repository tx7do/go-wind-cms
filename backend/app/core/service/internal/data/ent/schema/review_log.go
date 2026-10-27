package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/tx7do/go-crud/entgo/mixin"
)

// ReviewLog 发布审批留痕：submit/approve/reject 全程记录（谁、哪一级、何时、
// 意见）。多级审批的推进依据 = 同实体同批次内最大 level。批次概念：in_progress
// 从 false→true 即开启新批次（level 从 0 的 submit 起算）。
type ReviewLog struct {
	ent.Schema
}

func (ReviewLog) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "sys_review_logs"},
	}
}

// Fields of the ReviewLog.
func (ReviewLog) Fields() []ent.Field {
	return []ent.Field{
		field.String("entity_type").
			Comment("实体类型（post/page）").
			NotEmpty().
			MaxLen(16),
		field.Uint32("entity_id").
			Comment("实体 ID"),
		field.Enum("action").
			Comment("动作").
			NamedValues(
				"Submit", "SUBMIT",
				"Approve", "APPROVE",
				"Reject", "REJECT",
			).
			Default("SUBMIT"),
		field.Uint32("level").
			Comment("审批级别（submit=0；第 N 次 approve=N）"),
		field.String("note").
			Comment("意见（驳回理由/批准备注）").
			Optional().
			Nillable().
			MaxLen(512),
	}
}

// Edges of the ReviewLog.
func (ReviewLog) Edges() []ent.Edge {
	return nil
}

// Mixin of the ReviewLog.
func (ReviewLog) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.TimeAt{},
		mixin.AutoIncrementId{},
		mixin.TenantID[uint32]{},
		mixin.OperatorID{},
	}
}

// Indexes of the ReviewLog.
func (ReviewLog) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("entity_type", "entity_id"),
	}
}
