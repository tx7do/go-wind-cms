package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/tx7do/go-crud/entgo/mixin"
)

// FieldDefinitionTranslation 字段定义 label/description/placeholder 的多语言翻译
// （存储模式参照 SectionTranslation）。
type FieldDefinitionTranslation struct {
	ent.Schema
}

func (FieldDefinitionTranslation) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "field_definition_translations",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
		entsql.WithComments(true),
		schema.Comment("字段定义翻译表"),
	}
}

func (FieldDefinitionTranslation) Fields() []ent.Field {
	return []ent.Field{
		field.Uint32("field_definition_id").
			Comment("字段定义ID").
			Optional().
			Nillable(),

		field.String("language_code").
			Comment("语言代码（BCP 47，如 zh-CN）").
			NotEmpty().
			Optional().
			Nillable(),

		field.String("label").
			Comment("翻译后字段标签").
			Optional().
			Nillable(),

		field.String("description").
			Comment("翻译后字段描述").
			Optional().
			Nillable(),

		field.String("placeholder").
			Comment("翻译后占位符").
			Optional().
			Nillable(),
	}
}

func (FieldDefinitionTranslation) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.AutoIncrementId{},
		mixin.TimeAt{},
		mixin.OperatorID{},
	}
}

func (FieldDefinitionTranslation) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("field_definition_id", "language_code").Unique().StorageKey("uix_field_definition_translation_field_lang"),
		index.Fields("field_definition_id"),
	}
}
