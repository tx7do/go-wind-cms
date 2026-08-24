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

type CategoryService struct {
	contentV1.UnimplementedCategoryServiceServer

	categoryRepo      *data.CategoryRepo
	contentModelRepo  *data.ContentModelRepo
	log               *log.Helper
}

func NewCategoryService(ctx *bootstrap.Context, uc *data.CategoryRepo, contentModelRepo *data.ContentModelRepo) *CategoryService {
	return &CategoryService{
		log:              ctx.NewLoggerHelper("category/service/core-service"),
		categoryRepo:     uc,
		contentModelRepo: contentModelRepo,
	}
}

func (s *CategoryService) List(ctx context.Context, req *paginationV1.PagingRequest) (*contentV1.ListCategoryResponse, error) {
	return s.categoryRepo.List(ctx, req)
}

func (s *CategoryService) Get(ctx context.Context, req *contentV1.GetCategoryRequest) (*contentV1.Category, error) {
	return s.categoryRepo.Get(ctx, req)
}

func (s *CategoryService) Create(ctx context.Context, req *contentV1.CreateCategoryRequest) (*contentV1.Category, error) {
	// 自定义字段校验：按分类自身绑定的内容模型校验 custom_fields
	// 自定义字段校验：按分类自身绑定的内容模型校验 custom_fields
	// TenantId 不在此消息中（由 viewer 注入），租户归属校验由后端
	// 的 viewer 租户上下文兜底；此处传 0 跳过 ValidateValues 的跨租户检查
	if req != nil && req.Data != nil {
		if vErr := s.contentModelRepo.ValidateValues(ctx,
			req.Data.GetContentModelId(),
			req.Data.GetCustomFields(),
			0); vErr != nil {
			return nil, vErr
		}
	}
	return s.categoryRepo.Create(ctx, req)
}

func (s *CategoryService) Update(ctx context.Context, req *contentV1.UpdateCategoryRequest) (*contentV1.Category, error) {
	// 自定义字段校验：仅当本次请求携带 custom_fields 时校验。
	// 模型优先取请求绑定；未携带时查现库绑定，保证换绑后旧值也能被复检。
	if req != nil && req.Data != nil && req.Data.CustomFields != nil {
		modelID := req.Data.GetContentModelId()
		if modelID == 0 {
			if existing, gErr := s.categoryRepo.Get(ctx, &contentV1.GetCategoryRequest{
				QueryBy: &contentV1.GetCategoryRequest_Id{Id: req.GetId()},
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
	return s.categoryRepo.Update(ctx, req)
}

func (s *CategoryService) Delete(ctx context.Context, req *contentV1.DeleteCategoryRequest) (*emptypb.Empty, error) {
	err := s.categoryRepo.Delete(ctx, req)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CategoryService) TranslationExists(ctx context.Context, req *contentV1.CategoryTranslationExistsRequest) (*contentV1.CategoryTranslationExistsResponse, error) {
	exists, err := s.categoryRepo.TranslationExists(ctx, req.GetCategoryId(), req.GetLanguageCode())
	if err != nil {
		return nil, err
	}

	return &contentV1.CategoryTranslationExistsResponse{
		Exists: exists,
	}, nil
}

func (s *CategoryService) GetTranslation(ctx context.Context, req *contentV1.GetCategoryRequest) (*contentV1.CategoryTranslation, error) {
	return s.categoryRepo.GetTranslation(ctx, req)
}

func (s *CategoryService) CreateTranslation(ctx context.Context, req *contentV1.CreateCategoryTranslationRequest) (*contentV1.CategoryTranslation, error) {
	return s.categoryRepo.CreateTranslation(ctx, req)
}

func (s *CategoryService) UpdateTranslation(ctx context.Context, req *contentV1.UpdateCategoryTranslationRequest) (*contentV1.CategoryTranslation, error) {
	return s.categoryRepo.UpdateTranslation(ctx, req)
}

func (s *CategoryService) DeleteTranslation(ctx context.Context, req *contentV1.DeleteCategoryTranslationRequest) (*emptypb.Empty, error) {
	err := s.categoryRepo.DeleteTranslation(ctx, req)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
