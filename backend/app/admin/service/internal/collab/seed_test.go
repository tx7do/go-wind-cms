package collab

import (
	"strings"
	"testing"

	"github.com/reearth/ygo/crdt"
)

// seedDocFromHTML 的表驱动单测：HTML（受控子集）→ yjs 文档结构。
// 断言走 ToXML（节点名为 ProseMirror 节点名：paragraph/heading/bulletList/...，
// TipTap Collaboration 消费同一 fragment，结构即契约）；
// 文本 marks 走 YText delta 断言（marks 存 delta 属性，不进 XML 树）。

func seedToXML(t *testing.T, htmlContent string) string {
	t.Helper()
	doc := crdt.New()
	seedDocFromHTML(doc, htmlContent)
	return doc.GetXmlFragment("default").ToXML()
}

func seedDoc(t *testing.T, htmlContent string) *crdt.Doc {
	t.Helper()
	doc := crdt.New()
	seedDocFromHTML(doc, htmlContent)
	return doc
}

func fragText(t *testing.T, doc *crdt.Doc) string {
	t.Helper()
	frag := doc.GetXmlFragment("default")
	var b strings.Builder
	for _, child := range frag.Children() {
		if txt, ok := child.(*crdt.YXmlText); ok {
			b.WriteString(txt.ToXML())
		} else if elem, ok := child.(*crdt.YXmlElement); ok {
			b.WriteString(elem.ToXML())
		}
	}
	return b.String()
}

func TestSeedDocFromHTML_HeadingAndParagraph(t *testing.T) {
	got := seedToXML(t, "<h2>种子标题</h2><p>正文尾</p>")
	want := "<heading level=\"2\">种子标题</heading><paragraph>正文尾</paragraph>"
	if got != want {
		t.Fatalf("seed mismatch:\n got: %s\nwant: %s", got, want)
	}
}

func TestSeedDocFromHTML_BoldMarkViaDelta(t *testing.T) {
	doc := seedDoc(t, "<p>正文<b>加粗</b>尾</p>")
	frag := doc.GetXmlFragment("default")
	children := frag.Children()
	if len(children) != 1 {
		t.Fatalf("expected 1 paragraph, got %d", len(children))
	}
	p, ok := children[0].(*crdt.YXmlElement)
	if !ok {
		t.Fatalf("expected paragraph element, got %T", children[0])
	}
	// 段落内容为多个 YText（每 run 一个，独立事务插入），marks 在 delta 属性里
	var delta []crdt.Delta
	for _, c := range p.Children() {
		if tt, ok := c.(*crdt.YXmlText); ok {
			delta = append(delta, tt.ToDelta()...)
		}
	}
	if len(delta) != 3 {
		t.Fatalf("expected 3 delta runs, got %d: %+v", len(delta), delta)
	}
	if delta[0].Insert != "正文" || len(delta[0].Attributes) != 0 {
		t.Fatalf("run0 unexpected: %+v", delta[0])
	}
	if delta[1].Insert != "加粗" {
		t.Fatalf("run1 unexpected: %+v", delta[1])
	}
	if _, ok := delta[1].Attributes["bold"]; !ok {
		t.Fatalf("run1 missing bold mark: %+v", delta[1].Attributes)
	}
	if delta[2].Insert != "尾" {
		t.Fatalf("run2 unexpected: %+v", delta[2])
	}
	if _, ok := delta[2].Attributes["bold"]; ok {
		t.Fatalf("run2 must be unmarked: %+v", delta[2].Attributes)
	}
}

