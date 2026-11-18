package service

import (
	"context"
	"encoding/base64"
	"regexp"
	"strings"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/tx7do/go-utils/captcha"
	"github.com/tx7do/go-utils/crypto"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	appV1 "go-wind-cms/api/gen/go/app/service/v1"
	authenticationV1 "go-wind-cms/api/gen/go/authentication/service/v1"
	identityV1 "go-wind-cms/api/gen/go/identity/service/v1"
	entmiddleware "go-wind-cms/pkg/middleware/ent"

	"go-wind-cms/app/app/service/internal/data"
	"go-wind-cms/pkg/middleware/auth"
	"go-wind-cms/pkg/netutil"
)

// 验证码相关请求头（渐进式验证码：通过 header 传递以避免改动 proto 与四端生成代码）。
const (
	headerCaptchaID        = "X-Captcha-Id"
	headerCaptchaValue     = "X-Captcha-Value"
	headerCaptchaSupported = "X-Captcha-Supported"
)

// CaptchaEnabled 控制登录是否启用验证码校验（渐进式闸门的总开关）。
const CaptchaEnabled = true

// captchaRequiredThreshold 渐进式验证码阈值：同一 ip+用户名的登录失败计数
// 达到该值后，声明支持验证码的客户端（X-Captcha-Supported）强制验证码。
// 计数由登录失败自增（rateLimiter.CheckAndIncr）维护、登录成功 Reset，
// 滑动窗口 loginFailKeyTTL 过期自然清零。
const captchaRequiredThreshold = 1

// credentialFailReasons 计入登录失败数的凭据错误 reason 族。
// 验证码错误与基础设施错误不计；INVALID_TENANT 为 core 对齐 reason 后
// 自动生效的预留项。
var credentialFailReasons = map[string]struct{}{
	"INVALID_PASSWORD": {},
	"USER_NOT_FOUND":   {},
	"INVALID_TENANT":   {},
}

type AuthenticationService struct {
	appV1.AuthenticationServiceHTTPServer

	authenticationServiceClient authenticationV1.AuthenticationServiceClient
	tenantClient                identityV1.TenantServiceClient
	tenantResolve               entmiddleware.TenantResolver

	captchaClient *captcha.Captcha
	rateLimiter   *data.LoginRateLimiter

	log *log.Helper
}

