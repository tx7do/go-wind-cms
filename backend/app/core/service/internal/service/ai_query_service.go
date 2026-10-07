package service

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/sashabaranov/go-openai"
	"github.com/tx7do/go-crud/viewer"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/timestamppb"

	internalMessageV1 "go-wind-cms/api/gen/go/internal_message/service/v1"
	"go-wind-cms/app/core/service/internal/data/ent/user"
	appViewer "go-wind-cms/pkg/entgo/viewer"
	"go-wind-cms/pkg/task"

	entCrud "github.com/tx7do/go-crud/entgo"

	aiV1 "go-wind-cms/api/gen/go/ai/service/v1"
	"go-wind-cms/app/core/service/internal/data"
	"go-wind-cms/app/core/service/internal/data/ent"
)

// ── 智能问数安全边界 ────────────────────────────────────────────────

const (
	// aiQueryRowLimit 结果行硬上限（缺 LIMIT 自动补；超限的显式 LIMIT 也被改写）。
	aiQueryRowLimit = 100

	// aiQueryCallTimeout SQL 生成与结论生成的单次 LLM 调用兜底时长。
	aiQueryCallTimeout = 90 * time.Second
)

// aiQueryTableWhitelist 平台表白名单：NL→SQL 只允许触达这些表。
// 白名单按本仓 CMS 真实 schema 适配（列清单经 information_schema 核对）。
// 键为表名，值为给 LLM 的列说明（拼进提示词当数据字典）。
var aiQueryTableWhitelist = map[string]string{
	"posts":            "帖子主表：id, code(唯一编码), status(POST_STATUS_DRAFT草稿/POST_STATUS_PUBLISHED已发布/POST_STATUS_SCHEDULED定时/POST_STATUS_TRASHED回收站), author_id(作者用户ID), author_name, is_featured(是否推荐), publish_time, created_at；注意：标题和摘要在 post_translations 表（经 post_id 关联），posts 表没有 title 列",
	"post_translations": "帖子翻译表：post_id(帖子ID), language_code(语言代码), title(标题), summary(摘要), word_count, created_at；查帖子标题必须 JOIN 此表",
	"categories":        "分类表：id, code, status, parent_id(父分类), post_count(帖子数冗余), created_at；分类名在 category_translations（经 category_id 关联）",
	"category_translations": "分类翻译表：category_id(分类ID), language_code, name(分类名), description, created_at",
	"tags":            "标签表：id, code, status, post_count(帖子数冗余), created_at；标签名在 tag_translations（经 tag_id 关联）",
	"tag_translations": "标签翻译表：tag_id(标签ID), language_code, name(标签名), created_at",
	"post_tags":       "帖子-标签关联表：post_id, tag_id（查每帖标签或每标签帖子数用此表 JOIN）",
	"post_categories": "帖子-分类关联表：post_id, category_id",
	"comments":        "评论表：id, content(评论内容), author_name(评论者名), author_email, status(PUBLISHED=已发布/PENDING=待审/TRASHED), ip_address, created_at；帖子评论经 content_type 与 object_id(=帖子ID) 关联",
	"interaction_counters": "互动计数表：target_type(对象类型), target_id(对象ID), metric(指标如 likes), count(计数), created_at",
	"media_assets":   "媒体资产表：id, filename(文件名), type, mime_type, size(字节), url, width, height, created_at",
	"sys_users":      "用户表：id, username(登录名), nickname(昵称), email, status, last_login_at, created_at；用户与角色是多对多，经 sys_user_roles 关联，sys_users 表没有 role_id 列",
	"sys_roles":      "角色表：id, name(角色名), code(角色编码), status, created_at",
	"sys_user_roles": "用户-角色关联表：user_id, role_id",
	"sys_operation_audit_logs": "操作审计日志：id, user_id, username, action(CREATE/UPDATE/DELETE/READ/EXPORT/…), resource_type, success(bool), failure_reason, created_at",
	"sys_login_audit_logs":     "登录审计日志：id, user_id, username, status(SUCCESS/FAILED), ip_address, created_at",
	"sys_ai_usage_logs":        "AI 用量流水：id, user_id, tenant_id, model_name, prompt_tokens, completion_tokens, total_tokens, duration_ms, created_at",
}

