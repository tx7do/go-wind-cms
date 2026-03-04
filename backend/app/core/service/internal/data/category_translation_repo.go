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
	"go-wind-cms/app/core/service/internal/data/ent/categorytranslation"
	"go-wind-cms/app/core/service/internal/data/ent/predicate"

	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
)

type CategoryTranslationRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper

	mapper *mapper.CopierMapper[contentV1.CategoryTranslation, ent.CategoryTranslation]

	repository *entCrud.Repository[
		ent.CategoryTranslationQuery, ent.CategoryTranslationSelect,
		ent.CategoryTranslationCreate, ent.CategoryTranslationCreateBulk,
		ent.CategoryTranslationUpdate, ent.CategoryTranslationUpdateOne,
		ent.CategoryTranslationDelete,
		predicate.CategoryTranslation,
		contentV1.CategoryTranslation, ent.CategoryTranslation,
	]
}

func NewCategoryTranslationRepo(ctx *bootstrap.Context, entClient *entCrud.EntClient[*ent.Client]) *CategoryTranslationRepo {
	repo := &CategoryTranslationRepo{
		entClient: entClient,
		log:       ctx.NewLoggerHelper("category-translation/repo/core-service"),
		mapper:    mapper.NewCopierMapper[contentV1.CategoryTranslation, ent.CategoryTranslation](),
	}

	repo.init()

	return repo
}

func (r *CategoryTranslationRepo) init() {
	r.repository = entCrud.NewRepository[
		ent.CategoryTranslationQuery, ent.CategoryTranslationSelect,
		ent.CategoryTranslationCreate, ent.CategoryTranslationCreateBulk,
		ent.CategoryTranslationUpdate, ent.CategoryTranslationUpdateOne,
		ent.CategoryTranslationDelete,
		predicate.CategoryTranslation,
		contentV1.CategoryTranslation, ent.CategoryTranslation,
	](r.mapper)

	r.mapper.AppendConverters(copierutil.NewTimeStringConverterPair())
	r.mapper.AppendConverters(copierutil.NewTimeTimestamppbConverterPair())
}

func (r *CategoryTranslationRepo) CleanTranslations(
	ctx context.Context,
	tx *ent.Tx,
	categoryID uint32,
) error {
	if _, err := tx.CategoryTranslation.Delete().
		Where(
			categorytranslation.CategoryIDEQ(categoryID),
		).
		Exec(ctx); err != nil {
		r.log.Errorf("delete old category [%d] translations failed: %s", categoryID, err.Error())
		return contentV1.ErrorInternalServerError("delete old category translations failed")
	}
	return nil
}

func (r *CategoryTranslationRepo) ListTranslations(ctx context.Context, categoryID uint32) ([]*contentV1.CategoryTranslation, error) {
	q := r.entClient.Client().CategoryTranslation.Query().
		Where(
			categorytranslation.CategoryIDEQ(categoryID),
		)

	entities, err := q.
		All(ctx)
	if err != nil {
		r.log.Errorf("query translations by category id failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("query translations by category id failed")
	}

	var dtos []*contentV1.CategoryTranslation
	for _, entity := range entities {
		dtos = append(dtos, r.mapper.ToDTO(entity))
	}

	return dtos, nil
}

func (r *CategoryTranslationRepo) BatchCreate(ctx context.Context, tx *ent.Tx, items []*contentV1.CategoryTranslation) error {
	if len(items) == 0 {
		return nil
	}

	now := time.Now()

	builders := make([]*ent.CategoryTranslationCreate, 0, len(items))
	for _, data := range items {
		builder := tx.CategoryTranslation.Create().
			SetNillableCategoryID(data.CategoryId).
			SetNillableLanguageCode(data.LanguageCode).
			SetNillableName(data.Name).
			SetNillableSlug(data.Slug).
			SetNillableDescription(data.Description).
			SetNillableThumbnail(data.Thumbnail).
			SetNillableCoverImage(data.CoverImage).
			SetNillableTemplate(data.Template).
			SetNillableFullPath(data.FullPath).
			SetNillableMetaKeywords(data.MetaKeywords).
			SetNillableMetaDescription(data.MetaDescription).
			SetNillableSeoTitle(data.SeoTitle).
			SetNillableCreatedBy(data.CreatedBy).
			SetCreatedAt(now)

		builders = append(builders, builder)
	}

	err := tx.CategoryTranslation.CreateBulk(builders...).Exec(ctx)
	if err != nil {
		r.log.Errorf("batch create category translations failed: %s", err.Error())
		return contentV1.ErrorInternalServerError("batch create category translations failed")
	}

	return nil
}

func (r *CategoryTranslationRepo) CountByBaseSlug(ctx context.Context, baseSlug string) (int64, error) {
	count, err := r.entClient.Client().CategoryTranslation.Query().
		Where(
			categorytranslation.SlugHasPrefix(baseSlug),
		).
		Count(ctx)
	if err != nil {
		r.log.Errorf("count category translations by slug failed: %s", err.Error())
		return 0, contentV1.ErrorInternalServerError("count category translations by slug failed")
	}

	return int64(count), nil
}

// TranslationExists checks if a translation exists for the given category ID and language code.
func (r *CategoryTranslationRepo) TranslationExists(ctx context.Context, categoryId uint32, languageCode string) (bool, error) {
	count, err := r.entClient.Client().CategoryTranslation.Query().
		Where(
			categorytranslation.CategoryIDEQ(categoryId),
			categorytranslation.LanguageCodeEQ(languageCode),
		).
		Count(ctx)
	if err != nil {
		r.log.Errorf("count category translations by category id and language code failed: %s", err.Error())
		return false, contentV1.ErrorInternalServerError("count category translations by category id and language code failed")
	}

	return count > 0, nil
}

// ListAvailedLanguages lists the language codes of all translations available for the given category ID.
func (r *CategoryTranslationRepo) ListAvailedLanguages(ctx context.Context, categoryId uint32) ([]string, error) {
	entities, err := r.entClient.Client().CategoryTranslation.Query().
		Where(
			categorytranslation.CategoryIDEQ(categoryId),
		).
		Select(categorytranslation.FieldLanguageCode).
		Strings(ctx)
	if err != nil {
		r.log.Errorf("query available translation languages by category id failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("query available translation languages by category id failed")
	}

	return entities, nil
}
