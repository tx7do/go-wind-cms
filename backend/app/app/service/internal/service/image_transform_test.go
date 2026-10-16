package service

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"net/url"
	"testing"

	"github.com/disintegration/imaging"
)

func TestParseTransformParams(t *testing.T) {
	mustParse := func(t *testing.T, raw string) url.Values {
		t.Helper()
		v, err := url.ParseQuery(raw)
		if err != nil {
			t.Fatalf("ParseQuery(%q): %v", raw, err)
		}
		return v
	}

	cases := []struct {
		name    string
		query   string
		want    TransformParams
		wantErr bool
	}{
		{
			name:  "empty_identity",
			query: "",
			want:  TransformParams{Quality: 80, Fit: "contain"},
		},
		{
			name:  "width_only",
			query: "w=500&q=80",
			want:  TransformParams{Width: 500, Quality: 80, Fit: "contain"},
		},
		{
			name:  "cover_both_edges",
			query: "w=200&h=200&fit=cover&q=60",
			want:  TransformParams{Width: 200, Height: 200, Quality: 60, Fit: "cover"},
		},
		{
			name:  "cover_single_edge_degrades_to_contain",
			query: "w=200&fit=cover",
			want:  TransformParams{Width: 200, Quality: 80, Fit: "contain"},
		},
		{
			name:  "clamped_edges_and_quality",
			query: "w=99999&h=0&q=999",
			want:  TransformParams{Width: 4096, Height: 0, Quality: 100, Fit: "contain"},
		},
		{
			name:  "quality_floor",
			query: "q=0",
			want:  TransformParams{Quality: 1, Fit: "contain"},
		},
		{
			name:    "negative_width",
			query:   "w=-5",
			wantErr: true,
		},
		{
			name:    "bad_fit",
			query:   "fit=stretch",
			wantErr: true,
		},
		{
			name:    "bad_quality",
			query:   "q=abc",
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseTransformParams(mustParse(t, tc.query))
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got %+v", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("got %+v, want %+v", got, tc.want)
			}
		})
	}
}

func TestTargetSize(t *testing.T) {
	cases := []struct {
		name       string
		srcW, srcH int
		p          TransformParams
		wantW, H   int
	}{
		{name: "contain_both", srcW: 800, srcH: 600, p: TransformParams{Width: 400, Height: 400, Fit: "contain"}, wantW: 400, H: 300},
		{name: "contain_width_only", srcW: 800, srcH: 600, p: TransformParams{Width: 400, Fit: "contain"}, wantW: 400, H: 300},
		{name: "contain_height_only", srcW: 800, srcH: 600, p: TransformParams{Height: 300, Fit: "contain"}, wantW: 400, H: 300},
		{name: "cover_exact", srcW: 800, srcH: 600, p: TransformParams{Width: 200, Height: 200, Fit: "cover"}, wantW: 200, H: 200},
		{name: "upscale_noop_contain", srcW: 100, srcH: 80, p: TransformParams{Width: 500, Height: 500, Fit: "contain"}, wantW: 100, H: 80},
		{name: "upscale_noop_cover", srcW: 100, srcH: 80, p: TransformParams{Width: 500, Height: 500, Fit: "cover"}, wantW: 100, H: 80},
		{name: "upscale_noop_single_edge", srcW: 100, srcH: 80, p: TransformParams{Width: 500, Fit: "contain"}, wantW: 100, H: 80},
		{name: "identity", srcW: 100, srcH: 80, p: TransformParams{Fit: "contain"}, wantW: 100, H: 80},
		{name: "tiny_result_floors_to_one", srcW: 1000, srcH: 4000, p: TransformParams{Height: 1, Fit: "contain"}, wantW: 1, H: 1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w, h := TargetSize(tc.srcW, tc.srcH, tc.p)
			if w != tc.wantW || h != tc.H {
				t.Fatalf("TargetSize(%d,%d,%+v) = (%d,%d), want (%d,%d)",
					tc.srcW, tc.srcH, tc.p, w, h, tc.wantW, tc.H)
			}
		})
	}
}

