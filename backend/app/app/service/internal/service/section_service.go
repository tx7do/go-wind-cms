package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	appV1 "go-wind-cms/api/gen/go/app/service/v1"
	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
)

type SectionService struct {
	appV1.SectionServiceHTTPServer

	sectionServiceClient contentV1.SectionServiceClient
	log                  *log.Helper
}

func NewSectionService(ctx *bootstrap.Context, sectionServiceClient contentV1.SectionServiceClient) *SectionService {
	return &SectionService{
		log:                  ctx.NewLoggerHelper("section/service/app-service"),
		sectionServiceClient: sectionServiceClient,
	}
}

func (s *SectionService) List(ctx context.Context, req *paginationV1.PagingRequest) (*contentV1.ListSectionResponse, error) {
	return s.sectionServiceClient.List(ctx, req)
}

func (s *SectionService) Get(ctx context.Context, req *contentV1.GetSectionRequest) (*contentV1.Section, error) {
	return s.sectionServiceClient.Get(ctx, req)
}

// Create/Update/Delete 在 app（公开站点）服务上禁用：CMS 内容的写操作应经由 admin 服务。
func (s *SectionService) Create(_ context.Context, _ *contentV1.CreateSectionRequest) (*contentV1.Section, error) {
	return nil, contentV1.ErrorForbidden("content mutation is not allowed on the public app service")
}

func (s *SectionService) Update(_ context.Context, _ *contentV1.UpdateSectionRequest) (*contentV1.Section, error) {
	return nil, contentV1.ErrorForbidden("content mutation is not allowed on the public app service")
}

func (s *SectionService) Delete(_ context.Context, _ *contentV1.DeleteSectionRequest) (*emptypb.Empty, error) {
	return nil, contentV1.ErrorForbidden("content mutation is not allowed on the public app service")
}

func (s *SectionService) GetTranslation(ctx context.Context, req *contentV1.GetSectionRequest) (*contentV1.SectionTranslation, error) {
	return s.sectionServiceClient.GetTranslation(ctx, req)
}
