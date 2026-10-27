package service

import (
	"context"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/tx7do/go-crud/viewer"
	contentV1 "go-wind-cms/api/gen/go/content/service/v1"

	"go-wind-cms/app/core/service/internal/data"
)

// SearchSynonymService 搜索同义词管理（租户级词库）：
// 检索时查询侧按触发词扩展（OR 组），免索引重建、更新即时生效。
type SearchSynonymService struct {
	contentV1.UnimplementedSearchSynonymServiceServer

	repo *data.SearchSynonymRepo
	log  *log.Helper
}

func NewSearchSynonymService(ctx *bootstrap.Context, repo *data.SearchSynonymRepo) *SearchSynonymService {
	return &SearchSynonymService{
		log:  ctx.NewLoggerHelper("search-synonym/service/core-service"),
		repo: repo,
	}
}

func (s *SearchSynonymService) List(ctx context.Context, _ *emptypb.Empty) (*contentV1.ListSearchSynonymsResponse, error) {
	items, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	return &contentV1.ListSearchSynonymsResponse{Items: items}, nil
}

func (s *SearchSynonymService) Create(ctx context.Context, req *contentV1.CreateSearchSynonymRequest) (*contentV1.SearchSynonym, error) {
	if req == nil || req.Data == nil {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}
	term := strings.TrimSpace(req.GetData().GetTerm())
	if term == "" {
		return nil, contentV1.ErrorBadRequest("term is required")
	}
	if len(req.GetData().GetSynonyms()) == 0 {
		return nil, contentV1.ErrorBadRequest("at least one synonym is required")
	}
	req.Data.Term = term
	if req.Data.Enabled == nil {
		req.Data.Enabled = new(bool)
		*req.Data.Enabled = true
	}
	// 租户归属按操作者 viewer 显式写入（租户隔离 mixin 无创建钩子）
	var tenantID uint32
	if vc, exist := viewer.FromContext(ctx); exist && vc != nil {
		tenantID = uint32(vc.TenantID())
	}
	return s.repo.Create(ctx, tenantID, req.Data)
}

func (s *SearchSynonymService) Delete(ctx context.Context, req *contentV1.DeleteSearchSynonymRequest) (*emptypb.Empty, error) {
	if req.GetId() == 0 {
		return nil, contentV1.ErrorBadRequest("id is required")
	}
	if err := s.repo.Delete(ctx, req.GetId()); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
