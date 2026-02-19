package data

import (
	"context"

	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	entCrud "github.com/tx7do/go-crud/entgo"
	"github.com/tx7do/go-utils/copierutil"
	"github.com/tx7do/go-utils/mapper"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/navigation"
	"go-wind-cms/app/core/service/internal/data/ent/predicate"

	siteV1 "go-wind-cms/api/gen/go/site/service/v1"
)

type NavigationRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper

	mapper *mapper.CopierMapper[siteV1.Navigation, ent.Navigation]

	repository *entCrud.Repository[
		ent.NavigationQuery, ent.NavigationSelect,
		ent.NavigationCreate, ent.NavigationCreateBulk,
		ent.NavigationUpdate, ent.NavigationUpdateOne,
		ent.NavigationDelete,
		predicate.Navigation,
		siteV1.Navigation, ent.Navigation,
	]

	navigationItemRepo *NavigationItemRepo
}

func NewNavigationRepo(
	ctx *bootstrap.Context,
	entClient *entCrud.EntClient[*ent.Client],
	navigationItemRepo *NavigationItemRepo,
) *NavigationRepo {
	repo := &NavigationRepo{
		entClient:          entClient,
		log:                ctx.NewLoggerHelper("navigation/repo/core-service"),
		mapper:             mapper.NewCopierMapper[siteV1.Navigation, ent.Navigation](),
		navigationItemRepo: navigationItemRepo,
	}

	repo.init()

	return repo
}

func (r *NavigationRepo) init() {
	r.repository = entCrud.NewRepository[
		ent.NavigationQuery, ent.NavigationSelect,
		ent.NavigationCreate, ent.NavigationCreateBulk,
		ent.NavigationUpdate, ent.NavigationUpdateOne,
		ent.NavigationDelete,
		predicate.Navigation,
		siteV1.Navigation, ent.Navigation,
	](r.mapper)

	r.mapper.AppendConverters(copierutil.NewTimeStringConverterPair())
	r.mapper.AppendConverters(copierutil.NewTimeTimestamppbConverterPair())
}

func (r *NavigationRepo) IsExist(ctx context.Context, id uint32) (bool, error) {
	exist, err := r.entClient.Client().Navigation.Query().
		Where(navigation.IDEQ(id)).
		Exist(ctx)
	if err != nil {
		r.log.Errorf("query navigation exist failed: %s", err.Error())
		return false, siteV1.ErrorInternalServerError("query navigation exist failed")
	}
	return exist, nil
}

func (r *NavigationRepo) count(ctx context.Context, whereCond []func(s *sql.Selector)) (int, error) {
	builder := r.entClient.Client().Navigation.Query()
	if len(whereCond) != 0 {
		builder.Modify(whereCond...)
	}

	count, err := builder.Count(ctx)
	if err != nil {
		r.log.Errorf("query count failed: %s", err.Error())
	}

	return count, err
}

func (r *NavigationRepo) List(ctx context.Context, req *paginationV1.PagingRequest) (*siteV1.ListNavigationResponse, error) {
	if req == nil {
		return nil, siteV1.ErrorBadRequest("invalid parameter")
	}

	builder := r.entClient.Client().Navigation.Query()

	ret, err := r.repository.ListWithPaging(ctx, builder, builder.Clone(), req)
	if err != nil {
		return nil, err
	}
	if ret == nil {
		return &siteV1.ListNavigationResponse{Total: 0, Items: nil}, nil
	}

	for _, item := range ret.Items {
		navigationItems, err := r.navigationItemRepo.ListItems(ctx, item.GetId())
		if err != nil {
			r.log.Errorf("query navigation items failed: %s", err.Error())
			return nil, siteV1.ErrorInternalServerError("query navigation items failed")
		}
		item.Items = navigationItems
	}

	return &siteV1.ListNavigationResponse{
		Total: ret.Total,
		Items: ret.Items,
	}, nil
}

func (r *NavigationRepo) Get(ctx context.Context, req *siteV1.GetNavigationRequest) (*siteV1.Navigation, error) {
	if req == nil {
		return nil, siteV1.ErrorBadRequest("invalid parameter")
	}

	entity, err := r.entClient.Client().Navigation.Get(ctx, req.GetId())
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, siteV1.ErrorFileNotFound("navigation not found")
		}

		r.log.Errorf("query navigation failed: %s", err.Error())

		return nil, siteV1.ErrorInternalServerError("query navigation failed")
	}

	dto := r.mapper.ToDTO(entity)

	navigationItems, err := r.navigationItemRepo.ListItems(ctx, dto.GetId())
	if err != nil {
		r.log.Errorf("query navigation items failed: %s", err.Error())
		return nil, siteV1.ErrorInternalServerError("query navigation items failed")
	}
	dto.Items = navigationItems

	return dto, nil
}

