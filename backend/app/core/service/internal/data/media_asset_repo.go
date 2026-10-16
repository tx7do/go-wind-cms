package data

import (
	"context"
	"sort"

	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	entCrud "github.com/tx7do/go-crud/entgo"
	"github.com/tx7do/go-utils/copierutil"
	"github.com/tx7do/go-utils/mapper"
	"github.com/tx7do/go-utils/trans"
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

// FindByHash 按内容哈希查重：同租户内第一条未删除资源；无则返回空 asset。
// hash 服务端计算，同租户外不可枚举，仅内部可信链路调用。
func (r *MediaAssetRepo) FindByHash(ctx context.Context, fileHash string) (*mediaV1.MediaAsset, error) {
	conds := []predicate.MediaAsset{}
	if tid, hasTenant := maybeTenantFromViewer(ctx); hasTenant {
		conds = append(conds, mediaasset.TenantIDEQ(tid))
	}
	conds = append(conds, mediaasset.FileHashEQ(fileHash))

	row, err := r.entClient.Client().MediaAsset.Query().Where(conds...).First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return &mediaV1.MediaAsset{}, nil
		}
		r.log.Errorf("find media asset by hash failed: %s", err.Error())
		return nil, mediaV1.ErrorInternalServerError("find media asset by hash failed")
	}
	return r.mapper.ToDTO(row), nil
}

// ListDuplicatedGroups 同租户内同内容哈希 >1 条的重复分组。
// 组内按创建时间升序（首个为最早上传，应保留者）；组按条数降序。
func (r *MediaAssetRepo) ListDuplicatedGroups(ctx context.Context, limit int) ([]*mediaV1.MediaAssetDuplicateGroup, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}

	tenantConds := []predicate.MediaAsset{}
	if tid, hasTenant := maybeTenantFromViewer(ctx); hasTenant {
		tenantConds = append(tenantConds, mediaasset.TenantIDEQ(tid))
	}

	// 1) 找出出现 >1 次的哈希
	type hashCount struct {
		FileHash string `json:"file_hash"`
		Count    int    `json:"count"`
	}
	var counts []hashCount
	err := r.entClient.Client().MediaAsset.Query().
		Where(append(tenantConds,
			mediaasset.FileHashNEQ(""),
			mediaasset.FileHashNotNil(),
		)...).
		GroupBy(mediaasset.FieldFileHash).
		Aggregate(ent.Count()).
		Scan(ctx, &counts)
	if err != nil {
		r.log.Errorf("group media assets by hash failed: %s", err.Error())
		return nil, mediaV1.ErrorInternalServerError("list duplicated media assets failed")
	}

	dupHashes := make([]string, 0, len(counts))
	countByHash := make(map[string]int, len(counts))
	for _, c := range counts {
		if c.Count > 1 {
			dupHashes = append(dupHashes, c.FileHash)
			countByHash[c.FileHash] = c.Count
		}
	}
	if len(dupHashes) == 0 {
		return []*mediaV1.MediaAssetDuplicateGroup{}, nil
	}
	if len(dupHashes) > limit {
		// 组按条数降序取前 limit 组
		sort.Slice(dupHashes, func(i, j int) bool {
			return countByHash[dupHashes[i]] > countByHash[dupHashes[j]]
		})
		dupHashes = dupHashes[:limit]
	}

	// 2) 取组内资源（创建时间升序）
	rows, err := r.entClient.Client().MediaAsset.Query().
		Where(append(tenantConds, mediaasset.FileHashIn(dupHashes...))...).
		Order(ent.Asc(mediaasset.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		r.log.Errorf("query duplicated media assets failed: %s", err.Error())
		return nil, mediaV1.ErrorInternalServerError("list duplicated media assets failed")
	}

	groups := make([]*mediaV1.MediaAssetDuplicateGroup, 0, len(dupHashes))
	indexByHash := make(map[string]int, len(dupHashes))
	for _, row := range rows {
		if row.FileHash == nil {
			continue
		}
		hash := *row.FileHash
		idx, ok := indexByHash[hash]
		if !ok {
			groups = append(groups, &mediaV1.MediaAssetDuplicateGroup{
				FileHash: trans.Ptr(hash),
				Count:    trans.Ptr(uint32(countByHash[hash])),
			})
			idx = len(groups) - 1
			indexByHash[hash] = idx
		}
		groups[idx].Assets = append(groups[idx].Assets, r.mapper.ToDTO(row))
	}
	return groups, nil
}

