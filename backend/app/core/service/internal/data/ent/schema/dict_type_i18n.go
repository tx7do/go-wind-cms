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

// DictTypeI18n holds the schema definition for the DictTypeI18n entity.
type DictTypeI18n struct {
	ent.Schema
}

func (DictTypeI18n) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "sys_dict_type_i18n",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
		entsql.WithComments(true),
		schema.Comment("字典类型翻译表"),
	}
}

// Fields of the DictTypeI18n.
func (DictTypeI18n) Fields() []ent.Field {
	return []ent.Field{
		field.String("language_code").
			Comment("语言代码").
			NotEmpty().
			Immutable().
			Optional().
			Nillable(),

		field.String("type_name").
			Comment("字典类型名称").
			NotEmpty().
			Optional().
			Nillable(),
	}
}

// Mixin of the DictTypeI18n.
func (DictTypeI18n) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.AutoIncrementId{},
		mixin.TimeAt{},
		mixin.OperatorID{},
		mixin.Description{},
		mixin.TenantID[uint32]{},
	}
}

// Edges of the DictTypeI18n.
func (DictTypeI18n) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("dict_type", DictType.Type).
			Ref("i18ns").
			Unique(),
	}
}

// Indexes of the DictTypeI18n.
func (DictTypeI18n) Indexes() []ent.Index {
	return []ent.Index{
		//index.Fields("type_id", "language_code").
		//	Unique().
		//	StorageKey("idx_sys_dict_type_i18n_type_id_lang_code"),

		index.Fields("language_code").
			StorageKey("idx_sys_dict_type_i18n_lang_code"),
	}
}
