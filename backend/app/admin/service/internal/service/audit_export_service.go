package service

import (
	"context"
	"encoding/csv"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"github.com/xuri/excelize/v2"
	"google.golang.org/protobuf/types/known/timestamppb"

	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"

	aiV1 "go-wind-cms/api/gen/go/ai/service/v1"
	auditV1 "go-wind-cms/api/gen/go/audit/service/v1"
	identityV1 "go-wind-cms/api/gen/go/identity/service/v1"
	permissionV1 "go-wind-cms/api/gen/go/permission/service/v1"

	appMetadata "go-wind-cms/pkg/metadata"
	"go-wind-cms/pkg/middleware/auth"
)

// 服务端导出的单次行数上限：突破客户端聚合的 1 万行上限，同时防无界查询 OOM。
const exportMaxRows = 50000

const exportPageSize = 1000

// colDef 导出列：标题 + 取值函数。
type colDef struct {
	header string
	value  func(v any) string
}

// AuditExportService 审计日志/用量流水服务端导出（XLSX/CSV，二进制响应）。
// 手动注册路由（要写 Content-Disposition，proto 生成路由不支持）。
//
// 本仓两服务拓扑下的实现与上游不同：BFF 无 DB，逐页复用各域 List RPC 聚合行——
// 条件与列表页同源（query 参数透传给 PagingRequest.Query，contains 语义一致），
// "导出的就是当前搜索看到的"。租户视角由 core 侧 viewer 裁定。
type AuditExportService struct {
	loginClient            auditV1.LoginAuditLogServiceClient
	apiClient              auditV1.ApiAuditLogServiceClient
	operationClient        auditV1.OperationAuditLogServiceClient
	dataAccessClient       auditV1.DataAccessAuditLogServiceClient
	permissionClient       auditV1.PermissionAuditLogServiceClient
	policyEvaluationClient permissionV1.PolicyEvaluationLogServiceClient
	usageClient            aiV1.AiUsageLogServiceClient
	log                    *log.Helper

	// tokenChecker 手动路由的鉴权锚：导出端点未设 Operation，auth 中间件的
	// selector 不会应用（不会注入操作人），须自行验 Bearer token。
	tokenChecker auth.AccessTokenChecker
}

func NewAuditExportService(
	ctx *bootstrap.Context,
	loginClient auditV1.LoginAuditLogServiceClient,
	apiClient auditV1.ApiAuditLogServiceClient,
	operationClient auditV1.OperationAuditLogServiceClient,
	dataAccessClient auditV1.DataAccessAuditLogServiceClient,
	permissionClient auditV1.PermissionAuditLogServiceClient,
	policyEvaluationClient permissionV1.PolicyEvaluationLogServiceClient,
	usageClient aiV1.AiUsageLogServiceClient,
	tokenChecker auth.AccessTokenChecker,
) *AuditExportService {
	return &AuditExportService{
		loginClient:            loginClient,
		apiClient:              apiClient,
		operationClient:        operationClient,
		dataAccessClient:       dataAccessClient,
		permissionClient:       permissionClient,
		policyEvaluationClient: policyEvaluationClient,
		usageClient:            usageClient,
		log:                    ctx.NewLoggerHelper("audit-export/service/admin-service"),
		tokenChecker:           tokenChecker,
	}
}

