package data

import (
	"context"
	"strconv"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	kratosMetadata "github.com/go-kratos/kratos/v2/metadata"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	paginationFilter "github.com/tx7do/go-crud/pagination/filter"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"entgo.io/ent/dialect/sql"
	"google.golang.org/protobuf/types/known/timestamppb"

	entCrud "github.com/tx7do/go-crud/entgo"

	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/experiment"
	"go-wind-cms/app/core/service/internal/data/ent/experimentparticipant"
	"go-wind-cms/app/core/service/internal/data/ent/experimentvariant"
	"go-wind-cms/app/core/service/internal/data/ent/predicate"
)

// ExperimentRepo 内容实验仓储（实验/变体/参与记录三表）。
// 变体随实验整组托管（写时全删重建，对齐 page/sections 模式；
// 与仓内惯例一致不建 ent edge，子表按外键手工水合）。
type ExperimentRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper
}

func NewExperimentRepo(ctx *bootstrap.Context, entClient *entCrud.EntClient[*ent.Client]) *ExperimentRepo {
	return &ExperimentRepo{
		entClient: entClient,
		log:       ctx.NewLoggerHelper("experiment/repo/core-service"),
	}
}

// List 实验列表（时间倒序；支持 name 模糊 / status / pageId 过滤）。
func (r *ExperimentRepo) List(ctx context.Context, req *paginationV1.PagingRequest) ([]*contentV1.Experiment, int, error) {
	conds := r.tenantConds(ctx)
	if req != nil {
		if filterExpr, fErr := paginationFilter.ConvertFilterByPagingRequest(req); fErr == nil && filterExpr != nil {
			for _, cond := range filterExpr.GetConditions() {
				switch cond.GetField() {
				case "name":
					if v := cond.GetValue(); v != "" {
						conds = append(conds, experiment.NameContains(v))
					}
				case "status":
					if st, ok := statusStringToEnt(cond.GetValue()); ok {
						conds = append(conds, experiment.StatusEQ(st))
					}
				case "pageId", "page_id":
					if v, pErr := strconv.ParseUint(cond.GetValue(), 10, 32); pErr == nil && v > 0 {
						conds = append(conds, experiment.PageIDEQ(uint32(v)))
					}
				}
			}
		}
	}

	total, err := r.entClient.Client().Experiment.Query().Where(conds...).Count(ctx)
	if err != nil {
		r.log.Errorf("count experiments failed: %s", err.Error())
		return nil, 0, err
	}

	page, pageSize := uint32(1), uint32(50)
	if req != nil {
		if req.GetPage() > 0 {
			page = req.GetPage()
		}
		if req.GetPageSize() > 0 {
			pageSize = req.GetPageSize()
		}
	}
	if pageSize > 500 {
		pageSize = 500
	}

	rows, err := r.entClient.Client().Experiment.Query().Where(conds...).
		Order(ent.Desc(experiment.FieldID)).
		Limit(int(pageSize)).
		Offset(int((page - 1) * pageSize)).
		All(ctx)
	if err != nil {
		r.log.Errorf("list experiments failed: %s", err.Error())
		return nil, 0, err
	}

	items := make([]*contentV1.Experiment, 0, len(rows))
	for _, row := range rows {
		items = append(items, r.toDTO(ctx, row))
	}
	return items, total, nil
}

// ListRunningByPage 某页面上运行中的实验（app 分流入口，匿名链路可达）。
func (r *ExperimentRepo) ListRunningByPage(ctx context.Context, pageID uint32) ([]*contentV1.Experiment, error) {
	conds := append(r.tenantConds(ctx),
		experiment.StatusEQ(experiment.StatusRunning),
		experiment.PageIDEQ(pageID),
	)
	rows, err := r.entClient.Client().Experiment.Query().Where(conds...).
		Order(ent.Asc(experiment.FieldID)).
		All(ctx)
	if err != nil {
		r.log.Errorf("list running experiments by page (page_id=%d) failed: %s", pageID, err.Error())
		return nil, err
	}
	items := make([]*contentV1.Experiment, 0, len(rows))
	for _, row := range rows {
		items = append(items, r.toDTO(ctx, row))
	}
	return items, nil
}

// Get 实验详情（含变体）。
func (r *ExperimentRepo) Get(ctx context.Context, id uint32) (*contentV1.Experiment, error) {
	conds := append(r.tenantConds(ctx), experiment.IDEQ(id))
	row, err := r.entClient.Client().Experiment.Query().Where(conds...).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.NotFound("EXPERIMENT", "experiment not found")
		}
		r.log.Errorf("get experiment failed (id=%d): %s", id, err.Error())
		return nil, err
	}
	return r.toDTO(ctx, row), nil
}

