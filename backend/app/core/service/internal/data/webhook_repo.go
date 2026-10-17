package data

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/timestamppb"

	entCrud "github.com/tx7do/go-crud/entgo"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/tx7do/go-utils/trans"

	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/predicate"
	"go-wind-cms/app/core/service/internal/data/ent/webhook"
	"go-wind-cms/app/core/service/internal/data/ent/webhookdelivery"
	webhookV1 "go-wind-cms/api/gen/go/webhook/service/v1"
)

// WebhookRepo Webhook 订阅仓储 + 投递记录读写。
type WebhookRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper
}

func NewWebhookRepo(ctx *bootstrap.Context, entClient *entCrud.EntClient[*ent.Client]) *WebhookRepo {
	return &WebhookRepo{
		entClient: entClient,
		log:       ctx.NewLoggerHelper("webhook/repo/core-service"),
	}
}

// List 按 admin 管理列表查询（关键字过滤 description/url，全量不加分页由前端翻页时传 page/pageSize）。
func (r *WebhookRepo) List(ctx context.Context) ([]*webhookV1.Webhook, int, error) {
	conds := r.tenantConds(ctx)
	rows, err := r.entClient.Client().Webhook.Query().
		Where(conds...).
		Order(ent.Desc(webhook.FieldID)).
		All(ctx)
	if err != nil {
		r.log.Errorf("list webhooks failed: %s", err.Error())
		return nil, 0, err
	}
	items := make([]*webhookV1.Webhook, 0, len(rows))
	for _, row := range rows {
		items = append(items, r.toDTO(row))
	}
	return items, len(items), nil
}

// Get 取单个订阅。
func (r *WebhookRepo) Get(ctx context.Context, id uint32) (*webhookV1.Webhook, error) {
	conds := append(r.tenantConds(ctx), webhook.IDEQ(id))
	row, err := r.entClient.Client().Webhook.Query().
		Where(conds...).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.NotFound("WEBHOOK", "webhook not found")
		}
		r.log.Errorf("get webhook failed (id=%d): %s", id, err.Error())
		return nil, err
	}
	return r.toDTO(row), nil
}

// Create 新建订阅。
func (r *WebhookRepo) Create(ctx context.Context, dto *webhookV1.Webhook) (*webhookV1.Webhook, error) {
	builder := r.entClient.Client().Webhook.Create().
		SetURL(dto.GetUrl()).
		SetEnabled(dto.GetEnabled())
	if dto.GetSecret() != "" {
		builder.SetSecret(dto.GetSecret())
	}
	if len(dto.GetEvents()) > 0 {
		events := dto.GetEvents()
		builder.SetEvents(&events)
	}
	if dto.GetDescription() != "" {
		builder.SetDescription(dto.GetDescription())
	}

	if callerUserID, hasUser := viewerUserIDFromContext(ctx); hasUser {
		builder.SetCreatedBy(callerUserID).SetUpdatedBy(callerUserID)
	}
	if tid, hasTenant := maybeTenantFromViewer(ctx); hasTenant {
		builder.SetTenantID(tid)
	}
	row, err := builder.Save(ctx)
	if err != nil {
		r.log.Errorf("create webhook failed: %s", err.Error())
		return nil, err
	}
	return r.toDTO(row), nil
}

