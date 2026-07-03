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

type SectionService struct {
	contentV1.UnimplementedSectionServiceServer

	sectionRepo *data.SectionRepo
	log         *log.Helper
}

func NewSectionService(ctx *bootstrap.Context, uc *data.SectionRepo) *SectionService {
	return &SectionService{
		log:         ctx.NewLoggerHelper("section/service/core-service"),
		sectionRepo: uc,
	}
}

func (s *SectionService) List(ctx context.Context, req *paginationV1.PagingRequest) (*contentV1.ListSectionResponse, error) {
	return s.sectionRepo.List(ctx, req)
}

func (s *SectionService) Get(ctx context.Context, req *contentV1.GetSectionRequest) (*contentV1.Section, error) {
	return s.sectionRepo.Get(ctx, req)
}

func (s *SectionService) Create(ctx context.Context, req *contentV1.CreateSectionRequest) (*contentV1.Section, error) {
	return s.sectionRepo.Create(ctx, req)
}

func (s *SectionService) Update(ctx context.Context, req *contentV1.UpdateSectionRequest) (*contentV1.Section, error) {
	return s.sectionRepo.Update(ctx, req)
}

func (s *SectionService) Delete(ctx context.Context, req *contentV1.DeleteSectionRequest) (*emptypb.Empty, error) {
	err := s.sectionRepo.Delete(ctx, req)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *SectionService) TranslationExists(ctx context.Context, req *contentV1.SectionTranslationExistsRequest) (*contentV1.SectionTranslationExistsResponse, error) {
	exists, err := s.sectionRepo.TranslationExists(ctx, req.GetSectionId(), req.GetLanguageCode())
	if err != nil {
		return nil, err
	}

	return &contentV1.SectionTranslationExistsResponse{
		Exists: exists,
	}, nil
}

func (s *SectionService) GetTranslation(ctx context.Context, req *contentV1.GetSectionRequest) (*contentV1.SectionTranslation, error) {
	return s.sectionRepo.GetTranslation(ctx, req)
}

func (s *SectionService) CreateTranslation(ctx context.Context, req *contentV1.CreateSectionTranslationRequest) (*contentV1.SectionTranslation, error) {
	return s.sectionRepo.CreateTranslation(ctx, req)
}

func (s *SectionService) UpdateTranslation(ctx context.Context, req *contentV1.UpdateSectionTranslationRequest) (*contentV1.SectionTranslation, error) {
	return s.sectionRepo.UpdateTranslation(ctx, req)
}

func (s *SectionService) DeleteTranslation(ctx context.Context, req *contentV1.DeleteSectionTranslationRequest) (*emptypb.Empty, error) {
	err := s.sectionRepo.DeleteTranslation(ctx, req)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
