package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	monitorAlertV1 "go-wind-cms/api/gen/go/monitor_alert/service/v1"

	"go-wind-cms/pkg/middleware/auth"
)

// MonitorAlertService 监控告警规则管理（BFF 薄代理：规则数据、评估内核与指标采集都在 core）。
type MonitorAlertService struct {
	adminV1.MonitorAlertServiceHTTPServer

	log             *log.Helper
	monitorAlertCli monitorAlertV1.MonitorAlertServiceClient
}

func NewMonitorAlertService(
	ctx *bootstrap.Context,
	monitorAlertCli monitorAlertV1.MonitorAlertServiceClient,
) *MonitorAlertService {
	l := log.NewHelper(log.With(ctx.GetLogger(), "module", "monitor-alert/service/admin-service"))
	return &MonitorAlertService{
		log:             l,
		monitorAlertCli: monitorAlertCli,
	}
}

func (s *MonitorAlertService) ListMonitorAlertRule(ctx context.Context, req *paginationV1.PagingRequest) (*monitorAlertV1.ListMonitorAlertRuleResponse, error) {
	return s.monitorAlertCli.ListMonitorAlertRule(ctx, req)
}

func (s *MonitorAlertService) GetMonitorAlertRule(ctx context.Context, req *monitorAlertV1.GetMonitorAlertRuleRequest) (*monitorAlertV1.MonitorAlertRule, error) {
	return s.monitorAlertCli.GetMonitorAlertRule(ctx, req)
}

func (s *MonitorAlertService) CreateMonitorAlertRule(ctx context.Context, req *monitorAlertV1.CreateMonitorAlertRuleRequest) (*monitorAlertV1.MonitorAlertRule, error) {
	if req == nil || req.Data == nil {
		return nil, adminV1.ErrorBadRequest("invalid parameter")
	}

	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	req.Data.CreatedBy = trans.Ptr(operator.UserId)

	return s.monitorAlertCli.CreateMonitorAlertRule(ctx, req)
}

func (s *MonitorAlertService) UpdateMonitorAlertRule(ctx context.Context, req *monitorAlertV1.UpdateMonitorAlertRuleRequest) (*emptypb.Empty, error) {
	if req == nil || req.Data == nil {
		return nil, adminV1.ErrorBadRequest("invalid parameter")
	}

	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	req.Data.UpdatedBy = trans.Ptr(operator.UserId)

	return s.monitorAlertCli.UpdateMonitorAlertRule(ctx, req)
}

func (s *MonitorAlertService) DeleteMonitorAlertRule(ctx context.Context, req *monitorAlertV1.DeleteMonitorAlertRuleRequest) (*emptypb.Empty, error) {
	return s.monitorAlertCli.DeleteMonitorAlertRule(ctx, req)
}

func (s *MonitorAlertService) EvaluateMonitorAlerts(ctx context.Context, req *monitorAlertV1.EvaluateMonitorAlertsRequest) (*monitorAlertV1.EvaluateMonitorAlertsResponse, error) {
	return s.monitorAlertCli.EvaluateMonitorAlerts(ctx, req)
}
