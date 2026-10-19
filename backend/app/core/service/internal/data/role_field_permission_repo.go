package data

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	entCrud "github.com/tx7do/go-crud/entgo"

	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/rolefieldpermission"

	identityV1 "go-wind-cms/api/gen/go/identity/service/v1"
	permissionV1 "go-wind-cms/api/gen/go/permission/service/v1"
)

// maxFieldPermissionsPerRole 单角色字段权限条目上限：防止异常配置撑爆配置 UI。
const maxFieldPermissionsPerRole = 256

// RoleFieldPermissionRepo 角色字段权限配置仓储。
// 黑名单语义：记录角色在各资源上不可见的字段集；未配置即全字段可见。
type RoleFieldPermissionRepo struct {
	log       *log.Helper
	entClient *entCrud.EntClient[*ent.Client]
}

func NewRoleFieldPermissionRepo(ctx *bootstrap.Context, entClient *entCrud.EntClient[*ent.Client]) *RoleFieldPermissionRepo {
	return &RoleFieldPermissionRepo{
		log:       ctx.NewLoggerHelper("role-field-permission/repo/core-service"),
		entClient: entClient,
	}
}

// validateEntries 校验字段权限条目：资源/字段名非空、字段去重、总量不超上限。
// 资源名不在这里做注册表校验——裁剪器对未知资源天然不生效，垃圾条目无安全影响。
func (r *RoleFieldPermissionRepo) validateEntries(entries []*permissionV1.RoleFieldPermission) error {
	total := 0
	for _, entry := range entries {
		if entry.GetResource() == "" {
			return errors.BadRequest("ROLE_FIELD_PERMISSION", "field permission resource is required")
		}
		seen := make(map[string]struct{}, len(entry.GetHiddenFields()))
		for _, f := range entry.GetHiddenFields() {
			if f == "" {
				return errors.BadRequest("ROLE_FIELD_PERMISSION", "field permission field name is required")
			}
			seen[f] = struct{}{}
		}
		total += len(seen)
	}
	if total > maxFieldPermissionsPerRole {
		r.log.Errorf("field permission entries [%d] exceed cap [%d]", total, maxFieldPermissionsPerRole)
		return errors.BadRequest("ROLE_FIELD_PERMISSION", "too many field permission entries")
	}
	return nil
}

// ReplaceForRole 整体替换角色的字段权限配置：先清空既有配置，再插入新集合。
func (r *RoleFieldPermissionRepo) ReplaceForRole(
	ctx context.Context,
	roleID uint32,
	entries []*permissionV1.RoleFieldPermission,
) error {
	if err := r.validateEntries(entries); err != nil {
		return err
	}

	if _, err := r.entClient.Client().RoleFieldPermission.Delete().
		Where(rolefieldpermission.RoleIDEQ(roleID)).
		Exec(ctx); err != nil {
		r.log.Errorf("delete role [%d] field permissions failed: %s", roleID, err.Error())
		return errors.InternalServer("ROLE_FIELD_PERMISSION", "delete role field permissions failed")
	}

	if len(entries) == 0 {
		return nil
	}

	now := time.Now()
	for _, entry := range entries {
		seen := make(map[string]struct{}, len(entry.GetHiddenFields()))
		for _, f := range entry.GetHiddenFields() {
			if f == "" {
				continue
			}
			if _, ok := seen[f]; ok {
				continue
			}
			seen[f] = struct{}{}

			builder := r.entClient.Client().RoleFieldPermission.Create().
				SetRoleID(roleID).
				SetResource(entry.GetResource()).
				SetFieldName(f).
				SetCreatedAt(now)
			if err := builder.
				OnConflictColumns(
					rolefieldpermission.FieldRoleID,
					rolefieldpermission.FieldResource,
					rolefieldpermission.FieldFieldName,
				).
				UpdateNewValues().
				Exec(ctx); err != nil {
				r.log.Errorf("insert role field permission failed: %s", err.Error())
				return errors.InternalServer("ROLE_FIELD_PERMISSION", "insert role field permission failed")
			}
		}
	}

	return nil
}

