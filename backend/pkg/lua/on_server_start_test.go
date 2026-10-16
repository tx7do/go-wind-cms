package lua

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-kratos/kratos/v2/log"
)

func TestOnServerStartScript(t *testing.T) {
	// Create engine
	config := DefaultConfig()
	config.ScriptDir = "" // Don't auto-load
	engine := NewEngine(config, log.DefaultLogger)
	defer engine.Close()

	// 自写脚本文件：on_server_start.lua 从未入库（仓库无此 fixture），
	// 测试自带最小脚本以保持 LoadScriptFile + hook 执行的覆盖
	scriptPath := filepath.Join(t.TempDir(), "on_server_start.lua")
	script := `
function on_server_start(ctx)
  local config = ctx.data.config
  if config == nil or config.server == nil then
    error("config.server is required")
  end
  return { started = true, server = config.server }
end
`
	if werr := os.WriteFile(scriptPath, []byte(script), 0o644); werr != nil {
		t.Fatalf("write temp script failed: %v", werr)
	}
	err := engine.LoadScriptFile(context.Background(), scriptPath)
	if err != nil {
		t.Fatalf("Failed to load on_server_start.lua: %v", err)
	}

	// Execute the hook with test data
	ctx := &Context{
		ID:       "server_startup",
		HookName: "on_server_start",
		Data: map[string]interface{}{
			"config": map[string]interface{}{
				"server": "test-server",
				"port":   8080,
			},
			"service": map[string]interface{}{
				"name":    "go-wind-cms",
				"version": "1.0.0",
			},
		},
	}

	err = engine.ExecuteHook(context.Background(), "on_server_start", ctx)
	if err != nil {
		t.Fatalf("Failed to execute on_server_start hook: %v", err)
	}

	t.Log("✓ on_server_start script executed successfully")
	t.Log("✓ All logs from the script should be visible above")
}
