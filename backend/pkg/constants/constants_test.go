package constants

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── 角色码 ──

func TestRoleCodeConstantsPinned(t *testing.T) {
	// 角色码在 token 管理员标志、模板实例化、鉴权策略主体里按码值精确匹配，
	// 改名是破坏性变更（种子数据与既有库数据同时失效）。
	assert.Equal(t, ":", RoleCodeSpilt)
	assert.Equal(t, "platform:admin", PlatformAdminRoleCode)
	assert.Equal(t, "tenant:manager", TenantAdminRoleCode)
	assert.Equal(t, "tenant:user", TenantUserRoleCode)
	assert.Equal(t, "template:tenant:manager", TenantAdminTemplateRoleCode)
	assert.Equal(t, "platform:", PlatformRoleCodePrefix)
	assert.Equal(t, "tenant:", TenantRoleCodePrefix)
	assert.Equal(t, "template:", TemplateRoleCodePrefix)
}

func TestHasRoleCodePrefix(t *testing.T) {
	assert.True(t, HasRoleCodePrefix("platform:admin", "platform:"))
	assert.True(t, HasRoleCodePrefix("platform:", "platform:"), "exactly the prefix")
	assert.False(t, HasRoleCodePrefix("", "platform:"))
	assert.False(t, HasRoleCodePrefix("plat", "platform:"))
	assert.False(t, HasRoleCodePrefix("Platform:admin", "platform:"), "case sensitive")
	assert.True(t, HasRoleCodePrefix("anything", ""), "empty prefix matches")
}

func TestRoleCodeClassifiers(t *testing.T) {
	assert.True(t, IsPlatformRoleCode("platform:admin"))
	assert.False(t, IsPlatformRoleCode("tenant:manager"))
	assert.False(t, IsPlatformRoleCode("editor"))

	assert.True(t, IsTemplateRoleCode("template:tenant:manager"))
	assert.False(t, IsTemplateRoleCode("tenant:manager"))

	assert.True(t, IsReservedRoleCode("platform:admin"))
	assert.True(t, IsPlatformRoleCode("platform:custom-role"), "any platform: code is reserved")
	assert.True(t, IsReservedRoleCode("template:x"))
	assert.True(t, IsReservedRoleCode(TenantAdminRoleCode))
	assert.True(t, IsReservedRoleCode(TenantAdminTemplateRoleCode))
	assert.False(t, IsReservedRoleCode("tenant:user"), "tenant:user 码不在保留集合（仅 tenant:manager 保留）")
	assert.False(t, IsReservedRoleCode("editor"))
	assert.False(t, IsReservedRoleCode(""))
}

func TestExtractRoleCodeFromTemplate(t *testing.T) {
	assert.Equal(t, "tenant:manager", ExtractRoleCodeFromTemplate("template:tenant:manager"))
	assert.Equal(t, "plain", ExtractRoleCodeFromTemplate("plain"), "非模板码原样返回")
	assert.Equal(t, "", ExtractRoleCodeFromTemplate(""))
}

// ── 权限码 ──

func TestPermissionConstantsPinned(t *testing.T) {
	assert.Equal(t, "sys:", SystemPermissionCodePrefix)
	assert.Equal(t, "sys:access_backend", SystemAccessBackendPermissionCode)
	assert.Equal(t, "sys:access_app", SystemAccessAppPermissionCode)
	assert.Equal(t, "sys:manage_tenants", SystemManageTenantsPermissionCode)
	assert.Equal(t, "sys:audit_logs", SystemAuditLogsPermissionCode)
	assert.Equal(t, "sys:platform_admin", SystemPlatformAdminPermissionCode)
	assert.Equal(t, "sys:tenant_manager", SystemTenantManagerPermissionCode)
}

func TestIsReservedPermissionCode(t *testing.T) {
	assert.True(t, IsReservedPermissionCode("sys:access_backend"))
	assert.True(t, IsReservedPermissionCode("sys:custom"), "整个 sys: 命名空间保留")
	assert.False(t, IsReservedPermissionCode("post:create"))
	assert.False(t, IsReservedPermissionCode("sysx:foo"), "前缀须完整为 sys:")
	assert.False(t, IsReservedPermissionCode("sys"), "裸前缀不是权限码")
	assert.False(t, IsReservedPermissionCode(""))
}

func TestProtectedPermissionCodesAllReserved(t *testing.T) {
	require.NotEmpty(t, ProtectedPermissionCodes)
	for _, code := range ProtectedPermissionCodes {
		assert.True(t, IsReservedPermissionCode(code),
			"protected code %q must live in the reserved sys: namespace", code)
	}
	// 受保护集合必须覆盖登录闸门等核心码。
	assert.Contains(t, ProtectedPermissionCodes, SystemAccessBackendPermissionCode)
	assert.Contains(t, ProtectedPermissionCodes, SystemAccessAppPermissionCode)
}

// ── 口令策略配置键 ──

func TestPasswordPolicyConfigKeysPinned(t *testing.T) {
	assert.Equal(t, "sys.password.minLen", ConfigKeyPasswordMinLen)
	assert.Equal(t, "sys.password.maxAgeDays", ConfigKeyPasswordMaxAgeDays)
	assert.Equal(t, "sys.password.historyCount", ConfigKeyPasswordHistoryCount)
	assert.Equal(t, 8, DefaultPasswordMinLen)
	assert.Equal(t, 90, DefaultPasswordMaxAgeDays)
	assert.Equal(t, 3, DefaultPasswordHistoryCount)
}

// ── 租户关联模式 ──

func TestUserTenantRelation(t *testing.T) {
	// DefaultUserTenantRelationType 为一对一，IsTenantModeEnabled 据此编译期成立。
	assert.Equal(t, UserTenantRelationOneToOne, DefaultUserTenantRelationType)
	assert.True(t, IsTenantModeEnabled)
	assert.Equal(t, UserTenantRelationType(0), UserTenantRelationNone)
	assert.Equal(t, UserTenantRelationType(1), UserTenantRelationOneToOne)
	assert.Equal(t, UserTenantRelationType(2), UserTenantRelationOneToMany)
}
