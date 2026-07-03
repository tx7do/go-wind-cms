package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	contentV1 "go-wind-cms/api/gen/go/content/service/v1"

	"go-wind-cms/pkg/middleware/auth"
)

type SectionService struct {
	adminV1.SectionServiceHTTPServer

	sectionServiceClient contentV1.SectionServiceClient
	log                  *log.Helper
}

func NewSectionService(ctx *bootstrap.Context, sectionServiceClient contentV1.SectionServiceClient) *SectionService {
	return &SectionService{
		log:                  ctx.NewLoggerHelper("section/service/admin-service"),
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
	if req.Data == nil {
		return nil, adminV1.ErrorBadRequest("invalid parameter")
	}

	// 获取操作人信息
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.Data.CreatedBy = trans.Ptr(operator.UserId)

	return s.sectionServiceClient.Create(ctx, req)
}

func (s *SectionService) Update(ctx context.Context, req *contentV1.UpdateSectionRequest) (*contentV1.Section, error) {
	if req.Data == nil {
		return nil, adminV1.ErrorBadRequest("invalid parameter")
	}

	// 获取操作人信息
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.Data.Id = trans.Ptr(req.GetId())

	req.Data.UpdatedBy = trans.Ptr(operator.GetUserId())
	if req.UpdateMask != nil {
		req.UpdateMask.Paths = append(req.UpdateMask.Paths, "updated_by")
	}

	return s.sectionServiceClient.Update(ctx, req)
}

func (s *SectionService) Delete(ctx context.Context, req *contentV1.DeleteSectionRequest) (*emptypb.Empty, error) {
	return s.sectionServiceClient.Delete(ctx, req)
}
