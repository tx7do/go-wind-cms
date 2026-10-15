package service

import (
	"context"
	"strconv"

	"github.com/go-kratos/kratos/v2/log"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	kratosMetadata "github.com/go-kratos/kratos/v2/metadata"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	appV1 "go-wind-cms/api/gen/go/app/service/v1"
	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
	entmiddleware "go-wind-cms/pkg/middleware/ent"
)

// FormSubmissionService 表单提交（公开站点）。
//
// 匿名可调用；租户归属与评论/浏览上报链路同款：按请求 Host 解析租户后经
// 内网可信元数据 x-md-global-tenant-id 传给 core 显式归属。
// 校验（页面已发布/区块 FORM/字段集约束）在 core 侧执行。
type FormSubmissionService struct {
	appV1.FormSubmissionServiceHTTPServer

	formSubmissionClient contentV1.FormSubmissionServiceClient
	tenantResolve        entmiddleware.TenantResolver
	log                  *log.Helper
}

func NewFormSubmissionService(ctx *bootstrap.Context, formSubmissionClient contentV1.FormSubmissionServiceClient, tenantResolve entmiddleware.TenantResolver) *FormSubmissionService {
	return &FormSubmissionService{
		log:                  ctx.NewLoggerHelper("form-submission/service/app-service"),
		formSubmissionClient: formSubmissionClient,
		tenantResolve:        tenantResolve,
	}
}

func (s *FormSubmissionService) Submit(ctx context.Context, req *contentV1.SubmitFormSubmissionRequest) (*contentV1.FormSubmission, error) {
	if s.tenantResolve != nil {
		if hr, ok := khttp.RequestFromServerContext(ctx); ok && hr != nil && hr.Host != "" {
			if tid, err := s.tenantResolve.ResolveTenantIDByDomain(ctx, hr.Host); err == nil && tid > 0 {
				ctx = kratosMetadata.AppendToClientContext(ctx, "x-md-global-tenant-id", strconv.FormatUint(uint64(tid), 10))
			}
		}
	}
	return s.formSubmissionClient.Submit(ctx, req)
}
