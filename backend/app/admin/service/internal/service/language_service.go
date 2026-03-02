package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	dictV1 "go-wind-cms/api/gen/go/dict/service/v1"

	"go-wind-cms/pkg/middleware/auth"
)

type LanguageService struct {
	adminV1.LanguageServiceHTTPServer

	log *log.Helper

	languageServiceClient dictV1.LanguageServiceClient
}

func NewLanguageService(
	ctx *bootstrap.Context,
	languageServiceClient dictV1.LanguageServiceClient,
) *LanguageService {
	l := log.NewHelper(log.With(ctx.GetLogger(), "module", "language/service/admin-service"))
	return &LanguageService{
		log:                   l,
		languageServiceClient: languageServiceClient,
	}
}

func (s *LanguageService) List(ctx context.Context, req *paginationV1.PagingRequest) (*dictV1.ListLanguageResponse, error) {
	return s.languageServiceClient.List(ctx, req)
}

func (s *LanguageService) Get(ctx context.Context, req *dictV1.GetLanguageRequest) (*dictV1.Language, error) {
	return s.languageServiceClient.Get(ctx, req)
}

func (s *LanguageService) Create(ctx context.Context, req *dictV1.CreateLanguageRequest) (*emptypb.Empty, error) {
	if req.Data == nil {
		return nil, adminV1.ErrorBadRequest("invalid parameter")
	}

	// 获取操作人信息
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.Data.CreatedBy = trans.Ptr(operator.UserId)

	if _, err = s.languageServiceClient.Create(ctx, req); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *LanguageService) Update(ctx context.Context, req *dictV1.UpdateLanguageRequest) (*emptypb.Empty, error) {
	if req.Data == nil {
		return nil, adminV1.ErrorBadRequest("invalid parameter")
	}

	// 获取操作人信息
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.Data.Id = trans.Ptr(req.GetId())

	req.Data.UpdatedBy = trans.Ptr(operator.GetUserId())
	if req.UpdateMask != nil {
		req.UpdateMask.Paths = append(req.UpdateMask.Paths, "updated_by")
	}

	if _, err = s.languageServiceClient.Update(ctx, req); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *LanguageService) Delete(ctx context.Context, req *dictV1.DeleteLanguageRequest) (*emptypb.Empty, error) {
	return s.languageServiceClient.Delete(ctx, req)
}
