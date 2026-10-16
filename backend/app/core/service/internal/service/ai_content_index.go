package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/sashabaranov/go-openai"

	"go-wind-cms/app/core/service/internal/data"
	appViewer "go-wind-cms/pkg/entgo/viewer"
	"go-wind-cms/pkg/task"
)

// ── 帖子语义索引（asynq：ai.post.index） ─────────────────────────────
//
// 把帖子标题+摘要向量化写入 sys_ai_search_index（item_type='post'），
// 全局搜索语义分区可直达帖子编辑页。与菜单索引共用索引表：
//   - 菜单行：tenant_id=0（全局），权限裁剪按角色可见菜单集
//   - 帖子行：tenant_id=帖子归属租户，租户用户按 tenant 过滤、平台用户全量
// 由 PostService.enqueuePostReindex 的钩子在帖子增删改后入队（ai.post.index）。

// aiPostIndexBatch 重索引时分批送 /v1/embeddings 的条数。
const aiPostIndexBatch = 32

// aiPostIndexMax 单次全量重建最多索引的帖子数（防御无界查询；更多走增量）。
const aiPostIndexMax = 2000

// aiPostTextSQL 取单篇帖子的可索引文本（标题+摘要，标题优先中文翻译）。
const aiPostTextSQL = `SELECT p.tenant_id,
		COALESCE((SELECT pt.title || ' ' || COALESCE(pt.summary, '')
		          FROM post_translations pt
		          WHERE pt.post_id = p.id
		          ORDER BY (pt.language_code = 'zh-CN') DESC, pt.id
		          LIMIT 1), '') AS text
 FROM posts p WHERE p.id = $1`

// ensureSearchIndexColumns 幂等建表 + 补 tenant_id 列（老库升级）。
// 必须在任何按 tenant_id 过滤/写入的语句前调用。
func (s *AiContentService) ensureSearchIndexColumns(ctx context.Context) error {
	db := s.entClient.DB()
	if _, err := db.ExecContext(ctx, buildSearchIndexSQL); err != nil {
		return fmt.Errorf("create search index table: %w", err)
	}
	if _, err := db.ExecContext(ctx,
		`ALTER TABLE sys_ai_search_index ADD COLUMN IF NOT EXISTS tenant_id BIGINT NOT NULL DEFAULT 0`); err != nil {
		return fmt.Errorf("add tenant_id column: %w", err)
	}
	return nil
}

// IndexPostWorker asynq "ai.post.index" worker：op=index 重取帖子写索引，
// op=delete 移除索引行。入队由 PostService.enqueuePostReindex 钩子触发
// （帖子增删改后与 OpenSearch 同步双写）。
func (s *AiContentService) IndexPostWorker(taskType string, payload *task.AiPostIndexPayload) error {
	if payload == nil || payload.PostID == 0 {
		return nil
	}
	ctx := appViewer.NewSystemViewerContext(context.Background())

	if err := s.ensureSearchIndexColumns(ctx); err != nil {
		return err
	}

	if payload.Op == "delete" {
		_, err := s.entClient.DB().ExecContext(ctx,
			`DELETE FROM sys_ai_search_index WHERE item_type = 'post' AND item_id = $1`,
			payload.PostID)
		return err
	}

	// 取帖子（含租户/标题+摘要文本）；帖子已物理删除 → 等价 delete
	var tenantId int64
	var text string
	err := s.entClient.DB().QueryRowContext(ctx, aiPostTextSQL, payload.PostID).Scan(&tenantId, &text)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			_, dErr := s.entClient.DB().ExecContext(ctx,
				`DELETE FROM sys_ai_search_index WHERE item_type = 'post' AND item_id = $1`,
				payload.PostID)
			return dErr
		}
		return fmt.Errorf("fetch post for index: %w", err)
	}
	if strings.TrimSpace(text) == "" {
		return nil
	}

	provider, err := s.providerRepo.GetEnabledDefault(ctx)
	if err != nil {
		return err
	}
	client, err := newOpenAIClientForProvider(provider)
	if err != nil {
		return err
	}
	embCtx, cancel := context.WithTimeout(ctx, aiQueryCallTimeout)
	defer cancel()
	resp, embErr := client.CreateEmbeddings(embCtx, openai.EmbeddingRequest{
		Model: openai.EmbeddingModel(aiMenuEmbedModel),
		Input: []string{text},
	})
	if embErr != nil {
		return fmt.Errorf("embed post: %w", embErr)
	}
	if len(resp.Data) == 0 || len(resp.Data[0].Embedding) == 0 {
		return fmt.Errorf("empty embedding for post %d", payload.PostID)
	}
	if s.usageRepo != nil {
		_ = s.usageRepo.Create(ctx, &data.UsageRecord{
			TenantID:    uint32(tenantId),
			ModelName:   aiMenuEmbedModel,
			PromptTokens: uint32(resp.Usage.PromptTokens),
			TotalTokens: uint32(resp.Usage.TotalTokens),
		})
	}
	vec := vectorToLiteral(resp.Data[0].Embedding)

	if _, err = s.entClient.DB().ExecContext(ctx,
		`DELETE FROM sys_ai_search_index WHERE item_type = 'post' AND item_id = $1`,
		payload.PostID); err != nil {
		return err
	}
	route := fmt.Sprintf("/content/posts/edit/%d", payload.PostID)
	title := text
	// 按 rune 截断：len() 是字节长度，中文标题 200 字节只有 ~67 个字，
	// 用字节守卫 + rune 切片会越界 panic（曾致 ai.post.index 任务反复失败）
	if runes := []rune(title); len(runes) > 200 {
		title = string(runes[:200])
	}
	_, err = s.entClient.DB().ExecContext(ctx,
		`INSERT INTO sys_ai_search_index (created_at, item_type, item_id, title, route, tenant_id, embedding)
		 VALUES (NOW(), 'post', $1, $2, $3, $4, $5::vector)`,
		payload.PostID, title, route, tenantId, vec)
	return err
}

