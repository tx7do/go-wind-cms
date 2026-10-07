package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"

	aiV1 "go-wind-cms/api/gen/go/ai/service/v1"
)

// 异常信号阈值（24h 窗口）。
const (
	nightOpsAlertMin     = 1 // 深夜（本地 0-6 点）操作条数下限
	userFailedAlertMin   = 3 // 单用户失败操作条数下限
	loginFailAlertMin    = 3 // 单账号登录失败次数下限（疑似口令尝试）
	sensitiveOpsAlertMin = 1 // 敏感操作（DELETE/EXPORT/ASSIGN）条数下限
	sensitiveDetailMax   = 5 // 响应里最多携带的敏感操作明细条数
)

// 告警类型常量（前端 i18n 模板的 type 键，告警本身零文案）。
const (
	aiAlertBruteForce   = "BRUTE_FORCE"
	aiAlertNightOps     = "NIGHT_OPS"
	aiAlertFailedOps    = "FAILED_OPS"
	aiAlertSensitiveOps = "SENSITIVE_OPS"
)

// GetSafetyInsights 安全与异常洞察：规则预筛审计明细里的行为模式，
// 返回结构化告警事实（severity + type + facts），不含任何人类文案——
// 文案由前端 i18n 模板按界面语言渲染；LLM 总体评估按请求语言生成。
// 平台用户专属：明细日志为全平台数据，且会外发到模型端点。
func (s *AiContentService) GetSafetyInsights(ctx context.Context, req *aiV1.SafetyInsightsRequest) (*aiV1.SafetyInsightsResponse, error) {
	operatorId, tenantId := s.operatorIds(ctx)
	if tenantId != 0 {
		return nil, aiV1.ErrorForbidden("ai insights is platform-only")
	}
	_ = operatorId

	var alerts []*aiV1.AiInsightAlert
	addAlert := func(severity, alertType string, facts map[string]string) {
		alerts = append(alerts, &aiV1.AiInsightAlert{
			Severity: severity,
			Type:     alertType,
			Facts:    facts,
		})
	}

	// 1. 疑似口令尝试
	type failRow struct {
		Username string
		Fails    int
		LastIP   string
	}
	var failRows []failRow
	fails, err := s.entClient.DB().QueryContext(ctx,
		`SELECT COALESCE(username, ''), COUNT(*) AS fails, MAX(ip_address) AS last_ip
		 FROM sys_login_audit_logs
		 WHERE created_at >= NOW() - INTERVAL '24 hours' AND status = 'FAILED'
		 GROUP BY COALESCE(username, '') HAVING COUNT(*) >= $1
		 ORDER BY fails DESC`,
		loginFailAlertMin,
	)
	if err != nil {
		return nil, aiV1.ErrorInternalServerError("query login fails failed")
	}
	for fails.Next() {
		var row failRow
		if err = fails.Scan(&row.Username, &row.Fails, &row.LastIP); err != nil {
			continue
		}
		failRows = append(failRows, row)
	}
	_ = fails.Close()
	for _, acc := range failRows {
		addAlert("HIGH", aiAlertBruteForce, map[string]string{
			"username": acc.Username,
			"count":    fmt.Sprintf("%d", acc.Fails),
			"ip":       acc.LastIP,
		})
	}

	// 2. 用户行为模式（24h 聚合：总数/失败/深夜）
	type behaviorRow struct {
		Username string
		Total    int
		Failed   int
		Night    int
	}
	var behaviors []behaviorRow
	beh, err := s.entClient.DB().QueryContext(ctx,
		`SELECT user_id, MAX(username) AS username,
		        COUNT(*) AS total,
		        SUM(CASE WHEN success = false THEN 1 ELSE 0 END) AS failed,
		        SUM(CASE WHEN EXTRACT(HOUR FROM created_at AT TIME ZONE 'Asia/Shanghai') < 6 THEN 1 ELSE 0 END) AS night
		 FROM sys_operation_audit_logs
		 WHERE created_at >= NOW() - INTERVAL '24 hours' AND user_id IS NOT NULL
		 GROUP BY user_id
		 ORDER BY total DESC`,
	)
	if err != nil {
		return nil, aiV1.ErrorInternalServerError("query user behavior failed")
	}
	for beh.Next() {
		var uid any
		var row behaviorRow
		if err = beh.Scan(&uid, &row.Username, &row.Total, &row.Failed, &row.Night); err != nil {
			continue
		}
		behaviors = append(behaviors, row)
	}
	_ = beh.Close()
	for _, b := range behaviors {
		if b.Night >= nightOpsAlertMin {
			addAlert("MEDIUM", aiAlertNightOps, map[string]string{
				"username": b.Username,
				"count":    fmt.Sprintf("%d", b.Night),
			})
		}
		if b.Failed >= userFailedAlertMin {
			addAlert("MEDIUM", aiAlertFailedOps, map[string]string{
				"username": b.Username,
				"count":    fmt.Sprintf("%d", b.Failed),
				"total":    fmt.Sprintf("%d", b.Total),
			})
		}
	}

	// 3. 敏感操作明细
	type sensitiveRow struct {
		Username string
		Action   string
		Type     string
		At       time.Time
	}
	var sensitiveOps []sensitiveRow
	sen, err := s.entClient.DB().QueryContext(ctx,
		`SELECT COALESCE(username,''), action, resource_type, created_at
		 FROM sys_operation_audit_logs
		 WHERE created_at >= NOW() - INTERVAL '24 hours' AND action IN ('DELETE','EXPORT','ASSIGN')
		 ORDER BY created_at DESC
		 LIMIT $1`,
		sensitiveDetailMax,
	)
	if err != nil {
		return nil, aiV1.ErrorInternalServerError("query sensitive ops failed")
	}
	for sen.Next() {
		var row sensitiveRow
		if err = sen.Scan(&row.Username, &row.Action, &row.Type, &row.At); err != nil {
			continue
		}
		sensitiveOps = append(sensitiveOps, row)
	}
	_ = sen.Close()
	if len(sensitiveOps) >= sensitiveOpsAlertMin {
		alert := &aiV1.AiInsightAlert{
			Severity: "LOW",
			Type:     aiAlertSensitiveOps,
			Facts:    map[string]string{"count": fmt.Sprintf("%d", len(sensitiveOps))},
		}
		for _, op := range sensitiveOps {
			alert.Items = append(alert.Items, &aiV1.AiSensitiveOpItem{
				Username:     op.Username,
				Action:       op.Action,
				ResourceType: op.Type,
				CreatedAt:    op.At.Local().Format("01-02 15:04"),
			})
		}
		alerts = append(alerts, alert)
	}

	// 4. 总体评估：有告警才调 LLM，按请求语言输出；失败不阻断告警返回
	summary := ""
	if len(alerts) > 0 && s.providerRepo != nil {
		lang := req.GetLang()
		if lang == "" {
			lang = "zh-CN"
		}
		var alertFacts strings.Builder
		for _, a := range alerts {
			alertFacts.WriteString("- [" + a.Severity + "] " + a.Type)
			for k, v := range a.Facts {
				alertFacts.WriteString(fmt.Sprintf(" %s=%s", k, v))
			}
			alertFacts.WriteString("\n")
		}
		langLine := "用中文输出。"
		if strings.HasPrefix(lang, "en") {
			langLine = "Respond in English."
		}

		if provider, pErr := s.providerRepo.GetEnabledDefault(ctx); pErr == nil {
			if client, cErr := newOpenAIClientForProvider(provider); cErr == nil {
				cCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
				resp, gErr := client.CreateChatCompletion(cCtx, openai.ChatCompletionRequest{
					Model: ptrStrOr(provider.ModelName, "gpt-4o-mini"),
					Messages: []openai.ChatCompletionMessage{
						{Role: openai.ChatMessageRoleSystem, Content: "你是企业后台的安全运营助手。以下是从审计日志规则筛出的异常信号清单，" +
							"写一段 100 字以内的总体风险评估：概括风险等级与最需要优先处置的一项，" +
							"语气克制、只基于给定信号，不要编造未提及的信息。" + langLine},
						{Role: openai.ChatMessageRoleUser, Content: alertFacts.String()},
					},
				})
				cancel()
				if gErr != nil {
					s.log.Errorf("ai insights assessment failed, alerts still returned: %v", gErr)
				} else {
					summary = strings.TrimSpace(resp.Choices[0].Message.Content)
				}
			}
		}
	}

	return &aiV1.SafetyInsightsResponse{
		Alerts:  alerts,
		Summary: summary,
	}, nil
}