// ListHiddenFieldGroupsByRoleIDs 按角色ID列表取字段权限配置（按角色分组）。
// 登录聚合用：一次查询取全部角色的配置，避免 N+1。
func (r *RoleFieldPermissionRepo) ListHiddenFieldGroupsByRoleIDs(ctx context.Context, roleIDs []uint32) (map[uint32][]*permissionV1.RoleFieldPermission, error) {
	if len(roleIDs) == 0 {
		return map[uint32][]*permissionV1.RoleFieldPermission{}, nil
	}

	entities, err := r.entClient.Client().RoleFieldPermission.Query().
		Where(
			rolefieldpermission.RoleIDIn(roleIDs...),
		).
		Select(
			rolefieldpermission.FieldRoleID,
			rolefieldpermission.FieldResource,
			rolefieldpermission.FieldFieldName,
		).
		All(ctx)
	if err != nil {
		r.log.Errorf("query field permissions by role ids failed: %s", err.Error())
		return nil, errors.InternalServer("ROLE_FIELD_PERMISSION", "query role field permissions failed")
	}

	type groupKey struct {
		roleID   uint32
		resource string
	}

	index := make(map[groupKey]int)
	result := make(map[uint32][]*permissionV1.RoleFieldPermission, len(roleIDs))
	for _, e := range entities {
		key := groupKey{roleID: e.RoleID, resource: e.Resource}
		idx, ok := index[key]
		if !ok {
			result[e.RoleID] = append(result[e.RoleID], &permissionV1.RoleFieldPermission{
				Resource:     e.Resource,
				HiddenFields: []string{e.FieldName},
			})
			index[key] = len(result[e.RoleID]) - 1
			continue
		}
		result[e.RoleID][idx].HiddenFields = append(result[e.RoleID][idx].HiddenFields, e.FieldName)
	}

	return result, nil
}

// ListByRoleID 列出角色的字段权限配置（角色详情回显用）。
func (r *RoleFieldPermissionRepo) ListByRoleID(ctx context.Context, roleID uint32) ([]*permissionV1.RoleFieldPermission, error) {
	groups, err := r.ListHiddenFieldGroupsByRoleIDs(ctx, []uint32{roleID})
	if err != nil {
		return nil, err
	}
	return groups[roleID], nil
}

// ListHiddenFieldsByUser 聚合用户全部角色的隐藏字段为 "资源.字段" 串数组。
// 非租户上下文的用户（tenant_id=0 平台用户）一律空集：平台管理员字段全可见。
func (r *RoleFieldPermissionRepo) ListHiddenFieldsByUser(ctx context.Context, userRepo UserRepo, userRoleRepo *UserRoleRepo, userID uint32) ([]string, error) {
	// 平台用户（tenant_id=0）不做字段裁剪
	user, err := userRepo.Get(ctx, &identityV1.GetUserRequest{
		QueryBy: &identityV1.GetUserRequest_Id{Id: userID},
	})
	if err != nil {
		return []string{}, nil
	}
	if user.GetTenantId() == 0 {
		return []string{}, nil
	}

	roleIDs, err := userRoleRepo.ListRoleIDs(ctx, userID, false)
	if err != nil {
		r.log.Errorf("hidden fields: list role ids for user [%d] failed: %s", userID, err.Error())
		return nil, err
	}

	groups, err := r.ListHiddenFieldGroupsByRoleIDs(ctx, roleIDs)
	if err != nil {
		return nil, err
	}
	r.log.Infof("hidden fields: user [%d] roles [%v] groups [%d]", userID, roleIDs, len(groups))

	seen := make(map[string]struct{})
	var out []string
	for _, entries := range groups {
		for _, entry := range entries {
			resource := entry.GetResource()
			if resource == "" {
				continue
			}
			for _, f := range entry.GetHiddenFields() {
				if f == "" {
					continue
				}
				key := resource + "." + f
				if _, ok := seen[key]; ok {
					continue
				}
				seen[key] = struct{}{}
				out = append(out, key)
			}
		}
	}
	r.log.Infof("hidden fields: user [%d] aggregated [%d] entries", userID, len(out))
	return out, nil
}
