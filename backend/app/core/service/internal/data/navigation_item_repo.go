package data

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	entCrud "github.com/tx7do/go-crud/entgo"
	"github.com/tx7do/go-utils/copierutil"
	"github.com/tx7do/go-utils/mapper"
	"github.com/tx7do/go-utils/timeutil"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/navigationitem"
	"go-wind-cms/app/core/service/internal/data/ent/predicate"

	siteV1 "go-wind-cms/api/gen/go/site/service/v1"
)

type NavigationItemRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper

	mapper *mapper.CopierMapper[siteV1.NavigationItem, ent.NavigationItem]

	repository *entCrud.Repository[
		ent.NavigationItemQuery, ent.NavigationItemSelect,
		ent.NavigationItemCreate, ent.NavigationItemCreateBulk,
		ent.NavigationItemUpdate, ent.NavigationItemUpdateOne,
		ent.NavigationItemDelete,
		predicate.NavigationItem,
		siteV1.NavigationItem, ent.NavigationItem,
	]

	linkTypeConverter *mapper.EnumTypeConverter[siteV1.NavigationItem_LinkType, navigationitem.LinkType]
}

func NewNavigationItemRepo(ctx *bootstrap.Context, entClient *entCrud.EntClient[*ent.Client]) *NavigationItemRepo {
	repo := &NavigationItemRepo{
		entClient: entClient,
		log:       ctx.NewLoggerHelper("navigation-item/repo/core-service"),
		mapper:    mapper.NewCopierMapper[siteV1.NavigationItem, ent.NavigationItem](),
		linkTypeConverter: mapper.NewEnumTypeConverter[siteV1.NavigationItem_LinkType, navigationitem.LinkType](
			siteV1.NavigationItem_LinkType_name, siteV1.NavigationItem_LinkType_value,
		),
	}

	repo.init()

	return repo
}

func (r *NavigationItemRepo) init() {
	r.repository = entCrud.NewRepository[
		ent.NavigationItemQuery, ent.NavigationItemSelect,
		ent.NavigationItemCreate, ent.NavigationItemCreateBulk,
		ent.NavigationItemUpdate, ent.NavigationItemUpdateOne,
		ent.NavigationItemDelete,
		predicate.NavigationItem,
		siteV1.NavigationItem, ent.NavigationItem,
	](r.mapper)

	r.mapper.AppendConverters(copierutil.NewTimeStringConverterPair())
	r.mapper.AppendConverters(copierutil.NewTimeTimestamppbConverterPair())

	r.mapper.AppendConverters(r.linkTypeConverter.NewConverterPair())
}

func (r *NavigationItemRepo) IsExist(ctx context.Context, navigationID uint32) (bool, error) {
	exist, err := r.entClient.Client().NavigationItem.Query().
		Where(navigationitem.NavigationIDEQ(navigationID)).
		Exist(ctx)
	if err != nil {
		r.log.Errorf("query navigation item exist failed: %s", err.Error())
		return false, siteV1.ErrorInternalServerError("query navigation item exist failed")
	}
	return exist, nil
}

func (r *NavigationItemRepo) CleanItems(
	ctx context.Context,
	tx *ent.Tx,
	navigationID uint32,
) error {
	if _, err := tx.NavigationItem.Delete().
		Where(
			navigationitem.NavigationIDEQ(navigationID),
		).
		Exec(ctx); err != nil {
		r.log.Errorf("delete old navigation [%d] items failed: %s", navigationID, err.Error())
		return siteV1.ErrorInternalServerError("delete old navigation items failed")
	}
	return nil
}

func (r *NavigationItemRepo) Get(ctx context.Context, id uint32) (*siteV1.NavigationItem, error) {
	entity, err := r.entClient.Client().NavigationItem.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, siteV1.ErrorFileNotFound("navigation item not found")
		}

		r.log.Errorf("query navigation item failed: %s", err.Error())

		return nil, siteV1.ErrorInternalServerError("query navigation item failed")
	}

	return r.mapper.ToDTO(entity), nil
}

