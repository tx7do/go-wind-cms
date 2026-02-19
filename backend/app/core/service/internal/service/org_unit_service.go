package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/go-utils/aggregator"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	"go-wind-cms/app/core/service/internal/data"

	identityV1 "go-wind-cms/api/gen/go/identity/service/v1"
)

type OrgUnitService struct {
	identityV1.UnimplementedOrgUnitServiceServer

	log *log.Helper

	orgUnitRepo *data.OrgUnitRepo
	userRepo    data.UserRepo
}

func NewOrgUnitService(
	ctx *bootstrap.Context,
	organizationRepo *data.OrgUnitRepo,
	userRepo data.UserRepo,
) *OrgUnitService {
	return &OrgUnitService{
		log:         ctx.NewLoggerHelper("org-unit/service/core-service"),
		orgUnitRepo: organizationRepo,
		userRepo:    userRepo,
	}
}

// extractRelationIDs 提取关联数据ID
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

// fetchRelationInfo 获取关联数据
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

	users, err := s.userRepo.ListUsersByIds(ctx, userIds)
	if err != nil {
		log.Errorf("query users err: %v", err)
		return err
	}

	for _, user := range users {
		userSet[user.GetId()] = user
	}
	return nil
}

// bindRelations 绑定关联数据
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

// enrichRelations 填充关联数据
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
	resp, err := s.orgUnitRepo.List(ctx, req)
	if err != nil {
		return nil, err
	}

	_ = s.enrichRelations(ctx, resp.Items)

	return resp, nil
}

func (s *OrgUnitService) Count(ctx context.Context, req *paginationV1.PagingRequest) (*identityV1.CountOrgUnitResponse, error) {
	count, err := s.orgUnitRepo.Count(ctx, req)
	if err != nil {
		return nil, err
	}

	return &identityV1.CountOrgUnitResponse{
		Count: uint64(count),
	}, nil
}

func (s *OrgUnitService) Get(ctx context.Context, req *identityV1.GetOrgUnitRequest) (*identityV1.OrgUnit, error) {
	resp, err := s.orgUnitRepo.Get(ctx, req)
	if err != nil {
		return nil, err
	}

	fakeItems := []*identityV1.OrgUnit{resp}
	_ = s.enrichRelations(ctx, fakeItems)

	return resp, nil
}

func (s *OrgUnitService) Create(ctx context.Context, req *identityV1.CreateOrgUnitRequest) (*emptypb.Empty, error) {
	if req.Data == nil {
		return nil, identityV1.ErrorBadRequest("invalid parameter")
	}

	if err := s.orgUnitRepo.Create(ctx, req); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *OrgUnitService) Update(ctx context.Context, req *identityV1.UpdateOrgUnitRequest) (*emptypb.Empty, error) {
	if req.Data == nil {
		return nil, identityV1.ErrorBadRequest("invalid parameter")
	}

	if err := s.orgUnitRepo.Update(ctx, req); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *OrgUnitService) Delete(ctx context.Context, req *identityV1.DeleteOrgUnitRequest) (*emptypb.Empty, error) {
	if err := s.orgUnitRepo.Delete(ctx, req); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
