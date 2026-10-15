package service

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/timestamppb"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	"go-wind-cms/pkg/preview"
)

// PreviewService 草稿预览令牌签发：admin 侧签发短时效 HMAC 令牌，
// C 端公开 BFF 校验后放行非已发布内容。令牌无服务端状态，过期即失效。
type PreviewService struct {
	adminV1.PreviewServiceHTTPServer

	log *log.Helper
}

func NewPreviewService(ctx *bootstrap.Context) *PreviewService {
	return &PreviewService{
		log: ctx.NewLoggerHelper("preview/service/admin-service"),
	}
}

func (s *PreviewService) CreatePostPreviewToken(ctx context.Context, req *adminV1.CreatePostPreviewTokenRequest) (*adminV1.PreviewToken, error) {
	if req == nil || req.GetId() == 0 {
		return nil, adminV1.ErrorBadRequest("invalid parameter")
	}
	return s.mint(ctx, preview.KindPost, req.GetId())
}

func (s *PreviewService) CreatePagePreviewToken(ctx context.Context, req *adminV1.CreatePagePreviewTokenRequest) (*adminV1.PreviewToken, error) {
	if req == nil || req.GetId() == 0 {
		return nil, adminV1.ErrorBadRequest("invalid parameter")
	}
	return s.mint(ctx, preview.KindPage, req.GetId())
}

func (s *PreviewService) mint(ctx context.Context, kind string, id uint32) (*adminV1.PreviewToken, error) {
	token, exp, err := preview.Mint(kind, id, preview.DefaultTTL, time.Now())
	if err != nil {
		s.log.Errorf("mint preview token failed (kind=%s id=%d): %v", kind, id, err)
		return nil, adminV1.ErrorInternalServerError("mint preview token failed")
	}
	return &adminV1.PreviewToken{
		Token:     &token,
		ExpiresAt: timestamppb.New(exp),
	}, nil
}
