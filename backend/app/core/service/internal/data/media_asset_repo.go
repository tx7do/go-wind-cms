package data

import (
	"context"

	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	entCrud "github.com/tx7do/go-crud/entgo"
	"github.com/tx7do/go-utils/copierutil"
	"github.com/tx7do/go-utils/mapper"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/mediaasset"
	"go-wind-cms/app/core/service/internal/data/ent/predicate"

	mediaV1 "go-wind-cms/api/gen/go/media/service/v1"
)

type MediaAssetRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper

	mapper *mapper.CopierMapper[mediaV1.MediaAsset, ent.MediaAsset]

	repository *entCrud.Repository[
		ent.MediaAssetQuery, ent.MediaAssetSelect,
		ent.MediaAssetCreate, ent.MediaAssetCreateBulk,
		ent.MediaAssetUpdate, ent.MediaAssetUpdateOne,
		ent.MediaAssetDelete,
		predicate.MediaAsset,
		mediaV1.MediaAsset, ent.MediaAsset,
	]

	assetTypeConverter        *mapper.EnumTypeConverter[mediaV1.MediaAsset_AssetType, mediaasset.Type]
	processingStatusConverter *mapper.EnumTypeConverter[mediaV1.MediaAsset_ProcessingStatus, mediaasset.ProcessingStatus]

	mediaVariantRepo *MediaVariantRepo
}

func NewMediaAssetRepo(
	ctx *bootstrap.Context,
	entClient *entCrud.EntClient[*ent.Client],
	mediaVariantRepo *MediaVariantRepo,
) *MediaAssetRepo {
	repo := &MediaAssetRepo{
		entClient: entClient,
		log:       ctx.NewLoggerHelper("media-asset/repo/core-service"),
		mapper:    mapper.NewCopierMapper[mediaV1.MediaAsset, ent.MediaAsset](),
		assetTypeConverter: mapper.NewEnumTypeConverter[mediaV1.MediaAsset_AssetType, mediaasset.Type](
			mediaV1.MediaAsset_AssetType_name, mediaV1.MediaAsset_AssetType_value,
		),
		processingStatusConverter: mapper.NewEnumTypeConverter[mediaV1.MediaAsset_ProcessingStatus, mediaasset.ProcessingStatus](
			mediaV1.MediaAsset_ProcessingStatus_name, mediaV1.MediaAsset_ProcessingStatus_value,
		),
		mediaVariantRepo: mediaVariantRepo,
	}

	repo.init()

	return repo
}

func (r *MediaAssetRepo) init() {
	r.repository = entCrud.NewRepository[
		ent.MediaAssetQuery, ent.MediaAssetSelect,
		ent.MediaAssetCreate, ent.MediaAssetCreateBulk,
		ent.MediaAssetUpdate, ent.MediaAssetUpdateOne,
		ent.MediaAssetDelete,
		predicate.MediaAsset,
		mediaV1.MediaAsset, ent.MediaAsset,
	](r.mapper)

	r.mapper.AppendConverters(copierutil.NewTimeStringConverterPair())
	r.mapper.AppendConverters(copierutil.NewTimeTimestamppbConverterPair())

	r.mapper.AppendConverters(r.assetTypeConverter.NewConverterPair())
	r.mapper.AppendConverters(r.processingStatusConverter.NewConverterPair())
}

func (r *MediaAssetRepo) IsExist(ctx context.Context, id uint32) (bool, error) {
	exist, err := r.entClient.Client().MediaAsset.Query().
		Where(mediaasset.IDEQ(id)).
		Exist(ctx)
	if err != nil {
		r.log.Errorf("query media asset exist failed: %s", err.Error())
		return false, mediaV1.ErrorInternalServerError("query media asset exist failed")
	}
	return exist, nil
}

func (r *MediaAssetRepo) count(ctx context.Context, whereCond []func(s *sql.Selector)) (int, error) {
	builder := r.entClient.Client().MediaAsset.Query()
	if len(whereCond) != 0 {
		builder.Modify(whereCond...)
	}

	count, err := builder.Count(ctx)
	if err != nil {
		r.log.Errorf("query count failed: %s", err.Error())
	}

	return count, err
}

func (r *MediaAssetRepo) List(ctx context.Context, req *paginationV1.PagingRequest) (*mediaV1.ListMediaAssetResponse, error) {
	if req == nil {
		return nil, mediaV1.ErrorBadRequest("invalid parameter")
	}

	builder := r.entClient.Client().MediaAsset.Query()

	ret, err := r.repository.ListWithPaging(ctx, builder, builder.Clone(), req)
	if err != nil {
		return nil, err
	}
	if ret == nil {
		return &mediaV1.ListMediaAssetResponse{Total: 0, Items: nil}, nil
	}

	return &mediaV1.ListMediaAssetResponse{
		Total: ret.Total,
		Items: ret.Items,
	}, nil
}

