package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
)

// SearchSynonymService 搜索同义词管理 BFF 面：代理 core SearchSynonymService。
type SearchSynonymService struct {
	adminV1.SearchSynonymServiceHTTPServer

	log *log.Helper

	searchSynonymClient contentV1.SearchSynonymServiceClient
}

func NewSearchSynonymService(
	ctx *bootstrap.Context,
	searchSynonymClient contentV1.SearchSynonymServiceClient,
) *SearchSynonymService {
	return &SearchSynonymService{
		log:                 ctx.NewLoggerHelper("search-synonym/service/admin-service"),
		searchSynonymClient: searchSynonymClient,
	}
}

func (s *SearchSynonymService) List(ctx context.Context, _ *emptypb.Empty) (*contentV1.ListSearchSynonymsResponse, error) {
	return s.searchSynonymClient.List(ctx, &emptypb.Empty{})
}

func (s *SearchSynonymService) Create(ctx context.Context, req *contentV1.CreateSearchSynonymRequest) (*contentV1.SearchSynonym, error) {
	return s.searchSynonymClient.Create(ctx, req)
}

func (s *SearchSynonymService) Delete(ctx context.Context, req *contentV1.DeleteSearchSynonymRequest) (*emptypb.Empty, error) {
	return s.searchSynonymClient.Delete(ctx, req)
}
