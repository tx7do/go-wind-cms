package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	entCrud "github.com/tx7do/go-crud/entgo"

	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/commentlike"
	"go-wind-cms/app/core/service/internal/data/ent/interactioncounter"
	"go-wind-cms/app/core/service/internal/data/ent/postlike"
	"go-wind-cms/app/core/service/internal/data/ent/postwatch"

	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
	interactionV1 "go-wind-cms/api/gen/go/interaction/service/v1"
)

// InteractionRepo 是点赞/收藏 ledger 与 interaction_counter 计数表的唯一写入方。
//
// 设计要点：
//   - viewer 用户身份由 service 层从鉴权上下文提取后传入，repo 不接受客户端
//     传入的 user_id。
//   - Like/Unlike 在单个 ent.Tx 内同时操作 ledger 表与 interaction_counter 计数
//     表（upsert 计数行），保证原子性（沿用 post_repo 的跨表事务惯例）。
//   - 计数统一存于 interaction_counter 表（target_type/target_id/metric/count），
//     post/comment 表上的散落计数列已移除。
//   - 重复点赞/收藏走幂等：先 Exist 检查，已存在则 no-op，不依赖约束错误判断。
//   - Watch/Unwatch 仅维护 post_watch ledger，不触碰任何计数表（收藏无计数语义）。
//   - 跨租户隔离由 mixin + 查询带 tenant_id 保证。
type InteractionRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper
}

func NewInteractionRepo(ctx *bootstrap.Context, entClient *entCrud.EntClient[*ent.Client]) *InteractionRepo {
	return &InteractionRepo{
		entClient: entClient,
		log:       ctx.NewLoggerHelper("interaction/repo/core-service"),
	}
}

// txn 在事务内执行 fn，沿用 post_repo 的 defer rollback/commit 惯例。
func (r *InteractionRepo) txn(ctx context.Context, fn func(tx *ent.Tx) error) (err error) {
	var tx *ent.Tx
	tx, err = r.entClient.Client().Tx(ctx)
	if err != nil {
		r.log.Errorf("start transaction failed: %s", err.Error())
		return interactionV1.ErrorInternalServerError("start transaction failed")
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				r.log.Errorf("transaction rollback failed: %s", rollbackErr.Error())
			}
			return
		}
		if commitErr := tx.Commit(); commitErr != nil {
			r.log.Errorf("transaction commit failed: %s", commitErr.Error())
			err = interactionV1.ErrorInternalServerError("transaction commit failed")
		}
	}()
	return fn(tx)
}

// Like 点赞 post 或 comment。幂等：已点赞则 no-op。
// 返回操作后的点赞状态与最新计数。
func (r *InteractionRepo) Like(ctx context.Context, viewerUserID uint32, targetType interactionV1.TargetType, targetID uint32) (bool, int32, error) {
	tid, hasTenant := maybeTenantFromViewer(ctx)
	if !hasTenant || viewerUserID == 0 {
		return false, 0, interactionV1.ErrorUnauthorized("viewer identity required")
	}

	switch targetType {
	case interactionV1.TargetType_TARGET_TYPE_POST, interactionV1.TargetType_TARGET_TYPE_COMMENT:
		likeCount, liked, err := r.likeTarget(ctx, tid, viewerUserID, targetID, targetType)
		if err != nil {
			return false, 0, err
		}
		return liked, likeCount, nil

	default:
		return false, 0, interactionV1.ErrorBadRequest("invalid target type")
	}
}

// Unlike 取消点赞。幂等：未点赞则 no-op。
func (r *InteractionRepo) Unlike(ctx context.Context, viewerUserID uint32, targetType interactionV1.TargetType, targetID uint32) (bool, int32, error) {
	tid, hasTenant := maybeTenantFromViewer(ctx)
	if !hasTenant || viewerUserID == 0 {
		return false, 0, interactionV1.ErrorUnauthorized("viewer identity required")
	}

	switch targetType {
	case interactionV1.TargetType_TARGET_TYPE_POST, interactionV1.TargetType_TARGET_TYPE_COMMENT:
		likeCount, liked, err := r.unlikeTarget(ctx, tid, viewerUserID, targetID, targetType)
		if err != nil {
			return false, 0, err
		}
		return liked, likeCount, nil

	default:
		return false, 0, interactionV1.ErrorBadRequest("invalid target type")
	}
}

