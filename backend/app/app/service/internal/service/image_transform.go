package service

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"net/url"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	_ "golang.org/x/image/webp" // 注册 webp 解码（x/image 无 webp 编码，输出统一 JPEG）
)

// 即时图片变换：/app/v1/media-assets/{id}/content?w=&h=&q=&fit=。
// imgproxy/Cloudinary 的缩小子集，设计约束：
//   - 只做缩小不放大（放大请求原样直出原图），源限 jpeg/png/webp 栅格，
//     gif（丢帧）/svg（XSS 面）/其余 MIME 返回 415，走元数据 url 兜底；
//   - 输出统一 JPEG（有损压缩最通用；Go 纯实现无 webp 编码）；
//   - 参数钳位 + DecodeConfig 先行的解码炸弹防护（像素总量上限），
//     变换结果按参数键落 MinIO 派生对象，命中即免解码直出——
//     同一 (资产, 参数) 组合只在首次请求付一次 CPU；
//   - 参数空间经钳位有界、对象名由服务端拼装不可伪造；公网暴露建议叠
//     接入层限流防缓存填充。

const (
	// transformMaxEdge 单边上限（像素）。4K 屏一屏边长，再大无渲染意义。
	transformMaxEdge = 4096
	// transformMaxSrcPixels 源图像素总量上限（解码炸弹防护，≈8k×5k）。
	transformMaxSrcPixels = 40_000_000
	// defaultImageQuality q 参数缺省值（与预生成变体 medium 档一致）。
	defaultImageQuality = 80
	// derivedImageContentType 变换输出的统一 MIME。
	derivedImageContentType = "image/jpeg"

	fitContain = "contain" // 等比缩放进目标框（默认）
	fitCover   = "cover"   // 等比放大裁切到目标框（需同时给 w 与 h）
)

// TransformParams 即时变换参数（已钳位/归一）。
type TransformParams struct {
	Width   int
	Height  int
	Quality int
	Fit     string
}

// IsIdentity 未要求缩放（w/h 均缺省）——原样直出原图。
func (p TransformParams) IsIdentity() bool { return p.Width <= 0 && p.Height <= 0 }

// cacheKey 派生对象参数键（进对象名，须稳定且文件名安全）。
func (p TransformParams) cacheKey() string {
	return fmt.Sprintf("w%dh%dq%d-%s", p.Width, p.Height, p.Quality, p.Fit)
}

// ParseTransformParams 从查询串解析变换参数：w/h（像素，钳到 [0,4096]）、
// q（JPEG 质量，钳到 [1,100]，缺省 80）、fit（contain|cover，缺省 contain；
// cover 单边时退化为 contain——裁切需要完整目标框）。非法数值返回错误。
func ParseTransformParams(q url.Values) (TransformParams, error) {
	p := TransformParams{Quality: defaultImageQuality, Fit: fitContain}

	parseEdge := func(key string) (int, error) {
		v := q.Get(key)
		if v == "" {
			return 0, nil
		}
		n, err := strconv.Atoi(v)
		if err != nil || n < 0 {
			return 0, fmt.Errorf("invalid %s %q", key, v)
		}
		if n > transformMaxEdge {
			n = transformMaxEdge
		}
		return n, nil
	}
	var err error
	if p.Width, err = parseEdge("w"); err != nil {
		return TransformParams{}, err
	}
	if p.Height, err = parseEdge("h"); err != nil {
		return TransformParams{}, err
	}

	if v := q.Get("q"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return TransformParams{}, fmt.Errorf("invalid q %q", v)
		}
		switch {
		case n < 1:
			n = 1
		case n > 100:
			n = 100
		}
		p.Quality = n
	}

	switch f := q.Get("fit"); f {
	case "":
	case fitContain, fitCover:
		p.Fit = f
	default:
		return TransformParams{}, fmt.Errorf("invalid fit %q", f)
	}
	if p.Fit == fitCover && (p.Width <= 0 || p.Height <= 0) {
		p.Fit = fitContain
	}
	return p, nil
}

