package data

import (
	"context"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	entCrud "github.com/tx7do/go-crud/entgo"
	"github.com/tx7do/go-crud/pagination"
	paginationFilter "github.com/tx7do/go-crud/pagination/filter"

	"github.com/tx7do/go-utils/copierutil"
	"github.com/tx7do/go-utils/mapper"
	"github.com/tx7do/go-utils/trans"

	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/predicate"
	"go-wind-cms/app/core/service/internal/data/ent/tag"

	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
)

type TagRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper

	mapper *mapper.CopierMapper[contentV1.Tag, ent.Tag]

	repository *entCrud.Repository[
		ent.TagQuery, ent.TagSelect,
		ent.TagCreate, ent.TagCreateBulk,
		ent.TagUpdate, ent.TagUpdateOne,
		ent.TagDelete,
		predicate.Tag,
		contentV1.Tag, ent.Tag,
	]

	statusConverter *mapper.EnumTypeConverter[contentV1.Tag_TagStatus, tag.Status]

	tagTranslationRepo *TagTranslationRepo
}

func NewTagRepo(
	ctx *bootstrap.Context,
	entClient *entCrud.EntClient[*ent.Client],
	tagTranslationRepo *TagTranslationRepo,
) *TagRepo {
	repo := &TagRepo{
		entClient: entClient,
		log:       ctx.NewLoggerHelper("tag/repo/core-service"),
		mapper:    mapper.NewCopierMapper[contentV1.Tag, ent.Tag](),
		statusConverter: mapper.NewEnumTypeConverter[contentV1.Tag_TagStatus, tag.Status](
			contentV1.Tag_TagStatus_name, contentV1.Tag_TagStatus_value,
		),
		tagTranslationRepo: tagTranslationRepo,
	}

	repo.init()

	return repo
}

func (r *TagRepo) init() {
	r.repository = entCrud.NewRepository[
		ent.TagQuery, ent.TagSelect,
		ent.TagCreate, ent.TagCreateBulk,
		ent.TagUpdate, ent.TagUpdateOne,
		ent.TagDelete,
		predicate.Tag,
		contentV1.Tag, ent.Tag,
	](r.mapper)

	r.mapper.AppendConverters(copierutil.NewTimeStringConverterPair())
	r.mapper.AppendConverters(copierutil.NewTimeTimestamppbConverterPair())

	r.mapper.AppendConverters(r.statusConverter.NewConverterPair())
}

func (r *TagRepo) count(ctx context.Context, whereCond []func(s *sql.Selector)) (int, error) {
	builder := r.entClient.Client().Tag.Query()
	if len(whereCond) != 0 {
		builder.Modify(whereCond...)
	}

	count, err := builder.Count(ctx)
	if err != nil {
		r.log.Errorf("query count failed: %s", err.Error())
	}

	return count, err
}

func (r *TagRepo) IsExist(ctx context.Context, id uint32) (bool, error) {
	exist, err := r.entClient.Client().Tag.Query().
		Where(tag.IDEQ(id)).
		Exist(ctx)
	if err != nil {
		r.log.Errorf("query tag exist failed: %s", err.Error())
		return false, contentV1.ErrorInternalServerError("query tag exist failed")
	}
	return exist, nil
}

func (r *TagRepo) prepareTranslationMaskFields(req *paginationV1.PagingRequest) (
	excludeConditions []*paginationV1.FilterCondition,
	translationMaskFields []string,
	needQueryTranslation bool,
	err error,
) {
	var filterExpr *paginationV1.FilterExpr
	filterExpr, err = paginationFilter.ConvertFilterByPagingRequest(req)
	if err != nil {
		r.log.Errorf("convert filter by paging request failed: %s", err.Error())
		return
	}

	excludeFields := []string{"translations", "available_languages"}
	if req.FieldMask != nil && len(req.FieldMask.Paths) > 0 {
		for _, path := range req.FieldMask.Paths {
			path = strings.TrimSpace(path)

			if path == "translations" || strings.HasPrefix(path, "translations.") {
				needQueryTranslation = true
			}
			if strings.HasPrefix(path, "translations.") {
				excludeFields = append(excludeFields, path)
				path = strings.TrimPrefix(path, "translations.")
				if len(path) > 0 {
					translationMaskFields = append(translationMaskFields, path)
				}
			}
		}

		req.FieldMask = FilterViewMask(excludeFields, req.FieldMask)

	} else {
		needQueryTranslation = true
	}

	excludeConditions = pagination.FilterFields(filterExpr, []string{
		"locale",
	})
	req.FilteringType = &paginationV1.PagingRequest_FilterExpr{FilterExpr: filterExpr}

	return
}

