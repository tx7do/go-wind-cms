package service

import (
	"reflect"
	"testing"
)

// TestVectorToLiteral 向量字面量序列化：pgvector 接受的 '[a,b,...]' 形态，
// 空切片落 '[]'（由列定维与维度校验兜底，不会成 0 维行）。
func TestVectorToLiteral(t *testing.T) {
	cases := []struct {
		in   []float32
		want string
	}{
		{[]float32{0.5, -1.25, 2}, "[0.5,-1.25,2]"},
		{[]float32{1}, "[1]"},
		{nil, "[]"},
	}
	for _, c := range cases {
		if got := vectorToLiteral(c.in); got != c.want {
			t.Fatalf("vectorToLiteral(%v) = %q, want %q", c.in, got, c.want)
		}
	}
}

// TestSplitIntoChunks 切片窗口/重叠语义（生产参数 500/50 的缩小版）：
// step = size - overlap，尾窗截断到串尾；overlap >= size 时 step 回退为 size；
// 空串返回 nil；内容短于窗口时单整片。RAG 入库与重索引共用该切分，
// 边界漂移会静默改变切片数与检索粒度，这里钉住。
func TestSplitIntoChunks(t *testing.T) {
	got := splitIntoChunks("abcdefghijklmnopqrstuvwxyz", 10, 3)
	want := []string{
		"abcdefghij", // 起点 0：窗口 [0,10)
		"hijklmnopq", // 起点 7：窗口 [7,17)
		"opqrstuvwx", // 起点 14：窗口 [14,24)
		"vwxyz",      // 起点 21：窗口 [21,26)，尾窗截断到串尾
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("splitIntoChunks(26, size=10, overlap=3) = %v, want %v", got, want)
	}

	// overlap == size → step 回退为 size（无重叠推进）
	got = splitIntoChunks("abcdefg", 5, 5)
	want = []string{"abcde", "fg"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("splitIntoChunks(7, size=5, overlap=5) = %v, want %v", got, want)
	}

	// 内容短于窗口：单整片
	got = splitIntoChunks("abc", 10, 3)
	want = []string{"abc"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("splitIntoChunks(short) = %v, want %v", got, want)
	}

	// 空串：nil
	if got = splitIntoChunks("", 10, 3); got != nil {
		t.Fatalf("splitIntoChunks(empty) = %v, want nil", got)
	}
}
