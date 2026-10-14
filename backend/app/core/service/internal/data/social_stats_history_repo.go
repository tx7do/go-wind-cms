package data

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"

	entCrud "github.com/tx7do/go-crud/entgo"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/socialstatshistory"
)

// SocialStatsHistoryRepo 社交分发统计历史快照仓储。
// 行只通过 distribution_id 访问（租户过滤由父表分发记录承担），写入只增不改；
// 旧数据由 PruneOlderThan 定期裁剪。
type SocialStatsHistoryRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper
}

func NewSocialStatsHistoryRepo(ctx *bootstrap.Context, entClient *entCrud.EntClient[*ent.Client]) *SocialStatsHistoryRepo {
	return &SocialStatsHistoryRepo{
		entClient: entClient,
		log:       ctx.NewLoggerHelper("social-stats-history/repo/core-service"),
	}
}

// StatsRecord 一次统计快照。
type StatsRecord struct {
	DistributionID uint32
	Likes          uint32
	Replies        uint32
	Reposts        uint32
	RecordedAt     time.Time
}

// Record 追加一行快照（只增不改）。
func (r *SocialStatsHistoryRepo) Record(ctx context.Context, rec StatsRecord) error {
	builder := r.entClient.Client().SocialStatsHistory.Create().
		SetNillableDistributionID(&rec.DistributionID).
		SetNillableLikes(&rec.Likes).
		SetNillableReplies(&rec.Replies).
		SetNillableReposts(&rec.Reposts).
		SetCreatedAt(rec.RecordedAt)
	if err := builder.Exec(ctx); err != nil {
		r.log.Errorf("record stats history failed (distribution=%d): %s", rec.DistributionID, err.Error())
		return err
	}
	return nil
}

// PruneOlderThan 裁剪指定时间之前的快照（限批，防长事务）。
func (r *SocialStatsHistoryRepo) PruneOlderThan(ctx context.Context, before time.Time, limit int) (int, error) {
	ids, err := r.entClient.Client().SocialStatsHistory.Query().
		Where(socialstatshistory.CreatedAtLT(before)).
		Select(socialstatshistory.FieldID).
		IDs(ctx)
	if err != nil {
		r.log.Errorf("prune stats history: list ids failed: %s", err.Error())
		return 0, err
	}
	if len(ids) > limit {
		ids = ids[:limit]
	}
	if len(ids) == 0 {
		return 0, nil
	}
	n, err := r.entClient.Client().SocialStatsHistory.Delete().
		Where(socialstatshistory.IDIn(ids...)).
		Exec(ctx)
	if err != nil {
		r.log.Errorf("prune stats history failed: %s", err.Error())
		return 0, err
	}
	return n, nil
}
