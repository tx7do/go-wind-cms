package mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

var _ ent.Mixin = (*EditorType)(nil)

type EditorType struct{ mixin.Schema }

func (EditorType) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("editor_type").
			Comment("编辑器类型").
			NamedValues(
				"EditorTypeMarkdown", "EDITOR_TYPE_MARKDOWN",
				"EditorTypeRichText", "EDITOR_TYPE_RICH_TEXT",
				"EditorTypePlainText", "EDITOR_TYPE_PLAIN_TEXT",
				"EditorTypeCode", "EDITOR_TYPE_CODE",
				"EditorTypeJsonBlock", "EDITOR_TYPE_JSON_BLOCK",
				"EditorTypeVisualBuilder", "EDITOR_TYPE_VISUAL_BUILDER",
			).
			Default("EDITOR_TYPE_MARKDOWN").
			Optional().
			Nillable(),
	}
}
