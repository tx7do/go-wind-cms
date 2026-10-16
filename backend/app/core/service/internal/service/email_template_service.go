package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	"go-wind-cms/app/core/service/internal/data"
	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
)

// EmailTemplateService 邮件模板管理。
type EmailTemplateService struct {
	contentV1.UnimplementedEmailTemplateServiceServer

	emailTemplateRepo *data.EmailTemplateRepo
	log               *log.Helper
}

func NewEmailTemplateService(ctx *bootstrap.Context, emailTemplateRepo *data.EmailTemplateRepo) *EmailTemplateService {
	return &EmailTemplateService{
		log:               ctx.NewLoggerHelper("email-template/service/core-service"),
		emailTemplateRepo: emailTemplateRepo,
	}
}

func (s *EmailTemplateService) List(ctx context.Context, req *paginationV1.PagingRequest) (*contentV1.ListEmailTemplatesResponse, error) {
	items, total, err := s.emailTemplateRepo.List(ctx, req)
	if err != nil {
		return nil, err
	}
	return &contentV1.ListEmailTemplatesResponse{Items: items, Total: uint64(total)}, nil
}

func (s *EmailTemplateService) Get(ctx context.Context, req *contentV1.GetEmailTemplateRequest) (*contentV1.EmailTemplate, error) {
	if req == nil || req.GetId() == 0 {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}
	return s.emailTemplateRepo.Get(ctx, req.GetId())
}

func (s *EmailTemplateService) Create(ctx context.Context, req *contentV1.CreateEmailTemplateRequest) (*contentV1.EmailTemplate, error) {
	if req == nil || req.Data == nil || req.Data.GetCode() == "" || req.Data.GetName() == "" || req.Data.GetSubject() == "" {
		return nil, contentV1.ErrorBadRequest("code, name and subject are required")
	}
	return s.emailTemplateRepo.Create(ctx, req.Data)
}

func (s *EmailTemplateService) Update(ctx context.Context, req *contentV1.UpdateEmailTemplateRequest) (*contentV1.EmailTemplate, error) {
	if req == nil || req.Data == nil || req.GetId() == 0 {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}
	return s.emailTemplateRepo.Update(ctx, req.GetId(), req.Data)
}

func (s *EmailTemplateService) Delete(ctx context.Context, req *contentV1.DeleteEmailTemplateRequest) (*emptypb.Empty, error) {
	if req == nil || req.GetId() == 0 {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}
	if err := s.emailTemplateRepo.Delete(ctx, req.GetId()); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// FindEnabledByCode 按 code 取启用模板（出站 worker 渲染用；租户优先回退平台）。
func (s *EmailTemplateService) FindEnabledByCode(ctx context.Context, req *contentV1.FindEnabledEmailTemplateByCodeRequest) (*contentV1.EmailTemplate, error) {
	if req == nil || req.GetCode() == "" {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}
	tpl, err := s.emailTemplateRepo.FindEnabledByCode(ctx, req.GetCode())
	if err != nil {
		return nil, err
	}
	if tpl == nil {
		return &contentV1.EmailTemplate{}, nil
	}
	return tpl, nil
}