// aiQueryTenantTableWhitelist 租户表白名单：只含带 tenant_id 列、可安全暴露给
// 租户自助分析的本租户数据表（平台级表不在其中）。
var aiQueryTenantTableWhitelist = map[string]string{
	"posts":            "帖子表（本租户）：id, code, status(POST_STATUS_*), author_id, is_featured, created_at；标题在 post_translations（经 post_id 关联）",
	"post_translations": "帖子翻译表（本租户）：post_id, language_code, title, summary, created_at",
	"categories":       "分类表（本租户）：id, code, post_count, created_at；分类名在 category_translations",
	"category_translations": "分类翻译表（本租户）：category_id, language_code, name, created_at",
	"tags":             "标签表（本租户）：id, code, post_count, created_at；标签名在 tag_translations",
	"tag_translations": "标签翻译表（本租户）：tag_id, language_code, name, created_at",
	"comments":         "评论表（本租户）：id, content, author_name, status, created_at",
	"interaction_counters": "互动计数表（本租户）：target_type, target_id, metric, count, created_at",
	"media_assets":     "媒体资产表（本租户）：id, filename, mime_type, size, url, created_at",
	"sys_operation_audit_logs": "操作审计日志（本租户）：id, user_id, username, action, success, failure_reason, created_at",
	"sys_login_audit_logs":     "登录审计日志（本租户）：id, user_id, username, status(SUCCESS/FAILED), ip_address, created_at",
	"sys_ai_usage_logs":        "AI 用量流水（本租户）：id, user_id, model_name, prompt_tokens, completion_tokens, total_tokens, created_at",
}

// tenantIDPatternFor 构造"tenant_id = {id}"存在性校验的正则（容忍别名与空白）。
func tenantIDPatternFor(tenantID uint32) *regexp.Regexp {
	return regexp.MustCompile(`(?i)\btenant_id\s*=\s*` + fmt.Sprintf("%d", tenantID) + `\b`)
}

