package service

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/sashabaranov/go-openai"
	"github.com/tx7do/go-crud/viewer"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	entCrud "github.com/tx7do/go-crud/entgo"

	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"

	aiV1 "go-wind-cms/api/gen/go/ai/service/v1"
	identityV1 "go-wind-cms/api/gen/go/identity/service/v1"
	permissionV1 "go-wind-cms/api/gen/go/permission/service/v1"
	"go-wind-cms/app/core/service/internal/data"
	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/aiprovider"
	pkgAi "go-wind-cms/pkg/ai"
	"go-wind-cms/pkg/crypto"
)

const aiContentCallTimeout = 60 * time.Second

// 场景 → 系统提示词模板（%s = langName(lang)）。
// 1-4 与上游对齐；5-7 为 CMS 内容场景（帖子标题/摘要/正文）。
var aiContentScenePrompts = map[aiV1.ContentScene]string{
	aiV1.ContentScene_DESCRIPTION:  "你是企业后台的内容撰写助手。根据用户给出的主题，为后台管理系统的实体（角色/菜单/字典项/文件等）撰写一段简明的描述文字。要求：专业、简洁、准确，不超过 100 字。用%s输出，只输出正文。",
	aiV1.ContentScene_ANNOUNCEMENT: "你是企业后台的公告撰写助手。根据用户给出的主题，撰写一则面向全体员工的通知公告。要求：包含标题行、正文（背景/内容/要求），语气正式但不生硬，总长度不超过 300 字。用%s输出，只输出正文。",
	aiV1.ContentScene_REPLY:        "你是企业后台的客服回复助手。根据用户给出的问题或场景，撰写一段专业的回复文案。要求：语气友好、有条理、针对性强，不超过 200 字。用%s输出，只输出正文。",
	aiV1.ContentScene_GENERAL:      "你是企业后台的智能写作助手。根据用户的指令生成内容。要求：直接输出结果，不要解释你的思考过程。用%s输出。",
	aiV1.ContentScene_POST_TITLE:   "你是CMS编辑部的标题撰搞人。根据用户给出的文章主题或草稿要点，拟一个吸引人但不标题党的中文标题。要求：不超过 30 字，准确概括主题，可以直接输出一个标题，不要书名号、不要解释。用%s输出。",
	aiV1.ContentScene_POST_SUMMARY: "你是CMS编辑部的摘要撰搞人。根据用户给出的文章主题或正文要点，撰写一段文章摘要。要求：120 字以内，概括核心内容与读者价值，平实自然，不要以「本文」「这篇文章」开头。用%s输出，只输出摘要正文。",
	aiV1.ContentScene_POST_BODY:    "你是CMS的正文撰搞助手。根据用户给出的主题与要点撰写文章正文。要求：结构清晰（可用短段落或小标题），总长度控制在 800 字以内，直接输出正文，不要解释你的思考过程。用%s输出。",
	aiV1.ContentScene_SOCIAL_POST:  "你是社交媒体运营文案助手。根据用户给出的主题或文章要点，撰写一条适合发布到社交平台的短文案。要求：口语化、有吸引力、能引发点击或互动，控制在 200 字以内，可酌情带 1-3 个话题标签（#标签），不要解释。用%s输出，只输出文案。",
}

// AiContentService AI 内容生成（表单助手）。
type AiContentService struct {
	aiV1.UnimplementedAiContentServiceServer
	log          *log.Helper
	providerRepo *data.AiProviderRepo
	usageRepo    *data.AiUsageLogRepo
	knowledgeRepo *data.AiKnowledgeRepo
	entClient    *entCrud.EntClient[*ent.Client]
	menuRepo     *data.MenuRepo
	roleRepo     interface {
		GetRolesPermissionMenuIDs(ctx context.Context, roleIDs []uint32) ([]uint32, error)
	}
	userRepo interface {
		Get(ctx context.Context, req *identityV1.GetUserRequest) (*identityV1.User, error)
	}
}