// readCount 读取 interaction_counter 表中 (tenant, target_type, target_id, metric) 的当前计数。
// 不存在行则返回 0。
func (r *InteractionRepo) readCount(ctx context.Context, tid, targetID uint32, targetType interactionV1.TargetType, metric interactionV1.CounterMetric) (int32, error) {
	entity, err := r.entClient.Client().InteractionCounter.Query().
		Where(
			interactioncounter.TenantIDEQ(tid),
			interactioncounter.TargetTypeEQ(uint8(targetType)),
			interactioncounter.TargetIDEQ(targetID),
			interactioncounter.MetricEQ(uint8(metric)),
		).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return 0, nil
		}
		r.log.Errorf("query interaction_counter count failed: %s", err.Error())
		return 0, interactionV1.ErrorInternalServerError("query interaction_counter count failed")
	}
	if entity.Count == nil {
		return 0, nil
	}
	return int32(*entity.Count), nil
}

// likeTarget 在单个 tx 内：写 ledger（post_like 或 comment_like）+ 在 interaction_counter
// 表中 upsert 计数（LIKE 指标 +1）。幂等：若 ledger 行已存在则 no-op。
func (r *InteractionRepo) likeTarget(ctx context.Context, tid, viewerUserID, targetID uint32, targetType interactionV1.TargetType) (int32, bool, error) {
	exists, err := r.ledgerExists(ctx, tid, viewerUserID, targetID, targetType, interactionV1.CounterMetric_COUNTER_METRIC_LIKE)
	if err != nil {
		return 0, false, err
	}
	if exists {
		// 幂等：已点赞，no-op，但仍返回当前状态
		likeCount, qerr := r.readCount(ctx, tid, targetID, targetType, interactionV1.CounterMetric_COUNTER_METRIC_LIKE)
		if qerr != nil {
			return 0, false, qerr
		}
		return likeCount, true, nil
	}

	txErr := r.txn(ctx, func(tx *ent.Tx) error {
		if err := r.createLedgerRow(ctx, tx, tid, viewerUserID, targetID, targetType, interactionV1.CounterMetric_COUNTER_METRIC_LIKE); err != nil {
			return err
		}
		// 在 interaction_counter 表中 upsert 计数（+1）
		if err := r.adjustCounterRow(ctx, tx, tid, targetID, targetType, interactionV1.CounterMetric_COUNTER_METRIC_LIKE, 1); err != nil {
			return err
		}
		return nil
	})
	if txErr != nil {
		return 0, false, txErr
	}

	likeCount, err := r.readCount(ctx, tid, targetID, targetType, interactionV1.CounterMetric_COUNTER_METRIC_LIKE)
	if err != nil {
		return 0, false, err
	}
	return likeCount, true, nil
}

// unlikeTarget 在单个 tx 内：删 ledger + 在 interaction_counter 表中递减计数。
// 幂等：若 ledger 行不存在则 no-op。
func (r *InteractionRepo) unlikeTarget(ctx context.Context, tid, viewerUserID, targetID uint32, targetType interactionV1.TargetType) (int32, bool, error) {
	exists, err := r.ledgerExists(ctx, tid, viewerUserID, targetID, targetType, interactionV1.CounterMetric_COUNTER_METRIC_LIKE)
	if err != nil {
		return 0, false, err
	}
	if !exists {
		// 幂等：未点赞，no-op
		likeCount, qerr := r.readCount(ctx, tid, targetID, targetType, interactionV1.CounterMetric_COUNTER_METRIC_LIKE)
		if qerr != nil {
			return 0, false, qerr
		}
		return likeCount, false, nil
	}

	txErr := r.txn(ctx, func(tx *ent.Tx) error {
		affected, err := r.deleteLedgerRow(ctx, tx, tid, viewerUserID, targetID, targetType, interactionV1.CounterMetric_COUNTER_METRIC_LIKE)
		if err != nil {
			return err
		}
		if affected == 0 {
			// 并发删除，幂等处理
			return nil
		}
		// 在 interaction_counter 表中递减计数（-1）
		if err := r.adjustCounterRow(ctx, tx, tid, targetID, targetType, interactionV1.CounterMetric_COUNTER_METRIC_LIKE, -1); err != nil {
			return err
		}
		return nil
	})
	if txErr != nil {
		return 0, false, txErr
	}

	likeCount, err := r.readCount(ctx, tid, targetID, targetType, interactionV1.CounterMetric_COUNTER_METRIC_LIKE)
	if err != nil {
		return 0, false, err
	}
	return likeCount, false, nil
}

