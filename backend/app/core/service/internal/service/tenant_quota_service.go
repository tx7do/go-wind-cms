package service

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	identityV1 "go-wind-cms/api/gen/go/identity/service/v1"

	appViewer "go-wind-cms/pkg/entgo/viewer"

	"go-wind-cms/app/core/service/internal/data"
	"go-wind-cms/app/core/service/internal/data/ent/tenant"
	"go-wind-cms/app/core/service/internal/data/ent/tenantquota"

	"go-wind-cms/pkg/task"
)

// QuotaWatermarkThreshold 水位阈值（用量/上限达到该比例即命中）。
const QuotaWatermarkThreshold = 0.8

// TenantQuotaService 租户配额：按租户直配的配额 CRUD + 用量聚合 + 双闸门 + 水位扫描。
//
// 与上游（PlanQuota 挂 plan_id、租户经套餐边取配额）不同：CMS 无套餐表，配额行
// 直接挂 tenant_id。闸门 fail-open 语义与上游一致：配额缺失 = 未限售，拒绝只在
// "配额明确存在且已用满"时发生。
type TenantQuotaService struct {
	identityV1.UnimplementedTenantQuotaServiceServer

	log  *log.Helper
	repo *data.TenantQuotaRepo

	// internalMessages 站内信投递缝（水位告警出口）：双向引用，构造后注册
	//（与通知域 Notifier 缝同型；未注入时水位任务只扫描不投递并留告警日志）。
	internalMessages internalMessageSender
}

// internalMessageSender 水位告警投递所需的最小站内信面（便于测试桩替）。
type internalMessageSender interface {
	SendSystemNotification(ctx context.Context, title, content string, targetUserIDs []uint32) (uint32, error)
}

// RegisterInternalMessages 注册站内信投递缝（构造完成后调用，解双向引用）。
func (s *TenantQuotaService) RegisterInternalMessages(sender internalMessageSender) {
	s.internalMessages = sender
}

func NewTenantQuotaService(ctx *bootstrap.Context, repo *data.TenantQuotaRepo) *TenantQuotaService {
	return &TenantQuotaService{
		log:  ctx.NewLoggerHelper("tenant-quota/service/core-service"),
		repo: repo,
	}
}

// ─── CRUD ───

func (s *TenantQuotaService) List(ctx context.Context, req *paginationV1.PagingRequest) (*identityV1.ListTenantQuotaResponse, error) {
	return s.repo.List(ctx, req)
}

func (s *TenantQuotaService) Get(ctx context.Context, req *identityV1.GetTenantQuotaRequest) (*identityV1.TenantQuota, error) {
	if req == nil || req.GetId() == 0 {
		return nil, errors.BadRequest("TENANT_QUOTA", "id is required")
	}
	return s.repo.Get(ctx, req.GetId())
}

