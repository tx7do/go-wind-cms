package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
	mediaV1 "go-wind-cms/api/gen/go/media/service/v1"

	"go-wind-cms/app/core/service/internal/data"
)

type MediaAssetService struct {
	mediaV1.UnimplementedMediaAssetServiceServer

	mediaAssetRepo *data.MediaAssetRepo
	referenceRepo  *data.ContentReferenceRepo
	postRepo       *data.PostRepo
	pageRepo       *data.PageRepo
	log            *log.Helper
}

func NewMediaAssetService(ctx *bootstrap.Context, uc *data.MediaAssetRepo, referenceRepo *data.ContentReferenceRepo, postRepo *data.PostRepo, pageRepo *data.PageRepo) *MediaAssetService {
	return &MediaAssetService{
		log:            ctx.NewLoggerHelper("media-asset/service/core-service"),
		mediaAssetRepo: uc,
		referenceRepo:  referenceRepo,
		postRepo:       postRepo,
		pageRepo:       pageRepo,
	}
}

// ListDuplicatedMediaAssets 重复资源分组（同租户内同内容哈希 >1 条）。
func (s *MediaAssetService) ListDuplicatedMediaAssets(ctx context.Context, req *mediaV1.ListDuplicatedMediaAssetsRequest) (*mediaV1.ListDuplicatedMediaAssetsResponse, error) {
	groups, err := s.mediaAssetRepo.ListDuplicatedGroups(ctx, int(req.GetLimit()))
	if err != nil {
		return nil, err
	}
	return &mediaV1.ListDuplicatedMediaAssetsResponse{
		Groups: groups,
		Total:  uint64(len(groups)),
	}, nil
}

// ListMediaAssetReferences 媒体被内容引用的关系（引用图谱查询）。
// 源标题取文章/页面的首个非空翻译标题；源已删除的残留行原样返回（sourceTitle 空），
// 由前端提示"源已不存在"。
func (s *MediaAssetService) ListMediaAssetReferences(ctx context.Context, req *mediaV1.ListMediaAssetReferencesRequest) (*mediaV1.ListMediaAssetReferencesResponse, error) {
	if req == nil || req.GetMediaAssetId() == 0 {
		return nil, mediaV1.ErrorBadRequest("invalid parameter")
	}
	rows, err := s.referenceRepo.ListByTarget(ctx, "media_asset", req.GetMediaAssetId())
	if err != nil {
		return nil, mediaV1.ErrorInternalServerError("list media asset references failed")
	}

	items := make([]*mediaV1.MediaAssetReference, 0, len(rows))
	for _, row := range rows {
		item := &mediaV1.MediaAssetReference{
			MediaAssetId: trans.Ptr(req.GetMediaAssetId()),
			SourceType:   trans.Ptr(row.SourceType),
			SourceId:     trans.Ptr(row.SourceID),
			Field:        trans.Ptr(row.Field),
		}
		if row.UpdatedAt != nil {
			item.UpdatedAt = timestamppb.New(*row.UpdatedAt)
		}
		switch row.SourceType {
		case "post":
			if post, pErr := s.postRepo.Get(ctx, &contentV1.GetPostRequest{
				QueryBy: &contentV1.GetPostRequest_Id{Id: row.SourceID},
			}); pErr == nil && post != nil {
				item.SourceTitle = trans.Ptr(firstNonEmptyContentTitle(post.GetTranslations()))
			}
		case "page":
			if page, pErr := s.pageRepo.Get(ctx, &contentV1.GetPageRequest{
				QueryBy: &contentV1.GetPageRequest_Id{Id: row.SourceID},
			}); pErr == nil && page != nil {
				item.SourceTitle = trans.Ptr(firstNonEmptyContentTitle(page.GetTranslations()))
			}
		}
		items = append(items, item)
	}
	return &mediaV1.ListMediaAssetReferencesResponse{
		Items: items,
		Total: uint64(len(items)),
	}, nil
}

// firstNonEmptyContentTitle 从多语言翻译中取首个非空标题。
func firstNonEmptyContentTitle(trs any) string {
	switch t := trs.(type) {
	case []*contentV1.PostTranslation:
		for _, tr := range t {
			if tr != nil && tr.GetTitle() != "" {
				return tr.GetTitle()
			}
		}
	case []*contentV1.PageTranslation:
		for _, tr := range t {
			if tr != nil && tr.GetTitle() != "" {
				return tr.GetTitle()
			}
		}
	}
	return ""
}

func (s *MediaAssetService) List(ctx context.Context, req *paginationV1.PagingRequest) (*mediaV1.ListMediaAssetResponse, error) {
	return s.mediaAssetRepo.List(ctx, req)
}

func (s *MediaAssetService) Get(ctx context.Context, req *mediaV1.GetMediaAssetRequest) (*mediaV1.MediaAsset, error) {
	return s.mediaAssetRepo.Get(ctx, req)
}

// FindByHash 按内容哈希查重。
func (s *MediaAssetService) FindByHash(ctx context.Context, req *mediaV1.FindMediaAssetByHashRequest) (*mediaV1.FindMediaAssetByHashResponse, error) {
	if req == nil || req.GetFileHash() == "" {
		return nil, mediaV1.ErrorBadRequest("invalid parameter")
	}
	asset, err := s.mediaAssetRepo.FindByHash(ctx, req.GetFileHash())
	if err != nil {
		return nil, err
	}
	return &mediaV1.FindMediaAssetByHashResponse{Asset: asset}, nil
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
