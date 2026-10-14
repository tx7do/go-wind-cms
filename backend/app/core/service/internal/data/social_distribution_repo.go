package data

import (
	"context"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/tx7do/go-utils/copierutil"
	"github.com/tx7do/go-utils/mapper"

	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	entCrud "github.com/tx7do/go-crud/entgo"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/predicate"
	"go-wind-cms/app/core/service/internal/data/ent/socialdistribution"

	socialV1 "go-wind-cms/api/gen/go/social/service/v1"
)

// SocialDistributionRepo 分发记录仓储。
// 状态流转（SCHEDULED→PUBLISHING→PUBLISHED/FAILED）由本仓储的专用方法收敛，
// service/worker 不直接散写 status，避免竞态下状态机被绕过。
type SocialDistributionRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper

	mapper *mapper.CopierMapper[socialV1.SocialDistribution, ent.SocialDistribution]

	repository *entCrud.Repository[
		ent.SocialDistributionQuery, ent.SocialDistributionSelect,
		ent.SocialDistributionCreate, ent.SocialDistributionCreateBulk,
		ent.SocialDistributionUpdate, ent.SocialDistributionUpdateOne,
		ent.SocialDistributionDelete,
		predicate.SocialDistribution,
		socialV1.SocialDistribution, ent.SocialDistribution,
	]

	platformConverter *mapper.EnumTypeConverter[socialV1.Platform, socialdistribution.Platform]
	statusConverter   *mapper.EnumTypeConverter[socialV1.SocialDistribution_Status, socialdistribution.Status]
}

func NewSocialDistributionRepo(ctx *bootstrap.Context, entClient *entCrud.EntClient[*ent.Client]) *SocialDistributionRepo {
	repo := &SocialDistributionRepo{
		entClient: entClient,
		log:       ctx.NewLoggerHelper("social-distribution/repo/core-service"),
		mapper:    mapper.NewCopierMapper[socialV1.SocialDistribution, ent.SocialDistribution](),
		platformConverter: mapper.NewEnumTypeConverter[socialV1.Platform, socialdistribution.Platform](
			socialV1.Platform_name, socialV1.Platform_value,
		),
		statusConverter: mapper.NewEnumTypeConverter[socialV1.SocialDistribution_Status, socialdistribution.Status](
			socialV1.SocialDistribution_Status_name, socialV1.SocialDistribution_Status_value,
		),
	}

	repo.init()

	return repo
}

func (r *SocialDistributionRepo) init() {
	r.repository = entCrud.NewRepository[
		ent.SocialDistributionQuery, ent.SocialDistributionSelect,
		ent.SocialDistributionCreate, ent.SocialDistributionCreateBulk,
		ent.SocialDistributionUpdate, ent.SocialDistributionUpdateOne,
		ent.SocialDistributionDelete,
		predicate.SocialDistribution,
		socialV1.SocialDistribution, ent.SocialDistribution,
	](r.mapper)

	r.mapper.AppendConverters(copierutil.NewTimeStringConverterPair())
	r.mapper.AppendConverters(copierutil.NewTimeTimestamppbConverterPair())
	r.mapper.AppendConverters(r.platformConverter.NewConverterPair())
	r.mapper.AppendConverters(r.statusConverter.NewConverterPair())
}

// List 租户过滤的分发记录列表。
func (r *SocialDistributionRepo) List(ctx context.Context, req *paginationV1.PagingRequest) (*socialV1.ListSocialDistributionsResponse, error) {
	if req == nil {
		return nil, socialV1.ErrorBadRequest("invalid parameter")
	}

	builder := r.entClient.Client().SocialDistribution.Query()
	if tid, hasTenant := maybeTenantFromViewer(ctx); hasTenant {
		builder.Where(socialdistribution.TenantIDEQ(tid))
	}

	ret, err := r.repository.ListWithPaging(ctx, builder, builder.Clone(), req)
	if err != nil {
		return nil, err
	}
	if ret == nil {
		return &socialV1.ListSocialDistributionsResponse{Total: 0, Items: nil}, nil
	}

	return &socialV1.ListSocialDistributionsResponse{
		Total: ret.Total,
		Items: ret.Items,
	}, nil
}

