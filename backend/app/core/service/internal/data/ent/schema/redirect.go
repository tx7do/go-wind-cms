package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/tx7do/go-crud/entgo/mixin"
)

// Redirect holds the schema definition for the Redirect entity.
type Redirect struct {
	ent.Schema
}

func (Redirect) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "redirects",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
		entsql.WithComments(true),
		schema.Comment("URL 重定向表（SEO 迁移/短链管理）"),
	}
}

// Fields of the Redirect.
func (Redirect) Fields() []ent.Field {
	return []ent.Field{
		field.String("from_path").
			Comment("源路径（站内相对路径，如 /old-post）").
			NotEmpty(),
		field.String("to_url").
			Comment("目标地址（站内路径或完整 URL）").
			NotEmpty(),
		field.Uint32("status_code").
			Comment("HTTP 状态码（301 永久 / 302 临时 / 307 / 308）").
			Default(301).
			Optional().
			Nillable(),
		field.Bool("enabled").
			Comment("是否启用").
			Default(true).
			Optional().
			Nillable(),
		field.String("description").
			Comment("用途说明").
			Optional().
			Nillable(),
	}
}

// Mixin of the Redirect.
func (Redirect) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.AutoIncrementId{},
		mixin.TimeAt{},
		mixin.OperatorID{},
		mixin.TenantID[uint32]{},
	}
}

func (Redirect) Indexes() []ent.Index {
	return []ent.Index{
		// 源路径查重/查找；同租户内 from_path 唯一（服务层保证，索引加速查找）
		index.Fields("from_path"),
		index.Fields("enabled"),
	}
}
