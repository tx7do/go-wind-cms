package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/tx7do/go-crud/entgo/mixin"
)

// Tag holds the schema definition for the Tag entity.
type Tag struct {
	ent.Schema
}

func (Tag) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "tags",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
		entsql.WithComments(true),
		schema.Comment("标签表"),
	}
}

// Fields of the Tag.
func (Tag) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("status").
			Comment("标签状态").
			NamedValues(
				"TAG_STATUS_ACTIVE", "TAG_STATUS_ACTIVE",
				"TAG_STATUS_HIDDEN", "TAG_STATUS_HIDDEN",
				"TAG_STATUS_ARCHIVED", "TAG_STATUS_ARCHIVED",
			).
			Default("TAG_STATUS_ACTIVE").
			Optional().
			Nillable(),

		field.String("color").
			Comment("标签颜色").
			Optional().
			Nillable(),

		field.String("icon").
			Comment("标签图标").
			Optional().
			Nillable(),

		field.String("group").
			Comment("标签分组").
			Optional().
			Nillable(),

		field.Bool("is_featured").
			Comment("是否推荐").
			Default(false).
			Optional().
			Nillable(),

		field.Uint32("post_count").
			Comment("使用该标签的文章总数").
			Optional().
			Nillable(),
	}
}

// Mixin of the Tag.
func (Tag) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.AutoIncrementId{},
		mixin.TimeAt{},
		mixin.OperatorID{},
		mixin.SortOrder{},
	}
}

func (Tag) Indexes() []ent.Index {
	return []ent.Index{
		// 单字段索引，用于按标签状态查询
		index.Fields("status"),
		// 单字段索引，用于按标签分组查询
		index.Fields("group"),
		// 单字段索引，用于按是否推荐过滤
		index.Fields("is_featured"),
		// 复合索引，优化按状态和分组查询
		index.Fields("status", "group"),
		// 复合索引，优化按状态和推荐状态查询
		index.Fields("status", "is_featured"),
	}
}
