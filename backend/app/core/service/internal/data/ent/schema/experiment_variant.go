package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/tx7do/go-crud/entgo/mixin"
)

// ExperimentVariant holds the schema definition for the ExperimentVariant entity.
// 实验变体：一份 sections 快照（protojson），random 型按 weight 权重分流，
// rule 型按 conditions（受众条件，JSON）依序首中定向，未中回退权重分流。
type ExperimentVariant struct {
	ent.Schema
}

func (ExperimentVariant) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "experiment_variants",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
		entsql.WithComments(true),
		schema.Comment("内容实验变体表"),
	}
}

// Fields of the ExperimentVariant.
func (ExperimentVariant) Fields() []ent.Field {
	return []ent.Field{
		field.Uint32("experiment_id").
			Comment("所属实验 ID").
			Optional().
			Nillable(),
		field.String("name").
			Comment("变体名称（如 A / B；参与记录按名称快照）").
			NotEmpty(),
		field.Int("weight").
			Comment("分流权重（百分比 0-100，random 型按权重比例分流）").
			Default(50).
			Optional().
			Nillable(),
		field.Text("conditions").
			Comment("受众条件 JSON（rule 型；数组 [{field,op,value}] 全部满足才命中）").
			Optional(),
		field.Text("sections").
			Comment("该变体的页面区块快照（protojson 数组，命中后整页替换 sections）").
			Optional(),
		field.Uint32("sort_order").
			Comment("排序（rule 型即求值顺序）").
			Default(0).
			Optional().
			Nillable(),
	}
}

// Mixin of the ExperimentVariant.
func (ExperimentVariant) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.AutoIncrementId{},
		mixin.TimeAt{},
		mixin.OperatorID{},
		mixin.TenantID[uint32]{},
	}
}

func (ExperimentVariant) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("experiment_id"),
	}
}
