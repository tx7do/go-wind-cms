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

type PostService struct {
	adminV1.PostServiceHTTPServer

	postServiceClient contentV1.PostServiceClient
	collabDocStateCli contentV1.CollabDocStateServiceClient
	collabSvc         *collab.Service
	log               *log.Helper
}

func NewPostService(ctx *bootstrap.Context, postServiceClient contentV1.PostServiceClient, collabDocStateCli contentV1.CollabDocStateServiceClient, collabSvc *collab.Service) *PostService {
	return &PostService{
		log:               ctx.NewLoggerHelper("post/service/admin-service"),
		postServiceClient: postServiceClient,
		collabDocStateCli: collabDocStateCli,
		collabSvc:         collabSvc,
	}
}

func (s *PostService) List(ctx context.Context, req *paginationV1.PagingRequest) (*contentV1.ListPostResponse, error) {
	return s.postServiceClient.List(ctx, req)
}

func (s *PostService) Get(ctx context.Context, req *contentV1.GetPostRequest) (*contentV1.Post, error) {
	return s.postServiceClient.Get(ctx, req)
}

func (s *PostService) Create(ctx context.Context, req *contentV1.CreatePostRequest) (*contentV1.Post, error) {
	if req == nil || req.Data == nil {
		return nil, adminV1.ErrorBadRequest("invalid parameter")
	}

	// 获取操作人信息
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.Data.CreatedBy = trans.Ptr(operator.UserId)
	req.Data.AuthorId = trans.Ptr(operator.UserId)

	for i := range req.Data.Translations {
		req.Data.Translations[i].CreatedBy = trans.Ptr(operator.UserId)
	}

	return s.postServiceClient.Create(ctx, req)
}

func (s *PostService) Update(ctx context.Context, req *contentV1.UpdatePostRequest) (*contentV1.Post, error) {
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

	for i := range req.Data.Translations {
		req.Data.Translations[i].UpdatedBy = trans.Ptr(operator.UserId)
	}

	return s.postServiceClient.Update(ctx, req)
}

func (s *PostService) Delete(ctx context.Context, req *contentV1.DeletePostRequest) (*emptypb.Empty, error) {
	empty, err := s.postServiceClient.Delete(ctx, req)
	if err == nil && req.GetId() != 0 {
		// 联动清理协同：先踢活房间（evict 前 flush 最终状态），再清持久化行。
		// 预留 worker flush 窗口，避免踢房的最终写与删行竞态复活。
		if s.collabSvc != nil {
			s.collabSvc.CloseRoomsByPrefix(fmt.Sprintf("post:%d:", req.GetId()))
			time.Sleep(1500 * time.Millisecond)
		}
		if s.collabDocStateCli != nil {
			// 两段删除：踢房 flush 后清一次，3s 后再兜底一次
			//（ygo 持久化 worker 的最终 flush 可能晚于第一次删除）
			cleanup := func() {
				cctx, ccancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer ccancel()
				if _, derr := s.collabDocStateCli.DeleteDocStatesByPrefix(cctx, &contentV1.DeleteCollabDocStatesByPrefixRequest{
					RoomPrefix: fmt.Sprintf("post:%d:", req.GetId()),
				}); derr != nil {
					s.log.Warnf("cleanup collab doc states failed (post=%d): %v", req.GetId(), derr)
				}
			}
			cleanup()
			time.AfterFunc(3*time.Second, cleanup)
		}
	}
	return empty, err
}

func (s *PostService) TranslationExists(ctx context.Context, req *contentV1.PostTranslationExistsRequest) (*contentV1.PostTranslationExistsResponse, error) {
	return s.postServiceClient.TranslationExists(ctx, req)
}

func (s *PostService) ListVersions(ctx context.Context, req *contentV1.ListPostVersionsRequest) (*contentV1.ListPostVersionsResponse, error) {
	return s.postServiceClient.ListVersions(ctx, req)
}

func (s *PostService) GetVersion(ctx context.Context, req *contentV1.GetPostVersionRequest) (*contentV1.PostVersion, error) {
	return s.postServiceClient.GetVersion(ctx, req)
}

func (s *PostService) RestoreVersion(ctx context.Context, req *contentV1.RestorePostVersionRequest) (*contentV1.Post, error) {
	return s.postServiceClient.RestoreVersion(ctx, req)
}

func (s *PostService) SubmitReview(ctx context.Context, req *contentV1.SubmitPostReviewRequest) (*contentV1.Post, error) {
	return s.postServiceClient.SubmitReview(ctx, req)
}

func (s *PostService) ApproveReview(ctx context.Context, req *contentV1.ReviewPostRequest) (*contentV1.Post, error) {
	return s.postServiceClient.ApproveReview(ctx, req)
}

func (s *PostService) RejectReview(ctx context.Context, req *contentV1.RejectPostReviewRequest) (*contentV1.Post, error) {
	return s.postServiceClient.RejectReview(ctx, req)
}
