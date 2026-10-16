package service

import (
	"context"
	"encoding/json"
	"path"
	"regexp"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	contentV1 "go-wind-cms/api/gen/go/content/service/v1"

	"go-wind-cms/app/core/service/internal/data"
)

// ReferenceTrackingService 内容引用追踪：内容保存后从正文/缩略图抽取媒体引用，
// 重建 content_references 引用行并回写 media_assets.reference_count。
// 全链路 best-effort：失败仅记日志，不阻断内容保存（计数由下次保存修正）。
type ReferenceTrackingService struct {
	mediaAssetRepo *data.MediaAssetRepo
	referenceRepo  *data.ContentReferenceRepo
	log            *log.Helper
}

func NewReferenceTrackingService(
	ctx *bootstrap.Context,
	mediaAssetRepo *data.MediaAssetRepo,
	referenceRepo *data.ContentReferenceRepo,
) *ReferenceTrackingService {
	return &ReferenceTrackingService{
		mediaAssetRepo: mediaAssetRepo,
		referenceRepo:  referenceRepo,
		log:            ctx.NewLoggerHelper("reference-tracking/service/core-service"),
	}
}

// HTML 内嵌图片（TipTap 富文本产出的 <img src>；src 可单引号/双引号/无引号）
var reHTMLImg = regexp.MustCompile(`(?i)<img[^>]+src\s*=\s*("([^"]+)"|'([^']+)'|([^\s>]+))`)

// Markdown 图片 ![alt](src)
var reMarkdownImg = regexp.MustCompile(`!\[[^\]]*\]\(([^)\s]+)`)

// RemoveSource 源内容删除时清空其引用边并重算受影响目标计数。
func (s *ReferenceTrackingService) RemoveSource(ctx context.Context, sourceType string, sourceID uint32) {
	s.syncWithTargets(ctx, sourceType, sourceID, nil)
}

// SyncPostReferences 抽取并重建一篇文章的媒体引用。
// translations 取保存后回读的完整翻译集（部分更新也能拿到最终状态）。
func (s *ReferenceTrackingService) SyncPostReferences(ctx context.Context, postID uint32, translations []*contentV1.PostTranslation, thumbnailURL string) {
	s.sync(ctx, "post", postID, thumbnailURL, func() []string {
		bodies := make([]string, 0, len(translations))
		for _, tr := range translations {
			if tr != nil {
				bodies = append(bodies, tr.GetContent())
			}
		}
		return bodies
	})
}

// SyncPageReferences 抽取页面的媒体引用：各语言封面图 coverImage（field=cover）
// 与 sections 内容值中的图片（field=section）。
// section 内容值形状随区块类型而异（直接 URL / HTML / Markdown / JSON 图列表），
// 统一走通用值提取，不感知区块类型。
func (s *ReferenceTrackingService) SyncPageReferences(ctx context.Context, pageID uint32, translations []*contentV1.PageTranslation, sections []*contentV1.Section) {
	covers := make([]string, 0, len(translations))
	for _, tr := range translations {
		if tr != nil && tr.GetCoverImage() != "" {
			covers = append(covers, tr.GetCoverImage())
		}
	}
	sectionURLs := extractSectionContentURLs(sections)
	if len(covers) == 0 && len(sectionURLs) == 0 {
		s.syncWithTargets(ctx, "page", pageID, nil)
		return
	}

	candidates := make([]string, 0, len(covers)+len(sectionURLs))
	candidates = append(candidates, covers...)
	candidates = append(candidates, sectionURLs...)
	paths := make([]string, 0, len(candidates))
	for _, u := range candidates {
		if p := normalizeMediaPath(u); p != "" {
			paths = append(paths, p)
		}
	}
	pathByID, urlByID, err := s.mediaAssetRepo.FindByPathCandidates(ctx, paths, candidates)
	if err != nil {
		s.log.Errorf("resolve page refs for page [%d] failed: %s", pageID, err.Error())
		return
	}
	resolve := func(u string) uint32 {
		if id := urlByID[u]; id != 0 {
			return id
		}
		if p := normalizeMediaPath(u); p != "" {
			return pathByID[p]
		}
		return 0
	}

	seen := make(map[string]struct{}, len(candidates))
	targets := make([]data.RefTarget, 0, len(candidates))
	addRef := func(u, field string) {
		if _, dup := seen[u+field]; dup {
			return
		}
		seen[u+field] = struct{}{}
		if id := resolve(u); id != 0 {
			targets = append(targets, data.RefTarget{TargetType: "media_asset", TargetID: id, Field: field})
		}
	}
	for _, u := range covers {
		addRef(u, "cover")
	}
	for _, u := range sectionURLs {
		addRef(u, "section")
	}
	s.syncWithTargets(ctx, "page", pageID, targets)
}

