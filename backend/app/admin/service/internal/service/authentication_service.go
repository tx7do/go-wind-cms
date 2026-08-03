package service

import (
	"context"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
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

	return s.authenticationServiceClient.Login(ctx, req)
}

// Logout 登出
func (s *AuthenticationService) Logout(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

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

	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.ClientType = trans.Ptr(authenticationV1.ClientType_admin)
	req.UserId = trans.Ptr(operator.GetUserId())
	req.Jti = operator.Jti

	return s.authenticationServiceClient.RefreshToken(ctx, req)
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
