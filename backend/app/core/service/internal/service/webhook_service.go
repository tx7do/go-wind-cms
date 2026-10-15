package service

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/go-crud/viewer"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	appViewer "go-wind-cms/pkg/entgo/viewer"
	"go-wind-cms/pkg/task"

	"go-wind-cms/app/core/service/internal/data"
	webhookV1 "go-wind-cms/api/gen/go/webhook/service/v1"
)

// webhookHTTPTimeout 投递超时：接收端无响应时快速失败交由 asynq 退避重试。
const webhookHTTPTimeout = 10 * time.Second

// WebhookService 内容事件外发订阅管理 + 投递 worker。
type WebhookService struct {
	webhookV1.UnimplementedWebhookServiceServer

	webhookRepo *data.WebhookRepo
	taskService *TaskService
	log         *log.Helper

	httpClient *http.Client
}

func NewWebhookService(ctx *bootstrap.Context, webhookRepo *data.WebhookRepo, taskService *TaskService) *WebhookService {
	return &WebhookService{
		log:         ctx.NewLoggerHelper("webhook/service/core-service"),
		webhookRepo: webhookRepo,
		taskService: taskService,
		httpClient:  &http.Client{Timeout: webhookHTTPTimeout},
	}
}

// ─── CRUD ───

func (s *WebhookService) List(ctx context.Context, req *paginationV1.PagingRequest) (*webhookV1.ListWebhookResponse, error) {
	items, total, err := s.webhookRepo.List(ctx)
	if err != nil {
		return nil, errors.New(http.StatusInternalServerError, "WEBHOOK_LIST_FAILED", "list webhooks failed")
	}
	return &webhookV1.ListWebhookResponse{Items: items, Total: uint64(total)}, nil
}

func (s *WebhookService) Get(ctx context.Context, req *webhookV1.GetWebhookRequest) (*webhookV1.Webhook, error) {
	if req == nil || req.GetId() == 0 {
		return nil, errors.New(http.StatusBadRequest, "INVALID_PARAMETER", "invalid parameter")
	}
	return s.webhookRepo.Get(ctx, req.GetId())
}

func (s *WebhookService) Create(ctx context.Context, req *webhookV1.CreateWebhookRequest) (*webhookV1.Webhook, error) {
	if req == nil || req.Data == nil || req.Data.GetUrl() == "" {
		return nil, errors.New(http.StatusBadRequest, "INVALID_PARAMETER", "url is required")
	}
	// URL 仅接受 http/https，防 file:/gopher: 等协议被 Worker 端发起
	if u := req.Data.GetUrl(); !bytes.HasPrefix([]byte(u), []byte("http://")) && !bytes.HasPrefix([]byte(u), []byte("https://")) {
		return nil, errors.New(http.StatusBadRequest, "INVALID_PARAMETER", "url must start with http:// or https://")
	}
	return s.webhookRepo.Create(ctx, req.Data)
}

func (s *WebhookService) Update(ctx context.Context, req *webhookV1.UpdateWebhookRequest) (*webhookV1.Webhook, error) {
	if req == nil || req.Data == nil || req.GetId() == 0 {
		return nil, errors.New(http.StatusBadRequest, "INVALID_PARAMETER", "invalid parameter")
	}
	return s.webhookRepo.Update(ctx, req.GetId(), req.Data)
}

