package data

import (
	"context"
	"time"

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

func (r *SectionTranslationRepo) ListTranslations(ctx context.Context, sectionID uint32, locale string, viewMask *fieldmaskpb.FieldMask) ([]*contentV1.SectionTranslation, error) {
	builder := r.entClient.Client().SectionTranslation.Query().
		Where(
			sectiontranslation.SectionIDEQ(sectionID),
		)

	if len(locale) > 0 {
		builder.Where(
			sectiontranslation.LanguageCodeEQ(locale),
		)
	}

	if viewMask != nil {
		selectSelector, err := r.repository.BuildSelector(viewMask.GetPaths())
		if err != nil {
			r.log.Errorf("build section translation selector failed: %s", err.Error())
			return nil, contentV1.ErrorInternalServerError("build section translation selector failed")
		}
		if selectSelector != nil {
			builder.Modify(selectSelector)
		}
	}

	entities, err := builder.
		Order(ent.Asc(sectiontranslation.FieldID)).
		All(ctx)
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

// ListTranslationsWithFallback 优先返回指定语言的译文；该语言缺译时返回全部译文，
// 供前端按"匹配当前语言→回退第一条"兜底（只回传空数组会让前端无 fallback 可用）。
func (r *SectionTranslationRepo) ListTranslationsWithFallback(ctx context.Context, sectionID uint32, locale string, viewMask *fieldmaskpb.FieldMask) ([]*contentV1.SectionTranslation, error) {
	translations, err := r.ListTranslations(ctx, sectionID, locale, viewMask)
	if err != nil {
		return nil, err
	}
	if len(translations) == 0 && locale != "" {
		return r.ListTranslations(ctx, sectionID, "", viewMask)
	}
	return translations, nil
}

// st 必须传与调用方一致的事务/非事务客户端，保证批量插入运行在事务内
// （见 PostTranslationRepo.PrepareTranslation 注释）。
func (r *SectionTranslationRepo) newCreateBuilder(st *ent.SectionTranslationClient, data *contentV1.SectionTranslation) *ent.SectionTranslationCreate {
	now := time.Now()

	builder := st.Create().
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
		builder := r.newCreateBuilder(tx.SectionTranslation, data)
		builders = append(builders, builder)
	}

	err := tx.SectionTranslation.CreateBulk(builders...).Exec(ctx)
	if err != nil {
		r.log.Errorf("batch create section translations failed: %s", err.Error())
		return contentV1.ErrorInternalServerError("batch create section translations failed")
	}

	return nil
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
