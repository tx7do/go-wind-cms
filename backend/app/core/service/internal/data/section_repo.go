package data

import (
	"context"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	entCrud "github.com/tx7do/go-crud/entgo"

	"github.com/tx7do/go-utils/copierutil"
	"github.com/tx7do/go-utils/mapper"
	"github.com/tx7do/go-utils/trans"

	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/predicate"
	"go-wind-cms/app/core/service/internal/data/ent/section"

	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
)

type SectionRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper

	mapper *mapper.CopierMapper[contentV1.Section, ent.Section]

	repository *entCrud.Repository[
		ent.SectionQuery, ent.SectionSelect,
		ent.SectionCreate, ent.SectionCreateBulk,
		ent.SectionUpdate, ent.SectionUpdateOne,
		ent.SectionDelete,
		predicate.Section,
		contentV1.Section, ent.Section,
	]

	typeConverter *mapper.EnumTypeConverter[contentV1.SectionType, section.Type]

	sectionTranslationRepo *SectionTranslationRepo
}

func NewSectionRepo(
	ctx *bootstrap.Context,
	entClient *entCrud.EntClient[*ent.Client],
	sectionTranslationRepo *SectionTranslationRepo,
) *SectionRepo {
	repo := &SectionRepo{
		entClient: entClient,
		log:       ctx.NewLoggerHelper("section/repo/core-service"),
		mapper:    mapper.NewCopierMapper[contentV1.Section, ent.Section](),
		typeConverter: mapper.NewEnumTypeConverter[contentV1.SectionType, section.Type](
			contentV1.SectionType_name, contentV1.SectionType_value,
		),
		sectionTranslationRepo: sectionTranslationRepo,
	}

	repo.init()

	return repo
}

func (r *SectionRepo) init() {
	r.repository = entCrud.NewRepository[
		ent.SectionQuery, ent.SectionSelect,
		ent.SectionCreate, ent.SectionCreateBulk,
		ent.SectionUpdate, ent.SectionUpdateOne,
		ent.SectionDelete,
		predicate.Section,
		contentV1.Section, ent.Section,
	](r.mapper)

	r.mapper.AppendConverters(copierutil.NewTimeStringConverterPair())
	r.mapper.AppendConverters(copierutil.NewTimeTimestamppbConverterPair())

	r.mapper.AppendConverters(r.typeConverter.NewConverterPair())
}

func (r *SectionRepo) IsExist(ctx context.Context, id uint32) (bool, error) {
	exist, err := r.entClient.Client().Section.Query().
		Where(section.IDEQ(id)).
		Exist(ctx)
	if err != nil {
		r.log.Errorf("query section exist failed: %s", err.Error())
		return false, contentV1.ErrorInternalServerError("query section exist failed")
	}
	return exist, nil
}

func (r *SectionRepo) List(ctx context.Context, req *paginationV1.PagingRequest) (*contentV1.ListSectionResponse, error) {
	if req == nil {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}

	builder := r.entClient.Client().Section.Query()

	ret, err := r.repository.ListWithPaging(ctx, builder, builder.Clone(), req)
	if err != nil {
		return nil, err
	}
	if ret == nil {
		return &contentV1.ListSectionResponse{Total: 0, Items: nil}, nil
	}

	for _, item := range ret.Items {
		languages, err := r.sectionTranslationRepo.ListAvailedLanguages(ctx, item.GetId())
		if err != nil {
			r.log.Errorf("query availed languages failed: %s", err.Error())
			return nil, contentV1.ErrorInternalServerError("query availed languages failed")
		}
		item.AvailableLanguages = languages
	}

	return &contentV1.ListSectionResponse{
		Total: ret.Total,
		Items: ret.Items,
	}, nil
}

