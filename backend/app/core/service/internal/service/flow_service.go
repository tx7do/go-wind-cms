package service

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/go-crud/viewer"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"google.golang.org/protobuf/types/known/emptypb"

	appViewer "go-wind-cms/pkg/entgo/viewer"
	"go-wind-cms/pkg/task"

	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
	flowV1 "go-wind-cms/api/gen/go/flow/service/v1"

	"go-wind-cms/app/core/service/internal/data"
)

// flowAction 流动作（actions JSON 数组元素）。
//   - webhook：{"type":"webhook","url":"https://...","secret":"可选签名密钥"}
//   - notification：{"type":"notification","title":"...","content":"...",
//     "target_user_ids":[..]}（title/content 支持 {event}/{entity_id}/{rule_id} 占位符）
type flowAction struct {
	Type          string   `json:"type"`
	URL           string   `json:"url"`
	Secret        string   `json:"secret"`
	Title         string   `json:"title"`
	Content       string   `json:"content"`
	TargetUserIDs []uint32 `json:"target_user_ids"`
}

// FlowService 自动化 Flow 规则服务。
type FlowService struct {
	flowV1.UnimplementedFlowRuleServiceServer

	flowRepo           *data.FlowRepo
	postRepo           *data.PostRepo
	pageRepo           *data.PageRepo
	taskService        *TaskService
	internalMessageSvc *InternalMessageService

	httpClient *http.Client
	log        *log.Helper
}

func NewFlowService(ctx *bootstrap.Context, flowRepo *data.FlowRepo, postRepo *data.PostRepo, pageRepo *data.PageRepo, taskService *TaskService, internalMessageSvc *InternalMessageService) *FlowService {
	return &FlowService{
		log:                ctx.NewLoggerHelper("flow/service/core-service"),
		flowRepo:           flowRepo,
		postRepo:           postRepo,
		pageRepo:           pageRepo,
		taskService:        taskService,
		internalMessageSvc: internalMessageSvc,
		httpClient:         &http.Client{Timeout: 10 * time.Second},
	}
}

// ─── CRUD（admin 门面透传目标）───

func (s *FlowService) List(ctx context.Context, req *paginationV1.PagingRequest) (*flowV1.ListFlowRuleResponse, error) {
	return s.flowRepo.List(ctx, req)
}

func (s *FlowService) Get(ctx context.Context, req *flowV1.GetFlowRuleRequest) (*flowV1.FlowRule, error) {
	if req == nil || req.GetId() == 0 {
		return nil, errors.BadRequest("FLOW_RULE", "invalid parameter")
	}
	return s.flowRepo.Get(ctx, req.GetId())
}

func (s *FlowService) Create(ctx context.Context, req *flowV1.CreateFlowRuleRequest) (*flowV1.FlowRule, error) {
	if req == nil || req.Data == nil {
		return nil, errors.BadRequest("FLOW_RULE", "invalid parameter")
	}
	if err := s.validateActions(req.Data.GetActions()); err != nil {
		return nil, err
	}
	return s.flowRepo.Create(ctx, req.Data)
}

func (s *FlowService) Update(ctx context.Context, req *flowV1.UpdateFlowRuleRequest) (*flowV1.FlowRule, error) {
	if req == nil || req.Data == nil || req.GetId() == 0 {
		return nil, errors.BadRequest("FLOW_RULE", "invalid parameter")
	}
	if err := s.validateActions(req.Data.GetActions()); err != nil {
		return nil, err
	}
	return s.flowRepo.Update(ctx, req.GetId(), req.Data)
}

