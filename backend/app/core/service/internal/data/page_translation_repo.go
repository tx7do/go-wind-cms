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
	"go-wind-cms/app/core/service/internal/data/ent/pagetranslation"
	"go-wind-cms/app/core/service/internal/data/ent/predicate"

	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
)

type PageTranslationRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper

	mapper *mapper.CopierMapper[contentV1.PageTranslation, ent.PageTranslation]

	repository *entCrud.Repository[
		ent.PageTranslationQuery, ent.PageTranslationSelect,
		ent.PageTranslationCreate, ent.PageTranslationCreateBulk,
		ent.PageTranslationUpdate, ent.PageTranslationUpdateOne,
		ent.PageTranslationDelete,
		predicate.PageTranslation,
		contentV1.PageTranslation, ent.PageTranslation,
	]
}

func NewPageTranslationRepo(ctx *bootstrap.Context, entClient *entCrud.EntClient[*ent.Client]) *PageTranslationRepo {
	repo := &PageTranslationRepo{
		entClient: entClient,
		log:       ctx.NewLoggerHelper("page-translation/repo/core-service"),
		mapper:    mapper.NewCopierMapper[contentV1.PageTranslation, ent.PageTranslation](),
	}

	repo.init()

	return repo
}

func (r *PageTranslationRepo) init() {
	r.repository = entCrud.NewRepository[
		ent.PageTranslationQuery, ent.PageTranslationSelect,
		ent.PageTranslationCreate, ent.PageTranslationCreateBulk,
		ent.PageTranslationUpdate, ent.PageTranslationUpdateOne,
		ent.PageTranslationDelete,
		predicate.PageTranslation,
		contentV1.PageTranslation, ent.PageTranslation,
	](r.mapper)

	r.mapper.AppendConverters(copierutil.NewTimeStringConverterPair())
	r.mapper.AppendConverters(copierutil.NewTimeTimestamppbConverterPair())
}

func (r *PageTranslationRepo) CleanTranslations(
	ctx context.Context,
	tx *ent.Tx,
	pageID uint32,
) error {
	if _, err := tx.PageTranslation.Delete().
		Where(
			pagetranslation.PageIDEQ(pageID),
		).
		Exec(ctx); err != nil {
		r.log.Errorf("delete old page [%d] translations failed: %s", pageID, err.Error())
		return contentV1.ErrorInternalServerError("delete old page translations failed")
	}
	return nil
}

func (r *PageTranslationRepo) ListTranslations(ctx context.Context, pageID uint32) ([]*contentV1.PageTranslation, error) {
	q := r.entClient.Client().PageTranslation.Query().
		Where(
			pagetranslation.PageIDEQ(pageID),
		)

	entities, err := q.
		All(ctx)
	if err != nil {
		r.log.Errorf("query translations by page id failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("query translations by page id failed")
	}

	var dtos []*contentV1.PageTranslation
	for _, entity := range entities {
		dtos = append(dtos, r.mapper.ToDTO(entity))
	}

	return dtos, nil
}

func (r *PageTranslationRepo) BatchCreate(ctx context.Context, tx *ent.Tx, items []*contentV1.PageTranslation) error {
	if len(items) == 0 {
		return nil
	}

	now := time.Now()

	builders := make([]*ent.PageTranslationCreate, 0, len(items))
	for _, data := range items {
		builder := tx.PageTranslation.Create().
			SetNillablePageID(data.PageId).
			SetNillableLanguageCode(data.LanguageCode).
			SetNillableTitle(data.Title).
			SetNillableSlug(data.Slug).
			SetNillableSummary(data.Summary).
			SetNillableContent(data.Content).
			SetNillableOriginalContent(data.OriginalContent).
			SetNillableThumbnail(data.Thumbnail).
			SetNillableCoverImage(data.CoverImage).
			SetNillableWordCount(data.WordCount).
			SetNillableFullPath(data.FullPath).
			SetNillableMetaKeywords(data.MetaKeywords).
			SetNillableMetaDescription(data.MetaDescription).
			SetNillableSeoTitle(data.SeoTitle).
			SetNillableCreatedBy(data.CreatedBy).
			SetCreatedAt(now)

		builders = append(builders, builder)
	}

	err := tx.PageTranslation.CreateBulk(builders...).Exec(ctx)
	if err != nil {
		r.log.Errorf("batch create page translations failed: %s", err.Error())
		return contentV1.ErrorInternalServerError("batch create page translations failed")
	}

	return nil
}