func (s *TenantQuotaService) Create(ctx context.Context, req *identityV1.CreateTenantQuotaRequest) (*emptypb.Empty, error) {
	if req == nil || req.Data == nil {
		return nil, errors.BadRequest("TENANT_QUOTA", "invalid parameter")
	}
	operatorID, _ := viewerUserIDFromContext(ctx)
	if _, err := s.repo.Create(ctx, req.GetData(), operatorID); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *TenantQuotaService) Update(ctx context.Context, req *identityV1.UpdateTenantQuotaRequest) (*emptypb.Empty, error) {
	if req == nil || req.Data == nil {
		return nil, errors.BadRequest("TENANT_QUOTA", "invalid parameter")
	}
	if req.GetId() == 0 {
		return nil, errors.BadRequest("TENANT_QUOTA", "id is required")
	}
	operatorID, _ := viewerUserIDFromContext(ctx)
	if err := s.repo.Update(ctx, req.GetId(), req.GetData(), operatorID); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *TenantQuotaService) Delete(ctx context.Context, req *identityV1.DeleteTenantQuotaRequest) (*emptypb.Empty, error) {
	if req == nil || req.GetId() == 0 {
		return nil, errors.BadRequest("TENANT_QUOTA", "id is required")
	}
	if err := s.repo.Delete(ctx, req.GetId()); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// ─── 用量聚合 ───

// buildUsage 聚合租户用量与配额（GetUsage/GetMyTenantUsage 共用）。
// 单维度计量失败记日志后按 0 计——展示不因单表抖动整页失败（闸门另有 fail-open）。
func (s *TenantQuotaService) buildUsage(ctx context.Context, tenantID uint32) (*identityV1.TenantUsage, error) {
	sysCtx := appViewer.NewSystemViewerContext(ctx)
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)

	usage := &identityV1.TenantUsage{TenantId: &tenantID}

	// 配额行 + 租户基础信息（套餐名/管理员）
	quotas, err := s.repo.ListByTenant(sysCtx, tenantID)
	if err != nil {
		s.log.Errorf("get usage: list quotas for tenant %d failed: %s", tenantID, err.Error())
		return nil, errors.InternalServer("TENANT_QUOTA", "get usage failed")
	}
	for _, q := range quotas {
		usage.Quotas = append(usage.Quotas, &identityV1.QuotaUsage{
			QuotaType: quotaTypeEntToProto(q.QuotaType),
			QuotaValue: &q.QuotaValue,
		})
	}

	if t, terr := s.repo.GetTenantStatus(sysCtx, tenantID); terr != nil {
		s.log.Errorf("get usage: tenant %d not loadable: %s", tenantID, terr.Error())
	} else if t != nil {
		usage.SubscriptionPlan = t.SubscriptionPlan
	}

	if n, uerr := s.repo.CountUsersInTenant(sysCtx, tenantID); uerr != nil {
		s.log.Errorf("get usage: count users for tenant %d failed: %s", tenantID, uerr.Error())
	} else {
		usage.UserCount = &n
	}
	if n, serr := s.repo.SumStorageInTenant(sysCtx, tenantID); serr != nil {
		s.log.Errorf("get usage: sum storage for tenant %d failed: %s", tenantID, serr.Error())
	} else {
		usage.StorageUsedBytes = &n
	}
	if n, aerr := s.repo.CountApiCallsInTenant(sysCtx, tenantID); aerr != nil {
		s.log.Errorf("get usage: count api calls for tenant %d failed: %s", tenantID, aerr.Error())
	} else {
		usage.ApiCallCount = &n
	}
	if n, terr := s.repo.SumAiTokensInTenantSince(sysCtx, tenantID, monthStart); terr != nil {
		s.log.Errorf("get usage: sum ai tokens for tenant %d failed: %s", tenantID, terr.Error())
	} else {
		usage.AiTokensMonth = &n
	}

	return usage, nil
}

func (s *TenantQuotaService) GetUsage(ctx context.Context, req *identityV1.GetTenantUsageRequest) (*identityV1.TenantUsage, error) {
	if req == nil || req.GetTenantId() == 0 {
		return nil, errors.BadRequest("TENANT_QUOTA", "tenant id is required")
	}
	return s.buildUsage(ctx, req.GetTenantId())
}

func (s *TenantQuotaService) GetMyTenantUsage(ctx context.Context, req *identityV1.GetMyTenantUsageRequest) (*identityV1.TenantUsage, error) {
	if req == nil || req.GetTenantId() == 0 {
		return nil, errors.BadRequest("TENANT_QUOTA", "tenant id is required (injected by BFF)")
	}
	return s.buildUsage(ctx, req.GetTenantId())
}

// ─── 闸门（每请求网关闸 + 业务入口闸） ───

// CheckTenantAccess 租户访问闸门：租户必须存在且 status==ON；随后 API_CALL 配额。
// 全链 fail-open（租户查不到/计量失败放行仅告警——上游同语义：配额缺失=未限售，
// 基础设施抖动不打死租户业务）。拒绝只在"配额明确存在且已用满"时发生。
func (s *TenantQuotaService) CheckTenantAccess(ctx context.Context, req *identityV1.CheckTenantAccessRequest) (*identityV1.CheckTenantAccessResponse, error) {
	if req == nil || req.GetTenantId() == 0 {
		// 平台身份（tenant 0）与缺租户的调用不走租户闸门
		return &identityV1.CheckTenantAccessResponse{}, nil
	}

	sysCtx := appViewer.NewSystemViewerContext(ctx)

	t, err := s.repo.GetTenantStatus(sysCtx, req.GetTenantId())
	if err != nil {
		s.log.Errorf("tenant gate: tenant %d lookup failed: %s", req.GetTenantId(), err.Error())
		return &identityV1.CheckTenantAccessResponse{}, nil // fail-open
	}
	if t == nil {
		// 租户行不存在：放行交给下游（登录链/免租户接口可能合法无租户）
		return &identityV1.CheckTenantAccessResponse{}, nil
	}
	if t.Status != nil && *t.Status != tenant.StatusOn {
		return nil, errors.Forbidden("TENANT_QUOTA", "tenant is not active")
	}

	// API_CALL 配额（计量与 GetUsage 同源）
	limit, hasLimit, qerr := s.quotaLimitOf(sysCtx, req.GetTenantId(), identityV1.TenantQuota_API_CALL)
	if qerr != nil {
		s.log.Errorf("tenant gate: api call metering for tenant %d failed: %s", req.GetTenantId(), qerr.Error())
		return &identityV1.CheckTenantAccessResponse{}, nil // 计量失败 fail-open
	}
	if hasLimit {
		cnt, cerr := s.repo.CountApiCallsInTenant(sysCtx, req.GetTenantId())
		if cerr != nil {
			s.log.Errorf("tenant gate: api call count for tenant %d failed: %s", req.GetTenantId(), cerr.Error())
			return &identityV1.CheckTenantAccessResponse{}, nil // fail-open
		}
		if err := checkQuotaNotExceeded("API_CALL", limit, uint64(cnt), 0); err != nil {
			s.log.Warnf("tenant gate: tenant %d request %s %s blocked by API_CALL quota", req.GetTenantId(), req.GetMethod(), req.GetPath())
			return nil, err
		}
	}

	return &identityV1.CheckTenantAccessResponse{}, nil
}

