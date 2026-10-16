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
	"github.com/robfig/cron/v3"
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
	Seconds       int      `json:"seconds"` // delay 动作：延时秒数
}

// flowScheduleEvent 定时触发规则的事件名（无实体主体，entity_id=0）。
const flowScheduleEvent = "schedule"

// maxFlowDelaySeconds 延时动作上限（7 天——超过该量级应走外部调度）。
const maxFlowDelaySeconds = 7 * 24 * 3600

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
	if err := s.validateTriggerCron(req.Data); err != nil {
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
	if err := s.validateTriggerCron(req.Data); err != nil {
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
		case "delay":
			if a.Seconds < 1 || a.Seconds > maxFlowDelaySeconds {
				return errors.BadRequest("FLOW_RULE", fmt.Sprintf("actions[%d] delay seconds must be in [1,%d]", i, maxFlowDelaySeconds))
			}
		default:
			return errors.BadRequest("FLOW_RULE", fmt.Sprintf("actions[%d] unsupported type %q", i, a.Type))
		}
	}
	return nil
}

// validateTriggerCron 定时触发表达式校验：非空时必须是合法 5 段 cron
//（分 时 日 月 周）。定时规则忽略事件触发与状态条件（写入时即提示语义）。
func (s *FlowService) validateTriggerCron(dto *flowV1.FlowRule) error {
	cronExpr := strings.TrimSpace(dto.GetTriggerCron())
	if cronExpr == "" {
		return nil
	}
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	if _, err := parser.Parse(cronExpr); err != nil {
		return errors.BadRequest("FLOW_RULE", fmt.Sprintf("triggerCron is not a valid 5-field cron expression: %v", err))
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
// RuleID 为 0 表示按事件扇出（DispatchFlow 的入参形态）；ActionIndex > 0 表示
// 延时续跑（flow.resume 任务：跳过条件判定，从该下标继续动作链）。
// 该 handler 同时注册在 flow.execute 与 flow.resume 两个任务类型上。
func (s *FlowService) ExecuteWorker(taskType string, payload *task.FlowExecutePayload) error {
	if payload == nil || payload.Event == "" {
		return nil
	}
	ctx := appViewer.NewSystemViewerContext(context.Background())

	if payload.ActionIndex > 0 {
		return s.resumeExecute(ctx, payload)
	}

	// 定时触发（scheduleTickWorker 直派：RuleID=规则自身 + Event=schedule）：
	// 直接加载规则执行，不走事件扇出——定时规则已被 ListEnabledForEvent 的
	// TriggerCronIsNil 过滤排除，且其 trigger_event 是原标签而非 schedule。
	if payload.Event == flowScheduleEvent {
		return s.executeScheduled(ctx, payload)
	}

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
// 定时触发（schedule）无实体主体，跳过条件。
func (s *FlowService) evaluateCondition(ctx context.Context, rule *flowV1.FlowRule, payload *task.FlowExecutePayload) (bool, error) {
	if payload.Event == flowScheduleEvent {
		return true, nil
	}
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
	s.executeActionsFrom(ctx, rule, payload, 0)
}

// executeActionsFrom 从 start 下标执行动作链；遇 delay 动作时把剩余链
//（下标 i+1 起）交给 flow.resume 延迟任务后返回。resume 入队失败则放弃
// 剩余动作（已执行部分不回滚——动作语义是旁路通知/外呼，无事务性）。
func (s *FlowService) executeActionsFrom(ctx context.Context, rule *flowV1.FlowRule, payload *task.FlowExecutePayload, start int) {
	var actions []flowAction
	if err := json.Unmarshal([]byte(rule.GetActions()), &actions); err != nil {
		s.log.Errorf("flow execute: rule [%d] actions is not valid JSON: %v", rule.GetId(), err)
		return
	}
	for i := start; i < len(actions); i++ {
		a := actions[i]
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
		case "delay":
			secs := a.Seconds
			if secs < 1 {
				secs = 1
			}
			if secs > maxFlowDelaySeconds {
				secs = maxFlowDelaySeconds
			}
			if s.taskService == nil {
				s.log.Errorf("flow execute: rule [%d] action[%d] delay skipped: task service unavailable", rule.GetId(), i)
				return
			}
			cont := *payload
			cont.ActionIndex = i + 1
			if err := s.taskService.EnqueueFlowResume(&cont, time.Duration(secs)*time.Second); err != nil {
				s.log.Errorf("flow execute: rule [%d] action[%d] delay enqueue failed, remaining actions dropped: %v",
					rule.GetId(), i, err)
				return
			}
			s.log.Infof("flow execute: rule [%d] action[%d] delay %ds, resume at action[%d]",
				rule.GetId(), i, secs, i+1)
			return
		default:
			s.log.Warnf("flow execute: rule [%d] action[%d] unsupported type %q, skip", rule.GetId(), i, a.Type)
		}
	}
}

// resumeExecute 延时续跑：加载规则（已删/停用则放弃）→ 跳过条件 → 从
// ActionIndex 继续动作链。命中计数不重复累加（触发时已计）。
func (s *FlowService) resumeExecute(ctx context.Context, payload *task.FlowExecutePayload) error {
	rule, err := s.flowRepo.Get(ctx, payload.RuleID)
	if err != nil {
		if errors.IsNotFound(err) {
			s.log.Infof("flow resume: rule [%d] deleted, drop continuation", payload.RuleID)
			return nil
		}
		return err // 瞬时查询失败：交 asynq 退避重试
	}
	if !rule.GetEnabled() {
		s.log.Infof("flow resume: rule [%d] disabled, drop continuation", payload.RuleID)
		return nil
	}
	s.executeActionsFrom(ctx, rule, payload, payload.ActionIndex)
	return nil
}

// executeScheduled 定时规则直派执行：加载规则 →（无实体，跳过条件）→ 全量动作链。
func (s *FlowService) executeScheduled(ctx context.Context, payload *task.FlowExecutePayload) error {
	if payload.RuleID == 0 {
		return nil
	}
	rule, err := s.flowRepo.Get(ctx, payload.RuleID)
	if err != nil {
		if errors.IsNotFound(err) {
			s.log.Infof("flow schedule: rule [%d] deleted, drop fire", payload.RuleID)
			return nil
		}
		return err
	}
	if !rule.GetEnabled() {
		s.log.Infof("flow schedule: rule [%d] disabled, drop fire", payload.RuleID)
		return nil
	}
	rulePayload := *payload
	rulePayload.RuleID = rule.GetId()
	s.executeActions(ctx, rule, &rulePayload)
	_ = s.flowRepo.IncrMatchCount(ctx, rule.GetId())
	return nil
}

// ScheduleTickWorker 定时规则调度（flow.schedule.tick，每分钟）：解析各启用
// 规则的 cron，命中当前分钟的按其租户入队 flow.execute（Event 固定 schedule、
// EntityID=0、RuleID=规则自身——不走扇出）。tick 丢失的分钟不补跑（与 cron 语义
// 一致）；多实例部署会重复触发，本服务按单调度实例部署。
func (s *FlowService) ScheduleTickWorker(taskType string, payload *task.FlowScheduleTickPayload) error {
	ctx := appViewer.NewSystemViewerContext(context.Background())
	rules, err := s.flowRepo.ListEnabledCronRules(ctx)
	if err != nil {
		s.log.Errorf("flow schedule tick: list cron rules failed: %v", err)
		return err
	}
	if len(rules) == 0 {
		return nil
	}

	now := time.Now()
	minute := now.Truncate(time.Minute)
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	for _, r := range rules {
		sched, pErr := parser.Parse(r.Cron)
		if pErr != nil {
			s.log.Errorf("flow schedule tick: rule [%d] bad cron %q: %v", r.ID, r.Cron, pErr)
			continue
		}
		if !sched.Next(minute.Add(-time.Second)).Equal(minute) {
			continue
		}
		if s.taskService == nil {
			return nil
		}
		if err := s.taskService.EnqueueFlowExecute(&task.FlowExecutePayload{
			RuleID:   r.ID,
			Event:    flowScheduleEvent,
			TenantID: r.TenantID,
		}); err != nil {
			s.log.Errorf("flow schedule tick: enqueue rule [%d] failed: %v", r.ID, err)
			continue
		}
		s.log.Infof("flow schedule tick: rule [%d] cron=%q fired", r.ID, r.Cron)
	}
	return nil
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
