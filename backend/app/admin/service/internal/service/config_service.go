package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	configV1 "go-wind-cms/api/gen/go/config/service/v1"

	"go-wind-cms/pkg/middleware/auth"
)

// ConfigService 系统参数管理（BFF 薄代理：KV 数据/播种/读取器都在 core）。
type ConfigService struct {
	adminV1.ConfigServiceHTTPServer

	log           *log.Helper
	configServiceClient configV1.ConfigServiceClient
}

func NewConfigService(
	ctx *bootstrap.Context,
	configServiceClient configV1.ConfigServiceClient,
) *ConfigService {
	l := log.NewHelper(log.With(ctx.GetLogger(), "module", "config/service/admin-service"))
	return &ConfigService{
		log:                 l,
		configServiceClient: configServiceClient,
	}
}

// List 系统参数列表（透传 core）。
func (s *ConfigService) List(ctx context.Context, req *paginationV1.PagingRequest) (*configV1.ListConfigResponse, error) {
	return s.configServiceClient.List(ctx, req)
}

// Get 系统参数详情（透传 core）。
func (s *ConfigService) Get(ctx context.Context, req *configV1.GetConfigRequest) (*configV1.Config, error) {
	return s.configServiceClient.Get(ctx, req)
}

// Create 创建系统参数（操作人盖章后透传）。
func (s *ConfigService) Create(ctx context.Context, req *configV1.CreateConfigRequest) (*emptypb.Empty, error) {
	if req == nil || req.Data == nil {
		return nil, adminV1.ErrorBadRequest("invalid parameter")
	}
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	req.Data.CreatedBy = trans.Ptr(operator.UserId)
	return s.configServiceClient.Create(ctx, req)
}

// Update 更新系统参数（操作人盖章后透传）。
func (s *ConfigService) Update(ctx context.Context, req *configV1.UpdateConfigRequest) (*emptypb.Empty, error) {
	if req == nil || req.Data == nil {
		return nil, adminV1.ErrorBadRequest("invalid parameter")
	}
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	req.Data.UpdatedBy = trans.Ptr(operator.UserId)
	if req.UpdateMask != nil {
		req.UpdateMask.Paths = append(req.UpdateMask.Paths, "updated_by")
	}
	return s.configServiceClient.Update(ctx, req)
}

// Delete 删除系统参数（内置参数禁删守卫在 core；透传）。
func (s *ConfigService) Delete(ctx context.Context, req *configV1.DeleteConfigRequest) (*emptypb.Empty, error) {
	return s.configServiceClient.Delete(ctx, req)
}