func (s *FlowService) Delete(ctx context.Context, req *flowV1.DeleteFlowRuleRequest) (*emptypb.Empty, error) {
	if req == nil || req.GetId() == 0 {
		return nil, errors.BadRequest("FLOW_RULE", "invalid parameter")
	}
	if err := s.flowRepo.Delete(ctx, req.GetId()); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// validateActions 动作链 JSON 结构校验（创建/更新时把关，worker 侧容错跳过坏动作）。
func (s *FlowService) validateActions(actionsJSON string) error {
	if strings.TrimSpace(actionsJSON) == "" {
		return errors.BadRequest("FLOW_RULE", "actions is required")
	}
	var actions []flowAction
	if err := json.Unmarshal([]byte(actionsJSON), &actions); err != nil {
		return errors.BadRequest("FLOW_RULE", "actions is not valid JSON array")
	}
	for i, a := range actions {
		switch a.Type {
		case "webhook":
			if a.URL == "" {
				return errors.BadRequest("FLOW_RULE", fmt.Sprintf("actions[%d] webhook url is required", i))
			}
		case "notification":
			if strings.TrimSpace(a.Title) == "" || strings.TrimSpace(a.Content) == "" {
				return errors.BadRequest("FLOW_RULE", fmt.Sprintf("actions[%d] notification title/content are required", i))
			}
			if len(a.TargetUserIDs) == 0 {
				return errors.BadRequest("FLOW_RULE", fmt.Sprintf("actions[%d] notification target_user_ids are required", i))
			}
		default:
			return errors.BadRequest("FLOW_RULE", fmt.Sprintf("actions[%d] unsupported type %q", i, a.Type))
		}
	}
	return nil
}

// ─── 事件触发 ───

// DispatchFlow 内容事件触发流规则（内容写路径钩子，best-effort：入队失败仅记日志）。
func (s *FlowService) DispatchFlow(ctx context.Context, event string, entityID uint32) {
	if s.taskService == nil || entityID == 0 {
		return
	}
	var tenantID uint32
	if vc, exist := viewer.FromContext(ctx); exist && vc != nil {
		tenantID = uint32(vc.TenantID())
	}
	if err := s.taskService.EnqueueFlowExecute(&task.FlowExecutePayload{
		RuleID:   0, // 0=按事件扇出所有匹配规则
		Event:    event,
		EntityID: entityID,
		TenantID: tenantID,
	}); err != nil {
		s.log.Errorf("dispatch flow event failed (event=%s entity=%d): %v", event, entityID, err)
	}
}

// ─── asynq Worker ───

// ExecuteWorker 规则执行：加载事件下的启用规则 → 状态条件判定 → 逐动作执行。
// RuleID 为 0 表示按事件扇出（DispatchFlow 的入参形态）。
func (s *FlowService) ExecuteWorker(taskType string, payload *task.FlowExecutePayload) error {
	if payload == nil || payload.Event == "" {
		return nil
	}
	ctx := appViewer.NewSystemViewerContext(context.Background())

	rules, err := s.flowRepo.ListEnabledForEvent(ctx, payload.TenantID, payload.Event)
	if err != nil {
		s.log.Errorf("flow execute: list rules failed (event=%s): %v", payload.Event, err)
		return err // 查库失败：可重试
	}
	s.log.Infof("flow execute: event=%s entity=%d tenant=%d rules=%d",
		payload.Event, payload.EntityID, payload.TenantID, len(rules))

	for _, rule := range rules {
		if payload.RuleID != 0 && rule.GetId() != payload.RuleID {
			continue
		}
		matched, mErr := s.evaluateCondition(ctx, rule, payload)
		if mErr != nil {
			s.log.Errorf("flow execute: evaluate rule [%d] failed: %v", rule.GetId(), mErr)
			continue
		}
		if !matched {
			continue
		}
		// envelope 中的 rule_id 以实际命中的规则为准（入队任务 RuleID=0 表示按事件扇出）
		rulePayload := *payload
		rulePayload.RuleID = rule.GetId()
		s.executeActions(ctx, rule, &rulePayload)
		// 命中计数（观测用，best-effort）
		_ = s.flowRepo.IncrMatchCount(ctx, rule.GetId())
	}
	return nil
}

// evaluateCondition 条件判定。v1 仅支持状态等值（condition_status 空=无条件）。
func (s *FlowService) evaluateCondition(ctx context.Context, rule *flowV1.FlowRule, payload *task.FlowExecutePayload) (bool, error) {
	want := strings.TrimSpace(rule.GetConditionStatus())
	if want == "" {
		return true, nil
	}
	got, err := s.entityStatus(ctx, payload.Event, payload.EntityID)
	if err != nil {
		return false, err
	}
	return got == want, nil
}

// entityStatus 按事件族加载实体当前状态（供条件判定）。
func (s *FlowService) entityStatus(ctx context.Context, event string, entityID uint32) (string, error) {
	if strings.HasPrefix(event, "post.") {
		post, err := s.postRepo.Get(ctx, &contentV1.GetPostRequest{
			QueryBy: &contentV1.GetPostRequest_Id{Id: entityID},
		})
		if err != nil {
			return "", err
		}
		return post.GetStatus().String(), nil
	}
	if strings.HasPrefix(event, "page.") {
		page, err := s.pageRepo.Get(ctx, &contentV1.GetPageRequest{
			QueryBy: &contentV1.GetPageRequest_Id{Id: entityID},
		})
		if err != nil {
			return "", err
		}
		return page.GetStatus().String(), nil
	}
	return "", nil
}

// executeActions 逐动作执行；单动作失败记日志不中断其余动作。
func (s *FlowService) executeActions(ctx context.Context, rule *flowV1.FlowRule, payload *task.FlowExecutePayload) {
	var actions []flowAction
	if err := json.Unmarshal([]byte(rule.GetActions()), &actions); err != nil {
		s.log.Errorf("flow execute: rule [%d] actions is not valid JSON: %v", rule.GetId(), err)
		return
	}
	for i, a := range actions {
		switch a.Type {
		case "webhook":
			ok, httpStatus, errMsg := s.postWebhook(a, payload)
			if !ok {
				s.log.Errorf("flow execute: rule [%d] action[%d] webhook failed (http=%d): %s",
					rule.GetId(), i, httpStatus, errMsg)
			}
		case "notification":
			if s.internalMessageSvc == nil {
				s.log.Errorf("flow execute: rule [%d] action[%d] notification skipped: message service unavailable", rule.GetId(), i)
				continue
			}
			title := s.expandPlaceholders(a.Title, rule, payload)
			content := s.expandPlaceholders(a.Content, rule, payload)
			if _, nErr := s.internalMessageSvc.SendSystemNotification(ctx, title, content, a.TargetUserIDs); nErr != nil {
				s.log.Errorf("flow execute: rule [%d] action[%d] notification failed: %v", rule.GetId(), i, nErr)
			}
		default:
			s.log.Warnf("flow execute: rule [%d] action[%d] unsupported type %q, skip", rule.GetId(), i, a.Type)
		}
	}
}

// expandPlaceholders 展开通知模板占位符：{event}/{entity_id}/{rule_id}。
func (s *FlowService) expandPlaceholders(tpl string, rule *flowV1.FlowRule, payload *task.FlowExecutePayload) string {
	out := strings.ReplaceAll(tpl, "{event}", payload.Event)
	out = strings.ReplaceAll(out, "{entity_id}", strconv.FormatUint(uint64(payload.EntityID), 10))
	out = strings.ReplaceAll(out, "{rule_id}", strconv.FormatUint(uint64(rule.GetId()), 10))
	return out
}

// postWebhook 签名投递（与 webhook 领域同款 envelope：事件 + 实体 id + 规则 id，不含内容全文）。
func (s *FlowService) postWebhook(a flowAction, payload *task.FlowExecutePayload) (bool, int, string) {
	body, err := json.Marshal(map[string]any{
		"event":     payload.Event,
		"entity_id": payload.EntityID,
		"rule_id":   payload.RuleID,
		"tenant_id": payload.TenantID,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
	if err != nil {
		return false, 0, err.Error()
	}

	req, err := http.NewRequest(http.MethodPost, a.URL, bytes.NewReader(body))
	if err != nil {
		return false, 0, "invalid url: " + err.Error()
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "GoWind-Content-Hub-Flow/1.0")
	req.Header.Set("X-Flow-Event", payload.Event)
	if a.Secret != "" {
		mac := hmac.New(sha256.New, []byte(a.Secret))
		mac.Write(body)
		req.Header.Set("X-Webhook-Signature", "sha256="+hex.EncodeToString(mac.Sum(nil)))
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return false, 0, err.Error()
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return false, resp.StatusCode, fmt.Sprintf("unexpected status %d", resp.StatusCode)
	}
	return true, resp.StatusCode, ""
}