// CheckQuota 业务入口闸门：USER_LIMIT（创建用户前）/ STORAGE（上传前，带本次字节数）。
func (s *TenantQuotaService) CheckQuota(ctx context.Context, req *identityV1.CheckQuotaRequest) (*identityV1.CheckQuotaResponse, error) {
	if req == nil || req.GetTenantId() == 0 {
		return &identityV1.CheckQuotaResponse{}, nil
	}

	sysCtx := appViewer.NewSystemViewerContext(ctx)

	switch req.GetQuotaType() {
	case identityV1.TenantQuota_USER_LIMIT:
		limit, hasLimit, err := s.quotaLimitOf(sysCtx, req.GetTenantId(), identityV1.TenantQuota_USER_LIMIT)
		if err != nil {
			s.log.Errorf("quota gate: user count for tenant %d failed: %s", req.GetTenantId(), err.Error())
			return &identityV1.CheckQuotaResponse{}, nil // fail-open
		}
		if !hasLimit {
			return &identityV1.CheckQuotaResponse{}, nil
		}
		cnt, cerr := s.repo.CountUsersInTenant(sysCtx, req.GetTenantId())
		if cerr != nil {
			s.log.Errorf("quota gate: user count for tenant %d failed: %s", req.GetTenantId(), cerr.Error())
			return &identityV1.CheckQuotaResponse{}, nil // fail-open
		}
		if err := checkQuotaNotExceeded("USER_LIMIT", limit, cnt, 0); err != nil {
			return nil, err
		}

	case identityV1.TenantQuota_STORAGE:
		limit, hasLimit, err := s.quotaLimitOf(sysCtx, req.GetTenantId(), identityV1.TenantQuota_STORAGE)
		if err != nil {
			s.log.Errorf("quota gate: storage metering for tenant %d failed: %s", req.GetTenantId(), err.Error())
			return &identityV1.CheckQuotaResponse{}, nil // fail-open
		}
		if !hasLimit {
			return &identityV1.CheckQuotaResponse{}, nil
		}
		used, cerr := s.repo.SumStorageInTenant(sysCtx, req.GetTenantId())
		if cerr != nil {
			s.log.Errorf("quota gate: storage sum for tenant %d failed: %s", req.GetTenantId(), cerr.Error())
			return &identityV1.CheckQuotaResponse{}, nil // fail-open
		}
		if err := checkQuotaNotExceeded("STORAGE", limit, used, req.GetAdditional()); err != nil {
			return nil, err
		}

	default:
		// API_CALL 走 CheckTenantAccess；AI_TOKENS 当前只在水位告警，不做入口硬闸
		return &identityV1.CheckQuotaResponse{}, nil
	}

	return &identityV1.CheckQuotaResponse{}, nil
}

// quotaLimitOf 取租户指定类型的配额上限；无条目返回 hasLimit=false（未限售）。
func (s *TenantQuotaService) quotaLimitOf(ctx context.Context, tenantID uint32, qt identityV1.TenantQuota_QuotaType) (limit uint64, hasLimit bool, err error) {
	quotas, err := s.repo.ListByTenant(ctx, tenantID)
	if err != nil {
		return 0, false, err
	}
	for _, q := range quotas {
		if q.QuotaType == tenantquota.QuotaType(quotaTypeName(qt)) {
			return q.QuotaValue, true, nil
		}
	}
	return 0, false, nil
}

