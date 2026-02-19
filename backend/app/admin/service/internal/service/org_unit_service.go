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

func (s *OrgUnitService) extractRelationIDs(
	orgUnits []*identityV1.OrgUnit,
	userSet aggregator.ResourceMap[uint32, *identityV1.User],
) {
	for _, ou := range orgUnits {
		if ou == nil {
			continue
		}
		if id := ou.GetLeaderId(); id > 0 {
			userSet[id] = nil
		}
		if id := ou.GetContactUserId(); id > 0 {
			userSet[id] = nil
		}

		if len(ou.Children) > 0 {
			s.extractRelationIDs(ou.Children, userSet)
		}
	}
}

func (s *OrgUnitService) fetchRelationInfo(
	ctx context.Context,
	userSet aggregator.ResourceMap[uint32, *identityV1.User],
) error {
	if len(userSet) == 0 {
		return nil
	}

	userIds := make([]uint32, 0, len(userSet))
	for id := range userSet {
		userIds = append(userIds, id)
	}

	users, err := s.userServiceClient.List(ctx, &paginationV1.PagingRequest{
		NoPaging: trans.Ptr(true),
		FilteringType: &paginationV1.PagingRequest_Query{
			Query: fmt.Sprintf(`{"id__in": "[%s]"}`, strings.Join(
				sliceutil.Map(userIds, func(value uint32, _ int, _ []uint32) string {
					return strconv.FormatUint(uint64(value), 10)
				}),
				","),
			),
		},
	})
	if err != nil {
		log.Errorf("query users err: %v", err)
		return err
	}

	for _, user := range users.Items {
		userSet[user.GetId()] = user
	}

	return nil
}

func (s *OrgUnitService) bindRelations(
	orgUnits []*identityV1.OrgUnit,
	userSet aggregator.ResourceMap[uint32, *identityV1.User],
) {
	childrenFunc := func(ou *identityV1.OrgUnit) []*identityV1.OrgUnit { return ou.GetChildren() }

	// 回填 LeaderName
	aggregator.PopulateTree(
		orgUnits,
		userSet,
		func(ou *identityV1.OrgUnit) uint32 { return ou.GetLeaderId() },
		func(ou *identityV1.OrgUnit, user *identityV1.User) {
			ou.LeaderName = trans.Ptr(user.GetUsername())
		},
		childrenFunc,
	)

	// 回填 ContactUserName
	aggregator.PopulateTree(
		orgUnits,
		userSet,
		func(ou *identityV1.OrgUnit) uint32 { return ou.GetContactUserId() },
		func(ou *identityV1.OrgUnit, user *identityV1.User) {
			ou.ContactUserName = trans.Ptr(user.GetUsername())
		},
		childrenFunc,
	)
}

func (s *OrgUnitService) enrichRelations(ctx context.Context, orgUnits []*identityV1.OrgUnit) error {
	var userSet = make(aggregator.ResourceMap[uint32, *identityV1.User])
	s.extractRelationIDs(orgUnits, userSet)
	if err := s.fetchRelationInfo(ctx, userSet); err != nil {
		return err
	}
	s.bindRelations(orgUnits, userSet)
	return nil
}

func (s *OrgUnitService) List(ctx context.Context, req *paginationV1.PagingRequest) (*identityV1.ListOrgUnitResponse, error) {
	resp, err := s.orgUnitServiceClient.List(ctx, req)
	if err != nil {
		return nil, err
	}

	_ = s.enrichRelations(ctx, resp.Items)

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

	fakeItems := []*identityV1.OrgUnit{resp}
	_ = s.enrichRelations(ctx, fakeItems)

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
