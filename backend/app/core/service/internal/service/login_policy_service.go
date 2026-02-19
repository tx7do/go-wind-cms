package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	"go-wind-cms/app/core/service/internal/data"

	authenticationV1 "go-wind-cms/api/gen/go/authentication/service/v1"

	"go-wind-cms/pkg/middleware/auth"
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

	// 获取操作人信息
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.Data.CreatedBy = trans.Ptr(operator.UserId)

	if err = s.loginPolicyRepo.Create(ctx, req); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *LoginPolicyService) Update(ctx context.Context, req *authenticationV1.UpdateLoginPolicyRequest) (*emptypb.Empty, error) {
	if req == nil || req.Data == nil {
		return nil, authenticationV1.ErrorBadRequest("invalid request")
	}

	// 获取操作人信息
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.Data.UpdatedBy = trans.Ptr(operator.UserId)
	if req.UpdateMask != nil {
		req.UpdateMask.Paths = append(req.UpdateMask.Paths, "updated_by")
	}

	if err = s.loginPolicyRepo.Update(ctx, req); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
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