func (r *NavigationItemRepo) BatchCreate(ctx context.Context, tx *ent.Tx, items []*siteV1.NavigationItem) error {
	if len(items) == 0 {
		return nil
	}

	now := time.Now()

	builders := make([]*ent.NavigationItemCreate, 0, len(items))
	for _, data := range items {
		builder := tx.NavigationItem.Create().
			SetNillableNavigationID(data.NavigationId).
			SetNillableTitle(data.Title).
			SetNillableURL(data.Url).
			SetNillableIcon(data.Icon).
			SetNillableDescription(data.Description).
			SetNillableLinkType(r.linkTypeConverter.ToEntity(data.LinkType)).
			SetNillableObjectID(data.ObjectId).
			SetNillableSortOrder(data.SortOrder).
			SetNillableIsOpenNewTab(data.IsOpenNewTab).
			SetNillableIsInvalid(data.IsInvalid).
			SetNillableCSSClass(data.CssClass).
			SetNillableRequiredPermission(data.RequiredPermission).
			SetNillableParentID(data.ParentId).
			SetNillableCreatedBy(data.CreatedBy).
			SetCreatedAt(now)

		builders = append(builders, builder)
	}

	err := tx.NavigationItem.CreateBulk(builders...).Exec(ctx)
	if err != nil {
		r.log.Errorf("batch create navigation items failed: %s", err.Error())
		return siteV1.ErrorInternalServerError("batch create navigation items failed")
	}

	return nil
}

func (r *NavigationItemRepo) Upsert(ctx context.Context, data *siteV1.NavigationItem) error {
	if data == nil {
		return siteV1.ErrorBadRequest("invalid parameter")
	}

	var operatorID *uint32
	if data.UpdatedBy != nil {
		operatorID = data.UpdatedBy
	} else {
		operatorID = data.CreatedBy
	}
	if operatorID == nil {
		r.log.Errorf("operator ID is nil for upsert navigation item")
		return siteV1.ErrorInternalServerError("operator ID is required for upsert navigation item")
	}

	var now *time.Time
	if data.UpdatedAt != nil {
		t := timeutil.TimestamppbToTime(data.UpdatedAt)
		now = t
	} else if data.CreatedAt != nil {
		t := timeutil.TimestamppbToTime(data.CreatedAt)
		now = t
	}
	if now == nil {
		now = trans.Ptr(time.Now())
	}

	builder := r.entClient.Client().NavigationItem.Create().
		SetNillableNavigationID(data.NavigationId).
		SetNillableTitle(data.Title).
		SetNillableURL(data.Url).
		SetNillableIcon(data.Icon).
		SetNillableDescription(data.Description).
		SetNillableLinkType(r.linkTypeConverter.ToEntity(data.LinkType)).
		SetNillableObjectID(data.ObjectId).
		SetNillableSortOrder(data.SortOrder).
		SetNillableIsOpenNewTab(data.IsOpenNewTab).
		SetNillableIsInvalid(data.IsInvalid).
		SetNillableCSSClass(data.CssClass).
		SetNillableRequiredPermission(data.RequiredPermission).
		SetNillableParentID(data.ParentId).
		SetNillableCreatedBy(operatorID).
		SetNillableCreatedAt(now).
		OnConflictColumns(
			navigationitem.FieldNavigationID,
		).
		UpdateNewValues().
		SetUpdatedBy(*operatorID).
		SetUpdatedAt(*now)

	err := builder.Exec(ctx)
	if err != nil {
		r.log.Errorf("create navigation item failed: %s", err.Error())
		return siteV1.ErrorInternalServerError("create navigation item failed")
	}

	return nil
}

func (r *NavigationItemRepo) ListItems(ctx context.Context, navigationID uint32) ([]*siteV1.NavigationItem, error) {
	q := r.entClient.Client().NavigationItem.Query().
		Where(
			navigationitem.NavigationIDEQ(navigationID),
		)

	entities, err := q.
		All(ctx)
	if err != nil {
		r.log.Errorf("query navigation items by role id failed: %s", err.Error())
		return nil, siteV1.ErrorInternalServerError("query navigation items by navigation ID failed")
	}

	var dtos []*siteV1.NavigationItem
	for _, entity := range entities {
		dtos = append(dtos, r.mapper.ToDTO(entity))
	}

	return dtos, nil
}
