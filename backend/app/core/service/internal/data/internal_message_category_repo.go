package data

import (
	"context"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	entCrud "github.com/tx7do/go-crud/entgo"

	"github.com/tx7do/go-utils/copierutil"
	"github.com/tx7do/go-utils/mapper"

	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/internalmessagecategory"
	"go-wind-cms/app/core/service/internal/data/ent/predicate"

	internalMessageV1 "go-wind-cms/api/gen/go/internal_message/service/v1"
)

type InternalMessageCategoryRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper

	mapper *mapper.CopierMapper[internalMessageV1.InternalMessageCategory, ent.InternalMessageCategory]

	repository *entCrud.Repository[
		ent.InternalMessageCategoryQuery, ent.InternalMessageCategorySelect,
		ent.InternalMessageCategoryCreate, ent.InternalMessageCategoryCreateBulk,
		ent.InternalMessageCategoryUpdate, ent.InternalMessageCategoryUpdateOne,
		ent.InternalMessageCategoryDelete,
		predicate.InternalMessageCategory,
		internalMessageV1.InternalMessageCategory, ent.InternalMessageCategory,
	]
}

func NewInternalMessageCategoryRepo(ctx *bootstrap.Context, entClient *entCrud.EntClient[*ent.Client]) *InternalMessageCategoryRepo {
	repo := &InternalMessageCategoryRepo{
		log:       ctx.NewLoggerHelper("internal-message-category/repo/core-service"),
		entClient: entClient,
		mapper:    mapper.NewCopierMapper[internalMessageV1.InternalMessageCategory, ent.InternalMessageCategory](),
	}

	repo.init()

	return repo
}

func (r *InternalMessageCategoryRepo) init() {
	r.repository = entCrud.NewRepository[
		ent.InternalMessageCategoryQuery, ent.InternalMessageCategorySelect,
		ent.InternalMessageCategoryCreate, ent.InternalMessageCategoryCreateBulk,
		ent.InternalMessageCategoryUpdate, ent.InternalMessageCategoryUpdateOne,
		ent.InternalMessageCategoryDelete,
		predicate.InternalMessageCategory,
		internalMessageV1.InternalMessageCategory, ent.InternalMessageCategory,
	](r.mapper)

	r.mapper.AppendConverters(copierutil.NewTimeStringConverterPair())
	r.mapper.AppendConverters(copierutil.NewTimeTimestamppbConverterPair())
}

func (r *InternalMessageCategoryRepo) count(ctx context.Context, whereCond []func(s *sql.Selector)) (int, error) {
	builder := r.entClient.Client().InternalMessageCategory.Query()
	if len(whereCond) != 0 {
		builder.Modify(whereCond...)
	}

	count, err := builder.Count(ctx)
	if err != nil {
		r.log.Errorf("query count failed: %s", err.Error())
		return 0, internalMessageV1.ErrorInternalServerError("query count failed")
	}

	return count, nil
}

func (r *InternalMessageCategoryRepo) Count(ctx context.Context, req *paginationV1.PagingRequest) (int, error) {
	builder := r.entClient.Client().InternalMessageCategory.Query()

	whereSelectors, _, err := r.repository.BuildListSelectorWithPaging(builder, req)
	if len(whereSelectors) != 0 {
		builder.Modify(whereSelectors...)
	}

	count, err := builder.Count(ctx)
	if err != nil {
		r.log.Errorf("query count failed: %s", err.Error())
		return 0, internalMessageV1.ErrorInternalServerError("query count failed")
	}

	return count, nil
}

