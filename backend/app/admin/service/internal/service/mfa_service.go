package service

import (
	"context"
	"strconv"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/tx7do/go-utils/trans"

	"go-wind-cms/pkg/middleware/auth"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	authenticationV1 "go-wind-cms/api/gen/go/authentication/service/v1"
)

// MFAService 多因素认证 BFF 面：代理 core MFAService。
// me/mfa* 端点以会话操作者身份注入 user_id（core 不信任调用方自报身份）；
// challenge/verify 免鉴权透传（半登录态换正式令牌）。
type MFAService struct {
	adminV1.MFAServiceHTTPServer

	log *log.Helper

	mfaServiceClient authenticationV1.MFAServiceClient
}

func NewMFAService(
	ctx *bootstrap.Context,
	mfaServiceClient authenticationV1.MFAServiceClient,
) *MFAService {
	return &MFAService{
		log:              ctx.NewLoggerHelper("mfa/service/admin-service"),
		mfaServiceClient: mfaServiceClient,
	}
}

func (s *MFAService) GetMFAStatus(ctx context.Context, _ *authenticationV1.GetMFAStatusRequest) (*authenticationV1.GetMFAStatusResponse, error) {
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	return s.mfaServiceClient.GetMFAStatus(ctx, &authenticationV1.GetMFAStatusRequest{
		UserId: trans.Ptr(strconv.FormatUint(uint64(operator.GetUserId()), 10)),
	})
}

func (s *MFAService) StartEnrollMethod(ctx context.Context, req *authenticationV1.StartEnrollMethodRequest) (*authenticationV1.StartEnrollMethodResponse, error) {
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	if operator.GetUserId() == 0 {
		return nil, authenticationV1.ErrorBadRequest("operator not resolved")
	}
	req.UserId = trans.Ptr(operator.GetUserId())
	req.Account = trans.Ptr(operator.GetUsername())
	return s.mfaServiceClient.StartEnrollMethod(ctx, req)
}

func (s *MFAService) ConfirmEnrollMethod(ctx context.Context, req *authenticationV1.ConfirmEnrollMethodRequest) (*authenticationV1.ConfirmEnrollMethodResponse, error) {
	return s.mfaServiceClient.ConfirmEnrollMethod(ctx, req)
}

func (s *MFAService) DisableMFA(ctx context.Context, req *authenticationV1.DisableMFARequest) (*emptypb.Empty, error) {
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	if operator.GetUserId() == 0 {
		return nil, authenticationV1.ErrorBadRequest("operator not resolved")
	}
	req.UserId = trans.Ptr(operator.GetUserId())
	return s.mfaServiceClient.DisableMFA(ctx, req)
}

func (s *MFAService) VerifyMFAChallenge(ctx context.Context, req *authenticationV1.VerifyMFAChallengeRequest) (*authenticationV1.LoginResponse, error) {
	resp, err := s.mfaServiceClient.VerifyMFAChallenge(ctx, req)
	if err != nil {
		return nil, err
	}
	return &authenticationV1.LoginResponse{
		TokenType:        authenticationV1.TokenType_bearer,
		AccessToken:      resp.GetAccessToken(),
		ExpiresIn:        resp.GetExpiresIn(),
		RefreshToken:     resp.RefreshToken,
		RefreshExpiresIn: resp.RefreshExpiresIn,
		Jti:              resp.Jti,
	}, nil
}

func (s *MFAService) ResetUserMFA(ctx context.Context, req *authenticationV1.ResetUserMFARequest) (*emptypb.Empty, error) {
	return s.mfaServiceClient.ResetUserMFA(ctx, req)
}
