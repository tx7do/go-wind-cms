package service

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	ktransport "github.com/go-kratos/kratos/v2/transport"
	khttp "github.com/go-kratos/kratos/v2/transport/http"

	"github.com/tx7do/go-utils/captcha"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"go-wind-cms/pkg/middleware/logging"
	"google.golang.org/protobuf/types/known/emptypb"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	authenticationV1 "go-wind-cms/api/gen/go/authentication/service/v1"

	"go-wind-cms/app/admin/service/internal/data"
	"go-wind-cms/pkg/middleware/auth"
	"go-wind-cms/pkg/netutil"
)

// 验证码相关请求头（H5：登录强制验证码，通过 header 传递以避免改动 proto 与三套前端生成代码）。
const (
	headerCaptchaID    = "X-Captcha-Id"
	headerCaptchaValue = "X-Captcha-Value"
)

// CaptchaEnabled 控制登录是否启用验证码校验（渐进式闸门的总开关）。
// 开发/无 Redis 等环境可改为 false 跳过验证码校验，避免登录被验证码闸门阻断。
const CaptchaEnabled = true

// captchaRequiredThreshold 渐进式验证码阈值：同一 ip+用户名的登录失败计数
// 达到该值后，后续登录强制验证码。计数由登录失败自增（rateLimiter.CheckAndIncr）
// 维护、登录成功 Reset，滑动窗口 loginFailKeyTTL 过期自然清零——阈值取 1 即
// 「错一次下一次就要验证码」，与锁定闸门（阈值 5）之间留出验证码阶段。
const captchaRequiredThreshold = 1

// credentialFailReasons 计入登录失败数的凭据错误 reason 族。
// BFF 形态下凭证校验在 core，凭转发错误 reason 判定「这次登录失败该不该计数」：
// 验证码错误（CAPTCHA_INVALID/REQUIRED）不是凭据爆破，不计；core 不可用等
// 基础设施错误不计。INVALID_TENANT 为 core 对齐 reason 后自动生效的预留项。
var credentialFailReasons = map[string]struct{}{
	"INVALID_PASSWORD": {},
	"USER_NOT_FOUND":   {},
	"INVALID_TENANT":   {},
}

// refresh token cookie 相关常量（对齐 go-wind-admin 上游方案）。
// refresh token 以 HttpOnly Cookie 传输，Path 收窄到刷新端点，SameSite=Lax 阻断跨站 POST。
// refresh_exp 为非 HttpOnly 的过期时间戳 cookie，供前端定时器/启动静默恢复读取。
const (
	refreshTokenCookieName = "refresh_token"
	refreshExpCookieName   = "refresh_exp"
	refreshCookiePath      = "/admin/v1/refresh-token"
	defaultRefreshTTL      = 30 * 24 * time.Hour
)

// resolveCookieSecure 按请求的实际传输层判断是否加 Secure 属性：
// TLS 直连或经可信反代（X-Forwarded-Proto: https）→ true；明文 HTTP → false。
// 依据：带 Secure 的 cookie 在明文 HTTP 下会被浏览器拒收，导致 refresh token 无法落地。
func resolveCookieSecure(htr *khttp.Transport) bool {
	req := htr.Request()
	if req == nil {
		return true
	}
	if req.TLS != nil {
		return true
	}
	return strings.EqualFold(req.Header.Get("X-Forwarded-Proto"), "https")
}