func NewAiContentService(
	ctx *bootstrap.Context,
	providerRepo *data.AiProviderRepo,
	usageRepo *data.AiUsageLogRepo,
	knowledgeRepo *data.AiKnowledgeRepo,
	entClient *entCrud.EntClient[*ent.Client],
	menuRepo *data.MenuRepo,
	roleRepo *data.RoleRepo,
	userRepo data.UserRepo,
) *AiContentService {
	return &AiContentService{
		log:           ctx.NewLoggerHelper("ai_content/service/core-service"),
		providerRepo:  providerRepo,
		usageRepo:     usageRepo,
		knowledgeRepo: knowledgeRepo,
		entClient:     entClient,
		menuRepo:      menuRepo,
		roleRepo:      roleRepo,
		userRepo:      userRepo,
	}
}

// GenerateContent 按场景生成内容：登录用户均可使用。
// 前端在表单字段旁放 ✨ 按钮，点击后传入场景+主题，生成结果由前端自行填入字段。
func (s *AiContentService) GenerateContent(ctx context.Context, req *aiV1.GenerateContentRequest) (*aiV1.GenerateContentResponse, error) {
	topic := strings.TrimSpace(req.GetTopic())
	if topic == "" {
		return nil, aiV1.ErrorBadRequest("topic is required")
	}

	lang := req.GetLang()
	if lang == "" {
		lang = "zh-CN"
	}
	langDisplay := langName(lang)

	scene := req.GetScene()
	systemPrompt, ok := aiContentScenePrompts[scene]
	if !ok {
		systemPrompt = aiContentScenePrompts[aiV1.ContentScene_GENERAL]
	}
	systemPrompt = fmt.Sprintf(systemPrompt, langDisplay)

	// 补充上下文作为追加指令
	if extra := strings.TrimSpace(req.GetContext()); extra != "" {
		systemPrompt += "\n补充要求：" + extra
	}
	if req.GetMaxLength() > 0 {
		systemPrompt += fmt.Sprintf("\n输出长度不超过 %d 字。", req.GetMaxLength())
	}

	provider, err := s.providerRepo.GetEnabledDefault(ctx)
	if err != nil {
		return nil, err
	}
	client, err := newOpenAIClientForProvider(provider)
	if err != nil {
		return nil, err
	}

	callStart := time.Now()
	cCtx, cancel := context.WithTimeout(ctx, aiContentCallTimeout)
	defer cancel()
	resp, err := client.CreateChatCompletion(cCtx, openai.ChatCompletionRequest{
		Model: ptrStrOr(provider.ModelName, "gpt-4o-mini"),
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: systemPrompt},
			{Role: openai.ChatMessageRoleUser, Content: topic},
		},
	})
	if err != nil {
		s.log.Errorf("ai content generate failed: scene=%s: %v", scene.String(), err)
		return nil, aiV1.ErrorInternalServerError("ai content generate failed: %v", err)
	}

	content := strings.TrimSpace(resp.Choices[0].Message.Content)
	tokens := uint32(resp.Usage.TotalTokens)
	durationMs := uint32(time.Since(callStart).Milliseconds())

	// 记用量（尽力而为，不阻断返回）
	if s.usageRepo != nil {
		var userId, tenantId uint32
		if vc, exist := viewer.FromContext(ctx); exist && vc != nil {
			tenantId = uint32(vc.TenantID())
		}
		if uid, hasUser := viewerUserIDFromContext(ctx); hasUser {
			userId = uid
		}
		if uerr := s.usageRepo.Create(ctx, &data.UsageRecord{
			ProviderID:       provider.ID,
			UserID:           userId,
			TenantID:         tenantId,
			ModelName:        ptrStrOr(provider.ModelName, ""),
			PromptTokens:     uint32(resp.Usage.PromptTokens),
			CompletionTokens: uint32(resp.Usage.CompletionTokens),
			TotalTokens:      uint32(resp.Usage.TotalTokens),
			DurationMs:       durationMs,
		}); uerr != nil {
			s.log.Errorf("ai content usage log failed: %v", uerr)
		}
	}

	return &aiV1.GenerateContentResponse{
		Content:     content,
		TotalTokens: tokens,
	}, nil
}

