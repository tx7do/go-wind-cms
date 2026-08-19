package service

import (
	"context"
	"io"
	"net/http"
	"path"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	"go-wind-cms/app/core/service/internal/data"

	storageV1 "go-wind-cms/api/gen/go/storage/service/v1"

	"go-wind-cms/pkg/netutil"
	"go-wind-cms/pkg/oss"

	"github.com/tx7do/go-crud/viewer"
)

// FileTransferService 是 core 端的文件传输服务实现。
// 注意：该服务当前未在 grpc_server.go 中注册（仍是死代码），
// 实际生效的上传链路在 admin/app 端 FileTransferService（已正确使用
// auth.FromContext）。若日后在此启用，须参照 admin/app 端为 DownloadFile
// 补做 fileId 归属校验，并禁用裸 StorageObject/DownloadUrl selector。
type FileTransferService struct {
	storageV1.UnimplementedFileTransferServiceServer

	log *log.Helper

	mc       *oss.MinIOClient
	fileRepo *data.FileRepo
}

func NewFileTransferService(
	ctx *bootstrap.Context,
	mc *oss.MinIOClient,
	fileRepo *data.FileRepo,
) *FileTransferService {
	return &FileTransferService{
		log:      ctx.NewLoggerHelper("file-transfer/service/core-service"),
		mc:       mc,
		fileRepo: fileRepo,
	}
}

func parseKey(key string) (folder, filename, ext string) {
	if key == "" {
		return "", "", ""
	}

	// 统一去除前导斜杠，但保留中间路径
	key = strings.TrimPrefix(key, "/")

	// 如果以 '/' 结尾，则视为目录
	if strings.HasSuffix(key, "/") {
		f := strings.TrimSuffix(key, "/")
		return f, "", ""
	}

	// 目录部分
	dir := path.Dir(key)
	if dir == "." {
		dir = ""
	}

	base := path.Base(key)

	// 处理点文件（如 .env）：当且仅当只有一个前导点且没有其他点，视为无扩展名
	if strings.HasPrefix(base, ".") && strings.Count(base, ".") == 1 {
		return dir, base, ""
	}

	// 查找最后一个点作为扩展名分隔（点在开头不算）
	idx := strings.LastIndex(base, ".")
	if idx <= 0 {
		// 无扩展名或点在首位（已处理首位点情况）
		return dir, base, ""
	}

	name := base[:idx]
	ext = strings.ToLower(base[idx+1:])

	return dir, name, ext
}

// recordFile 记录文件元数据到数据库
func (s *FileTransferService) recordFile(
	ctx context.Context,
	tenantID, userID uint32,
	sourceFileName string,
	info minio.UploadInfo,
	downloadUrl string,
) error {

	dir, fileName, ext := parseKey(info.Key)

	if _, err := s.fileRepo.Create(ctx, &storageV1.CreateFileRequest{
		Data: &storageV1.File{
			Provider:      trans.Ptr(storageV1.OSSProvider_MINIO),
			BucketName:    trans.Ptr(info.Bucket),
			SaveFileName:  trans.Ptr(fileName + "." + ext),
			FileDirectory: trans.Ptr(dir),
			FileName:      trans.Ptr(sourceFileName),
			Extension:     trans.Ptr(ext),
			FileGuid:      trans.Ptr(uuid.New().String()),
			Size:          trans.Ptr(uint64(info.Size)),
			LinkUrl:       trans.Ptr(downloadUrl),
			CreatedBy:     trans.Ptr(userID),
			TenantId:      trans.Ptr(tenantID),
		},
	}); err != nil {
		s.log.Errorf("Failed to create file record: %v", err)
		return err
	}
	return nil
}

// directUploadFile 直接上传文件
func (s *FileTransferService) directUploadFile(ctx context.Context, req *storageV1.UploadFileRequest) (*storageV1.UploadFileResponse, error) {
	if req == nil || req.StorageObject == nil {
		return nil, storageV1.ErrorUploadFailed("unknown storageObject")
	}

	// 流式上传：reader 由 handler 通过 context 注入。该服务当前未注册，
	// 此处 reader 恒为 nil，保持与 admin/app 端一致的签名以便日后启用。
	reader, objectSize := oss.UploadReaderFromContext(ctx)
	if reader == nil || objectSize <= 0 {
		return nil, storageV1.ErrorUploadFailed("unknown fileData")
	}

	if req.GetMime() == "" {
		return nil, storageV1.ErrorUploadFailed("unknown mime type")
	}

	if req.GetSourceFileName() == "" {
		return nil, storageV1.ErrorUploadFailed("unknown source file name")
	}

	if req.StorageObject.BucketName == nil {
		req.StorageObject.BucketName = trans.Ptr(oss.ContentTypeToBucketName(req.GetMime()))
	}

	if req.StorageObject.ObjectName == nil {
		req.StorageObject.ObjectName = trans.Ptr(
			oss.EnsureObjectName(
				req.GetStorageObject().GetFileDirectory(),
				req.GetSourceFileName(),
				req.GetMime(),
				oss.GenerateFileNameTypeUUID,
			),
		)
	}

	info, _, downloadUrl, err := s.mc.UploadFile(
		ctx,
		req.GetStorageObject().GetBucketName(),
		req.GetStorageObject().GetObjectName(),
		req.GetMime(),
		reader, objectSize,
	)
	if err != nil {
		return nil, err
	}

	// 从 viewer 取操作者身份，而非信任请求体里的 tenantId/userId。
	// 与 file_repo.go Update/Delete 的取值范式一致。
	uid, hasUser := viewerUserIDFromContext(ctx)
	tid := uint32(0)
	hasTenant := false
	if vc, exist := viewer.FromContext(ctx); exist && vc != nil {
		tid = uint32(vc.TenantID())
		hasTenant = tid != 0
	}
	if !hasTenant || !hasUser {
		return nil, storageV1.ErrorUploadFailed("missing operator identity")
	}

	if err = s.recordFile(
		ctx,
		tid, uid,
		req.GetSourceFileName(),
		info, downloadUrl); err != nil {
		return nil, err
	}

	return &storageV1.UploadFileResponse{
		ObjectName: trans.Ptr(downloadUrl),
	}, nil
}