func (r *InternalMessageCategoryRepo) List(ctx context.Context, req *paginationV1.PagingRequest) (*internalMessageV1.ListInternalMessageCategoryResponse, error) {
	if req == nil {
		return nil, internalMessageV1.ErrorBadRequest("invalid parameter")
	}

	builder := r.entClient.Client().InternalMessageCategory.Query()

	// field mask 含 children 时以树形结构返回（参照 CategoryRepo 的 treeTravel）
	treeTravel := false
	if req.FieldMask != nil && len(req.FieldMask.Paths) > 0 {
		excludeFields := []string{"children"}
		for _, p := range req.FieldMask.Paths {
			if strings.TrimSpace(p) == "children" {
				treeTravel = true
			}
		}
		if treeTravel {
			req.FieldMask = FilterViewMask(excludeFields, req.FieldMask)
		}
	}

	ret, err := r.repository.ListWithPaging(ctx, builder, builder.Clone(), req)
	if err != nil {
		return nil, err
	}
	if ret == nil {
		return &internalMessageV1.ListInternalMessageCategoryResponse{Total: 0, Items: nil}, nil
	}

	if treeTravel {
		roots := r.buildTree(ret.Items, 0) // 根节点 ParentId 为 0
		return &internalMessageV1.ListInternalMessageCategoryResponse{
			Total: ret.Total,
			Items: roots,
		}, nil
	}

	return &internalMessageV1.ListInternalMessageCategoryResponse{
		Total: ret.Total,
		Items: ret.Items,
	}, nil
}

// buildTree 递归组装树形结构（参照 CategoryRepo.buildCategoryTree）
func (r *InternalMessageCategoryRepo) buildTree(items []*internalMessageV1.InternalMessageCategory, parentId uint32) []*internalMessageV1.InternalMessageCategory {
	var tree []*internalMessageV1.InternalMessageCategory
	for _, item := range items {
		if item.GetParentId() == parentId {
			// 递归查找子节点
			children := r.buildTree(items, item.GetId())
			item.Children = children
			tree = append(tree, item)
		}
	}
	return tree
}

func (r *InternalMessageCategoryRepo) IsExist(ctx context.Context, id uint32) (bool, error) {
	exist, err := r.entClient.Client().InternalMessageCategory.Query().
		Where(internalmessagecategory.IDEQ(id)).
		Exist(ctx)
	if err != nil {
		r.log.Errorf("query exist failed: %s", err.Error())
		return false, internalMessageV1.ErrorInternalServerError("query exist failed")
	}
	return exist, nil
}

func (r *InternalMessageCategoryRepo) Get(ctx context.Context, req *internalMessageV1.GetInternalMessageCategoryRequest) (*internalMessageV1.InternalMessageCategory, error) {
	if req == nil {
		return nil, internalMessageV1.ErrorBadRequest("invalid parameter")
	}

	builder := r.entClient.Client().InternalMessageCategory.Query()

	var whereCond []func(s *sql.Selector)
	switch req.QueryBy.(type) {
	default:
	case *internalMessageV1.GetInternalMessageCategoryRequest_Id:
		whereCond = append(whereCond, internalmessagecategory.IDEQ(req.GetId()))
	}

	// view mask 含 children 时递归填充子树（参照 CategoryRepo 的 treeTravel）
	treeTravel := false
	if req.ViewMask != nil && len(req.ViewMask.Paths) > 0 {
		excludeFields := []string{"children"}
		for _, p := range req.ViewMask.Paths {
			if strings.TrimSpace(p) == "children" {
				treeTravel = true
			}
		}
		if treeTravel {
			req.ViewMask = FilterViewMask(excludeFields, req.ViewMask)
		}
	}

	dto, err := r.repository.Get(ctx, builder, req.GetViewMask(), whereCond...)
	if err != nil {
		return nil, err
	}

	if treeTravel {
		return r.getWithChildren(ctx, dto.GetId())
	}

	return dto, err
}

