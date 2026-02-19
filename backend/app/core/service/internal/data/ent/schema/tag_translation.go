package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/tx7do/go-crud/entgo/mixin"
)

// TagTranslation holds the schema definition for the TagTranslation entity.
type TagTranslation struct {
	ent.Schema
}

func (TagTranslation) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "tag_translations",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
		entsql.WithComments(true),
		schema.Comment("标签翻译表"),
	}
}

// Fields of the TagTranslation.
func (TagTranslation) Fields() []ent.Field {
	return []ent.Field{
		field.Uint32("tag_id").
			Comment("关联的标签ID").
			Optional().
			Nillable(),

		field.String("language_code").
			Comment("语言代码").
			Optional().
			Nillable(),

		field.String("name").
			Comment("标签名称").
			Optional().
			Nillable(),

		field.String("slug").
			Comment("语言特定 slug").
			Optional().
			Nillable(),

		field.String("description").
			Comment("标签描述").
			Optional().
			Nillable(),

		field.String("cover_image").
			Comment("封面图").
			Optional().
			Nillable(),

		field.String("template").
			Comment("模板名称").
			Optional().
			Nillable(),

		field.String("full_path").
			Comment("完整路径").
			Optional().
			Nillable(),

		field.String("canonical_url").
			Comment("规范 URL").
			Optional().
			Nillable(),

		field.String("meta_keywords").
			Comment("SEO 关键词").
			Optional().
			Nillable(),

		field.String("meta_description").
			Comment("SEO 描述").
			Optional().
			Nillable(),

		field.String("seo_title").
			Comment("SEO 标题").
			Optional().
			Nillable(),
	}
}

// Mixin of the TagTranslation.
func (TagTranslation) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.AutoIncrementId{},
		mixin.TimeAt{},
		mixin.OperatorID{},
	}
}

func (TagTranslation) Indexes() []ent.Index {
	return []ent.Index{
		// 单字段索引，用于按关联标签查询翻译
		index.Fields("tag_id"),
		// 单字段索引，用于按语言代码查询翻译
		index.Fields("language_code"),
		// 单字段索引，优化URL友好别名的查询
		index.Fields("slug"),
		// 单字段索引，优化完整路径的查询
		index.Fields("full_path"),
		// 单字段索引，优化SEO关键词搜索
		index.Fields("meta_keywords"),
		// 复合索引，优化按标签和语言查询特定翻译版本
		index.Fields("tag_id", "language_code"),
		// 复合索引，优化按语言代码和slug查询
		index.Fields("language_code", "slug"),
	}
}