// Create 新建实验（连同变体）。
func (r *ExperimentRepo) Create(ctx context.Context, dto *contentV1.Experiment) (*contentV1.Experiment, error) {
	builder := r.entClient.Client().Experiment.Create().
		SetName(dto.GetName()).
		SetNillablePageID(dto.PageId).
		SetStatus(statusToEntOrDefault(dto.GetStatus())).
		SetAssignmentType(assignmentTypeToEntOrDefault(dto.GetAssignmentType()))
	if dto.GetDescription() != "" {
		builder.SetDescription(dto.GetDescription())
	}
	if callerUserID, hasUser := viewerUserIDFromContext(ctx); hasUser {
		builder.SetCreatedBy(callerUserID).SetUpdatedBy(callerUserID)
	}
	if tid, hasTenant := maybeTenantFromViewer(ctx); hasTenant {
		builder.SetTenantID(tid)
	} else if tid := r.tenantIDFromInternalMetadata(ctx); tid > 0 {
		builder.SetTenantID(tid)
	}

	row, err := builder.Save(ctx)
	if err != nil {
		r.log.Errorf("create experiment failed: %s", err.Error())
		return nil, err
	}
	if err := r.replaceVariants(ctx, row.ID, dto.GetVariants()); err != nil {
		return nil, err
	}
	return r.Get(ctx, row.ID)
}

// Update 按 mask 更新；variants 出现在 data 里即整组替换。
func (r *ExperimentRepo) Update(ctx context.Context, id uint32, dto *contentV1.Experiment, dtoHasVariants bool) (*contentV1.Experiment, error) {
	conds := append(r.tenantConds(ctx), experiment.IDEQ(id))
	exist, err := r.entClient.Client().Experiment.Query().Where(conds...).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.NotFound("EXPERIMENT", "experiment not found")
		}
		return nil, err
	}

	builder := r.entClient.Client().Experiment.UpdateOneID(exist.ID)
	if dto.Name != nil {
		builder.SetName(dto.GetName())
	}
	if dto.PageId != nil {
		builder.SetNillablePageID(dto.PageId)
	}
	if dto.Status != nil {
		builder.SetStatus(statusToEntOrDefault(dto.GetStatus()))
	}
	if dto.AssignmentType != nil {
		builder.SetAssignmentType(assignmentTypeToEntOrDefault(dto.GetAssignmentType()))
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

	if _, err := builder.Save(ctx); err != nil {
		r.log.Errorf("update experiment failed (id=%d): %s", id, err.Error())
		return nil, err
	}
	if dtoHasVariants {
		if err := r.replaceVariants(ctx, exist.ID, dto.GetVariants()); err != nil {
			return nil, err
		}
	}
	return r.Get(ctx, id)
}

// SetStatus 状态机流转（服务层已校验，这里落库）。
func (r *ExperimentRepo) SetStatus(ctx context.Context, id uint32, status contentV1.Experiment_Status) (*contentV1.Experiment, error) {
	conds := append(r.tenantConds(ctx), experiment.IDEQ(id))
	exist, err := r.entClient.Client().Experiment.Query().Where(conds...).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.NotFound("EXPERIMENT", "experiment not found")
		}
		return nil, err
	}
	builder := r.entClient.Client().Experiment.UpdateOneID(exist.ID).SetStatus(statusToEntOrDefault(status))
	if callerUserID, hasUser := viewerUserIDFromContext(ctx); hasUser {
		builder.SetUpdatedBy(callerUserID)
	}
	if _, err := builder.Save(ctx); err != nil {
		r.log.Errorf("set experiment status failed (id=%d): %s", id, err.Error())
		return nil, err
	}
	return r.Get(ctx, id)
}

// Delete 删除实验（变体与参与记录一并清理）。
func (r *ExperimentRepo) Delete(ctx context.Context, id uint32) error {
	conds := append(r.tenantConds(ctx), experiment.IDEQ(id))
	exist, err := r.entClient.Client().Experiment.Query().Where(conds...).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return errors.NotFound("EXPERIMENT", "experiment not found")
		}
		return err
	}
	if _, err := r.entClient.Client().ExperimentParticipant.Delete().
		Where(experimentparticipant.ExperimentIDEQ(exist.ID)).Exec(ctx); err != nil {
		r.log.Errorf("delete experiment participants failed (id=%d): %s", id, err.Error())
		return err
	}
	if _, err := r.entClient.Client().ExperimentVariant.Delete().
		Where(experimentvariant.ExperimentIDEQ(exist.ID)).Exec(ctx); err != nil {
		r.log.Errorf("delete experiment variants failed (id=%d): %s", id, err.Error())
		return err
	}
	if err := r.entClient.Client().Experiment.DeleteOneID(exist.ID).Exec(ctx); err != nil {
		r.log.Errorf("delete experiment failed (id=%d): %s", id, err.Error())
		return err
	}
	return nil
}

