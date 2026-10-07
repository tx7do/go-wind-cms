package service

import (
	"context"
	"io"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"github.com/tx7do/kratos-transport/transport/sse"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	aiV1 "go-wind-cms/api/gen/go/ai/service/v1"
	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	authenticationV1 "go-wind-cms/api/gen/go/authentication/service/v1"

	"go-wind-cms/pkg/middleware/auth"
)

// AI chat SSE 事件名（与前端 src/core/transport/sse/event.ts 注册表逐字对齐）。
const (
	sseEventAIChatChunk = "ai_chat_chunk"
	sseEventAIChatTool  = "ai_chat_tool"
)

// AiChatPublisher SSE 推送接口（与 InternalMessagePublisher 同型——同一个
// kratos sse.Server 实例同时实现两者，装配期注入）。
type AiChatPublisher interface {
	Publish(ctx context.Context, streamId sse.StreamID, event *sse.Event)
}

// noopAiChatPublisher SSE 未配置时的占位实现：对话仍完整可用（同步响应兜底），只是不推送 chunk。
type noopAiChatPublisher struct{}

func (noopAiChatPublisher) Publish(_ context.Context, _ sse.StreamID, _ *sse.Event) {}

// AiChatService AI 对话（BFF 编排）：消费 core 的 ChatStreamFrame 流——
// chunk/tool 帧转 protojson 经 SSE 推给发起用户（按其全部 accessToken 流），
// done 帧作为同步 POST 响应返回。
type AiChatService struct {
	adminV1.AiChatServiceHTTPServer

	aiChatServiceClient         aiV1.AiChatServiceClient
	authenticationServiceClient authenticationV1.AuthenticationServiceClient
	log                         *log.Helper

	aiChatPublisher AiChatPublisher
}

func NewAiChatService(
	ctx *bootstrap.Context,
	aiChatServiceClient aiV1.AiChatServiceClient,
	authenticationServiceClient authenticationV1.AuthenticationServiceClient,
) *AiChatService {
	return &AiChatService{
		log:                         ctx.NewLoggerHelper("ai_chat/service/admin-service"),
		aiChatServiceClient:         aiChatServiceClient,
		authenticationServiceClient: authenticationServiceClient,
		aiChatPublisher:             noopAiChatPublisher{},
	}
}

// RegisterAiChatPublisher SSE 启动时注入真身（与 RegisterInternalMessagePublisher 同模式）。
func (s *AiChatService) RegisterAiChatPublisher(publisher AiChatPublisher) {
	if publisher == nil {
		return
	}
	s.aiChatPublisher = publisher
}

// Chat 一轮对话。
func (s *AiChatService) Chat(ctx context.Context, req *aiV1.ChatRequest) (*aiV1.ChatResponse, error) {
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	// 发起用户的全部 SSE 流（同一账号多标签页各持一个 accessToken 流）；
	// 拿不到就降级为无 chunk 对话（同步响应兜底）。
	streamIds, err := s.authenticationServiceClient.GetAccessTokens(ctx, &authenticationV1.GetAccessTokensRequest{
		UserId:     operator.UserId,
		ClientType: authenticationV1.ClientType_admin,
	})
	if err != nil {
		s.log.Warnf("get user access tokens failed, degrade to no-chunk chat: %v", err)
		streamIds = nil
	}

	// 流式 gRPC 调用不走 unary 客户端中间件，operator 元数据不会自动随请求转发，
	// 显式把操作者身份写进请求体（core 信任内部跳，前端传入值在此被覆盖）。
	req.OperatorUserId = trans.Ptr(operator.UserId)
	req.OperatorTenantId = trans.Ptr(operator.GetTenantId())

	stream, err := s.aiChatServiceClient.Chat(ctx, req)
	if err != nil {
		return nil, err
	}

	for {
		frame, recvErr := stream.Recv()
		if recvErr != nil {
			if recvErr == io.EOF {
				return nil, aiV1.ErrorInternalServerError("ai chat stream ended without done frame")
			}
			return nil, recvErr
		}
		switch f := frame.Frame.(type) {
		case *aiV1.ChatStreamFrame_Chunk:
			s.publishToUserStreams(ctx, s.aiChatPublisher, streamIds, sseEventAIChatChunk, f.Chunk)
		case *aiV1.ChatStreamFrame_Tool:
			s.publishToUserStreams(ctx, s.aiChatPublisher, streamIds, sseEventAIChatTool, f.Tool)
		case *aiV1.ChatStreamFrame_Done:
			return f.Done, nil
		}
	}
}

// publishToUserStreams 向用户全部流推送一帧事件（protojson、camelCase 键）。
func (s *AiChatService) publishToUserStreams(ctx context.Context, publisher AiChatPublisher, streamIds *authenticationV1.GetAccessTokensResponse, eventName string, msg proto.Message) {
	if streamIds == nil || len(streamIds.AccessTokens) == 0 {
		return
	}
	data, err := protojson.Marshal(msg)
	if err != nil {
		s.log.Errorf("marshal %s event failed: %v", eventName, err)
		return
	}
	for _, streamId := range streamIds.AccessTokens {
		publisher.Publish(ctx, sse.StreamID(streamId), &sse.Event{
			Event: []byte(eventName),
			Data:  data,
		})
	}
}
