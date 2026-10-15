package service

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
)

// SeoService sitemap/robots 生成（面向搜索引擎的部署域工具）。
//
// URL 规则由环境变量配置（部署域与 C 端实现相关，后端不做硬编码假设）：
//
//	SITE_BASE_URL   站点对外基址（必配才生成 sitemap；未配置返回空集 sitemap）
//	SEO_POST_PATH   帖子 URL 模板，默认 "/{locale}/post/detail?id={id}"（react 静态导出姿势）
//	SEO_PAGE_PATH   页面 URL 模板，默认 "/{locale}/page/detail?slug={slug}"
//
// 占位符：{locale} {id} {slug} {updatedAt}(W3C date)。
// 每语言一条 <url>（取帖子/页面各语言的翻译行；无翻译回退默认语言）。
// 条数上限 1000（V1 单页拉取；超大型站点应改分页迭代或预生成）。
type SeoService struct {
	postClient contentV1.PostServiceClient
	pageClient contentV1.PageServiceClient
	log        *log.Helper
}

func NewSeoService(ctx *bootstrap.Context, postClient contentV1.PostServiceClient, pageClient contentV1.PageServiceClient) *SeoService {
	return &SeoService{
		log:        ctx.NewLoggerHelper("seo/service/app-service"),
		postClient: postClient,
		pageClient: pageClient,
	}
}

type sitemapURL struct {
	Loc        string
	LastMod    string // W3C date (2006-01-02)
}

// envPath 取路径模板并把反斜杠归一为斜杠（Windows 环境变量手滑容错）。
func envPath(key, def string) string {
	v := strings.ReplaceAll(os.Getenv(key), "\\", "/")
	if v == "" {
		return def
	}
	return v
}

func (s *SeoService) baseURL() string {
	return strings.TrimRight(os.Getenv("SITE_BASE_URL"), "/")
}

// Sitemap 生成 sitemap XML。站点基址未配置时返回带注释的空 sitemap（合法 XML，不 500）。
func (s *SeoService) Sitemap(ctx context.Context) (string, error) {
	base := s.baseURL()
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	b.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">` + "\n")
	if base == "" {
		b.WriteString("<!-- SITE_BASE_URL not configured -->\n")
		b.WriteString("</urlset>\n")
		return b.String(), nil
	}

	// 首页
	b.WriteString("  <url><loc>" + xmlEscape(base + "/") + "</loc></url>\n")

	postPath := envPath("SEO_POST_PATH", "/{locale}/post/detail?id={id}")
	pagePath := envPath("SEO_PAGE_PATH", "/{locale}/page/detail?slug={slug}")

	// 帖子：已发布（app List 已过滤），带翻译行生成多语言 URL
	posts, err := s.postClient.List(ctx, &paginationV1.PagingRequest{
		Page:     trans.Ptr(uint32(1)),
		PageSize: trans.Ptr(uint32(1000)),
		FieldMask: pathMask("id", "updated_at", "translations.language_code", "translations.slug"),
	})
	if err != nil {
		s.log.Errorf("sitemap: list posts failed: %v", err)
	} else {
		for _, p := range posts.GetItems() {
			lastMod := w3cDate(p.GetUpdatedAt())
			for _, tr := range p.GetTranslations() {
				loc := base + renderPath(postPath, map[string]string{
					"locale":    tr.GetLanguageCode(),
					"id":        fmt.Sprintf("%d", p.GetId()),
					"slug":      tr.GetSlug(),
					"updatedAt": lastMod,
				})
				b.WriteString("  <url><loc>" + xmlEscape(loc) + "</loc>" + lastModTag(lastMod) + "</url>\n")
			}
		}
	}

	// 页面：已发布（app List 已过滤）
	pages, err := s.pageClient.List(ctx, &paginationV1.PagingRequest{
		Page:     trans.Ptr(uint32(1)),
		PageSize: trans.Ptr(uint32(1000)),
		FieldMask: pathMask("id", "updated_at", "slug", "translations.language_code", "translations.slug"),
	})
	if err != nil {
		s.log.Errorf("sitemap: list pages failed: %v", err)
	} else {
		for _, pg := range pages.GetItems() {
			lastMod := w3cDate(pg.GetUpdatedAt())
			for _, tr := range pg.GetTranslations() {
				loc := base + renderPath(pagePath, map[string]string{
					"locale":    tr.GetLanguageCode(),
					"id":        fmt.Sprintf("%d", pg.GetId()),
					"slug":      firstNonEmpty(tr.GetSlug(), pg.GetSlug()),
					"updatedAt": lastMod,
				})
				b.WriteString("  <url><loc>" + xmlEscape(loc) + "</loc>" + lastModTag(lastMod) + "</url>\n")
			}
		}
	}

	b.WriteString("</urlset>\n")
	return b.String(), nil
}

// Robots 生成 robots.txt：全放行 + sitemap 地址（基址未配置时省略 Sitemap 行）。
func (s *SeoService) Robots(ctx context.Context) (string, error) {
	var b strings.Builder
	b.WriteString("User-agent: *\nAllow: /\n")
	if base := s.baseURL(); base != "" {
		b.WriteString("Sitemap: " + base + "/sitemap.xml\n")
	}
	return b.String(), nil
}

// ─── 内部 ───

func w3cDate(unix *timestamppb.Timestamp) string {
	if unix == nil || unix.AsTime().IsZero() {
		return ""
	}
	return unix.AsTime().Format("2006-01-02")
}

func lastModTag(lastMod string) string {
	if lastMod == "" {
		return ""
	}
	return "<lastmod>" + lastMod + "</lastmod>"
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

func renderPath(tpl string, vars map[string]string) string {
	out := tpl
	for k, v := range vars {
		out = strings.ReplaceAll(out, "{"+k+"}", v)
	}
	return out
}

func xmlEscape(s string) string {
	r := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", "'", "&#39;", `"`, "&quot;")
	return r.Replace(s)
}

// pathMask 构造 view_mask 的 FieldMask（仅取 sitemap 需要的最小字段集）。
func pathMask(paths ...string) *fieldmaskpb.FieldMask {
	return &fieldmaskpb.FieldMask{Paths: paths}
}
