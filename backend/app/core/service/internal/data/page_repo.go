package data

import (
	"context"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	entCrud "github.com/tx7do/go-crud/entgo"

	"github.com/tx7do/go-utils/copierutil"
	"github.com/tx7do/go-utils/mapper"
	"github.com/tx7do/go-utils/trans"

	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/page"
	"go-wind-cms/app/core/service/internal/data/ent/predicate"

	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
)

type PageRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper

	mapper *mapper.CopierMapper[contentV1.Page, ent.Page]

	repository *entCrud.Repository[
		ent.PageQuery, ent.PageSelect,
		ent.PageCreate, ent.PageCreateBulk,
		ent.PageUpdate, ent.PageUpdateOne,
		ent.PageDelete,
		predicate.Page,
		contentV1.Page, ent.Page,
	]

	statusConverter     *mapper.EnumTypeConverter[contentV1.Page_PageStatus, page.Status]
	typeConverter       *mapper.EnumTypeConverter[contentV1.Page_PageType, page.Type]
	editorTypeConverter *mapper.EnumTypeConverter[contentV1.EditorType, page.EditorType]

	pageTranslationRepo *PageTranslationRepo
}

func NewPageRepo(
	ctx *bootstrap.Context,
	entClient *entCrud.EntClient[*ent.Client],
	pageTranslationRepo *PageTranslationRepo,
) *PageRepo {
	repo := &PageRepo{
		entClient: entClient,
		log:       ctx.NewLoggerHelper("page/repo/core-service"),
		mapper:    mapper.NewCopierMapper[contentV1.Page, ent.Page](),
		statusConverter: mapper.NewEnumTypeConverter[contentV1.Page_PageStatus, page.Status](
			contentV1.Page_PageStatus_name, contentV1.Page_PageStatus_value,
		),
		typeConverter: mapper.NewEnumTypeConverter[contentV1.Page_PageType, page.Type](
			contentV1.Page_PageType_name, contentV1.Page_PageType_value,
		),
		editorTypeConverter: mapper.NewEnumTypeConverter[contentV1.EditorType, page.EditorType](
			contentV1.EditorType_name, contentV1.EditorType_value,
		),
		pageTranslationRepo: pageTranslationRepo,
	}

	repo.init()

	return repo
}

func (r *PageRepo) init() {
	r.repository = entCrud.NewRepository[
		ent.PageQuery, ent.PageSelect,
		ent.PageCreate, ent.PageCreateBulk,
		ent.PageUpdate, ent.PageUpdateOne,
		ent.PageDelete,
		predicate.Page,
		contentV1.Page, ent.Page,
	](r.mapper)

	r.mapper.AppendConverters(copierutil.NewTimeStringConverterPair())
	r.mapper.AppendConverters(copierutil.NewTimeTimestamppbConverterPair())

	r.mapper.AppendConverters(r.statusConverter.NewConverterPair())
	r.mapper.AppendConverters(r.typeConverter.NewConverterPair())
	r.mapper.AppendConverters(r.editorTypeConverter.NewConverterPair())
}

func (r *PageRepo) IsExist(ctx context.Context, id uint32) (bool, error) {
	exist, err := r.entClient.Client().Page.Query().
		Where(page.IDEQ(id)).
		Exist(ctx)
	if err != nil {
		r.log.Errorf("query page exist failed: %s", err.Error())
		return false, contentV1.ErrorInternalServerError("query page exist failed")
	}
	return exist, nil
}

func (r *PageRepo) buildPageTree(items []*contentV1.Page, parentId uint32) []*contentV1.Page {
	var tree []*contentV1.Page
	for _, item := range items {
		if item.GetParentId() == parentId {
			// 递归查找子节点
			children := r.buildPageTree(items, item.GetId())
			item.Children = children
			tree = append(tree, item)
		}
	}
	return tree
}