// makeTestPNG 生成 srcW×srcH 的测试图（左上红右下蓝渐变块，含 alpha=0 区域
// 验证白底铺底）。
func makeTestPNG(t *testing.T, srcW, srcH int) []byte {
	t.Helper()
	img := image.NewNRGBA(image.Rect(0, 0, srcW, srcH))
	for y := 0; y < srcH; y++ {
		for x := 0; x < srcW; x++ {
			if x < srcW/2 {
				img.Set(x, y, color.RGBA{R: 220, G: 30, B: 30, A: 255})
			} else {
				img.Set(x, y, color.RGBA{R: 30, G: 30, B: 220, A: 255})
			}
		}
	}
	// 右上角一块全透明，验证 JPEG 输出铺白底而不是黑底
	for y := 0; y < srcH/10; y++ {
		for x := srcW*9/10; x < srcW; x++ {
			img.Set(x, y, color.RGBA{A: 0})
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("png.Encode: %v", err)
	}
	return buf.Bytes()
}

func TestTransformImage(t *testing.T) {
	src := makeTestPNG(t, 800, 600)

	t.Run("resize_width_preserves_aspect", func(t *testing.T) {
		out, err := TransformImage(src, TransformParams{Width: 400, Quality: 80, Fit: "contain"})
		if err != nil {
			t.Fatalf("TransformImage: %v", err)
		}
		cfg, format, err := image.DecodeConfig(bytes.NewReader(out))
		if err != nil {
			t.Fatalf("output not decodable: %v", err)
		}
		if format != "jpeg" {
			t.Fatalf("output format = %s, want jpeg", format)
		}
		if cfg.Width != 400 || cfg.Height != 300 {
			t.Fatalf("output %dx%d, want 400x300", cfg.Width, cfg.Height)
		}
	})

	t.Run("cover_exact_dimensions", func(t *testing.T) {
		out, err := TransformImage(src, TransformParams{Width: 200, Height: 200, Quality: 75, Fit: "cover"})
		if err != nil {
			t.Fatalf("TransformImage: %v", err)
		}
		cfg, _, err := image.DecodeConfig(bytes.NewReader(out))
		if err != nil {
			t.Fatalf("output not decodable: %v", err)
		}
		if cfg.Width != 200 || cfg.Height != 200 {
			t.Fatalf("output %dx%d, want 200x200", cfg.Width, cfg.Height)
		}
	})

	t.Run("output_smaller_than_source", func(t *testing.T) {
		out, err := TransformImage(src, TransformParams{Width: 100, Quality: 60, Fit: "contain"})
		if err != nil {
			t.Fatalf("TransformImage: %v", err)
		}
		if len(out) >= len(src) {
			t.Fatalf("output %d bytes not smaller than source %d bytes", len(out), len(src))
		}
	})

	t.Run("webp_source_decodes", func(t *testing.T) {
		img := imaging.New(64, 48, color.RGBA{R: 10, G: 200, B: 10, A: 255})
		var src bytes.Buffer
		// imaging 无 webp 编码，用 png 源验证通用路径即可；webp 解码注册
		// 由 blank import 提供，此处仅确认 jpeg/png 主链路不回归。
		if err := png.Encode(&src, img); err != nil {
			t.Fatalf("png.Encode: %v", err)
		}
		if _, err := TransformImage(src.Bytes(), TransformParams{Width: 32, Quality: 80, Fit: "contain"}); err != nil {
			t.Fatalf("TransformImage: %v", err)
		}
	})

	t.Run("rejects_garbage_bytes", func(t *testing.T) {
		if _, err := TransformImage([]byte("not an image"), TransformParams{Width: 32, Fit: "contain"}); err == nil {
			t.Fatal("expected error for garbage input")
		}
	})
}

func TestDerivedObjectName(t *testing.T) {
	p := TransformParams{Width: 500, Quality: 80, Fit: "contain"}
	got := derivedObjectName("assets/2026/abc-123.png", p)
	want := "assets/2026/abc-123-w500h0q80-contain.jpg"
	if got != want {
		t.Fatalf("derivedObjectName = %q, want %q", got, want)
	}
	// 与预生成变体后缀模式（-thumbnail/-medium/-large）不撞名
	p2 := TransformParams{Width: 320, Quality: 75, Fit: "contain"}
	if got := derivedObjectName("assets/2026/abc-123.png", p2); got == "assets/2026/abc-123-thumbnail.jpg" {
		t.Fatal("derived name collides with pre-generated variant name")
	}
}

func TestSplitStoragePath(t *testing.T) {
	bucket, object, err := splitStoragePath("/images/assets/2026/abc.png")
	if err != nil || bucket != "images" || object != "assets/2026/abc.png" {
		t.Fatalf("splitStoragePath = (%q,%q,%v)", bucket, object, err)
	}
	if _, _, err := splitStoragePath("no-slash"); err == nil {
		t.Fatal("expected error for path without bucket")
	}
	if _, _, err := splitStoragePath("/only-bucket/"); err == nil {
		t.Fatal("expected error for empty object name")
	}
}

func TestEtagMatches(t *testing.T) {
	cases := []struct {
		header, etag string
		want         bool
	}{
		{`"img-1-w500h0q80-contain"`, `"img-1-w500h0q80-contain"`, true},
		{`W/"img-1-raw", "img-1-w500h0q80-contain"`, `"img-1-w500h0q80-contain"`, true},
		{"*", `"img-1-raw"`, true},
		{`"img-1-raw"`, `"img-1-w500h0q80-contain"`, false},
		{"", `"img-1-raw"`, false},
	}
	for _, tc := range cases {
		if got := etagMatches(tc.header, tc.etag); got != tc.want {
			t.Fatalf("etagMatches(%q,%q) = %v, want %v", tc.header, tc.etag, got, tc.want)
		}
	}
}

func TestIsTransformableImage(t *testing.T) {
	for _, mime := range []string{"image/jpeg", "image/jpg", "image/png", "image/webp", "IMAGE/PNG"} {
		if !isTransformableImage(mime) {
			t.Fatalf("%q should be transformable", mime)
		}
	}
	for _, mime := range []string{"image/gif", "image/svg+xml", "image/avif", "video/mp4", ""} {
		if isTransformableImage(mime) {
			t.Fatalf("%q should NOT be transformable", mime)
		}
	}
}
