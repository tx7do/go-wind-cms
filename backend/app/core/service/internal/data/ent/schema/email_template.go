package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/tx7do/go-crud/entgo/mixin"
)

// EmailTemplate holds the schema definition for the EmailTemplate entity.
// 邮件模板：系统通知出站的可视化模板。code 唯一定位（如 notify/system），
// 主题与正文支持 {{.Title}}/{{.Content}}/{{.SiteURL}} 占位符（text/template）。
// 渲染顺序：租户启用模板 → 平台启用模板 → 内置最小模板（buildNotifyEmailHTML）。
type EmailTemplate struct {
	ent.Schema
}

func (EmailTemplate) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "email_templates",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
		entsql.WithComments(true),
		schema.Comment("邮件模板表（系统通知出站可视化模板）"),
	}
}

// Fields of the EmailTemplate.
func (EmailTemplate) Fields() []ent.Field {
	return []ent.Field{
		field.String("code").
			Comment("模板代码（唯一键，如 notify/system）").
			NotEmpty(),
		field.String("name").
			Comment("模板名称（管理端展示）").
			NotEmpty(),
		field.String("subject").
			Comment("邮件主题模板（支持 {{.Title}} 占位符）").
			NotEmpty(),
		field.Text("body_html").
			Comment("邮件正文 HTML 模板（支持 {{.Title}}/{{.Content}}/{{.SiteURL}}）").
			Optional(),
		field.Bool("enabled").
			Comment("是否启用").
			Default(true).
			Optional().
			Nillable(),
		field.String("description").
			Comment("用途说明").
			Optional().
			Nillable(),
	}
}

// Mixin of the EmailTemplate.
func (EmailTemplate) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.AutoIncrementId{},
		mixin.TimeAt{},
		mixin.OperatorID{},
		mixin.TenantID[uint32]{},
	}
}

func (EmailTemplate) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("code"),
		index.Fields("enabled"),
	}
}