// Get 按 ID 取分发记录（管理面，租户过滤）。
func (r *SocialDistributionRepo) Get(ctx context.Context, req *socialV1.GetSocialDistributionRequest) (*socialV1.SocialDistribution, error) {
	if req == nil {
		return nil, socialV1.ErrorBadRequest("invalid parameter")
	}

	builder := r.entClient.Client().SocialDistribution.Query()

	var whereConds []func(*sql.Selector)
	if tid, hasTenant := maybeTenantFromViewer(ctx); hasTenant {
		whereConds = append(whereConds, func(sels *sql.Selector) {
			sels.Where(sql.EQ(socialdistribution.FieldTenantID, tid))
		})
	}

	dto, err := r.repository.Get(ctx, builder, req.GetViewMask(), whereConds...)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, socialV1.ErrorNotFound("social distribution not found")
		}
		r.log.Errorf("query social distribution failed: %s", err.Error())
		return nil, socialV1.ErrorInternalServerError("query social distribution failed")
	}

	return dto, nil
}

// GetByID 内部读取（worker 的 SystemViewer 上下文，跨租户）。
func (r *SocialDistributionRepo) GetByID(ctx context.Context, id uint32) (*socialV1.SocialDistribution, error) {
	entity, err := r.entClient.Client().SocialDistribution.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, socialV1.ErrorNotFound("social distribution not found")
		}
		r.log.Errorf("query social distribution by id failed: %s", err.Error())
		return nil, socialV1.ErrorInternalServerError("query social distribution failed")
	}
	return r.mapper.ToDTO(entity), nil
}

// Create 落库分发记录。状态/排期入队由 service 层决定。
func (r *SocialDistributionRepo) Create(ctx context.Context, req *socialV1.CreateSocialDistributionRequest) (*socialV1.SocialDistribution, error) {
	if req == nil || req.Data == nil {
		return nil, socialV1.ErrorBadRequest("invalid parameter")
	}

	builder := r.entClient.Client().SocialDistribution.Create().
		SetTenantID(socialTenantIDFromContext(ctx)).
		SetNillablePlatform(r.platformConverter.ToEntity(req.Data.Platform)).
		SetNillableAccountID(req.Data.AccountId).
		SetNillablePostID(req.Data.PostId).
		SetNillableOverrideText(req.Data.OverrideText).
		SetNillableStatus(r.statusConverter.ToEntity(req.Data.Status)).
		SetNillablePublishTime(timeTimestamppbToTime(req.Data.PublishTime)).
		SetNillableTitle(req.Data.Title).
		SetNillableCreatedBy(req.Data.CreatedBy).
		SetCreatedAt(time.Now())

	entity, err := builder.Save(ctx)
	if err != nil {
		r.log.Errorf("insert social distribution failed: %s", err.Error())
		return nil, socialV1.ErrorInternalServerError("insert social distribution failed")
	}

	return r.mapper.ToDTO(entity), nil
}

// Update 更新分发记录（租户守卫）。仅允许编辑文案/排期/关联等字段；
// 状态字段走专用流转方法，这里不开放 status 写入。
func (r *SocialDistributionRepo) Update(ctx context.Context, req *socialV1.UpdateSocialDistributionRequest) (*socialV1.SocialDistribution, error) {
	if req == nil || req.Data == nil {
		return nil, socialV1.ErrorBadRequest("invalid parameter")
	}

	// 不存在则创建
	if req.GetAllowMissing() {
		exist, err := r.entClient.Client().SocialDistribution.Query().
			Where(socialdistribution.IDEQ(req.GetId())).
			Exist(ctx)
		if err != nil {
			r.log.Errorf("query social distribution exist failed: %s", err.Error())
			return nil, socialV1.ErrorInternalServerError("query social distribution exist failed")
		}
		if !exist {
			req.Data.CreatedBy = req.Data.UpdatedBy
			req.Data.UpdatedBy = nil
			_, err = r.Create(ctx, &socialV1.CreateSocialDistributionRequest{Data: req.Data})
			return nil, err
		}
	}

	tid, hasTenant := maybeTenantFromViewer(ctx)
	callerUserID, hasUser := viewerUserIDFromContext(ctx)

	builder := r.entClient.Client().SocialDistribution.UpdateOneID(req.GetId())
	builder.Where(socialdistribution.IDEQ(req.GetId()))
	if hasTenant {
		builder.Where(socialdistribution.TenantIDEQ(tid))
	}

	result, err := r.repository.UpdateOne(ctx, builder, req.Data, req.GetUpdateMask(),
		func(dto *socialV1.SocialDistribution) {
			builder.
				SetNillablePlatform(r.platformConverter.ToEntity(req.Data.Platform)).
				SetNillableAccountID(req.Data.AccountId).
				SetNillablePostID(req.Data.PostId).
				SetNillableOverrideText(req.Data.OverrideText).
				SetNillablePublishTime(timeTimestamppbToTime(req.Data.PublishTime)).
				SetNillableTitle(req.Data.Title).
				SetUpdatedAt(time.Now())

			// updated_by 强制由服务端 viewer context 推导，忽略客户端传入值
			if hasUser {
				builder.SetUpdatedBy(callerUserID)
			}
		},
		func(s *sql.Selector) {
			s.Where(sql.EQ(socialdistribution.FieldID, req.GetId()))
		},
	)

	return result, err
}

