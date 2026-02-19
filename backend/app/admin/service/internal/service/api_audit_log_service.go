package service

import (
	"context"
	resourceV1 "go-wind-cms/api/gen/go/resource/service/v1"
	"sync"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	auditV1 "go-wind-cms/api/gen/go/audit/service/v1"
)

type ApiAuditLogService struct {
	adminV1.ApiAuditLogServiceHTTPServer

	log *log.Helper

	apiAuditLogServiceClient auditV1.ApiAuditLogServiceClient
	apiServiceClient         resourceV1.ApiServiceClient

	apis     []*resourceV1.Api
	apiMutex sync.RWMutex
}

func NewApiAuditLogService(
	ctx *bootstrap.Context,
	apiAuditLogServiceClient auditV1.ApiAuditLogServiceClient,
	apiServiceClient resourceV1.ApiServiceClient,
) *ApiAuditLogService {
	return &ApiAuditLogService{
		log:                      ctx.NewLoggerHelper("api-audit-log/service/admin-service"),
		apiAuditLogServiceClient: apiAuditLogServiceClient,
		apiServiceClient:         apiServiceClient,
	}
}

func (s *ApiAuditLogService) queryApis(ctx context.Context, path, method string) (*resourceV1.Api, error) {
	if len(s.apis) == 0 {
		s.apiMutex.Lock()
		apis, err := s.apiServiceClient.List(ctx, &paginationV1.PagingRequest{
			NoPaging: trans.Ptr(true),
		})
		if err != nil {
			s.apiMutex.Unlock()
			return nil, err
		}
		s.apis = apis.Items
		s.apiMutex.Unlock()
	}

	if len(s.apis) == 0 {
		return nil, adminV1.ErrorNotFound("no apis found")
	}

	for _, api := range s.apis {
		if api.GetPath() == path && api.GetMethod() == method {
			return api, nil
		}
	}

	return nil, nil
}

func (s *ApiAuditLogService) List(ctx context.Context, req *paginationV1.PagingRequest) (*auditV1.ListApiAuditLogResponse, error) {
	resp, err := s.apiAuditLogServiceClient.List(ctx, req)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(resp.Items); i++ {
		l := resp.Items[i]
		if l == nil {
			continue
		}
		a, _ := s.queryApis(ctx, l.GetPath(), l.GetHttpMethod())
		if a != nil {
			l.ApiDescription = a.Description
			l.ApiModule = a.ModuleDescription
		}
	}

	return resp, nil
}

func (s *ApiAuditLogService) Get(ctx context.Context, req *auditV1.GetApiAuditLogRequest) (*auditV1.ApiAuditLog, error) {
	resp, err := s.apiAuditLogServiceClient.Get(ctx, req)
	if err != nil {
		return nil, err
	}

	a, _ := s.queryApis(ctx, resp.GetPath(), resp.GetHttpMethod())
	if a != nil {
		resp.ApiDescription = a.Description
		resp.ApiModule = a.ModuleDescription
	}

	return resp, nil
}

func (s *ApiAuditLogService) Create(ctx context.Context, req *auditV1.CreateApiAuditLogRequest) (*emptypb.Empty, error) {
	if req.Data == nil {
		return nil, adminV1.ErrorBadRequest("invalid parameter")
	}

	return s.apiAuditLogServiceClient.Create(ctx, req)
}
