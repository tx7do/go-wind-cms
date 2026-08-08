package data

import (
	"context"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	entCrud "github.com/tx7do/go-crud/entgo"
	"github.com/tx7do/go-crud/pagination"
	paginationFilter "github.com/tx7do/go-crud/pagination/filter"

	"github.com/tx7do/go-utils/copierutil"
	"github.com/tx7do/go-utils/mapper"
	"github.com/tx7do/go-utils/trans"

	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/category"
	"go-wind-cms/app/core/service/internal/data/ent/predicate"

	"go-wind-cms/pkg/utils"

	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
)

type CategoryRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper

	mapper *mapper.CopierMapper[contentV1.Category, ent.Category]

	repository *entCrud.Repository[
		ent.CategoryQuery, ent.CategorySelect,
		ent.CategoryCreate, ent.CategoryCreateBulk,
		ent.CategoryUpdate, ent.CategoryUpdateOne,
		ent.CategoryDelete,
		predicate.Category,
		contentV1.Category, ent.Category,
	]

	statusConverter *mapper.EnumTypeConverter[contentV1.Category_CategoryStatus, category.Status]

	categoryTranslationRepo *CategoryTranslationRepo
}

func NewCategoryRepo(
	ctx *bootstrap.Context,
	entClient *entCrud.EntClient[*ent.Client],
	categoryTranslationRepo *CategoryTranslationRepo,
) *CategoryRepo {
	repo := &CategoryRepo{
		entClient: entClient,
		log:       ctx.NewLoggerHelper("category/repo/core-service"),
		mapper:    mapper.NewCopierMapper[contentV1.Category, ent.Category](),
		statusConverter: mapper.NewEnumTypeConverter[contentV1.Category_CategoryStatus, category.Status](
			contentV1.Category_CategoryStatus_name, contentV1.Category_CategoryStatus_value,
		),
		categoryTranslationRepo: categoryTranslationRepo,
	}

	repo.init()

	return repo
}

func (r *CategoryRepo) init() {
	r.repository = entCrud.NewRepository[
		ent.CategoryQuery, ent.CategorySelect,
		ent.CategoryCreate, ent.CategoryCreateBulk,
		ent.CategoryUpdate, ent.CategoryUpdateOne,
		ent.CategoryDelete,
		predicate.Category,
		contentV1.Category, ent.Category,
	](r.mapper)

	r.mapper.AppendConverters(copierutil.NewTimeStringConverterPair())
	r.mapper.AppendConverters(copierutil.NewTimeTimestamppbConverterPair())

	r.mapper.AppendConverters(r.statusConverter.NewConverterPair())
}

func (r *CategoryRepo) count(ctx context.Context, whereCond []func(s *sql.Selector)) (int, error) {
	builder := r.entClient.Client().Category.Query()
	if len(whereCond) != 0 {
		builder.Modify(whereCond...)
	}

	count, err := builder.Count(ctx)
	if err != nil {
		r.log.Errorf("query count failed: %s", err.Error())
	}

	return count, err
}

func (r *CategoryRepo) IsExist(ctx context.Context, id uint32) (bool, error) {
	exist, err := r.entClient.Client().Category.Query().
		Where(category.IDEQ(id)).
		Exist(ctx)
	if err != nil {
		r.log.Errorf("query category exist failed: %s", err.Error())
		return false, contentV1.ErrorInternalServerError("query category exist failed")
	}
	return exist, nil
}

func FilterViewMask(excludeFields []string, mask *fieldmaskpb.FieldMask) *fieldmaskpb.FieldMask {
	if mask == nil {
		return nil
	}
	if len(excludeFields) == 0 {
		return mask
	}

	excludeMap := make(map[string]struct{})
	for _, field := range excludeFields {
		excludeMap[field] = struct{}{}
	}

	var paths []string
	for _, path := range mask.Paths {
		path = strings.TrimSpace(path)
		if _, exist := excludeMap[path]; !exist {
			paths = append(paths, path)
		}
	}

	if len(paths) == 0 {
		return nil
	}

	return &fieldmaskpb.FieldMask{Paths: paths}
}

