package service

import (
	"context"

	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
)

// ComponentDefinitionService 组件定义管理门面（内容建模 Components）：
// 透传 core ContentModelService 的组件 CRUD，租户隔离由 core viewer 保证。
type ComponentDefinitionService struct {
	adminV1.UnimplementedComponentDefinitionServiceServer

	contentModelClient contentV1.ContentModelServiceClient
}

func NewComponentDefinitionService(ctx *bootstrap.Context, contentModelClient contentV1.ContentModelServiceClient) *ComponentDefinitionService {
	return &ComponentDefinitionService{
		contentModelClient: contentModelClient,
	}
}

func (s *ComponentDefinitionService) List(ctx context.Context, req *paginationV1.PagingRequest) (*contentV1.ListComponentDefinitionsResponse, error) {
	return s.contentModelClient.ListComponentDefinitions(ctx, req)
}

func (s *ComponentDefinitionService) Get(ctx context.Context, req *contentV1.GetComponentDefinitionRequest) (*contentV1.ComponentDefinition, error) {
	return s.contentModelClient.GetComponentDefinition(ctx, req)
}

func (s *ComponentDefinitionService) Create(ctx context.Context, req *contentV1.CreateComponentDefinitionRequest) (*contentV1.ComponentDefinition, error) {
	return s.contentModelClient.CreateComponentDefinition(ctx, req)
}

func (s *ComponentDefinitionService) Update(ctx context.Context, req *contentV1.UpdateComponentDefinitionRequest) (*contentV1.ComponentDefinition, error) {
	return s.contentModelClient.UpdateComponentDefinition(ctx, req)
}

func (s *ComponentDefinitionService) Delete(ctx context.Context, req *contentV1.DeleteComponentDefinitionRequest) (*emptypb.Empty, error) {
	return s.contentModelClient.DeleteComponentDefinition(ctx, req)
}

func (s *ComponentDefinitionService) SaveFieldDefinitions(ctx context.Context, req *contentV1.SaveComponentFieldDefinitionsRequest) (*contentV1.ListFieldDefinitionsResponse, error) {
	return s.contentModelClient.SaveComponentFieldDefinitions(ctx, req)
}

func (s *ComponentDefinitionService) ListFieldDefinitions(ctx context.Context, req *contentV1.ListFieldDefinitionsRequest) (*contentV1.ListFieldDefinitionsResponse, error) {
	return s.contentModelClient.ListFieldDefinitions(ctx, req)
}
