package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"

	"go-wind-cms/app/core/service/internal/data"

	auditV1 "go-wind-cms/api/gen/go/audit/service/v1"
)

type LoginAuditLogService struct {
	auditV1.UnimplementedLoginAuditLogServiceServer

	log *log.Helper

	repo *data.LoginAuditLogRepo
}

func NewLoginAuditLogService(ctx *bootstrap.Context, repo *data.LoginAuditLogRepo) *LoginAuditLogService {
	return &LoginAuditLogService{
		log:  ctx.NewLoggerHelper("login-audit-log/service/core-service"),
		repo: repo,
	}
}

func (s *LoginAuditLogService) List(ctx context.Context, req *paginationV1.PagingRequest) (*auditV1.ListLoginAuditLogResponse, error) {
	return s.repo.List(ctx, req)
}

func (s *LoginAuditLogService) Get(ctx context.Context, req *auditV1.GetLoginAuditLogRequest) (*auditV1.LoginAuditLog, error) {
	return s.repo.Get(ctx, req)
}

func (s *LoginAuditLogService) Create(ctx context.Context, req *auditV1.CreateLoginAuditLogRequest) (*emptypb.Empty, error) {
	if req == nil || req.Data == nil {
		return nil, auditV1.ErrorBadRequest("invalid parameter")
	}

	if err := s.repo.Create(ctx, req); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
