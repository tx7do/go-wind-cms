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
	log              *log.Helper
}

func NewContentModelService(ctx *bootstrap.Context, contentModelRepo *data.ContentModelRepo) *ContentModelService {
	return &ContentModelService{
		log:              ctx.NewLoggerHelper("content-model/service/core-service"),
		contentModelRepo: contentModelRepo,
	}
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
