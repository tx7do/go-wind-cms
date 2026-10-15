package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/tx7do/go-crud/entgo/mixin"
)

// FormSubmission holds the schema definition for the FormSubmission entity.
type FormSubmission struct {
	ent.Schema
}

func (FormSubmission) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "form_submissions",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
		entsql.WithComments(true),
		schema.Comment("表单提交表（页面 FORM 区块的访客提交数据，只增不改）"),
	}
}

// Fields of the FormSubmission.
func (FormSubmission) Fields() []ent.Field {
	return []ent.Field{
		field.Uint32("page_id").
			Comment("所属页面ID").
			Optional().
			Nillable(),

		field.Uint32("section_id").
			Comment("所属区块ID（FORM 类型 Section）").
			Optional().
			Nillable(),

		field.JSON("data", &map[string]string{}).
			Comment("提交数据（字段名→值；字段集须为该区块 config.fields 的子集）").
			Optional(),

		field.Uint32("user_id").
			Comment("提交者用户ID（匿名提交为 0）").
			Default(0).
			Optional().
			Nillable(),
	}
}

// Mixin of the FormSubmission.
func (FormSubmission) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.AutoIncrementId{},
		mixin.TimeAt{},
		mixin.TenantID[uint32]{},
	}
}

func (FormSubmission) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("page_id", "created_at"),
		index.Fields("section_id"),
	}
}
