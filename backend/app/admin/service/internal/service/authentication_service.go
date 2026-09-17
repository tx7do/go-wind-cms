package service

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	ktransport "github.com/go-kratos/kratos/v2/transport"
	"github.com/tx7do/go-utils/captcha"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	authenticationV1 "go-wind-cms/api/gen/go/authentication/service/v1"

	"go-wind-cms/pkg/middleware/auth"
	"go-wind-cms/pkg/netutil"
)

// 验证码相关请求头（H5：登录强制验证码，通过 header 传递以避免改动 proto 与三套前端生成代码）。
const (
	headerCaptchaID    = "X-Captcha-Id"
	headerCaptchaValue = "X-Captcha-Value"
)

// CaptchaEnabled 控制登录是否强制校验验证码。
// 开发/无 Redis 等环境可改为 false 跳过验证码校验，避免登录被 400 invalid or missing captcha 阻断。
const CaptchaEnabled = true

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
}

func NewAuthenticationService(
	ctx *bootstrap.Context,
	authenticationServiceClient authenticationV1.AuthenticationServiceClient,
	captchaClient *captcha.Captcha,
) *AuthenticationService {
	return &AuthenticationService{
		log:                         log.NewHelper(log.With(ctx.GetLogger(), "module", "user/service/admin-service")),
		authenticationServiceClient: authenticationServiceClient,
		captchaClient:               captchaClient,
	}
}

// verifyLoginCaptcha 校验登录请求携带的验证码。
// 验证码 id/value 通过 HTTP Header（X-Captcha-Id / X-Captcha-Value）传递。
// captchaClient.Verify 已是 verify-and-delete 单次有效语义。
// 注意：refresh_token / client_credentials 等非密码授权不走此校验（仅 password 授权调用）。
func (s *AuthenticationService) verifyLoginCaptcha(ctx context.Context) bool {
	if !CaptchaEnabled {
		// 验证码开关关闭，跳过校验
		return true
	}
	if s.captchaClient == nil {
		// captcha 未配置时 fail-open（仅记录告警），避免影响登录基本功能
		return true
	}
	header := netutil.HeaderFromContext(ctx)
	if header == nil {
		return false
	}
	captchaID := strings.TrimSpace(header.Get(headerCaptchaID))
	captchaValue := strings.TrimSpace(header.Get(headerCaptchaValue))
	if captchaID == "" || captchaValue == "" {
		return false
	}
	ok, err := s.captchaClient.Verify(ctx, captchaID, captchaValue)
	if err != nil {
		s.log.Errorf("verify captcha failed: %s", err.Error())
		return false
	}
	return ok
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

// Login 登录
func (s *AuthenticationService) Login(ctx context.Context, req *authenticationV1.LoginRequest) (*authenticationV1.LoginResponse, error) {
	if req == nil {
		return nil, authenticationV1.ErrorBadRequest("invalid request")
	}

	req.ClientType = trans.Ptr(authenticationV1.ClientType_admin)

	if req.GetGrantType() == authenticationV1.GrantType_refresh_token {
		operator, err := auth.FromContext(ctx)
		if err != nil {
			return nil, err
		}

		req.Jti = operator.Jti
		req.UserId = trans.Ptr(operator.GetUserId())
	} else if req.GetGrantType() == authenticationV1.GrantType_password {
		// ===== 强制验证码（仅密码授权；通过 HTTP Header 传递，避免改动 proto/前端生成代码）=====
		if !s.verifyLoginCaptcha(ctx) {
			return nil, authenticationV1.ErrorBadRequest("invalid or missing captcha")
		}
	}

	resp, err := s.authenticationServiceClient.Login(ctx, req)
	if err != nil {
		return nil, err
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
