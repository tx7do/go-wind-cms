package data

import (
	"context"
	"strconv"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/go-kratos/kratos/v2/log"
	entCrud "github.com/tx7do/go-crud/entgo"
	"github.com/tx7do/go-utils/copierutil"
	"github.com/tx7do/go-utils/mapper"
	"github.com/tx7do/go-utils/slug"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

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

func (r *CategoryTranslationRepo) ListTranslations(ctx context.Context, categoryID uint32, locale string, viewMask *fieldmaskpb.FieldMask) ([]*contentV1.CategoryTranslation, error) {
	builder := r.entClient.Client().CategoryTranslation.Query().
		Where(
			categorytranslation.CategoryIDEQ(categoryID),
		)

	if len(locale) > 0 {
		builder.Where(
			categorytranslation.LanguageCodeEQ(locale),
		)
	}

	if viewMask != nil {
		selectSelector, err := r.repository.BuildSelector(viewMask.GetPaths())
		if err != nil {
			r.log.Errorf("build category translation selector failed: %s", err.Error())
			return nil, contentV1.ErrorInternalServerError("build category translation selector failed")
		}
		if selectSelector != nil {
			builder.Modify(selectSelector)
		}
	}

	entities, err := builder.
		Order(ent.Asc(categorytranslation.FieldID)).
		All(ctx)
	if err != nil {
		r.log.Errorf("query translations list failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("query translations list failed")
	}

	var dtos []*contentV1.CategoryTranslation
	for _, entity := range entities {
		dtos = append(dtos, r.mapper.ToDTO(entity))
	}

	return dtos, nil
}

// ListTranslationsWithFallback 优先返回指定语言的译文；该语言缺译时返回全部译文，
// 供前端按"匹配当前语言→回退第一条"兜底（只回传空数组会让前端无 fallback 可用）。
func (r *CategoryTranslationRepo) ListTranslationsWithFallback(ctx context.Context, categoryID uint32, locale string, viewMask *fieldmaskpb.FieldMask) ([]*contentV1.CategoryTranslation, error) {
	translations, err := r.ListTranslations(ctx, categoryID, locale, viewMask)
	if err != nil {
		return nil, err
	}
	if len(translations) == 0 && locale != "" {
		return r.ListTranslations(ctx, categoryID, "", viewMask)
	}
	return translations, nil
}

// ct 必须传与调用方一致的事务/非事务客户端，事务内用非事务连接读同一张表
// 在 sqlite 内存库下会死锁（见 PostTranslationRepo.PrepareTranslation 注释）。
func (r *CategoryTranslationRepo) newCreateBuilder(ct *ent.CategoryTranslationClient, data *contentV1.CategoryTranslation) *ent.CategoryTranslationCreate {
	now := time.Now()

	builder := ct.Create().
		SetNillableCategoryID(data.CategoryId).
		SetNillableLanguageCode(data.LanguageCode).
		SetNillableName(data.Name).
		SetNillableSlug(data.Slug).
		SetNillableDescription(data.Description).
		SetNillableThumbnail(data.Thumbnail).
		SetNillableCoverImage(data.CoverImage).
		SetNillableFullPath(data.FullPath).
		SetNillableCreatedBy(data.CreatedBy).
		SetCreatedAt(now)

	if data.Seo != nil {
		builder.SetSeo(data.Seo)
	}

	return builder
}

func (r *CategoryTranslationRepo) BatchCreate(ctx context.Context, tx *ent.Tx, items []*contentV1.CategoryTranslation) error {
	if len(items) == 0 {
		return nil
	}

	builders := make([]*ent.CategoryTranslationCreate, 0, len(items))
	for _, data := range items {
		_ = r.PrepareTranslation(ctx, tx.CategoryTranslation, data)

		builder := r.newCreateBuilder(tx.CategoryTranslation, data)

		builders = append(builders, builder)
	}

	err := tx.CategoryTranslation.CreateBulk(builders...).Exec(ctx)
	if err != nil {
		r.log.Errorf("batch create category translations failed: %s", err.Error())
		return contentV1.ErrorInternalServerError("batch create category translations failed")
	}

	return nil
}

// UpsertTranslations 按 (category_id, language_code) 逐条 upsert：
// 已存在则原行就地更新（保留 created_by/created_at 与未提交字段），
// 不存在则插入；同语言历史重复行仅保留 id 最小的一条并更新，其余删除。
// Update 语义据此从"整表替换"改为按语言合并，删除译文请走 DeleteTranslation。
func (r *CategoryTranslationRepo) UpsertTranslations(ctx context.Context, tx *ent.Tx, categoryID uint32, items []*contentV1.CategoryTranslation) error {
	if len(items) == 0 {
		return nil
	}

	callerUserID, hasUser := viewerUserIDFromContext(ctx)

	for _, data := range items {
		if data == nil || len(data.GetLanguageCode()) == 0 {
			return contentV1.ErrorBadRequest("invalid parameter: language_code is required")
		}
		data.CategoryId = trans.Ptr(categoryID)

		existings, err := tx.CategoryTranslation.Query().
			Where(
				categorytranslation.CategoryIDEQ(categoryID),
				categorytranslation.LanguageCodeEQ(data.GetLanguageCode()),
			).
			Order(ent.Asc(categorytranslation.FieldID)).
			All(ctx)
		if err != nil {
			r.log.Errorf("query category translations for upsert failed: %s", err.Error())
			return contentV1.ErrorInternalServerError("query category translations for upsert failed")
		}

		if len(existings) == 0 {
			_ = r.PrepareTranslation(ctx, tx.CategoryTranslation, data)
			builder := r.newCreateBuilder(tx.CategoryTranslation, data)
			if _, err = builder.Save(ctx); err != nil {
				r.log.Errorf("upsert insert category translation failed: %s", err.Error())
				return contentV1.ErrorInternalServerError("upsert insert category translation failed")
			}
			continue
		}

		keep := existings[0]

		// 名称未变时保留原 slug，避免每次重发都因前缀计数叠加后缀
		nameChanged := (keep.Name == nil && data.GetName() != "") ||
			(keep.Name != nil && *keep.Name != data.GetName())
		if nameChanged {
			_ = r.PrepareTranslation(ctx, tx.CategoryTranslation, data)
		}

		upd := tx.CategoryTranslation.UpdateOneID(keep.ID)
		if tid, hasTenant := maybeTenantFromViewer(ctx); hasTenant {
			upd.Where(categorytranslation.TenantIDEQ(tid))
		}
		upd.
			SetNillableName(data.Name).
			SetNillableSlug(data.Slug).
			SetNillableDescription(data.Description).
			SetNillableThumbnail(data.Thumbnail).
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
			r.log.Errorf("upsert update category translation failed: %s", err.Error())
			return contentV1.ErrorInternalServerError("upsert update category translation failed")
		}

		if len(existings) > 1 {
			dupeIDs := make([]uint32, 0, len(existings)-1)
			for _, e := range existings[1:] {
				dupeIDs = append(dupeIDs, e.ID)
			}
			if _, err = tx.CategoryTranslation.Delete().Where(categorytranslation.IDIn(dupeIDs...)).Exec(ctx); err != nil {
				r.log.Errorf("clean duplicated category translations failed: %s", err.Error())
				return contentV1.ErrorInternalServerError("clean duplicated category translations failed")
			}
		}
	}

	return nil
}

func (r *CategoryTranslationRepo) CreateTranslation(ctx context.Context, data *contentV1.CategoryTranslation) (*contentV1.CategoryTranslation, error) {
	_ = r.PrepareTranslation(ctx, r.entClient.Client().CategoryTranslation, data)

	builder := r.newCreateBuilder(r.entClient.Client().CategoryTranslation, data)

	entity, err := builder.Save(ctx)
	if err != nil {
		r.log.Errorf("create category translation failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("create category translation failed")
	}

	return r.mapper.ToDTO(entity), nil
}

func (r *CategoryTranslationRepo) UpdateTranslation(ctx context.Context, id uint32, data *contentV1.CategoryTranslation, updateMask *fieldmaskpb.FieldMask) (*contentV1.CategoryTranslation, error) {
	if data == nil {
		return nil, nil
	}

	builder := r.entClient.Client().CategoryTranslation.UpdateOneID(id)
	// 租户作用域：仅更新本租户翻译，避免跨租户改他人翻译（按 hasTenant 条件加）
	if tid, hasTenant := maybeTenantFromViewer(ctx); hasTenant {
		builder.Where(categorytranslation.TenantIDEQ(tid))
	}
	callerUserID, hasUser := viewerUserIDFromContext(ctx)

	dto, err := r.repository.UpdateOne(ctx, builder, data, updateMask,
		func(dto *contentV1.CategoryTranslation) {
			builder.
				SetNillableCategoryID(dto.CategoryId).
				SetNillableLanguageCode(dto.LanguageCode).
				SetNillableName(dto.Name).
				SetNillableSlug(dto.Slug).
				SetNillableDescription(dto.Description).
				SetNillableThumbnail(dto.Thumbnail).
				SetNillableCoverImage(dto.CoverImage).
				SetNillableFullPath(dto.FullPath).
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
			s.Where(sql.EQ(categorytranslation.FieldID, id))
		},
	)
	if err != nil {
		r.log.Errorf("update category translation failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("update category translation failed")
	}

	return dto, nil
}

func (r *CategoryTranslationRepo) CountByBaseSlug(ctx context.Context, baseSlug string) (int64, error) {
	return r.countByBaseSlug(ctx, r.entClient.Client().CategoryTranslation, baseSlug)
}

func (r *CategoryTranslationRepo) countByBaseSlug(ctx context.Context, ct *ent.CategoryTranslationClient, baseSlug string) (int64, error) {
	c, err := ct.Query().
		Where(
			categorytranslation.SlugHasPrefix(baseSlug),
		).
		Count(ctx)
	if err != nil {
		r.log.Errorf("count category translations by slug failed: %s", err.Error())
		return 0, contentV1.ErrorInternalServerError("count category translations by slug failed")
	}

	return int64(c), nil
}

// TranslationExists checks if a translation exists for the given category ID and language code.
func (r *CategoryTranslationRepo) TranslationExists(ctx context.Context, categoryId uint32, languageCode string) (bool, error) {
	c, err := r.entClient.Client().CategoryTranslation.Query().
		Where(
			categorytranslation.CategoryIDEQ(categoryId),
			categorytranslation.LanguageCodeEQ(languageCode),
		).
		Count(ctx)
	if err != nil {
		r.log.Errorf("count category translations by category id and language code failed: %s", err.Error())
		return false, contentV1.ErrorInternalServerError("count category translations by category id and language code failed")
	}

	return c > 0, nil
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

func (r *CategoryTranslationRepo) GetTranslation(ctx context.Context, categoryId uint32, languageCode string) (*contentV1.CategoryTranslation, error) {
	// 取 id 最小的一条保证结果确定（防历史重复行）；缺译文不是错误，交由调用方按 nil 处理
	entity, err := r.entClient.Client().CategoryTranslation.Query().
		Where(
			categorytranslation.CategoryIDEQ(categoryId),
			categorytranslation.LanguageCodeEQ(languageCode),
		).
		Order(ent.Asc(categorytranslation.FieldID)).
		First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		r.log.Errorf("query category translation by category id and language code failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("query category translation by category id and language code failed")
	}

	return r.mapper.ToDTO(entity), nil
}

func (r *CategoryTranslationRepo) DeleteTranslation(ctx context.Context, req *contentV1.DeleteCategoryTranslationRequest) error {
	if req == nil {
		return nil
	}

	switch req.QueryBy.(type) {
	case *contentV1.DeleteCategoryTranslationRequest_Id:
		if req.GetId() == 0 {
			return contentV1.ErrorBadRequest("id is required for delete category translation")
		}
	case *contentV1.DeleteCategoryTranslationRequest_Identifier:
		// do nothing, use category id and language code to delete
	default:
		return contentV1.ErrorBadRequest("invalid query by for delete category translation")
	}

	builder := r.entClient.Client().CategoryTranslation.Delete()
	// 租户作用域：仅删除本租户翻译，避免跨租户删他人翻译（按 hasTenant 条件加）
	if tid, hasTenant := maybeTenantFromViewer(ctx); hasTenant {
		builder.Where(categorytranslation.TenantIDEQ(tid))
	}

	_, err := r.repository.Delete(ctx, builder, func(s *sql.Selector) {
		switch req.QueryBy.(type) {
		case *contentV1.DeleteCategoryTranslationRequest_Id:
			id := req.GetId()
			s.Where(sql.EQ(categorytranslation.FieldID, id))

		case *contentV1.DeleteCategoryTranslationRequest_Identifier:
			identifier := req.GetIdentifier()
			s.Where(
				sql.And(
					sql.EQ(categorytranslation.FieldCategoryID, identifier.GetCategoryId()),
					sql.EQ(categorytranslation.FieldLanguageCode, identifier.GetLanguageCode()),
				),
			)

		default:
			return
		}
	})
	if err != nil {
		r.log.Errorf("delete category translation failed: %s", err.Error())
		return contentV1.ErrorInternalServerError("delete category translation failed")
	}

	return nil
}

func (r *CategoryTranslationRepo) PrepareTranslation(ctx context.Context, ct *ent.CategoryTranslationClient, data *contentV1.CategoryTranslation) error {
	baseSlug := slug.Generate(data.GetName())
	slugCount, err := r.countByBaseSlug(ctx, ct, baseSlug)
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
