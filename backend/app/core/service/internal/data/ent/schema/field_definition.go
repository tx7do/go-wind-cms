package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/tx7do/go-crud/entgo/mixin"
)

// RelationConfig relation 类型字段的引用配置（typed-JSON 存储）。
// TargetEntityType 为白名单：post/page/category/tag。
type RelationConfig struct {
	TargetEntityType string `json:"target_entity_type"`
	AllowCrossTenant bool   `json:"allow_cross_tenant"`
	FilterCategoryID uint32 `json:"filter_category_id,omitempty"`
}

// FieldDefinition 字段定义：归属 ContentModel，声明一个自定义字段的
// 名称/类型/校验规则。值存于内容实体（Post/Page/Category）的 custom_fields
// JSON 列（键为 name，image/file/relation 值为字符串化引用，如 "post:456"）。
// 列模式克隆自 SiteSetting。
type FieldDefinition struct {
	ent.Schema
}

func (FieldDefinition) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "field_definitions",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
		entsql.WithComments(true),
		schema.Comment("字段定义表"),
	}
}

func (FieldDefinition) Fields() []ent.Field {
	return []ent.Field{
		field.Uint32("content_model_id").
			Comment("所属内容模型ID").
			Optional().
			Nillable(),

		field.String("name").
			Comment("字段名（语言无关，作为 custom_fields 的键）").
			NotEmpty().
			Optional().
			Nillable(),

		field.Enum("type").
			Comment("字段类型").
			NamedValues(
				"FIELD_TYPE_TEXT", "FIELD_TYPE_TEXT",
				"FIELD_TYPE_NUMBER", "FIELD_TYPE_NUMBER",
				"FIELD_TYPE_RICHTEXT", "FIELD_TYPE_RICHTEXT",
				"FIELD_TYPE_IMAGE", "FIELD_TYPE_IMAGE",
				"FIELD_TYPE_FILE", "FIELD_TYPE_FILE",
				"FIELD_TYPE_RELATION", "FIELD_TYPE_RELATION",
			).
			Default("FIELD_TYPE_TEXT").
			Optional().
			Nillable(),

		field.String("label").
			Comment("字段标签（默认语言）").
			Optional().
			Nillable(),

		field.String("description").
			Comment("字段描述（默认语言）").
			Optional().
			Nillable(),

		field.String("placeholder").
			Comment("输入框占位符（默认语言）").
			Optional().
			Nillable(),

		field.Bool("is_required").
			Comment("是否必填").
			Default(false).
			Optional().
			Nillable(),

		field.String("validation_regex").
			Comment("校验正则（仅 text 类适用）").
			Optional().
			Nillable(),

		field.JSON("options", &map[string]string{}).
			Comment("扩展选项（select 类等，预留）").
			Optional(),

		field.JSON("relation_config", &RelationConfig{}).
			Comment("relation 类型字段的引用配置（仅 FIELD_TYPE_RELATION 适用）").
			Optional(),
	}
}

func (FieldDefinition) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.AutoIncrementId{},
		mixin.TimeAt{},
		mixin.OperatorID{},
		mixin.SortOrder{},
		mixin.TenantID[uint32]{},
	}
}

func (FieldDefinition) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("content_model_id", "name").Unique().StorageKey("uix_field_definition_model_name"),
		index.Fields("content_model_id"),
	}
}
