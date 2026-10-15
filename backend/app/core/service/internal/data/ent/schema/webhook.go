package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/tx7do/go-crud/entgo/mixin"
)

// Webhook holds the schema definition for the Webhook entity.
type Webhook struct {
	ent.Schema
}

func (Webhook) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "webhooks",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
		entsql.WithComments(true),
		schema.Comment("Webhook 订阅表（内容事件外发，HMAC 签名投递）"),
	}
}

// Fields of the Webhook.
func (Webhook) Fields() []ent.Field {
	return []ent.Field{
		field.String("url").
			Comment("回调 URL（https:// 推荐，http 仅限内网调试）").
			NotEmpty(),
		field.String("secret").
			Comment("签名密钥（X-Webhook-Signature: sha256=HMAC(body)）").
			Optional().
			Nillable(),
		field.JSON("events", &[]string{}).
			Comment("订阅的事件类型（post.created/post.updated/post.published/post.deleted/page.*）；空数组=全部事件").
			Optional(),
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

// Mixin of the Webhook.
func (Webhook) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.AutoIncrementId{},
		mixin.TimeAt{},
		mixin.OperatorID{},
		mixin.TenantID[uint32]{},
	}
}

func (Webhook) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("enabled"),
	}
}
