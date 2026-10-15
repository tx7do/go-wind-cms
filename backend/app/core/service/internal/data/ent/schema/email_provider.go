package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/tx7do/go-crud/entgo/mixin"
)

// EmailProvider holds the schema definition for the EmailProvider entity.
type EmailProvider struct {
	ent.Schema
}

func (EmailProvider) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "email_providers",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
		entsql.WithComments(true),
		schema.Comment("邮件出站 SMTP 配置表（凭据 AES-GCM 加密落库；未配置时邮件链回退环境变量）"),
	}
}

// Fields of the EmailProvider.
func (EmailProvider) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Comment("配置名称（如 主 SMTP / 备用）").
			NotEmpty(),
		field.String("host").
			Comment("SMTP 主机").
			NotEmpty(),
		field.Uint32("port").
			Comment("端口（587/465/25）").
			Default(587).
			Optional().
			Nillable(),
		field.String("username").
			Comment("认证用户名").
			Optional().
			Nillable(),
		field.String("password").
			Comment("认证密码（AES-GCM 加密落库，接口永不回显）").
			Optional().
			Nillable(),
		field.String("from").
			Comment("发件人地址（From 头）").
			NotEmpty(),
		field.String("tls_mode").
			Comment("TLS 模式（starttls/ssl/none）").
			Default("starttls").
			Optional().
			Nillable(),
		field.Bool("enabled").
			Comment("是否启用（邮件链取第一个启用项）").
			Default(true).
			Optional().
			Nillable(),
		field.String("description").
			Comment("用途说明").
			Optional().
			Nillable(),
	}
}

// Mixin of the EmailProvider.
func (EmailProvider) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.AutoIncrementId{},
		mixin.TimeAt{},
		mixin.OperatorID{},
		mixin.TenantID[uint32]{},
	}
}

func (EmailProvider) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("enabled"),
	}
}
