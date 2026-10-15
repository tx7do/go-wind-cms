package service

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"net/http"
	"time"

	khttp "github.com/go-kratos/kratos/v2/transport/http"
	kratosTransport "github.com/go-kratos/kratos/v2/transport"
	"github.com/google/uuid"
	"google.golang.org/protobuf/encoding/protojson"

	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
)

// 访客标识 cookie（一年期、匿名 UUID；实验内再与实验 ID 派生稳定哈希，
// 不落任何原始身份）。
const experimentVisitorCookie = "gwc_vid"

// audienceCondition 受众条件（变体 conditions JSON 数组元素），全部满足才命中。
type audienceCondition struct {
	Field string `json:"field"` // locale / device / logged_in / param
	Name  string `json:"name"`  // field=param 时的查询参数名
	Op    string `json:"op"`    // eq（默认）/ ne / contains / exists
	Value string `json:"value"`
}

// applyExperiment 给已发布页面做运行中实验的变体分配（A/B 分流 + 受众定向）。
//
// 分配规则：实验内以 (experiment_id, visitor_id) 稳定哈希分桶（0-99）——
// assignment_type=random 直接按变体权重走累积区间；rule 型先依变体顺序求值
// 受众条件首中定向，未中回退权重分流（权重和不足 100 的部分回落默认渲染）。
// 命中后以变体 sections 快照（protojson）整页替换渲染内容，并幂等记录参与。
// 全程 best-effort：任何失败只告警并回落原页面，不影响正常取页。
func (s *PageService) applyExperiment(ctx context.Context, page *contentV1.Page, locale string) {
	if s.experimentClient == nil || page == nil || page.GetId() == 0 {
		return
	}
	exps, err := s.experimentClient.ListRunningByPage(ctx, &contentV1.ListRunningExperimentsByPageRequest{PageId: page.GetId()})
	if err != nil {
		s.log.Warnf("list running experiments failed (page=%d): %v", page.GetId(), err)
		return
	}
	if len(exps.GetItems()) == 0 {
		return
	}
	exp := exps.GetItems()[0]
	if len(exp.GetVariants()) == 0 {
		return
	}

	vid, isNew, err := s.visitorID(ctx)
	if err != nil {
		s.log.Warnf("resolve visitor id failed: %v", err)
		return
	}
	if isNew {
		setVisitorCookie(ctx, vid)
	}

	hash := fmt.Sprintf("%d:%s", exp.GetId(), vid)
	bucket := visitorBucket(hash)
	variant := assignVariant(exp, bucket, ctx, locale)
	if variant == nil || variant.GetSections() == "" {
		return
	}

	// 变体 sections 快照存 {"sections":[...]} 包装形态（protojson 需消息根）
	var wrapper contentV1.Page
	if err := (protojson.UnmarshalOptions{DiscardUnknown: true}).Unmarshal([]byte(variant.GetSections()), &wrapper); err != nil {
		s.log.Warnf("unmarshal variant sections failed (experiment=%d variant=%s): %v", exp.GetId(), variant.GetName(), err)
		return
	}
	page.Sections = wrapper.GetSections()

	if _, err := s.experimentClient.RecordParticipant(ctx, &contentV1.RecordExperimentParticipantRequest{
		ExperimentId: exp.GetId(),
		VisitorHash:  hash,
		VariantName:  variant.GetName(),
	}); err != nil {
		s.log.Warnf("record experiment participant failed (experiment=%d): %v", exp.GetId(), err)
	}
}

// visitorID 读（或种）访客标识 cookie。isNew 表示本次响应需回种。
func (s *PageService) visitorID(ctx context.Context) (string, bool, error) {
	if hr, ok := khttp.RequestFromServerContext(ctx); ok && hr != nil {
		if c, cErr := hr.Cookie(experimentVisitorCookie); cErr == nil && c.Value != "" {
			return c.Value, false, nil
		}
	}
	id, err := uuid.NewRandom()
	if err != nil {
		return "", false, fmt.Errorf("generate visitor uuid: %w", err)
	}
	return id.String(), true, nil
}

