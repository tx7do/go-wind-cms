package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"path"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/minio/minio-go/v7"

	"github.com/tx7do/go-utils/id"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	mediaV1 "go-wind-cms/api/gen/go/media/service/v1"
	storageV1 "go-wind-cms/api/gen/go/storage/service/v1"

	"go-wind-cms/pkg/middleware/auth"
	"go-wind-cms/pkg/oss"
)

type FileTransferService struct {
	adminV1.FileTransferServiceHTTPServer

	log *log.Helper

	mc *oss.MinIOClient

	fileServiceClient       storageV1.FileServiceClient
	mediaAssetServiceClient mediaV1.MediaAssetServiceClient
}

func NewFileTransferService(
	ctx *bootstrap.Context,
	mc *oss.MinIOClient,
	fileServiceClient storageV1.FileServiceClient,
	mediaAssetServiceClient mediaV1.MediaAssetServiceClient,
) *FileTransferService {
	return &FileTransferService{
		log:                     ctx.NewLoggerHelper("file-transfer/service/admin-service"),
		mc:                      mc,
		fileServiceClient:       fileServiceClient,
		mediaAssetServiceClient: mediaAssetServiceClient,
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

func (s *FileTransferService) hashFileContent(fileData []byte) string {
	sum := sha256.Sum256(fileData)          // sha256.Sum256 返回 [32]byte
	sha256Hex := hex.EncodeToString(sum[:]) // 转为十六进制字符串
	return sha256Hex
}

// recordFile 记录文件元数据到数据库
func (s *FileTransferService) recordFile(
	ctx context.Context,
	tenantID, userID uint32,
	fileData []byte,
	sourceFileName string,
	info minio.UploadInfo,
	downloadUrl string,
) (*storageV1.File, error) {

	sha256Hex := s.hashFileContent(fileData)

	dir, fileName, ext := parseKey(info.Key)
	//s.log.Debugf("Parsed file - Dir: %s, FileName: %s, Ext: %s", dir, fileName, ext)

	var err error
	var file *storageV1.File
	if file, err = s.fileServiceClient.Create(ctx, &storageV1.CreateFileRequest{
		Data: &storageV1.File{
			Provider:      trans.Ptr(storageV1.OSSProvider_MINIO),
			BucketName:    trans.Ptr(info.Bucket),
			SaveFileName:  trans.Ptr(fileName + "." + ext),
			ContentHash:   trans.Ptr(sha256Hex),
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
		return nil, err
	}
	return file, nil
}

// directUploadFile 直接上传文件
func (s *FileTransferService) directUploadFile(ctx context.Context, req *storageV1.UploadFileRequest) (*storageV1.UploadFileResponse, error) {
	if req.StorageObject == nil {
		return nil, storageV1.ErrorUploadFailed("unknown storageObject")
	}

	if req.GetFile() == nil {
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
				req.GetFile(),
				oss.GenerateFileNameTypeUUID,
			),
		)
	}

	info, _, downloadUrl, err := s.mc.UploadFile(
		ctx,
		req.GetStorageObject().GetBucketName(),
		req.GetStorageObject().GetObjectName(),
		req.GetMime(),
		req.GetFile(),
	)
	if err != nil {
		return nil, err
	}

	if _, err = s.recordFile(
		ctx,
		operator.GetTenantId(), operator.GetUserId(),
		req.GetFile(),
		req.GetSourceFileName(),
		info, downloadUrl); err != nil {
	}

	return &storageV1.UploadFileResponse{
		ObjectName: trans.Ptr(downloadUrl),
	}, err
}

// presignedUploadFile 预签名上传文件
func (s *FileTransferService) presignedUploadFile(ctx context.Context, req *storageV1.UploadFileRequest) (*storageV1.UploadFileResponse, error) {
	if req.StorageObject == nil {
		return nil, storageV1.ErrorUploadFailed("unknown storageObject")
	}

	if req.GetPresign() == nil {
		return nil, storageV1.ErrorUploadFailed("unknown presign data")
	}

	var contentType string
	if req.GetPresign().GetContentType() != "" {
		contentType = req.GetPresign().GetContentType()
	} else if req.GetMime() != "" {
		contentType = req.GetMime()
	} else {
		return nil, storageV1.ErrorUploadFailed("unknown content type for presign")
	}

	if req.GetSourceFileName() == "" {
		return nil, storageV1.ErrorUploadFailed("unknown source file name")
	}

	if req.StorageObject.BucketName == nil {
		req.StorageObject.BucketName = trans.Ptr(oss.ContentTypeToBucketName(contentType))
	}
	if req.StorageObject.ObjectName == nil {
		req.StorageObject.ObjectName = trans.Ptr(
			oss.EnsureObjectName(
				req.GetStorageObject().GetFileDirectory(),
				req.GetSourceFileName(),
				contentType,
				req.GetFile(),
				oss.GenerateFileNameTypeUUID,
			),
		)
	}

	var method storageV1.GetUploadPresignedUrlRequest_Method
	switch strings.ToLower(req.GetPresign().GetMethod()) {
	case "put":
		method = storageV1.GetUploadPresignedUrlRequest_Put
	case "post":
		method = storageV1.GetUploadPresignedUrlRequest_Post
	default:
		method = storageV1.GetUploadPresignedUrlRequest_Post
	}

	resp, err := s.mc.GetUploadPresignedUrl(
		ctx,
		&storageV1.GetUploadPresignedUrlRequest{
			ContentType:   trans.Ptr(contentType),
			ExpireSeconds: req.GetPresign().ExpireSeconds,
			Method:        method,
			BucketName:    req.StorageObject.BucketName,
			FileDirectory: req.StorageObject.FileDirectory,
		})
	if err != nil {
		return nil, err
	}

	// TODO : 记录文件元数据到数据库（待上传完成后再记录更合适？）

	return &storageV1.UploadFileResponse{
		PresignedUrl: trans.Ptr(resp.UploadUrl),
	}, nil
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
func (s *FileTransferService) downloadFileFromURL(ctx context.Context, downloadUrl string) (*storageV1.DownloadFileResponse, error) {
	if downloadUrl == "" {
		return nil, storageV1.ErrorDownloadFailed("empty download url")
	}

	// 如果需要支持断点续传，可在此构造请求并设置 Range 头
	httpReq, err := http.NewRequestWithContext(ctx, "GET", downloadUrl, nil)
	if err != nil {
		return nil, storageV1.ErrorDownloadFailed(err.Error())
	}
	// 示例：如果你要设置 Range（可选）
	// httpReq.Header.Set("Range", "bytes=100-")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, storageV1.ErrorDownloadFailed(err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		return nil, storageV1.ErrorDownloadFailed("unexpected status: " + resp.Status)
	}

	fileData, err := io.ReadAll(resp.Body)
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
	switch req.Selector.(type) {
	case *storageV1.DownloadFileRequest_FileId:
		resp, err := s.fileServiceClient.Get(ctx, &storageV1.GetFileRequest{
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

func (s *FileTransferService) UEditorUploadFile(ctx context.Context, req *storageV1.UEditorUploadRequest) (*storageV1.UEditorUploadResponse, error) {
	//s.log.Infof("上传文件： %s", req.GetFile())

	if req.File == nil {
		return nil, storageV1.ErrorUploadFailed("unknown file")
	}

	// 获取操作人信息
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.TenantId = trans.Ptr(operator.GetTenantId())
	req.UserId = trans.Ptr(operator.GetUserId())

	var bucketName string
	switch req.GetAction() {
	default:
		fallthrough
	case storageV1.UEditorAction_uploadFile.String():
		bucketName = "files"
	case storageV1.UEditorAction_uploadImage.String(), storageV1.UEditorAction_uploadScrawl.String(), storageV1.UEditorAction_catchImage.String():
		bucketName = "images"
	case storageV1.UEditorAction_uploadVideo.String():
		bucketName = "videos"
	}

	info, _, downloadUrl, err := s.mc.UploadFile(
		ctx,
		bucketName,
		req.GetSourceFileName(),
		req.GetMime(),
		req.GetFile(),
	)
	if err != nil {
		return nil, err
	}

	if _, err = s.recordFile(
		ctx,
		operator.GetTenantId(), operator.GetUserId(),
		req.GetFile(),
		req.GetSourceFileName(),
		info, downloadUrl); err != nil {
	}

	return &storageV1.UEditorUploadResponse{
		State:    trans.Ptr("SUCCESS"),
		Original: trans.Ptr(req.GetSourceFileName()),
		Title:    trans.Ptr(req.GetSourceFileName()),
		Url:      trans.Ptr(downloadUrl),
		Size:     trans.Ptr(int32(len(req.GetFile()))),
		Type:     trans.Ptr(path.Ext(req.GetSourceFileName())),
	}, nil
}

func (s *FileTransferService) mimeTypeToBucketName(mimeType string) string {
	if mimeType == "" {
		return "files"
	}

	if strings.HasPrefix(mimeType, "image/") {
		return "images"
	}

	if strings.HasPrefix(mimeType, "video/") {
		return "videos"
	}

	if strings.HasPrefix(mimeType, "audio/") {
		return "audios"
	}

	return "files"
}

func (s *FileTransferService) mimeTypeToAssetType(mimeType string) *mediaV1.MediaAsset_AssetType {
	if mimeType == "" {
		return trans.Ptr(mediaV1.MediaAsset_ASSET_TYPE_UNSPECIFIED)
	}

	mt := strings.ToLower(strings.TrimSpace(mimeType))
	if idx := strings.Index(mt, ";"); idx >= 0 {
		mt = strings.TrimSpace(mt[:idx])
	}

	if strings.HasPrefix(mt, "image/") {
		return trans.Ptr(mediaV1.MediaAsset_ASSET_TYPE_IMAGE)
	}

	if strings.HasPrefix(mt, "video/") {
		return trans.Ptr(mediaV1.MediaAsset_ASSET_TYPE_VIDEO)
	}

	if strings.HasPrefix(mt, "audio/") {
		return trans.Ptr(mediaV1.MediaAsset_ASSET_TYPE_AUDIO)
	}

	if s.isArchiveMimeType(mt) {
		return trans.Ptr(mediaV1.MediaAsset_ASSET_TYPE_ARCHIVE)
	}

	if s.isDocumentMimeType(mt) {
		return trans.Ptr(mediaV1.MediaAsset_ASSET_TYPE_DOCUMENT)
	}

	return trans.Ptr(mediaV1.MediaAsset_ASSET_TYPE_OTHER)
}

func (s *FileTransferService) isDocumentMimeType(mimeType string) bool {
	if strings.HasPrefix(mimeType, "application/pdf") {
		return true
	}

	if strings.HasPrefix(mimeType, "application/msword") ||
		strings.HasPrefix(mimeType, "application/vnd.openxmlformats-officedocument.wordprocessingml.document") {
		return true
	}

	if strings.HasPrefix(mimeType, "application/vnd.ms-excel") ||
		strings.HasPrefix(mimeType, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet") {
		return true
	}

	if strings.HasPrefix(mimeType, "application/vnd.ms-powerpoint") ||
		strings.HasPrefix(mimeType, "application/vnd.openxmlformats-officedocument.presentationml.presentation") {
		return true
	}

	return false
}

func (s *FileTransferService) isArchiveMimeType(mimeType string) bool {
	switch mimeType {
	case "application/zip",
		"application/x-zip-compressed",
		"multipart/x-zip",
		"application/x-7z-compressed",
		"application/x-tar",
		"application/gzip",
		"application/x-gzip",
		"application/x-bzip2",
		"application/x-bzip",
		"application/x-xz",
		"application/vnd.rar",
		"application/x-rar-compressed",
		"application/java-archive",
		"application/zstd",
		"application/x-zstd":
		return true
	default:
		return false
	}
}

func (s *FileTransferService) UploadMediaAsset(ctx context.Context, req *storageV1.UploadMediaAssetRequest) (*storageV1.UploadFileResponse, error) {
	if req.File == nil {
		return nil, storageV1.ErrorUploadFailed("unknown file")
	}

	// 获取操作人信息
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.TenantId = trans.Ptr(operator.GetTenantId())
	req.UserId = trans.Ptr(operator.GetUserId())

	var bucketName = s.mimeTypeToBucketName(req.GetMimeType())

	info, storagePath, downloadUrl, err := s.mc.UploadFile(
		ctx,
		bucketName,
		"",
		req.GetMimeType(),
		req.GetFile(),
	)
	if err != nil {
		return nil, err
	}

	var file *storageV1.File
	if file, err = s.recordFile(
		ctx,
		operator.GetTenantId(), operator.GetUserId(),
		req.GetFile(),
		req.GetSourceFileName(),
		info, downloadUrl,
	); err != nil {
		return nil, err
	}

	if _, err = s.mediaAssetServiceClient.Create(ctx, &mediaV1.CreateMediaAssetRequest{
		Data: &mediaV1.MediaAsset{
			FileId:           file.Id,
			AltText:          req.AltText,
			Title:            req.Title,
			Caption:          req.Caption,
			Url:              trans.Ptr(downloadUrl),
			StoragePath:      trans.Ptr(storagePath),
			Size:             trans.Ptr(uint64(info.Size)),
			MimeType:         req.MimeType,
			Filename:         req.SourceFileName,
			Type:             s.mimeTypeToAssetType(req.GetMimeType()),
			FileHash:         trans.Ptr(s.hashFileContent(req.File)),
			CreatedBy:        trans.Ptr(operator.GetUserId()),
			ProcessingStatus: trans.Ptr(mediaV1.MediaAsset_PROCESSING_STATUS_COMPLETED),
		},
	}); err != nil {
		return nil, err
	}

	return &storageV1.UploadFileResponse{
		ObjectName: trans.Ptr(downloadUrl),
	}, nil
}