// getWithChildren 递归查询节点及其子树（参照 CategoryRepo.getCategoryWithChildren，无翻译层）
func (r *InternalMessageCategoryRepo) getWithChildren(ctx context.Context, id uint32) (*internalMessageV1.InternalMessageCategory, error) {
	entity, err := r.entClient.Client().InternalMessageCategory.Query().
		Where(internalmessagecategory.IDEQ(id)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, internalMessageV1.ErrorNotFound("internal message category not found")
		}
		r.log.Errorf("query internal message category failed: %s", err.Error())
		return nil, internalMessageV1.ErrorInternalServerError("query internal message category failed")
	}

	dto := r.mapper.ToDTO(entity)

	// 查询子节点
	childrenEntities, err := r.entClient.Client().InternalMessageCategory.Query().
		Where(internalmessagecategory.ParentIDEQ(id)).
		All(ctx)
	if err != nil {
		r.log.Errorf("query children failed: %s", err.Error())
		return nil, internalMessageV1.ErrorInternalServerError("query children failed")
	}
	for _, child := range childrenEntities {
		childDTO, err := r.getWithChildren(ctx, child.ID)
		if err != nil {
			return nil, err
		}
		dto.Children = append(dto.Children, childDTO)
	}

	return dto, nil
}

// ListCategoriesByIds 根据ID列表获取分类列表
func (r *InternalMessageCategoryRepo) ListCategoriesByIds(ctx context.Context, ids []uint32) ([]*internalMessageV1.InternalMessageCategory, error) {
	if len(ids) == 0 {
		return []*internalMessageV1.InternalMessageCategory{}, nil
	}

	entities, err := r.entClient.Client().InternalMessageCategory.Query().
		Where(internalmessagecategory.IDIn(ids...)).
		All(ctx)
	if err != nil {
		r.log.Errorf("query internal message category by ids failed: %s", err.Error())
		return nil, internalMessageV1.ErrorInternalServerError("query internal message category by ids failed")
	}

	dtos := make([]*internalMessageV1.InternalMessageCategory, 0, len(entities))
	for _, entity := range entities {
		dto := r.mapper.ToDTO(entity)
		dtos = append(dtos, dto)
	}

	return dtos, nil
}

// ensureParentAcyclic 校验 parentId 不构成环：沿 parent 链上溯，命中 selfID 即拒绝。
// 数据成环会让 buildTree/getWithChildren 的递归无限循环，参照实现（content Category）
// 没有此防护，这里补上。
func (r *InternalMessageCategoryRepo) ensureParentAcyclic(ctx context.Context, selfID, parentID uint32) error {
	if parentID == 0 {
		return nil
	}
	if parentID == selfID {
		return internalMessageV1.ErrorBadRequest("parent cannot be itself")
	}
	current := parentID
	for current != 0 {
		entity, err := r.entClient.Client().InternalMessageCategory.Query().
			Where(internalmessagecategory.IDEQ(current)).
			Only(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				// 父节点不存在，交由外键约束/后续查询兜底
				return nil
			}
			r.log.Errorf("query parent for cycle check failed: %s", err.Error())
			return internalMessageV1.ErrorInternalServerError("query parent failed")
		}
		if entity.ParentID != nil {
			if *entity.ParentID == selfID {
				return internalMessageV1.ErrorBadRequest("parent cannot be itself or its descendant")
			}
			current = *entity.ParentID
		} else {
			return nil
		}
	}
	return nil
}

func (r *InternalMessageCategoryRepo) Create(ctx context.Context, req *internalMessageV1.CreateInternalMessageCategoryRequest) error {
	if req == nil || req.Data == nil {
		return internalMessageV1.ErrorBadRequest("invalid parameter")
	}

	if err := r.ensureParentAcyclic(ctx, req.GetData().GetId(), req.Data.GetParentId()); err != nil {
		return err
	}

	builder := r.entClient.Client().InternalMessageCategory.Create().
		SetNillableTenantID(req.Data.TenantId).
		SetNillableName(req.Data.Name).
		SetNillableCode(req.Data.Code).
		SetNillableIconURL(req.Data.IconUrl).
		SetNillableSortOrder(req.Data.SortOrder).
		SetNillableIsEnabled(req.Data.IsEnabled).
		SetNillableParentID(req.Data.ParentId).
		SetNillableDepth(req.Data.Depth).
		SetNillablePath(req.Data.Path).
		SetNillableCreatedBy(req.Data.CreatedBy).
		SetCreatedAt(time.Now())

	if req.Data.Id != nil {
		builder.SetID(req.GetData().GetId())
	}

	if err := builder.Exec(ctx); err != nil {
		r.log.Errorf("insert internal message category failed: %s", err.Error())
		return internalMessageV1.ErrorInternalServerError("insert internal message category failed")
	}

	return nil
}

