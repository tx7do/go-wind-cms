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
	"github.com/tx7do/go-utils/timeutil"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/notificationdelivery"
	"go-wind-cms/app/core/service/internal/data/ent/predicate"

	notifyV1 "go-wind-cms/api/gen/go/notify/service/v1"
)

// NotificationDeliveryRepo 通知投递台账仓储。
//
// 台账是"投递事实"的记录，生命周期只有三步：投递前 Create(SENDING)，投递后 MarkResult，
// 外加 SweepStaleSending 给"再也不会有结论"的行补一个终态。
// 因此本仓不提供 Update/Delete——改一条已发生的投递记录等于伪造事实；
// 清扫之所以不算伪造，见它的注释。
type NotificationDeliveryRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper

	mapper *mapper.CopierMapper[notifyV1.NotificationDelivery, ent.NotificationDelivery]

	channelConverter   *mapper.EnumTypeConverter[notifyV1.Channel, notificationdelivery.Channel]
	eventTypeConverter *mapper.EnumTypeConverter[notifyV1.EventType, notificationdelivery.EventType]
	statusConverter    *mapper.EnumTypeConverter[notifyV1.DeliveryStatus, notificationdelivery.Status]

	repository *entCrud.Repository[
		ent.NotificationDeliveryQuery, ent.NotificationDeliverySelect,
		ent.NotificationDeliveryCreate, ent.NotificationDeliveryCreateBulk,
		ent.NotificationDeliveryUpdate, ent.NotificationDeliveryUpdateOne,
		ent.NotificationDeliveryDelete,
		predicate.NotificationDelivery,
		notifyV1.NotificationDelivery, ent.NotificationDelivery,
	]
}

func NewNotificationDeliveryRepo(ctx *bootstrap.Context, entClient *entCrud.EntClient[*ent.Client]) *NotificationDeliveryRepo {
	repo := &NotificationDeliveryRepo{
		log:       ctx.NewLoggerHelper("notification-delivery/repo/core-service"),
		entClient: entClient,
		mapper:    mapper.NewCopierMapper[notifyV1.NotificationDelivery, ent.NotificationDelivery](),
		channelConverter: mapper.NewEnumTypeConverter[notifyV1.Channel, notificationdelivery.Channel](
			notifyV1.Channel_name, notifyV1.Channel_value,
		),
		eventTypeConverter: mapper.NewEnumTypeConverter[notifyV1.EventType, notificationdelivery.EventType](
			notifyV1.EventType_name, notifyV1.EventType_value,
		),
		statusConverter: mapper.NewEnumTypeConverter[notifyV1.DeliveryStatus, notificationdelivery.Status](
			notifyV1.DeliveryStatus_name, notifyV1.DeliveryStatus_value,
		),
	}

	repo.repository = entCrud.NewRepository[
		ent.NotificationDeliveryQuery, ent.NotificationDeliverySelect,
		ent.NotificationDeliveryCreate, ent.NotificationDeliveryCreateBulk,
		ent.NotificationDeliveryUpdate, ent.NotificationDeliveryUpdateOne,
		ent.NotificationDeliveryDelete,
		predicate.NotificationDelivery,
		notifyV1.NotificationDelivery, ent.NotificationDelivery,
	](repo.mapper)

	repo.mapper.AppendConverters(copierutil.NewTimeStringConverterPair())
	repo.mapper.AppendConverters(copierutil.NewTimeTimestamppbConverterPair())
	repo.mapper.AppendConverters(repo.channelConverter.NewConverterPair())
	repo.mapper.AppendConverters(repo.eventTypeConverter.NewConverterPair())
	repo.mapper.AppendConverters(repo.statusConverter.NewConverterPair())

	return repo
}

func (r *NotificationDeliveryRepo) List(ctx context.Context, req *paginationV1.PagingRequest) (*notifyV1.ListNotificationDeliveryResponse, error) {
	if req == nil {
		return nil, errors.BadRequest("NOTIFICATION_DELIVERY", "invalid parameter")
	}

	builder := r.entClient.Client().NotificationDelivery.Query()

	ret, err := r.repository.ListWithPaging(ctx, builder, builder.Clone(), req)
	if err != nil {
		r.log.Errorf("list notification deliveries failed: %s", err.Error())
		return nil, errors.InternalServer("NOTIFICATION_DELIVERY", "list notification deliveries failed")
	}
	if ret == nil {
		return &notifyV1.ListNotificationDeliveryResponse{Total: 0, Items: nil}, nil
	}

	return &notifyV1.ListNotificationDeliveryResponse{
		Total: ret.Total,
		Items: ret.Items,
	}, nil
}

