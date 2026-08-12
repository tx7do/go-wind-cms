package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/tx7do/go-crud/entgo/mixin"
)

// PostWatch holds the schema definition for the PostWatch entity.
//
// 收藏 ledger 表，记录 (tenant, user, post) 三元组的收藏关系。收藏不对应任何计数
// 缓存列，仅维护收藏列表本身。复合 unique 索引天然防重复收藏。
type PostWatch struct {
	ent.Schema
}

func (PostWatch) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "post_watches",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
		entsql.WithComments(true),
		schema.Comment("帖子收藏 ledger 表"),
	}
}

func (PostWatch) Fields() []ent.Field {
	return []ent.Field{
		field.Uint32("user_id").
			Comment("收藏用户ID").
			Optional().
			Nillable(),
		field.Uint32("post_id").
			Comment("被收藏帖子ID").
			Optional().
			Nillable(),
	}
}

func (PostWatch) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.AutoIncrementId{},
		mixin.TimeAt{},
		mixin.TenantID[uint32]{},
	}
}

func (PostWatch) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("tenant_id", "user_id", "post_id").
			Unique(),
	}
}
