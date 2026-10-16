package collab

import (
	"fmt"
	"strings"

	"github.com/reearth/ygo/crdt"
	"golang.org/x/net/html"
)

// 受控 HTML → yjs XML 转换器：服务端权威种子化用。
// 只处理本仓 TipTap schema 已知节点；未知块级标签递归降级、未知行内标签取文本。
// 文本 marks 按 y-tiptap 的 delta 约定写入（markName → mark attrs）；
// 元素属性按 y-tiptap 约定原样写值（level 为数字，非 JSON 字符串）。

// seedDocFromHTML 把 HTML 灌入 ydoc 的 "default" XmlFragment（TipTap Collaboration 字段名）。
// 仅当文档为空时调用（空判断在调用方）。html 解析失败按空内容处理（不破坏文档）。
// seedDocFromHTML 把 HTML 灌入 ydoc 的 "default" XmlFragment（TipTap Collaboration 字段名）。
// 仅当文档为空时调用（空判断在调用方）。返回错误供调用方记录（转换/事务失败不得静默）。
func seedDocFromHTML(doc *crdt.Doc, htmlContent string) error {
	if strings.TrimSpace(htmlContent) == "" {
		return nil
	}
	root, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return fmt.Errorf("parse html: %w", err)
	}
	body := findBody(root)
	if body == nil {
		return nil
	}
	// ygo 铁律：Transact 内调用 doc.GetXmlFragment 等读方法会死锁——
	// 引用必须在事务外取好
	frag := doc.GetXmlFragment("default")
	return doc.TransactE(func(txn *crdt.Transaction) error {
		idx := 0
		for n := body.FirstChild; n != nil; n = n.NextSibling {
			idx = insertBlockNode(frag, txn, idx, n)
		}
		return nil
	})
}

func findBody(n *html.Node) *html.Node {
	if n.Type == html.ElementNode && n.Data == "body" {
		return n
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if b := findBody(c); b != nil {
			return b
		}
	}
	return nil
}

// insertBlockNode 插入一个块级节点，返回新的插入位置。
func insertBlockNode(frag *crdt.YXmlFragment, txn *crdt.Transaction, idx int, n *html.Node) int {
	if n.Type == html.TextNode {
		if strings.TrimSpace(n.Data) == "" {
			return idx
		}
		return insertParagraphWithRuns(frag, txn, idx, []inlineRun{{text: n.Data}})
	}
	if n.Type != html.ElementNode {
		return idx
	}
	switch n.Data {
	case "p":
		return insertParagraphWithRuns(frag, txn, idx, runsOf(collectInline(n, nil)))
	case "h1", "h2", "h3", "h4", "h5", "h6":
		elem := crdt.NewYXmlElement("heading")
		elem.SetAttributeValue(txn, "level", int(n.Data[1]-'0'))
		fillInlineContainer(elem, txn, collectInline(n, nil))
		frag.InsertElement(txn, idx, elem)
		return idx + 1
	case "ul":
		return insertListInto(frag, txn, idx, n, "bulletList")
	case "ol":
		return insertListInto(frag, txn, idx, n, "orderedList")
	case "blockquote":
		elem := crdt.NewYXmlElement("blockquote")
		frag.InsertElement(txn, idx, elem)
		inner := 0
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			inner = insertIntoContainer(elem, txn, inner, c)
		}
		return idx + 1
	case "pre":
		return insertCodeBlock(frag, txn, idx, n)
	case "hr":
		frag.InsertElement(txn, idx, crdt.NewYXmlElement("horizontalRule"))
		return idx + 1
	case "table":
		// v1 种子化不支持表格：降级为逐行段落文本
		return insertParagraphWithRuns(frag, txn, idx, runsOf(collectInline(n, nil)))
	default:
		// div/section/span 等未知块级：递归子节点；无子块时当普通段落
		childBlocks := false
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode && isKnownBlock(c.Data) {
				childBlocks = true
				break
			}
		}
		if childBlocks {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				idx = insertBlockNode(frag, txn, idx, c)
			}
			return idx
		}
		return insertParagraphWithRuns(frag, txn, idx, runsOf(collectInline(n, nil)))
	}
}

