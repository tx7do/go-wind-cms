package server

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"github.com/tx7do/kratos-bootstrap/transport/sse"

	sseServer "github.com/tx7do/kratos-transport/transport/sse"

	authenticationV1 "go-wind-cms/api/gen/go/authentication/service/v1"
)

// NewSseServer creates a new SSE server.
func NewSseServer(
	ctx *bootstrap.Context,
	authenticationServiceClient authenticationV1.AuthenticationServiceClient,
) *sseServer.Server {
	cfg := ctx.GetConfig()

	if cfg == nil || cfg.Server == nil || cfg.Server.Sse == nil {
		return nil
	}

	l := ctx.NewLoggerHelper("sse-server/app-service")

	srv := sse.NewSseServer(cfg.Server.Sse,
		sseServer.WithSubscriberFunction(func(streamID sseServer.StreamID, sub *sseServer.Subscriber) {
			l.Infof("subscriber [%s] connected", streamID)
		}),
		sseServer.WithAuthorizeFunc(func(r *http.Request, token string) error {
			// SSE 订阅必须携带有效 access token，拒绝匿名连接。
			resp, err := authenticationServiceClient.ValidateToken(context.Background(), &authenticationV1.ValidateTokenRequest{
				ClientType:    authenticationV1.ClientType_app,
				Token:         token,
				TokenCategory: authenticationV1.TokenCategory_ACCESS,
			})
			if err != nil {
				log.Errorf("app sse token authentication failed: %s", err)
				return err
			}
			tokenHash := hashToken(token)
			if resp.GetIsBlocked() {
				log.Warnf("app sse token is blocked: %s", tokenHash)
				return authenticationV1.ErrorForbidden("token is blocked")
			}
			if !resp.GetIsValid() {
				log.Warnf("app sse token is invalid: %s", tokenHash)
				return authenticationV1.ErrorUnauthorized("invalid token")
			}
			return nil
		}),
	)

	return srv
}

// hashToken returns a short SHA-256 fingerprint of a token suitable for
// log correlation without exposing the raw credential.
func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])[:16]
}
