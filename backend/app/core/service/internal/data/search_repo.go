package data

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"strconv"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	elasticsearchCrud "github.com/tx7do/go-crud/elasticsearch"
)

// ============================================================================
// OpenSearch 搜索 Repo —— posts 索引的 ES 操作封装
//
// 安全模型（与 ent 层 TenantPrivacy 的关键差异）：
//
//   1. ES 无自动租户策略兜底。所有隔离由本 Repo 在查询 DSL 中手动注入
//      bool.filter 的 term{tenant_id} 实现，且不可被调用方绕过。
//   2. 搜索路径不接受 SystemViewer bypass：tenantID==0 直接返回空结果，
//      与 ent 层对 SystemViewer「跳过过滤」的语义有意不同。
//      原因：ent 层有 TenantPrivacy 自动兜底，platform context 跳过是安全的；
//      ES 层无兜底，一旦跳过即全租户泄漏，故必须强制。
//   3. ES 文档的 tenant_id 永远取自 DB 记录（reindex 路径），非 viewer。
//      reindex handler 用 SystemViewer 跨租户读 DB（特权），但写入 ES 的
//      tenant_id 取自 DB 记录字段值，保证数据归属准确。
//   4. 删除按 post_id 单条件 delete-by-query。post_id 是 ent 自增主键、全局唯一
//      （跨租户不重复），单条件即可精确定位、不存在跨租户碰撞风险。这覆盖了
//      软删/硬删/状态变更三种场景——尤其软删时 ent 查询返回 NotFound、无法取
//      tenant_id，单条件删除是唯一可行路径。tenant_id 过滤保留在搜索路径
//      （那里总有有效的 viewer tid）。
//   5. 搜索结果只回传 post_id / language / title，并用 WithSource 限制 ES
//      只返回这三字段——content / tenant_id / status 不回传调用方。
// ============================================================================

const (
	searchIndexName     = "posts"
	searchTemplateName  = "posts_template"
	searchTemplatePrio  = 100
	maxSearchPageSize   = 50
	maxSearchResultFrom = 10000
)

// PostDocument 是写入 ES posts 索引的文档结构。
// tenant_id 永远由 reindex 路径从 DB 记录填入，不接受调用方覆盖。
type PostDocument struct {
	TenantID string `json:"tenant_id"` // 必填，keyword，搜索强制 term 过滤
	PostID   string `json:"post_id"`   // 必填，keyword，删除按此 + tenant_id 定位
	Language string `json:"language"`  // 必填，keyword，搜索强制 term 过滤
	Status   string `json:"status"`    // 必填，keyword，搜索强制 term 过滤（前台仅 published）
	Title    string `json:"title"`     // text + smartcn，全文检索字段
	Summary  string `json:"summary"`   // text + smartcn，全文检索字段
	Content  string `json:"content"`   // text + smartcn，全文检索字段
}

// PostSearchHit 是搜索结果中的单条命中，只暴露最小字段集。
type PostSearchHit struct {
	PostID   string
	Language string
	Title    string
	Score    float64
}

// PostSearchResult 是搜索返回。
type PostSearchResult struct {
	Total int
	Hits  []PostSearchHit
}

type SearchRepo struct {
	esClient *elasticsearchCrud.Client
	log      *log.Helper
}

func NewSearchRepo(ctx *bootstrap.Context, esClient *elasticsearchCrud.Client) *SearchRepo {
	return &SearchRepo{
		log:      ctx.NewLoggerHelper("search/repo/core-service"),
		esClient: esClient,
	}
}

