package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	aiV1 "go-wind-cms/api/gen/go/ai/service/v1"
	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
)

// AiKnowledgeService 知识库管理（BFF 直通转发；
// 操作者归属与切片/向量链路在 core 侧完成）。
type AiKnowledgeService struct {
	adminV1.AiKnowledgeBaseServiceHTTPServer

	aiKnowledgeBaseServiceClient aiV1.AiKnowledgeBaseServiceClient
	log                          *log.Helper
}

func NewAiKnowledgeService(ctx *bootstrap.Context, aiKnowledgeBaseServiceClient aiV1.AiKnowledgeBaseServiceClient) *AiKnowledgeService {
	return &AiKnowledgeService{
		log:                          ctx.NewLoggerHelper("ai_knowledge/service/admin-service"),
		aiKnowledgeBaseServiceClient: aiKnowledgeBaseServiceClient,
	}
}

func (s *AiKnowledgeService) List(ctx context.Context, req *paginationV1.PagingRequest) (*aiV1.ListAiKnowledgeBaseResponse, error) {
	return s.aiKnowledgeBaseServiceClient.List(ctx, req)
}

func (s *AiKnowledgeService) Get(ctx context.Context, req *aiV1.GetAiKnowledgeBaseRequest) (*aiV1.AiKnowledgeBase, error) {
	return s.aiKnowledgeBaseServiceClient.Get(ctx, req)
}

func (s *AiKnowledgeService) Create(ctx context.Context, req *aiV1.CreateAiKnowledgeBaseRequest) (*emptypb.Empty, error) {
	return s.aiKnowledgeBaseServiceClient.Create(ctx, req)
}

func (s *AiKnowledgeService) Update(ctx context.Context, req *aiV1.UpdateAiKnowledgeBaseRequest) (*emptypb.Empty, error) {
	return s.aiKnowledgeBaseServiceClient.Update(ctx, req)
}

func (s *AiKnowledgeService) Delete(ctx context.Context, req *aiV1.DeleteAiKnowledgeBaseRequest) (*emptypb.Empty, error) {
	return s.aiKnowledgeBaseServiceClient.Delete(ctx, req)
}

func (s *AiKnowledgeService) UploadDoc(ctx context.Context, req *aiV1.UploadAiDocRequest) (*aiV1.UploadAiDocResponse, error) {
	return s.aiKnowledgeBaseServiceClient.UploadDoc(ctx, req)
}

func (s *AiKnowledgeService) UploadDocFile(ctx context.Context, req *aiV1.UploadAiDocFileRequest) (*aiV1.UploadAiDocResponse, error) {
	return s.aiKnowledgeBaseServiceClient.UploadDocFile(ctx, req)
}

func (s *AiKnowledgeService) ListDocs(ctx context.Context, req *aiV1.ListAiDocsRequest) (*aiV1.ListAiDocsResponse, error) {
	return s.aiKnowledgeBaseServiceClient.ListDocs(ctx, req)
}

func (s *AiKnowledgeService) DeleteDoc(ctx context.Context, req *aiV1.DeleteAiDocRequest) (*emptypb.Empty, error) {
	return s.aiKnowledgeBaseServiceClient.DeleteDoc(ctx, req)
}

func (s *AiKnowledgeService) Search(ctx context.Context, req *aiV1.SearchAiKnowledgeRequest) (*aiV1.SearchAiKnowledgeResponse, error) {
	return s.aiKnowledgeBaseServiceClient.Search(ctx, req)
}

func (s *AiKnowledgeService) ReindexDocs(ctx context.Context, req *aiV1.ReindexDocsRequest) (*emptypb.Empty, error) {
	return s.aiKnowledgeBaseServiceClient.ReindexDocs(ctx, req)
}
