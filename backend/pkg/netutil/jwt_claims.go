package netutil

import (
	"context"
	"errors"
	"strings"

	jwtutil "github.com/tx7do/go-utils/jwtutil"
)

// ParseUnverifiedBearerJWTClaims 从请求头解析 access token(可能是已过期的)
// payload,取 uid 与 jti。
//
// 安全性:payload 不做签名与时效校验,仅作为"刷新令牌绑定键"的来源;
// 真正的凭证校验由 core 侧 VerifyAndRevokeTokenPair 以 refresh token 值
// 完成(随机秘密,伪造 uid/jti 无法通过)。
// payload 解码委托 jwtutil.ParseJWTPayload（golang-jwt 不验签解 claims）。
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

	claims, err := jwtutil.ParseJWTPayload(token)
	if err != nil {
		return 0, "", errors.New("invalid jwt payload")
	}

	rawUID, _ := claims["uid"].(float64)
	uid = int64(rawUID)
	jti, _ = claims["jti"].(string)
	if uid == 0 || jti == "" {
		return 0, "", errors.New("missing uid or jti")
	}
	return uid, jti, nil
}
