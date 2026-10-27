package service

import (
	"context"
	"strconv"

	"github.com/go-kratos/kratos/v2/log"
	kratosMetadata "github.com/go-kratos/kratos/v2/metadata"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"github.com/tx7do/go-crud/viewer"

	"go-wind-cms/app/core/service/internal/data"
	appViewer "go-wind-cms/pkg/entgo/viewer"
	"go-wind-cms/pkg/task"
)

// ============================================================================
// SearchService —— OpenSearch 搜索与重索引的业务编排层
//
// 职责：
//   - Search：前台全文搜索入口，强制租户/语言/状态过滤，不接受 bypass
//   - ReindexPost：asynq worker handler，从 DB 取数据写入/删除 ES
//   - ReindexAll：周期全量重索引，修复崩溃窗口内的漂移
//
// 依赖：
//   - searchRepo：ES 操作（强制隔离）
//   - postRepo：reindex 时读 DB（含 tenant_id 真实值）
//
// 安全模型详见 data/search_repo.go 顶部注释。核心差异：
//   - 搜索路径：UserViewer，tid==0 返回空，绝不 bypass
//   - reindex 路径：注入 SystemViewer 跨租户读 DB，但 ES 文档 tenant_id 取自 DB
// ============================================================================

type SearchService struct {
	log        *log.Helper
	searchRepo *data.SearchRepo
	postRepo   *data.PostRepo

	synonymRepo *data.SearchSynonymRepo
}

func NewSearchService(
	ctx *bootstrap.Context,
	searchRepo *data.SearchRepo,
	postRepo *data.PostRepo,
	synonymRepo *data.SearchSynonymRepo,) *SearchService {
	return &SearchService{
		log:        ctx.NewLoggerHelper("search/service/core-service"),
		searchRepo: searchRepo,
		postRepo:   postRepo,
	
		synonymRepo: synonymRepo,
}
}

// Search 前台全文搜索。
//
// 安全：
//   - tenantID 优先取自 viewer（maybeTenantFromViewerForSearch），调用方无法覆盖；
//     viewer 无租户（如匿名 gRPC 链路，core 端拿不到 HTTP Host）时，回退读
//     x-md-global-tenant-id 元数据——该值由 app BFF 按请求 Host 解析后注入，
//     core 的 gRPC 端口仅在内网经注册中心暴露，外部调用方无法直接伪造该 header
//   - 两者皆无 → 返回空（不接受 SystemViewer bypass）
//   - 语言/状态由调用方传，但 SearchRepo 内部强制 term 过滤，不可绕过
func (s *SearchService) Search(
	ctx context.Context,
	query string,
	language string,
	status string,
	page int,
	pageSize int,
) (*data.PostSearchResult, error) {
	tenantID, hasTenant := maybeTenantFromViewerForSearch(ctx)
	if !hasTenant {
		tenantID = tenantIDFromGlobalMetadata(ctx)
		hasTenant = tenantID > 0
	}
	if !hasTenant {
		// 搜索路径不 bypass：无有效租户上下文 → 返回空
		return &data.PostSearchResult{}, nil
	}

	// 同义词查询侧扩展（租户词库；扩展词以 OR 组参与检索，原始查询不丢）。
	// 词库加载失败/无命中都不影响主查询——扩展是纯增益。
	if s.synonymRepo != nil {
		if expanded := s.synonymRepo.ExpandQuery(ctx, tenantID, query); len(expanded) > 0 {
			return s.searchRepo.SearchPostsExpanded(ctx, query, expanded, tenantID, language, status, page, pageSize)
		}
	}
	return s.searchRepo.SearchPosts(ctx, query, tenantID, language, status, page, pageSize)
}

// SearchPostsSuggest 标题前缀联想（租户解析同 Search；免鉴权公开链路）。
func (s *SearchService) SearchPostsSuggest(ctx context.Context, prefix string, language string, limit int) ([]string, error) {
	tenantID, hasTenant := maybeTenantFromViewerForSearch(ctx)
	if !hasTenant {
		tenantID = tenantIDFromGlobalMetadata(ctx)
		hasTenant = tenantID > 0
	}
	if !hasTenant {
		return []string{}, nil
	}
	return s.searchRepo.SuggestPostTitles(ctx, prefix, tenantID, language, limit)
}

// tenantIDFromGlobalMetadata 读取 BFF 注入的租户元数据（kratos metadata 中间件
// 透传的 x-md-global-* 键），供匿名搜索链路在 gRPC 侧恢复租户上下文。
func tenantIDFromGlobalMetadata(ctx context.Context) uint32 {
	md, ok := kratosMetadata.FromServerContext(ctx)
	if !ok {
		return 0
	}
	raw := md.Get("x-md-global-tenant-id")
	if raw == "" {
		return 0
	}
	tid, err := strconv.ParseUint(raw, 10, 32)
	if err != nil {
		return 0
	}
	return uint32(tid)
}

