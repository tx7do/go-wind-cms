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
	"go-wind-cms/app/core/service/internal/data/ent/posttranslation"
	"go-wind-cms/app/core/service/internal/data/ent/predicate"

	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
)

type PostTranslationRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper

	mapper *mapper.CopierMapper[contentV1.PostTranslation, ent.PostTranslation]

	repository *entCrud.Repository[
		ent.PostTranslationQuery, ent.PostTranslationSelect,
		ent.PostTranslationCreate, ent.PostTranslationCreateBulk,
		ent.PostTranslationUpdate, ent.PostTranslationUpdateOne,
		ent.PostTranslationDelete,
		predicate.PostTranslation,
		contentV1.PostTranslation, ent.PostTranslation,
	]
}

func NewPostTranslationRepo(ctx *bootstrap.Context, entClient *entCrud.EntClient[*ent.Client]) *PostTranslationRepo {
	repo := &PostTranslationRepo{
		entClient: entClient,
		mapper:    mapper.NewCopierMapper[contentV1.PostTranslation, ent.PostTranslation](),
		log:       ctx.NewLoggerHelper("post-translation/repo/core-service"),
	}

	repo.init()

	return repo
}

func (r *PostTranslationRepo) init() {
	r.repository = entCrud.NewRepository[
		ent.PostTranslationQuery, ent.PostTranslationSelect,
		ent.PostTranslationCreate, ent.PostTranslationCreateBulk,
		ent.PostTranslationUpdate, ent.PostTranslationUpdateOne,
		ent.PostTranslationDelete,
		predicate.PostTranslation,
		contentV1.PostTranslation, ent.PostTranslation,
	](r.mapper)

	r.mapper.AppendConverters(copierutil.NewTimeStringConverterPair())
	r.mapper.AppendConverters(copierutil.NewTimeTimestamppbConverterPair())
}

func (r *PostTranslationRepo) CleanTranslations(
	ctx context.Context,
	tx *ent.Tx,
	postID uint32,
) error {
	if _, err := tx.PostTranslation.Delete().
		Where(
			posttranslation.PostIDEQ(postID),
		).
		Exec(ctx); err != nil {
		r.log.Errorf("delete old post [%d] translations failed: %s", postID, err.Error())
		return contentV1.ErrorInternalServerError("delete old post translations failed")
	}
	return nil
}

func (r *PostTranslationRepo) ListTranslations(ctx context.Context, postID uint32) ([]*contentV1.PostTranslation, error) {
	q := r.entClient.Client().PostTranslation.Query().
		Where(
			posttranslation.PostIDEQ(postID),
		)

	entities, err := q.
		All(ctx)
	if err != nil {
		r.log.Errorf("query translations by post id failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("query translations by post id failed")
	}

	var dtos []*contentV1.PostTranslation
	for _, entity := range entities {
		dtos = append(dtos, r.mapper.ToDTO(entity))
	}

	return dtos, nil
}

func (r *PostTranslationRepo) BatchCreate(ctx context.Context, tx *ent.Tx, items []*contentV1.PostTranslation) error {
	if len(items) == 0 {
		return nil
	}

	now := time.Now()

	builders := make([]*ent.PostTranslationCreate, 0, len(items))
	for _, data := range items {
		builder := tx.PostTranslation.Create().
			SetNillablePostID(data.PostId).
			SetNillableLanguageCode(data.LanguageCode).
			SetNillableTitle(data.Title).
			SetNillableSlug(data.Slug).
			SetNillableSummary(data.Summary).
			SetNillableContent(data.Content).
			SetNillableOriginalContent(data.OriginalContent).
			SetNillableThumbnail(data.Thumbnail).
			SetNillableTemplate(data.Template).
			SetNillableWordCount(data.WordCount).
			SetNillableFullPath(data.FullPath).
			SetNillableMetaKeywords(data.MetaKeywords).
			SetNillableMetaDescription(data.MetaDescription).
			SetNillableSeoTitle(data.SeoTitle).
			SetNillableCreatedBy(data.CreatedBy).
			SetCreatedAt(now)

		builders = append(builders, builder)
	}

	err := tx.PostTranslation.CreateBulk(builders...).Exec(ctx)
	if err != nil {
		r.log.Errorf("batch create post translations failed: %s", err.Error())
		return contentV1.ErrorInternalServerError("batch create post translations failed")
	}

	return nil
}
