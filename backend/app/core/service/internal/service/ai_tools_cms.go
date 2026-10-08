package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/sashabaranov/go-openai"
)

// ── CMS 域工具：把内容数据暴露给对话的 Function Call 循环 ─────────────
//
// 工具只读：全部走固定 SQL（参数化、LIMIT 钳制），不接受模型拼 SQL。
// 租户边界：tenantId != 0 时强制 tenant_id 过滤。
// 表/列名与 ai_query_service.go 的白名单数据字典一致（经 information_schema 核对）。

// cmsToolDefs CMS 域工具清单（与内置工具合并后暴露给模型）。
func cmsToolDefs() []openai.Tool {
	return []openai.Tool{
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "query_recent_posts",
				Description: "查询帖子列表（含标题/状态/创建时间）。用户问「最近发了什么帖子」「列出草稿」等时调用。",
				Parameters: json.RawMessage(`{"type":"object","properties":{` +
					`"status":{"type":"string","description":"可选，帖子状态过滤：POST_STATUS_DRAFT/POST_STATUS_PUBLISHED/POST_STATUS_SCHEDULED/POST_STATUS_TRASHED"},` +
					`"limit":{"type":"string","description":"可选，返回条数，默认 10，最大 50"}},` +
					`"required":[]}`),
			},
		},
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "count_posts_by_status",
				Description: "按状态统计帖子数量（草稿/已发布/定时/回收站各多少）。用户问内容量、发布情况时调用。",
				Parameters:  json.RawMessage(`{"type":"object","properties":{},"required":[]}`),
			},
		},
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "get_content_overview",
				Description: "站点内容概览：用户总数、帖子总数、评论总数、媒体文件总数。用户问站点规模/整体数据时调用。",
				Parameters:  json.RawMessage(`{"type":"object","properties":{},"required":[]}`),
			},
		},
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "query_comments",
				Description: "查询评论统计（按状态分组：PUBLISHED已发布/PENDING待审/TRASHED回收站）。用户问评论情况、待审数量时调用。",
				Parameters:  json.RawMessage(`{"type":"object","properties":{},"required":[]}`),
			},
		},
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "query_top_posts",
				Description: "查询互动量最高的帖子 TOP N（按评论数排序）。用户问热门帖子/互动排行时调用。",
				Parameters: json.RawMessage(`{"type":"object","properties":{` +
					`"limit":{"type":"string","description":"可选，返回条数，默认 5，最大 20"}},` +
					`"required":[]}`),
			},
		},
	}
}

// execCMSTool 执行一个 CMS 工具调用。返回 ok=false 表示非 CMS 工具（回退内置执行器）。
// 未知参数/执行错误返回错误文本（回传模型让它自行纠正，而非中断整轮）。
func (s *AiChatService) execCMSTool(ctx context.Context, tenantId uint32, name, args string) (string, bool, error) {
	var p map[string]any
	if args != "" {
		_ = json.Unmarshal([]byte(args), &p)
	}
	limit := 10
	if v, ok := p["limit"].(string); ok {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 50 {
			limit = n
		}
	}
	status, _ := p["status"].(string)

	db := s.entClient.DB()
	switch name {
	case "query_recent_posts":
		sqlText := `SELECT p.id, COALESCE(MAX(pt.title), '') AS title, p.status,
			TO_CHAR(MAX(p.created_at), 'YYYY-MM-DD HH24:MI') AS created_at
			FROM posts p LEFT JOIN post_translations pt ON pt.post_id = p.id`
		var args []any
		where := ""
		if status != "" {
			where += " WHERE p.status = $" + strconv.Itoa(len(args)+1)
			args = append(args, status)
		}
		if tenantId != 0 {
			where += " WHERE p.tenant_id = $" + strconv.Itoa(len(args)+1)
			args = append(args, tenantId)
		}
		sqlText += where + " GROUP BY p.id, p.status ORDER BY MAX(p.created_at) DESC LIMIT $" + strconv.Itoa(len(args)+1)
		args = append(args, limit)
		out, err := s.queryRowsAsString(ctx, db, sqlText, args...)
		if err != nil {
			return "", true, err
		}
		if out == "" {
			return "没有符合条件的帖子", true, nil
		}
		return out, true, nil

	case "count_posts_by_status":
		out, err := s.queryRowsAsString(ctx, db,
			`SELECT status, COUNT(*) FROM posts GROUP BY status ORDER BY status`)
		if err != nil {
			return "", true, err
		}
		return out, true, nil

	case "get_content_overview":
		out, err := s.queryRowsAsString(ctx, db,
			`SELECT '用户' AS item, COUNT(*) FROM sys_users
			 UNION ALL SELECT '帖子', COUNT(*) FROM posts
			 UNION ALL SELECT '评论', COUNT(*) FROM comments
			 UNION ALL SELECT '媒体文件', COUNT(*) FROM media_assets`)
		if err != nil {
			return "", true, err
		}
		return out, true, nil

	case "query_comments":
		out, err := s.queryRowsAsString(ctx, db,
			`SELECT status, COUNT(*) FROM comments GROUP BY status ORDER BY status`)
		if err != nil {
			return "", true, err
		}
		return out, true, nil

	case "query_top_posts":
		topN := 5
		if v, ok := p["limit"].(string); ok {
			if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 20 {
				topN = n
			}
		}
		out, err := s.queryRowsAsString(ctx, db,
			`SELECT c.object_id AS post_id, COALESCE(MAX(pt.title), '') AS title, COUNT(*) AS comment_count
			 FROM comments c LEFT JOIN post_translations pt ON pt.post_id = c.object_id
			 WHERE c.content_type = 'POST' OR c.content_type = 'post'
			 GROUP BY c.object_id ORDER BY comment_count DESC LIMIT $1`, topN)
		if err != nil {
			return "", true, err
		}
		return out, true, nil

	default:
		return "", false, fmt.Errorf("unknown cms tool: %s", name)
	}
}

// queryRowsAsString 执行查询并把结果集渲染成 LLM 可读的竖线分隔文本。
func (s *AiChatService) queryRowsAsString(ctx context.Context, db *sql.DB, query string, args ...any) (string, error) {
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return "", err
	}
	defer func() { _ = rows.Close() }()

	columns, err := rows.Columns()
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	sb.WriteString(strings.Join(columns, " | ") + "\n")
	count := 0
	for rows.Next() {
		values := make([]interface{}, len(columns))
		ptrs := make([]interface{}, len(columns))
		for i := range values {
			ptrs[i] = &values[i]
		}
		if err = rows.Scan(ptrs...); err != nil {
			return "", err
		}
		cells := make([]string, len(columns))
		for i, v := range values {
			cells[i] = stringifyCellValue(v)
		}
		sb.WriteString(strings.Join(cells, " | ") + "\n")
		count++
		if count >= 50 {
			break
		}
	}
	if err = rows.Err(); err != nil {
		return "", err
	}
	return strings.TrimRight(sb.String(), "\n"), nil
}
