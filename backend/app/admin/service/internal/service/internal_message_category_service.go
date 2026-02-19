package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	internalMessageV1 "go-wind-cms/api/gen/go/internal_message/service/v1"

	"go-wind-cms/pkg/middleware/auth"
)

type InternalMessageCategoryService struct {
	adminV1.InternalMessageCategoryServiceHTTPServer

	log *log.Helper

	repo internalMessageV1.InternalMessageCategoryServiceClient
}

func NewInternalMessageCategoryService(
	ctx *bootstrap.Context,
	repo internalMessageV1.InternalMessageCategoryServiceClient,
) *InternalMessageCategoryService {
	l := log.NewHelper(log.With(ctx.GetLogger(), "module", "internal-message-category/service/admin-service"))
	return &InternalMessageCategoryService{
		log:  l,
		repo: repo,
	}
}

func (s *InternalMessageCategoryService) List(ctx context.Context, req *paginationV1.PagingRequest) (*internalMessageV1.ListInternalMessageCategoryResponse, error) {
	return s.repo.List(ctx, req)
}

func (s *InternalMessageCategoryService) Get(ctx context.Context, req *internalMessageV1.GetInternalMessageCategoryRequest) (*internalMessageV1.InternalMessageCategory, error) {
	return s.repo.Get(ctx, req)
}

func (s *InternalMessageCategoryService) Create(ctx context.Context, req *internalMessageV1.CreateInternalMessageCategoryRequest) (*emptypb.Empty, error) {
	if req.Data == nil {
		return nil, adminV1.ErrorBadRequest("invalid parameter")
	}

	// 获取操作人信息
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.Data.CreatedBy = trans.Ptr(operator.UserId)

	return s.repo.Create(ctx, req)
}

func (s *InternalMessageCategoryService) Update(ctx context.Context, req *internalMessageV1.UpdateInternalMessageCategoryRequest) (*emptypb.Empty, error) {
	if req.Data == nil {
		return nil, adminV1.ErrorBadRequest("invalid parameter")
	}

	// 获取操作人信息
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.Data.UpdatedBy = trans.Ptr(operator.GetUserId())
	if req.UpdateMask != nil {
		req.UpdateMask.Paths = append(req.UpdateMask.Paths, "updated_by")
	}

	return s.repo.Update(ctx, req)
}

func (s *InternalMessageCategoryService) Delete(ctx context.Context, req *internalMessageV1.DeleteInternalMessageCategoryRequest) (*emptypb.Empty, error) {
	return s.repo.Delete(ctx, req)
}
