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

type OrgUnitService struct {
	adminV1.OrgUnitServiceHTTPServer

	log *log.Helper

	orgUnitServiceClient identityV1.OrgUnitServiceClient
	userServiceClient    identityV1.UserServiceClient
}

func NewOrgUnitService(
	ctx *bootstrap.Context,
	orgUnitServiceClient identityV1.OrgUnitServiceClient,
	userServiceClient identityV1.UserServiceClient,
) *OrgUnitService {
	return &OrgUnitService{
		log:                  ctx.NewLoggerHelper("org-unit/service/admin-service"),
		orgUnitServiceClient: orgUnitServiceClient,
		userServiceClient:    userServiceClient,
	}
}

func (s *OrgUnitService) List(ctx context.Context, req *paginationV1.PagingRequest) (*identityV1.ListOrgUnitResponse, error) {
	resp, err := s.orgUnitServiceClient.List(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *OrgUnitService) Count(ctx context.Context, req *paginationV1.PagingRequest) (*identityV1.CountOrgUnitResponse, error) {
	return s.orgUnitServiceClient.Count(ctx, req)
}

func (s *OrgUnitService) Get(ctx context.Context, req *identityV1.GetOrgUnitRequest) (*identityV1.OrgUnit, error) {
	resp, err := s.orgUnitServiceClient.Get(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *OrgUnitService) Create(ctx context.Context, req *identityV1.CreateOrgUnitRequest) (*emptypb.Empty, error) {
	if req.Data == nil {
		return nil, adminV1.ErrorBadRequest("invalid parameter")
	}

	// 获取操作人信息
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.Data.CreatedBy = trans.Ptr(operator.UserId)

	return s.orgUnitServiceClient.Create(ctx, req)
}

func (s *OrgUnitService) Update(ctx context.Context, req *identityV1.UpdateOrgUnitRequest) (*emptypb.Empty, error) {
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

	return s.orgUnitServiceClient.Update(ctx, req)
}

func (s *OrgUnitService) Delete(ctx context.Context, req *identityV1.DeleteOrgUnitRequest) (*emptypb.Empty, error) {
	return s.orgUnitServiceClient.Delete(ctx, req)
}
