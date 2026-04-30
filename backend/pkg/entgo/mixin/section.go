package mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"

	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
)

var _ ent.Mixin = (*Sections)(nil)

type Sections struct{ mixin.Schema }

func (Sections) Fields() []ent.Field {
	return []ent.Field{
		field.JSON("sections", []*contentV1.Section{}).
			Comment("页面内容区块（模块化）").
			Optional(),
	}
}
