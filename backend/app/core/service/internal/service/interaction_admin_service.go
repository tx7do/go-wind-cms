package service

import (
	"context"
	"strconv"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/go-crud/viewer"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"github.com/tx7do/go-utils/trans"

	"go-wind-cms/app/core/service/internal/data"

	auditV1 "go-wind-cms/api/gen/go/audit/service/v1"
	interactionV1 "go-wind-cms/api/gen/go/interaction/service/v1"
)

// InteractionAdminService 是点赞/收藏计数子系统的运营清数据面。
//
// 与 viewer-facing 的 InteractionService 相对，本服务面向后台运营：
//   - 接受 target_id / user_id 作为入参（运营跨用户操作别人数据）。
//   - 身份从运营上下文（OperatorMetadata 注入的 viewer context）取，
//     调用方须为平台管理员（IsPlatformContext || IsSystemContext），否则 403。
//   - 每次操作写 OperationAuditLog（Action=DELETE，记录操作人/目标/成败）。
//
// 注：HasPermission 当前为空实现，无法做细粒度权限点，只能做平台管理员硬 gate。
type InteractionAdminService struct {
	interactionV1.UnimplementedInteractionAdminServiceServer

	interactionRepo      *data.InteractionRepo
	operationAuditLogRepo *data.OperationAuditLogRepo
	log                   *log.Helper
}

func NewInteractionAdminService(
	ctx *bootstrap.Context,
	interactionRepo *data.InteractionRepo,
	operationAuditLogRepo *data.OperationAuditLogRepo,
) *InteractionAdminService {
	return &InteractionAdminService{
		interactionRepo:       interactionRepo,
		operationAuditLogRepo: operationAuditLogRepo,
		log:                   ctx.NewLoggerHelper("interaction-admin/service/core-service"),
	}
}

// requireAdminOperator 校验调用者为平台/系统管理员，返回操作人的 (tenantID, userID)。
// 非管理员上下文一律 403。这是仓库首个“硬性平台管理员 gate”
// （现有 IsPlatformContext||IsSystemContext 用法均为 bypass 放行）。
func (s *InteractionAdminService) requireAdminOperator(ctx context.Context) (operatorTenantID, operatorUserID uint32, err error) {
	vc, exist := viewer.FromContext(ctx)
	if !exist || vc == nil {
		return 0, 0, interactionV1.ErrorUnauthorized("operator identity required")
	}
	if !(vc.IsPlatformContext() || vc.IsSystemContext()) {
		return 0, 0, interactionV1.ErrorForbidden("platform admin only")
	}
	return uint32(vc.TenantID()), uint32(vc.UserID()), nil
}

// writeAudit 写一条 OperationAuditLog（Action=DELETE）。resourceID 标识被操作目标。
func (s *InteractionAdminService) writeAudit(ctx context.Context, operatorTenantID, operatorUserID uint32, resourceType, resourceID string, success bool) {
	req := &auditV1.CreateOperationAuditLogRequest{
		Data: &auditV1.OperationAuditLog{
			TenantId:     trans.Ptr(operatorTenantID),
			UserId:       trans.Ptr(operatorUserID),
			Action:       auditV1.OperationAuditLog_DELETE.Enum(),
			ResourceType: trans.Ptr(resourceType),
			ResourceId:   trans.Ptr(resourceID),
			Success:      trans.Ptr(success),
		},
	}
	if err := s.operationAuditLogRepo.Create(ctx, req); err != nil {
		s.log.Errorf("write interaction purge audit failed: %s", err.Error())
	}
}

func (s *InteractionAdminService) PurgeTargetInteractions(ctx context.Context, req *interactionV1.PurgeTargetInteractionsRequest) (*interactionV1.PurgeTargetInteractionsResponse, error) {
	opTid, opUid, err := s.requireAdminOperator(ctx)
	if err != nil {
		return nil, err
	}

	affected, err := s.interactionRepo.PurgeTargetInteractions(ctx, req.GetTargetType(), req.GetTargetId())
	s.writeAudit(ctx, opTid, opUid, "interaction_counter",
		strconv.FormatUint(uint64(req.GetTargetId()), 10), err == nil)
	if err != nil {
		return nil, err
	}
	return &interactionV1.PurgeTargetInteractionsResponse{AffectedRows: affected}, nil
}

func (s *InteractionAdminService) PurgeUserInteractions(ctx context.Context, req *interactionV1.PurgeUserInteractionsRequest) (*interactionV1.PurgeUserInteractionsResponse, error) {
	opTid, opUid, err := s.requireAdminOperator(ctx)
	if err != nil {
		return nil, err
	}

	affected, err := s.interactionRepo.PurgeUserInteractions(ctx, req.GetUserId())
	s.writeAudit(ctx, opTid, opUid, "interaction_user_ledger",
		strconv.FormatUint(uint64(req.GetUserId()), 10), err == nil)
	if err != nil {
		return nil, err
	}
	return &interactionV1.PurgeUserInteractionsResponse{AffectedRows: affected}, nil
}

func (s *InteractionAdminService) ResetCounter(ctx context.Context, req *interactionV1.ResetCounterRequest) (*interactionV1.ResetCounterResponse, error) {
	opTid, opUid, err := s.requireAdminOperator(ctx)
	if err != nil {
		return nil, err
	}

	recount, err := s.interactionRepo.ResetCounter(ctx, req.GetTargetType(), req.GetTargetId(), req.GetMetric())
	s.writeAudit(ctx, opTid, opUid, "interaction_counter_reset",
		strconv.FormatUint(uint64(req.GetTargetId()), 10), err == nil)
	if err != nil {
		return nil, err
	}
	return &interactionV1.ResetCounterResponse{Recount: recount}, nil
}
