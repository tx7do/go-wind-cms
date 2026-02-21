package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/tx7do/go-crud/entgo/mixin"
)

// Site holds the schema definition for the Site entity.
type Site struct {
	ent.Schema
}

func (Site) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "sites",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
		entsql.WithComments(true),
		schema.Comment("站点表"),
	}
}

// Fields of the Site.
func (Site) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Comment("站点名称").
			Optional().
			Nillable(),

		field.String("slug").
			Comment("站点标识").
			Optional().
			Nillable(),

		field.String("domain").
			Comment("主域名").
			Optional().
			Nillable(),

		field.Strings("alternate_domains").
			Comment("备用域名列表").
			Optional(),

		field.Bool("is_default").
			Comment("是否为默认站点").
			Default(false).
			Optional().
			Nillable(),

		field.Enum("status").
			Comment("站点状态").
			NamedValues(
				"SiteStatusActive", "SITE_STATUS_ACTIVE",
				"SiteStatusInactive", "SITE_STATUS_INACTIVE",
				"SiteStatusMaintenance", "SITE_STATUS_MAINTENANCE",
			).
			Default("SITE_STATUS_ACTIVE").
			Optional().
			Nillable(),

		field.String("default_locale").
			Comment("默认语言").
			Optional().
			Nillable(),

		field.String("template").
			Comment("站点模板").
			Optional().
			Nillable(),

		field.String("theme").
			Comment("主题名称").
			Optional().
			Nillable(),

		field.Uint64("visit_count").
			Comment("访问次数").
			Optional().
			Nillable(),
	}
}

// Mixin of the Site.
func (Site) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.AutoIncrementId{},
		mixin.TimeAt{},
		mixin.OperatorID{},
		mixin.TenantID[uint32]{},
	}
}

func (Site) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("tenant_id", "slug").
			Unique(),
		index.Fields("tenant_id", "domain").
			Unique(),
		index.Fields("tenant_id", "is_default"),
		index.Fields("tenant_id", "status"),
	}
}
