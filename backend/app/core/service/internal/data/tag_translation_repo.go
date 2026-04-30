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

func (r *TagTranslationRepo) ListTranslations(ctx context.Context, tagID uint32, locale string, viewMask *fieldmaskpb.FieldMask) ([]*contentV1.TagTranslation, error) {
	builder := r.entClient.Client().TagTranslation.Query().
		Where(
			tagtranslation.TagIDEQ(tagID),
		)

	if len(locale) > 0 {
		builder.Where(
			tagtranslation.LanguageCodeEQ(locale),
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
		r.log.Errorf("query translations by tag id failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("query translations by tag id failed")
	}

	var dtos []*contentV1.TagTranslation
	for _, entity := range entities {
		dtos = append(dtos, r.mapper.ToDTO(entity))
	}

	return dtos, nil
}

func (r *TagTranslationRepo) newCreateBuilder(tt *ent.TagTranslationClient, data *contentV1.TagTranslation) *ent.TagTranslationCreate {
	builder := tt.Create().
		SetNillableTagID(data.TagId).
		SetNillableLanguageCode(data.LanguageCode).
		SetNillableName(data.Name).
		SetNillableSlug(data.Slug).
		SetNillableDescription(data.Description).
		SetNillableCoverImage(data.CoverImage).
		SetNillableFullPath(data.FullPath).
		SetNillableCreatedBy(data.CreatedBy).
		SetCreatedAt(time.Now())

	if data.Seo != nil {
		builder.SetSeo(data.Seo)
	}

	return builder
}

func (r *TagTranslationRepo) BatchCreate(ctx context.Context, tx *ent.Tx, items []*contentV1.TagTranslation) error {
	if len(items) == 0 {
		return nil
	}

	builders := make([]*ent.TagTranslationCreate, 0, len(items))
	for _, data := range items {
		_ = r.PrepareTranslation(ctx, data)

		builder := r.newCreateBuilder(r.entClient.Client().TagTranslation, data)

		builders = append(builders, builder)
	}

	err := tx.TagTranslation.CreateBulk(builders...).Exec(ctx)
	if err != nil {
		r.log.Errorf("batch create tag translations failed: %s", err.Error())
		return contentV1.ErrorInternalServerError("batch create tag translations failed")
	}

	return nil
}

func (r *TagTranslationRepo) CreateTranslation(ctx context.Context, data *contentV1.TagTranslation) (*contentV1.TagTranslation, error) {

	_ = r.PrepareTranslation(ctx, data)

	builder := r.newCreateBuilder(r.entClient.Client().TagTranslation, data)

	entity, err := builder.Save(ctx)
	if err != nil {
		r.log.Errorf("create tag translation failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("create tag translation failed")
	}

	return r.mapper.ToDTO(entity), nil
}

func (r *TagTranslationRepo) UpdateTranslation(ctx context.Context, id uint32, data *contentV1.TagTranslation, updateMask *fieldmaskpb.FieldMask) (*contentV1.TagTranslation, error) {
	if data == nil {
		return nil, nil
	}

	builder := r.entClient.Client().TagTranslation.UpdateOneID(id)

	dto, err := r.repository.UpdateOne(ctx, builder, data, updateMask,
		func(dto *contentV1.TagTranslation) {
			builder.
				SetNillableTagID(data.TagId).
				SetNillableLanguageCode(data.LanguageCode).
				SetNillableName(data.Name).
				SetNillableSlug(data.Slug).
				SetNillableDescription(data.Description).
				SetNillableCoverImage(data.CoverImage).
				SetNillableFullPath(data.FullPath).
				SetNillableUpdatedBy(data.UpdatedBy).
				SetUpdatedAt(time.Now())

			if data.Seo != nil {
				builder.SetSeo(data.Seo)
			}
		},
		func(s *sql.Selector) {
			s.Where(sql.EQ(tagtranslation.FieldID, id))
		},
	)
	if err != nil {
		r.log.Errorf("update tag translation failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("update tag translation failed")
	}

	return dto, nil
}

func (r *TagTranslationRepo) CountByBaseSlug(ctx context.Context, baseSlug string) (int64, error) {
	c, err := r.entClient.Client().TagTranslation.Query().
		Where(
			tagtranslation.SlugHasPrefix(baseSlug),
		).
		Count(ctx)
	if err != nil {
		r.log.Errorf("count tag translations by slug failed: %s", err.Error())
		return 0, contentV1.ErrorInternalServerError("count tag translations by slug failed")
	}

	return int64(c), nil
}

// TranslationExists checks if a translation exists for the given tag ID and language code.
func (r *TagTranslationRepo) TranslationExists(ctx context.Context, tagId uint32, languageCode string) (bool, error) {
	c, err := r.entClient.Client().TagTranslation.Query().
		Where(
			tagtranslation.TagIDEQ(tagId),
			tagtranslation.LanguageCodeEQ(languageCode),
		).
		Count(ctx)
	if err != nil {
		r.log.Errorf("count tag translations by tag id and language code failed: %s", err.Error())
		return false, contentV1.ErrorInternalServerError("count tag translations by tag id and language code failed")
	}

	return c > 0, nil
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

func (r *TagTranslationRepo) GetTranslation(ctx context.Context, tagId uint32, languageCode string) (*contentV1.TagTranslation, error) {
	entity, err := r.entClient.Client().TagTranslation.Query().
		Where(
			tagtranslation.TagIDEQ(tagId),
			tagtranslation.LanguageCodeEQ(languageCode),
		).
		Only(ctx)
	if err != nil {
		r.log.Errorf("query tag translation by tag id and language code failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("query tag translation by tag id and language code failed")
	}

	return r.mapper.ToDTO(entity), nil
}

func (r *TagTranslationRepo) DeleteTranslation(ctx context.Context, req *contentV1.DeleteTagTranslationRequest) error {
	if req.QueryBy == nil {
		return contentV1.ErrorBadRequest("invalid parameter: query_by is required")
	}

	switch req.QueryBy.(type) {
	case *contentV1.DeleteTagTranslationRequest_Id:
		if req.GetId() == 0 {
			return contentV1.ErrorBadRequest("invalid parameter: id must be greater than 0")
		}

	case *contentV1.DeleteTagTranslationRequest_Identifier:
		if req.GetIdentifier() == nil {
			return contentV1.ErrorBadRequest("invalid parameter: identifier is required")
		}
		if req.GetIdentifier().GetTagId() == 0 {
			return contentV1.ErrorBadRequest("invalid parameter: tag_id must be greater than 0")
		}
		if len(req.GetIdentifier().GetLanguageCode()) == 0 {
			return contentV1.ErrorBadRequest("invalid parameter: language_code is required")
		}

	default:
		return contentV1.ErrorBadRequest("invalid parameter: unsupported query_by type")
	}

	builder := r.entClient.Client().TagTranslation.Delete()

	_, err := r.repository.Delete(ctx, builder, func(s *sql.Selector) {
		switch req.QueryBy.(type) {
		case *contentV1.DeleteTagTranslationRequest_Id:
			id := req.GetId()
			s.Where(sql.EQ(tagtranslation.FieldID, id))

		case *contentV1.DeleteTagTranslationRequest_Identifier:
			identifier := req.GetIdentifier()
			s.Where(
				sql.And(
					sql.EQ(tagtranslation.FieldTagID, identifier.GetTagId()),
					sql.EQ(tagtranslation.FieldLanguageCode, identifier.GetLanguageCode()),
				),
			)

		default:
			return
		}
	})
	if err != nil {
		r.log.Errorf("delete tag translation failed: %s", err.Error())
		return contentV1.ErrorInternalServerError("delete tag translation failed")
	}

	return nil
}

func (r *TagTranslationRepo) PrepareTranslation(ctx context.Context, data *contentV1.TagTranslation) error {
	baseSlug := slug.Generate(data.GetName())
	slugCount, err := r.CountByBaseSlug(ctx, baseSlug)
	if err != nil {
		r.log.Errorf("count slug failed: %s", err.Error())
		return contentV1.ErrorInternalServerError("count slug failed")
	}

	if slugCount > 0 {
		baseSlug = slug.Generate(data.GetName()) + "-" + strconv.Itoa(int(slugCount))
	}

	data.Slug = trans.Ptr(baseSlug)

	return nil
}