// ServeExport 审计日志导出：GET /admin/v1/audit-logs:export
// ?type=login|api|operation|data_access|permission|policy_evaluation&format=xlsx|csv&query=<列表同源条件>&maxRows=
func (s *AuditExportService) ServeExport(w http.ResponseWriter, r *http.Request) error {
	ctx, ok := s.authorize(w, r)
	if !ok {
		return nil
	}

	q := r.URL.Query()
	switch q.Get("type") {
	case "login":
		s.exportType(ctx, w, q, "login-audit-logs",
			func(ctx context.Context, page int64) ([]any, uint64, error) {
				resp, err := s.loginClient.List(ctx, pagingReq(page, q.Get("query")))
				if err != nil {
					return nil, 0, err
				}
				return toAnySlice(resp.GetItems()), resp.GetTotal(), nil
			},
			[]colDef{
				{"用户名", func(v any) string { return v.(*auditV1.LoginAuditLog).GetUsername() }},
				{"IP", func(v any) string { return v.(*auditV1.LoginAuditLog).GetIpAddress() }},
				{"动作", func(v any) string { return v.(*auditV1.LoginAuditLog).GetActionType().String() }},
				{"风险等级", func(v any) string { return v.(*auditV1.LoginAuditLog).GetRiskLevel().String() }},
				{"登录方式", func(v any) string { return v.(*auditV1.LoginAuditLog).GetLoginMethod().String() }},
				{"失败原因", func(v any) string { return v.(*auditV1.LoginAuditLog).GetFailureReason() }},
				{"时间", func(v any) string { return tsStr(v.(*auditV1.LoginAuditLog).CreatedAt) }},
			},
		)
	case "api":
		s.exportType(ctx, w, q, "api-audit-logs",
			func(ctx context.Context, page int64) ([]any, uint64, error) {
				resp, err := s.apiClient.List(ctx, pagingReq(page, q.Get("query")))
				if err != nil {
					return nil, 0, err
				}
				return toAnySlice(resp.GetItems()), resp.GetTotal(), nil
			},
			[]colDef{
				{"模块", func(v any) string { return v.(*auditV1.ApiAuditLog).GetApiModule() }},
				{"操作", func(v any) string { return v.(*auditV1.ApiAuditLog).GetApiOperation() }},
				{"方法", func(v any) string { return v.(*auditV1.ApiAuditLog).GetHttpMethod() }},
				{"路径", func(v any) string { return v.(*auditV1.ApiAuditLog).GetPath() }},
				{"状态码", func(v any) string { return strconv.FormatUint(uint64(v.(*auditV1.ApiAuditLog).GetStatusCode()), 10) }},
				{"耗时(ms)", func(v any) string { return strconv.FormatUint(uint64(v.(*auditV1.ApiAuditLog).GetLatencyMs()), 10) }},
				{"成功", func(v any) string { return boolStr(v.(*auditV1.ApiAuditLog).Success) }},
				{"原因", func(v any) string { return v.(*auditV1.ApiAuditLog).GetReason() }},
				{"时间", func(v any) string { return tsStr(v.(*auditV1.ApiAuditLog).CreatedAt) }},
			},
		)
	case "operation":
		s.exportType(ctx, w, q, "operation-audit-logs",
			func(ctx context.Context, page int64) ([]any, uint64, error) {
				resp, err := s.operationClient.List(ctx, pagingReq(page, q.Get("query")))
				if err != nil {
					return nil, 0, err
				}
				return toAnySlice(resp.GetItems()), resp.GetTotal(), nil
			},
			[]colDef{
				{"用户名", func(v any) string { return v.(*auditV1.OperationAuditLog).GetUsername() }},
				{"资源类型", func(v any) string { return v.(*auditV1.OperationAuditLog).GetResourceType() }},
				{"资源ID", func(v any) string { return v.(*auditV1.OperationAuditLog).GetResourceId() }},
				{"动作", func(v any) string { return v.(*auditV1.OperationAuditLog).GetAction().String() }},
				{"成功", func(v any) string { return boolStr(v.(*auditV1.OperationAuditLog).Success) }},
				{"失败原因", func(v any) string { return v.(*auditV1.OperationAuditLog).GetFailureReason() }},
				{"IP", func(v any) string { return v.(*auditV1.OperationAuditLog).GetIpAddress() }},
				{"时间", func(v any) string { return tsStr(v.(*auditV1.OperationAuditLog).CreatedAt) }},
			},
		)
	case "data_access":
		s.exportType(ctx, w, q, "data-access-audit-logs",
			func(ctx context.Context, page int64) ([]any, uint64, error) {
				resp, err := s.dataAccessClient.List(ctx, pagingReq(page, q.Get("query")))
				if err != nil {
					return nil, 0, err
				}
				return toAnySlice(resp.GetItems()), resp.GetTotal(), nil
			},
			[]colDef{
				{"用户名", func(v any) string { return v.(*auditV1.DataAccessAuditLog).GetUsername() }},
				{"数据源", func(v any) string { return v.(*auditV1.DataAccessAuditLog).GetDataSource() }},
				{"表名", func(v any) string { return v.(*auditV1.DataAccessAuditLog).GetTableName() }},
				{"访问类型", func(v any) string { return v.(*auditV1.DataAccessAuditLog).GetAccessType().String() }},
				{"影响行数", func(v any) string { return strconv.FormatUint(uint64(v.(*auditV1.DataAccessAuditLog).GetAffectedRows()), 10) }},
				{"成功", func(v any) string { return boolStr(v.(*auditV1.DataAccessAuditLog).Success) }},
				{"时间", func(v any) string { return tsStr(v.(*auditV1.DataAccessAuditLog).CreatedAt) }},
			},
		)
	case "permission":
		s.exportType(ctx, w, q, "permission-audit-logs",
			func(ctx context.Context, page int64) ([]any, uint64, error) {
				resp, err := s.permissionClient.List(ctx, pagingReq(page, q.Get("query")))
				if err != nil {
					return nil, 0, err
				}
				return toAnySlice(resp.GetItems()), resp.GetTotal(), nil
			},
			[]colDef{
				{"操作者", func(v any) string { return v.(*auditV1.PermissionAuditLog).GetOperatorName() }},
				{"目标类型", func(v any) string { return v.(*auditV1.PermissionAuditLog).GetTargetType() }},
				{"目标名称", func(v any) string { return v.(*auditV1.PermissionAuditLog).GetTargetName() }},
				{"动作", func(v any) string { return v.(*auditV1.PermissionAuditLog).GetAction().String() }},
				{"IP", func(v any) string { return v.(*auditV1.PermissionAuditLog).GetIpAddress() }},
				{"原因", func(v any) string { return v.(*auditV1.PermissionAuditLog).GetReason() }},
				{"时间", func(v any) string { return tsStr(v.(*auditV1.PermissionAuditLog).CreatedAt) }},
			},
		)
	case "policy_evaluation":
		s.exportType(ctx, w, q, "policy-evaluation-logs",
			func(ctx context.Context, page int64) ([]any, uint64, error) {
				resp, err := s.policyEvaluationClient.List(ctx, pagingReq(page, q.Get("query")))
				if err != nil {
					return nil, 0, err
				}
				return toAnySlice(resp.GetItems()), resp.GetTotal(), nil
			},
			[]colDef{
				{"用户ID", func(v any) string { return strconv.FormatUint(uint64(v.(*permissionV1.PolicyEvaluationLog).GetUserId()), 10) }},
				{"请求路径", func(v any) string { return v.(*permissionV1.PolicyEvaluationLog).GetRequestPath() }},
				{"方法", func(v any) string { return v.(*permissionV1.PolicyEvaluationLog).GetRequestMethod() }},
				{"结果", func(v any) string { return boolStr(v.(*permissionV1.PolicyEvaluationLog).Result) }},
				{"IP", func(v any) string { return v.(*permissionV1.PolicyEvaluationLog).GetIpAddress() }},
				{"时间", func(v any) string { return tsStr(v.(*permissionV1.PolicyEvaluationLog).CreatedAt) }},
			},
		)
	default:
		http.Error(w, "unknown export type", http.StatusBadRequest)
	}
	return nil
}

