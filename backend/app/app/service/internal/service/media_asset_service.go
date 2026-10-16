package service

import (
	"context"

	"github.com/tx7do/kratos-bootstrap/bootstrap"

	mediaV1 "go-wind-cms/api/gen/go/media/service/v1"
)

// MediaAssetService 公开媒体资产门面（C 端渲染 customFields IMAGE 引用换 URL；
// 透传 core，租户隔离由 core viewer 保证）。
type MediaAssetService struct {
	mediaV1.UnimplementedMediaAssetServiceServer

	mediaClient mediaV1.MediaAssetServiceClient
}

func NewMediaAssetService(ctx *bootstrap.Context, mediaClient mediaV1.MediaAssetServiceClient) *MediaAssetService {
	return &MediaAssetService{
		mediaClient: mediaClient,
	}
}

func (s *MediaAssetService) Get(ctx context.Context, req *mediaV1.GetMediaAssetRequest) (*mediaV1.MediaAsset, error) {
	return s.mediaClient.Get(ctx, req)
}
