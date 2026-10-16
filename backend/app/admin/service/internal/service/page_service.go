package service

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
	"go-wind-cms/app/admin/service/internal/collab"

	"go-wind-cms/pkg/middleware/auth"
)

type PageService struct {
	adminV1.PageServiceHTTPServer

	pageServiceClient contentV1.PageServiceClient
	collabDocStateCli contentV1.CollabDocStateServiceClient
	collabSvc         *collab.Service
	log               *log.Helper
}

func NewPageService(ctx *bootstrap.Context, pageServiceClient contentV1.PageServiceClient, collabDocStateCli contentV1.CollabDocStateServiceClient, collabSvc *collab.Service) *PageService {
	return &PageService{
		log:               ctx.NewLoggerHelper("page/service/admin-service"),
		pageServiceClient: pageServiceClient,
		collabDocStateCli: collabDocStateCli,
		collabSvc:         collabSvc,
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
	if err == nil && req.GetId() != 0 {
		// 联动清理协同：先踢活房间（evict 前 flush 最终状态），再清持久化行。
		// 预留 worker flush 窗口，避免踢房的最终写与删行竞态复活。
		if s.collabSvc != nil {
			s.collabSvc.CloseRoomsByPrefix(fmt.Sprintf("page:%d:", req.GetId()))
			time.Sleep(1500 * time.Millisecond)
		}
		if s.collabDocStateCli != nil {
			// 两段删除：踢房 flush 后清一次，3s 后再兜底一次
			//（ygo 持久化 worker 的最终 flush 可能晚于第一次删除）
			cleanup := func() {
				cctx, ccancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer ccancel()
				if _, derr := s.collabDocStateCli.DeleteDocStatesByPrefix(cctx, &contentV1.DeleteCollabDocStatesByPrefixRequest{
					RoomPrefix: fmt.Sprintf("page:%d:", req.GetId()),
				}); derr != nil {
					s.log.Warnf("cleanup collab doc states failed (page=%d): %v", req.GetId(), derr)
				}
			}
			cleanup()
			time.AfterFunc(3*time.Second, cleanup)
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
