package mdprotect

import (
	"strings"
	"testing"
)

// benchmarkSource 构造一份含各类可保护片段的 Markdown（约 4KB）。
func benchmarkSource() string {
	var sb strings.Builder
	for i := 0; i < 12; i++ {
		sb.WriteString("## 章节\n\n这是一段包含 `inline_code_var` 的中文正文，附一个 [链接](https://example.com/docs/page?q=1)。\n\n")
		sb.WriteString("```go\nfunc main() { fmt.Println(\"hello 中文注释\") }\n```\n\n")
		sb.WriteString("图片示例：![封面](https://cdn.example.com/img/cover.png) 与裸地址 https://example.com/wiki/术语 。\n\n")
	}
	return sb.String()
}

// BenchmarkProtect 遮蔽阶段：翻译前的正则扫描与占位符替换。
func BenchmarkProtect(b *testing.B) {
	src := benchmarkSource()
	b.ReportAllocs()
	b.SetBytes(int64(len(src)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Protect(src)
	}
}

// BenchmarkRestoreRoundTrip 遮蔽→（模拟翻译）→恢复的完整链路。
func BenchmarkRestoreRoundTrip(b *testing.B) {
	src := benchmarkSource()
	masked, restore := Protect(src)
	// 模拟翻译引擎：仅改写非占位符文本。
	translated := strings.ReplaceAll(masked, "中文", "Chinese")
	b.ReportAllocs()
	b.SetBytes(int64(len(src)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = restore(translated)
	}
}

// BenchmarkProtectPlainText 无可保护片段的纯文本（恒等变换路径）。
func BenchmarkProtectPlainText(b *testing.B) {
	src := strings.Repeat("这是一段没有任何可保护片段的纯中文正文。", 50)
	b.ReportAllocs()
	b.SetBytes(int64(len(src)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Protect(src)
	}
}

// TestBenchmarkSourcesCovered 保证 benchmark 数据确实命中各保护规则。
func TestBenchmarkSourcesCovered(t *testing.T) {
	src := benchmarkSource()
	masked, restore := Protect(src)
	restored := restore(masked)
	if restored != src {
		t.Fatal("round trip must be identity on untranslated input")
	}
	if strings.Contains(masked, "```") || strings.Contains(masked, "https://") {
		t.Fatal("code fences and urls must be masked")
	}
}