func (s *WebhookService) Delete(ctx context.Context, req *webhookV1.DeleteWebhookRequest) (*emptypb.Empty, error) {
	if req == nil || req.GetId() == 0 {
		return nil, errors.New(http.StatusBadRequest, "INVALID_PARAMETER", "invalid parameter")
	}
	if err := s.webhookRepo.Delete(ctx, req.GetId()); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// TestPing 同步发送 ping 事件并返回投递结果（管理端"测试"按钮用）。
func (s *WebhookService) TestPing(ctx context.Context, req *webhookV1.TestWebhookPingRequest) (*webhookV1.WebhookDelivery, error) {
	if req == nil || req.GetId() == 0 {
		return nil, errors.New(http.StatusBadRequest, "INVALID_PARAMETER", "invalid parameter")
	}
	wh, err := s.webhookRepo.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	var tenantID uint32
	if vc, exist := viewer.FromContext(ctx); exist && vc != nil {
		tenantID = uint32(vc.TenantID())
	}

	deliveryID, err := s.webhookRepo.CreateDelivery(ctx, wh.GetId(), "ping", 0, tenantID)
	if err != nil {
		return nil, errors.New(http.StatusInternalServerError, "WEBHOOK_PING_FAILED", "create delivery record failed")
	}

	success, httpStatus, errMsg, duration := s.post(wh.GetUrl(), wh.GetSecret(), "ping", wh.GetId(), 0, tenantID)
	_ = s.webhookRepo.UpdateDeliveryResult(ctx, deliveryID, success, httpStatus, errMsg, duration)

	deliveries, _, err := s.webhookRepo.ListDeliveries(ctx, wh.GetId(), 1)
	if err != nil || len(deliveries) == 0 {
		// 回查失败时返回最小结果，不吞掉投递本身的结果
		return &webhookV1.WebhookDelivery{
			Id:         &deliveryID,
			WebhookId:  wh.Id,
			Event:      trans.Ptr("ping"),
			Status:     trans.Ptr(deliveryStatusString(success)),
			HttpStatus: trans.Ptr(int32(httpStatus)),
			Error:      trans.Ptr(errMsg),
			DurationMs: &duration,
		}, nil
	}
	return deliveries[0], nil
}

// ListDeliveries 投递记录。
func (s *WebhookService) ListDeliveries(ctx context.Context, req *webhookV1.ListWebhookDeliveriesRequest) (*webhookV1.ListWebhookDeliveriesResponse, error) {
	if req == nil {
		return nil, errors.New(http.StatusBadRequest, "INVALID_PARAMETER", "invalid parameter")
	}
	items, total, err := s.webhookRepo.ListDeliveries(ctx, req.GetWebhookId(), int(req.GetLimit()))
	if err != nil {
		return nil, errors.New(http.StatusInternalServerError, "WEBHOOK_DELIVERIES_FAILED", "list webhook deliveries failed")
	}
	return &webhookV1.ListWebhookDeliveriesResponse{Items: items, Total: uint64(total)}, nil
}

// ─── 事件发布（由 PostService/PageService 写路径调用） ───

// DispatchEvent 扇出一个内容事件（best-effort：调度器未注入或入队失败仅记日志，
// 不阻断内容保存主链路）。
func (s *WebhookService) DispatchEvent(ctx context.Context, event string, entityID uint32) {
	if s.taskService == nil || entityID == 0 {
		return
	}
	var tenantID uint32
	if vc, exist := viewer.FromContext(ctx); exist && vc != nil {
		tenantID = uint32(vc.TenantID())
	}
	if err := s.taskService.EnqueueWebhookFanout(&task.WebhookFanoutPayload{
		Event:    event,
		EntityID: entityID,
		TenantID: tenantID,
	}); err != nil {
		s.log.Errorf("dispatch webhook event failed (event=%s entity=%d): %v", event, entityID, err)
	}
}

// ─── asynq Workers ───

// FanoutWorker 按事件列出启用的订阅，逐个入队单订阅投递任务（互不牵连重试）。
func (s *WebhookService) FanoutWorker(taskType string, payload *task.WebhookFanoutPayload) error {
	if payload == nil || payload.Event == "" {
		return nil
	}
	ctx := appViewer.NewSystemViewerContext(context.Background())

	webhooks, err := s.webhookRepo.ListEnabledForEvent(ctx, payload.TenantID, payload.Event)
	if err != nil {
		s.log.Errorf("webhook fanout: list webhooks failed (event=%s): %v", payload.Event, err)
		return err // 查库失败：可重试
	}

	// 顺手裁剪 90 天前的投递记录（限批，防止表无限膨胀；best-effort）
	if pruned, perr := s.webhookRepo.PruneDeliveriesOlderThan(ctx, time.Now().Add(-90*24*time.Hour), 1000); perr == nil && pruned != 0 {
		s.log.Infof("webhook fanout: pruned %d stale delivery rows", pruned)
	}
	for _, wh := range webhooks {
		if err = s.taskService.EnqueueWebhookDeliver(&task.WebhookDeliverPayload{
			WebhookID: wh.GetId(),
			Event:     payload.Event,
			EntityID:  payload.EntityID,
			TenantID:  payload.TenantID,
		}); err != nil {
			s.log.Errorf("webhook fanout: enqueue deliver failed (webhook=%d event=%s): %v",
				wh.GetId(), payload.Event, err)
		}
	}
	if len(webhooks) > 0 {
		s.log.Infof("webhook fanout: event=%s entity=%d -> %d subscriber(s)",
			payload.Event, payload.EntityID, len(webhooks))
	}
	return nil
}

// DeliverWorker 单订阅投递：落记录 → 签名 POST → 回写结果。
// 失败返回 err 交 asynq MaxRetry 退避重试；每次尝试落一条 delivery 记录。
func (s *WebhookService) DeliverWorker(taskType string, payload *task.WebhookDeliverPayload) error {
	if payload == nil || payload.WebhookID == 0 {
		return nil
	}
	ctx := appViewer.NewSystemViewerContext(context.Background())

	wh, err := s.webhookRepo.Get(ctx, payload.WebhookID)
	if err != nil {
		s.log.Errorf("webhook deliver: webhook %d not loadable: %v", payload.WebhookID, err)
		return nil // 订阅已删除：不重试
	}
	if !wh.GetEnabled() {
		s.log.Infof("webhook deliver: webhook %d disabled, skip", payload.WebhookID)
		return nil
	}

	deliveryID, dErr := s.webhookRepo.CreateDelivery(ctx, payload.WebhookID, payload.Event, payload.EntityID, payload.TenantID)
	if dErr != nil {
		return dErr // 记录失败：可重试
	}

	success, httpStatus, errMsg, duration := s.post(wh.GetUrl(), wh.GetSecret(), payload.Event, payload.WebhookID, payload.EntityID, payload.TenantID)
	_ = s.webhookRepo.UpdateDeliveryResult(ctx, deliveryID, success, httpStatus, errMsg, duration)

	if !success {
		err = fmt.Errorf("webhook deliver failed (webhook=%d http=%d): %s", payload.WebhookID, httpStatus, errMsg)
		s.log.Errorf("%v", err)
		return err // asynq 退避重试
	}
	s.log.Infof("webhook deliver: webhook=%d event=%s entity=%d ok (http=%d, %dms)",
		payload.WebhookID, payload.Event, payload.EntityID, httpStatus, duration)
	return nil
}

// post 执行签名投递。body 为最小 envelope（不含内容全文）。
func (s *WebhookService) post(url, secret, event string, webhookID, entityID, tenantID uint32) (bool, int, string, uint32) {
	body, _ := json.Marshal(map[string]any{
		"event":      event,
		"webhook_id": webhookID,
		"entity_id":  entityID,
		"tenant_id":  tenantID,
		"timestamp":  time.Now().UTC().Format(time.RFC3339),
	})

	start := time.Now()
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return false, 0, "invalid url: " + err.Error(), 0
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "GoWind-Content-Hub-Webhook/1.0")
	req.Header.Set("X-Webhook-Event", event)
	req.Header.Set("X-Webhook-Delivery", strconv.FormatUint(uint64(start.UnixNano()), 16))
	if secret != "" {
		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write(body)
		req.Header.Set("X-Webhook-Signature", "sha256="+hex.EncodeToString(mac.Sum(nil)))
	}

	resp, err := s.httpClient.Do(req)
	duration := uint32(time.Since(start).Milliseconds())
	if err != nil {
		return false, 0, err.Error(), duration
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			s.log.Warnf("webhook: close response body failed: %v", cerr)
		}
	}()
	// 丢弃 body，仅消费以复用连接；限制读取量防超大响应
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 4096))

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return false, resp.StatusCode, fmt.Sprintf("non-2xx status %d", resp.StatusCode), duration
	}
	return true, resp.StatusCode, "", duration
}

func deliveryStatusString(success bool) string {
	if success {
		return "DELIVERY_STATUS_SUCCESS"
	}
	return "DELIVERY_STATUS_FAILED"
}
