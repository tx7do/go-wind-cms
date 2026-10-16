package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/tx7do/go-crud/entgo/mixin"
)

// ContentReference holds the schema definition for the ContentReference entity.
type ContentReference struct {
	ent.Schema
}

func (ContentReference) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "content_references",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
		entsql.WithComments(true),
		schema.Comment("内容引用关系表（引用图谱的边：文章/页面 → 媒体等目标资源）"),
	}
}

// Fields of the ContentReference.
func (ContentReference) Fields() []ent.Field {
	return []ent.Field{
		field.String("source_type").
			Comment("源内容类型（post/page）").
			NotEmpty().
			MaxLen(32),
		field.Uint32("source_id").
			Comment("源内容 ID"),
		field.String("target_type").
			Comment("目标资源类型（media_asset）").
			NotEmpty().
			MaxLen(32),
		field.Uint32("target_id").
			Comment("目标资源 ID"),
		field.String("field").
			Comment("引用位置（body=正文内嵌，thumbnail=缩略图字段）").
			NotEmpty().
			MaxLen(32).
			Default("body"),
	}
}

// Mixin of the ContentReference.
func (ContentReference) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.AutoIncrementId{},
		mixin.TimeAt{},
		mixin.OperatorID{},
		mixin.TenantID[uint32]{},
	}
}

// Indexes of the ContentReference.
func (ContentReference) Indexes() []ent.Index {
	return []ent.Index{
		// 源维度：重存内容时按源整体重建引用行
		index.Fields("source_type", "source_id"),
		// 目标维度：查"哪些内容在用此媒体"
		index.Fields("target_type", "target_id"),
		// 幂等：同一源同一字段同一目标唯一
		index.Fields("source_type", "source_id", "target_type", "target_id", "field").
			Unique(),
	}
}