func (r *TagRepo) List(ctx context.Context, req *paginationV1.PagingRequest) (*contentV1.ListTagResponse, error) {
	if req == nil {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}

	builder := r.entClient.Client().Tag.Query()

	excludeConditions, translationMaskFields, needQueryTranslation, err := r.prepareTranslationMaskFields(req)
	if err != nil {
		return nil, err
	}

	ret, err := r.repository.ListWithPaging(ctx, builder, builder.Clone(), req)
	if err != nil {
		return nil, err
	}
	if ret == nil {
		return &contentV1.ListTagResponse{Total: 0, Items: nil}, nil
	}

	if needQueryTranslation {
		viewMask := &fieldmaskpb.FieldMask{
			Paths: translationMaskFields,
		}

		var locale string
		if len(excludeConditions) > 0 {
			for _, cond := range excludeConditions {
				if cond.GetField() == "locale" {
					locale = cond.GetValue()
					break
				}
			}
		}

		for _, item := range ret.Items {
			translations, err := r.tagTranslationRepo.ListTranslations(ctx, item.GetId(), locale, viewMask)
			if err != nil {
				r.log.Errorf("query translations failed: %s", err.Error())
				return nil, contentV1.ErrorInternalServerError("query translations failed")
			}
			item.Translations = translations
		}
	}

	for _, item := range ret.Items {
		languages, err := r.tagTranslationRepo.ListAvailedLanguages(ctx, item.GetId())
		if err != nil {
			r.log.Errorf("query availed languages failed: %s", err.Error())
			return nil, contentV1.ErrorInternalServerError("query availed languages failed")
		}
		item.AvailableLanguages = languages
	}

	return &contentV1.ListTagResponse{
		Total: ret.Total,
		Items: ret.Items,
	}, nil
}

