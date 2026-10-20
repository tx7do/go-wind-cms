package service

import (
	"context"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/v9"
	"github.com/tx7do/go-crud/viewer"
	"github.com/tx7do/go-utils/timeutil"
	"github.com/tx7do/go-utils/geoip/geolite"
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

func normalizeLoginVerifyError(err error) error {
	switch {
	case authenticationV1.IsUserNotFound(err),
		authenticationV1.IsUserFreeze(err),
		authenticationV1.IsInvalidPassword(err):
		return authenticationV1.ErrorInvalidPassword("invalid username or password")
	default:
		return err
	}
}

type AuthenticationService struct {
	authenticationV1.UnimplementedAuthenticationServiceServer

	userRepo   data.UserRepo
	roleRepo   *data.RoleRepo
	tenantRepo *data.TenantRepo

	permissionRepo *data.PermissionRepo

	userCredentialRepo *data.UserCredentialRepo

	authenticator *data.Authenticator

	// redisClient 用于 OIDC 授权码 state 的暂存（防 CSRF，一次性消费）
	redisClient *redis.Client

	vcodeCache *data.VCodeCache
	geoClient         *geolite.Client
	loginPolicyRepo   *data.LoginPolicyRepo

	// 构造期留 nil，装配期由 wiring 在通知域就绪后 RegisterNotifier 注入
	//（忘记密码的发码邮件走这条缝）。
	notifier Notifier

	log *log.Helper
}

// RegisterNotifier 注入通知域投递缝（忘记密码发码用）。
func (s *AuthenticationService) RegisterNotifier(notifier Notifier) {
	s.notifier = notifier
}

func NewAuthenticationService(
	ctx *bootstrap.Context,
	authenticator *data.Authenticator,
	userCredentialRepo *data.UserCredentialRepo,
	userRepo data.UserRepo,
	roleRepo *data.RoleRepo,
	tenantRepo *data.TenantRepo,
	permissionRepo *data.PermissionRepo,
	redisClient *redis.Client,
	vcodeCache *data.VCodeCache,
	loginPolicyRepo *data.LoginPolicyRepo,
	geoClient *geolite.Client,
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
		redisClient:        redisClient,
		vcodeCache:         vcodeCache,
		geoClient:          geoClient,
		loginPolicyRepo:    loginPolicyRepo,
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

	case authenticationV1.GrantType_authorization_code:
		// OIDC 授权码模式：code 由 IdP 颁发，本服务向 IdP 换取身份后映射/建号
		return s.doGrantTypeOidcAuthorizationCode(ctx, req)

	default:
		return nil, authenticationV1.ErrorInvalidGrantType("invalid grant type")
	}
}

// checkLoginPolicies 拉取租户登录策略并按当前上下文匹配（两段式：
// 密码校验前 userId=0 只查全局条目；凭证校验通过后带 userId 查全局+定向）。
// 策略读取失败 fail-open（与上游一致：登录基础设施抖动不打死登录）。
func (s *AuthenticationService) checkLoginPolicies(ctx context.Context, tenantID, userId uint32, clientIP, deviceId string) (bool, string) {
	policies, err := s.loginPolicyRepo.ListForLogin(ctx, tenantID)
	if err != nil {
		s.log.Errorf("list login policies failed for tenant [%d]: %s", tenantID, err.Error())
		return false, ""
	}
	return data.MatchLoginPolicy(policies, userId, clientIP, deviceId, s.resolveRegionFromIP(clientIP), time.Now())
}

