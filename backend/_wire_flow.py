import io
import re

# ── PostService：加 flowService 字段/构造 + 事件钩子 ──
p = 'app/core/service/internal/service/post_service.go'
s = io.open(p, encoding='utf-8').read()

old = '''	referenceTracker      *ReferenceTrackingService
	log                   *log.Helper
}'''
new = '''	referenceTracker      *ReferenceTrackingService
	flowSvc               *FlowService
	log                   *log.Helper
}'''
assert s.count(old) == 1, 'post struct'
s = s.replace(old, new)

old = 'func NewPostService(ctx *bootstrap.Context, uc *data.PostRepo, searchService *SearchService, taskService *TaskService, contentModelRepo *data.ContentModelRepo, postVersionRepo *data.PostVersionRepo, webhookService *WebhookService, internalMessageSvc *InternalMessageService, referenceTracker *ReferenceTrackingService) *PostService {'
new = 'func NewPostService(ctx *bootstrap.Context, uc *data.PostRepo, searchService *SearchService, taskService *TaskService, contentModelRepo *data.ContentModelRepo, postVersionRepo *data.PostVersionRepo, webhookService *WebhookService, internalMessageSvc *InternalMessageService, referenceTracker *ReferenceTrackingService, flowSvc *FlowService) *PostService {'
assert s.count(old) == 1, 'post ctor'
s = s.replace(old, new)

old = '''		contentModelRepo:    contentModelRepo,
		referenceTracker:    referenceTracker,
	}
}'''
new = '''		contentModelRepo:    contentModelRepo,
		referenceTracker:    referenceTracker,
		flowSvc:             flowSvc,
	}
}'''
assert s.count(old) == 1, 'post ctor body'
s = s.replace(old, new)

old = '''	// Webhook 事件外发（best-effort）
	s.dispatchWebhookEvent(ctx, task.WebhookEventPostCreated, dto.GetId(), dto.GetStatus())'''
new = '''	// Webhook 事件外发（best-effort）
	s.dispatchWebhookEvent(ctx, task.WebhookEventPostCreated, dto.GetId(), dto.GetStatus())

	// Flow 规则触发（best-effort）
	if s.flowSvc != nil {
		s.flowSvc.DispatchFlow(ctx, task.WebhookEventPostCreated, dto.GetId())
	}'''
assert s.count(old) == 1, 'post create hook'
s = s.replace(old, new)

old = '''	if dto.GetStatus() == contentV1.Post_POST_STATUS_PUBLISHED && prevStatus != contentV1.Post_POST_STATUS_PUBLISHED {
		s.dispatchWebhookEvent(ctx, task.WebhookEventPostPublished, dto.GetId(), dto.GetStatus())
	} else {
		s.dispatchWebhookEvent(ctx, task.WebhookEventPostUpdated, dto.GetId(), dto.GetStatus())
	}'''
new = '''	if dto.GetStatus() == contentV1.Post_POST_STATUS_PUBLISHED && prevStatus != contentV1.Post_POST_STATUS_PUBLISHED {
		s.dispatchWebhookEvent(ctx, task.WebhookEventPostPublished, dto.GetId(), dto.GetStatus())
	} else {
		s.dispatchWebhookEvent(ctx, task.WebhookEventPostUpdated, dto.GetId(), dto.GetStatus())
	}

	// Flow 规则触发：发布转移记 published，其余记 updated（best-effort）
	if s.flowSvc != nil {
		if dto.GetStatus() == contentV1.Post_POST_STATUS_PUBLISHED && prevStatus != contentV1.Post_POST_STATUS_PUBLISHED {
			s.flowSvc.DispatchFlow(ctx, task.WebhookEventPostPublished, dto.GetId())
		} else {
			s.flowSvc.DispatchFlow(ctx, task.WebhookEventPostUpdated, dto.GetId())
		}
	}'''
assert s.count(old) == 1, 'post update hook'
s = s.replace(old, new)

old = '''	// Webhook 事件外发（best-effort）
	s.dispatchWebhookEvent(ctx, task.WebhookEventPostDeleted, req.GetId(), 0)

	return &emptypb.Empty{}, nil
}'''
new = '''	// Webhook 事件外发（best-effort）
	s.dispatchWebhookEvent(ctx, task.WebhookEventPostDeleted, req.GetId(), 0)

	// Flow 规则触发（best-effort）
	if s.flowSvc != nil {
		s.flowSvc.DispatchFlow(ctx, task.WebhookEventPostDeleted, req.GetId())
	}

	return &emptypb.Empty{}, nil
}'''
assert s.count(old) == 1, 'post delete hook'
s = s.replace(old, new)
io.open(p, 'w', encoding='utf-8', newline='').write(s)
print('post service flow hooks done')

# ── PageService ──
p = 'app/core/service/internal/service/page_service.go'
s = io.open(p, encoding='utf-8').read()

old = '''	referenceTracker  *ReferenceTrackingService
	log               *log.Helper
}'''
new = '''	referenceTracker  *ReferenceTrackingService
	flowSvc           *FlowService
	log               *log.Helper
}'''
assert s.count(old) == 1, 'page struct'
s = s.replace(old, new)