// RecordParticipant 幂等落参与记录（唯一索引冲突忽略）。
// 变体名以快照存储，历史归因不受变体改名影响。
func (r *ExperimentRepo) RecordParticipant(ctx context.Context, experimentID uint32, visitorHash, variantName string) error {
	conds := append(r.tenantConds(ctx), experiment.IDEQ(experimentID))
	exist, err := r.entClient.Client().Experiment.Query().Where(conds...).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return errors.NotFound("EXPERIMENT", "experiment not found")
		}
		return err
	}

	builder := r.entClient.Client().ExperimentParticipant.Create().
		SetExperimentID(exist.ID).
		SetVisitorHash(visitorHash).
		SetVariantName(variantName)
	if exist.TenantID != nil {
		builder.SetTenantID(*exist.TenantID)
	}

	err = builder.OnConflict(
		sql.ConflictColumns(
			experimentparticipant.FieldExperimentID,
			experimentparticipant.FieldVisitorHash,
		),
	).DoNothing().Exec(ctx)
	if err != nil {
		r.log.Errorf("record experiment participant failed (experiment=%d): %s", experimentID, err.Error())
		return err
	}
	return nil
}

// GetVariantStats 按变体名快照计数（空实验返回空列表）。
func (r *ExperimentRepo) GetVariantStats(ctx context.Context, experimentID uint32) (*contentV1.ExperimentVariantStatsResponse, error) {
	var rows []struct {
		VariantName string `json:"variant_name"`
		Count       int    `json:"count"`
	}
	err := r.entClient.Client().ExperimentParticipant.Query().
		Where(experimentparticipant.ExperimentIDEQ(experimentID)).
		GroupBy(experimentparticipant.FieldVariantName).
		Aggregate(ent.As(ent.Count(), "count")).
		Scan(ctx, &rows)
	if err != nil {
		r.log.Errorf("get experiment variant stats failed (experiment=%d): %s", experimentID, err.Error())
		return nil, err
	}
	resp := &contentV1.ExperimentVariantStatsResponse{
		Items: make([]*contentV1.ExperimentVariantStatsResponse_VariantCount, 0, len(rows)),
	}
	for _, row := range rows {
		resp.Items = append(resp.Items, &contentV1.ExperimentVariantStatsResponse_VariantCount{
			VariantName:  row.VariantName,
			Participants: uint64(row.Count),
		})
	}
	return resp, nil
}

// replaceVariants 变体整组替换（写时全删重建）。
func (r *ExperimentRepo) replaceVariants(ctx context.Context, experimentID uint32, variants []*contentV1.ExperimentVariant) error {
	if _, err := r.entClient.Client().ExperimentVariant.Delete().
		Where(experimentvariant.ExperimentIDEQ(experimentID)).Exec(ctx); err != nil {
		r.log.Errorf("clear experiment variants failed (experiment=%d): %s", experimentID, err.Error())
		return err
	}
	for _, v := range variants {
		if v == nil {
			continue
		}
		builder := r.entClient.Client().ExperimentVariant.Create().
			SetExperimentID(experimentID).
			SetName(v.GetName()).
			SetNillableSortOrder(v.SortOrder)
		if v.Weight != nil {
			builder.SetWeight(int(*v.Weight))
		}
		if v.GetConditions() != "" {
			builder.SetConditions(v.GetConditions())
		}
		if v.GetSections() != "" {
			builder.SetSections(v.GetSections())
		}
		if _, err := builder.Save(ctx); err != nil {
			r.log.Errorf("create experiment variant failed (experiment=%d name=%s): %s", experimentID, v.GetName(), err.Error())
			return err
		}
	}
	return nil
}

// listVariants 实验的变体（sort_order,id 升序；rule 型即求值顺序）。
func (r *ExperimentRepo) listVariants(ctx context.Context, experimentID uint32) ([]*contentV1.ExperimentVariant, error) {
	rows, err := r.entClient.Client().ExperimentVariant.Query().
		Where(experimentvariant.ExperimentIDEQ(experimentID)).
		Order(ent.Asc(experimentvariant.FieldSortOrder), ent.Asc(experimentvariant.FieldID)).
		All(ctx)
	if err != nil {
		r.log.Errorf("query experiment variants failed (experiment=%d): %s", experimentID, err.Error())
		return nil, err
	}
	items := make([]*contentV1.ExperimentVariant, 0, len(rows))
	for _, v := range rows {
		items = append(items, &contentV1.ExperimentVariant{
			Id:           &v.ID,
			ExperimentId: v.ExperimentID,
			Name:         &v.Name,
			Weight:       weightToInt32(v.Weight),
			Conditions:   nillableString(v.Conditions),
			Sections:     nillableString(v.Sections),
			SortOrder:    v.SortOrder,
		})
	}
	return items, nil
}