// ─── 纯判定（可测；闸门语义在此单点） ───

// quotaTypeName proto 枚举的技术名（与 ent NamedValues 存储值逐字一致）。
func quotaTypeName(qt identityV1.TenantQuota_QuotaType) string {
	return identityV1.TenantQuota_QuotaType_name[int32(qt)]
}

// quotaTypeEntToProto ent 枚举 → proto 枚举（存储名逐字一致，按名反查）。
func quotaTypeEntToProto(qt tenantquota.QuotaType) *identityV1.TenantQuota_QuotaType {
	v, ok := identityV1.TenantQuota_QuotaType_value[string(qt)]
	if !ok {
		return identityV1.TenantQuota_QUOTA_TYPE_UNSPECIFIED.Enum()
	}
	return identityV1.TenantQuota_QuotaType(v).Enum()
}

// checkQuotaNotExceeded 纯判定：配额存在且超限 → 403（消息带类型与数值供前端展示）。
//   - STORAGE（容量型）：used+本次上传 > limit 即超；
//   - USER_LIMIT / API_CALL（计数型）：used >= limit 即超。
func checkQuotaNotExceeded(quotaType string, limit, used, additional uint64) error {
	exceeded := false
	if quotaType == "STORAGE" {
		exceeded = used+additional > limit
	} else {
		exceeded = used >= limit
	}
	if !exceeded {
		return nil
	}
	return errors.Forbidden("TENANT_QUOTA", fmt.Sprintf("plan quota exceeded: %s limit %d reached (current %d)", quotaType, limit, used))
}

// ─── 水位扫描（只读 + 提前告警，与硬执行互补） ───

// QuotaWatermarkHit 单租户单配额类型的水位命中（只读扫描结果，无任何拒绝语义）。
type QuotaWatermarkHit struct {
	TenantID   uint32
	TenantName string
	// AdminUserID 租户管理员（sys_tenants.admin_user_id；0 = 未设置，投递侧跳过）。
	AdminUserID uint32
	// QuotaType 配额类型技术名（proto 枚举名，与用量页/闸门报错口径一致，不做翻译）。
	QuotaType string
	Used      uint64
	Limit     uint64
	// RatioPct 用量百分数（四舍五入；limit==0 且 used>0 视为 100）。
	RatioPct uint64
}

// quotaWatermarkExceeded 纯判定（可测）：limit==0 且 used>0 视为命中（配 0 上限
// 且已有用量 = 事实超限）；否则 used >= 阈值×limit 即命中。
func quotaWatermarkExceeded(used, limit uint64) bool {
	if limit == 0 {
		return used > 0
	}
	return float64(used) >= QuotaWatermarkThreshold*float64(limit)
}

// ratioPct 纯函数（可测）：用量百分数。
func ratioPct(used, limit uint64) uint64 {
	if limit == 0 {
		if used > 0 {
			return 100
		}
		return 0
	}
	return uint64(math.Round(float64(used) / float64(limit) * 100))
}

