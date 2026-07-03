package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/tx7do/go-crud/entgo/mixin"
)

// SectionTranslation holds the schema definition for the SectionTranslation entity.
type SectionTranslation struct {
	ent.Schema
}

func (SectionTranslation) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "section_translations",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
		entsql.WithComments(true),
		schema.Comment("区块翻译表"),
	}
}

// Fields of the SectionTranslation.
func (SectionTranslation) Fields() []ent.Field {
	return []ent.Field{
		field.Uint32("section_id").
			Comment("关联的区块ID").
			Optional().
			Nillable(),

		field.String("language_code").
			Comment("语言代码").
			Optional().
			Nillable(),

		field.JSON("content", &map[string]string{}).
			Comment("区块内容（键值对，适配不同类型区块，随语言变化）").
			Optional(),
	}
}

// Mixin of the SectionTranslation.
func (SectionTranslation) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.AutoIncrementId{},
		mixin.TimeAt{},
		mixin.OperatorID{},
	}
}

func (SectionTranslation) Indexes() []ent.Index {
	return []ent.Index{
		// 单字段索引，用于按区块查询翻译
		index.Fields("section_id"),
		// 单字段索引，用于按语言代码查询翻译
		index.Fields("language_code"),
		// 复合索引，优化按区块和语言查询特定翻译版本
		index.Fields("section_id", "language_code").Unique(),
	}
}
