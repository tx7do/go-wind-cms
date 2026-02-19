package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/tx7do/go-crud/entgo/mixin"
)

// PageTranslation holds the schema definition for the PageTranslation entity.
type PageTranslation struct {
	ent.Schema
}

func (PageTranslation) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "page_translations",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
		entsql.WithComments(true),
		schema.Comment("页面翻译表"),
	}
}

// Fields of the PageTranslation.
func (PageTranslation) Fields() []ent.Field {
	return []ent.Field{
		field.Uint32("page_id").
			Comment("关联的页面ID").
			Optional().
			Nillable(),

		field.String("language_code").
			Comment("语言代码").
			Optional().
			Nillable(),

		field.String("title").
			Comment("页面标题").
			Optional().
			Nillable(),

		field.String("slug").
			Comment("语言特定 slug").
			Optional().
			Nillable(),

		field.String("summary").
			Comment("页面摘要").
			Optional().
			Nillable(),

		field.String("content").
			Comment("页面内容").
			Optional().
			Nillable(),

		field.String("original_content").
			Comment("原始内容").
			Optional().
			Nillable(),

		field.String("thumbnail").
			Comment("缩略图").
			Optional().
			Nillable(),

		field.String("cover_image").
			Comment("封面图").
			Optional().
			Nillable(),

		field.String("full_path").
			Comment("完整路径").
			Optional().
			Nillable(),

		field.Uint32("word_count").
			Comment("当前语言版本的字数").
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

// Mixin of the PageTranslation.
func (PageTranslation) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.AutoIncrementId{},
		mixin.TimeAt{},
		mixin.OperatorID{},
	}
}

func (PageTranslation) Indexes() []ent.Index {
	return []ent.Index{
		// 单字段索引，用于按关联页面查询翻译
		index.Fields("page_id"),
		// 单字段索引，用于按语言代码查询翻译
		index.Fields("language_code"),
		// 单字段索引，优化URL友好别名的查询
		index.Fields("slug"),
		// 复合索引，优化按页面和语言查询特定翻译版本
		index.Fields("page_id", "language_code"),
		// 复合索引，优化按语言代码和slug查询
		index.Fields("language_code", "slug"),
		// 单字段索引，优化SEO关键词搜索
		index.Fields("meta_keywords"),
	}
}
