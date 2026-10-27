package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	"go-wind-cms/app/core/service/internal/data"

	"go-wind-cms/pkg/task"

	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
)

type PageService struct {
	contentV1.UnimplementedPageServiceServer

	pageRepo         *data.PageRepo
	pageVersionRepo  *data.PageVersionRepo
	webhookService   *WebhookService
	notifySvc        *InternalMessageService
	contentModelRepo *data.ContentModelRepo
	referenceTracker *ReferenceTrackingService
	flowSvc          *FlowService
	reviewLogRepo    *data.ReviewLogRepo
	reviewPolicy     *ReviewPolicy
	log              *log.Helper
}

func NewPageService(ctx *bootstrap.Context, uc *data.PageRepo, contentModelRepo *data.ContentModelRepo, pageVersionRepo *data.PageVersionRepo, webhookService *WebhookService, notifySvc *InternalMessageService, referenceTracker *ReferenceTrackingService, flowSvc *FlowService,
	reviewLogRepo *data.ReviewLogRepo, reviewPolicy *ReviewPolicy) *PageService {
	return &PageService{
		log:              ctx.NewLoggerHelper("page/service/core-service"),
		pageRepo:         uc,
		pageVersionRepo:  pageVersionRepo,
		webhookService:   webhookService,
		notifySvc:        notifySvc,
		contentModelRepo: contentModelRepo,
		referenceTracker: referenceTracker,
		flowSvc:          flowSvc,
		reviewLogRepo:    reviewLogRepo,
		reviewPolicy:     reviewPolicy,
	}
}

func (s *PageService) List(ctx context.Context, req *paginationV1.PagingRequest) (*contentV1.ListPageResponse, error) {
	return s.pageRepo.List(ctx, req)
}

func (s *PageService) Get(ctx context.Context, req *contentV1.GetPageRequest) (*contentV1.Page, error) {
	return s.pageRepo.Get(ctx, req)
}

func (s *PageService) Create(ctx context.Context, req *contentV1.CreatePageRequest) (*contentV1.Page, error) {
	// 自定义字段校验：按页面自身绑定的内容模型校验 custom_fields
	// TenantId 由 viewer 注入（不在此消息中），此处传 0 跳过跨租户检查
	if req != nil && req.Data != nil {
		if vErr := s.contentModelRepo.ValidateValues(ctx,
			req.Data.GetContentModelId(),
			req.Data.GetCustomFields(),
			0); vErr != nil {
			return nil, vErr
		}
	}
	dto, err := s.pageRepo.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	// 引用追踪钩子：回读完整页面抽取媒体引用（best-effort）
	s.syncPageReferences(ctx, dto.GetId())
	// 版本快照钩子：创建即 v1（best-effort，失败不阻断保存）
	s.snapshotPageVersion(ctx, dto.GetId())
	// Webhook 事件外发（best-effort）
	if s.webhookService != nil && dto.GetId() != 0 {
		s.webhookService.DispatchEvent(ctx, task.WebhookEventPageCreated, dto.GetId())
	}
	// Flow 规则触发（best-effort；带触发时刻快照供字段级条件）
	if s.flowSvc != nil && dto.GetId() != 0 {
		s.flowSvc.DispatchFlow(ctx, task.WebhookEventPageCreated, dto.GetId(), PageEventSnapshot(dto, ""))
	}
	return dto, nil
}

