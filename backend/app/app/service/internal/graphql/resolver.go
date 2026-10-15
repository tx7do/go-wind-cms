// Package graphql C 端只读内容 GraphQL API 的解析器实现。
//
// 数据源复用 app BFF 的 REST 同源服务（internal/service），继承其
// "仅已发布 + 翻译回退"语义；本包不直连 core，也不自行鉴权。
// 契约：backend/api/graphql/schema.graphql；生成产物：api/gen/go/graphql。
package graphql

// THIS CODE WILL BE UPDATED WITH SCHEMA CHANGES. PREVIOUS IMPLEMENTATION FOR SCHEMA CHANGES WILL BE KEPT IN THE COMMENT SECTION. IMPLEMENTATION FOR UNCHANGED SCHEMA WILL BE KEPT.

import (
	"context"
	"math"
	"sort"
	"strconv"
	"time"

	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/go-utils/trans"
	"google.golang.org/protobuf/types/known/timestamppb"

	appV1 "go-wind-cms/api/gen/go/app/service/v1"
	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
	graphqlV1 "go-wind-cms/api/gen/go/graphql"
	"go-wind-cms/app/app/service/internal/service"
)

// Resolver GraphQL 根解析器：持有 app BFF 各内容服务。
type Resolver struct {
	postSvc     *service.PostService
	pageSvc     *service.PageService
	categorySvc *service.CategoryService
	tagSvc      *service.TagService
}

// NewResolver 构造根解析器（由 graphql 服务器装配注入）。
func NewResolver(
	postService *service.PostService,
	pageService *service.PageService,
	categoryService *service.CategoryService,
	tagService *service.TagService,
) *Resolver {
	return &Resolver{
		postSvc:     postService,
		pageSvc:     pageService,
		categorySvc: categoryService,
		tagSvc:      tagService,
	}
}

// Post is the resolver for the post field.
func (r *queryResolver) Post(ctx context.Context, id *string, slug *string, locale *string) (*graphqlV1.Post, error) {
	hasSlug := slug != nil && *slug != ""
	hasID := id != nil && *id != ""
	if !hasSlug && !hasID {
		return nil, nil
	}

	loc := optStrVal(locale)
	var (
		p   *contentV1.Post
		err error
	)
	if hasSlug {
		p, err = r.postSvc.GetBySlug(ctx, &appV1.GetPostBySlugRequest{Slug: *slug, Locale: loc})
	} else {
		var pid uint32
		if pid, err = parseID(*id); err == nil {
			p, err = r.postSvc.Get(ctx, &contentV1.GetPostRequest{
				QueryBy: &contentV1.GetPostRequest_Id{Id: pid},
			})
		}
	}
	if err != nil {
		return nil, err
	}
	return mapPost(p, loc), nil
}

// Posts is the resolver for the posts field.
func (r *queryResolver) Posts(ctx context.Context, page *int, pageSize *int, locale *string) (*graphqlV1.PostList, error) {
	resp, err := r.postSvc.List(ctx, pagingRequest(page, pageSize))
	if err != nil {
		return nil, err
	}
	loc := optStrVal(locale)
	items := make([]*graphqlV1.Post, 0, len(resp.GetItems()))
	for _, p := range resp.GetItems() {
		if g := mapPost(p, loc); g != nil {
			items = append(items, g)
		}
	}
	return &graphqlV1.PostList{
		Items:    items,
		Total:    clampInt(resp.GetTotal()),
		Page:     normPage(page),
		PageSize: normPageSize(pageSize),
	}, nil
}

// Page is the resolver for the page field.
func (r *queryResolver) Page(ctx context.Context, id *string, slug *string, locale *string) (*graphqlV1.Page, error) {
	hasSlug := slug != nil && *slug != ""
	hasID := id != nil && *id != ""
	if !hasSlug && !hasID {
		return nil, nil
	}

	loc := optStrVal(locale)
	var (
		pg  *contentV1.Page
		err error
	)
	if hasSlug {
		pg, err = r.pageSvc.GetBySlug(ctx, &appV1.GetPageBySlugRequest{Slug: *slug, Locale: loc})
	} else {
		var pid uint32
		if pid, err = parseID(*id); err == nil {
			pg, err = r.pageSvc.Get(ctx, &contentV1.GetPageRequest{
				QueryBy: &contentV1.GetPageRequest_Id{Id: pid},
			})
		}
	}
	if err != nil {
		return nil, err
	}
	return mapPage(pg, loc), nil
}