// ledgerExists 查询 (tenant, viewer, target) 的 ledger 行是否存在。
// metric 决定查哪张 ledger 表：LIKE→post_like/comment_like，WATCH→post_watch。
func (r *InteractionRepo) ledgerExists(ctx context.Context, tid, viewerUserID, targetID uint32, targetType interactionV1.TargetType, metric interactionV1.CounterMetric) (bool, error) {
	switch metric {
	case interactionV1.CounterMetric_COUNTER_METRIC_LIKE:
		switch targetType {
		case interactionV1.TargetType_TARGET_TYPE_POST:
			exists, err := r.entClient.Client().PostLike.Query().
				Where(
					postlike.TenantIDEQ(tid),
					postlike.UserIDEQ(viewerUserID),
					postlike.PostIDEQ(targetID),
				).
				Exist(ctx)
			if err != nil {
				r.log.Errorf("query post_like exist failed: %s", err.Error())
				return false, interactionV1.ErrorInternalServerError("query post_like exist failed")
			}
			return exists, nil

		case interactionV1.TargetType_TARGET_TYPE_COMMENT:
			exists, err := r.entClient.Client().CommentLike.Query().
				Where(
					commentlike.TenantIDEQ(tid),
					commentlike.UserIDEQ(viewerUserID),
					commentlike.CommentIDEQ(targetID),
				).
				Exist(ctx)
			if err != nil {
				r.log.Errorf("query comment_like exist failed: %s", err.Error())
				return false, interactionV1.ErrorInternalServerError("query comment_like exist failed")
			}
			return exists, nil

		default:
			return false, interactionV1.ErrorBadRequest("invalid target type")
		}

	case interactionV1.CounterMetric_COUNTER_METRIC_WATCH:
		// WATCH 仅作用于 post
		exists, err := r.entClient.Client().PostWatch.Query().
			Where(
				postwatch.TenantIDEQ(tid),
				postwatch.UserIDEQ(viewerUserID),
				postwatch.PostIDEQ(targetID),
			).
			Exist(ctx)
		if err != nil {
			r.log.Errorf("query post_watch exist failed: %s", err.Error())
			return false, interactionV1.ErrorInternalServerError("query post_watch exist failed")
		}
		return exists, nil

	default:
		return false, interactionV1.ErrorBadRequest("invalid metric")
	}
}

// createLedgerRow 在 tx 内创建 ledger 行。
// metric 决定写哪张 ledger 表。
func (r *InteractionRepo) createLedgerRow(ctx context.Context, tx *ent.Tx, tid, viewerUserID, targetID uint32, targetType interactionV1.TargetType, metric interactionV1.CounterMetric) error {
	switch metric {
	case interactionV1.CounterMetric_COUNTER_METRIC_LIKE:
		switch targetType {
		case interactionV1.TargetType_TARGET_TYPE_POST:
			if _, err := tx.PostLike.Create().
				SetTenantID(tid).
				SetUserID(viewerUserID).
				SetPostID(targetID).
				Save(ctx); err != nil {
				r.log.Errorf("insert post_like failed: %s", err.Error())
				return interactionV1.ErrorInternalServerError("insert post_like failed")
			}
			return nil

		case interactionV1.TargetType_TARGET_TYPE_COMMENT:
			if _, err := tx.CommentLike.Create().
				SetTenantID(tid).
				SetUserID(viewerUserID).
				SetCommentID(targetID).
				Save(ctx); err != nil {
				r.log.Errorf("insert comment_like failed: %s", err.Error())
				return interactionV1.ErrorInternalServerError("insert comment_like failed")
			}
			return nil

		default:
			return interactionV1.ErrorBadRequest("invalid target type")
		}

	case interactionV1.CounterMetric_COUNTER_METRIC_WATCH:
		if _, err := tx.PostWatch.Create().
			SetTenantID(tid).
			SetUserID(viewerUserID).
			SetPostID(targetID).
			Save(ctx); err != nil {
			r.log.Errorf("insert post_watch failed: %s", err.Error())
			return interactionV1.ErrorInternalServerError("insert post_watch failed")
		}
		return nil

	default:
		return interactionV1.ErrorBadRequest("invalid metric")
	}
}

