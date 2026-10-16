package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/tx7do/go-crud/entgo/mixin"
)

// CollabDocState holds the schema definition for the CollabDocState entity.
// 协同编辑文档状态（ygo 持久化适配的后端存储）：按房间存最新合并后的
// V1 update（ygo 侧已防抖合并），服务重启后据此恢复协同会话，消除
// "重启/全员断线丢最后几秒输入" 的窗口。内容事实源仍是 core 内容表
// （显式保存链路），本表只是协同会话的续命存储。
type CollabDocState struct {
	ent.Schema
}

func (CollabDocState) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "collab_doc_states",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
		entsql.WithComments(true),
		schema.Comment("协同编辑文档状态表（房间 → 最新 CRDT update）"),
	}
}

// Fields of the CollabDocState.
func (CollabDocState) Fields() []ent.Field {
	return []ent.Field{
		field.String("room").
			Comment("房间名（如 post:12:zh-CN）").
			NotEmpty(),
		field.Bytes("doc_update").
			Comment("最新合并后的 V1 update（ygo 防抖合并产物，整量替换）").
			Optional(),
	}
}

// Mixin of the CollabDocState.
func (CollabDocState) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.AutoIncrementId{},
		mixin.TimeAt{},
	}
}

func (CollabDocState) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("room").Unique(),
	}
}