func (s *PageService) Update(ctx context.Context, req *contentV1.UpdatePageRequest) (*contentV1.Page, error) {
	// 自定义字段校验：仅当本次请求携带 custom_fields 时校验。
	// 模型优先取请求绑定；未携带时查现库绑定，保证换绑后旧值也能被复检。
	if req != nil && req.Data != nil && req.Data.CustomFields != nil {
		modelID := req.Data.GetContentModelId()
		if modelID == 0 {
			if existing, gErr := s.pageRepo.Get(ctx, &contentV1.GetPageRequest{
				QueryBy: &contentV1.GetPageRequest_Id{Id: req.GetId()},
			}); gErr == nil && existing != nil {
				modelID = existing.GetContentModelId()
			}
		}
		if vErr := s.contentModelRepo.ValidateValues(ctx,
			modelID,
			req.Data.GetCustomFields(),
			0); vErr != nil {
			return nil, vErr
		}
	}
	// 记录更新前状态：供 Webhook 事件判定是否为"发布"转移
	var prevStatus contentV1.Page_PageStatus
	if req != nil {
		if existing, gErr := s.pageRepo.Get(ctx, &contentV1.GetPageRequest{
			QueryBy: &contentV1.GetPageRequest_Id{Id: req.GetId()},
		}); gErr == nil && existing != nil {
			prevStatus = existing.GetStatus()
			// 审批闸门：审核中的页面不允许经 Update 直改状态（与 post 同款）
			if req.Data != nil && req.Data.Status != nil &&
				existing.GetInProgress() && req.Data.GetStatus() != prevStatus {
				return nil, contentV1.ErrorBadRequest(
					"page is under review: status changes must go through the review flow (approve/reject)")
			}
		}
	}

	dto, err := s.pageRepo.Update(ctx, req)
	if err != nil {
		return nil, err
	}
	// 引用追踪钩子：回读完整页面抽取媒体引用（best-effort）
	s.syncPageReferences(ctx, dto.GetId())
	// 版本快照钩子：每次保存落一版（best-effort，失败不阻断保存）
	s.snapshotPageVersion(ctx, dto.GetId())
	// Webhook 事件外发：转为 PUBLISHED 记 published，其余记 updated（best-effort）
	if s.webhookService != nil && dto.GetId() != 0 {
		if dto.GetStatus() == contentV1.Page_PAGE_STATUS_PUBLISHED && prevStatus != contentV1.Page_PAGE_STATUS_PUBLISHED {
			s.webhookService.DispatchEvent(ctx, task.WebhookEventPagePublished, dto.GetId())
		} else {
			s.webhookService.DispatchEvent(ctx, task.WebhookEventPageUpdated, dto.GetId())
		}
	}
	// Flow 规则触发：同 webhook 的 published/updated 判定（best-effort；快照带 prev_status）
	if s.flowSvc != nil && dto.GetId() != 0 {
		if dto.GetStatus() == contentV1.Page_PAGE_STATUS_PUBLISHED && prevStatus != contentV1.Page_PAGE_STATUS_PUBLISHED {
			s.flowSvc.DispatchFlow(ctx, task.WebhookEventPagePublished, dto.GetId(), PageEventSnapshot(dto, prevStatus.String()))
		} else {
			s.flowSvc.DispatchFlow(ctx, task.WebhookEventPageUpdated, dto.GetId(), PageEventSnapshot(dto, prevStatus.String()))
		}
	}
	return dto, nil
}

func (s *PageService) Delete(ctx context.Context, req *contentV1.DeletePageRequest) (*emptypb.Empty, error) {
	err := s.pageRepo.Delete(ctx, req)
	if err != nil {
		return nil, err
	}
	// 引用追踪钩子：清空该页面的引用边并重算媒体引用计数（best-effort）
	if s.referenceTracker != nil {
		s.referenceTracker.RemoveSource(ctx, "page", req.GetId())
	}
	// Webhook 事件外发（best-effort）
	if s.webhookService != nil && req.GetId() != 0 {
		s.webhookService.DispatchEvent(ctx, task.WebhookEventPageDeleted, req.GetId())
	}
	// Flow 规则触发（best-effort；快照在删除前取好，删除后实体不可回查）
	if s.flowSvc != nil && req.GetId() != 0 {
		var flowSnapshot map[string]string
		if existing, gErr := s.pageRepo.Get(ctx, &contentV1.GetPageRequest{
			QueryBy: &contentV1.GetPageRequest_Id{Id: req.GetId()},
		}); gErr == nil {
			flowSnapshot = PageEventSnapshot(existing, "")
		}
		s.flowSvc.DispatchFlow(ctx, task.WebhookEventPageDeleted, req.GetId(), flowSnapshot)
	}
	return &emptypb.Empty{}, nil
}

