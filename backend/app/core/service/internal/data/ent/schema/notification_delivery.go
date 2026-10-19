package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/tx7do/go-crud/entgo/mixin"
)

// NotificationDelivery holds the schema definition for the NotificationDelivery entity.
// 通知投递台账：一条通知 × 一个渠道 × 一个收件人 = 一行。
//
// 与 notification_channels 同为平台级全局表（无租户隔离）：渠道花名册本就是平台配的，
// 台账记录"平台用了它发没发出去"。收件地址跨租户，读侧只开放给平台视图
//（service 层 IsPlatformContext 守卫）。
type NotificationDelivery struct {
	ent.Schema
}

func (NotificationDelivery) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "notification_deliveries",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
		entsql.WithComments(true),
		schema.Comment("通知投递台账表"),
	}
}

// Fields of the NotificationDelivery.
func (NotificationDelivery) Fields() []ent.Field {
	return []ent.Field{
		// 枚举列一律 Optional().Nillable()：copier 注册的枚举转换是"指针↔指针对"，
		// 值型枚举列（无 Nillable）配可选指针 DTO 字段会失配并让 DTO 落成零值枚举
		// ——上游在 channel.type 上踩过这个坑（queryTypeByIDs 只能批量回填绕开）。
		field.Enum("event_type").
			Comment("业务事件类型（成员名与 proto 的 EventType 枚举逐字同名，EnumTypeConverter 按名字配对）").
			NamedValues(
				"InternalMessage", "INTERNAL_MESSAGE",
				"ChannelTestEmail", "CHANNEL_TEST_EMAIL",
				"SystemNotify", "SYSTEM_NOTIFY",
				"ChannelTestExt", "CHANNEL_TEST_EXT",
				"MonitorAlert", "MONITOR_ALERT",
				"PasswordResetCode", "PASSWORD_RESET_CODE",
				"ChannelTestSms", "CHANNEL_TEST_SMS",
			).
			Optional().
			Nillable(),

		field.Enum("channel").
			Comment("投递渠道（EXT：扩展渠道大类，具体型别在渠道行的 type 列）").
			NamedValues(
				"Email", "EMAIL",
				"Webhook", "WEBHOOK",
				"Internal", "INTERNAL",
				"Sms", "SMS",
				"Ext", "EXT",
			).
			Optional().
			Nillable(),

		field.Uint32("channel_id").
			Comment("实际选中的渠道配置ID").
			Optional().
			Nillable(),

		field.Uint32("recipient_user_id").
			Comment("收件用户ID（直发模式可为空）").
			Optional().
			Nillable(),

		// 产生本条投递的业务对象主键，含义由 event_type 决定（INTERNAL_MESSAGE → internal_messages.id）。
		// 台账不存正文快照（防敏感正文明文集邮），
		// 这一列是"这条投递发的是什么"的唯一回跳入口。
		field.Uint32("related_id").
			Comment("关联业务对象ID（按事件类型解释）").
			Optional().
			Nillable(),

		field.String("target").
			Comment("投递目标（脱敏后的地址）").
			Optional().
			Nillable(),

		field.Enum("status").
			Comment("投递状态").
			NamedValues(
				"Sending", "SENDING",
				"Sent", "SENT",
				"Failed", "FAILED",
				"Skipped", "SKIPPED",
			).
			Default("SENDING").
			Optional().
			Nillable(),

		field.String("last_error").
			Comment("最近一次失败原因").
			Optional().
			Nillable(),

		// 一次派发意图的标识：调用方不传时服务端生成 UUID。唯一性落在 (request_id, channel)
		// 而不是单列 —— 台账的行粒度是"一条通知 × 一个渠道 × 一个收件人"，
		// 将来一次意图若扇出到多渠道，单列唯一会直接拒绝第二行。
		field.String("request_id").
			Comment("派发请求ID（幂等锚）").
			Optional().
			Nillable().
			MaxLen(64),

		// 已尝试次数，含首次。同步路径恒为 1；异步路径每次尝试先 +1 再拨号，
		// 所以"attempts>0 但 status 还是 SENDING"就是"拨过号但没回写结论"的唯一线索。
		field.Uint32("attempts").
			Comment("实际尝试投递次数（含首次）").
			Default(0).
			Optional().
			Nillable(),

		field.String("receipt_provider_id").
			Comment("服务商侧消息/请求 ID（阿里云 RequestId、IM 响应 ID 等，拿得到就带）；跨系统对账用").
			Optional().
			Nillable(),
		field.Time("sent_at").
			Comment("投递完成时间").
			Optional().
			Nillable(),
	}
}

// Mixin of the NotificationDelivery.
func (NotificationDelivery) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.AutoIncrementId{},
		mixin.TimeAt{},
		mixin.OperatorID{},
	}
}

// Indexes of the NotificationDelivery.
func (NotificationDelivery) Indexes() []ent.Index {
	return []ent.Index{
		// 台账页默认按时间倒序翻页，并按事件类型/状态过滤
		index.Fields("created_at").
			StorageKey("idx_notification_delivery_created_at"),
		index.Fields("event_type", "created_at").
			StorageKey("idx_notification_delivery_event_created_at"),
		index.Fields("status", "created_at").
			StorageKey("idx_notification_delivery_status_created_at"),
		index.Fields("recipient_user_id").
			StorageKey("idx_notification_delivery_recipient_user"),

		// "这条消息到底发出去了没"——按业务对象反查投递事实
		index.Fields("event_type", "related_id").
			StorageKey("idx_notification_delivery_event_related"),

		// 幂等锚：唯一而不是普通索引才是它的价值所在——
		// 重复入队/重复调用在写入这一步就被 DB 挡下，而不是靠读侧自觉。
		index.Fields("request_id", "channel").Unique().
			StorageKey("uidx_notification_delivery_request_channel"),
	}
}
