package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/tx7do/go-crud/entgo/mixin"
)

// FlowRule holds the schema definition for the FlowRule entity.
type FlowRule struct {
	ent.Schema
}

func (FlowRule) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "flow_rules",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
		entsql.WithComments(true),
		schema.Comment("自动化 Flow 规则（内容事件 → 条件 → 动作链）"),
	}
}

// Fields of the FlowRule.
func (FlowRule) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Comment("规则名称").
			NotEmpty().
			MaxLen(128),
		field.String("trigger_event").
			Comment("触发事件（post.published/post.updated/page.published 等内容事件族）").
			NotEmpty().
			MaxLen(64),
		field.String("condition_status").
			Comment("条件：实体状态等值匹配（如 POST_STATUS_PUBLISHED）；空=无条件；与 conditions 同时存在时两者都须命中").
			Optional().
			Nillable().
			MaxLen(64),
		field.String("conditions").
			Comment("结构化条件 JSON 数组 [{field,op,value}]，全部命中(AND)才触发；field 取快照字段名(status/prev_status/title/slug/author_id/category_id 或 cf.<自定义字段键>)，op ∈ eq/ne/contains/in（in 值为 | 分隔）；空=无条件").
			Optional().
			Nillable(),
		field.String("trigger_cron").
			Comment("定时触发 cron 表达式（5 段：分 时 日 月 周；非空=定时规则，忽略事件触发与状态条件）").
			Optional().
			Nillable().
			MaxLen(64),
		field.String("actions").
			Comment("动作链 JSON 数组，元素形如 {'type':'webhook','url':'https://...','secret':'可选签名密钥'}").
			NotEmpty().
			Default("[]"),
		field.String("else_actions").
			Comment("条件不满足分支的动作链 JSON（结构与 actions 相同）；空=不执行").
			Optional().
			Nillable(),
		field.Bool("enabled").
			Comment("是否启用").
			Default(true).
			Optional().
			Nillable(),
		field.Uint32("match_count").
			Comment("累计命中次数（观测用）").
			Default(0).
			Optional().
			Nillable(),
	}
}

// Mixin of the FlowRule.
func (FlowRule) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.AutoIncrementId{},
		mixin.TimeAt{},
		mixin.OperatorID{},
		mixin.TenantID[uint32]{},
	}
}

// Indexes of the FlowRule.
func (FlowRule) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("trigger_event", "enabled"),
	}
}