// newOpenAIClientForProvider 从 provider 表行构造客户端（api_key 密文就地解密，
// 不出 service 层）。
func newOpenAIClientForProvider(provider *ent.AiProvider) (*openai.Client, error) {
	cfg := &pkgAi.ClientConfig{
		ModelType:      modelTypeToProtoInt(provider.ModelType),
		ModelName:      ptrStrOr(provider.ModelName, ""),
		BaseUrl:        ptrStrOr(provider.BaseURL, ""),
		Organization:   ptrStrOr(provider.Organization, ""),
		LocalHost:      ptrStrOr(provider.LocalHost, ""),
		LocalPort:      int32(ptrIntOr(provider.LocalPort, 0)),
		TimeoutSeconds: int32(ptrIntOr(provider.TimeoutSeconds, 0)),
	}

	var apiKey string
	if pkgAi.ModelType(cfg.ModelType) == pkgAi.ModelTypeCloud {
		encrypted := ptrStrOr(provider.APIKey, "")
		if encrypted == "" {
			return nil, aiV1.ErrorBadRequest("ai provider api key is not configured")
		}
		plain, err := crypto.DecryptIfNeeded(encrypted)
		if err != nil {
			return nil, aiV1.ErrorInternalServerError("decrypt api key failed")
		}
		apiKey = plain
	}

	client, err := pkgAi.NewClient(cfg, apiKey)
	if err != nil {
		return nil, aiV1.ErrorInternalServerError("create ai client failed")
	}
	return client, nil
}

// modelTypeToProtoInt ent 枚举（string 底层）→ proto ModelType 数值。
// 两端枚举命名逐字一致（LOCAL/CLOUD），这里只做值域对齐。
func modelTypeToProtoInt(mt *aiprovider.ModelType) int32 {
	if mt == nil {
		return int32(aiV1.AiProvider_MODEL_TYPE_UNSPECIFIED)
	}
	switch *mt {
	case aiprovider.ModelTypeLOCAL:
		return int32(aiV1.AiProvider_LOCAL)
	case aiprovider.ModelTypeCLOUD:
		return int32(aiV1.AiProvider_CLOUD)
	default:
		return int32(aiV1.AiProvider_MODEL_TYPE_UNSPECIFIED)
	}
}

// langName BCP-47 → 提示词里的语言名。
func langName(lang string) string {
	if strings.HasPrefix(lang, "en") {
		return "English"
	}
	return "中文"
}

func ptrStrOr(p *string, def string) string {
	if p == nil || *p == "" {
		return def
	}
	return *p
}

func ptrIntOr(p *int, def int) int {
	if p == nil {
		return def
	}
	return *p
}

// ── 语义搜索（pgvector 菜单索引；RAG 之外的第二处向量消费方） ─────────

// aiMenuEmbedModel 菜单索引使用的 embedding 模型（与维度常量配套，换模型须重索引）。
const aiMenuEmbedModel = "text-embedding-3-small"

// buildSearchIndexSQL 菜单语义搜索索引表。embedding 定维（见 data.AiEmbeddingDimensions）：
// 混维向量会让 <=> 比较在查询期整体报错，定维把失败提前到写入期。
var buildSearchIndexSQL = fmt.Sprintf(
	"CREATE TABLE IF NOT EXISTS sys_ai_search_index ("+
		"id BIGSERIAL PRIMARY KEY, "+
		"created_at TIMESTAMPTZ DEFAULT NOW(), "+
		"item_type VARCHAR(20) NOT NULL DEFAULT 'menu', "+
		"item_id BIGINT NOT NULL, "+
		"title TEXT NOT NULL, "+
		"route TEXT NOT NULL DEFAULT '', "+
		"embedding vector(%d))",
	data.AiEmbeddingDimensions)

// menuIndexDoc 待索引菜单：真实菜单 id + 标题 + 路由。
type menuIndexDoc struct {
	menuId uint32
	title  string
	route  string
}

// collectMenuIndexDocs 从菜单列表收集可索引文档：跳过无 id / 无 path /
// 无标题的菜单；标题取 meta.title，为空回退 name。
func collectMenuIndexDocs(items []*permissionV1.Menu) []menuIndexDoc {
	docs := make([]menuIndexDoc, 0, len(items))
	for _, m := range items {
		if m.GetId() == 0 || m.GetPath() == "" {
			continue
		}
		title := ""
		if mt := m.GetMeta(); mt != nil {
			title = mt.GetTitle()
		}
		if title == "" {
			title = m.GetName()
		}
		if title == "" {
			continue
		}
		docs = append(docs, menuIndexDoc{menuId: m.GetId(), title: title, route: m.GetPath()})
	}
	return docs
}

