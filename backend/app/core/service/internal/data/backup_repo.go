package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	entCrud "github.com/tx7do/go-crud/entgo"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	"go-wind-cms/app/core/service/internal/data/ent"
)

// BackupRepo 负责数据库核心业务表的全量导出与空库恢复，供异步备份任务使用。
//
// 它把对 ent client 的直接依赖收敛在 data 层，避免 service 层（如 TaskService）
// 直接持有 *entCrud.EntClient，从而维持分层边界。
type BackupRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper
}

// NewBackupRepo 创建备份仓库。
func NewBackupRepo(ctx *bootstrap.Context, entClient *entCrud.EntClient[*ent.Client]) *BackupRepo {
	return &BackupRepo{
		entClient: entClient,
		log:       ctx.NewLoggerHelper("backup/repo/core-service"),
	}
}

// IsConfigured 返回备份所需依赖（ent client）是否就绪。
func (r *BackupRepo) IsConfigured() bool {
	return r != nil && r.entClient != nil
}

// coreTableExport 核心表导出集：身份/权限/组织/站点导航配置。
//
// 与上游（8 张 sys_ 表）相比按 CMS 表结构重定义——CMS 把凭证、用户-角色绑定、
// 角色-权限绑定等拆成了独立表，缺任何一张恢复后都无法登录或失去授权关系：
// - sys_user_credentials：登录凭据（密码哈希随备份保真，恢复后管理员能直接登录）
// - sys_user_roles / sys_user_org_units / sys_user_positions：用户绑定三元组
// - sys_role_permissions / sys_role_metadata / sys_role_field_permissions：角色配置
// - sys_permission_* 五张：RBAC 配置（组/API 授权/菜单授权/策略）
// - sys_membership* 四张：租户成员关系
// - sys_apis：API 闸门表；sys_login_policies：登录策略；sys_menus：菜单树
//
// 审计日志/内容业务表等大体量或可再生表不纳入（内容走内容自身的运维通道）。
// 键为备份 JSON 的逻辑名，查询为对应 ent 实体（JSON tag 与列名一致）。
func (r *BackupRepo) coreTableExports(ctx context.Context) []struct {
	name  string
	query func() (any, error)
} {
	client := r.entClient.Client()
	return []struct {
		name  string
		query func() (any, error)
	}{
		{"tenants", func() (any, error) { return client.Tenant.Query().All(ctx) }},
		{"org_units", func() (any, error) { return client.OrgUnit.Query().All(ctx) }},
		{"positions", func() (any, error) { return client.Position.Query().All(ctx) }},
		{"users", func() (any, error) { return client.User.Query().All(ctx) }},
		{"user_credentials", func() (any, error) { return client.UserCredential.Query().All(ctx) }},
		{"roles", func() (any, error) { return client.Role.Query().All(ctx) }},
		{"role_permissions", func() (any, error) { return client.RolePermission.Query().All(ctx) }},
		{"role_metadata", func() (any, error) { return client.RoleMetadata.Query().All(ctx) }},
		{"role_field_permissions", func() (any, error) { return client.RoleFieldPermission.Query().All(ctx) }},
		{"permissions", func() (any, error) { return client.Permission.Query().All(ctx) }},
		{"permission_groups", func() (any, error) { return client.PermissionGroup.Query().All(ctx) }},
		{"permission_apis", func() (any, error) { return client.PermissionApi.Query().All(ctx) }},
		{"permission_menus", func() (any, error) { return client.PermissionMenu.Query().All(ctx) }},
		{"permission_policies", func() (any, error) { return client.PermissionPolicy.Query().All(ctx) }},
		{"user_roles", func() (any, error) { return client.UserRole.Query().All(ctx) }},
		{"user_org_units", func() (any, error) { return client.UserOrgUnit.Query().All(ctx) }},
		{"user_positions", func() (any, error) { return client.UserPosition.Query().All(ctx) }},
		{"memberships", func() (any, error) { return client.Membership.Query().All(ctx) }},
		{"membership_org_units", func() (any, error) { return client.MembershipOrgUnit.Query().All(ctx) }},
		{"membership_positions", func() (any, error) { return client.MembershipPosition.Query().All(ctx) }},
		{"membership_roles", func() (any, error) { return client.MembershipRole.Query().All(ctx) }},
		{"menus", func() (any, error) { return client.Menu.Query().All(ctx) }},
		{"apis", func() (any, error) { return client.Api.Query().All(ctx) }},
		{"login_policies", func() (any, error) { return client.LoginPolicy.Query().All(ctx) }},
	}
}