old = 'func NewPageService(ctx *bootstrap.Context, uc *data.PageRepo, contentModelRepo *data.ContentModelRepo, pageVersionRepo *data.PageVersionRepo, webhookService *WebhookService, notifySvc *InternalMessageService, referenceTracker *ReferenceTrackingService) *PageService {'
new = 'func NewPageService(ctx *bootstrap.Context, uc *data.PageRepo, contentModelRepo *data.ContentModelRepo, pageVersionRepo *data.PageVersionRepo, webhookService *WebhookService, notifySvc *InternalMessageService, referenceTracker *ReferenceTrackingService, flowSvc *FlowService) *PageService {'
assert s.count(old) == 1, 'page ctor'
s = s.replace(old, new)

old = '''		contentModelRepo: contentModelRepo,
		referenceTracker: referenceTracker,
	}
}'''
new = '''		contentModelRepo: contentModelRepo,
		referenceTracker: referenceTracker,
		flowSvc:          flowSvc,
	}
}'''
assert s.count(old) == 1, 'page ctor body'
s = s.replace(old, new)

old = '''	// Webhook 事件外发（best-effort）
	if s.webhookService != nil && dto.GetId() != 0 {
		s.webhookService.DispatchEvent(ctx, task.WebhookEventPageCreated, dto.GetId())
	}
	return dto, nil'''
new = '''	// Webhook 事件外发（best-effort）
	if s.webhookService != nil && dto.GetId() != 0 {
		s.webhookService.DispatchEvent(ctx, task.WebhookEventPageCreated, dto.GetId())
	}
	// Flow 规则触发（best-effort）
	if s.flowSvc != nil && dto.GetId() != 0 {
		s.flowSvc.DispatchFlow(ctx, task.WebhookEventPageCreated, dto.GetId())
	}
	return dto, nil'''
assert s.count(old) == 1, 'page create hook'
s = s.replace(old, new)

old = '''	// Webhook 事件外发：转为 PUBLISHED 记 published，其余记 updated（best-effort）
	if s.webhookService != nil && dto.GetId() != 0 {
		if dto.GetStatus() == contentV1.Page_PAGE_STATUS_PUBLISHED && prevStatus != contentV1.Page_PAGE_STATUS_PUBLISHED {
			s.webhookService.DispatchEvent(ctx, task.WebhookEventPagePublished, dto.GetId())
		} else {
			s.webhookService.DispatchEvent(ctx, task.WebhookEventPageUpdated, dto.GetId())
		}
	}
	return dto, nil'''
new = '''	// Webhook 事件外发：转为 PUBLISHED 记 published，其余记 updated（best-effort）
	if s.webhookService != nil && dto.GetId() != 0 {
		if dto.GetStatus() == contentV1.Page_PAGE_STATUS_PUBLISHED && prevStatus != contentV1.Page_PAGE_STATUS_PUBLISHED {
			s.webhookService.DispatchEvent(ctx, task.WebhookEventPagePublished, dto.GetId())
		} else {
			s.webhookService.DispatchEvent(ctx, task.WebhookEventPageUpdated, dto.GetId())
		}
	}
	// Flow 规则触发：同 webhook 的 published/updated 判定（best-effort）
	if s.flowSvc != nil && dto.GetId() != 0 {
		if dto.GetStatus() == contentV1.Page_PAGE_STATUS_PUBLISHED && prevStatus != contentV1.Page_PAGE_STATUS_PUBLISHED {
			s.flowSvc.DispatchFlow(ctx, task.WebhookEventPagePublished, dto.GetId())
		} else {
			s.flowSvc.DispatchFlow(ctx, task.WebhookEventPageUpdated, dto.GetId())
		}
	}
	return dto, nil'''
assert s.count(old) == 1, 'page update hook'
s = s.replace(old, new)

old = '''	// Webhook 事件外发（best-effort）
	if s.webhookService != nil && req.GetId() != 0 {
		s.webhookService.DispatchEvent(ctx, task.WebhookEventPageDeleted, req.GetId())
	}
	return &emptypb.Empty{}, nil'''
new = '''	// Webhook 事件外发（best-effort）
	if s.webhookService != nil && req.GetId() != 0 {
		s.webhookService.DispatchEvent(ctx, task.WebhookEventPageDeleted, req.GetId())
	}
	// Flow 规则触发（best-effort）
	if s.flowSvc != nil && req.GetId() != 0 {
		s.flowSvc.DispatchFlow(ctx, task.WebhookEventPageDeleted, req.GetId())
	}
	return &emptypb.Empty{}, nil'''
assert s.count(old) == 1, 'page delete hook'
s = s.replace(old, new)
io.open(p, 'w', encoding='utf-8', newline='').write(s)
print('page service flow hooks done')

# ── asynq server：注册 flow.execute worker ──
p = 'app/core/service/internal/server/asynq_server.go'
s = io.open(p, encoding='utf-8').read()

old = 'webhookService *service.WebhookService, internalMessageService *service.InternalMessageService) *asynq.Server {'
new = 'webhookService *service.WebhookService, internalMessageService *service.InternalMessageService, flowService *service.FlowService) *asynq.Server {'
assert s.count(old) == 1, 'asynq param'
s = s.replace(old, new)

m = re.search(r'if err = asynq\.RegisterSubscriber\(srv, task\.WebhookDeliverTaskType, webhookService\.DeliverWorker\); err != nil \{\n\t\treturn err\n\t\}', s)
assert m, 'deliver block'
block = m.group(0)
s = s.replace(block, block + '\n\tif err = asynq.RegisterSubscriber(srv, task.FlowExecuteTaskType, flowService.ExecuteWorker); err != nil {\n\t\treturn err\n\t}', 1)
io.open(p, 'w', encoding='utf-8', newline='').write(s)
print('asynq flow worker registered')