// Pages is the resolver for the pages field.
func (r *queryResolver) Pages(ctx context.Context, page *int, pageSize *int, locale *string) (*graphqlV1.PageList, error) {
	resp, err := r.pageSvc.List(ctx, pagingRequest(page, pageSize))
	if err != nil {
		return nil, err
	}
	loc := optStrVal(locale)
	items := make([]*graphqlV1.Page, 0, len(resp.GetItems()))
	for _, pg := range resp.GetItems() {
		if g := mapPage(pg, loc); g != nil {
			items = append(items, g)
		}
	}
	return &graphqlV1.PageList{
		Items:    items,
		Total:    clampInt(resp.GetTotal()),
		Page:     normPage(page),
		PageSize: normPageSize(pageSize),
	}, nil
}

// Categories is the resolver for the categories field.
func (r *queryResolver) Categories(ctx context.Context, page *int, pageSize *int, locale *string) (*graphqlV1.CategoryList, error) {
	resp, err := r.categorySvc.List(ctx, pagingRequest(page, pageSize))
	if err != nil {
		return nil, err
	}
	loc := optStrVal(locale)
	items := make([]*graphqlV1.Category, 0, len(resp.GetItems()))
	for _, c := range resp.GetItems() {
		if g := mapCategory(c, loc); g != nil {
			items = append(items, g)
		}
	}
	return &graphqlV1.CategoryList{
		Items:    items,
		Total:    clampInt(resp.GetTotal()),
		Page:     normPage(page),
		PageSize: normPageSize(pageSize),
	}, nil
}

// Tags is the resolver for the tags field.
func (r *queryResolver) Tags(ctx context.Context, page *int, pageSize *int, locale *string) (*graphqlV1.TagList, error) {
	resp, err := r.tagSvc.List(ctx, pagingRequest(page, pageSize))
	if err != nil {
		return nil, err
	}
	loc := optStrVal(locale)
	items := make([]*graphqlV1.Tag, 0, len(resp.GetItems()))
	for _, tg := range resp.GetItems() {
		if g := mapTag(tg, loc); g != nil {
			items = append(items, g)
		}
	}
	return &graphqlV1.TagList{
		Items:    items,
		Total:    clampInt(resp.GetTotal()),
		Page:     normPage(page),
		PageSize: normPageSize(pageSize),
	}, nil
}

// Query returns graphqlV1.QueryResolver implementation.
func (r *Resolver) Query() graphqlV1.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }

// ── 内部映射工具 ──

// pickTranslation 翻译挑选：精确匹配 locale，缺省回退第一条
//（与前端 server-api pickTranslation 的回退顺序一致）。
func pickTranslation[T interface{ GetLanguageCode() string }](translations []T, locale string) T {
	var zero T
	if len(translations) == 0 {
		return zero
	}
	if locale != "" {
		for i := range translations {
			if translations[i].GetLanguageCode() == locale {
				return translations[i]
			}
		}
	}
	return translations[0]
}

func mapPost(p *contentV1.Post, locale string) *graphqlV1.Post {
	if p == nil {
		return nil
	}
	tr := pickTranslation(p.GetTranslations(), locale)
	out := &graphqlV1.Post{
		ID:                 formatID(p.GetId()),
		Status:             p.GetStatus().String(),
		Thumbnail:          optStr(p.GetThumbnail()),
		AuthorName:         optStr(p.GetAuthorName()),
		PublishTime:        optStr(rfc3339(p.GetPublishTime())),
		DisallowComment:    optBool(p.DisallowComment),
		IsFeatured:         optBool(p.IsFeatured),
		AvailableLanguages: p.GetAvailableLanguages(),
		Slug:               optStr(p.GetCode()),
	}
	if tr != nil {
		out.Title = optStr(tr.GetTitle())
		out.Summary = optStr(tr.GetSummary())
		out.Content = optStr(tr.GetContent())
		out.Slug = firstNonEmptyStr(tr.GetSlug(), p.GetCode())
	}
	return out
}

