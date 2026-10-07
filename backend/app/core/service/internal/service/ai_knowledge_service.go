package service

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/sashabaranov/go-openai"
	"github.com/tx7do/go-crud/viewer"
	"github.com/tx7do/go-utils/doctext"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"

	appViewer "go-wind-cms/pkg/entgo/viewer"
	"go-wind-cms/pkg/task"

	aiV1 "go-wind-cms/api/gen/go/ai/service/v1"
	"go-wind-cms/app/core/service/internal/data"
	"go-wind-cms/app/core/service/internal/data/ent"
)

const (
	// ragChunkRunes 切片长度（字符）：中文友好。
	ragChunkRunes = 500
	// ragChunkOverlap 切片重叠（字符）：保证跨片语义连续。
	ragChunkOverlap = 50
	// ragTopK 检索注入上下文的默认片段数。
	ragTopK = 3
)

// AiKnowledgeService 知识库管理：文档入库（切片→向量化）+ 检索 + Chat 注入。
// 操作者身份经 viewer 上下文（unary 中间件注入）获取。
type AiKnowledgeService struct {
	aiV1.UnimplementedAiKnowledgeBaseServiceServer
	log  *log.Helper
	repo *data.AiKnowledgeRepo

	providerRepo *data.AiProviderRepo
	usageLogRepo *data.AiUsageLogRepo
}

func NewAiKnowledgeService(
	ctx *bootstrap.Context,
	repo *data.AiKnowledgeRepo,
	providerRepo *data.AiProviderRepo,
	usageLogRepo *data.AiUsageLogRepo,
) *AiKnowledgeService {
	s := &AiKnowledgeService{
		log:          ctx.NewLoggerHelper("ai_knowledge/service/core-service"),
		repo:         repo,
		providerRepo: providerRepo,
		usageLogRepo: usageLogRepo,
	}

	// 启动期补建 pgvector 扩展与 embedding 列（幂等；失败只降级 RAG，不阻断服务）
	if err := repo.MigrateVectorColumn(ctx.Context()); err != nil {
		s.log.Errorf("rag vector column migration failed: %v", err)
	}

	return s
}

// ── 知识库 CRUD ────────────────────────────────────────────────────

func (s *AiKnowledgeService) List(ctx context.Context, req *paginationV1.PagingRequest) (*aiV1.ListAiKnowledgeBaseResponse, error) {
	items, total, err := s.repo.List(ctx, req)
	if err != nil {
		return nil, err
	}
	s.repo.FillDocCounts(ctx, items)
	return &aiV1.ListAiKnowledgeBaseResponse{Items: items, Total: total}, nil
}

func (s *AiKnowledgeService) Get(ctx context.Context, req *aiV1.GetAiKnowledgeBaseRequest) (*aiV1.AiKnowledgeBase, error) {
	return s.repo.Get(ctx, req)
}

