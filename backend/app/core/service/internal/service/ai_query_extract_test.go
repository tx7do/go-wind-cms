package service

import "testing"

// 复现线上"no SQL statement found"：DeepSeek 在 SQL 前后带说明文字与围栏时的提取。
func TestExtractSQL_PrefixedFence(t *testing.T) {
	raw := "您这个问题我没办法直接帮您查——我这边没有连接到任何数据库或日志系统，看不到您环境里的实际记录。\n\n" +
		"如果您是在自己系统上排查，可以参考下面这个思路：\n" +
		"```sql\n" +
		"SELECT id, username, ip_address, created_at\n" +
		"FROM sys_login_audit_logs\n" +
		"WHERE status = 'FAILED'\n" +
		"  AND created_at >= NOW() - INTERVAL '24 hours'\n" +
		"ORDER BY created_at DESC\n" +
		"LIMIT 100;\n" +
		"```\n" +
		"几点说明，避免你直接跑出问题。"

	got, err := sanitizeSQLFor(raw, aiQueryTableWhitelist, nil)
	if err != nil {
		t.Fatalf("sanitizeSQL rejected: %v", err)
	}
	if !aiQuerySelectPattern.MatchString(got) {
		t.Fatalf("extracted SQL invalid: %q", got)
	}
	t.Logf("extracted: %s", got)
}
