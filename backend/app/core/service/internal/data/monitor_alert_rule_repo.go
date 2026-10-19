package data

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	entCrud "github.com/tx7do/go-crud/entgo"
	"github.com/tx7do/go-utils/copierutil"
	"github.com/tx7do/go-utils/mapper"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/monitoralertrule"

	monitorAlertV1 "go-wind-cms/api/gen/go/monitor_alert/service/v1"
	notifyV1 "go-wind-cms/api/gen/go/notify/service/v1"
)

// MonitorAlertRuleRepo 监控告警规则仓储：指标阈值 → 触发通知，平台全局。
type MonitorAlertRuleRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper

	mapper           *mapper.CopierMapper[monitorAlertV1.MonitorAlertRule, ent.MonitorAlertRule]
	metricConverter  *mapper.EnumTypeConverter[monitorAlertV1.MonitorMetric, monitoralertrule.Metric]
	opConverter      *mapper.EnumTypeConverter[monitorAlertV1.AlertOp, monitoralertrule.Op]
	channelConverter *mapper.EnumTypeConverter[notifyV1.Channel, monitoralertrule.Channel]
}

func NewMonitorAlertRuleRepo(ctx *bootstrap.Context, entClient *entCrud.EntClient[*ent.Client]) *MonitorAlertRuleRepo {
	repo := &MonitorAlertRuleRepo{
		log:              ctx.NewLoggerHelper("monitor-alert/repo/core-service"),
		entClient:        entClient,
		mapper:           mapper.NewCopierMapper[monitorAlertV1.MonitorAlertRule, ent.MonitorAlertRule](),
		metricConverter:  mapper.NewEnumTypeConverter[monitorAlertV1.MonitorMetric, monitoralertrule.Metric](monitorAlertV1.MonitorMetric_name, monitorAlertV1.MonitorMetric_value),
		opConverter:      mapper.NewEnumTypeConverter[monitorAlertV1.AlertOp, monitoralertrule.Op](monitorAlertV1.AlertOp_name, monitorAlertV1.AlertOp_value),
		channelConverter: mapper.NewEnumTypeConverter[notifyV1.Channel, monitoralertrule.Channel](notifyV1.Channel_name, notifyV1.Channel_value),
	}

	repo.mapper.AppendConverters(copierutil.NewTimeStringConverterPair())
	repo.mapper.AppendConverters(copierutil.NewTimeTimestamppbConverterPair())
	repo.mapper.AppendConverters(repo.metricConverter.NewConverterPair())
	repo.mapper.AppendConverters(repo.opConverter.NewConverterPair())
	repo.mapper.AppendConverters(repo.channelConverter.NewConverterPair())

	return repo
}

func (r *MonitorAlertRuleRepo) List(ctx context.Context, req *paginationV1.PagingRequest) (*monitorAlertV1.ListMonitorAlertRuleResponse, error) {
	if req == nil {
		return nil, errors.BadRequest("MONITOR_ALERT", "invalid parameter")
	}

	builder := r.entClient.Client().MonitorAlertRule.Query()

	total, err := builder.Clone().Count(ctx)
	if err != nil {
		r.log.Errorf("list monitor alert rules failed: %s", err.Error())
		return nil, errors.InternalServer("MONITOR_ALERT", "list monitor alert rules failed")
	}

	entities, err := builder.
		Order(ent.Desc(monitoralertrule.FieldID)).
		All(ctx)
	if err != nil {
		r.log.Errorf("list monitor alert rules failed: %s", err.Error())
		return nil, errors.InternalServer("MONITOR_ALERT", "list monitor alert rules failed")
	}

	dtos := make([]*monitorAlertV1.MonitorAlertRule, 0, len(entities))
	for _, entity := range entities {
		dtos = append(dtos, r.mapper.ToDTO(entity))
	}
	return &monitorAlertV1.ListMonitorAlertRuleResponse{
		Total: uint64(total),
		Items: dtos,
	}, nil
}

// ListEnabled 取全部启用规则（评估器每轮用，量小不分页）。
func (r *MonitorAlertRuleRepo) ListEnabled(ctx context.Context) ([]*monitorAlertV1.MonitorAlertRule, error) {
	entities, err := r.entClient.Client().MonitorAlertRule.Query().
		Where(monitoralertrule.IsEnabledEQ(true)).
		All(ctx)
	if err != nil {
		r.log.Errorf("list enabled monitor alert rules failed: %s", err.Error())
		return nil, errors.InternalServer("MONITOR_ALERT", "list monitor alert rules failed")
	}

	dtos := make([]*monitorAlertV1.MonitorAlertRule, 0, len(entities))
	for _, entity := range entities {
		dtos = append(dtos, r.mapper.ToDTO(entity))
	}
	return dtos, nil
}

