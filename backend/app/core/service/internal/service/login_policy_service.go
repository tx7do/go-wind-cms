package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/emptypb"

	"go-wind-cms/app/core/service/internal/data"

	authenticationV1 "go-wind-cms/api/gen/go/authentication/service/v1"
)

type LoginPolicyService struct {
	authenticationV1.UnimplementedLoginPolicyServiceServer

	log *log.Helper

	loginPolicyRepo *data.LoginPolicyRepo
}

func NewLoginPolicyService(ctx *bootstrap.Context, repo *data.LoginPolicyRepo) *LoginPolicyService {
	return &LoginPolicyService{
		log:             ctx.NewLoggerHelper("login-policy/service/core-service"),
		loginPolicyRepo: repo,
	}
}

func (s *LoginPolicyService) List(ctx context.Context, req *paginationV1.PagingRequest) (*authenticationV1.ListLoginPolicyResponse, error) {
	return s.loginPolicyRepo.List(ctx, req)
}

func (s *LoginPolicyService) Get(ctx context.Context, req *authenticationV1.GetLoginPolicyRequest) (*authenticationV1.LoginPolicy, error) {
	return s.loginPolicyRepo.Get(ctx, req)
}

func (s *LoginPolicyService) Create(ctx context.Context, req *authenticationV1.CreateLoginPolicyRequest) (*emptypb.Empty, error) {
	if req == nil || req.Data == nil {
		return nil, authenticationV1.ErrorBadRequest("invalid request")
	}

	// 值格式校验：配错的值会让匹配器静默不命中——尤其白名单，存在但永不命中 = 全员被锁
	if err := s.validatePolicyValue(req.Data.GetMethod(), req.Data.GetValue()); err != nil {
		return nil, err
	}

	if err := s.loginPolicyRepo.Create(ctx, req); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *LoginPolicyService) Update(ctx context.Context, req *authenticationV1.UpdateLoginPolicyRequest) (*emptypb.Empty, error) {
	if req == nil || req.Data == nil {
		return nil, authenticationV1.ErrorBadRequest("invalid request")
	}

	// 掩码更新只在动 value 时校验；method 未随请求携带时按存量行的方法判定
	//（掩码更新常只带 value/reason，method 为零值——直接校验会拿 UNSPECIFIED 误炸）。
	if policyValueInMask(req.GetUpdateMask()) {
		method := req.Data.GetMethod()
		if method == authenticationV1.LoginPolicy_LOGIN_RESTRICTION_METHOD_UNSPECIFIED {
			old, err := s.loginPolicyRepo.Get(ctx, &authenticationV1.GetLoginPolicyRequest{
				QueryBy: &authenticationV1.GetLoginPolicyRequest_Id{Id: req.GetId()},
			})
			if err != nil {
				return nil, err
			}
			method, _ = old.GetMethod(), old.GetValue()
		}
		if err := s.validatePolicyValue(method, req.Data.GetValue()); err != nil {
			return nil, err
		}
	}

	if err := s.loginPolicyRepo.Update(ctx, req); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// validatePolicyValue 统一的策略值格式校验入口（method 零值即拒绝）。
func (s *LoginPolicyService) validatePolicyValue(method authenticationV1.LoginPolicy_Method, value string) error {
	if method == authenticationV1.LoginPolicy_LOGIN_RESTRICTION_METHOD_UNSPECIFIED {
		return authenticationV1.ErrorBadRequest("policy method is required")
	}
	if err := data.ValidateLoginPolicyValue(method.String(), value); err != nil {
		return authenticationV1.ErrorBadRequest("invalid policy value: %s", err.Error())
	}
	return nil
}

// policyValueInMask 判断更新掩码是否包含 value 字段（nil 掩码视为全字段）。
func policyValueInMask(mask *fieldmaskpb.FieldMask) bool {
	if mask == nil {
		return true
	}
	for _, p := range mask.GetPaths() {
		if p == "value" {
			return true
		}
	}
	return false
}

func (s *LoginPolicyService) Delete(ctx context.Context, req *authenticationV1.DeleteLoginPolicyRequest) (*emptypb.Empty, error) {
	if req == nil {
		return nil, authenticationV1.ErrorBadRequest("invalid request")
	}

	if err := s.loginPolicyRepo.Delete(ctx, req); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
