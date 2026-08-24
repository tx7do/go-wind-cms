package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/tx7do/go-crud/entgo/mixin"
)

// ContentModelTranslation 内容模型名称/描述的多语言翻译
// （存储模式参照 SectionTranslation：按 (content_model_id, language_code) 唯一）。
type ContentModelTranslation struct {
	ent.Schema
}

func (ContentModelTranslation) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "content_model_translations",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
		entsql.WithComments(true),
		schema.Comment("内容模型翻译表"),
	}
}

func (ContentModelTranslation) Fields() []ent.Field {
	return []ent.Field{
		field.Uint32("content_model_id").
			Comment("内容模型ID").
			Optional().
			Nillable(),

		field.String("language_code").
			Comment("语言代码（BCP 47，如 zh-CN）").
			NotEmpty().
			Optional().
			Nillable(),

		field.String("name").
			Comment("翻译后模型名称").
			Optional().
			Nillable(),

		field.String("description").
			Comment("翻译后模型描述").
			Optional().
			Nillable(),
	}
}

func (ContentModelTranslation) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.AutoIncrementId{},
		mixin.TimeAt{},
		mixin.OperatorID{},
	}
}

func (ContentModelTranslation) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("content_model_id", "language_code").Unique().StorageKey("uix_content_model_translation_model_lang"),
		index.Fields("content_model_id"),
	}
}
