package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	flowV1 "go-wind-cms/api/gen/go/flow/service/v1"
)

// FlowRuleService 自动化 Flow 规则管理（admin 门面，透传 core）。
type FlowRuleService struct {
	adminV1.FlowRuleServiceHTTPServer

	log        *log.Helper
	flowClient flowV1.FlowRuleServiceClient
}

func NewFlowRuleService(ctx *bootstrap.Context, flowClient flowV1.FlowRuleServiceClient) *FlowRuleService {
	return &FlowRuleService{
		log:        ctx.NewLoggerHelper("flow/service/admin-service"),
		flowClient: flowClient,
	}
}

func (s *FlowRuleService) List(ctx context.Context, req *paginationV1.PagingRequest) (*flowV1.ListFlowRuleResponse, error) {
	return s.flowClient.List(ctx, req)
}

func (s *FlowRuleService) Get(ctx context.Context, req *flowV1.GetFlowRuleRequest) (*flowV1.FlowRule, error) {
	return s.flowClient.Get(ctx, req)
}

func (s *FlowRuleService) Create(ctx context.Context, req *flowV1.CreateFlowRuleRequest) (*flowV1.FlowRule, error) {
	return s.flowClient.Create(ctx, req)
}

func (s *FlowRuleService) Update(ctx context.Context, req *flowV1.UpdateFlowRuleRequest) (*flowV1.FlowRule, error) {
	return s.flowClient.Update(ctx, req)
}

func (s *FlowRuleService) Delete(ctx context.Context, req *flowV1.DeleteFlowRuleRequest) (*emptypb.Empty, error) {
	return s.flowClient.Delete(ctx, req)
}