func mapPage(pg *contentV1.Page, locale string) *graphqlV1.Page {
	if pg == nil {
		return nil
	}
	tr := pickTranslation(pg.GetTranslations(), locale)
	out := &graphqlV1.Page{
		ID:                 formatID(pg.GetId()),
		Status:             pg.GetStatus().String(),
		Slug:               optStr(pg.GetSlug()),
		Sections:           mapSections(pg.GetSections(), locale),
		AvailableLanguages: pg.GetAvailableLanguages(),
	}
	if tr != nil {
		out.Title = optStr(tr.GetTitle())
		out.CoverImage = optStr(tr.GetCoverImage())
		out.Slug = firstNonEmptyStr(tr.GetSlug(), pg.GetSlug())
	}
	return out
}

func mapSections(sections []*contentV1.Section, locale string) []*graphqlV1.Section {
	out := make([]*graphqlV1.Section, 0, len(sections))
	for _, s := range sections {
		if s == nil {
			continue
		}
		item := &graphqlV1.Section{
			ID:        formatID(s.GetId()),
			Type:      s.GetType().String(),
			Name:      optStr(s.GetName()),
			SortOrder: int(s.GetSortOrder()),
			Content:   []*graphqlV1.KeyValue{},
		}
		if str := pickTranslation(s.GetTranslations(), locale); str != nil {
			for k, v := range str.GetContent() {
				item.Content = append(item.Content, &graphqlV1.KeyValue{Key: k, Value: v})
			}
			sort.Slice(item.Content, func(i, j int) bool {
				return item.Content[i].Key < item.Content[j].Key
			})
		}
		out = append(out, item)
	}
	return out
}

func mapCategory(c *contentV1.Category, locale string) *graphqlV1.Category {
	if c == nil {
		return nil
	}
	out := &graphqlV1.Category{
		ID: formatID(c.GetId()),
	}
	if tr := pickTranslation(c.GetTranslations(), locale); tr != nil {
		out.Name = optStr(tr.GetName())
		out.Slug = optStr(tr.GetSlug())
		out.Description = optStr(tr.GetDescription())
		out.CoverImage = firstNonEmptyStr(tr.GetCoverImage(), c.GetThumbnail())
	}
	return out
}

func mapTag(tg *contentV1.Tag, locale string) *graphqlV1.Tag {
	if tg == nil {
		return nil
	}
	out := &graphqlV1.Tag{
		ID: formatID(tg.GetId()),
	}
	if tr := pickTranslation(tg.GetTranslations(), locale); tr != nil {
		out.Name = optStr(tr.GetName())
		out.Slug = optStr(tr.GetSlug())
	}
	return out
}

// pagingRequest 分页参数：缺省 1/20，pageSize 钳制 1..100。
func pagingRequest(page, pageSize *int) *paginationV1.PagingRequest {
	return &paginationV1.PagingRequest{
		Page:     trans.Ptr(uint32(normPage(page))),
		PageSize: trans.Ptr(uint32(normPageSize(pageSize))),
	}
}

func normPage(page *int) int {
	if page != nil && *page > 0 {
		return *page
	}
	return 1
}

func normPageSize(pageSize *int) int {
	if pageSize != nil && *pageSize > 0 {
		if *pageSize > 100 {
			return 100
		}
		return *pageSize
	}
	return 20
}

func parseID(raw string) (uint32, error) {
	v, err := strconv.ParseUint(raw, 10, 32)
	if err != nil {
		return 0, contentV1.ErrorBadRequest("invalid id: %s", raw)
	}
	return uint32(v), nil
}

func formatID(id uint32) string {
	return strconv.FormatUint(uint64(id), 10)
}

func clampInt(v uint64) int {
	if v > math.MaxInt32 {
		return math.MaxInt32
	}
	return int(v)
}

func optStr(s string) *string {
	if s == "" {
		return nil
	}
	return trans.Ptr(s)
}

func optStrVal(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func firstNonEmptyStr(vals ...string) *string {
	for _, v := range vals {
		if v != "" {
			return trans.Ptr(v)
		}
	}
	return nil
}

func optBool(b *bool) *bool {
	if b == nil {
		return nil
	}
	return trans.Ptr(*b)
}

func rfc3339(ts *timestamppb.Timestamp) string {
	if ts == nil || !ts.IsValid() {
		return ""
	}
	return ts.AsTime().UTC().Format(time.RFC3339)
}
