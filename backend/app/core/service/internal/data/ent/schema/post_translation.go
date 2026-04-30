package schema

import (
	appMixin "go-wind-cms/pkg/entgo/mixin"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/tx7do/go-crud/entgo/mixin"
)

// PostTranslation holds the schema definition for the PostTranslation entity.
type PostTranslation struct {
	ent.Schema
}

func (PostTranslation) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "post_translations",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
		entsql.WithComments(true),
		schema.Comment("帖子翻译表"),
	}
}

// Fields of the PostTranslation.
func (PostTranslation) Fields() []ent.Field {
	return []ent.Field{
		field.Uint32("post_id").
			Comment("关联的帖子ID").
			Optional().
			Nillable(),

		field.String("language_code").
			Comment("语言代码").
			Optional().
			Nillable(),

		field.String("title").
			Comment("帖子标题").
			Optional().
			Nillable(),

		field.String("slug").
			Comment("语言特定 slug").
			Optional().
			Nillable(),

		field.String("summary").
			Comment("帖子摘要").
			Optional().
			Nillable(),

		field.String("content").
			Comment("帖子内容").
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

		field.String("full_path").
			Comment("完整路径").
			Optional().
			Nillable(),

		field.Uint32("word_count").
			Comment("当前语言版本的字数").
			Default(0).
			Optional().
			Nillable(),
	}
}

// Mixin of the PostTranslation.
func (PostTranslation) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.AutoIncrementId{},
		mixin.TimeAt{},
		mixin.OperatorID{},
		appMixin.Seo{},
	}
}

func (PostTranslation) Indexes() []ent.Index {
	return []ent.Index{
		// 单字段索引，用于按关联帖子查询翻译
		index.Fields("post_id"),
		// 单字段索引，用于按语言代码查询翻译
		index.Fields("language_code"),
		// 单字段索引，优化URL友好别名的查询
		index.Fields("slug"),
		// 单字段索引，优化完整路径的查询
		index.Fields("full_path"),
		// 复合唯一索引，确保同一帖子的同一语言下 slug 唯一（这是最重要的约束）
		index.Fields("post_id", "language_code", "slug").
			Unique(),
		// 复合索引，优化按帖子和语言查询特定翻译版本
		index.Fields("post_id", "language_code"),
		// 复合索引，优化按语言代码和slug查询（用于 URL 路由）
		index.Fields("language_code", "slug"),
	}
}