// Delete 删除分发记录（租户守卫）。已发布记录删除不影响平台侧帖子。
func (r *SocialDistributionRepo) Delete(ctx context.Context, req *socialV1.DeleteSocialDistributionRequest) error {
	tid, hasTenant := maybeTenantFromViewer(ctx)
	delBuilder := r.entClient.Client().SocialDistribution.Delete()
	delBuilder.Where(socialdistribution.IDEQ(req.GetId()))
	if hasTenant {
		delBuilder.Where(socialdistribution.TenantIDEQ(tid))
	}
	_, err := delBuilder.Exec(ctx)
	if err != nil {
		r.log.Errorf("delete social distribution failed: %s", err.Error())
	}
	return err
}

// ── 状态流转专用方法 ──

// tryLockPublishing 状态闸门：仅 DRAFT/SCHEDULED/FAILED 可进入 PUBLISHING。
// 返回 false 表示当前状态不可发布（并发点击/重复发布的守卫）。
func (r *SocialDistributionRepo) TryLockPublishing(ctx context.Context, id uint32) (bool, error) {
	n, err := r.entClient.Client().SocialDistribution.Update().
		Where(
			socialdistribution.IDEQ(id),
			socialdistribution.StatusIn(
				socialdistribution.StatusStatusDraft,
				socialdistribution.StatusStatusScheduled,
				socialdistribution.StatusStatusFailed,
			),
		).
		SetStatus(socialdistribution.StatusStatusPublishing).
		ClearErrorMessage().
		Save(ctx)
	if err != nil {
		r.log.Errorf("lock publishing failed (id=%d): %s", id, err.Error())
		return false, socialV1.ErrorInternalServerError("lock publishing failed")
	}
	return n > 0, nil
}

// markPublishing worker 开工前置标记（从 SCHEDULED 直入时可见；幂等）。
func (r *SocialDistributionRepo) MarkPublishing(ctx context.Context, id uint32) {
	_, _ = r.entClient.Client().SocialDistribution.Update().
		Where(
			socialdistribution.IDEQ(id),
			socialdistribution.StatusEQ(socialdistribution.StatusStatusScheduled),
		).
		SetStatus(socialdistribution.StatusStatusPublishing).
		Save(ctx)
}

// markPublished 发布成功回写。
func (r *SocialDistributionRepo) MarkPublished(ctx context.Context, id uint32, externalID, externalURI string) error {
	builder := r.entClient.Client().SocialDistribution.UpdateOneID(id).
		SetStatus(socialdistribution.StatusStatusPublished).
		SetPublishedAt(time.Now()).
		ClearErrorMessage().
		SetUpdatedAt(time.Now())
	if externalID != "" {
		builder.SetExternalID(externalID)
	}
	if externalURI != "" {
		builder.SetExternalURI(externalURI)
	}
	if _, err := builder.Save(ctx); err != nil {
		r.log.Errorf("mark published failed (id=%d): %s", id, err.Error())
		return socialV1.ErrorInternalServerError("mark published failed")
	}
	return nil
}

// MarkFailed 发布失败回写并累计尝试次数（error_message 覆盖为最近一次原因）。
// 返回累计尝试次数：worker 据此判定 asynq 重试是否已穷尽（次数达到
// 1+MaxRetry 即最终失败，触发站内通知）。
func (r *SocialDistributionRepo) MarkFailed(ctx context.Context, id uint32, message string) (uint32, error) {
	entity, err := r.entClient.Client().SocialDistribution.UpdateOneID(id).
		SetStatus(socialdistribution.StatusStatusFailed).
		SetErrorMessage(message).
		AddPublishAttempts(1).
		SetUpdatedAt(time.Now()).
		Save(ctx)
	if err != nil {
		r.log.Errorf("mark failed failed (id=%d): %s", id, err.Error())
		return 0, socialV1.ErrorInternalServerError("mark failed failed")
	}
	var attempts uint32
	if entity != nil && entity.PublishAttempts != nil {
		attempts = *entity.PublishAttempts
	}
	return attempts, nil
}

