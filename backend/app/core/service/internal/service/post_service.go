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
	"google.golang.org/protobuf/encoding/protojson"
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

	postRepo              *data.PostRepo
	postVersionRepo       *data.PostVersionRepo
	searchService         *SearchService
	taskService           *TaskService
	webhookService        *WebhookService
	internalMessageSvc    *InternalMessageService
	contentModelRepo      *data.ContentModelRepo
	referenceTracker      *ReferenceTrackingService
	flowSvc               *FlowService
	log                   *log.Helper
}

func NewPostService(ctx *bootstrap.Context, uc *data.PostRepo, searchService *SearchService, taskService *TaskService, contentModelRepo *data.ContentModelRepo, postVersionRepo *data.PostVersionRepo, webhookService *WebhookService, internalMessageSvc *InternalMessageService, referenceTracker *ReferenceTrackingService, flowSvc *FlowService) *PostService {
	return &PostService{
		log:                 ctx.NewLoggerHelper("post/service/core-service"),
		postRepo:            uc,
		postVersionRepo:     postVersionRepo,
		searchService:       searchService,
		taskService:         taskService,
		webhookService:      webhookService,
		internalMessageSvc:  internalMessageSvc,
		contentModelRepo:    contentModelRepo,
		referenceTracker:    referenceTracker,
		flowSvc:             flowSvc,
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

	// 引用追踪钩子：回读完整文章（含翻译）抽取媒体引用（best-effort）
	s.syncPostReferences(ctx, dto.GetId())

	// 版本快照钩子：创建即 v1（best-effort，失败不阻断保存）
	s.snapshotPostVersion(ctx, dto.GetId())

	// Webhook 事件外发（best-effort）
	s.dispatchWebhookEvent(ctx, task.WebhookEventPostCreated, dto.GetId(), dto.GetStatus())

	// Flow 规则触发（best-effort；带触发时刻快照供字段级条件）
	if s.flowSvc != nil {
		s.flowSvc.DispatchFlow(ctx, task.WebhookEventPostCreated, dto.GetId(), PostEventSnapshot(dto, ""))
	}

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

	// 记录更新前状态：供 Webhook 事件判定是否为"发布"转移
	var prevStatus contentV1.Post_PostStatus
	if req != nil && req.Data != nil && req.Data.Status != nil {
		if existing, gErr := s.postRepo.Get(ctx, &contentV1.GetPostRequest{
			QueryBy: &contentV1.GetPostRequest_Id{Id: req.GetId()},
		}); gErr == nil && existing != nil {
			prevStatus = existing.GetStatus()
		}
	}

	dto, err := s.postRepo.Update(ctx, req)
	if err != nil {
		return nil, err
	}

	// 双写钩子：更新后入队 ES 重索引（worker 会从 DB 取最新数据 upsert ES）。
	s.enqueuePostReindex(ctx, dto.GetId(), "index")

	// 引用追踪钩子：回读完整文章（含翻译）抽取媒体引用（best-effort）
	s.syncPostReferences(ctx, dto.GetId())

	// 版本快照钩子：每次保存落一版（best-effort，失败不阻断保存）
	s.snapshotPostVersion(ctx, dto.GetId())

	// Webhook 事件外发：转为 PUBLISHED 记 published，其余记 updated（best-effort）
	if dto.GetStatus() == contentV1.Post_POST_STATUS_PUBLISHED && prevStatus != contentV1.Post_POST_STATUS_PUBLISHED {
		s.dispatchWebhookEvent(ctx, task.WebhookEventPostPublished, dto.GetId(), dto.GetStatus())
	} else {
		s.dispatchWebhookEvent(ctx, task.WebhookEventPostUpdated, dto.GetId(), dto.GetStatus())
	}

	// Flow 规则触发：发布转移记 published，其余记 updated（best-effort；快照带 prev_status）
	if s.flowSvc != nil {
		if dto.GetStatus() == contentV1.Post_POST_STATUS_PUBLISHED && prevStatus != contentV1.Post_POST_STATUS_PUBLISHED {
			s.flowSvc.DispatchFlow(ctx, task.WebhookEventPostPublished, dto.GetId(), PostEventSnapshot(dto, prevStatus.String()))
		} else {
			s.flowSvc.DispatchFlow(ctx, task.WebhookEventPostUpdated, dto.GetId(), PostEventSnapshot(dto, prevStatus.String()))
		}
	}

	// 定时发布钩子：以 SCHEDULED 终态保存（含改期）时重新入队到点转正任务。
	if err := s.enqueueScheduledPublish(dto); err != nil {
		return nil, err
	}

	return dto, nil
}

func (s *PostService) Delete(ctx context.Context, req *contentV1.DeletePostRequest) (*emptypb.Empty, error) {
	// Flow 快照在删除前取好（删除后实体不可回查；best-effort）
	var flowSnapshot map[string]string
	if s.flowSvc != nil {
		if existing, gErr := s.postRepo.Get(ctx, &contentV1.GetPostRequest{
			QueryBy: &contentV1.GetPostRequest_Id{Id: req.GetId()},
		}); gErr == nil {
			flowSnapshot = PostEventSnapshot(existing, "")
		}
	}
	err := s.postRepo.Delete(ctx, req)
	if err != nil {
		return nil, err
	}

	// 双写钩子：删除后入队 ES 删除（worker 按 post_id 删 ES 所有语言文档）。
	s.enqueuePostReindex(ctx, req.GetId(), "delete")

	// 引用追踪钩子：清空该文章的引用边并重算媒体引用计数（best-effort）
	if s.referenceTracker != nil {
		s.referenceTracker.RemoveSource(ctx, "post", req.GetId())
	}

	// Webhook 事件外发（best-effort）
	s.dispatchWebhookEvent(ctx, task.WebhookEventPostDeleted, req.GetId(), 0)

	// Flow 规则触发（best-effort；删除后实体不可回查，快照在删除前取好）
	if s.flowSvc != nil {
		s.flowSvc.DispatchFlow(ctx, task.WebhookEventPostDeleted, req.GetId(), flowSnapshot)
	}

	return &emptypb.Empty{}, nil
}

// syncPostReferences 回读完整文章后同步媒体引用（保存钩子，best-effort）。
// 回读而非用请求值：部分更新（update_mask）只带变更字段，回读拿最终态。
func (s *PostService) syncPostReferences(ctx context.Context, postID uint32) {
	if s.referenceTracker == nil {
		return
	}
	post, err := s.postRepo.Get(ctx, &contentV1.GetPostRequest{
		QueryBy: &contentV1.GetPostRequest_Id{Id: postID},
	})
	if err != nil || post == nil {
		s.log.Errorf("reload post [%d] for reference tracking failed: %v", postID, err)
		return
	}
	s.referenceTracker.SyncPostReferences(ctx, postID, post.GetTranslations(), post.GetThumbnail())
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

	// 发布事件外发：webhook 订阅与 Flow 规则同手工发布口径（此前定时转正漏发，顺带修复）
	s.dispatchWebhookEvent(ctx, task.WebhookEventPostPublished, payload.PostID, updated.GetStatus())
	if s.flowSvc != nil {
		s.flowSvc.DispatchFlow(ctx, task.WebhookEventPostPublished, payload.PostID, PostEventSnapshot(updated, "POST_STATUS_SCHEDULED"))
	}

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

// ─── 版本化：保存自动落快照，恢复走常规 Update 链路再落新版本 ───

// ListVersions 列出帖子版本（版本号倒序，仅元数据）。
func (s *PostService) ListVersions(ctx context.Context, req *contentV1.ListPostVersionsRequest) (*contentV1.ListPostVersionsResponse, error) {
	if req == nil || req.GetPostId() == 0 {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}
	items, err := s.postVersionRepo.ListByPost(ctx, req.GetPostId())
	if err != nil {
		return nil, contentV1.ErrorInternalServerError("list post versions failed")
	}
	return &contentV1.ListPostVersionsResponse{
		Items: items,
		Total: uint64(len(items)),
	}, nil
}

// GetVersion 取版本详情（含快照 JSON）。
func (s *PostService) GetVersion(ctx context.Context, req *contentV1.GetPostVersionRequest) (*contentV1.PostVersion, error) {
	if req == nil || req.GetId() == 0 {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}
	return s.postVersionRepo.Get(ctx, req.GetId())
}

// RestoreVersion 恢复帖子到指定版本：快照反序列化后走常规 Update 链路
// （翻译按语言 upsert、分类标签替换、ES 重索引钩子），随后重新快照为新版本。
func (s *PostService) RestoreVersion(ctx context.Context, req *contentV1.RestorePostVersionRequest) (*contentV1.Post, error) {
	if req == nil || req.GetPostId() == 0 || req.GetVersionId() == 0 {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}

	ver, err := s.postVersionRepo.Get(ctx, req.GetVersionId())
	if err != nil {
		return nil, err
	}
	if ver.GetPostId() != req.GetPostId() {
		return nil, contentV1.ErrorBadRequest("version does not belong to this post")
	}

	var snapshot contentV1.Post
	if uErr := (protojson.UnmarshalOptions{DiscardUnknown: true}).
		Unmarshal([]byte(ver.GetContent()), &snapshot); uErr != nil {
		s.log.Errorf("restore post version: corrupted snapshot (version=%d): %v", ver.GetId(), uErr)
		return nil, contentV1.ErrorInternalServerError("corrupted version snapshot")
	}

	// 快照里的 password_hash 已在落版本时剔除；author_id 类归属字段不受 Update 影响
	snapshot.Id = nil
	if _, err = s.postRepo.Update(ctx, &contentV1.UpdatePostRequest{
		Id:   req.GetPostId(),
		Data: &snapshot,
	}); err != nil {
		return nil, err
	}

	// 恢复动作本身落一个新版本，保证版本链完整可回溯
	s.snapshotPostVersion(ctx, req.GetPostId())

	// ES/AI 索引随恢复内容刷新
	s.enqueuePostReindex(ctx, req.GetPostId(), "index")

	return s.postRepo.Get(ctx, &contentV1.GetPostRequest{
		QueryBy: &contentV1.GetPostRequest_Id{Id: req.GetPostId()},
	})
}

// snapshotPostVersion 写入当前帖子完整状态快照。
// best-effort：版本化失败只记日志不阻断保存——内容保存是主链路，版本是附加保障。
func (s *PostService) snapshotPostVersion(ctx context.Context, postID uint32) {
	if postID == 0 {
		return
	}
	full, err := s.postRepo.Get(ctx, &contentV1.GetPostRequest{
		QueryBy: &contentV1.GetPostRequest_Id{Id: postID},
	})
	if err != nil {
		s.log.Errorf("snapshot post version: load post %d failed: %v", postID, err)
		return
	}
	if _, err = s.postVersionRepo.CreateSnapshot(ctx, postID, full, ""); err != nil {
		s.log.Errorf("snapshot post version: create snapshot failed (post=%d): %v", postID, err)
	}
}

// dispatchWebhookEvent Webhook 事件外发统一入口（best-effort，失败只记日志）。
func (s *PostService) dispatchWebhookEvent(ctx context.Context, event string, postID uint32, _ contentV1.Post_PostStatus) {
	if s.webhookService == nil || postID == 0 {
		return
	}
	s.webhookService.DispatchEvent(ctx, event, postID)
}


// ─── 发布审批工作流：DRAFT →（提交审核）审核中 →（通过）PUBLISHED /（驳回）DRAFT ───
//
// 权限卡点在 admin 门面侧的接口 RBAC：接口管理为 review/approve/reject 三个
// 端点绑定可调用角色即完成分权（无需新的权限码机制）。状态机本身在此处强校验。

// loadPostForReview 取帖并校验可流转状态；返回 (post, error)。
func (s *PostService) loadPostForReview(ctx context.Context, postID uint32) (*contentV1.Post, error) {
	post, err := s.postRepo.Get(ctx, &contentV1.GetPostRequest{
		QueryBy: &contentV1.GetPostRequest_Id{Id: postID},
	})
	if err != nil {
		return nil, err
	}
	return post, nil
}

// SubmitReview 提交审核：DRAFT（非审核中）→ in_progress=true。
func (s *PostService) SubmitReview(ctx context.Context, req *contentV1.SubmitPostReviewRequest) (*contentV1.Post, error) {
	if req == nil || req.GetPostId() == 0 {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}
	post, err := s.loadPostForReview(ctx, req.GetPostId())
	if err != nil {
		return nil, err
	}
	if post.GetStatus() != contentV1.Post_POST_STATUS_DRAFT {
		return nil, contentV1.ErrorBadRequest("only draft posts can be submitted for review")
	}
	if post.GetInProgress() {
		return nil, contentV1.ErrorBadRequest("post is already under review")
	}

	updated, err := s.postRepo.Update(ctx, &contentV1.UpdatePostRequest{
		Id:   req.GetPostId(),
		Data: &contentV1.Post{InProgress: trans.Ptr(true)},
		UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"in_progress"}},
	})
	if err != nil {
		return nil, err
	}
	return updated, nil
}