func (r *InternalMessageCategoryRepo) Update(ctx context.Context, req *internalMessageV1.UpdateInternalMessageCategoryRequest) error {
	if req == nil || req.Data == nil {
		return internalMessageV1.ErrorBadRequest("invalid parameter")
	}

	// 如果不存在则创建
	if req.GetAllowMissing() {
		exist, err := r.IsExist(ctx, req.GetId())
		if err != nil {
			return err
		}
		if !exist {
			createReq := &internalMessageV1.CreateInternalMessageCategoryRequest{Data: req.Data}
			createReq.Data.CreatedBy = createReq.Data.UpdatedBy
			createReq.Data.UpdatedBy = nil
			return r.Create(ctx, createReq)
		}
	}

	tid, hasTenant := maybeTenantFromViewer(ctx)
	callerUserID, hasUser := viewerUserIDFromContext(ctx)

	// 变更 parentId 时校验不成环
	if req.Data.ParentId != nil {
		if err := r.ensureParentAcyclic(ctx, req.GetId(), req.Data.GetParentId()); err != nil {
			return err
		}
	}

	builder := r.entClient.Client().Debug().InternalMessageCategory.Update()
	builder.Where(internalmessagecategory.IDEQ(req.GetId()))
	if hasTenant {
		builder.Where(internalmessagecategory.TenantIDEQ(tid))
	}
	err := r.repository.UpdateX(ctx, builder, req.Data, req.GetUpdateMask(),
		func(dto *internalMessageV1.InternalMessageCategory) {
			builder.
				SetNillableName(req.Data.Name).
				SetNillableCode(req.Data.Code).
				SetNillableIconURL(req.Data.IconUrl).
				SetNillableSortOrder(req.Data.SortOrder).
				SetNillableIsEnabled(req.Data.IsEnabled).
				SetNillableParentID(req.Data.ParentId).
				SetNillableDepth(req.Data.Depth).
				SetNillablePath(req.Data.Path).
				SetUpdatedAt(time.Now())

			// updated_by 强制由服务端 viewer context 推导，忽略客户端传入值
			if hasUser {
				builder.SetUpdatedBy(callerUserID)
			}
		},
		func(s *sql.Selector) {
			s.Where(sql.EQ(internalmessagecategory.FieldID, req.GetId()))
		},
	)

	return err
}

func (r *InternalMessageCategoryRepo) Delete(ctx context.Context, req *internalMessageV1.DeleteInternalMessageCategoryRequest) error {
	if req == nil {
		return internalMessageV1.ErrorBadRequest("invalid parameter")
	}

	childrenIds, err := entCrud.QueryAllChildrenIds(ctx, r.entClient, "internal_message_categories", req.GetId())
	if err != nil {
		r.log.Errorf("query child internal message categories failed: %s", err.Error())
		return internalMessageV1.ErrorInternalServerError("query child internal message categories failed")
	}
	childrenIds = append(childrenIds, req.GetId())

	//r.log.Info("internal message category childrenIds to delete: ", childrenIds)

	var ids []any
	for _, id := range childrenIds {
		ids = append(ids, id)
	}

	builder := r.entClient.Client().Debug().InternalMessageCategory.Delete()

	deletePreds := []predicate.InternalMessageCategory{
		func(s *sql.Selector) {
			s.Where(sql.In(internalmessagecategory.FieldID, ids...))
		},
	}
	tid, hasTenant := maybeTenantFromViewer(ctx)
	if hasTenant {
		deletePreds = append(deletePreds, internalmessagecategory.TenantIDEQ(tid))
	}

	_, err = r.repository.Delete(ctx, builder, deletePreds...)
	if err != nil {
		r.log.Errorf("delete internal message categories failed: %s", err.Error())
		return internalMessageV1.ErrorInternalServerError("delete internal message categories failed")
	}

	return nil
}