// EnsureIndexTemplate 幂等创建 posts 索引模板。
// 模板绑定 smartcn 分词器到 title/summary/content，并定义 keyword 字段。
// 在 posts 索引首次自动创建时，模板 mapping 会被应用。
func (r *SearchRepo) EnsureIndexTemplate(ctx context.Context) error {
	if r.esClient == nil {
		return errors.New("elasticsearch client is nil")
	}

	properties := map[string]any{
		"tenant_id": map[string]any{"type": "keyword"},
		"post_id":   map[string]any{"type": "keyword"},
		"language":  map[string]any{"type": "keyword"},
		"status":    map[string]any{"type": "keyword"},
		"title":     map[string]any{"type": "text", "analyzer": "smartcn"},
		"summary":   map[string]any{"type": "text", "analyzer": "smartcn"},
		"content":   map[string]any{"type": "text", "analyzer": "smartcn"},
	}

	templateBody := map[string]any{
		"index_patterns": []string{searchIndexName},
		"priority":       searchTemplatePrio,
		"template": map[string]any{
			"mappings": map[string]any{
				// AUD9-L3: 严格化 mapping，禁止动态字段。
				// 仅显式声明的 properties 可被索引；任何未声明字段被静默忽略，
				// 防止脏数据或注入的额外字段进入索引。
				"dynamic":    false,
				"properties": properties,
			},
		},
	}

	bodyBytes, err := json.Marshal(templateBody)
	if err != nil {
		r.log.Errorf("marshal index template body failed: %v", err)
		return err
	}

	if err := r.esClient.CreateIndexTemplate(ctx, searchTemplateName, string(bodyBytes)); err != nil {
		r.log.Errorf("create index template failed: %v", err)
		return err
	}

	r.log.Infof("index template %s ensured for index %s", searchTemplateName, searchIndexName)
	return nil
}

// IndexPost 将一篇帖子的某个语言翻译 upsert 到 ES。
// 文档 id = {post_id}_{language}，同一帖子的每种语言各一个 ES 文档。
// doc.TenantID 必须取自 DB 记录，调用方不可覆盖。
func (r *SearchRepo) IndexPost(ctx context.Context, doc *PostDocument) error {
	if r.esClient == nil {
		return errors.New("elasticsearch client is nil")
	}
	if doc == nil {
		return errors.New("nil post document")
	}
	if doc.TenantID == "" || doc.PostID == "" || doc.Language == "" {
		return errors.New("post document missing mandatory field (tenant_id/post_id/language)")
	}

	docID := doc.PostID + "_" + doc.Language
	if err := r.esClient.InsertDocument(ctx, searchIndexName, docID, doc); err != nil {
		r.log.Errorf("index post document failed (post_id=%s lang=%s): %v", doc.PostID, doc.Language, err)
		return err
	}
	return nil
}

// DeletePost 删除指定帖子在 ES 中的所有语言文档。
// 按 post_id 单条件 delete-by-query。post_id 是 ent 自增主键、全局唯一，
// 单条件即可精确定位，覆盖软删/硬删/状态变更三种场景。
func (r *SearchRepo) DeletePost(ctx context.Context, postID uint32) error {
	if r.esClient == nil {
		return errors.New("elasticsearch client is nil")
	}
	if postID == 0 {
		return errors.New("delete post requires non-zero postID")
	}

	pidStr := strconv.FormatUint(uint64(postID), 10)

	// delete-by-query body：bool.filter 单条件（post_id 全局唯一）
	queryBody := map[string]any{
		"query": map[string]any{
			"bool": map[string]any{
				"filter": []any{
					map[string]any{"term": map[string]any{"post_id": pidStr}},
				},
			},
		},
	}

	bodyBytes, err := json.Marshal(queryBody)
	if err != nil {
		r.log.Errorf("marshal delete-by-query body failed: %v", err)
		return err
	}

	resp, err := r.esClient.Client.DeleteByQuery(
		[]string{searchIndexName},
		bytes.NewReader(bodyBytes),
		r.esClient.Client.DeleteByQuery.WithContext(ctx),
	)
	if err != nil {
		r.log.Errorf("delete-by-query failed (post_id=%s): %v", pidStr, err)
		return err
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			r.log.Warnf("close delete-by-query response body failed: %v", closeErr)
		}
	}()

	if resp.IsError() {
		bodyBytes, _ := io.ReadAll(resp.Body)
		r.log.Errorf("delete-by-query error [%d] (post_id=%s): %s",
			resp.StatusCode, pidStr, string(bodyBytes))
		return errors.New("delete-by-query failed")
	}

	r.log.Infof("deleted ES documents for post_id=%s", pidStr)
	return nil
}

