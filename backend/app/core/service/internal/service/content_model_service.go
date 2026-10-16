package service

import (
	"context"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	"go-wind-cms/app/core/service/internal/data"

	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
)

// contentRefPattern RELATION 字段引用值的形态："post:12"/"page:7"。
// category/tag 引用 v1 不解析（无公开详情路由语义）。
var contentRefPattern = regexp.MustCompile(`^(post|page):(\d+)$`)

// maxResolveRefs 单次解析引用数上限（RELATION 字段值个数有限，防批量拖库）。
const maxResolveRefs = 20

// ContentModelService 内容模型服务：模型 CRUD + 字段定义查询 + 字段值校验。
type ContentModelService struct {
	contentV1.UnimplementedContentModelServiceServer

	contentModelRepo *data.ContentModelRepo
	postRepo         *data.PostRepo
	pageRepo         *data.PageRepo
	log              *log.Helper
}

func NewContentModelService(ctx *bootstrap.Context, contentModelRepo *data.ContentModelRepo, postRepo *data.PostRepo, pageRepo *data.PageRepo) *ContentModelService {
	return &ContentModelService{
		log:              ctx.NewLoggerHelper("content-model/service/core-service"),
		contentModelRepo: contentModelRepo,
		postRepo:         postRepo,
		pageRepo:         pageRepo,
	}
}

// GetModelByPost 按文章取其分类链绑定的内容模型（C 端渲染 custom_fields 用）。
// 文章→分类→模型的解析在服务端完成，调用方无法指定模型或借道越权；
// 文章不存在/未绑定模型返回 NotFound/空模型（由门面决定公开语义）。
func (s *ContentModelService) GetModelByPost(ctx context.Context, req *contentV1.GetModelByPostRequest) (*contentV1.ContentModel, error) {
	if req == nil || req.GetPostId() == 0 {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}
	post, err := s.postRepo.Get(ctx, &contentV1.GetPostRequest{
		QueryBy: &contentV1.GetPostRequest_Id{Id: req.GetPostId()},
	})
	if err != nil {
		return nil, err
	}
	modelID := s.contentModelRepo.ResolveModelIDByCategories(ctx, post.GetCategoryIds())
	if modelID == 0 {
		return &contentV1.ContentModel{}, nil
	}
	return s.contentModelRepo.Get(ctx, &contentV1.GetContentModelRequest{
		QueryBy: &contentV1.GetContentModelRequest_Id{Id: modelID},
	})
}

func (s *ContentModelService) List(ctx context.Context, req *paginationV1.PagingRequest) (*contentV1.ListContentModelResponse, error) {
	return s.contentModelRepo.List(ctx, req)
}

func (s *ContentModelService) Count(ctx context.Context, req *paginationV1.PagingRequest) (*contentV1.CountContentModelResponse, error) {
	count, err := s.contentModelRepo.Count(ctx, req)
	if err != nil {
		return nil, err
	}
	return &contentV1.CountContentModelResponse{Count: uint64(count)}, nil
}

func (s *ContentModelService) Get(ctx context.Context, req *contentV1.GetContentModelRequest) (*contentV1.ContentModel, error) {
	return s.contentModelRepo.Get(ctx, req)
}

func (s *ContentModelService) Create(ctx context.Context, req *contentV1.CreateContentModelRequest) (*contentV1.ContentModel, error) {
	return s.contentModelRepo.Create(ctx, req)
}

func (s *ContentModelService) Update(ctx context.Context, req *contentV1.UpdateContentModelRequest) (*contentV1.ContentModel, error) {
	return s.contentModelRepo.Update(ctx, req)
}

func (s *ContentModelService) Delete(ctx context.Context, req *contentV1.DeleteContentModelRequest) (*emptypb.Empty, error) {
	if err := s.contentModelRepo.Delete(ctx, req); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *ContentModelService) ListFieldDefinitions(ctx context.Context, req *contentV1.ListFieldDefinitionsRequest) (*contentV1.ListFieldDefinitionsResponse, error) {
	return s.contentModelRepo.ListFieldDefinitions(ctx, req)
}

