package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	appViewer "go-wind-cms/pkg/entgo/viewer"
	"go-wind-cms/pkg/social"

	"go-wind-cms/app/core/service/internal/data"

	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
	socialV1 "go-wind-cms/api/gen/go/social/service/v1"

	"go-wind-cms/pkg/crypto"
	_ "go-wind-cms/pkg/social/bluesky"    // 注册 Bluesky 平台适配器
	_ "go-wind-cms/pkg/social/devto"      // 注册 DEV.to 平台适配器
	_ "go-wind-cms/pkg/social/mastodon"   // 注册 Mastodon 平台适配器
	"go-wind-cms/pkg/task"
)

// 未来的兜底扫描周期（P1）：改期后的旧 DELAY 任务到点时若发现 publish_time
// 已被推迟超过该窗口则自弃，交由改期时新入队的任务执行。
const socialStaleTaskSkipWindow = 2 * time.Minute

// SocialService 社交分发服务（core）。
// 职责边界：账号凭据的验证/加解密、分发的排期入队、发布 worker 的状态机推进、
// 统计回拉；平台 HTTP 细节全部在 pkg/social/<platform> 适配器内。
// socialPublishMaxAttempts 发布投递总次数上限（首次 + asynq MaxRetry 3），
// 失败次数达到该值即判定最终失败并通知操作者。
const socialPublishMaxAttempts = 4

type SocialService struct {
	socialV1.UnimplementedSocialServiceServer

	log *log.Helper

	accountRepo        *data.SocialAccountRepo
	distributionRepo   *data.SocialDistributionRepo
	statsHistoryRepo   *data.SocialStatsHistoryRepo
	postRepo           *data.PostRepo
	taskService        *TaskService
	internalMessageSvc *InternalMessageService
}

func NewSocialService(
	ctx *bootstrap.Context,
	accountRepo *data.SocialAccountRepo,
	distributionRepo *data.SocialDistributionRepo,
	statsHistoryRepo *data.SocialStatsHistoryRepo,
	postRepo *data.PostRepo,
	taskService *TaskService,
	internalMessageSvc *InternalMessageService,
) *SocialService {
	return &SocialService{
		log:                ctx.NewLoggerHelper("social/service/core-service"),
		accountRepo:        accountRepo,
		distributionRepo:   distributionRepo,
		statsHistoryRepo:   statsHistoryRepo,
		postRepo:           postRepo,
		taskService:        taskService,
		internalMessageSvc: internalMessageSvc,
	}
}

// ── 社交账号 ──

func (s *SocialService) ListSocialAccounts(ctx context.Context, req *paginationV1.PagingRequest) (*socialV1.ListSocialAccountsResponse, error) {
	resp, err := s.accountRepo.List(ctx, req)
	if err != nil {
		return nil, err
	}
	sanitizeAccounts(resp.GetItems())
	return resp, nil
}

func (s *SocialService) GetSocialAccount(ctx context.Context, req *socialV1.GetSocialAccountRequest) (*socialV1.SocialAccount, error) {
	dto, err := s.accountRepo.Get(ctx, req)
	if err != nil {
		return nil, err
	}
	sanitizeAccounts([]*socialV1.SocialAccount{dto})
	return dto, nil
}

func (s *SocialService) CreateSocialAccount(ctx context.Context, req *socialV1.CreateSocialAccountRequest) (*socialV1.SocialAccount, error) {
	if req == nil || req.Data == nil {
		return nil, socialV1.ErrorBadRequest("invalid parameter")
	}
	data := req.Data

	if data.GetPlatform() == socialV1.Platform_PLATFORM_UNSPECIFIED {
		return nil, socialV1.ErrorBadRequest("platform is required")
	}
	provider, err := social.GetProvider(data.GetPlatform().String())
	if err != nil {
		return nil, socialV1.ErrorPlatformUnsupported("platform not supported yet")
	}

	// 凭据必须可解析（JSON 对象）
	creds, err := social.ParseCredentials(data.GetCredentials())
	if err != nil {
		return nil, socialV1.ErrorCredentialsMissing("credentials must be a JSON object string")
	}
	// Bluesky 的 handle 列是唯一事实源：列值覆盖凭据内的 handle 键
	if data.GetPlatform() == socialV1.Platform_PLATFORM_BLUESKY && data.GetHandle() != "" {
		creds["handle"] = data.GetHandle()
	}
	// 保存即验证：杜绝坏凭据入库（发布时才暴雷）
	if err = provider.ValidateCredentials(ctx, creds); err != nil {
		return nil, s.mapProviderError(err, "validate credentials")
	}

	mergedJSON, merr := social.MarshalCredentials(creds)
	if merr != nil {
		return nil, socialV1.ErrorInternalServerError("marshal credentials failed")
	}
	encrypted, err := crypto.EncryptIfNeeded(mergedJSON)
	if err != nil {
		s.log.Errorf("encrypt credentials failed: %v", err)
		return nil, socialV1.ErrorInternalServerError("encrypt credentials failed")
	}
	data.Credentials = &encrypted
	data.CredentialHint = trans.Ptr(buildCredentialHint(creds))

	created, err := s.accountRepo.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	sanitizeAccounts([]*socialV1.SocialAccount{created})
	return created, nil
}

