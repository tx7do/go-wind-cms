package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/go-crud/viewer"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	accesskeyV1 "go-wind-cms/api/gen/go/access_key/service/v1"
	authenticationV1 "go-wind-cms/api/gen/go/authentication/service/v1"

	"go-wind-cms/app/core/service/internal/data"
	"go-wind-cms/app/core/service/internal/data/ent/accesskey"
)

// AK 前缀与长度：AK 公开可见（16 字节 hex），SK 高熵（32 字节 hex）只在创建响应中出现一次。
const (
	accessKeyPrefix = "ak-"
	secretPrefix    = "sk-"
	accessKeyBytes  = 8
	secretBytes     = 32
)

// AccessKeyService OpenAPI 访问凭证（AK/SK）管理 + 令牌交换。
type AccessKeyService struct {
	accesskeyV1.UnimplementedAccessKeyServiceServer

	log           *log.Helper
	repo          *data.AccessKeyRepo
	authenticator *data.Authenticator
	rateLimiter   *data.AccessKeyRateLimiter
}

func NewAccessKeyService(
	ctx *bootstrap.Context,
	repo *data.AccessKeyRepo,
	authenticator *data.Authenticator,
	rateLimiter *data.AccessKeyRateLimiter,
) *AccessKeyService {
	return &AccessKeyService{
		log:           ctx.NewLoggerHelper("access-key/service/core-service"),
		repo:          repo,
		authenticator: authenticator,
		rateLimiter:   rateLimiter,
	}
}

func (s *AccessKeyService) List(ctx context.Context, req *paginationV1.PagingRequest) (*accesskeyV1.ListAccessKeyResponse, error) {
	return s.repo.List(ctx, req)
}

func (s *AccessKeyService) Get(ctx context.Context, req *accesskeyV1.GetAccessKeyRequest) (*accesskeyV1.AccessKey, error) {
	return s.repo.Get(ctx, req)
}

// Create 生成 AK/SK 并落库（SK 只存 SHA-256 摘要），明文 secret 仅本次响应返回。
func (s *AccessKeyService) Create(ctx context.Context, req *accesskeyV1.CreateAccessKeyRequest) (*accesskeyV1.CreateAccessKeyResponse, error) {
	if req == nil || req.Data == nil {
		return nil, errors.BadRequest("ACCESS_KEY", "invalid parameter")
	}

	// 凭证归属创建者的租户：机器令牌将继承该租户的隔离
	var operatorID, tenantID uint32
	if vc, exist := viewer.FromContext(ctx); exist && vc != nil {
		operatorID = uint32(vc.UserID())
		tenantID = uint32(vc.TenantID())
	}
	req.Data.TenantId = trans.Ptr(tenantID)

	accessKey, err := randomToken(accessKeyPrefix, accessKeyBytes)
	if err != nil {
		s.log.Errorf("generate access key failed: %s", err.Error())
		return nil, errors.InternalServer("ACCESS_KEY", "generate access key failed")
	}
	secret, err := randomToken(secretPrefix, secretBytes)
	if err != nil {
		s.log.Errorf("generate secret failed: %s", err.Error())
		return nil, errors.InternalServer("ACCESS_KEY", "generate secret failed")
	}
	secretHash := hashSecret(secret)

	entity, err := s.repo.Create(ctx, req.Data, accessKey, secretHash, operatorID)
	if err != nil {
		return nil, err
	}

	return &accesskeyV1.CreateAccessKeyResponse{
		Data: &accesskeyV1.AccessKey{
			Id:        trans.Ptr(uint32(entity.ID)),
			Name:      entity.Name,
			AccessKey: entity.AccessKey,
			Status:    statusProto(entity.Status),
			ExpiresAt: timestampOf(entity.ExpiresAt),
			TenantId:  entity.TenantID,
			CreatedAt: timestampOf(entity.CreatedAt),
		},
		Secret: secret,
	}, nil
}