// deleteLedgerRow 在 tx 内删除 ledger 行，返回受影响行数。
// metric 决定删哪张 ledger 表。
func (r *InteractionRepo) deleteLedgerRow(ctx context.Context, tx *ent.Tx, tid, viewerUserID, targetID uint32, targetType interactionV1.TargetType, metric interactionV1.CounterMetric) (int, error) {
	switch metric {
	case interactionV1.CounterMetric_COUNTER_METRIC_LIKE:
		switch targetType {
		case interactionV1.TargetType_TARGET_TYPE_POST:
			affected, err := tx.PostLike.Delete().
				Where(
					postlike.TenantIDEQ(tid),
					postlike.UserIDEQ(viewerUserID),
					postlike.PostIDEQ(targetID),
				).
				Exec(ctx)
			if err != nil {
				r.log.Errorf("delete post_like failed: %s", err.Error())
				return 0, interactionV1.ErrorInternalServerError("delete post_like failed")
			}
			return affected, nil

		case interactionV1.TargetType_TARGET_TYPE_COMMENT:
			affected, err := tx.CommentLike.Delete().
				Where(
					commentlike.TenantIDEQ(tid),
					commentlike.UserIDEQ(viewerUserID),
					commentlike.CommentIDEQ(targetID),
				).
				Exec(ctx)
			if err != nil {
				r.log.Errorf("delete comment_like failed: %s", err.Error())
				return 0, interactionV1.ErrorInternalServerError("delete comment_like failed")
			}
			return affected, nil

		default:
			return 0, interactionV1.ErrorBadRequest("invalid target type")
		}

	case interactionV1.CounterMetric_COUNTER_METRIC_WATCH:
		affected, err := tx.PostWatch.Delete().
			Where(
				postwatch.TenantIDEQ(tid),
				postwatch.UserIDEQ(viewerUserID),
				postwatch.PostIDEQ(targetID),
			).
			Exec(ctx)
		if err != nil {
			r.log.Errorf("delete post_watch failed: %s", err.Error())
			return 0, interactionV1.ErrorInternalServerError("delete post_watch failed")
		}
		return affected, nil

	default:
		return 0, interactionV1.ErrorBadRequest("invalid metric")
	}
}

// adjustCounterRow 在 tx 内 upsert interaction_counter 表的计数行。
//
// 行不存在时按 delta 创建（仅当 delta>0）；存在时按 delta 递增，归 0 则删行。
// 复合 unique 索引 (tenant, target_type, target_id, metric) 保证单目标单 metric 仅一行。
func (r *InteractionRepo) adjustCounterRow(ctx context.Context, tx *ent.Tx, tid, targetID uint32, targetType interactionV1.TargetType, metric interactionV1.CounterMetric, delta int64) error {
	existing, err := tx.InteractionCounter.Query().
		Where(
			interactioncounter.TenantIDEQ(tid),
			interactioncounter.TargetTypeEQ(uint8(targetType)),
			interactioncounter.TargetIDEQ(targetID),
			interactioncounter.MetricEQ(uint8(metric)),
		).
		Only(ctx)
	if err != nil {
		if !ent.IsNotFound(err) {
			r.log.Errorf("query interaction_counter row failed: %s", err.Error())
			return interactionV1.ErrorInternalServerError("query interaction_counter row failed")
		}
		// 行不存在
		if delta <= 0 {
			// 无行可递减，幂等 no-op
			return nil
		}
		if _, cerr := tx.InteractionCounter.Create().
			SetTenantID(tid).
			SetTargetType(uint8(targetType)).
			SetTargetID(targetID).
			SetMetric(uint8(metric)).
			SetCount(delta).
			Save(ctx); cerr != nil {
			r.log.Errorf("create interaction_counter row failed: %s", cerr.Error())
			return interactionV1.ErrorInternalServerError("create interaction_counter row failed")
		}
		return nil
	}

	// 行存在，递增
	newVal := *existing.Count + delta
	if newVal <= 0 {
		// 计数归 0，删行
		if derr := tx.InteractionCounter.DeleteOneID(existing.ID).Exec(ctx); derr != nil {
			r.log.Errorf("delete interaction_counter row failed: %s", derr.Error())
			return interactionV1.ErrorInternalServerError("delete interaction_counter row failed")
		}
		return nil
	}
	if _, uerr := tx.InteractionCounter.UpdateOneID(existing.ID).SetCount(newVal).Save(ctx); uerr != nil {
		r.log.Errorf("update interaction_counter row failed: %s", uerr.Error())
		return interactionV1.ErrorInternalServerError("update interaction_counter row failed")
	}
	return nil
}

