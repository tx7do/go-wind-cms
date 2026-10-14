package service

import (
	"context"
	"strconv"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"github.com/tx7do/go-crud/viewer"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	appViewer "go-wind-cms/pkg/entgo/viewer"

	"go-wind-cms/app/core/service/internal/data"
	"go-wind-cms/pkg/task"

	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
)

// postPublishRescheduleTolerance 定时发布的改期容差：
// 旧任务到点时若 publish_time 已被推迟超过该窗口则自弃（与 social 分发同款语义）。
const postPublishRescheduleTolerance = 2 * time.Minute

type PostService struct {
	contentV1.UnimplementedPostServiceServer

	postRepo          *data.PostRepo
	searchService     *SearchService
	taskService       *TaskService
	contentModelRepo  *data.ContentModelRepo
	log               *log.Helper
}

func NewPostService(ctx *bootstrap.Context, uc *data.PostRepo, searchService *SearchService, taskService *TaskService, contentModelRepo *data.ContentModelRepo) *PostService {
	return &PostService{
		log:              ctx.NewLoggerHelper("post/service/core-service"),
		postRepo:         uc,
		searchService:    searchService,
		taskService:      taskService,
		contentModelRepo: contentModelRepo,
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

	// 自定义字段校验：按所属分类绑定的内容模型校验 custom_fields
	// TenantId 由 viewer 注入（不在此消息中），此处传 0 跳过跨租户检查
	if vErr := s.contentModelRepo.ValidateValues(ctx,
		s.contentModelRepo.ResolveModelIDByCategories(ctx, req.Data.GetCategoryIds()),
		req.Data.GetCustomFields(),
		0); vErr != nil {
		return nil, vErr
	}

	dto, err := s.postRepo.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	// 双写钩子：事务提交成功后，入队 ES 重索引。
	// best-effort：失败仅记日志，不回滚 DB；漏掉的文档由周期 ReindexAll 修复。
	s.enqueuePostReindex(ctx, dto.GetId(), "index")

	// 定时发布钩子：SCHEDULED 落库后入队到点转正任务（失败显式报错）。
	if err := s.enqueueScheduledPublish(dto); err != nil {
		return nil, err
	}

	return dto, nil
}

func (s *PostService) Update(ctx context.Context, req *contentV1.UpdatePostRequest) (*contentV1.Post, error) {
	// 自定义字段校验：仅当本次请求携带 custom_fields 时校验（部分更新不带则跳过）。
	// 分类取请求值；未携带分类时查现库归属，保证换绑模型后旧值也能被复检。
	if req != nil && req.Data != nil && req.Data.CustomFields != nil {
		categoryIds := req.Data.GetCategoryIds()
		if len(categoryIds) == 0 {
			if existing, gErr := s.postRepo.Get(ctx, &contentV1.GetPostRequest{
				QueryBy: &contentV1.GetPostRequest_Id{Id: req.GetId()},
			}); gErr == nil && existing != nil {
				categoryIds = existing.GetCategoryIds()
			}
		}
		if vErr := s.contentModelRepo.ValidateValues(ctx,
			s.contentModelRepo.ResolveModelIDByCategories(ctx, categoryIds),
			req.Data.GetCustomFields(),
			0); vErr != nil {
			return nil, vErr
		}
	}

	dto, err := s.postRepo.Update(ctx, req)
	if err != nil {
		return nil, err
	}

	// 双写钩子：更新后入队 ES 重索引（worker 会从 DB 取最新数据 upsert ES）。
	s.enqueuePostReindex(ctx, dto.GetId(), "index")

	// 定时发布钩子：以 SCHEDULED 终态保存（含改期）时重新入队到点转正任务。
	if err := s.enqueueScheduledPublish(dto); err != nil {
		return nil, err
	}

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
	// AI 语义索引（pgvector）同步入队：全局搜索可命中帖子标题/摘要。
	// best-effort 同上。
	if err := s.taskService.EnqueueAiPostIndex(&task.AiPostIndexPayload{
		PostID:   postID,
		TenantID: tenantID,
		Op:       op,
	}); err != nil {
		s.log.Errorf("enqueue ai post index failed (post_id=%d op=%s): %v", postID, op, err)
	}
}

// enqueueScheduledPublish 文章保存后的定时发布钩子：终态为 SCHEDULED 时入队
// post.publish 任务（publish_time 在未来按时刻 DELAY，为空或已过立即执行）。
//
// 每次以 SCHEDULED 状态保存都会重新入队；改期后旧任务到点自检
// publish_time 已被推迟则自弃，不会双发（详见 PublishScheduledPostWorker）。
// 入队失败返回错误：文章已落库为 SCHEDULED 而无任务执行等于永远不发布，
// 不能像 reindex 那样静默（没有周期兜底）；调用方据此向用户报错。
func (s *PostService) enqueueScheduledPublish(post *contentV1.Post) error {
	if post == nil || post.GetStatus() != contentV1.Post_POST_STATUS_SCHEDULED {
		return nil
	}

	var at time.Time
	if pt := post.GetPublishTime(); pt != nil && pt.IsValid() {
		at = pt.AsTime()
		if !at.After(time.Now()) {
			at = time.Time{} // 排期已过 → 立即转正
		}
	}

	if err := s.taskService.EnqueuePostPublish(&task.PostPublishPayload{
		PostID: post.GetId(),
	}, at); err != nil {
		s.log.Errorf("schedule post publish failed (post=%d): %v", post.GetId(), err)
		return contentV1.ErrorInternalServerError("schedule post publish failed")
	}
	return nil
}

// PublishScheduledPostWorker 是 asynq "post.publish" 任务的 handler：
// 把到点的 SCHEDULED 文章流转为 PUBLISHED，并复用既有双写钩子刷新 ES/AI 索引。
//
// 幂等与自弃规则：
//   - 状态已非 SCHEDULED（已发布/退回草稿）→ 跳过
//   - publish_time 仍晚于 now+容差（改期后的旧任务）→ 自弃，由新任务执行
//   - 转正失败返回 err 交 asynq 重试
func (s *PostService) PublishScheduledPostWorker(taskType string, payload *task.PostPublishPayload) error {
	if payload == nil || payload.PostID == 0 {
		s.log.Warnf("post publish worker: nil payload")
		return nil
	}
	s.log.Infof("post publish worker: post=%d tenant=%d", payload.PostID, payload.TenantID)

	ctx := appViewer.NewSystemViewerContext(context.Background())

	post, err := s.postRepo.Get(ctx, &contentV1.GetPostRequest{
		QueryBy: &contentV1.GetPostRequest_Id{Id: payload.PostID},
	})
	if err != nil {
		s.log.Errorf("post publish worker: post %d not loadable: %v", payload.PostID, err)
		return nil // 文章不存在/已删：不重试
	}

	if post.GetStatus() != contentV1.Post_POST_STATUS_SCHEDULED {
		s.log.Infof("post publish worker: post %d status=%s (not scheduled), skip",
			payload.PostID, post.GetStatus())
		return nil
	}
	if pt := post.GetPublishTime(); pt != nil && pt.IsValid() &&
		pt.AsTime().After(time.Now().Add(postPublishRescheduleTolerance)) {
		s.log.Infof("post publish worker: post %d rescheduled to %s, skip stale task",
			payload.PostID, pt.AsTime())
		return nil
	}

	// 转正：走 repo.Update（与手动发布同一条写入路径），再复用双写钩子刷索引。
	// publish_time 缺失时补为当前时间，避免 C 端按发布时间排序/展示拿到空值。
	updateMask := []string{"status"}
	data := &contentV1.Post{Status: trans.Ptr(contentV1.Post_POST_STATUS_PUBLISHED)}
	if post.GetPublishTime() == nil {
		data.PublishTime = timestamppb.Now()
		updateMask = append(updateMask, "publishTime")
	}

	updated, err := s.postRepo.Update(ctx, &contentV1.UpdatePostRequest{
		Id:         payload.PostID,
		Data:       data,
		UpdateMask: &fieldmaskpb.FieldMask{Paths: updateMask},
	})
	if err != nil {
		s.log.Errorf("post publish worker: publish post %d failed: %v", payload.PostID, err)
		return err // 暂时性失败：asynq 重试
	}

	// 状态变化影响 ES/AI 索引可见性（C 端只查 PUBLISHED），复用既有重索引链。
	s.enqueuePostReindex(ctx, payload.PostID, "index")

	s.log.Infof("post publish worker: post %d published (publishTime=%s)",
		payload.PostID, updated.GetPublishTime())
	return nil
}

// PublishSweepWorker 是 asynq "post.scheduled.sweep" 周期任务的 handler：
// 兜底扫描到期（publish_time 为空或已过）仍未转正的 SCHEDULED 文章并重新入队。
// DELAY 任务在 Redis 丢失/崩溃窗口内可能静默消失，本扫描是自愈网；
// 重新入队的任务仍走 PublishScheduledPostWorker 的状态自检，重复入队无害。
func (s *PostService) PublishSweepWorker(taskType string, payload *task.PublishSweepPayload) error {
	ctx := appViewer.NewSystemViewerContext(context.Background())

	ids, err := s.postRepo.ListDueScheduledPostIDs(ctx, time.Now())
	if err != nil {
		s.log.Errorf("post publish sweep: list due posts failed: %v", err)
		return err
	}
	for _, id := range ids {
		if err := s.taskService.EnqueuePostPublish(&task.PostPublishPayload{PostID: id}, time.Time{}); err != nil {
			s.log.Errorf("post publish sweep: re-enqueue post %d failed: %v", id, err)
		}
	}
	if len(ids) != 0 {
		s.log.Infof("post publish sweep: re-enqueued %d due scheduled posts", len(ids))
	}
	return nil
}
