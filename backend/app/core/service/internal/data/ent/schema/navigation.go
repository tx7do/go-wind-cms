package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/tx7do/go-crud/entgo/mixin"
)

// Navigation holds the schema definition for the Navigation entity.
type Navigation struct {
	ent.Schema
}

func (Navigation) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "navigations",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
		entsql.WithComments(true),
		schema.Comment("导航表"),
	}
}

// Fields of the Navigation.
func (Navigation) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Comment("导航名称").
			Optional().
			Nillable(),

		field.String("location").
			Comment("渲染位置标识").
			Optional().
			Nillable(),

		field.String("locale").
			Comment("关联的语言区域").
			Optional().
			Nillable(),

		field.Bool("is_active").
			Comment("是否启用").
			Default(true).
			Optional().
			Nillable(),
	}
}

// Mixin of the Navigation.
func (Navigation) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.AutoIncrementId{},
		mixin.TimeAt{},
		mixin.OperatorID{},
	}
}

func (Navigation) Indexes() []ent.Index {
	return []ent.Index{
		// 复合索引，优化按渲染位置和语言区域查询导航
		index.Fields("location", "locale"),
		// 单字段索引，用于按启用状态查询导航
		index.Fields("is_active"),
		// 单字段索引，用于按语言区域查询导航
		index.Fields("locale"), // 复合索引，优化按渲染位置、语言区域和启用状态查询
		index.Fields("location", "locale", "is_active"),
	}
}
