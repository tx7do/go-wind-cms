package oss

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestContentTypeToBucketName(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"image jpeg", "image/jpeg", "images"},
		{"image uppercase", "IMAGE/JPEG", "images"},
		{"image svg xml", "image/svg+xml", "images"},
		{"text plain with charset", "text/plain; charset=utf-8", "docs"},
		{"text html mixed case", "Text/HTML; Charset=UTF-8", "docs"},
		{"application pdf", "application/pdf", "docs"},
		{"application json", "application/json", "docs"},
		{"application word (openxml)", "application/vnd.openxmlformats-officedocument.wordprocessingml.document", "docs"},
		{"application ms-excel", "application/vnd.ms-excel", "docs"},
		{"video mp4", "video/mp4", "videos"},
		{"audio mpeg", "audio/mpeg", "audios"},
		{"application octet-stream", "application/octet-stream", "files"},
		{"unknown no slash", "invalid-type", "files"},
		{"application unknown subtype", "application/unknown", "files"},
		{"empty input", "", "files"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := ContentTypeToBucketName(tt.in)
			if got != tt.want {
				t.Errorf("ContentTypeToBucketName(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestFileExtensionToBucketName(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"image mime", "image/jpeg", BucketImages},
		{"jpg lowercase", "jpg", BucketImages},
		{"jpg uppercase with dot", ".JPG", BucketImages},
		{"png with dot", ".png", BucketImages},
		{"video mp4", "mp4", BucketVideos},
		{"video MP4 uppercase with dot", ".MP4", BucketVideos},
		{"audio mp3", "mp3", BucketAudios},
		{"js extension", "js", BucketDocs},
		{"json extension", "json", BucketDocs},
		{"pdf extension", "pdf", BucketDocs},
		{"text mime with charset", "text/plain; charset=utf-8", BucketDocs},
		{"application octet-stream mime", "application/octet-stream", BucketFiles},
		{"unknown extension", "unknownext", BucketFiles},
		{"empty input", "", BucketFiles},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := FileExtensionToBucketName(tt.in)
			if got != tt.want {
				t.Errorf("FileExtensionToBucketName(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestContentTypeToFileExtension_Rewritten(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"image jpeg", "image/jpeg", ".jpg"},
		{"image jpeg uppercase", "IMAGE/JPEG", ".jpg"},
		{"image svg xml", "image/svg+xml", ".svg"},
		{"image x-icon", "image/x-icon", ".ico"},
		{"text plain with charset", "text/plain; charset=utf-8", ".txt"},
		{"text html", "text/html", ".html"},
		{"text css", "text/css", ".css"},
		{"js text/javascript", "text/javascript", ".js"},
		{"js application/javascript", "application/javascript", ".js"},
		{"js x-javascript", "application/x-javascript", ".js"},
		{"lua text/x-lua", "text/x-lua", ".lua"},
		{"lua application/x-lua", "application/x-lua", ".lua"},
		{"python text/x-python", "text/x-python", ".py"},
		{"python application/x-python", "application/x-python", ".py"},
		{"shell text/x-shellscript", "text/x-shellscript", ".sh"},
		{"shell application/x-sh", "application/x-sh", ".sh"},
		{"video mp4", "video/mp4", ".mp4"},
		{"audio mpeg", "audio/mpeg", ".mp3"},
		{"application pdf", "application/pdf", ".pdf"},
		{"application json", "application/json", ".json"},
		{"application zip", "application/zip", ".zip"},
		{"application x-tar", "application/x-tar", ".tar"},
		{"application gzip", "application/gzip", ".gz"},
		{"application 7z", "application/x-7z-compressed", ".7z"},
		{"octet stream", "application/octet-stream", ".bin"},
		{"text xml fallback (mime.ExtensionsByType)", "text/xml", ".xml"},
		{"unknown application subtype", "application/unknown", ""},
		{"invalid mime without slash", "not-a-mime", ""},
		{"empty input", "", ""},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := ContentTypeToFileExtension(tt.in)
			if got != tt.want {
				t.Errorf("ContentTypeToFileExtension(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestGeneraUUIDFileName_NoExt(t *testing.T) {
	got := GeneraUUIDFileName("")
	if len(got) != 32 {
		t.Fatalf("expected length 32, got %d (%q)", len(got), got)
	}
	matched, _ := regexp.MatchString(`^[0-9a-f]{32}$`, got)
	if !matched {
		t.Fatalf("expected 32 lowercase hex chars, got %q", got)
	}
}

func TestGeneraUUIDFileName_WithDotExt_Lowercase(t *testing.T) {
	got := GeneraUUIDFileName(".jpg")
	if !strings.HasSuffix(got, ".jpg") {
		t.Fatalf("expected suffix .jpg, got %q", got)
	}
	if len(got) != 32+1+3 {
		t.Fatalf("expected total length %d, got %d (%q)", 36, len(got), got)
	}
	matched, _ := regexp.MatchString(`^[0-9a-f]{32}\.jpg$`, got)
	if !matched {
		t.Fatalf("unexpected format: %q", got)
	}
}

func TestGeneraUUIDFileName_WithExt_NoLeadingDot_And_CasePreserved(t *testing.T) {
	got := GeneraUUIDFileName("txt")
	if !strings.HasSuffix(got, ".txt") {
		t.Fatalf("expected suffix .txt, got %q", got)
	}

	got2 := GeneraUUIDFileName(".JPG")
	if !strings.HasSuffix(got2, ".JPG") {
		t.Fatalf("expected suffix .JPG (case preserved), got %q", got2)
	}
	matched, _ := regexp.MatchString(`^[0-9a-f]{32}\.JPG$`, got2)
	if !matched {
		t.Fatalf("unexpected format for uppercase ext: %q", got2)
	}
}

func TestGeneraUUIDFileName_Uniqueness(t *testing.T) {
	a := GeneraUUIDFileName(".png")
	b := GeneraUUIDFileName(".png")
	if a == b {
		t.Fatalf("expected two generated names to differ, both: %q", a)
	}
	t.Logf("Generated names: %q and %q", a, b)
}

func TestGeneraContentSHA265FileName_NoExt(t *testing.T) {
	t.Parallel()
	// SHA256("hello") 已知值
	want := "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
	got := GeneraContentSHA265FileName([]byte("hello"), "")
	if got != want {
		t.Fatalf("GeneraContentSHA265FileName(%q, %q) = %q, want %q", "hello", "", got, want)
	}
}

func TestGeneraContentSHA265FileName_WithDotExt(t *testing.T) {
	t.Parallel()
	got := GeneraContentSHA265FileName([]byte("some data"), ".txt")
	if !strings.HasSuffix(got, ".txt") {
		t.Fatalf("expected suffix .txt, got %q", got)
	}
	matched, _ := regexp.MatchString(`^[0-9a-f]{64}\.txt$`, got)
	if !matched {
		t.Fatalf("unexpected format: %q", got)
	}
	t.Logf("Generated name: %q", got)
}

func TestGeneraContentSHA265FileName_WithExt_NoLeadingDot_And_CasePreserved(t *testing.T) {
	t.Parallel()
	got := GeneraContentSHA265FileName([]byte("some data"), "JPG")
	if !strings.HasSuffix(got, ".JPG") {
		t.Fatalf("expected suffix .JPG, got %q", got)
	}
	matched, _ := regexp.MatchString(`^[0-9a-f]{64}\.JPG$`, got)
	if !matched {
		t.Fatalf("unexpected format for uppercase ext: %q", got)
	}
}

func TestGeneraContentSHA265FileName_DifferentContentDifferentName(t *testing.T) {
	t.Parallel()
	a := GeneraContentSHA265FileName([]byte("first"), ".bin")
	b := GeneraContentSHA265FileName([]byte("second"), ".bin")
	if a == b {
		t.Fatalf("expected different names for different content, both: %q", a)
	}
	// 验证哈希长度和后缀
	if len(strings.TrimSuffix(a, ".bin")) != 64 || len(strings.TrimSuffix(b, ".bin")) != 64 {
		t.Fatalf("expected 64-hex hash prefix, got %q and %q", a, b)
	}
}

func TestGenerateHMACContentFileName_DeterministicAndMatchesManual(t *testing.T) {
	t.Parallel()
	secret := []byte("test-secret-123")
	content := []byte("hello hmac content")
	ext := ".bin"

	// 调用被测函数
	got := GenerateHMACContentFileName(secret, content, ext)
	got2 := GenerateHMACContentFileName(secret, content, ext)
	if got != got2 {
		t.Fatalf("expected deterministic results, got %q and %q", got, got2)
	}

	// 手工计算期望值（与函数实现一致的方式）
	sum := sha256.Sum256(content)
	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write(sum[:])
	expected := hex.EncodeToString(mac.Sum(nil))
	expectedWithExt := expected + ".bin"

	if got != expectedWithExt {
		t.Fatalf("expected %q, got %q", expectedWithExt, got)
	}
}

func TestGenerateHMACContentFileName_NoExt_Format(t *testing.T) {
	t.Parallel()
	secret := []byte("s")
	content := []byte("payload")
	got := GenerateHMACContentFileName(secret, content, "")
	// 仅 64 个小写十六进制字符
	matched, _ := regexp.MatchString(`^[0-9a-f]{64}$`, got)
	if !matched {
		t.Fatalf("unexpected format without ext: %q", got)
	}
	t.Logf("Generated name without ext: %q", got)
}

func TestGenerateHMACContentFileName_WithExt_CasePreserved(t *testing.T) {
	t.Parallel()
	secret := []byte("s2")
	content := []byte("payload")
	got := GenerateHMACContentFileName(secret, content, ".txt")
	if !strings.HasSuffix(got, ".txt") {
		t.Fatalf("expected suffix .txt, got %q", got)
	}
	matched, _ := regexp.MatchString(`^[0-9a-f]{64}\.txt$`, got)
	if !matched {
		t.Fatalf("unexpected format with .txt ext: %q", got)
	}

	got2 := GenerateHMACContentFileName(secret, content, "JPG")
	if !strings.HasSuffix(got2, ".JPG") {
		t.Fatalf("expected suffix .JPG (case preserved), got %q", got2)
	}
	matched2, _ := regexp.MatchString(`^[0-9a-f]{64}\.JPG$`, got2)
	if !matched2 {
		t.Fatalf("unexpected format with JPG ext: %q", got2)
	}
}

func TestGenerateHMACContentFileName_DifferentSecretOrContent(t *testing.T) {
	t.Parallel()
	secretA := []byte("a")
	secretB := []byte("b")
	contentA := []byte("one")
	contentB := []byte("two")

	a := GenerateHMACContentFileName(secretA, contentA, ".bin")
	b := GenerateHMACContentFileName(secretB, contentA, ".bin")
	if a == b {
		t.Fatalf("expected different names for different secrets, both: %q", a)
	}
	t.Logf("Generated names with different secrets: %q and %q", a, b)

	c := GenerateHMACContentFileName(secretA, contentB, ".bin")
	if a == c {
		t.Fatalf("expected different names for different content, both: %q", a)
	}
	t.Logf("Generated names with different content: %q and %q", a, c)
}

func TestGeneraTimeBaseFileName_NoExt(t *testing.T) {
	got := GeneraTimeBaseFileName("")

	// 形如: <timestamp>_<8hex>
	parts := strings.SplitN(got, "_", 2)
	if len(parts) != 2 {
		t.Fatalf("unexpected format, want timestamp_suffix, got: %q", got)
	}

	tsStr, suffix := parts[0], parts[1]

	// timestamp 可解析为整数（可能是秒/毫秒/纳秒），将其归一为秒并接近当前时间（允许 5s 误差）
	tsRaw, err := strconv.ParseInt(tsStr, 10, 64)
	if err != nil {
		t.Fatalf("timestamp not integer: %q", tsStr)
	}

	// 归一化到秒：如果数值很大，可能是毫秒或纳秒
	ts := tsRaw
	switch {
	case tsRaw > 1e15:
		// very large -> nanoseconds
		ts = tsRaw / 1e9
	case tsRaw > 1e12:
		// medium large -> milliseconds
		ts = tsRaw / 1e3
	default:
		// already seconds
		ts = tsRaw
	}

	now := time.Now().Unix()
	if diff := now - ts; diff < -5 || diff > 5 {
		t.Fatalf("timestamp %d not within 5s of now %d", ts, now)
	}

	// suffix 8 个小写十六进制字符
	matched, _ := regexp.MatchString(`^[0-9a-f]{8}$`, suffix)
	if !matched {
		t.Fatalf("suffix not 8 hex chars: %q", suffix)
	}
}

func TestGeneraTimeBaseFileName_WithDotExt(t *testing.T) {
	t.Parallel()
	got := GeneraTimeBaseFileName(".txt")
	// 应以 ".txt" 结尾
	if !strings.HasSuffix(got, ".txt") {
		t.Fatalf("expected suffix .txt, got %q", got)
	}

	// 去掉扩展名后检查前缀格式
	base := strings.TrimSuffix(got, ".txt")
	parts := strings.SplitN(base, "_", 2)
	if len(parts) != 2 {
		t.Fatalf("unexpected format before ext, got: %q", base)
	}
	// 校验 timestamp 可解析
	if _, err := strconv.ParseInt(parts[0], 10, 64); err != nil {
		t.Fatalf("timestamp not integer: %q", parts[0])
	}
	// 校验后缀 8 hex
	matched, _ := regexp.MatchString(`^[0-9a-f]{8}$`, parts[1])
	if !matched {
		t.Fatalf("suffix not 8 hex chars: %q", parts[1])
	}
}

func TestGeneraTimeBaseFileName_WithExt_NoLeadingDot_And_CasePreserved(t *testing.T) {
	t.Parallel()
	got := GeneraTimeBaseFileName("JPG")
	// 应以 ".JPG" 结尾（大小写保留）
	if !strings.HasSuffix(got, ".JPG") {
		t.Fatalf("expected suffix .JPG, got %q", got)
	}
	// 验证前缀格式
	base := strings.TrimSuffix(got, ".JPG")
	parts := strings.SplitN(base, "_", 2)
	if len(parts) != 2 {
		t.Fatalf("unexpected format before ext, got: %q", base)
	}
	matched, _ := regexp.MatchString(`^[0-9a-f]{8}$`, parts[1])
	if !matched {
		t.Fatalf("suffix not 8 hex chars: %q", parts[1])
	}
}

func TestGeneraTimeBaseFileName_Uniqueness(t *testing.T) {
	t.Parallel()
	a := GeneraTimeBaseFileName(".bin")
	b := GeneraTimeBaseFileName(".bin")
	if a == b {
		t.Fatalf("expected two generated names to differ, both: %q", a)
	}
	t.Logf("Generated names: %q and %q", a, b)
}
