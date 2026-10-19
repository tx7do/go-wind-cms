package service

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/stretchr/testify/require"
)

func TestParseBackupRetentionDays(t *testing.T) {
	require.Equal(t, backupRetentionDefault, parseBackupRetentionDays(""), "未设回默认")
	require.Equal(t, backupRetentionDefault, parseBackupRetentionDays("abc"), "非法回默认")
	require.Equal(t, 7, parseBackupRetentionDays("7"))
	require.Equal(t, 0, parseBackupRetentionDays("0"), "0=显式关闭清理")
	require.Equal(t, -1, parseBackupRetentionDays(" -1 "), "负数=显式关闭清理")
}

func TestExpiredBackupKeys(t *testing.T) {
	now := time.Now()
	cutoff := now.AddDate(0, 0, -30)

	metas := []backupObjectMeta{
		{Key: "gowind-cms/2026/01/01/old-000000.json.gz", LastModified: now.AddDate(0, 0, -40)},
		{Key: "gowind-cms/2026/07/01/recent-120000.json.gz", LastModified: now.AddDate(0, 0, -1)},
		{Key: "gowind-cms/2026/06/30/boundary-235959.json.gz", LastModified: cutoff},
	}

	// 契约只看时间：项目前缀隔离是 listObjects 扫描侧的事（见编排用例），
	// 喂进来的对象只要早于 cutoff 就在删除清单里。
	expired := expiredBackupKeys(metas, cutoff)
	require.Len(t, expired, 1, "只清早于 cutoff 的（cutoff 当刻不算过期）")
	require.Equal(t, "gowind-cms/2026/01/01/old-000000.json.gz", expired[0])
}

// fakeBackupCleaner 纯内存假件：钉住 cleanupOldBackups 的扫描/删除编排。
type fakeBackupCleaner struct {
	metas    []backupObjectMeta
	removed  []string
	lastScan string
	listErr  error
}

func (f *fakeBackupCleaner) listObjects(_ context.Context, _, prefix string) ([]backupObjectMeta, error) {
	f.lastScan = prefix
	if f.listErr != nil {
		return nil, f.listErr
	}
	// 模拟 MinIO ListObjectsOptions.Prefix 语义：只返回命中前缀的对象。
	// 真实环境里"不误删他项目对象"的隔离保证就住在这一层。
	out := make([]backupObjectMeta, 0, len(f.metas))
	for _, m := range f.metas {
		if strings.HasPrefix(m.Key, prefix) {
			out = append(out, m)
		}
	}
	return out, nil
}

func (f *fakeBackupCleaner) removeObject(_ context.Context, _, key string) error {
	f.removed = append(f.removed, key)
	return nil
}

func TestCleanupOldBackupsScansProjectPrefix(t *testing.T) {
	svc := &TaskService{log: log.NewHelper(log.DefaultLogger)}
	now := time.Now()

	cleaner := &fakeBackupCleaner{metas: []backupObjectMeta{
		{Key: "gowind-cms/2026/01/01/old.json.gz", LastModified: now.AddDate(0, 0, -40)},
		{Key: "gowind-cms/2026/07/01/new.json.gz", LastModified: now.AddDate(0, 0, -1)},
		{Key: "other/2025/01/01/foreign.json.gz", LastModified: now.AddDate(0, 0, -400)},
	}}

	svc.cleanupOldBackups(context.Background(), cleaner)

	// 假件按前缀过滤返回（模拟 MinIO ListObjects 的 Prefix 语义）：
	// other/ 前缀的对象根本不该出现在扫描结果里，故即便它最老也不会被删。
	require.Equal(t, backupObjectPrefix, cleaner.lastScan, "只扫描本项目前缀（共享桶不误删他项目对象）")
	require.Equal(t, []string{"gowind-cms/2026/01/01/old.json.gz"}, cleaner.removed)
}
