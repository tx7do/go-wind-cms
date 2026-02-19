package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/tx7do/go-crud/entgo/mixin"
)

// CategoryTranslation holds the schema definition for the CategoryTranslation entity.
type CategoryTranslation struct {
	ent.Schema
}

func (CategoryTranslation) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "category_translations",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
		entsql.WithComments(true),
		schema.Comment("分类翻译表"),
	}
}

// Fields of the CategoryTranslation.
func (CategoryTranslation) Fields() []ent.Field {
	return []ent.Field{
		field.Uint32("category_id").
			Comment("关联的分类ID").
			Optional().
			Nillable(),

		field.Bool("is_nav").
			Comment("是否显示在导航菜单").
			Optional().
			Nillable(),

		field.String("language_code").
			Comment("语言代码").
			Optional().
			Nillable(),

		field.String("name").
			Comment("分类名称").
			Optional().
			Nillable(),

		field.String("slug").
			Comment("分类别名").
			Optional().
			Nillable(),

		field.String("description").
			Comment("分类描述").
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

		field.String("template").
			Comment("模板名称").
			Optional().
			Nillable(),

		field.String("full_path").
			Comment("完整路径").
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

// Mixin of the CategoryTranslation.
func (CategoryTranslation) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.AutoIncrementId{},
		mixin.TimeAt{},
		mixin.OperatorID{},
	}
}

func (CategoryTranslation) Indexes() []ent.Index {
	return []ent.Index{
		// 复合唯一索引，确保同一分类的同一语言只有一条翻译记录
		index.Fields("category_id", "language_code").Unique(),
		// 单字段索引，用于按分类ID查询其所有语言的翻译
		index.Fields("category_id"),
		// 单字段索引，用于按语言代码查询翻译
		index.Fields("language_code"),
		// 单字段索引，用于按导航菜单状态查询
		index.Fields("is_nav"),
		// 单字段索引，优化SEO相关的搜索查询
		index.Fields("slug"),
	}
}
