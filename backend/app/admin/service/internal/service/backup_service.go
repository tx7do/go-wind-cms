package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	taskV1 "go-wind-cms/api/gen/go/task/service/v1"
)

// BackupService 备份档案 BFF 面（只读）：代理 core BackupArchiveService。
// 恢复无入口（有意防误触）；下载为 bytes 内联（presigned 全层禁用）。
type BackupService struct {
	adminV1.BackupServiceHTTPServer

	log *log.Helper

	backupArchiveServiceClient taskV1.BackupArchiveServiceClient
}

func NewBackupService(
	ctx *bootstrap.Context,
	backupArchiveServiceClient taskV1.BackupArchiveServiceClient,
) *BackupService {
	return &BackupService{
		log:                        ctx.NewLoggerHelper("backup/service/admin-service"),
		backupArchiveServiceClient: backupArchiveServiceClient,
	}
}

func (s *BackupService) ListBackups(ctx context.Context, req *taskV1.ListBackupsRequest) (*taskV1.ListBackupsResponse, error) {
	return s.backupArchiveServiceClient.ListBackups(ctx, req)
}

func (s *BackupService) DownloadBackup(ctx context.Context, req *taskV1.DownloadBackupRequest) (*taskV1.DownloadBackupResponse, error) {
	return s.backupArchiveServiceClient.DownloadBackup(ctx, req)
}
