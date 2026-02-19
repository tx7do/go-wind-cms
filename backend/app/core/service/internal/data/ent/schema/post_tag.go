package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/tx7do/go-crud/entgo/mixin"
)

// PostTag holds the schema definition for the PostTag entity.
type PostTag struct {
	ent.Schema
}

func (PostTag) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "post_tags",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
		entsql.WithComments(true),
		schema.Comment("帖子标签关联表"),
	}
}

// Fields of the PostTag.
func (PostTag) Fields() []ent.Field {
	return []ent.Field{
		field.Uint32("post_id").
			Comment("帖子ID").
			Nillable(),

		field.Uint32("tag_id").
			Comment("标签ID").
			Nillable(),
	}
}

// Mixin of the PostTag.
func (PostTag) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.AutoIncrementId{},
		mixin.CreatedAt{},
	}
}

func (PostTag) Indexes() []ent.Index {
	return []ent.Index{
		// 复合唯一索引，确保同一篇帖子的同一个标签只关联一次
		index.Fields("post_id", "tag_id").Unique(),
		// 单字段索引，用于按帖子查询其所有标签
		index.Fields("post_id"),
		// 单字段索引，用于按标签查询其关联的所有帖子
		index.Fields("tag_id"),
	}
}