// ServeAiUsageExport AI 用量流水导出：GET /admin/v1/ai/usage-logs:export
// 租户视角与 /admin/v1/ai/usage-logs 列表一致（core 侧按操作者裁定）。
func (s *AuditExportService) ServeAiUsageExport(w http.ResponseWriter, r *http.Request) error {
	ctx, ok := s.authorize(w, r)
	if !ok {
		return nil
	}
	q := r.URL.Query()
	s.exportType(ctx, w, q, "ai-usage-logs",
		func(ctx context.Context, page int64) ([]any, uint64, error) {
			resp, err := s.usageClient.List(ctx, pagingReq(page, q.Get("query")))
			if err != nil {
				return nil, 0, err
			}
			return toAnySlice(resp.GetItems()), resp.GetTotal(), nil
		},
		[]colDef{
			{"模型", func(v any) string { return v.(*aiV1.AiUsageLog).GetModelName() }},
			{"输入 tokens", func(v any) string { return strconv.FormatUint(uint64(v.(*aiV1.AiUsageLog).GetPromptTokens()), 10) }},
			{"输出 tokens", func(v any) string { return strconv.FormatUint(uint64(v.(*aiV1.AiUsageLog).GetCompletionTokens()), 10) }},
			{"总 tokens", func(v any) string { return strconv.FormatUint(uint64(v.(*aiV1.AiUsageLog).GetTotalTokens()), 10) }},
			{"耗时(ms)", func(v any) string { return strconv.FormatUint(uint64(v.(*aiV1.AiUsageLog).GetDurationMs()), 10) }},
			{"时间", func(v any) string { return tsStr(v.(*aiV1.AiUsageLog).CreatedAt) }},
		},
	)
	return nil
}

