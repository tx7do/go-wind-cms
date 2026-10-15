package service

import (
	"context"
	"encoding/json"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	"go-wind-cms/app/core/service/internal/data"
	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
)

// ExperimentService 内容实验（A/B 分流 + 受众定向）。
// 状态机：draft→running；running↔paused；running/paused→ended；ended 终态。
// 开跑（→running）要求至少一个变体。分流的访客侧逻辑在 app BFF。
type ExperimentService struct {
	contentV1.UnimplementedExperimentServiceServer

	experimentRepo *data.ExperimentRepo
	log            *log.Helper
}

func NewExperimentService(ctx *bootstrap.Context, experimentRepo *data.ExperimentRepo) *ExperimentService {
	return &ExperimentService{
		log:            ctx.NewLoggerHelper("experiment/service/core-service"),
		experimentRepo: experimentRepo,
	}
}

func (s *ExperimentService) List(ctx context.Context, req *paginationV1.PagingRequest) (*contentV1.ListExperimentsResponse, error) {
	items, total, err := s.experimentRepo.List(ctx, req)
	if err != nil {
		return nil, err
	}
	return &contentV1.ListExperimentsResponse{Items: items, Total: uint64(total)}, nil
}

func (s *ExperimentService) Get(ctx context.Context, req *contentV1.GetExperimentRequest) (*contentV1.Experiment, error) {
	if req == nil || req.GetId() == 0 {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}
	return s.experimentRepo.Get(ctx, req.GetId())
}

func (s *ExperimentService) Create(ctx context.Context, req *contentV1.CreateExperimentRequest) (*contentV1.Experiment, error) {
	if req == nil || req.Data == nil || req.Data.GetName() == "" {
		return nil, contentV1.ErrorBadRequest("name is required")
	}
	if err := validateVariants(req.Data.GetVariants()); err != nil {
		return nil, err
	}
	return s.experimentRepo.Create(ctx, req.Data)
}

func (s *ExperimentService) Update(ctx context.Context, req *contentV1.UpdateExperimentRequest) (*contentV1.Experiment, error) {
	if req == nil || req.Data == nil || req.GetId() == 0 {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}
	if err := validateVariants(req.Data.GetVariants()); err != nil {
		return nil, err
	}
	// 变体为指针切片：字段出现（含空数组）即视为整组替换语义
	return s.experimentRepo.Update(ctx, req.GetId(), req.Data, req.Data.Variants != nil)
}

func (s *ExperimentService) SetStatus(ctx context.Context, req *contentV1.SetExperimentStatusRequest) (*contentV1.Experiment, error) {
	if req == nil || req.GetId() == 0 || req.GetStatus() == contentV1.Experiment_EXPERIMENT_STATUS_UNSPECIFIED {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}
	current, err := s.experimentRepo.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	if !canTransition(current.GetStatus(), req.GetStatus()) {
		return nil, contentV1.ErrorBadRequest("invalid status transition: %s -> %s", current.GetStatus(), req.GetStatus())
	}
	if req.GetStatus() == contentV1.Experiment_EXPERIMENT_STATUS_RUNNING && len(current.GetVariants()) == 0 {
		return nil, contentV1.ErrorBadRequest("experiment needs at least one variant before running")
	}
	return s.experimentRepo.SetStatus(ctx, req.GetId(), req.GetStatus())
}

func (s *ExperimentService) Delete(ctx context.Context, req *contentV1.DeleteExperimentRequest) (*emptypb.Empty, error) {
	if req == nil || req.GetId() == 0 {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}
	if err := s.experimentRepo.Delete(ctx, req.GetId()); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// ListRunningByPage 页面上运行中的实验（app BFF 分流入口）。
func (s *ExperimentService) ListRunningByPage(ctx context.Context, req *contentV1.ListRunningExperimentsByPageRequest) (*contentV1.ListExperimentsResponse, error) {
	if req == nil || req.GetPageId() == 0 {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}
	items, err := s.experimentRepo.ListRunningByPage(ctx, req.GetPageId())
	if err != nil {
		return nil, err
	}
	return &contentV1.ListExperimentsResponse{Items: items}, nil
}

// RecordParticipant 记录访客分配（幂等）。
func (s *ExperimentService) RecordParticipant(ctx context.Context, req *contentV1.RecordExperimentParticipantRequest) (*emptypb.Empty, error) {
	if req == nil || req.GetExperimentId() == 0 || req.GetVisitorHash() == "" || req.GetVariantName() == "" {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}
	if err := s.experimentRepo.RecordParticipant(ctx, req.GetExperimentId(), req.GetVisitorHash(), req.GetVariantName()); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// GetVariantStats 变体参与统计。
func (s *ExperimentService) GetVariantStats(ctx context.Context, req *contentV1.GetExperimentVariantStatsRequest) (*contentV1.ExperimentVariantStatsResponse, error) {
	if req == nil || req.GetExperimentId() == 0 {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}
	return s.experimentRepo.GetVariantStats(ctx, req.GetExperimentId())
}

// canTransition 实验状态机。
func canTransition(from, to contentV1.Experiment_Status) bool {
	switch from {
	case contentV1.Experiment_EXPERIMENT_STATUS_DRAFT:
		return to == contentV1.Experiment_EXPERIMENT_STATUS_RUNNING || to == contentV1.Experiment_EXPERIMENT_STATUS_ENDED
	case contentV1.Experiment_EXPERIMENT_STATUS_RUNNING:
		return to == contentV1.Experiment_EXPERIMENT_STATUS_PAUSED || to == contentV1.Experiment_EXPERIMENT_STATUS_ENDED
	case contentV1.Experiment_EXPERIMENT_STATUS_PAUSED:
		return to == contentV1.Experiment_EXPERIMENT_STATUS_RUNNING || to == contentV1.Experiment_EXPERIMENT_STATUS_ENDED
	}
	return false
}

// validateVariants 变体基础校验：名称必填、权重范围、受众条件/区块快照
// 必须是合法 JSON（条件求值与 protojson 解析在消费侧，坏 JSON 提前拦下）。
func validateVariants(variants []*contentV1.ExperimentVariant) error {
	seen := make(map[string]bool, len(variants))
	totalWeight := 0
	for _, v := range variants {
		if v == nil {
			continue
		}
		if v.GetName() == "" {
			return contentV1.ErrorBadRequest("variant name is required")
		}
		if seen[v.GetName()] {
			return contentV1.ErrorBadRequest("duplicated variant name: %s", v.GetName())
		}
		seen[v.GetName()] = true
		if v.Weight != nil && (*v.Weight < 0 || *v.Weight > 100) {
			return contentV1.ErrorBadRequest("variant weight must be within [0, 100]")
		}
		totalWeight += int(v.GetWeight())
		if cond := v.GetConditions(); cond != "" {
			if !json.Valid([]byte(cond)) {
				return contentV1.ErrorBadRequest("variant conditions is not valid JSON")
			}
		}
		if sec := v.GetSections(); sec != "" {
			if !json.Valid([]byte(sec)) {
				return contentV1.ErrorBadRequest("variant sections is not valid JSON")
			}
		}
	}
	// random 型权重总和应不超过 100（不足部分视为不分配/回落默认渲染）
	if totalWeight > 100 {
		return contentV1.ErrorBadRequest("variant weights sum exceeds 100")
	}
	return nil
}
