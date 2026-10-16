package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	entCrud "github.com/tx7do/go-crud/entgo"
	"github.com/tx7do/go-utils/copierutil"
	"github.com/tx7do/go-utils/mapper"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	flowV1 "go-wind-cms/api/gen/go/flow/service/v1"

	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/flowrule"
	"go-wind-cms/app/core/service/internal/data/ent/predicate"
)

// FlowRepo 自动化 Flow 规则仓储。
type FlowRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper

	mapper *mapper.CopierMapper[flowV1.FlowRule, ent.FlowRule]
}

func NewFlowRepo(ctx *bootstrap.Context, entClient *entCrud.EntClient[*ent.Client]) *FlowRepo {
	r := &FlowRepo{
		entClient: entClient,
		log:       ctx.NewLoggerHelper("flow/repo/core-service"),
		mapper:    mapper.NewCopierMapper[flowV1.FlowRule, ent.FlowRule](),
	}
	r.mapper.AppendConverters(copierutil.NewTimeStringConverterPair())
	r.mapper.AppendConverters(copierutil.NewTimeTimestamppbConverterPair())
	return r
}

func (r *FlowRepo) List(ctx context.Context, req *paginationV1.PagingRequest) (*flowV1.ListFlowRuleResponse, error) {
	// 列表按租户过滤（管理面）；平台管理员 tenant 0 不限
	builder := r.entClient.Client().FlowRule.Query()
	if tid, hasTenant := maybeTenantFromViewer(ctx); hasTenant {
		builder.Where(flowrule.TenantIDEQ(tid))
	}
	total, err := builder.Clone().Count(ctx)
	if err != nil {
		r.log.Errorf("count flow rules failed: %s", err.Error())
		return nil, errors.InternalServer("FLOW_RULE", "list flow rules failed")
	}

	rows, err := builder.
		Order(ent.Desc(flowrule.FieldID)).
		All(ctx)
	if err != nil {
		r.log.Errorf("list flow rules failed: %s", err.Error())
		return nil, errors.InternalServer("FLOW_RULE", "list flow rules failed")
	}
	items := make([]*flowV1.FlowRule, 0, len(rows))
	for _, row := range rows {
		items = append(items, r.mapper.ToDTO(row))
	}
	return &flowV1.ListFlowRuleResponse{Items: items, Total: uint64(total)}, nil
}

func (r *FlowRepo) Get(ctx context.Context, id uint32) (*flowV1.FlowRule, error) {
	row, err := r.entClient.Client().FlowRule.Query().
		Where(flowrule.IDEQ(id)).
		First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.NotFound("FLOW_RULE", "flow rule not found")
		}
		r.log.Errorf("get flow rule failed: %s", err.Error())
		return nil, errors.InternalServer("FLOW_RULE", "get flow rule failed")
	}
	return r.mapper.ToDTO(row), nil
}

func (r *FlowRepo) Create(ctx context.Context, dto *flowV1.FlowRule) (*flowV1.FlowRule, error) {
	builder := r.entClient.Client().FlowRule.Create().
		SetName(dto.GetName()).
		SetTriggerEvent(dto.GetTriggerEvent()).
		SetActions(dto.GetActions())
	if dto.ConditionStatus != nil {
		builder.SetConditionStatus(dto.GetConditionStatus())
	}
	if dto.TriggerCron != nil {
		builder.SetTriggerCron(dto.GetTriggerCron())
	}
	if dto.Enabled != nil {
		builder.SetEnabled(dto.GetEnabled())
	}
	row, err := builder.Save(ctx)
	if err != nil {
		r.log.Errorf("create flow rule failed: %s", err.Error())
		return nil, errors.InternalServer("FLOW_RULE", "create flow rule failed")
	}
	return r.mapper.ToDTO(row), nil
}

