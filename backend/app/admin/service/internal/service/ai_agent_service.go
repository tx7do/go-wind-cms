package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	sse "github.com/tx7do/kratos-transport/transport/sse"
	"google.golang.org/protobuf/types/known/emptypb"

	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	aiV1 "go-wind-cms/api/gen/go/ai/service/v1"
)

// AiAgentService AI Agent 编排管理（BFF 薄代理：应用 CRUD 与运行台账都在 core；
// 试运行 Run 由 core 的图执行引擎切片接入）。
// AiAgentRunPublisher SSE 发布者接口（sse_server 装配期注入；nil = 无流式观感）。
type AiAgentRunPublisher interface {
	Publish(ctx context.Context, streamId sse.StreamID, event *sse.Event)
}

type AiAgentService struct {
	adminV1.AiAgentServiceHTTPServer

	aiAgentServiceClient aiV1.AiAgentServiceClient
	log                  *log.Helper

	agentRunPublisher AiAgentRunPublisher
}

// RegisterAiAgentRunPublisher SSE 启动时注入真身（与 chat/站内信同模式）。
func (s *AiAgentService) RegisterAiAgentRunPublisher(p AiAgentRunPublisher) {
	if p != nil {
		s.agentRunPublisher = p
	}
}

func NewAiAgentService(ctx *bootstrap.Context, aiAgentServiceClient aiV1.AiAgentServiceClient) *AiAgentService {
	return &AiAgentService{
		log:                  ctx.NewLoggerHelper("ai-agent/service/admin-service"),
		aiAgentServiceClient: aiAgentServiceClient,
	}
}

func (s *AiAgentService) List(ctx context.Context, req *paginationV1.PagingRequest) (*aiV1.ListAiAgentAppsResponse, error) {
	return s.aiAgentServiceClient.List(ctx, req)
}

func (s *AiAgentService) Get(ctx context.Context, req *aiV1.GetAiAgentAppRequest) (*aiV1.AiAgentApp, error) {
	return s.aiAgentServiceClient.Get(ctx, req)
}

func (s *AiAgentService) Create(ctx context.Context, req *aiV1.CreateAiAgentAppRequest) (*aiV1.AiAgentApp, error) {
	return s.aiAgentServiceClient.Create(ctx, req)
}

func (s *AiAgentService) Update(ctx context.Context, req *aiV1.UpdateAiAgentAppRequest) (*emptypb.Empty, error) {
	return s.aiAgentServiceClient.Update(ctx, req)
}

func (s *AiAgentService) Delete(ctx context.Context, req *aiV1.DeleteAiAgentAppRequest) (*emptypb.Empty, error) {
	return s.aiAgentServiceClient.Delete(ctx, req)
}

// Run 试运行（透传 core）。执行引擎切片未接入前，core 侧为 gRPC
// Unimplemented——透传让它以规范的 501 呈现，而不是本层嵌入接口的空指针 502。
// ApproveRun 审批裁决运行（透传 core）。
func (s *AiAgentService) ApproveRun(ctx context.Context, req *aiV1.ApproveAgentRunRequest) (*aiV1.AiAgentRun, error) {
	return s.aiAgentServiceClient.ApproveRun(ctx, req)
}

func (s *AiAgentService) Run(ctx context.Context, req *aiV1.RunAiAgentAppRequest) (*aiV1.AiAgentRun, error) {
	return s.aiAgentServiceClient.Run(ctx, req)
}

func (s *AiAgentService) ListRuns(ctx context.Context, req *paginationV1.PagingRequest) (*aiV1.ListAiAgentRunsResponse, error) {
	return s.aiAgentServiceClient.ListRuns(ctx, req)
}

func (s *AiAgentService) GetRun(ctx context.Context, req *aiV1.GetAiAgentRunRequest) (*aiV1.AiAgentRun, error) {
	return s.aiAgentServiceClient.GetRun(ctx, req)
}
