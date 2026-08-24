package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	contentV1 "go-wind-cms/api/gen/go/content/service/v1"

	"go-wind-cms/pkg/middleware/auth"
)

// ContentModelService 内容模型管理（admin BFF 转发层）。
// 纯转发至 core ContentModelService，写入操作盖章 operator 为 created_by。
type ContentModelService struct {
	adminV1.ContentModelServiceHTTPServer

	contentModelServiceClient contentV1.ContentModelServiceClient
	log                       *log.Helper
}

func NewContentModelService(ctx *bootstrap.Context, contentModelServiceClient contentV1.ContentModelServiceClient) *ContentModelService {
	return &ContentModelService{
		log:                       ctx.NewLoggerHelper("content-model/service/admin-service"),
		contentModelServiceClient: contentModelServiceClient,
	}
}

func (s *ContentModelService) List(ctx context.Context, req *paginationV1.PagingRequest) (*contentV1.ListContentModelResponse, error) {
	return s.contentModelServiceClient.List(ctx, req)
}

func (s *ContentModelService) Get(ctx context.Context, req *contentV1.GetContentModelRequest) (*contentV1.ContentModel, error) {
	return s.contentModelServiceClient.Get(ctx, req)
}

func (s *ContentModelService) Create(ctx context.Context, req *contentV1.CreateContentModelRequest) (*contentV1.ContentModel, error) {
	if req == nil || req.Data == nil {
		return nil, adminV1.ErrorBadRequest("invalid parameter")
	}

	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	req.Data.CreatedBy = trans.Ptr(operator.UserId)

	return s.contentModelServiceClient.Create(ctx, req)
}

func (s *ContentModelService) Update(ctx context.Context, req *contentV1.UpdateContentModelRequest) (*contentV1.ContentModel, error) {
	if req == nil || req.Data == nil {
		return nil, adminV1.ErrorBadRequest("invalid parameter")
	}

	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	req.Data.UpdatedBy = trans.Ptr(operator.UserId)

	return s.contentModelServiceClient.Update(ctx, req)
}

func (s *ContentModelService) Delete(ctx context.Context, req *contentV1.DeleteContentModelRequest) (*emptypb.Empty, error) {
	return s.contentModelServiceClient.Delete(ctx, req)
}

func (s *ContentModelService) ListFieldDefinitions(ctx context.Context, req *contentV1.ListFieldDefinitionsRequest) (*contentV1.ListFieldDefinitionsResponse, error) {
	return s.contentModelServiceClient.ListFieldDefinitions(ctx, req)
}
