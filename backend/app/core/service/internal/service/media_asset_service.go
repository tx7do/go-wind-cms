package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	"go-wind-cms/app/core/service/internal/data"

	mediaV1 "go-wind-cms/api/gen/go/media/service/v1"
)

type MediaAssetService struct {
	mediaV1.UnimplementedMediaAssetServiceServer

	mediaAssetRepo *data.MediaAssetRepo
	log            *log.Helper
}

func NewMediaAssetService(ctx *bootstrap.Context, uc *data.MediaAssetRepo) *MediaAssetService {
	return &MediaAssetService{
		log:            ctx.NewLoggerHelper("media-asset/service/core-service"),
		mediaAssetRepo: uc,
	}
}

func (s *MediaAssetService) List(ctx context.Context, req *paginationV1.PagingRequest) (*mediaV1.ListMediaAssetResponse, error) {
	return s.mediaAssetRepo.List(ctx, req)
}

func (s *MediaAssetService) Get(ctx context.Context, req *mediaV1.GetMediaAssetRequest) (*mediaV1.MediaAsset, error) {
	return s.mediaAssetRepo.Get(ctx, req)
}

func (s *MediaAssetService) Create(ctx context.Context, req *mediaV1.CreateMediaAssetRequest) (*mediaV1.MediaAsset, error) {
	return s.mediaAssetRepo.Create(ctx, req)
}

func (s *MediaAssetService) Update(ctx context.Context, req *mediaV1.UpdateMediaAssetRequest) (*mediaV1.MediaAsset, error) {
	return s.mediaAssetRepo.Update(ctx, req)
}

func (s *MediaAssetService) Delete(ctx context.Context, req *mediaV1.DeleteMediaAssetRequest) (*emptypb.Empty, error) {
	_, err := s.mediaAssetRepo.Delete(ctx, req)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
