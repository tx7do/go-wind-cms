package service

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/go-crud/viewer"
	"github.com/tx7do/go-utils/timeutil"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	"go-wind-cms/app/core/service/internal/data"
	"go-wind-cms/app/core/service/internal/data/ent/privacy"

	authenticationV1 "go-wind-cms/api/gen/go/authentication/service/v1"
	identityV1 "go-wind-cms/api/gen/go/identity/service/v1"
	"go-wind-cms/pkg/constants"
	"go-wind-cms/pkg/metadata"
)

type AuthenticationService struct {
	authenticationV1.UnimplementedAuthenticationServiceServer

	userRepo   data.UserRepo
	roleRepo   *data.RoleRepo
	tenantRepo *data.TenantRepo

	permissionRepo *data.PermissionRepo

	userCredentialRepo *data.UserCredentialRepo

	authenticator *data.Authenticator

	log *log.Helper
}

func NewAuthenticationService(
	ctx *bootstrap.Context,
	authenticator *data.Authenticator,
	userCredentialRepo *data.UserCredentialRepo,
	userRepo data.UserRepo,
	roleRepo *data.RoleRepo,
	tenantRepo *data.TenantRepo,
	permissionRepo *data.PermissionRepo,
) *AuthenticationService {
	l := log.NewHelper(log.With(ctx.GetLogger(), "module", "authn/service/core-service"))
	return &AuthenticationService{
		log:                l,
		userRepo:           userRepo,
		userCredentialRepo: userCredentialRepo,
		tenantRepo:         tenantRepo,
		roleRepo:           roleRepo,
		permissionRepo:     permissionRepo,
		authenticator:      authenticator,
	}
}

// Login 登录
func (s *AuthenticationService) Login(ctx context.Context, req *authenticationV1.LoginRequest) (*authenticationV1.LoginResponse, error) {
	// 没有 viewer 信息，使用空的 NoopContext
	ctx = viewer.WithContext(ctx, viewer.NewNoopContext())
	// 绕过隐私保护中间件
	ctx = privacy.DecisionContext(ctx, privacy.Allow)
	ctx, _ = metadata.NewContext(ctx, metadata.NewUserOperator(0, 0, 0, identityV1.DataScope_ALL))

	switch req.GetGrantType() {
	case authenticationV1.GrantType_password:
		return s.doGrantTypePassword(ctx, req)

	case authenticationV1.GrantType_refresh_token:
		return s.doGrantTypeRefreshToken(ctx, req)

	case authenticationV1.GrantType_client_credentials:
		return s.doGrantTypeClientCredentials(ctx, req)

	default:
		return nil, authenticationV1.ErrorInvalidGrantType("invalid grant type")
	}
}

// containsPermission 检查权限代码列表中是否包含指定权限代码
func containsPermission(perms []string, target string) bool {
	for _, p := range perms {
		if p == target {
			return true
		}
	}
	return false
}

// authorizeAndEnrichUserTokenPayloadUserTenantRelationOneToOne 一对一用户-租户关系的授权与丰富
func (s *AuthenticationService) authorizeAndEnrichUserTokenPayloadUserTenantRelationOneToOne(ctx context.Context, userID, tenantID uint32, tokenPayload *authenticationV1.UserTokenPayload) error {
	hasBackendAccess := false

	if tenantID > 0 {
		// 检查租户状态
		tenant, _ := s.tenantRepo.Get(ctx, &identityV1.GetTenantRequest{
			QueryBy: &identityV1.GetTenantRequest_Id{Id: tenantID},
		})
		if tenant == nil || tenant.GetStatus() != identityV1.Tenant_ON {
			return authenticationV1.ErrorForbidden("insufficient authority")
		}
	}

	// 获取角色 ID 列表
	roleIDs, err := s.userRepo.ListRoleIDsByUserID(ctx, userID)
	if err != nil || len(roleIDs) == 0 {
		s.log.Errorf("get roles by user [%d] failed [%v]", userID, err)
		return authenticationV1.ErrorForbidden("insufficient authority")
	}

	// 获取权限 ID 列表
	permissionIDs, err := s.roleRepo.ListPermissionIDsByRoleIDs(ctx, roleIDs)
	if err != nil || len(permissionIDs) == 0 {
		s.log.Errorf("get permissions by role ids failed [%v]", err)
		return authenticationV1.ErrorForbidden("insufficient authority")
	}

	// 获取权限代码列表
	permissionCodes, err := s.permissionRepo.ListPermissionCodesByIds(ctx, permissionIDs)
	if err != nil || len(permissionCodes) == 0 {
		s.log.Errorf("get permission codes by ids failed [%v]", err)
		return authenticationV1.ErrorForbidden("insufficient authority")
	}

	// 检查是否包含系统访问后台权限
	if containsPermission(permissionCodes, constants.SystemAccessBackendPermissionCode) {
		hasBackendAccess = true
	}

	// 授权决策
	if !hasBackendAccess {
		s.log.Errorf("user [%d] has no backend access permission", userID)
		return authenticationV1.ErrorForbidden("insufficient authority")
	}

	// 获取角色代码列表
	roleCodes, err := s.roleRepo.ListRoleCodesByIds(ctx, roleIDs)
	if err != nil || len(roleCodes) == 0 {
		s.log.Errorf("list role codes by role ids failed [%v]", err)
		return authenticationV1.ErrorForbidden("insufficient authority")
	}
	tokenPayload.Roles = roleCodes

	return nil
}