// ExportCoreTables 导出核心业务表的全量记录。
//
// 任一核心表导出失败即返回 error：备份是"全量恢复"语义，缺表会造成静默还原残缺
// 数据，因此宁可不产出也不产出残缺备份。调用方据此决定是否重试。
//
// 返回的 map 键为逻辑表名，值为可被 json 序列化的 ent 实体切片。
func (r *BackupRepo) ExportCoreTables(ctx context.Context) (map[string]any, error) {
	result := make(map[string]any)

	for _, t := range r.coreTableExports(ctx) {
		rows, err := t.query()
		if err != nil {
			// 核心表导出失败直接返回错误，避免产出残缺备份被当作成功
			r.log.Errorf("backup: export table %q failed: %s", t.name, err.Error())
			return nil, fmt.Errorf("export table %q failed: %w", t.name, err)
		}
		result[t.name] = rows
	}

	return result, nil
}

// restoreTableOrder 恢复顺序：被引用表先于引用表（外键依赖），与导出表集一致。
// 键为备份 JSON 的逻辑名（ExportCoreTables 的键）。
var restoreTableOrder = []string{
	"tenants",
	"org_units",
	"positions",
	"users",
	"user_credentials",
	"roles",
	"role_permissions",
	"role_metadata",
	"role_field_permissions",
	"permissions",
	"permission_groups",
	"permission_apis",
	"permission_menus",
	"permission_policies",
	"user_roles",
	"user_org_units",
	"user_positions",
	"memberships",
	"membership_org_units",
	"membership_positions",
	"membership_roles",
	"menus",
	"apis",
	"login_policies",
}

// restorePhysicalTables 逻辑名 → 物理表名（ent 迁移生成的物理表名）。
var restorePhysicalTables = map[string]string{
	"tenants":          "sys_tenants",
	"org_units":        "sys_org_units",
	"positions":        "sys_positions",
	"users":            "sys_users",
	"user_credentials": "sys_user_credentials",
	"roles":            "sys_roles",
	"role_permissions": "sys_role_permissions",
	"role_metadata":    "sys_role_metadata",
	// 表名无 sys_ 前缀：schema 注解即如此（存量不一致，勿"顺手修正"——改了会和既有库对不上）
	"role_field_permissions": "role_field_permissions",
	"permissions":            "sys_permissions",
	"permission_groups":      "sys_permission_groups",
	"permission_apis":        "sys_permission_apis",
	"permission_menus":       "sys_permission_menus",
	"permission_policies":    "sys_permission_policies",
	"user_roles":             "sys_user_roles",
	"user_org_units":         "sys_user_org_units",
	"user_positions":         "sys_user_positions",
	"memberships":            "sys_memberships",
	"membership_org_units":   "sys_membership_org_units",
	"membership_positions":   "sys_membership_positions",
	"membership_roles":       "sys_membership_roles",
	"menus":                  "sys_menus",
	"apis":                   "sys_apis",
	"login_policies":         "sys_login_policies",
}

// RestoreCoreTables 把备份 JSON 的 data 段恢复进核心表，返回各表行数。
//
// 安全性：**仅允许向空库恢复**——任一目标表非空即整体中止（disaster recovery
// 语义是"全新库 + 恢复 + 启动"，不是对活库做合并；活库合并涉及 ID 冲突/
// 外键顺序/序列回填，超出工具职责）。恢复保留原 ID，按 ID 升序逐表插入，
// 单事务：任一表失败整体回滚，不留半恢复状态。
//
// 实现取表无关的动态 INSERT（备份 JSON 的键即列名——ent 实体 JSON tag 与
// 列名一致，且不含边字段），避免逐表手写字段映射漏列。表名走固定白名单，
// 列名经引号包裹，无注入面。
//
// 与上游的差异：PostgreSQL 下恢复完成后对各表做 setval 序列对齐（best-effort）——
// 显式 ID 插入不会推进 IDENTITY/serial 序列，不对齐则恢复后首个 insert 撞主键。
func (r *BackupRepo) RestoreCoreTables(ctx context.Context, rawData json.RawMessage) (map[string]int64, error) {
	var tables map[string]json.RawMessage
	if err := json.Unmarshal(rawData, &tables); err != nil {
		return nil, fmt.Errorf("decode backup data failed: %w", err)
	}

	db := r.entClient.DB()
	dialect := r.entClient.Driver().Dialect()

	// 前置校验：全部目标表必须为空（任一非空即中止）
	for _, name := range restoreTableOrder {
		count, err := countRestoreTable(ctx, db, restorePhysicalTables[name])
		if err != nil {
			return nil, err
		}
		if count != 0 {
			return nil, fmt.Errorf("table %q is not empty (%d rows): restore only allowed into empty tables", name, count)
		}
	}

	// 原生 SQL 事务（动态 INSERT 不经 ent builder）
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin restore transaction failed: %w", err)
	}
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	counts := make(map[string]int64, len(restoreTableOrder))
	for _, name := range restoreTableOrder {
		raw, ok := tables[name]
		if !ok || strings.TrimSpace(string(raw)) == "" || string(raw) == "null" {
			counts[name] = 0
			continue
		}
		n, rerr := restoreOneTable(ctx, tx, restorePhysicalTables[name], raw)
		if rerr != nil {
			return nil, fmt.Errorf("restore table %q failed: %w", name, rerr)
		}
		counts[name] = n
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit restore failed: %w", err)
	}
	committed = true

	// 序列对齐（事务外、best-effort）：只对 PostgreSQL 有意义。
	if dialect == "postgres" || dialect == "postgresql" {
		for _, name := range restoreTableOrder {
			if counts[name] == 0 {
				continue
			}
			r.alignPostgresSequence(ctx, db, restorePhysicalTables[name])
		}
	}

	return counts, nil
}

