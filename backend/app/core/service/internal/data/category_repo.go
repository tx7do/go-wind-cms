package data

import (
	"context"
	siteV1 "go-wind-cms/api/gen/go/site/service/v1"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	entCrud "github.com/tx7do/go-crud/entgo"

	"github.com/tx7do/go-utils/copierutil"
	"github.com/tx7do/go-utils/mapper"

	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/category"
	"go-wind-cms/app/core/service/internal/data/ent/predicate"

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

func (r *CategoryRepo) List(ctx context.Context, req *paginationV1.PagingRequest) (*contentV1.ListCategoryResponse, error) {
	if req == nil {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}

	builder := r.entClient.Client().Category.Query()

	ret, err := r.repository.ListWithPaging(ctx, builder, builder.Clone(), req)
	if err != nil {
		return nil, err
	}
	if ret == nil {
		return &contentV1.ListCategoryResponse{Total: 0, Items: nil}, nil
	}

	for _, item := range ret.Items {
		translations, err := r.categoryTranslationRepo.ListTranslations(ctx, item.GetId())
		if err != nil {
			r.log.Errorf("query translations failed: %s", err.Error())
			return nil, contentV1.ErrorInternalServerError("query translations failed")
		}
		item.Translations = translations
	}

	return &contentV1.ListCategoryResponse{
		Total: ret.Total,
		Items: ret.Items,
	}, nil
}

func (r *CategoryRepo) Get(ctx context.Context, req *contentV1.GetCategoryRequest) (*contentV1.Category, error) {
	if req == nil {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}

	entity, err := r.entClient.Client().Category.Get(ctx, req.GetId())
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, contentV1.ErrorFileNotFound("category not found")
		}

		r.log.Errorf("query category failed: %s", err.Error())

		return nil, contentV1.ErrorInternalServerError("query category failed")
	}

	dto := r.mapper.ToDTO(entity)

	translations, err := r.categoryTranslationRepo.ListTranslations(ctx, dto.GetId())
	if err != nil {
		r.log.Errorf("query translations failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("query translations failed")
	}
	dto.Translations = translations

	return dto, nil
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
		SetNillablePostCount(req.Data.PostCount).
		SetNillableDirectPostCount(req.Data.DirectPostCount).
		SetNillableParentID(req.Data.ParentId).
		SetNillableDepth(req.Data.Depth).
		SetNillablePath(req.Data.Path).
		SetNillableCreatedBy(req.Data.CreatedBy).
		SetCreatedAt(time.Now())

	var entity *ent.Category
	if entity, err = builder.Save(ctx); err != nil {
		r.log.Errorf("insert category failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("insert category failed")
	}

	if req.Data.Translations != nil {
		if err = r.categoryTranslationRepo.CleanTranslations(ctx, tx, entity.ID); err != nil {
			r.log.Errorf("clean translations failed: %s", err.Error())
			return nil, contentV1.ErrorInternalServerError("clean translations failed")
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

	if req.Data.Translations != nil {
		if err = r.categoryTranslationRepo.CleanTranslations(ctx, tx, req.GetId()); err != nil {
			r.log.Errorf("clean translations failed: %s", err.Error())
			return nil, contentV1.ErrorInternalServerError("clean translations failed")
		}
		if err = r.categoryTranslationRepo.BatchCreate(ctx, tx, req.Data.GetTranslations()); err != nil {
			r.log.Errorf("batch insert translations failed: %s", err.Error())
			return nil, contentV1.ErrorInternalServerError("batch insert translations failed")
		}
	}

	builder := tx.Category.UpdateOneID(req.GetId())
	result, err := r.repository.UpdateOne(ctx, builder, req.Data, req.GetUpdateMask(),
		func(dto *contentV1.Category) {
			builder.
				SetNillableStatus(r.statusConverter.ToEntity(req.Data.Status)).
				SetNillableSortOrder(req.Data.SortOrder).
				SetNillableIsNav(req.Data.IsNav).
				SetNillableIcon(req.Data.Icon).
				SetNillablePostCount(req.Data.PostCount).
				SetNillableDirectPostCount(req.Data.DirectPostCount).
				SetNillableParentID(req.Data.ParentId).
				SetNillableDepth(req.Data.Depth).
				SetNillablePath(req.Data.Path).
				SetNillableUpdatedBy(req.Data.UpdatedBy).
				SetUpdatedAt(time.Now())
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
		return siteV1.ErrorInternalServerError("start transaction failed")
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

	if err = tx.Category.
		DeleteOneID(req.GetId()).
		Exec(ctx); err != nil {
		r.log.Errorf("delete one data failed: %s", err.Error())
	}

	if err = r.categoryTranslationRepo.CleanTranslations(ctx, tx, req.GetId()); err != nil {
		r.log.Errorf("clean translations failed: %s", err.Error())
		return contentV1.ErrorInternalServerError("clean translations failed")
	}

	return err
}
