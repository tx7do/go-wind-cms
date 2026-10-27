package service

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pkgcrypto "go-wind-cms/pkg/crypto"

	authenticationV1 "go-wind-cms/api/gen/go/authentication/service/v1"
	identityV1 "go-wind-cms/api/gen/go/identity/service/v1"
	"go-wind-cms/app/core/service/internal/data/ent/privacy"
	"go-wind-cms/pkg/totp"
)

// otpauthIssuer 绑定 URI 的 issuer 段（验证器 App 里展示的条目来源名）。
const otpauthIssuer = "GoWind Content Hub"

// MFAService TOTP 两步验证：绑定（分步注册）、解绑、状态查询、登录挑战校验、管理员救援重置。
// secret 落 sys_user_credentials（ENCRYPTION_KEY 加密）；待确认 secret 暂存 Redis（一次性）。
// 登录挑战换令牌复用 AuthenticationService 的签发链路（同一包内直接持引用）。
type MFAService struct {
	authenticationV1.UnimplementedMFAServiceServer

	auth *AuthenticationService

	log *log.Helper
}

func NewMFAService(
	ctx *bootstrap.Context,
	authService *AuthenticationService,
) *MFAService {
	return &MFAService{
		log:  ctx.NewLoggerHelper("mfa/service/core-service"),
		auth: authService,
	}
}

// enrollPayload 绑定流程的 Redis 暂存体（operation_id → payload）。
type enrollPayload struct {
	UserID   uint32 `json:"uid"`
	TenantID uint32 `json:"tid"`
	Secret   string `json:"secret"`
}

// loadUser 解析目标用户（显式 user_id，不依赖调用方上下文）。
func (s *MFAService) loadUser(ctx context.Context, userID uint32) (*identityV1.User, error) {
	if userID == 0 {
		return nil, authenticationV1.ErrorBadRequest("user id is required")
	}
	return s.auth.userRepo.Get(ctx, &identityV1.GetUserRequest{
		QueryBy: &identityV1.GetUserRequest_Id{Id: userID},
	})
}

// GetMFAStatus 查询用户 MFA 总览。
func (s *MFAService) GetMFAStatus(ctx context.Context, req *authenticationV1.GetMFAStatusRequest) (*authenticationV1.GetMFAStatusResponse, error) {
	userID := parseUint32(req.GetUserId())
	if userID == 0 {
		return nil, authenticationV1.ErrorBadRequest("user id is required")
	}
	enabled, err := s.auth.userCredentialRepo.HasTOTPCredential(ctx, userID)
	if err != nil {
		return nil, err
	}
	resp := &authenticationV1.GetMFAStatusResponse{
		Enabled: enabled,
	}
	if enabled {
		resp.Enrolled = []*authenticationV1.EnrolledMethod{{
			Id:      "totp",
			Method:  authenticationV1.MFAMethod_TOTP,
			Display: "Authenticator App",
			Enabled: true,
		}}
	}
	return resp, nil
}

// ListEnrolledMethods 列出已注册的 MFA 凭证（v1 仅 TOTP 一种）。
func (s *MFAService) ListEnrolledMethods(ctx context.Context, req *authenticationV1.ListEnrolledMethodsRequest) (*authenticationV1.ListEnrolledMethodsResponse, error) {
	status, err := s.GetMFAStatus(ctx, &authenticationV1.GetMFAStatusRequest{UserId: req.UserId})
	if err != nil {
		return nil, err
	}
	return &authenticationV1.ListEnrolledMethodsResponse{Items: status.GetEnrolled()}, nil
}

// StartEnrollMethod 开始绑定（v1 仅 TOTP）：生成 secret 暂存 Redis，返回绑定 URI。
func (s *MFAService) StartEnrollMethod(ctx context.Context, req *authenticationV1.StartEnrollMethodRequest) (*authenticationV1.StartEnrollMethodResponse, error) {
	if req.GetMethod() != authenticationV1.MFAMethod_TOTP {
		return nil, authenticationV1.ErrorBadRequest("only totp enrollment is supported")
	}
	user, err := s.loadUser(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}
	if s.auth.redisClient == nil {
		return nil, authenticationV1.ErrorInternalServerError("enroll store unavailable")
	}

	secret, err := totp.GenerateSecret()
	if err != nil {
		s.log.Errorf("generate totp secret failed: %v", err)
		return nil, authenticationV1.ErrorInternalServerError("generate secret failed")
	}

	payload, err := json.Marshal(enrollPayload{UserID: user.GetId(), TenantID: user.GetTenantId(), Secret: secret})
	if err != nil {
		return nil, authenticationV1.ErrorInternalServerError("marshal enroll payload failed")
	}
	operationID := mfaRandomToken()
	if err := s.auth.redisClient.Set(ctx, mfaEnrollOpKey(operationID), payload, mfaEnrollTTL).Err(); err != nil {
		s.log.Errorf("store enroll operation failed: %v", err)
		return nil, authenticationV1.ErrorInternalServerError("store enroll operation failed")
	}

	account := req.GetAccount()
	if account == "" {
		account = user.GetUsername()
	}
	return &authenticationV1.StartEnrollMethodResponse{
		Result: &authenticationV1.StartEnrollMethodResponse_Totp{
			Totp: &authenticationV1.TOTPResult{
				Secret:     secret,
				OtpAuthUrl: totp.OTPAuthURI(otpauthIssuer, account, secret),
			},
		},
		OperationId: operationID,
		ExpiresAt:   timestamppb.New(time.Now().Add(mfaEnrollTTL)),
	}, nil
}

