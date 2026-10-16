package service

import (
	"context"
	"fmt"
	"time"
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
	experimentClient     contentV1.ExperimentServiceClient
	tenantResolve        entmiddleware.TenantResolver
	log                  *log.Helper
}

func NewFormSubmissionService(ctx *bootstrap.Context, formSubmissionClient contentV1.FormSubmissionServiceClient, experimentClient contentV1.ExperimentServiceClient, tenantResolve entmiddleware.TenantResolver) *FormSubmissionService {
	return &FormSubmissionService{
		log:                  ctx.NewLoggerHelper("form-submission/service/app-service"),
		formSubmissionClient: formSubmissionClient,
		experimentClient:     experimentClient,
		tenantResolve:        tenantResolve,
	}
}

func (s *FormSubmissionService) Submit(ctx context.Context, req *contentV1.SubmitFormSubmissionRequest) (*contentV1.FormSubmission, error) {
	tenantMD := ""
	if s.tenantResolve != nil {
		if hr, ok := khttp.RequestFromServerContext(ctx); ok && hr != nil && hr.Host != "" {
			if tid, err := s.tenantResolve.ResolveTenantIDByDomain(ctx, hr.Host); err == nil && tid > 0 {
				tenantMD = strconv.FormatUint(uint64(tid), 10)
				ctx = kratosMetadata.AppendToClientContext(ctx, "x-md-global-tenant-id", tenantMD)
			}
		}
	}
	resp, err := s.formSubmissionClient.Submit(ctx, req)
	if err == nil && tenantMD != "" {
		// A/B 实验转化归因：访客（gwc_vid cookie）在页面实验中参与过，
		// 其表单提交计为该变体的转化。后台协程 + 独立超时执行，
		// gRPC 陈旧连接挂起也不拖累提交响应（best-effort）。
		hr, ok := khttp.RequestFromServerContext(ctx)
		if ok && hr != nil {
			if vid, cerr := hr.Cookie(experimentVisitorCookie); cerr == nil && vid.Value != "" {
				vidValue := vid.Value
				pageID := req.GetPageId()
				md := tenantMD
				go s.attributeConversion(md, pageID, vidValue)
			}
		}
	}
	return resp, err
}

// attributeConversion 表单提交 → 实验转化归因（后台协程执行）。
// 页面上运行中的首个实验若已记录该访客的参与，则标记转化。
func (s *FormSubmissionService) attributeConversion(tenantMD string, pageID uint32, vidValue string) {
	if s.experimentClient == nil || pageID == 0 || vidValue == "" {
		return
	}
	// 独立于请求生命周期：Background + 超时 + 租户元数据
	//（派生请求 ctx 会在请求结束时被取消，归因必然失败）
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if tenantMD != "" {
		ctx = kratosMetadata.AppendToClientContext(ctx, "x-md-global-tenant-id", tenantMD)
	}
	exps, err := s.experimentClient.ListRunningByPage(ctx, &contentV1.ListRunningExperimentsByPageRequest{PageId: pageID})
	if err != nil || len(exps.GetItems()) == 0 {
		if err != nil {
			s.log.Warnf("experiment conversion: list running failed (page=%d): %v", pageID, err)
		}
		return
	}
	exp := exps.GetItems()[0]
	visitorHash := fmt.Sprintf("%d:%s", exp.GetId(), vidValue)
	if _, err := s.experimentClient.RecordExperimentConversion(ctx, &contentV1.RecordExperimentConversionRequest{
		ExperimentId: exp.GetId(),
		VisitorHash:  visitorHash,
	}); err != nil {
		s.log.Warnf("experiment conversion attribution failed (experiment=%d): %v", exp.GetId(), err)
	}
}