// ApproveReview 审核通过：审核中 → PUBLISHED（publish_time 缺省补当前时间），
// 复用发布链路钩子（ES/AI 重索引、Webhook、快照）并通知作者。
func (s *PostService) ApproveReview(ctx context.Context, req *contentV1.ReviewPostRequest) (*contentV1.Post, error) {
	if req == nil || req.GetPostId() == 0 {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}
	post, err := s.loadPostForReview(ctx, req.GetPostId())
	if err != nil {
		return nil, err
	}
	if !post.GetInProgress() {
		return nil, contentV1.ErrorBadRequest("post is not under review")
	}

	data := &contentV1.Post{
		Status:     trans.Ptr(contentV1.Post_POST_STATUS_PUBLISHED),
		InProgress: trans.Ptr(false),
	}
	mask := []string{"status", "in_progress"}
	if post.GetPublishTime() == nil {
		data.PublishTime = timestamppb.Now()
		mask = append(mask, "publishTime")
	}

	updated, err := s.postRepo.Update(ctx, &contentV1.UpdatePostRequest{
		Id:         req.GetPostId(),
		Data:       data,
		UpdateMask: &fieldmaskpb.FieldMask{Paths: mask},
	})
	if err != nil {
		return nil, err
	}

	// 发布链路钩子：索引可见性变化 + Webhook published 事件（repo.Update 不带这些，这里补齐）
	s.enqueuePostReindex(ctx, req.GetPostId(), "index")
	s.dispatchWebhookEvent(ctx, task.WebhookEventPostPublished, req.GetPostId(), updated.GetStatus())

	// Flow 规则触发（best-effort；快照 prev_status 取审核前状态）
	if s.flowSvc != nil {
		s.flowSvc.DispatchFlow(ctx, task.WebhookEventPostPublished, req.GetPostId(), PostEventSnapshot(updated, post.GetStatus().String()))
	}

	s.notifyReviewResult(ctx, post, true, "")
	return updated, nil
}

