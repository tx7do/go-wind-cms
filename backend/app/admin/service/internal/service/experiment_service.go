package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
)

// ExperimentService 内容实验管理门面（透传 core）。
type ExperimentService struct {
	adminV1.ExperimentServiceHTTPServer

	experimentServiceClient contentV1.ExperimentServiceClient
	log                     *log.Helper
}

func NewExperimentService(ctx *bootstrap.Context, experimentServiceClient contentV1.ExperimentServiceClient) *ExperimentService {
	return &ExperimentService{
		log:                     ctx.NewLoggerHelper("experiment/service/admin-service"),
		experimentServiceClient: experimentServiceClient,
	}
}

func (s *ExperimentService) List(ctx context.Context, req *paginationV1.PagingRequest) (*contentV1.ListExperimentsResponse, error) {
	return s.experimentServiceClient.List(ctx, req)
}

func (s *ExperimentService) Get(ctx context.Context, req *contentV1.GetExperimentRequest) (*contentV1.Experiment, error) {
	return s.experimentServiceClient.Get(ctx, req)
}

func (s *ExperimentService) Create(ctx context.Context, req *contentV1.CreateExperimentRequest) (*contentV1.Experiment, error) {
	return s.experimentServiceClient.Create(ctx, req)
}

func (s *ExperimentService) Update(ctx context.Context, req *contentV1.UpdateExperimentRequest) (*contentV1.Experiment, error) {
	return s.experimentServiceClient.Update(ctx, req)
}

func (s *ExperimentService) Delete(ctx context.Context, req *contentV1.DeleteExperimentRequest) (*emptypb.Empty, error) {
	return s.experimentServiceClient.Delete(ctx, req)
}

func (s *ExperimentService) SetStatus(ctx context.Context, req *contentV1.SetExperimentStatusRequest) (*contentV1.Experiment, error) {
	return s.experimentServiceClient.SetStatus(ctx, req)
}

func (s *ExperimentService) GetVariantStats(ctx context.Context, req *contentV1.GetExperimentVariantStatsRequest) (*contentV1.ExperimentVariantStatsResponse, error) {
	return s.experimentServiceClient.GetVariantStats(ctx, req)
}