// TargetSize 计算输出尺寸（纯函数）。只缩小：目标不小于源时返回源尺寸，
// 调用方据此跳过解码直出原图。
func TargetSize(srcW, srcH int, p TransformParams) (int, int) {
	if srcW <= 0 || srcH <= 0 {
		return srcW, srcH
	}
	switch {
	case p.Width > 0 && p.Height > 0:
		if p.Fit == fitCover {
			if p.Width >= srcW && p.Height >= srcH {
				return srcW, srcH
			}
			return p.Width, p.Height
		}
		scale := min(float64(p.Width)/float64(srcW), float64(p.Height)/float64(srcH))
		if scale >= 1 {
			return srcW, srcH
		}
		return roundPositive(float64(srcW)*scale), roundPositive(float64(srcH)*scale)
	case p.Width > 0:
		if p.Width >= srcW {
			return srcW, srcH
		}
		return p.Width, roundPositive(float64(srcH) * float64(p.Width) / float64(srcW))
	case p.Height > 0:
		if p.Height >= srcH {
			return srcW, srcH
		}
		return roundPositive(float64(srcW) * float64(p.Height) / float64(srcH)), p.Height
	default:
		return srcW, srcH
	}
}

func roundPositive(v float64) int {
	if v < 1 {
		return 1
	}
	return int(v + 0.5)
}

// TransformImage 变换图片字节为 JPEG（调用方已保证需要缩放：
// TargetSize 结果与源尺寸不同的场合）。含两道防护：
// DecodeConfig 先行拒绝超像素源（全量解码前），参数已在入口钳位。
func TransformImage(src []byte, p TransformParams) ([]byte, error) {
	cfg, _, err := image.DecodeConfig(bytes.NewReader(src))
	if err != nil {
		return nil, fmt.Errorf("probe image: %w", err)
	}
	if cfg.Width <= 0 || cfg.Height <= 0 {
		return nil, fmt.Errorf("bad image dimensions %dx%d", cfg.Width, cfg.Height)
	}
	if cfg.Width*cfg.Height > transformMaxSrcPixels {
		return nil, fmt.Errorf("source image too large: %dx%d", cfg.Width, cfg.Height)
	}

	img, err := imaging.Decode(bytes.NewReader(src), imaging.AutoOrientation(true))
	if err != nil {
		return nil, fmt.Errorf("decode image: %w", err)
	}

	tw, th := TargetSize(img.Bounds().Dx(), img.Bounds().Dy(), p)
	var out image.Image
	if p.Fit == fitCover {
		out = imaging.Fill(img, tw, th, imaging.Center, imaging.Lanczos)
	} else {
		out = imaging.Resize(img, tw, th, imaging.Lanczos)
	}
	return encodeJPEG(flattenWhite(out), p.Quality)
}

// flattenWhite JPEG 无 alpha 通道：透明像素铺白底，避免变黑。
func flattenWhite(img image.Image) image.Image {
	b := img.Bounds()
	dst := image.NewRGBA(b)
	draw.Draw(dst, b, image.NewUniform(color.White), image.Point{}, draw.Src)
	draw.Draw(dst, b, img, b.Min, draw.Over)
	return dst
}

func encodeJPEG(img image.Image, quality int) ([]byte, error) {
	var buf bytes.Buffer
	if err := imaging.Encode(&buf, img, imaging.JPEG, imaging.JPEGQuality(quality)); err != nil {
		return nil, fmt.Errorf("encode jpeg: %w", err)
	}
	return buf.Bytes(), nil
}

// isTransformableImage 可即时变换的源 MIME（gif 动图丢帧、svg 是 XSS 面，
// 均不支持——请求方走元数据 url 直出）。
func isTransformableImage(mimeType string) bool {
	switch strings.ToLower(strings.TrimSpace(mimeType)) {
	case "image/jpeg", "image/jpg", "image/png", "image/webp":
		return true
	}
	return false
}

// derivedObjectName 派生缓存对象名：原对象名去扩展名 + 参数键 + .jpg，
// 与预生成变体（-thumbnail/-medium/-large）后缀模式不同，不会撞名。
func derivedObjectName(objectName string, p TransformParams) string {
	base := objectName
	if dot := strings.LastIndex(base, "."); dot > 0 {
		base = base[:dot]
	}
	return fmt.Sprintf("%s-%s.jpg", base, p.cacheKey())
}

// splitStoragePath 媒体资产 storagePath 形如 "/{bucket}/{objectName}"
// （JoinObjectUrl 空 endpoint 的产物），拆出桶与对象名。
func splitStoragePath(storagePath string) (bucket, objectName string, err error) {
	trimmed := strings.TrimPrefix(storagePath, "/")
	slash := strings.Index(trimmed, "/")
	if slash <= 0 || slash == len(trimmed)-1 {
		return "", "", fmt.Errorf("bad storage path %q", storagePath)
	}
	return trimmed[:slash], trimmed[slash+1:], nil
}