func (r *PageRepo) List(ctx context.Context, req *paginationV1.PagingRequest) (*contentV1.ListPageResponse, error) {
	if req == nil {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}

	builder := r.entClient.Client().Page.Query()

	ret, err := r.repository.ListWithPaging(ctx, builder, builder.Clone(), req)
	if err != nil {
		return nil, err
	}
	if ret == nil {
		return &contentV1.ListPageResponse{Total: 0, Items: nil}, nil
	}

	//for _, item := range ret.Items {
	//	translations, err := r.pageTranslationRepo.ListTranslations(ctx, item.GetId())
	//	if err != nil {
	//		r.log.Errorf("query translations failed: %s", err.Error())
	//		return nil, contentV1.ErrorInternalServerError("query translations failed")
	//	}
	//	item.Translations = translations
	//}

	for _, item := range ret.Items {
		languages, err := r.pageTranslationRepo.ListAvailedLanguages(ctx, item.GetId())
		if err != nil {
			r.log.Errorf("query availed languages failed: %s", err.Error())
			return nil, contentV1.ErrorInternalServerError("query availed languages failed")
		}
		item.AvailableLanguages = languages
	}

	return &contentV1.ListPageResponse{
		Total: ret.Total,
		Items: ret.Items,
	}, nil
}

func (r *PageRepo) Get(ctx context.Context, req *contentV1.GetPageRequest) (*contentV1.Page, error) {
	if req == nil {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}

	builder := r.entClient.Client().Page.Query()

	switch req.QueryBy.(type) {
	case *contentV1.GetPageRequest_Id:
		builder.Where(page.IDEQ(req.GetId()))
	case *contentV1.GetPageRequest_Slug:
		builder.Where(page.SlugEQ(req.GetSlug()))
	default:
		return nil, contentV1.ErrorBadRequest("invalid query_by value")
	}

	entity, err := builder.Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, contentV1.ErrorFileNotFound("page not found")
		}

		r.log.Errorf("query page failed: %s", err.Error())

		return nil, contentV1.ErrorInternalServerError("query page failed")
	}

	dto := r.mapper.ToDTO(entity)

	translations, err := r.pageTranslationRepo.ListTranslations(ctx, dto.GetId())
	if err != nil {
		r.log.Errorf("query translations failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("query translations failed")
	}
	dto.Translations = translations

	languages, err := r.pageTranslationRepo.ListAvailedLanguages(ctx, dto.GetId())
	if err != nil {
		r.log.Errorf("query availed languages failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("query availed languages failed")
	}
	dto.AvailableLanguages = languages

	return dto, nil
}

func (r *PageRepo) Create(ctx context.Context, req *contentV1.CreatePageRequest) (dto *contentV1.Page, err error) {
	if req == nil || req.Data == nil {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}

	var tx *ent.Tx
	tx, err = r.entClient.Client().Tx(ctx)
	if err != nil {
		r.log.Errorf("start transaction failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("start transaction failed")
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				r.log.Errorf("transaction rollback failed: %s", rollbackErr.Error())
			}
			return
		}
		if commitErr := tx.Commit(); commitErr != nil {
			r.log.Errorf("transaction commit failed: %s", commitErr.Error())
			err = contentV1.ErrorInternalServerError("transaction commit failed")
		}
	}()

	builder := tx.Page.Create().
		SetNillableStatus(r.statusConverter.ToEntity(req.Data.Status)).
		SetNillableType(r.typeConverter.ToEntity(req.Data.Type)).
		SetNillableEditorType(r.editorTypeConverter.ToEntity(req.Data.EditorType)).
		SetNillableSlug(req.Data.Slug).
		SetNillableAuthorID(req.Data.AuthorId).
		SetNillableAuthorName(req.Data.AuthorName).
		SetNillableDisallowComment(req.Data.DisallowComment).
		SetNillableRedirectURL(req.Data.RedirectUrl).
		SetNillableShowInNavigation(req.Data.ShowInNavigation).
		SetNillableSortOrder(req.Data.SortOrder).
		SetNillableTemplate(req.Data.Template).
		SetNillableIsCustomTemplate(req.Data.IsCustomTemplate).
		SetNillableVisits(req.Data.Visits).
		SetNillableParentID(req.Data.ParentId).
		SetNillableDepth(req.Data.Depth).
		SetNillablePath(req.Data.Path).
		SetNillableCreatedBy(req.Data.CreatedBy).
		SetCreatedAt(time.Now())

	if req.Data.CustomFields != nil {
		builder.SetCustomFields(trans.Ptr(req.Data.GetCustomFields()))
	}

	var entity *ent.Page
	if entity, err = builder.Save(ctx); err != nil {
		r.log.Errorf("insert page failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("insert page failed")
	}

	if len(req.Data.Translations) > 0 {
		if err = r.pageTranslationRepo.CleanTranslations(ctx, tx, entity.ID); err != nil {
			r.log.Errorf("clean translations failed: %s", err.Error())
			return nil, contentV1.ErrorInternalServerError("clean translations failed")
		}

		for i := range req.Data.Translations {
			req.Data.Translations[i].PageId = trans.Ptr(entity.ID)

			_ = r.pageTranslationRepo.PrepareTranslation(ctx, req.Data.Translations[i])
		}

		if err = r.pageTranslationRepo.BatchCreate(ctx, tx, req.Data.GetTranslations()); err != nil {
			r.log.Errorf("batch insert translations failed: %s", err.Error())
			return nil, contentV1.ErrorInternalServerError("batch insert translations failed")
		}
	}

	return r.mapper.ToDTO(entity), nil
}

func (r *PageRepo) Update(ctx context.Context, req *contentV1.UpdatePageRequest) (dto *contentV1.Page, err error) {
	if req == nil || req.Data == nil {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}

	// 如果不存在则创建
	if req.GetAllowMissing() {
		exist, err := r.IsExist(ctx, req.GetId())
		if err != nil {
			return nil, err
		}
		if !exist {
			req.Data.CreatedBy = req.Data.UpdatedBy
			req.Data.UpdatedBy = nil
			_, err = r.Create(ctx, &contentV1.CreatePageRequest{Data: req.Data})
			return nil, err
		}
	}

	var tx *ent.Tx
	tx, err = r.entClient.Client().Tx(ctx)
	if err != nil {
		r.log.Errorf("start transaction failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("start transaction failed")
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				r.log.Errorf("transaction rollback failed: %s", rollbackErr.Error())
			}
			return
		}
		if commitErr := tx.Commit(); commitErr != nil {
			r.log.Errorf("transaction commit failed: %s", commitErr.Error())
			err = contentV1.ErrorInternalServerError("transaction commit failed")
		}
	}()

	if len(req.Data.Translations) > 0 {
		//if err = r.pageTranslationRepo.CleanTranslations(ctx, tx, req.GetId()); err != nil {
		//	r.log.Errorf("clean translations failed: %s", err.Error())
		//	return nil, contentV1.ErrorInternalServerError("clean translations failed")
		//}

		for i := range req.Data.Translations {
			req.Data.Translations[i].PageId = trans.Ptr(req.GetId())
		}

		if err = r.pageTranslationRepo.BatchCreate(ctx, tx, req.Data.GetTranslations()); err != nil {
			r.log.Errorf("batch insert translations failed: %s", err.Error())
			return nil, contentV1.ErrorInternalServerError("batch insert translations failed")
		}
	}

	builder := tx.Page.UpdateOneID(req.GetId())
	result, err := r.repository.UpdateOne(ctx, builder, req.Data, req.GetUpdateMask(),
		func(dto *contentV1.Page) {
			builder.
				SetNillableStatus(r.statusConverter.ToEntity(req.Data.Status)).
				SetNillableType(r.typeConverter.ToEntity(req.Data.Type)).
				SetNillableEditorType(r.editorTypeConverter.ToEntity(req.Data.EditorType)).
				SetNillableSlug(req.Data.Slug).
				SetNillableAuthorID(req.Data.AuthorId).
				SetNillableAuthorName(req.Data.AuthorName).
				SetNillableDisallowComment(req.Data.DisallowComment).
				SetNillableRedirectURL(req.Data.RedirectUrl).
				SetNillableShowInNavigation(req.Data.ShowInNavigation).
				SetNillableSortOrder(req.Data.SortOrder).
				SetNillableTemplate(req.Data.Template).
				SetNillableIsCustomTemplate(req.Data.IsCustomTemplate).
				SetNillableVisits(req.Data.Visits).
				SetNillableParentID(req.Data.ParentId).
				SetNillableDepth(req.Data.Depth).
				SetNillablePath(req.Data.Path).
				SetNillableUpdatedBy(req.Data.UpdatedBy).
				SetUpdatedAt(time.Now())

			if req.Data.CustomFields != nil {
				builder.SetCustomFields(trans.Ptr(req.Data.GetCustomFields()))
			}
		},
		func(s *sql.Selector) {
			s.Where(sql.EQ(page.FieldID, req.GetId()))
		},
	)

	return result, err
}

