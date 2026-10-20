package data

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	entCrud "github.com/tx7do/go-crud/entgo"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	appViewer "go-wind-cms/pkg/entgo/viewer"

	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/aiusagelog"
	"go-wind-cms/app/core/service/internal/data/ent/apiauditlog"
	"go-wind-cms/app/core/service/internal/data/ent/file"
	"go-wind-cms/app/core/service/internal/data/ent/tenant"
	"go-wind-cms/app/core/service/internal/data/ent/tenantquota"
	"go-wind-cms/app/core/service/internal/data/ent/user"

	identityV1 "go-wind-cms/api/gen/go/identity/service/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TenantQuotaRepo 租户配额仓储：配额 CRUD + 用量计量聚合 + 闸门查询。
//
// 计量口径与 GetUsage/水位扫描三方同源（同一张表、同一个查询）：
//   - USER_LIMIT: sys_users 按租户 COUNT
//   - STORAGE:    files 按租户 SUM(size)
//   - API_CALL:   sys_api_audit_logs 按租户 COUNT（活表现存行数即口径，审计归档会回落）
//   - AI_TOKENS:  sys_ai_usage_logs 本月 SUM(total_tokens)
type TenantQuotaRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper
}

func NewTenantQuotaRepo(ctx *bootstrap.Context, entClient *entCrud.EntClient[*ent.Client]) *TenantQuotaRepo {
	return &TenantQuotaRepo{
		entClient: entClient,
		log:       ctx.NewLoggerHelper("tenant-quota/repo/core-service"),
	}
}

// ─── CRUD ───

func (r *TenantQuotaRepo) List(ctx context.Context, req *paginationV1.PagingRequest) (*identityV1.ListTenantQuotaResponse, error) {
	if req == nil {
		return nil, errors.BadRequest("TENANT_QUOTA", "invalid parameter")
	}

	builder := r.entClient.Client().TenantQuota.Query()

	total, err := builder.Clone().Count(ctx)
	if err != nil {
		r.log.Errorf("list tenant quotas failed: %s", err.Error())
		return nil, errors.InternalServer("TENANT_QUOTA", "list tenant quotas failed")
	}

	if req.GetPage() > 0 || req.GetPageSize() > 0 {
		pageSize := int(req.GetPageSize())
		if pageSize <= 0 {
			pageSize = 20
		}
		offset := int(req.GetPage()-1) * pageSize
		if offset < 0 {
			offset = 0
		}
		builder = builder.Limit(pageSize).Offset(offset)
	}

	entities, err := builder.
		Order(ent.Desc(tenantquota.FieldID)).
		All(ctx)
	if err != nil {
		r.log.Errorf("list tenant quotas failed: %s", err.Error())
		return nil, errors.InternalServer("TENANT_QUOTA", "list tenant quotas failed")
	}

	dtos := make([]*identityV1.TenantQuota, 0, len(entities))
	for _, e := range entities {
		dtos = append(dtos, quotaEntityToProto(e))
	}
	return &identityV1.ListTenantQuotaResponse{Total: uint64(total), Items: dtos}, nil
}

func (r *TenantQuotaRepo) Get(ctx context.Context, id uint32) (*identityV1.TenantQuota, error) {
	if id == 0 {
		return nil, errors.BadRequest("TENANT_QUOTA", "id is required")
	}

	entity, err := r.entClient.Client().TenantQuota.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.NotFound("TENANT_QUOTA", "tenant quota not found")
		}
		r.log.Errorf("get tenant quota [%d] failed: %s", id, err.Error())
		return nil, errors.InternalServer("TENANT_QUOTA", "get tenant quota failed")
	}
	return quotaEntityToProto(entity), nil
}