func (s *AiKnowledgeService) Create(ctx context.Context, req *aiV1.CreateAiKnowledgeBaseRequest) (*emptypb.Empty, error) {
	if req.Data == nil {
		return nil, aiV1.ErrorBadRequest("invalid parameter")
	}
	userId, tenantId := s.operator(ctx)

	req.Data.CreatedBy = trans.Ptr(userId)
	req.Data.UserId = trans.Ptr(userId)
	req.Data.TenantId = trans.Ptr(tenantId)

	// provider 必须存在且启用（向量化要走它的端点）
	if req.Data.ProviderId == nil || req.Data.EmbeddingModel == nil || *req.Data.EmbeddingModel == "" {
		return nil, aiV1.ErrorBadRequest("provider_id and embedding_model are required")
	}
	if _, err := s.providerRepo.GetEntityByID(ctx, req.Data.GetProviderId()); err != nil {
		return nil, err
	}

	if _, err := s.repo.Create(ctx, req.Data); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *AiKnowledgeService) Update(ctx context.Context, req *aiV1.UpdateAiKnowledgeBaseRequest) (*emptypb.Empty, error) {
	if req.Data == nil {
		return nil, aiV1.ErrorBadRequest("invalid parameter")
	}
	userId, _ := s.operator(ctx)
	req.Data.Id = trans.Ptr(req.GetId())
	req.Data.UpdatedBy = trans.Ptr(userId)
	if req.UpdateMask != nil {
		req.UpdateMask.Paths = append(req.UpdateMask.Paths, "updated_by")
	}
	if err := s.repo.Update(ctx, req); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *AiKnowledgeService) Delete(ctx context.Context, req *aiV1.DeleteAiKnowledgeBaseRequest) (*emptypb.Empty, error) {
	if err := s.repo.Delete(ctx, req.GetId()); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// ── 文档入库：切片 → embedding → 落库 ───────────────────────────────

// UploadDoc 纯文本入库入口。
func (s *AiKnowledgeService) UploadDoc(ctx context.Context, req *aiV1.UploadAiDocRequest) (*aiV1.UploadAiDocResponse, error) {
	userId, tenantId := s.operator(ctx)
	if req.GetBaseId() == 0 || req.GetName() == "" || req.GetContent() == "" {
		return nil, aiV1.ErrorBadRequest("base_id, name and content are required")
	}

	base, err := s.repo.GetEntityByID(ctx, req.GetBaseId())
	if err != nil {
		return nil, err
	}

	doc, err := s.ingestDoc(ctx, userId, tenantId, base, req.GetName(), req.GetContent())
	if err != nil {
		return nil, err
	}
	return &aiV1.UploadAiDocResponse{Doc: doc, ChunkCount: doc.GetChunkCount()}, nil
}

// UploadDocFile 文件入库入口：按扩展名抽取文本（txt/md/docx/pdf）后走同一入库链。
func (s *AiKnowledgeService) UploadDocFile(ctx context.Context, req *aiV1.UploadAiDocFileRequest) (*aiV1.UploadAiDocResponse, error) {
	userId, tenantId := s.operator(ctx)
	if req.GetBaseId() == 0 || req.GetFileName() == "" || len(req.GetContentBase64()) == 0 {
		return nil, aiV1.ErrorBadRequest("base_id, fileName and contentBase64 are required")
	}

	base, err := s.repo.GetEntityByID(ctx, req.GetBaseId())
	if err != nil {
		return nil, err
	}

	content, err := doctext.Extract(req.GetFileName(), req.GetContentBase64())
	if err != nil {
		return nil, aiV1.ErrorBadRequest("extract text failed: %v", err)
	}

	docName := req.GetDocName()
	if docName == "" {
		docName = strings.TrimSuffix(req.GetFileName(), filepath.Ext(req.GetFileName()))
	}

	doc, err := s.ingestDoc(ctx, userId, tenantId, base, docName, content)
	if err != nil {
		return nil, err
	}
	return &aiV1.UploadAiDocResponse{Doc: doc, ChunkCount: doc.GetChunkCount()}, nil
}

// ingestDoc 切片 → 向量化 → 落库的共用主链（纯文本与文件入口共用）。
// 失败时文档行标 FAILED（错误信息保留），不让一次 embedding 故障污染整个请求语义。
func (s *AiKnowledgeService) ingestDoc(ctx context.Context, userId, tenantId uint32, base *ent.AiKnowledgeBase, docName, content string) (*aiV1.AiDoc, error) {
	if base.EmbeddingModel == nil || base.ProviderID == nil {
		return nil, aiV1.ErrorBadRequest("knowledge base embedding config is incomplete")
	}

	// 1. 切片
	chunks := splitIntoChunks(content, ragChunkRunes, ragChunkOverlap)

	// 2. 建文档行（先 READY 假设成功，失败置 FAILED）
	doc, err := s.repo.CreateDoc(ctx, &aiV1.AiDoc{
		BaseId:   &base.ID,
		Name:     trans.Ptr(docName),
		Status:   trans.Ptr("READY"),
		UserId:   &userId,
		TenantId: trans.Ptr(tenantId),
	})
	if err != nil {
		return nil, err
	}

	// 3. 向量化（一次请求全部切片，OpenAI 兼容 /v1/embeddings；消耗记入用量流水）
	vectors, err := embedTextsForBase(ctx, s.providerRepo, s.usageLogRepo, s.log, userId, tenantId, base, chunks)
	if err != nil {
		errMsg := err.Error()
		_ = s.repo.UpdateDocStatus(ctx, doc.ID, "FAILED", uint32(0), errMsg)
		return nil, aiV1.ErrorInternalServerError("embed doc failed: %v", err)
	}

	// 4. 切片落库
	if err = s.repo.InsertChunks(ctx, tenantId, doc.ID, chunks, vectors); err != nil {
		errMsg := err.Error()
		_ = s.repo.UpdateDocStatus(ctx, doc.ID, "FAILED", uint32(0), errMsg)
		return nil, err
	}

	chunkCount := uint32(len(chunks))
	if err = s.repo.UpdateDocStatus(ctx, doc.ID, "READY", chunkCount, ""); err != nil {
		return nil, err
	}

	doc, err = s.repo.GetDocByID(ctx, doc.ID)
	if err != nil {
		return nil, err
	}
	return s.repo.ToDocDTO(doc), nil
}

func (s *AiKnowledgeService) ListDocs(ctx context.Context, req *aiV1.ListAiDocsRequest) (*aiV1.ListAiDocsResponse, error) {
	if req.GetBaseId() == 0 {
		return nil, aiV1.ErrorBadRequest("base_id is required")
	}
	entities, err := s.repo.ListDocs(ctx, req.GetBaseId())
	if err != nil {
		return nil, err
	}
	dtos := make([]*aiV1.AiDoc, 0, len(entities))
	for _, entity := range entities {
		dtos = append(dtos, s.repo.ToDocDTO(entity))
	}
	return &aiV1.ListAiDocsResponse{Items: dtos, Total: uint64(len(dtos))}, nil
}

func (s *AiKnowledgeService) DeleteDoc(ctx context.Context, req *aiV1.DeleteAiDocRequest) (*emptypb.Empty, error) {
	if err := s.repo.DeleteDoc(ctx, req.GetId()); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// ── 检索 ───────────────────────────────────────────────────────────

func (s *AiKnowledgeService) Search(ctx context.Context, req *aiV1.SearchAiKnowledgeRequest) (*aiV1.SearchAiKnowledgeResponse, error) {
	userId, tenantId := s.operator(ctx)
	if req.GetBaseId() == 0 || req.GetQuery() == "" {
		return nil, aiV1.ErrorBadRequest("base_id and query are required")
	}

	topK := int(req.GetTopK())
	hits, err := searchKnowledgeBase(ctx, s.repo, s.providerRepo, s.usageLogRepo, s.log, userId, tenantId, req.GetBaseId(), req.GetQuery(), topK)
	if err != nil {
		return nil, err
	}

	resp := &aiV1.SearchAiKnowledgeResponse{Hits: make([]*aiV1.KnowledgeHit, 0, len(hits))}
	for _, hit := range hits {
		resp.Hits = append(resp.Hits, &aiV1.KnowledgeHit{
			Doc:        &aiV1.AiDoc{Id: &hit.DocID},
			ChunkIndex: hit.ChunkIndex,
			Content:    hit.Content,
			Score:      hit.Score,
		})
	}
	return resp, nil
}

// searchKnowledgeBase 包级检索：向量化 query 并按余弦相似度取 topK 片段。
// tenantId=0（平台用户）可检索任意库，租户用户经 SQL 的 tenant 过滤兜底；
// chat 主链路（RAG 注入）与本 service 共用。向量化消耗记入用量流水。
func searchKnowledgeBase(ctx context.Context, knowledgeRepo *data.AiKnowledgeRepo, providerRepo *data.AiProviderRepo, usageLogRepo *data.AiUsageLogRepo, log *log.Helper, userId, tenantId uint32, baseId uint32, query string, topK int) ([]*data.ChunkHit, error) {
	base, err := knowledgeRepo.GetEntityByID(ctx, baseId)
	if err != nil {
		return nil, err
	}

	vecLiteral, err := embedOneForBase(ctx, providerRepo, usageLogRepo, log, userId, tenantId, base, query)
	if err != nil {
		return nil, aiV1.ErrorInternalServerError("embed query failed: %v", err)
	}

	return knowledgeRepo.SearchChunks(ctx, tenantId, baseId, vecLiteral, topK)
}

// ── 向量化 ─────────────────────────────────────────────────────────

// embedOneForBase 单文本向量化，返回 pgvector 字面量（包级，供检索与 chat 复用）。
func embedOneForBase(ctx context.Context, providerRepo *data.AiProviderRepo, usageLogRepo *data.AiUsageLogRepo, log *log.Helper, userId, tenantId uint32, base *ent.AiKnowledgeBase, text string) (string, error) {
	vectors, err := embedTextsForBase(ctx, providerRepo, usageLogRepo, log, userId, tenantId, base, []string{text})
	if err != nil {
		return "", err
	}
	return vectors[0], nil
}

// embedTextsForBase 批量向量化，返回与 texts 对齐的 pgvector 字面量数组。
// embedding 消耗与 chat 同口径记入用量流水；记量失败不阻断向量化，仅记日志。
func embedTextsForBase(ctx context.Context, providerRepo *data.AiProviderRepo, usageLogRepo *data.AiUsageLogRepo, log *log.Helper, userId, tenantId uint32, base *ent.AiKnowledgeBase, texts []string) ([]string, error) {
	provider, err := providerRepo.GetEntityByID(ctx, derefUint32(base.ProviderID))
	if err != nil {
		return nil, err
	}

	client, err := newOpenAIClientForProvider(provider)
	if err != nil {
		return nil, err
	}

	callStart := time.Now()
	resp, err := client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Model: openai.EmbeddingModel(ptrStrOr(base.EmbeddingModel, "text-embedding-3-small")),
		Input: texts,
	})
	if err != nil {
		return nil, err
	}
	if len(resp.Data) != len(texts) {
		return nil, fmt.Errorf("embedding count mismatch: got %d want %d", len(resp.Data), len(texts))
	}

	if usageLogRepo != nil {
		modelName := ptrStrOr(base.EmbeddingModel, "text-embedding-3-small")
		if uerr := usageLogRepo.Create(ctx, &data.UsageRecord{
			UserID:       userId,
			TenantID:     tenantId,
			ModelName:    modelName,
			PromptTokens: uint32(resp.Usage.PromptTokens),
			TotalTokens:  uint32(resp.Usage.TotalTokens),
			DurationMs:   uint32(time.Since(callStart).Milliseconds()),
		}); uerr != nil {
			log.Errorf("record embedding usage failed: %v", uerr)
		}
	}

	literals := make([]string, 0, len(resp.Data))
	for _, item := range resp.Data {
		literals = append(literals, vectorToLiteral(item.Embedding))
	}
	return literals, nil
}

