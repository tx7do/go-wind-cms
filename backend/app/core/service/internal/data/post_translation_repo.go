package data

import (
	"context"
	"strconv"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	entCrud "github.com/tx7do/go-crud/entgo"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	"github.com/tx7do/go-utils/copierutil"
	"github.com/tx7do/go-utils/mapper"
	"github.com/tx7do/go-utils/slug"
	"github.com/tx7do/go-utils/trans"

	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/posttranslation"
	"go-wind-cms/app/core/service/internal/data/ent/predicate"

	contentV1 "go-wind-cms/api/gen/go/content/service/v1"

	"go-wind-cms/pkg/content/count"
	"go-wind-cms/pkg/content/summary"
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

func (r *PostTranslationRepo) ListTranslations(ctx context.Context, postID uint32, locale string, viewMask *fieldmaskpb.FieldMask) ([]*contentV1.PostTranslation, error) {
	builder := r.entClient.Client().PostTranslation.Query().
		Where(
			posttranslation.PostIDEQ(postID),
		)

	if len(locale) > 0 {
		builder.Where(
			posttranslation.LanguageCodeEQ(locale),
		)
	}

	if viewMask != nil {
		selectSelector, err := r.repository.BuildSelector(viewMask.GetPaths())
		if err != nil {
			r.log.Errorf("build post translation selector failed: %s", err.Error())
			return nil, contentV1.ErrorInternalServerError("build post translation selector failed")
		}
		if selectSelector != nil {
			builder.Modify(selectSelector)
		}
	}

	entities, err := builder.
		All(ctx)
	if err != nil {
		r.log.Errorf("query post translations failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("query post translations failed")
	}

	var dtos []*contentV1.PostTranslation
	for _, entity := range entities {
		dtos = append(dtos, r.mapper.ToDTO(entity))
	}

	return dtos, nil
}

func (r *PostTranslationRepo) GetTranslation(ctx context.Context, postID uint32, languageCode string) (*contentV1.PostTranslation, error) {
	q := r.entClient.Client().PostTranslation.Query().
		Where(
			posttranslation.PostIDEQ(postID),
			posttranslation.LanguageCodeEQ(languageCode),
		)

	entity, err := q.Only(ctx)
	if err != nil {
		r.log.Errorf("query translation by post id and language code failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("query translation by post id and language code failed")
	}

	return r.mapper.ToDTO(entity), nil
}

func (r *PostTranslationRepo) newCreateBuilder(pt *ent.PostTranslationClient, data *contentV1.PostTranslation) *ent.PostTranslationCreate {
	now := time.Now()

	builder := pt.Create().
		SetNillablePostID(data.PostId).
		SetNillableLanguageCode(data.LanguageCode).
		SetNillableTitle(data.Title).
		SetNillableSlug(data.Slug).
		SetNillableSummary(data.Summary).
		SetNillableContent(data.Content).
		SetNillableOriginalContent(data.OriginalContent).
		SetNillableThumbnail(data.Thumbnail).
		SetNillableWordCount(data.WordCount).
		SetNillableFullPath(data.FullPath).
		SetNillableCreatedBy(data.CreatedBy).
		SetCreatedAt(now)

	if data.Seo != nil {
		builder.SetSeo(data.Seo)
	}

	return builder
}

func (r *PostTranslationRepo) BatchCreate(ctx context.Context, tx *ent.Tx, items []*contentV1.PostTranslation) error {
	if len(items) == 0 {
		return nil
	}

	builders := make([]*ent.PostTranslationCreate, 0, len(items))
	for _, data := range items {
		_ = r.PrepareTranslation(ctx, data)
		builder := r.newCreateBuilder(tx.PostTranslation, data)
		builders = append(builders, builder)
	}

	err := tx.PostTranslation.CreateBulk(builders...).Exec(ctx)
	if err != nil {
		r.log.Errorf("batch create post translations failed: %s", err.Error())
		return contentV1.ErrorInternalServerError("batch create post translations failed")
	}

	return nil
}

// PrepareTranslation 预处理帖子翻译数据，生成slug、摘要、字数等信息
func (r *PostTranslationRepo) PrepareTranslation(ctx context.Context, data *contentV1.PostTranslation) error {
	baseSlug := slug.Generate(data.GetTitle())
	slugCount, err := r.CountByBaseSlug(ctx, baseSlug)
	if err != nil {
		r.log.Errorf("count slug failed: %s", err.Error())
		return contentV1.ErrorInternalServerError("count slug failed")
	}

	if slugCount > 0 {
		baseSlug = slug.Generate(data.GetTitle()) + "-" + strconv.Itoa(int(slugCount))
	}
	data.Slug = trans.Ptr(baseSlug)

	if len(data.GetSummary()) == 0 {
		sm := summary.GenerateSummaryByRule(data.GetContent(), 100, true)
		data.Summary = trans.Ptr(sm)
	}

	counter := count.NewContentCounter(data.GetContent())
	data.WordCount = trans.Ptr(uint32(counter.RawChars()))

	return nil
}

func (r *PostTranslationRepo) CreateTranslation(ctx context.Context, data *contentV1.PostTranslation) (*contentV1.PostTranslation, error) {
	if data == nil {
		return nil, nil
	}

	_ = r.PrepareTranslation(ctx, data)

	builder := r.newCreateBuilder(r.entClient.Client().PostTranslation, data)

	var entity *ent.PostTranslation
	var err error
	if entity, err = builder.Save(ctx); err != nil {
		r.log.Errorf("create post translation failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("create post translation failed")
	}

	return r.mapper.ToDTO(entity), nil
}

func (r *PostTranslationRepo) UpdateTranslation(ctx context.Context, id uint32, data *contentV1.PostTranslation, updateMask *fieldmaskpb.FieldMask) (*contentV1.PostTranslation, error) {
	if data == nil {
		return nil, nil
	}

	builder := r.entClient.Client().PostTranslation.UpdateOneID(id)

	dto, err := r.repository.UpdateOne(ctx, builder, data, updateMask,
		func(dto *contentV1.PostTranslation) {
			builder.
				SetNillablePostID(data.PostId).
				SetNillableLanguageCode(data.LanguageCode).
				SetNillableTitle(data.Title).
				SetNillableSlug(data.Slug).
				SetNillableSummary(data.Summary).
				SetNillableContent(data.Content).
				SetNillableOriginalContent(data.OriginalContent).
				SetNillableThumbnail(data.Thumbnail).
				SetNillableWordCount(data.WordCount).
				SetNillableFullPath(data.FullPath).
				SetNillableUpdatedBy(data.UpdatedBy).
				SetUpdatedAt(time.Now())

			if data.Seo != nil {
				builder.SetSeo(data.Seo)
			}
		},
		func(s *sql.Selector) {
			s.Where(sql.EQ(posttranslation.FieldID, id))
		},
	)
	if err != nil {
		r.log.Errorf("update post translation failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("update post translation failed")
	}

	return dto, nil
}

// CountByBaseSlug counts the number of post translations with the given base slug (case-insensitive).
func (r *PostTranslationRepo) CountByBaseSlug(ctx context.Context, baseSlug string) (int64, error) {
	c, err := r.entClient.Client().PostTranslation.Query().
		Where(
			posttranslation.SlugHasPrefix(baseSlug),
		).
		Count(ctx)
	if err != nil {
		r.log.Errorf("count post translations by slug failed: %s", err.Error())
		return 0, contentV1.ErrorInternalServerError("count post translations by slug failed")
	}

	return int64(c), nil
}

// TranslationExists checks if a translation exists for the given post ID and language code.
func (r *PostTranslationRepo) TranslationExists(ctx context.Context, postId uint32, languageCode string) (bool, error) {
	c, err := r.entClient.Client().PostTranslation.Query().
		Where(
			posttranslation.PostIDEQ(postId),
			posttranslation.LanguageCodeEQ(languageCode),
		).
		Count(ctx)
	if err != nil {
		r.log.Errorf("count post translations by post id and language code failed: %s", err.Error())
		return false, contentV1.ErrorInternalServerError("count post translations by post id and language code failed")
	}

	return c > 0, nil
}

// ListAvailedLanguages lists the language codes of all translations available for the given post ID.
func (r *PostTranslationRepo) ListAvailedLanguages(ctx context.Context, postId uint32) ([]string, error) {
	entities, err := r.entClient.Client().PostTranslation.Query().
		Where(
			posttranslation.PostIDEQ(postId),
		).
		Select(posttranslation.FieldLanguageCode).
		Strings(ctx)
	if err != nil {
		r.log.Errorf("query available translation languages by post id failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("query available translation languages by post id failed")
	}

	return entities, nil
}

func (r *PostTranslationRepo) DeleteTranslation(ctx context.Context, req *contentV1.DeletePostTranslationRequest) error {
	if req.QueryBy == nil {
		return contentV1.ErrorBadRequest("invalid parameter: query_by is required")
	}

	switch req.QueryBy.(type) {
	case *contentV1.DeletePostTranslationRequest_Id:
		if req.GetId() == 0 {
			return contentV1.ErrorBadRequest("invalid parameter: id must be greater than 0")
		}

	case *contentV1.DeletePostTranslationRequest_Identifier:
		if req.GetIdentifier() == nil {
			return contentV1.ErrorBadRequest("invalid parameter: identifier is required")
		}
		if req.GetIdentifier().GetPostId() == 0 {
			return contentV1.ErrorBadRequest("invalid parameter: post_id must be greater than 0")
		}
		if len(req.GetIdentifier().GetLanguageCode()) == 0 {
			return contentV1.ErrorBadRequest("invalid parameter: language_code is required")
		}

	default:
		return contentV1.ErrorBadRequest("invalid parameter: unsupported query_by type")
	}

	builder := r.entClient.Client().PostTranslation.Delete()

	_, err := r.repository.Delete(ctx, builder, func(s *sql.Selector) {
		switch req.QueryBy.(type) {
		case *contentV1.DeletePostTranslationRequest_Id:
			id := req.GetId()
			s.Where(sql.EQ(posttranslation.FieldID, id))

		case *contentV1.DeletePostTranslationRequest_Identifier:
			identifier := req.GetIdentifier()
			s.Where(
				sql.And(
					sql.EQ(posttranslation.FieldPostID, identifier.GetPostId()),
					sql.EQ(posttranslation.FieldLanguageCode, identifier.GetLanguageCode()),
				),
			)

		default:
			return
		}
	})
	if err != nil {
		r.log.Errorf("delete post translation failed: %s", err.Error())
		return contentV1.ErrorInternalServerError("delete post translation failed")
	}

	return nil
}
