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
	"go-wind-cms/app/core/service/internal/data/ent/notificationrule"
	"go-wind-cms/app/core/service/internal/data/ent/predicate"

	notifyV1 "go-wind-cms/api/gen/go/notify/service/v1"
)

// NotificationRuleRepo 通知路由规则仓储：事件类型 → 渠道 + 是否异步。
//
// 这张表是 SendDirect 每次派发都要读一次的运行时真相，
// 所以 GetByEventType 刻意不带缓存：路由改了立刻生效，比省一次索引查询重要
//（单表单行、走 event_type 唯一索引）。
type NotificationRuleRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper

	mapper             *mapper.CopierMapper[notifyV1.NotificationRule, ent.NotificationRule]
	channelConverter   *mapper.EnumTypeConverter[notifyV1.Channel, notificationrule.Channel]
	eventTypeConverter *mapper.EnumTypeConverter[notifyV1.EventType, notificationrule.EventType]

	repository *entCrud.Repository[
		ent.NotificationRuleQuery, ent.NotificationRuleSelect,
		ent.NotificationRuleCreate, ent.NotificationRuleCreateBulk,
		ent.NotificationRuleUpdate, ent.NotificationRuleUpdateOne,
		ent.NotificationRuleDelete,
		predicate.NotificationRule,
		notifyV1.NotificationRule, ent.NotificationRule,
	]
}

func NewNotificationRuleRepo(ctx *bootstrap.Context, entClient *entCrud.EntClient[*ent.Client]) *NotificationRuleRepo {
	repo := &NotificationRuleRepo{
		log:       ctx.NewLoggerHelper("notification-rule/repo/core-service"),
		entClient: entClient,
		mapper:    mapper.NewCopierMapper[notifyV1.NotificationRule, ent.NotificationRule](),
		channelConverter: mapper.NewEnumTypeConverter[notifyV1.Channel, notificationrule.Channel](
			notifyV1.Channel_name, notifyV1.Channel_value,
		),
		eventTypeConverter: mapper.NewEnumTypeConverter[notifyV1.EventType, notificationrule.EventType](
			notifyV1.EventType_name, notifyV1.EventType_value,
		),
	}

	repo.repository = entCrud.NewRepository[
		ent.NotificationRuleQuery, ent.NotificationRuleSelect,
		ent.NotificationRuleCreate, ent.NotificationRuleCreateBulk,
		ent.NotificationRuleUpdate, ent.NotificationRuleUpdateOne,
		ent.NotificationRuleDelete,
		predicate.NotificationRule,
		notifyV1.NotificationRule, ent.NotificationRule,
	](repo.mapper)

	repo.mapper.AppendConverters(copierutil.NewTimeStringConverterPair())
	repo.mapper.AppendConverters(copierutil.NewTimeTimestamppbConverterPair())
	repo.mapper.AppendConverters(repo.channelConverter.NewConverterPair())
	repo.mapper.AppendConverters(repo.eventTypeConverter.NewConverterPair())

	return repo
}

func (r *NotificationRuleRepo) List(ctx context.Context, req *paginationV1.PagingRequest) (*notifyV1.ListNotificationRuleResponse, error) {
	if req == nil {
		return nil, errors.BadRequest("NOTIFICATION_RULE", "invalid parameter")
	}

	builder := r.entClient.Client().NotificationRule.Query()

	ret, err := r.repository.ListWithPaging(ctx, builder, builder.Clone(), req)
	if err != nil {
		r.log.Errorf("list notification rules failed: %s", err.Error())
		return nil, errors.InternalServer("NOTIFICATION_RULE", "list notification rules failed")
	}
	if ret == nil {
		return &notifyV1.ListNotificationRuleResponse{Total: 0, Items: nil}, nil
	}

	return &notifyV1.ListNotificationRuleResponse{
		Total: ret.Total,
		Items: ret.Items,
	}, nil
}

func (r *NotificationRuleRepo) Get(ctx context.Context, id uint32) (*notifyV1.NotificationRule, error) {
	if id == 0 {
		return nil, errors.BadRequest("NOTIFICATION_RULE", "id is required")
	}

	entity, err := r.entClient.Client().NotificationRule.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.NotFound("NOTIFICATION_RULE", "notification rule not found")
		}
		r.log.Errorf("get notification rule [%d] failed: %s", id, err.Error())
		return nil, errors.InternalServer("NOTIFICATION_RULE", "get notification rule failed")
	}

	return r.mapper.ToDTO(entity), nil
}

// GetByEventType 按事件类型取路由规则；没有该行时返回 (nil, nil)。
//
// 为什么不返回 NotFound 错误：调用方要区分"没有这条路由"（配置事实，回给业务侧的
// 是一句人话）与"查库失败"（基础设施故障，必须原样上抛）。把两者都压成 error 会让
// SendDirect 在 DB 抖动时报出"没有路由规则"，把人往错的方向引。
func (r *NotificationRuleRepo) GetByEventType(ctx context.Context, eventType notifyV1.EventType) (*notifyV1.NotificationRule, error) {
	entity, err := r.entClient.Client().NotificationRule.Query().
		Where(notificationrule.EventTypeEQ(r.eventTypeString(eventType))).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		r.log.Errorf("query notification rule of event type [%s] failed: %s", eventType.String(), err.Error())
		return nil, errors.InternalServer("NOTIFICATION_RULE", "query notification rule failed")
	}

	return r.mapper.ToDTO(entity), nil
}

