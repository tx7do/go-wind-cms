package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/tx7do/go-crud/entgo/mixin"
)

// MediaVariant holds the schema definition for the MediaVariant entity.
type MediaVariant struct {
	ent.Schema
}

func (MediaVariant) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "media_variants",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
		entsql.WithComments(true),
		schema.Comment("媒体变体表，关联媒体文件和标签"),
	}
}

// Fields of the MediaVariant.
func (MediaVariant) Fields() []ent.Field {
	return []ent.Field{
		field.Uint32("media_id").
			Comment("媒体资源ID").
			Nillable(),

		field.Uint32("file_id").
			Comment("文件ID").
			Nillable(),

		field.Uint32("variant_name").
			Comment("变体名称").
			Optional().
			Nillable(),
	}
}

// Mixin of the MediaVariant.
func (MediaVariant) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.AutoIncrementId{},
		mixin.CreatedAt{},
	}
}

func (MediaVariant) Indexes() []ent.Index {
	return []ent.Index{
		// 复合唯一索引，确保同一媒体资源的同一个文件变体只关联一次
		index.Fields("media_id", "file_id").Unique(),
		// 单字段索引，用于按媒体资源查询其所有变体
		index.Fields("media_id"),
		// 单字段索引，用于按文件查询其所有变体
		index.Fields("file_id"),
	}
}