// Create 新建配额行；同一租户同一类型唯一（库层唯一索引兜底，先查给出可读错误）。
func (r *TenantQuotaRepo) Create(ctx context.Context, data *identityV1.TenantQuota, operatorID uint32) (uint32, error) {
	if data == nil || data.GetTenantId() == 0 || data.GetQuotaType() == identityV1.TenantQuota_QUOTA_TYPE_UNSPECIFIED {
		return 0, errors.BadRequest("TENANT_QUOTA", "tenant id and quota type are required")
	}

	exist, err := r.entClient.Client().TenantQuota.Query().
		Where(
			tenantquota.TenantIDEQ(data.GetTenantId()),
			tenantquota.QuotaTypeEQ(quotaTypeProtoToEnt(data.GetQuotaType())),
		).
		Exist(ctx)
	if err != nil {
		r.log.Errorf("check tenant quota exist failed: %s", err.Error())
		return 0, errors.InternalServer("TENANT_QUOTA", "create tenant quota failed")
	}
	if exist {
		return 0, errors.BadRequest("TENANT_QUOTA", "quota type already exists for this tenant")
	}

	builder := r.entClient.Client().TenantQuota.Create().
		SetTenantID(data.GetTenantId()).
		SetQuotaType(quotaTypeProtoToEnt(data.GetQuotaType())).
		SetQuotaValue(data.GetQuotaValue()).
		SetCreatedAt(time.Now())
	if operatorID != 0 {
		builder.SetCreatedBy(operatorID)
	}

	created, err := builder.Save(ctx)
	if err != nil {
		r.log.Errorf("insert tenant quota failed: %s", err.Error())
		return 0, errors.InternalServer("TENANT_QUOTA", "insert tenant quota failed")
	}
	return created.ID, nil
}

// Update 更新配额值；tenant_id/quota_type 创建后不可改（改配额口径=换行重建）。
func (r *TenantQuotaRepo) Update(ctx context.Context, id uint32, data *identityV1.TenantQuota, operatorID uint32) error {
	if data == nil {
		return errors.BadRequest("TENANT_QUOTA", "invalid parameter")
	}
	if id == 0 {
		return errors.BadRequest("TENANT_QUOTA", "id is required")
	}

	builder := r.entClient.Client().TenantQuota.UpdateOneID(id).
		SetQuotaValue(data.GetQuotaValue()).
		SetUpdatedAt(time.Now())
	if operatorID != 0 {
		builder.SetUpdatedBy(operatorID)
	}

	if err := builder.Exec(ctx); err != nil {
		if ent.IsNotFound(err) {
			return errors.NotFound("TENANT_QUOTA", "tenant quota not found")
		}
		r.log.Errorf("update tenant quota [%d] failed: %s", id, err.Error())
		return errors.InternalServer("TENANT_QUOTA", "update tenant quota failed")
	}
	return nil
}

func (r *TenantQuotaRepo) Delete(ctx context.Context, id uint32) error {
	if id == 0 {
		return errors.BadRequest("TENANT_QUOTA", "id is required")
	}
	if err := r.entClient.Client().TenantQuota.DeleteOneID(id).Exec(ctx); err != nil {
		if ent.IsNotFound(err) {
			return errors.NotFound("TENANT_QUOTA", "tenant quota not found")
		}
		r.log.Errorf("delete tenant quota [%d] failed: %s", id, err.Error())
		return errors.InternalServer("TENANT_QUOTA", "delete tenant quota failed")
	}
	return nil
}

// ListByTenant 取租户全部配额行（GetUsage / 闸门 / 水位扫描共用）。
func (r *TenantQuotaRepo) ListByTenant(ctx context.Context, tenantID uint32) ([]*ent.TenantQuota, error) {
	return r.entClient.Client().TenantQuota.Query().
		Where(tenantquota.TenantIDEQ(tenantID)).
		All(ctx)
}

// ListAll 取全部配额行（水位扫描按租户归集用；平台级配置表行数有限）。
func (r *TenantQuotaRepo) ListAll(ctx context.Context) ([]*ent.TenantQuota, error) {
	return r.entClient.Client().TenantQuota.Query().
		Order(ent.Desc(tenantquota.FieldID)).
		All(ctx)
}

// ─── 计量聚合（与 GetUsage / 闸门 / 水位扫描同源） ───

func (r *TenantQuotaRepo) CountUsersInTenant(ctx context.Context, tenantID uint32) (uint64, error) {
	count, err := r.entClient.Client().User.Query().
		Where(user.TenantIDEQ(tenantID)).
		Count(ctx)
	return uint64(count), err
}

