package service

import (
	"context"
	"io"
	"net/http"
	"path"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/minio/minio-go/v7"

	"github.com/tx7do/go-utils/id"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	appV1 "go-wind-cms/api/gen/go/app/service/v1"
	storageV1 "go-wind-cms/api/gen/go/storage/service/v1"

	"go-wind-cms/pkg/middleware/auth"
	"go-wind-cms/pkg/netutil"
	"go-wind-cms/pkg/oss"
)

type FileTransferService struct {
	appV1.FileTransferServiceHTTPServer

	log *log.Helper

	mc                *oss.MinIOClient
	fileServiceClient storageV1.FileServiceClient
}

func NewFileTransferService(
	ctx *bootstrap.Context,
	mc *oss.MinIOClient,
	fileServiceClient storageV1.FileServiceClient,
) *FileTransferService {
	return &FileTransferService{
		log:               ctx.NewLoggerHelper("file-transfer/service/app-service"),
		mc:                mc,
		fileServiceClient: fileServiceClient,
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

	if _, err := s.fileServiceClient.Create(ctx, &storageV1.CreateFileRequest{
		Data: &storageV1.File{
			Provider:      trans.Ptr(storageV1.OSSProvider_MINIO),
			BucketName:    trans.Ptr(info.Bucket),
			SaveFileName:  trans.Ptr(fileName + "." + ext),
			FileDirectory: trans.Ptr(dir),
			FileName:      trans.Ptr(sourceFileName),
			Extension:     trans.Ptr(ext),
			FileGuid:      trans.Ptr(id.NewGUIDv7(false)),
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

	// 流式上传：reader 由 handler 通过 context 注入，不再从 proto 的 []byte 字段取。
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

	// 获取操作人信息
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
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

	if err = s.recordFile(
		ctx,
		operator.GetTenantId(), operator.GetUserId(),
		req.GetSourceFileName(),
		info, downloadUrl); err != nil {
		// 元数据写入失败，回滚已上传的对象，避免孤儿文件
		if delErr := s.mc.DeleteFile(ctx, req.GetStorageObject().GetBucketName(), req.GetStorageObject().GetObjectName()); delErr != nil {
			s.log.Errorf("cleanup orphaned object after recordFile failure failed: %s", delErr.Error())
		}
		return nil, err
	}

	return &storageV1.UploadFileResponse{
		ObjectName: trans.Ptr(downloadUrl),
	}, nil
}

// presignedUploadFile 预签名上传文件
//
// 预签名直传路径已禁用：与 core 端保持一致。该路径无法在服务端可靠记录
// 文件元数据（sourceFileName/tenantId/userId/sha256 无法从 MinIO 事件通知
// 获得，且 x-amz-meta-* 客户端可伪造），会产生不入库的孤儿对象。当前业务
// 上传统一走 directUploadFile（服务端中转，已正确落库）。待有预签名直传
// 刚需时，需引入 MinIO 事件通知 + 待确认表 + 回调端点 + 定时清理的完整闭环。
func (s *FileTransferService) presignedUploadFile(ctx context.Context, req *storageV1.UploadFileRequest) (*storageV1.UploadFileResponse, error) {
	_ = req
	return nil, storageV1.ErrorUploadFailed("presigned upload is not implemented, use direct upload instead")
}

// UploadFile 上传文件
func (s *FileTransferService) UploadFile(ctx context.Context, req *storageV1.UploadFileRequest) (*storageV1.UploadFileResponse, error) {
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
func (s *FileTransferService) downloadFileFromURL(ctx context.Context, downloadContext string) (*storageV1.DownloadFileResponse, error) {
	parsedURL, err := netutil.ValidateURL(downloadContext)
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
	// 下载须凭 fileId 走归属校验：StorageObject/DownloadUrl 两个 selector
	// 无法关联到 file 元数据，对外部调用者关闭，避免越权下载与 SSRF。
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	switch req.Selector.(type) {
	case *storageV1.DownloadFileRequest_FileId:
		resp, err := s.fileServiceClient.Get(ctx, &storageV1.GetFileRequest{
			QueryBy: &storageV1.GetFileRequest_Id{Id: req.GetFileId()},
		})
		if err != nil {
			return nil, storageV1.ErrorDownloadFailed("file not found")
		}

		// 归属校验：文件 tenant 必须等于调用者 tenant，否则拒绝。
		if resp.GetTenantId() != operator.GetTenantId() {
			return nil, storageV1.ErrorDownloadFailed("forbidden: file does not belong to caller's tenant")
		}

		req.Selector = &storageV1.DownloadFileRequest_StorageObject{
			StorageObject: &storageV1.StorageObject{
				BucketName: resp.BucketName,
				ObjectName: trans.Ptr(resp.GetFileDirectory() + resp.GetSaveFileName()),
			},
		}

		return s.mc.DownloadFile(ctx, req)

	case *storageV1.DownloadFileRequest_StorageObject:
		return nil, storageV1.ErrorDownloadFailed("storageObject selector is not allowed for external callers, use fileId")

	case *storageV1.DownloadFileRequest_DownloadUrl:
		return nil, storageV1.ErrorDownloadFailed("downloadUrl selector is disabled")

	default:
		return nil, storageV1.ErrorDownloadFailed("unknown download selector")
	}
}
