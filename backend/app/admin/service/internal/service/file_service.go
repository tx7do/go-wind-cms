package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	storageV1 "go-wind-cms/api/gen/go/storage/service/v1"

	"go-wind-cms/pkg/middleware/auth"
)

type FileService struct {
	adminV1.FileServiceHTTPServer

	log *log.Helper

	fileServiceClient storageV1.FileServiceClient
}

func NewFileService(ctx *bootstrap.Context, fileServiceClient storageV1.FileServiceClient) *FileService {
	l := log.NewHelper(log.With(ctx.GetLogger(), "module", "file/service/admin-service"))
	return &FileService{
		log:               l,
		fileServiceClient: fileServiceClient,
	}
}

func (s *FileService) List(ctx context.Context, req *paginationV1.PagingRequest) (*storageV1.ListFileResponse, error) {
	return s.fileServiceClient.List(ctx, req)
}

func (s *FileService) Get(ctx context.Context, req *storageV1.GetFileRequest) (*storageV1.File, error) {
	return s.fileServiceClient.Get(ctx, req)
}

func (s *FileService) Create(ctx context.Context, req *storageV1.CreateFileRequest) (*emptypb.Empty, error) {
	if req == nil || req.Data == nil {
		return nil, adminV1.ErrorBadRequest("invalid request")
	}

	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.Data.CreatedBy = trans.Ptr(operator.GetUserId())

	return s.fileServiceClient.Create(ctx, req)
}

func (s *FileService) Update(ctx context.Context, req *storageV1.UpdateFileRequest) (*emptypb.Empty, error) {
	if req == nil || req.Data == nil {
		return nil, adminV1.ErrorBadRequest("invalid request")
	}

	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.Data.UpdatedBy = trans.Ptr(operator.GetUserId())
	if req.UpdateMask != nil {
		req.UpdateMask.Paths = append(req.UpdateMask.Paths, "updated_by")
	}

	return s.fileServiceClient.Update(ctx, req)
}

func (s *FileService) Delete(ctx context.Context, req *storageV1.DeleteFileRequest) (*emptypb.Empty, error) {
	if req == nil {
		return nil, adminV1.ErrorBadRequest("invalid request")
	}

	return s.fileServiceClient.Delete(ctx, req)
}
