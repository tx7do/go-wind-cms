package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/tx7do/go-crud/entgo/mixin"
)

// Experiment holds the schema definition for the Experiment entity.
// 内容实验（A/B 分流 + 受众定向）：挂在某页面上的运行时实验，访客按
// 随机权重或受众规则被分配到变体，app 侧以变体 sections 覆盖页面渲染。
type Experiment struct {
	ent.Schema
}

func (Experiment) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "experiments",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
		entsql.WithComments(true),
		schema.Comment("内容实验表（A/B 测试与受众定向）"),
	}
}

// Fields of the Experiment.
func (Experiment) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Comment("实验名称").
			NotEmpty(),
		field.Uint32("page_id").
			Comment("目标页面 ID（实验作用于该页面的 sections）").
			Optional().
			Nillable(),
		field.Enum("status").
			Comment("状态：draft 草稿 / running 运行中 / paused 暂停 / ended 已结束").
			Values("draft", "running", "paused", "ended").
			Default("draft").
			Optional().
			Nillable(),
		field.Enum("assignment_type").
			Comment("分配方式：random 随机权重分流（A/B）/ rule 受众规则定向（个性化）").
			Values("random", "rule").
			Default("random").
			Optional().
			Nillable(),
		field.String("description").
			Comment("实验说明").
			Optional().
			Nillable(),
	}
}

// Mixin of the Experiment.
func (Experiment) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.AutoIncrementId{},
		mixin.TimeAt{},
		mixin.OperatorID{},
		mixin.TenantID[uint32]{},
	}
}

func (Experiment) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("page_id"),
		index.Fields("status"),
	}
}
