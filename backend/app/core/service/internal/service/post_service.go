package service

import (
	"context"
	"strconv"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"github.com/tx7do/go-crud/viewer"
	"google.golang.org/protobuf/types/known/emptypb"

	"go-wind-cms/app/core/service/internal/data"
	"go-wind-cms/pkg/task"

	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
)

type PostService struct {
	contentV1.UnimplementedPostServiceServer

	postRepo      *data.PostRepo
	searchService *SearchService
	taskService   *TaskService
	log           *log.Helper
}

func NewPostService(ctx *bootstrap.Context, uc *data.PostRepo, searchService *SearchService, taskService *TaskService) *PostService {
	return &PostService{
		log:           ctx.NewLoggerHelper("post/service/core-service"),
		postRepo:      uc,
		searchService: searchService,
		taskService:   taskService,
	}
}

func (s *PostService) List(ctx context.Context, req *paginationV1.PagingRequest) (*contentV1.ListPostResponse, error) {
	return s.postRepo.List(ctx, req)
}

// SearchPosts 前台全文搜索。
//
// 安全：
//   - tenant_id 由 SearchService.Search 从 viewer 注入，客户端无法指定
//   - status 硬编码为 PUBLISHED——前台仅检索已发布内容，不接受客户端传 status
//   - language 由客户端传，但 SearchRepo 内部强制 term 过滤
//   - 响应只含 post_id / language / title，不含 content / tenant_id / status
func (s *PostService) SearchPosts(ctx context.Context, req *contentV1.SearchPostsRequest) (*contentV1.SearchPostsResponse, error) {
	if req == nil {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}

	// status 固定 PUBLISHED，不接受客户端覆盖
	const publishedStatus = "POST_STATUS_PUBLISHED"

	result, err := s.searchService.Search(
		ctx,
		req.GetQuery(),
		req.GetLanguage(),
		publishedStatus,
		int(req.GetPage()),
		int(req.GetPageSize()),
	)
	if err != nil {
		s.log.Errorf("search posts failed: %v", err)
		return nil, contentV1.ErrorInternalServerError("search posts failed")
	}

	// 转换 data.PostSearchResult → contentV1.SearchPostsResponse
	// ES 文档中 post_id 存为 string（keyword），转回 uint32 用于 proto 响应
	resp := &contentV1.SearchPostsResponse{
		Total: int32(result.Total),
		Items: make([]*contentV1.SearchPostHit, 0, len(result.Hits)),
	}
	for _, hit := range result.Hits {
		pid, err := strconv.ParseUint(hit.PostID, 10, 32)
		if err != nil {
			s.log.Warnf("search result: invalid post_id %q, skipping", hit.PostID)
			continue
		}
		resp.Items = append(resp.Items, &contentV1.SearchPostHit{
			PostId:   uint32(pid),
			Language: hit.Language,
			Title:    hit.Title,
		})
	}
	return resp, nil
}

func (s *PostService) Get(ctx context.Context, req *contentV1.GetPostRequest) (*contentV1.Post, error) {
	return s.postRepo.Get(ctx, req)
}

func (s *PostService) Create(ctx context.Context, req *contentV1.CreatePostRequest) (*contentV1.Post, error) {
	if req == nil || req.Data == nil {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}
	if req.Data.Status == nil {
		// 新建文章默认为草稿状态，发布须显式 Update 切换状态，
		// 避免创建即发布绕过编辑审核流程
		req.Data.Status = trans.Ptr(contentV1.Post_POST_STATUS_DRAFT)
	}

	dto, err := s.postRepo.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	// 双写钩子：事务提交成功后，入队 ES 重索引。
	// best-effort：失败仅记日志，不回滚 DB；漏掉的文档由周期 ReindexAll 修复。
	s.enqueuePostReindex(ctx, dto.GetId(), "index")

	return dto, nil
}

func (s *PostService) Update(ctx context.Context, req *contentV1.UpdatePostRequest) (*contentV1.Post, error) {
	dto, err := s.postRepo.Update(ctx, req)
	if err != nil {
		return nil, err
	}

	// 双写钩子：更新后入队 ES 重索引（worker 会从 DB 取最新数据 upsert ES）。
	s.enqueuePostReindex(ctx, dto.GetId(), "index")

	return dto, nil
}

