package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	appViewer "go-wind-cms/pkg/entgo/viewer"
	"go-wind-cms/pkg/constants"

	configV1 "go-wind-cms/api/gen/go/config/service/v1"

	"strings"

	"go-wind-cms/app/core/service/internal/data"
)

// ConfigService 系统参数（平台全局动态 KV）管理 + 启动期内置键播种。
type ConfigService struct {
	configV1.UnimplementedConfigServiceServer

	log        *log.Helper
	configRepo *data.ConfigRepo
}

func NewConfigService(ctx *bootstrap.Context, configRepo *data.ConfigRepo) *ConfigService {
	svc := &ConfigService{
		log:        ctx.NewLoggerHelper("config/service/core-service"),
		configRepo: configRepo,
	}

	// 播种内置平台参数（等保口令策略阈值）：按键缺一补一、不覆盖既有值。
	svc.seedDefaults()

	return svc
}

func (s *ConfigService) seedDefaults() {
	ctx := appViewer.NewSystemViewerContext(context.Background())
	if err := s.configRepo.SeedDefaults(ctx, constants.DefaultConfigs); err != nil {
		s.log.Errorf("seed default configs failed: %s", err.Error())
	}
}

func (s *ConfigService) List(ctx context.Context, req *paginationV1.PagingRequest) (*configV1.ListConfigResponse, error) {
	return s.configRepo.List(ctx, req)
}

func (s *ConfigService) Get(ctx context.Context, req *configV1.GetConfigRequest) (*configV1.Config, error) {
	return s.configRepo.Get(ctx, req)
}

func (s *ConfigService) Create(ctx context.Context, req *configV1.CreateConfigRequest) (*emptypb.Empty, error) {
	if req == nil || req.Data == nil {
		return nil, errors.BadRequest("CONFIG", "invalid parameter")
	}
	if strings.TrimSpace(req.Data.GetKey()) == "" {
		return nil, errors.BadRequest("CONFIG", "config key is required")
	}

	if operatorID, ok := viewerUserIDFromContext(ctx); ok {
		req.Data.CreatedBy = trans.Ptr(operatorID)
	}

	if err := s.configRepo.Create(ctx, req); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *ConfigService) Update(ctx context.Context, req *configV1.UpdateConfigRequest) (*emptypb.Empty, error) {
	if req == nil || req.Data == nil {
		return nil, errors.BadRequest("CONFIG", "invalid parameter")
	}

	req.Data.Id = trans.Ptr(req.GetId())

	if operatorID, ok := viewerUserIDFromContext(ctx); ok {
		req.Data.UpdatedBy = trans.Ptr(operatorID)
	}

	if err := s.configRepo.Update(ctx, req); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *ConfigService) Delete(ctx context.Context, req *configV1.DeleteConfigRequest) (*emptypb.Empty, error) {
	if req == nil || req.GetId() == 0 {
		return nil, errors.BadRequest("CONFIG", "id is required")
	}
	if err := s.configRepo.Delete(ctx, req); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