// vectorToLiteral float32 向量 → pgvector 字面量 '[0.1,0.2,...]'。
func vectorToLiteral(vec []float32) string {
	var sb strings.Builder
	sb.WriteByte('[')
	for i, v := range vec {
		if i > 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(fmt.Sprintf("%g", v))
	}
	sb.WriteByte(']')
	return sb.String()
}

// splitIntoChunks 按字符切片，带重叠保证跨片语义连续。
func splitIntoChunks(content string, size, overlap int) []string {
	runes := []rune(content)
	if len(runes) == 0 {
		return nil
	}
	step := size - overlap
	if step <= 0 {
		step = size
	}
	chunks := make([]string, 0, len(runes)/step+1)
	for start := 0; start < len(runes); start += step {
		end := start + size
		if end > len(runes) {
			end = len(runes)
		}
		chunks = append(chunks, string(runes[start:end]))
		if end == len(runes) {
			break
		}
	}
	return chunks
}

// derefUint32 nil 安全取值。
func derefUint32(p *uint32) uint32 {
	if p == nil {
		return 0
	}
	return *p
}

// operator 从 viewer 上下文取操作者（unary 中间件注入）。
func (s *AiKnowledgeService) operator(ctx context.Context) (uint32, uint32) {
	var userId, tenantId uint32
	if vc, exist := viewer.FromContext(ctx); exist && vc != nil {
		tenantId = uint32(vc.TenantID())
	}
	if uid, ok := viewerUserIDFromContext(ctx); ok {
		userId = uid
	}
	return userId, tenantId
}

