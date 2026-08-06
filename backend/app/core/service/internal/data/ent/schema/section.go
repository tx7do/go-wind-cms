package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/tx7do/go-crud/entgo/mixin"
)

// Section holds the schema definition for the Section entity.
type Section struct {
	ent.Schema
}

func (Section) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "sections",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
		entsql.WithComments(true),
		schema.Comment("页面区块表"),
	}
}

// Fields of the Section.
func (Section) Fields() []ent.Field {
	return []ent.Field{
		field.Uint32("page_id").
			Comment("所属页面ID").
			Optional().
			Nillable(),

		field.Enum("type").
			Comment("区块类型").
			NamedValues(
				"SectionTypeRichText", "SECTION_TYPE_RICH_TEXT",
				"SectionTypeMarkdown", "SECTION_TYPE_MARKDOWN",
				"SectionTypeTitle", "SECTION_TYPE_TITLE",
				"SectionTypeImage", "SECTION_TYPE_IMAGE",
				"SectionTypeGallery", "SECTION_TYPE_GALLERY",
				"SectionTypeVideo", "SECTION_TYPE_VIDEO",
				"SectionTypeButton", "SECTION_TYPE_BUTTON",
				"SectionTypeDivider", "SECTION_TYPE_DIVIDER",
				"SectionTypeSpacer", "SECTION_TYPE_SPACER",
				"SectionTypeCode", "SECTION_TYPE_CODE",
				"SectionTypeHtml", "SECTION_TYPE_HTML",
				"SectionTypeForm", "SECTION_TYPE_FORM",
				"SectionTypeCarousel", "SECTION_TYPE_CAROUSEL",
				"SectionTypeCustom", "SECTION_TYPE_CUSTOM",
			).
			Default("SECTION_TYPE_RICH_TEXT").
			Optional().
			Nillable(),

		field.String("name").
			Comment("区块名称（后台标识用，语言无关）").
			Optional().
			Nillable(),

		field.JSON("config", &map[string]string{}).
			Comment("区块样式/布局配置（边距、CSS class、列数等，语言无关）").
			Optional(),
	}
}

// Mixin of the Section.
func (Section) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.AutoIncrementId{},
		mixin.TimeAt{},
		mixin.OperatorID{},
		mixin.SortOrder{},
		mixin.TenantID[uint32]{},
	}
}

func (Section) Indexes() []ent.Index {
	return []ent.Index{
		// 单字段索引，用于按页面查询其所有区块
		index.Fields("page_id"),
		// 复合索引，优化按页面查询并排序
		index.Fields("page_id", "sort_order"),
	}
}
