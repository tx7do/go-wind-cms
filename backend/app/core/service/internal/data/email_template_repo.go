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
	"go-wind-cms/app/core/service/internal/data/ent/emailtemplate"
	"go-wind-cms/app/core/service/internal/data/ent/predicate"
	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
)

// EmailTemplateRepo 邮件模板仓储。
type EmailTemplateRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper
}

func NewEmailTemplateRepo(ctx *bootstrap.Context, entClient *entCrud.EntClient[*ent.Client]) *EmailTemplateRepo {
	return &EmailTemplateRepo{
		entClient: entClient,
		log:       ctx.NewLoggerHelper("email-template/repo/core-service"),
	}
}

// List 模板列表（时间倒序；支持 code/name 模糊、enabled 过滤）。
func (r *EmailTemplateRepo) List(ctx context.Context, req *paginationV1.PagingRequest) ([]*contentV1.EmailTemplate, int, error) {
	conds := r.tenantConds(ctx)
	if req != nil {
		if filterExpr, fErr := paginationFilter.ConvertFilterByPagingRequest(req); fErr == nil && filterExpr != nil {
			for _, cond := range filterExpr.GetConditions() {
				switch cond.GetField() {
				case "code":
					if v := cond.GetValue(); v != "" {
						conds = append(conds, emailtemplate.CodeContains(v))
					}
				case "name":
					if v := cond.GetValue(); v != "" {
						conds = append(conds, emailtemplate.NameContains(v))
					}
				case "enabled":
					conds = append(conds, emailtemplate.EnabledEQ(cond.GetValue() == "true"))
				}
			}
		}
	}

	total, err := r.entClient.Client().EmailTemplate.Query().Where(conds...).Count(ctx)
	if err != nil {
		r.log.Errorf("count email templates failed: %s", err.Error())
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

	rows, err := r.entClient.Client().EmailTemplate.Query().Where(conds...).
		Order(ent.Desc(emailtemplate.FieldID)).
		Limit(int(pageSize)).
		Offset(int((page - 1) * pageSize)).
		All(ctx)
	if err != nil {
		r.log.Errorf("list email templates failed: %s", err.Error())
		return nil, 0, err
	}

	items := make([]*contentV1.EmailTemplate, 0, len(rows))
	for _, row := range rows {
		items = append(items, r.toDTO(row))
	}
	return items, total, nil
}

// Get 模板详情。
func (r *EmailTemplateRepo) Get(ctx context.Context, id uint32) (*contentV1.EmailTemplate, error) {
	conds := append(r.tenantConds(ctx), emailtemplate.IDEQ(id))
	row, err := r.entClient.Client().EmailTemplate.Query().Where(conds...).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.NotFound("EMAIL_TEMPLATE", "email template not found")
		}
		r.log.Errorf("get email template failed (id=%d): %s", id, err.Error())
		return nil, err
	}
	return r.toDTO(row), nil
}

// Create 新建（同租户 code 查重）。
func (r *EmailTemplateRepo) Create(ctx context.Context, dto *contentV1.EmailTemplate) (*contentV1.EmailTemplate, error) {
	dup, err := r.entClient.Client().EmailTemplate.Query().Where(append(
		r.tenantConds(ctx),
		emailtemplate.CodeEQ(dto.GetCode()),
	)...).Count(ctx)
	if err != nil {
		r.log.Errorf("check email template duplicate failed: %s", err.Error())
		return nil, err
	}
	if dup > 0 {
		return nil, errors.Conflict("EMAIL_TEMPLATE", "code already exists")
	}

	builder := r.entClient.Client().EmailTemplate.Create().
		SetCode(dto.GetCode()).
		SetName(dto.GetName()).
		SetSubject(dto.GetSubject()).
		SetEnabled(dto.GetEnabled())
	if dto.GetBodyHtml() != "" {
		builder.SetBodyHTML(dto.GetBodyHtml())
	}
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
		r.log.Errorf("create email template failed: %s", err.Error())
		return nil, err
	}
	return r.toDTO(row), nil
}