func (r *MonitorAlertRuleRepo) Get(ctx context.Context, id uint32) (*monitorAlertV1.MonitorAlertRule, error) {
	if id == 0 {
		return nil, errors.BadRequest("MONITOR_ALERT", "id is required")
	}

	entity, err := r.entClient.Client().MonitorAlertRule.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.NotFound("MONITOR_ALERT", "monitor alert rule not found")
		}
		r.log.Errorf("get monitor alert rule [%d] failed: %s", id, err.Error())
		return nil, errors.InternalServer("MONITOR_ALERT", "get monitor alert rule failed")
	}

	return r.mapper.ToDTO(entity), nil
}

// Create 建一条规则。
func (r *MonitorAlertRuleRepo) Create(ctx context.Context, req *monitorAlertV1.MonitorAlertRule, operatorID uint32) (uint32, error) {
	if req == nil {
		return 0, errors.BadRequest("MONITOR_ALERT", "invalid request")
	}

	builder := r.entClient.Client().MonitorAlertRule.Create().
		SetNillableName(req.Name).
		SetNillableMetric(r.metricConverter.ToEntity(req.Metric)).
		SetNillableOp(r.opConverter.ToEntity(req.Op)).
		SetNillableThreshold(req.Threshold).
		SetNillableCooldownMinutes(req.CooldownMinutes).
		SetNillableChannel(r.channelConverter.ToEntity(req.Channel)).
		SetNillableTarget(req.Target).
		SetNillableIsEnabled(req.IsEnabled).
		SetNillableRemark(req.Remark).
		SetCreatedAt(time.Now())

	if operatorID != 0 {
		builder.SetCreatedBy(operatorID)
	}

	created, err := builder.Save(ctx)
	if err != nil {
		r.log.Errorf("insert monitor alert rule failed: %s", err.Error())
		return 0, errors.InternalServer("MONITOR_ALERT", "insert monitor alert rule failed")
	}

	return created.ID, nil
}

// Update 更新规则；评估器状态列（last_firing/last_value/last_alerted_at）不在
// 可编辑集合里——它们只走 MarkFiring，掩码更新永不触碰。
func (r *MonitorAlertRuleRepo) Update(ctx context.Context, id uint32, req *monitorAlertV1.MonitorAlertRule, operatorID uint32) error {
	if req == nil {
		return errors.BadRequest("MONITOR_ALERT", "invalid request")
	}
	if id == 0 {
		return errors.BadRequest("MONITOR_ALERT", "id is required")
	}

	builder := r.entClient.Client().MonitorAlertRule.UpdateOneID(id)
	if req.Name != nil {
		builder.SetNillableName(req.Name)
	}
	if req.Metric != nil {
		builder.SetNillableMetric(r.metricConverter.ToEntity(req.Metric))
	}
	if req.Op != nil {
		builder.SetNillableOp(r.opConverter.ToEntity(req.Op))
	}
	if req.Threshold != nil {
		builder.SetNillableThreshold(req.Threshold)
	}
	if req.CooldownMinutes != nil {
		builder.SetNillableCooldownMinutes(req.CooldownMinutes)
	}
	if req.Channel != nil {
		builder.SetNillableChannel(r.channelConverter.ToEntity(req.Channel))
	}
	if req.Target != nil {
		builder.SetNillableTarget(req.Target)
	}
	if req.IsEnabled != nil {
		builder.SetNillableIsEnabled(req.IsEnabled)
	}
	if req.Remark != nil {
		builder.SetNillableRemark(req.Remark)
	}
	if operatorID != 0 {
		builder.SetUpdatedBy(operatorID)
	}
	builder.SetUpdatedAt(time.Now())

	if _, err := builder.Save(ctx); err != nil {
		r.log.Errorf("update monitor alert rule [%d] failed: %s", id, err.Error())
		return errors.InternalServer("MONITOR_ALERT", "update monitor alert rule failed")
	}
	return nil
}

// MarkFiring 评估器回写规则状态：firing 状态、本次指标值、告警时间。
// 独立成一个写路径（不走 Update 的掩码语义）：评估回写是系统行为，
// 与管理员的编辑互不覆盖——掩码更新永不触碰这三列，本方法只写这三列。
func (r *MonitorAlertRuleRepo) MarkFiring(ctx context.Context, id uint32, firing bool, value *float64, alertedAt *time.Time) error {
	builder := r.entClient.Client().MonitorAlertRule.UpdateOneID(id).
		SetLastFiring(firing).
		SetNillableLastValue(value)
	builder.SetNillableLastAlertedAt(alertedAt)
	if _, err := builder.Save(ctx); err != nil {
		r.log.Errorf("mark monitor alert rule [%d] firing=%v failed: %s", id, firing, err.Error())
		return errors.InternalServer("MONITOR_ALERT", "update monitor alert rule failed")
	}
	return nil
}

func (r *MonitorAlertRuleRepo) Delete(ctx context.Context, id uint32) error {
	if id == 0 {
		return errors.BadRequest("MONITOR_ALERT", "id is required")
	}

	if err := r.entClient.Client().MonitorAlertRule.DeleteOneID(id).Exec(ctx); err != nil {
		r.log.Errorf("delete monitor alert rule [%d] failed: %s", id, err.Error())
		return errors.InternalServer("MONITOR_ALERT", "delete monitor alert rule failed")
	}

	return nil
}
