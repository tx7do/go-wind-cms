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
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/go-utils/trans"
	authenticationV1 "go-wind-cms/api/gen/go/authentication/service/v1"

	"go-wind-cms/pkg/middleware/auth"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	bootstrapMcp "github.com/tx7do/kratos-bootstrap/transport/mcp"
	mcptransport "github.com/tx7do/kratos-transport/transport/mcp"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
	socialV1 "go-wind-cms/api/gen/go/social/service/v1"
	appViewer "go-wind-cms/pkg/entgo/viewer"

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
// 配置来源：config 的 server.mcp 段（conf.Server_MCP，经 bootstrap 的
// transport/mcp 适配器映射）优先；未配置回退环境变量 MCP_ADDR（兼容既有
// 部署，等价 HTTP + Recovery + 默认名/版本）。两处都没有（或显式写
// off/disabled）则不启动 MCP 服务器。生产建议仅绑定内网地址或置于反向
// 代理之后。
//   - server.mcp.address / MCP_ADDR  监听地址（如 ":6602"）
//   - MCP_TOKEN                      鉴权 Token（必填；缺失拒绝启动）。
//     凭据只走环境变量，不入配置文件
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
	commentService *service.CommentService,
) *mcptransport.Server {
	// 解析监听地址：yaml server.mcp.address 优先，env MCP_ADDR 兜底
	var mcpCfg *conf.Server_MCP
	if cfg := ctx.GetConfig(); cfg != nil {
		mcpCfg = cfg.GetServer().GetMcp()
	}
	var addr string
	if mcpCfg != nil {
		addr = strings.TrimSpace(mcpCfg.GetAddress())
	}
	if addr == "" {
		addr = strings.TrimSpace(os.Getenv("MCP_ADDR"))
		if addr != "" {
			// env 兜底保持旧行为：HTTP + Recovery，其余能力关闭
			mcpCfg = &conf.Server_MCP{
				ServerName:    "GoWind CMS",
				ServerVersion: "1.1.0",
				ServerType:    "HTTP",
				Address:       addr,
				Recovery:      true,
			}
		}
	}
	switch strings.ToLower(addr) {
	case "", "off", "disabled", "false", "0":
		// 未配置或显式停用：不启动（MCP 为可选组件）
		return nil
	}

	token := strings.TrimSpace(os.Getenv("MCP_TOKEN"))
	if token == "" {
		log.Warnf("MCP server enabled at %s but MCP_TOKEN empty: refusing to start MCP server (auth required)", addr)
		return nil
	}

	srv := bootstrapMcp.NewMcpServer(mcpCfg,
		mcptransport.WithMCPHTTPServerOptions(
			server.WithHTTPContextFunc(func(ctx context.Context, r *http.Request) context.Context {
				return context.WithValue(ctx, mcpTokenCtxKey{}, r.Header.Get("Authorization"))
			}),
		),
	)
	if srv == nil {
		return nil
	}

	registerMcpTools(srv, token, postService, socialService, commentService)

	log.Infof("MCP server configured at %s (tools: list_recent_posts/get_post/create_draft_post/list_comments/get_social_summary/list_social_distributions/publish_distribution)", addr)
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

// mcpAuthz 鉴权通过后注入机器操作者身份：BFF 层服务用 auth.FromContext 取
// 操作人（写入审计字段/租户归属），MCP 调用没有用户会话，统一以系统账号
// （userId=0、平台租户）落库——与系统通知 created_by=0 同口径。
func mcpAuthz(ctx context.Context, token string) (context.Context, error) {
	if err := authorizeMcp(ctx, token); err != nil {
		return nil, err
	}
	return auth.NewContext(ctx, &authenticationV1.UserTokenPayload{
		UserId:   0,
		TenantId: trans.Ptr(uint32(0)),
		Username: trans.Ptr("mcp"),
	}), nil
}

// registerMcpTools 注册全部 MCP 工具。
func registerMcpTools(
	srv *mcptransport.Server,
	token string,
	postService *service.PostService,
	socialService *service.SocialService,
	commentService *service.CommentService,
) {
	// ── list_recent_posts ──
	listPostsTool := mcp.NewTool("list_recent_posts",
		mcp.WithDescription("列出 CMS 最近的文章（标题/状态/发布时间/ID）"),
		mcp.WithNumber("limit", mcp.Description("返回条数，缺省 10，至多 50")),
	)
	_ = srv.RegisterHandler(listPostsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		ctx, aerr := mcpAuthz(ctx, token)
		if aerr != nil {
			return mcp.NewToolResultError(aerr.Error()), nil
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
		ctx, aerr := mcpAuthz(ctx, token)
		if aerr != nil {
			return mcp.NewToolResultError(aerr.Error()), nil
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
		ctx, aerr := mcpAuthz(ctx, token)
		if aerr != nil {
			return mcp.NewToolResultError(aerr.Error()), nil
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
		ctx, aerr := mcpAuthz(ctx, token)
		if aerr != nil {
			return mcp.NewToolResultError(aerr.Error()), nil
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

	// ── get_post ──
	getPostTool := mcp.NewTool("get_post",
		mcp.WithDescription("按 ID 读文章：标题/状态/发布时间/正文（Markdown 或 HTML）"),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("文章 ID")),
	)
	_ = srv.RegisterHandler(getPostTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		ctx, aerr := mcpAuthz(ctx, token)
		if aerr != nil {
			return mcp.NewToolResultError(aerr.Error()), nil
		}
		idF, err := request.RequireFloat("id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		post, gerr := postService.Get(ctx, &contentV1.GetPostRequest{
			QueryBy: &contentV1.GetPostRequest_Id{Id: uint32(idF)},
		})
		if gerr != nil {
			return mcp.NewToolResultError("get post failed: " + gerr.Error()), nil
		}
		title, content := "", ""
		for _, tr := range post.GetTranslations() {
			if title == "" && tr.GetTitle() != "" {
				title = tr.GetTitle()
			}
			if content == "" && tr.GetContent() != "" {
				content = tr.GetContent()
			}
		}
		var b strings.Builder
		fmt.Fprintf(&b, "#%d [%s] %s\n", post.GetId(), post.GetStatus(), title)
		if pt := formatDateMcp(post.GetPublishTime()); pt != "-" {
			b.WriteString("发布时间：" + pt + "\n")
		}
		b.WriteString("\n" + truncateRunesMcp(content, 4000))
		return mcp.NewToolResultText(b.String()), nil
	})

	// ── create_draft_post ──
	createDraftTool := mcp.NewTool("create_draft_post",
		mcp.WithDescription("创建一篇草稿文章（写操作；AI 起稿入口，落平台租户，状态 DRAFT，可在管理后台继续编辑/发布）"),
		mcp.WithString("title", mcp.Required(), mcp.Description("标题")),
		mcp.WithString("content", mcp.Required(), mcp.Description("正文（Markdown 或 HTML）")),
		mcp.WithString("summary", mcp.Description("摘要（可选）")),
	)
	_ = srv.RegisterHandler(createDraftTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		ctx, aerr := mcpAuthz(ctx, token)
		if aerr != nil {
			return mcp.NewToolResultError(aerr.Error()), nil
		}
		title, _ := request.RequireString("title")
		content, _ := request.RequireString("content")
		if strings.TrimSpace(title) == "" || strings.TrimSpace(content) == "" {
			return mcp.NewToolResultError("title and content are required"), nil
		}
		summary := request.GetString("summary", "")
		// MCP 调用无操作者会话：SystemViewer 落平台租户（tid=0）的 DRAFT
		sysCtx := appViewer.NewSystemViewerContext(ctx)
		post, cerr := postService.Create(sysCtx, &contentV1.CreatePostRequest{
			Data: &contentV1.Post{
				Status: contentV1.Post_POST_STATUS_DRAFT.Enum(),
				Translations: []*contentV1.PostTranslation{{
					LanguageCode: trans.Ptr("zh-CN"),
					Title:        trans.Ptr(title),
					Content:      trans.Ptr(content),
					Summary:      trans.Ptr(summary),
				}},
			},
		})
		if cerr != nil {
			return mcp.NewToolResultError("create draft failed: " + cerr.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("草稿已创建：#%d《%s》（管理后台可继续编辑与发布）", post.GetId(), title)), nil
	})

	// ── list_comments ──
	listCommentsTool := mcp.NewTool("list_comments",
		mcp.WithDescription("列出最近评论（作者/内容摘要/状态/所属对象）"),
		mcp.WithString("status",
			mcp.Description("按状态筛选（STATUS_PENDING/STATUS_APPROVED/STATUS_REJECTED），留空不过滤"),
			mcp.Enum("STATUS_PENDING", "STATUS_APPROVED", "STATUS_REJECTED"),
		),
		mcp.WithNumber("limit", mcp.Description("返回条数，缺省 10，至多 50")),
	)
	_ = srv.RegisterHandler(listCommentsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		ctx, aerr := mcpAuthz(ctx, token)
		if aerr != nil {
			return mcp.NewToolResultError(aerr.Error()), nil
		}
		req := &paginationV1.PagingRequest{
			Page:     trans.Ptr(uint32(1)),
			PageSize: trans.Ptr(uint32(clampLimit(request, 10))),
		}
		if status := request.GetString("status", ""); status != "" {
			req.FilteringType = &paginationV1.PagingRequest_Query{
				Query: fmt.Sprintf(`{"status":%q}`, status),
			}
		}
		resp, lerr := commentService.List(ctx, req)
		if lerr != nil {
			return mcp.NewToolResultError("list comments failed: " + lerr.Error()), nil
		}
		var b strings.Builder
		fmt.Fprintf(&b, "共 %d 条评论，显示最近 %d 条：\n", resp.GetTotal(), len(resp.GetItems()))
		for _, c := range resp.GetItems() {
			author := c.GetAuthorName()
			if author == "" {
				author = fmt.Sprintf("user#%d", c.GetAuthorId())
			}
			fmt.Fprintf(&b, "- [%s] %s（对象 #%d）：%s\n",
				c.GetStatus(), author, c.GetObjectId(), truncateRunesMcp(c.GetContent(), 60))
		}
		return mcp.NewToolResultText(b.String()), nil
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