func (s *SocialService) UpdateSocialAccount(ctx context.Context, req *socialV1.UpdateSocialAccountRequest) (*socialV1.SocialAccount, error) {
	if req == nil || req.Data == nil {
		return nil, socialV1.ErrorBadRequest("invalid parameter")
	}
	data := req.Data

	// 传入了新凭据（非空 JSON）→ 增量合并后验证并加密；是否真正落库由 update_mask 决定。
	// 动态凭据表单只回传填写的键：未填的旧键从已存凭据保留（部分更新不丢密钥）。
	if data.GetCredentials() != "" {
		ensureMaskPath(req, "credentials")
		ensureMaskPath(req, "credential_hint")

		provider, err := social.GetProvider(data.GetPlatform().String())
		if err != nil {
			return nil, socialV1.ErrorPlatformUnsupported("platform not supported yet")
		}
		newCreds, err := social.ParseCredentials(data.GetCredentials())
		if err != nil {
			return nil, socialV1.ErrorCredentialsMissing("credentials must be a JSON object string")
		}

		existing, err := s.accountRepo.Get(ctx, &socialV1.GetSocialAccountRequest{Id: req.GetId()})
		if err != nil {
			return nil, err
		}
		merged := social.Credentials{}
		if raw := existing.GetCredentials(); raw != "" {
			plain, derr := crypto.DecryptIfNeeded(raw)
			if derr != nil {
				s.log.Errorf("decrypt credentials failed (account=%d): %v", req.GetId(), derr)
				return nil, socialV1.ErrorInternalServerError("decrypt credentials failed")
			}
			if oldCreds, perr := social.ParseCredentials(plain); perr == nil {
				merged = oldCreds
			}
		}
		for k, v := range newCreds {
			if v != "" {
				merged[k] = v
			}
		}
		// Bluesky 的 handle 列是唯一事实源：列值覆盖凭据内的 handle 键
		if data.GetPlatform() == socialV1.Platform_PLATFORM_BLUESKY && data.GetHandle() != "" {
			merged["handle"] = data.GetHandle()
		}

		if err = provider.ValidateCredentials(ctx, merged); err != nil {
			return nil, s.mapProviderError(err, "validate credentials")
		}
		mergedJSON, merr := social.MarshalCredentials(merged)
		if merr != nil {
			return nil, socialV1.ErrorInternalServerError("marshal credentials failed")
		}
		encrypted, err := crypto.EncryptIfNeeded(mergedJSON)
		if err != nil {
			s.log.Errorf("encrypt credentials failed: %v", err)
			return nil, socialV1.ErrorInternalServerError("encrypt credentials failed")
		}
		data.Credentials = &encrypted
		data.CredentialHint = trans.Ptr(buildCredentialHint(merged))
	} else if maskHasPath(req.GetUpdateMask(), "handle") {
		// 只改 handle 没重输密码 → 同步改写已存凭据里的 handle 键（handle 列是唯一事实源，
		// 不改会让登录标识与展示 handle 脱节）。凭据密文无法在 BFF 侧改写，必须由 core 完成。
		existing, err := s.accountRepo.Get(ctx, &socialV1.GetSocialAccountRequest{Id: req.GetId()})
		if err != nil {
			return nil, err
		}
		if raw := existing.GetCredentials(); raw != "" {
			plain, derr := crypto.DecryptIfNeeded(raw)
			if derr != nil {
				s.log.Errorf("decrypt credentials failed (account=%d): %v", req.GetId(), derr)
				return nil, socialV1.ErrorInternalServerError("decrypt credentials failed")
			}
			if creds, perr := social.ParseCredentials(plain); perr == nil {
				if data.GetHandle() != "" {
					creds["handle"] = data.GetHandle()
				}
				if merged, merr := social.MarshalCredentials(creds); merr == nil {
					if enc, eerr := crypto.EncryptIfNeeded(merged); eerr == nil {
						data.Credentials = &enc
						ensureMaskPath(req, "credentials")
					}
				}
			}
		}
	}

	updated, err := s.accountRepo.Update(ctx, req)
	if err != nil {
		return nil, err
	}
	sanitizeAccounts([]*socialV1.SocialAccount{updated})
	return updated, nil
}

