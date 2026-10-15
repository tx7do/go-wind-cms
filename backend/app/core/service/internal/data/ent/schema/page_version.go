package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/tx7do/go-crud/entgo/mixin"
)

// PageVersion holds the schema definition for the PageVersion entity.
type PageVersion struct {
	ent.Schema
}

func (PageVersion) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "page_versions",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
		entsql.WithComments(true),
		schema.Comment("页面版本表（内容快照，保存/恢复时自动落版本）"),
	}
}

// Fields of the PageVersion.
func (PageVersion) Fields() []ent.Field {
	return []ent.Field{
		field.Uint32("page_id").
			Comment("关联的页面ID").
			Optional().
			Nillable(),

		field.Uint32("version_number").
			Comment("版本号（每页递增，1 起；创建为 1，之后每次保存 +1）").
			Default(1).
			Optional().
			Nillable(),

		field.Text("content").
			Comment("快照 JSON（Page DTO 全量，含 translations/sections）").
			Optional().
			Nillable(),

		field.String("comment").
			Comment("版本说明（自动保存为空）").
			Optional().
			Nillable(),
	}
}

// Mixin of the PageVersion.
func (PageVersion) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.AutoIncrementId{},
		mixin.TimeAt{},
		mixin.OperatorID{},
		mixin.TenantID[uint32]{},
	}
}

func (PageVersion) Indexes() []ent.Index {
	return []ent.Index{
		// 按页面查版本列表；version_number 非唯一索引：并发保存竞态下宁可重号也不阻断保存
		index.Fields("page_id"),
		index.Fields("page_id", "version_number"),
	}
}
