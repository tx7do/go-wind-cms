package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	entCrud "github.com/tx7do/go-crud/entgo"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/reviewlog"
)

// ReviewLogRepo 发布审批留痕仓储（sys_review_logs）。
type ReviewLogRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper
}

func NewReviewLogRepo(ctx *bootstrap.Context, entClient *entCrud.EntClient[*ent.Client]) *ReviewLogRepo {
	return &ReviewLogRepo{
		log:       ctx.NewLoggerHelper("review-log/repo/core-service"),
		entClient: entClient,
	}
}

// Create 落一条留痕（operatorID/tenantID 由调用方显式传入——写入方在服务层从 viewer 取）。
func (r *ReviewLogRepo) Create(ctx context.Context, entityType string, entityID uint32, action string, level uint32, note string, operatorID, tenantID uint32) error {
	builder := r.entClient.Client().ReviewLog.Create().
		SetEntityType(entityType).
		SetEntityID(entityID).
		SetAction(reviewlog.Action(action)).
		SetLevel(level).
		SetCreatedBy(operatorID)
	if note != "" {
		builder.SetNote(note)
	}
	if tenantID > 0 {
		builder.SetTenantID(tenantID)
	}
	if err := builder.Exec(ctx); err != nil {
		r.log.Errorf("create review log failed: %s", err.Error())
		return contentV1.ErrorInternalServerError("create review log failed")
	}
	return nil
}

// CurrentLevel 取实体当前审批推进到的级别（本批次内最大 APPROVE level；无记录=0）。
func (r *ReviewLogRepo) CurrentLevel(ctx context.Context, entityType string, entityID uint32) (uint32, error) {
	rows, err := r.entClient.Client().ReviewLog.Query().
		Where(
			reviewlog.EntityTypeEQ(entityType),
			reviewlog.EntityIDEQ(entityID),
			reviewlog.ActionEQ(reviewlog.ActionApprove),
		).
		All(ctx)
	if err != nil {
		r.log.Errorf("query review logs failed: %s", err.Error())
		return 0, contentV1.ErrorInternalServerError("query review logs failed")
	}
	var max uint32
	for _, row := range rows {
		if row.Level > max {
			max = row.Level
		}
	}
	return max, nil
}

// ListByEntity 按实体取全部留痕（created_at 正序，id 兜底排序）。
func (r *ReviewLogRepo) ListByEntity(ctx context.Context, entityType string, entityID uint32) ([]*contentV1.ReviewLog, error) {
	rows, err := r.entClient.Client().ReviewLog.Query().
		Where(
			reviewlog.EntityTypeEQ(entityType),
			reviewlog.EntityIDEQ(entityID),
		).
		Order(ent.Asc(reviewlog.FieldID)).
		All(ctx)
	if err != nil {
		r.log.Errorf("list review logs failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("list review logs failed")
	}
	out := make([]*contentV1.ReviewLog, 0, len(rows))
	for _, e := range rows {
		dto := &contentV1.ReviewLog{
			Id:         &e.ID,
			EntityType: e.EntityType,
			EntityId:   e.EntityID,
			Action:     e.Action.String(),
			Level:      e.Level,
			OperatorId: e.CreatedBy,
		}
		if e.Note != nil {
			dto.Note = e.Note
		}
		out = append(out, dto)
	}
	return out, nil
}