func (r *SectionRepo) Get(ctx context.Context, req *contentV1.GetSectionRequest) (*contentV1.Section, error) {
	if req == nil {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}

	builder := r.entClient.Client().Section.Query()

	switch req.QueryBy.(type) {
	case *contentV1.GetSectionRequest_Id:
		builder.Where(section.IDEQ(req.GetId()))
	default:
		return nil, contentV1.ErrorBadRequest("invalid query_by value")
	}

	entity, err := builder.Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, contentV1.ErrorFileNotFound("section not found")
		}

		r.log.Errorf("query section failed: %s", err.Error())

		return nil, contentV1.ErrorInternalServerError("query section failed")
	}

	dto := r.mapper.ToDTO(entity)

	translations, err := r.sectionTranslationRepo.ListTranslations(ctx, dto.GetId())
	if err != nil {
		r.log.Errorf("query translations failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("query translations failed")
	}
	dto.Translations = translations

	languages, err := r.sectionTranslationRepo.ListAvailedLanguages(ctx, dto.GetId())
	if err != nil {
		r.log.Errorf("query availed languages failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("query availed languages failed")
	}
	dto.AvailableLanguages = languages

	return dto, nil
}

func (r *SectionRepo) Create(ctx context.Context, req *contentV1.CreateSectionRequest) (dto *contentV1.Section, err error) {
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

	builder := tx.Section.Create().
		SetNillablePageID(req.Data.PageId).
		SetNillableType(r.typeConverter.ToEntity(req.Data.Type)).
		SetNillableName(req.Data.Name).
		SetNillableSortOrder(req.Data.SortOrder).
		SetNillableCreatedBy(req.Data.CreatedBy).
		SetCreatedAt(time.Now())

	if req.Data.Config != nil {
		builder.SetConfig(trans.Ptr(req.Data.GetConfig()))
	}

	var entity *ent.Section
	if entity, err = builder.Save(ctx); err != nil {
		r.log.Errorf("insert section failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("insert section failed")
	}

	if len(req.Data.Translations) > 0 {
		if err = r.sectionTranslationRepo.CleanTranslations(ctx, tx, entity.ID); err != nil {
			r.log.Errorf("clean translations failed: %s", err.Error())
			return nil, contentV1.ErrorInternalServerError("clean translations failed")
		}

		for i := range req.Data.Translations {
			req.Data.Translations[i].SectionId = trans.Ptr(entity.ID)
		}

		if err = r.sectionTranslationRepo.BatchCreate(ctx, tx, req.Data.GetTranslations()); err != nil {
			r.log.Errorf("batch insert translations failed: %s", err.Error())
			return nil, contentV1.ErrorInternalServerError("batch insert translations failed")
		}
	}

	return r.mapper.ToDTO(entity), nil
}

func (r *SectionRepo) Update(ctx context.Context, req *contentV1.UpdateSectionRequest) (dto *contentV1.Section, err error) {
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
			_, err = r.Create(ctx, &contentV1.CreateSectionRequest{Data: req.Data})
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
		for i := range req.Data.Translations {
			req.Data.Translations[i].SectionId = trans.Ptr(req.GetId())
		}

		if err = r.sectionTranslationRepo.BatchCreate(ctx, tx, req.Data.GetTranslations()); err != nil {
			r.log.Errorf("batch insert translations failed: %s", err.Error())
			return nil, contentV1.ErrorInternalServerError("batch insert translations failed")
		}
	}

	builder := tx.Section.UpdateOneID(req.GetId())
	result, err := r.repository.UpdateOne(ctx, builder, req.Data, req.GetUpdateMask(),
		func(dto *contentV1.Section) {
			builder.
				SetNillablePageID(req.Data.PageId).
				SetNillableType(r.typeConverter.ToEntity(req.Data.Type)).
				SetNillableName(req.Data.Name).
				SetNillableSortOrder(req.Data.SortOrder).
				SetNillableUpdatedBy(req.Data.UpdatedBy).
				SetUpdatedAt(time.Now())

			if req.Data.Config != nil {
				builder.SetConfig(trans.Ptr(req.Data.GetConfig()))
			}
		},
		func(s *sql.Selector) {
			s.Where(sql.EQ(section.FieldID, req.GetId()))
		},
	)

	return result, err
}

func (r *SectionRepo) Delete(ctx context.Context, req *contentV1.DeleteSectionRequest) (err error) {
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

	if err = tx.Section.
		DeleteOneID(req.GetId()).
		Exec(ctx); err != nil {
		r.log.Errorf("delete one data failed: %s", err.Error())
		return contentV1.ErrorInternalServerError("delete one data failed")
	}

	if err = r.sectionTranslationRepo.CleanTranslations(ctx, tx, req.GetId()); err != nil {
		r.log.Errorf("clean translations failed: %s", err.Error())
		return contentV1.ErrorInternalServerError("clean translations failed")
	}

	return nil
}

func (r *SectionRepo) TranslationExists(ctx context.Context, sectionId uint32, languageCode string) (bool, error) {
	return r.sectionTranslationRepo.TranslationExists(ctx, sectionId, languageCode)
}

func (r *SectionRepo) CreateTranslation(ctx context.Context, req *contentV1.CreateSectionTranslationRequest) (*contentV1.SectionTranslation, error) {
	if req == nil || req.Data == nil {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}

	if len(req.Data.GetLanguageCode()) == 0 {
		return nil, contentV1.ErrorBadRequest("language code is required")
	}

	if req.GetSectionId() == 0 {
		return nil, contentV1.ErrorBadRequest("section id is required")
	}

	req.Data.SectionId = trans.Ptr(req.GetSectionId())

	return r.sectionTranslationRepo.CreateTranslation(ctx, req.Data)
}

func (r *SectionRepo) UpdateTranslation(ctx context.Context, req *contentV1.UpdateSectionTranslationRequest) (*contentV1.SectionTranslation, error) {
	if req == nil || req.Data == nil {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}

	if len(req.Data.GetLanguageCode()) == 0 {
		return nil, contentV1.ErrorBadRequest("language code is required")
	}

	if req.Data.GetSectionId() == 0 {
		return nil, contentV1.ErrorBadRequest("section id is required")
	}

	if exist, err := r.TranslationExists(ctx, req.Data.GetSectionId(), req.Data.GetLanguageCode()); err != nil {
		return nil, err
	} else if !exist {
		if req.GetAllowMissing() {
			return r.CreateTranslation(ctx, &contentV1.CreateSectionTranslationRequest{
				Data:      req.Data,
				SectionId: req.Data.GetSectionId(),
			})
		}

		return nil, contentV1.ErrorFileNotFound("translation not found")
	}

	return r.sectionTranslationRepo.UpdateTranslation(ctx, req.GetId(), req.Data, req.GetUpdateMask())
}

func (r *SectionRepo) GetTranslation(ctx context.Context, req *contentV1.GetSectionRequest) (*contentV1.SectionTranslation, error) {
	if req == nil {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}

	return r.sectionTranslationRepo.GetTranslation(ctx, req.GetId(), req.GetLocale())
}

func (r *SectionRepo) ListTranslations(ctx context.Context, sectionID uint32) ([]*contentV1.SectionTranslation, error) {
	return r.sectionTranslationRepo.ListTranslations(ctx, sectionID)
}

func (r *SectionRepo) DeleteTranslation(ctx context.Context, req *contentV1.DeleteSectionTranslationRequest) error {
	return r.sectionTranslationRepo.DeleteTranslation(ctx, req)
}

func (r *SectionRepo) CleanTranslations(ctx context.Context, tx *ent.Tx, sectionID uint32) error {
	return r.sectionTranslationRepo.CleanTranslations(ctx, tx, sectionID)
}

// CleanByPageID 删除指定页面下的所有 section 及其翻译。
// 在页面删除时调用，防止 section 成为孤儿记录。
func (r *SectionRepo) CleanByPageID(ctx context.Context, tx *ent.Tx, pageID uint32) error {
	if pageID == 0 {
		return nil
	}

	// 先查出所有 section ID，用于清理翻译
	sectionIDs, err := tx.Section.Query().
		Where(section.PageIDEQ(pageID)).
		IDs(ctx)
	if err != nil {
		r.log.Errorf("query sections by page id failed: %s", err.Error())
		return contentV1.ErrorInternalServerError("query sections by page id failed")
	}

	// 清理每个 section 的翻译
	for _, sid := range sectionIDs {
		if err := r.sectionTranslationRepo.CleanTranslations(ctx, tx, sid); err != nil {
			r.log.Errorf("clean section translations failed for section %d: %s", sid, err.Error())
			return contentV1.ErrorInternalServerError("clean section translations failed")
		}
	}

	// 删除所有 section
	if _, err := tx.Section.Delete().
		Where(section.PageIDEQ(pageID)).
		Exec(ctx); err != nil {
		r.log.Errorf("delete sections by page id failed: %s", err.Error())
		return contentV1.ErrorInternalServerError("delete sections by page id failed")
	}

	return nil
}