func (s *SocialService) DeleteSocialAccount(ctx context.Context, req *socialV1.DeleteSocialAccountRequest) (*emptypb.Empty, error) {
	if req == nil {
		return nil, socialV1.ErrorBadRequest("invalid parameter")
	}

	// 引用守卫：名下还有未完结分发时禁止删除（先删/改绑分发）
	if n, err := s.distributionRepo.CountUnfinishedByAccount(ctx, req.GetId()); err != nil {
		return nil, socialV1.ErrorInternalServerError("check account references failed")
	} else if n > 0 {
		return nil, socialV1.ErrorConflict("account has %d unfinished distributions", n)
	}

	if err := s.accountRepo.Delete(ctx, req); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// ── 内容分发 ──

func (s *SocialService) ListSocialDistributions(ctx context.Context, req *paginationV1.PagingRequest) (*socialV1.ListSocialDistributionsResponse, error) {
	resp, err := s.distributionRepo.List(ctx, req)
	if err != nil {
		return nil, err
	}
	s.enrichDistributions(ctx, resp.GetItems())
	return resp, nil
}

func (s *SocialService) GetSocialDistribution(ctx context.Context, req *socialV1.GetSocialDistributionRequest) (*socialV1.SocialDistribution, error) {
	dto, err := s.distributionRepo.Get(ctx, req)
	if err != nil {
		return nil, err
	}
	s.enrichDistributions(ctx, []*socialV1.SocialDistribution{dto})
	return dto, nil
}

func (s *SocialService) CreateSocialDistribution(ctx context.Context, req *socialV1.CreateSocialDistributionRequest) (*socialV1.SocialDistribution, error) {
	if req == nil || req.Data == nil {
		return nil, socialV1.ErrorBadRequest("invalid parameter")
	}
	data := req.Data

	if data.GetPlatform() == socialV1.Platform_PLATFORM_UNSPECIFIED {
		return nil, socialV1.ErrorBadRequest("platform is required")
	}
	if data.GetAccountId() == 0 {
		return nil, socialV1.ErrorBadRequest("accountId is required")
	}
	if strings.TrimSpace(data.GetOverrideText()) == "" {
		return nil, socialV1.ErrorBadRequest("overrideText is required")
	}
	// DEV.to 以"文章"为单位必须有标题：显式给 title，或关联文章（worker 会取文章标题兜底）
	if data.GetPlatform() == socialV1.Platform_PLATFORM_DEVTO &&
		strings.TrimSpace(data.GetTitle()) == "" && data.GetPostId() == 0 {
		return nil, socialV1.ErrorBadRequest("title or postId is required for DEV.to")
	}

	// 账号必须在本租户内、平台匹配且启用（Get 走租户过滤，跨租户引用直接 404）
	acct, err := s.accountRepo.Get(ctx, &socialV1.GetSocialAccountRequest{Id: data.GetAccountId()})
	if err != nil {
		return nil, socialV1.ErrorAccountNotFound("social account not found")
	}
	if acct.GetPlatform() != data.GetPlatform() {
		return nil, socialV1.ErrorBadRequest("platform mismatch with account")
	}
	if !acct.GetIsEnabled() {
		return nil, socialV1.ErrorAccountDisabled("social account is disabled")
	}

	// 客户端只允许进入 DRAFT / SCHEDULED（缺省 SCHEDULED 即排期发布）
	status := data.GetStatus()
	if status == socialV1.SocialDistribution_STATUS_UNSPECIFIED {
		status = socialV1.SocialDistribution_STATUS_SCHEDULED
	}
	if status != socialV1.SocialDistribution_STATUS_DRAFT &&
		status != socialV1.SocialDistribution_STATUS_SCHEDULED {
		return nil, socialV1.ErrorBadRequest("invalid initial status")
	}
	data.Status = &status

	created, err := s.distributionRepo.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	// SCHEDULED → 排期入队（publish_time 为空或已过 → 立即）
	if created.GetStatus() == socialV1.SocialDistribution_STATUS_SCHEDULED {
		if err = s.enqueuePublish(created); err != nil {
			_, _ = s.distributionRepo.MarkFailed(ctx, created.GetId(), "enqueue failed: "+err.Error())
			return nil, socialV1.ErrorInternalServerError("schedule publish failed")
		}
	}

	s.enrichDistributions(ctx, []*socialV1.SocialDistribution{created})
	return created, nil
}

func (s *SocialService) UpdateSocialDistribution(ctx context.Context, req *socialV1.UpdateSocialDistributionRequest) (*socialV1.SocialDistribution, error) {
	if req == nil || req.Data == nil {
		return nil, socialV1.ErrorBadRequest("invalid parameter")
	}

	existing, err := s.distributionRepo.Get(ctx, &socialV1.GetSocialDistributionRequest{Id: req.GetId()})
	if err != nil {
		return nil, err
	}
	switch existing.GetStatus() {
	case socialV1.SocialDistribution_STATUS_PUBLISHED, socialV1.SocialDistribution_STATUS_PUBLISHING:
		return nil, socialV1.ErrorNotPublishable("published/publishing record is immutable")
	}

	// 改期时不允许把 publish_time 清空（SCHEDULED 记录会被无声搁置）
	maskHasPublishTime := maskHasPath(req.GetUpdateMask(), "publish_time")
	if maskHasPublishTime && req.Data.GetPublishTime() == nil &&
		existing.GetStatus() == socialV1.SocialDistribution_STATUS_SCHEDULED {
		return nil, socialV1.ErrorBadRequest("publish_time cannot be cleared while scheduled")
	}

	// 换绑账号 → 与 Create 同样的账号校验
	if maskHasPath(req.GetUpdateMask(), "account_id") || maskHasPath(req.GetUpdateMask(), "platform") {
		acct, err := s.accountRepo.Get(ctx, &socialV1.GetSocialAccountRequest{Id: req.Data.GetAccountId()})
		if err != nil {
			return nil, socialV1.ErrorAccountNotFound("social account not found")
		}
		if req.Data.GetPlatform() != socialV1.Platform_PLATFORM_UNSPECIFIED && acct.GetPlatform() != req.Data.GetPlatform() {
			return nil, socialV1.ErrorBadRequest("platform mismatch with account")
		}
		if !acct.GetIsEnabled() {
			return nil, socialV1.ErrorAccountDisabled("social account is disabled")
		}
	}

	updated, err := s.distributionRepo.Update(ctx, req)
	if err != nil {
		return nil, err
	}

	// 改期 → 重新入队（旧 DELAY 任务到点自检 publish_time 已推迟会自弃）
	if maskHasPublishTime && updated.GetStatus() == socialV1.SocialDistribution_STATUS_SCHEDULED {
		if err = s.enqueuePublish(updated); err != nil {
			_, _ = s.distributionRepo.MarkFailed(ctx, updated.GetId(), "enqueue failed: "+err.Error())
			return nil, socialV1.ErrorInternalServerError("reschedule publish failed")
		}
	}

	s.enrichDistributions(ctx, []*socialV1.SocialDistribution{updated})
	return updated, nil
}

func (s *SocialService) DeleteSocialDistribution(ctx context.Context, req *socialV1.DeleteSocialDistributionRequest) (*emptypb.Empty, error) {
	if err := s.distributionRepo.Delete(ctx, req); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// PublishSocialDistribution 立即发布：状态闸门 + 零时刻入队，实际平台调用在 worker。
func (s *SocialService) PublishSocialDistribution(ctx context.Context, req *socialV1.PublishSocialDistributionRequest) (*socialV1.SocialDistribution, error) {
	if req == nil || req.GetId() == 0 {
		return nil, socialV1.ErrorBadRequest("invalid parameter")
	}

	d, err := s.distributionRepo.Get(ctx, &socialV1.GetSocialDistributionRequest{Id: req.GetId()})
	if err != nil {
		return nil, err
	}
	if d.GetStatus() == socialV1.SocialDistribution_STATUS_PUBLISHED {
		return nil, socialV1.ErrorNotPublishable("already published")
	}

	// 快速失败预检：账号缺失/停用没必要入队等 worker 报错
	acct, err := s.accountRepo.GetByID(ctx, d.GetAccountId())
	if err != nil {
		return nil, socialV1.ErrorAccountNotFound("social account not found")
	}
	if !acct.GetIsEnabled() {
		return nil, socialV1.ErrorAccountDisabled("social account is disabled")
	}

	locked, err := s.distributionRepo.TryLockPublishing(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	if !locked {
		return nil, socialV1.ErrorNotPublishable("not in publishable state")
	}

	if err = s.enqueuePublish(d); err != nil {
		_, _ = s.distributionRepo.MarkFailed(ctx, req.GetId(), "enqueue failed: "+err.Error())
		return nil, socialV1.ErrorInternalServerError("enqueue publish failed")
	}

	return s.distributionRepo.GetByID(ctx, req.GetId())
}

// RefreshSocialDistributionStats 回拉平台互动统计（同步调用，页面刷新按钮）。
func (s *SocialService) RefreshSocialDistributionStats(ctx context.Context, req *socialV1.RefreshSocialDistributionStatsRequest) (*socialV1.SocialDistribution, error) {
	if req == nil || req.GetId() == 0 {
		return nil, socialV1.ErrorBadRequest("invalid parameter")
	}

	d, err := s.distributionRepo.Get(ctx, &socialV1.GetSocialDistributionRequest{Id: req.GetId()})
	if err != nil {
		return nil, err
	}
	if d.GetStatus() != socialV1.SocialDistribution_STATUS_PUBLISHED || d.GetExternalId() == "" {
		return nil, socialV1.ErrorNotPublishable("stats available only for published records")
	}

	creds, provider, err := s.resolveProvider(ctx, d)
	if err != nil {
		return nil, err
	}

	stats, err := provider.FetchStats(ctx, creds, d.GetExternalId(), d.GetExternalUri())
	if err != nil {
		return nil, s.mapProviderError(err, "fetch stats")
	}

	now := time.Now()
	if err = s.distributionRepo.UpdateStats(ctx, req.GetId(), stats.Likes, stats.Replies, stats.Reposts, now); err != nil {
		return nil, err
	}
	s.recordStatsHistory(ctx, req.GetId(), stats, now)

	return s.distributionRepo.GetByID(ctx, req.GetId())
}

// ── 发布 worker（asynq 订阅者）──

// PublishDistributionWorker 是 asynq "social.distribution.publish" 任务的 handler。
//
// 上下文与安全模型（与 search.reindex 一致）：
//   - SystemViewer 跨租户读 DB（payload 不含敏感字段，凭据在 worker 内解密）
//   - 永久性失败（记录/账号不存在、凭据无效、文案超长、账号停用）标记 FAILED 后
//     返回 nil，不烧 asynq 重试
//   - 暂时性失败（网络/平台 5xx/限流）标记 FAILED 后返回 err，由 asynq MaxRetry 重试
//     （每次重试会重新加载记录并以 FAILED 状态重新进入可发布集合）
func (s *SocialService) PublishDistributionWorker(taskType string, payload *task.SocialDistributionPublishPayload) error {
	if payload == nil || payload.DistributionID == 0 {
		s.log.Warnf("social publish worker: nil payload")
		return nil
	}
	// payload.TenantID 仅供排查展示，DB 记录的 tenant_id 为准（与 search.reindex 口径一致）
	s.log.Infof("social publish worker: distribution=%d tenant=%d", payload.DistributionID, payload.TenantID)

	ctx := appViewer.NewSystemViewerContext(context.Background())

	d, err := s.distributionRepo.GetByID(ctx, payload.DistributionID)
	if err != nil {
		s.log.Errorf("social publish worker: distribution %d not loadable: %v", payload.DistributionID, err)
		return nil // 记录不存在/已删：不重试
	}

	// 状态闸门
	if d.GetStatus() == socialV1.SocialDistribution_STATUS_PUBLISHED {
		s.log.Infof("social publish worker: distribution %d already published, skip", d.GetId())
		return nil
	}
	if pt := d.GetPublishTime(); pt != nil &&
		d.GetStatus() == socialV1.SocialDistribution_STATUS_SCHEDULED &&
		pt.AsTime().After(time.Now().Add(socialStaleTaskSkipWindow)) {
		// 改期后的旧任务：到点时发现排期已被推迟 → 自弃，由改期时新入队的任务执行
		s.log.Infof("social publish worker: distribution %d rescheduled to %s, skip stale task",
			d.GetId(), pt.AsTime())
		return nil
	}

	// 账号侧永久性检查
	acct, err := s.accountRepo.GetByID(ctx, d.GetAccountId())
	if err != nil {
		_, _ = s.distributionRepo.MarkFailed(ctx, d.GetId(), "social account missing")
		return nil
	}
	if !acct.GetIsEnabled() {
		_, _ = s.distributionRepo.MarkFailed(ctx, d.GetId(), "social account is disabled")
		return nil
	}
	if acct.GetPlatform() != d.GetPlatform() {
		_, _ = s.distributionRepo.MarkFailed(ctx, d.GetId(), "platform mismatch with account")
		return nil
	}

	provider, err := social.GetProvider(d.GetPlatform().String())
	if err != nil {
		_, _ = s.distributionRepo.MarkFailed(ctx, d.GetId(), err.Error())
		return nil
	}

	creds, err := s.decryptCredentials(acct)
	if err != nil {
		s.distributionRepo.MarkFailed(ctx, d.GetId(), "credentials unavailable: "+err.Error())
		return nil
	}

	s.distributionRepo.MarkPublishing(ctx, d.GetId())

	// 标题：显式 title 优先，缺省取关联文章标题（DEV.to 等平台需要）
	title := d.GetTitle()
	if title == "" && d.GetPostId() != 0 {
		if post, perr := s.postRepo.Get(ctx, &contentV1.GetPostRequest{
			QueryBy: &contentV1.GetPostRequest_Id{Id: d.GetPostId()},
		}); perr == nil {
			title = firstPostTitle(post)
		}
	}

	res, err := provider.Publish(ctx, creds, &social.PublishInput{
		Text:  d.GetOverrideText(),
		Title: title,
	})
	if err != nil {
		attempts, merr := s.distributionRepo.MarkFailed(ctx, d.GetId(), err.Error())
		if merr != nil {
			return merr
		}
		s.log.Errorf("social publish worker: publish distribution %d failed (attempt=%d): %v", d.GetId(), attempts, err)
		final := attempts >= socialPublishMaxAttempts
		if errors.Is(err, social.ErrInvalidCredentials) || errors.Is(err, social.ErrContentTooLong) {
			s.notifyPublishResult(d, acct.GetName(), false, err.Error()) // 永久性失败：即最终失败
			return nil                                                   // 重试无意义
		}
		if final {
			s.notifyPublishResult(d, acct.GetName(), false, err.Error()) // 重试穷尽：最终失败
		}
		return err // 暂时性失败：asynq 重试
	}

	if err = s.distributionRepo.MarkPublished(ctx, d.GetId(), res.ExternalID, res.ExternalURI); err != nil {
		return err
	}
	s.notifyPublishResult(d, acct.GetName(), true, res.ExternalURI)

	s.log.Infof("social publish worker: distribution %d published to %s", d.GetId(), res.ExternalURI)
	return nil
}

// notifyPublishResult 发布结果站内信（借鉴 go-wind-admin notifications 的
// "结果触达"设计）：成功即时通知；失败只在最终失败时通知一次（重试期间不打扰，
// 状态在列表页可见）。通知失败不影响发布结果，仅记日志。
func (s *SocialService) notifyPublishResult(d *socialV1.SocialDistribution, accountName string, success bool, detail string) {
	if s.internalMessageSvc == nil || d.GetCreatedBy() == 0 {
		return
	}

	platform := strings.TrimPrefix(d.GetPlatform().String(), "PLATFORM_")
	account := accountName
	if account == "" {
		account = fmt.Sprintf("#%d", d.GetAccountId())
	}

	var title, content string
	if success {
		title = "社交发布成功"
		content = fmt.Sprintf("已成功发布到 %s（账号：%s）。", platform, account)
		if detail != "" {
			content += "\n查看：" + detail
		}
	} else {
		title = "社交发布失败"
		content = fmt.Sprintf("发布到 %s（账号：%s）失败：%s。\n可在「社交分发」列表中修复后点击「立即发布」重试。",
			platform, account, detail)
	}

	s.internalMessageSvc.NotifyUser(context.Background(), d.GetCreatedBy(), title, content)
}

// ── 内部工具 ──

// enqueuePublish 把分发记录入队发布：publish_time 在未来 → 按时刻 DELAY；
// 为空或已过 → 零时刻立即执行。
func (s *SocialService) enqueuePublish(d *socialV1.SocialDistribution) error {
	var at time.Time
	if pt := d.GetPublishTime(); pt != nil && pt.IsValid() {
		at = pt.AsTime()
		if !at.After(time.Now()) {
			at = time.Time{} // 排期已过 → 立即
		}
	}
	return s.taskService.EnqueueSocialDistributionPublish(&task.SocialDistributionPublishPayload{
		DistributionID: d.GetId(),
	}, at)
}

// resolveProvider 解密凭据并解析平台适配器。
func (s *SocialService) resolveProvider(ctx context.Context, d *socialV1.SocialDistribution) (social.Credentials, social.Provider, error) {
	acct, err := s.accountRepo.GetByID(ctx, d.GetAccountId())
	if err != nil {
		return nil, nil, socialV1.ErrorAccountNotFound("social account not found")
	}
	provider, err := social.GetProvider(d.GetPlatform().String())
	if err != nil {
		return nil, nil, socialV1.ErrorPlatformUnsupported("platform not supported yet")
	}
	creds, err := s.decryptCredentials(acct)
	if err != nil {
		return nil, nil, err
	}
	return creds, provider, nil
}

// decryptCredentials 解密账号凭据（只在本服务内存中出现，不落日志）。
func (s *SocialService) decryptCredentials(acct *socialV1.SocialAccount) (social.Credentials, error) {
	raw := acct.GetCredentials()
	if raw == "" {
		return nil, socialV1.ErrorCredentialsMissing("account credentials missing")
	}
	plain, err := crypto.DecryptIfNeeded(raw)
	if err != nil {
		s.log.Errorf("decrypt credentials failed (account=%d): %v", acct.GetId(), err)
		return nil, socialV1.ErrorInternalServerError("decrypt credentials failed")
	}
	creds, err := social.ParseCredentials(plain)
	if err != nil {
		return nil, socialV1.ErrorCredentialsMissing("account credentials malformed")
	}
	return creds, nil
}

// mapProviderError 把适配器的语义化错误映射为 proto 错误。
func (s *SocialService) mapProviderError(err error, action string) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, social.ErrInvalidCredentials):
		return socialV1.ErrorCredentialsInvalid("platform rejected the credentials")
	case errors.Is(err, social.ErrContentTooLong):
		return socialV1.ErrorContentTooLong("content exceeds platform length limit")
	case errors.Is(err, social.ErrRateLimited):
		return socialV1.ErrorTooManyRequests("platform rate limited")
	case errors.Is(err, social.ErrNotFound):
		return socialV1.ErrorNotFound("platform content not found")
	case errors.Is(err, social.ErrStatsUnsupported):
		return socialV1.ErrorStatsUnsupported("platform does not support engagement stats")
	default:
		s.log.Errorf("social %s failed: %v", action, err)
		return socialV1.ErrorPlatformError("platform error")
	}
}