// SearchPosts 执行前台全文搜索。
//
// 安全保证：
//   - tenantID==0 → 返回空（不接受 SystemViewer bypass）
//   - language / status 空 → 返回空
//   - 查询 DSL 必带 bool.filter 的 term{tenant_id} + term{language} + term{status}
//     调用方无法覆盖这三个过滤条件
//   - bool.must 的 multi_match 仅作用于 title/summary/content
//   - WithSource 限制 ES 只回传 post_id / language / title
func (r *SearchRepo) SearchPosts(
	ctx context.Context,
	query string,
	tenantID uint32,
	language string,
	status string,
	page int,
	pageSize int,
) (*PostSearchResult, error) {
	result := &PostSearchResult{}

	if r.esClient == nil {
		return result, errors.New("elasticsearch client is nil")
	}

	// 强制不可绕过的租户/语言/状态过滤
	if tenantID == 0 {
		return result, nil
	}
	if language == "" || status == "" {
		return result, nil
	}
	if query == "" {
		return result, nil
	}

	// 分页边界
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > maxSearchPageSize {
		pageSize = maxSearchPageSize
	}
	if page < 0 {
		page = 0
	}
	from := page * pageSize
	if from > maxSearchResultFrom {
		from = maxSearchResultFrom
	}

	tidStr := strconv.FormatUint(uint64(tenantID), 10)

	// 构建 DSL：filter（访问控制，不参与评分）+ must（相关性评分）
	dsl := map[string]any{
		"from": from,
		"size": pageSize,
		"query": map[string]any{
			"bool": map[string]any{
				"filter": []any{
					map[string]any{"term": map[string]any{"tenant_id": tidStr}},
					map[string]any{"term": map[string]any{"language": language}},
					map[string]any{"term": map[string]any{"status": status}},
				},
				"must": []any{
					map[string]any{
						"multi_match": map[string]any{
							"query":  query,
							"fields": []string{"title", "summary", "content"},
						},
					},
				},
			},
		},
	}

	bodyBytes, err := json.Marshal(dsl)
	if err != nil {
		r.log.Errorf("marshal search DSL failed: %v", err)
		return result, err
	}

	// 调用 raw ES client（绕过 go-crud Search 的 Lucene query string 封装，
	// 因为后者不支持 multi_match 且注入风险高）
	resp, err := r.esClient.Client.Search(
		r.esClient.Client.Search.WithContext(ctx),
		r.esClient.Client.Search.WithIndex(searchIndexName),
		r.esClient.Client.Search.WithBody(bytes.NewReader(bodyBytes)),
		// 仅回传最小字段集，content/tenant_id/status 不返回
		r.esClient.Client.Search.WithSource("post_id", "language", "title"),
	)
	if err != nil {
		r.log.Errorf("search failed: %v", err)
		return result, err
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			r.log.Warnf("close search response body failed: %v", closeErr)
		}
	}()

	if resp.IsError() {
		errBody, _ := io.ReadAll(resp.Body)
		r.log.Errorf("search error [%d]: %s", resp.StatusCode, string(errBody))
		return result, errors.New("search request failed")
	}

	var searchResult elasticsearchCrud.SearchResult
	if err := json.NewDecoder(resp.Body).Decode(&searchResult); err != nil {
		r.log.Errorf("decode search result failed: %v", err)
		return result, err
	}

	hits := searchResult.Hits.Hits
	result.Total = searchResult.Hits.Total.Value
	result.Hits = make([]PostSearchHit, 0, len(hits))

	for _, hit := range hits {
		// hit.Source 是 json.RawMessage，仅含 post_id/language/title（因 WithSource 过滤）
		var src struct {
			PostID   string `json:"post_id"`
			Language string `json:"language"`
			Title    string `json:"title"`
		}
		if err := json.Unmarshal(hit.Source, &src); err != nil {
			r.log.Warnf("unmarshal search hit source failed: %v", err)
			continue
		}
		result.Hits = append(result.Hits, PostSearchHit{
			PostID:   src.PostID,
			Language: src.Language,
			Title:    src.Title,
			Score:    hit.Score,
		})
	}

	return result, nil
}
