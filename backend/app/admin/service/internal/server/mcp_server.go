package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	mcptransport "github.com/tx7do/kratos-transport/transport/mcp"

	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
	socialV1 "go-wind-cms/api/gen/go/social/service/v1"

	"go-wind-cms/app/admin/service/internal/service"
)

// NewMcpServer MCP 服务器：把 CMS 的查询与发布触发能力以 MCP 工具暴露给
// AI 代理（Claude 等）。StreamableHTTP 传输。
//
// 鉴权模型：HTTP 层经 WithHTTPContextFunc 注入 Authorization 头（上游
// v1.3.4 的 WithMCPHTTPServerOptions 透传），各工具 handler 内用
// authorizeMcp 校验——tools/list 元数据对未鉴权请求可见（仅工具名与
// 描述），调用必须携带正确 Token。
//
// 环境变量：
//   - MCP_ADDR  监听地址（如 ":6602"）；未设置则不启动 MCP 服务器。
//     生产建议仅绑定内网地址或置于反向代理之后
//   - MCP_TOKEN 鉴权 Token（必填；缺失拒绝启动）
//
// 工具集：
//   - list_recent_posts           最近文章（标题/状态/ID）
//   - get_social_summary          分发 KPI 汇总
//   - list_social_distributions   分发记录（可按状态筛选）
//   - publish_distribution        触发一条分发立即发布（写操作）
func NewMcpServer(
	ctx *bootstrap.Context,
	postService *service.PostService,
	socialService *service.SocialService,
) *mcptransport.Server {
	addr := strings.TrimSpace(os.Getenv("MCP_ADDR"))
	if addr == "" {
		return nil
	}
	token := strings.TrimSpace(os.Getenv("MCP_TOKEN"))
	if token == "" {
		log.Warnf("MCP_ADDR set but MCP_TOKEN empty: refusing to start MCP server (auth required)")
		return nil
	}

	srv := mcptransport.NewServer(
		mcptransport.WithServerName("GoWind CMS"),
		mcptransport.WithServerVersion("1.1.0"),
		mcptransport.WithMCPServeType(mcptransport.ServerTypeHTTP),
		mcptransport.WithMCPServeAddress(addr),
		mcptransport.WithMCPServerOptions(server.WithRecovery()),
		mcptransport.WithMCPHTTPServerOptions(
			server.WithHTTPContextFunc(func(ctx context.Context, r *http.Request) context.Context {
				return context.WithValue(ctx, mcpTokenCtxKey{}, r.Header.Get("Authorization"))
			}),
		),
	)

	registerMcpTools(srv, token, postService, socialService)

	log.Infof("MCP server configured at %s (tools: list_recent_posts/get_social_summary/list_social_distributions/publish_distribution)", addr)
	return srv
}

// mcpTokenCtxKey Authorization 头的 ctx 注入键。
type mcpTokenCtxKey struct{}

// authorizeMcp 校验 MCP 调用者的 Bearer Token（与启动时 MCP_TOKEN 比对）。
func authorizeMcp(ctx context.Context, token string) error {
	if v, ok := ctx.Value(mcpTokenCtxKey{}).(string); ok && v == "Bearer "+token {
		return nil
	}
	return fmt.Errorf("unauthorized: missing or invalid bearer token")
}

