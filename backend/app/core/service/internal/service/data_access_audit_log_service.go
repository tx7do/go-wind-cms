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

type DataAccessAuditLogService struct {
	auditV1.UnimplementedDataAccessAuditLogServiceServer

	log *log.Helper

	dataAccessAuditLogRepo *data.DataAccessAuditLogRepo
}

func NewDataAccessAuditLogService(ctx *bootstrap.Context, repo *data.DataAccessAuditLogRepo) *DataAccessAuditLogService {
	return &DataAccessAuditLogService{
		log:                    ctx.NewLoggerHelper("data-access-audit-log/service/core-service"),
		dataAccessAuditLogRepo: repo,
	}
}

func (s *DataAccessAuditLogService) List(ctx context.Context, req *paginationV1.PagingRequest) (*auditV1.ListDataAccessAuditLogResponse, error) {
	return s.dataAccessAuditLogRepo.List(ctx, req)
}

func (s *DataAccessAuditLogService) Get(ctx context.Context, req *auditV1.GetDataAccessAuditLogRequest) (*auditV1.DataAccessAuditLog, error) {
	return s.dataAccessAuditLogRepo.Get(ctx, req)
}

func (s *DataAccessAuditLogService) Create(ctx context.Context, req *auditV1.CreateDataAccessAuditLogRequest) (*emptypb.Empty, error) {
	if req == nil || req.Data == nil {
		return nil, auditV1.ErrorBadRequest("invalid parameter")
	}

	if err := s.dataAccessAuditLogRepo.Create(ctx, req); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