// syncPageReferences 回读完整页面后同步媒体引用（保存钩子，best-effort）。
func (s *PageService) syncPageReferences(ctx context.Context, pageID uint32) {
	if s.referenceTracker == nil {
		return
	}
	page, err := s.pageRepo.Get(ctx, &contentV1.GetPageRequest{
		QueryBy: &contentV1.GetPageRequest_Id{Id: pageID},
	})
	if err != nil || page == nil {
		s.log.Errorf("reload page [%d] for reference tracking failed: %v", pageID, err)
		return
	}
	s.referenceTracker.SyncPageReferences(ctx, pageID, page.GetTranslations(), page.GetSections())
}

func (s *PageService) TranslationExists(ctx context.Context, req *contentV1.PageTranslationExistsRequest) (*contentV1.PageTranslationExistsResponse, error) {
	exists, err := s.pageRepo.TranslationExists(ctx, req.GetPageId(), req.GetLanguageCode())
	if err != nil {
		return nil, err
	}

	return &contentV1.PageTranslationExistsResponse{
		Exists: exists,
	}, nil
}

func (s *PageService) GetTranslation(ctx context.Context, req *contentV1.GetPageRequest) (*contentV1.PageTranslation, error) {
	return s.pageRepo.GetTranslation(ctx, req)
}

func (s *PageService) CreateTranslation(ctx context.Context, req *contentV1.CreatePageTranslationRequest) (*contentV1.PageTranslation, error) {
	return s.pageRepo.CreateTranslation(ctx, req)
}

func (s *PageService) UpdateTranslation(ctx context.Context, req *contentV1.UpdatePageTranslationRequest) (*contentV1.PageTranslation, error) {
	return s.pageRepo.UpdateTranslation(ctx, req)
}

