package data

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	entCrud "github.com/tx7do/go-crud/entgo"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/searchsynonym"
)

// SearchSynonymRepo 搜索同义词仓储（sys_search_synonyms；租户级词库）。
type SearchSynonymRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper
}

func NewSearchSynonymRepo(ctx *bootstrap.Context, entClient *entCrud.EntClient[*ent.Client]) *SearchSynonymRepo {
	return &SearchSynonymRepo{
		log:       ctx.NewLoggerHelper("search-synonym/repo/core-service"),
		entClient: entClient,
	}
}

// synonymsFromJSON 解析扩展词 JSON 数组。
func synonymsFromJSON(raw string) []string {
	var arr []string
	if err := json.Unmarshal([]byte(raw), &arr); err != nil {
		return nil
	}
	out := make([]string, 0, len(arr))
	for _, v := range arr {
		if t := strings.TrimSpace(v); t != "" {
			out = append(out, t)
		}
	}
	return out
}

// toProto ent 行 → proto。
func (r *SearchSynonymRepo) toProto(e *ent.SearchSynonym) *contentV1.SearchSynonym {
	return &contentV1.SearchSynonym{
		Id:       trans.Ptr(e.ID),
		Term:     e.Term,
		Synonyms: synonymsFromJSON(e.Synonyms),
		Enabled:  trans.Ptr(e.Enabled == nil || *e.Enabled),
	}
}

// List 全量（本租户，按 term 排序；量小不分页）。
func (r *SearchSynonymRepo) List(ctx context.Context) ([]*contentV1.SearchSynonym, error) {
	rows, err := r.entClient.Client().SearchSynonym.Query().
		Order(ent.Asc(searchsynonym.FieldTerm)).
		All(ctx)
	if err != nil {
		r.log.Errorf("list search synonyms failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("list search synonyms failed")
	}
	out := make([]*contentV1.SearchSynonym, 0, len(rows))
	for i := range rows {
		out = append(out, r.toProto(rows[i]))
	}
	return out, nil
}

// Create 创建（term 重复由唯一索引兜底）。
func (r *SearchSynonymRepo) Create(ctx context.Context, tenantID uint32, dto *contentV1.SearchSynonym) (*contentV1.SearchSynonym, error) {
	raw, err := json.Marshal(dto.GetSynonyms())
	if err != nil {
		return nil, contentV1.ErrorBadRequest("invalid synonyms payload")
	}
	// 租户归属：优先显式传入值（服务层从 viewer 取），缺失时由调用方兜底
	now := time.Now()
	builder := r.entClient.Client().SearchSynonym.Create().
		SetTerm(dto.GetTerm()).
		SetSynonyms(string(raw)).
		SetNillableEnabled(dto.Enabled).
		SetNillableCreatedAt(&now).
		SetNillableUpdatedAt(&now)
	if tenantID > 0 {
		builder.SetTenantID(tenantID)
	}
	row, err := builder.Save(ctx)
	if err != nil {
		r.log.Errorf("create search synonym failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("create search synonym failed")
	}
	return r.toProto(row), nil
}

// Delete 删除。
func (r *SearchSynonymRepo) Delete(ctx context.Context, id uint32) error {
	err := r.entClient.Client().SearchSynonym.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return contentV1.ErrorNotFound("search synonym not found")
		}
		r.log.Errorf("delete search synonym failed: %s", err.Error())
		return contentV1.ErrorInternalServerError("delete search synonym failed")
	}
	return nil
}

// EnabledTermMap 本租户启用词表：触发词 → 扩展词列表（检索查询侧扩展用）。
// 词表量小（几十条），每次搜索直查可接受。
func (r *SearchSynonymRepo) EnabledTermMap(ctx context.Context, tenantID uint32) map[string][]string {
	rows, err := r.entClient.Client().SearchSynonym.Query().
		Where(
			searchsynonym.TenantIDEQ(tenantID),
			searchsynonym.EnabledEQ(true),
		).
		All(ctx)
	if err != nil {
		r.log.Errorf("load enabled synonyms failed: %s", err.Error())
		return nil
	}
	out := make(map[string][]string, len(rows))
	for i := range rows {
		term := strings.TrimSpace(rows[i].Term)
		if term == "" {
			continue
		}
		out[strings.ToLower(term)] = synonymsFromJSON(rows[i].Synonyms)
	}
	return out
}

// ExpandQuery 查询侧同义词扩展：把查询串按词表展开成候选查询组
// （原始查询永远在首位）。命中不到词表时返回 nil（无需扩展）。
func (r *SearchSynonymRepo) ExpandQuery(ctx context.Context, tenantID uint32, query string) []string {
	if r.entClient == nil {
		return nil
	}
	termMap := r.EnabledTermMap(ctx, tenantID)
	// 平台级（tenant 0）词库全局生效，与租户词库合并——匿名搜索链路的租户
	// 由 Host 解析，平台管理员维护的词库应对所有租户站点可用。
	for term, syns := range r.EnabledTermMap(ctx, 0) {
		if _, dup := termMap[term]; !dup {
			termMap[term] = syns
		}
	}
	if len(termMap) == 0 {
		return nil
	}
	q := strings.ToLower(strings.TrimSpace(query))
	if q == "" {
		return nil
	}
	// 触发词可能是多字词：先整串匹配，再按空白拆词逐个匹配
	candidates := []string{q}
	candidates = append(candidates, strings.Fields(q)...)
	var expanded []string
	seen := map[string]bool{q: true}
	for _, cand := range candidates {
		for _, syn := range termMap[cand] {
			if !seen[syn] {
				seen[syn] = true
				expanded = append(expanded, syn)
			}
		}
	}
	return expanded
}
