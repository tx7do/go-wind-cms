package data

import (
	"context"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	entCrud "github.com/tx7do/go-crud/entgo"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	"github.com/tx7do/go-utils/copierutil"
	"github.com/tx7do/go-utils/mapper"
	"github.com/tx7do/go-utils/trans"

	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/predicate"
	"go-wind-cms/app/core/service/internal/data/ent/sectiontranslation"

	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
)

type SectionTranslationRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper

	mapper *mapper.CopierMapper[contentV1.SectionTranslation, ent.SectionTranslation]

	repository *entCrud.Repository[
		ent.SectionTranslationQuery, ent.SectionTranslationSelect,
		ent.SectionTranslationCreate, ent.SectionTranslationCreateBulk,
		ent.SectionTranslationUpdate, ent.SectionTranslationUpdateOne,
		ent.SectionTranslationDelete,
		predicate.SectionTranslation,
		contentV1.SectionTranslation, ent.SectionTranslation,
	]
}

func NewSectionTranslationRepo(ctx *bootstrap.Context, entClient *entCrud.EntClient[*ent.Client]) *SectionTranslationRepo {
	repo := &SectionTranslationRepo{
		entClient: entClient,
		log:       ctx.NewLoggerHelper("section-translation/repo/core-service"),
		mapper:    mapper.NewCopierMapper[contentV1.SectionTranslation, ent.SectionTranslation](),
	}

	repo.init()

	return repo
}

func (r *SectionTranslationRepo) init() {
	r.repository = entCrud.NewRepository[
		ent.SectionTranslationQuery, ent.SectionTranslationSelect,
		ent.SectionTranslationCreate, ent.SectionTranslationCreateBulk,
		ent.SectionTranslationUpdate, ent.SectionTranslationUpdateOne,
		ent.SectionTranslationDelete,
		predicate.SectionTranslation,
		contentV1.SectionTranslation, ent.SectionTranslation,
	](r.mapper)

	r.mapper.AppendConverters(copierutil.NewTimeStringConverterPair())
	r.mapper.AppendConverters(copierutil.NewTimeTimestamppbConverterPair())
}

func (r *SectionTranslationRepo) CleanTranslations(
	ctx context.Context,
	tx *ent.Tx,
	sectionID uint32,
) error {
	if _, err := tx.SectionTranslation.Delete().
		Where(
			sectiontranslation.SectionIDEQ(sectionID),
		).
		Exec(ctx); err != nil {
		r.log.Errorf("delete old section [%d] translations failed: %s", sectionID, err.Error())
		return contentV1.ErrorInternalServerError("delete old section translations failed")
	}
	return nil
}

