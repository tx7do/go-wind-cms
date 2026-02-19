package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"google.golang.org/protobuf/types/known/emptypb"

	"go-wind-cms/app/core/service/internal/data"

	resourceV1 "go-wind-cms/api/gen/go/resource/service/v1"
)

type ApiService struct {
	resourceV1.UnimplementedApiServiceServer

	log *log.Helper

	apiRepo *data.ApiRepo
}

func NewApiService(
	ctx *bootstrap.Context,
	repo *data.ApiRepo,
) *ApiService {
	svc := &ApiService{
		log:     ctx.NewLoggerHelper("api/service/core-service"),
		apiRepo: repo,
	}

	svc.init()

	return svc
}

func (s *ApiService) init() {
}

func (s *ApiService) List(ctx context.Context, req *paginationV1.PagingRequest) (*resourceV1.ListApiResponse, error) {
	return s.apiRepo.List(ctx, req)
}

func (s *ApiService) Get(ctx context.Context, req *resourceV1.GetApiRequest) (*resourceV1.Api, error) {
	return s.apiRepo.Get(ctx, req)
}

func (s *ApiService) Create(ctx context.Context, req *resourceV1.CreateApiRequest) (*emptypb.Empty, error) {
	if req.Data == nil {
		return nil, resourceV1.ErrorBadRequest("invalid parameter")
	}

	if err := s.apiRepo.Create(ctx, req); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *ApiService) Update(ctx context.Context, req *resourceV1.UpdateApiRequest) (*emptypb.Empty, error) {
	if req.Data == nil {
		return nil, resourceV1.ErrorBadRequest("invalid parameter")
	}

	if err := s.apiRepo.Update(ctx, req); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *ApiService) Delete(ctx context.Context, req *resourceV1.DeleteApiRequest) (*emptypb.Empty, error) {
	if err := s.apiRepo.Delete(ctx, req); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
