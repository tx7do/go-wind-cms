package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
)

// EmailTemplateService 邮件模板管理门面（透传 core）。
type EmailTemplateService struct {
	adminV1.EmailTemplateServiceHTTPServer

	emailTemplateServiceClient contentV1.EmailTemplateServiceClient
	log                        *log.Helper
}

func NewEmailTemplateService(ctx *bootstrap.Context, emailTemplateServiceClient contentV1.EmailTemplateServiceClient) *EmailTemplateService {
	return &EmailTemplateService{
		log:                        ctx.NewLoggerHelper("email-template/service/admin-service"),
		emailTemplateServiceClient: emailTemplateServiceClient,
	}
}

func (s *EmailTemplateService) List(ctx context.Context, req *paginationV1.PagingRequest) (*contentV1.ListEmailTemplatesResponse, error) {
	return s.emailTemplateServiceClient.List(ctx, req)
}

func (s *EmailTemplateService) Get(ctx context.Context, req *contentV1.GetEmailTemplateRequest) (*contentV1.EmailTemplate, error) {
	return s.emailTemplateServiceClient.Get(ctx, req)
}

func (s *EmailTemplateService) Create(ctx context.Context, req *contentV1.CreateEmailTemplateRequest) (*contentV1.EmailTemplate, error) {
	return s.emailTemplateServiceClient.Create(ctx, req)
}

func (s *EmailTemplateService) Update(ctx context.Context, req *contentV1.UpdateEmailTemplateRequest) (*contentV1.EmailTemplate, error) {
	return s.emailTemplateServiceClient.Update(ctx, req)
}

func (s *EmailTemplateService) Delete(ctx context.Context, req *contentV1.DeleteEmailTemplateRequest) (*emptypb.Empty, error) {
	return s.emailTemplateServiceClient.Delete(ctx, req)
}