// authorizeAndEnrichUserTokenPayload 授权并丰富用户令牌载荷
func (s *AuthenticationService) authorizeAndEnrichUserTokenPayload(ctx context.Context, userID, tenantID uint32, tokenPayload *authenticationV1.UserTokenPayload) error {
	switch constants.DefaultUserTenantRelationType {
	case constants.UserTenantRelationOneToOne:
		return s.authorizeAndEnrichUserTokenPayloadUserTenantRelationOneToOne(ctx, userID, tenantID, tokenPayload)

	case constants.UserTenantRelationOneToMany:
		s.log.Errorf("user-tenant relation type one-to-many is not implemented yet")
		return authenticationV1.ErrorServiceUnavailable("user-tenant relation type one-to-many is not implemented yet")

	default:
		s.log.Errorf("unsupported user-tenant relation type: %d", constants.DefaultUserTenantRelationType)
		return authenticationV1.ErrorServiceUnavailable("unsupported user-tenant relation type")
	}
}

// resolveUserAuthority 解析用户权限信息
func (s *AuthenticationService) resolveUserAuthority(ctx context.Context, user *identityV1.User, tokenPayload *authenticationV1.UserTokenPayload) error {
	if user.GetStatus() != identityV1.User_NORMAL {
		s.log.Errorf("user [%d] is [%v]", user.GetId(), user.GetStatus())
		return authenticationV1.ErrorForbidden("user is disabled")
	}

	if err := s.authorizeAndEnrichUserTokenPayload(ctx, user.GetId(), user.GetTenantId(), tokenPayload); err != nil {
		return err
	}

	return nil
}

// doGrantTypePassword 处理授权类型 - 密码
func (s *AuthenticationService) doGrantTypePassword(ctx context.Context, req *authenticationV1.LoginRequest) (*authenticationV1.LoginResponse, error) {
	var err error
	if _, err = s.userCredentialRepo.VerifyCredential(ctx, &authenticationV1.VerifyCredentialRequest{
		IdentityType: authenticationV1.UserCredential_USERNAME,
		Identifier:   req.GetUsername(),
		Credential:   req.GetPassword(),
		NeedDecrypt:  true,
	}); err != nil {
		return nil, err
	}

	// 获取用户信息
	var user *identityV1.User
	user, err = s.userRepo.Get(ctx, &identityV1.GetUserRequest{
		QueryBy: &identityV1.GetUserRequest_Username{Username: req.GetUsername()},
	})
	if err != nil {
		return nil, err
	}

	tokenPayload := &authenticationV1.UserTokenPayload{
		UserId:   user.GetId(),
		TenantId: user.TenantId,
		Username: user.Username,
		ClientId: req.ClientId,
		DeviceId: req.DeviceId,
	}

	// 验证权限
	if err = s.resolveUserAuthority(ctx, user, tokenPayload); err != nil {
		return nil, err
	}

	roleCodes, err := s.roleRepo.ListRoleCodesByIds(ctx, user.GetRoleIds())
	if err != nil {
		s.log.Errorf("get user role codes failed [%s]", err.Error())
	}
	if roleCodes != nil {
		user.Roles = roleCodes
	}

	// 生成令牌
	accessToken, refreshToken, err := s.authenticator.CreateUserToken(ctx, req.GetClientType(), tokenPayload)
	if err != nil {
		return nil, err
	}

	return &authenticationV1.LoginResponse{
		TokenType:        authenticationV1.TokenType_bearer,
		AccessToken:      accessToken,
		RefreshToken:     trans.Ptr(refreshToken),
		ExpiresIn:        int64(s.authenticator.GetAccessTokenExpires(req.GetClientType()).Seconds()),
		RefreshExpiresIn: trans.Ptr(int64(s.authenticator.GetRefreshTokenExpires(req.GetClientType()).Seconds())),
	}, nil
}