func (r *CategoryRepo) prepareTranslationMaskFields(req *paginationV1.PagingRequest) (
	excludeConditions []*paginationV1.FilterCondition,
	translationMaskFields []string,
	needQueryTranslation bool,
	treeTravel bool,
	err error,
) {
	var filterExpr *paginationV1.FilterExpr
	filterExpr, err = paginationFilter.ConvertFilterByPagingRequest(req)
	if err != nil {
		r.log.Errorf("convert filter by paging request failed: %s", err.Error())
		return
	}

	excludeFields := []string{"translations", "available_languages"}
	if req.FieldMask != nil && len(req.FieldMask.Paths) > 0 {
		for _, path := range req.FieldMask.Paths {
			path = strings.TrimSpace(path)

			if path == "translations" || strings.HasPrefix(path, "translations.") {
				needQueryTranslation = true
			}
			if strings.HasPrefix(path, "translations.") {
				excludeFields = append(excludeFields, path)

				path = strings.TrimPrefix(path, "translations.")
				if len(path) > 0 {
					translationMaskFields = append(translationMaskFields, path)
				}
			}
			if path == "children" {
				treeTravel = true
				excludeFields = append(excludeFields, path)
			}
		}

		req.FieldMask = FilterViewMask(excludeFields, req.FieldMask)

	} else {
		needQueryTranslation = true
	}

	excludeConditions = pagination.FilterFields(filterExpr, []string{
		"locale",
	})
	req.FilteringType = &paginationV1.PagingRequest_FilterExpr{FilterExpr: filterExpr}

	return
}

func (r *CategoryRepo) buildCategoryTree(items []*contentV1.Category, parentId uint32) []*contentV1.Category {
	var tree []*contentV1.Category
	for _, item := range items {
		if item.GetParentId() == parentId {
			// 递归查找子节点
			children := r.buildCategoryTree(items, item.GetId())
			item.Children = children
			tree = append(tree, item)
		}
	}
	return tree
}

func (r *CategoryRepo) List(ctx context.Context, req *paginationV1.PagingRequest) (*contentV1.ListCategoryResponse, error) {
	if req == nil {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}

	builder := r.entClient.Client().Category.Query()

	excludeConditions, translationMaskFields, needQueryTranslation, treeTravel, err := r.prepareTranslationMaskFields(req)
	if err != nil {
		return nil, err
	}

	ret, err := r.repository.ListWithPaging(ctx, builder, builder.Clone(), req)
	if err != nil {
		return nil, err
	}
	if ret == nil {
		return &contentV1.ListCategoryResponse{Total: 0, Items: nil}, nil
	}

	if needQueryTranslation {
		viewMask := &fieldmaskpb.FieldMask{
			Paths: translationMaskFields,
		}

		var locale string
		if len(excludeConditions) > 0 {
			for _, cond := range excludeConditions {
				if cond.GetField() == "locale" {
					locale = cond.GetValue()
					break
				}
			}
		}

		for _, item := range ret.Items {
			translations, err := r.categoryTranslationRepo.ListTranslations(ctx, item.GetId(), locale, viewMask)
			if err != nil {
				r.log.Errorf("query translations failed: %s", err.Error())
				return nil, contentV1.ErrorInternalServerError("query translations failed")
			}
			item.Translations = translations
		}
	}

	for _, item := range ret.Items {
		languages, err := r.categoryTranslationRepo.ListAvailedLanguages(ctx, item.GetId())
		if err != nil {
			r.log.Errorf("query availed languages failed: %s", err.Error())
			return nil, contentV1.ErrorInternalServerError("query availed languages failed")
		}
		item.AvailableLanguages = languages
	}

	if treeTravel {
		roots := r.buildCategoryTree(ret.Items, 0) // 假设根节点ParentId为0
		return &contentV1.ListCategoryResponse{
			Total: ret.Total,
			Items: roots,
		}, nil
	}

	return &contentV1.ListCategoryResponse{
		Total: ret.Total,
		Items: ret.Items,
	}, nil
}

