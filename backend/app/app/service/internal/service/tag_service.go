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

type TagService struct {
	appV1.TagServiceHTTPServer

	tagClient contentV1.TagServiceClient
	log       *log.Helper
}

func NewTagService(ctx *bootstrap.Context, tagClient contentV1.TagServiceClient) *TagService {
	return &TagService{
		log:       ctx.NewLoggerHelper("tag/service/app-service"),
		tagClient: tagClient,
	}
}

func (s *TagService) List(ctx context.Context, req *paginationV1.PagingRequest) (*contentV1.ListTagResponse, error) {
	return s.tagClient.List(ctx, req)
}

func (s *TagService) Get(ctx context.Context, req *contentV1.GetTagRequest) (*contentV1.Tag, error) {
	return s.tagClient.Get(ctx, req)
}

// Create/Update/Delete 在 app（公开站点）服务上禁用：CMS 内容的写操作应经由 admin 服务。
func (s *TagService) Create(_ context.Context, _ *contentV1.CreateTagRequest) (*contentV1.Tag, error) {
	return nil, contentV1.ErrorForbidden("content mutation is not allowed on the public app service")
}

func (s *TagService) Update(_ context.Context, _ *contentV1.UpdateTagRequest) (*contentV1.Tag, error) {
	return nil, contentV1.ErrorForbidden("content mutation is not allowed on the public app service")
}

func (s *TagService) Delete(_ context.Context, _ *contentV1.DeleteTagRequest) (*emptypb.Empty, error) {
	return nil, contentV1.ErrorForbidden("content mutation is not allowed on the public app service")
}

func (s *TagService) GetTranslation(ctx context.Context, req *contentV1.GetTagRequest) (*contentV1.TagTranslation, error) {
	return s.tagClient.GetTranslation(ctx, req)
}