// resolveRegionFromIP 把客户端 IP 解析为归属地串（取省，空=解析失败）。
// REGION 策略值与此串精确相等（大小写不敏感）即命中。
// geolite 客户端内嵌库，失败不阻断登录（REGION 维度按未判定处理）。
func (s *AuthenticationService) resolveRegionFromIP(clientIP string) string {
	if clientIP == "" || s.geoClient == nil {
		return ""
	}
	res, err := s.geoClient.Query(clientIP)
	if err != nil {
		return ""
	}
	return res.Province
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
func (s *AuthenticationService) authorizeAndEnrichUserTokenPayloadUserTenantRelationOneToOne(ctx context.Context, userID, tenantID uint32, clientType authenticationV1.ClientType, tokenPayload *authenticationV1.UserTokenPayload) error {
	// 登录所需权限按客户端类型区分:
	//   - admin(默认): 要求 sys:access_backend,保护管理后台
	//   - app:        要求 sys:access_app,C 端用户无需后台权限
	requiredPermissionCode := constants.SystemAccessBackendPermissionCode
	if clientType == authenticationV1.ClientType_app {
		requiredPermissionCode = constants.SystemAccessAppPermissionCode
	}
	hasRequiredAccess := false

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

	// 检查是否包含所需的访问权限
	if containsPermission(permissionCodes, requiredPermissionCode) {
		hasRequiredAccess = true
	}

	// 授权决策
	if !hasRequiredAccess {
		s.log.Errorf("user [%d] has no [%s] permission", userID, requiredPermissionCode)
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
func (s *AuthenticationService) authorizeAndEnrichUserTokenPayload(ctx context.Context, userID, tenantID uint32, clientType authenticationV1.ClientType, tokenPayload *authenticationV1.UserTokenPayload) error {
	switch constants.DefaultUserTenantRelationType {
	default:
		fallthrough
	case constants.UserTenantRelationOneToOne:
		return s.authorizeAndEnrichUserTokenPayloadUserTenantRelationOneToOne(ctx, userID, tenantID, clientType, tokenPayload)

	// OneToMany 与 OneToOne 的授权流程一致：userRepo.ListRoleIDsByUserID 已按关系类型
	// 分流（OneToMany 走 membershipRepo.ListMembershipRoleIDs），故复用同一实现。
	case constants.UserTenantRelationOneToMany:
		return s.authorizeAndEnrichUserTokenPayloadUserTenantRelationOneToOne(ctx, userID, tenantID, clientType, tokenPayload)
	}
}

// resolveUserAuthority 解析用户权限信息
func (s *AuthenticationService) resolveUserAuthority(ctx context.Context, user *identityV1.User, clientType authenticationV1.ClientType, tokenPayload *authenticationV1.UserTokenPayload) error {
	if user.GetStatus() != identityV1.User_NORMAL {
		s.log.Errorf("user [%d] is [%v]", user.GetId(), user.GetStatus())
		return authenticationV1.ErrorForbidden("user is disabled")
	}

	if err := s.authorizeAndEnrichUserTokenPayload(ctx, user.GetId(), user.GetTenantId(), clientType, tokenPayload); err != nil {
		return err
	}

	return nil
}

// doGrantTypePassword 处理授权类型 - 密码
func (s *AuthenticationService) doGrantTypePassword(ctx context.Context, req *authenticationV1.LoginRequest) (*authenticationV1.LoginResponse, error) {
	if req == nil {
		return nil, authenticationV1.ErrorBadRequest("invalid request")
	}
	// ===== 租户解析：tenant_code 留空视为平台（tenant 0），非空则按编号定位租户 =====
	// 解析后的 tenantID 限定后续凭证查询范围，消除同名 identifier 跨租户歧义。
	var tenantID uint32 = 0
	if code := req.GetTenantCode(); strings.TrimSpace(code) != "" {
		tenant, _ := s.tenantRepo.Get(ctx, &identityV1.GetTenantRequest{
			QueryBy: &identityV1.GetTenantRequest_Code{Code: code},
		})
		// 查不到、或租户非启用状态，统一返回同一文案，防止通过返回差异枚举有效租户编号
		if tenant == nil || tenant.GetStatus() != identityV1.Tenant_ON {
			return nil, authenticationV1.ErrorBadRequest("invalid tenant")
		}
		tenantID = tenant.GetId()
	}

	// ===== 登录策略第一段：密码校验前，按租户全局条目判定（IP/TIME/REGION 等）=====
	if blocked, reason := s.checkLoginPolicies(ctx, tenantID, 0, req.GetClientIp(), req.GetDeviceId()); blocked {
		s.log.Warnf("login blocked by policy (global) for tenant [%d]: %s", tenantID, reason)
		return nil, authenticationV1.ErrorForbidden("login blocked: %s", reason)
	}

	// ===== 凭证校验：在解析出的 tenant 范围内查单条凭证并校验密码 =====
	// 支持用户名/邮箱两种身份:登录表单的"邮箱"tab 传 email 不传 username,
	// 此时按 EMAIL 凭证查找(同一用户可同时持有用户名与邮箱两类凭证)。
	identityType := authenticationV1.UserCredential_USERNAME
	identifier := req.GetUsername()
	if strings.TrimSpace(identifier) == "" && strings.TrimSpace(req.GetEmail()) != "" {
		identityType = authenticationV1.UserCredential_EMAIL
		identifier = strings.TrimSpace(req.GetEmail())
	}

	var matchedUserID uint32
	var err error
	matchedUserID, err = s.userCredentialRepo.FindUserCredential(ctx, tenantID, identityType, identifier, req.GetPassword(), true)
	if err != nil && tenantID == 0 {
		// 未携带 tenant_code 时（C 端登录表单没有租户输入），先按平台（tenant 0）
		// 精确查找；UserNotFound 再跨租户回退一次——仅当标识符全局唯一命中时放行，
		// 并以凭证归属租户为准，消除"租户用户在 C 端永远无法登录"的问题。
		// 同名歧义仍要求调用方显式传 tenant_code。身份鉴别已通过后才更新 tenantID。
		var globalUserID, globalTenantID uint32
		if authenticationV1.IsUserNotFound(err) {
			globalUserID, globalTenantID, err = s.userCredentialRepo.FindUserCredentialAcrossTenants(ctx, identityType, identifier, req.GetPassword(), true)
		}
		if err == nil {
			matchedUserID, tenantID = globalUserID, globalTenantID
		}
	}
	if err != nil {
		// 服务端日志保留真实原因（USER_NOT_FOUND / USER_FREEZE / INVALID_PASSWORD），便于运维排查
		s.log.Errorf("verify user credential failed for identifier [%s]: %s", identifier, err.Error())

		return nil, normalizeLoginVerifyError(err)
	}

	// 获取用户信息（按凭证归属的 user_id 精确查找，避免同 identifier 多租户歧义）
	var user *identityV1.User
	user, err = s.userRepo.Get(ctx, &identityV1.GetUserRequest{
		QueryBy: &identityV1.GetUserRequest_Id{Id: matchedUserID},
	})
	if err != nil {
		s.log.Errorf("get user by id [%d] failed [%s]", matchedUserID, err.Error())
		return nil, err
	}

	// 纵深防御：凭证行的 tenant 必须与用户行的 tenant 一致，否则拒绝登录
	if user.GetTenantId() != tenantID {
		s.log.Errorf("tenant mismatch for user [%d]: credential tenant [%d] vs user tenant [%d]",
			matchedUserID, tenantID, user.GetTenantId())
		return nil, authenticationV1.ErrorBadRequest("invalid tenant")
	}

	// ===== 登录策略第二段：凭证校验通过后，带 userId 查全局+定向条目 =====
	if blocked, reason := s.checkLoginPolicies(ctx, tenantID, user.GetId(), req.GetClientIp(), req.GetDeviceId()); blocked {
		s.log.Warnf("login blocked by policy for user [%d]: %s", user.GetId(), reason)
		return nil, authenticationV1.ErrorForbidden("login blocked: %s", reason)
	}

	tokenPayload := &authenticationV1.UserTokenPayload{
		UserId:   user.GetId(),
		TenantId: user.TenantId,
		Username: user.Username,
		ClientId: req.ClientId,
		DeviceId: req.DeviceId,
		Ip:       trans.Ptr(req.GetClientIp()),
		UserAgent: trans.Ptr(req.GetUserAgent()),
	}

	// 验证权限
	if err = s.resolveUserAuthority(ctx, user, req.GetClientType(), tokenPayload); err != nil {
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
	accessToken, refreshToken, newJti, err := s.authenticator.CreateUserToken(ctx, req.GetClientType(), tokenPayload)
	if err != nil {
		return nil, err
	}

	return &authenticationV1.LoginResponse{
		TokenType:        authenticationV1.TokenType_bearer,
		AccessToken:      accessToken,
		RefreshToken:     trans.Ptr(refreshToken),
		Jti:              trans.Ptr(newJti),
		ExpiresIn:        int64(s.authenticator.GetAccessTokenExpires(req.GetClientType()).Seconds()),
		RefreshExpiresIn: trans.Ptr(int64(s.authenticator.GetRefreshTokenExpires(req.GetClientType()).Seconds())),
	}, nil
}

// doGrantTypeRefreshToken 处理授权类型 - 刷新令牌
func (s *AuthenticationService) doGrantTypeRefreshToken(ctx context.Context, req *authenticationV1.LoginRequest) (*authenticationV1.LoginResponse, error) {
	if req == nil {
		return nil, authenticationV1.ErrorBadRequest("invalid request")
	}

	// 首先验证刷新令牌——在任何数据库查询之前。
	// user_id 和 jti 必须来自已验证的令牌绑定，而非请求体中不可信的输入。
	if err := s.authenticator.VerifyRefreshToken(ctx, req.GetClientType(), req.GetUserId(), req.GetJti(), req.GetRefreshToken()); err != nil {
		s.log.Errorf("verify refresh token failed for user [%d]: [%s]", req.GetUserId(), err)
		return nil, authenticationV1.ErrorIncorrectRefreshToken("invalid refresh token")
	}

	// 令牌验证通过后才加载用户信息
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
		Ip:       trans.Ptr(req.GetClientIp()),
		UserAgent: trans.Ptr(req.GetUserAgent()),
	}

	// 刷新轮换继承首次登录时间：读旧会话元数据的 loginAt，轮换后回写到新元数据
	var inheritedLoginAt int64
	if oldMeta, mErr := s.authenticator.GetSessionMeta(ctx, req.GetClientType(), req.GetUserId(), req.GetJti()); mErr == nil && oldMeta != nil {
		inheritedLoginAt = oldMeta.LoginAt
	}

	// 解析用户权限信息
	err = s.resolveUserAuthority(ctx, user, req.GetClientType(), tokenPayload)
	if err != nil {
		s.log.Errorf("resolve user [%d] authority failed [%s]", user.GetId(), err.Error())
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
	accessToken, refreshToken, newJti, err := s.authenticator.CreateUserToken(ctx, req.GetClientType(), tokenPayload)
	if err != nil {
		return nil, authenticationV1.ErrorServiceUnavailable("generate token failed")
	}

	// 刷新轮换继承首次登录时间：新元数据的 loginAt 回写为旧会话的登录时间
	if inheritedLoginAt > 0 {
		_ = s.authenticator.UpdateSessionLoginAt(ctx, req.GetClientType(), req.GetUserId(), newJti, inheritedLoginAt)
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
func (s *AuthenticationService) RegisterUser(ctx context.Context, req *authenticationV1.RegisterUserRequest) (resp *authenticationV1.RegisterUserResponse, err error) {
	if req == nil {
		return nil, authenticationV1.ErrorBadRequest("invalid request")
	}

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

	// 用户和凭证创建必须在同一事务中，避免凭证创建失败时产生孤立用户行
	tx, err := s.userRepo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { s.userRepo.FinishTx(tx, err) }()

	// 查找租户的普通用户角色，分配给新注册用户（否则用户无角色将无法登录）。
	// 注意不能授租户管理员角色:注册用户只应有 C 端访问能力。
	var roleID *uint32
	if tenantId != nil {
		tenantRole, roleErr := s.roleRepo.GetTenantRoleByCode(ctx, *tenantId, constants.TenantUserRoleCode)
		if roleErr != nil {
			s.log.Errorf("get tenant user role for registration failed: %v", roleErr)
		} else if tenantRole != nil && tenantRole.Id != nil {
			roleID = tenantRole.Id
		}
	}

	user, err := s.userRepo.CreateWithTx(ctx, tx, &identityV1.User{
		TenantId: tenantId,
		Username: trans.Ptr(req.Username),
		Email:    req.Email,
		Status:   trans.Ptr(identityV1.User_NORMAL),
		RoleId:   roleID,
	})
	if err != nil {
		s.log.Errorf("create user error: %v", err)
		return nil, err
	}

	if err = s.userCredentialRepo.CreateWithTx(ctx, tx, &authenticationV1.UserCredential{
		UserId:   user.Id,
		TenantId: user.TenantId,

		IdentityType: authenticationV1.UserCredential_USERNAME.Enum(),
		Identifier:   trans.Ptr(req.GetUsername()),

		CredentialType: authenticationV1.UserCredential_PASSWORD_HASH.Enum(),
		Credential:     trans.Ptr(req.GetPassword()),

		IsPrimary: trans.Ptr(true),
		Status:    authenticationV1.UserCredential_ENABLED.Enum(),
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
		blockedUntil = time.Now().Add(req.GetDuration().AsDuration())
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
	if req == nil {
		return nil, authenticationV1.ErrorBadRequest("invalid request")
	}
	if err := s.authenticator.RevokeTokenByJti(ctx, req.ClientType, req.GetUserId(), req.GetJti()); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// WhoAmI 返回当前已认证调用者的身份信息。
// 身份从 viewer 上下文中解析，不接受客户端传入的任何标识，避免越权查询他人身份。
func (s *AuthenticationService) WhoAmI(ctx context.Context, _ *emptypb.Empty) (*authenticationV1.WhoAmIResponse, error) {
	userID, hasUser := viewerUserIDFromContext(ctx)
	if !hasUser {
		return nil, authenticationV1.ErrorUnauthorized("missing authentication context")
	}

	// 通过 viewer 上下文中的 user_id 查询当前用户，取其用户名
	user, err := s.userRepo.Get(ctx, &identityV1.GetUserRequest{
		QueryBy: &identityV1.GetUserRequest_Id{
			Id: userID,
		},
	})
	if err != nil {
		return nil, err
	}

	return &authenticationV1.WhoAmIResponse{
		UserId:   user.GetId(),
		Username: user.GetUsername(),
	}, nil
}
