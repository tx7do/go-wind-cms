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

func (s *SectionService) Create(ctx context.Context, req *contentV1.CreateSectionRequest) (*contentV1.Section, error) {
	return s.sectionServiceClient.Create(ctx, req)
}

func (s *SectionService) Update(ctx context.Context, req *contentV1.UpdateSectionRequest) (*contentV1.Section, error) {
	return s.sectionServiceClient.Update(ctx, req)
}

func (s *SectionService) Delete(ctx context.Context, req *contentV1.DeleteSectionRequest) (*emptypb.Empty, error) {
	return s.sectionServiceClient.Delete(ctx, req)
}

func (s *SectionService) GetTranslation(ctx context.Context, req *contentV1.GetSectionRequest) (*contentV1.SectionTranslation, error) {
	return s.sectionServiceClient.GetTranslation(ctx, req)
}