// extractSectionContentURLs 遍历区块内容键值对收集候选图片 URL。
// 三路收集：值本身形似 URL、值内嵌 HTML/Markdown 图片、值为 JSON 数组/对象时
// 递归取字符串（gallery/carousel 的图列表）。
func extractSectionContentURLs(sections []*contentV1.Section) []string {
	seen := make(map[string]struct{})
	collect := func(str string) {
		for _, u := range extractImageURLs(str) {
			seen[u] = struct{}{}
		}
		if isURLLike(str) {
			seen[str] = struct{}{}
		}
	}
	var walk func(v any)
	walk = func(v any) {
		switch tv := v.(type) {
		case string:
			collect(tv)
			trimmed := strings.TrimSpace(tv)
			if strings.HasPrefix(trimmed, "[") || strings.HasPrefix(trimmed, "{") {
				var decoded any
				if err := json.Unmarshal([]byte(trimmed), &decoded); err == nil {
					walk(decoded)
				}
			}
		case []any:
			for _, item := range tv {
				walk(item)
			}
		case map[string]any:
			for _, item := range tv {
				walk(item)
			}
		}
	}
	for _, sec := range sections {
		if sec == nil {
			continue
		}
		for _, tr := range sec.GetTranslations() {
			if tr == nil {
				continue
			}
			for _, v := range tr.GetContent() {
				walk(v)
			}
		}
	}
	out := make([]string, 0, len(seen))
	for u := range seen {
		out = append(out, u)
	}
	return out
}

// isURLLike 值是否形似媒体 URL（http/https 绝对地址，与上传链路落库的 url 形态一致）。
func isURLLike(v string) bool {
	return strings.HasPrefix(v, "http://") || strings.HasPrefix(v, "https://")
}

// sync 抽取 → URL 解析为资产 ID → 重建引用行 → 回写引用计数。
func (s *ReferenceTrackingService) sync(ctx context.Context, sourceType string, sourceID uint32, thumbnailURL string, bodies func() []string) {
	// 1) 正文内嵌图片（全语言正文合并抽取）；缩略图单独记边（field=thumbnail）
	// 注意：抽取为空也要走到底——重建为空引用集正是"去图后计数归零"的清除语义。
	seenURLs := make(map[string]struct{})
	for _, body := range bodies() {
		for _, rawURL := range extractImageURLs(body) {
			seenURLs[rawURL] = struct{}{}
		}
	}

	candidateURLs := make([]string, 0, len(seenURLs)+1)
	for u := range seenURLs {
		candidateURLs = append(candidateURLs, u)
	}
	thumbnailIsBody := false
	if thumbnailURL != "" {
		if _, inBody := seenURLs[thumbnailURL]; inBody {
			thumbnailIsBody = true
		} else {
			candidateURLs = append(candidateURLs, thumbnailURL)
		}
	}
	if len(candidateURLs) == 0 {
		s.syncWithTargets(ctx, sourceType, sourceID, nil)
		return
	}

	// 2) URL → 资产 ID：storage_path 精确（剥协议/host/bucket 后的路径）与 url 精确双匹配
	paths := make([]string, 0, len(candidateURLs))
	for _, u := range candidateURLs {
		if p := normalizeMediaPath(u); p != "" {
			paths = append(paths, p)
		}
	}
	pathByID, urlByID, err := s.mediaAssetRepo.FindByPathCandidates(ctx, paths, candidateURLs)
	if err != nil {
		s.log.Errorf("resolve media refs for %s [%d] failed: %s", sourceType, sourceID, err.Error())
		return
	}

	resolve := func(u string) uint32 {
		if id := urlByID[u]; id != 0 {
			return id
		}
		if p := normalizeMediaPath(u); p != "" {
			return pathByID[p]
		}
		return 0
	}

	targets := make([]data.RefTarget, 0, len(candidateURLs)+1)
	for u := range seenURLs {
		if id := resolve(u); id != 0 {
			targets = append(targets, data.RefTarget{TargetType: "media_asset", TargetID: id, Field: "body"})
		}
	}
	if thumbnailURL != "" && !thumbnailIsBody {
		if id := resolve(thumbnailURL); id != 0 {
			targets = append(targets, data.RefTarget{TargetType: "media_asset", TargetID: id, Field: "thumbnail"})
		}
	}
	s.syncWithTargets(ctx, sourceType, sourceID, targets)
}

