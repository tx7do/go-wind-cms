package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/tx7do/go-crud/entgo/mixin"
)

// SiteSetting holds the schema definition for the SiteSetting entity.
type SiteSetting struct {
	ent.Schema
}

func (SiteSetting) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "site_settings",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
		entsql.WithComments(true),
		schema.Comment("站点设置表"),
	}
}

// Fields of the SiteSetting.
func (SiteSetting) Fields() []ent.Field {
	return []ent.Field{
		field.Uint32("site_id").
			Comment("站点ID").
			Optional().
			Nillable(),

		field.String("locale").
			Comment("语言区域").
			Optional().
			Nillable(),

		field.String("group").
			Comment("配置分组标识").
			Optional().
			Nillable(),

		field.String("key").
			Comment("配置键").
			Optional().
			Nillable(),

		field.String("value").
			Comment("配置值").
			Optional().
			Nillable(),

		field.Enum("type").
			Comment("配置项类型").
			NamedValues(
				"SETTING_TYPE_TEXT", "SETTING_TYPE_TEXT",
				"SETTING_TYPE_TEXTAREA", "SETTING_TYPE_TEXTAREA",
				"SETTING_TYPE_NUMBER", "SETTING_TYPE_NUMBER",
				"SETTING_TYPE_BOOLEAN", "SETTING_TYPE_BOOLEAN",
				"SETTING_TYPE_URL", "SETTING_TYPE_URL",
				"SETTING_TYPE_EMAIL", "SETTING_TYPE_EMAIL",
				"SETTING_TYPE_IMAGE", "SETTING_TYPE_IMAGE",
				"SETTING_TYPE_SELECT", "SETTING_TYPE_SELECT",
				"SETTING_TYPE_JSON", "SETTING_TYPE_JSON",
			).
			Default("SETTING_TYPE_TEXT").
			Optional().
			Nillable(),

		field.String("label").
			Comment("配置项标签").
			Optional().
			Nillable(),

		field.String("description").
			Comment("配置项描述").
			Optional().
			Nillable(),

		field.String("placeholder").
			Comment("输入框占位符").
			Optional().
			Nillable(),

		field.JSON("options", &map[string]string{}).
			Comment("下拉选项").
			Optional(),

		field.Bool("is_required").
			Comment("是否必填").
			Default(false).
			Optional().
			Nillable(),

		field.String("validation_regex").
			Comment("验证正则表达式").
			Optional().
			Nillable(),
	}
}

// Mixin of the SiteSetting.
func (SiteSetting) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.AutoIncrementId{},
		mixin.TimeAt{},
		mixin.OperatorID{},
	}
}

func (SiteSetting) Indexes() []ent.Index {
	return []ent.Index{
		// 单字段索引，用于按站点ID查询配置
		index.Fields("site_id"),
		// 单字段索引，用于按语言区域查询配置
		index.Fields("locale"),
		// 单字段索引，用于按配置分组查询
		index.Fields("group"),
		// 单字段索引，用于按配置键查询
		index.Fields("key"),
		// 单字段索引，用于按配置类型查询
		index.Fields("type"),
		// 复合索引，优化按站点和语言查询配置
		index.Fields("site_id", "locale"),
		// 复合索引，优化按站点、分组和键查询特定配置
		index.Fields("site_id", "group", "key"),
		// 复合索引，优化按语言、分组和键查询配置
		index.Fields("locale", "group", "key"),
	}
}
