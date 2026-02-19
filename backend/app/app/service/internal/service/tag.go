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

func (s *TagService) Create(ctx context.Context, req *contentV1.CreateTagRequest) (*contentV1.Tag, error) {
	return s.tagClient.Create(ctx, req)
}

func (s *TagService) Update(ctx context.Context, req *contentV1.UpdateTagRequest) (*contentV1.Tag, error) {
	return s.tagClient.Update(ctx, req)
}

func (s *TagService) Delete(ctx context.Context, req *contentV1.DeleteTagRequest) (*emptypb.Empty, error) {
	return s.tagClient.Delete(ctx, req)
}