// aiQueryForbiddenPatterns 生成 SQL 的硬拒绝模式：写操作与危险语句。
var aiQueryForbiddenPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)\b(INSERT|UPDATE|DELETE|DROP|TRUNCATE|ALTER|CREATE|GRANT|REVOKE|COPY|VACUUM|REINDEX|CALL|DO|SET)\b`),
	regexp.MustCompile(`(?i)\b(pg_read_file|pg_ls_dir|pg_sleep|lo_import|lo_export|dblink)\b`),
	regexp.MustCompile(`\;`),     // 多语句
	regexp.MustCompile(`--|/\*`), // 注释（防注释截断绕过）
}

// aiQuerySelectPattern 必须以 SELECT 或 WITH 开头（含前导空白/括号容忍）。
var aiQuerySelectPattern = regexp.MustCompile(`(?is)^\s*\(?\s*(SELECT|WITH)\b`)

// ── 服务 ───────────────────────────────────────────────────────────

// AiQueryService 智能问数：NL → 只读 SQL → 结构化结果 →（可选）自然语言结论。
// 安全边界见 ai_query.proto：只读校验/表白名单/LIMIT 钳制/只读事务四重防线。
type AiQueryService struct {
	aiV1.UnimplementedAiQueryServiceServer
	log              *log.Helper
	providerRepo     *data.AiProviderRepo
	usageRepo        *data.AiUsageLogRepo
	entClient        *entCrud.EntClient[*ent.Client]
	internalMsgRepo  *data.InternalMessageRepo
	internalMessages *InternalMessageService
}

func NewAiQueryService(
	ctx *bootstrap.Context,
	providerRepo *data.AiProviderRepo,
	usageRepo *data.AiUsageLogRepo,
	entClient *entCrud.EntClient[*ent.Client],
	internalMsgRepo *data.InternalMessageRepo,
	internalMessages *InternalMessageService,
) *AiQueryService {
	return &AiQueryService{
		log:              ctx.NewLoggerHelper("ai_query/service/core-service"),
		providerRepo:     providerRepo,
		usageRepo:        usageRepo,
		entClient:        entClient,
		internalMsgRepo:  internalMsgRepo,
		internalMessages: internalMessages,
	}
}

// Ask 智能问数主链路：NL→SQL（第一次调用）→ 护栏校验 → 只读执行 → NL 结论（第二次调用）。
func (s *AiQueryService) Ask(ctx context.Context, req *aiV1.AskAiQueryRequest) (*aiV1.AskAiQueryResponse, error) {
	// 双白名单：平台用户查全平台表；租户用户查本租户数据表子集，
	// 且生成的 SQL 必须自带 tenant_id = 自身租户 过滤（后置校验 fail-closed）。
	tenantID, _ := s.operatorTenant(ctx)
	whitelist := aiQueryTableWhitelist
	tenantPattern := (*regexp.Regexp)(nil)
	if tenantID != 0 {
		whitelist = aiQueryTenantTableWhitelist
		tenantPattern = tenantIDPatternFor(tenantID)
	}
	question := strings.TrimSpace(req.GetQuestion())
	if question == "" {
		return nil, aiV1.ErrorBadRequest("question is required")
	}

	provider, err := s.providerRepo.GetEnabledDefault(ctx)
	if err != nil {
		return nil, err
	}
	client, err := newOpenAIClientForProvider(provider)
	if err != nil {
		return nil, err
	}

	lang := req.GetLang()
	if lang == "" {
		lang = "zh-CN"
	}

	// 1. NL → SQL（携带对话历史以消解追问里的指代，如"那只看 admin 的"）
	history := make([]aiV1.AiQueryHistoryItem, 0, len(req.GetHistory()))
	for _, h := range req.GetHistory() {
		history = append(history, aiV1.AiQueryHistoryItem{
			Question:      h.GetQuestion(),
			Sql:           h.GetSql(),
			ResultSummary: h.GetResultSummary(),
		})
	}
	sqlGenerated, genTokens, err := s.generateSQL(ctx, client, provider, question, lang, history, whitelist, tenantID)
	if err != nil {
		return nil, err
	}

	// 2. 四重护栏：只读语句 / 白名单表 / LIMIT 钳制 / 租户谓词复核
	cleaned, err := sanitizeSQLFor(sqlGenerated, whitelist, tenantPattern)
	if err != nil {
		return nil, aiV1.ErrorBadRequest("generated sql rejected: %v", err)
	}

	// 3. 只读事务执行
	columns, rows, execErr := s.executeReadOnly(ctx, cleaned)

	resp := &aiV1.AskAiQueryResponse{
		Sql:     cleaned,
		Columns: columns,
		Rows:    rows,
	}
	if execErr != nil {
		// 执行失败把数据库错误回给前端（便于用户改问题重试），SQL 照样返回
		s.log.Errorf("ai query execute failed: %v", execErr)
		resp.ErrorMessage = trans.Ptr(execErr.Error())
		return resp, nil
	}
	resp.RowCount = uint32(len(rows))

	// 4.（可选）结果 → 自然语言结论
	if !req.GetWithAnswer() || len(rows) == 0 {
		return resp, nil
	}
	answer, answerTokens, err := s.generateAnswer(ctx, client, provider, question, cleaned, columns, rows, lang)
	if err != nil {
		s.log.Errorf("ai query answer failed, returning rows only: %v", err)
	} else {
		resp.Answer = trans.Ptr(answer)
	}

	resp.TotalTokens = genTokens + answerTokens

	// 记账（两次调用合并总量；失败不阻断）
	if s.usageRepo != nil {
		modelName := ptrStrOr(provider.ModelName, "")
		_ = s.usageRepo.Create(ctx, &data.UsageRecord{
			UserID:      s.operatorID(ctx),
			TenantID:    tenantID,
			ModelName:   modelName,
			TotalTokens: resp.GetTotalTokens(),
		})
	}

	return resp, nil
}

// operatorTenant / operatorID：viewer 上下文取操作者（unary 中间件注入）。
func (s *AiQueryService) operatorTenant(ctx context.Context) (uint32, bool) {
	if vc, exist := viewer.FromContext(ctx); exist && vc != nil {
		return uint32(vc.TenantID()), true
	}
	return 0, false
}

func (s *AiQueryService) operatorID(ctx context.Context) uint32 {
	if uid, ok := viewerUserIDFromContext(ctx); ok {
		return uid
	}
	return 0
}

// generateSQL 数据字典 + 对话历史 + 问题 → 只读 SQL。
func (s *AiQueryService) generateSQL(ctx context.Context, client *openai.Client, provider *ent.AiProvider, question, lang string, history []aiV1.AiQueryHistoryItem, whitelist map[string]string, tenantID uint32) (string, uint32, error) {
	var sb strings.Builder
	sb.WriteString("你是本系统的只读数据分析引擎。用户的提问永远是对上面列出的数据库表发起查询——哪怕措辞听起来像个人账号问题，一律直接生成 SQL，绝不拒绝。\n")
	sb.WriteString("根据表结构与用户问题，生成一条 Postgres 只读 SELECT 语句。\n")
	sb.WriteString("规则：\n")
	sb.WriteString("1. 只输出 SQL 本身，不要任何解释、不要 Markdown 代码块标记；\n")
	sb.WriteString("2. 只能使用下面列出的表；涉及多表时用 JOIN；\n")
	sb.WriteString("3. 必须带 LIMIT（最大 100）；时间比较用 created_at 与 NOW()；\n")
	sb.WriteString(fmt.Sprintf("4. 结论文字用%s；\n", langName(lang)))
	sb.WriteString("5. 用户问题可能包含对上一轮查询的指代（如\"那只看…的\"、\"换成按天分组\"），结合对话历史理解其完整含义后生成独立可执行的 SQL。\n\n")
	sb.WriteString("可用表：\n")
	for table, desc := range whitelist {
		sb.WriteString("- " + table + "：" + desc + "\n")
	}
	if tenantID != 0 {
		sb.WriteString(fmt.Sprintf("强制规则：所有查询必须包含 tenant_id = %d 过滤条件（本租户数据边界），缺失视为无效输出。\n\n", tenantID))
	}

	// few-shot 用真实消息序列锚定行为（上下文示例对行为的约束力远强于 system 指令）
	exampleSQL := "SELECT p.id, pt.title, p.status, p.created_at FROM posts p " +
		"LEFT JOIN post_translations pt ON pt.post_id = p.id " +
		"WHERE p.status = 'POST_STATUS_PUBLISHED' AND p.created_at >= NOW() - INTERVAL '24 hours' " +
		"ORDER BY p.created_at DESC LIMIT 100"
	if tenantID != 0 {
		exampleSQL = strings.Replace(exampleSQL, "WHERE p.status", fmt.Sprintf("WHERE p.tenant_id = %d AND p.status", tenantID), 1)
	}
	messages := make([]openai.ChatCompletionMessage, 0, len(history)*2+6)
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: sb.String(),
	})
	// few-shot 1：CMS 高频问法——标题在翻译表，必须 JOIN（锚定本仓 schema 的最大陷阱）
	messages = append(messages,
		openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: "最近 24 小时发布的帖子有哪些？",
		},
		openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: exampleSQL,
		},
	)
	// few-shot 2：帖子-标签是多对多，经 post_tags 关联
	messages = append(messages,
		openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: "统计每个标签下的帖子数量",
		},
		openai.ChatCompletionMessage{
			Role: openai.ChatMessageRoleAssistant,
			Content: "SELECT t.id AS tag_id, tt.name AS tag_name, COUNT(pt.post_id) AS post_count " +
				"FROM tags t " +
				"LEFT JOIN tag_translations tt ON tt.tag_id = t.id " +
				"LEFT JOIN post_tags pt ON pt.tag_id = t.id " +
				"GROUP BY t.id, tt.name ORDER BY post_count DESC LIMIT 100",
		},
	)
	// 下标访问而非按值 range：AiQueryHistoryItem 内嵌 protoimpl.MessageState
	// （含 sync.Mutex），按值拷贝会连带复制锁。
	for i := range history {
		messages = append(messages,
			openai.ChatCompletionMessage{Role: openai.ChatMessageRoleUser, Content: history[i].GetQuestion()},
			openai.ChatCompletionMessage{Role: openai.ChatMessageRoleAssistant, Content: history[i].GetSql()},
		)
	}
	messages = append(messages, openai.ChatCompletionMessage{Role: openai.ChatMessageRoleUser, Content: question})

	qCtx, cancel := context.WithTimeout(ctx, aiQueryCallTimeout)
	defer cancel()
	resp, err := client.CreateChatCompletion(qCtx, openai.ChatCompletionRequest{
		Model:    ptrStrOr(provider.ModelName, "gpt-4o-mini"),
		Messages: messages,
	})
	if err != nil {
		return "", 0, aiV1.ErrorInternalServerError("ai query generate failed: %v", err)
	}
	tokens := uint32(resp.Usage.TotalTokens)
	if len(resp.Choices) == 0 {
		return "", tokens, aiV1.ErrorInternalServerError("ai query generate empty")
	}
	raw := resp.Choices[0].Message.Content
	s.log.Infof("ai query raw model output: %s", raw)
	return raw, tokens, nil
}

// sanitizeSQLFor 带白名单与租户谓词校验的四重护栏：
// SQL 提取 → 只读校验 → 危险模式拒绝 → 白名单表校验（含表名纠正）→ LIMIT 钳制。
func sanitizeSQLFor(raw string, whitelist map[string]string, tenantPattern *regexp.Regexp) (string, error) {
	sqlText := extractSQL(raw)
	if sqlText == "" {
		return "", fmt.Errorf("no SQL statement found in model output")
	}

	if !aiQuerySelectPattern.MatchString(sqlText) {
		return "", fmt.Errorf("only SELECT/WITH statements are allowed")
	}
	for _, pat := range aiQueryForbiddenPatterns {
		if pat.MatchString(sqlText) {
			return "", fmt.Errorf("forbidden pattern detected")
		}
	}
	// 白名单：提取 FROM/JOIN 后的表名逐一校验；模型偶发丢前缀或写单复数，依候选序列纠正
	tablePattern := regexp.MustCompile(`(?i)\b(FROM|JOIN)\s+([a-zA-Z_][a-zA-Z0-9_]*)`)
	matches := tablePattern.FindAllStringSubmatch(sqlText, -1)
	for _, m := range matches {
		table := strings.ToLower(m[2])
		if _, ok := whitelist[table]; ok {
			continue
		}
		corrected, ok := correctTableName(table, whitelist)
		if !ok {
			return "", fmt.Errorf("table %q is not in the whitelist", m[2])
		}
		fixPattern := regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(table) + `\b`)
		sqlText = fixPattern.ReplaceAllStringFunc(sqlText, func(string) string {
			return corrected
		})
	}
	if tenantPattern != nil && !tenantPattern.MatchString(sqlText) {
		return "", fmt.Errorf("missing mandatory tenant_id filter")
	}
	// LIMIT 钳制：无 LIMIT 补 aiQueryRowLimit；LIMIT N>上限改写
	lower := strings.ToLower(sqlText)
	if !regexp.MustCompile(`(?i)\blimit\s+\d+`).MatchString(lower) {
		sqlText = fmt.Sprintf("SELECT * FROM (%s) AS _q LIMIT %d", sqlText, aiQueryRowLimit)
	} else {
		reLimit := regexp.MustCompile(`(?i)\blimit\s+(\d+)`)
		sqlText = reLimit.ReplaceAllStringFunc(sqlText, func(m string) string {
			n := 0
			fmt.Sscanf(strings.TrimSpace(m[5:]), "%d", &n)
			if n > aiQueryRowLimit {
				return fmt.Sprintf("LIMIT %d", aiQueryRowLimit)
			}
			return m
		})
	}
	return sqlText, nil
}

// executeReadOnly 只读事务执行（数据库层兜底：READ ONLY 事务里写操作必败）。
func (s *AiQueryService) executeReadOnly(ctx context.Context, sqlText string) ([]string, []*aiV1.AiQueryRow, error) {
	db := s.entClient.DB()
	txCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	tx, err := db.BeginTx(txCtx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, nil, fmt.Errorf("begin read-only tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	resultRows, err := tx.QueryContext(txCtx, sqlText)
	if err != nil {
		return nil, nil, err
	}
	defer func() { _ = resultRows.Close() }()

	columns, err := resultRows.Columns()
	if err != nil {
		return nil, nil, err
	}

	out := make([]*aiV1.AiQueryRow, 0, 16)
	for resultRows.Next() {
		values := make([]interface{}, len(columns))
		ptrs := make([]interface{}, len(columns))
		for i := range values {
			ptrs[i] = &values[i]
		}
		if err = resultRows.Scan(ptrs...); err != nil {
			return nil, nil, err
		}
		row := &aiV1.AiQueryRow{Values: make([]string, 0, len(columns))}
		for _, v := range values {
			row.Values = append(row.Values, stringifyCellValue(v))
		}
		out = append(out, row)
	}
	if err = resultRows.Err(); err != nil {
		return nil, nil, err
	}
	return columns, out, nil
}

// generateAnswer 结果集 → 自然语言结论（第二次调用）。
func (s *AiQueryService) generateAnswer(ctx context.Context, client *openai.Client, provider *ent.AiProvider, question, sqlText string, columns []string, rows []*aiV1.AiQueryRow, lang string) (string, uint32, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("用户问题：%s\n\n执行的 SQL：%s\n\n查询结果（列：%s）：\n", question, sqlText, strings.Join(columns, ", ")))
	maxRows := len(rows)
	if maxRows > 20 {
		maxRows = 20 // 结论只看前 20 行，避免超长
	}
	for _, r := range rows[:maxRows] {
		sb.WriteString(strings.Join(r.Values, " | "))
		sb.WriteString("\n")
	}

	langLine := "用中文回答。"
	if strings.HasPrefix(lang, "en") {
		langLine = "Answer in English."
	}

	aCtx, cancel := context.WithTimeout(ctx, aiQueryCallTimeout)
	defer cancel()
	resp, err := client.CreateChatCompletion(aCtx, openai.ChatCompletionRequest{
		Model: ptrStrOr(provider.ModelName, "gpt-4o-mini"),
		Messages: []openai.ChatCompletionMessage{
			{
				Role: openai.ChatMessageRoleSystem,
				Content: "你是数据分析助手。根据用户问题、SQL 与查询结果，" +
					"给出简明的结论：直接回答问题本身，引用具体数字；" +
					"结果为空就明确说明没有匹配数据。" + langLine,
			},
			{Role: openai.ChatMessageRoleUser, Content: sb.String()},
		},
	})
	if err != nil {
		return "", 0, err
	}
	tokens := uint32(resp.Usage.TotalTokens)
	if len(resp.Choices) == 0 {
		return "", tokens, fmt.Errorf("empty answer")
	}
	return strings.TrimSpace(resp.Choices[0].Message.Content), tokens, nil
}

// aiQueryExtractPattern 提取用：非锚定定位第一条 SELECT/WITH。
var aiQueryExtractPattern = regexp.MustCompile(`(?is)\b(SELECT|WITH)\b`)

// extractSQL 从模型输出中提取 SQL（剥围栏、截分号、剥行注释）。
func extractSQL(raw string) string {
	s := stripSQLFence(raw)
	loc := aiQueryExtractPattern.FindStringIndex(s)
	if loc == nil {
		return ""
	}
	s = s[loc[0]:]
	if idx := strings.Index(s, ";"); idx >= 0 {
		s = s[:idx]
	}
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		if idx := strings.Index(line, "--"); idx >= 0 {
			lines[i] = line[:idx]
		}
	}
	return strings.TrimSpace(strings.Join(lines, "\n"))
}

// stripSQLFence 剥掉模型可能输出的 Markdown 代码块围栏。
func stripSQLFence(s string) string {
	s = strings.TrimSpace(s)
	for _, fence := range []string{"```sql", "```SQL", "```"} {
		if strings.HasPrefix(s, fence) {
			s = strings.TrimPrefix(s, fence)
			if idx := strings.LastIndex(s, "```"); idx >= 0 {
				s = s[:idx]
			}
			break
		}
	}
	return strings.TrimSpace(s)
}

