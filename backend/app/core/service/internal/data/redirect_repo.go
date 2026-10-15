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
	"google.golang.org/protobuf/types/known/timestamppb"

	entCrud "github.com/tx7do/go-crud/entgo"

	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/predicate"
	"go-wind-cms/app/core/service/internal/data/ent/redirect"
	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
)

// RedirectRepo URL 重定向仓储。
type RedirectRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper
}

func NewRedirectRepo(ctx *bootstrap.Context, entClient *entCrud.EntClient[*ent.Client]) *RedirectRepo {
	return &RedirectRepo{
		entClient: entClient,
		log:       ctx.NewLoggerHelper("redirect/repo/core-service"),
	}
}

// List 重定向列表（时间倒序；支持 from_path 模糊 / enabled 过滤）。
func (r *RedirectRepo) List(ctx context.Context, req *paginationV1.PagingRequest) ([]*contentV1.Redirect, int, error) {
	conds := r.tenantConds(ctx)
	if req != nil {
		if filterExpr, fErr := paginationFilter.ConvertFilterByPagingRequest(req); fErr == nil && filterExpr != nil {
			for _, cond := range filterExpr.GetConditions() {
				switch cond.GetField() {
				case "fromPath", "from_path":
					if v := cond.GetValue(); v != "" {
						conds = append(conds, redirect.FromPathContains(v))
					}
				case "enabled":
					conds = append(conds, redirect.EnabledEQ(cond.GetValue() == "true"))
				}
			}
		}
	}

	total, err := r.entClient.Client().Redirect.Query().Where(conds...).Count(ctx)
	if err != nil {
		r.log.Errorf("count redirects failed: %s", err.Error())
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

	rows, err := r.entClient.Client().Redirect.Query().Where(conds...).
		Order(ent.Desc(redirect.FieldID)).
		Limit(int(pageSize)).
		Offset(int((page - 1) * pageSize)).
		All(ctx)
	if err != nil {
		r.log.Errorf("list redirects failed: %s", err.Error())
		return nil, 0, err
	}

	items := make([]*contentV1.Redirect, 0, len(rows))
	for _, row := range rows {
		items = append(items, r.toDTO(row))
	}
	return items, total, nil
}

// ListEnabled 全量启用项（公开查询，C 端路由中间件消费）。
func (r *RedirectRepo) ListEnabled(ctx context.Context) ([]*contentV1.Redirect, error) {
	conds := append(r.tenantConds(ctx), redirect.EnabledEQ(true))
	rows, err := r.entClient.Client().Redirect.Query().Where(conds...).
		Order(ent.Desc(redirect.FieldID)).
		All(ctx)
	if err != nil {
		r.log.Errorf("list enabled redirects failed: %s", err.Error())
		return nil, err
	}
	items := make([]*contentV1.Redirect, 0, len(rows))
	for _, row := range rows {
		items = append(items, r.toDTO(row))
	}
	return items, nil
}

// Get 重定向详情。
func (r *RedirectRepo) Get(ctx context.Context, id uint32) (*contentV1.Redirect, error) {
	conds := append(r.tenantConds(ctx), redirect.IDEQ(id))
	row, err := r.entClient.Client().Redirect.Query().Where(conds...).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.NotFound("REDIRECT", "redirect not found")
		}
		r.log.Errorf("get redirect failed (id=%d): %s", id, err.Error())
		return nil, err
	}
	return r.toDTO(row), nil
}

// Create 新建（同租户 from_path 查重）。
func (r *RedirectRepo) Create(ctx context.Context, dto *contentV1.Redirect) (*contentV1.Redirect, error) {
	dup, err := r.entClient.Client().Redirect.Query().Where(append(
		r.tenantConds(ctx),
		redirect.FromPathEQ(dto.GetFromPath()),
	)...).Count(ctx)
	if err != nil {
		r.log.Errorf("check redirect duplicate failed: %s", err.Error())
		return nil, err
	}
	if dup > 0 {
		return nil, errors.Conflict("REDIRECT", "from_path already exists")
	}

	builder := r.entClient.Client().Redirect.Create().
		SetFromPath(dto.GetFromPath()).
		SetToURL(dto.GetToUrl()).
		SetNillableStatusCode(dto.StatusCode).
		SetEnabled(dto.GetEnabled())
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
		r.log.Errorf("create redirect failed: %s", err.Error())
		return nil, err
	}
	return r.toDTO(row), nil
}

// Update 按 mask 更新。
func (r *RedirectRepo) Update(ctx context.Context, id uint32, dto *contentV1.Redirect) (*contentV1.Redirect, error) {
	conds := append(r.tenantConds(ctx), redirect.IDEQ(id))
	exist, err := r.entClient.Client().Redirect.Query().Where(conds...).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.NotFound("REDIRECT", "redirect not found")
		}
		return nil, err
	}

	builder := r.entClient.Client().Redirect.UpdateOneID(exist.ID)
	if dto.FromPath != nil {
		builder.SetFromPath(dto.GetFromPath())
	}
	if dto.ToUrl != nil {
		builder.SetToURL(dto.GetToUrl())
	}
	if dto.StatusCode != nil {
		builder.SetNillableStatusCode(dto.StatusCode)
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
		r.log.Errorf("update redirect failed (id=%d): %s", id, err.Error())
		return nil, err
	}
	return r.toDTO(row), nil
}

// Delete 删除。
func (r *RedirectRepo) Delete(ctx context.Context, id uint32) error {
	conds := append(r.tenantConds(ctx), redirect.IDEQ(id))
	count, err := r.entClient.Client().Redirect.Delete().Where(conds...).Exec(ctx)
	if err != nil {
		r.log.Errorf("delete redirect failed (id=%d): %s", id, err.Error())
		return err
	}
	if count == 0 {
		return errors.NotFound("REDIRECT", "redirect not found")
	}
	return nil
}

func (r *RedirectRepo) tenantConds(ctx context.Context) []predicate.Redirect {
	tid, hasTenant := maybeTenantFromViewer(ctx)
	if !hasTenant {
		return nil
	}
	return []predicate.Redirect{redirect.TenantIDEQ(tid)}
}

// tenantIDFromInternalMetadata 从内网可信元数据 x-md-global-tenant-id 读租户。
func (r *RedirectRepo) tenantIDFromInternalMetadata(ctx context.Context) uint32 {
	if md, ok := kratosMetadata.FromServerContext(ctx); ok {
		if v := md.Get("x-md-global-tenant-id"); v != "" {
			if tid, err := strconv.ParseUint(v, 10, 32); err == nil && tid > 0 {
				return uint32(tid)
			}
		}
	}
	return 0
}

func (r *RedirectRepo) toDTO(row *ent.Redirect) *contentV1.Redirect {
	enabled := true
	if row.Enabled != nil {
		enabled = *row.Enabled
	}
	dto := &contentV1.Redirect{
		Id:          &row.ID,
		FromPath:    &row.FromPath,
		ToUrl:       &row.ToURL,
		Enabled:     &enabled,
		Description: row.Description,
		CreatedBy:   row.CreatedBy,
		UpdatedBy:   row.UpdatedBy,
	}
	if row.StatusCode != nil {
		dto.StatusCode = row.StatusCode
	}
	if row.CreatedAt != nil {
		dto.CreatedAt = timestamppb.New(*row.CreatedAt)
	}
	if row.UpdatedAt != nil {
		dto.UpdatedAt = timestamppb.New(*row.UpdatedAt)
	}
	return dto
}