// 递归查询子节点
func (r *CategoryRepo) getCategoryWithChildren(ctx context.Context, id uint32, locale string, viewMask *fieldmaskpb.FieldMask, translationMaskFields []string) (*contentV1.Category, error) {
	entity, err := r.entClient.Client().Category.Query().Where(category.IDEQ(id)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, contentV1.ErrorFileNotFound("category not found")
		}
		r.log.Errorf("query category failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("query category failed")
	}
	dto := r.mapper.ToDTO(entity)

	// 查询子节点
	childrenEntities, err := r.entClient.Client().Category.Query().Where(category.ParentIDEQ(id)).All(ctx)
	if err != nil {
		r.log.Errorf("query children failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("query children failed")
	}
	for _, child := range childrenEntities {
		childDTO, err := r.getCategoryWithChildren(ctx, child.ID, locale, viewMask, translationMaskFields)
		if err != nil {
			return nil, err
		}
		dto.Children = append(dto.Children, childDTO)
	}

	// 查询翻译和可用语言（可复用原有逻辑）
	translations, err := r.categoryTranslationRepo.ListTranslations(
		ctx,
		dto.GetId(),
		locale,
		&fieldmaskpb.FieldMask{Paths: translationMaskFields},
	)
	if err == nil {
		dto.Translations = translations
	}
	languages, err := r.categoryTranslationRepo.ListAvailedLanguages(ctx, dto.GetId())
	if err == nil {
		dto.AvailableLanguages = languages
	}

	return dto, nil
}

func (r *CategoryRepo) Get(ctx context.Context, req *contentV1.GetCategoryRequest) (*contentV1.Category, error) {
	if req == nil {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}

	treeTravel := false
	var translationMaskFields []string
	excludeFields := []string{"translations", "available_languages"}
	if req.ViewMask != nil && len(req.ViewMask.Paths) > 0 {
		for _, path := range req.ViewMask.Paths {
			path = strings.TrimSpace(path)

			if strings.HasPrefix(path, "translations.") {
				excludeFields = append(excludeFields, path)

				path = strings.TrimPrefix(path, "translations.")
				if len(path) > 0 {
					translationMaskFields = append(translationMaskFields, path)
				}
			}
			if path == "children" {
				treeTravel = true
				excludeFields = append(excludeFields, path)
			}
		}
		req.ViewMask = FilterViewMask(excludeFields, req.ViewMask)
	}

	if !treeTravel {
		var dto *contentV1.Category

		switch req.QueryBy.(type) {
		case *contentV1.GetCategoryRequest_Id:
			entity, err := r.entClient.Client().Category.Query().Where(category.IDEQ(req.GetId())).Only(ctx)
			if err != nil {
				if ent.IsNotFound(err) {
					return nil, contentV1.ErrorFileNotFound("category not found")
				}
				r.log.Errorf("query category failed: %s", err.Error())
				return nil, contentV1.ErrorInternalServerError("query category failed")
			}
			dto = r.mapper.ToDTO(entity)

		case *contentV1.GetCategoryRequest_Code:
			entity, err := r.entClient.Client().Category.Query().Where(category.CodeEQ(req.GetCode())).Only(ctx)
			if err != nil {
				if ent.IsNotFound(err) {
					return nil, contentV1.ErrorFileNotFound("category not found")
				}
				r.log.Errorf("query category failed: %s", err.Error())
				return nil, contentV1.ErrorInternalServerError("query category failed")
			}
			dto = r.mapper.ToDTO(entity)

		default:
			return nil, contentV1.ErrorBadRequest("invalid query_by value")
		}

		translations, err := r.categoryTranslationRepo.ListTranslations(
			ctx,
			dto.GetId(),
			req.GetLocale(),
			&fieldmaskpb.FieldMask{Paths: translationMaskFields},
		)
		if err != nil {
			r.log.Errorf("query translations failed: %s", err.Error())
			return nil, contentV1.ErrorInternalServerError("query translations failed")
		}
		dto.Translations = translations

		languages, err := r.categoryTranslationRepo.ListAvailedLanguages(ctx, dto.GetId())
		if err != nil {
			r.log.Errorf("query availed languages failed: %s", err.Error())
			return nil, contentV1.ErrorInternalServerError("query availed languages failed")
		}
		dto.AvailableLanguages = languages

		return dto, nil
	}

	var id uint32
	switch req.QueryBy.(type) {
	case *contentV1.GetCategoryRequest_Id:
		id = req.GetId()
	case *contentV1.GetCategoryRequest_Code:
		entity, err := r.entClient.Client().Category.Query().Where(category.CodeEQ(req.GetCode())).Only(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				return nil, contentV1.ErrorFileNotFound("category not found")
			}
			r.log.Errorf("query category failed: %s", err.Error())
			return nil, contentV1.ErrorInternalServerError("query category failed")
		}
		id = entity.ID
	default:
		return nil, contentV1.ErrorBadRequest("invalid query_by value")
	}

	return r.getCategoryWithChildren(ctx, id, req.GetLocale(), req.GetViewMask(), translationMaskFields)
}

