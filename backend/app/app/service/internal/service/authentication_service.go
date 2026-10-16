package service

import (
	"context"
	"encoding/base64"
	"regexp"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/tx7do/go-utils/crypto"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	appV1 "go-wind-cms/api/gen/go/app/service/v1"
	authenticationV1 "go-wind-cms/api/gen/go/authentication/service/v1"
	identityV1 "go-wind-cms/api/gen/go/identity/service/v1"
	entmiddleware "go-wind-cms/pkg/middleware/ent"

	"go-wind-cms/pkg/middleware/auth"
	"go-wind-cms/pkg/netutil"
)

type AuthenticationService struct {
	appV1.AuthenticationServiceHTTPServer

	authenticationServiceClient authenticationV1.AuthenticationServiceClient
	tenantClient                identityV1.TenantServiceClient
	tenantResolve               entmiddleware.TenantResolver

	log *log.Helper
}

func NewAuthenticationService(
	ctx *bootstrap.Context,
	authenticationServiceClient authenticationV1.AuthenticationServiceClient,
	tenantClient identityV1.TenantServiceClient,
	tenantResolve entmiddleware.TenantResolver,
) *AuthenticationService {
	return &AuthenticationService{
		log:                         ctx.NewLoggerHelper("authn/service/app-service"),
		authenticationServiceClient: authenticationServiceClient,
		tenantClient:                tenantClient,
		tenantResolve:               tenantResolve,
	}
}

// Login 登陆
// GetOidcAuthorizeUrl 获取 OIDC 授权跳转地址（透传 core；state 由 core 生成并落 Redis）。
// 免鉴权端点（登录前置步骤）。redirect_uri 由前端按自身回调页传入。
func (s *AuthenticationService) GetOidcAuthorizeUrl(ctx context.Context, req *authenticationV1.GetOidcAuthorizeUrlRequest) (*authenticationV1.GetOidcAuthorizeUrlResponse, error) {
	return s.authenticationServiceClient.GetOidcAuthorizeUrl(ctx, req)
}

func (s *AuthenticationService) Login(ctx context.Context, req *authenticationV1.LoginRequest) (*authenticationV1.LoginResponse, error) {
	if req == nil {
		return nil, authenticationV1.ErrorBadRequest("invalid request")
	}

	req.ClientType = trans.Ptr(authenticationV1.ClientType_app)

	if req.GetGrantType() == authenticationV1.GrantType_refresh_token {
		operator, err := auth.FromContext(ctx)
		if err != nil {
			return nil, err
		}

		req.Jti = operator.Jti
		req.UserId = trans.Ptr(operator.GetUserId())
	}

	return s.authenticationServiceClient.Login(ctx, req)
}

// Logout 登出
func (s *AuthenticationService) Logout(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	return s.authenticationServiceClient.Logout(ctx, &authenticationV1.LogoutRequest{
		ClientType: authenticationV1.ClientType_app,
		UserId:     operator.GetUserId(),
	})
}

// RefreshToken 刷新认证令牌
func (s *AuthenticationService) RefreshToken(ctx context.Context, req *authenticationV1.LoginRequest) (*authenticationV1.LoginResponse, error) {
	if req == nil {
		return nil, authenticationV1.ErrorBadRequest("invalid request")
	}

	// 常规路径:access token 仍有效时,auth 中间件已注入 operator。
	// 兜底路径:access token 已过期(本操作在白名单中,operator 不会被注入)时,
	// 从过期 JWT 的 payload 自行解码 uid/jti 作为刷新令牌的绑定键——真正的
	// 凭证校验由 core 以 refresh token 值完成,过期 access token 不作为凭证。
	operator, err := auth.FromContext(ctx)
	if err == nil {
		req.UserId = trans.Ptr(operator.GetUserId())
		req.Jti = operator.Jti
	} else if uid, jti, uerr := netutil.ParseUnverifiedBearerJWTClaims(ctx); uerr == nil {
		req.UserId = trans.Ptr(uid)
		req.Jti = trans.Ptr(jti)
	}

	req.ClientType = trans.Ptr(authenticationV1.ClientType_app)

	return s.authenticationServiceClient.RefreshToken(ctx, req)
}

// Register C 端注册。租户归属按请求 Host 解析(前端注册表单没有租户输入,
// 调用方传入的 tenant_code 一律丢弃),解析不出租户则拒绝;租户编号由 id 反查
// code 后传给 core(其按 code 定位租户,并把租户管理员角色分配给新用户,
// 否则无角色用户无法登录)。
func (s *AuthenticationService) Register(ctx context.Context, req *authenticationV1.RegisterUserRequest) (*authenticationV1.RegisterUserResponse, error) {
	if req == nil || strings.TrimSpace(req.GetUsername()) == "" || strings.TrimSpace(req.GetPassword()) == "" {
		return nil, authenticationV1.ErrorBadRequest("invalid username or password")
	}
	// 服务端兜底校验,口径与前端一致:用户名 3-20 位字母/数字/下划线,密码至少 6 位
	if len(req.GetUsername()) < 3 || len(req.GetUsername()) > 20 ||
		!regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(req.GetUsername()) {
		return nil, authenticationV1.ErrorBadRequest("invalid username format")
	}
	if len(req.GetPassword()) < 6 {
		return nil, authenticationV1.ErrorBadRequest("password too weak")
	}
	if s.tenantResolve == nil || s.tenantClient == nil {
		return nil, authenticationV1.ErrorServiceUnavailable("register unavailable")
	}

	var tenantID uint32
	if hr, ok := khttp.RequestFromServerContext(ctx); ok && hr != nil && hr.Host != "" {
		if tid, err := s.tenantResolve.ResolveTenantIDByDomain(ctx, hr.Host); err == nil {
			tenantID = tid
		}
	}
	if tenantID == 0 {
		return nil, authenticationV1.ErrorBadRequest("cannot resolve tenant, registration unavailable")
	}

	tenant, err := s.tenantClient.Get(ctx, &identityV1.GetTenantRequest{
		QueryBy: &identityV1.GetTenantRequest_Id{Id: tenantID},
	})
	if err != nil || tenant == nil || tenant.GetStatus() != identityV1.Tenant_ON || tenant.GetCode() == "" {
		return nil, authenticationV1.ErrorBadRequest("tenant unavailable")
	}

	// 租户归属由服务端决定,调用方不可覆盖
	req.TenantCode = tenant.GetCode()

	// 密码口径对齐:前端提交的是 AES(base64) 密文,core 凭证创建期望明文
	//(入库前 bcrypt);而登录校验会先 AES 解密再比对。若此处不解密,注册会
	// 存下 bcrypt(密文),登录时永远对不上。
	raw, err := base64.StdEncoding.DecodeString(req.GetPassword())
	if err != nil {
		return nil, authenticationV1.ErrorBadRequest("invalid password encoding")
	}
	plain, err := crypto.AesDecrypt(raw, crypto.DefaultAESKey, nil)
	if err != nil {
		return nil, authenticationV1.ErrorBadRequest("invalid password encoding")
	}
	req.Password = string(plain)

	return s.authenticationServiceClient.RegisterUser(ctx, req)
}
