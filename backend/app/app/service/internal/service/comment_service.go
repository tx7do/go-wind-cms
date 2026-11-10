package service

import (
	"context"
	"strconv"

	authnEngine "github.com/tx7do/kratos-authn/engine"
	"github.com/go-kratos/kratos/v2/log"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	kratosMetadata "github.com/go-kratos/kratos/v2/metadata"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"github.com/tx7do/go-utils/trans"
	"google.golang.org/protobuf/types/known/emptypb"

	appV1 "go-wind-cms/api/gen/go/app/service/v1"
	commentV1 "go-wind-cms/api/gen/go/comment/service/v1"

	"go-wind-cms/pkg/middleware/auth"
	entmiddleware "go-wind-cms/pkg/middleware/ent"
)

type CommentService struct {
	appV1.CommentServiceHTTPServer

	commentClient      commentV1.CommentServiceClient
	tenantResolve      entmiddleware.TenantResolver
	accessTokenChecker auth.AccessTokenChecker
	log                *log.Helper
}

func NewCommentService(
	ctx *bootstrap.Context,
	commentClient commentV1.CommentServiceClient,
	tenantResolve entmiddleware.TenantResolver,
	accessTokenChecker auth.AccessTokenChecker,
) *CommentService {
	return &CommentService{
		log:                ctx.NewLoggerHelper("comment/service/app-service"),
		commentClient:      commentClient,
		tenantResolve:      tenantResolve,
		accessTokenChecker: accessTokenChecker,
	}
}


func (s *CommentService) List(ctx context.Context, req *paginationV1.PagingRequest) (*commentV1.ListCommentResponse, error) {
	resp, err := s.commentClient.List(ctx, req)
	if err != nil {
		return nil, err
	}
	// 公开端点仅返回已批准评论，过滤待审核/拒绝/垃圾等状态
	if resp != nil {
		filtered := make([]*commentV1.Comment, 0, len(resp.GetItems()))
		for _, c := range resp.GetItems() {
			if c != nil && c.GetStatus() == commentV1.Comment_STATUS_APPROVED {
				filtered = append(filtered, c)
			}
		}
		resp.Items = filtered
		resp.Total = uint64(len(filtered))
	}
	return resp, nil
}

func (s *CommentService) Get(ctx context.Context, req *commentV1.GetCommentRequest) (*commentV1.Comment, error) {
	resp, err := s.commentClient.Get(ctx, req)
	if err != nil {
		return nil, err
	}
	// 公开端点仅返回已批准评论
	if resp == nil || resp.GetStatus() != commentV1.Comment_STATUS_APPROVED {
		return nil, commentV1.ErrorNotFound("comment not found")
	}
	return resp, nil
}