// BuildPostSearchIndex 全量重建帖子语义索引（分批 embedding，单事务替换）。
// 由 RebuildSearchIndex RPC 与菜单索引一起触发。返回索引的帖子数。
func (s *AiContentService) BuildPostSearchIndex(ctx context.Context) (int, error) {
	db := s.entClient.DB()
	if err := s.ensureSearchIndexColumns(ctx); err != nil {
		return 0, err
	}

	rows, err := db.QueryContext(ctx,
		`SELECT p.id, p.tenant_id,
		        COALESCE(MAX(pt.title) || ' ' || COALESCE(MAX(pt.summary), ''), '') AS text
		 FROM posts p LEFT JOIN post_translations pt ON pt.post_id = p.id
		 GROUP BY p.id, p.tenant_id
		 ORDER BY p.id DESC
		 LIMIT $1`, aiPostIndexMax)
	if err != nil {
		return 0, fmt.Errorf("fetch posts for index: %w", err)
	}
	type postIdxRow struct {
		id       uint32
		tenantId int64
		text     string
	}
	var posts []postIdxRow
	for rows.Next() {
		var r postIdxRow
		if scanErr := rows.Scan(&r.id, &r.tenantId, &r.text); scanErr != nil {
			continue
		}
		_ = r.tenantId
		if strings.TrimSpace(r.text) != "" {
			posts = append(posts, r)
		}
	}
	_ = rows.Close()
	if len(posts) == 0 {
		return 0, nil
	}

	provider, err := s.providerRepo.GetEnabledDefault(ctx)
	if err != nil {
		return 0, err
	}
	client, err := newOpenAIClientForProvider(provider)
	if err != nil {
		return 0, err
	}

	vectors := make([]string, len(posts))
	for start := 0; start < len(posts); start += aiPostIndexBatch {
		end := start + aiPostIndexBatch
		if end > len(posts) {
			end = len(posts)
		}
		texts := make([]string, 0, end-start)
		for _, p := range posts[start:end] {
			texts = append(texts, p.text)
		}
		resp, embErr := client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
			Model: openai.EmbeddingModel(aiMenuEmbedModel),
			Input: texts,
		})
		if embErr != nil {
			return 0, fmt.Errorf("embed batch: %w", embErr)
		}
		if placeErr := placeEmbeddingBatch(start, resp.Data, vectors); placeErr != nil {
			return 0, placeErr
		}
	}
	for i, p := range posts {
		if vectors[i] == "" {
			return 0, fmt.Errorf("post %d missing embedding", p.id)
		}
	}

	tx, txErr := db.BeginTx(ctx, nil)
	if txErr != nil {
		return 0, txErr
	}
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()
	if _, delErr := tx.ExecContext(ctx, `DELETE FROM sys_ai_search_index WHERE item_type = 'post'`); delErr != nil {
		return 0, delErr
	}
	const rowsPerStmt = 100
	for start := 0; start < len(posts); start += rowsPerStmt {
		end := start + rowsPerStmt
		if end > len(posts) {
			end = len(posts)
		}
		var sb strings.Builder
		var args []any
		sb.WriteString("INSERT INTO sys_ai_search_index (created_at, item_type, item_id, title, route, tenant_id, embedding) VALUES ")
		for j := start; j < end; j++ {
			if j > start {
				sb.WriteByte(',')
			}
			base := (j - start) * 5
			_, _ = fmt.Fprintf(&sb, "(NOW(),'post',$%d,$%d,$%d,$%d,$%d::vector)",
				base+1, base+2, base+3, base+4, base+5)
			args = append(args, posts[j].id, posts[j].text,
				fmt.Sprintf("/content/posts/edit/%d", posts[j].id), posts[j].tenantId, vectors[j])
		}
		if _, insErr := tx.ExecContext(ctx, sb.String(), args...); insErr != nil {
			return 0, insErr
		}
	}
	if cErr := tx.Commit(); cErr != nil {
		return 0, cErr
	}
	committed = true
	return len(posts), nil
}
