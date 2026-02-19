package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/tx7do/go-crud/entgo/mixin"
)

// Category holds the schema definition for the Category entity.
type Category struct {
	ent.Schema
}

func (Category) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "categories",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
		entsql.WithComments(true),
		schema.Comment("类别表"),
	}
}

// Fields of the Category.
func (Category) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("status").
			Comment("分类状态").
			NamedValues(
				"CategoryStatusActive", "CATEGORY_STATUS_ACTIVE",
				"CategoryStatusHidden", "CATEGORY_STATUS_HIDDEN",
				"CategoryStatusArchived", "CATEGORY_STATUS_ARCHIVED",
			).
			Optional().
			Nillable(),

		field.Int32("depth").
			Comment("分类层级深度").
			Optional().
			Nillable(),

		field.Bool("is_nav").
			Comment("是否显示在导航菜单").
			Optional().
			Nillable(),

		field.String("icon").
			Comment("分类图标").
			Optional().
			Nillable(),

		field.Uint32("post_count").
			Comment("该分类下的文章总数").
			Optional().
			Nillable(),

		field.Uint32("direct_post_count").
			Comment("该分类下的直接文章数").
			Optional().
			Nillable(),

		field.JSON("custom_fields", &map[string]string{}).
			Comment("自定义字段").
			Optional(),
	}
}

// Mixin of the Category.
func (Category) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.AutoIncrementId{},
		mixin.TimeAt{},
		mixin.OperatorID{},
		mixin.SortOrder{},
		mixin.TreePath{},
		mixin.Tree[Category]{},
	}
}

func (Category) Indexes() []ent.Index {
	return []ent.Index{
		// 单字段索引，用于按状态查询分类
		index.Fields("status"),
		// 单字段索引，用于查询导航菜单中的分类
		index.Fields("is_nav"),
		// 单字段索引，用于树形关系中按父节点查询子分类
		index.Fields("parent_id"),
	}
}
