package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/tx7do/go-crud/entgo/mixin"
)

// ComponentDefinition 组件定义：具名、可复用的字段组（Components/可重复字段组/
// 动态区的基座，对标 Strapi Components）。字段定义复用 field_definitions 表
// （component_id 可空外键），本表只承载组件身份（code/name）。
// 详见 docs/content-modeling-components-design.md。
type ComponentDefinition struct {
	ent.Schema
}

func (ComponentDefinition) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "component_definitions",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
		entsql.WithComments(true),
		schema.Comment("组件定义表（内容建模 Components）"),
	}
}

func (ComponentDefinition) Fields() []ent.Field {
	return []ent.Field{
		field.String("code").
			Comment("组件编码（语言无关引用键，custom_fields 值里 _component 判别用）").
			NotEmpty().
			MaxLen(64),
		field.String("name").
			Comment("组件显示名").
			NotEmpty().
			MaxLen(128),
		field.String("description").
			Comment("组件描述").
			Optional().
			Nillable(),
	}
}

func (ComponentDefinition) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.AutoIncrementId{},
		mixin.TimeAt{},
		mixin.OperatorID{},
		mixin.TenantID[uint32]{},
	}
}

func (ComponentDefinition) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("code").Unique().StorageKey("uix_component_definition_tenant_code"),
	}
}
