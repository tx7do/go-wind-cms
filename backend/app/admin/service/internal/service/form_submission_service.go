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

// FormSubmissionService 表单提交管理门面（透传 core）。
type FormSubmissionService struct {
	adminV1.FormSubmissionServiceHTTPServer

	formSubmissionServiceClient contentV1.FormSubmissionServiceClient
	log                         *log.Helper
}

func NewFormSubmissionService(ctx *bootstrap.Context, formSubmissionServiceClient contentV1.FormSubmissionServiceClient) *FormSubmissionService {
	return &FormSubmissionService{
		log:                         ctx.NewLoggerHelper("form-submission/service/admin-service"),
		formSubmissionServiceClient: formSubmissionServiceClient,
	}
}

func (s *FormSubmissionService) List(ctx context.Context, req *paginationV1.PagingRequest) (*contentV1.ListFormSubmissionsResponse, error) {
	return s.formSubmissionServiceClient.List(ctx, req)
}

func (s *FormSubmissionService) Get(ctx context.Context, req *contentV1.GetFormSubmissionRequest) (*contentV1.FormSubmission, error) {
	return s.formSubmissionServiceClient.Get(ctx, req)
}

func (s *FormSubmissionService) Delete(ctx context.Context, req *contentV1.DeleteFormSubmissionRequest) (*emptypb.Empty, error) {
	return s.formSubmissionServiceClient.Delete(ctx, req)
}