func (r *SectionTranslationRepo) ListTranslations(ctx context.Context, sectionID uint32) ([]*contentV1.SectionTranslation, error) {
	q := r.entClient.Client().SectionTranslation.Query().
		Where(
			sectiontranslation.SectionIDEQ(sectionID),
		)

	entities, err := q.All(ctx)
	if err != nil {
		r.log.Errorf("query translations by section id failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("query translations by section id failed")
	}

	var dtos []*contentV1.SectionTranslation
	for _, entity := range entities {
		dtos = append(dtos, r.mapper.ToDTO(entity))
	}

	return dtos, nil
}

func (r *SectionTranslationRepo) newCreateBuilder(data *contentV1.SectionTranslation) *ent.SectionTranslationCreate {
	now := time.Now()

	builder := r.entClient.Client().SectionTranslation.Create().
		SetNillableSectionID(data.SectionId).
		SetNillableLanguageCode(data.LanguageCode).
		SetNillableCreatedBy(data.CreatedBy).
		SetCreatedAt(now)

	if data.Content != nil {
		builder.SetContent(trans.Ptr(data.GetContent()))
	}

	return builder
}

func (r *SectionTranslationRepo) BatchCreate(ctx context.Context, tx *ent.Tx, items []*contentV1.SectionTranslation) error {
	if len(items) == 0 {
		return nil
	}

	builders := make([]*ent.SectionTranslationCreate, 0, len(items))
	for _, data := range items {
		builder := r.newCreateBuilder(data)
		builders = append(builders, builder)
	}

	err := tx.SectionTranslation.CreateBulk(builders...).Exec(ctx)
	if err != nil {
		r.log.Errorf("batch create section translations failed: %s", err.Error())
		return contentV1.ErrorInternalServerError("batch create section translations failed")
	}

	return nil
}

func (r *SectionTranslationRepo) CreateTranslation(ctx context.Context, data *contentV1.SectionTranslation) (*contentV1.SectionTranslation, error) {
	builder := r.newCreateBuilder(data)

	entity, err := builder.Save(ctx)
	if err != nil {
		r.log.Errorf("create section translation failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("create section translation failed")
	}

	dto := r.mapper.ToDTO(entity)

	return dto, nil
}

func (r *SectionTranslationRepo) UpdateTranslation(ctx context.Context, id uint32, data *contentV1.SectionTranslation, updateMask *fieldmaskpb.FieldMask) (*contentV1.SectionTranslation, error) {
	if data == nil {
		return nil, nil
	}

	builder := r.entClient.Client().SectionTranslation.UpdateOneID(id)
	// 租户作用域：仅更新本租户翻译，避免跨租户改他人翻译（按 hasTenant 条件加）
	if tid, hasTenant := maybeTenantFromViewer(ctx); hasTenant {
		builder.Where(sectiontranslation.TenantIDEQ(tid))
	}

	dto, err := r.repository.UpdateOne(ctx, builder, data, updateMask,
		func(dto *contentV1.SectionTranslation) {
			builder.
				SetNillableSectionID(data.SectionId).
				SetNillableLanguageCode(data.LanguageCode).
				SetNillableUpdatedBy(data.UpdatedBy).
				SetUpdatedAt(time.Now())

			if data.Content != nil {
				builder.SetContent(trans.Ptr(data.GetContent()))
			}
		},
		func(s *sql.Selector) {
			s.Where(sql.EQ(sectiontranslation.FieldID, id))
		},
	)
	if err != nil {
		r.log.Errorf("update section translation failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("update section translation failed")
	}

	return dto, nil
}

// TranslationExists checks if a translation exists for the given section ID and language code.
func (r *SectionTranslationRepo) TranslationExists(ctx context.Context, sectionId uint32, languageCode string) (bool, error) {
	c, err := r.entClient.Client().SectionTranslation.Query().
		Where(
			sectiontranslation.SectionIDEQ(sectionId),
			sectiontranslation.LanguageCodeEQ(languageCode),
		).
		Count(ctx)
	if err != nil {
		r.log.Errorf("count section translations by section id and language code failed: %s", err.Error())
		return false, contentV1.ErrorInternalServerError("count section translations by section id and language code failed")
	}

	return c > 0, nil
}

// ListAvailedLanguages lists the language codes of all translations available for the given section ID.
func (r *SectionTranslationRepo) ListAvailedLanguages(ctx context.Context, sectionId uint32) ([]string, error) {
	entities, err := r.entClient.Client().SectionTranslation.Query().
		Where(
			sectiontranslation.SectionIDEQ(sectionId),
		).
		Select(sectiontranslation.FieldLanguageCode).
		Strings(ctx)
	if err != nil {
		r.log.Errorf("query available translation languages by section id failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("query available translation languages by section id failed")
	}

	return entities, nil
}

func (r *SectionTranslationRepo) GetTranslation(ctx context.Context, sectionId uint32, languageCode string) (*contentV1.SectionTranslation, error) {
	entity, err := r.entClient.Client().SectionTranslation.Query().
		Where(
			sectiontranslation.SectionIDEQ(sectionId),
			sectiontranslation.LanguageCodeEQ(languageCode),
		).
		Only(ctx)
	if err != nil {
		r.log.Errorf("query section translation by section id and language code failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("query section translation by section id and language code failed")
	}

	dto := r.mapper.ToDTO(entity)

	return dto, nil
}

func (r *SectionTranslationRepo) DeleteTranslation(ctx context.Context, req *contentV1.DeleteSectionTranslationRequest) error {
	if req.QueryBy == nil {
		return contentV1.ErrorBadRequest("invalid parameter: query_by is required")
	}

	switch req.QueryBy.(type) {
	case *contentV1.DeleteSectionTranslationRequest_Id:
		if req.GetId() == 0 {
			return contentV1.ErrorBadRequest("invalid parameter: id must be greater than 0")
		}

	case *contentV1.DeleteSectionTranslationRequest_Identifier:
		if req.GetIdentifier() == nil {
			return contentV1.ErrorBadRequest("invalid parameter: identifier is required")
		}
		if req.GetIdentifier().GetSectionId() == 0 {
			return contentV1.ErrorBadRequest("invalid parameter: section_id must be greater than 0")
		}
		if len(req.GetIdentifier().GetLanguageCode()) == 0 {
			return contentV1.ErrorBadRequest("invalid parameter: language_code is required")
		}

	default:
		return contentV1.ErrorBadRequest("invalid parameter: unsupported query_by type")
	}

	builder := r.entClient.Client().SectionTranslation.Delete()
	// 租户作用域：仅删除本租户翻译，避免跨租户删他人翻译（按 hasTenant 条件加）
	if tid, hasTenant := maybeTenantFromViewer(ctx); hasTenant {
		builder.Where(sectiontranslation.TenantIDEQ(tid))
	}

	_, err := r.repository.Delete(ctx, builder, func(s *sql.Selector) {
		switch req.QueryBy.(type) {
		case *contentV1.DeleteSectionTranslationRequest_Id:
			id := req.GetId()
			s.Where(sql.EQ(sectiontranslation.FieldID, id))

		case *contentV1.DeleteSectionTranslationRequest_Identifier:
			identifier := req.GetIdentifier()
			s.Where(
				sql.And(
					sql.EQ(sectiontranslation.FieldSectionID, identifier.GetSectionId()),
					sql.EQ(sectiontranslation.FieldLanguageCode, identifier.GetLanguageCode()),
				),
			)

		default:
			return
		}
	})
	if err != nil {
		r.log.Errorf("delete section translation failed: %s", err.Error())
		return contentV1.ErrorInternalServerError("delete section translation failed")
	}

	return nil
}
