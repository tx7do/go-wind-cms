package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/tx7do/go-utils/trans"

	authenticationV1 "go-wind-cms/api/gen/go/authentication/service/v1"
	identityV1 "go-wind-cms/api/gen/go/identity/service/v1"
	"go-wind-cms/pkg/totp"
)

// MFA 两步验证（TOTP）登录挑战：
//   - 密码/SSO 验证通过但用户已启用 TOTP 因子时，Login 不发令牌，改为签发一次性
//     挑战凭据 mfa_operation_id（Redis 暂存，短时效），客户端凭它 + 6 位动态码走
//     MFAService.VerifyMFAChallenge（或 Login 内联回传）换正式令牌；
//   - 挑战凭据一次校验即消费（无论成败），失败须重新走完整登录，天然限速。

const (
	// mfaChallengeTTL 挑战凭据有效期（半登录态，短时效）
	mfaChallengeTTL = 5 * time.Minute
	// mfaEnrollTTL 绑定流程中待确认 secret 的暂存有效期
	mfaEnrollTTL = 10 * time.Minute
)

// mfaChallengePayload 挑战凭据承载的半登录态（通过校验后据此重走授权签发令牌）。
type mfaChallengePayload struct {
	UserID     uint32 `json:"uid"`
	TenantID   uint32 `json:"tid"`
	Username   string `json:"username"`
	ClientID   string `json:"client_id,omitempty"`
	DeviceID   string `json:"device_id,omitempty"`
	IP         string `json:"ip,omitempty"`
	UserAgent  string `json:"user_agent,omitempty"`
	ClientType int32  `json:"client_type"`
}

func mfaChallengeKey(operationID string) string {
	return "gowind-cms:mfa:challenge:" + operationID
}

func mfaEnrollOpKey(operationID string) string {
	return "gowind-cms:mfa:enroll:" + operationID
}

func mfaRandomToken() string {
	buf := make([]byte, 32) // 256-bit
	_, _ = rand.Read(buf)
	return hex.EncodeToString(buf)
}

// issueMFAChallenge 签发一次性挑战凭据。
func (s *AuthenticationService) issueMFAChallenge(ctx context.Context, payload *mfaChallengePayload) (string, error) {
	if s.redisClient == nil {
		return "", authenticationV1.ErrorInternalServerError("mfa challenge store unavailable")
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return "", authenticationV1.ErrorInternalServerError("marshal mfa challenge failed")
	}
	operationID := mfaRandomToken()
	if err := s.redisClient.Set(ctx, mfaChallengeKey(operationID), raw, mfaChallengeTTL).Err(); err != nil {
		s.log.Errorf("store mfa challenge failed: %v", err)
		return "", authenticationV1.ErrorInternalServerError("store mfa challenge failed")
	}
	return operationID, nil
}

// consumeMFAChallenge 取出并一次性消费挑战凭据（不存在/已用过均拒绝）。
func (s *AuthenticationService) consumeMFAChallenge(ctx context.Context, operationID string) (*mfaChallengePayload, error) {
	if operationID == "" {
		return nil, authenticationV1.ErrorBadRequest("mfa operation id is required")
	}
	if s.redisClient == nil {
		return nil, authenticationV1.ErrorInternalServerError("mfa challenge store unavailable")
	}
	raw, err := s.redisClient.GetDel(ctx, mfaChallengeKey(operationID)).Result()
	if err != nil || raw == "" {
		return nil, authenticationV1.ErrorBadRequest("invalid or expired mfa challenge")
	}
	var payload mfaChallengePayload
	if err := json.Unmarshal([]byte(raw), &payload); err != nil || payload.UserID == 0 {
		return nil, authenticationV1.ErrorBadRequest("invalid mfa challenge")
	}
	return &payload, nil
}

