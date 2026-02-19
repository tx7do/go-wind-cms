package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	permissionV1 "go-wind-cms/api/gen/go/permission/service/v1"
)

type PolicyEvaluationLogService struct {
	adminV1.PolicyEvaluationLogServiceHTTPServer

	log *log.Helper

	policyEvaluationLogServiceClient permissionV1.PolicyEvaluationLogServiceClient
}

func NewPolicyEvaluationLogService(
	ctx *bootstrap.Context,
	policyEvaluationLogServiceClient permissionV1.PolicyEvaluationLogServiceClient,
) *PolicyEvaluationLogService {
	return &PolicyEvaluationLogService{
		log:                              ctx.NewLoggerHelper("policy-evaluation-log/service/admin-service"),
		policyEvaluationLogServiceClient: policyEvaluationLogServiceClient,
	}
}

func (s *PolicyEvaluationLogService) List(ctx context.Context, req *paginationV1.PagingRequest) (*permissionV1.ListPolicyEvaluationLogResponse, error) {
	resp, err := s.policyEvaluationLogServiceClient.List(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *PolicyEvaluationLogService) Get(ctx context.Context, req *permissionV1.GetPolicyEvaluationLogRequest) (*permissionV1.PolicyEvaluationLog, error) {
	resp, err := s.policyEvaluationLogServiceClient.Get(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *PolicyEvaluationLogService) Create(ctx context.Context, req *permissionV1.CreatePolicyEvaluationLogRequest) (*emptypb.Empty, error) {
	if req.Data == nil {
		return nil, adminV1.ErrorBadRequest("invalid parameter")
	}

	return s.policyEvaluationLogServiceClient.Create(ctx, req)
}
