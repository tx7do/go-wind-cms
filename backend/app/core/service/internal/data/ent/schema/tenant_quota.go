package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/tx7do/go-crud/entgo/mixin"
)

// TenantQuota holds the schema definition for the TenantQuota entity.
// 租户配额：配额按租户直配（CMS 无套餐表，与上游 PlanQuota 挂 plan_id 不同）。
// 平台级配置表，不挂 TenantID mixin（配额管理是平台管理员动作，闸门读取走系统旁路）。
type TenantQuota struct {
	ent.Schema
}

func (TenantQuota) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "sys_tenant_quotas",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
		entsql.WithComments(true),
		schema.Comment("租户配额表（按租户直配，配额类型唯一）"),
	}
}

// Fields of the TenantQuota.
func (TenantQuota) Fields() []ent.Field {
	return []ent.Field{
		field.Uint32("tenant_id").Comment("租户ID").Positive(),
		field.Enum("quota_type").
			Comment("配额类型").
			NamedValues(
				"QuotaTypeUnspecified", "QUOTA_TYPE_UNSPECIFIED",
				"UserLimit", "USER_LIMIT",
				"Storage", "STORAGE",
				"ApiCall", "API_CALL",
				"AiTokens", "AI_TOKENS",
			),
		field.Uint64("quota_value").Comment("配额值（USER_LIMIT=人数、STORAGE=字节、API_CALL=次数、AI_TOKENS=token 数）"),
	}
}

// Mixin of the TenantQuota.
func (TenantQuota) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.AutoIncrementId{},
		mixin.TimeAt{},
		mixin.OperatorID{},
	}
}

func (TenantQuota) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("tenant_id", "quota_type").Unique().StorageKey("idx_sys_tenant_quotas_tenant_type"),
	}
}