// enforceMFAOnLogin 密码/SSO 验证通过后的两步验证闸门。
// 返回值：(challengeResponse, nil) = 用户启用 MFA，调用方直接返回挑战响应；
// (nil, nil) = 未启用 MFA，继续签发令牌；(nil, err) = 内联挑战校验失败。
func (s *AuthenticationService) enforceMFAOnLogin(ctx context.Context, req *authenticationV1.LoginRequest, user *identityV1.User) (*authenticationV1.LoginResponse, error) {
	enabled, err := s.userCredentialRepo.HasTOTPCredential(ctx, user.GetId())
	if err != nil {
		return nil, err
	}
	if !enabled {
		return nil, nil
	}

	// 第二步：挑战凭据 + 动态码随 Login 一并回传（免二次 RPC 的内联路径）
	if req.GetMfaOperationId() != "" {
		return s.verifyMFAChallengeAndMint(ctx, req.GetMfaOperationId(), req.GetTotpCode(), req)
	}

	// 第一步：身份验证已通过，签发挑战（响应不含任何令牌字段）
	payload := &mfaChallengePayload{
		UserID:     user.GetId(),
		TenantID:   user.GetTenantId(),
		Username:   user.GetUsername(),
		ClientID:   req.GetClientId(),
		DeviceID:   req.GetDeviceId(),
		IP:         req.GetClientIp(),
		UserAgent:  req.GetUserAgent(),
		ClientType: int32(req.GetClientType()),
	}
	operationID, err := s.issueMFAChallenge(ctx, payload)
	if err != nil {
		return nil, err
	}
	return &authenticationV1.LoginResponse{
		MfaOperationId: trans.Ptr(operationID),
	}, nil
}

// verifyMFAChallengeAndMint 消费挑战凭据、校验动态码，通过后签发正式令牌对。
// 挑战凭据无论校验成败均已消费——失败须重新完整登录（限速防爆破）。
func (s *AuthenticationService) verifyMFAChallengeAndMint(ctx context.Context, operationID, totpCode string, req *authenticationV1.LoginRequest) (*authenticationV1.LoginResponse, error) {
	payload, err := s.consumeMFAChallenge(ctx, operationID)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.Get(ctx, &identityV1.GetUserRequest{
		QueryBy: &identityV1.GetUserRequest_Id{Id: payload.UserID},
	})
	if err != nil {
		s.log.Errorf("mfa: get user [%d] failed [%s]", payload.UserID, err.Error())
		return nil, err
	}

	secret, err := s.userCredentialRepo.GetTOTPCredential(ctx, payload.UserID)
	if err != nil {
		s.log.Errorf("mfa: load totp factor for user [%d] failed [%s]", payload.UserID, err.Error())
		return nil, authenticationV1.ErrorBadRequest("invalid mfa code")
	}
	if !totp.Verify(secret, totpCode, time.Now()) {
		return nil, authenticationV1.ErrorBadRequest("invalid mfa code")
	}

	tokenPayload := &authenticationV1.UserTokenPayload{
		UserId:    user.GetId(),
		TenantId:  user.TenantId,
		Username:  user.Username,
		ClientId:  trans.Ptr(payload.ClientID),
		DeviceId:  trans.Ptr(payload.DeviceID),
		Ip:        trans.Ptr(payload.IP),
		UserAgent: trans.Ptr(payload.UserAgent),
	}
	// 客户端类型以挑战签发时刻为准（签发即身份已验证，不接受第二步改写）
	clientType := authenticationV1.ClientType(payload.ClientType)
	if err = s.resolveUserAuthority(ctx, user, clientType, tokenPayload); err != nil {
		return nil, err
	}

	accessToken, refreshToken, newJti, err := s.authenticator.CreateUserToken(ctx, clientType, tokenPayload)
	if err != nil {
		return nil, err
	}

	return &authenticationV1.LoginResponse{
		TokenType:        authenticationV1.TokenType_bearer,
		AccessToken:      accessToken,
		RefreshToken:     trans.Ptr(refreshToken),
		Jti:              trans.Ptr(newJti),
		ExpiresIn:        int64(s.authenticator.GetAccessTokenExpires(clientType).Seconds()),
		RefreshExpiresIn: trans.Ptr(int64(s.authenticator.GetRefreshTokenExpires(clientType).Seconds())),
	}, nil
}
