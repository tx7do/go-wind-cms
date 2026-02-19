package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/tx7do/go-crud/entgo/mixin"
)

// PostCategory holds the schema definition for the PostCategory entity.
type PostCategory struct {
	ent.Schema
}

func (PostCategory) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "post_categories",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
		entsql.WithComments(true),
		schema.Comment("帖子分类关联表"),
	}
}

// Fields of the PostCategory.
func (PostCategory) Fields() []ent.Field {
	return []ent.Field{
		field.Uint32("post_id").
			Comment("帖子ID").
			Nillable(),

		field.Uint32("category_id").
			Comment("分类ID").
			Nillable(),
	}
}

// Mixin of the PostCategory.
func (PostCategory) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.AutoIncrementId{},
		mixin.CreatedAt{},
	}
}

func (PostCategory) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("post_id", "category_id").Unique(),
		index.Fields("post_id"),
		index.Fields("category_id"),
	}
}