// Create 创建评论。
//
// 游客评论策略:本操作在认证白名单中(游客无 bearer 也能到达),身份在此
// 按可选认证解析——携带有效 token 为登录用户,否则为游客,并按站点设置
// (site_settings)的开关决定放行/拒绝:
//   - enable_comments=false          → 评论功能整体关闭
//   - allow_guest_comments=false      → 仅登录用户可评论(游客要求先登录)
// 游客评论必须提供昵称与邮箱,入库 author_type=GUEST、created_by=0;
// 登录用户入库 author_id/author_type=USER、created_by=本人。
func (s *CommentService) Create(ctx context.Context, req *commentV1.CreateCommentRequest) (*commentV1.Comment, error) {
	if req == nil || req.Data == nil {
		return nil, commentV1.ErrorBadRequest("invalid parameter")
	}

	// ── 可选认证:白名单操作不会注入 operator,手动校验 bearer(存在才解析) ──
	var (
		isLogin       bool
		loginUserID   uint64
		loginUsername string
	)
	if s.accessTokenChecker != nil {
		if token, err := authnEngine.AuthFromMD(ctx, authnEngine.BearerWord, authnEngine.ContextTypeKratosMetaData); err == nil {
			valid, payload := s.accessTokenChecker.IsValidAccessToken(ctx, token, false)
			if valid && payload != nil {
				isLogin = true
				loginUserID = payload.GetUserId()
				loginUsername = payload.GetUsername()
			}
		}
	}

	// ── 站点评论策略开关(enable_comments / allow_guest_comments)
	// 由 core 端在入库前校验:设置键直接查库,绕开 BFF→core 列表过滤不可靠的问题。

	// 游客/匿名链路 BFF→core 不携带租户上下文,按 Host 解析租户后经
	// x-md-global-tenant-id 元数据传给 core,供其显式落库归属(内网 gRPC 可信);
	// 解析不出租户则不注入,评论归属平台租户(与匿名内容链路一致)。
	if s.tenantResolve != nil {
		if hr, ok := khttp.RequestFromServerContext(ctx); ok && hr != nil && hr.Host != "" {
			if tid, terr := s.tenantResolve.ResolveTenantIDByDomain(ctx, hr.Host); terr == nil && tid > 0 {
				ctx = kratosMetadata.AppendToClientContext(ctx, "x-md-global-tenant-id", strconv.FormatUint(uint64(tid), 10))
			}
		}
	}

	if isLogin {
		// 登录用户:归属真实身份;昵称未填时回落用户名
		req.Data.AuthorId = trans.Ptr(loginUserID)
		req.Data.AuthorType = trans.Ptr(commentV1.Comment_AUTHOR_TYPE_USER)
		req.Data.CreatedBy = trans.Ptr(loginUserID)
		if req.Data.AuthorName == nil || *req.Data.AuthorName == "" {
			req.Data.AuthorName = trans.Ptr(loginUsername)
		}
	} else {
		// 游客:必须提供昵称与邮箱;身份不与服务账号挂钩
		if req.Data.GetAuthorName() == "" || req.Data.GetAuthorEmail() == "" {
			return nil, commentV1.ErrorBadRequest("author name and email are required for guest comments")
		}
		req.Data.AuthorType = trans.Ptr(commentV1.Comment_AUTHOR_TYPE_GUEST)
		req.Data.AuthorId = trans.Ptr(uint64(0))
		req.Data.CreatedBy = trans.Ptr(uint64(0))
	}

	// 统一待审核入库,审核通过后前台可见
	req.Data.Status = trans.Ptr(commentV1.Comment_STATUS_PENDING)

	return s.commentClient.Create(ctx, req)
}


// ensureCommentOwner 校验调用者是否为目标评论的作者，防止任意登录用户改/删他人评论（IDOR）。
func (s *CommentService) ensureCommentOwner(ctx context.Context, commentID uint64) error {
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return commentV1.ErrorUnauthorized("authentication required")
	}

	comment, err := s.commentClient.Get(ctx, &commentV1.GetCommentRequest{
		Id: commentID,
	})
	if err != nil {
		return err
	}

	if comment.GetCreatedBy() != operator.GetUserId() {
		return commentV1.ErrorForbidden("you can only modify your own comments")
	}
	return nil
}

func (s *CommentService) Update(ctx context.Context, req *commentV1.UpdateCommentRequest) (*commentV1.Comment, error) {
	if err := s.ensureCommentOwner(ctx, req.GetId()); err != nil {
		return nil, err
	}

	// 获取操作人信息，强制以服务端身份覆盖 UpdatedBy，审计归属真实
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	req.Data.UpdatedBy = trans.Ptr(operator.GetUserId())

	return s.commentClient.Update(ctx, req)
}

func (s *CommentService) Delete(ctx context.Context, req *commentV1.DeleteCommentRequest) (*emptypb.Empty, error) {
	if err := s.ensureCommentOwner(ctx, req.GetId()); err != nil {
		return nil, err
	}
	return s.commentClient.Delete(ctx, req)
}