// GetCounts 批量查询 interaction_counter 表中指定 (target_type, target_ids, metrics) 的计数。
// 返回 target_id → metric → count 的嵌套 map。未记录的 (target, metric) 不出现在响应中。
func (r *InteractionRepo) GetCounts(ctx context.Context, targetType interactionV1.TargetType, targetIDs []uint32, metrics []interactionV1.CounterMetric) (map[uint32]*interactionV1.CountMap, error) {
	tid, hasTenant := maybeTenantFromViewer(ctx)
	if !hasTenant {
		return nil, interactionV1.ErrorUnauthorized("viewer identity required")
	}
	if len(targetIDs) == 0 || len(metrics) == 0 {
		return map[uint32]*interactionV1.CountMap{}, nil
	}

	// 转换 metrics 为 uint8 谓词列表
	metricPreds := make([]uint8, 0, len(metrics))
	for _, m := range metrics {
		if m == interactionV1.CounterMetric_COUNTER_METRIC_UNSPECIFIED {
			continue
		}
		metricPreds = append(metricPreds, uint8(m))
	}
	if len(metricPreds) == 0 {
		return map[uint32]*interactionV1.CountMap{}, nil
	}

	rows, err := r.entClient.Client().InteractionCounter.Query().
		Where(
			interactioncounter.TenantIDEQ(tid),
			interactioncounter.TargetTypeEQ(uint8(targetType)),
			interactioncounter.TargetIDIn(targetIDs...),
			interactioncounter.MetricIn(metricPreds...),
		).
		All(ctx)
	if err != nil {
		r.log.Errorf("query interaction_counter batch failed: %s", err.Error())
		return nil, interactionV1.ErrorInternalServerError("query interaction_counter batch failed")
	}

	result := make(map[uint32]*interactionV1.CountMap, len(rows))
	for _, row := range rows {
		if row.TargetID == nil || row.Metric == nil || row.Count == nil {
			continue
		}
		tidVal := *row.TargetID
		metricVal := interactionV1.CounterMetric(*row.Metric)
		countVal := *row.Count
		cm, ok := result[tidVal]
		if !ok {
			cm = &interactionV1.CountMap{}
			result[tidVal] = cm
		}
		cm.Counts = append(cm.Counts, &interactionV1.MetricCount{
			Metric: metricVal,
			Count:  countVal,
		})
	}
	return result, nil
}

// Watch 收藏 post。幂等：已收藏则 no-op。
// 在单个 tx 内写 post_watch ledger + 在 interaction_counter 表中 upsert 计数（WATCH 指标 +1）。
// 返回操作后的收藏状态与最新收藏计数。
func (r *InteractionRepo) Watch(ctx context.Context, viewerUserID uint32, postID uint32) (bool, int32, error) {
	tid, hasTenant := maybeTenantFromViewer(ctx)
	if !hasTenant || viewerUserID == 0 {
		return false, 0, interactionV1.ErrorUnauthorized("viewer identity required")
	}

	exists, err := r.ledgerExists(ctx, tid, viewerUserID, postID, interactionV1.TargetType_TARGET_TYPE_POST, interactionV1.CounterMetric_COUNTER_METRIC_WATCH)
	if err != nil {
		return false, 0, err
	}
	if exists {
		// 幂等：已收藏，no-op，但仍返回当前状态与计数
		watchCount, qerr := r.readCount(ctx, tid, postID, interactionV1.TargetType_TARGET_TYPE_POST, interactionV1.CounterMetric_COUNTER_METRIC_WATCH)
		if qerr != nil {
			return false, 0, qerr
		}
		return true, watchCount, nil
	}

	txErr := r.txn(ctx, func(tx *ent.Tx) error {
		if err := r.createLedgerRow(ctx, tx, tid, viewerUserID, postID, interactionV1.TargetType_TARGET_TYPE_POST, interactionV1.CounterMetric_COUNTER_METRIC_WATCH); err != nil {
			return err
		}
		// 在 interaction_counter 表中 upsert 计数（WATCH +1）
		if err := r.adjustCounterRow(ctx, tx, tid, postID, interactionV1.TargetType_TARGET_TYPE_POST, interactionV1.CounterMetric_COUNTER_METRIC_WATCH, 1); err != nil {
			return err
		}
		return nil
	})
	if txErr != nil {
		return false, 0, txErr
	}

	watchCount, err := r.readCount(ctx, tid, postID, interactionV1.TargetType_TARGET_TYPE_POST, interactionV1.CounterMetric_COUNTER_METRIC_WATCH)
	if err != nil {
		return false, 0, err
	}
	return true, watchCount, nil
}

