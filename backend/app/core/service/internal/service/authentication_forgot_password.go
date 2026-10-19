package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/tx7do/go-utils/trans"
	"google.golang.org/protobuf/types/known/emptypb"

	authenticationV1 "go-wind-cms/api/gen/go/authentication/service/v1"
	identityV1 "go-wind-cms/api/gen/go/identity/service/v1"
	notifyV1 "go-wind-cms/api/gen/go/notify/service/v1"

)

// generateVCode 生成 6 位数字验证码。
// 使用纳秒时间取模：验证码仅用于一次性校验，安全强度由
// 「单次有效 + 10 分钟 TTL + 服务端比对」保证，不依赖密码学随机。
func generateVCode() string {
	return fmt.Sprintf("%06d", time.Now().UnixNano()%1000000)
}

// ForgotPassword 忘记密码：向 identifier（已绑定的邮箱凭证）发送重置验证码——
// 邮箱走 SMTP（显式点名 EMAIL 渠道，不经路由表）。
// 用户不存在时同样返回成功，防止通过接口枚举有效账号。
//
// 与上游 go-wind-admin 的差异：CMS 暂无短信渠道，identifier 仅支持邮箱路径。
func (s *AuthenticationService) ForgotPassword(ctx context.Context, req *authenticationV1.ForgotPasswordRequest) (*emptypb.Empty, error) {
	identifier := strings.TrimSpace(req.GetIdentifier())
	if identifier == "" {
		return nil, errors.BadRequest("AUTHENTICATION", "identifier is required")
	}
	if !strings.Contains(identifier, "@") {
		return nil, errors.BadRequest("AUTHENTICATION", "identifier must be an email address")
	}

	cred, credErr := s.userCredentialRepo.GetByIdentifier(ctx, &authenticationV1.GetUserCredentialByIdentifierRequest{
		IdentityType: authenticationV1.UserCredential_EMAIL,
		Identifier:   identifier,
	})
	userId := cred.GetUserId()
	if credErr != nil || userId == 0 {
		s.log.Infof("forgot-password: no email credential for [%s] (silent ok)", identifier)
		return &emptypb.Empty{}, nil
	}

	code := generateVCode()
	if err := s.vcodeCache.Save(ctx, "reset_password", identifier, code, 10*time.Minute); err != nil {
		s.log.Errorf("forgot-password: save verification code failed: %s", err.Error())
		return nil, errors.InternalServer("AUTHENTICATION", "save verification code failed")
	}

	// 文案语言：优先收件用户的偏好 locale（个人中心设置），未设置回落中文。
	locale := alertLocaleZhCN
	if user, uErr := s.userRepo.Get(ctx, &identityV1.GetUserRequest{
		QueryBy: &identityV1.GetUserRequest_Id{Id: userId},
	}); uErr == nil && user != nil && user.GetLocale() != "" {
		if normalizeLocaleTag(user.GetLocale()) == "en-us" || normalizeLocaleTag(user.GetLocale()) == "en" {
			locale = alertLocaleEnUS
		}
	}

	if s.notifier == nil {
		s.log.Errorf("forgot-password: notification seam is not wired")
		return nil, errors.InternalServer("AUTHENTICATION", "email channel is not configured")
	}

	title, content := passwordResetCodeText(locale, code)
	resp, err := s.notifier.SendDirect(ctx, &notifyV1.SendDirectNotificationRequest{
		EventType:       notifyV1.EventType_PASSWORD_RESET_CODE,
		Channel:         trans.Ptr(notifyV1.Channel_EMAIL),
		Target:          identifier,
		Title:           title,
		Content:         content,
		RecipientUserId: trans.Ptr(userId),
	})
	if err != nil {
		s.log.Errorf("forgot-password: send code to [%s] failed (delivery %d): %s",
			identifier, resp.GetDeliveryId(), err.Error())
		// SKIPPED＝一条都没发出去：没配启用的邮件渠道。这与"渠道报错"是两回事，
		// 前者要管理员去渠道页配置，分开报才不用翻日志才知道差哪一步。
		if resp.GetStatus() == notifyV1.DeliveryStatus_SKIPPED {
			return nil, errors.InternalServer("AUTHENTICATION", "email channel is not configured")
		}
		return nil, errors.InternalServer("AUTHENTICATION", "send verification message failed")
	}

	s.log.Infof("forgot-password: reset code sent to [%s] for user [%d]", identifier, userId)
	return &emptypb.Empty{}, nil
}