// Update 仅允许改名称/状态/过期时间。AK 与 secret 摘要不可变：
// 仓储侧按可编辑字段白名单应用掩码，非法路径（access_key/secret_hash/tenant_id）
// 无论怎么拼掩码都落不进库（上游 e852daba：黑名单放行曾致掩码管线整体报错）。
func (s *AccessKeyService) Update(ctx context.Context, req *accesskeyV1.UpdateAccessKeyRequest) (*emptypb.Empty, error) {
	if req == nil || req.Data == nil {
		return nil, errors.BadRequest("ACCESS_KEY", "invalid parameter")
	}
	if req.GetId() == 0 {
		return nil, errors.BadRequest("ACCESS_KEY", "id is required")
	}

	operatorID, _ := viewerUserIDFromContext(ctx)
	if err := s.repo.Update(ctx, req, operatorID); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *AccessKeyService) Delete(ctx context.Context, req *accesskeyV1.DeleteAccessKeyRequest) (*emptypb.Empty, error) {
	if req == nil || req.GetId() == 0 {
		return nil, errors.BadRequest("ACCESS_KEY", "id is required")
	}
	if err := s.repo.Delete(ctx, req); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// IssueToken 令牌交换（免鉴权）：校验 AK/SK 后签发租户作用域的机器令牌。
// 令牌 userId=0（机器身份）、tenant=凭证归属租户、角色仅 machine——下游的租户隔离
// 照常生效；授权引擎启用后机器令牌默认无任何角色权限。
func (s *AccessKeyService) IssueToken(ctx context.Context, req *accesskeyV1.IssueTokenRequest) (*accesskeyV1.IssueTokenResponse, error) {
	if req == nil || req.GetAccessKey() == "" || req.GetSecret() == "" {
		return nil, errors.BadRequest("ACCESS_KEY", "access key and secret are required")
	}

	// 交换端点免鉴权，做按 IP+AK 的尝试限流（防爆破；client_ip 由 BFF 注入）
	clientIP := req.GetClientIp()
	if locked, lerr := s.rateLimiter.IsLocked(ctx, clientIP, req.GetAccessKey()); lerr == nil && locked {
		return nil, errors.BadRequest("ACCESS_KEY", "too many attempts, try again later")
	}

	entity, err := s.repo.GetByAccessKeyBySystem(ctx, req.GetAccessKey())
	if err != nil {
		// 不区分"不存在"与"密钥错误"，避免泄露 AK 是否有效
		if _, _, _, cerr := s.rateLimiter.CheckAndIncr(ctx, clientIP, req.GetAccessKey()); cerr != nil {
			s.log.Errorf("access key rate limiter incr failed: %s", cerr.Error())
		}
		return nil, errors.BadRequest("ACCESS_KEY", "invalid access key or secret")
	}

	if entity.Status == nil || *entity.Status != accesskey.StatusOn {
		return nil, errors.BadRequest("ACCESS_KEY", "access key is disabled")
	}
	if entity.ExpiresAt != nil && time.Now().After(*entity.ExpiresAt) {
		return nil, errors.BadRequest("ACCESS_KEY", "access key is expired")
	}

	sum := sha256.Sum256([]byte(req.GetSecret()))
	given := hex.EncodeToString(sum[:])
	if entity.SecretHash == nil ||
		subtle.ConstantTimeCompare([]byte(given), []byte(*entity.SecretHash)) != 1 {
		if _, _, _, ierr := s.rateLimiter.CheckAndIncr(ctx, clientIP, req.GetAccessKey()); ierr != nil {
			s.log.Errorf("access key rate limiter incr failed: %s", ierr.Error())
		}
		return nil, errors.BadRequest("ACCESS_KEY", "invalid access key or secret")
	}

	// 交换成功：清除该 IP+AK 的失败计数
	s.rateLimiter.Reset(ctx, clientIP, req.GetAccessKey())

	payload := &authenticationV1.UserTokenPayload{
		UserId:   0, // 机器身份：无对应用户行
		TenantId: entity.TenantID,
		Username: trans.Ptr("ak:" + req.GetAccessKey()),
		Roles:    []string{"machine"}, // authz 中间件要求非空 subject；noop 引擎全放行，启用策略引擎后默认无策略=fail-closed
	}
	accessToken, expires, err := s.authenticator.CreateMachineToken(ctx, payload)
	if err != nil {
		s.log.Errorf("issue machine token failed: %s", err.Error())
		return nil, errors.InternalServer("ACCESS_KEY", "issue token failed")
	}

	// 尽力而为刷新最近使用时间（不阻塞响应）
	s.repo.TouchLastUsedBySystem(context.WithoutCancel(ctx), entity.ID)

	return &accesskeyV1.IssueTokenResponse{
		AccessToken: accessToken,
		ExpiresIn:   uint32(expires.Seconds()),
		TokenType:   "bearer",
	}, nil
}

// ResetSecret 重置密钥（轮换）：生成新 SK 并更新摘要，新 SK 明文仅本次返回。
// 注意：旧 SK 的交换被立即阻断（摘要已变），但此前已签发的机器令牌在过期前仍有效。
func (s *AccessKeyService) ResetSecret(ctx context.Context, req *accesskeyV1.ResetAccessKeySecretRequest) (*accesskeyV1.CreateAccessKeyResponse, error) {
	if req == nil || req.GetId() == 0 {
		return nil, errors.BadRequest("ACCESS_KEY", "id is required")
	}

	secret, err := randomToken(secretPrefix, secretBytes)
	if err != nil {
		s.log.Errorf("generate secret failed: %s", err.Error())
		return nil, errors.InternalServer("ACCESS_KEY", "generate secret failed")
	}
	secretHash := hashSecret(secret)

	if err = s.repo.UpdateSecretHash(ctx, req.GetId(), secretHash); err != nil {
		return nil, err
	}

	entity, err := s.repo.Get(ctx, &accesskeyV1.GetAccessKeyRequest{
		QueryBy: &accesskeyV1.GetAccessKeyRequest_Id{Id: req.GetId()},
	})
	if err != nil {
		return nil, err
	}

	return &accesskeyV1.CreateAccessKeyResponse{
		Data:   entity,
		Secret: secret,
	}, nil
}

func randomToken(prefix string, nBytes int) (string, error) {
	buf := make([]byte, nBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("read random bytes: %w", err)
	}
	return prefix + hex.EncodeToString(buf), nil
}

func hashSecret(secret string) string {
	sum := sha256.Sum256([]byte(secret))
	return hex.EncodeToString(sum[:])
}

func statusProto(st *accesskey.Status) *accesskeyV1.AccessKey_Status {
	if st == nil {
		return nil
	}
	if *st == accesskey.StatusOff {
		v := accesskeyV1.AccessKey_OFF
		return &v
	}
	v := accesskeyV1.AccessKey_ON
	return &v
}

func timestampOf(t *time.Time) *timestamppb.Timestamp {
	if t == nil || t.IsZero() {
		return nil
	}
	return timestamppb.New(*t)
}