// ConfirmEnrollMethod 确认绑定：校验动态码通过后写凭据表启用。
func (s *MFAService) ConfirmEnrollMethod(ctx context.Context, req *authenticationV1.ConfirmEnrollMethodRequest) (*authenticationV1.ConfirmEnrollMethodResponse, error) {
	if req.GetMethod() != authenticationV1.MFAMethod_TOTP {
		return nil, authenticationV1.ErrorBadRequest("only totp enrollment is supported")
	}
	operationID := req.GetOperationId()
	if operationID == "" {
		return nil, authenticationV1.ErrorBadRequest("operation id is required")
	}
	if s.auth.redisClient == nil {
		return nil, authenticationV1.ErrorInternalServerError("enroll store unavailable")
	}

	// 一次性取出（失败也消费，须重新发起绑定）
	raw, err := s.auth.redisClient.GetDel(ctx, mfaEnrollOpKey(operationID)).Result()
	if err != nil || raw == "" {
		return nil, authenticationV1.ErrorBadRequest("invalid or expired enroll operation")
	}
	var payload enrollPayload
	if err := json.Unmarshal([]byte(raw), &payload); err != nil || payload.UserID == 0 || payload.Secret == "" {
		return nil, authenticationV1.ErrorBadRequest("invalid enroll operation")
	}

	if !totp.Verify(payload.Secret, req.GetTotpCode(), time.Now()) {
		return nil, authenticationV1.ErrorBadRequest("invalid mfa code")
	}

	encrypted, err := pkgcrypto.EncryptIfNeeded(payload.Secret)
	if err != nil {
		s.log.Errorf("encrypt totp secret failed: %v", err)
		return nil, authenticationV1.ErrorInternalServerError("encrypt secret failed")
	}
	if err := s.auth.userCredentialRepo.SaveTOTPCredential(ctx, payload.UserID, payload.TenantID, encrypted); err != nil {
		return nil, err
	}

	return &authenticationV1.ConfirmEnrollMethodResponse{
		Success:      true,
		CredentialId: "totp",
	}, nil
}

// DisableMFA 解除绑定：v1 支持 totp_code 验证（会话劫持者无动态码即无法关闭）。
func (s *MFAService) DisableMFA(ctx context.Context, req *authenticationV1.DisableMFARequest) (*emptypb.Empty, error) {
	switch req.GetVerifier().(type) {
	case *authenticationV1.DisableMFARequest_TotpCode:
	case *authenticationV1.DisableMFARequest_Password:
		return nil, authenticationV1.ErrorBadRequest("password verifier is not supported, use totp_code")
	default:
		return nil, authenticationV1.ErrorBadRequest("totp_code verifier is required")
	}
	userID := req.GetUserId()
	if userID == 0 {
		return nil, authenticationV1.ErrorBadRequest("user id is required")
	}

	secret, err := s.auth.userCredentialRepo.GetTOTPCredential(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !totp.Verify(secret, req.GetTotpCode(), time.Now()) {
		return nil, authenticationV1.ErrorBadRequest("invalid mfa code")
	}
	if _, err := s.auth.userCredentialRepo.DeleteTOTPCredentials(ctx, userID); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// VerifyMFAChallenge 登录挑战校验：消费一次性挑战凭据 + 校验动态码，签发正式令牌对。
func (s *MFAService) VerifyMFAChallenge(ctx context.Context, req *authenticationV1.VerifyMFAChallengeRequest) (*authenticationV1.VerifyMFAChallengeResponse, error) {
	if req.GetTotpCode() == "" {
		return nil, authenticationV1.ErrorBadRequest("totp code is required")
	}
	resp, err := s.auth.verifyMFAChallengeAndMint(ctx, req.GetOperationId(), req.GetTotpCode(), nil)
	if err != nil {
		return nil, err
	}
	return &authenticationV1.VerifyMFAChallengeResponse{
		Success:          true,
		TokenType:        trans.Ptr(resp.GetTokenType().String()),
		AccessToken:      trans.Ptr(resp.GetAccessToken()),
		ExpiresIn:        trans.Ptr(resp.GetExpiresIn()),
		RefreshToken:     resp.RefreshToken,
		RefreshExpiresIn: resp.RefreshExpiresIn,
		Jti:              resp.Jti,
	}, nil
}

// ResetUserMFA 管理员救援重置：清空目标用户全部 TOTP 因子（跨租户，系统权限）。
func (s *MFAService) ResetUserMFA(ctx context.Context, req *authenticationV1.ResetUserMFARequest) (*emptypb.Empty, error) {
	if req.GetUserId() == 0 {
		return nil, authenticationV1.ErrorBadRequest("user id is required")
	}
	// 平台管理员清理任意租户用户的因子：显式系统权限，绕过租户行级谓词
	if _, err := s.auth.userCredentialRepo.DeleteTOTPCredentials(privacy.DecisionContext(ctx, privacy.Allow), req.GetUserId()); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// parseUint32 宽松解析数字串（非法返回 0）。
func parseUint32(v string) uint32 {
	v = strings.TrimSpace(v)
	var n uint64
	for _, c := range v {
		if c < '0' || c > '9' {
			return 0
		}
		n = n*10 + uint64(c-'0')
		if n > 0xFFFFFFFF {
			return 0
		}
	}
	return uint32(n)
}