// setVisitorCookie 经 kratos reply header 回种访客 cookie（一年期）。
func setVisitorCookie(ctx context.Context, vid string) {
	sc, ok := kratosTransport.FromServerContext(ctx)
	if !ok {
		return
	}
	cookie := (&http.Cookie{
		Name:     experimentVisitorCookie,
		Value:    vid,
		Path:     "/",
		MaxAge:   int((365 * 24 * time.Hour).Seconds()),
		HttpOnly: false,
		SameSite: http.SameSiteLaxMode,
	}).String()
	sc.ReplyHeader().Set("Set-Cookie", cookie)
}

// visitorBucket 访客派生串 → 0-99 稳定桶位（同一访客在同一实验内恒定）。
func visitorBucket(visitorHash string) int {
	h := fnv.New32a()
	_, _ = h.Write([]byte(visitorHash))
	return int(h.Sum32() % 100)
}

// assignVariant 按实验分配方式选变体；nil 表示不命中（默认渲染）。
func assignVariant(exp *contentV1.Experiment, bucket int, ctx context.Context, locale string) *contentV1.ExperimentVariant {
	variants := exp.GetVariants()

	if exp.GetAssignmentType() == contentV1.Experiment_ASSIGNMENT_TYPE_RULE {
		for _, v := range variants {
			if v.GetConditions() != "" && matchConditions(v.GetConditions(), ctx, locale) {
				return v
			}
		}
		// 未中回退权重分流
	}

	cumulative := 0
	for _, v := range variants {
		cumulative += int(v.GetWeight())
		if bucket < cumulative {
			return v
		}
	}
	return nil
}

// matchConditions 受众条件求值：JSON 数组内全部条件满足才命中。
func matchConditions(conditionsJSON string, ctx context.Context, locale string) bool {
	var conditions []audienceCondition
	if err := json.Unmarshal([]byte(conditionsJSON), &conditions); err != nil {
		return false
	}
	if len(conditions) == 0 {
		return false
	}
	for _, c := range conditions {
		if !matchCondition(c, ctx, locale) {
			return false
		}
	}
	return true
}

func matchCondition(c audienceCondition, ctx context.Context, locale string) bool {
	value, present := conditionFieldValue(c, ctx, locale)
	switch c.Op {
	case "exists":
		return present
	case "ne":
		return present && value != c.Value
	case "contains":
		return present && (containsFold(value, c.Value))
	case "", "eq":
		return present && value == c.Value
	}
	return false
}

func conditionFieldValue(c audienceCondition, ctx context.Context, locale string) (string, bool) {
	switch c.Field {
	case "locale":
		return locale, locale != ""
	case "device":
		if hr, ok := khttp.RequestFromServerContext(ctx); ok && hr != nil {
			ua := hr.Header.Get("User-Agent")
			if isMobileUA(ua) {
				return "mobile", true
			}
			return "desktop", true
		}
		return "", false
	case "logged_in":
		if hr, ok := khttp.RequestFromServerContext(ctx); ok && hr != nil {
			if hr.Header.Get("Authorization") != "" {
				return "true", true
			}
			return "false", true
		}
		return "", false
	case "param":
		if c.Name == "" {
			return "", false
		}
		if hr, ok := khttp.RequestFromServerContext(ctx); ok && hr != nil {
			v := hr.URL.Query().Get(c.Name)
			return v, v != ""
		}
		return "", false
	}
	return "", false
}

// isMobileUA 粗粒度移动端识别（UA 关键词，够受众定向用）。
func isMobileUA(ua string) bool {
	for _, kw := range []string{"Mobile", "Android", "iPhone", "iPad", "HarmonyOS"} {
		if containsFold(ua, kw) {
			return true
		}
	}
	return false
}

func containsFold(s, sub string) bool {
	if len(sub) == 0 {
		return true
	}
	sLower := []rune(s)
	subLower := []rune(sub)
	toLower := func(r rune) rune {
		if r >= 'A' && r <= 'Z' {
			return r + ('a' - 'A')
		}
		return r
	}
	for i := range sLower {
		if i+len(subLower) > len(sLower) {
			break
		}
		match := true
		for j := range subLower {
			if toLower(sLower[i+j]) != toLower(subLower[j]) {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}
