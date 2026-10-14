package data

import (
	"context"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/go-kratos/kratos/v2/log"

	entCrud "github.com/tx7do/go-crud/entgo"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/predicate"
	"go-wind-cms/app/core/service/internal/data/ent/socialinboxitem"

	socialV1 "go-wind-cms/api/gen/go/social/service/v1"
)

// SocialInboxRepo 社交收件箱仓储。
// 写入以 (platform, external_id) 去重（先查后插；唯一索引兜底），
// 读取按租户隔离（viewer 过滤由调用方以条件传入）。
type SocialInboxRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper
}

func NewSocialInboxRepo(ctx *bootstrap.Context, entClient *entCrud.EntClient[*ent.Client]) *SocialInboxRepo {
	return &SocialInboxRepo{
		entClient: entClient,
		log:       ctx.NewLoggerHelper("social-inbox/repo/core-service"),
	}
}

// CreateIfAbsent 平台侧条目不存在时才落库（回拉幂等）。返回是否新建。
func (r *SocialInboxRepo) CreateIfAbsent(ctx context.Context, req *socialV1.SocialInboxItem, tenantID uint32) (bool, error) {
	exists, err := r.entClient.Client().SocialInboxItem.Query().
		Where(
			socialinboxitem.PlatformEQ(socialinboxitem.Platform(*req.Platform)),
			socialinboxitem.ExternalIDEQ(req.GetExternalId()),
		).
		Exist(ctx)
	if err != nil {
		r.log.Errorf("inbox exists check failed: %s", err.Error())
		return false, err
	}
	if exists {
		return false, nil
	}

	builder := r.entClient.Client().SocialInboxItem.Create().
		SetTenantID(tenantID).
		SetNillablePlatform(r.platformPtrToEntity(req.Platform)).
		SetNillableAccountID(req.AccountId).
		SetNillableType(req.Type).
		SetNillableExternalID(req.ExternalId).
		SetNillableAuthorName(req.AuthorName).
		SetNillableAuthorHandle(req.AuthorHandle).
		SetNillableContent(req.Content).
		SetNillableURL(req.Url).
		SetNillableOccurredAt(timeTimestamppbToTime(req.OccurredAt)).
		SetCreatedAt(time.Now())

	if err = builder.Exec(ctx); err != nil {
		// 并发回拉撞唯一索引：视为已存在
		r.log.Warnf("inbox create failed (maybe raced): %s", err.Error())
		return false, err
	}
	return true, nil
}

// platformPtrToEntity proto 枚举指针 → ent 枚举指针（nil 安全）。
func (r *SocialInboxRepo) platformPtrToEntity(p *socialV1.Platform) *socialinboxitem.Platform {
	if p == nil {
		return nil
	}
	v := socialinboxitem.Platform(*p)
	return &v
}

// List 收件箱列表（租户隔离；可选账号/未读过滤；按发生时间倒序）。
func (r *SocialInboxRepo) List(ctx context.Context, tenantID uint32, accountID uint32, unreadOnly bool, limit int) ([]*ent.SocialInboxItem, uint64, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}

	var conds []predicate.SocialInboxItem
	conds = append(conds, socialinboxitem.TenantIDEQ(tenantID))
	if accountID != 0 {
		conds = append(conds, socialinboxitem.AccountIDEQ(accountID))
	}
	if unreadOnly {
		conds = append(conds, socialinboxitem.ReadAtIsNil())
	}

	rows, err := r.entClient.Client().SocialInboxItem.Query().
		Where(conds...).
		Order(socialinboxitem.ByOccurredAt(sql.OrderDesc())).
		Limit(limit).
		All(ctx)
	if err != nil {
		r.log.Errorf("inbox list failed: %s", err.Error())
		return nil, 0, err
	}

	unread, err := r.entClient.Client().SocialInboxItem.Query().
		Where(append(conds[:len(conds)-1], socialinboxitem.ReadAtIsNil())...).
		Count(ctx)
	if err != nil {
		r.log.Errorf("inbox unread count failed: %s", err.Error())
	}

	return rows, uint64(unread), nil
}

// MarkRead 标记已读（租户守卫；已读幂等）。
func (r *SocialInboxRepo) MarkRead(ctx context.Context, tenantID uint32, id uint32, now time.Time) error {
	builder := r.entClient.Client().SocialInboxItem.Update().
		Where(socialinboxitem.IDEQ(id)).
		SetReadAt(now).
		SetUpdatedAt(now)
	if tenantID != 0 {
		builder.Where(socialinboxitem.TenantIDEQ(tenantID))
	}
	if _, err := builder.Save(ctx); err != nil {
		r.log.Errorf("inbox mark read failed (id=%d): %s", id, err.Error())
		return err
	}
	return nil
}