// sanitizeAccounts 清空凭据字段（读取路径永不外泄，只回脱敏 hint）。
func sanitizeAccounts(items []*socialV1.SocialAccount) {
	for _, item := range items {
		if item != nil {
			item.Credentials = nil
		}
	}
}

// enrichDistributions 填充展示字段：账号名称 + 文章标题（best-effort，失败留空）。
func (s *SocialService) enrichDistributions(ctx context.Context, items []*socialV1.SocialDistribution) {
	if len(items) == 0 {
		return
	}

	accountNames := make(map[uint32]string)
	postTitles := make(map[uint32]string)
	for _, item := range items {
		if item == nil {
			continue
		}
		if id := item.GetAccountId(); id != 0 {
			if _, ok := accountNames[id]; !ok {
				accountNames[id] = ""
				if acct, err := s.accountRepo.GetByID(ctx, id); err == nil {
					accountNames[id] = acct.GetName()
				}
			}
			item.AccountName = trans.Ptr(accountNames[id])
		}
		if id := item.GetPostId(); id != 0 {
			if _, ok := postTitles[id]; !ok {
				postTitles[id] = ""
				if post, err := s.postRepo.Get(ctx, &contentV1.GetPostRequest{
					QueryBy: &contentV1.GetPostRequest_Id{Id: id},
				}); err == nil {
					postTitles[id] = firstPostTitle(post)
				}
			}
			item.PostTitle = trans.Ptr(postTitles[id])
		}
	}
}

