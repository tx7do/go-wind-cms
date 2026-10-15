package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/emptypb"

	"go-wind-cms/app/core/service/internal/data"

	"go-wind-cms/pkg/task"

	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
)

type PageService struct {
	contentV1.UnimplementedPageServiceServer

	pageRepo         *data.PageRepo
	pageVersionRepo  *data.PageVersionRepo
	webhookService   *WebhookService
	contentModelRepo *data.ContentModelRepo
	log              *log.Helper
}

func NewPageService(ctx *bootstrap.Context, uc *data.PageRepo, contentModelRepo *data.ContentModelRepo, pageVersionRepo *data.PageVersionRepo, webhookService *WebhookService) *PageService {
	return &PageService{
		log:              ctx.NewLoggerHelper("page/service/core-service"),
		pageRepo:         uc,
		pageVersionRepo:  pageVersionRepo,
		webhookService:   webhookService,
		contentModelRepo: contentModelRepo,
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
	// 版本快照钩子：创建即 v1（best-effort，失败不阻断保存）
	s.snapshotPageVersion(ctx, dto.GetId())
	// Webhook 事件外发（best-effort）
	if s.webhookService != nil && dto.GetId() != 0 {
		s.webhookService.DispatchEvent(ctx, task.WebhookEventPageCreated, dto.GetId())
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
		}
	}

	dto, err := s.pageRepo.Update(ctx, req)
	if err != nil {
		return nil, err
	}
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
	return dto, nil
}

func (s *PageService) Delete(ctx context.Context, req *contentV1.DeletePageRequest) (*emptypb.Empty, error) {
	err := s.pageRepo.Delete(ctx, req)
	if err != nil {
		return nil, err
	}
	// Webhook 事件外发（best-effort）
	if s.webhookService != nil && req.GetId() != 0 {
		s.webhookService.DispatchEvent(ctx, task.WebhookEventPageDeleted, req.GetId())
	}
	return &emptypb.Empty{}, nil
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