// Unwatch 取消收藏 post。幂等。
// 在单个 tx 内删 post_watch ledger + 在 interaction_counter 表中递减计数。
func (r *InteractionRepo) Unwatch(ctx context.Context, viewerUserID uint32, postID uint32) (bool, int32, error) {
	tid, hasTenant := maybeTenantFromViewer(ctx)
	if !hasTenant || viewerUserID == 0 {
		return false, 0, interactionV1.ErrorUnauthorized("viewer identity required")
	}

	exists, err := r.ledgerExists(ctx, tid, viewerUserID, postID, interactionV1.TargetType_TARGET_TYPE_POST, interactionV1.CounterMetric_COUNTER_METRIC_WATCH)
	if err != nil {
		return false, 0, err
	}
	if !exists {
		// 幂等：未收藏，no-op
		watchCount, qerr := r.readCount(ctx, tid, postID, interactionV1.TargetType_TARGET_TYPE_POST, interactionV1.CounterMetric_COUNTER_METRIC_WATCH)
		if qerr != nil {
			return false, 0, qerr
		}
		return false, watchCount, nil
	}

	txErr := r.txn(ctx, func(tx *ent.Tx) error {
		affected, err := r.deleteLedgerRow(ctx, tx, tid, viewerUserID, postID, interactionV1.TargetType_TARGET_TYPE_POST, interactionV1.CounterMetric_COUNTER_METRIC_WATCH)
		if err != nil {
			return err
		}
		if affected == 0 {
			// 并发删除，幂等处理
			return nil
		}
		// 在 interaction_counter 表中递减计数（WATCH -1）
		if err := r.adjustCounterRow(ctx, tx, tid, postID, interactionV1.TargetType_TARGET_TYPE_POST, interactionV1.CounterMetric_COUNTER_METRIC_WATCH, -1); err != nil {
			return err
		}
		return nil
	})
	if txErr != nil {
		return false, 0, txErr
	}

	watchCount, err := r.readCount(ctx, tid, postID, interactionV1.TargetType_TARGET_TYPE_POST, interactionV1.CounterMetric_COUNTER_METRIC_WATCH)
	if err != nil {
		return false, 0, err
	}
	return false, watchCount, nil
}

// GetInteractionStatus 批量查询当前 viewer 对指定目标的 {liked, watched} 状态。
// comment 目标的 watched 永远为 false（收藏仅对 post）。
func (r *InteractionRepo) GetInteractionStatus(ctx context.Context, viewerUserID uint32, targetType interactionV1.TargetType, targetIDs []uint32) (map[uint32]*interactionV1.InteractionStatus, error) {
	tid, hasTenant := maybeTenantFromViewer(ctx)
	if !hasTenant || viewerUserID == 0 {
		return nil, interactionV1.ErrorUnauthorized("viewer identity required")
	}

	result := make(map[uint32]*interactionV1.InteractionStatus, len(targetIDs))
	for _, id := range targetIDs {
		result[id] = &interactionV1.InteractionStatus{Liked: false, Watched: false}
	}

	switch targetType {
	case interactionV1.TargetType_TARGET_TYPE_POST:
		// 查 post_like
		likedRows, err := r.entClient.Client().PostLike.Query().
			Where(
				postlike.TenantIDEQ(tid),
				postlike.UserIDEQ(viewerUserID),
				postlike.PostIDIn(targetIDs...),
			).
			All(ctx)
		if err != nil {
			r.log.Errorf("query post_like status failed: %s", err.Error())
			return nil, interactionV1.ErrorInternalServerError("query post_like status failed")
		}
		for _, row := range likedRows {
			if row.PostID != nil {
				if _, ok := result[*row.PostID]; ok {
					result[*row.PostID].Liked = true
				}
			}
		}

		// 查 post_watch
		watchedRows, err := r.entClient.Client().PostWatch.Query().
			Where(
				postwatch.TenantIDEQ(tid),
				postwatch.UserIDEQ(viewerUserID),
				postwatch.PostIDIn(targetIDs...),
			).
			All(ctx)
		if err != nil {
			r.log.Errorf("query post_watch status failed: %s", err.Error())
			return nil, interactionV1.ErrorInternalServerError("query post_watch status failed")
		}
		for _, row := range watchedRows {
			if row.PostID != nil {
				if _, ok := result[*row.PostID]; ok {
					result[*row.PostID].Watched = true
				}
			}
		}

	case interactionV1.TargetType_TARGET_TYPE_COMMENT:
		likedRows, err := r.entClient.Client().CommentLike.Query().
			Where(
				commentlike.TenantIDEQ(tid),
				commentlike.UserIDEQ(viewerUserID),
				commentlike.CommentIDIn(targetIDs...),
			).
			All(ctx)
		if err != nil {
			r.log.Errorf("query comment_like status failed: %s", err.Error())
			return nil, interactionV1.ErrorInternalServerError("query comment_like status failed")
		}
		for _, row := range likedRows {
			if row.CommentID != nil {
				if _, ok := result[*row.CommentID]; ok {
					result[*row.CommentID].Liked = true
				}
			}
		}
		// watched 对 comment 永远 false，已在初始化时设好

	default:
		return nil, interactionV1.ErrorBadRequest("invalid target type")
	}

	return result, nil
}