// firstPostTitle 从文章的多语言翻译中取第一个非空标题（展示字段，best-effort）。
func firstPostTitle(post *contentV1.Post) string {
	for _, tr := range post.GetTranslations() {
		if t := tr.GetTitle(); t != "" {
			return t
		}
	}
	return ""
}

// buildCredentialHint 生成凭据脱敏提示（密文永不回显，这里只给最后 4 位）。
func buildCredentialHint(creds social.Credentials) string {
	for _, key := range []string{"app_password", "api_key", "access_token"} {
		if v := creds.Get(key); v != "" && len(v) >= 4 {
			return fmt.Sprintf("%s ***%s", key, v[len(v)-4:])
		}
	}
	return "credentials ***"
}

// maskHasPath 判断 update_mask 是否包含指定字段。
func maskHasPath(mask *fieldmaskpb.FieldMask, path string) bool {
	if mask == nil {
		return false
	}
	for _, p := range mask.Paths {
		if p == path {
			return true
		}
	}
	return false
}

// ensureMaskPath 保证 update_mask 包含指定字段（更新凭据时强制带上）。
func ensureMaskPath(req *socialV1.UpdateSocialAccountRequest, path string) {
	if req.GetUpdateMask() == nil {
		req.UpdateMask = &fieldmaskpb.FieldMask{}
	}
	if !maskHasPath(req.UpdateMask, path) {
		req.UpdateMask.Paths = append(req.UpdateMask.Paths, path)
	}
}

