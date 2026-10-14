package social

import (
	"html"
	"strings"
)

// PlainTextToHTML 把纯文本转成基础 HTML 段落（富文本平台的正文用）：
// 先 HTML 转义，再按空行分段落 <p>，段内单换行转 <br/>。
// 适用于微信公众号草稿正文、WordPress 文章内容等以 HTML 提交的平台。
func PlainTextToHTML(text string) string {
	text = strings.TrimSpace(text)
	if text == "" {
		return ""
	}

	blocks := strings.Split(text, "\n\n")
	var b strings.Builder
	for _, block := range blocks {
		block = strings.TrimSpace(block)
		if block == "" {
			continue
		}
		b.WriteString("<p>")
		b.WriteString(html.EscapeString(block))
		b.WriteString("</p>")
	}

	out := b.String()
	// 段内的单换行（\n）此时位于转义后的文本里，转成 <br/>
	return strings.ReplaceAll(out, "\n", "<br/>")
}
