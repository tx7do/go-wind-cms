package data

import (
	"context"

	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	entCrud "github.com/tx7do/go-crud/entgo"

	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/api"
	"go-wind-cms/app/core/service/internal/data/ent/predicate"

	"github.com/tx7do/go-utils/copierutil"
	"github.com/tx7do/go-utils/mapper"

	permissionV1 "go-wind-cms/api/gen/go/permission/service/v1"
)

type ApiRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper

	mapper         *mapper.CopierMapper[permissionV1.Api, ent.Api]
	scopeConverter *mapper.EnumTypeConverter[permissionV1.Api_Scope, api.Scope]

	repository *entCrud.Repository[
		ent.APIQuery, ent.APISelect,
		ent.APICreate, ent.APICreateBulk,
		ent.APIUpdate, ent.APIUpdateOne,
		ent.APIDelete,
		predicate.Api,
		permissionV1.Api, ent.Api,
	]
}

func NewApiRepo(ctx *bootstrap.Context, entClient *entCrud.EntClient[*ent.Client]) *ApiRepo {
	repo := &ApiRepo{
		log:       ctx.NewLoggerHelper("api/repo/core-service"),
		entClient: entClient,
		mapper:    mapper.NewCopierMapper[permissionV1.Api, ent.Api](),
		scopeConverter: mapper.NewEnumTypeConverter[permissionV1.Api_Scope, api.Scope](
			permissionV1.Api_Scope_name, permissionV1.Api_Scope_value,
		),
	}

	repo.init()

	return repo
}

func (r *ApiRepo) init() {
	r.repository = entCrud.NewRepository[
		ent.APIQuery, ent.APISelect,
		ent.APICreate, ent.APICreateBulk,
		ent.APIUpdate, ent.APIUpdateOne,
		ent.APIDelete,
		predicate.Api,
		permissionV1.Api, ent.Api,
	](r.mapper)

	r.mapper.AppendConverters(copierutil.NewTimeStringConverterPair())
	r.mapper.AppendConverters(copierutil.NewTimeTimestamppbConverterPair())

	r.mapper.AppendConverters(r.scopeConverter.NewConverterPair())
}

func (r *ApiRepo) Count(ctx context.Context, req *paginationV1.PagingRequest) (*permissionV1.CountApiResponse, error) {
	builder := r.entClient.Client().Api.Query()

	whereSelectors, _, selErr := r.repository.BuildListSelectorWithPaging(builder, req)
	if selErr != nil {
		r.log.Errorf("build list selector failed: %s", selErr.Error())
		return nil, permissionV1.ErrorInternalServerError("build list selector failed")
	}
	if len(whereSelectors) != 0 {
		builder.Modify(whereSelectors...)
	}

	count, err := builder.Count(ctx)
	if err != nil {
		r.log.Errorf("query api count failed: %s", err.Error())
		return nil, permissionV1.ErrorInternalServerError("query api count failed")
	}

	return &permissionV1.CountApiResponse{
		Count: uint64(count),
	}, nil
}

func (r *ApiRepo) List(ctx context.Context, req *paginationV1.PagingRequest) (*permissionV1.ListApiResponse, error) {
	if req == nil {
		return nil, permissionV1.ErrorBadRequest("invalid parameter")
	}

	builder := r.entClient.Client().Api.Query()

	ret, err := r.repository.ListWithPaging(ctx, builder, builder.Clone(), req)
	if err != nil {
		return nil, err
	}
	if ret == nil {
		return &permissionV1.ListApiResponse{Total: 0, Items: nil}, nil
	}

	return &permissionV1.ListApiResponse{
		Total: ret.Total,
		Items: ret.Items,
	}, nil
}

func (r *ApiRepo) IsExist(ctx context.Context, id uint32) (bool, error) {
	exist, err := r.entClient.Client().Api.Query().
		Where(api.IDEQ(id)).
		Exist(ctx)
	if err != nil {
		r.log.Errorf("query exist failed: %s", err.Error())
		return false, permissionV1.ErrorInternalServerError("query exist failed")
	}
	return exist, nil
}

func (r *ApiRepo) Get(ctx context.Context, req *permissionV1.GetApiRequest) (*permissionV1.Api, error) {
	if req == nil {
		return nil, permissionV1.ErrorBadRequest("invalid parameter")
	}

	builder := r.entClient.Client().Api.Query()

	var whereCond []func(s *sql.Selector)
	switch req.QueryBy.(type) {
	default:
	case *permissionV1.GetApiRequest_Id:
		whereCond = append(whereCond, api.IDEQ(req.GetId()))
	}

	dto, err := r.repository.Get(ctx, builder, req.GetViewMask(), whereCond...)
	if err != nil {
		return nil, err
	}

	return dto, err
}

