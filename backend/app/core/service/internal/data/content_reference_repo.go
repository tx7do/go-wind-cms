package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	entCrud "github.com/tx7do/go-crud/entgo"

	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/contentreference"
)

// RefTarget 引用目标（一条边指向的资源 + 所在字段）。
type RefTarget struct {
	TargetType string // media_asset
	TargetID   uint32
	Field      string // body / thumbnail
}

// ContentReferenceRepo 内容引用关系仓储（引用图谱的边）。
type ContentReferenceRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper
}

func NewContentReferenceRepo(ctx *bootstrap.Context, entClient *entCrud.EntClient[*ent.Client]) *ContentReferenceRepo {
	return &ContentReferenceRepo{
		entClient: entClient,
		log:       ctx.NewLoggerHelper("content-reference/repo/core-service"),
	}
}

// ReplaceSourceReferences 重建某源内容的全部引用行：先删旧行（含跨租户残留，
// 以源定位为准），再批量插入新行（(源,字段,目标) 唯一，重复目标自动去重）。
// 返回受影响的目标 ID 集合（新旧行并集，供调用方重算目标侧计数）。
func (r *ContentReferenceRepo) ReplaceSourceReferences(ctx context.Context, sourceType string, sourceID uint32, targets []RefTarget) (map[uint32]struct{}, error) {
	client := r.entClient.Client()
	affected := make(map[uint32]struct{})

	// 旧行目标并入受影响集合
	oldRows, err := client.ContentReference.Query().
		Where(
			contentreference.SourceTypeEQ(sourceType),
			contentreference.SourceIDEQ(sourceID),
		).
		All(ctx)
	if err != nil {
		r.log.Errorf("query old references failed: %s", err.Error())
		return nil, err
	}
	for _, row := range oldRows {
		affected[row.TargetID] = struct{}{}
	}

	if _, err = client.ContentReference.Delete().
		Where(
			contentreference.SourceTypeEQ(sourceType),
			contentreference.SourceIDEQ(sourceID),
		).
		Exec(ctx); err != nil {
		r.log.Errorf("delete old references failed: %s", err.Error())
		return nil, err
	}

	if len(targets) == 0 {
		return affected, nil
	}

	// 同一 (field, target) 去重后批量插入
	seen := make(map[RefTarget]struct{}, len(targets))
	bulk := make([]*ent.ContentReferenceCreate, 0, len(targets))
	for _, t := range targets {
		if t.TargetID == 0 || t.TargetType == "" {
			continue
		}
		key := t
		if _, dup := seen[key]; dup {
			continue
		}
		seen[key] = struct{}{}
		affected[t.TargetID] = struct{}{}
		bulk = append(bulk, client.ContentReference.Create().
			SetSourceType(sourceType).
			SetSourceID(sourceID).
			SetTargetType(t.TargetType).
			SetTargetID(t.TargetID).
			SetField(t.Field),
		)
	}
	if len(bulk) == 0 {
		return affected, nil
	}
	if _, err = client.ContentReference.CreateBulk(bulk...).Save(ctx); err != nil {
		r.log.Errorf("bulk create references failed: %s", err.Error())
		return nil, err
	}
	return affected, nil
}

// ListByTarget 查目标资源（如某媒体）被引用的全部行。
func (r *ContentReferenceRepo) ListByTarget(ctx context.Context, targetType string, targetID uint32) ([]*ent.ContentReference, error) {
	rows, err := r.entClient.Client().ContentReference.Query().
		Where(
			contentreference.TargetTypeEQ(targetType),
			contentreference.TargetIDEQ(targetID),
		).
		Order(ent.Desc(contentreference.FieldUpdatedAt)).
		All(ctx)
	if err != nil {
		r.log.Errorf("list references by target failed: %s", err.Error())
		return nil, err
	}
	return rows, nil
}

// CountByTargetIDs 批量统计各目标被引用次数。
// 受影响目标数有限（单次保存涉及数个），逐 ID 计数即可，避免 GroupBy 多列 Scan 的映射歧义。
func (r *ContentReferenceRepo) CountByTargetIDs(ctx context.Context, targetType string, targetIDs []uint32) (map[uint32]int, error) {
	counts := make(map[uint32]int, len(targetIDs))
	for _, id := range targetIDs {
		n, cErr := r.entClient.Client().ContentReference.Query().
			Where(
				contentreference.TargetTypeEQ(targetType),
				contentreference.TargetIDEQ(id),
			).
			Count(ctx)
		if cErr != nil {
			r.log.Errorf("count references for target [%d] failed: %s", id, cErr.Error())
			return counts, cErr
		}
		counts[id] = n
	}
	return counts, nil
}
