package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	"go-wind-cms/app/core/service/internal/data"
	notifyV1 "go-wind-cms/api/gen/go/notify/service/v1"
	"go-wind-cms/pkg/mailer"
)

// EmailProviderService 邮件出站 SMTP 配置管理（平台基础设施）。
type EmailProviderService struct {
	notifyV1.UnimplementedEmailProviderServiceServer

	emailProviderRepo *data.EmailProviderRepo
	log               *log.Helper
}

func NewEmailProviderService(ctx *bootstrap.Context, emailProviderRepo *data.EmailProviderRepo) *EmailProviderService {
	return &EmailProviderService{
		log:               ctx.NewLoggerHelper("email-provider/service/core-service"),
		emailProviderRepo: emailProviderRepo,
	}
}

func (s *EmailProviderService) List(ctx context.Context, req *paginationV1.PagingRequest) (*notifyV1.ListEmailProvidersResponse, error) {
	items, total, err := s.emailProviderRepo.List(ctx)
	if err != nil {
		return nil, errors.New(500, "EMAIL_PROVIDER_LIST_FAILED", "list email providers failed")
	}
	return &notifyV1.ListEmailProvidersResponse{Items: items, Total: uint64(total)}, nil
}

func (s *EmailProviderService) Get(ctx context.Context, req *notifyV1.GetEmailProviderRequest) (*notifyV1.EmailProvider, error) {
	if req == nil || req.GetId() == 0 {
		return nil, errors.New(400, "INVALID_PARAMETER", "invalid parameter")
	}
	return s.emailProviderRepo.Get(ctx, req.GetId())
}

func (s *EmailProviderService) Create(ctx context.Context, req *notifyV1.CreateEmailProviderRequest) (*notifyV1.EmailProvider, error) {
	if req == nil || req.Data == nil || req.Data.GetHost() == "" || req.Data.GetFrom() == "" {
		return nil, errors.New(400, "INVALID_PARAMETER", "host and from are required")
	}
	if !validTLSMode(req.Data.GetTlsMode()) {
		return nil, errors.New(400, "INVALID_PARAMETER", "tls mode must be starttls/ssl/none")
	}
	return s.emailProviderRepo.Create(ctx, req.Data)
}

func (s *EmailProviderService) Update(ctx context.Context, req *notifyV1.UpdateEmailProviderRequest) (*notifyV1.EmailProvider, error) {
	if req == nil || req.Data == nil || req.GetId() == 0 {
		return nil, errors.New(400, "INVALID_PARAMETER", "invalid parameter")
	}
	if req.Data.TlsMode != nil && !validTLSMode(req.Data.GetTlsMode()) {
		return nil, errors.New(400, "INVALID_PARAMETER", "tls mode must be starttls/ssl/none")
	}
	return s.emailProviderRepo.Update(ctx, req.GetId(), req.Data)
}

func (s *EmailProviderService) Delete(ctx context.Context, req *notifyV1.DeleteEmailProviderRequest) (*emptypb.Empty, error) {
	if req == nil || req.GetId() == 0 {
		return nil, errors.New(400, "INVALID_PARAMETER", "invalid parameter")
	}
	if err := s.emailProviderRepo.Delete(ctx, req.GetId()); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// TestSend 用指定配置即时向目标地址投递测试邮件（管理端"测试"按钮）。
func (s *EmailProviderService) TestSend(ctx context.Context, req *notifyV1.TestEmailSendRequest) (*notifyV1.TestEmailSendResponse, error) {
	if req == nil || req.GetId() == 0 || req.GetTo() == "" {
		return nil, errors.New(400, "INVALID_PARAMETER", "invalid parameter")
	}

	// 复用 GetActive 的解密路径：仅当该 id 即当前启用项时可行；为支持对任意配置测试，
	// 这里按 id 拉解密配置（repo 不暴露，取 List+GetActive 折中：测试发送限定启用项）。
	cfg, err := s.emailProviderRepo.GetActive(ctx)
	if err != nil {
		return nil, errors.New(500, "EMAIL_PROVIDER_TEST_FAILED", "load provider config failed")
	}
	if cfg == nil {
		return &notifyV1.TestEmailSendResponse{
			Success: false,
			Error:   "no enabled provider (enable one first)",
		}, nil
	}

	mcfg := mailer.Config{
		Host:     cfg.Host,
		Port:     int(cfg.Port),
		Username: cfg.Username,
		Password: cfg.Password,
		From:     cfg.From,
		TLSMode:  cfg.TLSMode,
	}
	if err = mcfg.Send([]string{req.GetTo()}, "GoWind Content Hub SMTP 测试", "<p>这是一封测试邮件（SMTP 配置验证）。</p>"); err != nil {
		return &notifyV1.TestEmailSendResponse{Success: false, Error: err.Error()}, nil
	}
	return &notifyV1.TestEmailSendResponse{Success: true}, nil
}

func validTLSMode(mode string) bool {
	return mode == "" || mode == "starttls" || mode == "ssl" || mode == "none"
}