// updateStats 回拉统计快照回写。
func (r *SocialDistributionRepo) UpdateStats(ctx context.Context, id uint32, likes, replies, reposts uint32, syncedAt time.Time) error {
	_, err := r.entClient.Client().SocialDistribution.UpdateOneID(id).
		SetLikes(likes).
		SetReplies(replies).
		SetReposts(reposts).
		SetStatsSyncedAt(syncedAt).
		Save(ctx)
	if err != nil {
		r.log.Errorf("update stats failed (id=%d): %s", id, err.Error())
		return socialV1.ErrorInternalServerError("update stats failed")
	}
	return nil
}

// ListStatsSyncCandidates 统计周期回拉用：近 since 天内发布、external_id 就绪、
// 距上次回拉超过 throttle 间隔的已发布分发记录（至多 limit 条）。
// SystemViewer 上下文调用（跨租户扫描）。
func (r *SocialDistributionRepo) ListStatsSyncCandidates(ctx context.Context, since, throttleBefore time.Time, limit int) ([]*socialV1.SocialDistribution, error) {
	entities, err := r.entClient.Client().SocialDistribution.Query().
		Where(
			socialdistribution.StatusEQ(socialdistribution.StatusStatusPublished),
			socialdistribution.ExternalIDNotNil(),
			socialdistribution.ExternalIDNEQ(""),
			socialdistribution.PublishedAtGTE(since),
			socialdistribution.Or(
				socialdistribution.StatsSyncedAtIsNil(),
				socialdistribution.StatsSyncedAtLTE(throttleBefore),
			),
		).
		Limit(limit).
		All(ctx)
	if err != nil {
		r.log.Errorf("list stats sync candidates failed: %s", err.Error())
		return nil, err
	}

	dtos := make([]*socialV1.SocialDistribution, 0, len(entities))
	for _, entity := range entities {
		dtos = append(dtos, r.mapper.ToDTO(entity))
	}
	return dtos, nil
}


// ListDueScheduledIDs 兜底扫描用：到期（publish_time 为空或已过）仍未发布的
// SCHEDULED 分发记录 ID。SystemViewer 上下文调用（跨租户扫描）。
func (r *SocialDistributionRepo) ListDueScheduledIDs(ctx context.Context, now time.Time) ([]uint32, error) {
	ids, err := r.entClient.Client().SocialDistribution.Query().
		Where(
			socialdistribution.StatusEQ(socialdistribution.StatusStatusScheduled),
			socialdistribution.Or(
				socialdistribution.PublishTimeIsNil(),
				socialdistribution.PublishTimeLTE(now),
			),
		).
		Select(socialdistribution.FieldID).
		IDs(ctx)
	if err != nil {
		r.log.Errorf("list due scheduled distributions failed: %s", err.Error())
		return nil, err
	}
	return ids, nil
}


// countUnfinishedByAccount 统计某账号名下未完结（DRAFT/SCHEDULED/PUBLISHING）的分发数，
// 用于删除账号前的引用守卫。
func (r *SocialDistributionRepo) CountUnfinishedByAccount(ctx context.Context, accountID uint32) (int, error) {
	count, err := r.entClient.Client().SocialDistribution.Query().
		Where(
			socialdistribution.AccountIDEQ(accountID),
			socialdistribution.StatusIn(
				socialdistribution.StatusStatusDraft,
				socialdistribution.StatusStatusScheduled,
				socialdistribution.StatusStatusPublishing,
			),
		).
		Count(ctx)
	if err != nil {
		r.log.Errorf("count unfinished distributions failed (account=%d): %s", accountID, err.Error())
		return 0, err
	}
	return count, nil
}

// TouchStatsSyncedAt 仅推进统计回拉时间戳（平台不支持统计时用来节流跳过）。
func (r *SocialDistributionRepo) TouchStatsSyncedAt(ctx context.Context, id uint32, syncedAt time.Time) error {
	_, err := r.entClient.Client().SocialDistribution.UpdateOneID(id).
		SetStatsSyncedAt(syncedAt).
		Save(ctx)
	if err != nil {
		r.log.Errorf("touch stats synced_at failed (id=%d): %s", id, err.Error())
		return socialV1.ErrorInternalServerError("touch stats synced_at failed")
	}
	return nil
}

// timeTimestamppbToTime Timestamp → time.Time（nil 安全）。
func timeTimestamppbToTime(ts *timestamppb.Timestamp) *time.Time {
	if ts == nil || !ts.IsValid() {
		return nil
	}
	t := ts.AsTime()
	return &t
}
