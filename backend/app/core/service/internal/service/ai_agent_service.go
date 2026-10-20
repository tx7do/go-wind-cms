package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/v9"
	"github.com/sashabaranov/go-openai"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/go-crud/viewer"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/emptypb"

	aiV1 "go-wind-cms/api/gen/go/ai/service/v1"

	"go-wind-cms/app/core/service/internal/data"
	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/aiagentrun"
	"go-wind-cms/pkg/sseevent"
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
	providerRepo  *data.AiProviderRepo
	knowledgeRepo *data.AiKnowledgeRepo
	usageLogRepo  *data.AiUsageLogRepo
	rdb           *redis.Client
}

func NewAiAgentService(
	ctx *bootstrap.Context,
	appRepo *data.AiAgentAppRepo,
	runRepo *data.AiAgentRunRepo,
	providerRepo *data.AiProviderRepo,
	knowledgeRepo *data.AiKnowledgeRepo,
	usageLogRepo *data.AiUsageLogRepo,
	rdb *redis.Client,
) *AiAgentService {
	return &AiAgentService{
		log:           ctx.NewLoggerHelper("ai-agent/service/core-service"),
		appRepo:       appRepo,
		runRepo:       runRepo,
		providerRepo:  providerRepo,
		knowledgeRepo: knowledgeRepo,
		usageLogRepo:  usageLogRepo,
		rdb:           rdb,
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
	if err := s.validateGraphOrBadRequest(req.Data.GetGraph()); err != nil {
		return nil, err
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
	if err := s.validateGraphOrBadRequest(req.Data.GetGraph()); err != nil {
		return nil, err
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

// ─── 试运行（图执行引擎，切片②） ───

// validateGraphOrBadRequest 图非空时必须能解析且过静态校验：存一张画布都画不出的图，
// 试运行时才炸等于把配置错误藏到使用时刻。
func (s *AiAgentService) validateGraphOrBadRequest(raw string) error {
	if raw == "" {
		return nil // 允许先建壳后画图
	}
	g, err := parseAgentGraph(raw)
	if err != nil {
		return errors.BadRequest("AI_AGENT", err.Error())
	}
	if err = validateAgentGraph(g); err != nil {
		return errors.BadRequest("AI_AGENT", err.Error())
	}
	return nil
}

// Run 试运行：校验（不过则 400、不建行）→ 建 RUNNING 行 → 同步执行 → 回写终态。
// 逐节点轨迹帧经 Redis Pub/Sub 发给 admin 服务转 SSE（发起用户的全部在线流）；
// 丢帧可容忍——同步响应携带完整 trace 与终态，以响应为准。
// 与上游的差异：上游 monolith 直接推 SSE；CMS 的 SSE 在 admin，core 经既有
// 「core 发 Redis、admin 转推」模式（系统通知同款）。
func (s *AiAgentService) Run(ctx context.Context, req *aiV1.RunAiAgentAppRequest) (*aiV1.AiAgentRun, error) {
	if req == nil || req.GetId() == 0 {
		return nil, errors.BadRequest("AI_AGENT", "id is required")
	}

	var operatorID, tenantID uint32
	if vc, exist := viewer.FromContext(ctx); exist && vc != nil {
		operatorID = uint32(vc.UserID())
		tenantID = uint32(vc.TenantID())
	}

	app, err := s.appRepo.GetEntityByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	if !app.Enabled {
		return nil, errors.BadRequest("AI_AGENT", "ai agent app is disabled")
	}

	g, err := parseAgentGraph(derefAgentStr(app.Graph))
	if err != nil {
		return nil, errors.BadRequest("AI_AGENT", err.Error())
	}
	if err = validateAgentGraph(g); err != nil {
		return nil, errors.BadRequest("AI_AGENT", err.Error())
	}

	// AI_TOKENS 配额：CMS 的租户配额（tenant-quota 批）对该维度只做水位告警、
	// 不做入口硬闸——与上游 checkAiTokenQuota 的前置拒绝有意不同，水位兜底。

	inputJSON, _ := json.Marshal(req.GetInput())
	runID, err := s.runRepo.CreateRunning(ctx, app.ID, string(inputJSON), operatorID, tenantID)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	seq := uint32(0)
	publishFrame := func(phase, nodeID, summary, status string) {
		s.publishAgentFrame(ctx, operatorID, &aiV1.AgentRunEvent{
			RunId:   runID,
			AppId:   app.ID,
			Seq:     seq,
			Phase:   phase,
			NodeId:  nodeID,
			Summary: summary,
			Status:  status,
		})
		seq++
	}

	output, trace, usage, execErr := executeAgentGraph(ctx, g, req.GetInput(), &agentEngineDeps{
		searchKnowledge: s.searchKnowledgeForAgent(operatorID, tenantID),
		runLLM:          s.runLLMForAgent(operatorID, tenantID),
		onStep:          publishFrame,
	})
	durationMs := uint32(time.Since(start).Milliseconds())

	traceJSON, _ := json.Marshal(trace)

	status := aiagentrun.StatusSucceeded
	errMsg := ""
	if execErr != nil {
		status = aiagentrun.StatusFailed
		errMsg = execErr.Error()
		publishFrame("error", "", errMsg, "")
		s.log.Errorf("ai agent run [%d] app [%d] failed: %s", runID, app.ID, execErr.Error())
	} else {
		publishFrame("done", "", output, "")
	}

	if err = s.runRepo.SettleResult(ctx, runID, status, output, errMsg, string(traceJSON),
		data.AgentRunTokens{Prompt: usage.Prompt, Completion: usage.Completion, Total: usage.Total},
		durationMs); err != nil {
		return nil, err
	}

	s.log.Infof("ai agent run [%d] app [%d] %s in %dms (tokens %d) by user [%d]",
		runID, app.ID, status.String(), durationMs, usage.Total, operatorID)

	return s.runRepo.Get(ctx, runID)
}

// publishAgentFrame 把一帧轨迹事件发往 Redis（admin 订阅后按发起用户转 SSE）。
// 尽力而为：Redis 不可用/发布失败只留日志，不影响执行与同步响应。
func (s *AiAgentService) publishAgentFrame(ctx context.Context, userID uint32, frame *aiV1.AgentRunEvent) {
	if s.rdb == nil {
		return
	}
	data, err := protojson.Marshal(frame)
	if err != nil {
		s.log.Errorf("marshal agent run event failed, skip push: %v", err)
		return
	}
	payload, _ := json.Marshal(sseevent.AgentRunEvent{UserID: userID, Frame: data})
	if err := s.rdb.Publish(ctx, sseevent.RedisChannelAgentRun, payload).Err(); err != nil {
		s.log.Warnf("publish agent run frame failed: %s", err.Error())
	}
}

// searchKnowledgeForAgent 知识检索节点的真实现：复用 chat 同一条 RAG 链路
// （embed + 余弦检索 + embed 用量记账），命中片段按序拼接为上下文。
func (s *AiAgentService) searchKnowledgeForAgent(userID, tenantID uint32) func(ctx context.Context, node *agentGraphNode, query string) (string, int, error) {
	return func(ctx context.Context, node *agentGraphNode, query string) (string, int, error) {
		baseID := agentDataUint32(node.Data, "baseId")
		if baseID == 0 {
			return "", 0, fmt.Errorf("knowledge node has no baseId configured")
		}
		topK := agentDataInt(node.Data, "topK", 4)
		hits, err := searchKnowledgeBase(ctx, s.knowledgeRepo, s.providerRepo, s.usageLogRepo, s.log, userID, tenantID, baseID, query, topK)
		if err != nil {
			return "", 0, err
		}
		var sb strings.Builder
		for i, h := range hits {
			if i > 0 {
				sb.WriteString("\n\n")
			}
			sb.WriteString(h.Content)
		}
		return sb.String(), len(hits), nil
	}
}

// runLLMForAgent LLM 节点的真实现：provider 解析（节点点名，留空取默认启用）→
// 客户端构造（密文解密就地）→ aiToolLoop 多轮工具调用 → 逐节点落用量流水
//（conversation_id=0 的非对话调用，与 chat 共用 AI_TOKENS 配额口径）。
func (s *AiAgentService) runLLMForAgent(userID, tenantID uint32) func(ctx context.Context, node *agentGraphNode, system, user string) (string, agentUsage, error) {
	return func(ctx context.Context, node *agentGraphNode, system, user string) (string, agentUsage, error) {
		var provider *ent.AiProvider
		var err error
		if providerID := agentDataUint32(node.Data, "providerId"); providerID > 0 {
			provider, err = s.providerRepo.GetEntityByID(ctx, providerID)
			if err != nil {
				return "", agentUsage{}, err
			}
			if provider.IsEnabled != nil && !*provider.IsEnabled {
				return "", agentUsage{}, fmt.Errorf("ai provider [%d] is disabled", providerID)
			}
		} else {
			provider, err = s.providerRepo.GetEnabledDefault(ctx)
			if err != nil {
				return "", agentUsage{}, err
			}
		}

		client, err := newOpenAIClientForProvider(provider)
		if err != nil {
			return "", agentUsage{}, err
		}

		model := agentDataString(node.Data, "model")
		if model == "" {
			model = derefAgentStr(provider.ModelName)
		}

		var temperature *float32
		if t, ok := agentDataFloat(node.Data, "temperature"); ok {
			t32 := float32(t)
			temperature = &t32
		}

		var messages []openai.ChatCompletionMessage
		if system != "" {
			messages = append(messages, openai.ChatCompletionMessage{Role: openai.ChatMessageRoleSystem, Content: system})
		}
		messages = append(messages, openai.ChatCompletionMessage{Role: openai.ChatMessageRoleUser, Content: user})

		loop := &aiToolLoop{
			client:      client,
			model:       model,
			tools:       aiBuiltinToolDefs(),
			exec:        execAiTool,
			maxRounds:   aiToolMaxRounds,
			temperature: temperature,
		}

		start := time.Now()
		answer, usage, err := loop.run(ctx, messages)
		durationMs := uint32(time.Since(start).Milliseconds())

		// 逐 LLM 节点记一条用量流水：agent 运行在用量台账里是"非对话调用"，
		// model_name 记节点实际用的模型，水位扫描与 AI_TOKENS 配额同源。
		if usage.TotalTokens > 0 {
			if uerr := s.usageLogRepo.Create(ctx, &data.UsageRecord{
				ProviderID:       provider.ID,
				UserID:           userID,
				TenantID:         tenantID,
				ModelName:        model,
				PromptTokens:     uint32(usage.PromptTokens),
				CompletionTokens: uint32(usage.CompletionTokens),
				TotalTokens:      uint32(usage.TotalTokens),
				DurationMs:       durationMs,
			}); uerr != nil {
				// 记账失败不回滚答案：运行本身已成功，台账缺一行可补，答案丢了不可补。
				s.log.Errorf("ai agent run: write usage log failed (model %s): %s", model, uerr.Error())
			}
		}

		if err != nil {
			return "", agentUsage{}, err
		}
		return answer, agentUsage{
			Prompt:     uint32(usage.PromptTokens),
			Completion: uint32(usage.CompletionTokens),
			Total:      uint32(usage.TotalTokens),
		}, nil
	}
}

func derefAgentStr(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}