func (r *CategoryRepo) Create(ctx context.Context, req *contentV1.CreateCategoryRequest) (dto *contentV1.Category, err error) {
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

	builder := tx.Category.Create().
		SetNillableStatus(r.statusConverter.ToEntity(req.Data.Status)).
		SetNillableSortOrder(req.Data.SortOrder).
		SetNillableIsNav(req.Data.IsNav).
		SetNillableIcon(req.Data.Icon).
		SetNillableCode(req.Data.Code).
		SetPostCount(0).
		SetDirectPostCount(0).
		SetNillableParentID(req.Data.ParentId).
		SetNillableDepth(req.Data.Depth).
		SetNillablePath(req.Data.Path).
		SetNillableCreatedBy(req.Data.CreatedBy).
		SetCreatedAt(time.Now())

	if req.Data.CustomFields != nil {
		builder.SetCustomFields(trans.Ptr(req.Data.GetCustomFields()))
	}

	var entity *ent.Category
	if entity, err = builder.Save(ctx); err != nil {
		r.log.Errorf("insert category failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("insert category failed")
	}

	if len(req.Data.Translations) > 0 {
		if err = r.categoryTranslationRepo.CleanTranslations(ctx, tx, entity.ID); err != nil {
			r.log.Errorf("clean translations failed: %s", err.Error())
			return nil, contentV1.ErrorInternalServerError("clean translations failed")
		}

		for i := range req.Data.Translations {
			req.Data.Translations[i].CategoryId = trans.Ptr(entity.ID)
		}

		if err = r.categoryTranslationRepo.BatchCreate(ctx, tx, req.Data.GetTranslations()); err != nil {
			r.log.Errorf("batch insert translations failed: %s", err.Error())
			return nil, contentV1.ErrorInternalServerError("batch insert translations failed")
		}
	}

	return r.mapper.ToDTO(entity), nil
}

func (r *CategoryRepo) Update(ctx context.Context, req *contentV1.UpdateCategoryRequest) (dto *contentV1.Category, err error) {
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
			_, err = r.Create(ctx, &contentV1.CreateCategoryRequest{Data: req.Data})
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
		if err = r.categoryTranslationRepo.CleanTranslations(ctx, tx, req.GetId()); err != nil {
			r.log.Errorf("clean translations failed: %s", err.Error())
			return nil, contentV1.ErrorInternalServerError("clean translations failed")
		}

		for i := range req.Data.Translations {
			req.Data.Translations[i].CategoryId = trans.Ptr(req.GetId())
		}

		if err = r.categoryTranslationRepo.BatchCreate(ctx, tx, req.Data.GetTranslations()); err != nil {
			r.log.Errorf("batch insert translations failed: %s", err.Error())
			return nil, contentV1.ErrorInternalServerError("batch insert translations failed")
		}
	}

	tid, hasTenant := maybeTenantFromViewer(ctx)
	callerUserID, hasUser := viewerUserIDFromContext(ctx)
	// 剔除关联字段：translations / available_languages 是子表数据，非本表列，
	// 放入 updateMask 会触发列不存在或误更新，关联由专用翻译接口单独维护。
	if req.UpdateMask != nil {
		req.UpdateMask.Paths = utils.FilterBlacklist(req.UpdateMask.GetPaths(), []string{
			"translations", "available_languages",
		})
	}
	builder := tx.Category.UpdateOneID(req.GetId())
	builder.Where(category.IDEQ(req.GetId()))
	if hasTenant {
		builder.Where(category.TenantIDEQ(tid))
	}
	result, err := r.repository.UpdateOne(ctx, builder, req.Data, req.GetUpdateMask(),
		func(dto *contentV1.Category) {
			builder.
				SetNillableStatus(r.statusConverter.ToEntity(req.Data.Status)).
				SetNillableSortOrder(req.Data.SortOrder).
				SetNillableIsNav(req.Data.IsNav).
				SetNillableIcon(req.Data.Icon).
				SetNillableCode(req.Data.Code).
				SetNillableParentID(req.Data.ParentId).
				SetNillableDepth(req.Data.Depth).
				SetNillablePath(req.Data.Path).
				SetUpdatedAt(time.Now())

			// updated_by 强制由服务端 viewer context 推导，忽略客户端传入值
			if hasUser {
				builder.SetUpdatedBy(callerUserID)
			}

			if req.Data.CustomFields != nil {
				builder.SetCustomFields(trans.Ptr(req.Data.GetCustomFields()))
			}
		},
		func(s *sql.Selector) {
			s.Where(sql.EQ(category.FieldID, req.GetId()))
		},
	)

	return result, err
}

