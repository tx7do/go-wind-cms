package service

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	mediaV1 "go-wind-cms/api/gen/go/media/service/v1"
	storageV1 "go-wind-cms/api/gen/go/storage/service/v1"

	"go-wind-cms/pkg/oss"
)

// MediaAssetService 公开媒体资产门面（C 端渲染 customFields IMAGE 引用换 URL；
// 透传 core，租户隔离由 core viewer 保证）。
type MediaAssetService struct {
	mediaV1.UnimplementedMediaAssetServiceServer

	mediaClient mediaV1.MediaAssetServiceClient
	mc          *oss.MinIOClient
	log         *log.Helper
}

func NewMediaAssetService(ctx *bootstrap.Context, mediaClient mediaV1.MediaAssetServiceClient, mc *oss.MinIOClient) *MediaAssetService {
	return &MediaAssetService{
		mediaClient: mediaClient,
		mc:          mc,
		log:         ctx.NewLoggerHelper("media-asset/service/app-service"),
	}
}

func (s *MediaAssetService) Get(ctx context.Context, req *mediaV1.GetMediaAssetRequest) (*mediaV1.MediaAsset, error) {
	return s.mediaClient.Get(ctx, req)
}

// ServeContent 输出媒体资产图片字节：即时变换端点
// GET /app/v1/media-assets/{id}/content?w=&h=&q=&fit= 的服务层实现。
//
// 可见性复用 core Get（匿名 viewer 按请求 Host 注入、软删资产 404、跨租户隔离），
// 拿到元数据后才经 MinIO 读字节；变换结果按参数键落派生对象缓存，
// 命中免解码直出。参数语义见 image_transform.go。
func (s *MediaAssetService) ServeContent(ctx context.Context, w http.ResponseWriter, r *http.Request, assetID uint32) error {
	p, err := ParseTransformParams(r.URL.Query())
	if err != nil {
		return mediaV1.ErrorBadRequest("invalid transform params: %v", err)
	}

	asset, err := s.mediaClient.Get(ctx, &mediaV1.GetMediaAssetRequest{Id: assetID})
	if err != nil {
		return err
	}

	mimeType := asset.GetMimeType()
	if !isTransformableImage(mimeType) {
		return mediaV1.ErrorUnsupportedMediaType("media asset #%d is not a transformable image (%s)", assetID, mimeType)
	}
	bucket, objectName, err := splitStoragePath(asset.GetStoragePath())
	if err != nil {
		return mediaV1.ErrorNotFound("media asset #%d storage missing", assetID)
	}

	cacheKey := "raw"
	if !p.IsIdentity() {
		cacheKey = p.cacheKey()
	}
	w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	// 变换请求先查派生缓存：命中则不必读原图，更不必解码
	if !p.IsIdentity() {
		derived := derivedObjectName(objectName, p)
		if data, readErr := s.readObject(ctx, bucket, derived); readErr == nil {
			serveImageBytes(w, r, fmt.Sprintf(`"img-%d-%s"`, assetID, cacheKey), derivedImageContentType, data)
			return nil
		}
	}

	src, err := s.readObject(ctx, bucket, objectName)
	if err != nil {
		return mediaV1.ErrorDownloadFailed("read media asset #%d failed", assetID)
	}

	contentType := mimeType
	if !p.IsIdentity() {
		// 只缩小：目标不小于源时直出原图（省一次解码/重编码）
		cfg, _, probeErr := image.DecodeConfig(bytes.NewReader(src))
		if probeErr != nil {
			return mediaV1.ErrorUnprocessableEntity("media asset #%d is not a decodable image", assetID)
		}
		tw, th := TargetSize(cfg.Width, cfg.Height, p)
		if tw == cfg.Width && th == cfg.Height {
			p = TransformParams{} // 回落 raw 直出
		} else {
			out, transformErr := TransformImage(src, p)
			if transformErr != nil {
				return mediaV1.ErrorUnprocessableEntity("transform media asset #%d failed: %v", assetID, transformErr)
			}
			// 派生缓存写入 best-effort：失败只记日志，不阻塞本次响应
			if _, _, _, upErr := s.mc.UploadFile(ctx, bucket, derivedObjectName(objectName, p), derivedImageContentType, bytes.NewReader(out), int64(len(out))); upErr != nil {
				s.log.Errorf("upload derived image object failed: %v", upErr)
			}
			src = out
			contentType = derivedImageContentType
			cacheKey = p.cacheKey()
		}
	}

	serveImageBytes(w, r, fmt.Sprintf(`"img-%d-%s"`, assetID, cacheKey), contentType, src)
	return nil
}

// serveImageBytes 写图片字节响应：immutable 缓存 + ETag/304 协商。
func serveImageBytes(w http.ResponseWriter, r *http.Request, etag, contentType string, data []byte) {
	w.Header().Set("ETag", etag)
	if match := r.Header.Get("If-None-Match"); match != "" && etagMatches(match, etag) {
		w.WriteHeader(http.StatusNotModified)
		return
	}
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

// etagMatches If-None-Match 弱比较：逗号分隔的标签列表，`*` 匹配任意。
func etagMatches(ifNoneMatch, etag string) bool {
	for _, candidate := range strings.Split(ifNoneMatch, ",") {
		candidate = strings.TrimSpace(candidate)
		if candidate == "*" || candidate == etag {
			return true
		}
	}
	return false
}

// readObject 读整个对象字节（图片 ≤ 上传上限 100MiB，与变体预生成同款内存策略）。
func (s *MediaAssetService) readObject(ctx context.Context, bucket, objectName string) ([]byte, error) {
	resp, err := s.mc.DownloadFile(ctx, &storageV1.DownloadFileRequest{
		Selector: &storageV1.DownloadFileRequest_StorageObject{
			StorageObject: &storageV1.StorageObject{
				BucketName: &bucket,
				ObjectName: &objectName,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return resp.GetFile(), nil
}
