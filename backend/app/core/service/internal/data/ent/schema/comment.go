package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/tx7do/go-crud/entgo/mixin"
)

// Comment holds the schema definition for the Comment entity.
type Comment struct {
	ent.Schema
}

func (Comment) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "comments",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
		entsql.WithComments(true),
		schema.Comment("评论表"),
	}
}

// Fields of the Comment.
func (Comment) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("content_type").
			Comment("内容类型").
			NamedValues(
				"ContentTypePost", "CONTENT_TYPE_POST",
				"ContentTypePage", "CONTENT_TYPE_PAGE",
				"ContentTypeProduct", "CONTENT_TYPE_PRODUCT",
			).
			Optional().
			Nillable(),

		field.Uint32("object_id").
			Comment("对象ID").
			Optional().
			Nillable(),

		field.String("content").
			Comment("评论内容").
			Optional().
			Nillable(),

		field.Uint32("author_id").
			Comment("评论作者ID，0表示游客").
			Default(0).
			Optional().
			Nillable(),

		field.String("author_name").
			Comment("评论作者名称").
			Optional().
			Nillable(),

		field.String("author_email").
			Comment("评论作者邮箱").
			Optional().
			Nillable(),

		field.String("author_url").
			Comment("评论作者网址").
			Optional().
			Nillable(),

		field.Enum("author_type").
			Comment("作者类型").
			NamedValues(
				"AuthorTypeGuest", "AUTHOR_TYPE_GUEST",
				"AuthorTypeUser", "AUTHOR_TYPE_USER",
				"AuthorTypeAdmin", "AUTHOR_TYPE_ADMIN",
				"AuthorTypeModerator", "AUTHOR_TYPE_MODERATOR",
			).
			Optional().
			Nillable(),

		field.Enum("status").
			Comment("评论状态").
			NamedValues(
				"StatusPending", "STATUS_PENDING",
				"StatusApproved", "STATUS_APPROVED",
				"StatusRejected", "STATUS_REJECTED",
				"StatusSpam", "STATUS_SPAM",
			).
			Optional().
			Nillable(),

		field.Uint32("like_count").
			Comment("点赞数").
			Default(0).
			Optional().
			Nillable(),

		field.Uint32("dislike_count").
			Comment("点踩数").
			Default(0).
			Optional().
			Nillable(),

		field.Uint32("reply_count").
			Comment("回复数").
			Default(0).
			Optional().
			Nillable(),

		field.String("ip_address").
			Comment("评论者 IP").
			Optional().
			Nillable(),

		field.String("location").
			Comment("评论者地理位置").
			Optional().
			Nillable(),

		field.String("user_agent").
			Comment("User-Agent").
			Optional().
			Nillable(),

		field.String("detected_language").
			Comment("自动检测的语言代码").
			Optional().
			Nillable(),

		field.Bool("is_spam").
			Comment("是否标记为垃圾评论").
			Default(false).
			Optional().
			Nillable(),

		field.Bool("is_sticky").
			Comment("是否置顶评论").
			Default(false).
			Optional().
			Nillable(),

		field.Uint32("reply_to_id").
			Comment("回复的评论ID").
			Optional().
			Nillable(),
	}
}

// Mixin of the Comment.
func (Comment) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.AutoIncrementId{},
		mixin.TimeAt{},
		mixin.OperatorID{},
		mixin.Tree[Comment]{},
	}
}

func (Comment) Indexes() []ent.Index {
	return []ent.Index{
		// 复合索引，优化按内容类型和对象ID查询评论
		index.Fields("content_type", "object_id"),
		// 单字段索引，用于按评论状态查询
		index.Fields("status"),
		// 单字段索引，用于按作者ID查询其所有评论
		index.Fields("author_id"),
		// 单字段索引，用于树形结构中按父节点查询子评论
		index.Fields("reply_to_id"),
		// 单字段索引，用于垃圾评论过滤
		index.Fields("is_spam"),
		// 单字段索引，用于置顶评论查询
		index.Fields("is_sticky"),
	}
}
