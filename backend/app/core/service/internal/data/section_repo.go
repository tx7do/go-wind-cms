package data

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	entCrud "github.com/tx7do/go-crud/entgo"

	"github.com/tx7do/go-utils/copierutil"
	"github.com/tx7do/go-utils/mapper"
	"github.com/tx7do/go-utils/trans"

	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/section"

	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
)

// SectionRepo 管理 Page 的嵌套子部件 Section。
// 无独立服务入口：所有读写经 PageRepo 以页面为单位整体替换，
// 对齐 content_model.fields 惯例。
type SectionRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper

	mapper *mapper.CopierMapper[contentV1.Section, ent.Section]

	typeConverter *mapper.EnumTypeConverter[contentV1.SectionType, section.Type]

	sectionTranslationRepo *SectionTranslationRepo
}

func NewSectionRepo(
	ctx *bootstrap.Context,
	entClient *entCrud.EntClient[*ent.Client],
	sectionTranslationRepo *SectionTranslationRepo,
) *SectionRepo {
	repo := &SectionRepo{
		entClient: entClient,
		log:       ctx.NewLoggerHelper("section/repo/core-service"),
		mapper:    mapper.NewCopierMapper[contentV1.Section, ent.Section](),
		typeConverter: mapper.NewEnumTypeConverter[contentV1.SectionType, section.Type](
			contentV1.SectionType_name, contentV1.SectionType_value,
		),
		sectionTranslationRepo: sectionTranslationRepo,
	}

	repo.init()

	return repo
}

func (r *SectionRepo) init() {
	r.mapper.AppendConverters(copierutil.NewTimeStringConverterPair())
	r.mapper.AppendConverters(copierutil.NewTimeTimestamppbConverterPair())

	r.mapper.AppendConverters(r.typeConverter.NewConverterPair())
}

// BatchCreateByPage 在给定页面下批量创建 section 及其翻译。
// page_id 强制设为父页面值；tenant_id 由 viewer context 推导，
// 与 Section.Query/Delete 路径的 tenant 处理一致（对齐 content_model.batchCreateFields，
// 后者依赖 viewer 而非父实体传值）。调用方须在自身事务内调用，本方法不再开事务。
func (r *SectionRepo) BatchCreateByPage(
	ctx context.Context,
	tx *ent.Tx,
	pageID uint32,
	sections []*contentV1.Section,
) error {
	tid, hasTenant := maybeTenantFromViewer(ctx)
	var tenantPtr *uint32
	if hasTenant {
		tenantPtr = trans.Ptr(tid)
	}

	for _, s := range sections {
		if s == nil {
			continue
		}

		s.PageId = trans.Ptr(pageID)

		builder := tx.Section.Create().
			SetNillablePageID(s.PageId).
			SetNillableTenantID(tenantPtr).
			SetNillableType(r.typeConverter.ToEntity(s.Type)).
			SetNillableName(s.Name).
			SetNillableSortOrder(s.SortOrder).
			SetNillableCreatedBy(s.CreatedBy).
			SetCreatedAt(time.Now())

		if s.Config != nil {
			builder.SetConfig(trans.Ptr(s.GetConfig()))
		}

		created, err := builder.Save(ctx)
		if err != nil {
			r.log.Errorf("insert section failed: %s", err.Error())
			return contentV1.ErrorInternalServerError("insert section failed")
		}

		if len(s.Translations) > 0 {
			if err = r.sectionTranslationRepo.CleanTranslations(ctx, tx, created.ID); err != nil {
				r.log.Errorf("clean section translations failed: %s", err.Error())
				return contentV1.ErrorInternalServerError("clean section translations failed")
			}

			for i := range s.Translations {
				s.Translations[i].SectionId = trans.Ptr(created.ID)
			}

			if err = r.sectionTranslationRepo.BatchCreate(ctx, tx, s.GetTranslations()); err != nil {
				r.log.Errorf("batch insert section translations failed: %s", err.Error())
				return contentV1.ErrorInternalServerError("batch insert section translations failed")
			}
		}
	}
	return nil
}

