package data

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	entCrud "github.com/tx7do/go-crud/entgo"
	"github.com/tx7do/go-utils/copierutil"
	"github.com/tx7do/go-utils/mapper"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/aiagentrun"
	"go-wind-cms/app/core/service/internal/data/ent/predicate"

	aiV1 "go-wind-cms/api/gen/go/ai/service/v1"
)

type AiAgentRunRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper

	mapper          *mapper.CopierMapper[aiV1.AiAgentRun, ent.AiAgentRun]
	statusConverter *mapper.EnumTypeConverter[aiV1.AiAgentRun_Status, aiagentrun.Status]

	repository *entCrud.Repository[
		ent.AiAgentRunQuery, ent.AiAgentRunSelect,
		ent.AiAgentRunCreate, ent.AiAgentRunCreateBulk,
		ent.AiAgentRunUpdate, ent.AiAgentRunUpdateOne,
		ent.AiAgentRunDelete,
		predicate.AiAgentRun,
		aiV1.AiAgentRun, ent.AiAgentRun,
	]
}

func NewAiAgentRunRepo(ctx *bootstrap.Context, entClient *entCrud.EntClient[*ent.Client]) *AiAgentRunRepo {
	repo := &AiAgentRunRepo{
		log:       ctx.NewLoggerHelper("ai-agent-run/repo/core-service"),
		entClient: entClient,
		mapper:    mapper.NewCopierMapper[aiV1.AiAgentRun, ent.AiAgentRun](),
		statusConverter: mapper.NewEnumTypeConverter[aiV1.AiAgentRun_Status, aiagentrun.Status](
			aiV1.AiAgentRun_Status_name, aiV1.AiAgentRun_Status_value,
		),
	}
	repo.init()
	return repo
}

func (r *AiAgentRunRepo) init() {
	r.repository = entCrud.NewRepository[
		ent.AiAgentRunQuery, ent.AiAgentRunSelect,
		ent.AiAgentRunCreate, ent.AiAgentRunCreateBulk,
		ent.AiAgentRunUpdate, ent.AiAgentRunUpdateOne,
		ent.AiAgentRunDelete,
		predicate.AiAgentRun,
		aiV1.AiAgentRun, ent.AiAgentRun,
	](r.mapper)

	r.mapper.AppendConverters(copierutil.NewTimeStringConverterPair())
	r.mapper.AppendConverters(copierutil.NewTimeTimestamppbConverterPair())
	r.mapper.AppendConverters(r.statusConverter.NewConverterPair())
}

// List 分页查询运行台账（appId eq 过滤走通用机制；status 由读路径回填）。
func (r *AiAgentRunRepo) List(ctx context.Context, req *paginationV1.PagingRequest) (*aiV1.ListAiAgentRunsResponse, error) {
	if req == nil {
		return nil, errors.BadRequest("AI_AGENT", "invalid parameter")
	}

	builder := r.entClient.Client().AiAgentRun.Query()
	ret, err := r.repository.ListWithPaging(ctx, builder, builder.Clone(), req)
	if err != nil {
		r.log.Errorf("list ai agent runs failed: %s", err.Error())
		return nil, errors.InternalServer("AI_AGENT", "list ai agent runs failed")
	}
	if ret == nil {
		return &aiV1.ListAiAgentRunsResponse{Total: 0}, nil
	}

	// status 是值型枚举列（Default 无 Nillable）→ DTO 可选指针，copier 失配（同渠道 Type 的教训）。
	entities, err := r.entClient.Client().AiAgentRun.Query().
		Select(aiagentrun.FieldID, aiagentrun.FieldStatus).
		All(ctx)
	if err != nil {
		r.log.Errorf("backfill status failed: %s", err.Error())
	} else {
		statusByID := make(map[uint32]aiagentrun.Status, len(entities))
		for _, e := range entities {
			statusByID[e.ID] = e.Status
		}
		for _, it := range ret.Items {
			if s, ok := statusByID[it.GetId()]; ok {
				sv := s
				if p := r.statusConverter.ToDTO(&sv); p != nil {
					it.Status = p
				}
			}
		}
	}

	return &aiV1.ListAiAgentRunsResponse{Total: ret.Total, Items: ret.Items}, nil
}

// Get 按 ID 查询运行记录。
func (r *AiAgentRunRepo) Get(ctx context.Context, id uint32) (*aiV1.AiAgentRun, error) {
	entity, err := r.entClient.Client().AiAgentRun.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.NotFound("AI_AGENT", "ai agent run not found")
		}
		r.log.Errorf("get ai agent run [%d] failed: %s", id, err.Error())
		return nil, errors.InternalServer("AI_AGENT", "get ai agent run failed")
	}
	dto := r.mapper.ToDTO(entity)
	sv := entity.Status
	if p := r.statusConverter.ToDTO(&sv); p != nil {
		dto.Status = p
	}
	return dto, nil
}

// CreateRunning 建一条 RUNNING 运行行（Run 流程第一步，执行引擎批接入）。
func (r *AiAgentRunRepo) CreateRunning(ctx context.Context, appID uint32, inputJSON string, operatorID uint32, tenantID uint32) (uint32, error) {
	created, err := r.entClient.Client().AiAgentRun.Create().
		SetAppID(appID).
		SetStatus(aiagentrun.StatusRunning).
		SetNillableInput(trans.Ptr(inputJSON)).
		SetTrace("[]").
		SetNillableTenantID(trans.Ptr(tenantID)).
		SetCreatedBy(operatorID).
		SetCreatedAt(time.Now()).
		Save(ctx)
	if err != nil {
		r.log.Errorf("insert ai agent run failed: %s", err.Error())
		return 0, errors.InternalServer("AI_AGENT", "insert ai agent run failed")
	}
	return created.ID, nil
}

// AgentRunTokens 运行用量三件套（repo 落参用，避免一串同形参数）。
type AgentRunTokens struct {
	Prompt     uint32
	Completion uint32
	Total      uint32
}

// SettleResult 回写运行结论（终态 + 产出/错误/轨迹/用量/耗时）。
func (r *AiAgentRunRepo) SettleResult(
	ctx context.Context,
	id uint32,
	status aiagentrun.Status,
	output string,
	errMsg string,
	traceJSON string,
	usage AgentRunTokens,
	durationMs uint32,
) error {
	builder := r.entClient.Client().AiAgentRun.UpdateOneID(id).
		SetStatus(status).
		SetTrace(traceJSON).
		SetPromptTokens(usage.Prompt).
		SetCompletionTokens(usage.Completion).
		SetTotalTokens(usage.Total).
		SetDurationMs(durationMs)
	if output != "" {
		builder.SetOutput(output)
	}
	if errMsg != "" {
		builder.SetError(errMsg)
	}
	if err := builder.Exec(ctx); err != nil {
		r.log.Errorf("settle ai agent run [%d] failed: %s", id, err.Error())
		return errors.InternalServer("AI_AGENT", "settle ai agent run failed")
	}
	return nil
}