// stringifyCellValue 数据库原始值 → 展示字符串（time/[]byte/nil 特判）。
func stringifyCellValue(v interface{}) string {
	switch val := v.(type) {
	case nil:
		return ""
	case time.Time:
		return val.Format("2006-01-02 15:04:05")
	case []byte:
		return string(val)
	default:
		return fmt.Sprintf("%v", val)
	}
}

// correctTableName 表名纠正：依候选序列匹配白名单（前缀 / 复数 s / 组合）。
func correctTableName(table string, whitelist map[string]string) (string, bool) {
	if _, ok := whitelist[table]; ok {
		return table, true
	}
	for _, candidate := range []string{
		"sys_" + table,
		table + "s",
		"sys_" + table + "s",
		strings.TrimSuffix(table, "s"),
		strings.TrimSuffix(table, "s") + "es",
	} {
		if _, ok := whitelist[candidate]; ok {
			return candidate, true
		}
	}
	return "", false
}

// ── 定时问数（asynq：ai.query.run） ─────────────────────────────────

// RunScheduledQuery 执行一次问数并返回格式化的文本结果（平台白名单）。
func (s *AiQueryService) RunScheduledQuery(ctx context.Context, question, lang string) (string, error) {
	provider, err := s.providerRepo.GetEnabledDefault(ctx)
	if err != nil {
		return "", err
	}
	client, err := newOpenAIClientForProvider(provider)
	if err != nil {
		return "", err
	}

	var history []aiV1.AiQueryHistoryItem
	sqlText, _, err := s.generateSQL(ctx, client, provider, question, lang, history, aiQueryTableWhitelist, 0)
	if err != nil {
		return "", err
	}
	cleaned, err := sanitizeSQLFor(sqlText, aiQueryTableWhitelist, nil)
	if err != nil {
		return "", err
	}
	_, rows, execErr := s.executeReadOnly(ctx, cleaned)
	if execErr != nil {
		return "", fmt.Errorf("execute: %w", execErr)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("问题：%s\n", question))
	sb.WriteString(fmt.Sprintf("SQL：%s\n", cleaned))
	if len(rows) == 0 {
		sb.WriteString("结果：无匹配数据")
	} else {
		sb.WriteString(fmt.Sprintf("结果（%d 行）：\n", len(rows)))
		for _, r := range rows {
			sb.WriteString("  " + strings.Join(r.Values, " | ") + "\n")
		}
	}
	return sb.String(), nil
}