func (r *TagRepo) Get(ctx context.Context, req *contentV1.GetTagRequest) (*contentV1.Tag, error) {
	if req == nil {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}

	builder := r.entClient.Client().Tag.Query()

	switch req.QueryBy.(type) {
	case *contentV1.GetTagRequest_Id:
		builder.Where(tag.IDEQ(req.GetId()))

	case *contentV1.GetTagRequest_Code:
		builder.Where(tag.CodeEQ(req.GetCode()))

	default:
		return nil, contentV1.ErrorBadRequest("invalid query field")
	}

	entity, err := builder.Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, contentV1.ErrorFileNotFound("tag not found")
		}

		r.log.Errorf("query tag failed: %s", err.Error())

		return nil, contentV1.ErrorInternalServerError("query tag failed")
	}

	dto := r.mapper.ToDTO(entity)

	translations, err := r.tagTranslationRepo.ListTranslations(ctx, dto.GetId(), "", nil)
	if err != nil {
		r.log.Errorf("query translations failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("query translations failed")
	}
	dto.Translations = translations

	languages, err := r.tagTranslationRepo.ListAvailedLanguages(ctx, dto.GetId())
	if err != nil {
		r.log.Errorf("query availed languages failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("query availed languages failed")
	}
	dto.AvailableLanguages = languages

	return dto, nil
}

func (r *TagRepo) Create(ctx context.Context, req *contentV1.CreateTagRequest) (dto *contentV1.Tag, err error) {
	if req == nil || req.Data == nil {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}

	var tx *ent.Tx
	tx, err = r.entClient.Client().Tx(ctx)
	if err != nil {
		r.log.Errorf("start transaction failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("start transaction failed")
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				r.log.Errorf("transaction rollback failed: %s", rollbackErr.Error())
			}
			return
		}
		if commitErr := tx.Commit(); commitErr != nil {
			r.log.Errorf("transaction commit failed: %s", commitErr.Error())
			err = contentV1.ErrorInternalServerError("transaction commit failed")
		}
	}()

	builder := tx.Tag.Create().
		SetNillableStatus(r.statusConverter.ToEntity(req.Data.Status)).
		SetNillableColor(req.Data.Color).
		SetNillableIcon(req.Data.Icon).
		SetNillableGroup(req.Data.Group).
		SetNillableCode(req.Data.Code).
		SetNillableSortOrder(req.Data.SortOrder).
		SetNillableIsFeatured(req.Data.IsFeatured).
		SetPostCount(0).
		SetNillableCreatedBy(req.Data.CreatedBy).
		SetCreatedAt(time.Now())

	var entity *ent.Tag
	if entity, err = builder.Save(ctx); err != nil {
		r.log.Errorf("insert tag failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("insert tag failed")
	}

	if len(req.Data.Translations) > 0 {
		if err = r.tagTranslationRepo.CleanTranslations(ctx, tx, entity.ID); err != nil {
			r.log.Errorf("clean translations failed: %s", err.Error())
			return nil, contentV1.ErrorInternalServerError("clean translations failed")
		}

		for i := range req.Data.Translations {
			req.Data.Translations[i].TagId = trans.Ptr(entity.ID)
		}

		if err = r.tagTranslationRepo.BatchCreate(ctx, tx, req.Data.GetTranslations()); err != nil {
			r.log.Errorf("batch insert translations failed: %s", err.Error())
			return nil, contentV1.ErrorInternalServerError("batch insert translations failed")
		}
	}

	return r.mapper.ToDTO(entity), nil
}

func (r *TagRepo) Update(ctx context.Context, req *contentV1.UpdateTagRequest) (dto *contentV1.Tag, err error) {
	if req == nil || req.Data == nil {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
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
			_, err = r.Create(ctx, &contentV1.CreateTagRequest{Data: req.Data})
			return nil, err
		}
	}

	var tx *ent.Tx
	tx, err = r.entClient.Client().Tx(ctx)
	if err != nil {
		r.log.Errorf("start transaction failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("start transaction failed")
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				r.log.Errorf("transaction rollback failed: %s", rollbackErr.Error())
			}
			return
		}
		if commitErr := tx.Commit(); commitErr != nil {
			r.log.Errorf("transaction commit failed: %s", commitErr.Error())
			err = contentV1.ErrorInternalServerError("transaction commit failed")
		}
	}()

	if len(req.Data.Translations) > 0 {
		//if err = r.tagTranslationRepo.CleanTranslations(ctx, tx, req.GetId()); err != nil {
		//	r.log.Errorf("clean translations failed: %s", err.Error())
		//	return nil, contentV1.ErrorInternalServerError("clean translations failed")
		//}

		for i := range req.Data.Translations {
			req.Data.Translations[i].TagId = trans.Ptr(req.GetId())
		}

		if err = r.tagTranslationRepo.BatchCreate(ctx, tx, req.Data.GetTranslations()); err != nil {
			r.log.Errorf("batch insert translations failed: %s", err.Error())
			return nil, contentV1.ErrorInternalServerError("batch insert translations failed")
		}
	}

	builder := tx.Tag.UpdateOneID(req.GetId())
	result, err := r.repository.UpdateOne(ctx, builder, req.Data, req.GetUpdateMask(),
		func(dto *contentV1.Tag) {
			builder.
				SetNillableStatus(r.statusConverter.ToEntity(req.Data.Status)).
				SetNillableColor(req.Data.Color).
				SetNillableIcon(req.Data.Icon).
				SetNillableGroup(req.Data.Group).
				SetNillableCode(req.Data.Code).
				SetNillableSortOrder(req.Data.SortOrder).
				SetNillableIsFeatured(req.Data.IsFeatured).
				SetNillableUpdatedBy(req.Data.UpdatedBy).
				SetUpdatedAt(time.Now())
		},
		func(s *sql.Selector) {
			s.Where(sql.EQ(tag.FieldID, req.GetId()))
		},
	)

	return result, err
}