// ── 重索引任务（asynq：ai.doc.reindex） ─────────────────────────────

// reindexEmbedBatchSize 每批送 /v1/embeddings 的切片数。
const reindexEmbedBatchSize = 32

// ReindexDocs 重索引任务 worker：对指定（或全部）知识库的既有切片用当前
// embedding 模型重算向量。切片文本不变，只换 embedding 列——场景是管理员
// 更换了知识库的 embedding 模型（或 provider 端点）后需要重建。
// asynq ctx 不携带 viewer，租户隔离会拒绝无 viewer 查询；重索引是平台操作，
// 用系统查看器。
func (s *AiKnowledgeService) ReindexDocs(taskType string, data *task.AiDocReindexPayload) error {
	ctx := appViewer.NewSystemViewerContext(context.Background())

	var baseIds []uint32
	if data != nil && data.BaseID > 0 {
		baseIds = append(baseIds, data.BaseID)
	} else {
		bases, err := s.repo.ListBases(ctx)
		if err != nil {
			return err
		}
		for _, b := range bases {
			baseIds = append(baseIds, b.ID)
		}
	}

	for _, baseId := range baseIds {
		updated, err := s.reindexBase(ctx, baseId)
		if err != nil {
			s.log.Errorf("reindex base %d failed: %v", baseId, err)
			continue // 单库失败不阻断其余库
		}
		s.log.Infof("reindex base %d done: %d chunks re-embedded", baseId, updated)
	}
	return nil
}

