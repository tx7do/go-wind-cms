package service

import (
	"context"
	"fmt"
	"io"
	"path"
	"sort"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/minio/minio-go/v7"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/timestamppb"

	"go-wind-cms/pkg/oss"

	taskV1 "go-wind-cms/api/gen/go/task/service/v1"
)

// maxBackupDownloadBytes 下载上限（64MB——超出说明备份异常膨胀，管理面不应整载）。
const maxBackupDownloadBytes = 64 << 20

// BackupArchiveService 备份档案只读视图：列出/下载 backups 桶里本项目前缀的
// 备份对象。备份写入在 backup asynq 任务（TaskService.AsyncBackup）；恢复无
// RPC 入口（有意防误触，走一次性程序 + 演练文档 docs/backup-restore-drill.md）。
type BackupArchiveService struct {
	taskV1.UnimplementedBackupArchiveServiceServer

	log         *log.Helper
	minioClient *oss.MinIOClient
}

func NewBackupArchiveService(
	ctx *bootstrap.Context,
	minIOClient *oss.MinIOClient,
) *BackupArchiveService {
	return &BackupArchiveService{
		log:         ctx.NewLoggerHelper("backup-archive/service/core-service"),
		minioClient: minIOClient,
	}
}

// guardBackupKey 校验对象键属于本项目备份前缀（防路径穿越与跨项目误取——
// dev MinIO 多项目共用桶，键越界一律拒绝）。
func guardBackupKey(objectKey string) error {
	clean := path.Clean("/" + strings.TrimLeft(objectKey, "/"))[1:]
	if !strings.HasPrefix(clean, backupObjectPrefix) || strings.Contains(clean, "..") {
		return fmt.Errorf("invalid backup object key")
	}
	return nil
}

// ListBackups 列出备份档案（时间倒序）。
// minioSDK 取底层 minio SDK client（nil 容错）。
func (s *BackupArchiveService) minioSDK() *minio.Client {
	if s.minioClient == nil {
		return nil
	}
	return s.minioClient.GetClient()
}

func (s *BackupArchiveService) ListBackups(ctx context.Context, req *taskV1.ListBackupsRequest) (*taskV1.ListBackupsResponse, error) {
	client := s.minioSDK()
	if client == nil {
		return nil, fmt.Errorf("oss client unavailable")
	}

	opts := minio.ListObjectsOptions{Recursive: true, Prefix: backupObjectPrefix}
	items := make([]*taskV1.BackupArchive, 0, 16)
	for obj := range client.ListObjects(ctx, backupBucket, opts) {
		if obj.Err != nil {
			s.log.Errorf("list backup objects failed: %s", obj.Err.Error())
			return nil, fmt.Errorf("list backup objects failed")
		}
		items = append(items, &taskV1.BackupArchive{
			ObjectKey:    obj.Key,
			DisplayName:  backupDisplayName(obj.Key),
			SizeBytes:    uint64(obj.Size),
			LastModified: timestamppb.New(obj.LastModified),
		})
	}
	sort.Slice(items, func(i, j int) bool {
		ti, tj := items[i].GetLastModified(), items[j].GetLastModified()
		if ti == nil || tj == nil {
			return items[i].GetObjectKey() > items[j].GetObjectKey()
		}
		return ti.AsTime().After(tj.AsTime())
	})
	return &taskV1.ListBackupsResponse{Items: items}, nil
}

// DownloadBackup 下载备份档案（bytes 内联；presigned 全层禁用，走后端鉴权流）。
func (s *BackupArchiveService) DownloadBackup(ctx context.Context, req *taskV1.DownloadBackupRequest) (*taskV1.DownloadBackupResponse, error) {
	client := s.minioSDK()
	if client == nil {
		return nil, fmt.Errorf("oss client unavailable")
	}
	objectKey := req.GetObjectKey()
	if err := guardBackupKey(objectKey); err != nil {
		return nil, err
	}

	obj, err := client.GetObject(ctx, backupBucket, objectKey, minio.GetObjectOptions{})
	if err != nil {
		s.log.Errorf("get backup object failed: %s", err.Error())
		return nil, fmt.Errorf("get backup object failed")
	}
	defer obj.Close()

	// 备份是 gzip JSON，量级 MB 级；一次性读入（管理面下载低频，不设流式），
	// 超上限即拒绝——异常膨胀的备份不应被管理页整载进内存。
	buf, err := io.ReadAll(io.LimitReader(obj, maxBackupDownloadBytes+1))
	if err != nil {
		return nil, fmt.Errorf("read backup object failed")
	}
	if len(buf) == 0 {
		return nil, fmt.Errorf("backup object not found")
	}
	if len(buf) > maxBackupDownloadBytes {
		return nil, fmt.Errorf("backup object exceeds %d bytes", maxBackupDownloadBytes)
	}

	return &taskV1.DownloadBackupResponse{
		ContentType: "application/gzip",
		File:        buf,
	}, nil
}

// backupDisplayName 对象键 → 展示名（去 gowind-cms/ 前缀与日期目录，保留文件名）。
func backupDisplayName(objectKey string) string {
	name := strings.TrimPrefix(objectKey, backupObjectPrefix)
	if idx := strings.LastIndex(name, "/"); idx >= 0 {
		name = name[idx+1:]
	}
	return name
}