// eventTypeString 把 proto 枚举换成实体枚举字面量（两者同名，见 ent schema 的 NamedValues）。
// 不在受支持取值内（含 UNSPECIFIED）返回空串，让唯一索引那一侧的查询自然落空。
func (r *NotificationRuleRepo) eventTypeString(eventType notifyV1.EventType) notificationrule.EventType {
	if name, ok := notifyV1.EventType_name[int32(eventType)]; ok && eventType != notifyV1.EventType_EVENT_TYPE_UNSPECIFIED {
		return notificationrule.EventType(name)
	}
	return ""
}

// Count 规则总行数：播种守卫用（只在空表播，删行不复活是特性）。
func (r *NotificationRuleRepo) Count(ctx context.Context) (int, error) {
	count, err := r.entClient.Client().NotificationRule.Query().Count(ctx)
	if err != nil {
		r.log.Errorf("count notification rules failed: %s", err.Error())
		return 0, errors.InternalServer("NOTIFICATION_RULE", "count notification rules failed")
	}
	return count, nil
}

// Create 建一条规则；operatorID 为 0 表示系统播种（没有操作人）。
//
// created_at 显式写入：mixin 的列没有 DB 默认值，留给 NULL 会让
// "按创建时间排序"的列表页出现一排空行。
func (r *NotificationRuleRepo) Create(ctx context.Context, req *notifyV1.NotificationRule, operatorID uint32) (uint32, error) {
	if req == nil {
		return 0, errors.BadRequest("NOTIFICATION_RULE", "invalid request")
	}
	if req.GetEventType() == notifyV1.EventType_EVENT_TYPE_UNSPECIFIED {
		return 0, errors.BadRequest("NOTIFICATION_RULE", "event_type is required")
	}
	if req.GetChannel() == notifyV1.Channel_CHANNEL_UNSPECIFIED {
		return 0, errors.BadRequest("NOTIFICATION_RULE", "channel is required")
	}

	exist, err := r.entClient.Client().NotificationRule.Query().
		Where(notificationrule.EventTypeEQ(r.eventTypeString(req.GetEventType()))).
		Exist(ctx)
	if err != nil {
		r.log.Errorf("check duplicated event type failed: %s", err.Error())
		return 0, errors.InternalServer("NOTIFICATION_RULE", "check notification rule failed")
	}
	if exist {
		return 0, errors.BadRequest("NOTIFICATION_RULE", "event type "+req.GetEventType().String()+" already has a routing rule")
	}

	builder := r.entClient.Client().NotificationRule.Create().
		SetNillableEventType(r.eventTypeConverter.ToEntity(req.EventType)).
		SetNillableChannel(r.channelConverter.ToEntity(req.Channel)).
		SetNillableIsAsync(req.IsAsync).
		SetNillableIsEnabled(req.IsEnabled).
		SetNillableRemark(req.Remark).
		SetCreatedAt(time.Now())

	if operatorID != 0 {
		builder.SetCreatedBy(operatorID)
	}

	created, err := builder.Save(ctx)
	if err != nil {
		r.log.Errorf("insert notification rule failed: %s", err.Error())
		return 0, errors.InternalServer("NOTIFICATION_RULE", "insert notification rule failed")
	}

	return created.ID, nil
}

func (r *NotificationRuleRepo) Update(ctx context.Context, req *notifyV1.UpdateNotificationRuleRequest, operatorID uint32) (*notifyV1.NotificationRule, error) {
	if req == nil || req.Data == nil {
		return nil, errors.BadRequest("NOTIFICATION_RULE", "invalid request")
	}
	if req.GetId() == 0 {
		return nil, errors.BadRequest("NOTIFICATION_RULE", "id is required")
	}

	exist, err := r.entClient.Client().NotificationRule.Query().
		Where(notificationrule.IDEQ(req.GetId())).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.NotFound("NOTIFICATION_RULE", "notification rule not found")
		}
		r.log.Errorf("get notification rule [%d] failed: %s", req.GetId(), err.Error())
		return nil, errors.InternalServer("NOTIFICATION_RULE", "update notification rule failed")
	}

	// event_type 刻意不在更新面里：一事件一行的唯一锚不随编辑漂移（与上游一致）
	builder := r.entClient.Client().NotificationRule.UpdateOneID(exist.ID)
	if req.Data.Channel != nil {
		builder.SetNillableChannel(r.channelConverter.ToEntity(req.Data.Channel))
	}
	if req.Data.IsAsync != nil {
		builder.SetNillableIsAsync(req.Data.IsAsync)
	}
	if req.Data.IsEnabled != nil {
		builder.SetNillableIsEnabled(req.Data.IsEnabled)
	}
	if req.Data.Remark != nil {
		builder.SetNillableRemark(req.Data.Remark)
	}
	if operatorID != 0 {
		builder.SetUpdatedBy(operatorID)
	}
	builder.SetUpdatedAt(time.Now())

	if _, err := builder.Save(ctx); err != nil {
		r.log.Errorf("update notification rule [%d] failed: %s", req.GetId(), err.Error())
		return nil, errors.InternalServer("NOTIFICATION_RULE", "update notification rule failed")
	}

	return r.Get(ctx, req.GetId())
}

func (r *NotificationRuleRepo) Delete(ctx context.Context, id uint32) error {
	if id == 0 {
		return errors.BadRequest("NOTIFICATION_RULE", "id is required")
	}

	if err := r.entClient.Client().NotificationRule.DeleteOneID(id).Exec(ctx); err != nil {
		if ent.IsNotFound(err) {
			return errors.NotFound("NOTIFICATION_RULE", "notification rule not found")
		}
		r.log.Errorf("delete notification rule [%d] failed: %s", id, err.Error())
		return errors.InternalServer("NOTIFICATION_RULE", "delete notification rule failed")
	}

	return nil
}
