package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/go-crud/viewer"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	aiV1 "go-wind-cms/api/gen/go/ai/service/v1"

	"go-wind-cms/app/core/service/internal/data"
)

// AiAgentService AI Agent 编排应用管理（可视化画布定义的图应用）。
//
// 本切片（①）承载应用 CRUD 与运行台账读取；图执行引擎与试运行
// （Run/ListRuns 流式）在执行引擎切片接入——Run 走嵌入的
// Unimplemented 默认实现（501），proto 契约先行稳定。
type AiAgentService struct {
	aiV1.UnimplementedAiAgentServiceServer

	log           *log.Helper
	appRepo       *data.AiAgentAppRepo
	runRepo       *data.AiAgentRunRepo
}

func NewAiAgentService(ctx *bootstrap.Context, appRepo *data.AiAgentAppRepo, runRepo *data.AiAgentRunRepo) *AiAgentService {
	return &AiAgentService{
		log:     ctx.NewLoggerHelper("ai-agent/service/core-service"),
		appRepo: appRepo,
		runRepo: runRepo,
	}
}

// ─── 应用 CRUD ───

func (s *AiAgentService) List(ctx context.Context, req *paginationV1.PagingRequest) (*aiV1.ListAiAgentAppsResponse, error) {
	return s.appRepo.List(ctx, req)
}

func (s *AiAgentService) Get(ctx context.Context, req *aiV1.GetAiAgentAppRequest) (*aiV1.AiAgentApp, error) {
	if req == nil || req.GetId() == 0 {
		return nil, errors.BadRequest("AI_AGENT", "id is required")
	}
	return s.appRepo.Get(ctx, req.GetId())
}

// Create 创建应用。租户归属按操作者 viewer 显式写入（租户隔离 mixin 无创建钩子，
// 与 AiConversation 同口径）；平台管理员建在平台租户（tenant 0）。
func (s *AiAgentService) Create(ctx context.Context, req *aiV1.CreateAiAgentAppRequest) (*aiV1.AiAgentApp, error) {
	if req == nil || req.Data == nil {
		return nil, errors.BadRequest("AI_AGENT", "invalid parameter")
	}

	var operatorID, tenantID uint32
	if vc, exist := viewer.FromContext(ctx); exist && vc != nil {
		operatorID = uint32(vc.UserID())
		tenantID = uint32(vc.TenantID())
	}
	req.Data.TenantId = trans.Ptr(tenantID)

	id, err := s.appRepo.Create(ctx, req.GetData(), operatorID)
	if err != nil {
		return nil, err
	}
	return s.appRepo.Get(ctx, id)
}

func (s *AiAgentService) Update(ctx context.Context, req *aiV1.UpdateAiAgentAppRequest) (*emptypb.Empty, error) {
	if req == nil || req.Data == nil {
		return nil, errors.BadRequest("AI_AGENT", "invalid parameter")
	}
	if req.GetId() == 0 {
		return nil, errors.BadRequest("AI_AGENT", "id is required")
	}

	operatorID, _ := viewerUserIDFromContext(ctx)
	if err := s.appRepo.Update(ctx, req.GetId(), req.GetData(), req.GetUpdateMask().GetPaths(), operatorID); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// Delete 删除应用（运行台账保留——台账是历史事实，不随应用级联消失）。
func (s *AiAgentService) Delete(ctx context.Context, req *aiV1.DeleteAiAgentAppRequest) (*emptypb.Empty, error) {
	if req == nil || req.GetId() == 0 {
		return nil, errors.BadRequest("AI_AGENT", "id is required")
	}
	if err := s.appRepo.Delete(ctx, req.GetId()); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// ─── 运行台账（只读；写入在执行引擎切片） ───

func (s *AiAgentService) ListRuns(ctx context.Context, req *paginationV1.PagingRequest) (*aiV1.ListAiAgentRunsResponse, error) {
	return s.runRepo.List(ctx, req)
}

func (s *AiAgentService) GetRun(ctx context.Context, req *aiV1.GetAiAgentRunRequest) (*aiV1.AiAgentRun, error) {
	if req == nil || req.GetId() == 0 {
		return nil, errors.BadRequest("AI_AGENT", "id is required")
	}
	return s.runRepo.Get(ctx, req.GetId())
}
