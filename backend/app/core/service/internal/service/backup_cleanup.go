package service

import (
	"context"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"

	"go-wind-cms/pkg/oss"
)

// backupRetentionEnv 备份保留期环境变量（天）。
// 运维面走环境变量而非配置表：清理策略与部署形态（磁盘/桶容量）强相关而与业务租户无关。
// 非正值（0/负数）显式关闭清理——运维显式关掉比悄悄清空更安全。
const backupRetentionEnv = "BACKUP_RETENTION_DAYS"

// backupRetentionDefault 备份保留期默认值（天）。
const backupRetentionDefault = 30

// backupObjectPrefix 备份对象键前缀。
// 与上游（整桶扫描）不同：dev/小部署的 MinIO 常多项目共用一个 backups 桶，
// 整桶按保留期删除会误删其它项目的备份——本仓对象带项目前缀，清理只扫该前缀。
const backupObjectPrefix = "gowind-cms/"

// backupObjectMeta 清理所需的备份对象元数据（key + 修改时间）。
type backupObjectMeta struct {
	Key          string
	LastModified time.Time
}

// backupCleaner 抽象对象存储的"列+删"。清理逻辑只依赖此接口：
// 单测用假件枚举对象，不依赖真实 MinIO。
type backupCleaner interface {
	listObjects(ctx context.Context, bucket, prefix string) ([]backupObjectMeta, error)
	removeObject(ctx context.Context, bucket, key string) error
}

// minioBackupCleaner 基于 MinIOClient 暴露的 SDK 客户端的清理适配器。
type minioBackupCleaner struct {
	client *minio.Client
}

func newMinioBackupCleaner(mc *oss.MinIOClient) *minioBackupCleaner {
	if mc == nil {
		return nil
	}
	return &minioBackupCleaner{client: mc.GetClient()}
}

func (c *minioBackupCleaner) listObjects(ctx context.Context, bucket, prefix string) ([]backupObjectMeta, error) {
	if c.client == nil {
		return nil, nil
	}
	opts := minio.ListObjectsOptions{Recursive: true, Prefix: prefix}
	var metas []backupObjectMeta
	for obj := range c.client.ListObjects(ctx, bucket, opts) {
		if obj.Err != nil {
			return nil, obj.Err
		}
		metas = append(metas, backupObjectMeta{Key: obj.Key, LastModified: obj.LastModified})
	}
	return metas, nil
}

func (c *minioBackupCleaner) removeObject(ctx context.Context, bucket, key string) error {
	return c.client.RemoveObject(ctx, bucket, key, minio.RemoveObjectOptions{})
}

// parseBackupRetentionDays 解析保留期环境变量；未设/非法回默认 30。
// 非正值（0/负数）原样返回，视为"关闭清理"。
func parseBackupRetentionDays(raw string) int {
	if raw == "" {
		return backupRetentionDefault
	}
	v, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil {
		return backupRetentionDefault
	}
	return v
}

// expiredBackupKeys 挑出 LastModified 早于 cutoff 的对象 key（纯函数，可测）。
func expiredBackupKeys(metas []backupObjectMeta, cutoff time.Time) []string {
	var expired []string
	for _, meta := range metas {
		if meta.LastModified.Before(cutoff) {
			expired = append(expired, meta.Key)
		}
	}
	return expired
}

// cleanupOldBackups 清理 backups 桶中本项目前缀下超过保留期的备份对象（best-effort）。
//
// 在 AsyncBackup 上传成功后调用：备份本体已成功，清理是顺手维护——
// 单个对象删除失败只留日志并继续删其余，整段失败也不影响任务结论
// （把备份任务标 FAILED 会让告警指向错误的故障点）。错误全部 Errorf 留痕。
// 关闭方式：BACKUP_RETENTION_DAYS<=0（显式关闭，不悄悄清空）。
func (s *TaskService) cleanupOldBackups(ctx context.Context, cleaner backupCleaner) {
	retentionDays := parseBackupRetentionDays(os.Getenv(backupRetentionEnv))
	if retentionDays <= 0 {
		s.log.Infof("backup cleanup: disabled (%s<=%d)", backupRetentionEnv, retentionDays)
		return
	}

	metas, err := cleaner.listObjects(ctx, backupBucket, backupObjectPrefix)
	if err != nil {
		s.log.Errorf("backup cleanup: list objects failed: %s", err.Error())
		return
	}

	cutoff := time.Now().AddDate(0, 0, -retentionDays)
	expired := expiredBackupKeys(metas, cutoff)
	if len(expired) == 0 {
		s.log.Infof("backup cleanup: %d objects scanned, none older than %dd", len(metas), retentionDays)
		return
	}

	var removed, failed int
	for _, key := range expired {
		if err := cleaner.removeObject(ctx, backupBucket, key); err != nil {
			failed++
			s.log.Errorf("backup cleanup: remove %s failed: %s", key, err.Error())
			continue
		}
		removed++
	}
	s.log.Infof("backup cleanup: %d objects scanned, %d expired, %d removed, %d failed (retention=%dd)",
		len(metas), len(expired), removed, failed, retentionDays)
}