// Update 按 mask 更新。
func (r *EmailTemplateRepo) Update(ctx context.Context, id uint32, dto *contentV1.EmailTemplate) (*contentV1.EmailTemplate, error) {
	conds := append(r.tenantConds(ctx), emailtemplate.IDEQ(id))
	exist, err := r.entClient.Client().EmailTemplate.Query().Where(conds...).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.NotFound("EMAIL_TEMPLATE", "email template not found")
		}
		return nil, err
	}

	builder := r.entClient.Client().EmailTemplate.UpdateOneID(exist.ID)
	if dto.Code != nil {
		builder.SetCode(dto.GetCode())
	}
	if dto.Name != nil {
		builder.SetName(dto.GetName())
	}
	if dto.Subject != nil {
		builder.SetSubject(dto.GetSubject())
	}
	if dto.BodyHtml != nil {
		if dto.GetBodyHtml() == "" {
			builder.ClearBodyHTML()
		} else {
			builder.SetBodyHTML(dto.GetBodyHtml())
		}
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
		r.log.Errorf("update email template failed (id=%d): %s", id, err.Error())
		return nil, err
	}
	return r.toDTO(row), nil
}

// Delete 删除。
func (r *EmailTemplateRepo) Delete(ctx context.Context, id uint32) error {
	conds := append(r.tenantConds(ctx), emailtemplate.IDEQ(id))
	count, err := r.entClient.Client().EmailTemplate.Delete().Where(conds...).Exec(ctx)
	if err != nil {
		r.log.Errorf("delete email template failed (id=%d): %s", id, err.Error())
		return err
	}
	if count == 0 {
		return errors.NotFound("EMAIL_TEMPLATE", "email template not found")
	}
	return nil
}

// FindEnabledByCode 按代码取启用模板（出站 worker 渲染用；无租户口径过滤，
// code 全局唯一约定，租户优先由服务层先查租户再回退平台）。
func (r *EmailTemplateRepo) FindEnabledByCode(ctx context.Context, code string) (*contentV1.EmailTemplate, error) {
	preds := []predicate.EmailTemplate{emailtemplate.CodeEQ(code), emailtemplate.EnabledEQ(true)}
	// 有租户视角先按租户查，查不到回退平台（tenant_id=0）
	if tid, hasTenant := maybeTenantFromViewer(ctx); hasTenant {
		row, err := r.entClient.Client().EmailTemplate.Query().
			Where(append(preds, emailtemplate.TenantIDEQ(tid))...).First(ctx)
		if err == nil {
			return r.toDTO(row), nil
		}
		if !ent.IsNotFound(err) {
			r.log.Errorf("find tenant email template by code failed (code=%s): %s", code, err.Error())
		}
	}
	row, err := r.entClient.Client().EmailTemplate.Query().
		Where(append(preds, emailtemplate.TenantIDEQ(0))...).First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		r.log.Errorf("find platform email template by code failed (code=%s): %s", code, err.Error())
		return nil, err
	}
	return r.toDTO(row), nil
}

func (r *EmailTemplateRepo) toDTO(row *ent.EmailTemplate) *contentV1.EmailTemplate {
	if row == nil {
		return nil
	}
	dto := &contentV1.EmailTemplate{
		Id:          &row.ID,
		Code:        &row.Code,
		Name:        &row.Name,
		Subject:     &row.Subject,
		BodyHtml:    nillableString(row.BodyHTML),
		Enabled:     nillableBoolPtr(row.Enabled),
		Description: row.Description,
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
	return dto
}

func (r *EmailTemplateRepo) tenantConds(ctx context.Context) []predicate.EmailTemplate {
	tid, hasTenant := maybeTenantFromViewer(ctx)
	if !hasTenant {
		return nil
	}
	return []predicate.EmailTemplate{emailtemplate.TenantIDEQ(tid)}
}

// tenantIDFromInternalMetadata 从内网可信元数据 x-md-global-tenant-id 读租户。
func (r *EmailTemplateRepo) tenantIDFromInternalMetadata(ctx context.Context) uint32 {
	if md, ok := kratosMetadata.FromServerContext(ctx); ok {
		if v := md.Get("x-md-global-tenant-id"); v != "" {
			if tid, err := strconv.ParseUint(v, 10, 32); err == nil && tid > 0 {
				return uint32(tid)
			}
		}
	}
	return 0
}

func nillableBoolPtr(b *bool) *bool {
	if b == nil {
		return nil
	}
	v := *b
	return &v
}
