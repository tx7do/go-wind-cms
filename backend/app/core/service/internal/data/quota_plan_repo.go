package data

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	entCrud "github.com/tx7do/go-crud/entgo"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/timestamppb"

	identityV1 "go-wind-cms/api/gen/go/identity/service/v1"
	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/quotaplan"
)

// QuotaPlanRepo 配额套餐模板仓储（sys_quota_plans；平台级命名模板）。
type QuotaPlanRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper
}

func NewQuotaPlanRepo(ctx *bootstrap.Context, entClient *entCrud.EntClient[*ent.Client]) *QuotaPlanRepo {
	return &QuotaPlanRepo{
		log:       ctx.NewLoggerHelper("quota-plan/repo/core-service"),
		entClient: entClient,
	}
}

// quotaPlanItemDTO items JSON 数组元素（与 proto QuotaPlanItem 同构，quotaType 为枚举名字符串）。
type quotaPlanItemDTO struct {
	QuotaType string `json:"quotaType"`
	Limit     int64  `json:"limit"`
}

// toProto ent 行 → proto QuotaPlan（items JSON 反序列化）。
func (r *QuotaPlanRepo) toProto(e *ent.QuotaPlan) *identityV1.QuotaPlan {
	dto := &identityV1.QuotaPlan{
		Id:          trans.Ptr(uint32(e.ID)),
		Name:        e.Name,
		Description: e.Description,
	}
	if e.CreatedAt != nil {
		dto.CreatedAt = timestamppb.New(*e.CreatedAt)
	}
	if e.UpdatedAt != nil {
		dto.UpdatedAt = timestamppb.New(*e.UpdatedAt)
	}
	if e.Items != "" {
		var items []quotaPlanItemDTO
		if err := json.Unmarshal([]byte(e.Items), &items); err == nil {
			for _, it := range items {
				enumVal, ok := identityV1.TenantQuota_QuotaType_value[it.QuotaType]
				if !ok {
					continue // 未知类型跳过（旧数据容忍）
				}
				dto.Items = append(dto.GetItems(), &identityV1.QuotaPlanItem{
					QuotaType: identityV1.TenantQuota_QuotaType(enumVal),
					Limit:     it.Limit,
				})
			}
		} else {
			r.log.Errorf("parse quota plan items failed: %s", err.Error())
		}
	}
	return dto
}

// itemsToJSON proto items → JSON 串（校验后调用）。
func QuotaPlanItemsToJSON(items []*identityV1.QuotaPlanItem) (string, error) {
	dtos := make([]quotaPlanItemDTO, 0, len(items))
	seen := map[string]bool{}
	for _, it := range items {
		name := it.GetQuotaType().String()
		if name == "" || seen[name] {
			continue
		}
		seen[name] = true
		dtos = append(dtos, quotaPlanItemDTO{QuotaType: name, Limit: it.GetLimit()})
	}
	raw, err := json.Marshal(dtos)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

// List 全量模板（平台级数量少，不分页）。
func (r *QuotaPlanRepo) List(ctx context.Context) ([]*identityV1.QuotaPlan, error) {
	rows, err := r.entClient.Client().QuotaPlan.Query().
		Order(ent.Desc(quotaplan.FieldID)).
		All(ctx)
	if err != nil {
		r.log.Errorf("list quota plans failed: %s", err.Error())
		return nil, identityV1.ErrorInternalServerError("list quota plans failed")
	}
	out := make([]*identityV1.QuotaPlan, 0, len(rows))
	for i := range rows {
		out = append(out, r.toProto(rows[i]))
	}
	return out, nil
}

// Get 按 ID 取模板。
func (r *QuotaPlanRepo) Get(ctx context.Context, id uint32) (*identityV1.QuotaPlan, error) {
	e, err := r.entClient.Client().QuotaPlan.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, identityV1.ErrorNotFound("quota plan not found")
		}
		r.log.Errorf("get quota plan failed: %s", err.Error())
		return nil, identityV1.ErrorInternalServerError("get quota plan failed")
	}
	return r.toProto(e), nil
}

// Create 创建模板（重名报错由唯一索引兜底）。
func (r *QuotaPlanRepo) Create(ctx context.Context, dto *identityV1.QuotaPlan, operatorID uint32) (uint32, error) {
	itemsJSON, err := QuotaPlanItemsToJSON(dto.GetItems())
	if err != nil {
		return 0, identityV1.ErrorBadRequest("invalid quota plan items")
	}
	now := time.Now()
	builder := r.entClient.Client().QuotaPlan.Create().
		SetName(dto.GetName()).
		SetItems(itemsJSON).
		SetNillableDescription(dto.Description).
		SetNillableCreatedAt(&now).
		SetNillableUpdatedAt(&now)
	_ = operatorID
	row, err := builder.Save(ctx)
	if err != nil {
		r.log.Errorf("create quota plan failed: %s", err.Error())
		return 0, identityV1.ErrorInternalServerError("create quota plan failed")
	}
	return row.ID, nil
}

// Delete 删除模板（已套用的租户配额是落库快照，不受影响）。
func (r *QuotaPlanRepo) Delete(ctx context.Context, id uint32) error {
	err := r.entClient.Client().QuotaPlan.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return identityV1.ErrorNotFound("quota plan not found")
		}
		r.log.Errorf("delete quota plan failed: %s", err.Error())
		return identityV1.ErrorInternalServerError("delete quota plan failed")
	}
	return nil
}
