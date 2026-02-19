package ueditor

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	storageV1 "go-wind-cms/api/gen/go/storage/service/v1"

	"go-wind-cms/pkg/oss"
)

const (
	StateOK = "SUCCESS"
)

type Service struct {
	log *log.Helper

	mc *oss.MinIOClient
}

func NewUEditorService(ctx *bootstrap.Context, mc *oss.MinIOClient) *Service {
	l := log.NewHelper(log.With(ctx.GetLogger(), "module", "ueditor/service/file-service"))
	return &Service{
		log: l,
		mc:  mc,
	}
}

func (s *Service) UEditorAPI(ctx context.Context, req *storageV1.UEditorRequest) (*storageV1.UEditorResponse, error) {
	//s.log.Infof("UEditorAPI [%s]", req.GetAction())

	switch req.GetAction() {
	default:
		fallthrough
	case storageV1.UEditorAction_config.String():
		return &storageV1.UEditorResponse{}, nil

	case storageV1.UEditorAction_listFile.String():
		return s.mc.ListFileForUEditor(ctx, "files", "")

	case storageV1.UEditorAction_listImage.String():
		return s.mc.ListFileForUEditor(ctx, "images", "")
	}
}
