package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/tx7do/go-crud/entgo/mixin"
)

// PostLike holds the schema definition for the PostLike entity.
//
// 点赞 ledger 表，记录 (tenant, user, post) 三元组的点赞关系。计数本身不存于此，
// 仅作为 InteractionService 递增 post.likes 缓存的唯一写入依据。复合 unique 索引
// 天然防重复点赞，并为按用户/按帖子查询提供索引。
type PostLike struct {
	ent.Schema
}

func (PostLike) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "post_likes",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
		entsql.WithComments(true),
		schema.Comment("帖子点赞 ledger 表"),
	}
}

func (PostLike) Fields() []ent.Field {
	return []ent.Field{
		field.Uint32("user_id").
			Comment("点赞用户ID").
			Optional().
			Nillable(),
		field.Uint32("post_id").
			Comment("被点赞帖子ID").
			Optional().
			Nillable(),
	}
}

func (PostLike) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.AutoIncrementId{},
		mixin.TimeAt{},
		mixin.TenantID[uint32]{},
	}
}

func (PostLike) Indexes() []ent.Index {
	return []ent.Index{
		// 复合 unique 索引：天然防重复点赞，兼作 (user, post) 查询索引
		index.Fields("tenant_id", "user_id", "post_id").
			Unique(),
	}
}
