package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/tx7do/go-utils/trans"

	authenticationV1 "go-wind-cms/api/gen/go/authentication/service/v1"
	identityV1 "go-wind-cms/api/gen/go/identity/service/v1"
	permissionV1 "go-wind-cms/api/gen/go/permission/service/v1"
)

// OIDC 环境变量配置（授权码模式，core 服务读取）：
//   - OIDC_ISSUER          Issuer 标识（如 https://idp.example.com），用于用户映射
//   - OIDC_CLIENT_ID       客户端 ID
//   - OIDC_CLIENT_SECRET   客户端密钥
//   - OIDC_AUTHORIZE_URL   授权端点（浏览器跳转取 code）
//   - OIDC_TOKEN_URL       令牌端点（code 换 access_token）
//   - OIDC_USERINFO_URL    用户信息端点（取 sub/email/preferred_username）
//   - OIDC_SCOPES          授权 scope，空格分隔（默认 "openid profile email"）
//   - OIDC_AUTO_PROVISION  首次登录自动建用户（默认 true；false 时未映射用户拒绝登录）
//   - OIDC_DEFAULT_ROLE_CODE  首登自动建号时绑定的默认角色 code（可选，平台租户内解析）；
//     未设置则建号不带角色——授权链会在登录时以 insufficient authority 拒绝，
//     运维需配置该角色或预先做用户映射
//
// 任一必填缺失即视为 OIDC 未启用：GetOidcAuthorizeUrl 返回未配置错误，
// 授权码登录同样拒绝——前端据此隐藏 SSO 入口，不影响密码登录。
//
// 用户映射：按 (issuer, sub) 查 ENTERPRISE_SSO 凭证（identifier=issuer|sub），
// 命中则登录既有用户；未命中且 AUTO_PROVISION 开启则按 userinfo 建用户+凭证。
// 登录成功签发本仓 JWT，与密码登录同一令牌链（admin 侧 refresh token 经
// HttpOnly Cookie 下发的门面逻辑不变）。

// oidcStateTTL 授权 state 的暂存时长：覆盖人 completing IdP 登录的常见耗时。
const oidcStateTTL = 10 * time.Minute

// oidcConfig OIDC 授权码模式所需的最小配置集。
type oidcConfig struct {
	Issuer       string
	ClientID     string
	ClientSecret string
	AuthorizeURL string
	TokenURL     string
	UserinfoURL  string
	Scopes       string
}

// loadOidcConfig 读取环境变量；未配置（任一必填为空）返回 false。
func loadOidcConfig(issuerOverride string) (*oidcConfig, bool) {
	cfg := &oidcConfig{
		Issuer:       strings.TrimSpace(firstNonEmptyStr(issuerOverride, os.Getenv("OIDC_ISSUER"))),
		ClientID:     strings.TrimSpace(os.Getenv("OIDC_CLIENT_ID")),
		ClientSecret: os.Getenv("OIDC_CLIENT_SECRET"),
		AuthorizeURL: strings.TrimSpace(os.Getenv("OIDC_AUTHORIZE_URL")),
		TokenURL:     strings.TrimSpace(os.Getenv("OIDC_TOKEN_URL")),
		UserinfoURL:  strings.TrimSpace(os.Getenv("OIDC_USERINFO_URL")),
		Scopes:       strings.TrimSpace(firstNonEmptyStr(os.Getenv("OIDC_SCOPES"), "openid profile email")),
	}
	if cfg.Issuer == "" || cfg.ClientID == "" || cfg.ClientSecret == "" ||
		cfg.AuthorizeURL == "" || cfg.TokenURL == "" || cfg.UserinfoURL == "" {
		return nil, false
	}
	return cfg, true
}

// oidcTokenResponse IdP 令牌端点响应（取子集）。
type oidcTokenResponse struct {
	AccessToken string `json:"access_token"`
	IDToken     string `json:"id_token"`
	TokenType   string `json:"token_type"`
}

