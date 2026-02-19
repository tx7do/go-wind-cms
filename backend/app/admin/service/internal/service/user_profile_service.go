package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	authenticationV1 "go-wind-cms/api/gen/go/authentication/service/v1"
	identityV1 "go-wind-cms/api/gen/go/identity/service/v1"
	permissionV1 "go-wind-cms/api/gen/go/permission/service/v1"

	"go-wind-cms/pkg/middleware/auth"
)

type UserProfileService struct {
	adminV1.UserProfileServiceHTTPServer

	log *log.Helper

	userServiceClient     identityV1.UserServiceClient
	tenantServiceClient   identityV1.TenantServiceClient
	orgUnitServiceClient  identityV1.OrgUnitServiceClient
	positionServiceClient identityV1.PositionServiceClient

	roleServiceClient permissionV1.RoleServiceClient

	userCredentialServiceClient authenticationV1.UserCredentialServiceClient
}

func NewUserProfileService(
	ctx *bootstrap.Context,
	userServiceClient identityV1.UserServiceClient,
	tenantServiceClient identityV1.TenantServiceClient,
	orgUnitServiceClient identityV1.OrgUnitServiceClient,
	positionServiceClient identityV1.PositionServiceClient,
	roleServiceClient permissionV1.RoleServiceClient,
	userCredentialServiceClient authenticationV1.UserCredentialServiceClient,
) *UserProfileService {
	return &UserProfileService{
		log:                         ctx.NewLoggerHelper("user-profile/service/admin-service"),
		userServiceClient:           userServiceClient,
		tenantServiceClient:         tenantServiceClient,
		orgUnitServiceClient:        orgUnitServiceClient,
		positionServiceClient:       positionServiceClient,
		roleServiceClient:           roleServiceClient,
		userCredentialServiceClient: userCredentialServiceClient,
	}
}

func (s *UserProfileService) GetUser(ctx context.Context, _ *emptypb.Empty) (*identityV1.User, error) {
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := s.userServiceClient.Get(ctx, &identityV1.GetUserRequest{
		QueryBy: &identityV1.GetUserRequest_Id{
			Id: operator.UserId,
		},
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *UserProfileService) UpdateUser(ctx context.Context, req *identityV1.UpdateUserRequest) (*emptypb.Empty, error) {
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.Data.Id = trans.Ptr(operator.UserId)
	req.Id = operator.UserId

	return s.userServiceClient.Update(ctx, req)
}

func (s *UserProfileService) ChangePassword(ctx context.Context, req *identityV1.ChangePasswordRequest) (*emptypb.Empty, error) {
	// 获取操作人信息
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	return s.userCredentialServiceClient.ChangeCredential(ctx, &authenticationV1.ChangeCredentialRequest{
		IdentityType:  authenticationV1.UserCredential_USERNAME,
		Identifier:    operator.GetUsername(),
		OldCredential: req.GetOldPassword(),
		NewCredential: req.GetNewPassword(),
	})
}
