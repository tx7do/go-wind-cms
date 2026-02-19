package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/tx7do/go-crud/entgo/mixin"
)

// Page holds the schema definition for the Page entity.
type Page struct {
	ent.Schema
}

func (Page) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "pages",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
		entsql.WithComments(true),
		schema.Comment("页面表"),
	}
}

// Fields of the Page.
func (Page) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("status").
			Comment("页面状态").
			NamedValues(
				"PageStatusDraft", "PAGE_STATUS_DRAFT",
				"PageStatusPublished", "PAGE_STATUS_PUBLISHED",
				"PageStatusArchived", "PAGE_STATUS_ARCHIVED",
			).
			Default("PAGE_STATUS_DRAFT").
			Optional().
			Nillable(),

		field.Enum("type").
			Comment("页面类型").
			NamedValues(
				"PageTypeDefault", "PAGE_TYPE_DEFAULT",
				"PageTypeHome", "PAGE_TYPE_HOME",
				"PageTypeError404", "PAGE_TYPE_ERROR_404",
				"PageTypeError500", "PAGE_TYPE_ERROR_500",
				"PageTypeCustom", "PAGE_TYPE_CUSTOM",
			).
			Default("PAGE_TYPE_HOME").
			Optional().
			Nillable(),

		field.Enum("editor_type").
			Comment("编辑器类型").
			NamedValues(
				"EditorTypeMarkdown", "EDITOR_TYPE_MARKDOWN",
				"EditorTypeRichText", "EDITOR_TYPE_RICH_TEXT",
				"EditorTypeHtml", "EDITOR_TYPE_HTML",
				"EditorTypeJsonBlock", "EDITOR_TYPE_JSON_BLOCK",
			).
			Default("EDITOR_TYPE_MARKDOWN").
			Optional().
			Nillable(),

		field.String("slug").
			Comment("页面别名").
			Optional().
			Nillable(),

		field.Uint32("author_id").
			Comment("评论作者ID，0表示游客").
			Default(0).
			Optional().
			Nillable(),

		field.String("author_name").
			Comment("评论作者名称").
			Optional().
			Nillable(),

		field.Bool("disallow_comment").
			Comment("是否禁止评论").
			Default(false).
			Optional().
			Nillable(),

		field.String("redirect_url").
			Comment("重定向 URL").
			Optional().
			Nillable(),

		field.Bool("show_in_navigation").
			Comment("是否在主导航中显示").
			Default(false).
			Optional().
			Nillable(),

		field.String("template").
			Comment("页面模板名称").
			Optional().
			Nillable(),

		field.Bool("is_custom_template").
			Comment("是否使用自定义模板代码").
			Default(false).
			Optional().
			Nillable(),

		field.Uint32("visits").
			Comment("页面访问次数").
			Default(0).
			Optional().
			Nillable(),

		field.JSON("custom_fields", &map[string]string{}).
			Comment("自定义字段").
			Optional(),

		field.String("custom_head").
			Comment("自定义 head 代码").
			Optional().
			Nillable(),

		field.String("custom_foot").
			Comment("自定义 foot 代码").
			Optional().
			Nillable(),

		field.Int32("depth").
			Comment("页面层级深度").
			Optional().
			Nillable(),
	}
}

// Mixin of the Page.
func (Page) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.AutoIncrementId{},
		mixin.TimeAt{},
		mixin.OperatorID{},
		mixin.SortOrder{},
		mixin.TreePath{},
		mixin.Tree[Page]{},
	}
}

func (Page) Indexes() []ent.Index {
	return []ent.Index{
		// 单字段索引，用于按页面状态查询
		index.Fields("status"),
		// 单字段索引，用于按页面类型查询
		index.Fields("type"),
		// 单字段索引，用于按作者ID查询其所有页面
		index.Fields("author_id"),
		// 单字段索引，用于按页面别名查询
		index.Fields("slug"),
		// 单字段索引，用于按是否显示在导航中查询
		index.Fields("show_in_navigation"),
		// 单字段索引，用于树形结构中按父节点查询子页面
		index.Fields("parent_id"),
		// 单字段索引，用于按禁止评论状态过滤
		index.Fields("disallow_comment"),
		// 复合索引，优化按状态和类型查询
		index.Fields("status", "type"),
		// 复合索引，优化按状态和是否显示在导航中查询
		index.Fields("status", "show_in_navigation"),
	}
}
