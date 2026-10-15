package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/tx7do/go-crud/entgo/mixin"
)

// WebhookDelivery holds the schema definition for the WebhookDelivery entity.
type WebhookDelivery struct {
	ent.Schema
}

func (WebhookDelivery) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "webhook_deliveries",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
		entsql.WithComments(true),
		schema.Comment("Webhook 投递记录表（每次尝试一行，只增不改）"),
	}
}

// Fields of the WebhookDelivery.
func (WebhookDelivery) Fields() []ent.Field {
	return []ent.Field{
		field.Uint32("webhook_id").
			Comment("关联的 Webhook ID").
			Optional().
			Nillable(),
		field.String("event").
			Comment("事件类型（post.published 等）").
			Optional().
			Nillable(),
		field.Uint32("entity_id").
			Comment("事件主体 ID（帖子/页面）").
			Default(0).
			Optional().
			Nillable(),
		field.Enum("status").
			Comment("投递状态").
			NamedValues(
				"DeliveryStatusSuccess", "DELIVERY_STATUS_SUCCESS",
				"DeliveryStatusFailed", "DELIVERY_STATUS_FAILED",
			).
			Default("DELIVERY_STATUS_FAILED").
			Optional().
			Nillable(),
		field.Int("http_status").
			Comment("响应 HTTP 状态码（0=未达 HTTP 层）").
			Default(0).
			Optional().
			Nillable(),
		field.String("error").
			Comment("失败原因（网络错误或非 2xx）").
			Optional().
			Nillable(),
		field.Uint32("duration_ms").
			Comment("耗时（毫秒）").
			Default(0).
			Optional().
			Nillable(),
	}
}

// Mixin of the WebhookDelivery.
func (WebhookDelivery) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.AutoIncrementId{},
		mixin.TimeAt{},
		mixin.TenantID[uint32]{},
	}
}

func (WebhookDelivery) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("webhook_id", "created_at"),
		index.Fields("event"),
	}
}
