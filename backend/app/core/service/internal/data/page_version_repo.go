package data

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"

	entCrud "github.com/tx7do/go-crud/entgo"

	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/pageversion"
	"go-wind-cms/app/core/service/internal/data/ent/predicate"
	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
)

// PageVersionRepo 页面版本快照仓储。行只增不改，由 PageService 在
// Create/Update 成功后写入完整 DTO 快照（含 translations/sections）；
// 恢复 = 读快照走常规 Update 链路，自然生成新版本。
type PageVersionRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper
}

func NewPageVersionRepo(ctx *bootstrap.Context, entClient *entCrud.EntClient[*ent.Client]) *PageVersionRepo {
	return &PageVersionRepo{
		entClient: entClient,
		log:       ctx.NewLoggerHelper("page-version/repo/core-service"),
	}
}

// pageSnapshotMaxBytes 快照体积上限（约 2MB）：超限拒绝入库，防止异常大页面拖垮版本表。
const pageSnapshotMaxBytes = 2 << 20

// CreateSnapshot 为页面写入一条版本快照，返回新版本号。
func (r *PageVersionRepo) CreateSnapshot(ctx context.Context, pageID uint32, dto *contentV1.Page, comment string) (uint32, error) {
	if dto == nil {
		return 0, nil
	}

	raw, err := protojson.Marshal(dto)
	if err != nil {
		r.log.Errorf("marshal page snapshot failed (page=%d): %s", pageID, err.Error())
		return 0, err
	}
	if len(raw) > pageSnapshotMaxBytes {
		r.log.Errorf("page snapshot too large (page=%d, %d bytes), skip versioning", pageID, len(raw))
		return 0, nil
	}

	next, err := r.nextVersionNumber(ctx, pageID)
	if err != nil {
		return 0, err
	}

	builder := r.entClient.Client().PageVersion.Create().
		SetNillablePageID(&pageID).
		SetNillableVersionNumber(&next).
		SetContent(string(raw))
	if comment != "" {
		builder.SetComment(comment)
	}
	callerUserID, hasUser := viewerUserIDFromContext(ctx)
	if hasUser {
		builder.SetCreatedBy(callerUserID)
	}

	if err = builder.Exec(ctx); err != nil {
		r.log.Errorf("create page version failed (page=%d): %s", pageID, err.Error())
		return 0, err
	}
	return next, nil
}

// tenantGuard 租户过滤：viewer 携带租户时仅限本租户版本。
func (r *PageVersionRepo) tenantGuard(ctx context.Context) ([]predicate.PageVersion, bool) {
	tid, hasTenant := maybeTenantFromViewer(ctx)
	if !hasTenant {
		return nil, false
	}
	return []predicate.PageVersion{pageversion.TenantIDEQ(tid)}, true
}

// nextVersionNumber 取该页下一版本号（max+1）。并发竞态下可能出现重号，
// (page_id, version_number) 为非唯一索引，宁可重号也不阻断保存。
func (r *PageVersionRepo) nextVersionNumber(ctx context.Context, pageID uint32) (uint32, error) {
	conds := []predicate.PageVersion{pageversion.PageIDEQ(pageID)}
	if tenantConds, has := r.tenantGuard(ctx); has {
		conds = append(conds, tenantConds...)
	}
	max, err := r.entClient.Client().PageVersion.Query().
		Where(conds...).
		Aggregate(ent.Max(pageversion.FieldVersionNumber)).
		Int(ctx)
	if err != nil {
		r.log.Errorf("query max version number failed (page=%d): %s", pageID, err.Error())
		return 0, err
	}
	if max < 0 {
		max = 0
	}
	return uint32(max) + 1, nil
}

// ListByPage 取页面的版本列表（版本号倒序，仅元数据，不含快照正文）。
func (r *PageVersionRepo) ListByPage(ctx context.Context, pageID uint32) ([]*contentV1.PageVersion, error) {
	conds := []predicate.PageVersion{pageversion.PageIDEQ(pageID)}
	if tenantConds, has := r.tenantGuard(ctx); has {
		conds = append(conds, tenantConds...)
	}
	rows, err := r.entClient.Client().PageVersion.Query().
		Where(conds...).
		Order(ent.Desc(pageversion.FieldVersionNumber, pageversion.FieldID)).
		All(ctx)
	if err != nil {
		r.log.Errorf("list page versions failed (page=%d): %s", pageID, err.Error())
		return nil, err
	}

	items := make([]*contentV1.PageVersion, 0, len(rows))
	for _, row := range rows {
		items = append(items, pageVersionToDTO(row.ID, row.PageID, row.VersionNumber, nil, row.Comment, row.CreatedBy, row.CreatedAt))
	}
	return items, nil
}

// Get 取版本详情（含快照正文）。
func (r *PageVersionRepo) Get(ctx context.Context, id uint32) (*contentV1.PageVersion, error) {
	conds := []predicate.PageVersion{pageversion.IDEQ(id)}
	if tenantConds, has := r.tenantGuard(ctx); has {
		conds = append(conds, tenantConds...)
	}
	row, err := r.entClient.Client().PageVersion.Query().
		Where(conds...).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, contentV1.ErrorFileNotFound("page version not found")
		}
		r.log.Errorf("get page version failed (id=%d): %s", id, err.Error())
		return nil, err
	}
	return pageVersionToDTO(row.ID, row.PageID, row.VersionNumber, row.Content, row.Comment, row.CreatedBy, row.CreatedAt), nil
}

func pageVersionToDTO(id uint32, parentID, versionNumber *uint32, content, comment *string, createdBy *uint32, createdAt *time.Time) *contentV1.PageVersion {
	dto := &contentV1.PageVersion{
		Id:      &id,
		Content: content,
	}
	if parentID != nil {
		dto.PageId = parentID
	}
	if versionNumber != nil {
		dto.VersionNumber = versionNumber
	}
	if comment != nil {
		dto.Comment = comment
	}
	if createdBy != nil {
		dto.CreatedBy = createdBy
	}
	if createdAt != nil {
		dto.CreatedAt = timestamppb.New(*createdAt)
	}
	return dto
}
