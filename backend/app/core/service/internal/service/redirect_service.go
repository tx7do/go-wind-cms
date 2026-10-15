package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	"go-wind-cms/app/core/service/internal/data"
	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
)

// RedirectService URL 重定向管理。
type RedirectService struct {
	contentV1.UnimplementedRedirectServiceServer

	redirectRepo *data.RedirectRepo
	log          *log.Helper
}

func NewRedirectService(ctx *bootstrap.Context, redirectRepo *data.RedirectRepo) *RedirectService {
	return &RedirectService{
		log:          ctx.NewLoggerHelper("redirect/service/core-service"),
		redirectRepo: redirectRepo,
	}
}

func (s *RedirectService) List(ctx context.Context, req *paginationV1.PagingRequest) (*contentV1.ListRedirectsResponse, error) {
	items, total, err := s.redirectRepo.List(ctx, req)
	if err != nil {
		return nil, errors.New(500, "REDIRECT_LIST_FAILED", "list redirects failed")
	}
	return &contentV1.ListRedirectsResponse{Items: items, Total: uint64(total)}, nil
}

func (s *RedirectService) Get(ctx context.Context, req *contentV1.GetRedirectRequest) (*contentV1.Redirect, error) {
	if req == nil || req.GetId() == 0 {
		return nil, errors.New(400, "INVALID_PARAMETER", "invalid parameter")
	}
	return s.redirectRepo.Get(ctx, req.GetId())
}

func (s *RedirectService) Create(ctx context.Context, req *contentV1.CreateRedirectRequest) (*contentV1.Redirect, error) {
	if req == nil || req.Data == nil || req.Data.GetFromPath() == "" || req.Data.GetToUrl() == "" {
		return nil, errors.New(400, "INVALID_PARAMETER", "from_path and to_url are required")
	}
	return s.redirectRepo.Create(ctx, req.Data)
}

func (s *RedirectService) Update(ctx context.Context, req *contentV1.UpdateRedirectRequest) (*contentV1.Redirect, error) {
	if req == nil || req.Data == nil || req.GetId() == 0 {
		return nil, errors.New(400, "INVALID_PARAMETER", "invalid parameter")
	}
	return s.redirectRepo.Update(ctx, req.GetId(), req.Data)
}

func (s *RedirectService) Delete(ctx context.Context, req *contentV1.DeleteRedirectRequest) (*emptypb.Empty, error) {
	if req == nil || req.GetId() == 0 {
		return nil, errors.New(400, "INVALID_PARAMETER", "invalid parameter")
	}
	if err := s.redirectRepo.Delete(ctx, req.GetId()); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