// ReplaceByPage 整体替换给定页面下的 section：先清旧（含翻译），再建新。
// 对齐 content_model.replaceFields。调用方须在自身事务内调用。
func (r *SectionRepo) ReplaceByPage(
	ctx context.Context,
	tx *ent.Tx,
	pageID uint32,
	sections []*contentV1.Section,
) error {
	if err := r.CleanByPageID(ctx, tx, pageID); err != nil {
		return err
	}
	return r.BatchCreateByPage(ctx, tx, pageID, sections)
}

// ListByPage 列出给定页面下的所有 section（含各自翻译/可用语言），
// 供 PageRepo.Get/List 水合 Page.Sections。对齐 content_model.fillModelAssociations。
func (r *SectionRepo) ListByPage(ctx context.Context, pageID uint32) ([]*contentV1.Section, error) {
	entities, err := r.entClient.Client().Section.Query().
		Where(section.PageIDEQ(pageID)).
		Order(section.BySortOrder(), section.ByID()).
		All(ctx)
	if err != nil {
		r.log.Errorf("query sections by page id failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("query sections by page id failed")
	}

	items := make([]*contentV1.Section, 0, len(entities))
	for _, entity := range entities {
		dto := r.mapper.ToDTO(entity)

		// 整体替换语义下，父页面持有全部语言翻译；详情/列表均回退全量译文，
		// 前端按"匹配当前语言→回退第一条"兜底
		translations, tErr := r.sectionTranslationRepo.ListTranslationsWithFallback(ctx, dto.GetId(), "", nil)
		if tErr != nil {
			r.log.Errorf("query section translations failed: %s", tErr.Error())
			return nil, contentV1.ErrorInternalServerError("query section translations failed")
		}
		dto.Translations = translations

		languages, lErr := r.sectionTranslationRepo.ListAvailedLanguages(ctx, dto.GetId())
		if lErr != nil {
			r.log.Errorf("query section availed languages failed: %s", lErr.Error())
			return nil, contentV1.ErrorInternalServerError("query section availed languages failed")
		}
		dto.AvailableLanguages = languages

		items = append(items, dto)
	}
	return items, nil
}

func (r *SectionRepo) CleanTranslations(ctx context.Context, tx *ent.Tx, sectionID uint32) error {
	return r.sectionTranslationRepo.CleanTranslations(ctx, tx, sectionID)
}

// CleanByPageID 删除指定页面下的所有 section 及其翻译。
// 在页面删除时调用，防止 section 成为孤儿记录。
func (r *SectionRepo) CleanByPageID(ctx context.Context, tx *ent.Tx, pageID uint32) error {
	if pageID == 0 {
		return nil
	}

	// 先查出所有 section ID，用于清理翻译
	queryBuilder := tx.Section.Query().Where(section.PageIDEQ(pageID))
	if tid, hasTenant := maybeTenantFromViewer(ctx); hasTenant {
		queryBuilder.Where(section.TenantIDEQ(tid))
	}
	sectionIDs, err := queryBuilder.IDs(ctx)
	if err != nil {
		r.log.Errorf("query sections by page id failed: %s", err.Error())
		return contentV1.ErrorInternalServerError("query sections by page id failed")
	}

	// 清理每个 section 的翻译
	for _, sid := range sectionIDs {
		if err := r.sectionTranslationRepo.CleanTranslations(ctx, tx, sid); err != nil {
			r.log.Errorf("clean section translations failed for section %d: %s", sid, err.Error())
			return contentV1.ErrorInternalServerError("clean section translations failed")
		}
	}

	// 删除所有 section
	delBuilder := tx.Section.Delete().Where(section.PageIDEQ(pageID))
	if tid, hasTenant := maybeTenantFromViewer(ctx); hasTenant {
		delBuilder.Where(section.TenantIDEQ(tid))
	}
	if _, err := delBuilder.Exec(ctx); err != nil {
		r.log.Errorf("delete sections by page id failed: %s", err.Error())
		return contentV1.ErrorInternalServerError("delete sections by page id failed")
	}

	return nil
}
