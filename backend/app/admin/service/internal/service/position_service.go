package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	identityV1 "go-wind-cms/api/gen/go/identity/service/v1"

	"go-wind-cms/pkg/middleware/auth"
)

type PositionService struct {
	adminV1.PositionServiceHTTPServer

	log *log.Helper

	positionServiceClient identityV1.PositionServiceClient
	orgUnitServiceClient  identityV1.OrgUnitServiceClient
}

func NewPositionService(
	ctx *bootstrap.Context,
	positionServiceClient identityV1.PositionServiceClient,
	orgUnitServiceClient identityV1.OrgUnitServiceClient,
) *PositionService {
	return &PositionService{
		log:                   ctx.NewLoggerHelper("position/service/admin-service"),
		positionServiceClient: positionServiceClient,
		orgUnitServiceClient:  orgUnitServiceClient,
	}
}

func (s *PositionService) List(ctx context.Context, req *paginationV1.PagingRequest) (*identityV1.ListPositionResponse, error) {
	resp, err := s.positionServiceClient.List(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *PositionService) Count(ctx context.Context, req *paginationV1.PagingRequest) (*identityV1.CountPositionResponse, error) {
	return s.positionServiceClient.Count(ctx, req)
}

func (s *PositionService) Get(ctx context.Context, req *identityV1.GetPositionRequest) (*identityV1.Position, error) {
	resp, err := s.positionServiceClient.Get(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *PositionService) Create(ctx context.Context, req *identityV1.CreatePositionRequest) (*emptypb.Empty, error) {
	if req.Data == nil {
		return nil, adminV1.ErrorBadRequest("invalid parameter")
	}

	// 获取操作人信息
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.Data.CreatedBy = trans.Ptr(operator.UserId)

	return s.positionServiceClient.Create(ctx, req)
}

func (s *PositionService) Update(ctx context.Context, req *identityV1.UpdatePositionRequest) (*emptypb.Empty, error) {
	if req.Data == nil {
		return nil, adminV1.ErrorBadRequest("invalid parameter")
	}

	// 获取操作人信息
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.Data.UpdatedBy = trans.Ptr(operator.UserId)
	if req.UpdateMask != nil {
		req.UpdateMask.Paths = append(req.UpdateMask.Paths, "updated_by")
	}

	return s.positionServiceClient.Update(ctx, req)
}

func (s *PositionService) Delete(ctx context.Context, req *identityV1.DeletePositionRequest) (*emptypb.Empty, error) {
	return s.positionServiceClient.Delete(ctx, req)
}
