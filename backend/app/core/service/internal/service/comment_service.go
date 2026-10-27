package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	"go-wind-cms/app/core/service/internal/data"

	commentV1 "go-wind-cms/api/gen/go/comment/service/v1"
	"go-wind-cms/pkg/task"
)

type CommentService struct {
	commentV1.UnimplementedCommentServiceServer

	commentRepo *data.CommentRepo
	flowSvc     *FlowService
	configRepo  *data.ConfigRepo
	aiDeps      *CommentAIModerationDeps
	log         *log.Helper
}

func NewCommentService(
	ctx *bootstrap.Context,
	uc *data.CommentRepo,
	flowSvc *FlowService,
	configRepo *data.ConfigRepo,
	aiDeps *CommentAIModerationDeps,
) *CommentService {
	return &CommentService{
		log:         ctx.NewLoggerHelper("comment/service/core-service"),
		commentRepo: uc,
		flowSvc:     flowSvc,
		configRepo:  configRepo,
		aiDeps:      aiDeps,
	}
}

func (s *CommentService) List(ctx context.Context, req *paginationV1.PagingRequest) (*commentV1.ListCommentResponse, error) {
	return s.commentRepo.List(ctx, req, true)
}

func (s *CommentService) Count(ctx context.Context, req *paginationV1.PagingRequest) (*commentV1.CountCommentResponse, error) {
	count, err := s.commentRepo.Count(ctx, req)
	if err != nil {
		return nil, err
	}

	return &commentV1.CountCommentResponse{
		Count: uint64(count),
	}, nil
}

func (s *CommentService) Get(ctx context.Context, req *commentV1.GetCommentRequest) (*commentV1.Comment, error) {
	return s.commentRepo.Get(ctx, req)
}

func (s *CommentService) Create(ctx context.Context, req *commentV1.CreateCommentRequest) (*commentV1.Comment, error) {
	// AI 反垃圾闸门（系统参数开关，默认关）：判为垃圾即拒绝创建。
	// AI 链路任何错误都 fail-open（放行）——可用性优先，垃圾漏网可由人工审核兜底。
	if err := s.moderateCommentAI(ctx, req.GetData().GetContent()); err != nil {
		return nil, err
	}
	dto, err := s.commentRepo.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	// Flow 规则触发（best-effort；评论创建事件，快照带状态/内容摘要）
	if s.flowSvc != nil && dto.GetId() != 0 {
		s.flowSvc.DispatchFlow(ctx, task.WebhookEventCommentCreated, dto.GetId(), CommentEventSnapshot(dto))
	}
	return dto, nil
}

func (s *CommentService) Update(ctx context.Context, req *commentV1.UpdateCommentRequest) (*commentV1.Comment, error) {
	return s.commentRepo.Update(ctx, req)
}

func (s *CommentService) Delete(ctx context.Context, req *commentV1.DeleteCommentRequest) (*emptypb.Empty, error) {
	err := s.commentRepo.Delete(ctx, req)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