// placeEmbeddingBatch 把一批 embedding 响应按其 index（对位请求输入序号）
// 写入 vectors 的 [batchStart, batchStart+len) 段：拒绝越界 index 与零维向量。
func placeEmbeddingBatch(batchStart int, data []openai.Embedding, vectors []string) error {
	for _, item := range data {
		if item.Index < 0 || batchStart+item.Index >= len(vectors) || len(item.Embedding) == 0 {
			return fmt.Errorf("embedding response malformed: batch=%d item=%d", batchStart, item.Index)
		}
		vectors[batchStart+item.Index] = vectorToLiteral(item.Embedding)
	}
	return nil
}

// paginationV1PagingNoLimit 构造不分页的 PagingRequest。
func paginationV1PagingNoLimit() *paginationV1.PagingRequest {
	return &paginationV1.PagingRequest{NoPaging: trans.Ptr(true)}
}

// BuildMenuSearchIndex 向量化全部菜单标题并原子替换索引表里的菜单行。
// 顺序：先离线算完全部 embedding（不持事务调外部 API），再单事务 DELETE+分批
// 多行 INSERT——任一环节失败即回滚，旧索引原样保留。
func (s *AiContentService) BuildMenuSearchIndex(ctx context.Context) (uint32, error) {
	db := s.entClient.DB()
	if _, err := db.ExecContext(ctx, buildSearchIndexSQL); err != nil {
		return 0, fmt.Errorf("create search index table: %w", err)
	}
	if _, err := data.EnsureVectorColumnDim(ctx, db, "sys_ai_search_index"); err != nil {
		return 0, fmt.Errorf("ensure embedding column: %w", err)
	}
	operatorId, tenantId := s.operatorIds(ctx)
	menus, err := s.menuRepo.List(ctx, paginationV1PagingNoLimit(), true)
	if err != nil {
		return 0, err
	}

	docs := collectMenuIndexDocs(menus.GetItems())
	if len(docs) == 0 {
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

	// 离线分批算完全部向量；消耗记入用量流水（记量失败不阻断）。
	vectors := make([]string, len(docs))
	embCtx, embCancel := context.WithTimeout(ctx, 60*time.Second)
	defer embCancel()
	for start := 0; start < len(docs); start += 32 {
		end := min(start+32, len(docs))
		batch := make([]string, 0, end-start)
		for _, d := range docs[start:end] {
			batch = append(batch, d.title)
		}
		callStart := time.Now()
		resp, embErr := client.CreateEmbeddings(embCtx, openai.EmbeddingRequest{
			Model: openai.EmbeddingModel(aiMenuEmbedModel),
			Input: batch,
		})
		if embErr != nil {
			return 0, fmt.Errorf("embed batch: %w", embErr)
		}
		if placeErr := placeEmbeddingBatch(start, resp.Data, vectors); placeErr != nil {
			return 0, placeErr
		}
		if s.usageRepo != nil {
			_ = callStart
			if uerr := s.usageRepo.Create(ctx, &data.UsageRecord{
				UserID:       operatorId,
				TenantID:     tenantId,
				ModelName:    aiMenuEmbedModel,
				PromptTokens: uint32(resp.Usage.PromptTokens),
				TotalTokens:  uint32(resp.Usage.TotalTokens),
				DurationMs:   uint32(time.Since(callStart).Milliseconds()),
			}); uerr != nil {
				s.log.Errorf("record embedding usage failed: %v", uerr)
			}
		}
	}
	for i := range vectors {
		if vectors[i] == "" {
			return 0, fmt.Errorf("menu %d missing embedding response", docs[i].menuId)
		}
	}

	tx, txErr := db.BeginTx(ctx, nil)
	if txErr != nil {
		return 0, fmt.Errorf("begin tx: %w", txErr)
	}
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	if _, delErr := tx.ExecContext(ctx, "DELETE FROM sys_ai_search_index WHERE item_type = 'menu'"); delErr != nil {
		return 0, fmt.Errorf("clear search index: %w", delErr)
	}

	const rowsPerStmt = 100
	for start := 0; start < len(docs); start += rowsPerStmt {
		end := start + rowsPerStmt
		if end > len(docs) {
			end = len(docs)
		}
		var sb strings.Builder
		var args []any
		sb.WriteString("INSERT INTO sys_ai_search_index (item_type, item_id, title, route, embedding) VALUES ")
		for j := start; j < end; j++ {
			if j > start {
				sb.WriteByte(',')
			}
			base := (j - start) * 5
			_, _ = fmt.Fprintf(&sb, "($%d,$%d,$%d,$%d,$%d::vector)", base+1, base+2, base+3, base+4, base+5)
			args = append(args, "menu", int64(docs[j].menuId), docs[j].title, docs[j].route, vectors[j])
		}
		if _, insErr := tx.ExecContext(ctx, sb.String(), args...); insErr != nil {
			return 0, fmt.Errorf("insert search index rows: %w", insErr)
		}
	}

	if cErr := tx.Commit(); cErr != nil {
		return 0, fmt.Errorf("commit search index: %w", cErr)
	}
	committed = true
	return uint32(len(docs)), nil
}

// SemanticSearch 语义搜索：向量化 query 并按余弦相似度检索菜单索引。
// 非平台用户按其可见菜单集过滤（fail-closed），防止越权菜单经搜索泄露。
func (s *AiContentService) SemanticSearch(ctx context.Context, req *aiV1.SemanticSearchRequest) (*aiV1.SemanticSearchResponse, error) {
	operatorId, tenantId := s.operatorIds(ctx)
	query := strings.TrimSpace(req.GetQuery())
	if query == "" {
		return &aiV1.SemanticSearchResponse{}, nil
	}
	limit := int(req.GetLimit())
	if limit <= 0 {
		limit = 8
	}

	// 权限裁剪前置：非平台用户按其可见菜单集过滤索引行（item_id = menu_id）。
	allowedMenuIDs, filterErr := s.visibleMenuIDs(ctx, operatorId, tenantId)
	if filterErr != nil {
		s.log.Errorf("semantic search: resolve visible menus failed: %v", filterErr)
		return nil, aiV1.ErrorInternalServerError("resolve accessible menus failed")
	}
	if tenantId != 0 && len(allowedMenuIDs) == 0 {
		return &aiV1.SemanticSearchResponse{Items: []*aiV1.SemanticSearchItem{}}, nil
	}

	// 懒自建：索引表为空（首次使用/换库）时先重建再检索，避免"永远静默降级"。
	// 仅平台路径触发（租户首查等平台侧已建好）；构建失败降级为空结果不报错。
	if tenantId == 0 {
		if err := s.ensureSearchIndexColumns(ctx); err == nil {
			var idxCount int
			if cntErr := s.entClient.DB().QueryRowContext(ctx,
				"SELECT COUNT(*) FROM sys_ai_search_index").Scan(&idxCount); cntErr == nil && idxCount == 0 {
				if _, bErr := s.BuildMenuSearchIndex(ctx); bErr != nil {
					s.log.Errorf("lazy build search index failed, degrade to empty: %v", bErr)
				} else if _, bErr := s.BuildPostSearchIndex(ctx); bErr != nil {
					s.log.Errorf("lazy build post index failed: %v", bErr)
				}
			}
		}
	}

	provider, err := s.providerRepo.GetEnabledDefault(ctx)
	if err != nil {
		return nil, err
	}
	client, err := newOpenAIClientForProvider(provider)
	if err != nil {
		return nil, err
	}

	embCtx, embCancel := context.WithTimeout(ctx, 30*time.Second)
	defer embCancel()
	embResp, embErr := client.CreateEmbeddings(embCtx, openai.EmbeddingRequest{
		Model: openai.EmbeddingModel(aiMenuEmbedModel),
		Input: []string{query},
	})
	if embErr != nil {
		return nil, aiV1.ErrorInternalServerError("embed query failed: %v", embErr)
	}
	if len(embResp.Data) == 0 || len(embResp.Data[0].Embedding) == 0 {
		return nil, aiV1.ErrorInternalServerError("embed query returned empty embedding")
	}
	queryVec := vectorToLiteral(embResp.Data[0].Embedding)

	if s.usageRepo != nil {
		if uerr := s.usageRepo.Create(ctx, &data.UsageRecord{
			UserID:       operatorId,
			TenantID:     tenantId,
			ModelName:    aiMenuEmbedModel,
			PromptTokens: uint32(embResp.Usage.PromptTokens),
			TotalTokens:  uint32(embResp.Usage.TotalTokens),
		}); uerr != nil {
			s.log.Errorf("record embedding usage failed: %v", uerr)
		}
	}

	db := s.entClient.DB()
	var rows *sql.Rows
	var dbErr error
	if tenantId == 0 {
		rows, dbErr = db.QueryContext(ctx,
			"SELECT title, route, item_type, 1 - (embedding <=> $1::vector) AS score FROM sys_ai_search_index WHERE item_type IN ('menu','post') ORDER BY embedding <=> $1::vector LIMIT $2",
			queryVec, limit,
		)
	} else {
		// 菜单行按角色可见菜单集裁剪；帖子行按本租户过滤
		placeholders := make([]string, len(allowedMenuIDs))
		args := make([]any, 0, len(allowedMenuIDs)+3)
		args = append(args, queryVec)
		for i, id := range allowedMenuIDs {
			placeholders[i] = fmt.Sprintf("$%d", i+2)
			args = append(args, id)
		}
		args = append(args, limit)
		args = append(args, tenantId)
		sqlText := fmt.Sprintf(
			"SELECT title, route, item_type, 1 - (embedding <=> $1::vector) AS score FROM sys_ai_search_index WHERE "+
				"((item_type = 'menu' AND item_id IN (%s)) OR (item_type = 'post' AND tenant_id = $%d)) "+
				"ORDER BY embedding <=> $1::vector LIMIT $%d",
			strings.Join(placeholders, ", "), len(args), len(args)+1,
		)
		rows, dbErr = db.QueryContext(ctx, sqlText, args...)
	}
	if dbErr != nil {
		s.log.Errorf("semantic search query failed: %s", dbErr.Error())
		return nil, aiV1.ErrorInternalServerError("semantic search failed")
	}
	defer func() { _ = rows.Close() }()

	resp := &aiV1.SemanticSearchResponse{Items: make([]*aiV1.SemanticSearchItem, 0, limit)}
	for rows.Next() {
		var title, route, itemType string
		var score float64
		if err = rows.Scan(&title, &route, &itemType, &score); err != nil {
			continue
		}
		resp.Items = append(resp.Items, &aiV1.SemanticSearchItem{
			Title: title, Route: route, ItemType: itemType, Score: float32(score),
		})
	}
	return resp, rows.Err()
}

// RebuildSearchIndex 重建搜索索引 RPC。
func (s *AiContentService) RebuildSearchIndex(ctx context.Context, _ *emptypb.Empty) (*aiV1.RebuildSearchIndexResponse, error) {
	count, err := s.BuildMenuSearchIndex(ctx)
	if err != nil {
		return nil, aiV1.ErrorInternalServerError("rebuild search index failed: %v", err)
	}
	// 帖子索引失败不回滚菜单索引——部分可用优于全不可用
	postCount, pErr := s.BuildPostSearchIndex(ctx)
	if pErr != nil {
		s.log.Errorf("rebuild post index failed: %v", pErr)
	}
	return &aiV1.RebuildSearchIndexResponse{IndexedCount: count + uint32(postCount)}, nil
}

// visibleMenuIDs 解析调用者可见的菜单 ID 集。
// 平台用户返回 nil（不过滤）；普通用户经 用户→角色→权限→菜单 链路解析，
// 解析失败返回错误（fail-closed）。
func (s *AiContentService) visibleMenuIDs(ctx context.Context, operatorId, tenantId uint32) ([]uint32, error) {
	if tenantId == 0 {
		return nil, nil
	}
	user, err := s.userRepo.Get(ctx, &identityV1.GetUserRequest{
		QueryBy: &identityV1.GetUserRequest_Id{Id: operatorId},
	})
	if err != nil {
		return nil, err
	}
	if user == nil || len(user.GetRoleIds()) == 0 {
		return []uint32{}, nil
	}
	return s.roleRepo.GetRolesPermissionMenuIDs(ctx, user.GetRoleIds())
}

// operatorIds 从 viewer 上下文取操作者（unary 中间件注入）。
func (s *AiContentService) operatorIds(ctx context.Context) (uint32, uint32) {
	var userId, tenantId uint32
	if vc, exist := viewer.FromContext(ctx); exist && vc != nil {
		tenantId = uint32(vc.TenantID())
	}
	if uid, ok := viewerUserIDFromContext(ctx); ok {
		userId = uid
	}
	return userId, tenantId
}