// Update 按 mask 更新（与页面/文章仓库同款：nil 字段不覆盖）。
func (r *WebhookRepo) Update(ctx context.Context, id uint32, dto *webhookV1.Webhook) (*webhookV1.Webhook, error) {
	conds := append(r.tenantConds(ctx), webhook.IDEQ(id))
	exist, err := r.entClient.Client().Webhook.Query().Where(conds...).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.NotFound("WEBHOOK", "webhook not found")
		}
		return nil, err
	}

	builder := r.entClient.Client().Webhook.UpdateOneID(exist.ID)
	if dto.Url != nil {
		builder.SetURL(dto.GetUrl())
	}
	if dto.Secret != nil {
		if dto.GetSecret() == "" {
			builder.ClearSecret()
		} else {
			builder.SetSecret(dto.GetSecret())
		}
	}
	if dto.Events != nil {
		if len(dto.GetEvents()) == 0 {
			builder.ClearEvents()
		} else {
			events := dto.GetEvents()
			builder.SetEvents(&events)
		}
	}
	if dto.Enabled != nil {
		builder.SetEnabled(dto.GetEnabled())
	}
	if dto.Description != nil {
		if dto.GetDescription() == "" {
			builder.ClearDescription()
		} else {
			builder.SetDescription(dto.GetDescription())
		}
	}
	if callerUserID, hasUser := viewerUserIDFromContext(ctx); hasUser {
		builder.SetUpdatedBy(callerUserID)
	}

	row, err := builder.Save(ctx)
	if err != nil {
		r.log.Errorf("update webhook failed (id=%d): %s", id, err.Error())
		return nil, err
	}
	return r.toDTO(row), nil
}

// Delete 删除订阅。
func (r *WebhookRepo) Delete(ctx context.Context, id uint32) error {
	conds := append(r.tenantConds(ctx), webhook.IDEQ(id))
	count, err := r.entClient.Client().Webhook.Delete().Where(conds...).Exec(ctx)
	if err != nil {
		r.log.Errorf("delete webhook failed (id=%d): %s", id, err.Error())
		return err
	}
	if count == 0 {
		return errors.NotFound("WEBHOOK", "webhook not found")
	}
	return nil
}

// ListEnabledForEvent 列出（租户内）订阅了指定事件且启用的 Webhook；events 为空视为订阅全部。
func (r *WebhookRepo) ListEnabledForEvent(ctx context.Context, tenantID uint32, event string) ([]*webhookV1.Webhook, error) {
	conds := []predicate.Webhook{
		webhook.EnabledEQ(true),
		webhook.TenantIDEQ(tenantID),
	}
	rows, err := r.entClient.Client().Webhook.Query().
		Where(conds...).
		All(ctx)
	if err != nil {
		r.log.Errorf("list webhooks for event failed (event=%s): %s", event, err.Error())
		return nil, err
	}
	matched := make([]*webhookV1.Webhook, 0, len(rows))
	for _, row := range rows {
		var events []string
		if row.Events != nil {
			events = *row.Events
		}
		if len(events) == 0 {
			matched = append(matched, r.toDTO(row))
			continue
		}
		for _, e := range events {
			if e == event {
				matched = append(matched, r.toDTO(row))
				break
			}
		}
	}
	return matched, nil
}

// ─── 投递记录 ───

// CreateDelivery 落一条投递记录（pending 由调用方以 status 表达：先 FAILED 占位，完成后 Update）。
func (r *WebhookRepo) CreateDelivery(ctx context.Context, webhookID uint32, event string, entityID uint32, tenantID uint32) (uint32, error) {
	builder := r.entClient.Client().WebhookDelivery.Create().
		SetNillableWebhookID(&webhookID).
		SetNillableEvent(&event).
		SetNillableEntityID(&entityID).
		SetNillableStatus(ptrDeliveryStatus("DELIVERY_STATUS_FAILED"))
	if tenantID > 0 {
		builder.SetTenantID(tenantID)
	}
	row, err := builder.Save(ctx)
	if err != nil {
		r.log.Errorf("create webhook delivery failed (webhook=%d): %s", webhookID, err.Error())
		return 0, err
	}
	return row.ID, nil
}

// UpdateDeliveryResult 回写投递结果。
func (r *WebhookRepo) UpdateDeliveryResult(ctx context.Context, id uint32, success bool, httpStatus int, errMsg string, durationMs uint32) error {
	builder := r.entClient.Client().WebhookDelivery.UpdateOneID(id)
	if success {
		builder.SetNillableStatus(ptrDeliveryStatus("DELIVERY_STATUS_SUCCESS"))
	} else {
		builder.SetNillableStatus(ptrDeliveryStatus("DELIVERY_STATUS_FAILED"))
	}
	builder.SetHTTPStatus(httpStatus)
	if errMsg != "" {
		builder.SetError(errMsg)
	}
	builder.SetDurationMs(durationMs)
	if err := builder.Exec(ctx); err != nil {
		r.log.Errorf("update webhook delivery failed (id=%d): %s", id, err.Error())
		return err
	}
	return nil
}

