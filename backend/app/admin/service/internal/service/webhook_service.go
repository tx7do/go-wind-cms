package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	adminV1 "go-wind-cms/api/gen/go/admin/service/v1"
	webhookV1 "go-wind-cms/api/gen/go/webhook/service/v1"
)

// WebhookService Webhook 订阅管理门面（透传 core）。
type WebhookService struct {
	adminV1.WebhookServiceHTTPServer

	webhookServiceClient webhookV1.WebhookServiceClient
	log                  *log.Helper
}

func NewWebhookService(ctx *bootstrap.Context, webhookServiceClient webhookV1.WebhookServiceClient) *WebhookService {
	return &WebhookService{
		log:                  ctx.NewLoggerHelper("webhook/service/admin-service"),
		webhookServiceClient: webhookServiceClient,
	}
}

func (s *WebhookService) List(ctx context.Context, req *paginationV1.PagingRequest) (*webhookV1.ListWebhookResponse, error) {
	return s.webhookServiceClient.List(ctx, req)
}

func (s *WebhookService) Get(ctx context.Context, req *webhookV1.GetWebhookRequest) (*webhookV1.Webhook, error) {
	return s.webhookServiceClient.Get(ctx, req)
}

func (s *WebhookService) Create(ctx context.Context, req *webhookV1.CreateWebhookRequest) (*webhookV1.Webhook, error) {
	return s.webhookServiceClient.Create(ctx, req)
}

func (s *WebhookService) Update(ctx context.Context, req *webhookV1.UpdateWebhookRequest) (*webhookV1.Webhook, error) {
	return s.webhookServiceClient.Update(ctx, req)
}

func (s *WebhookService) Delete(ctx context.Context, req *webhookV1.DeleteWebhookRequest) (*emptypb.Empty, error) {
	return s.webhookServiceClient.Delete(ctx, req)
}

func (s *WebhookService) TestPing(ctx context.Context, req *webhookV1.TestWebhookPingRequest) (*webhookV1.WebhookDelivery, error) {
	return s.webhookServiceClient.TestPing(ctx, req)
}

func (s *WebhookService) ListDeliveries(ctx context.Context, req *webhookV1.ListWebhookDeliveriesRequest) (*webhookV1.ListWebhookDeliveriesResponse, error) {
	return s.webhookServiceClient.ListDeliveries(ctx, req)
}