func (r *NotificationDeliveryRepo) Get(ctx context.Context, id uint32) (*notifyV1.NotificationDelivery, error) {
	if id == 0 {
		return nil, errors.BadRequest("NOTIFICATION_DELIVERY", "id is required")
	}

	entity, err := r.entClient.Client().NotificationDelivery.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.NotFound("NOTIFICATION_DELIVERY", "notification delivery not found")
		}
		r.log.Errorf("get notification delivery [%d] failed: %s", id, err.Error())
		return nil, errors.InternalServer("NOTIFICATION_DELIVERY", "get notification delivery failed")
	}

	return r.mapper.ToDTO(entity), nil
}

// Create 落一条投递台账（status 由 req 决定，未给时列默认 SENDING）。
//
// attempts 不从 req 取：它是结果列，只有 MarkAttempted / MarkResult 写得动
//（新建时靠列默认 0），否则"台账不提供 Update"这条约束就从旁边漏了个洞。
func (r *NotificationDeliveryRepo) Create(ctx context.Context, req *notifyV1.NotificationDelivery) (*notifyV1.NotificationDelivery, error) {
	if req == nil {
		return nil, errors.BadRequest("NOTIFICATION_DELIVERY", "invalid parameter")
	}

	builder := r.entClient.Client().NotificationDelivery.Create().
		SetNillableEventType(r.eventTypeConverter.ToEntity(req.EventType)).
		SetNillableChannel(r.channelConverter.ToEntity(req.Channel)).
		SetNillableChannelID(req.ChannelId).
		SetNillableRecipientUserID(req.RecipientUserId).
		SetNillableRelatedID(req.RelatedId).
		SetNillableRequestID(req.RequestId).
		SetNillableTarget(req.Target).
		SetNillableStatus(r.statusConverter.ToEntity(req.Status)).
		SetNillableSentAt(timeutil.TimestamppbToTime(req.SentAt)).
		SetNillableCreatedBy(req.CreatedBy).
		SetNillableCreatedAt(timeutil.TimestamppbToTime(req.CreatedAt))

	entity, err := builder.Save(ctx)
	if err != nil {
		r.log.Errorf("insert notification delivery failed: %s", err.Error())
		return nil, errors.InternalServer("NOTIFICATION_DELIVERY", "insert notification delivery failed")
	}

	return r.mapper.ToDTO(entity), nil
}

// DeliveryOutcome 投递结果回写的内容。
type DeliveryOutcome struct {
	Status    notifyV1.DeliveryStatus
	LastError string
	// ChannelID 实际选中的渠道配置；nil 表示不修改（自选渠道时 Create 阶段还不知道）。
	ChannelID *uint32
	// Attempts 已尝试次数；nil 表示不修改（异步路径在拨号前用 MarkAttempted 写过，回写时就不再碰）。
	Attempts *uint32
	// SentAt 投递完成时间；FAILED/SKIPPED 传 nil，列保持空。
	SentAt *time.Time
}

// MarkResult 回写投递结果。
//
// 只按主键定位、只覆盖结果字段：台账其余部分是投递意图的快照，不应被结果回写篡改。
func (r *NotificationDeliveryRepo) MarkResult(ctx context.Context, id uint32, outcome *DeliveryOutcome) error {
	if id == 0 || outcome == nil {
		return errors.BadRequest("NOTIFICATION_DELIVERY", "id and outcome are required")
	}

	update := r.entClient.Client().NotificationDelivery.UpdateOneID(id).
		SetStatus(*r.statusConverter.ToEntity(trans.Ptr(outcome.Status)))
	if outcome.LastError != "" {
		update.SetLastError(outcome.LastError)
	}
	if outcome.ChannelID != nil {
		update.SetChannelID(*outcome.ChannelID)
	}
	if outcome.Attempts != nil {
		update.SetAttempts(*outcome.Attempts)
	}
	if outcome.SentAt != nil {
		update.SetSentAt(*outcome.SentAt)
	}

	if err := update.Exec(ctx); err != nil {
		r.log.Errorf("mark notification delivery [%d] result failed: %s", id, err.Error())
		return errors.InternalServer("NOTIFICATION_DELIVERY", "mark notification delivery result failed")
	}

	return nil
}