func (s *PageService) DeleteTranslation(ctx context.Context, req *contentV1.DeletePageTranslationRequest) (*emptypb.Empty, error) {
	err := s.pageRepo.DeleteTranslation(ctx, req)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// ─── 版本化：保存自动落快照，恢复走常规 Update 链路再落新版本 ───

// ListVersions 列出页面版本（版本号倒序，仅元数据）。
func (s *PageService) ListVersions(ctx context.Context, req *contentV1.ListPageVersionsRequest) (*contentV1.ListPageVersionsResponse, error) {
	if req == nil || req.GetPageId() == 0 {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}
	items, err := s.pageVersionRepo.ListByPage(ctx, req.GetPageId())
	if err != nil {
		return nil, contentV1.ErrorInternalServerError("list page versions failed")
	}
	return &contentV1.ListPageVersionsResponse{
		Items: items,
		Total: uint64(len(items)),
	}, nil
}

// GetVersion 取版本详情（含快照 JSON）。
func (s *PageService) GetVersion(ctx context.Context, req *contentV1.GetPageVersionRequest) (*contentV1.PageVersion, error) {
	if req == nil || req.GetId() == 0 {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}
	return s.pageVersionRepo.Get(ctx, req.GetId())
}

// RestoreVersion 恢复页面到指定版本：快照反序列化后走常规 Update 链路
// （翻译/sections 随 Update 整体替换），随后重新快照为新版本。
func (s *PageService) RestoreVersion(ctx context.Context, req *contentV1.RestorePageVersionRequest) (*contentV1.Page, error) {
	if req == nil || req.GetPageId() == 0 || req.GetVersionId() == 0 {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}

	ver, err := s.pageVersionRepo.Get(ctx, req.GetVersionId())
	if err != nil {
		return nil, err
	}
	if ver.GetPageId() != req.GetPageId() {
		return nil, contentV1.ErrorBadRequest("version does not belong to this page")
	}

	var snapshot contentV1.Page
	if uErr := (protojson.UnmarshalOptions{DiscardUnknown: true}).
		Unmarshal([]byte(ver.GetContent()), &snapshot); uErr != nil {
		s.log.Errorf("restore page version: corrupted snapshot (version=%d): %v", ver.GetId(), uErr)
		return nil, contentV1.ErrorInternalServerError("corrupted version snapshot")
	}

	snapshot.Id = nil
	if _, err = s.pageRepo.Update(ctx, &contentV1.UpdatePageRequest{
		Id:   req.GetPageId(),
		Data: &snapshot,
	}); err != nil {
		return nil, err
	}

	// 恢复动作本身落一个新版本，保证版本链完整可回溯
	s.snapshotPageVersion(ctx, req.GetPageId())

	return s.pageRepo.Get(ctx, &contentV1.GetPageRequest{
		QueryBy: &contentV1.GetPageRequest_Id{Id: req.GetPageId()},
	})
}

// snapshotPageVersion 写入当前页面完整状态快照（含 translations/sections）。
// best-effort：版本化失败只记日志不阻断保存——内容保存是主链路，版本是附加保障。
func (s *PageService) snapshotPageVersion(ctx context.Context, pageID uint32) {
	if pageID == 0 {
		return
	}
	full, err := s.pageRepo.Get(ctx, &contentV1.GetPageRequest{
		QueryBy: &contentV1.GetPageRequest_Id{Id: pageID},
	})
	if err != nil {
		s.log.Errorf("snapshot page version: load page %d failed: %v", pageID, err)
		return
	}
	if _, err = s.pageVersionRepo.CreateSnapshot(ctx, pageID, full, ""); err != nil {
		s.log.Errorf("snapshot page version: create snapshot failed (page=%d): %v", pageID, err)
	}
}

// ─── 页面审批工作流：DRAFT →（提交审核）审核中 →（通过）PUBLISHED /（驳回）DRAFT ───
//
// 权限卡点在 admin 门面侧的接口 RBAC：接口管理为 review/approve/reject 三个
// 端点绑定可调用角色即完成分权（无需新的权限码机制）。状态机本身在此处强校验。

// loadPageForReview 取页面并校验可流转状态；返回 (page, error)。
func (s *PageService) loadPageForReview(ctx context.Context, pageID uint32) (*contentV1.Page, error) {
	page, err := s.pageRepo.Get(ctx, &contentV1.GetPageRequest{
		QueryBy: &contentV1.GetPageRequest_Id{Id: pageID},
	})
	if err != nil {
		return nil, err
	}
	return page, nil
}

// SubmitReview 提交审核：DRAFT（非审核中）→ in_progress=true。
func (s *PageService) SubmitReview(ctx context.Context, req *contentV1.SubmitPageReviewRequest) (*contentV1.Page, error) {
	if req == nil || req.GetPageId() == 0 {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}
	page, err := s.loadPageForReview(ctx, req.GetPageId())
	if err != nil {
		return nil, err
	}
	if page.GetStatus() != contentV1.Page_PAGE_STATUS_DRAFT {
		return nil, contentV1.ErrorBadRequest("only draft pages can be submitted for review")
	}
	if page.GetInProgress() {
		return nil, contentV1.ErrorBadRequest("page is already under review")
	}

	updated, err := s.pageRepo.Update(ctx, &contentV1.UpdatePageRequest{
		Id:         req.GetPageId(),
		Data:       &contentV1.Page{InProgress: trans.Ptr(true)},
		UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"inProgress"}},
	})
	if err != nil {
		return nil, err
	}
	// 审批留痕：开启新批次（level=0）
	if s.reviewLogRepo != nil {
		if operatorID, tenantID, ok := operatorTenantUserID(ctx); ok {
			_ = s.reviewLogRepo.Create(ctx, "page", req.GetPageId(), "SUBMIT", 0, "", operatorID, tenantID)
		}
	}
	return updated, nil
}

