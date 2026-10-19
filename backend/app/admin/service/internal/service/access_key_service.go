package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	ktransport "github.com/go-kratos/kratos/v2/transport"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"

	accesskeyV1 "go-wind-cms/api/gen/go/access_key/service/v1"
	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"

	"go-wind-cms/pkg/middleware/logging"
)

// AccessKeyService OpenAPI 访问凭证（AK/SK）管理（BFF 薄代理：生成/校验/签发都在 core）。
// IssueToken 免鉴权（白名单端点）——它本身即是认证；调用方 IP 在此从 HTTP 层注入。
type AccessKeyService struct {
	adminV1.AccessKeyServiceHTTPServer

	log          *log.Helper
	accessKeyCli accesskeyV1.AccessKeyServiceClient
}

func NewAccessKeyService(
	ctx *bootstrap.Context,
	accessKeyCli accesskeyV1.AccessKeyServiceClient,
) *AccessKeyService {
	l := log.NewHelper(log.With(ctx.GetLogger(), "module", "access-key/service/admin-service"))
	return &AccessKeyService{
		log:          l,
		accessKeyCli: accessKeyCli,
	}
}

// List 访问凭证列表（透传 core）。
func (s *AccessKeyService) List(ctx context.Context, req *paginationV1.PagingRequest) (*accesskeyV1.ListAccessKeyResponse, error) {
	return s.accessKeyCli.List(ctx, req)
}

// Get 访问凭证详情（透传 core）。
func (s *AccessKeyService) Get(ctx context.Context, req *accesskeyV1.GetAccessKeyRequest) (*accesskeyV1.AccessKey, error) {
	return s.accessKeyCli.Get(ctx, req)
}

// Create 创建凭证：core 生成 AK/SK，明文 secret 仅本次响应返回。
func (s *AccessKeyService) Create(ctx context.Context, req *accesskeyV1.CreateAccessKeyRequest) (*accesskeyV1.CreateAccessKeyResponse, error) {
	return s.accessKeyCli.Create(ctx, req)
}

// Update 更新凭证（名称/状态/过期时间；AK 与摘要不可变，core 白名单保护）。
func (s *AccessKeyService) Update(ctx context.Context, req *accesskeyV1.UpdateAccessKeyRequest) (*emptypb.Empty, error) {
	return s.accessKeyCli.Update(ctx, req)
}

// Delete 删除凭证（透传 core）。
func (s *AccessKeyService) Delete(ctx context.Context, req *accesskeyV1.DeleteAccessKeyRequest) (*emptypb.Empty, error) {
	return s.accessKeyCli.Delete(ctx, req)
}

// ResetSecret 重置密钥（轮换）：新 SK 明文仅本次返回，旧 SK 立即失效于交换。
func (s *AccessKeyService) ResetSecret(ctx context.Context, req *accesskeyV1.ResetAccessKeySecretRequest) (*accesskeyV1.CreateAccessKeyResponse, error) {
	return s.accessKeyCli.ResetSecret(ctx, req)
}

// IssueToken 令牌交换（免鉴权）：AK/SK 换机器 JWT。
// core 无法感知 HTTP 层，防爆破限流所需的调用方 IP 在此注入（同 Login 的 client_ip 先例）。
func (s *AccessKeyService) IssueToken(ctx context.Context, req *accesskeyV1.IssueTokenRequest) (*accesskeyV1.IssueTokenResponse, error) {
	if req != nil && req.ClientIp == nil {
		if tr, ok := ktransport.FromServerContext(ctx); ok {
			if htr, hok := tr.(*khttp.Transport); hok && htr.Request() != nil {
				req.ClientIp = trans.Ptr(logging.ClientRealIP(htr.Request()))
			}
		}
	}
	return s.accessKeyCli.IssueToken(ctx, req)
}
