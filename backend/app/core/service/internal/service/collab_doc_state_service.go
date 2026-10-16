package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	"go-wind-cms/app/core/service/internal/data"
	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
)

// CollabDocStateService 协同文档状态（collab ws 引擎的持久化后端）。
// 房间 → 最新 CRDT update 的键值存取；内容事实源仍是内容表，
// 本服务只做协同会话续命存储与内容删除联动清理。
type CollabDocStateService struct {
	contentV1.UnimplementedCollabDocStateServiceServer

	collabDocStateRepo *data.CollabDocStateRepo
	log                *log.Helper
}

func NewCollabDocStateService(ctx *bootstrap.Context, collabDocStateRepo *data.CollabDocStateRepo) *CollabDocStateService {
	return &CollabDocStateService{
		log:                ctx.NewLoggerHelper("collab-doc-state/service/core-service"),
		collabDocStateRepo: collabDocStateRepo,
	}
}

func (s *CollabDocStateService) GetDocState(ctx context.Context, req *contentV1.GetCollabDocStateRequest) (*contentV1.CollabDocState, error) {
	if req == nil || req.GetRoom() == "" {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}
	return s.collabDocStateRepo.Get(ctx, req.GetRoom())
}

func (s *CollabDocStateService) StoreDocState(ctx context.Context, req *contentV1.StoreCollabDocStateRequest) (*emptypb.Empty, error) {
	if req == nil || req.GetRoom() == "" || len(req.GetDocUpdate()) == 0 {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}
	if err := s.collabDocStateRepo.Store(ctx, req.GetRoom(), req.GetDocUpdate()); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CollabDocStateService) DeleteDocStatesByPrefix(ctx context.Context, req *contentV1.DeleteCollabDocStatesByPrefixRequest) (*emptypb.Empty, error) {
	if req == nil || req.GetRoomPrefix() == "" {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}
	if err := s.collabDocStateRepo.DeleteByPrefix(ctx, req.GetRoomPrefix()); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
