package mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"

	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
)

var _ ent.Mixin = (*Seo)(nil)

type Seo struct{ mixin.Schema }

func (Seo) Fields() []ent.Field {
	return []ent.Field{
		field.JSON("seo", &contentV1.SeoMeta{}).
			Comment("SEO 结构化元数据").
			Optional(),
	}
}