// AsyncAiQueryRun 定时问数任务 handler：执行问题 → 站内信推送结果给平台侧用户。
// asynq ctx 无 viewer，走系统查看器。
func (s *AiQueryService) AsyncAiQueryRun(taskType string, data *task.AiQueryRunPayload) error {
	ctx := appViewer.NewSystemViewerContext(context.Background())
	question := data.Question
	if question == "" {
		return fmt.Errorf("question is required")
	}

	result, err := s.RunScheduledQuery(ctx, question, "zh-CN")
	if err != nil {
		return err
	}

	title := "AI 问数结果：" + question
	if err = s.pushInternalMessage(ctx, title, result); err != nil {
		s.log.Errorf("ai query run: push result failed: %v", err)
	}
	return nil
}

// pushInternalMessage 创建站内信（平台侧，NOTIFICATION 类型）并给每个
// 平台用户（tenant_id=0）建收件行。
func (s *AiQueryService) pushInternalMessage(ctx context.Context, title, content string) error {
	now := time.Now()
	msg, err := s.internalMsgRepo.Create(ctx, &internalMessageV1.CreateInternalMessageRequest{
		Data: &internalMessageV1.InternalMessage{
			Title:     &title,
			Content:   &content,
			Status:    trans.Ptr(internalMessageV1.InternalMessage_PUBLISHED),
			Type:      trans.Ptr(internalMessageV1.InternalMessage_NOTIFICATION),
			CreatedBy: trans.Ptr(uint32(0)),
			CreatedAt: timestamppb.New(now),
		},
	})
	if err != nil {
		return err
	}

	// 投递给平台侧（tenant_id=0）用户
	users, err := s.entClient.Client().User.Query().
		Where(user.TenantIDEQ(0), user.DeletedAtIsNil()).
		All(ctx)
	if err != nil {
		return err
	}
	for _, u := range users {
		if err = s.internalMessages.sendNotification(ctx, msg.GetId(), u.ID, 0, &now, title, content); err != nil {
			s.log.Errorf("push to user %d failed: %v", u.ID, err)
		}
	}
	return nil
}
