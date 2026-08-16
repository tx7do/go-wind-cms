package viewer

import (
	"github.com/tx7do/go-crud/viewer"
)

// AnonymousTenantViewer 描述一个未登录但已解析出 tenant 归属的匿名访问者。
//
// 由 app BFF 的 ent 中间件在白名单（公开）请求中，按请求 Host 解析出 tenant_id
// 后注入。与 SystemViewer 的关键差异：
//   - TenantID() 返回解析到的真实 tenant_id（>0），使显式 maybeTenantFromViewer
//     查询能按 tenant 隔离，而非拿不到 tenant 返回空/401。
//   - IsSystemContext()==false，使 ent TenantPrivacy rule 不再 bypass，按 tenant
//     隔离过滤（堵住 SystemViewer 跨租户可见的洞）。
//   - UserID()/OrgUnitID()==0、HasPermission()==false、DataScope()==[]，
//     无写权限、无个人身份，仅可读所属 tenant 的公开数据。
//
// 解析失败（Host 未配 domain 等）不回退此 viewer，而由中间件注入 viewer.NewNoopContext
// fail-closed。
type AnonymousTenantViewer struct {
	tid     uint64
	traceID string
}

// NewAnonymousTenantViewer 构造一个绑定指定 tenant 的只读匿名 viewer。
func NewAnonymousTenantViewer(tid uint64, traceID string) viewer.Context {
	return AnonymousTenantViewer{tid: tid, traceID: traceID}
}

// UserID 返回当前用户ID（匿名访问者无用户身份）
func (v AnonymousTenantViewer) UserID() uint64 {
	return 0
}

// TenantID 返回解析到的租户ID
func (v AnonymousTenantViewer) TenantID() uint64 {
	return v.tid
}

// OrgUnitID 返回当前身份挂载的组织单元 ID（匿名无）
func (v AnonymousTenantViewer) OrgUnitID() uint64 {
	return 0
}

// Permissions 返回当前 Viewer 的权限列表（匿名无权限）
func (v AnonymousTenantViewer) Permissions() []string {
	return []string{}
}

// Roles 返回当前 Viewer 的角色列表（匿名无角色）
func (v AnonymousTenantViewer) Roles() []string {
	return []string{}
}

// DataScope 返回当前身份的数据权限范围（匿名无）
func (v AnonymousTenantViewer) DataScope() []viewer.DataScope {
	return []viewer.DataScope{}
}

// TraceID 返回当前请求的 Trace ID（用于日志跟踪）
func (v AnonymousTenantViewer) TraceID() string {
	return v.traceID
}

// HasPermission 判断是否具有某个动作/资源的权限（匿名一律拒绝）
func (v AnonymousTenantViewer) HasPermission(_, _ string) bool {
	return false
}

// IsPlatformContext 当前是否处于平台管理视图（匿名非平台）
func (v AnonymousTenantViewer) IsPlatformContext() bool {
	return false
}

// IsTenantContext 当前是否处于租户业务视图（已解析到 tenant，是）
func (v AnonymousTenantViewer) IsTenantContext() bool {
	return v.tid > 0
}

// IsSystemContext 判断是否为系统后台任务（匿名非系统）
func (v AnonymousTenantViewer) IsSystemContext() bool {
	return false
}

// ShouldAudit 返回是否需要记录审计日志
func (v AnonymousTenantViewer) ShouldAudit() bool {
	return false
}