// reindexBase 对单个知识库分批重算全部切片向量，返回更新的切片数。
func (s *AiKnowledgeService) reindexBase(ctx context.Context, baseId uint32) (int, error) {
	base, err := s.repo.GetEntityByID(ctx, baseId)
	if err != nil {
		return 0, err
	}

	refs, err := s.repo.ListChunkRefsByBase(ctx, baseId)
	if err != nil {
		return 0, err
	}
	if len(refs) == 0 {
		return 0, nil
	}

	updated := 0
	for start := 0; start < len(refs); start += reindexEmbedBatchSize {
		end := start + reindexEmbedBatchSize
		if end > len(refs) {
			end = len(refs)
		}
		batch := refs[start:end]

		texts := make([]string, 0, len(batch))
		for _, ref := range batch {
			texts = append(texts, ref.Content)
		}
		// 重索引是平台维护操作（无操作者）：消耗按知识库归属租户计量、userId=0。
		vectors, err := embedTextsForBase(ctx, s.providerRepo, s.usageLogRepo, s.log, 0, derefUint32(base.TenantID), base, texts)
		if err != nil {
			return updated, err
		}
		for i, ref := range batch {
			if err = s.repo.UpdateChunkEmbedding(ctx, ref.ID, vectors[i]); err != nil {
				return updated, err
			}
			updated++
		}
	}
	return updated, nil
}
