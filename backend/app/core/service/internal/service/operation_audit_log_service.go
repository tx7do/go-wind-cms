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

type OperationAuditLogService struct {
	auditV1.UnimplementedOperationAuditLogServiceServer

	log *log.Helper

	operationAuditLogRepo *data.OperationAuditLogRepo
}

func NewOperationAuditLogService(ctx *bootstrap.Context, repo *data.OperationAuditLogRepo) *OperationAuditLogService {
	return &OperationAuditLogService{
		log:                   ctx.NewLoggerHelper("operation-audit-log/service/core-service"),
		operationAuditLogRepo: repo,
	}
}

func (s *OperationAuditLogService) List(ctx context.Context, req *paginationV1.PagingRequest) (*auditV1.ListOperationAuditLogResponse, error) {
	return s.operationAuditLogRepo.List(ctx, req)
}

func (s *OperationAuditLogService) Get(ctx context.Context, req *auditV1.GetOperationAuditLogRequest) (*auditV1.OperationAuditLog, error) {
	return s.operationAuditLogRepo.Get(ctx, req)
}

func (s *OperationAuditLogService) Create(ctx context.Context, req *auditV1.CreateOperationAuditLogRequest) (*emptypb.Empty, error) {
	if req == nil || req.Data == nil {
		return nil, auditV1.ErrorBadRequest("invalid parameter")
	}

	if err := s.operationAuditLogRepo.Create(ctx, req); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