// presignedUploadFile 预签名上传文件
func (s *FileTransferService) presignedUploadFile(ctx context.Context, req *storageV1.UploadFileRequest) (*storageV1.UploadFileResponse, error) {
	// 预签名上传路径已禁用：该路径无法在服务端可靠记录文件元数据
	// （sourceFileName/tenantId/userId/sha256 无法从 MinIO 事件通知获得，
	// 且 x-amz-meta-* 客户端可伪造）。当前业务上传统一走 directUploadFile
	// （服务端中转，已正确落库）。待有预签名直传刚需时，需引入
	// MinIO 事件通知 + 待确认表 + 回调端点 + 定时清理的完整闭环。
	_ = req
	return nil, storageV1.ErrorUploadFailed("presigned upload is not implemented, use direct upload instead")
}

// UploadFile 上传文件
func (s *FileTransferService) UploadFile(ctx context.Context, req *storageV1.UploadFileRequest) (*storageV1.UploadFileResponse, error) {
	if req == nil {
		return nil, storageV1.ErrorUploadFailed("invalid request")
	}
	switch req.Source.(type) {
	case *storageV1.UploadFileRequest_File:
		return s.directUploadFile(ctx, req)

	case *storageV1.UploadFileRequest_Presign:
		return s.presignedUploadFile(ctx, req)

	default:
		return nil, storageV1.ErrorUploadFailed("unknown upload source")
	}
}

// downloadFileFromURL 从指定的 URL 下载文件内容
func (s *FileTransferService) downloadFileFromURL(ctx context.Context, downloadUrl string) (*storageV1.DownloadFileResponse, error) {
	parsedURL, err := netutil.ValidateURL(downloadUrl)
	if err != nil {
		return nil, storageV1.ErrorDownloadFailed(err.Error())
	}

	httpReq, err := http.NewRequestWithContext(ctx, "GET", parsedURL.String(), nil)
	if err != nil {
		return nil, storageV1.ErrorDownloadFailed(err.Error())
	}

	client := netutil.SafeHTTPClient()
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, storageV1.ErrorDownloadFailed(err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		return nil, storageV1.ErrorDownloadFailed("unexpected status: " + resp.Status)
	}

	fileData, err := io.ReadAll(netutil.LimitReader(resp.Body))
	if err != nil {
		return nil, storageV1.ErrorDownloadFailed(err.Error())
	}

	return &storageV1.DownloadFileResponse{
		Content: &storageV1.DownloadFileResponse_File{
			File: fileData,
		},
	}, nil
}

// DownloadFile 下载文件
func (s *FileTransferService) DownloadFile(ctx context.Context, req *storageV1.DownloadFileRequest) (*storageV1.DownloadFileResponse, error) {
	if req == nil {
		return nil, storageV1.ErrorDownloadFailed("invalid request")
	}
	switch req.Selector.(type) {
	case *storageV1.DownloadFileRequest_FileId:
		resp, err := s.fileRepo.Get(ctx, &storageV1.GetFileRequest{
			QueryBy: &storageV1.GetFileRequest_Id{Id: req.GetFileId()},
		})
		if err != nil {
			return nil, storageV1.ErrorDownloadFailed("file not found")
		}

		req.Selector = &storageV1.DownloadFileRequest_StorageObject{
			StorageObject: &storageV1.StorageObject{
				BucketName: resp.BucketName,
				ObjectName: trans.Ptr(resp.GetFileDirectory() + resp.GetSaveFileName()),
			},
		}

		return s.mc.DownloadFile(ctx, req)

	case *storageV1.DownloadFileRequest_StorageObject:
		return s.mc.DownloadFile(ctx, req)

	case *storageV1.DownloadFileRequest_DownloadUrl:
		return s.downloadFileFromURL(ctx, req.GetDownloadUrl())

	default:
		return nil, storageV1.ErrorDownloadFailed("unknown download selector")
	}
}
