package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/tx7do/go-crud/entgo/mixin"
)

// SearchSynonym 搜索同义词（租户级词库）：检索时查询侧按触发词扩展（OR 组），
// 免索引重建、更新即时生效。
type SearchSynonym struct {
	ent.Schema
}

func (SearchSynonym) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "sys_search_synonyms"},
	}
}

// Fields of the SearchSynonym.
func (SearchSynonym) Fields() []ent.Field {
	return []ent.Field{
		field.String("term").
			Comment("触发词（搜索词命中即扩展）").
			NotEmpty().
			MaxLen(64),
		field.String("synonyms").
			Comment("扩展词 JSON 字符串数组（与触发词一起以 OR 参与检索）").
			NotEmpty().
			Default("[]"),
		field.Bool("enabled").
			Comment("是否启用").
			Default(true).
			Optional().
			Nillable(),
	}
}

// Edges of the SearchSynonym.
func (SearchSynonym) Edges() []ent.Edge {
	return nil
}

// Mixin of the SearchSynonym.
func (SearchSynonym) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.TimeAt{},
		mixin.AutoIncrementId{},
		mixin.TenantID[uint32]{},
	}
}

// Indexes of the SearchSynonym.
func (SearchSynonym) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("term").Unique(),
	}
}
