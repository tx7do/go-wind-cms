package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	dictV1 "go-wind-cms/api/gen/go/dict/service/v1"

	"go-wind-cms/pkg/middleware/auth"
)

type DictTypeService struct {
	adminV1.DictTypeServiceHTTPServer

	log *log.Helper

	dictTypeServiceClient dictV1.DictTypeServiceClient
}

func NewDictTypeService(
	ctx *bootstrap.Context,
	dictTypeServiceClient dictV1.DictTypeServiceClient,
) *DictTypeService {
	l := log.NewHelper(log.With(ctx.GetLogger(), "module", "dict/service/admin-service"))
	return &DictTypeService{
		log:                   l,
		dictTypeServiceClient: dictTypeServiceClient,
	}
}

func (s *DictTypeService) List(ctx context.Context, req *paginationV1.PagingRequest) (*dictV1.ListDictTypeResponse, error) {
	return s.dictTypeServiceClient.List(ctx, req)
}

func (s *DictTypeService) Get(ctx context.Context, req *dictV1.GetDictTypeRequest) (*dictV1.DictType, error) {
	return s.dictTypeServiceClient.Get(ctx, req)
}

func (s *DictTypeService) Create(ctx context.Context, req *dictV1.CreateDictTypeRequest) (*emptypb.Empty, error) {
	if req.Data == nil {
		return nil, adminV1.ErrorBadRequest("invalid parameter")
	}

	// 获取操作人信息
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.Data.CreatedBy = trans.Ptr(operator.UserId)

	if _, err = s.dictTypeServiceClient.Create(ctx, req); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *DictTypeService) Update(ctx context.Context, req *dictV1.UpdateDictTypeRequest) (*emptypb.Empty, error) {
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

	if _, err = s.dictTypeServiceClient.Update(ctx, req); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *DictTypeService) Delete(ctx context.Context, req *dictV1.DeleteDictTypeRequest) (*emptypb.Empty, error) {
	return s.dictTypeServiceClient.Delete(ctx, req)
}
