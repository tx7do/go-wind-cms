package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/tx7do/go-crud/entgo/mixin"
)

// QuotaPlan 配额套餐模板：平台级命名模板（名称唯一），items 为配额项 JSON 数组
//（[{quotaType:"USER_LIMIT",limit:10}, ...]，quotaType 复用 TenantQuota 枚举名）。
// 套用到租户 = 按模板逐项 upsert 该租户的 sys_tenant_quotas（模板未含的类型不动）。
type QuotaPlan struct {
	ent.Schema
}

func (QuotaPlan) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "sys_quota_plans"},
	}
}

// Fields of the QuotaPlan.
func (QuotaPlan) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Comment("模板名").
			NotEmpty().
			MaxLen(64),
		field.String("description").
			Comment("描述").
			Optional().
			Nillable().
			MaxLen(255),
		field.String("items").
			Comment("配额项 JSON 数组 [{quotaType,limit}]，quotaType 复用 TenantQuota 枚举名").
			NotEmpty().
			Default("[]"),
	}
}

// Edges of the QuotaPlan.
func (QuotaPlan) Edges() []ent.Edge {
	return nil
}

// Mixins of the QuotaPlan.
func (QuotaPlan) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.TimeAt{},
		mixin.AutoIncrementId{},
		mixin.TenantID[uint32]{},
	}
}

// Indexes of the QuotaPlan.
func (QuotaPlan) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").Unique(),
	}
}