// PublishSweepWorker 是 asynq "social.distribution.sweep" 周期任务的 handler：
// 兜底扫描到期仍未发布的 SCHEDULED 分发记录并重新入队（自愈网，语义同 post sweep）。
func (s *SocialService) PublishSweepWorker(taskType string, payload *task.PublishSweepPayload) error {
	ctx := appViewer.NewSystemViewerContext(context.Background())

	ids, err := s.distributionRepo.ListDueScheduledIDs(ctx, time.Now())
	if err != nil {
		s.log.Errorf("social publish sweep: list due distributions failed: %v", err)
		return err
	}
	for _, id := range ids {
		if err := s.taskService.EnqueueSocialDistributionPublish(&task.SocialDistributionPublishPayload{
			DistributionID: id,
		}, time.Time{}); err != nil {
			s.log.Errorf("social publish sweep: re-enqueue distribution %d failed: %v", id, err)
		}
	}
	if len(ids) != 0 {
		s.log.Infof("social publish sweep: re-enqueued %d due scheduled distributions", len(ids))
	}
	return nil
}

// 统计周期回拉的扫描窗口与节流参数
const (
	// socialStatsSyncWindow 只回拉近 30 天内发布的记录（更老的统计基本不再变化）
	socialStatsSyncWindow = 30 * 24 * time.Hour
	// socialStatsSyncThrottle 同一条记录两次回拉的最小间隔
	socialStatsSyncThrottle = 2 * time.Hour
	// socialStatsSyncBatchLimit 单次批量回拉的记录上限（防长事务/平台限流）
	socialStatsSyncBatchLimit = 50
)

