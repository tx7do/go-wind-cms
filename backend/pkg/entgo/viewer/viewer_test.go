package viewer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tx7do/go-crud/viewer"

	identityV1 "go-wind-cms/api/gen/go/identity/service/v1"
)

// ── ConvertDataScope ──

func TestConvertDataScope(t *testing.T) {
	tests := []struct {
		name      string
		in        identityV1.DataScope
		wantScope viewer.ScopeType
	}{
		{"ALL", identityV1.DataScope_ALL, viewer.ScopeTypeAll},
		{"UNIT_ONLY", identityV1.DataScope_UNIT_ONLY, viewer.ScopeTypeUnit},
		{"UNIT_AND_CHILD", identityV1.DataScope_UNIT_AND_CHILD, viewer.ScopeTypeUnit},
		{"SELF", identityV1.DataScope_SELF, viewer.ScopeTypeSelf},
		{"unknown falls to NONE", identityV1.DataScope(99), viewer.ScopeTypeNone},
		{"zero falls to default", identityV1.DataScope(0), viewer.ScopeTypeNone},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := ConvertDataScope(tt.in)
			assert.Equal(t, tt.wantScope, ds.ScopeType)
			assert.Empty(t, ds.TargetIDs)
		})
	}
}

// ── UserViewer ──

func TestNewUserViewer(t *testing.T) {
	vc := NewUserViewer(42, 7, 3, "trace-1", identityV1.DataScope_ALL)
	requireVC(t, vc, expect{
		uid: 42, tid: 7, ouid: 3, trace: "trace-1",
		tenant: true, platform: false, system: false,
		scope: viewer.ScopeTypeAll,
	})
}

func TestNewUserViewerTenantZeroIsPlatform(t *testing.T) {
	// tid==0 → 平台管理视图；tid>0 → 租户业务视图。isPlatform 与 isTenant 互斥。
	vc := NewUserViewer(42, 0, 0, "", identityV1.DataScope_UNIT_AND_CHILD)
	requireVC(t, vc, expect{
		uid: 42, tid: 0, ouid: 0, trace: "",
		tenant: false, platform: true, system: false,
		scope: viewer.ScopeTypeUnit,
	})
}

func TestNewUserViewerWithDataScopes(t *testing.T) {
	scopes := []viewer.DataScope{
		{ScopeType: viewer.ScopeTypeUnit, TargetIDs: []int64{1, 2, 3}},
	}
	vc := NewUserViewerWithDataScopes(1, 2, 3, "t", scopes)
	got := vc.DataScope()
	assert.Equal(t, scopes, got, "WithDataScopes 透传（broker Subscriber 重建 viewer 的路径）")
}

func TestUserViewerNoPermissionByDefault(t *testing.T) {
	// HasPermission 恒 false：细粒度权限判定不走 viewer（走 permission service）。
	vc := NewUserViewer(42, 7, 3, "", identityV1.DataScope_ALL)
	assert.False(t, vc.HasPermission("update", "user"))
	assert.False(t, vc.HasPermission("", ""))
	assert.Empty(t, vc.Permissions())
	assert.Empty(t, vc.Roles())
	assert.False(t, vc.ShouldAudit())
}

// ── AnonymousTenantViewer ──

func TestAnonymousTenantViewerSemantics(t *testing.T) {
	vc := NewAnonymousTenantViewer(9, "trace-anon")
	requireVC(t, vc, expect{
		uid: 0, tid: 9, ouid: 0, trace: "trace-anon",
		tenant: true, platform: false, system: false,
		scope: "", // 无 datascope
	})

	// 与 SystemViewer 的关键差异：无 bypass 权限、无个人身份。
	assert.False(t, vc.HasPermission("read", "post"), "匿名一律拒绝")
	assert.Empty(t, vc.Permissions())
	assert.Empty(t, vc.Roles())
	assert.Empty(t, vc.DataScope())
	assert.False(t, vc.ShouldAudit())
}

func TestAnonymousTenantViewerZeroTenant(t *testing.T) {
	// tid=0 时 IsTenantContext 为 false（构造方约束 tid>0，防御面钉测）。
	vc := NewAnonymousTenantViewer(0, "")
	assert.False(t, vc.IsTenantContext())
	assert.False(t, vc.IsPlatformContext())
	assert.False(t, vc.IsSystemContext())
}

// ── SystemViewer ──

func TestSystemViewerSemantics(t *testing.T) {
	vc := NewSystemViewer()
	requireVC(t, vc, expect{
		uid: 0, tid: 0, ouid: 0, trace: "",
		tenant: false, platform: true, system: true,
	})

	assert.True(t, vc.HasPermission("anything", "anywhere"), "系统上下文全量 bypass")
	assert.False(t, vc.ShouldAudit(), "系统后台任务不进审计")
	assert.Empty(t, vc.Permissions())
}

func TestSystemViewerContext(t *testing.T) {
	ctx := NewSystemViewerContext(context.Background())
	got, ok := viewer.FromContext(ctx)
	require.True(t, ok)
	_, isSystem := got.(SystemViewer)
	assert.True(t, isSystem)
	assert.True(t, got.IsSystemContext())
}

// ── 三种 viewer 的隔离语义对照（TenantPrivacy rule 的判定输入）──

func TestViewerIsolationMatrix(t *testing.T) {
	matrix := []struct {
		name     string
		vc       viewer.Context
		tenant   bool
		platform bool
		system   bool
	}{
		{"UserViewer(tenant)", NewUserViewer(1, 7, 0, "", identityV1.DataScope_ALL), true, false, false},
		{"UserViewer(platform)", NewUserViewer(1, 0, 0, "", identityV1.DataScope_ALL), false, true, false},
		{"AnonymousTenantViewer", NewAnonymousTenantViewer(7, ""), true, false, false},
		{"SystemViewer", NewSystemViewer(), false, true, true},
	}
	for _, m := range matrix {
		t.Run(m.name, func(t *testing.T) {
			assert.Equal(t, m.tenant, m.vc.IsTenantContext())
			assert.Equal(t, m.platform, m.vc.IsPlatformContext())
			assert.Equal(t, m.system, m.vc.IsSystemContext())
		})
	}
}

// ── 断言辅助 ──

type expect struct {
	uid, tid, ouid int64
	trace          string
	tenant         bool
	platform       bool
	system         bool
	scope          viewer.ScopeType
}

func requireVC(t *testing.T, vc viewer.Context, want expect) {
	t.Helper()
	assert.Equal(t, want.uid, vc.UserID())
	assert.Equal(t, want.tid, vc.TenantID())
	assert.Equal(t, want.ouid, vc.OrgUnitID())
	assert.Equal(t, want.trace, vc.TraceID())
	assert.Equal(t, want.tenant, vc.IsTenantContext())
	assert.Equal(t, want.platform, vc.IsPlatformContext())
	assert.Equal(t, want.system, vc.IsSystemContext())
	if want.scope != "" {
		scopes := vc.DataScope()
		if assert.Len(t, scopes, 1) {
			assert.Equal(t, want.scope, scopes[0].ScopeType)
		}
	}
}
