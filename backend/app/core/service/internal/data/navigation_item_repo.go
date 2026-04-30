package data

import (
	"context"
	"sort"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	entCrud "github.com/tx7do/go-crud/entgo"

	"github.com/tx7do/go-utils/copierutil"
	"github.com/tx7do/go-utils/mapper"
	"github.com/tx7do/go-utils/timeutil"
	"github.com/tx7do/go-utils/trans"

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

func (r *NavigationItemRepo) List(ctx context.Context, req *paginationV1.PagingRequest) (*siteV1.ListNavigationItemResponse, error) {
	if req == nil {
		return nil, siteV1.ErrorBadRequest("invalid parameter")
	}

	builder := r.entClient.Client().NavigationItem.Query()

	ret, err := r.repository.ListWithPaging(ctx, builder, builder.Clone(), req)
	if err != nil {
		return nil, err
	}
	if ret == nil {
		return &siteV1.ListNavigationItemResponse{Total: 0, Items: nil}, nil
	}

	return &siteV1.ListNavigationItemResponse{
		Total: ret.Total,
		Items: ret.Items,
	}, nil
}

func (r *NavigationItemRepo) Get(ctx context.Context, req *siteV1.GetNavigationItemRequest) (*siteV1.NavigationItem, error) {
	builder := r.entClient.Client().NavigationItem.Query()

	switch req.QueryBy.(type) {
	case *siteV1.GetNavigationItemRequest_Id:
		builder = builder.Where(navigationitem.IDEQ(req.GetId()))
	case *siteV1.GetNavigationItemRequest_Title:
		builder = builder.Where(navigationitem.TitleEQ(req.GetTitle()))
	default:
		return nil, siteV1.ErrorBadRequest("invalid query by parameter")
	}

	entity, err := builder.Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, siteV1.ErrorFileNotFound("navigation item not found")
		}

		r.log.Errorf("query navigation item failed: %s", err.Error())

		return nil, siteV1.ErrorInternalServerError("query navigation item failed")
	}

	return r.mapper.ToDTO(entity), nil
}

func (r *NavigationItemRepo) newCreateBuilder(tx *ent.Tx, data *siteV1.NavigationItem) *ent.NavigationItemCreate {
	now := time.Now()
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
		SetNillableRequiredPermission(data.RequiredPermission).
		SetNillableParentID(data.ParentId).
		SetNillableCreatedBy(data.CreatedBy).
		SetCreatedAt(now)
	return builder
}

func (r *NavigationItemRepo) Create(ctx context.Context, data *siteV1.NavigationItem) (dto *siteV1.NavigationItem, err error) {
	if data == nil {
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

	builder := r.newCreateBuilder(tx, data)

	var entity *ent.NavigationItem
	entity, err = builder.Save(ctx)
	if err != nil {
		r.log.Errorf("create navigation item failed: %s", err.Error())
		return nil, siteV1.ErrorInternalServerError("create navigation item failed")
	}

	return r.mapper.ToDTO(entity), nil
}

func (r *NavigationItemRepo) BatchCreate(ctx context.Context, tx *ent.Tx, items []*siteV1.NavigationItem) error {
	if len(items) == 0 {
		return nil
	}

	builders := make([]*ent.NavigationItemCreate, 0, len(items))
	for _, data := range items {
		builder := r.newCreateBuilder(tx, data)
		builders = append(builders, builder)
	}

	err := tx.NavigationItem.CreateBulk(builders...).Exec(ctx)
	if err != nil {
		r.log.Errorf("batch create navigation items failed: %s", err.Error())
		return siteV1.ErrorInternalServerError("batch create navigation items failed")
	}

	return nil
}

func (r *NavigationItemRepo) newUpdateOneBuilder(tx *ent.Tx, data *siteV1.NavigationItem) *ent.NavigationItemUpdateOne {
	now := time.Now()
	builder := tx.NavigationItem.UpdateOneID(*data.Id).
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
		SetNillableRequiredPermission(data.RequiredPermission).
		SetNillableParentID(data.ParentId).
		SetNillableUpdatedBy(data.CreatedBy).
		SetUpdatedAt(now)
	return builder
}

func (r *NavigationItemRepo) Update(ctx context.Context, req *siteV1.UpdateNavigationItemRequest) (dto *siteV1.NavigationItem, err error) {
	if req == nil || req.Data == nil {
		return nil, siteV1.ErrorBadRequest("invalid parameter")
	}

	// 如果不存在则创建
	if req.GetAllowMissing() {
		var exist bool
		exist, err = r.IsExist(ctx, req.GetId())
		if err != nil {
			return nil, err
		}
		if !exist {
			req.Data.CreatedBy = req.Data.UpdatedBy
			req.Data.UpdatedBy = nil
			_, err = r.Create(ctx, req.Data)
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

	builder := r.newUpdateOneBuilder(tx, req.Data)

	var entity *ent.NavigationItem
	entity, err = builder.Save(ctx)
	if err != nil {
		r.log.Errorf("update navigation item failed: %s", err.Error())
		return nil, siteV1.ErrorInternalServerError("update navigation item failed")
	}

	return r.mapper.ToDTO(entity), nil
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

func (r *NavigationItemRepo) buildNavigationItemTree(items []*siteV1.NavigationItem, parentId uint32) []*siteV1.NavigationItem {
	var tree []*siteV1.NavigationItem
	for _, item := range items {
		if item.GetParentId() == parentId {
			// 递归查找子节点
			children := r.buildNavigationItemTree(items, item.GetId())
			item.Children = children
			tree = append(tree, item)
		}
	}
	return tree
}

func (r *NavigationItemRepo) ListItems(ctx context.Context, navigationID uint32, treeTravel bool) ([]*siteV1.NavigationItem, error) {
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

	if !treeTravel {
		dtos := make([]*siteV1.NavigationItem, 0, len(entities))
		for _, entity := range entities {
			dtos = append(dtos, r.mapper.ToDTO(entity))
		}
		return dtos, nil
	}

	grouped := make(map[uint32][]*siteV1.NavigationItem)
	for _, entity := range entities {
		pid := uint32(0)
		if entity.ParentID != nil {
			pid = *entity.ParentID
		}
		grouped[pid] = append(grouped[pid], r.mapper.ToDTO(entity))
	}

	sortByOrder := func(items []*siteV1.NavigationItem) {
		sort.SliceStable(items, func(i, j int) bool {
			if items[i].SortOrder == nil && items[j].SortOrder == nil {
				return false
			}
			if items[i].SortOrder == nil {
				return false
			}
			if items[j].SortOrder == nil {
				return true
			}
			return *items[i].SortOrder < *items[j].SortOrder
		})
	}

	var buildTree func(parentID uint32) []*siteV1.NavigationItem
	buildTree = func(parentID uint32) []*siteV1.NavigationItem {
		items := grouped[parentID]
		sortByOrder(items)
		for _, item := range items {
			item.Children = buildTree(item.GetId())
		}
		return items
	}

	return buildTree(0), nil
}

func (r *NavigationItemRepo) Delete(ctx context.Context, req *siteV1.DeleteNavigationItemRequest) (err error) {
	if req == nil {
		return siteV1.ErrorBadRequest("invalid parameter")
	}

	if err = r.entClient.Client().Navigation.DeleteOneID(req.GetId()).Exec(ctx); err != nil {
		r.log.Errorf("delete one data failed: %s", err.Error())
	}
	return err
}