// ListDeliveries 投递记录（时间倒序）。
func (r *WebhookRepo) ListDeliveries(ctx context.Context, webhookID uint32, limit int) ([]*webhookV1.WebhookDelivery, int, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	var conds []predicate.WebhookDelivery
	if webhookID > 0 {
		conds = append(conds, webhookdelivery.WebhookIDEQ(webhookID))
	}
	rows, err := r.entClient.Client().WebhookDelivery.Query().
		Where(conds...).
		Order(ent.Desc(webhookdelivery.FieldID)).
		Limit(limit).
		All(ctx)
	if err != nil {
		r.log.Errorf("list webhook deliveries failed: %s", err.Error())
		return nil, 0, err
	}
	items := make([]*webhookV1.WebhookDelivery, 0, len(rows))
	for _, row := range rows {
		dto := &webhookV1.WebhookDelivery{
			Id:         &row.ID,
			Event:      row.Event,
			EntityId:   row.EntityID,
			DurationMs: row.DurationMs,
		}
		if row.HTTPStatus != nil {
			dto.HttpStatus = trans.Ptr(int32(*row.HTTPStatus))
		}
		if row.CreatedAt != nil {
			dto.CreatedAt = timestamppb.New(*row.CreatedAt)
		}
		if row.Status != nil {
			dto.Status = trans.Ptr(string(*row.Status))
		}
		if row.Error != nil {
			dto.Error = row.Error
		}
		items = append(items, dto)
	}
	return items, len(items), nil
}

// PruneDeliveriesOlderThan 裁剪指定时间之前的投递记录（限批，防长事务）。
// 表只增不改：旧数据由周期任务定期裁剪，防止无限膨胀。
func (r *WebhookRepo) PruneDeliveriesOlderThan(ctx context.Context, before time.Time, limit int) (int, error) {
	ids, err := r.entClient.Client().WebhookDelivery.Query().
		Where(webhookdelivery.CreatedAtLT(before)).
		IDs(ctx)
	if err != nil {
		r.log.Errorf("prune webhook deliveries: list ids failed: %s", err.Error())
		return 0, err
	}
	if len(ids) > limit {
		ids = ids[:limit]
	}
	if len(ids) == 0 {
		return 0, nil
	}
	n, err := r.entClient.Client().WebhookDelivery.Delete().
		Where(webhookdelivery.IDIn(ids...)).
		Exec(ctx)
	if err != nil {
		r.log.Errorf("prune webhook deliveries failed: %s", err.Error())
		return 0, err
	}
	return n, nil
}

// ─── 内部 ───

func (r *WebhookRepo) tenantConds(ctx context.Context) []predicate.Webhook {
	tid, hasTenant := maybeTenantFromViewer(ctx)
	if !hasTenant {
		return nil
	}
	return []predicate.Webhook{webhook.TenantIDEQ(tid)}
}

func (r *WebhookRepo) toDTO(row *ent.Webhook) *webhookV1.Webhook {
	enabled := true
	if row.Enabled != nil {
		enabled = *row.Enabled
	}
	dto := &webhookV1.Webhook{
		Id:          &row.ID,
		Url:         &row.URL,
		Enabled:     &enabled,
		Secret:      row.Secret,
		Description: row.Description,
		CreatedBy:   row.CreatedBy,
		UpdatedBy:   row.UpdatedBy,
	}
	if row.Events != nil {
		dto.Events = *row.Events
	} else {
		dto.Events = []string{}
	}
	if row.CreatedAt != nil {
		dto.CreatedAt = timestamppb.New(*row.CreatedAt)
	}
	if row.UpdatedAt != nil {
		dto.UpdatedAt = timestamppb.New(*row.UpdatedAt)
	}
	return dto
}

func ptrDeliveryStatus(v string) *webhookdelivery.Status {
	s := webhookdelivery.Status(v)
	return &s
}