func TestSeedDocFromHTML_LinkMark(t *testing.T) {
	doc := seedDoc(t, `<p><a href="https://example.com">链接</a></p>`)
	frag := doc.GetXmlFragment("default")
	children := frag.Children()
	if len(children) != 1 {
		t.Fatalf("expected 1 paragraph")
	}
	p := children[0].(*crdt.YXmlElement)
	var txt *crdt.YXmlText
	for _, c := range p.Children() {
		if tt, ok := c.(*crdt.YXmlText); ok {
			txt = tt
			break
		}
	}
	if txt == nil {
		t.Fatalf("no YXmlText")
	}
	delta := txt.ToDelta()
	if len(delta) != 1 || delta[0].Insert != "链接" {
		t.Fatalf("unexpected delta: %+v", delta)
	}
	link, ok := delta[0].Attributes["link"].(map[string]any)
	if !ok || link["href"] != "https://example.com" {
		t.Fatalf("missing link attrs: %+v", delta[0].Attributes)
	}
}

func TestSeedDocFromHTML_InlineElements(t *testing.T) {
	got := seedToXML(t, `<p>前<br><img src="https://x/a.png" alt="图">后</p>`)
	// 属性写入顺序随 map 遍历不定，按包含断言
	for _, part := range []string{
		"<paragraph>前<hardBreak></hardBreak><image",
		`src="https://x/a.png"`,
		`alt="图"`,
		"后</paragraph>",
	} {
		if !strings.Contains(got, part) {
			t.Fatalf("missing %q in %s", part, got)
		}
	}
}

func TestSeedDocFromHTML_AllHeadingLevels(t *testing.T) {
	got := seedToXML(t, "<h1>a</h1><h3>b</h3><h6>c</h6>")
	for _, level := range []string{"1", "3", "6"} {
		if !strings.Contains(got, `level="`+level+`"`) {
			t.Fatalf("missing heading level %s in %s", level, got)
		}
	}
}

func TestSeedDocFromHTML_Lists(t *testing.T) {
	got := seedToXML(t, "<ul><li>无序</li><li>嵌套<ol><li>有序</li></ol></li></ul>")
	want := "<bulletList>" +
		"<listItem><paragraph>无序</paragraph></listItem>" +
		"<listItem><paragraph>嵌套</paragraph><orderedList><listItem><paragraph>有序</paragraph></listItem></orderedList></listItem>" +
		"</bulletList>"
	if got != want {
		t.Fatalf("seed mismatch:\n got: %s\nwant: %s", got, want)
	}
}

func TestSeedDocFromHTML_BlockquoteAndCode(t *testing.T) {
	got := seedToXML(t, "<blockquote><p>引用</p></blockquote><pre><code class=\"language-go\">fmt.Println(1)</code></pre>")
	want := "<blockquote><paragraph>引用</paragraph></blockquote>" +
		"<codeBlock language=\"go\">fmt.Println(1)</codeBlock>"
	if got != want {
		t.Fatalf("seed mismatch:\n got: %s\nwant: %s", got, want)
	}
}

func TestSeedDocFromHTML_UnknownBlockDegradesToParagraph(t *testing.T) {
	got := seedToXML(t, "<div><section>降级文本</section></div>")
	if got != "<paragraph>降级文本</paragraph>" {
		t.Fatalf("unexpected degrade: %s", got)
	}
}

func TestSeedDocFromHTML_TableDegradesToText(t *testing.T) {
	got := seedToXML(t, "<table><tr><td>单元格</td></tr></table>")
	if !strings.Contains(got, "单元格") {
		t.Fatalf("table text lost: %s", got)
	}
}

func TestSeedDocFromHTML_EmptyAndWhitespace(t *testing.T) {
	if got := seedToXML(t, ""); got != "" {
		t.Fatalf("empty html should seed nothing, got %s", got)
	}
	if got := seedToXML(t, "<p>   </p><p></p>"); got != "" {
		t.Fatalf("whitespace-only should seed nothing, got %s", got)
	}
}

func TestSeedDocFromHTML_StrayTextBetweenBlocks(t *testing.T) {
	got := seedToXML(t, "<p>段一</p>游离文本")
	if !strings.Contains(got, "<paragraph>游离文本</paragraph>") {
		t.Fatalf("stray text should become paragraph: %s", got)
	}
}
