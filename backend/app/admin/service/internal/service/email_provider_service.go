package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	notifyV1 "go-wind-cms/api/gen/go/notify/service/v1"
)

// EmailProviderService 邮件出站 SMTP 配置管理门面（透传 core）。
type EmailProviderService struct {
	adminV1.EmailProviderServiceHTTPServer

	emailProviderServiceClient notifyV1.EmailProviderServiceClient
	log                        *log.Helper
}

func NewEmailProviderService(ctx *bootstrap.Context, emailProviderServiceClient notifyV1.EmailProviderServiceClient) *EmailProviderService {
	return &EmailProviderService{
		log:                        ctx.NewLoggerHelper("email-provider/service/admin-service"),
		emailProviderServiceClient: emailProviderServiceClient,
	}
}

func (s *EmailProviderService) List(ctx context.Context, req *paginationV1.PagingRequest) (*notifyV1.ListEmailProvidersResponse, error) {
	return s.emailProviderServiceClient.List(ctx, req)
}

func (s *EmailProviderService) Get(ctx context.Context, req *notifyV1.GetEmailProviderRequest) (*notifyV1.EmailProvider, error) {
	return s.emailProviderServiceClient.Get(ctx, req)
}

func (s *EmailProviderService) Create(ctx context.Context, req *notifyV1.CreateEmailProviderRequest) (*notifyV1.EmailProvider, error) {
	return s.emailProviderServiceClient.Create(ctx, req)
}

func (s *EmailProviderService) Update(ctx context.Context, req *notifyV1.UpdateEmailProviderRequest) (*notifyV1.EmailProvider, error) {
	return s.emailProviderServiceClient.Update(ctx, req)
}

func (s *EmailProviderService) Delete(ctx context.Context, req *notifyV1.DeleteEmailProviderRequest) (*emptypb.Empty, error) {
	return s.emailProviderServiceClient.Delete(ctx, req)
}

func (s *EmailProviderService) TestSend(ctx context.Context, req *notifyV1.TestEmailSendRequest) (*notifyV1.TestEmailSendResponse, error) {
	return s.emailProviderServiceClient.TestSend(ctx, req)
}