// setRefreshCookies 将 refresh token 及其过期时间戳写入 Set-Cookie 响应头。
// refresh_token cookie 为 HttpOnly（JS 不可读）+ Path 收窄到刷新端点；
// refresh_exp cookie 非 HttpOnly（前端定时器/启动恢复可读）+ Path=/。
func setRefreshCookies(ctx context.Context, refreshToken string, refreshExpiresInSeconds int64) {
	tr, ok := ktransport.FromServerContext(ctx)
	if !ok {
		return
	}
	htr, hok := tr.(*khttp.Transport)
	if !hok {
		return
	}
	header := htr.ReplyHeader()
	cookieSecure := resolveCookieSecure(htr)

	rtCookie := &http.Cookie{
		Name:     refreshTokenCookieName,
		Value:    refreshToken,
		Path:     refreshCookiePath,
		MaxAge:   int(refreshExpiresInSeconds),
		Secure:   cookieSecure,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	header.Add("Set-Cookie", rtCookie.String())

	expCookie := &http.Cookie{
		Name:     refreshExpCookieName,
		Value:    fmt.Sprintf("%d", time.Now().Unix()+refreshExpiresInSeconds),
		Path:     "/",
		MaxAge:   int(refreshExpiresInSeconds),
		Secure:   cookieSecure,
		HttpOnly: false,
		SameSite: http.SameSiteLaxMode,
	}
	header.Add("Set-Cookie", expCookie.String())
}

// clearRefreshCookies 写入 Max-Age=0 的清除 cookie。
// 两个 cookie 的 Path 不同，清除时须各自匹配原 Path，否则浏览器拒绝删除。
func clearRefreshCookies(ctx context.Context) {
	tr, ok := ktransport.FromServerContext(ctx)
	if !ok {
		return
	}
	htr, hok := tr.(*khttp.Transport)
	if !hok {
		return
	}
	header := htr.ReplyHeader()
	cookieSecure := resolveCookieSecure(htr)

	for name, path := range map[string]string{
		refreshTokenCookieName: refreshCookiePath,
		refreshExpCookieName:   "/",
	} {
		cookie := &http.Cookie{
			Name:     name,
			Value:    "",
			Path:     path,
			MaxAge:   -1,
			Secure:   cookieSecure,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		}
		header.Add("Set-Cookie", cookie.String())
	}
}

// refreshCookieFromRequest 从请求 cookie 中取 refresh token 值。
// 页面刷新后 access/refresh token 均已脱离内存，绑定键 (uid/jti) 由前端从
// 持久化的 access token payload 解出随请求体回传，刷新凭证本身只走 HttpOnly cookie。
func refreshCookieFromRequest(ctx context.Context) string {
	tr, ok := ktransport.FromServerContext(ctx)
	if !ok {
		return ""
	}
	htr, hok := tr.(*khttp.Transport)
	if !hok || htr.Request() == nil {
		return ""
	}
	if c, err := htr.Request().Cookie(refreshTokenCookieName); err == nil {
		return c.Value
	}
	return ""
}

type AuthenticationService struct {
	adminV1.AuthenticationServiceHTTPServer

	log *log.Helper

	authenticationServiceClient authenticationV1.AuthenticationServiceClient

	captchaClient *captcha.Captcha

	rateLimiter *data.LoginRateLimiter
}

func NewAuthenticationService(
	ctx *bootstrap.Context,
	authenticationServiceClient authenticationV1.AuthenticationServiceClient,
	captchaClient *captcha.Captcha,
	rateLimiter *data.LoginRateLimiter,
) *AuthenticationService {
	return &AuthenticationService{
		log:                         log.NewHelper(log.With(ctx.GetLogger(), "module", "user/service/admin-service")),
		authenticationServiceClient: authenticationServiceClient,
		captchaClient:               captchaClient,
		rateLimiter:                 rateLimiter,
	}
}

// GetOidcAuthorizeUrl 获取 OIDC 授权跳转地址（透传 core；state 由 core 生成并落 Redis）。
// 免鉴权端点（登录前置步骤），不走验证码。
func (s *AuthenticationService) GetOidcAuthorizeUrl(ctx context.Context, req *authenticationV1.GetOidcAuthorizeUrlRequest) (*authenticationV1.GetOidcAuthorizeUrlResponse, error) {
	return s.authenticationServiceClient.GetOidcAuthorizeUrl(ctx, req)
}

// verifyLoginCaptcha 渐进式校验登录请求携带的验证码。
// 验证码 id/value 通过 HTTP Header（X-Captcha-Id / X-Captcha-Value）传递。
// 仅当同一 ip+用户名的登录失败计数达到 captchaRequiredThreshold 时强制——
// 计数由登录失败自增维护、登录成功 Reset，未达阈值时不校验也不消费验证码。
// captchaClient.Verify 已是 verify-and-delete 单次有效语义。
// 注意：refresh_token / client_credentials 等非密码授权不走此校验（仅 password 授权调用）。
func (s *AuthenticationService) verifyLoginCaptcha(ctx context.Context, clientIP, username string) error {
	if !CaptchaEnabled {
		// 验证码开关关闭，跳过校验
		return nil
	}
	if s.captchaClient == nil {
		// captcha 未配置时 fail-open（仅记录告警），避免影响登录基本功能
		return nil
	}
	// 渐进式闸门：失败计数读的是限流器的只读快照；Redis 不可用时返回 0，
	// 与锁定闸门的 fail-open 语义一致（限流是防御增强，不应故障时堵死登录）。
	if s.rateLimiter == nil || s.rateLimiter.FailureAttempts(ctx, clientIP, username) < captchaRequiredThreshold {
		return nil
	}
	header := netutil.HeaderFromContext(ctx)
	if header == nil {
		return adminV1.ErrorCaptchaRequired("captcha required")
	}
	captchaID := strings.TrimSpace(header.Get(headerCaptchaID))
	captchaValue := strings.TrimSpace(header.Get(headerCaptchaValue))
	if captchaID == "" || captchaValue == "" {
		return adminV1.ErrorCaptchaRequired("captcha required")
	}
	ok, err := s.captchaClient.Verify(ctx, captchaID, captchaValue)
	if err != nil {
		s.log.Errorf("verify captcha failed: %s", err.Error())
		// 校验异常也销毁该验证码（不可信状态下让单码不再可用，缩小在线试猜窗口）。
		_ = s.captchaClient.Delete(ctx, captchaID)
		return adminV1.ErrorCaptchaInvalid("invalid captcha")
	}
	if !ok {
		// 校验失败（含不匹配/过期不存在）销毁验证码——失败后单码仍可用会在
		// 10 分钟 TTL 内形成在线试猜窗口。
		_ = s.captchaClient.Delete(ctx, captchaID)
		return adminV1.ErrorCaptchaInvalid("invalid captcha")
	}
	return nil
}

// GetTenantBranding 租户白标查询（透传 core：租户数据在 core）。
func (s *AuthenticationService) GetTenantBranding(ctx context.Context, req *authenticationV1.GetTenantBrandingRequest) (*authenticationV1.GetTenantBrandingResponse, error) {
	return s.authenticationServiceClient.GetTenantBranding(ctx, req)
}

// ForgotPassword 忘记密码（透传 core：凭据、验证码缓存与通知缝都在 core）。
func (s *AuthenticationService) ForgotPassword(ctx context.Context, req *authenticationV1.ForgotPasswordRequest) (*emptypb.Empty, error) {
	return s.authenticationServiceClient.ForgotPassword(ctx, req)
}

// ResetPasswordByCode 凭邮箱验证码重置密码（透传 core）。
func (s *AuthenticationService) ResetPasswordByCode(ctx context.Context, req *authenticationV1.ResetPasswordByCodeRequest) (*emptypb.Empty, error) {
	return s.authenticationServiceClient.ResetPasswordByCode(ctx, req)
}

func (s *AuthenticationService) GenerateCaptcha(ctx context.Context, _ *emptypb.Empty) (*authenticationV1.GenerateCaptchaResponse, error) {
	captchaId, captchaImage, answer, err := s.captchaClient.Generate()
	if err != nil {
		s.log.Errorf("generate captcha failed: %s", err.Error())
		return nil, authenticationV1.ErrorInternalServerError("generate captcha failed")
	}

	// Generate() 只生成验证码但不落盘，必须手动 Save 到 Redis，否则 Verify 时查不到。
	if err = s.captchaClient.Save(ctx, captchaId, answer); err != nil {
		s.log.Errorf("save captcha failed: %s", err.Error())
		return nil, authenticationV1.ErrorInternalServerError("save captcha failed")
	}

	return &authenticationV1.GenerateCaptchaResponse{
		CaptchaId:   captchaId,
		ImageBase64: captchaImage,
	}, nil
}

func (s *AuthenticationService) VerifyCaptcha(ctx context.Context, req *authenticationV1.VerifyCaptchaRequest) (*authenticationV1.VerifyCaptchaResponse, error) {
	ok, err := s.captchaClient.Verify(ctx, req.GetCaptchaId(), req.GetUserInput())
	if err != nil {
		s.log.Errorf("verify captcha failed: %s", err.Error())
		return nil, authenticationV1.ErrorInternalServerError("verify captcha failed")
	}

	return &authenticationV1.VerifyCaptchaResponse{
		Valid: ok,
	}, nil
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

// Login 登录
func (s *AuthenticationService) Login(ctx context.Context, req *authenticationV1.LoginRequest) (*authenticationV1.LoginResponse, error) {
	if req == nil {
		return nil, authenticationV1.ErrorBadRequest("invalid request")
	}

	req.ClientType = trans.Ptr(authenticationV1.ClientType_admin)

	// 会话元数据：登录/刷新的客户端信息由 BFF 从 HTTP 请求注入
	//（core 无法感知 HTTP 层；在线会话列表的 IP/UA 即来源于此）。
	if tr, ok := ktransport.FromServerContext(ctx); ok {
		if htr, hok := tr.(*khttp.Transport); hok && htr.Request() != nil {
			httpReq := htr.Request()
			if req.ClientIp == nil {
				req.ClientIp = trans.Ptr(logging.ClientRealIP(httpReq))
			}
			if req.UserAgent == nil {
				req.UserAgent = trans.Ptr(httpReq.UserAgent())
			}
		}
	}

	if req.GetGrantType() == authenticationV1.GrantType_refresh_token {
		operator, err := auth.FromContext(ctx)
		if err != nil {
			return nil, err
		}

		req.Jti = operator.Jti
		req.UserId = trans.Ptr(operator.GetUserId())
	} else if req.GetGrantType() == authenticationV1.GrantType_password {
		clientIP := req.GetClientIp()
		username := req.GetUsername()

		// ===== H5 闸门 1：登录限流前置（IP + 用户名双维度，5 次失败锁 15 分钟）=====
		// 计数在下方转发失败的凭据错误分支自增、登录成功 Reset。
		if s.rateLimiter != nil {
			if locked, lerr := s.rateLimiter.IsLocked(ctx, clientIP, username); lerr != nil {
				s.log.Errorf("login rate limiter pre-check failed: %s", lerr.Error())
			} else if locked {
				s.log.Warnf("login blocked by rate limiter: ip=%s, username=%s", clientIP, username)
				return nil, adminV1.ErrorBadRequest("too many login failures, please try again later")
			}
		}

		// ===== H5 闸门 2：渐进式验证码（同一 ip+用户名失败计数达阈值后强制；
		// 通过 HTTP Header 传递，避免改动 proto/前端生成代码）=====
		// 未达阈值直接放行；达到阈值后缺失返回 CAPTCHA_REQUIRED（登录页按此
		// reason 点亮验证码字段）、校验不通过返回 CAPTCHA_INVALID。
		if err := s.verifyLoginCaptcha(ctx, clientIP, username); err != nil {
			return nil, err
		}
	}

	resp, err := s.authenticationServiceClient.Login(ctx, req)
	if err != nil {
		// H5：登录失败（凭据族错误）时自增失败计数（按 IP + 用户名双维度），
		// 供渐进式验证码与锁定闸门使用；验证码错误/基础设施错误不计数。
		if req.GetGrantType() == authenticationV1.GrantType_password && isCredentialFailure(err) && s.rateLimiter != nil {
			if _, _, _, cerr := s.rateLimiter.CheckAndIncr(ctx, req.GetClientIp(), req.GetUsername()); cerr != nil {
				s.log.Errorf("login rate limiter incr failed: %s", cerr.Error())
			}
		}
		return nil, err
	}

	// H5：登录成功后清零失败计数（回到免验证码态）
	if req.GetGrantType() == authenticationV1.GrantType_password && s.rateLimiter != nil {
		s.rateLimiter.Reset(ctx, req.GetClientIp(), req.GetUsername())
	}

	// refresh token 经 HttpOnly Cookie 下发，不再放入响应体（对齐 go-wind-admin 上游）
	if rt := resp.GetRefreshToken(); rt != "" {
		expiresIn := resp.GetRefreshExpiresIn()
		if expiresIn <= 0 {
			expiresIn = int64(defaultRefreshTTL.Seconds())
		}
		setRefreshCookies(ctx, rt, expiresIn)
		resp.RefreshToken = nil
	}

	return resp, nil
}

// Logout 登出
func (s *AuthenticationService) Logout(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	// 清除 refresh token 相关 cookie（对齐 go-wind-admin 上游）
	clearRefreshCookies(ctx)

	return s.authenticationServiceClient.Logout(ctx, &authenticationV1.LogoutRequest{
		ClientType: authenticationV1.ClientType_admin,
		UserId:     operator.GetUserId(),
	})
}

// RefreshToken 刷新令牌
func (s *AuthenticationService) RefreshToken(ctx context.Context, req *authenticationV1.LoginRequest) (*authenticationV1.LoginResponse, error) {
	if req == nil {
		return nil, authenticationV1.ErrorBadRequest("invalid request")
	}

	// 刷新凭证本体优先取 HttpOnly cookie（页面刷新后前端已无任何内存态可依赖）。
	if req.GetRefreshToken() == "" {
		if cv := refreshCookieFromRequest(ctx); cv != "" {
			req.RefreshToken = trans.Ptr(cv)
		}
	}

	// 绑定键 (uid/jti) 回填，优先级从高到低：
	// 1) 请求体（前端从持久化 access token payload 解出后回传——页面刷新后
	//    内存已空，这是唯一来源；uid/jti 本身非机密，凭证校验靠 refresh token 值）；
	// 2) auth 中间件注入的 operator（access token 仍有效的在途刷新）；
	// 3) 过期 JWT 的未验签 payload（与 2 同源，仅为容忍头存在但已过期的场景）。
	if req.GetUserId() == 0 || req.GetJti() == "" {
		if operator, oerr := auth.FromContext(ctx); oerr == nil {
			if req.GetUserId() == 0 {
				req.UserId = trans.Ptr(operator.GetUserId())
			}
			if req.GetJti() == "" {
				req.Jti = operator.Jti
			}
		} else if uid, jti, uerr := netutil.ParseUnverifiedBearerJWTClaims(ctx); uerr == nil {
			if req.GetUserId() == 0 {
				req.UserId = trans.Ptr(uid)
			}
			if req.GetJti() == "" {
				req.Jti = trans.Ptr(jti)
			}
		}
	}

	req.ClientType = trans.Ptr(authenticationV1.ClientType_admin)

	resp, err := s.authenticationServiceClient.RefreshToken(ctx, req)
	if err != nil {
		return nil, err
	}

	// 轮换：下发新 refresh cookie，响应体不再携带 refresh token
	if rt := resp.GetRefreshToken(); rt != "" {
		expiresIn := resp.GetRefreshExpiresIn()
		if expiresIn <= 0 {
			expiresIn = int64(defaultRefreshTTL.Seconds())
		}
		setRefreshCookies(ctx, rt, expiresIn)
		resp.RefreshToken = nil
	}

	return resp, nil
}

func (s *AuthenticationService) WhoAmI(ctx context.Context, _ *emptypb.Empty) (*authenticationV1.WhoAmIResponse, error) {
	// 获取操作人信息
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	return &authenticationV1.WhoAmIResponse{
		UserId:   operator.GetUserId(),
		Username: operator.GetUsername(),
	}, nil
}