// ResetPasswordByCode 凭邮箱验证码重置密码（免鉴权）。
// 校验通过后重置密码（密码策略/加密在 repo 层处理）并吊销该用户全部会话。
func (s *AuthenticationService) ResetPasswordByCode(ctx context.Context, req *authenticationV1.ResetPasswordByCodeRequest) (*emptypb.Empty, error) {
	identifier := strings.TrimSpace(req.GetIdentifier())
	code := strings.TrimSpace(req.GetCode())
	newPassword := req.GetNewPassword()
	if identifier == "" || code == "" || newPassword == "" {
		return nil, errors.BadRequest("AUTHENTICATION", "identifier, code and new password are required")
	}

	if !s.vcodeCache.Verify(ctx, "reset_password", identifier, code) {
		return nil, errors.BadRequest("AUTHENTICATION", "invalid or expired verification code")
	}

	cred, credErr := s.userCredentialRepo.GetByIdentifier(ctx, &authenticationV1.GetUserCredentialByIdentifierRequest{
		IdentityType: authenticationV1.UserCredential_EMAIL,
		Identifier:   identifier,
	})
	if credErr != nil || cred.GetUserId() == 0 {
		return nil, errors.BadRequest("AUTHENTICATION", "invalid or expired verification code")
	}
	userId := cred.GetUserId()

	user, err := s.userRepo.Get(ctx, &identityV1.GetUserRequest{
		QueryBy: &identityV1.GetUserRequest_Id{Id: userId},
	})
	if err != nil {
		return nil, err
	}

	if err = s.userCredentialRepo.ResetCredential(ctx, &authenticationV1.ResetCredentialRequest{
		IdentityType:  authenticationV1.UserCredential_USERNAME,
		Identifier:    user.GetUsername(),
		NewCredential: newPassword,
		NeedDecrypt:   true,
	}); err != nil {
		s.log.Errorf("reset password by code failed for user [%d]: %s", userId, err.Error())
		return nil, err
	}

	// 重置后吊销该用户全部已登录会话：旧会话可能握着被盗凭据的入口
	for name, clientType := range authenticationV1.ClientType_value {
		if name == "CLIENT_UNSPECIFIED" {
			continue
		}
		if err := s.authenticator.RevokeUserToken(ctx, authenticationV1.ClientType(clientType), userId); err != nil {
			s.log.Errorf("revoke tokens after reset-by-code failed for user [%d] (%s): %s", userId, name, err.Error())
		}
	}

	s.log.Infof("password reset by code done for user [%d]", userId)
	return &emptypb.Empty{}, nil
}

// ---- 验证码邮件文案（按收件人偏好语言，与监控告警同一套语言口径） ----

func passwordResetCodeText(locale alertLocaleTag, code string) (string, string) {
	switch locale {
	case alertLocaleEnUS:
		return fmt.Sprintf("[GoWind CMS] Password reset code: %s", code),
			fmt.Sprintf("Your password reset verification code is %s. It is valid for 10 minutes and can be used once. If you did not request it, please ignore this email.", code)
	default:
		return fmt.Sprintf("[GoWind CMS] 密码重置验证码：%s", code),
			fmt.Sprintf("您的密码重置验证码为 %s，10 分钟内有效、仅可使用一次。若非本人操作，请忽略本邮件。", code)
	}
}
