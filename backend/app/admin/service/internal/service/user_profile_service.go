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

	// 自服务更新仅允许修改个人资料字段，排除租户/身份/权限/审计等敏感字段，
	// 防止自服务用户越权改租户、改用户名、改角色绑定等。
	if req.UpdateMask != nil && len(req.UpdateMask.Paths) > 0 {
		allowed := map[string]struct{}{
			"nickname": {},
			"avatar":   {},
			"realname": {},
			"mobile":   {},
			"email":    {},
			"remark":   {},
			"gender":   {},
		}
		filtered := req.UpdateMask.Paths[:0]
		for _, p := range req.UpdateMask.Paths {
			if _, ok := allowed[p]; ok {
				filtered = append(filtered, p)
			}
		}
		req.UpdateMask.Paths = filtered
	}

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

// BindContact 绑定手机号码/邮箱 — 尚未实现
func (s *UserProfileService) BindContact(_ context.Context, _ *identityV1.BindContactRequest) (*emptypb.Empty, error) {
	return nil, adminV1.ErrorNotImplemented("contact binding is not implemented")
}

// VerifyContact 验证手机号码/邮箱 — 尚未实现
func (s *UserProfileService) VerifyContact(_ context.Context, _ *identityV1.VerifyContactRequest) (*emptypb.Empty, error) {
	return nil, adminV1.ErrorNotImplemented("contact verification is not implemented")
}
