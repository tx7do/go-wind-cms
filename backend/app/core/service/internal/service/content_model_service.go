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

// ContentModelService 内容模型服务：模型 CRUD + 字段定义查询 + 字段值校验。
type ContentModelService struct {
	contentV1.UnimplementedContentModelServiceServer

	contentModelRepo *data.ContentModelRepo
	postRepo         *data.PostRepo
	log              *log.Helper
}

func NewContentModelService(ctx *bootstrap.Context, contentModelRepo *data.ContentModelRepo, postRepo *data.PostRepo) *ContentModelService {
	return &ContentModelService{
		log:              ctx.NewLoggerHelper("content-model/service/core-service"),
		contentModelRepo: contentModelRepo,
		postRepo:         postRepo,
	}
}

// GetModelByPost 按文章取其分类链绑定的内容模型（C 端渲染 custom_fields 用）。
// 文章→分类→模型的解析在服务端完成，调用方无法指定模型或借道越权；
// 文章不存在/未绑定模型返回 NotFound/空模型（由门面决定公开语义）。
func (s *ContentModelService) GetModelByPost(ctx context.Context, req *contentV1.GetModelByPostRequest) (*contentV1.ContentModel, error) {
	if req == nil || req.GetPostId() == 0 {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}
	post, err := s.postRepo.Get(ctx, &contentV1.GetPostRequest{
		QueryBy: &contentV1.GetPostRequest_Id{Id: req.GetPostId()},
	})
	if err != nil {
		return nil, err
	}
	modelID := s.contentModelRepo.ResolveModelIDByCategories(ctx, post.GetCategoryIds())
	if modelID == 0 {
		return &contentV1.ContentModel{}, nil
	}
	return s.contentModelRepo.Get(ctx, &contentV1.GetContentModelRequest{
		QueryBy: &contentV1.GetContentModelRequest_Id{Id: modelID},
	})
}

func (s *ContentModelService) List(ctx context.Context, req *paginationV1.PagingRequest) (*contentV1.ListContentModelResponse, error) {
	return s.contentModelRepo.List(ctx, req)
}

func (s *ContentModelService) Count(ctx context.Context, req *paginationV1.PagingRequest) (*contentV1.CountContentModelResponse, error) {
	count, err := s.contentModelRepo.Count(ctx, req)
	if err != nil {
		return nil, err
	}
	return &contentV1.CountContentModelResponse{Count: uint64(count)}, nil
}

func (s *ContentModelService) Get(ctx context.Context, req *contentV1.GetContentModelRequest) (*contentV1.ContentModel, error) {
	return s.contentModelRepo.Get(ctx, req)
}

func (s *ContentModelService) Create(ctx context.Context, req *contentV1.CreateContentModelRequest) (*contentV1.ContentModel, error) {
	return s.contentModelRepo.Create(ctx, req)
}

func (s *ContentModelService) Update(ctx context.Context, req *contentV1.UpdateContentModelRequest) (*contentV1.ContentModel, error) {
	return s.contentModelRepo.Update(ctx, req)
}

func (s *ContentModelService) Delete(ctx context.Context, req *contentV1.DeleteContentModelRequest) (*emptypb.Empty, error) {
	if err := s.contentModelRepo.Delete(ctx, req); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *ContentModelService) ListFieldDefinitions(ctx context.Context, req *contentV1.ListFieldDefinitionsRequest) (*contentV1.ListFieldDefinitionsResponse, error) {
	return s.contentModelRepo.ListFieldDefinitions(ctx, req)
}
