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
	"go-wind-cms/app/core/service/internal/data/ent/postversion"
	"go-wind-cms/app/core/service/internal/data/ent/predicate"
	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
)

// PostVersionRepo 帖子版本快照仓储。行只增不改，由 PostService 在
// Create/Update 成功后写入完整 DTO 快照（不含 password_hash）；
// 恢复 = 读快照走常规 Update 链路，自然生成新版本。
type PostVersionRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper
}

func NewPostVersionRepo(ctx *bootstrap.Context, entClient *entCrud.EntClient[*ent.Client]) *PostVersionRepo {
	return &PostVersionRepo{
		entClient: entClient,
		log:       ctx.NewLoggerHelper("post-version/repo/core-service"),
	}
}

// postSnapshotMaxBytes 快照体积上限（约 2MB）：超限拒绝入库，防止异常大正文拖垮版本表。
const postSnapshotMaxBytes = 2 << 20

// CreateSnapshot 为帖子写入一条版本快照，返回新版本号。
func (r *PostVersionRepo) CreateSnapshot(ctx context.Context, postID uint32, dto *contentV1.Post, comment string) (uint32, error) {
	if dto == nil {
		return 0, nil
	}

	// 快照不落密码哈希（敏感字段，恢复也不应回写）
	dto.PasswordHash = nil

	raw, err := protojson.Marshal(dto)
	if err != nil {
		r.log.Errorf("marshal post snapshot failed (post=%d): %s", postID, err.Error())
		return 0, err
	}
	if len(raw) > postSnapshotMaxBytes {
		r.log.Errorf("post snapshot too large (post=%d, %d bytes), skip versioning", postID, len(raw))
		return 0, nil
	}

	next, err := r.nextVersionNumber(ctx, postID)
	if err != nil {
		return 0, err
	}

	builder := r.entClient.Client().PostVersion.Create().
		SetNillablePostID(&postID).
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
		r.log.Errorf("create post version failed (post=%d): %s", postID, err.Error())
		return 0, err
	}
	return next, nil
}

// tenantGuard 租户过滤：viewer 携带租户时仅限本租户版本。
func (r *PostVersionRepo) tenantGuard(ctx context.Context) ([]predicate.PostVersion, bool) {
	tid, hasTenant := maybeTenantFromViewer(ctx)
	if !hasTenant {
		return nil, false
	}
	return []predicate.PostVersion{postversion.TenantIDEQ(tid)}, true
}

// nextVersionNumber 取该帖下一版本号（max+1）。并发竞态下可能出现重号，
// (post_id, version_number) 为非唯一索引，宁可重号也不阻断保存。
func (r *PostVersionRepo) nextVersionNumber(ctx context.Context, postID uint32) (uint32, error) {
	conds := []predicate.PostVersion{postversion.PostIDEQ(postID)}
	if tenantConds, has := r.tenantGuard(ctx); has {
		conds = append(conds, tenantConds...)
	}
	// 取列后 Go 内求最大：MAX() 在空集合返回 NULL，Aggregate.Int 扫描会报错，
	// 导致每个帖子的第一个版本永远创建失败（空表鸡生蛋问题）。
	nums, err := r.entClient.Client().PostVersion.Query().
		Where(conds...).
		Select(postversion.FieldVersionNumber).
		Ints(ctx)
	if err != nil {
		r.log.Errorf("query version numbers failed (post=%d): %s", postID, err.Error())
		return 0, err
	}
	max := 0
	for _, n := range nums {
		if n > max {
			max = n
		}
	}
	return uint32(max) + 1, nil
}

// ListByPost 取帖子的版本列表（版本号倒序，仅元数据，不含快照正文）。
func (r *PostVersionRepo) ListByPost(ctx context.Context, postID uint32) ([]*contentV1.PostVersion, error) {
	conds := []predicate.PostVersion{postversion.PostIDEQ(postID)}
	if tenantConds, has := r.tenantGuard(ctx); has {
		conds = append(conds, tenantConds...)
	}
	rows, err := r.entClient.Client().PostVersion.Query().
		Where(conds...).
		Order(ent.Desc(postversion.FieldVersionNumber, postversion.FieldID)).
		All(ctx)
	if err != nil {
		r.log.Errorf("list post versions failed (post=%d): %s", postID, err.Error())
		return nil, err
	}

	items := make([]*contentV1.PostVersion, 0, len(rows))
	for _, row := range rows {
		items = append(items, postVersionToDTO(row.ID, row.PostID, row.VersionNumber, nil, row.Comment, row.CreatedBy, row.CreatedAt))
	}
	return items, nil
}

// Get 取版本详情（含快照正文）。
func (r *PostVersionRepo) Get(ctx context.Context, id uint32) (*contentV1.PostVersion, error) {
	conds := []predicate.PostVersion{postversion.IDEQ(id)}
	if tenantConds, has := r.tenantGuard(ctx); has {
		conds = append(conds, tenantConds...)
	}
	row, err := r.entClient.Client().PostVersion.Query().
		Where(conds...).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, contentV1.ErrorFileNotFound("post version not found")
		}
		r.log.Errorf("get post version failed (id=%d): %s", id, err.Error())
		return nil, err
	}
	return postVersionToDTO(row.ID, row.PostID, row.VersionNumber, row.Content, row.Comment, row.CreatedBy, row.CreatedAt), nil
}

func postVersionToDTO(id uint32, parentID, versionNumber *uint32, content, comment *string, createdBy *uint32, createdAt *time.Time) *contentV1.PostVersion {
	dto := &contentV1.PostVersion{
		Id:      &id,
		Content: content,
	}
	if parentID != nil {
		dto.PostId = parentID
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