func (r *PageRepo) Delete(ctx context.Context, req *contentV1.DeletePageRequest) (err error) {
	if req == nil {
		return contentV1.ErrorBadRequest("invalid parameter")
	}

	var tx *ent.Tx
	tx, err = r.entClient.Client().Tx(ctx)
	if err != nil {
		r.log.Errorf("start transaction failed: %s", err.Error())
		return contentV1.ErrorInternalServerError("start transaction failed")
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				r.log.Errorf("transaction rollback failed: %s", rollbackErr.Error())
			}
			return
		}
		if commitErr := tx.Commit(); commitErr != nil {
			r.log.Errorf("transaction commit failed: %s", commitErr.Error())
			err = contentV1.ErrorInternalServerError("transaction commit failed")
		}
	}()

	if err = tx.Page.
		DeleteOneID(req.GetId()).
		Exec(ctx); err != nil {
		r.log.Errorf("delete one data failed: %s", err.Error())
	}

	if err = r.pageTranslationRepo.CleanTranslations(ctx, tx, req.GetId()); err != nil {
		r.log.Errorf("clean translations failed: %s", err.Error())
		return contentV1.ErrorInternalServerError("clean translations failed")
	}

	return err
}

func (r *PageRepo) TranslationExists(ctx context.Context, pageId uint32, languageCode string) (bool, error) {
	return r.pageTranslationRepo.TranslationExists(ctx, pageId, languageCode)
}