// MarkAttempted 在拨号之前把"第 attempts 次尝试已经开始"写进台账。
//
// 为什么不跟结果一起写：进程在 SMTP 握手中途被杀，这一行会永远停在 SENDING，
// 而 attempts 是那次"发了一半"唯一留得下来的痕迹（没有它，重试到第几次、有没有真的拨过号
// 全都看不出来，异步派发就成了黑盒）。
func (r *NotificationDeliveryRepo) MarkAttempted(ctx context.Context, id uint32, attempts uint32) error {
	if id == 0 {
		return errors.BadRequest("NOTIFICATION_DELIVERY", "id is required")
	}

	if err := r.entClient.Client().NotificationDelivery.UpdateOneID(id).
		SetAttempts(attempts).
		Exec(ctx); err != nil {
		r.log.Errorf("mark notification delivery [%d] attempted failed: %s", id, err.Error())
		return errors.InternalServer("NOTIFICATION_DELIVERY", "mark notification delivery attempt failed")
	}

	return nil
}

// SweepStaleSending 把 created_at 早于 before 却仍停在 SENDING 的台账行落 FAILED，返回本次清扫条数。
//
// 这一条为什么不算"伪造事实"：MarkResult 拒绝的是**改写已有结论**，而这里的行压根没有结论——
// 它既没发出去也没失败记录，留着 SENDING 只会让台账页出现一排"永远在发"的行。
// 落的是"超期未结算"这一确实发生过的事实，attempts 与 last_error 里最后一次真实报错都保留。
//
// 为什么用"查 ID → UpdateMany 再带 status 谓词"而不是逐行 MarkResult：
// 逐行写法把"这一行现在是什么状态"读进内存后无条件覆盖，撞上正在收尾的异步 handler
// 就会把刚写好的 SENT 改回 FAILED。把 status==SENDING 放进 WHERE，比较就交给 DB 在同一条
// 语句里做，抢跑的那一行自然不命中、也就不会被扫。
//
// 判据用 created_at 而不是 updated_at：台账这张表的 updated_at 从来没被写过
//（mixin 的列是 Optional 无默认，MarkResult/MarkAttempted 都不碰它），拿它当锚等于把所有行
// 判成"从未更新"；而生成的 (status, created_at) 索引正好服务这条查询。
// 异步派发的尝试预算本身也有界（4 次 × 30s + 退避 ≈ 4 分钟），
// 从"意图产生"起算的超期阈值只要明显大于它就是安全的。
func (r *NotificationDeliveryRepo) SweepStaleSending(ctx context.Context, before time.Time, reason string, limit int) (int, error) {
	if limit <= 0 {
		return 0, errors.BadRequest("NOTIFICATION_DELIVERY", "limit must be positive")
	}

	// 比较侧不转 UTC（与上游不同）：本仓 sqlite 驱动把时间写成本地时区的文本（+0800），
	// Postgres timestamptz 按瞬间比较无所谓，但 sqlite 是字典序比较——比较侧一旦转成
	// UTC 渲染（+0000），"一小时前(本地 21:xx)" 对 "15 分钟前(UTC 06:xx)" 的字典序反而更大，
	// 清扫会一行都扫不到。两侧都用本地时间渲染，字典序才与时间序一致。
	// （上游写侧经 timestamppb 往返、渲染成 UTC，所以它必须转；本仓写侧不经那条路。）

	ids, err := r.entClient.Client().NotificationDelivery.Query().
		Where(
			notificationdelivery.StatusEQ(notificationdelivery.StatusSending),
			notificationdelivery.CreatedAtLT(before),
		).
		Limit(limit).
		IDs(ctx)
	if err != nil {
		r.log.Errorf("scan stale notification deliveries failed: %s", err.Error())
		return 0, errors.InternalServer("NOTIFICATION_DELIVERY", "scan stale notification deliveries failed")
	}
	if len(ids) == 0 {
		return 0, nil
	}

	count, err := r.entClient.Client().NotificationDelivery.Update().
		Where(
			notificationdelivery.IDIn(ids...),
			notificationdelivery.StatusEQ(notificationdelivery.StatusSending),
		).
		SetStatus(notificationdelivery.StatusFailed).
		SetLastError(reason).
		Save(ctx)
	if err != nil {
		r.log.Errorf("sweep stale notification deliveries failed: %s", err.Error())
		return 0, errors.InternalServer("NOTIFICATION_DELIVERY", "sweep stale notification deliveries failed")
	}

	return count, nil
}
