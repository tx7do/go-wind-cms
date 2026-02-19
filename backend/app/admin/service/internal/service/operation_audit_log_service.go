package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	auditV1 "go-wind-cms/api/gen/go/audit/service/v1"
)

type OperationAuditLogService struct {
	adminV1.OperationAuditLogServiceHTTPServer

	log *log.Helper

	operationAuditLogServiceClient auditV1.OperationAuditLogServiceClient
}

func NewOperationAuditLogService(ctx *bootstrap.Context, operationAuditLogServiceClient auditV1.OperationAuditLogServiceClient) *OperationAuditLogService {
	return &OperationAuditLogService{
		log:                            ctx.NewLoggerHelper("operation-audit-log/service/admin-service"),
		operationAuditLogServiceClient: operationAuditLogServiceClient,
	}
}

func (s *OperationAuditLogService) List(ctx context.Context, req *paginationV1.PagingRequest) (*auditV1.ListOperationAuditLogResponse, error) {
	return s.operationAuditLogServiceClient.List(ctx, req)
}

func (s *OperationAuditLogService) Get(ctx context.Context, req *auditV1.GetOperationAuditLogRequest) (*auditV1.OperationAuditLog, error) {
	return s.operationAuditLogServiceClient.Get(ctx, req)
}

func (s *OperationAuditLogService) Create(ctx context.Context, req *auditV1.CreateOperationAuditLogRequest) (*emptypb.Empty, error) {
	if req == nil || req.Data == nil {
		return nil, adminV1.ErrorBadRequest("invalid parameter")
	}

	return s.operationAuditLogServiceClient.Create(ctx, req)
}
