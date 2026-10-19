package service

import (
	"context"
	"database/sql"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/v9"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	entCrud "github.com/tx7do/go-crud/entgo"

	monitorAlertV1 "go-wind-cms/api/gen/go/monitor_alert/service/v1"
	notifyV1 "go-wind-cms/api/gen/go/notify/service/v1"

	"go-wind-cms/app/core/service/internal/data"
	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/pkg/task"
)

// MonitorAlertService 监控告警：规则 CRUD + 周期评估 + 告警通知分发。
//
// 评估内核 evaluateOnce 被两个入口共享：周期扫描任务（每 5 分钟）与管理页的
// 「立即评估」。「告警→恢复」成对通知：越限发告警、回到阈值内发恢复，
// 持续越限按 cooldown_minutes 重发（不是只发一次然后装死，也不是每次扫描都轰炸）。
//
// 与上游 go-wind-admin 的差异：指标采集在本进程内直接完成（runtime/DB 连接池/
// Redis DBSIZE），不依赖独立的 server-monitor / redis-cache 监控域；告警文案
// 固定中文（CMS 用户暂无偏好语言字段，接入用户 locale 后可升级）。
type MonitorAlertService struct {
	monitorAlertV1.UnimplementedMonitorAlertServiceServer

	repo      *data.MonitorAlertRuleRepo
	entClient *entCrud.EntClient[*ent.Client]
	rdb       *redis.Client
	log       *log.Helper

	// 构造期留 nil，装配期由 wiring 在服务就绪后 RegisterNotifier 注入
	notifier Notifier
}

func NewMonitorAlertService(
	ctx *bootstrap.Context,
	repo *data.MonitorAlertRuleRepo,
	entClient *entCrud.EntClient[*ent.Client],
	rdb *redis.Client,
) *MonitorAlertService {
	svc := &MonitorAlertService{
		repo:      repo,
		entClient: entClient,
		rdb:       rdb,
		log:       ctx.NewLoggerHelper("monitor-alert/service/core-service"),
		// notifier 默认"未装配"占位：装配期由 RegisterNotifier 换成 NotificationService
		// （与 InternalMessageService 同模式）。
	}
	return svc
}

// RegisterNotifier 装配通知出口（NotificationService）。
func (s *MonitorAlertService) RegisterNotifier(notifier Notifier) {
	s.notifier = notifier
}

// ==== CRUD（平台管理员，经 requireNotificationPlatformContext 同一条缝） ====

func (s *MonitorAlertService) ListMonitorAlertRule(ctx context.Context, req *paginationV1.PagingRequest) (*monitorAlertV1.ListMonitorAlertRuleResponse, error) {
	if err := requireNotificationPlatformContext(ctx, "monitor alert rules"); err != nil {
		return nil, err
	}
	return s.repo.List(ctx, req)
}

func (s *MonitorAlertService) GetMonitorAlertRule(ctx context.Context, req *monitorAlertV1.GetMonitorAlertRuleRequest) (*monitorAlertV1.MonitorAlertRule, error) {
	if req == nil || req.GetId() == 0 {
		return nil, errors.BadRequest("MONITOR_ALERT", "id is required")
	}
	if err := requireNotificationPlatformContext(ctx, "monitor alert rules"); err != nil {
		return nil, err
	}
	return s.repo.Get(ctx, req.GetId())
}

func (s *MonitorAlertService) CreateMonitorAlertRule(ctx context.Context, req *monitorAlertV1.CreateMonitorAlertRuleRequest) (*monitorAlertV1.MonitorAlertRule, error) {
	if req == nil || req.GetData() == nil {
		return nil, errors.BadRequest("MONITOR_ALERT", "invalid parameter")
	}
	if err := requireNotificationPlatformContext(ctx, "monitor alert rules"); err != nil {
		return nil, err
	}
	if err := validateAlertRule(req.GetData()); err != nil {
		return nil, err
	}

	operatorID, _ := viewerUserIDFromContext(ctx)
	id, err := s.repo.Create(ctx, req.GetData(), operatorID)
	if err != nil {
		return nil, err
	}

	return s.repo.Get(ctx, id)
}