// oidcUserInfo IdP 用户信息端点响应（取子集）。
type oidcUserInfo struct {
	Sub               string `json:"sub"`
	Email             string `json:"email"`
	PreferredUsername string `json:"preferred_username"`
	Name              string `json:"name"`
}

// GetOidcAuthorizeUrl 生成 IdP 授权跳转地址：服务端生成随机 state 并暂存
// Redis（一次性、10 分钟过期），前端拿到 authorize_url 后整页跳转 IdP。
// IdP 回调携带 code+state，前端再以 GrantType_authorization_code 调 Login。
func (s *AuthenticationService) GetOidcAuthorizeUrl(ctx context.Context, req *authenticationV1.GetOidcAuthorizeUrlRequest) (*authenticationV1.GetOidcAuthorizeUrlResponse, error) {
	// 授权跳转入口不支持多 issuer 选择（单 IdP，issuer 由环境变量决定）
	cfg, ok := loadOidcConfig("")
	if !ok {
		return nil, authenticationV1.ErrorInternalServerError("oidc login is not configured")
	}
	if req.GetRedirectUri() == "" {
		return nil, authenticationV1.ErrorBadRequest("redirect_uri is required")
	}
	if s.redisClient == nil {
		return nil, authenticationV1.ErrorInternalServerError("oidc state store unavailable")
	}

	state, err := randomState()
	if err != nil {
		s.log.Errorf("oidc: generate state failed: %v", err)
		return nil, authenticationV1.ErrorInternalServerError("generate state failed")
	}
	if err = s.redisClient.Set(ctx, oidcStateKey(state), "1", oidcStateTTL).Err(); err != nil {
		s.log.Errorf("oidc: store state failed: %v", err)
		return nil, authenticationV1.ErrorInternalServerError("store state failed")
	}

	authorizeURL := fmt.Sprintf("%s?response_type=code&client_id=%s&redirect_uri=%s&scope=%s&state=%s",
		cfg.AuthorizeURL,
		url.QueryEscape(cfg.ClientID),
		url.QueryEscape(req.GetRedirectUri()),
		url.QueryEscape(cfg.Scopes),
		url.QueryEscape(state),
	)
	return &authenticationV1.GetOidcAuthorizeUrlResponse{
		AuthorizeUrl: trans.Ptr(authorizeURL),
	}, nil
}

