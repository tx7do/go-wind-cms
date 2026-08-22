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

// pt 必须传与调用方一致的事务/非事务客户端，事务内用非事务连接读同一张表
// 在 sqlite 内存库下会死锁（见 PostTranslationRepo.PrepareTranslation 注释）。
func (r *PageTranslationRepo) newCreateBuilder(pt *ent.PageTranslationClient, data *contentV1.PageTranslation) *ent.PageTranslationCreate {
	now := time.Now()

	builder := pt.Create().
		SetNillablePageID(data.PageId).
		SetNillableLanguageCode(data.LanguageCode).
		SetNillableTitle(data.Title).
		SetNillableSlug(data.Slug).
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

func (r *PageTranslationRepo) BatchCreate(ctx context.Context, tx *ent.Tx, items []*contentV1.PageTranslation) error {
	if len(items) == 0 {
		return nil
	}

	builders := make([]*ent.PageTranslationCreate, 0, len(items))
	for _, data := range items {
		_ = r.PrepareTranslation(ctx, tx.PageTranslation, data)

		builder := r.newCreateBuilder(tx.PageTranslation, data)

		builders = append(builders, builder)
	}

	err := tx.PageTranslation.CreateBulk(builders...).Exec(ctx)
	if err != nil {
		r.log.Errorf("batch create page translations failed: %s", err.Error())
		return contentV1.ErrorInternalServerError("batch create page translations failed")
	}

	return nil
}

// UpsertTranslations 按 (page_id, language_code) 逐条 upsert：
// 已存在则原行就地更新（保留 created_by/created_at 与未提交字段），
// 不存在则插入；同语言历史重复行仅保留 id 最小的一条并更新，其余删除。
// Update 语义据此从"整表替换"改为按语言合并，删除译文请走 DeleteTranslation。
func (r *PageTranslationRepo) UpsertTranslations(ctx context.Context, tx *ent.Tx, pageID uint32, items []*contentV1.PageTranslation) error {
	if len(items) == 0 {
		return nil
	}

	callerUserID, hasUser := viewerUserIDFromContext(ctx)

	for _, data := range items {
		if data == nil || len(data.GetLanguageCode()) == 0 {
			return contentV1.ErrorBadRequest("invalid parameter: language_code is required")
		}
		data.PageId = trans.Ptr(pageID)

		existings, err := tx.PageTranslation.Query().
			Where(
				pagetranslation.PageIDEQ(pageID),
				pagetranslation.LanguageCodeEQ(data.GetLanguageCode()),
			).
			Order(ent.Asc(pagetranslation.FieldID)).
			All(ctx)
		if err != nil {
			r.log.Errorf("query page translations for upsert failed: %s", err.Error())
			return contentV1.ErrorInternalServerError("query page translations for upsert failed")
		}

		if len(existings) == 0 {
			_ = r.PrepareTranslation(ctx, tx.PageTranslation, data)
			builder := r.newCreateBuilder(tx.PageTranslation, data)
			if _, err = builder.Save(ctx); err != nil {
				r.log.Errorf("upsert insert page translation failed: %s", err.Error())
				return contentV1.ErrorInternalServerError("upsert insert page translation failed")
			}
			continue
		}

		keep := existings[0]

		// 标题未变时保留原 slug，避免每次重发都因前缀计数叠加后缀
		titleChanged := (keep.Title == nil && data.GetTitle() != "") ||
			(keep.Title != nil && *keep.Title != data.GetTitle())
		if titleChanged {
			_ = r.PrepareTranslation(ctx, tx.PageTranslation, data)
		}

		upd := tx.PageTranslation.UpdateOneID(keep.ID)
		if tid, hasTenant := maybeTenantFromViewer(ctx); hasTenant {
			upd.Where(pagetranslation.TenantIDEQ(tid))
		}
		upd.
			SetNillableTitle(data.Title).
			SetNillableSlug(data.Slug).
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
			r.log.Errorf("upsert update page translation failed: %s", err.Error())
			return contentV1.ErrorInternalServerError("upsert update page translation failed")
		}

		if len(existings) > 1 {
			dupeIDs := make([]uint32, 0, len(existings)-1)
			for _, e := range existings[1:] {
				dupeIDs = append(dupeIDs, e.ID)
			}
			if _, err = tx.PageTranslation.Delete().Where(pagetranslation.IDIn(dupeIDs...)).Exec(ctx); err != nil {
				r.log.Errorf("clean duplicated page translations failed: %s", err.Error())
				return contentV1.ErrorInternalServerError("clean duplicated page translations failed")
			}
		}
	}

	return nil
}

