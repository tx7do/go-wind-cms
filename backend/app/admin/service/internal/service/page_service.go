package service

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	contentV1 "go-wind-cms/api/gen/go/content/service/v1"

	"go-wind-cms/pkg/middleware/auth"
)

type PageService struct {
	adminV1.PageServiceHTTPServer

	pageServiceClient contentV1.PageServiceClient
	collabDocStateCli contentV1.CollabDocStateServiceClient
	log               *log.Helper
}

func NewPageService(ctx *bootstrap.Context, pageServiceClient contentV1.PageServiceClient, collabDocStateCli contentV1.CollabDocStateServiceClient) *PageService {
	return &PageService{
		log:               ctx.NewLoggerHelper("page/service/admin-service"),
		pageServiceClient: pageServiceClient,
	}
}

func (s *PageService) List(ctx context.Context, req *paginationV1.PagingRequest) (*contentV1.ListPageResponse, error) {
	return s.pageServiceClient.List(ctx, req)
}

func (s *PageService) Get(ctx context.Context, req *contentV1.GetPageRequest) (*contentV1.Page, error) {
	return s.pageServiceClient.Get(ctx, req)
}

func (s *PageService) Create(ctx context.Context, req *contentV1.CreatePageRequest) (*contentV1.Page, error) {
	if req == nil || req.Data == nil {
		return nil, adminV1.ErrorBadRequest("invalid parameter")
	}

	// 获取操作人信息
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.Data.CreatedBy = trans.Ptr(operator.UserId)

	return s.pageServiceClient.Create(ctx, req)
}

func (s *PageService) Update(ctx context.Context, req *contentV1.UpdatePageRequest) (*contentV1.Page, error) {
	if req == nil || req.Data == nil {
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

	return s.pageServiceClient.Update(ctx, req)
}

func (s *PageService) Delete(ctx context.Context, req *contentV1.DeletePageRequest) (*emptypb.Empty, error) {
	empty, err := s.pageServiceClient.Delete(ctx, req)
	if err == nil && req.GetId() != 0 && s.collabDocStateCli != nil {
		// 联动清理协同房间状态（page:{id}:*，best-effort）
		if _, derr := s.collabDocStateCli.DeleteDocStatesByPrefix(ctx, &contentV1.DeleteCollabDocStatesByPrefixRequest{
			RoomPrefix: fmt.Sprintf("page:%d:", req.GetId()),
		}); derr != nil {
			s.log.Warnf("cleanup collab doc states failed (page=%d): %v", req.GetId(), derr)
		}
	}
	return empty, err
}

func (s *PageService) ListVersions(ctx context.Context, req *contentV1.ListPageVersionsRequest) (*contentV1.ListPageVersionsResponse, error) {
	return s.pageServiceClient.ListVersions(ctx, req)
}

func (s *PageService) GetVersion(ctx context.Context, req *contentV1.GetPageVersionRequest) (*contentV1.PageVersion, error) {
	return s.pageServiceClient.GetVersion(ctx, req)
}

func (s *PageService) RestoreVersion(ctx context.Context, req *contentV1.RestorePageVersionRequest) (*contentV1.Page, error) {
	return s.pageServiceClient.RestoreVersion(ctx, req)
}

func (s *PageService) SubmitReview(ctx context.Context, req *contentV1.SubmitPageReviewRequest) (*contentV1.Page, error) {
	return s.pageServiceClient.SubmitReview(ctx, req)
}

func (s *PageService) ApproveReview(ctx context.Context, req *contentV1.ReviewPageRequest) (*contentV1.Page, error) {
	return s.pageServiceClient.ApproveReview(ctx, req)
}

func (s *PageService) RejectReview(ctx context.Context, req *contentV1.RejectPageReviewRequest) (*contentV1.Page, error) {
	return s.pageServiceClient.RejectReview(ctx, req)
}