// doGrantTypeOidcAuthorizationCode OIDC 授权码登录：state 校验 → code 换
// access_token → userinfo → 用户映射/自动建号 → 签发本仓 JWT。
func (s *AuthenticationService) doGrantTypeOidcAuthorizationCode(ctx context.Context, req *authenticationV1.LoginRequest) (*authenticationV1.LoginResponse, error) {
	cfg, ok := loadOidcConfig(req.GetOidcIssuer())
	if !ok {
		return nil, authenticationV1.ErrorInternalServerError("oidc login is not configured")
	}
	if req.GetCode() == "" {
		return nil, authenticationV1.ErrorBadRequest("oidc code is required")
	}

	// ===== state 防 CSRF：须为本服务签发且未消费（Redis 一次性取出）=====
	if err := s.oidcConsumeState(ctx, req.GetOidcState()); err != nil {
		return nil, err
	}

	// ===== code 换 access_token =====
	accessToken, err := s.oidcExchangeCode(ctx, cfg, req.GetCode(), req.GetRedirectUri())
	if err != nil {
		s.log.Errorf("oidc: exchange code failed: %v", err)
		return nil, authenticationV1.ErrorBadRequest("oidc code exchange failed")
	}

	// ===== userinfo 取身份 =====
	userinfo, err := s.oidcFetchUserinfo(ctx, cfg.UserinfoURL, accessToken)
	if err != nil {
		s.log.Errorf("oidc: fetch userinfo failed: %v", err)
		return nil, authenticationV1.ErrorBadRequest("oidc userinfo failed")
	}
	if userinfo.Sub == "" {
		return nil, authenticationV1.ErrorBadRequest("oidc userinfo missing sub")
	}

	// ===== 用户映射：ENTERPRISE_SSO 凭证 identifier = issuer|sub =====
	// GetByIdentifier 未命中返回通用 404（reason=NOT_FOUND），按 kratos 的
	// 404 判定（而非具体 reason）识别"未映射"，避免遗漏。
	identifier := cfg.Issuer + "|" + userinfo.Sub
	credential, credErr := s.userCredentialRepo.GetByIdentifier(ctx, &authenticationV1.GetUserCredentialByIdentifierRequest{
		IdentityType: authenticationV1.UserCredential_ENTERPRISE_SSO,
		Identifier:   identifier,
	})

	var userID uint32
	switch {
	case credErr == nil && credential != nil && credential.GetUserId() != 0:
		userID = credential.GetUserId()
	case kerrors.IsNotFound(credErr) || (credential != nil && credential.GetUserId() == 0):
		// 未映射：按 AUTO_PROVISION 决定是否自动建号
		if strings.EqualFold(strings.TrimSpace(os.Getenv("OIDC_AUTO_PROVISION")), "false") {
			return nil, authenticationV1.ErrorForbidden("oidc user not provisioned")
		}
		userID, err = s.oidcProvisionUser(ctx, cfg.Issuer, userinfo)
		if err != nil {
			s.log.Errorf("oidc: provision user failed: %v", err)
			return nil, err
		}
	default:
		s.log.Errorf("oidc: lookup credential failed: %v", credErr)
		return nil, credErr
	}

	user, err := s.userRepo.Get(ctx, &identityV1.GetUserRequest{
		QueryBy: &identityV1.GetUserRequest_Id{Id: userID},
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
	if err = s.resolveUserAuthority(ctx, user, req.GetClientType(), tokenPayload); err != nil {
		return nil, err
	}

	accessTokenJWT, refreshToken, _, err := s.authenticator.CreateUserToken(ctx, req.GetClientType(), tokenPayload)
	if err != nil {
		return nil, err
	}

	return &authenticationV1.LoginResponse{
		TokenType:        authenticationV1.TokenType_bearer,
		AccessToken:      accessTokenJWT,
		RefreshToken:     trans.Ptr(refreshToken),
		ExpiresIn:        int64(s.authenticator.GetAccessTokenExpires(req.GetClientType()).Seconds()),
		RefreshExpiresIn: trans.Ptr(int64(s.authenticator.GetRefreshTokenExpires(req.GetClientType()).Seconds())),
	}, nil
}

// oidcStateKey 授权 state 的 Redis 键。
func oidcStateKey(state string) string {
	return "gowind-cms:oidc:state:" + state
}

// oidcConsumeState 校验并一次性消费 state（不存在/已用过/为空均拒绝）。
// Redis 不可用按配置错误拒绝——跳过校验等于放弃 CSRF 防护。
func (s *AuthenticationService) oidcConsumeState(ctx context.Context, state string) error {
	if state == "" {
		return authenticationV1.ErrorBadRequest("oidc state is required")
	}
	if s.redisClient == nil {
		return authenticationV1.ErrorInternalServerError("oidc state store unavailable")
	}
	stored, err := s.redisClient.Get(ctx, oidcStateKey(state)).Result()
	if err != nil || stored == "" {
		return authenticationV1.ErrorBadRequest("invalid oidc state")
	}
	s.redisClient.Del(ctx, oidcStateKey(state))
	return nil
}

// randomState 生成 256-bit 随机 state（hex）。
func randomState() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

// oidcExchangeCode 授权码换令牌（form-urlencoded POST）。
func (s *AuthenticationService) oidcExchangeCode(ctx context.Context, cfg *oidcConfig, code, redirectURI string) (string, error) {
	form := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"client_id":     {cfg.ClientID},
		"client_secret": {cfg.ClientSecret},
	}
	if redirectURI != "" {
		form.Set("redirect_uri", redirectURI)
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, cfg.TokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	httpReq.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token endpoint status %d: %s", resp.StatusCode, string(body[:min(len(body), 200)]))
	}
	var tok oidcTokenResponse
	if err := json.Unmarshal(body, &tok); err != nil {
		return "", err
	}
	if tok.AccessToken == "" {
		return "", fmt.Errorf("token endpoint returned empty access_token")
	}
	return tok.AccessToken, nil
}

// oidcFetchUserinfo Bearer 调 userinfo 端点。
func (s *AuthenticationService) oidcFetchUserinfo(ctx context.Context, userinfoURL, accessToken string) (*oidcUserInfo, error) {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, userinfoURL, nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+accessToken)
	httpReq.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("userinfo endpoint status %d", resp.StatusCode)
	}
	var info oidcUserInfo
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// oidcProvisionUser 首次 OIDC 登录自动建号：用户 + ENTERPRISE_SSO 凭证。
// 用户名取 preferred_username（冲突追加 sub 后缀保证唯一）；平台租户（tenant 0）。
func (s *AuthenticationService) oidcProvisionUser(ctx context.Context, issuer string, userinfo *oidcUserInfo) (uint32, error) {
	username := firstNonEmptyStr(userinfo.PreferredUsername, userinfo.Email, userinfo.Sub)
	// 用户名冲突：追加 OIDC sub 后缀保证唯一
	if exist, err := s.userRepo.Get(ctx, &identityV1.GetUserRequest{
		QueryBy: &identityV1.GetUserRequest_Username{Username: username},
	}); err == nil && exist != nil && exist.GetId() != 0 {
		username = fmt.Sprintf("%s_%s", username, userinfo.Sub)
	}

	created, err := s.userRepo.Create(ctx, &identityV1.CreateUserRequest{
		Data: &identityV1.User{
			TenantId: trans.Ptr(uint32(0)),
			Username: trans.Ptr(username),
			Nickname: trans.Ptr(firstNonEmptyStr(userinfo.Name, username)),
			Email:    trans.Ptr(userinfo.Email),
			Status:   identityV1.User_NORMAL.Enum(),
		},
	})
	if err != nil {
		return 0, err
	}
	userID := created.GetId()

	if err = s.userCredentialRepo.Create(ctx, &authenticationV1.CreateUserCredentialRequest{
		Data: &authenticationV1.UserCredential{
			TenantId:       trans.Ptr(uint32(0)),
			UserId:         trans.Ptr(userID),
			IdentityType:   authenticationV1.UserCredential_ENTERPRISE_SSO.Enum(),
			Identifier:     trans.Ptr(issuer + "|" + userinfo.Sub),
			CredentialType: authenticationV1.UserCredential_OPENID_CONNECT_ID_TOKEN.Enum(),
			Credential:     trans.Ptr("oidc"),
			IsPrimary:      trans.Ptr(true),
			Status:         authenticationV1.UserCredential_ENABLED.Enum(),
		},
	}); err != nil {
		return 0, err
	}

	// 默认角色绑定（可选）：OIDC_DEFAULT_ROLE_CODE 指定时在平台级角色
	// （tenant_id IS NULL）内解析并绑定；未设置/解析失败不阻断建号
	//（登录会被授权链拒绝，角色可后续手工补绑）。
	if roleCode := strings.TrimSpace(os.Getenv("OIDC_DEFAULT_ROLE_CODE")); roleCode != "" {
		role, rerr := s.roleRepo.GetPlatformRoleByCode(ctx, roleCode)
		if rerr != nil || role == nil || role.GetId() == 0 {
			s.log.Warnf("oidc: default role [%s] not found in platform tenant, provisioned user [%d] has no roles", roleCode, userID)
		} else if aerr := s.userRepo.AssignUserRole(ctx, &permissionV1.UserRole{
			UserId:    trans.Ptr(userID),
			TenantId:  trans.Ptr(uint32(0)),
			RoleId:    trans.Ptr(role.GetId()),
			IsPrimary: trans.Ptr(true),
			Status:    permissionV1.UserRole_ACTIVE.Enum(),
		}); aerr != nil {
			s.log.Errorf("oidc: assign default role [%s] to user [%d] failed: %v", roleCode, userID, aerr)
			return 0, aerr
		}
	}
	return userID, nil
}

func firstNonEmptyStr(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}