// FindByPathCandidates 按存储路径候选（storage_path 精确 / url 精确）批量找资源。
// 返回 (storagePath→id) 与 (url→id) 两张映射，供引用抽取把正文 URL 解析成资产 ID。
func (r *MediaAssetRepo) FindByPathCandidates(ctx context.Context, paths, urls []string) (map[string]uint32, map[string]uint32, error) {
	pathByID := make(map[string]uint32, len(paths))
	urlByID := make(map[string]uint32, len(urls))
	if len(paths) == 0 && len(urls) == 0 {
		return pathByID, urlByID, nil
	}

	conds := []predicate.MediaAsset{}
	if tid, hasTenant := maybeTenantFromViewer(ctx); hasTenant {
		conds = append(conds, mediaasset.TenantIDEQ(tid))
	}
	if len(paths) > 0 && len(urls) > 0 {
		conds = append(conds, mediaasset.Or(mediaasset.StoragePathIn(paths...), mediaasset.URLIn(urls...)))
	} else if len(paths) > 0 {
		conds = append(conds, mediaasset.StoragePathIn(paths...))
	} else {
		conds = append(conds, mediaasset.URLIn(urls...))
	}

	rows, err := r.entClient.Client().MediaAsset.Query().Where(conds...).All(ctx)
	if err != nil {
		r.log.Errorf("find media assets by path candidates failed: %s", err.Error())
		return nil, nil, mediaV1.ErrorInternalServerError("find media assets by path failed")
	}
	for _, row := range rows {
		if row.StoragePath != nil {
			pathByID[*row.StoragePath] = row.ID
		}
		if row.URL != nil {
			urlByID[*row.URL] = row.ID
		}
	}
	return pathByID, urlByID, nil
}

// UpdateReferenceCounts 按给定计数回写 reference_count（引用追踪维护）。
func (r *MediaAssetRepo) UpdateReferenceCounts(ctx context.Context, counts map[uint32]int) error {
	for id, n := range counts {
		if _, err := r.entClient.Client().MediaAsset.UpdateOneID(id).
			SetReferenceCount(uint32(n)).
			Save(ctx); err != nil {
			if ent.IsNotFound(err) {
				continue
			}
			r.log.Errorf("update reference count for asset [%d] failed: %s", id, err.Error())
			return err
		}
	}
	return nil
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
		SetReferenceCount(0).
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

	tid, hasTenant := maybeTenantFromViewer(ctx)
	callerUserID, hasUser := viewerUserIDFromContext(ctx)
	builder := r.entClient.Client().MediaAsset.UpdateOneID(req.GetId())
	builder.Where(mediaasset.IDEQ(req.GetId()))
	if hasTenant {
		builder.Where(mediaasset.TenantIDEQ(tid))
	}
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
				SetNillableIsPrivate(req.Data.IsPrivate).
				SetNillableFileID(req.Data.FileId).
				SetNillableFolderID(req.Data.FolderId).
				SetUpdatedAt(time.Now())

			// updated_by 强制由服务端 viewer context 推导，忽略客户端传入值
			if hasUser {
				builder.SetUpdatedBy(callerUserID)
			}
		},
		func(s *sql.Selector) {
			s.Where(sql.EQ(mediaasset.FieldID, req.GetId()))
		},
	)

	return result, err
}

func (r *MediaAssetRepo) Delete(ctx context.Context, req *mediaV1.DeleteMediaAssetRequest) (bool, error) {
	tid, hasTenant := maybeTenantFromViewer(ctx)
	delBuilder := r.entClient.Client().MediaAsset.Delete()
	delBuilder.Where(mediaasset.IDEQ(req.GetId()))
	if hasTenant {
		delBuilder.Where(mediaasset.TenantIDEQ(tid))
	}
	_, err := delBuilder.Exec(ctx)
	if err != nil {
		r.log.Errorf("delete one data failed: %s", err.Error())
	}

	return err == nil, err
}