func (r *PageTranslationRepo) CreateTranslation(ctx context.Context, data *contentV1.PageTranslation) (*contentV1.PageTranslation, error) {
	if err := r.PrepareTranslation(ctx, r.entClient.Client().PageTranslation, data); err != nil {
		return nil, err
	}

	builder := r.newCreateBuilder(r.entClient.Client().PageTranslation, data)

	entity, err := builder.Save(ctx)
	if err != nil {
		r.log.Errorf("create page translation failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("create page translation failed")
	}

	dto := r.mapper.ToDTO(entity)

	return dto, nil
}

func (r *PageTranslationRepo) UpdateTranslation(ctx context.Context, id uint32, data *contentV1.PageTranslation, updateMask *fieldmaskpb.FieldMask) (*contentV1.PageTranslation, error) {
	if data == nil {
		return nil, nil
	}

	builder := r.entClient.Client().PageTranslation.UpdateOneID(id)
	// 租户作用域：仅更新本租户翻译，避免跨租户改他人翻译（按 hasTenant 条件加）
	if tid, hasTenant := maybeTenantFromViewer(ctx); hasTenant {
		builder.Where(pagetranslation.TenantIDEQ(tid))
	}
	callerUserID, hasUser := viewerUserIDFromContext(ctx)

	dto, err := r.repository.UpdateOne(ctx, builder, data, updateMask,
		func(dto *contentV1.PageTranslation) {
			builder.
				SetNillableTitle(data.Title).
				SetNillableThumbnail(data.Thumbnail).
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
			s.Where(sql.EQ(pagetranslation.FieldID, id))
		},
	)
	if err != nil {
		r.log.Errorf("update page translation failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("update page translation failed")
	}

	return dto, nil
}

func (r *PageTranslationRepo) CountByBaseSlug(ctx context.Context, baseSlug string) (int64, error) {
	return r.countByBaseSlug(ctx, r.entClient.Client().PageTranslation, baseSlug)
}

func (r *PageTranslationRepo) countByBaseSlug(ctx context.Context, pt *ent.PageTranslationClient, baseSlug string) (int64, error) {
	c, err := pt.Query().
		Where(
			pagetranslation.SlugHasPrefix(baseSlug),
		).
		Count(ctx)
	if err != nil {
		r.log.Errorf("count page translations by slug failed: %s", err.Error())
		return 0, contentV1.ErrorInternalServerError("count page translations by slug failed")
	}

	return int64(c), nil
}

// TranslationExists checks if a translation exists for the given page ID and language code.
func (r *PageTranslationRepo) TranslationExists(ctx context.Context, pageId uint32, languageCode string) (bool, error) {
	c, err := r.entClient.Client().PageTranslation.Query().
		Where(
			pagetranslation.PageIDEQ(pageId),
			pagetranslation.LanguageCodeEQ(languageCode),
		).
		Count(ctx)
	if err != nil {
		r.log.Errorf("count page translations by page id and language code failed: %s", err.Error())
		return false, contentV1.ErrorInternalServerError("count page translations by page id and language code failed")
	}

	return c > 0, nil
}

// ListAvailedLanguages lists the language codes of all translations available for the given page ID.
func (r *PageTranslationRepo) ListAvailedLanguages(ctx context.Context, pageId uint32) ([]string, error) {
	entities, err := r.entClient.Client().PageTranslation.Query().
		Where(
			pagetranslation.PageIDEQ(pageId),
		).
		Select(pagetranslation.FieldLanguageCode).
		Strings(ctx)
	if err != nil {
		r.log.Errorf("query available translation languages by page id failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("query available translation languages by page id failed")
	}

	return entities, nil
}

func (r *PageTranslationRepo) PrepareTranslation(ctx context.Context, pt *ent.PageTranslationClient, data *contentV1.PageTranslation) error {
	baseSlug := slug.Generate(data.GetTitle())
	slugCount, err := r.countByBaseSlug(ctx, pt, baseSlug)
	if err != nil {
		r.log.Errorf("count slug failed: %s", err.Error())
		return contentV1.ErrorInternalServerError("count slug failed")
	}

	if slugCount > 0 {
		baseSlug = slug.Generate(data.GetTitle()) + "-" + strconv.Itoa(int(slugCount))
	}
	data.Slug = trans.Ptr(baseSlug)

	return nil
}

func (r *PageTranslationRepo) GetTranslation(ctx context.Context, pageId uint32, languageCode string) (*contentV1.PageTranslation, error) {
	// 同 (page_id, language_code) 历史上可能存在重复行（旧 Update 直插遗留），
	// 取 id 最小的一条保证结果确定；缺译文不是错误，交由调用方按 nil 处理
	entity, err := r.entClient.Client().PageTranslation.Query().
		Where(
			pagetranslation.PageIDEQ(pageId),
			pagetranslation.LanguageCodeEQ(languageCode),
		).
		Order(ent.Asc(pagetranslation.FieldID)).
		First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		r.log.Errorf("query page translation by page id and language code failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("query page translation by page id and language code failed")
	}

	dto := r.mapper.ToDTO(entity)

	return dto, nil
}

func (r *PageTranslationRepo) DeleteTranslation(ctx context.Context, req *contentV1.DeletePageTranslationRequest) error {
	if req.QueryBy == nil {
		return contentV1.ErrorBadRequest("invalid parameter: query_by is required")
	}

	switch req.QueryBy.(type) {
	case *contentV1.DeletePageTranslationRequest_Id:
		if req.GetId() == 0 {
			return contentV1.ErrorBadRequest("invalid parameter: id must be greater than 0")
		}

	case *contentV1.DeletePageTranslationRequest_Identifier:
		if req.GetIdentifier() == nil {
			return contentV1.ErrorBadRequest("invalid parameter: identifier is required")
		}
		if req.GetIdentifier().GetPageId() == 0 {
			return contentV1.ErrorBadRequest("invalid parameter: page_id must be greater than 0")
		}
		if len(req.GetIdentifier().GetLanguageCode()) == 0 {
			return contentV1.ErrorBadRequest("invalid parameter: language_code is required")
		}

	default:
		return contentV1.ErrorBadRequest("invalid parameter: unsupported query_by type")
	}

	builder := r.entClient.Client().PageTranslation.Delete()
	// 租户作用域：仅删除本租户翻译，避免跨租户删他人翻译（按 hasTenant 条件加）
	if tid, hasTenant := maybeTenantFromViewer(ctx); hasTenant {
		builder.Where(pagetranslation.TenantIDEQ(tid))
	}

	_, err := r.repository.Delete(ctx, builder, func(s *sql.Selector) {
		switch req.QueryBy.(type) {
		case *contentV1.DeletePageTranslationRequest_Id:
			id := req.GetId()
			s.Where(sql.EQ(pagetranslation.FieldID, id))

		case *contentV1.DeletePageTranslationRequest_Identifier:
			identifier := req.GetIdentifier()
			s.Where(
				sql.And(
					sql.EQ(pagetranslation.FieldPageID, identifier.GetPageId()),
					sql.EQ(pagetranslation.FieldLanguageCode, identifier.GetLanguageCode()),
				),
			)

		default:
			return
		}
	})
	if err != nil {
		r.log.Errorf("delete page translation failed: %s", err.Error())
		return contentV1.ErrorInternalServerError("delete page translation failed")
	}

	return nil
}