func (r *TenantQuotaRepo) SumStorageInTenant(ctx context.Context, tenantID uint32) (uint64, error) {
	var rows []struct {
		Total uint64 `json:"total"`
	}
	err := r.entClient.Client().File.Query().
		Where(file.TenantIDEQ(tenantID)).
		Aggregate(ent.As(ent.Sum(file.FieldSize), "total")).
		Scan(ctx, &rows)
	if err != nil {
		return 0, err
	}
	if len(rows) > 0 {
		return rows[0].Total, nil
	}
	return 0, nil
}

func (r *TenantQuotaRepo) CountApiCallsInTenant(ctx context.Context, tenantID uint32) (uint64, error) {
	cnt, err := r.entClient.Client().ApiAuditLog.Query().
		Where(apiauditlog.TenantIDEQ(tenantID)).
		Count(ctx)
	return uint64(cnt), err
}

func (r *TenantQuotaRepo) SumAiTokensInTenantSince(ctx context.Context, tenantID uint32, since time.Time) (uint64, error) {
	var rows []struct {
		Total uint64 `json:"total"`
	}
	err := r.entClient.Client().AiUsageLog.Query().
		Where(
			aiusagelog.TenantIDEQ(tenantID),
			aiusagelog.CreatedAtGTE(since),
		).
		Aggregate(ent.As(ent.Sum(aiusagelog.FieldTotalTokens), "total")).
		Scan(ctx, &rows)
	if err != nil {
		return 0, err
	}
	if len(rows) > 0 {
		return rows[0].Total, nil
	}
	return 0, nil
}

// GetTenantStatus 查租户状态与基础信息（闸门第 1 步；未命中返回 nil）。
func (r *TenantQuotaRepo) GetTenantStatus(ctx context.Context, tenantID uint32) (*ent.Tenant, error) {
	svCtx := appViewer.NewSystemViewerContext(ctx)
	t, err := r.entClient.Client().Tenant.Query().
		Where(tenant.IDEQ(tenantID)).
		Only(svCtx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return t, nil
}

// GetUserByID 查用户行（水位告警读管理员偏好语言用；调用方需传系统旁路上下文）。
func (r *TenantQuotaRepo) GetUserByID(ctx context.Context, id uint32) (*ent.User, error) {
	u, err := r.entClient.Client().User.Query().
		Where(user.IDEQ(id)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return u, nil
}

// ─── 映射 ───

func quotaEntityToProto(e *ent.TenantQuota) *identityV1.TenantQuota {
	if e == nil {
		return nil
	}
	dto := &identityV1.TenantQuota{
		Id:         &e.ID,
		TenantId:   &e.TenantID,
		QuotaValue: &e.QuotaValue,
	}
	switch e.QuotaType {
	case tenantquota.QuotaTypeUserLimit:
		dto.QuotaType = identityV1.TenantQuota_USER_LIMIT.Enum()
	case tenantquota.QuotaTypeStorage:
		dto.QuotaType = identityV1.TenantQuota_STORAGE.Enum()
	case tenantquota.QuotaTypeApiCall:
		dto.QuotaType = identityV1.TenantQuota_API_CALL.Enum()
	case tenantquota.QuotaTypeAiTokens:
		dto.QuotaType = identityV1.TenantQuota_AI_TOKENS.Enum()
	default:
		dto.QuotaType = identityV1.TenantQuota_QUOTA_TYPE_UNSPECIFIED.Enum()
	}
	if e.CreatedAt != nil {
		dto.CreatedAt = timestamppb.New(*e.CreatedAt)
	}
	if e.UpdatedAt != nil {
		dto.UpdatedAt = timestamppb.New(*e.UpdatedAt)
	}
	return dto
}

func quotaTypeProtoToEnt(qt identityV1.TenantQuota_QuotaType) tenantquota.QuotaType {
	switch qt {
	case identityV1.TenantQuota_USER_LIMIT:
		return tenantquota.QuotaTypeUserLimit
	case identityV1.TenantQuota_STORAGE:
		return tenantquota.QuotaTypeStorage
	case identityV1.TenantQuota_API_CALL:
		return tenantquota.QuotaTypeApiCall
	case identityV1.TenantQuota_AI_TOKENS:
		return tenantquota.QuotaTypeAiTokens
	default:
		return tenantquota.QuotaTypeQuotaTypeUnspecified
	}
}
