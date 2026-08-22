// Package mdprotect 在机器翻译前后保护 Markdown 中不应被翻译的片段
// （围栏代码块、行内代码、链接/图片的 URL、裸 URL），防止翻译引擎
// 破坏代码与超链接结构。用法：
//
//	masked, restore := mdprotect.Protect(content)
//	translated, err := engine.Translate(masked, from, to)
//	result := restore(translated)
//
// 局限：占位符为 [[[n]]] 形式的 token，主流引擎会原样保留；若引擎
// 改写/丢弃占位符（如把其中数字翻译掉），恢复会失败并原样保留该 token。
package mdprotect

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// 占位符形如 [[[0]]]。恢复时容忍引擎在占位符内部插入的空白。
var tokenRe = regexp.MustCompile(`\[\[\[\s*(\d+)\s*\]\]\]`)

var (
	// 围栏代码块（``` 或 ~~~），跨行
	fencedCodeRe = regexp.MustCompile("(?s)```.*?```|~~~.*?~~~")
	// 行内代码，不跨行
	inlineCodeRe = regexp.MustCompile("`[^`\n]+`")
	// markdown 链接/图片的目标部分：](url)，保留 ]( ) 结构只遮 URL
	linkTargetRe = regexp.MustCompile(`\]\(\s*(https?://[^)\s]+)\s*\)`)
	// 裸 URL（自动链接、正文引用等）
	bareURLRe = regexp.MustCompile(`https?://[^\s)>\]]+`)
)

type protector struct {
	segments []string
}

// Protect 返回遮蔽后的文本与恢复函数。对不含可保护片段的纯文本，
// 遮蔽是恒等变换。
func Protect(src string) (string, func(string) string) {
	p := &protector{}
	masked := p.mask(src)
	return masked, p.restore
}

// add 记录一个受保护片段并返回其占位符。
// 占位符编号必须取"追加前"的下标（len-1 在追加后），封装成单一方法
// 避免调用方 append/编号顺序写反造成错位。
func (p *protector) add(segment string) string {
	p.segments = append(p.segments, segment)
	return fmt.Sprintf("[[[%d]]]", len(p.segments)-1)
}

func (p *protector) mask(src string) string {
	// 顺序：先整体块（围栏代码），再行内片段，最后 URL。
	// 已替换出的占位符不含反引号与 URL 字符，不会被后续规则再次命中。
	out := fencedCodeRe.ReplaceAllStringFunc(src, func(m string) string {
		return p.add(m)
	})

	out = inlineCodeRe.ReplaceAllStringFunc(out, func(m string) string {
		return p.add(m)
	})

	out = linkTargetRe.ReplaceAllStringFunc(out, func(m string) string {
		sub := linkTargetRe.FindStringSubmatch(m)
		if len(sub) < 2 {
			return m
		}
		return "](" + p.add(sub[1]) + ")"
	})

	out = bareURLRe.ReplaceAllStringFunc(out, func(m string) string {
		return p.add(m)
	})

	return out
}

func (p *protector) restore(translated string) string {
	return tokenRe.ReplaceAllStringFunc(translated, func(m string) string {
		sub := tokenRe.FindStringSubmatch(m)
		if len(sub) < 2 {
			return m
		}
		idx, err := strconv.Atoi(strings.TrimSpace(sub[1]))
		if err != nil || idx < 0 || idx >= len(p.segments) {
			return m
		}
		return p.segments[idx]
	})
}
