package server

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	bConfig "github.com/tx7do/kratos-bootstrap/config"
)

// writeCollabConfig 落一份临时 server.yaml 并完成扩展注册 + 配置装载，
// 返回可用于 NewCollab/collabExtConfig 的引导上下文。
func writeCollabConfig(t *testing.T, collabSection string) *bootstrap.Context {
	t.Helper()
	dir := t.TempDir()
	body := "server:\n  rest:\n    addr: \"0.0.0.0:6600\"\n" + collabSection
	if err := os.WriteFile(filepath.Join(dir, "server.yaml"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	ctx := bootstrap.NewContextWithParam(context.Background(), nil, nil, log.DefaultLogger)
	RegisterCollabConfigExt(ctx)
	if err := bConfig.LoadBootstrapConfig(dir); err != nil {
		t.Fatal(err)
	}
	return ctx
}

// TestCollabConfigScan server.collab 段须能经整树 Scan 落进扩展消息。
// 机制级回归：注册晚于 LoadBootstrapConfig 时会静默失效（无任何报错）。
func TestCollabConfigScan(t *testing.T) {
	ctx := writeCollabConfig(t, "  collab:\n    addr: \":7702\"\n    allowed_origins:\n      - \"https://admin.example.com\"\n")

	got := collabExtConfig(ctx)
	if got.GetAddr() != ":7702" {
		t.Fatalf("addr = %q, want \":7702\"", got.GetAddr())
	}
	if origins := got.GetAllowedOrigins(); len(origins) != 1 || origins[0] != "https://admin.example.com" {
		t.Fatalf("allowed_origins = %v", origins)
	}
}

// TestNewCollabYamlFirst 配置段优先于环境变量。
func TestNewCollabYamlFirst(t *testing.T) {
	t.Setenv("COLLAB_ADDR", ":7803")
	ctx := writeCollabConfig(t, "  collab:\n    addr: \":7702\"\n")

	svc, srv := NewCollab(ctx, nil, nil, nil)
	if svc == nil {
		t.Fatal("service must not be nil")
	}
	if srv == nil || srv.addr != ":7702" {
		t.Fatalf("addr = %v, want yaml \":7702\" to win over env :7803", srv)
	}
}

// TestNewCollabEnvFallback 未配置 yaml 段时回退环境变量。
func TestNewCollabEnvFallback(t *testing.T) {
	t.Setenv("COLLAB_ADDR", ":7803")
	ctx := writeCollabConfig(t, "")

	_, srv := NewCollab(ctx, nil, nil, nil)
	if srv == nil || srv.addr != ":7803" {
		t.Fatalf("addr = %v, want env \":7803\"", srv)
	}
}

// TestNewCollabDisabled 配置段显式停用：server 为 nil、service 仍返回
// （供删除踢房 no-op）。
func TestNewCollabDisabled(t *testing.T) {
	ctx := writeCollabConfig(t, "  collab:\n    addr: \"off\"\n")

	svc, srv := NewCollab(ctx, nil, nil, nil)
	if srv != nil {
		t.Fatalf("server should be nil when disabled, got %v", srv)
	}
	if svc == nil {
		t.Fatal("service must survive disable")
	}
}
