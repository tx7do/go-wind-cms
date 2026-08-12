package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/tx7do/go-crud/entgo/mixin"
)

// CommentLike holds the schema definition for the CommentLike entity.
//
// 评论点赞 ledger 表，记录 (tenant, user, comment) 三元组的点赞关系。计数本身不存于此，
// 仅作为 InteractionService 递增 comment.like_count 缓存的唯一写入依据。复合 unique
// 索引天然防重复点赞，并兼作 (user, comment) 查询索引。
type CommentLike struct {
	ent.Schema
}

func (CommentLike) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "comment_likes",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
		entsql.WithComments(true),
		schema.Comment("评论点赞 ledger 表"),
	}
}

func (CommentLike) Fields() []ent.Field {
	return []ent.Field{
		field.Uint32("user_id").
			Comment("点赞用户ID").
			Optional().
			Nillable(),
		field.Uint32("comment_id").
			Comment("被点赞评论ID").
			Optional().
			Nillable(),
	}
}

func (CommentLike) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.AutoIncrementId{},
		mixin.TimeAt{},
		mixin.TenantID[uint32]{},
	}
}

func (CommentLike) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("tenant_id", "user_id", "comment_id").
			Unique(),
	}
}
