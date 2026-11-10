package netutil

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
)

// ParseUnverifiedBearerJWTClaims 从请求头解析 access token(可能是已过期的)
// payload,取 uid 与 jti。
//
// 安全性:payload 不做签名与时效校验,仅作为"刷新令牌绑定键"的来源;
// 真正的凭证校验由 core 侧 VerifyAndRevokeTokenPair 以 refresh token 值
// 完成(随机秘密,伪造 uid/jti 无法通过)。
func ParseUnverifiedBearerJWTClaims(ctx context.Context) (uid int64, jti string, err error) {
	header := HeaderFromContext(ctx)
	if header == nil {
		return 0, "", errors.New("no http header in context")
	}
	authz := header.Get("Authorization")
	token, ok := strings.CutPrefix(authz, "Bearer ")
	if !ok || token == "" {
		return 0, "", errors.New("no bearer token")
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return 0, "", errors.New("invalid jwt format")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return 0, "", errors.New("invalid jwt payload encoding")
	}

	var claims struct {
		UID int64  `json:"uid"`
		Jti string `json:"jti"`
	}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return 0, "", errors.New("invalid jwt claims")
	}
	if claims.UID == 0 || claims.Jti == "" {
		return 0, "", errors.New("missing uid or jti")
	}
	return claims.UID, claims.Jti, nil
}
