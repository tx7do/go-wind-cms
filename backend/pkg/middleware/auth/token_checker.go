package auth

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	authenticationV1 "go-wind-cms/api/gen/go/authentication/service/v1"
)

type TokenChecker struct {
	log *log.Helper

	authenticationServiceClient authenticationV1.AuthenticationServiceClient
	clientType                  authenticationV1.ClientType
}

func NewTokenChecker(
	ctx *bootstrap.Context,
	authenticationServiceClient authenticationV1.AuthenticationServiceClient,
	clientType authenticationV1.ClientType,
) AccessTokenChecker {
	return &TokenChecker{
		log:                         log.NewHelper(log.With(ctx.GetLogger(), "module", "token-checker/auth/middleware")),
		authenticationServiceClient: authenticationServiceClient,
		clientType:                  clientType,
	}
}

// IsValidAccessToken checks if the access token is valid for the given user ID.
func (tc *TokenChecker) IsValidAccessToken(ctx context.Context, accessToken string, skipRedis bool) (bool, *authenticationV1.UserTokenPayload) {
	resp, err := tc.authenticationServiceClient.ValidateToken(ctx, &authenticationV1.ValidateTokenRequest{
		Token:         accessToken,
		TokenCategory: authenticationV1.TokenCategory_ACCESS,
		ClientType:    tc.clientType,
		SkipRedis:     trans.Ptr(skipRedis),
	})
	if err != nil {
		return false, nil
	}

	if !resp.IsValid {
		return false, nil
	}

	return true, resp.Payload
}

// IsBlockedAccessToken checks if the access token is blocked for the given user ID.
func (tc *TokenChecker) IsBlockedAccessToken(ctx context.Context, accessToken string) bool {
	resp, err := tc.authenticationServiceClient.ValidateToken(ctx, &authenticationV1.ValidateTokenRequest{
		Token:         accessToken,
		TokenCategory: authenticationV1.TokenCategory_ACCESS,
		ClientType:    tc.clientType,
		SkipRedis:     trans.Ptr(true),
	})
	if err != nil {
		return true
	}
	return !resp.IsValid
}
