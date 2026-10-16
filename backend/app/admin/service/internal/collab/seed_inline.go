package collab

import (
	"strings"

	"github.com/reearth/ygo/crdt"
	"golang.org/x/net/html"
)

// 行内内容收集与灌入：把行内 DOM 归一为 runs（文本 + mark 属性）与
// 行内元素（image/hardBreak），按 y-tiptap 的 YText delta 约定写入。

type inlineRun struct {
	text  string
	marks map[string]any
}

type inlineItem struct {
	run       *inlineRun
	elem      *crdt.YXmlElement // image / hardBreak 等行内元素
	elemAttrs map[string]any    // 行内元素属性（须在事务内写入）
}

func isKnownBlock(tag string) bool {
	switch tag {
	case "p", "h1", "h2", "h3", "h4", "h5", "h6", "ul", "ol", "blockquote", "pre", "hr", "table":
		return true
	}
	return false
}

// collectInline 收集行内内容（含 marks 聚合与行内元素）。
func collectInline(n *html.Node, marks map[string]any) []inlineItem {
	var items []inlineItem
	appendRun := func(text string, m map[string]any) {
		if text == "" {
			return
		}
		// 合并相邻同 marks 的文本
		if len(items) > 0 && items[len(items)-1].run != nil && sameMarks(items[len(items)-1].run.marks, m) {
			items[len(items)-1].run.text += text
			return
		}
		items = append(items, inlineItem{run: &inlineRun{text: text, marks: m}})
	}

	var walk func(n *html.Node, m map[string]any)
	walk = func(n *html.Node, m map[string]any) {
		switch n.Type {
		case html.TextNode:
			appendRun(n.Data, m)
			return
		case html.ElementNode:
		default:
			return
		}

		switch n.Data {
		case "br":
			items = append(items, inlineItem{elem: crdt.NewYXmlElement("hardBreak")})
			return
		case "img":
			// 元素在事务内插入后再写属性（此处不在事务中，属性暂存）
			items = append(items, inlineItem{elem: crdt.NewYXmlElement("image"), elemAttrs: imgAttrs(n)})
			return
		case "strong", "b":
			m = withMark(m, "bold", map[string]any{})
		case "em", "i":
			m = withMark(m, "italic", map[string]any{})
		case "u":
			m = withMark(m, "underline", map[string]any{})
		case "s", "del", "strike":
			m = withMark(m, "strike", map[string]any{})
		case "code":
			m = withMark(m, "code", map[string]any{})
		case "a":
			attrs := map[string]any{}
			for _, a := range n.Attr {
				if a.Key == "href" && a.Val != "" {
					attrs["href"] = a.Val
				}
			}
			m = withMark(m, "link", attrs)
		case "script", "style":
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c, m)
		}
	}
	walk(n, marks)
	return items
}

// fillInlineContainer 把 runs 灌入块级容器（paragraph 等）。
func fillInlineContainer(elem *crdt.YXmlElement, txn *crdt.Transaction, items []inlineItem) {
	i := 0
	offset := 0
	var txt *crdt.YXmlText
	ensureText := func() *crdt.YXmlText {
		if txt == nil {
			txt = crdt.NewYXmlText()
			elem.InsertText(txn, i, txt)
			i++
		}
		return txt
	}
	for _, it := range items {
		if it.elem != nil {
			elem.InsertElement(txn, i, it.elem)
			for k, v := range it.elemAttrs {
				it.elem.SetAttributeValue(txn, k, v)
			}
			i++
			txt = nil
			offset = 0
			continue
		}
		t := ensureText()
		t.Insert(txn, offset, it.run.text, it.run.marks)
		offset += len([]rune(it.run.text))
	}
}

// imgAttrs 提取 image 行内元素属性。
func imgAttrs(n *html.Node) map[string]any {
	attrs := map[string]any{}
	for _, a := range n.Attr {
		switch a.Key {
		case "src", "alt", "title":
			if a.Val != "" {
				attrs[a.Key] = a.Val
			}
		}
	}
	return attrs
}

// withMark 追加 mark（拷贝后追加，避免污染兄弟分支共享的 map）。
func withMark(marks map[string]any, name string, attrs map[string]any) map[string]any {
	next := make(map[string]any, len(marks)+1)
	for k, v := range marks {
		next[k] = v
	}
	next[name] = attrs
	return next
}

func sameMarks(a, b map[string]any) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		bv, ok := b[k]
		if !ok {
			return false
		}
		av, aok := v.(map[string]any)
		bv2, bok := bv.(map[string]any)
		if aok != bok {
			return false
		}
		if aok && !sameStringMap(av, bv2) {
			return false
		}
	}
	return true
}

func sameStringMap(a, b map[string]any) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}

func collectText(n *html.Node) string {
	var b strings.Builder
	var walk func(n *html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.TextNode {
			b.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(n)
	return b.String()
}

// runsOf 把 inlineItem 列表里的文本项展开为 runs（供 insertParagraphWithRuns）。
func runsOf(items []inlineItem) []inlineRun {
	var runs []inlineRun
	for _, it := range items {
		if it.run != nil {
			runs = append(runs, *it.run)
		}
	}
	return runs
}

// runsToItems runs → 行内项（逆变换）。
func runsToItems(runs []inlineRun) []inlineItem {
	items := make([]inlineItem, 0, len(runs))
	for i := range runs {
		items = append(items, inlineItem{run: &runs[i]})
	}
	return items
}
