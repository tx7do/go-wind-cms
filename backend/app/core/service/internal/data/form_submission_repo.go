package data

import (
	"context"
	"strconv"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/timestamppb"

	paginationFilter "github.com/tx7do/go-crud/pagination/filter"
	kratosMetadata "github.com/go-kratos/kratos/v2/metadata"

	entCrud "github.com/tx7do/go-crud/entgo"

	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/formsubmission"
	"go-wind-cms/app/core/service/internal/data/ent/predicate"
	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
)

// FormSubmissionRepo 表单提交仓储（只增读删，无更新——提交数据不可篡改）。
type FormSubmissionRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper
}

func NewFormSubmissionRepo(ctx *bootstrap.Context, entClient *entCrud.EntClient[*ent.Client]) *FormSubmissionRepo {
	return &FormSubmissionRepo{
		entClient: entClient,
		log:       ctx.NewLoggerHelper("form-submission/repo/core-service"),
	}
}

// Create 落一条提交。
func (r *FormSubmissionRepo) Create(ctx context.Context, pageID, sectionID uint32, data map[string]string, userID uint32) (*contentV1.FormSubmission, error) {
	builder := r.entClient.Client().FormSubmission.Create().
		SetNillablePageID(&pageID).
		SetNillableSectionID(&sectionID).
		SetNillableUserID(&userID)
	if data != nil {
		builder.SetData(&data)
	}
	if tid, hasTenant := maybeTenantFromViewer(ctx); hasTenant {
		builder.SetTenantID(tid)
	} else if tid := r.tenantIDFromInternalMetadata(ctx); tid > 0 {
		builder.SetTenantID(tid)
	}
	row, err := builder.Save(ctx)
	if err != nil {
		r.log.Errorf("create form submission failed: %s", err.Error())
		return nil, err
	}
	return r.toDTO(row), nil
}

// List 提交列表（时间倒序；支持 page_id/section_id 过滤，经 PagingRequest 的 query/filterExpr 传入）。
func (r *FormSubmissionRepo) List(ctx context.Context, req *paginationV1.PagingRequest) ([]*contentV1.FormSubmission, int, error) {
	conds := r.tenantConds(ctx)

	if req != nil {
		if filterExpr, fErr := paginationFilter.ConvertFilterByPagingRequest(req); fErr == nil && filterExpr != nil {
			for _, cond := range filterExpr.GetConditions() {
				switch cond.GetField() {
				case "pageId", "page_id":
					if v, ok := condValueUint(cond.GetValue()); ok {
						conds = append(conds, formsubmission.PageIDEQ(v))
					}
				case "sectionId", "section_id":
					if v, ok := condValueUint(cond.GetValue()); ok {
						conds = append(conds, formsubmission.SectionIDEQ(v))
					}
				}
			}
		}
	}

	total, err := r.entClient.Client().FormSubmission.Query().Where(conds...).Count(ctx)
	if err != nil {
		r.log.Errorf("count form submissions failed: %s", err.Error())
		return nil, 0, err
	}

	page, pageSize := uint32(1), uint32(20)
	if req != nil {
		if req.GetPage() > 0 {
			page = req.GetPage()
		}
		if req.GetPageSize() > 0 {
			pageSize = req.GetPageSize()
		}
	}
	if pageSize > 200 {
		pageSize = 200
	}

	rows, err := r.entClient.Client().FormSubmission.Query().Where(conds...).
		Order(ent.Desc(formsubmission.FieldID)).
		Limit(int(pageSize)).
		Offset(int((page - 1) * pageSize)).
		All(ctx)
	if err != nil {
		r.log.Errorf("list form submissions failed: %s", err.Error())
		return nil, 0, err
	}

	items := make([]*contentV1.FormSubmission, 0, len(rows))
	for _, row := range rows {
		items = append(items, r.toDTO(row))
	}
	return items, total, nil
}

// condValueUint 把过滤条件值解析为 uint32（仅接受十进制正整数）。
func condValueUint(v string) (uint32, bool) {
	if v == "" {
		return 0, false
	}
	n, err := strconv.ParseUint(v, 10, 32)
	if err != nil || n == 0 {
		return 0, false
	}
	return uint32(n), true
}

// Get 提交详情。
func (r *FormSubmissionRepo) Get(ctx context.Context, id uint32) (*contentV1.FormSubmission, error) {
	conds := append(r.tenantConds(ctx), formsubmission.IDEQ(id))
	row, err := r.entClient.Client().FormSubmission.Query().Where(conds...).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.NotFound("FORM_SUBMISSION", "form submission not found")
		}
		r.log.Errorf("get form submission failed (id=%d): %s", id, err.Error())
		return nil, err
	}
	return r.toDTO(row), nil
}

// Delete 删除提交。
func (r *FormSubmissionRepo) Delete(ctx context.Context, id uint32) error {
	conds := append(r.tenantConds(ctx), formsubmission.IDEQ(id))
	count, err := r.entClient.Client().FormSubmission.Delete().Where(conds...).Exec(ctx)
	if err != nil {
		r.log.Errorf("delete form submission failed (id=%d): %s", id, err.Error())
		return err
	}
	if count == 0 {
		return errors.NotFound("FORM_SUBMISSION", "form submission not found")
	}
	return nil
}

func (r *FormSubmissionRepo) tenantConds(ctx context.Context) []predicate.FormSubmission {
	tid, hasTenant := maybeTenantFromViewer(ctx)
	if !hasTenant {
		return nil
	}
	return []predicate.FormSubmission{formsubmission.TenantIDEQ(tid)}
}

// tenantIDFromInternalMetadata 从内网可信元数据 x-md-global-tenant-id 读租户
//（匿名提交链路：app BFF 按 Host 解析后注入；与 comment_repo 同口径）。
func (r *FormSubmissionRepo) tenantIDFromInternalMetadata(ctx context.Context) uint32 {
	if md, ok := kratosMetadata.FromServerContext(ctx); ok {
		if v := md.Get("x-md-global-tenant-id"); v != "" {
			if tid, err := strconv.ParseUint(v, 10, 32); err == nil && tid > 0 {
				return uint32(tid)
			}
		}
	}
	return 0
}

func (r *FormSubmissionRepo) toDTO(row *ent.FormSubmission) *contentV1.FormSubmission {
	dto := &contentV1.FormSubmission{
		Id:      &row.ID,
		PageId:  row.PageID,
		SectionId: row.SectionID,
		UserId:  row.UserID,
	}
	if row.Data != nil {
		dto.Data = *row.Data
	}
	if row.CreatedAt != nil {
		dto.CreatedAt = timestamppb.New(*row.CreatedAt)
	}
	return dto
}