func (s *MonitorAlertService) UpdateMonitorAlertRule(ctx context.Context, req *monitorAlertV1.UpdateMonitorAlertRuleRequest) (*emptypb.Empty, error) {
	if req == nil || req.GetData() == nil {
		return nil, errors.BadRequest("MONITOR_ALERT", "invalid parameter")
	}
	if err := requireNotificationPlatformContext(ctx, "monitor alert rules"); err != nil {
		return nil, err
	}
	if err := validateAlertRule(req.GetData()); err != nil {
		return nil, err
	}

	operatorID, _ := viewerUserIDFromContext(ctx)
	if err := s.repo.Update(ctx, req.GetId(), req.GetData(), operatorID); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *MonitorAlertService) DeleteMonitorAlertRule(ctx context.Context, req *monitorAlertV1.DeleteMonitorAlertRuleRequest) (*emptypb.Empty, error) {
	if err := requireNotificationPlatformContext(ctx, "monitor alert rules"); err != nil {
		return nil, err
	}
	if err := s.repo.Delete(ctx, req.GetId()); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// EvaluateMonitorAlerts 手动触发一轮评估（与周期扫描同一内核）。
func (s *MonitorAlertService) EvaluateMonitorAlerts(ctx context.Context, _ *monitorAlertV1.EvaluateMonitorAlertsRequest) (*monitorAlertV1.EvaluateMonitorAlertsResponse, error) {
	if err := requireNotificationPlatformContext(ctx, "monitor alert rules"); err != nil {
		return nil, err
	}

	outcomes := s.evaluateOnce(ctx)

	resp := &monitorAlertV1.EvaluateMonitorAlertsResponse{Outcomes: make([]*monitorAlertV1.EvaluateMonitorAlertsResponse_RuleOutcome, 0, len(outcomes))}
	for _, o := range outcomes {
		resp.Outcomes = append(resp.Outcomes, &monitorAlertV1.EvaluateMonitorAlertsResponse_RuleOutcome{
			RuleId:       o.ruleID,
			Name:         o.name,
			Firing:       o.firing,
			CurrentValue: o.currentValue,
			Notified:     o.notified,
			Reason:       trans.Ptr(o.reason),
		})
	}
	return resp, nil
}

// ==== 评估内核 ====

type alertOutcome struct {
	ruleID       uint32
	name         string
	firing       bool
	currentValue *float64
	notified     bool
	reason       string
}

// validateAlertRule 建/改规则的公共校验：非布尔指标必须有 op；渠道/目标合法。
func validateAlertRule(data *monitorAlertV1.MonitorAlertRule) error {
	if data == nil {
		return errors.BadRequest("MONITOR_ALERT", "invalid parameter")
	}
	if data.GetMetric() == monitorAlertV1.MonitorMetric_MONITOR_METRIC_UNSPECIFIED {
		return errors.BadRequest("MONITOR_ALERT", "metric is required")
	}
	if data.GetChannel() != notifyV1.Channel_EMAIL && data.GetChannel() != notifyV1.Channel_WEBHOOK {
		return errors.BadRequest("MONITOR_ALERT", "channel must be EMAIL or WEBHOOK")
	}
	if len(strings.TrimSpace(data.GetTarget())) < 3 {
		return errors.BadRequest("MONITOR_ALERT", "target is required")
	}
	// DB_PING_FAIL 是布尔指标（越限即 ping 失败），op/threshold 无意义可缺省；
	// 其余指标必须有比较运算与阈值。
	if data.GetMetric() != monitorAlertV1.MonitorMetric_DB_PING_FAIL {
		if data.GetOp() == monitorAlertV1.AlertOp_ALERT_OP_UNSPECIFIED {
			return errors.BadRequest("MONITOR_ALERT", "op is required for metric "+data.GetMetric().String())
		}
		if data.GetThreshold() == 0 {
			return errors.BadRequest("MONITOR_ALERT", "threshold is required for metric "+data.GetMetric().String())
		}
	}
	return nil
}

// metricValue 取规则的当前指标值。全部在本进程内采集：
//   - GO_*：core 自身 runtime；
//   - DB_*：ent 连接池的 *sql.DB（Stats/Ping）；
//   - REDIS_DB_SIZE：asynq/缓存同一 Redis 的 DBSIZE。
func (s *MonitorAlertService) metricValue(ctx context.Context, rule *monitorAlertV1.MonitorAlertRule) (float64, error) {
	switch rule.GetMetric() {
	case monitorAlertV1.MonitorMetric_GO_GOROUTINES:
		return float64(runtime.NumGoroutine()), nil
	case monitorAlertV1.MonitorMetric_GO_MEM_ALLOC_MB:
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		return float64(m.Alloc) / 1024 / 1024, nil
	case monitorAlertV1.MonitorMetric_DB_OPEN_CONNECTIONS,
		monitorAlertV1.MonitorMetric_DB_PING_FAIL:
		sqlDB, err := s.sqlDB()
		if err != nil {
			return 0, fmt.Errorf("database handle unavailable: %w", err)
		}
		if rule.GetMetric() == monitorAlertV1.MonitorMetric_DB_OPEN_CONNECTIONS {
			return float64(sqlDB.Stats().OpenConnections), nil
		}
		if err := sqlDB.PingContext(ctx); err != nil {
			return 1, nil // ping 失败 = 布尔指标越限
		}
		return 0, nil
	case monitorAlertV1.MonitorMetric_REDIS_DB_SIZE:
		if s.rdb == nil {
			return 0, fmt.Errorf("redis client unavailable")
		}
		n, err := s.rdb.DBSize(ctx).Result()
		if err != nil {
			return 0, fmt.Errorf("redis dbsize failed: %w", err)
		}
		return float64(n), nil
	default:
		return 0, fmt.Errorf("unsupported metric %s", rule.GetMetric().String())
	}
}

// sqlDB 取 ent 客户端底层的 *sql.DB（连接池统计与 ping 用）。
func (s *MonitorAlertService) sqlDB() (*sql.DB, error) {
	if s.entClient == nil {
		return nil, fmt.Errorf("ent client is nil")
	}
	db := s.entClient.DB()
	if db == nil {
		return nil, fmt.Errorf("sql db is nil")
	}
	return db, nil
}

// isBreached 阈值判定。DB_PING_FAIL 是布尔指标、绕过比较：ping 失败（值 1）即越限，
// ping 正常（值 0）永不越限——不绕过的话，未填 op/threshold 的规则会以零值兜底
// 走 GE 比较（0 >= 0），把"一切正常"判成"持续越限"。
func isBreached(rule *monitorAlertV1.MonitorAlertRule, value float64) bool {
	if rule.GetMetric() == monitorAlertV1.MonitorMetric_DB_PING_FAIL {
		return value >= 1
	}
	switch rule.GetOp() {
	case monitorAlertV1.AlertOp_LE:
		return value <= rule.GetThreshold()
	case monitorAlertV1.AlertOp_GE:
		fallthrough
	default:
		return value >= rule.GetThreshold()
	}
}

// evaluateOnce 评估一轮：逐条启用规则求值 → 决定通知 → 回写状态。
//
// 通知时机（cooldown 以 last_alerted_at 为锚）：
//   - 未越限且上次越限 → 「恢复」通知（不等冷却：恢复通知晚到没有意义）；
//   - 越限且（上次未越限 或 冷却已过）→ 「告警」通知；
//   - 越限但冷却未到 → 只回写 last_value，不重发。
//
// 采集失败轮：跳过该规则且**不清 firing 状态**——采集抖动不是恢复，
// 把它当恢复会发假通知，把它当持续越限重发会发重复告警，两个都比"跳过一轮"糟。
func (s *MonitorAlertService) evaluateOnce(ctx context.Context) []alertOutcome {
	rules, err := s.repo.ListEnabled(ctx)
	if err != nil {
		s.log.Errorf("monitor alert scan: list rules failed: %s", err.Error())
		return nil
	}
	if len(rules) == 0 {
		return nil
	}

	outcomes := make([]alertOutcome, 0, len(rules))
	now := time.Now()
	for _, rule := range rules {
		outcome := alertOutcome{ruleID: rule.GetId(), name: rule.GetName()}

		value, err := s.metricValue(ctx, rule)
		if err != nil {
			// 采集失败：保留上次 firing 态，下一轮再试
			outcome.firing = rule.GetLastFiring()
			outcome.reason = "sample collection failed, skipped"
			outcomes = append(outcomes, outcome)
			continue
		}
		outcome.currentValue = &value
		outcome.firing = isBreached(rule, value)

		switch {
		case outcome.firing && !rule.GetLastFiring():
			outcome.notified = s.sendAlert(ctx, rule, value, false, now)
			if !outcome.notified {
				outcome.reason = "alert dispatch failed"
			}
		case outcome.firing && rule.GetLastFiring() && cooldownElapsed(rule, now):
			outcome.notified = s.sendAlert(ctx, rule, value, false, now)
			outcome.reason = "cooldown elapsed, re-alerted"
		case !outcome.firing && rule.GetLastFiring():
			outcome.notified = s.sendAlert(ctx, rule, value, true, now)
			if !outcome.notified {
				outcome.reason = "resolve dispatch failed"
			}
		default:
			outcome.reason = "no state change"
		}

		// 状态回写：firing 如实反映本次求值（否则"恢复"判定会漂）；
		// 冷却锚 last_alerted_at 只在通知真正发出时推进——发失败下一轮会立即重试，
		// 而不是被假锚按进冷却窗口里漏发。
		alertedAt := &now
		if !outcome.notified {
			alertedAt = nil
		}
		if err := s.repo.MarkFiring(ctx, rule.GetId(), outcome.firing, outcome.currentValue, alertedAt); err != nil {
			s.log.Errorf("monitor alert scan: mark rule [%d] failed: %s", rule.GetId(), err.Error())
		}
		outcomes = append(outcomes, outcome)
	}

	return outcomes
}

func cooldownElapsed(rule *monitorAlertV1.MonitorAlertRule, now time.Time) bool {
	last := rule.GetLastAlertedAt()
	if last == nil {
		return true
	}
	return now.Sub(last.AsTime()) >= time.Duration(rule.GetCooldownMinutes())*time.Minute
}

// sendAlert 发告警/恢复通知，返回是否真的发出去了（冷却锚只按真发推进）。
// 失败不中断评估循环——一次 SMTP/Webhook 失败不该拦住其余规则，
// 本轮的缺口由下一轮扫描补发（firing 状态没变，下次仍会命中发送分支）。
func (s *MonitorAlertService) sendAlert(ctx context.Context, rule *monitorAlertV1.MonitorAlertRule, value float64, resolved bool, now time.Time) bool {
	if s.notifier == nil {
		s.log.Errorf("monitor alert [%s] dispatch skipped: notification seam is not wired", rule.GetName())
		return false
	}
	title, content := alertText(rule, value, resolved, now)
	resp, err := s.notifier.SendDirect(ctx, &notifyV1.SendDirectNotificationRequest{
		EventType:      notifyV1.EventType_MONITOR_ALERT,
		Channel:        rule.Channel,
		Target:         rule.GetTarget(),
		Title:          title,
		Content:        content,
		RelatedId:      rule.Id,
		OperatorUserId: trans.Ptr(uint32(0)),
	})
	if err != nil {
		s.log.Errorf("monitor alert [%s] notify (%s → %s) failed: %s",
			rule.GetName(), rule.GetChannel().String(), rule.GetTarget(), err.Error())
		return false
	}
	s.log.Infof("monitor alert [%s] %s: value=%.2f (delivery=%d)",
		rule.GetName(), map[bool]string{true: "RESOLVED", false: "FIRING"}[resolved], value, resp.GetDeliveryId())
	return true
}

// alertText 告警/恢复文案（固定中文；指标与比较符用枚举字面量，与闸门报错口径一致）。
func alertText(rule *monitorAlertV1.MonitorAlertRule, value float64, resolved bool, now time.Time) (string, string) {
	ts := now.Format("2006-01-02 15:04:05")
	if resolved {
		return fmt.Sprintf("[已恢复] %s", rule.GetName()),
			fmt.Sprintf("监控指标 %s 已回到阈值内（当前值 %.2f，阈值 %s %.2f），时间 %s。",
				rule.GetMetric().String(), value, rule.GetOp().String(), rule.GetThreshold(), ts)
	}
	return fmt.Sprintf("[告警] %s", rule.GetName()),
		fmt.Sprintf("监控指标 %s 当前值 %.2f，越过阈值 %s %.2f，时间 %s。持续越限时每 %d 分钟重发一次。",
			rule.GetMetric().String(), value, rule.GetOp().String(), rule.GetThreshold(), ts, rule.GetCooldownMinutes())
}

// ==== 周期任务 ====

// AsyncMonitorAlertScan 周期扫描 handler（系统级常驻任务，SystemContext 上下文）。
func (s *MonitorAlertService) AsyncMonitorAlertScan(taskType string, _ *task.MonitorAlertScanTaskData) error {
	ctx := context.Background()
	s.log.Infof("[%s] scan start", taskType)

	outcomes := s.evaluateOnce(ctx)

	var notified int
	for _, o := range outcomes {
		if o.notified {
			notified++
		}
	}
	s.log.Infof("[%s] scan done: %d rules, %d notified", taskType, len(outcomes), notified)
	return nil
}
