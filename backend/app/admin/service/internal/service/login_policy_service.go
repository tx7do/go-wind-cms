package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	authenticationV1 "go-wind-cms/api/gen/go/authentication/service/v1"

	"go-wind-cms/pkg/middleware/auth"
)

type LoginPolicyService struct {
	adminV1.LoginPolicyServiceHTTPServer

	log *log.Helper

	loginPolicyServiceClient authenticationV1.LoginPolicyServiceClient
}

func NewLoginPolicyService(ctx *bootstrap.Context, loginPolicyServiceClient authenticationV1.LoginPolicyServiceClient) *LoginPolicyService {
	return &LoginPolicyService{
		log:                      ctx.NewLoggerHelper("admlogin-policy/service/admin-service"),
		loginPolicyServiceClient: loginPolicyServiceClient,
	}
}

func (s *LoginPolicyService) List(ctx context.Context, req *paginationV1.PagingRequest) (*authenticationV1.ListLoginPolicyResponse, error) {
	return s.loginPolicyServiceClient.List(ctx, req)
}

func (s *LoginPolicyService) Get(ctx context.Context, req *authenticationV1.GetLoginPolicyRequest) (*authenticationV1.LoginPolicy, error) {
	return s.loginPolicyServiceClient.Get(ctx, req)
}

func (s *LoginPolicyService) Create(ctx context.Context, req *authenticationV1.CreateLoginPolicyRequest) (*emptypb.Empty, error) {
	if req == nil || req.Data == nil {
		return nil, adminV1.ErrorBadRequest("invalid request")
	}

	// 获取操作人信息
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.Data.CreatedBy = trans.Ptr(operator.UserId)

	return s.loginPolicyServiceClient.Create(ctx, req)
}

func (s *LoginPolicyService) Update(ctx context.Context, req *authenticationV1.UpdateLoginPolicyRequest) (*emptypb.Empty, error) {
	if req == nil || req.Data == nil {
		return nil, adminV1.ErrorBadRequest("invalid request")
	}

	// 获取操作人信息
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.Data.UpdatedBy = trans.Ptr(operator.GetUserId())
	if req.UpdateMask != nil {
		req.UpdateMask.Paths = append(req.UpdateMask.Paths, "updated_by")
	}

	return s.loginPolicyServiceClient.Update(ctx, req)
}

func (s *LoginPolicyService) Delete(ctx context.Context, req *authenticationV1.DeleteLoginPolicyRequest) (*emptypb.Empty, error) {
	if req == nil {
		return nil, adminV1.ErrorBadRequest("invalid request")
	}

	return s.loginPolicyServiceClient.Delete(ctx, req)
}