func (r *ExperimentRepo) toDTO(ctx context.Context, row *ent.Experiment) *contentV1.Experiment {
	if row == nil {
		return nil
	}
	dto := &contentV1.Experiment{
		Id:          &row.ID,
		Name:        &row.Name,
		PageId:      row.PageID,
		Description: row.Description,
	}
	if row.Status != nil {
		dto.Status = entStatusToProto(*row.Status)
	}
	if row.AssignmentType != nil {
		dto.AssignmentType = entAssignmentToProto(*row.AssignmentType)
	}
	if row.CreatedAt != nil {
		dto.CreatedAt = timestamppb.New(*row.CreatedAt)
	}
	if row.UpdatedAt != nil {
		dto.UpdatedAt = timestamppb.New(*row.UpdatedAt)
	}
	if row.CreatedBy != nil {
		dto.CreatedBy = row.CreatedBy
	}
	if row.UpdatedBy != nil {
		dto.UpdatedBy = row.UpdatedBy
	}
	variants, err := r.listVariants(ctx, row.ID)
	if err != nil {
		r.log.Warnf("hydrate experiment variants failed (experiment=%d): %s", row.ID, err.Error())
	} else {
		dto.Variants = variants
	}
	return dto
}

func (r *ExperimentRepo) tenantConds(ctx context.Context) []predicate.Experiment {
	tid, hasTenant := maybeTenantFromViewer(ctx)
	if !hasTenant {
		return nil
	}
	return []predicate.Experiment{experiment.TenantIDEQ(tid)}
}

// tenantIDFromInternalMetadata 从内网可信元数据 x-md-global-tenant-id 读租户。
func (r *ExperimentRepo) tenantIDFromInternalMetadata(ctx context.Context) uint32 {
	if md, ok := kratosMetadata.FromServerContext(ctx); ok {
		if v := md.Get("x-md-global-tenant-id"); v != "" {
			if tid, err := strconv.ParseUint(v, 10, 32); err == nil && tid > 0 {
				return uint32(tid)
			}
		}
	}
	return 0
}

// ── 枚举映射 ──

func statusStringToEnt(v string) (experiment.Status, bool) {
	switch v {
	case "draft", "EXPERIMENT_STATUS_DRAFT":
		return experiment.StatusDraft, true
	case "running", "EXPERIMENT_STATUS_RUNNING":
		return experiment.StatusRunning, true
	case "paused", "EXPERIMENT_STATUS_PAUSED":
		return experiment.StatusPaused, true
	case "ended", "EXPERIMENT_STATUS_ENDED":
		return experiment.StatusEnded, true
	}
	return experiment.StatusDraft, false
}

func statusToEntOrDefault(v contentV1.Experiment_Status) experiment.Status {
	switch v {
	case contentV1.Experiment_EXPERIMENT_STATUS_DRAFT:
		return experiment.StatusDraft
	case contentV1.Experiment_EXPERIMENT_STATUS_RUNNING:
		return experiment.StatusRunning
	case contentV1.Experiment_EXPERIMENT_STATUS_PAUSED:
		return experiment.StatusPaused
	case contentV1.Experiment_EXPERIMENT_STATUS_ENDED:
		return experiment.StatusEnded
	}
	return experiment.StatusDraft
}

func entStatusToProto(v experiment.Status) *contentV1.Experiment_Status {
	switch v {
	case experiment.StatusRunning:
		return contentV1.Experiment_EXPERIMENT_STATUS_RUNNING.Enum()
	case experiment.StatusPaused:
		return contentV1.Experiment_EXPERIMENT_STATUS_PAUSED.Enum()
	case experiment.StatusEnded:
		return contentV1.Experiment_EXPERIMENT_STATUS_ENDED.Enum()
	}
	return contentV1.Experiment_EXPERIMENT_STATUS_DRAFT.Enum()
}

func assignmentTypeToEntOrDefault(v contentV1.Experiment_AssignmentType) experiment.AssignmentType {
	if v == contentV1.Experiment_ASSIGNMENT_TYPE_RULE {
		return experiment.AssignmentTypeRule
	}
	return experiment.AssignmentTypeRandom
}

func entAssignmentToProto(v experiment.AssignmentType) *contentV1.Experiment_AssignmentType {
	if v == experiment.AssignmentTypeRule {
		return contentV1.Experiment_ASSIGNMENT_TYPE_RULE.Enum()
	}
	return contentV1.Experiment_ASSIGNMENT_TYPE_RANDOM.Enum()
}

// weightToInt32 ent int 权重 → proto *int32。
func weightToInt32(v *int) *int32 {
	if v == nil {
		return nil
	}
	w := int32(*v)
	return &w
}
