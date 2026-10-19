package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/tx7do/go-crud/entgo/mixin"
)

// RoleFieldPermission 角色字段权限：角色在某资源上隐藏的字段集（黑名单语义）。
// 一行 = (角色, 资源, 字段)；未配置即全字段可见。租户继承角色表的租户归属。
type RoleFieldPermission struct {
	ent.Schema
}

func (RoleFieldPermission) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "role_field_permissions",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
		entsql.WithComments(true),
		schema.Comment("角色字段权限表"),
	}
}

// Fields of the RoleFieldPermission.
func (RoleFieldPermission) Fields() []ent.Field {
	return []ent.Field{
		field.Uint32("role_id").
			Comment("角色ID").
			Positive(),

		field.String("resource").
			Comment("资源名（proto 消息名，如 User）").
			NotEmpty(),

		field.String("field_name").
			Comment("隐藏的字段名（proto 字段 json_name，如 email）").
			NotEmpty(),
	}
}

// Mixin of the RoleFieldPermission.
func (RoleFieldPermission) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.AutoIncrementId{},
		mixin.TimeAt{},
		mixin.OperatorID{},
	}
}

// Indexes of the RoleFieldPermission.
func (RoleFieldPermission) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("role_id", "resource", "field_name").
			Unique().
			StorageKey("uidx_role_field_permission"),
		index.Fields("role_id").StorageKey("idx_role_field_perm_role"),
	}
}