// insertIntoContainer blockquote 内部的递归插入（块级降级为段落）。
func insertIntoContainer(elem *crdt.YXmlElement, txn *crdt.Transaction, idx int, n *html.Node) int {
	if n.Type == html.TextNode {
		if strings.TrimSpace(n.Data) == "" {
			return idx
		}
		p := crdt.NewYXmlElement("paragraph")
		elem.InsertElement(txn, idx, p)
		fillInlineContainer(p, txn, runsToItems([]inlineRun{{text: n.Data}}))
		return idx + 1
	}
	if n.Type != html.ElementNode {
		return idx
	}
	switch n.Data {
	case "p":
		p := crdt.NewYXmlElement("paragraph")
		elem.InsertElement(txn, idx, p)
		fillInlineContainer(p, txn, collectInline(n, nil))
		return idx + 1
	case "ul", "ol":
		listType := "bulletList"
		if n.Data == "ol" {
			listType = "orderedList"
		}
		list := crdt.NewYXmlElement(listType)
		elem.InsertElement(txn, idx, list)
		fillListItems(list, txn, n)
		return idx + 1
	default:
		p := crdt.NewYXmlElement("paragraph")
		elem.InsertElement(txn, idx, p)
		fillInlineContainer(p, txn, collectInline(n, nil))
		return idx + 1
	}
}

func insertListInto(frag *crdt.YXmlFragment, txn *crdt.Transaction, idx int, n *html.Node, listType string) int {
	list := crdt.NewYXmlElement(listType)
	frag.InsertElement(txn, idx, list)
	fillListItems(list, txn, n)
	return idx + 1
}

func fillListItems(list *crdt.YXmlElement, txn *crdt.Transaction, n *html.Node) {
	i := 0
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type != html.ElementNode || c.Data != "li" {
			continue
		}
		li := crdt.NewYXmlElement("listItem")
		list.InsertElement(txn, i, li)
		fillInlineContainer(li, txn, collectInline(c, nil))
		// 嵌套列表
		j := 1
		for cc := c.FirstChild; cc != nil; cc = cc.NextSibling {
			if cc.Type == html.ElementNode && (cc.Data == "ul" || cc.Data == "ol") {
				listType := "bulletList"
				if cc.Data == "ol" {
					listType = "orderedList"
				}
				sub := crdt.NewYXmlElement(listType)
				li.InsertElement(txn, j, sub)
				fillListItems(sub, txn, cc)
				j++
			}
		}
		i++
	}
}

func insertCodeBlock(frag *crdt.YXmlFragment, txn *crdt.Transaction, idx int, n *html.Node) int {
	code := crdt.NewYXmlElement("codeBlock")
	language := ""
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "code" {
			for _, a := range c.Attr {
				if a.Key == "class" && strings.HasPrefix(a.Val, "language-") {
					language = strings.TrimPrefix(a.Val, "language-")
				}
			}
			code.SetAttributeValue(txn, "language", language)
			text := collectText(c)
			if text != "" {
				txt := crdt.NewYXmlText()
				code.InsertText(txn, 0, txt)
				txt.Insert(txn, 0, text, nil)
			}
			frag.InsertElement(txn, idx, code)
			return idx + 1
		}
	}
	// 无 code 子元素的 pre：整块文本
	code.SetAttributeValue(txn, "language", language)
	if text := collectText(n); text != "" {
		txt := crdt.NewYXmlText()
		code.InsertText(txn, 0, txt)
		txt.Insert(txn, 0, text, nil)
	}
	frag.InsertElement(txn, idx, code)
	return idx + 1
}

// insertParagraphWithRuns 插入一个以 runs 为内容的段落。
func insertParagraphWithRuns(frag *crdt.YXmlFragment, txn *crdt.Transaction, idx int, runs []inlineRun) int {
	return insertParagraphWithItems(frag, txn, idx, runsToItems(runs))
}

// insertParagraphWithItems 插入一个以行内项（文本/行内元素混合）为内容的段落。
func insertParagraphWithItems(frag *crdt.YXmlFragment, txn *crdt.Transaction, idx int, items []inlineItem) int {
	p := crdt.NewYXmlElement("paragraph")
	frag.InsertElement(txn, idx, p)
	fillInlineContainer(p, txn, items)
	return idx + 1
}