func (r *FlowRepo) Update(ctx context.Context, id uint32, dto *flowV1.FlowRule) (*flowV1.FlowRule, error) {
	row, err := r.entClient.Client().FlowRule.Query().
		Where(flowrule.IDEQ(id)).
		First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.NotFound("FLOW_RULE", "flow rule not found")
		}
		return nil, errors.InternalServer("FLOW_RULE", "update flow rule failed")
	}

	builder := r.entClient.Client().FlowRule.UpdateOneID(row.ID)
	if dto.Name != nil {
		builder.SetName(dto.GetName())
	}
	if dto.TriggerEvent != nil {
		builder.SetTriggerEvent(dto.GetTriggerEvent())
	}
	if dto.ConditionStatus != nil {
		builder.SetConditionStatus(dto.GetConditionStatus())
	}
	if dto.TriggerCron != nil {
		builder.SetTriggerCron(dto.GetTriggerCron())
	}
	if dto.Actions != nil {
		builder.SetActions(dto.GetActions())
	}
	if dto.Enabled != nil {
		builder.SetEnabled(dto.GetEnabled())
	}
	updated, err := builder.Save(ctx)
	if err != nil {
		r.log.Errorf("update flow rule failed: %s", err.Error())
		return nil, errors.InternalServer("FLOW_RULE", "update flow rule failed")
	}
	return r.mapper.ToDTO(updated), nil
}

func (r *FlowRepo) Delete(ctx context.Context, id uint32) error {
	if err := r.entClient.Client().FlowRule.DeleteOneID(id).Exec(ctx); err != nil {
		if ent.IsNotFound(err) {
			return errors.NotFound("FLOW_RULE", "flow rule not found")
		}
		r.log.Errorf("delete flow rule failed: %s", err.Error())
		return errors.InternalServer("FLOW_RULE", "delete flow rule failed")
	}
	return nil
}

// ListEnabledForEvent 列出某事件下启用的规则（worker 用，SystemViewer 全租户）。
func (r *FlowRepo) ListEnabledForEvent(ctx context.Context, tenantID uint32, event string) ([]*flowV1.FlowRule, error) {
	// 定时规则（trigger_cron 非空）只走 cron 调度，不参与事件扇出，避免双重触发
	conds := []predicate.FlowRule{
		flowrule.EnabledEQ(true),
		flowrule.TriggerEventEQ(event),
		flowrule.TriggerCronIsNil(),
	}
	if tenantID > 0 {
		conds = append(conds, flowrule.TenantIDEQ(tenantID))
	}
	rows, err := r.entClient.Client().FlowRule.Query().Where(conds...).All(ctx)
	if err != nil {
		r.log.Errorf("list flow rules for event failed: %s", err.Error())
		return nil, err
	}
	dtos := make([]*flowV1.FlowRule, 0, len(rows))
	for _, row := range rows {
		dtos = append(dtos, r.mapper.ToDTO(row))
	}
	return dtos, nil
}

// FlowCronRule 定时规则最小投影（调度 tick 用；FlowRule DTO 不含租户，故单列）。
type FlowCronRule struct {
	ID       uint32
	TenantID uint32
	Cron     string
	Event    string // trigger_event 原值（观测/占位符参考，触发事件固定 "schedule"）
}

// ListEnabledCronRules 列出全部启用中的定时规则（跨租户；SystemViewer）。
func (r *FlowRepo) ListEnabledCronRules(ctx context.Context) ([]FlowCronRule, error) {
	rows, err := r.entClient.Client().FlowRule.Query().
		Where(
			flowrule.EnabledEQ(true),
			flowrule.TriggerCronNotNil(),
		).
		All(ctx)
	if err != nil {
		r.log.Errorf("list cron flow rules failed: %s", err.Error())
		return nil, err
	}
	out := make([]FlowCronRule, 0, len(rows))
	for _, row := range rows {
		if row.TenantID == nil || row.TriggerCron == nil {
			continue
		}
		out = append(out, FlowCronRule{
			ID:       row.ID,
			TenantID: *row.TenantID,
			Cron:     *row.TriggerCron,
			Event:    row.TriggerEvent,
		})
	}
	return out, nil
}

// IncrMatchCount 命中计数 +1（观测用，调用方 best-effort 吞错）。
func (r *FlowRepo) IncrMatchCount(ctx context.Context, id uint32) error {
	if _, err := r.entClient.Client().FlowRule.UpdateOneID(id).
		AddMatchCount(1).
		Save(ctx); err != nil {
		r.log.Errorf("incr match count for flow rule [%d] failed: %s", id, err.Error())
		return err
	}
	return nil
}