// registerMcpTools 注册全部 MCP 工具。
func registerMcpTools(
	srv *mcptransport.Server,
	token string,
	postService *service.PostService,
	socialService *service.SocialService,
) {
	authz := func(ctx context.Context) error {
		return authorizeMcp(ctx, token)
	}

	// ── list_recent_posts ──
	listPostsTool := mcp.NewTool("list_recent_posts",
		mcp.WithDescription("列出 CMS 最近的文章（标题/状态/发布时间/ID）"),
		mcp.WithNumber("limit", mcp.Description("返回条数，缺省 10，至多 50")),
	)
	_ = srv.RegisterHandler(listPostsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if err := authz(ctx); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		resp, err := postService.List(ctx, &paginationV1.PagingRequest{
			Page:     trans.Ptr(uint32(1)),
			PageSize: trans.Ptr(uint32(clampLimit(request, 10))),
		})
		if err != nil {
			return mcp.NewToolResultError("list posts failed: " + err.Error()), nil
		}

		var b strings.Builder
		fmt.Fprintf(&b, "共 %d 篇文章，显示最近 %d 篇：\n", resp.GetTotal(), len(resp.GetItems()))
		for _, p := range resp.GetItems() {
			fmt.Fprintf(&b, "- [%s] #%d %s（创建于 %s）\n",
				p.GetStatus(), p.GetId(), firstPostTitleMcp(p), formatDateMcp(p.GetCreatedAt()))
		}
		return mcp.NewToolResultText(b.String()), nil
	})

	// ── get_social_summary ──
	summaryTool := mcp.NewTool("get_social_summary",
		mcp.WithDescription("获取社交分发 KPI 汇总：已发布总数、累计赞/评/转、近 7 天发布数"),
	)
	_ = srv.RegisterHandler(summaryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if err := authz(ctx); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		resp, err := socialService.GetSocialDistributionSummary(ctx, &emptypb.Empty{})
		if err != nil {
			return mcp.NewToolResultError("get summary failed: " + err.Error()), nil
		}
		text := fmt.Sprintf("已发布：%d\n累计点赞：%d\n累计回复：%d\n累计转发：%d\n近 7 天发布：%d",
			resp.GetTotalPublished(), resp.GetTotalLikes(), resp.GetTotalReplies(),
			resp.GetTotalReposts(), resp.GetPublishedLast_7Days())
		return mcp.NewToolResultText(text), nil
	})

	// ── list_social_distributions ──
	listDistTool := mcp.NewTool("list_social_distributions",
		mcp.WithDescription("列出社交分发记录（平台/状态/文案/链接）"),
		mcp.WithString("status",
			mcp.Description("按状态筛选（STATUS_DRAFT/STATUS_SCHEDULED/STATUS_PUBLISHING/STATUS_PUBLISHED/STATUS_FAILED），留空不过滤"),
			mcp.Enum("STATUS_DRAFT", "STATUS_SCHEDULED", "STATUS_PUBLISHING", "STATUS_PUBLISHED", "STATUS_FAILED"),
		),
		mcp.WithNumber("limit", mcp.Description("返回条数，缺省 10，至多 50")),
	)
	_ = srv.RegisterHandler(listDistTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if err := authz(ctx); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		req := &paginationV1.PagingRequest{
			Page:     trans.Ptr(uint32(1)),
			PageSize: trans.Ptr(uint32(clampLimit(request, 10))),
		}
		if status := request.GetString("status", ""); status != "" {
			// go-crud 裸键值即 EQ 精确匹配
			req.FilteringType = &paginationV1.PagingRequest_Query{
				Query: fmt.Sprintf(`{"status":%q}`, status),
			}
		}

		resp, err := socialService.ListSocialDistributions(ctx, req)
		if err != nil {
			return mcp.NewToolResultError("list distributions failed: " + err.Error()), nil
		}

		var b strings.Builder
		fmt.Fprintf(&b, "共 %d 条分发，显示最近 %d 条：\n", resp.GetTotal(), len(resp.GetItems()))
		for _, d := range resp.GetItems() {
			line := fmt.Sprintf("- [%s] %s → %s（账号：%s）",
				d.GetStatus(), strings.TrimPrefix(d.GetPlatform().String(), "PLATFORM_"),
				truncateRunesMcp(d.GetOverrideText(), 40), d.GetAccountName())
			if d.GetExternalUri() != "" {
				line += " " + d.GetExternalUri()
			}
			if d.GetErrorMessage() != "" {
				line += " 失败原因：" + d.GetErrorMessage()
			}
			fmt.Fprintln(&b, line)
		}
		return mcp.NewToolResultText(b.String()), nil
	})

	// ── publish_distribution ──
	publishTool := mcp.NewTool("publish_distribution",
		mcp.WithDescription("触发一条社交分发立即发布（写操作；仅草稿/已排期/失败状态可触发）"),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("分发记录 ID")),
	)
	_ = srv.RegisterHandler(publishTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if err := authz(ctx); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		idF, err := request.RequireFloat("id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		resp, perr := socialService.PublishSocialDistribution(ctx, &socialV1.PublishSocialDistributionRequest{
			Id: uint32(idF),
		})
		if perr != nil {
			return mcp.NewToolResultError("publish failed: " + perr.Error()), nil
		}
		text := fmt.Sprintf("已加入发布队列（当前状态：%s）。", resp.GetStatus())
		if resp.GetExternalUri() != "" {
			text += " 链接：" + resp.GetExternalUri()
		}
		return mcp.NewToolResultText(text), nil
	})
}

// ── 内部工具 ──

// clampLimit 从工具请求取 limit 参数（缺省 10，钳制 1..50）。
func clampLimit(request mcp.CallToolRequest, dft int) int {
	v, err := request.RequireFloat("limit")
	if err != nil || v <= 0 {
		return dft
	}
	n := int(v)
	if n > 50 {
		n = 50
	}
	return n
}

// firstPostTitleMcp 从文章的多语言翻译中取第一个非空标题。
func firstPostTitleMcp(post *contentV1.Post) string {
	for _, tr := range post.GetTranslations() {
		if t := tr.GetTitle(); t != "" {
			return t
		}
	}
	return ""
}

// formatDateMcp Timestamp → 展示时间。
func formatDateMcp(ts *timestamppb.Timestamp) string {
	if ts == nil || !ts.IsValid() {
		return "-"
	}
	return ts.AsTime().Format("2006-01-02 15:04")
}

// truncateRunesMcp 按 rune 截断。
func truncateRunesMcp(s string, max int) string {
	s = strings.TrimSpace(s)
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max]) + "…"
}
