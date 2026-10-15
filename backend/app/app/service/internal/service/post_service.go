package service

import (
	"context"
	"strconv"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	kratosMetadata "github.com/go-kratos/kratos/v2/metadata"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/tx7do/go-utils/trans"

	appV1 "go-wind-cms/api/gen/go/app/service/v1"
	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
	entmiddleware "go-wind-cms/pkg/middleware/ent"
	"go-wind-cms/pkg/preview"
)

type PostService struct {
	appV1.PostServiceHTTPServer

	postClient    contentV1.PostServiceClient
	tenantResolve entmiddleware.TenantResolver
	log           *log.Helper
}

func NewPostService(ctx *bootstrap.Context, postClient contentV1.PostServiceClient, tenantResolve entmiddleware.TenantResolver) *PostService {
	return &PostService{
		log:           ctx.NewLoggerHelper("post/service/app-service"),
		postClient:    postClient,
		tenantResolve: tenantResolve,
	}
}

func (s *PostService) List(ctx context.Context, req *paginationV1.PagingRequest) (*contentV1.ListPostResponse, error) {
	resp, err := s.postClient.List(ctx, req)
	if err != nil {
		return nil, err
	}
	// 公开端点仅返回已发布文章，过滤草稿/归档/私有等状态
	if resp != nil {
		filtered := make([]*contentV1.Post, 0, len(resp.GetItems()))
		for _, p := range resp.GetItems() {
			if p != nil && p.GetStatus() == contentV1.Post_POST_STATUS_PUBLISHED {
				filtered = append(filtered, p)
			}
		}
		resp.Items = filtered
		resp.Total = uint64(len(filtered))
	}
	return resp, nil
}

func (s *PostService) Get(ctx context.Context, req *contentV1.GetPostRequest) (*contentV1.Post, error) {
	resp, err := s.postClient.Get(ctx, req)
	if err != nil {
		return nil, err
	}
	// 公开端点仅返回已发布文章，草稿/归档/私有等状态按未找到处理；
	// 持有效预览令牌（admin 签发的短时效 HMAC）时放行草稿供后台预览
	if resp == nil || resp.GetStatus() != contentV1.Post_POST_STATUS_PUBLISHED {
		if !s.hasValidPreviewToken(ctx, resp.GetId()) {
			return nil, contentV1.ErrorNotFound("post not found")
		}
	}
	return resp, nil
}

// hasValidPreviewToken 校验请求 query 里的 preview_token：令牌有效且指向该文章
// 时放行；令牌存在但无效按"不存在"处理，不暴露令牌有效性细节。
// GetBySlug 按 slug 获取已发布帖子（SEO 路由 /post/{slug} 的数据源）。
func (s *PostService) GetBySlug(ctx context.Context, req *appV1.GetPostBySlugRequest) (*contentV1.Post, error) {
	if req == nil || req.GetSlug() == "" {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}
	// 帖子的 SEO 别名走 code 字段（唯一编码，demo 数据即 slug 形态）；
	// core GetPostRequest 的 query_by 不含 slug oneof，Code 语义等同。
	coreReq := &contentV1.GetPostRequest{
		QueryBy: &contentV1.GetPostRequest_Code{Code: req.GetSlug()},
	}
	if req.GetLocale() != "" {
		coreReq.Locale = trans.Ptr(req.GetLocale())
	}
	resp, err := s.postClient.Get(ctx, coreReq)
	if err != nil {
		return nil, err
	}
	// 公开 SEO 路由仅返回已发布文章
	if resp == nil || resp.GetStatus() != contentV1.Post_POST_STATUS_PUBLISHED {
		return nil, contentV1.ErrorNotFound("post not found")
	}
	return resp, nil
}
func (s *PostService) hasValidPreviewToken(ctx context.Context, targetID uint32) bool {
	token := rawQueryValue(ctx, "preview_token")
	if token == "" {
		return false
	}
	id, err := preview.Verify(token, preview.KindPost, time.Now())
	if err != nil {
		s.log.Warnf("preview token rejected (kind=post): %v", err)
		return false
	}
	return id != 0 && id == targetID
}

// Create/Update/Delete 在 app（公开站点）服务上禁用：CMS 内容的写操作应经由 admin 服务，
// 公开站点登录用户不应直接创建/修改/删除文章。RBAC 为故意的 noop，故在此显式拒绝。
func (s *PostService) Create(_ context.Context, _ *contentV1.CreatePostRequest) (*contentV1.Post, error) {
	return nil, contentV1.ErrorForbidden("content mutation is not allowed on the public app service")
}

func (s *PostService) Update(_ context.Context, _ *contentV1.UpdatePostRequest) (*contentV1.Post, error) {
	return nil, contentV1.ErrorForbidden("content mutation is not allowed on the public app service")
}

func (s *PostService) Delete(_ context.Context, _ *contentV1.DeletePostRequest) (*emptypb.Empty, error) {
	return nil, contentV1.ErrorForbidden("content mutation is not allowed on the public app service")
}

func (s *PostService) GetTranslation(ctx context.Context, req *contentV1.GetPostRequest) (*contentV1.PostTranslation, error) {
	return s.postClient.GetTranslation(ctx, req)
}

// SearchPosts 全文搜索帖子。
//
// core 端搜索对租户做强制过滤并硬编码 status=PUBLISHED，仅返回 postId/language/
// title 最小字段集。tenant_id 由服务端上下文决定，调用方无法指定或绕过。
// 匿名经路线2 解析 Host 得到 AnonymousTenantViewer，但 BFF→core 是 gRPC 调用、
// 拿不到 HTTP Host，故此处按请求 Host 解析租户后经 x-md-global-tenant-id 传给
// core 兜底（仅本服务内网 gRPC 可达，header 不可被外部调用方伪造注入）；
// 解析不出租户则不注入，core 端维持 fail-closed（返回空）。
func (s *PostService) SearchPosts(ctx context.Context, req *contentV1.SearchPostsRequest) (*contentV1.SearchPostsResponse, error) {
	if s.tenantResolve != nil {
		if hr, ok := khttp.RequestFromServerContext(ctx); ok && hr != nil && hr.Host != "" {
			if tid, err := s.tenantResolve.ResolveTenantIDByDomain(ctx, hr.Host); err == nil && tid > 0 {
				ctx = kratosMetadata.AppendToClientContext(ctx, "x-md-global-tenant-id", strconv.FormatUint(uint64(tid), 10))
			}
		}
	}
	return s.postClient.SearchPosts(ctx, req)
}
