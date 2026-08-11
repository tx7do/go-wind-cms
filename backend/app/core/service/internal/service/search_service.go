package service

import (
	"context"
	"strconv"

	"github.com/go-kratos/kratos/v2/log"
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
}

func NewSearchService(
	ctx *bootstrap.Context,
	searchRepo *data.SearchRepo,
	postRepo *data.PostRepo,
) *SearchService {
	return &SearchService{
		log:        ctx.NewLoggerHelper("search/service/core-service"),
		searchRepo: searchRepo,
		postRepo:   postRepo,
	}
}

// Search 前台全文搜索。
//
// 安全：
//   - tenantID 取自 viewer（maybeTenantFromViewer），调用方无法覆盖
//   - tid==0 → 返回空（不接受 SystemViewer bypass）
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
		// 搜索路径不 bypass：无有效租户上下文 → 返回空
		return &data.PostSearchResult{}, nil
	}

	return s.searchRepo.SearchPosts(ctx, query, tenantID, language, status, page, pageSize)
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
