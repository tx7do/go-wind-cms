package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	"go-wind-cms/app/core/service/internal/data"

	dictV1 "go-wind-cms/api/gen/go/dict/service/v1"
)

type DictEntryService struct {
	dictV1.UnimplementedDictEntryServiceServer

	log *log.Helper

	dictEntryRepo *data.DictEntryRepo
}

func NewDictEntryService(
	ctx *bootstrap.Context,
	dictEntryRepo *data.DictEntryRepo,
) *DictEntryService {
	return &DictEntryService{
		log:           ctx.NewLoggerHelper("dict-entry/service/core-service"),
		dictEntryRepo: dictEntryRepo,
	}
}

func (s *DictEntryService) List(ctx context.Context, req *paginationV1.PagingRequest) (*dictV1.ListDictEntryResponse, error) {
	return s.dictEntryRepo.List(ctx, req)
}

func (s *DictEntryService) Count(ctx context.Context, req *paginationV1.PagingRequest) (*dictV1.CountDictEntryResponse, error) {
	count, err := s.dictEntryRepo.Count(ctx, req)
	if err != nil {
		return nil, err
	}

	return &dictV1.CountDictEntryResponse{
		Count: uint64(count),
	}, nil
}

func (s *DictEntryService) Create(ctx context.Context, req *dictV1.CreateDictEntryRequest) (*emptypb.Empty, error) {
	if req == nil || req.Data == nil {
		return nil, dictV1.ErrorBadRequest("invalid parameter")
	}

	if err := s.dictEntryRepo.Create(ctx, req); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *DictEntryService) Update(ctx context.Context, req *dictV1.UpdateDictEntryRequest) (*emptypb.Empty, error) {
	if req == nil || req.Data == nil {
		return nil, dictV1.ErrorBadRequest("invalid parameter")
	}

	if err := s.dictEntryRepo.Update(ctx, req); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *DictEntryService) Delete(ctx context.Context, req *dictV1.DeleteDictEntryRequest) (*emptypb.Empty, error) {
	if err := s.dictEntryRepo.BatchDelete(ctx, req.GetIds()); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *DictEntryService) ListByTypeCode(ctx context.Context, req *dictV1.ListDictEntryByTypeCodeRequest) (*dictV1.ListDictEntryByTypeCodeResponse, error) {
	if req.GetTypeCode() == "" {
		return nil, dictV1.ErrorBadRequest("invalid parameter")
	}

	return s.dictEntryRepo.ListByTypeCode(ctx, req)
}
