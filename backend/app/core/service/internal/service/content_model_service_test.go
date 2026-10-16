package service

import (
	"strings"
	"testing"

	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
)

func TestContentRefPattern(t *testing.T) {
	ok := map[string][2]string{
		"post:12":    {"post", "12"},
		"page:345":   {"page", "345"},
		"  post:7  ": {"post", "7"}, // TrimSpace 后由调用方处理，此处验证模式本体
	}
	for raw, want := range ok {
		m := contentRefPattern.FindStringSubmatch(strings.TrimSpace(raw))
		if m == nil || m[1] != want[0] || m[2] != want[1] {
			t.Fatalf("contentRefPattern(%q) = %v, want %v", raw, m, want)
		}
	}
	for _, bad := range []string{"category:1", "tag:2", "post:", ":12", "post:abc", "widget:12", "post:12x", "media_asset:3"} {
		if contentRefPattern.MatchString(strings.TrimSpace(bad)) {
			t.Fatalf("%q should not match", bad)
		}
	}
}

func TestPickReferenceTranslation(t *testing.T) {
	trs := []*contentV1.PostTranslation{
		{LanguageCode: strPtr("en-US"), Title: strPtr("Hello"), Slug: strPtr("hello")},
		{LanguageCode: strPtr("zh-CN"), Title: strPtr("你好"), Slug: strPtr("ni-hao")},
	}

	if title, slug := pickReferenceTranslation("zh-CN", trs); title != "你好" || slug != "ni-hao" {
		t.Fatalf("lang match = (%q,%q)", title, slug)
	}
	if title, _ := pickReferenceTranslation("ja-JP", trs); title != "Hello" {
		t.Fatalf("fallback = %q, want first non-empty Hello", title)
	}
	// 某语言只有标题没有 slug：标题命中该语言，slug 独立回退
	trs2 := []*contentV1.PageTranslation{
		{LanguageCode: strPtr("en-US"), Title: strPtr("Only"), Slug: strPtr("only")},
		{LanguageCode: strPtr("zh-CN"), Title: strPtr("只有标题")},
	}
	if title, slug := pickReferenceTranslation("zh-CN", trs2); title != "只有标题" || slug != "only" {
		t.Fatalf("partial translation = (%q,%q), want (只有标题,only)", title, slug)
	}
	if title, slug := pickReferenceTranslation("", []*contentV1.PostTranslation(nil)); title != "" || slug != "" {
		t.Fatalf("empty translations = (%q,%q)", title, slug)
	}
}

func strPtr(s string) *string { return &s }