func (r *TagRepo) Delete(ctx context.Context, req *contentV1.DeleteTagRequest) (err error) {
	if req == nil {
		return contentV1.ErrorBadRequest("invalid parameter")
	}

	var tx *ent.Tx
	tx, err = r.entClient.Client().Tx(ctx)
	if err != nil {
		r.log.Errorf("start transaction failed: %s", err.Error())
		return contentV1.ErrorInternalServerError("start transaction failed")
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				r.log.Errorf("transaction rollback failed: %s", rollbackErr.Error())
			}
			return
		}
		if commitErr := tx.Commit(); commitErr != nil {
			r.log.Errorf("transaction commit failed: %s", commitErr.Error())
			err = contentV1.ErrorInternalServerError("transaction commit failed")
		}
	}()

	if err = tx.Tag.
		DeleteOneID(req.GetId()).
		Exec(ctx); err != nil {
		r.log.Errorf("delete one data failed: %s", err.Error())
		return contentV1.ErrorInternalServerError("delete one data failed")
	}

	if err = r.tagTranslationRepo.CleanTranslations(ctx, tx, req.GetId()); err != nil {
		r.log.Errorf("clean translations failed: %s", err.Error())
		return contentV1.ErrorInternalServerError("clean translations failed")
	}

	return nil
}

func (r *TagRepo) TranslationExists(ctx context.Context, tagId uint32, languageCode string) (bool, error) {
	return r.tagTranslationRepo.TranslationExists(ctx, tagId, languageCode)
}

func (r *TagRepo) CreateTranslation(ctx context.Context, req *contentV1.CreateTagTranslationRequest) (*contentV1.TagTranslation, error) {
	if req == nil || req.Data == nil {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}

	if len(req.Data.GetLanguageCode()) == 0 {
		return nil, contentV1.ErrorBadRequest("language code is required")
	}

	if req.GetTagId() == 0 {
		return nil, contentV1.ErrorBadRequest("post id is required")
	}

	req.Data.TagId = trans.Ptr(req.GetTagId())

	return r.tagTranslationRepo.CreateTranslation(ctx, req.Data)
}

func (r *TagRepo) UpdateTranslation(ctx context.Context, req *contentV1.UpdateTagTranslationRequest) (*contentV1.TagTranslation, error) {
	if req == nil || req.Data == nil {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}

	if len(req.Data.GetLanguageCode()) == 0 {
		return nil, contentV1.ErrorBadRequest("language code is required")
	}

	if req.Data.GetTagId() == 0 {
		return nil, contentV1.ErrorBadRequest("post id is required")
	}

	if exist, err := r.TranslationExists(ctx, req.Data.GetTagId(), req.Data.GetLanguageCode()); err != nil {
		return nil, err
	} else if !exist {
		if req.GetAllowMissing() {
			return r.CreateTranslation(ctx, &contentV1.CreateTagTranslationRequest{
				Data:  req.Data,
				TagId: req.Data.GetTagId(),
			})
		}

		return nil, contentV1.ErrorFileNotFound("translation not found")
	}

	return r.tagTranslationRepo.UpdateTranslation(ctx, req.GetId(), req.Data, req.GetUpdateMask())
}

func (r *TagRepo) GetTranslation(ctx context.Context, req *contentV1.GetTagRequest) (*contentV1.TagTranslation, error) {
	if req == nil {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}

	return r.tagTranslationRepo.GetTranslation(ctx, req.GetId(), req.GetLocale())
}

func (r *TagRepo) ListTranslations(ctx context.Context, tagID uint32) ([]*contentV1.TagTranslation, error) {
	return r.tagTranslationRepo.ListTranslations(ctx, tagID, "", nil)
}

func (r *TagRepo) DeleteTranslation(ctx context.Context, req *contentV1.DeleteTagTranslationRequest) error {
	return r.tagTranslationRepo.DeleteTranslation(ctx, req)
}

func (r *TagRepo) CleanTranslations(ctx context.Context, tx *ent.Tx, tagID uint32) error {
	return r.tagTranslationRepo.CleanTranslations(ctx, tx, tagID)
}
