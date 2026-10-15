package service

import (
	"context"
	"encoding/json"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/go-crud/viewer"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	"go-wind-cms/app/core/service/internal/data"
	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
)

// 单条提交的大小约束：字段数/值长防滥用（公开端点，无鉴权）。
const (
	formMaxFields     = 32
	formMaxValueBytes = 4096
)

// FormSubmissionService 页面 FORM 区块的访客提交。
type FormSubmissionService struct {
	contentV1.UnimplementedFormSubmissionServiceServer

	formRepo   *data.FormSubmissionRepo
	pageRepo   *data.PageRepo
	log        *log.Helper
}

func NewFormSubmissionService(ctx *bootstrap.Context, formRepo *data.FormSubmissionRepo, pageRepo *data.PageRepo) *FormSubmissionService {
	return &FormSubmissionService{
		log:      ctx.NewLoggerHelper("form-submission/service/core-service"),
		formRepo: formRepo,
		pageRepo: pageRepo,
	}
}

// Submit 提交表单（匿名可调用）。
//
// 校验链：页面存在且已发布 → 区块属于该页面且为 FORM 类型 → 数据键集 ⊆
// config.fields 定义 → required 字段非空 → 单值大小受限。
// 租户归属与 RecordView 同款：登录态 viewer 优先，否则内网元数据回落。
func (s *FormSubmissionService) Submit(ctx context.Context, req *contentV1.SubmitFormSubmissionRequest) (*contentV1.FormSubmission, error) {
	if req == nil || req.GetPageId() == 0 || req.GetSectionId() == 0 {
		return nil, errors.New(400, "INVALID_PARAMETER", "invalid parameter")
	}

	// 页面必须存在且已发布（公开端点不接受向草稿页面投递数据）
	page, err := s.pageRepo.Get(ctx, &contentV1.GetPageRequest{
		QueryBy: &contentV1.GetPageRequest_Id{Id: req.GetPageId()},
	})
	if err != nil {
		return nil, errors.New(404, "FORM_PAGE_NOT_FOUND", "page not found")
	}
	if page.GetStatus() != contentV1.Page_PAGE_STATUS_PUBLISHED {
		return nil, errors.New(404, "FORM_PAGE_NOT_FOUND", "page not found")
	}

	// 区块须属于该页面且为 FORM 类型
	var section *contentV1.Section
	for _, sec := range page.GetSections() {
		if sec.GetId() == req.GetSectionId() {
			section = sec
			break
		}
	}
	if section == nil || section.GetType() != contentV1.SectionType_SECTION_TYPE_FORM {
		return nil, errors.New(400, "FORM_SECTION_MISMATCH", "section is not a form of this page")
	}

	fields, err := parseFormFields(section.GetConfig())
	if err != nil {
		return nil, errors.New(400, "FORM_FIELDS_INVALID", "form field definition is invalid")
	}

	if vErr := validateFormData(fields, req.GetData()); vErr != nil {
		return nil, vErr
	}

	// 提交者：登录态取 viewer 用户，匿名 0
	var userID uint32
	if vc, exist := viewer.FromContext(ctx); exist && vc != nil {
		userID = uint32(vc.UserID())
	}

	return s.formRepo.Create(ctx, req.GetPageId(), req.GetSectionId(), req.GetData(), userID)
}

// formField 区块字段定义（config.fields JSON 数组元素）。
type formField struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Required bool   `json:"required"`
}

// parseFormFields 解析区块 config.fields（JSON 数组字符串）；缺省视为无字段（空表单拒绝提交数据）。
func parseFormFields(config map[string]string) ([]formField, error) {
	raw := ""
	if config != nil {
		raw = config["fields"]
	}
	fields := make([]formField, 0, 8)
	if raw == "" {
		return fields, nil
	}
	if err := json.Unmarshal([]byte(raw), &fields); err != nil {
		return nil, err
	}
	for _, f := range fields {
		if f.Name == "" {
			return nil, errors.New(400, "FORM_FIELDS_INVALID", "field name is required")
		}
	}
	return fields, nil
}

// validateFormData 校验提交数据：键集 ⊆ 字段定义、required 非空、大小受限。
func validateFormData(fields []formField, data map[string]string) error {
	if len(data) > formMaxFields {
		return errors.New(400, "FORM_DATA_INVALID", "too many fields")
	}
	allowed := make(map[string]formField, len(fields))
	for _, f := range fields {
		allowed[f.Name] = f
	}
	for name, value := range data {
		def, ok := allowed[name]
		if !ok {
			return errors.New(400, "FORM_FIELD_UNKNOWN", "unknown field: "+name)
		}
		if len(value) > formMaxValueBytes {
			return errors.New(400, "FORM_DATA_INVALID", "field value too large: "+name)
		}
		if def.Required && value == "" {
			return errors.New(400, "FORM_REQUIRED_MISSING", "required field is empty: "+name)
		}
	}
	for _, f := range fields {
		if f.Required && data[f.Name] == "" {
			return errors.New(400, "FORM_REQUIRED_MISSING", "required field is empty: "+f.Name)
		}
	}
	return nil
}

// ─── admin 管理面 ───

func (s *FormSubmissionService) List(ctx context.Context, req *paginationV1.PagingRequest) (*contentV1.ListFormSubmissionsResponse, error) {
	items, total, err := s.formRepo.List(ctx, req)
	if err != nil {
		return nil, errors.New(500, "FORM_SUBMISSION_LIST_FAILED", "list form submissions failed")
	}
	return &contentV1.ListFormSubmissionsResponse{Items: items, Total: uint64(total)}, nil
}

func (s *FormSubmissionService) Get(ctx context.Context, req *contentV1.GetFormSubmissionRequest) (*contentV1.FormSubmission, error) {
	if req == nil || req.GetId() == 0 {
		return nil, errors.New(400, "INVALID_PARAMETER", "invalid parameter")
	}
	return s.formRepo.Get(ctx, req.GetId())
}

func (s *FormSubmissionService) Delete(ctx context.Context, req *contentV1.DeleteFormSubmissionRequest) (*emptypb.Empty, error) {
	if req == nil || req.GetId() == 0 {
		return nil, errors.New(400, "INVALID_PARAMETER", "invalid parameter")
	}
	if err := s.formRepo.Delete(ctx, req.GetId()); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