// ScanQuotaWatermarks 扫描全部 ON 租户的配额水位，返回命中清单。
// 只读：不拒绝任何请求、不写任何业务行；投递由站内信完成。
func (s *TenantQuotaService) ScanQuotaWatermarks(ctx context.Context, monthStart time.Time) ([]QuotaWatermarkHit, error) {
	sysCtx := appViewer.NewSystemViewerContext(ctx)

	quotas, err := s.repo.ListAll(sysCtx)
	if err != nil {
		s.log.Errorf("quota watermark scan: list quotas failed: %s", err.Error())
		return nil, err
	}
	if len(quotas) == 0 {
		return nil, nil
	}

	// 按租户归集配额；只扫描 status==ON 的租户（停用/冻结租户已被闸门拦截，告警是噪音）
	quotaByTenant := make(map[uint32][]*tenantquotaRow, 4)
	tenantIDs := make([]uint32, 0, 8)
	seen := map[uint32]bool{}
	for _, q := range quotas {
		if !seen[q.TenantID] {
			seen[q.TenantID] = true
			tenantIDs = append(tenantIDs, q.TenantID)
		}
	}
	sort.Slice(tenantIDs, func(i, j int) bool { return tenantIDs[i] < tenantIDs[j] })
	for _, q := range quotas {
		quotaByTenant[q.TenantID] = append(quotaByTenant[q.TenantID], &tenantquotaRow{typ: string(q.QuotaType), value: q.QuotaValue})
	}

	var hits []QuotaWatermarkHit
	for _, tenantID := range tenantIDs {
		t, terr := s.repo.GetTenantStatus(sysCtx, tenantID)
		if terr != nil {
			s.log.Errorf("quota watermark scan: tenant %d lookup failed: %s", tenantID, terr.Error())
			continue
		}
		if t == nil || (t.Status != nil && *t.Status != tenant.StatusOn) {
			continue
		}

		tenantName := ""
		if t.Name != nil {
			tenantName = *t.Name
		}
		adminUserID := uint32(0)
		if t.AdminUserID != nil {
			adminUserID = *t.AdminUserID
		}

		// 各维度用量。单维度计量失败记日志（带原始错误）后跳过该维度：
		// 扫描继续、该租户该维度不产生命中。
		usage := map[string]uint64{}
		if n, uerr := s.repo.CountUsersInTenant(sysCtx, tenantID); uerr != nil {
			s.log.Errorf("quota watermark scan: count users for tenant %d failed: %s", tenantID, uerr.Error())
		} else {
			usage["USER_LIMIT"] = n
		}
		if n, serr := s.repo.SumStorageInTenant(sysCtx, tenantID); serr != nil {
			s.log.Errorf("quota watermark scan: sum storage for tenant %d failed: %s", tenantID, serr.Error())
		} else {
			usage["STORAGE"] = n
		}
		if n, aerr := s.repo.CountApiCallsInTenant(sysCtx, tenantID); aerr != nil {
			s.log.Errorf("quota watermark scan: count api calls for tenant %d failed: %s", tenantID, aerr.Error())
		} else {
			usage["API_CALL"] = n
		}
		if n, terr := s.repo.SumAiTokensInTenantSince(sysCtx, tenantID, monthStart); terr != nil {
			s.log.Errorf("quota watermark scan: sum ai tokens for tenant %d failed: %s", tenantID, terr.Error())
		} else {
			usage["AI_TOKENS"] = n
		}

		for _, q := range quotaByTenant[tenantID] {
			if q.typ == "QUOTA_TYPE_UNSPECIFIED" {
				continue
			}
			used, ok := usage[q.typ]
			if !ok {
				continue
			}
			if quotaWatermarkExceeded(used, q.value) {
				hits = append(hits, QuotaWatermarkHit{
					TenantID:    tenantID,
					TenantName:  tenantName,
					AdminUserID: adminUserID,
					QuotaType:   q.typ,
					Used:        used,
					Limit:       q.value,
					RatioPct:    ratioPct(used, q.value),
				})
			}
		}
	}

	s.log.Infof("quota watermark scan: %d tenants scanned, %d watermark hits", len(tenantIDs), len(hits))
	return hits, nil
}

// tenantquotaRow 扫描用轻量行（类型技术名 + 上限）。
type tenantquotaRow struct {
	typ   string
	value uint64
}

// ─── 水位告警 worker（asynq：tenant_quota_watermark_scan，每日 09:00） ───