func (r *CategoryRepo) Delete(ctx context.Context, req *contentV1.DeleteCategoryRequest) (err error) {
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

	tid, hasTenant := maybeTenantFromViewer(ctx)
	delBuilder := tx.Category.Delete()
	delBuilder.Where(category.IDEQ(req.GetId()))
	if hasTenant {
		delBuilder.Where(category.TenantIDEQ(tid))
	}
	if _, err = delBuilder.Exec(ctx); err != nil {
		r.log.Errorf("delete one data failed: %s", err.Error())
		return contentV1.ErrorInternalServerError("delete one data failed")
	}

	if err = r.categoryTranslationRepo.CleanTranslations(ctx, tx, req.GetId()); err != nil {
		r.log.Errorf("clean translations failed: %s", err.Error())
		return contentV1.ErrorInternalServerError("clean translations failed")
	}

	return nil
}

func (r *CategoryRepo) TranslationExists(ctx context.Context, categoryId uint32, languageCode string) (bool, error) {
	return r.categoryTranslationRepo.TranslationExists(ctx, categoryId, languageCode)
}

func (r *CategoryRepo) CreateTranslation(ctx context.Context, req *contentV1.CreateCategoryTranslationRequest) (*contentV1.CategoryTranslation, error) {
	if req == nil || req.Data == nil {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}

	if len(req.Data.GetLanguageCode()) == 0 {
		return nil, contentV1.ErrorBadRequest("language code is required")
	}

	if req.GetCategoryId() == 0 {
		return nil, contentV1.ErrorBadRequest("post id is required")
	}

	req.Data.CategoryId = trans.Ptr(req.GetCategoryId())

	return r.categoryTranslationRepo.CreateTranslation(ctx, req.Data)
}

func (r *CategoryRepo) UpdateTranslation(ctx context.Context, req *contentV1.UpdateCategoryTranslationRequest) (*contentV1.CategoryTranslation, error) {
	if req == nil || req.Data == nil {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}

	if len(req.Data.GetLanguageCode()) == 0 {
		return nil, contentV1.ErrorBadRequest("language code is required")
	}

	if req.Data.GetCategoryId() == 0 {
		return nil, contentV1.ErrorBadRequest("post id is required")
	}

	if exist, err := r.TranslationExists(ctx, req.Data.GetCategoryId(), req.Data.GetLanguageCode()); err != nil {
		return nil, err
	} else if !exist {
		if req.GetAllowMissing() {
			return r.CreateTranslation(ctx, &contentV1.CreateCategoryTranslationRequest{
				Data:       req.Data,
				CategoryId: req.Data.GetCategoryId(),
			})
		}

		return nil, contentV1.ErrorFileNotFound("translation not found")
	}

	return r.categoryTranslationRepo.UpdateTranslation(ctx, req.GetId(), req.Data, req.GetUpdateMask())
}

func (r *CategoryRepo) GetTranslation(ctx context.Context, req *contentV1.GetCategoryRequest) (*contentV1.CategoryTranslation, error) {
	if req == nil {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}

	return r.categoryTranslationRepo.GetTranslation(ctx, req.GetId(), req.GetLocale())
}

func (r *CategoryRepo) ListTranslations(ctx context.Context, categoryID uint32) ([]*contentV1.CategoryTranslation, error) {
	return r.categoryTranslationRepo.ListTranslations(ctx, categoryID, "", nil)
}

func (r *CategoryRepo) DeleteTranslation(ctx context.Context, req *contentV1.DeleteCategoryTranslationRequest) error {
	return r.categoryTranslationRepo.DeleteTranslation(ctx, req)
}

func (r *CategoryRepo) CleanTranslations(ctx context.Context, tx *ent.Tx, categoryID uint32) error {
	return r.categoryTranslationRepo.CleanTranslations(ctx, tx, categoryID)
}