func (s *PostService) Delete(ctx context.Context, req *contentV1.DeletePostRequest) (*emptypb.Empty, error) {
	err := s.postRepo.Delete(ctx, req)
	if err != nil {
		return nil, err
	}

	// 双写钩子：删除后入队 ES 删除（worker 按 post_id 删 ES 所有语言文档）。
	s.enqueuePostReindex(ctx, req.GetId(), "delete")

	return &emptypb.Empty{}, nil
}

func (s *PostService) TranslationExists(ctx context.Context, req *contentV1.PostTranslationExistsRequest) (*contentV1.PostTranslationExistsResponse, error) {
	exists, err := s.postRepo.TranslationExists(ctx, req.GetPostId(), req.GetLanguageCode())
	if err != nil {
		return nil, err
	}

	return &contentV1.PostTranslationExistsResponse{
		Exists: exists,
	}, nil
}

func (s *PostService) GetTranslation(ctx context.Context, req *contentV1.GetPostRequest) (*contentV1.PostTranslation, error) {
	return s.postRepo.GetTranslation(ctx, req)
}

func (s *PostService) CreateTranslation(ctx context.Context, req *contentV1.CreatePostTranslationRequest) (*contentV1.PostTranslation, error) {
	dto, err := s.postRepo.CreateTranslation(ctx, req)
	if err != nil {
		return nil, err
	}

	// 翻译内容变更影响 ES 文档（title/summary/content），入队重索引。
	if dto != nil {
		s.enqueuePostReindex(ctx, dto.GetPostId(), "index")
	}

	return dto, nil
}

func (s *PostService) UpdateTranslation(ctx context.Context, req *contentV1.UpdatePostTranslationRequest) (*contentV1.PostTranslation, error) {
	dto, err := s.postRepo.UpdateTranslation(ctx, req)
	if err != nil {
		return nil, err
	}

	if dto != nil {
		s.enqueuePostReindex(ctx, dto.GetPostId(), "index")
	}

	return dto, nil
}

func (s *PostService) DeleteTranslation(ctx context.Context, req *contentV1.DeletePostTranslationRequest) (*emptypb.Empty, error) {
	err := s.postRepo.DeleteTranslation(ctx, req)
	if err != nil {
		return nil, err
	}

	// 翻译被删除后，ES 中对应语言文档应移除。入队重索引：worker 会重取该 post
	// 的所有翻译，被删的翻译不再返回 → docs 减少 → worker 删除 ES 中残留文档。
	//
	// 仅 Identifier 模式可拿到 postID 入队；Id 模式（按翻译自身 id 删除）拿不到
	// postID，跳过即时同步，依赖周期 ReindexAll 兜底（最长 1 小时滞后）。
	if req != nil {
		if identifier := req.GetIdentifier(); identifier != nil {
			s.enqueuePostReindex(ctx, identifier.GetPostId(), "index")
		}
	}

	return &emptypb.Empty{}, nil
}

// enqueuePostReindex 是双写钩子的统一入口。
//
// 从 viewer context 取 tenant_id（仅用于 payload 日志辅助），构造
// SearchReindexPayload 并调 TaskService.EnqueueSearchReindex 入队。
//
// 安全：
//   - payload 的 TenantID 仅日志用，ES 文档 tenant_id 由 worker 从 DB 取
//   - op=index 时 worker 重取该 post 所有翻译 upsert ES
//   - op=delete 时 worker 按 post_id 删 ES 所有语言文档
//   - best-effort：入队失败仅记日志，不阻断主业务
func (s *PostService) enqueuePostReindex(ctx context.Context, postID uint32, op string) {
	if postID == 0 {
		return
	}

	// tenant_id 仅用于日志，取自 viewer（与 internal_message_service 同模式）
	var tenantID uint32
	if vc, exist := viewer.FromContext(ctx); exist && vc != nil {
		tenantID = uint32(vc.TenantID())
	}

	payload := &task.SearchReindexPayload{
		Entity:   "post",
		ID:       postID,
		TenantID: tenantID,
		Op:       op,
	}

	if err := s.taskService.EnqueueSearchReindex(payload); err != nil {
		s.log.Errorf("enqueue post reindex failed (post_id=%d op=%s): %v", postID, op, err)
	}
}
