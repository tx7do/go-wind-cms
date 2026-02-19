package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/tx7do/go-crud/entgo/mixin"
)

// DictType holds the schema definition for the DictType entity.
type DictType struct {
	ent.Schema
}

func (DictType) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "sys_dict_types",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
		entsql.WithComments(true),
		schema.Comment("字典类型表"),
	}
}

// Fields of the DictType.
func (DictType) Fields() []ent.Field {
	return []ent.Field{
		field.String("type_code").
			Comment("字典类型唯一代码").
			NotEmpty().
			Immutable().
			Optional().
			Nillable(),
	}
}

// Mixin of the DictType.
func (DictType) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.AutoIncrementId{},
		mixin.TimeAt{},
		mixin.OperatorID{},
		mixin.IsEnabled{},
		mixin.SortOrder{},
		mixin.TenantID[uint32]{},
	}
}

// Edges of the DictType.
func (DictType) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("entries", DictEntry.Type).
			Required().
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}).
			StorageKey(edge.Column("type_id")),

		edge.To("i18ns", DictTypeI18n.Type).
			Required().
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}).
			StorageKey(edge.Column("type_id")),
	}
}

// Indexes of the DictType.
func (DictType) Indexes() []ent.Index {
	return []ent.Index{
		// 租户级唯一：同一租户下 type_code 唯一
		index.Fields("tenant_id", "type_code").
			Unique().
			StorageKey("uix_sys_dict_types_tenant_type_code"),

		// 支持按租户快速筛选
		index.Fields("tenant_id").
			StorageKey("idx_sys_dict_types_tenant_id"),

		// 按启用状态过滤
		index.Fields("is_enabled").
			StorageKey("idx_sys_dict_types_is_enabled"),

		// 按排序值查询/排序优化
		index.Fields("sort_order").
			StorageKey("idx_sys_dict_types_sort_order"),
	}
}
