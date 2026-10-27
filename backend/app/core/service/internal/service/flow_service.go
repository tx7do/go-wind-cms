package service

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/robfig/cron/v3"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/go-crud/viewer"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	appViewer "go-wind-cms/pkg/entgo/viewer"
	"go-wind-cms/pkg/task"

	commentV1 "go-wind-cms/api/gen/go/comment/service/v1"
	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
	flowV1 "go-wind-cms/api/gen/go/flow/service/v1"
	identityV1 "go-wind-cms/api/gen/go/identity/service/v1"
	mediaV1 "go-wind-cms/api/gen/go/media/service/v1"

	"go-wind-cms/app/core/service/internal/data"
)

// flowAction 流动作（actions JSON 数组元素）。
//   - webhook：{"type":"webhook","url":"https://...","secret":"可选签名密钥"}
//   - notification：{"type":"notification","title":"...","content":"...",
//     "target_user_ids":[..]}（title/content 支持 {event}/{entity_id}/{rule_id} 占位符）
type flowAction struct {
	Type    string `json:"type"`
	URL     string `json:"url"`
	Secret  string `json:"secret"`
	Title   string `json:"title"`
	Content string `json:"content"`
	// 英文变体（模板 lite）：可选。两者齐备时，偏好语言为英文的收件人的
	// 邮件触达按 EN 文案渲染（站内信仍共享主文案单行）；任一为空即忽略变体。
	TitleEn       string   `json:"title_en,omitempty"`
	ContentEn     string   `json:"content_en,omitempty"`
	TargetUserIDs []uint32 `json:"target_user_ids"`
	Seconds       int      `json:"seconds"` // delay 动作：延时秒数
	// http 动作：自定义 HTTP 外呼（webhook 动作的泛化形态）
	Method  string            `json:"method"`  // GET/POST（缺省 POST）
	Headers map[string]string `json:"headers"` // 自定义请求头（≤8 个）
	Body    string            `json:"body"`    // body 模板，支持 {field} 快照占位符；GET 忽略
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
	if err := s.validateOptionalActions(req.Data.GetElseActions()); err != nil {
		return nil, err
	}
	if err := s.validateTriggerCron(req.Data); err != nil {
		return nil, err
	}
	if err := s.validateConditions(req.Data.GetConditions()); err != nil {
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
	if err := s.validateOptionalActions(req.Data.GetElseActions()); err != nil {
		return nil, err
	}
	if err := s.validateTriggerCron(req.Data); err != nil {
		return nil, err
	}
	if err := s.validateConditions(req.Data.GetConditions()); err != nil {
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
		case "http":
			if a.URL == "" {
				return errors.BadRequest("FLOW_RULE", fmt.Sprintf("actions[%d] http url is required", i))
			}
			if m := strings.ToUpper(strings.TrimSpace(a.Method)); m != "" && m != "GET" && m != "POST" {
				return errors.BadRequest("FLOW_RULE", fmt.Sprintf("actions[%d] http method must be GET or POST", i))
			}
			if len(a.Headers) > 8 {
				return errors.BadRequest("FLOW_RULE", fmt.Sprintf("actions[%d] http headers exceed 8", i))
			}
			if len(a.Body) > maxFlowHTTPBodyBytes {
				return errors.BadRequest("FLOW_RULE", fmt.Sprintf("actions[%d] http body exceeds 64KB", i))
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

// validateOptionalActions 可空动作链校验（else 分支用）：空串=无分支；
// 非空则与主链同一动作 schema。
func (s *FlowService) validateOptionalActions(actionsJSON string) error {
	if strings.TrimSpace(actionsJSON) == "" {
		return nil
	}
	return s.validateActions(actionsJSON)
}

// validateTriggerCron 定时触发表达式校验：非空时必须是合法 5 段 cron
// （分 时 日 月 周）。定时规则忽略事件触发与状态条件（写入时即提示语义）。
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

// ─── 事件快照构造（值统一字符串；字段名即条件 field 取值域）───

const flowSnapshotContentLimit = 200

// PostEventSnapshot 文章事件快照。
func PostEventSnapshot(post *contentV1.Post, prevStatus string) map[string]string {
	if post == nil {
		return nil
	}
	snap := map[string]string{
		"status":    post.GetStatus().String(),
		"author_id": strconv.FormatUint(uint64(post.GetAuthorId()), 10),
	}
	if prevStatus != "" {
		snap["prev_status"] = prevStatus
	}
	if ids := post.GetCategoryIds(); len(ids) > 0 {
		parts := make([]string, 0, len(ids))
		for _, id := range ids {
			parts = append(parts, strconv.FormatUint(uint64(id), 10))
		}
		snap["category_id"] = strings.Join(parts, "|")
	}
	for _, tr := range post.GetTranslations() {
		if snap["title"] == "" && tr.GetTitle() != "" {
			snap["title"] = tr.GetTitle()
		}
		if snap["slug"] == "" && tr.GetSlug() != "" {
			snap["slug"] = tr.GetSlug()
		}
	}
	for key, value := range post.GetCustomFields() {
		snap["cf."+key] = value
	}
	return snap
}

// PageEventSnapshot 页面事件快照。
func PageEventSnapshot(page *contentV1.Page, prevStatus string) map[string]string {
	if page == nil {
		return nil
	}
	snap := map[string]string{
		"status": page.GetStatus().String(),
	}
	if prevStatus != "" {
		snap["prev_status"] = prevStatus
	}
	for _, tr := range page.GetTranslations() {
		if snap["title"] == "" && tr.GetTitle() != "" {
			snap["title"] = tr.GetTitle()
		}
		if snap["slug"] == "" && tr.GetSlug() != "" {
			snap["slug"] = tr.GetSlug()
		}
	}
	for key, value := range page.GetCustomFields() {
		snap["cf."+key] = value
	}
	return snap
}

// CommentEventSnapshot 评论事件快照（content 截断到 200 字符防 payload 膨胀）。
func CommentEventSnapshot(comment *commentV1.Comment) map[string]string {
	if comment == nil {
		return nil
	}
	content := comment.GetContent()
	if len(content) > flowSnapshotContentLimit {
		content = content[:flowSnapshotContentLimit]
	}
	return map[string]string{
		"status":     comment.GetStatus().String(),
		"content":    content,
		"object_id":  strconv.FormatUint(uint64(comment.GetObjectId()), 10),
		"author_id":  strconv.FormatUint(uint64(comment.GetAuthorId()), 10),
		"authorType": comment.GetAuthorType().String(),
	}
}

// MediaEventSnapshot 媒体资产事件快照。
func MediaEventSnapshot(asset *mediaV1.MediaAsset) map[string]string {
	if asset == nil {
		return nil
	}
	return map[string]string{
		"asset_type": asset.GetType().String(),
		"mime_type":  asset.GetMimeType(),
		"file_name":  asset.GetFilename(),
		"size_bytes": strconv.FormatUint(uint64(asset.GetSize()), 10),
	}
}

// UserEventSnapshot 用户事件快照。
func UserEventSnapshot(user *identityV1.User) map[string]string {
	if user == nil {
		return nil
	}
	return map[string]string{
		"status":   user.GetStatus().String(),
		"username": user.GetUsername(),
		"email":    user.GetEmail(),
	}
}

// DispatchFlow 内容事件触发流规则（内容写路径钩子，best-effort：入队失败仅记日志）。
// snapshot 为触发时刻的实体字段快照（值统一字符串），供字段级条件与占位符渲染；
// 旧调用方可传 nil（条件判定回退 DB 只支持 status）。
func (s *FlowService) DispatchFlow(ctx context.Context, event string, entityID uint32, snapshot map[string]string) {
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
		Snapshot: s.enrichSnapshot(ctx, event, entityID, snapshot),
	}); err != nil {
		s.log.Errorf("dispatch flow event failed (event=%s entity=%d): %v", event, entityID, err)
	}
}

// enrichSnapshot 快照补全：触发点 DTO 常不带翻译（title/slug 只在 translations
// 里），此处 best-effort 回查补齐（已删除实体回查失败则保留原样——删除事件的
// title 类条件判不命中，status/prev_status 不受影响）。
func (s *FlowService) enrichSnapshot(ctx context.Context, event string, entityID uint32, snapshot map[string]string) map[string]string {
	if snapshot == nil {
		snapshot = map[string]string{}
	}
	if snapshot["title"] != "" && snapshot["slug"] != "" {
		return snapshot
	}
	switch {
	case strings.HasPrefix(event, "post."):
		if post, err := s.postRepo.Get(ctx, &contentV1.GetPostRequest{
			QueryBy: &contentV1.GetPostRequest_Id{Id: entityID},
		}); err == nil {
			for _, tr := range post.GetTranslations() {
				if snapshot["title"] == "" && tr.GetTitle() != "" {
					snapshot["title"] = tr.GetTitle()
				}
				if snapshot["slug"] == "" && tr.GetSlug() != "" {
					snapshot["slug"] = tr.GetSlug()
				}
			}
		}
	case strings.HasPrefix(event, "page."):
		if page, err := s.pageRepo.Get(ctx, &contentV1.GetPageRequest{
			QueryBy: &contentV1.GetPageRequest_Id{Id: entityID},
		}); err == nil {
			for _, tr := range page.GetTranslations() {
				if snapshot["title"] == "" && tr.GetTitle() != "" {
					snapshot["title"] = tr.GetTitle()
				}
				if snapshot["slug"] == "" && tr.GetSlug() != "" {
					snapshot["slug"] = tr.GetSlug()
				}
			}
		}
	}
	return snapshot
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
		// envelope 中的 rule_id 以实际命中的规则为准（入队任务 RuleID=0 表示按事件扇出）
		rulePayload := *payload
		rulePayload.RuleID = rule.GetId()
		if matched {
			s.executeActions(ctx, rule, &rulePayload)
		} else {
			// 条件不满足：走 else 分支（未配置则无事发生）
			s.executeChain(ctx, rule, &rulePayload, flowChainElse)
			continue
		}
		// 命中计数（观测用，best-effort；else 分支不计入"命中"）
		_ = s.flowRepo.IncrMatchCount(ctx, rule.GetId())
	}
	return nil
}

// flowCondition 结构化条件（conditions JSON 数组元素）。
type flowCondition struct {
	Field string `json:"field"` // 快照字段名（status/prev_status/title/slug/author_id/category_id/cf.<key> 等）
	Op    string `json:"op"`    // eq / ne / contains / in（in 的 value 用 | 分隔多个候选）
	Value string `json:"value"`
}

// flowConditionOps 条件算子表（op 大小写不敏感）。
var flowConditionOps = map[string]func(got, want string) bool{
	"eq":       func(g, w string) bool { return g == w },
	"ne":       func(g, w string) bool { return g != w },
	"contains": func(g, w string) bool { return strings.Contains(g, w) },
	"in": func(g, w string) bool {
		for _, part := range strings.Split(w, "|") {
			if strings.TrimSpace(part) == g {
				return true
			}
		}
		return false
	},
}

// validateConditions 结构化条件 JSON 校验（创建/更新时把关）。
func (s *FlowService) validateConditions(conditionsJSON string) error {
	conds, err := parseFlowConditions(conditionsJSON)
	if err != nil {
		return errors.BadRequest("FLOW_RULE", "conditions is not valid JSON array")
	}
	if len(conds) > 16 {
		return errors.BadRequest("FLOW_RULE", "conditions exceeds 16 entries")
	}
	for i, c := range conds {
		if strings.TrimSpace(c.Field) == "" {
			return errors.BadRequest("FLOW_RULE", fmt.Sprintf("conditions[%d] field is required", i))
		}
		if _, ok := flowConditionOps[strings.ToLower(strings.TrimSpace(c.Op))]; !ok {
			return errors.BadRequest("FLOW_RULE", fmt.Sprintf("conditions[%d] unsupported op %q (eq/ne/contains/in)", i, c.Op))
		}
	}
	return nil
}

// parseFlowConditions 解析条件 JSON；空串返回 nil（无条件）。
func parseFlowConditions(conditionsJSON string) ([]flowCondition, error) {
	trimmed := strings.TrimSpace(conditionsJSON)
	if trimmed == "" {
		return nil, nil
	}
	var conds []flowCondition
	if err := json.Unmarshal([]byte(trimmed), &conds); err != nil {
		return nil, err
	}
	return conds, nil
}

// evaluateCondition 条件判定 v2：旧标量状态条件与结构化字段条件并存时
// 两者都须命中（AND）。字段值优先取触发时刻快照（兼容删除事件与执行时
// 状态漂移两类时序问题）；快照缺 status 时回退 DB 回查，缺其余字段判不命中
// （宁缺勿滥）。定时触发（schedule）无实体主体，跳过条件。
func (s *FlowService) evaluateCondition(ctx context.Context, rule *flowV1.FlowRule, payload *task.FlowExecutePayload) (bool, error) {
	if payload.Event == flowScheduleEvent {
		return true, nil
	}

	lookup := func(field string) (string, bool, error) {
		if v, ok := payload.Snapshot[field]; ok {
			return v, true, nil
		}
		if field == "status" {
			got, err := s.entityStatus(ctx, payload.Event, payload.EntityID)
			return got, err == nil, err
		}
		return "", false, nil
	}

	if want := strings.TrimSpace(rule.GetConditionStatus()); want != "" {
		got, ok, err := lookup("status")
		if err != nil {
			return false, err
		}
		if !ok || got != want {
			return false, nil
		}
	}

	conds, err := parseFlowConditions(rule.GetConditions())
	if err != nil {
		return false, err
	}
	for _, c := range conds {
		match := flowConditionOps[strings.ToLower(strings.TrimSpace(c.Op))]
		got, ok, err := lookup(c.Field)
		if err != nil {
			return false, err
		}
		if !ok || !match(got, c.Value) {
			return false, nil
		}
	}
	return true, nil
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

// 动作链名（延时续跑定位用）。
const (
	flowChainMain = "main"
	flowChainElse = "else"
)

// executeActions 逐动作执行主链；单动作失败记日志不中断其余动作。
func (s *FlowService) executeActions(ctx context.Context, rule *flowV1.FlowRule, payload *task.FlowExecutePayload) {
	s.executeActionsFrom(ctx, rule, payload, 0)
}

// executeChain 按链名执行动作链（main=命中分支 / else=条件不满足分支）。
func (s *FlowService) executeChain(ctx context.Context, rule *flowV1.FlowRule, payload *task.FlowExecutePayload, chain string) {
	switch chain {
	case flowChainElse:
		if strings.TrimSpace(rule.GetElseActions()) == "" {
			return
		}
		s.executeActionsJSON(ctx, rule, payload, rule.GetElseActions(), 0, flowChainElse)
	default:
		s.executeActionsFrom(ctx, rule, payload, 0)
	}
}

// executeActionsJSON 解析并执行指定动作链 JSON（else 分支入口）。
func (s *FlowService) executeActionsJSON(ctx context.Context, rule *flowV1.FlowRule, payload *task.FlowExecutePayload, actionsJSON string, start int, chain string) {
	var actions []flowAction
	if err := json.Unmarshal([]byte(actionsJSON), &actions); err != nil {
		s.log.Errorf("flow execute: rule [%d] %s actions is not valid JSON: %v", rule.GetId(), chain, err)
		return
	}
	s.runActionSlice(ctx, rule, payload, actions, start, chain)
}

// executeActionsFrom 从 start 下标执行主链；遇 delay 动作时把剩余链
// （下标 i+1 起）交给 flow.resume 延迟任务后返回。resume 入队失败则放弃
// 剩余动作（已执行部分不回滚——动作语义是旁路通知/外呼，无事务性）。
func (s *FlowService) executeActionsFrom(ctx context.Context, rule *flowV1.FlowRule, payload *task.FlowExecutePayload, start int) {
	s.executeActionsJSON(ctx, rule, payload, rule.GetActions(), start, flowChainMain)
}

// runActionSlice 执行动作切片（主链/else 链共用）；遇 delay 落延迟续跑任务。
func (s *FlowService) runActionSlice(ctx context.Context, rule *flowV1.FlowRule, payload *task.FlowExecutePayload, actions []flowAction, start int, chain string) {
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
			var alt *AltNotificationText
			if a.TitleEn != "" && a.ContentEn != "" {
				alt = &AltNotificationText{
					Title:   s.expandPlaceholders(a.TitleEn, rule, payload),
					Content: s.expandPlaceholders(a.ContentEn, rule, payload),
				}
			}
			if _, nErr := s.internalMessageSvc.SendSystemNotificationLocalized(ctx, title, content, alt, a.TargetUserIDs); nErr != nil {
				s.log.Errorf("flow execute: rule [%d] action[%d] notification failed: %v", rule.GetId(), i, nErr)
			}
		case "http":
			ok, httpStatus, errMsg := s.doHTTPAction(a, rule, payload)
			if !ok {
				s.log.Errorf("flow execute: rule [%d] action[%d] http failed (http=%d): %s",
					rule.GetId(), i, httpStatus, errMsg)
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
			cont.Chain = chain
			if err := s.taskService.EnqueueFlowResume(&cont, time.Duration(secs)*time.Second); err != nil {
				s.log.Errorf("flow execute: rule [%d] %s action[%d] delay enqueue failed, remaining actions dropped: %v",
					rule.GetId(), chain, i, err)
				return
			}
			s.log.Infof("flow execute: rule [%d] %s action[%d] delay %ds, resume at action[%d]",
				rule.GetId(), chain, i, secs, i+1)
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
	if payload.Chain == flowChainElse {
		// else 链续跑：延时在条件不满足分支内
		s.executeActionsJSON(ctx, rule, payload, rule.GetElseActions(), payload.ActionIndex, flowChainElse)
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

// expandPlaceholders 展开通知模板占位符：{event}/{entity_id}/{rule_id} 与
// 事件快照字段 {status}/{title}/{slug} 等（{field} 形态，未知字段原样保留）。
func (s *FlowService) expandPlaceholders(tpl string, rule *flowV1.FlowRule, payload *task.FlowExecutePayload) string {
	out := strings.ReplaceAll(tpl, "{event}", payload.Event)
	out = strings.ReplaceAll(out, "{entity_id}", strconv.FormatUint(uint64(payload.EntityID), 10))
	out = strings.ReplaceAll(out, "{rule_id}", strconv.FormatUint(uint64(rule.GetId()), 10))
	for field, value := range payload.Snapshot {
		out = strings.ReplaceAll(out, "{"+field+"}", value)
	}
	return out
}

// maxFlowHTTPBodyBytes http 动作 body 模板上限。
const maxFlowHTTPBodyBytes = 64 * 1024

// doHTTPAction 自定义 HTTP 外呼：POST/GET + 自定义头 + {field} 占位符 body 模板；
// secret 非空时按 webhook 同款 HMAC-SHA256 签名（X-Webhook-Signature）。
func (s *FlowService) doHTTPAction(a flowAction, rule *flowV1.FlowRule, payload *task.FlowExecutePayload) (bool, int, string) {
	method := strings.ToUpper(strings.TrimSpace(a.Method))
	if method == "" {
		method = http.MethodPost
	}
	if method != http.MethodGet && method != http.MethodPost {
		return false, 0, "unsupported method (GET/POST only)"
	}
	if len(a.Headers) > 8 {
		return false, 0, "too many headers (max 8)"
	}

	var bodyReader io.Reader
	if method == http.MethodPost {
		body := s.expandPlaceholders(a.Body, rule, payload)
		if body == "" {
			body = s.defaultFlowEnvelope(rule, payload)
		}
		if len(body) > maxFlowHTTPBodyBytes {
			return false, 0, "body exceeds 64KB"
		}
		bodyReader = strings.NewReader(body)
	}

	req, err := http.NewRequest(method, a.URL, bodyReader)
	if err != nil {
		return false, 0, "invalid url: " + err.Error()
	}
	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("User-Agent", "GoWind-Content-Hub-Flow/1.0")
	req.Header.Set("X-Flow-Event", payload.Event)
	for k, v := range a.Headers {
		if strings.EqualFold(k, "user-agent") || strings.EqualFold(k, "host") || strings.EqualFold(k, "content-length") {
			continue // 保留受控值，防头注入覆盖
		}
		req.Header.Set(k, v)
	}
	if a.Secret != "" {
		mac := hmac.New(sha256.New, []byte(a.Secret))
		mac.Write([]byte(s.defaultFlowEnvelope(rule, payload)))
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

// defaultFlowEnvelope 与 webhook 动作同构的默认 envelope（含快照对象）。
func (s *FlowService) defaultFlowEnvelope(rule *flowV1.FlowRule, payload *task.FlowExecutePayload) string {
	envelope := map[string]any{
		"event":     payload.Event,
		"entity_id": payload.EntityID,
		"rule_id":   rule.GetId(),
		"tenant_id": payload.TenantID,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}
	if len(payload.Snapshot) > 0 {
		envelope["snapshot"] = payload.Snapshot
	}
	raw, err := json.Marshal(envelope)
	if err != nil {
		return "{}"
	}
	return string(raw)
}

// postWebhook 签名投递（与 webhook 领域同款 envelope：事件 + 实体 id + 规则 id，不含内容全文）。
func (s *FlowService) postWebhook(a flowAction, payload *task.FlowExecutePayload) (bool, int, string) {
	envelope := map[string]any{
		"event":     payload.Event,
		"entity_id": payload.EntityID,
		"rule_id":   payload.RuleID,
		"tenant_id": payload.TenantID,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}
	if len(payload.Snapshot) > 0 {
		envelope["snapshot"] = payload.Snapshot
	}
	body, err := json.Marshal(envelope)
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