// AsyncTenantQuotaWatermarkScan 水位扫描任务 handler：
// 扫描 → 命中租户归组（一租户一条消息，多维度命中合并为同一张清单）→
// 按租户管理员偏好语言渲染 → 站内信投递（系统身份，逐收件人落收件箱）。
// 投递失败语义：单租户失败记日志后继续其余租户、不向上报错重试——asynq 重试会
// 整任务重跑，重复投递已成功的租户；全部尝试都失败才报错（此时无任何收件行
// 落库，重试无重复副作用）。
func (s *TenantQuotaService) AsyncTenantQuotaWatermarkScan(taskType string, payload *task.TenantQuotaWatermarkTaskData) error {
	ctx := appViewer.NewSystemViewerContext(context.Background())

	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)

	hits, err := s.ScanQuotaWatermarks(ctx, monthStart)
	if err != nil {
		s.log.Errorf("quota watermark: scan failed: %s", err.Error())
		return err
	}
	if len(hits) == 0 {
		s.log.Infof("quota watermark: no tenant over watermark, nothing to deliver")
		return nil
	}

	grouped := make(map[uint32][]QuotaWatermarkHit)
	adminOf := make(map[uint32]uint32)
	nameOf := make(map[uint32]string)
	for _, h := range hits {
		grouped[h.TenantID] = append(grouped[h.TenantID], h)
		adminOf[h.TenantID] = h.AdminUserID
		nameOf[h.TenantID] = h.TenantName
	}

	if s.internalMessages == nil {
		s.log.Warnf("quota watermark: %d tenants over watermark but internal message sender not registered, no notification sent", len(grouped))
		return nil
	}

	attempted, failed := 0, 0
	for tenantID, tenantHits := range grouped {
		adminID := adminOf[tenantID]
		if adminID == 0 {
			// 无可投递的管理员：留运维日志（平台应当知道的事），不计入投递尝试。
			s.log.Warnf("quota watermark: tenant %d [%s] over watermark in %d dimension(s) but has no admin user, no notification sent",
				tenantID, nameOf[tenantID], len(tenantHits))
			continue
		}
		attempted++

		locale := s.adminLocaleOf(ctx, adminID)
		title := watermarkTitle(locale, nameOf[tenantID])
		content := watermarkBody(locale, nameOf[tenantID], tenantHits)
		if _, derr := s.internalMessages.SendSystemNotification(ctx, title, content, []uint32{adminID}); derr != nil {
			s.log.Errorf("quota watermark: deliver to admin %d of tenant %d failed: %s", adminID, tenantID, derr.Error())
			failed++
		}
	}
	if attempted > 0 && failed == attempted {
		return fmt.Errorf("quota watermark: all %d tenant notifications failed", failed)
	}
	s.log.Infof("quota watermark: %d tenants over watermark, %d delivered", len(grouped), attempted-failed)
	return nil
}

// adminLocaleOf 查租户管理员的偏好语言（未设置/未识别回落中文）。
func (s *TenantQuotaService) adminLocaleOf(ctx context.Context, adminUserID uint32) string {
	svCtx := appViewer.NewSystemViewerContext(ctx)
	u, err := s.repo.GetUserByID(svCtx, adminUserID)
	if err != nil || u == nil || u.Locale == nil || *u.Locale == "" {
		return "zh-cn"
	}
	l := strings.ToLower(strings.TrimSpace(*u.Locale))
	if _, ok := watermarkTables[l]; ok {
		return l
	}
	return "zh-cn"
}

// ─── 水位告警文案（站内信，按租户管理员偏好语言渲染） ───
//
// 语言标签归一按小写匹配（与监控告警双语同规矩：normalizeLocaleTag 归一后
// switch 分支必须写小写）。配额类型名不做翻译（USER_LIMIT 等技术名与用量页/
// 闸门报错口径一致）。未设置或未识别回落中文。

type watermarkCopy struct {
	title     string // %s=租户名
	quotaHead string // %d=水位百分比
	quotaLine string // %s=类型名 %d=已用 %d=上限 %d=百分比
	footer    string
}

var watermarkTables = map[string]watermarkCopy{
	"zh-cn": {
		title:     "租户配额水位告警（%s）",
		quotaHead: "以下配额的用量已达到 %d%% 水位：\n",
		quotaLine: "  - %s：已用 %d / 上限 %d（%d%%）\n",
		footer:    "请尽快检查用量并与平台管理员联系调整配额；用量到达上限后，相关操作将被拒绝。\n",
	},
	"en-us": {
		title:     "Tenant Quota Watermark Alert (%s)",
		quotaHead: "The usage of the following quotas has reached the %d%% watermark:\n",
		quotaLine: "  - %s: used %d / limit %d (%d%%)\n",
		footer:    "Please review your usage soon and contact the platform administrator to adjust the quotas; once a limit is reached, the related operations will be rejected.\n",
	},
}

func watermarkCopyOf(locale string) watermarkCopy {
	if cp, ok := watermarkTables[strings.ToLower(strings.TrimSpace(locale))]; ok {
		return cp
	}
	return watermarkTables["zh-cn"]
}

func watermarkTitle(locale, tenantName string) string {
	return fmt.Sprintf(watermarkCopyOf(locale).title, tenantName)
}

func watermarkBody(locale, tenantName string, hits []QuotaWatermarkHit) string {
	cp := watermarkCopyOf(locale)
	var sb strings.Builder
	sb.WriteString(tenantName + "\n")
	sb.WriteString(fmt.Sprintf(cp.quotaHead, int(math.Round(QuotaWatermarkThreshold*100))))
	for _, h := range hits {
		sb.WriteString(fmt.Sprintf(cp.quotaLine, h.QuotaType, h.Used, h.Limit, h.RatioPct))
	}
	sb.WriteString(cp.footer)
	return sb.String()
}