// doGrantTypeAuthorizationCode 处理授权类型 - 刷新令牌
func (s *AuthenticationService) doGrantTypeRefreshToken(ctx context.Context, req *authenticationV1.LoginRequest) (*authenticationV1.LoginResponse, error) {
	// 获取用户信息
	user, err := s.userRepo.Get(ctx, &identityV1.GetUserRequest{
		QueryBy: &identityV1.GetUserRequest_Id{
			Id: req.GetUserId(),
		},
	})
	if err != nil {
		return nil, err
	}

	tokenPayload := &authenticationV1.UserTokenPayload{
		UserId:   user.GetId(),
		TenantId: user.TenantId,
		Username: user.Username,
		ClientId: req.ClientId,
		DeviceId: req.DeviceId,
	}

	// 解析用户权限信息
	err = s.resolveUserAuthority(ctx, user, tokenPayload)
	if err != nil {
		s.log.Errorf("resolve user [%d] authority failed [%s]", user.GetId(), err.Error())
		return nil, err
	}

	// 验证刷新令牌
	if err = s.authenticator.VerifyRefreshToken(ctx, req.GetClientType(), req.GetUserId(), req.GetJti(), req.GetRefreshToken()); err != nil {
		s.log.Errorf("verify refresh token failed for user [%d]: [%s]", req.GetUserId(), err)
		return nil, authenticationV1.ErrorIncorrectRefreshToken("invalid refresh token")
	}

	roleCodes, err := s.roleRepo.ListRoleCodesByIds(ctx, user.GetRoleIds())
	if err != nil {
		s.log.Errorf("get user role codes failed [%s]", err.Error())
	}
	if roleCodes != nil {
		user.Roles = roleCodes
	}

	// 生成令牌
	accessToken, refreshToken, err := s.authenticator.CreateUserToken(ctx, req.GetClientType(), tokenPayload)
	if err != nil {
		return nil, authenticationV1.ErrorServiceUnavailable("generate token failed")
	}

	return &authenticationV1.LoginResponse{
		TokenType:        authenticationV1.TokenType_bearer,
		AccessToken:      accessToken,
		RefreshToken:     trans.Ptr(refreshToken),
		ExpiresIn:        int64(s.authenticator.GetAccessTokenExpires(req.GetClientType()).Seconds()),
		RefreshExpiresIn: trans.Ptr(int64(s.authenticator.GetRefreshTokenExpires(req.GetClientType()).Seconds())),
	}, nil
}

// doGrantTypeClientCredentials 处理授权类型 - 客户端凭据
func (s *AuthenticationService) doGrantTypeClientCredentials(_ context.Context, _ *authenticationV1.LoginRequest) (*authenticationV1.LoginResponse, error) {
	return nil, authenticationV1.ErrorInvalidGrantType("invalid grant type")
}

