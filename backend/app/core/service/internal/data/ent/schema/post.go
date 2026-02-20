package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/tx7do/go-crud/entgo/mixin"
)

// Post holds the schema definition for the Post entity.
type Post struct {
	ent.Schema
}

func (Post) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "posts",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
		entsql.WithComments(true),
		schema.Comment("帖子表"),
	}
}

// Fields of the Post.
func (Post) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("status").
			Comment("帖子状态").
			NamedValues(
				"PostStatusDraft", "POST_STATUS_DRAFT",
				"PostStatusPublished", "POST_STATUS_PUBLISHED",
				"PostStatusScheduled", "POST_STATUS_SCHEDULED",
				"PostStatusTrashed", "POST_STATUS_TRASHED",
			).
			Default("POST_STATUS_DRAFT").
			Optional().
			Nillable(),

		field.Enum("editor_type").
			Comment("编辑器类型").
			NamedValues(
				"EditorTypeMarkdown", "EDITOR_TYPE_MARKDOWN",
				"EditorTypeRichText", "EDITOR_TYPE_RICH_TEXT",
				"EditorTypeHtml", "EDITOR_TYPE_HTML",
				"EditorTypeJsonBlock", "EDITOR_TYPE_JSON_BLOCK",
				"EditorTypePlainText", "EDITOR_TYPE_PLAIN_TEXT",
				"EditorTypeCode", "EDITOR_TYPE_CODE",
				"EditorTypeWysiwyg", "EDITOR_TYPE_WYSIWYG",
				"EditorTypeVisualBuilder", "EDITOR_TYPE_VISUAL_BUILDER",
				"EditorTypeSlate", "EDITOR_TYPE_SLATE",
				"EditorTypeProsemirror", "EDITOR_TYPE_PROSEMIRROR",
			).
			Default("EDITOR_TYPE_MARKDOWN").
			Optional().
			Nillable(),

		field.String("slug").
			Comment("链接别名").
			Optional().
			Nillable(),

		field.Bool("disallow_comment").
			Comment("不允许评论").
			Default(false).
			Optional().
			Nillable(),

		field.Bool("in_progress").
			Comment("审核中").
			Default(false).
			Optional().
			Nillable(),

		field.Bool("auto_summary").
			Comment("是否自动生成摘要").
			Default(true).
			Optional().
			Nillable(),

		field.Bool("is_featured").
			Comment("是否推荐").
			Default(false).
			Optional().
			Nillable(),

		field.Int32("visits").
			Comment("帖子访问次数").
			Default(0).
			Optional().
			Nillable(),

		field.Int32("likes").
			Comment("帖子点赞次数").
			Default(0).
			Optional().
			Nillable(),

		field.Int32("comment_count").
			Comment("帖子评论数").
			Default(0).
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

		field.String("password_hash").
			Comment("密码哈希").
			Optional().
			Nillable(),

		field.JSON("custom_fields", &map[string]string{}).
			Comment("自定义字段").
			Optional(),

		field.JSON("category_ids", &[]uint32{}).
			Comment("关联的分类ID列表").
			Optional(),

		field.JSON("tag_ids", &[]uint32{}).
			Comment("关联的标签ID列表").
			Optional(),
	}
}

// Mixin of the Post.
func (Post) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.AutoIncrementId{},
		mixin.TimeAt{},
		mixin.OperatorID{},
		mixin.SortOrder{},
	}
}

func (Post) Indexes() []ent.Index {
	return []ent.Index{
		// 单字段索引，用于按帖子状态查询
		index.Fields("status"),
		// 单字段索引，用于按编辑器类型查询
		index.Fields("editor_type"),
		// 单字段索引，用于按链接别名查询
		index.Fields("slug"),
		// 单字段索引，用于按作者ID查询其所有帖子
		index.Fields("author_id"),
		// 单字段索引，用于按是否推荐过滤
		index.Fields("is_featured"),
		// 单字段索引，用于按审核状态过滤
		index.Fields("in_progress"),
		// 单字段索引，用于按评论禁用状态过滤
		index.Fields("disallow_comment"),
		// 复合索引，优化按状态和作者查询
		index.Fields("status", "author_id"),
		// 复合索引，优化按状态和推荐状态查询
		index.Fields("status", "is_featured"),
		// 复合索引，优化按状态和审核状态查询
		index.Fields("status", "in_progress"),
	}
}