// ListWatchedPosts 列出当前 viewer 收藏的 post。
// 查 post_watch 拿到分页后的 post_ids，再逐个调 PostRepo.Get 复用其完整的
// 附带查询（translations/tags/categories/view_mask）逻辑。
func (r *InteractionRepo) ListWatchedPosts(ctx context.Context, viewerUserID uint32, postRepo *PostRepo, req *paginationV1.PagingRequest) (*contentV1.ListPostResponse, error) {
	tid, hasTenant := maybeTenantFromViewer(ctx)
	if !hasTenant || viewerUserID == 0 {
		return nil, interactionV1.ErrorUnauthorized("viewer identity required")
	}

	// 注：post_watch 的 (tenant_id, user_id, post_id) unique 索引兼作查询索引，
	// 此处按 (tenant, user) 过滤，分页由 PostWatchQuery 的 Limit/Offset 处理。
	q := r.entClient.Client().PostWatch.Query().
		Where(
			postwatch.TenantIDEQ(tid),
			postwatch.UserIDEQ(viewerUserID),
		)

	// 应用分页（Page 为 1-based）
	page := req.GetPage()
	if page == 0 {
		page = 1
	}
	limit := int(req.GetPageSize())
	if limit <= 0 {
		limit = 20
	}
	offset := int(page-1) * limit
	q = q.Limit(limit).Offset(offset)

	watchRows, err := q.All(ctx)
	if err != nil {
		r.log.Errorf("query post_watch list failed: %s", err.Error())
		return nil, interactionV1.ErrorInternalServerError("query post_watch list failed")
	}

	items := make([]*contentV1.Post, 0, len(watchRows))
	for _, row := range watchRows {
		if row.PostID == nil {
			continue
		}
		// 复用 PostRepo.Get 的附带查询。收藏列表请求不携带 locale/view_mask，
		// 故此处 GetPostRequest 不设置二者（默认全字段、默认语言）。
		postEntity, gerr := postRepo.Get(ctx, &contentV1.GetPostRequest{
			QueryBy: &contentV1.GetPostRequest_Id{
				Id: *row.PostID,
			},
		})
		if gerr != nil {
			r.log.Errorf("fetch watched post %d failed: %s", *row.PostID, gerr.Error())
			continue
		}
		items = append(items, postEntity)
	}

	// total 为该 viewer 收藏总数（用于前端分页元数据）
	total, cerr := r.entClient.Client().PostWatch.Query().
		Where(
			postwatch.TenantIDEQ(tid),
			postwatch.UserIDEQ(viewerUserID),
		).
		Count(ctx)
	if cerr != nil {
		r.log.Errorf("count post_watch failed: %s", cerr.Error())
		return nil, interactionV1.ErrorInternalServerError("count post_watch failed")
	}

	return &contentV1.ListPostResponse{
		Items: items,
		Total: uint64(total),
	}, nil
}