// ReindexPost 是 asynq "search.reindex" 任务的 worker handler。
//
// 签名遵循 (taskType string, payload *T) error 模式（参考 TaskService.AsyncBackup）。
//
// 安全：
//   - 注入 SystemViewer 跨租户读 DB（特权）
//   - ES 文档 tenant_id 取自 DB 记录（PostRepo.GetReindexDocuments 内部从 ent 取）
//   - 非убликовated / tenant_id==0 / 空翻译 → 跳过或删除
func (s *SearchService) ReindexPost(taskType string, payload *task.SearchReindexPayload) error {
	if payload == nil {
		s.log.Warnf("reindex post: nil payload")
		return nil
	}

	// 注意：payload.TenantID 仅用于此日志行，ES 文档的 tenant_id 取自 DB 记录
	// （PostRepo.GetReindexDocuments 内部从 ent 取），不可用 payload 值覆盖。
	// Deprecated: payload.TenantID 仅供日志展示，后续应从 payload 移除以避免误用。
	s.log.Infof("reindex post: entity=%s id=%d tenant=%d op=%s",
		payload.Entity, payload.ID, payload.TenantID, payload.Op)

	if payload.Entity != "post" {
		s.log.Warnf("reindex post: unsupported entity %s", payload.Entity)
		return nil
	}

	// 注入 SystemViewer：跨租户读 DB 的特权上下文
	ctx := appViewer.NewSystemViewerContext(context.Background())

	switch payload.Op {
	case "delete":
		// post 被删除（软删/硬删/状态变更），按 post_id 删 ES 所有语言文档
		if err := s.searchRepo.DeletePost(ctx, payload.ID); err != nil {
			s.log.Errorf("reindex delete post %d failed: %v", payload.ID, err)
			return err
		}
		return nil

	case "index":
		// post 创建/更新，从 DB 取最新数据写入 ES
		docs, err := s.postRepo.GetReindexDocuments(ctx, payload.ID)
		if err != nil {
			s.log.Errorf("reindex get documents for post %d failed: %v", payload.ID, err)
			return err
		}

		// docs 为 nil 表示 post 不存在/非убликовated/tenant_id==0
		// 此时若 ES 中有残留文档（如状态从ublished变draft），需删除
		if len(docs) == 0 {
			if err := s.searchRepo.DeletePost(ctx, payload.ID); err != nil {
				s.log.Errorf("reindex cleanup post %d failed: %v", payload.ID, err)
				return err
			}
			return nil
		}

		// 确保 posts 索引模板存在（幂等）
		if err := s.searchRepo.EnsureIndexTemplate(ctx); err != nil {
			s.log.Errorf("ensure index template failed: %v", err)
			return err
		}

		// 逐条 upsert ES 文档（每个语言一个文档）
		for i := range docs {
			doc := &data.PostDocument{
				TenantID: strconv.FormatUint(uint64(docs[i].TenantID), 10),
				PostID:   strconv.FormatUint(uint64(docs[i].PostID), 10),
				Language: docs[i].Language,
				Status:   docs[i].Status,
				Title:    docs[i].Title,
				Summary:  docs[i].Summary,
				Content:  docs[i].Content,
			}
			if err := s.searchRepo.IndexPost(ctx, doc); err != nil {
				s.log.Errorf("index post document %d lang=%s failed: %v",
					docs[i].PostID, docs[i].Language, err)
				return err
			}
		}
		return nil

	default:
		s.log.Warnf("reindex post: unknown op %s", payload.Op)
		return nil
	}
}

// ReindexAll 周期全量重索引，修复崩溃窗口内的漂移。
// 遍历所有 PUBLISHED post，逐条 GetReindexDocuments + IndexPost。
func (s *SearchService) ReindexAll() error {
	ctx := appViewer.NewSystemViewerContext(context.Background())

	if err := s.searchRepo.EnsureIndexTemplate(ctx); err != nil {
		s.log.Errorf("ensure index template failed: %v", err)
		return err
	}

	postIDs, err := s.postRepo.ListPublishedPostIDs(ctx)
	if err != nil {
		s.log.Errorf("list published post ids failed: %v", err)
		return err
	}

	s.log.Infof("reindex all: %d published posts to process", len(postIDs))

	var failed int
	for _, postID := range postIDs {
		docs, err := s.postRepo.GetReindexDocuments(ctx, postID)
		if err != nil {
			s.log.Errorf("reindex all: get documents for post %d failed: %v", postID, err)
			failed++
			continue
		}
		for i := range docs {
			doc := &data.PostDocument{
				TenantID: strconv.FormatUint(uint64(docs[i].TenantID), 10),
				PostID:   strconv.FormatUint(uint64(docs[i].PostID), 10),
				Language: docs[i].Language,
				Status:   docs[i].Status,
				Title:    docs[i].Title,
				Summary:  docs[i].Summary,
				Content:  docs[i].Content,
			}
			if err := s.searchRepo.IndexPost(ctx, doc); err != nil {
				s.log.Errorf("reindex all: index post %d lang=%s failed: %v",
					docs[i].PostID, docs[i].Language, err)
				failed++
			}
		}
	}

	if failed > 0 {
		s.log.Warnf("reindex all completed with %d/%d failures", failed, len(postIDs))
	} else {
		s.log.Infof("reindex all completed: %d posts indexed", len(postIDs))
	}
	return nil
}

// maybeTenantFromViewerForSearch 搜索专用租户提取。
// 与 internal_message_service.go 的 senderTenantID 取法一致，直接从 viewer
// context 取 tenant ID。tid==0 → hasTenant=false，触发搜索返回空（不 bypass）。
func maybeTenantFromViewerForSearch(ctx context.Context) (tenantID uint32, hasTenant bool) {
	vc, exist := viewer.FromContext(ctx)
	if !exist || vc == nil {
		return 0, false
	}
	tid := uint32(vc.TenantID())
	if tid == 0 {
		return 0, false
	}
	return tid, true
}