// StatsSyncWorker 是 asynq "social.distribution.stats.sync" 周期任务的 handler：
// 批量回拉近 30 天内发布、距上次回拉超 2 小时的分发记录的互动统计。
// 单条失败只记日志并继续，不中断批次（下个周期会再次尝试）。
func (s *SocialService) StatsSyncWorker(taskType string, payload *task.StatsSyncPayload) error {
	ctx := appViewer.NewSystemViewerContext(context.Background())
	now := time.Now()

	items, err := s.distributionRepo.ListStatsSyncCandidates(ctx,
		now.Add(-socialStatsSyncWindow),
		now.Add(-socialStatsSyncThrottle),
		socialStatsSyncBatchLimit)
	if err != nil {
		s.log.Errorf("stats sync: list candidates failed: %v", err)
		return err
	}
	if len(items) == 0 {
		return nil
	}

	okCount, failCount := 0, 0
	for _, item := range items {
		creds, provider, rerr := s.resolveProvider(ctx, item)
		if rerr != nil {
			s.log.Warnf("stats sync: resolve provider failed (distribution=%d): %v", item.GetId(), rerr)
			failCount++
			continue
		}

		stats, perr := provider.FetchStats(ctx, creds, item.GetExternalId(), item.GetExternalUri())
		if perr != nil {
			if errors.Is(perr, social.ErrStatsUnsupported) {
				// 平台不提供统计：仅推进节流时间戳，避免每周期重复尝试与告警噪音
				if uerr := s.distributionRepo.TouchStatsSyncedAt(ctx, item.GetId(), now); uerr != nil {
					s.log.Errorf("stats sync: touch synced_at failed (distribution=%d): %v", item.GetId(), uerr)
				}
				continue
			}
			s.log.Warnf("stats sync: fetch stats failed (distribution=%d platform=%s): %v",
				item.GetId(), item.GetPlatform(), perr)
			failCount++
			continue
		}

		if uerr := s.distributionRepo.UpdateStats(ctx, item.GetId(),
			stats.Likes, stats.Replies, stats.Reposts, now); uerr != nil {
			s.log.Errorf("stats sync: update stats failed (distribution=%d): %v", item.GetId(), uerr)
			failCount++
			continue
		}
		s.recordStatsHistory(ctx, item.GetId(), stats, now)
		okCount++
	}

	// 顺手裁剪 90 天前的历史快照（限批，防止表无限膨胀）
	if pruned, perr := s.statsHistoryRepo.PruneOlderThan(ctx, now.Add(-90*24*time.Hour), 1000); perr == nil && pruned != 0 {
		s.log.Infof("stats sync: pruned %d stale history rows", pruned)
	}

	s.log.Infof("stats sync: batch done (candidates=%d ok=%d fail=%d)", len(items), okCount, failCount)
	return nil
}

// recordStatsHistory 追加一行统计快照（趋势数据从今天开始积累；失败不影响主流程）。
func (s *SocialService) recordStatsHistory(ctx context.Context, distributionID uint32, stats *social.Stats, now time.Time) {
	if s.statsHistoryRepo == nil || stats == nil {
		return
	}
	if err := s.statsHistoryRepo.Record(ctx, data.StatsRecord{
		DistributionID: distributionID,
		Likes:          stats.Likes,
		Replies:        stats.Replies,
		Reposts:        stats.Reposts,
		RecordedAt:     now,
	}); err != nil {
		s.log.Errorf("record stats history failed (distribution=%d): %v", distributionID, err)
	}
}