func (r *NavigationRepo) Create(ctx context.Context, req *siteV1.CreateNavigationRequest) (dto *siteV1.Navigation, err error) {
	if req == nil || req.Data == nil {
		return nil, siteV1.ErrorBadRequest("invalid parameter")
	}

	var tx *ent.Tx
	tx, err = r.entClient.Client().Tx(ctx)
	if err != nil {
		r.log.Errorf("start transaction failed: %s", err.Error())
		return nil, siteV1.ErrorInternalServerError("start transaction failed")
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
			err = siteV1.ErrorInternalServerError("transaction commit failed")
		}
	}()

	builder := tx.Navigation.Create().
		SetNillableName(req.Data.Name).
		SetNillableLocation(req.Data.Location).
		SetNillableIsActive(req.Data.IsActive).
		SetNillableLocale(req.Data.Locale).
		SetNillableCreatedBy(req.Data.CreatedBy).
		SetCreatedAt(time.Now())

	var entity *ent.Navigation
	if entity, err = builder.Save(ctx); err != nil {
		r.log.Errorf("insert navigation failed: %s", err.Error())
		return nil, siteV1.ErrorInternalServerError("insert navigation failed")
	}

	if req.Data.Items != nil {
		if err = r.navigationItemRepo.CleanItems(ctx, tx, entity.ID); err != nil {
			r.log.Errorf("clean navigation items failed: %s", err.Error())
			return nil, siteV1.ErrorInternalServerError("clean navigation items failed")
		}
		if err = r.navigationItemRepo.BatchCreate(ctx, tx, req.Data.GetItems()); err != nil {
			r.log.Errorf("batch insert navigation items failed: %s", err.Error())
			return nil, siteV1.ErrorInternalServerError("batch insert navigation items failed")
		}
	}

	return r.mapper.ToDTO(entity), nil
}

func (r *NavigationRepo) Update(ctx context.Context, req *siteV1.UpdateNavigationRequest) (dto *siteV1.Navigation, err error) {
	if req == nil || req.Data == nil {
		return nil, siteV1.ErrorBadRequest("invalid parameter")
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
			_, err = r.Create(ctx, &siteV1.CreateNavigationRequest{Data: req.Data})
			return nil, err
		}
	}

	var tx *ent.Tx
	tx, err = r.entClient.Client().Tx(ctx)
	if err != nil {
		r.log.Errorf("start transaction failed: %s", err.Error())
		return nil, siteV1.ErrorInternalServerError("start transaction failed")
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
			err = siteV1.ErrorInternalServerError("transaction commit failed")
		}
	}()

	if req.Data.Items != nil {
		if err = r.navigationItemRepo.CleanItems(ctx, tx, req.GetId()); err != nil {
			r.log.Errorf("clean navigation items failed: %s", err.Error())
			return nil, siteV1.ErrorInternalServerError("clean navigation items failed")
		}
		if err = r.navigationItemRepo.BatchCreate(ctx, tx, req.Data.GetItems()); err != nil {
			r.log.Errorf("batch insert navigation items failed: %s", err.Error())
			return nil, siteV1.ErrorInternalServerError("batch insert navigation items failed")
		}
	}

	builder := tx.Navigation.UpdateOneID(req.GetId())
	result, err := r.repository.UpdateOne(ctx, builder, req.Data, req.GetUpdateMask(),
		func(dto *siteV1.Navigation) {
			builder.
				SetNillableName(req.Data.Name).
				SetNillableLocation(req.Data.Location).
				SetNillableIsActive(req.Data.IsActive).
				SetNillableLocale(req.Data.Locale).
				SetNillableUpdatedBy(req.Data.UpdatedBy).
				SetUpdatedAt(time.Now())
		},
		func(s *sql.Selector) {
			s.Where(sql.EQ(navigation.FieldID, req.GetId()))
		},
	)

	return result, err
}

func (r *NavigationRepo) Delete(ctx context.Context, req *siteV1.DeleteNavigationRequest) (err error) {
	if req == nil {
		return siteV1.ErrorBadRequest("invalid parameter")
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
			err = siteV1.ErrorInternalServerError("transaction commit failed")
		}
	}()

	if err = r.entClient.Client().Navigation.DeleteOneID(req.GetId()).Exec(ctx); err != nil {
		r.log.Errorf("delete one data failed: %s", err.Error())
	}

	if err = r.navigationItemRepo.CleanItems(ctx, tx, req.GetId()); err != nil {
		r.log.Errorf("clean navigation items failed: %s", err.Error())
		return siteV1.ErrorInternalServerError("clean navigation items failed")
	}

	return err
}