// syncWithTargets 用给定的引用目标集重建引用行并回写目标侧计数。
// targets 为空表示清空该源的引用（源内容删除时）。
func (s *ReferenceTrackingService) syncWithTargets(ctx context.Context, sourceType string, sourceID uint32, targets []data.RefTarget) {
	// 重建引用行（旧目标并入受影响集合，删除的引用计数也会归零）
	affected, err := s.referenceRepo.ReplaceSourceReferences(ctx, sourceType, sourceID, targets)
	if err != nil {
		s.log.Errorf("replace references for %s [%d] failed: %s", sourceType, sourceID, err.Error())
		return
	}
	if len(affected) == 0 {
		return
	}

	// 回写 reference_count
	ids := make([]uint32, 0, len(affected))
	for id := range affected {
		ids = append(ids, id)
	}
	counts, err := s.referenceRepo.CountByTargetIDs(ctx, "media_asset", ids)
	if err != nil {
		s.log.Errorf("count references for %s [%d] failed: %s", sourceType, sourceID, err.Error())
		return
	}
	if err = s.mediaAssetRepo.UpdateReferenceCounts(ctx, counts); err != nil {
		s.log.Errorf("update reference counts for %s [%d] failed: %s", sourceType, sourceID, err.Error())
	}
}

// extractImageURLs 从正文（HTML 与 Markdown 混合容错）抽取全部图片 URL。
func extractImageURLs(body string) []string {
	if body == "" {
		return nil
	}
	urls := make([]string, 0, 8)
	for _, m := range reHTMLImg.FindAllStringSubmatch(body, -1) {
		for _, g := range m[2:] {
			if g != "" {
				urls = append(urls, g)
				break
			}
		}
	}
	for _, m := range reMarkdownImg.FindAllStringSubmatch(body, -1) {
		urls = append(urls, m[1])
	}
	return urls
}

// normalizeMediaPath 从完整 URL 剥出存储路径候选：去 query/fragment、去协议与 host、
// 去 MinIO bucket 前缀段，得到与 media_assets.storage_path 对齐的候选。
// 非 http(s) URL（data:、相对路径）原样返回（去掉查询串）。
func normalizeMediaPath(rawURL string) string {
	u := strings.TrimSpace(rawURL)
	if u == "" {
		return ""
	}
	if i := strings.IndexAny(u, "?#"); i >= 0 {
		u = u[:i]
	}
	if strings.HasPrefix(u, "http://") || strings.HasPrefix(u, "https://") {
		rest := u[len("https://"):]
		if i := strings.Index(rest, "/"); i >= 0 {
			u = rest[i+1:]
		} else {
			return ""
		}
	}
	u = strings.TrimPrefix(u, "/")
	// MinIO 公有直链形如 /{bucket}/{object...}：剥第一段 bucket。
	// storage_path 保存的是纯 objectName（上传链路约定），故剥离后精确匹配。
	if i := strings.Index(u, "/"); i > 0 {
		u = u[i+1:]
	}
	u = path.Clean(u)
	if u == "." || u == "/" {
		return ""
	}
	return u
}
