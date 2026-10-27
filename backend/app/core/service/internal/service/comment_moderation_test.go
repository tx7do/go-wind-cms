package service

import (
	"testing"
)

func TestParseCommentVerdict(t *testing.T) {
	cases := []struct {
		name    string
		text    string
		spam    bool
		wantErr bool
	}{
		{name: "plain json", text: `{"spam":true,"reason":"ad"}`, spam: true},
		{name: "fenced json", text: "```json\n{\"spam\":false,\"reason\":\"ok\"}\n```", spam: false},
		{name: "prose around json", text: `判定如下 {"spam":true,"reason":"seo"} 完毕`, spam: true},
		{name: "missing fields", text: `{}`, spam: false},
		{name: "invalid", text: `not json at all`, wantErr: true},
	}
	for _, c := range cases {
		v, err := parseCommentVerdict(c.text)
		if c.wantErr {
			if err == nil {
				t.Errorf("%s: want error, got %+v", c.name, v)
			}
			continue
		}
		if err != nil {
			t.Errorf("%s: unexpected error %v", c.name, err)
			continue
		}
		if v.Spam != c.spam {
			t.Errorf("%s: spam = %v, want %v", c.name, v.Spam, c.spam)
		}
	}
	// 长内容截断上限只影响送审文本（这里验证常量语义存在）
	if commentAIModerationMaxChars <= 0 || commentAIModerationMaxChars < 100 {
		t.Fatalf("moderation max chars unreasonable: %d", commentAIModerationMaxChars)
	}
}