func (r *PageRepo) CreateTranslation(ctx context.Context, req *contentV1.CreatePageTranslationRequest) (*contentV1.PageTranslation, error) {
	if req == nil || req.Data == nil {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}

	if len(req.Data.GetLanguageCode()) == 0 {
		return nil, contentV1.ErrorBadRequest("language code is required")
	}

	if req.GetPageId() == 0 {
		return nil, contentV1.ErrorBadRequest("post id is required")
	}

	req.Data.PageId = trans.Ptr(req.GetPageId())

	return r.pageTranslationRepo.CreateTranslation(ctx, req.Data)
}

func (r *PageRepo) UpdateTranslation(ctx context.Context, req *contentV1.UpdatePageTranslationRequest) (*contentV1.PageTranslation, error) {
	if req == nil || req.Data == nil {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}

	if len(req.Data.GetLanguageCode()) == 0 {
		return nil, contentV1.ErrorBadRequest("language code is required")
	}

	if req.Data.GetPageId() == 0 {
		return nil, contentV1.ErrorBadRequest("post id is required")
	}

	if exist, err := r.TranslationExists(ctx, req.Data.GetPageId(), req.Data.GetLanguageCode()); err != nil {
		return nil, err
	} else if !exist {
		if req.GetAllowMissing() {
			return r.CreateTranslation(ctx, &contentV1.CreatePageTranslationRequest{
				Data:   req.Data,
				PageId: req.Data.GetPageId(),
			})
		}

		return nil, contentV1.ErrorFileNotFound("translation not found")
	}

	return r.pageTranslationRepo.UpdateTranslation(ctx, req.GetId(), req.Data, req.GetUpdateMask())
}

func (r *PageRepo) GetTranslation(ctx context.Context, req *contentV1.GetPageRequest) (*contentV1.PageTranslation, error) {
	if req == nil {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}

	return r.pageTranslationRepo.GetTranslation(ctx, req.GetId(), req.GetLocale())
}

func (r *PageRepo) ListTranslations(ctx context.Context, pageID uint32) ([]*contentV1.PageTranslation, error) {
	return r.pageTranslationRepo.ListTranslations(ctx, pageID)
}

func (r *PageRepo) DeleteTranslation(ctx context.Context, req *contentV1.DeletePageTranslationRequest) error {
	return r.pageTranslationRepo.DeleteTranslation(ctx, req)
}

func (r *PageRepo) CleanTranslations(ctx context.Context, tx *ent.Tx, pageID uint32) error {
	return r.pageTranslationRepo.CleanTranslations(ctx, tx, pageID)
}
