package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/tx7do/go-crud/entgo/mixin"
)

// ExperimentParticipant holds the schema definition for the ExperimentParticipant entity.
// 实验参与记录：访客首次被分配时落一条（(experiment_id, visitor_hash) 幂等），
// 变体名快照存储（变体改名不影响历史归因）。
type ExperimentParticipant struct {
	ent.Schema
}

func (ExperimentParticipant) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "experiment_participants",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
		entsql.WithComments(true),
		schema.Comment("内容实验参与记录表"),
	}
}

// Fields of the ExperimentParticipant.
func (ExperimentParticipant) Fields() []ent.Field {
	return []ent.Field{
		field.Uint32("experiment_id").
			Comment("实验 ID").
			Optional().
			Nillable(),
		field.String("visitor_hash").
			Comment("访客标识哈希（cookie vid + 实验 ID 派生，不落原始身份）").
			NotEmpty(),
		field.String("variant_name").
			Comment("被分配的变体名称快照").
			NotEmpty(),
		field.Bool("converted").
			Comment("是否转化（如表单提交归因）").
			Default(false).
			Optional().
			Nillable(),
		field.Time("converted_at").
			Comment("转化时间").
			Optional().
			Nillable(),
	}
}

// Mixin of the ExperimentParticipant.
func (ExperimentParticipant) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.AutoIncrementId{},
		mixin.TimeAt{},
		mixin.TenantID[uint32]{},
	}
}

func (ExperimentParticipant) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("experiment_id", "visitor_hash").Unique(),
		index.Fields("experiment_id", "variant_name"),
		index.Fields("experiment_id", "converted"),
	}
}
