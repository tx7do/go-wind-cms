package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	"go-wind-cms/app/core/service/internal/data"

	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
)

type PageService struct {
	contentV1.UnimplementedPageServiceServer

	pageRepo         *data.PageRepo
	contentModelRepo *data.ContentModelRepo
	log              *log.Helper
}

func NewPageService(ctx *bootstrap.Context, uc *data.PageRepo, contentModelRepo *data.ContentModelRepo) *PageService {
	return &PageService{
		log:              ctx.NewLoggerHelper("page/service/core-service"),
		pageRepo:         uc,
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
	return s.pageRepo.Create(ctx, req)
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
	return s.pageRepo.Update(ctx, req)
}

func (s *PageService) Delete(ctx context.Context, req *contentV1.DeletePageRequest) (*emptypb.Empty, error) {
	err := s.pageRepo.Delete(ctx, req)
	if err != nil {
		return nil, err
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
