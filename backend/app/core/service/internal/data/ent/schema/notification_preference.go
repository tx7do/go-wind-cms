package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"

	"github.com/tx7do/go-crud/entgo/mixin"
)

// NotificationPreference 用户通知偏好，每用户一行：实时推送静音时段 + 消息分类退订。
//
// 刻意不挂 TenantID mixin：偏好跟着人走而不是跟着租户走，读方既有用户自己的上下文
//（个人中心），也有跨租户上下文（系统通知在 SystemContext 下逐收件人投递）。
// 挂了租户 mixin 后投递链读偏好会被租户过滤按发送方租户拦截，其他租户收件人的
// 偏好永远查不到；去租户化后唯一访问锚是 user_id（服务端从操作人/收件人钉定，不接受
// 客户端传入），越权面不因此扩大。全链 fail-open：读不到偏好 = 不静音不退订。
type NotificationPreference struct {
	ent.Schema
}

func (NotificationPreference) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "notification_preferences",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
		entsql.WithComments(true),
		schema.Comment("用户通知偏好表（静音时段 + 分类退订）"),
	}
}

// Fields of the NotificationPreference.
func (NotificationPreference) Fields() []ent.Field {
	return []ent.Field{
		field.Uint32("user_id").
			Comment("用户ID（唯一，每用户一行）").
			Unique(),

		field.Bool("quiet_enabled").
			Comment("是否启用实时推送静音时段（只抑制 SSE 实时推送，收件行照常落库）").
			Default(false),

		field.Int32("quiet_start_minute").
			Comment("静音开始：自当日 00:00 起的分钟数（0-1439）").
			Default(1320), // 22:00

		field.Int32("quiet_end_minute").
			Comment("静音结束：自当日 00:00 起的分钟数（0-1439）；跨零点窗口 start > end 合法").
			Default(480), // 08:00

		field.JSON("muted_category_ids", []uint32{}).
			Comment("退订的站内信分类ID列表（预留：v1 无全员广播，无执行点，仅落库）").
			Optional(),
	}
}

// Mixin of the NotificationPreference.
func (NotificationPreference) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.AutoIncrementId{},
		mixin.TimeAt{},
		mixin.OperatorID{},
	}
}
