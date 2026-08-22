package mdprotect

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 模拟翻译引擎：把正文文本"翻译"为固定串，占位符原样保留
func fakeTranslate(masked string) string {
	var sb strings.Builder
	for _, line := range strings.Split(masked, "\n") {
		if line == "" {
			continue
		}
		sb.WriteString("[译]")
		sb.WriteString(line)
		sb.WriteString("\n")
	}
	return sb.String()
}

func TestProtectPlainTextInputChanged(t *testing.T) {
	masked, restore := Protect("Hello World")
	assert.Equal(t, "Hello World", masked)
	assert.Equal(t, "Hello World", restore(masked))
}

func TestProtectFencedCodeBlock(t *testing.T) {
	src := "前文\n```go\nfmt.Println(\"hello\")\n```\n后文"

	masked, restore := Protect(src)
	assert.NotContains(t, masked, "fmt.Println", "代码块应被遮蔽")
	assert.Contains(t, masked, "[[[0]]]")

	got := restore(fakeTranslate(masked))
	assert.Contains(t, got, "```go\nfmt.Println(\"hello\")\n```", "代码块应原样恢复")
}

func TestProtectInlineCode(t *testing.T) {
	src := "使用 `kubectl get pods` 命令查看"

	masked, restore := Protect(src)
	assert.NotContains(t, masked, "kubectl", "行内代码应被遮蔽")

	got := restore(fakeTranslate(masked))
	assert.Contains(t, got, "`kubectl get pods`")
}

func TestProtectLinkURLButNotText(t *testing.T) {
	src := "参见 [官方文档](https://example.com/docs) 了解详情"

	masked, restore := Protect(src)
	assert.NotContains(t, masked, "https://example.com", "URL 应被遮蔽")
	assert.Contains(t, masked, "官方文档", "链接文字应保留可翻译")
	assert.Contains(t, masked, "]([[[0]]])", "应保留 ]( ) 结构只遮 URL")

	got := restore(fakeTranslate(masked))
	assert.Contains(t, got, "https://example.com/docs", "URL 应原样恢复")
}

func TestProtectImageURL(t *testing.T) {
	src := "![截图](https://example.com/img.png)"

	masked, restore := Protect(src)
	assert.NotContains(t, masked, "example.com")

	got := restore(fakeTranslate(masked))
	assert.Contains(t, got, "![截图]")
	assert.Contains(t, got, "(https://example.com/img.png)")
}

func TestProtectBareURL(t *testing.T) {
	src := "详见 https://example.com/page 说明"

	masked, restore := Protect(src)
	assert.NotContains(t, masked, "https://example.com")

	got := restore(fakeTranslate(masked))
	assert.Contains(t, got, "https://example.com/page")
}

func TestRestoreToleratesSpacesInToken(t *testing.T) {
	masked, restore := Protect("`code`")
	// 引擎在占位符内插入空白时应仍能恢复
	tampered := strings.ReplaceAll(masked, "[[[0]]]", "[[[ 0 ]]]")
	assert.Equal(t, "`code`", restore(tampered))
}

func TestRestoreLeavesUnknownTokenIntact(t *testing.T) {
	_, restore := Protect("`code`")
	// 引擎改写出未知 token 时原样保留，不影响其余内容
	assert.Equal(t, "译 [[[9]]] 文", restore("译 [[[9]]] 文"))
}

func TestProtectMultipleSegments(t *testing.T) {
	src := "先 `a` 再 ```py\nx=1\n``` 然后 [l](https://x.io) 最后 https://y.io"

	masked, restore := Protect(src)
	assert.Equal(t, 4, strings.Count(masked, "[[["), "四类片段各一个占位符")

	got := restore(masked)
	for _, want := range []string{"`a`", "```py\nx=1\n```", "https://x.io", "https://y.io"} {
		assert.Contains(t, got, want)
	}
}

func TestFencedBlockBeforeInlineCode(t *testing.T) {
	// 围栏块内的反引号不应被行内代码规则拆开处理
	src := "```\ninline `tick` inside\n```\n后文 `out` 外"

	masked, restore := Protect(src)
	got := restore(masked)
	assert.Contains(t, got, "```\ninline `tick` inside\n```")
	assert.Contains(t, got, "`out`")
}

func TestNoPanicOnWeirdInput(t *testing.T) {
	cases := []string{"", "```", "``", "`", "[[", "]]]", "((( ", "a`b`c`d"}
	for _, c := range cases {
		masked, restore := Protect(c)
		require.NotPanics(t, func() {
			_ = restore(masked)
		})
	}
}