func NewAuthenticationService(
	ctx *bootstrap.Context,
	authenticationServiceClient authenticationV1.AuthenticationServiceClient,
	tenantClient identityV1.TenantServiceClient,
	tenantResolve entmiddleware.TenantResolver,
	captchaClient *captcha.Captcha,
	rateLimiter *data.LoginRateLimiter,
) *AuthenticationService {
	return &AuthenticationService{
		log:                         ctx.NewLoggerHelper("authn/service/app-service"),
		authenticationServiceClient: authenticationServiceClient,
		tenantClient:                tenantClient,
		tenantResolve:               tenantResolve,
		captchaClient:               captchaClient,
		rateLimiter:                 rateLimiter,
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

	if req.GetGrantType() == authenticationV1.GrantType_password {
		clientIP := req.GetClientIp()
		username := req.GetUsername()

		// ===== 闸门 1：登录限流前置（IP + 用户名双维度，5 次失败锁 15 分钟）=====
		// 全端生效——未声明验证码支持的端（taro/flutter）无免验证码豁免，
		// 锁定是防爆破主力，不因端能力而放行。
		if s.rateLimiter != nil {
			if locked, lerr := s.rateLimiter.IsLocked(ctx, clientIP, username); lerr != nil {
				s.log.Errorf("login rate limiter pre-check failed: %s", lerr.Error())
			} else if locked {
				s.log.Warnf("login blocked by rate limiter: ip=%s, username=%s", clientIP, username)
				return nil, appV1.ErrorBadRequest("too many login failures, please try again later")
			}
		}

		// ===== 闸门 2：渐进式验证码（仅声明支持验证码的客户端）=====
		// taro/flutter 未实现验证码 UI，上报 X-Captcha-Supported 免受验证码 gate
		// （否则错一次密码就会在 15 分钟窗口内被 CAPTCHA_REQUIRED 卡死无法登录）。
		if s.captchaSupported(ctx) {
			if err := s.verifyLoginCaptcha(ctx, clientIP, username); err != nil {
				return nil, err
			}
		}
	}

	resp, err := s.authenticationServiceClient.Login(ctx, req)
	if err != nil {
		// 登录失败（凭据族错误）时自增失败计数（按 IP + 用户名双维度），
		// 供渐进式验证码与锁定闸门使用；验证码错误/基础设施错误不计数。
		if req.GetGrantType() == authenticationV1.GrantType_password && isCredentialFailure(err) && s.rateLimiter != nil {
			if _, _, _, cerr := s.rateLimiter.CheckAndIncr(ctx, req.GetClientIp(), req.GetUsername()); cerr != nil {
				s.log.Errorf("login rate limiter incr failed: %s", cerr.Error())
			}
		}
		return nil, err
	}

	// 登录成功后清零失败计数（回到免验证码态）
	if req.GetGrantType() == authenticationV1.GrantType_password && s.rateLimiter != nil {
		s.rateLimiter.Reset(ctx, req.GetClientIp(), req.GetUsername())
	}

	return resp, nil
}

// GenerateCaptcha 生成图形验证码（透传 captchaClient；免鉴权端点，见 rest_server 白名单）。
func (s *AuthenticationService) GenerateCaptcha(ctx context.Context, _ *emptypb.Empty) (*authenticationV1.GenerateCaptchaResponse, error) {
	if s.captchaClient == nil {
		return nil, appV1.ErrorServiceUnavailable("captcha unavailable")
	}
	captchaId, captchaImage, answer, err := s.captchaClient.Generate()
	if err != nil {
		s.log.Errorf("generate captcha failed: %s", err.Error())
		return nil, appV1.ErrorInternalServerError("generate captcha failed")
	}

	// Generate() 只生成验证码但不落盘，必须手动 Save 到 Redis，否则 Verify 时查不到。
	if err = s.captchaClient.Save(ctx, captchaId, answer); err != nil {
		s.log.Errorf("save captcha failed: %s", err.Error())
		return nil, appV1.ErrorInternalServerError("save captcha failed")
	}

	return &authenticationV1.GenerateCaptchaResponse{
		CaptchaId:   captchaId,
		ImageBase64: captchaImage,
	}, nil
}

// captchaSupported 判断客户端是否声明支持验证码 UI（X-Captcha-Supported: 1）。
// 声明缺失（老前端/未适配端）视为不支持——渐进式验证码只对明确支持的端生效。
func (s *AuthenticationService) captchaSupported(ctx context.Context) bool {
	if !CaptchaEnabled || s.captchaClient == nil {
		return false
	}
	header := netutil.HeaderFromContext(ctx)
	return header != nil && strings.TrimSpace(header.Get(headerCaptchaSupported)) == "1"
}

// verifyLoginCaptcha 渐进式校验登录请求携带的验证码。
// 验证码 id/value 通过 HTTP Header（X-Captcha-Id / X-Captcha-Value）传递。
// 仅当同一 ip+用户名的登录失败计数达到 captchaRequiredThreshold 时强制——
// 计数由登录失败自增维护、登录成功 Reset，未达阈值时不校验也不消费验证码。
// captchaClient.Verify 已是 verify-and-delete 单次有效语义。
func (s *AuthenticationService) verifyLoginCaptcha(ctx context.Context, clientIP, username string) error {
	// 渐进式闸门：失败计数读的是限流器的只读快照；Redis 不可用时返回 0，
	// 与锁定闸门的 fail-open 语义一致（限流是防御增强，不应故障时堵死登录）。
	if s.rateLimiter == nil || s.rateLimiter.FailureAttempts(ctx, clientIP, username) < captchaRequiredThreshold {
		return nil
	}
	header := netutil.HeaderFromContext(ctx)
	if header == nil {
		return appV1.ErrorCaptchaRequired("captcha required")
	}
	captchaID := strings.TrimSpace(header.Get(headerCaptchaID))
	captchaValue := strings.TrimSpace(header.Get(headerCaptchaValue))
	if captchaID == "" || captchaValue == "" {
		return appV1.ErrorCaptchaRequired("captcha required")
	}
	ok, err := s.captchaClient.Verify(ctx, captchaID, captchaValue)
	if err != nil {
		s.log.Errorf("verify captcha failed: %s", err.Error())
		// 校验异常也销毁该验证码（不可信状态下让单码不再可用，缩小在线试猜窗口）。
		_ = s.captchaClient.Delete(ctx, captchaID)
		return appV1.ErrorCaptchaInvalid("invalid captcha")
	}
	if !ok {
		// 校验失败（含不匹配/过期不存在）销毁验证码——失败后单码仍可用会在
		// 10 分钟 TTL 内形成在线试猜窗口。
		_ = s.captchaClient.Delete(ctx, captchaID)
		return appV1.ErrorCaptchaInvalid("invalid captcha")
	}
	return nil
}

// isCredentialFailure 判断转发错误是否为凭据校验失败（该计入登录失败数的类型）。
func isCredentialFailure(err error) bool {
	se := errors.FromError(err)
	if se == nil {
		return false
	}
	_, ok := credentialFailReasons[se.Reason]
	return ok
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

	var tenantID int64
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