// GetApiByEndpoint 根据路径和方法获取API资源
func (r *ApiRepo) GetApiByEndpoint(ctx context.Context, path, method string) (*permissionV1.Api, error) {
	if path == "" || method == "" {
		return nil, permissionV1.ErrorBadRequest("invalid parameter")
	}

	entity, err := r.entClient.Client().Api.Query().
		Where(
			api.PathEQ(path),
			api.MethodEQ(method),
		).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, permissionV1.ErrorNotFound("api not found")
		}

		r.log.Errorf("query one data failed: %s", err.Error())

		return nil, permissionV1.ErrorInternalServerError("query data failed")
	}

	return r.mapper.ToDTO(entity), nil
}

// GetApiByIDs 根据ID列表获取API资源
func (r *ApiRepo) GetApiByIDs(ctx context.Context, ids []uint32) ([]*permissionV1.Api, error) {
	if len(ids) == 0 {
		return nil, permissionV1.ErrorBadRequest("invalid parameter")
	}

	entities, err := r.entClient.Client().Api.Query().
		Where(
			api.IDIn(ids...),
		).
		All(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, permissionV1.ErrorNotFound("api not found")
		}

		r.log.Errorf("query one data failed: %s", err.Error())

		return nil, permissionV1.ErrorInternalServerError("query data failed")
	}

	dtos := make([]*permissionV1.Api, 0, len(entities))
	for _, entity := range entities {
		dto := r.mapper.ToDTO(entity)
		dtos = append(dtos, dto)
	}

	return dtos, nil
}

func (r *ApiRepo) Create(ctx context.Context, req *permissionV1.CreateApiRequest) error {
	if req == nil || req.Data == nil {
		return permissionV1.ErrorBadRequest("invalid parameter")
	}

	builder := r.newApiCreate(req.Data)

	if err := builder.Exec(ctx); err != nil {
		r.log.Errorf("insert api failed: %s", err.Error())
		return permissionV1.ErrorInternalServerError("insert api failed")
	}

	return nil
}

func (r *ApiRepo) newApiCreate(api *permissionV1.Api) *ent.APICreate {
	builder := r.entClient.Client().Api.Create().
		SetNillableDescription(api.Description).
		SetNillableModule(api.Module).
		SetNillableModuleDescription(api.ModuleDescription).
		SetNillableOperation(api.Operation).
		SetNillablePath(api.Path).
		SetNillableMethod(api.Method).
		SetNillableScope(r.scopeConverter.ToEntity(api.Scope)).
		SetNillableCreatedBy(api.CreatedBy).
		SetCreatedAt(time.Now())

	if api.Id != nil {
		builder.SetID(api.GetId())
	}

	return builder
}

func (r *ApiRepo) BatchCreate(ctx context.Context, apis []*permissionV1.Api) error {
	if len(apis) == 0 {
		return nil
	}

	bulk := make([]*ent.APICreate, 0, len(apis))
	for _, dto := range apis {
		builder := r.newApiCreate(dto)
		bulk = append(bulk, builder)
	}

	bulkBuilder := r.entClient.Client().Api.CreateBulk(bulk...)

	if err := bulkBuilder.Exec(ctx); err != nil {
		r.log.Errorf("batch insert apis failed: %s", err.Error())
		return permissionV1.ErrorInternalServerError("batch insert apis failed")
	}

	return nil
}

// derefString nil 安全解引用（ent 字段为 nillable 指针）。
func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// SyncApis 以 (path, method) 为业务键对账式同步接口表：
//   - 已存在的端点：module/operation/description 有变化则按 ID 原位更新
//     （不换 ID——角色/权限绑定引用既有 ID，重排会扯断绑定）
//   - 缺失的端点：显式分配 ID（从 max(id)+1 续号）创建，并在创建后把 PG 序列
//     对齐到 max(id)——历史上接口行全部以显式 ID 插入、序列从未启用，这是
//     "手工单建撞主键"历史遗留的根治
//   - 文档中已移除的端点不删除（同理保护绑定），仅在其上标注
//
// 该方法替代原先"清空才可用"的全量重插实现。
func (r *ApiRepo) SyncApis(ctx context.Context, apis []*permissionV1.Api) error {
	if len(apis) == 0 {
		return nil
	}

	rows, err := r.entClient.Client().Api.Query().All(ctx)
	if err != nil {
		r.log.Errorf("sync apis: load existing failed: %s", err.Error())
		return permissionV1.ErrorInternalServerError("load existing apis failed")
	}

	type key struct {
		path   string
		method string
	}
	existing := make(map[key]*ent.Api, len(rows))
	maxID := uint32(0)
	for _, row := range rows {
		existing[key{path: derefString(row.Path), method: derefString(row.Method)}] = row
		if row.ID > maxID {
			maxID = row.ID
		}
	}

	created, updated := 0, 0
	now := time.Now()
	for _, dto := range apis {
		k := key{path: dto.GetPath(), method: dto.GetMethod()}
		if row, ok := existing[k]; ok {
			// 已存在：仅在元数据变化时原位更新
			if derefString(row.Module) == dto.GetModule() &&
				derefString(row.Operation) == dto.GetOperation() &&
				derefString(row.Description) == dto.GetDescription() &&
				derefString(row.ModuleDescription) == dto.GetModuleDescription() {
				continue
			}
			if _, uerr := r.entClient.Client().Api.UpdateOneID(row.ID).
				SetNillableModule(dto.Module).
				SetNillableModuleDescription(dto.ModuleDescription).
				SetNillableOperation(dto.Operation).
				SetNillableDescription(dto.Description).
				SetUpdatedAt(now).
				Save(ctx); uerr != nil {
				r.log.Errorf("sync apis: update api %d failed: %s", row.ID, uerr.Error())
				return permissionV1.ErrorInternalServerError("update api failed")
			}
			updated++
			continue
		}

		// 缺失：显式续号创建（不依赖 PG 序列的历史位置）
		maxID++
		if _, cerr := r.newApiCreate(dto).SetID(maxID).SetCreatedAt(now).Save(ctx); cerr != nil {
			r.log.Errorf("sync apis: create api %s %s failed: %s", dto.GetMethod(), dto.GetPath(), cerr.Error())
			return permissionV1.ErrorInternalServerError("create api failed")
		}
		existing[k] = nil // 防同批次重复
		created++
	}

	// PG 序列对齐：修复"显式 ID 行导致序列落后、无 ID 单建撞主键"的历史遗留
	if created > 0 {
		if sxerr := r.entClient.Exec(ctx,
			"SELECT setval(pg_get_serial_sequence('sys_apis','id'), GREATEST(COALESCE((SELECT MAX(id) FROM sys_apis), 0), 1), true)",
			nil, nil); sxerr != nil {
			// 序列对齐失败不影响同步结果，仅影响后续无 ID 单建
			r.log.Warnf("sync apis: align id sequence failed: %s", sxerr.Error())
		}
	}

	r.log.Infof("sync apis: done (total=%d created=%d updated=%d)", len(apis), created, updated)
	return nil
}