// ResolveContentReferences 批量解析 RELATION 字段引用值为目标内容摘要。
// 公开语义：post/page 仅解析 PUBLISHED（草稿/定时/软删/不存在/非法形态
// 静默跳过，调用方据缺失项降级渲染）；租户隔离由 viewer 贯穿仓储查询。
// allow_cross_tenant 配置 v1 不放行跨租户（公开面宁缺毋滥）。
func (s *ContentModelService) ResolveContentReferences(ctx context.Context, req *contentV1.ResolveContentReferencesRequest) (*contentV1.ResolveContentReferencesResponse, error) {
	if req == nil {
		return &contentV1.ResolveContentReferencesResponse{}, nil
	}

	type refWant struct {
		entityType string
		id         uint32
	}
	seen := make(map[string]struct{}, len(req.GetRefs()))
	wants := make([]refWant, 0, len(req.GetRefs()))
	for _, raw := range req.GetRefs() {
		if len(wants) >= maxResolveRefs {
			break
		}
		m := contentRefPattern.FindStringSubmatch(strings.TrimSpace(raw))
		if m == nil {
			continue
		}
		if _, dup := seen[m[0]]; dup {
			continue
		}
		seen[m[0]] = struct{}{}
		id, _ := strconv.ParseUint(m[2], 10, 32)
		wants = append(wants, refWant{entityType: m[1], id: uint32(id)})
	}

	lang := req.GetLanguageCode()
	items := make([]*contentV1.ContentReferenceSummary, 0, len(wants))
	for _, w := range wants {
		switch w.entityType {
		case "post":
			p, err := s.postRepo.Get(ctx, &contentV1.GetPostRequest{
				QueryBy: &contentV1.GetPostRequest_Id{Id: w.id},
			})
			if err != nil || p.GetStatus() != contentV1.Post_POST_STATUS_PUBLISHED {
				continue
			}
			title, slug := pickReferenceTranslation(lang, p.GetTranslations())
			if title == "" && slug == "" {
				continue
			}
			items = append(items, &contentV1.ContentReferenceSummary{
				Ref:              trans.Ptr(w.entityType + ":" + strconv.FormatUint(uint64(w.id), 10)),
				TargetEntityType: trans.Ptr(w.entityType),
				TargetId:         trans.Ptr(w.id),
				Title:            trans.Ptr(title),
				Slug:             trans.Ptr(slug),
				UrlPath:          trans.Ptr("/post/" + slug),
			})
		case "page":
			pg, err := s.pageRepo.Get(ctx, &contentV1.GetPageRequest{
				QueryBy: &contentV1.GetPageRequest_Id{Id: w.id},
			})
			if err != nil || pg.GetStatus() != contentV1.Page_PAGE_STATUS_PUBLISHED {
				continue
			}
			title, slug := pickReferenceTranslation(lang, pg.GetTranslations())
			if slug == "" {
				slug = pg.GetSlug()
			}
			if title == "" && slug == "" {
				continue
			}
			items = append(items, &contentV1.ContentReferenceSummary{
				Ref:              trans.Ptr(w.entityType + ":" + strconv.FormatUint(uint64(w.id), 10)),
				TargetEntityType: trans.Ptr(w.entityType),
				TargetId:         trans.Ptr(w.id),
				Title:            trans.Ptr(title),
				Slug:             trans.Ptr(slug),
				UrlPath:          trans.Ptr("/page/" + slug),
			})
		}
	}
	return &contentV1.ResolveContentReferencesResponse{Items: items}, nil
}

// translationView 解析标题/slug 所需的翻译字段视图（泛型适配 Post/Page 翻译）。
type translationView interface {
	GetLanguageCode() string
	GetTitle() string
	GetSlug() string
}

// pickReferenceTranslation 取语言偏好匹配的标题/slug；无匹配回退首个非空。
// slug 与标题独立回退（某语言只有标题没有 slug 时不丢链接）。
func pickReferenceTranslation[T translationView](lang string, trs []T) (title, slug string) {
	for _, tr := range trs {
		if title == "" && tr.GetTitle() != "" {
			title = tr.GetTitle()
		}
		if slug == "" && tr.GetSlug() != "" {
			slug = tr.GetSlug()
		}
		if lang != "" && tr.GetLanguageCode() == lang && tr.GetTitle() != "" {
			title = tr.GetTitle()
			if tr.GetSlug() != "" {
				slug = tr.GetSlug()
			}
			break
		}
	}
	return title, slug
}
