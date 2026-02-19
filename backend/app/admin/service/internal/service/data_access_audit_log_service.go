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

type DataAccessAuditLogService struct {
	adminV1.DataAccessAuditLogServiceHTTPServer

	log *log.Helper

	dataAccessAuditLogServiceClient auditV1.DataAccessAuditLogServiceClient
}

func NewDataAccessAuditLogService(ctx *bootstrap.Context, dataAccessAuditLogServiceClient auditV1.DataAccessAuditLogServiceClient) *DataAccessAuditLogService {
	return &DataAccessAuditLogService{
		log:                             ctx.NewLoggerHelper("data-access-audit-log/service/admin-service"),
		dataAccessAuditLogServiceClient: dataAccessAuditLogServiceClient,
	}
}

func (s *DataAccessAuditLogService) List(ctx context.Context, req *paginationV1.PagingRequest) (*auditV1.ListDataAccessAuditLogResponse, error) {
	return s.dataAccessAuditLogServiceClient.List(ctx, req)
}

func (s *DataAccessAuditLogService) Get(ctx context.Context, req *auditV1.GetDataAccessAuditLogRequest) (*auditV1.DataAccessAuditLog, error) {
	return s.dataAccessAuditLogServiceClient.Get(ctx, req)
}

func (s *DataAccessAuditLogService) Create(ctx context.Context, req *auditV1.CreateDataAccessAuditLogRequest) (*emptypb.Empty, error) {
	if req == nil || req.Data == nil {
		return nil, adminV1.ErrorBadRequest("invalid parameter")
	}

	return s.dataAccessAuditLogServiceClient.Create(ctx, req)
}
