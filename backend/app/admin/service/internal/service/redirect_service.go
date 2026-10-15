package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
)

// RedirectService URL 重定向管理门面（透传 core）。
type RedirectService struct {
	adminV1.RedirectServiceHTTPServer

	redirectServiceClient contentV1.RedirectServiceClient
	log                   *log.Helper
}

func NewRedirectService(ctx *bootstrap.Context, redirectServiceClient contentV1.RedirectServiceClient) *RedirectService {
	return &RedirectService{
		log:                   ctx.NewLoggerHelper("redirect/service/admin-service"),
		redirectServiceClient: redirectServiceClient,
	}
}

func (s *RedirectService) List(ctx context.Context, req *paginationV1.PagingRequest) (*contentV1.ListRedirectsResponse, error) {
	return s.redirectServiceClient.List(ctx, req)
}

func (s *RedirectService) Get(ctx context.Context, req *contentV1.GetRedirectRequest) (*contentV1.Redirect, error) {
	return s.redirectServiceClient.Get(ctx, req)
}

func (s *RedirectService) Create(ctx context.Context, req *contentV1.CreateRedirectRequest) (*contentV1.Redirect, error) {
	return s.redirectServiceClient.Create(ctx, req)
}

func (s *RedirectService) Update(ctx context.Context, req *contentV1.UpdateRedirectRequest) (*contentV1.Redirect, error) {
	return s.redirectServiceClient.Update(ctx, req)
}

func (s *RedirectService) Delete(ctx context.Context, req *contentV1.DeleteRedirectRequest) (*emptypb.Empty, error) {
	return s.redirectServiceClient.Delete(ctx, req)
}
