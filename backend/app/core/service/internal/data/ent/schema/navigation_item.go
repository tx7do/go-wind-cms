package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/tx7do/go-crud/entgo/mixin"
)

// NavigationItem holds the schema definition for the NavigationItem entity.
type NavigationItem struct {
	ent.Schema
}

func (NavigationItem) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "navigation_items",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
		entsql.WithComments(true),
		schema.Comment("导航项表"),
	}
}

// Fields of the NavigationItem.
func (NavigationItem) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("link_type").
			Comment("链接类型").
			NamedValues(
				"LinkTypeCustom", "LINK_TYPE_CUSTOM",
				"LinkTypePost", "LINK_TYPE_POST",
				"LinkTypePage", "LINK_TYPE_PAGE",
				"LinkTypeCategory", "LINK_TYPE_CATEGORY",
				"LinkTypeExternal", "LINK_TYPE_EXTERNAL",
			).
			Default("LINK_TYPE_POST").
			Optional().
			Nillable(),

		field.Uint32("navigation_id").
			Comment("所属导航菜单ID").
			Optional().
			Nillable(),

		field.String("title").
			Comment("显示文本").
			Optional().
			Nillable(),

		field.String("url").
			Comment("目标 URL").
			Optional().
			Nillable(),

		field.Uint32("object_id").
			Comment("关联对象ID").
			Optional().
			Nillable(),

		field.String("icon").
			Comment("图标").
			Optional().
			Nillable(),

		field.String("description").
			Comment("描述文本").
			Optional().
			Nillable(),

		field.Bool("is_open_new_tab").
			Comment("是否在新标签页打开").
			Default(false).
			Optional().
			Nillable(),

		field.Bool("is_invalid").
			Comment("是否无效").
			Default(false).
			Optional().
			Nillable(),

		field.String("css_class").
			Comment("额外的 CSS 类").
			Optional().
			Nillable(),

		field.String("required_permission").
			Comment("访问权限标识").
			Optional().
			Nillable(),
	}
}

// Mixin of the NavigationItem.
func (NavigationItem) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.AutoIncrementId{},
		mixin.TimeAt{},
		mixin.OperatorID{},
		mixin.SortOrder{},
		mixin.Tree[NavigationItem]{},
	}
}

func (NavigationItem) Indexes() []ent.Index {
	return []ent.Index{
		// 单字段索引，用于按所属导航菜单查询导航项
		index.Fields("navigation_id"),
		// 单字段索引，用于按链接类型查询
		index.Fields("link_type"),
		// 单字段索引，用于树形结构中按父节点查询子项
		index.Fields("parent_id"),
		// 单字段索引，用于按关联对象ID查询
		index.Fields("object_id"),
		// 单字段索引，用于按无效状态过滤
		index.Fields("is_invalid"),
		// 复合索引，优化按导航菜单和链接类型查询
		index.Fields("navigation_id", "link_type"),
		// 复合索引，优化按导航菜单和无效状态查询
		index.Fields("navigation_id", "is_invalid"),
	}
}