// Logout 登出
func (s *AuthenticationService) Logout(ctx context.Context, req *authenticationV1.LogoutRequest) (*emptypb.Empty, error) {
	if err := s.authenticator.RevokeUserToken(ctx, req.GetClientType(), req.GetUserId()); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// RegisterUser 注册用户
func (s *AuthenticationService) RegisterUser(ctx context.Context, req *authenticationV1.RegisterUserRequest) (*authenticationV1.RegisterUserResponse, error) {
	var err error

	var tenantId *uint32
	if constants.IsTenantModeEnabled {
		var tenant *identityV1.Tenant
		tenant, err = s.tenantRepo.Get(ctx, &identityV1.GetTenantRequest{
			QueryBy: &identityV1.GetTenantRequest_Code{Code: req.GetTenantCode()},
		})
		if err != nil {
			s.log.Errorf("get tenant by code [%s] failed: %v", req.GetTenantCode(), err)
			return nil, authenticationV1.ErrorServiceUnavailable("failed to get tenant information")
		}

		if tenant != nil {
			tenantId = tenant.Id
		}
	}

	user, err := s.userRepo.Create(ctx, &identityV1.CreateUserRequest{
		Data: &identityV1.User{
			TenantId: tenantId,
			Username: trans.Ptr(req.Username),
			Email:    req.Email,
			Status:   trans.Ptr(identityV1.User_NORMAL),
		},
	})
	if err != nil {
		s.log.Errorf("create user error: %v", err)
		return nil, err
	}

	if err = s.userCredentialRepo.Create(ctx, &authenticationV1.CreateUserCredentialRequest{
		Data: &authenticationV1.UserCredential{
			UserId:   user.Id,
			TenantId: user.TenantId,

			IdentityType: authenticationV1.UserCredential_USERNAME.Enum(),
			Identifier:   trans.Ptr(req.GetUsername()),

			CredentialType: authenticationV1.UserCredential_PASSWORD_HASH.Enum(),
			Credential:     trans.Ptr(req.GetPassword()),

			IsPrimary: trans.Ptr(true),
			Status:    authenticationV1.UserCredential_ENABLED.Enum(),
		},
	}); err != nil {
		s.log.Errorf("create user credentials error: %v", err)
		return nil, err
	}

	return &authenticationV1.RegisterUserResponse{
		UserId: user.GetId(),
	}, nil
}

// RefreshToken 刷新令牌
func (s *AuthenticationService) RefreshToken(ctx context.Context, req *authenticationV1.LoginRequest) (*authenticationV1.LoginResponse, error) {
	// 校验授权类型
	if req.GetGrantType() != authenticationV1.GrantType_refresh_token {
		return nil, authenticationV1.ErrorInvalidGrantType("invalid grant type")
	}

	return s.doGrantTypeRefreshToken(ctx, req)
}

// ValidateToken 验证令牌
func (s *AuthenticationService) ValidateToken(ctx context.Context, req *authenticationV1.ValidateTokenRequest) (*authenticationV1.ValidateTokenResponse, error) {
	return s.authenticator.Authenticate(ctx, req)
}

func (s *AuthenticationService) GetAccessTokens(ctx context.Context, req *authenticationV1.GetAccessTokensRequest) (*authenticationV1.GetAccessTokensResponse, error) {
	accessTokens := s.authenticator.GetAccessTokens(ctx, req.GetClientType(), req.GetUserId())
	return &authenticationV1.GetAccessTokensResponse{
		AccessTokens: accessTokens,
	}, nil
}

func (s *AuthenticationService) BlockToken(ctx context.Context, req *authenticationV1.BlockTokenRequest) (*authenticationV1.BlockTokenResponse, error) {
	if err := s.authenticator.BlockToken(ctx, req); err != nil {
		return nil, err
	}

	var blockedUntil time.Time
	if req.GetDuration().Seconds > 0 {
		blockedUntil = time.Now().Add(req.GetDuration().AsDuration() * time.Second)
	}

	return &authenticationV1.BlockTokenResponse{
		BlockedUntil: timeutil.TimeToTimestamppb(trans.Ptr(blockedUntil)),
	}, nil
}

func (s *AuthenticationService) UnblockToken(ctx context.Context, req *authenticationV1.UnblockTokenRequest) (*emptypb.Empty, error) {
	if err := s.authenticator.UnblockToken(ctx, req); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *AuthenticationService) RevokeTokenById(ctx context.Context, req *authenticationV1.RevokeTokenByIdRequest) (*emptypb.Empty, error) {
	if err := s.authenticator.RevokeTokenByJti(ctx, req.ClientType, req.GetUserId(), req.GetJti()); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