func (r *MediaAssetRepo) Get(ctx context.Context, req *mediaV1.GetMediaAssetRequest) (*mediaV1.MediaAsset, error) {
	if req == nil {
		return nil, mediaV1.ErrorBadRequest("invalid parameter")
	}

	entity, err := r.entClient.Client().MediaAsset.Get(ctx, req.GetId())
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, mediaV1.ErrorFileNotFound("media asset not found")
		}

		r.log.Errorf("query media asset failed: %s", err.Error())

		return nil, mediaV1.ErrorInternalServerError("query media asset failed")
	}

	return r.mapper.ToDTO(entity), nil
}

func (r *MediaAssetRepo) Create(ctx context.Context, req *mediaV1.CreateMediaAssetRequest) (*mediaV1.MediaAsset, error) {
	if req == nil || req.Data == nil {
		return nil, mediaV1.ErrorBadRequest("invalid parameter")
	}

	builder := r.entClient.Client().MediaAsset.Create().
		SetNillableType(r.assetTypeConverter.ToEntity(req.Data.Type)).
		SetNillableFilename(req.Data.Filename).
		SetNillableMimeType(req.Data.MimeType).
		SetNillableSize(req.Data.Size).
		SetNillableStoragePath(req.Data.StoragePath).
		SetNillableURL(req.Data.Url).
		SetNillableWidth(req.Data.Width).
		SetNillableHeight(req.Data.Height).
		SetNillableDuration(req.Data.Duration).
		SetNillableAltText(req.Data.AltText).
		SetNillableTitle(req.Data.Title).
		SetNillableCaption(req.Data.Caption).
		SetNillableProcessingStatus(r.processingStatusConverter.ToEntity(req.Data.ProcessingStatus)).
		SetNillableProcessingError(req.Data.ProcessingError).
		SetNillableFileHash(req.Data.FileHash).
		SetNillableReferenceCount(req.Data.ReferenceCount).
		SetNillableIsPrivate(req.Data.IsPrivate).
		SetNillableFileID(req.Data.FileId).
		SetNillableFolderID(req.Data.FolderId).
		SetNillableCreatedBy(req.Data.CreatedBy).
		SetCreatedAt(time.Now())

	var err error
	var entity *ent.MediaAsset
	if entity, err = builder.Save(ctx); err != nil {
		r.log.Errorf("insert media asset failed: %s", err.Error())
		return nil, mediaV1.ErrorInternalServerError("insert media asset failed")
	}

	return r.mapper.ToDTO(entity), nil
}

func (r *MediaAssetRepo) Update(ctx context.Context, req *mediaV1.UpdateMediaAssetRequest) (*mediaV1.MediaAsset, error) {
	if req == nil || req.Data == nil {
		return nil, mediaV1.ErrorBadRequest("invalid parameter")
	}

	// 如果不存在则创建
	if req.GetAllowMissing() {
		exist, err := r.IsExist(ctx, req.GetId())
		if err != nil {
			return nil, err
		}
		if !exist {
			req.Data.CreatedBy = req.Data.UpdatedBy
			req.Data.UpdatedBy = nil
			_, err = r.Create(ctx, &mediaV1.CreateMediaAssetRequest{Data: req.Data})
			return nil, err
		}
	}

	builder := r.entClient.Client().MediaAsset.UpdateOneID(req.GetId())
	result, err := r.repository.UpdateOne(ctx, builder, req.Data, req.GetUpdateMask(),
		func(dto *mediaV1.MediaAsset) {
			builder.
				SetNillableType(r.assetTypeConverter.ToEntity(req.Data.Type)).
				SetNillableFilename(req.Data.Filename).
				SetNillableMimeType(req.Data.MimeType).
				SetNillableSize(req.Data.Size).
				SetNillableStoragePath(req.Data.StoragePath).
				SetNillableURL(req.Data.Url).
				SetNillableWidth(req.Data.Width).
				SetNillableHeight(req.Data.Height).
				SetNillableDuration(req.Data.Duration).
				SetNillableAltText(req.Data.AltText).
				SetNillableTitle(req.Data.Title).
				SetNillableCaption(req.Data.Caption).
				SetNillableProcessingStatus(r.processingStatusConverter.ToEntity(req.Data.ProcessingStatus)).
				SetNillableProcessingError(req.Data.ProcessingError).
				SetNillableFileHash(req.Data.FileHash).
				SetNillableReferenceCount(req.Data.ReferenceCount).
				SetNillableIsPrivate(req.Data.IsPrivate).
				SetNillableFileID(req.Data.FileId).
				SetNillableFolderID(req.Data.FolderId).
				SetNillableUpdatedBy(req.Data.UpdatedBy).
				SetUpdatedAt(time.Now())
		},
		func(s *sql.Selector) {
			s.Where(sql.EQ(mediaasset.FieldID, req.GetId()))
		},
	)

	return result, err
}

func (r *MediaAssetRepo) Delete(ctx context.Context, req *mediaV1.DeleteMediaAssetRequest) (bool, error) {
	err := r.entClient.Client().MediaAsset.
		DeleteOneID(req.GetId()).
		Exec(ctx)
	if err != nil {
		r.log.Errorf("delete one data failed: %s", err.Error())
	}

	return err == nil, err
}
