package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	aiV1 "go-wind-cms/api/gen/go/ai/service/v1"
	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"

	"go-wind-cms/pkg/crypto"
	"go-wind-cms/pkg/middleware/auth"
)

// AiProviderService AI 模型提供商管理（BFF 包装）：
// 操作人注入 + api_key 就地加密，其余透传 core。
type AiProviderService struct {
	adminV1.AiProviderServiceHTTPServer

	aiProviderServiceClient aiV1.AiProviderServiceClient
	log                     *log.Helper
}

func NewAiProviderService(ctx *bootstrap.Context, aiProviderServiceClient aiV1.AiProviderServiceClient) *AiProviderService {
	return &AiProviderService{
		log:                     ctx.NewLoggerHelper("ai_provider/service/admin-service"),
		aiProviderServiceClient: aiProviderServiceClient,
	}
}

func (s *AiProviderService) List(ctx context.Context, req *paginationV1.PagingRequest) (*aiV1.ListAiProviderResponse, error) {
	return s.aiProviderServiceClient.List(ctx, req)
}

func (s *AiProviderService) Get(ctx context.Context, req *aiV1.GetAiProviderRequest) (*aiV1.AiProvider, error) {
	return s.aiProviderServiceClient.Get(ctx, req)
}

// Create 创建提供商；api_key 明文经全局加密器（未启用则明文落库——与全局加密
// 既有语义一致）加密后落库，hint 只存脱敏形态。
func (s *AiProviderService) Create(ctx context.Context, req *aiV1.CreateAiProviderRequest) (*emptypb.Empty, error) {
	if req == nil || req.Data == nil {
		return nil, adminV1.ErrorBadRequest("invalid parameter")
	}
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	req.Data.CreatedBy = trans.Ptr(operator.UserId)
	req.Data.UpdatedBy = trans.Ptr(operator.UserId)

	if err = s.sealApiKey(req.Data); err != nil {
		return nil, err
	}

	return s.aiProviderServiceClient.Create(ctx, req)
}

func (s *AiProviderService) Update(ctx context.Context, req *aiV1.UpdateAiProviderRequest) (*emptypb.Empty, error) {
	if req == nil || req.Data == nil {
		return nil, adminV1.ErrorBadRequest("invalid parameter")
	}
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	req.Data.UpdatedBy = trans.Ptr(operator.GetUserId())

	if err = s.sealApiKey(req.Data); err != nil {
		return nil, err
	}

	return s.aiProviderServiceClient.Update(ctx, req)
}

func (s *AiProviderService) Delete(ctx context.Context, req *aiV1.DeleteAiProviderRequest) (*emptypb.Empty, error) {
	return s.aiProviderServiceClient.Delete(ctx, req)
}

// sealApiKey 把明文 api_key 加密为密文并生成脱敏提示；空值跳过（更新未改密钥）。
func (s *AiProviderService) sealApiKey(data *aiV1.AiProvider) error {
	plain := data.GetApiKey()
	if plain == "" {
		data.ApiKey = nil
		return nil
	}

	encrypted, err := crypto.EncryptIfNeeded(plain)
	if err != nil {
		s.log.Errorf("encrypt ai provider api key failed: %v", err)
		return adminV1.ErrorInternalServerError("encrypt api key failed")
	}
	data.ApiKey = trans.Ptr(encrypted)
	data.ApiKeyHint = trans.Ptr(maskApiKey(plain))
	return nil
}

// maskApiKey 脱敏：保留前 3 后 4，其余以 *** 概括；过短则整体 ***。
func maskApiKey(key string) string {
	if len(key) <= 7 {
		return "***"
	}
	return key[:3] + "***" + key[len(key)-4:]
}
