package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	socialV1 "go-wind-cms/api/gen/go/social/service/v1"

	"go-wind-cms/pkg/middleware/auth"
)

// SocialService 社交分发管理面（BFF 薄转发）。
// 凭据的验证/加解密在 core service 完成，这里只补操作人并透传。
type SocialService struct {
	adminV1.SocialServiceHTTPServer

	socialServiceClient socialV1.SocialServiceClient
	log                 *log.Helper
}

func NewSocialService(ctx *bootstrap.Context, socialServiceClient socialV1.SocialServiceClient) *SocialService {
	return &SocialService{
		log:                 ctx.NewLoggerHelper("social/service/admin-service"),
		socialServiceClient: socialServiceClient,
	}
}

// ── 社交账号 ──

func (s *SocialService) ListSocialAccounts(ctx context.Context, req *paginationV1.PagingRequest) (*socialV1.ListSocialAccountsResponse, error) {
	return s.socialServiceClient.ListSocialAccounts(ctx, req)
}

func (s *SocialService) GetSocialAccount(ctx context.Context, req *socialV1.GetSocialAccountRequest) (*socialV1.SocialAccount, error) {
	return s.socialServiceClient.GetSocialAccount(ctx, req)
}

func (s *SocialService) CreateSocialAccount(ctx context.Context, req *socialV1.CreateSocialAccountRequest) (*socialV1.SocialAccount, error) {
	if req == nil || req.Data == nil {
		return nil, adminV1.ErrorBadRequest("invalid parameter")
	}

	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.Data.CreatedBy = trans.Ptr(operator.UserId)

	return s.socialServiceClient.CreateSocialAccount(ctx, req)
}

func (s *SocialService) UpdateSocialAccount(ctx context.Context, req *socialV1.UpdateSocialAccountRequest) (*socialV1.SocialAccount, error) {
	if req == nil || req.Data == nil {
		return nil, adminV1.ErrorBadRequest("invalid parameter")
	}

	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.Data.Id = trans.Ptr(req.GetId())
	req.Data.UpdatedBy = trans.Ptr(operator.GetUserId())
	if req.UpdateMask != nil {
		req.UpdateMask.Paths = append(req.UpdateMask.Paths, "updated_by")
	}

	return s.socialServiceClient.UpdateSocialAccount(ctx, req)
}

func (s *SocialService) DeleteSocialAccount(ctx context.Context, req *socialV1.DeleteSocialAccountRequest) (*emptypb.Empty, error) {
	return s.socialServiceClient.DeleteSocialAccount(ctx, req)
}

// ── 内容分发 ──

func (s *SocialService) ListSocialDistributions(ctx context.Context, req *paginationV1.PagingRequest) (*socialV1.ListSocialDistributionsResponse, error) {
	return s.socialServiceClient.ListSocialDistributions(ctx, req)
}

func (s *SocialService) GetSocialDistribution(ctx context.Context, req *socialV1.GetSocialDistributionRequest) (*socialV1.SocialDistribution, error) {
	return s.socialServiceClient.GetSocialDistribution(ctx, req)
}

func (s *SocialService) CreateSocialDistribution(ctx context.Context, req *socialV1.CreateSocialDistributionRequest) (*socialV1.SocialDistribution, error) {
	if req == nil || req.Data == nil {
		return nil, adminV1.ErrorBadRequest("invalid parameter")
	}

	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.Data.CreatedBy = trans.Ptr(operator.UserId)

	return s.socialServiceClient.CreateSocialDistribution(ctx, req)
}

func (s *SocialService) UpdateSocialDistribution(ctx context.Context, req *socialV1.UpdateSocialDistributionRequest) (*socialV1.SocialDistribution, error) {
	if req == nil || req.Data == nil {
		return nil, adminV1.ErrorBadRequest("invalid parameter")
	}

	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.Data.Id = trans.Ptr(req.GetId())
	req.Data.UpdatedBy = trans.Ptr(operator.GetUserId())
	if req.UpdateMask != nil {
		req.UpdateMask.Paths = append(req.UpdateMask.Paths, "updated_by")
	}

	return s.socialServiceClient.UpdateSocialDistribution(ctx, req)
}

func (s *SocialService) DeleteSocialDistribution(ctx context.Context, req *socialV1.DeleteSocialDistributionRequest) (*emptypb.Empty, error) {
	return s.socialServiceClient.DeleteSocialDistribution(ctx, req)
}

func (s *SocialService) PublishSocialDistribution(ctx context.Context, req *socialV1.PublishSocialDistributionRequest) (*socialV1.SocialDistribution, error) {
	return s.socialServiceClient.PublishSocialDistribution(ctx, req)
}

func (s *SocialService) RefreshSocialDistributionStats(ctx context.Context, req *socialV1.RefreshSocialDistributionStatsRequest) (*socialV1.SocialDistribution, error) {
	return s.socialServiceClient.RefreshSocialDistributionStats(ctx, req)
}

func (s *SocialService) ListSocialDistributionStatsHistory(ctx context.Context, req *socialV1.ListSocialDistributionStatsHistoryRequest) (*socialV1.ListSocialDistributionStatsHistoryResponse, error) {
	return s.socialServiceClient.ListSocialDistributionStatsHistory(ctx, req)
}
