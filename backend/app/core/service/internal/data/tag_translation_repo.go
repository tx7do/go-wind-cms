package data

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	entCrud "github.com/tx7do/go-crud/entgo"
	"github.com/tx7do/go-utils/copierutil"
	"github.com/tx7do/go-utils/mapper"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/predicate"
	"go-wind-cms/app/core/service/internal/data/ent/tagtranslation"

	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
)

type TagTranslationRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper

	mapper *mapper.CopierMapper[contentV1.TagTranslation, ent.TagTranslation]

	repository *entCrud.Repository[
		ent.TagTranslationQuery, ent.TagTranslationSelect,
		ent.TagTranslationCreate, ent.TagTranslationCreateBulk,
		ent.TagTranslationUpdate, ent.TagTranslationUpdateOne,
		ent.TagTranslationDelete,
		predicate.TagTranslation,
		contentV1.TagTranslation, ent.TagTranslation,
	]
}

func NewTagTranslationRepo(ctx *bootstrap.Context, entClient *entCrud.EntClient[*ent.Client]) *TagTranslationRepo {
	repo := &TagTranslationRepo{
		entClient: entClient,
		log:       ctx.NewLoggerHelper("tag-translation/repo/core-service"),
		mapper:    mapper.NewCopierMapper[contentV1.TagTranslation, ent.TagTranslation](),
	}

	repo.init()

	return repo
}

func (r *TagTranslationRepo) init() {
	r.repository = entCrud.NewRepository[
		ent.TagTranslationQuery, ent.TagTranslationSelect,
		ent.TagTranslationCreate, ent.TagTranslationCreateBulk,
		ent.TagTranslationUpdate, ent.TagTranslationUpdateOne,
		ent.TagTranslationDelete,
		predicate.TagTranslation,
		contentV1.TagTranslation, ent.TagTranslation,
	](r.mapper)

	r.mapper.AppendConverters(copierutil.NewTimeStringConverterPair())
	r.mapper.AppendConverters(copierutil.NewTimeTimestamppbConverterPair())
}

func (r *TagTranslationRepo) CleanTranslations(
	ctx context.Context,
	tx *ent.Tx,
	tagID uint32,
) error {
	if _, err := tx.TagTranslation.Delete().
		Where(
			tagtranslation.TagIDEQ(tagID),
		).
		Exec(ctx); err != nil {
		r.log.Errorf("delete old tag [%d] translations failed: %s", tagID, err.Error())
		return contentV1.ErrorInternalServerError("delete old tag translations failed")
	}
	return nil
}

func (r *TagTranslationRepo) ListTranslations(ctx context.Context, tagID uint32) ([]*contentV1.TagTranslation, error) {
	q := r.entClient.Client().TagTranslation.Query().
		Where(
			tagtranslation.TagIDEQ(tagID),
		)

	entities, err := q.
		All(ctx)
	if err != nil {
		r.log.Errorf("query translations by tag id failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("query translations by tag id failed")
	}

	var dtos []*contentV1.TagTranslation
	for _, entity := range entities {
		dtos = append(dtos, r.mapper.ToDTO(entity))
	}

	return dtos, nil
}

func (r *TagTranslationRepo) BatchCreate(ctx context.Context, tx *ent.Tx, items []*contentV1.TagTranslation) error {
	if len(items) == 0 {
		return nil
	}

	now := time.Now()

	builders := make([]*ent.TagTranslationCreate, 0, len(items))
	for _, data := range items {
		builder := tx.TagTranslation.Create().
			SetNillableTagID(data.TagId).
			SetNillableLanguageCode(data.LanguageCode).
			SetNillableName(data.Name).
			SetNillableSlug(data.Slug).
			SetNillableDescription(data.Description).
			SetNillableCoverImage(data.CoverImage).
			SetNillableTemplate(data.Template).
			SetNillableFullPath(data.FullPath).
			SetNillableMetaKeywords(data.MetaKeywords).
			SetNillableMetaDescription(data.MetaDescription).
			SetNillableSeoTitle(data.SeoTitle).
			SetNillableCanonicalURL(data.CanonicalUrl).
			SetNillableCreatedBy(data.CreatedBy).
			SetCreatedAt(now)

		builders = append(builders, builder)
	}

	err := tx.TagTranslation.CreateBulk(builders...).Exec(ctx)
	if err != nil {
		r.log.Errorf("batch create tag translations failed: %s", err.Error())
		return contentV1.ErrorInternalServerError("batch create tag translations failed")
	}

	return nil
}

func (r *TagTranslationRepo) CountByBaseSlug(ctx context.Context, baseSlug string) (int64, error) {
	count, err := r.entClient.Client().TagTranslation.Query().
		Where(
			tagtranslation.SlugHasPrefix(baseSlug),
		).
		Count(ctx)
	if err != nil {
		r.log.Errorf("count tag translations by slug failed: %s", err.Error())
		return 0, contentV1.ErrorInternalServerError("count tag translations by slug failed")
	}

	return int64(count), nil
}

// TranslationExists checks if a translation exists for the given tag ID and language code.
func (r *TagTranslationRepo) TranslationExists(ctx context.Context, tagId uint32, languageCode string) (bool, error) {
	count, err := r.entClient.Client().TagTranslation.Query().
		Where(
			tagtranslation.TagIDEQ(tagId),
			tagtranslation.LanguageCodeEQ(languageCode),
		).
		Count(ctx)
	if err != nil {
		r.log.Errorf("count tag translations by tag id and language code failed: %s", err.Error())
		return false, contentV1.ErrorInternalServerError("count tag translations by tag id and language code failed")
	}

	return count > 0, nil
}

// ListAvailedLanguages lists the language codes of all translations available for the given tag ID.
func (r *TagTranslationRepo) ListAvailedLanguages(ctx context.Context, tagId uint32) ([]string, error) {
	entities, err := r.entClient.Client().TagTranslation.Query().
		Where(
			tagtranslation.TagIDEQ(tagId),
		).
		Select(tagtranslation.FieldLanguageCode).
		Strings(ctx)
	if err != nil {
		r.log.Errorf("query available translation languages by tag id failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("query available translation languages by tag id failed")
	}

	return entities, nil
}
