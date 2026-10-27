package service

import (
	"testing"
)

func TestExtractExternalLinks(t *testing.T) {
	html := `<p>see <a href="https://a.example.com/x">a</a> and <a href='https://b.example.com/y'>b</a>
	         also <a href="/relative">rel</a> dup <a href="https://a.example.com/x">a</a></p>`
	got := extractExternalLinks(html)
	if len(got) != 2 || got[0] != "https://a.example.com/x" || got[1] != "https://b.example.com/y" {
		t.Fatalf("html links = %v", got)
	}

	md := `看 [这个](https://c.example.com/z)。正文裸文本 URL 不抽取。`
	got = extractExternalLinks(md)
	if len(got) != 1 || got[0] != "https://c.example.com/z" {
		t.Fatalf("md links = %v", got)
	}

	if got := extractExternalLinks(""); got != nil {
		t.Fatalf("empty content should give nil")
	}
}