func (r *ApiRepo) Update(ctx context.Context, req *permissionV1.UpdateApiRequest) error {
	if req == nil || req.Data == nil {
		return permissionV1.ErrorBadRequest("invalid parameter")
	}

	// 如果不存在则创建
	if req.GetAllowMissing() {
		exist, err := r.IsExist(ctx, req.GetId())
		if err != nil {
			return err
		}
		if !exist {
			createReq := &permissionV1.CreateApiRequest{Data: req.Data}
			createReq.Data.CreatedBy = createReq.Data.UpdatedBy
			createReq.Data.UpdatedBy = nil
			return r.Create(ctx, createReq)
		}
	}

	builder := r.entClient.Client().Debug().Api.Update()
	// 租户作用域：仅更新本租户 API，避免跨租户改他人数据（按 hasTenant 条件加）
	if tid, hasTenant := maybeTenantFromViewer(ctx); hasTenant {
		builder.Where(api.TenantIDEQ(tid))
	}
	callerUserID, hasUser := viewerUserIDFromContext(ctx)
	err := r.repository.UpdateX(ctx, builder, req.Data, req.GetUpdateMask(),
		func(dto *permissionV1.Api) {
			builder.
				SetNillableDescription(req.Data.Description).
				SetNillableModule(req.Data.Module).
				SetNillableModuleDescription(req.Data.ModuleDescription).
				SetNillableOperation(req.Data.Operation).
				SetNillablePath(req.Data.Path).
				SetNillableMethod(req.Data.Method).
				SetNillableScope(r.scopeConverter.ToEntity(req.Data.Scope)).
				SetUpdatedAt(time.Now())

			// updated_by 强制取调用者身份，保证审计归属真实，忽略客户端值
			if hasUser {
				builder.SetUpdatedBy(callerUserID)
			}
		},
		func(s *sql.Selector) {
			s.Where(sql.EQ(api.FieldID, req.GetId()))
		},
	)

	return err
}

func (r *ApiRepo) Delete(ctx context.Context, req *permissionV1.DeleteApiRequest) error {
	if req == nil {
		return permissionV1.ErrorBadRequest("invalid parameter")
	}

	builder := r.entClient.Client().Debug().Api.Delete()
	// 租户作用域：仅删除本租户 API，避免跨租户删他人数据（按 hasTenant 条件加）
	if tid, hasTenant := maybeTenantFromViewer(ctx); hasTenant {
		builder.Where(api.TenantIDEQ(tid))
	}

	_, err := r.repository.Delete(ctx, builder, func(s *sql.Selector) {
		s.Where(sql.EQ(api.FieldID, req.GetId()))
	})
	if err != nil {
		r.log.Errorf("delete api failed: %s", err.Error())
		return permissionV1.ErrorInternalServerError("delete api failed")
	}

	return nil
}

// Truncate 清空表数据
func (r *ApiRepo) Truncate(ctx context.Context) error {
	if _, err := r.entClient.Client().Api.Delete().Exec(ctx); err != nil {
		r.log.Errorf("failed to truncate apis table: %s", err.Error())
		return permissionV1.ErrorInternalServerError("truncate failed")
	}
	return nil
}