// RejectReview 审核驳回：审核中 → DRAFT，带理由通知作者。
func (s *PostService) RejectReview(ctx context.Context, req *contentV1.RejectPostReviewRequest) (*contentV1.Post, error) {
	if req == nil || req.GetPostId() == 0 {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}
	post, err := s.loadPostForReview(ctx, req.GetPostId())
	if err != nil {
		return nil, err
	}
	if !post.GetInProgress() {
		return nil, contentV1.ErrorBadRequest("post is not under review")
	}

	updated, err := s.postRepo.Update(ctx, &contentV1.UpdatePostRequest{
		Id: req.GetPostId(),
		Data: &contentV1.Post{
			InProgress: trans.Ptr(false),
		},
		UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"in_progress"}},
	})
	if err != nil {
		return nil, err
	}

	s.notifyReviewResult(ctx, post, false, req.GetReason())
	return updated, nil
}

// notifyReviewResult 审核结果通知作者（站内信→SSE/邮件出站链自动跟进；best-effort）。
func (s *PostService) notifyReviewResult(ctx context.Context, post *contentV1.Post, approved bool, reason string) {
	if s.internalMessageSvc == nil || post.GetAuthorId() == 0 {
		return
	}
	title := ""
	if len(post.GetTranslations()) > 0 {
		title = post.GetTranslations()[0].GetTitle()
	}
	var content string
	if approved {
		content = "您提交的《" + title + "》已审核通过并发布。"
	} else {
		content = "您提交的《" + title + "》审核未通过。"
		if reason != "" {
			content += "理由：" + reason
		}
	}
	s.internalMessageSvc.NotifyUser(ctx, post.GetAuthorId(), title, content)
}