// ── 公共样板 ────────────────────────────────────────────────────────

// authorize 手动路由的鉴权：验 Bearer token 并把操作者身份写入出站上下文
// （core 侧 unary 中间件据此构建 viewer，租户隔离随之生效）。
func (s *AuditExportService) authorize(w http.ResponseWriter, r *http.Request) (context.Context, bool) {
	authHeader := r.Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" || token == authHeader {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return nil, false
	}
	valid, payload := s.tokenChecker.IsValidAccessToken(r.Context(), token, false)
	if !valid || payload == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return nil, false
	}

	ctx, err := appMetadata.NewContext(r.Context(), appMetadata.NewUserOperator(
		uint64(payload.GetUserId()), uint64(payload.GetTenantId()), 0, identityV1.DataScope_ALL,
	))
	if err != nil {
		http.Error(w, "identity error", http.StatusInternalServerError)
		return nil, false
	}
	return ctx, true
}

// exportType 通用导出：逐页聚合 → 按 format 渲染。
func (s *AuditExportService) exportType(
	ctx context.Context,
	w http.ResponseWriter,
	q url.Values,
	fileBase string,
	fetchPage func(ctx context.Context, page int64) ([]any, uint64, error),
	cols []colDef,
) {
	maxRows := exportMaxRows
	if mr, err := strconv.Atoi(q.Get("maxRows")); err == nil && mr > 0 && mr < exportMaxRows {
		maxRows = mr
	}

	headers := make([]string, len(cols))
	for i, c := range cols {
		headers[i] = c.header
	}

	var rows [][]string
	var total uint64
	for page := int64(1); int64(len(rows)) < int64(maxRows); page++ {
		items, t, err := fetchPage(ctx, page)
		if err != nil {
			s.log.Errorf("export [%s] page %d query failed: %v", fileBase, page, err)
			http.Error(w, "query failed", http.StatusInternalServerError)
			return
		}
		total = t
		if len(items) == 0 {
			break
		}
		for _, item := range items {
			if int64(len(rows)) >= int64(maxRows) {
				break
			}
			row := make([]string, len(cols))
			for i, c := range cols {
				row[i] = c.value(item)
			}
			rows = append(rows, row)
		}
		if uint64(page*exportPageSize) >= total {
			break
		}
	}

	format := strings.ToLower(q.Get("format"))
	if format == "csv" {
		writeCSV(w, fileBase, headers, rows)
		return
	}
	writeXLSX(w, fileBase, headers, rows)
}

func writeCSV(w http.ResponseWriter, fileBase string, headers []string, rows [][]string) {
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.csv"`, fileBase))
	_, _ = w.Write([]byte{0xEF, 0xBB, 0xBF})
	cw := csv.NewWriter(w)
	_ = cw.Write(headers)
	for _, row := range rows {
		_ = cw.Write(row)
	}
	cw.Flush()
}

func writeXLSX(w http.ResponseWriter, fileBase string, headers []string, rows [][]string) {
	f := excelize.NewFile()
	sheet := "Sheet1"
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		_ = f.SetCellValue(sheet, cell, h)
	}
	for r, row := range rows {
		for i, v := range row {
			cell, _ := excelize.CoordinatesToCellName(i+1, r+2)
			_ = f.SetCellValue(sheet, cell, v)
		}
	}
	body, err := f.WriteToBuffer()
	if err != nil {
		http.Error(w, "render xlsx failed", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.xlsx"`, fileBase))
	_, _ = w.Write(body.Bytes())
}

func pagingReq(page int64, query string) *paginationV1.PagingRequest {
	return &paginationV1.PagingRequest{
		Page:     trans.Ptr(uint32(page)),
		PageSize: trans.Ptr(uint32(exportPageSize)),
		FilteringType: &paginationV1.PagingRequest_Query{
			Query: query,
		},
	}
}

func toAnySlice[T any](xs []T) []any {
	out := make([]any, len(xs))
	for i, x := range xs {
		out[i] = x
	}
	return out
}

func boolStr(p *bool) string {
	if p == nil || !*p {
		return "否"
	}
	return "是"
}

func tsStr(t *timestamppb.Timestamp) string {
	if t == nil || t.AsTime().IsZero() {
		return ""
	}
	return t.AsTime().Format("2006-01-02 15:04:05")
}
