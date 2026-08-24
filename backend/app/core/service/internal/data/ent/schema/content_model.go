package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/tx7do/go-crud/entgo/mixin"
)

// ContentModel 内容模型：一组字段定义的命名集合，可绑定到 Category，
// 其下的 Post/Page 继承模型字段（值存于内容实体的 custom_fields 列）。
type ContentModel struct {
	ent.Schema
}

func (ContentModel) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "content_models",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
		entsql.WithComments(true),
		schema.Comment("内容模型表"),
	}
}

func (ContentModel) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Comment("模型名称（语言无关，后端标识）").
			NotEmpty().
			Optional().
			Nillable(),

		field.String("code").
			Comment("模型编码，租户内唯一").
			NotEmpty().
			Optional().
			Nillable(),

		field.String("description").
			Comment("模型描述（语言无关）").
			Optional().
			Nillable(),
	}
}

func (ContentModel) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.AutoIncrementId{},
		mixin.TimeAt{},
		mixin.OperatorID{},
		mixin.SortOrder{},
		mixin.TenantID[uint32]{},
	}
}

func (ContentModel) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("tenant_id", "code").Unique().StorageKey("uix_content_model_tenant_code"),
	}
}