// alignPostgresSequence 把表的 ID 序列对齐到当前 max(id)（best-effort，失败仅留日志）。
// 序列名经 pg_get_serial_sequence 取（IDENTITY 与 serial 两种形态都覆盖）。
func (r *BackupRepo) alignPostgresSequence(ctx context.Context, db *sql.DB, table string) {
	var seq sql.NullString
	if err := db.QueryRowContext(ctx,
		`SELECT pg_get_serial_sequence('"`+table+`"', 'id')`).Scan(&seq); err != nil || !seq.Valid || seq.String == "" {
		r.log.Errorf("backup: query serial sequence of %q failed: %v", table, err)
		return
	}
	if _, err := db.ExecContext(ctx,
		fmt.Sprintf(`SELECT setval('%s', GREATEST((SELECT COALESCE(MAX(id), 0) FROM "%s"), 1))`, seq.String, table)); err != nil {
		r.log.Errorf("backup: align sequence of %q failed: %v", table, err)
		return
	}
	r.log.Infof("backup: sequence of %q aligned", table)
}

// restoreOneTable 动态 INSERT 单表：行按 id 升序稳定排序。
func restoreOneTable(ctx context.Context, tx *sql.Tx, name string, raw json.RawMessage) (int64, error) {
	var rows []map[string]any
	if err := json.Unmarshal(raw, &rows); err != nil {
		return 0, fmt.Errorf("decode rows failed: %w", err)
	}
	if len(rows) == 0 {
		return 0, nil
	}

	// 行内列集取并集（ent 实体 JSON 键一致，但个别行可能缺 omitempty 字段）。
	// "edges" 键丢弃：ent 实体序列化时带空 edges 对象，非表列。
	colSet := map[string]bool{}
	for _, row := range rows {
		for col := range row {
			if col == "edges" {
				continue
			}
			colSet[col] = true
		}
	}
	cols := make([]string, 0, len(colSet))
	for col := range colSet {
		cols = append(cols, col)
	}
	sort.Strings(cols)

	// 按 id 升序稳定排序（被引用行先行；无 id 键保持原序）
	sort.SliceStable(rows, func(i, j int) bool {
		a, aok := rows[i]["id"].(float64)
		b, bok := rows[j]["id"].(float64)
		if aok && bok {
			return a < b
		}
		return false
	})

	quotedCols := make([]string, 0, len(cols))
	for _, c := range cols {
		quotedCols = append(quotedCols, `"`+c+`"`)
	}
	placeholders := strings.TrimSuffix(strings.Repeat("?, ", len(cols)), ", ")
	insertSQL := fmt.Sprintf(`INSERT INTO %q (%s) VALUES (%s)`,
		name, strings.Join(quotedCols, ", "), placeholders)

	var total int64
	for _, row := range rows {
		args := make([]any, 0, len(cols))
		for _, c := range cols {
			args = append(args, row[c])
		}
		if _, err := tx.ExecContext(ctx, insertSQL, args...); err != nil {
			return total, fmt.Errorf("insert row (id=%v) failed: %w", row["id"], err)
		}
		total++
	}
	return total, nil
}

func countRestoreTable(ctx context.Context, db *sql.DB, name string) (int, error) {
	// name 走固定白名单（restoreTableOrder），无注入面
	var count int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM "`+name+`"`).Scan(&count); err != nil {
		return 0, fmt.Errorf("count table %q failed: %w", name, err)
	}
	return count, nil
}
