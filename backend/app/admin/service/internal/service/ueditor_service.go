package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	storageV1 "go-wind-cms/api/gen/go/storage/service/v1"

	"go-wind-cms/pkg/oss"
	"go-wind-cms/pkg/ueditor"
)

type UEditorService struct {
	adminV1.UEditorServiceHTTPServer

	log *log.Helper

	mc *ueditor.Service
}

func NewUEditorService(ctx *bootstrap.Context, mc *oss.MinIOClient) *UEditorService {
	l := log.NewHelper(log.With(ctx.GetLogger(), "module", "ueditor/service/admin-service"))
	return &UEditorService{
		log: l,
		mc:  ueditor.NewUEditorService(ctx, mc),
	}
}

func (s *UEditorService) UEditorAPI(ctx context.Context, req *storageV1.UEditorRequest) (*storageV1.UEditorResponse, error) {
	return s.mc.UEditorAPI(ctx, req)
}
