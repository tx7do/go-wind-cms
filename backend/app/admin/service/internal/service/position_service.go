package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/go-utils/aggregator"
	"github.com/tx7do/go-utils/sliceutil"
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

func (s *PositionService) extractRelationIDs(
	positions []*identityV1.Position,
	orgUnitSet aggregator.ResourceMap[uint32, *identityV1.OrgUnit],
) {
	for _, p := range positions {
		if p.GetOrgUnitId() > 0 {
			orgUnitSet[p.GetOrgUnitId()] = nil
		}
	}
}

func (s *PositionService) fetchRelationInfo(
	ctx context.Context,
	orgUnitSet aggregator.ResourceMap[uint32, *identityV1.OrgUnit],
) error {
	if len(orgUnitSet) > 0 {
		orgUnitIds := make([]uint32, 0, len(orgUnitSet))
		for id := range orgUnitSet {
			orgUnitIds = append(orgUnitIds, id)
		}

		orgUnits, err := s.orgUnitServiceClient.List(ctx, &paginationV1.PagingRequest{
			NoPaging: trans.Ptr(true),
			FilteringType: &paginationV1.PagingRequest_Query{
				Query: fmt.Sprintf(`{"id__in": "[%s]"}`, strings.Join(
					sliceutil.Map(orgUnitIds, func(value uint32, _ int, _ []uint32) string {
						return strconv.FormatUint(uint64(value), 10)
					}),
					","),
				),
			},
		})
		if err != nil {
			s.log.Errorf("query orgUnits err: %v", err)
			return err
		}

		for _, orgUnit := range orgUnits.Items {
			orgUnitSet[orgUnit.GetId()] = orgUnit
		}
	}

	return nil
}

func (s *PositionService) bindRelations(
	positions []*identityV1.Position,
	orgUnitSet aggregator.ResourceMap[uint32, *identityV1.OrgUnit],
) {
	aggregator.Populate(
		positions,
		orgUnitSet,
		func(ou *identityV1.Position) uint32 { return ou.GetOrgUnitId() },
		func(ou *identityV1.Position, org *identityV1.OrgUnit) {
			ou.OrgUnitName = org.Name
		},
	)
}

func (s *PositionService) enrichRelations(ctx context.Context, positions []*identityV1.Position) error {
	var orgUnitSet = make(aggregator.ResourceMap[uint32, *identityV1.OrgUnit])
	s.extractRelationIDs(positions, orgUnitSet)
	if err := s.fetchRelationInfo(ctx, orgUnitSet); err != nil {
		return err
	}
	s.bindRelations(positions, orgUnitSet)
	return nil
}

func (s *PositionService) List(ctx context.Context, req *paginationV1.PagingRequest) (*identityV1.ListPositionResponse, error) {
	resp, err := s.positionServiceClient.List(ctx, req)
	if err != nil {
		return nil, err
	}

	_ = s.enrichRelations(ctx, resp.Items)

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

	fakeItems := []*identityV1.Position{resp}
	_ = s.enrichRelations(ctx, fakeItems)

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
