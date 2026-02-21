package service

import (
	"context"
	"go-wind-cms/pkg/middleware/auth"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	mediaV1 "go-wind-cms/api/gen/go/media/service/v1"
)

type MediaAssetService struct {
	adminV1.MediaAssetServiceHTTPServer

	mediaAssetServiceClient mediaV1.MediaAssetServiceClient
	log                     *log.Helper
}

func NewMediaAssetService(ctx *bootstrap.Context, mediaAssetServiceClient mediaV1.MediaAssetServiceClient) *MediaAssetService {
	return &MediaAssetService{
		log:                     ctx.NewLoggerHelper("media-asset/service/admin-service"),
		mediaAssetServiceClient: mediaAssetServiceClient,
	}
}

func (s *MediaAssetService) List(ctx context.Context, req *paginationV1.PagingRequest) (*mediaV1.ListMediaAssetResponse, error) {
	return s.mediaAssetServiceClient.List(ctx, req)
}

func (s *MediaAssetService) Get(ctx context.Context, req *mediaV1.GetMediaAssetRequest) (*mediaV1.MediaAsset, error) {
	return s.mediaAssetServiceClient.Get(ctx, req)
}

func (s *MediaAssetService) Create(ctx context.Context, req *mediaV1.CreateMediaAssetRequest) (*mediaV1.MediaAsset, error) {
	return s.mediaAssetServiceClient.Create(ctx, req)
}

func (s *MediaAssetService) Update(ctx context.Context, req *mediaV1.UpdateMediaAssetRequest) (*mediaV1.MediaAsset, error) {
	if req.Data == nil {
		return nil, adminV1.ErrorBadRequest("invalid parameter")
	}

	// 获取操作人信息
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.Data.UpdatedBy = trans.Ptr(operator.GetUserId())
	if req.UpdateMask != nil {
		req.UpdateMask.Paths = append(req.UpdateMask.Paths, "updated_by")
	}

	return s.mediaAssetServiceClient.Update(ctx, req)
}

func (s *MediaAssetService) Delete(ctx context.Context, req *mediaV1.DeleteMediaAssetRequest) (*emptypb.Empty, error) {
	return s.mediaAssetServiceClient.Delete(ctx, req)
}
