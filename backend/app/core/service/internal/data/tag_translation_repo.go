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
	delBuilder := tx.TagTranslation.Delete().Where(tagtranslation.TagIDEQ(tagID))
	// 租户作用域：仅清除本租户翻译，避免跨租户删他人翻译（按 hasTenant 条件加）
	if tid, hasTenant := maybeTenantFromViewer(ctx); hasTenant {
		delBuilder.Where(tagtranslation.TenantIDEQ(tid))
	}
	if _, err := delBuilder.Exec(ctx); err != nil {
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
		Order(ent.Asc(tagtranslation.FieldID)).
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

// ListTranslationsWithFallback 优先返回指定语言的译文；该语言缺译时返回全部译文，
// 供前端按"匹配当前语言→回退第一条"兜底（只回传空数组会让前端无 fallback 可用）。
func (r *TagTranslationRepo) ListTranslationsWithFallback(ctx context.Context, tagID uint32, locale string, viewMask *fieldmaskpb.FieldMask) ([]*contentV1.TagTranslation, error) {
	translations, err := r.ListTranslations(ctx, tagID, locale, viewMask)
	if err != nil {
		return nil, err
	}
	if len(translations) == 0 && locale != "" {
		return r.ListTranslations(ctx, tagID, "", viewMask)
	}
	return translations, nil
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
		_ = r.PrepareTranslation(ctx, tx.TagTranslation, data)

		builder := r.newCreateBuilder(tx.TagTranslation, data)

		builders = append(builders, builder)
	}

	err := tx.TagTranslation.CreateBulk(builders...).Exec(ctx)
	if err != nil {
		r.log.Errorf("batch create tag translations failed: %s", err.Error())
		return contentV1.ErrorInternalServerError("batch create tag translations failed")
	}

	return nil
}

// UpsertTranslations 按 (tag_id, language_code) 逐条 upsert：
// 已存在则原行就地更新（保留 created_by/created_at 与未提交字段），
// 不存在则插入；同语言历史重复行仅保留 id 最小的一条并更新，其余删除。
// Update 语义据此从"整表替换"改为按语言合并，删除译文请走 DeleteTranslation。
func (r *TagTranslationRepo) UpsertTranslations(ctx context.Context, tx *ent.Tx, tagID uint32, items []*contentV1.TagTranslation) error {
	if len(items) == 0 {
		return nil
	}

	callerUserID, hasUser := viewerUserIDFromContext(ctx)

	for _, data := range items {
		if data == nil || len(data.GetLanguageCode()) == 0 {
			return contentV1.ErrorBadRequest("invalid parameter: language_code is required")
		}
		data.TagId = trans.Ptr(tagID)

		existings, err := tx.TagTranslation.Query().
			Where(
				tagtranslation.TagIDEQ(tagID),
				tagtranslation.LanguageCodeEQ(data.GetLanguageCode()),
			).
			Order(ent.Asc(tagtranslation.FieldID)).
			All(ctx)
		if err != nil {
			r.log.Errorf("query tag translations for upsert failed: %s", err.Error())
			return contentV1.ErrorInternalServerError("query tag translations for upsert failed")
		}

		if len(existings) == 0 {
			_ = r.PrepareTranslation(ctx, tx.TagTranslation, data)
			builder := r.newCreateBuilder(tx.TagTranslation, data)
			if _, err = builder.Save(ctx); err != nil {
				r.log.Errorf("upsert insert tag translation failed: %s", err.Error())
				return contentV1.ErrorInternalServerError("upsert insert tag translation failed")
			}
			continue
		}

		keep := existings[0]

		// 名称未变时保留原 slug，避免每次重发都因前缀计数叠加后缀
		nameChanged := (keep.Name == nil && data.GetName() != "") ||
			(keep.Name != nil && *keep.Name != data.GetName())
		if nameChanged {
			_ = r.PrepareTranslation(ctx, tx.TagTranslation, data)
		}

		upd := tx.TagTranslation.UpdateOneID(keep.ID)
		if tid, hasTenant := maybeTenantFromViewer(ctx); hasTenant {
			upd.Where(tagtranslation.TenantIDEQ(tid))
		}
		upd.
			SetNillableName(data.Name).
			SetNillableSlug(data.Slug).
			SetNillableDescription(data.Description).
			SetNillableCoverImage(data.CoverImage).
			SetNillableFullPath(data.FullPath).
			SetUpdatedAt(time.Now())
		if hasUser {
			upd.SetUpdatedBy(callerUserID)
		}
		if data.Seo != nil {
			upd.SetSeo(data.Seo)
		}
		if _, err = upd.Save(ctx); err != nil {
			r.log.Errorf("upsert update tag translation failed: %s", err.Error())
			return contentV1.ErrorInternalServerError("upsert update tag translation failed")
		}

		if len(existings) > 1 {
			dupeIDs := make([]uint32, 0, len(existings)-1)
			for _, e := range existings[1:] {
				dupeIDs = append(dupeIDs, e.ID)
			}
			if _, err = tx.TagTranslation.Delete().Where(tagtranslation.IDIn(dupeIDs...)).Exec(ctx); err != nil {
				r.log.Errorf("clean duplicated tag translations failed: %s", err.Error())
				return contentV1.ErrorInternalServerError("clean duplicated tag translations failed")
			}
		}
	}

	return nil
}

func (r *TagTranslationRepo) CreateTranslation(ctx context.Context, data *contentV1.TagTranslation) (*contentV1.TagTranslation, error) {

	_ = r.PrepareTranslation(ctx, r.entClient.Client().TagTranslation, data)

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
	// 租户作用域：仅更新本租户翻译，避免跨租户改他人翻译（按 hasTenant 条件加）
	if tid, hasTenant := maybeTenantFromViewer(ctx); hasTenant {
		builder.Where(tagtranslation.TenantIDEQ(tid))
	}
	callerUserID, hasUser := viewerUserIDFromContext(ctx)

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
				SetUpdatedAt(time.Now())

			// updated_by 强制由服务端 viewer context 推导，忽略客户端传入值
			if hasUser {
				builder.SetUpdatedBy(callerUserID)
			}

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
	return r.countByBaseSlug(ctx, r.entClient.Client().TagTranslation, baseSlug)
}

func (r *TagTranslationRepo) countByBaseSlug(ctx context.Context, tt *ent.TagTranslationClient, baseSlug string) (int64, error) {
	c, err := tt.Query().
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
	// 取 id 最小的一条保证结果确定（防历史重复行）；缺译文不是错误，交由调用方按 nil 处理
	entity, err := r.entClient.Client().TagTranslation.Query().
		Where(
			tagtranslation.TagIDEQ(tagId),
			tagtranslation.LanguageCodeEQ(languageCode),
		).
		Order(ent.Asc(tagtranslation.FieldID)).
		First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
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
	// 租户作用域：仅删除本租户翻译，避免跨租户删他人翻译（按 hasTenant 条件加）
	if tid, hasTenant := maybeTenantFromViewer(ctx); hasTenant {
		builder.Where(tagtranslation.TenantIDEQ(tid))
	}

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

// tt 必须传与调用方一致的事务/非事务客户端（见 PostTranslationRepo.PrepareTranslation 注释）。
func (r *TagTranslationRepo) PrepareTranslation(ctx context.Context, tt *ent.TagTranslationClient, data *contentV1.TagTranslation) error {
	baseSlug := slug.Generate(data.GetName())
	slugCount, err := r.countByBaseSlug(ctx, tt, baseSlug)
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