// ApproveReview 审核通过：审核中 → PUBLISHED，dispatch page.published Webhook 并通知作者。
func (s *PageService) ApproveReview(ctx context.Context, req *contentV1.ReviewPageRequest) (*contentV1.Page, error) {
	if req == nil || req.GetPageId() == 0 {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}
	page, err := s.loadPageForReview(ctx, req.GetPageId())
	if err != nil {
		return nil, err
	}
	if !page.GetInProgress() {
		return nil, contentV1.ErrorBadRequest("page is not under review")
	}

	// ===== 多级审批推进（与 post 同款语义）=====
	operatorID, tenantID, _ := operatorTenantUserID(ctx)
	totalLevels := s.reviewPolicy.Levels(ctx, tenantID)
	currentLevel, err := s.reviewLogRepo.CurrentLevel(ctx, "page", req.GetPageId())
	if err != nil {
		return nil, err
	}
	nextLevel := currentLevel + 1
	if err = s.reviewPolicy.CheckApprover(ctx, tenantID, operatorID, int(nextLevel)); err != nil {
		return nil, err
	}

	finalApproval := int(nextLevel) >= totalLevels
	if !finalApproval {
		updated, uerr := s.pageRepo.Update(ctx, &contentV1.UpdatePageRequest{
			Id: req.GetPageId(),
			Data: &contentV1.Page{
				InProgress: trans.Ptr(true),
			},
			UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"inProgress"}},
		})
		if uerr != nil {
			return nil, uerr
		}
		if s.reviewLogRepo != nil {
			_ = s.reviewLogRepo.Create(ctx, "page", req.GetPageId(), "APPROVE", nextLevel, "", operatorID, tenantID)
		}
		s.notifyPageReviewResult(ctx, page, true, "")
		return updated, nil
	}

	updated, err := s.pageRepo.Update(ctx, &contentV1.UpdatePageRequest{
		Id: req.GetPageId(),
		Data: &contentV1.Page{
			Status:     trans.Ptr(contentV1.Page_PAGE_STATUS_PUBLISHED),
			InProgress: trans.Ptr(false),
		},
		UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"status", "inProgress"}},
	})
	if err != nil {
		return nil, err
	}
	if s.reviewLogRepo != nil {
		_ = s.reviewLogRepo.Create(ctx, "page", req.GetPageId(), "APPROVE", nextLevel, "", operatorID, tenantID)
	}

	// 发布链路钩子：Webhook published 事件（页面无 ES 索引，无重索引链）
	if s.webhookService != nil {
		s.webhookService.DispatchEvent(ctx, task.WebhookEventPagePublished, req.GetPageId())
	}
	// Flow 规则触发（best-effort；快照 prev_status 取审核前状态）
	if s.flowSvc != nil {
		s.flowSvc.DispatchFlow(ctx, task.WebhookEventPagePublished, req.GetPageId(), PageEventSnapshot(updated, page.GetStatus().String()))
	}

	s.notifyPageReviewResult(ctx, page, true, "")
	return updated, nil
}

// RejectReview 审核驳回：审核中 → DRAFT，带理由通知作者。
func (s *PageService) RejectReview(ctx context.Context, req *contentV1.RejectPageReviewRequest) (*contentV1.Page, error) {
	if req == nil || req.GetPageId() == 0 {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}
	page, err := s.loadPageForReview(ctx, req.GetPageId())
	if err != nil {
		return nil, err
	}
	if !page.GetInProgress() {
		return nil, contentV1.ErrorBadRequest("page is not under review")
	}

	updated, err := s.pageRepo.Update(ctx, &contentV1.UpdatePageRequest{
		Id: req.GetPageId(),
		Data: &contentV1.Page{
			InProgress: trans.Ptr(false),
		},
		UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"inProgress"}},
	})
	if err != nil {
		return nil, err
	}
	if s.reviewLogRepo != nil {
		if operatorID, tenantID, ok := operatorTenantUserID(ctx); ok {
			_ = s.reviewLogRepo.Create(ctx, "page", req.GetPageId(), "REJECT", 0, req.GetReason(), operatorID, tenantID)
		}
	}

	s.notifyPageReviewResult(ctx, page, false, req.GetReason())
	return updated, nil
}

// notifyPageReviewResult 审核结果通知作者（站内信→SSE/邮件出站链自动跟进；best-effort）。
func (s *PageService) notifyPageReviewResult(ctx context.Context, page *contentV1.Page, approved bool, reason string) {
	if s.notifySvc == nil || page.GetAuthorId() == 0 {
		return
	}
	title := ""
	for _, tr := range page.GetTranslations() {
		if tr.GetTitle() != "" {
			title = tr.GetTitle()
			break
		}
	}
	var content string
	if approved {
		content = "您提交的页面《" + title + "》已审核通过并发布。"
	} else {
		content = "您提交的页面《" + title + "》审核未通过。"
		if reason != "" {
			content += "理由：" + reason
		}
	}
	s.notifySvc.NotifyUser(ctx, page.GetAuthorId(), title, content)
}
