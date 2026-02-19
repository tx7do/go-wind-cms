package oss

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"mime"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type GenerateFileNameType string

const (
	GenerateFileNameTypeUUID          GenerateFileNameType = "uuid"
	GenerateFileNameTypeContentSHA256 GenerateFileNameType = "content_sha256"
	GenerateFileNameTypeHMACContent   GenerateFileNameType = "hmac_content"
	GenerateFileNameTypeTimeBase      GenerateFileNameType = "time_base"
)

const (
	BucketImages = "images"
	BucketVideos = "videos"
	BucketAudios = "audios"
	BucketDocs   = "docs"
	BucketFiles  = "files"
)

var staticHMACSecret = []byte("0123456789abcdef0123456789abcdef") // 32 bytes secret for HMAC

// ContentTypeToBucketName 根据文件类型获取存储桶名称
func ContentTypeToBucketName(contentType string) string {
	if strings.TrimSpace(contentType) == "" {
		return BucketFiles
	}

	// 解析 media type，忽略参数（如 charset）
	mt, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		mt = strings.ToLower(strings.TrimSpace(contentType))
	} else {
		mt = strings.ToLower(mt)
	}

	parts := strings.SplitN(mt, "/", 2)
	if len(parts) != 2 {
		return BucketFiles
	}

	main := parts[0]
	sub := parts[1]

	switch main {
	case "image":
		return BucketImages
	case "video":
		return BucketVideos
	case "audio":
		return BucketAudios
	case "text":
		return BucketDocs
	case "application":
		// 常见文档/办公类型映射到 docs，其余为 files
		switch sub {
		case "pdf", "json":
			return BucketDocs
		default:
			if strings.HasPrefix(sub, "vnd.ms-") ||
				strings.Contains(sub, "officedocument") ||
				strings.Contains(sub, "word") ||
				strings.Contains(sub, "excel") ||
				strings.Contains(sub, "powerpoint") {
				return BucketDocs
			}
			return BucketFiles
		}
	default:
		return BucketFiles
	}
}

// FileExtensionToBucketName 根据文件后缀获取存储桶名称
func FileExtensionToBucketName(ext string) string {
	if strings.TrimSpace(ext) == "" {
		return BucketFiles
	}

	s := strings.TrimSpace(ext)

	// 如果看起来像 MIME 类型，复用 ContentTypeToBucketName
	if strings.Contains(s, "/") {
		return ContentTypeToBucketName(s)
	}

	e := strings.ToLower(strings.TrimPrefix(s, "."))

	switch e {
	// images
	case "jpg", "jpeg", "png", "gif", "webp", "bmp", "ico", "svg", "tif", "tiff", "heic":
		return BucketImages
	// videos
	case "mp4", "webm", "mov", "mkv", "avi", "flv", "mpeg", "mpg":
		return BucketVideos
	// audios
	case "mp3", "wav", "ogg", "m4a", "flac", "aac":
		return BucketAudios
	// text / docs
	case "txt", "html", "htm", "css", "js", "csv", "md", "xml", "json":
		return BucketDocs
	// common office / documents
	case "pdf", "doc", "docx", "xls", "xlsx", "ppt", "pptx":
		return BucketDocs
	// archives / binaries -> 默认 files（可按需改为 archives）
	case "zip", "tar", "gz", "tgz", "7z", "rar", "bz2", "bin", "exe":
		return BucketFiles
	default:
		return BucketFiles
	}
}

// ContentTypeToFileExtension 根据文件类型获取文件后缀
func ContentTypeToFileExtension(contentType string) string {
	if strings.TrimSpace(contentType) == "" {
		return ""
	}

	// 忽略参数（如 charset）并转为小写
	mt, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		mt = strings.ToLower(strings.TrimSpace(contentType))
	} else {
		mt = strings.ToLower(mt)
	}

	// 常见类型映射
	switch mt {
	// images
	case "image/jpeg", "image/jpg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	case "image/bmp":
		return ".bmp"
	case "image/x-icon", "image/vnd.microsoft.icon":
		return ".ico"
	case "image/svg+xml":
		return ".svg"

	// video
	case "video/mp4":
		return ".mp4"
	case "video/webm":
		return ".webm"
	case "video/quicktime":
		return ".mov"
	case "video/x-matroska", "video/mkv":
		return ".mkv"

	// audio
	case "audio/mpeg":
		return ".mp3"
	case "audio/wav", "audio/x-wav":
		return ".wav"
	case "audio/ogg", "audio/vorbis":
		return ".ogg"
	case "audio/mp4":
		return ".m4a"

	// text
	case "text/plain":
		return ".txt"
	case "text/html":
		return ".html"
	case "text/css":
		return ".css"
	case "text/csv":
		return ".csv"
	case "text/xml":
		return ".xml"

	// JavaScript
	case "text/javascript", "application/javascript", "application/x-javascript":
		return ".js"

	// Lua
	case "text/x-lua", "application/x-lua":
		return ".lua"

	// Python
	case "text/x-python", "application/x-python", "text/python":
		return ".py"

	// Shell scripts
	case "text/x-shellscript", "application/x-sh", "application/x-shellscript", "text/x-sh", "text/x-bash", "application/x-bash":
		return ".sh"

	// application / documents / archives
	case "application/pdf":
		return ".pdf"
	case "application/json":
		return ".json"
	case "application/zip":
		return ".zip"
	case "application/x-tar":
		return ".tar"
	case "application/gzip", "application/x-gzip":
		return ".gz"
	case "application/x-7z-compressed", "application/7z":
		return ".7z"
	case "application/msword":
		return ".doc"
	case "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
		return ".docx"
	case "application/vnd.ms-excel":
		return ".xls"
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		return ".xlsx"
	case "application/vnd.ms-powerpoint":
		return ".ppt"
	case "application/vnd.openxmlformats-officedocument.presentationml.presentation":
		return ".pptx"
	case "application/octet-stream":
		return ".bin"
	}

	// 回退：使用标准库尝试获取扩展名
	exts, _ := mime.ExtensionsByType(mt)
	if len(exts) > 0 {
		ext := exts[0]
		if !strings.HasPrefix(ext, ".") {
			return "." + ext
		}
		return ext
	}

	// 未知类型，返回空串（调用方可据此决定）
	return ""
}

// DetectFileType 根据文件内容推断 MIME 类型和扩展名（带前导点）
// 输入：文件内容的字节片（可为整个文件或前若干字节）
// 输出：mimeType（例如 "image/png"）和 ext（例如 ".png"，若未知则为空字符串）
func DetectFileType(fileContent []byte) (mimeType, ext string) {
	// 使用前 512 字节进行 http.DetectContentType 判断
	head := fileContent
	if len(head) > 512 {
		head = head[:512]
	}
	mimeType = http.DetectContentType(head)

	// 常见魔术头优先判定
	switch {
	case len(fileContent) >= 8 && bytes.Equal(fileContent[:8], []byte("\x89PNG\r\n\x1a\n")):
		return "image/png", ".png"
	case len(fileContent) >= 3 && bytes.Equal(fileContent[:3], []byte{0xff, 0xd8, 0xff}):
		// JPG 可以以 FF D8 FF 开头
		return "image/jpeg", ".jpg"
	case len(fileContent) >= 6 && (bytes.Equal(fileContent[:6], []byte("GIF87a")) || bytes.Equal(fileContent[:6], []byte("GIF89a"))):
		return "image/gif", ".gif"
	case len(fileContent) >= 5 && bytes.Equal(fileContent[:5], []byte("%PDF-")):
		return "application/pdf", ".pdf"
	case len(fileContent) >= 4 && bytes.Equal(fileContent[:4], []byte("PK\x03\x04")):
		// zip / docx / xlsx / jar 等基于 zip 的格式，默认返回 .zip
		return "application/zip", ".zip"
	case len(fileContent) >= 12 && bytes.Equal(fileContent[4:8], []byte("ftyp")):
		// ftyp box 常见于 MP4/IS0 media files；更细分需要读取 brand 字段
		return "video/mp4", ".mp4"
	case len(fileContent) >= 3 && bytes.Equal(fileContent[:3], []byte("ID3")):
		return "audio/mpeg", ".mp3"
	case len(fileContent) >= 2 && fileContent[0] == 0xFF && (fileContent[1]&0xE0) == 0xE0:
		// 另一种 MP3 帧头判定（0xFF Ex）
		return "audio/mpeg", ".mp3"
	case len(fileContent) >= 12 && bytes.Equal(fileContent[:4], []byte("RIFF")) && bytes.Equal(fileContent[8:12], []byte("WAVE")):
		return "audio/wav", ".wav"
	case len(fileContent) >= 2 && bytes.Equal(fileContent[:2], []byte("BM")):
		return "image/bmp", ".bmp"
	}

	// 若魔术头未命中，尝试从 http.DetectContentType 的 MIME 类型获取扩展名
	if mimeType != "" {
		exts, _ := mime.ExtensionsByType(mimeType)
		if len(exts) > 0 {
			ext = exts[0]
			// mime.ExtensionsByType 可能返回多种扩展，取第一个
			return mimeType, ext
		}
	}

	// 额外的简单映射（在 mime 映射不可用时）
	switch {
	case strings.HasPrefix(mimeType, "image/"):
		// 常见 image 类型再尝试映射
		if mimeType == "image/svg+xml" {
			return mimeType, ".svg"
		}
	case strings.HasPrefix(mimeType, "video/"):
		if mimeType == "video/mp4" {
			return mimeType, ".mp4"
		}
	case strings.HasPrefix(mimeType, "audio/"):
		if mimeType == "audio/mpeg" {
			return mimeType, ".mp3"
		}
	}

	// 未知扩展，返回检测到的 mimeType（ext 为空）
	return mimeType, ""
}

// GeneraUUIDFileName 生成基于 UUID 的文件名
// - 描述：生成去掉横线的 UUID（128-bit 随机/唯一）。
// - 优点：几乎无碰撞、不依赖内容、生成快、安全性高（不泄露内容指纹）。
// - 缺点：名字较长、不可按时间排序、不能用于内容去重。
// - 适用：高并发、无需去重、希望最小碰撞风险的通用场景。
func GeneraUUIDFileName(fileExt string) string {
	// 生成 UUID 并去除横线
	name := strings.ReplaceAll(uuid.New().String(), "-", "")

	// 确保扩展名没有开头的点
	cleanExt := strings.TrimPrefix(fileExt, ".")
	if cleanExt == "" {
		return name
	}

	return fmt.Sprintf("%s.%s", name, cleanExt)
}

// GeneraContentSHA265FileName 生成基于内容 SHA256 的文件名
// - 描述：内容的 SHA‑256 十六进制摘要（明文哈希）。
// - 优点：确定性、极低碰撞、便于去重和校验。
// - 缺点：容易被离线匹配（泄露内容指纹）；需要读取并哈希整个文件（IO/CPU 开销）。
// - 适用：去重、缓存命中、内容可验证，但不适合对隐私敏感的公开键。
func GeneraContentSHA265FileName(fileContent []byte, fileExt string) string {
	// 计算内容的 SHA256 摘要并编码为十六进制字符串
	sum := sha256.Sum256(fileContent)
	hash := hex.EncodeToString(sum[:])

	// 清理扩展名中的前导点
	cleanExt := strings.TrimPrefix(fileExt, ".")
	if cleanExt == "" {
		return hash
	}

	return fmt.Sprintf("%s.%s", hash, cleanExt)
}

// GenerateHMACContentFileName 生成基于 HMAC 内容的文件名
// - 描述：先对内容做 SHA‑256，再用 secret 做 HMAC 生成文件名。
// - 优点：对同样内容一致（可去重），但在没有 secret 的情况下不可被离线匹配；抗探测性好（用于隐藏内容指纹）。
// - 缺点：需要管理 secret（泄露则失效）；计算开销中等（一次哈希 + HMAC）。
// - 适用：需要去重同时又要防止哈希被猜测/探测的场景。
func GenerateHMACContentFileName(secret []byte, fileContent []byte, fileExt string) string {
	// 先计算内容 SHA256
	sum := sha256.Sum256(fileContent)

	// 用 secret 做 HMAC
	mac := hmac.New(sha256.New, secret)
	mac.Write(sum[:])
	h := mac.Sum(nil)
	name := hex.EncodeToString(h)

	cleanExt := strings.TrimPrefix(fileExt, ".")
	if cleanExt == "" {
		return name
	}
	return fmt.Sprintf("%s.%s", name, cleanExt)
}

// GeneraTimeBaseFileName 生成基于时间戳的文件名
// - 描述：时间戳（秒级）+ 短随机后缀（当前实现 8 hex）。
// - 优点：可读、按时间排序、生成成本低。
// - 缺点：同一秒高并发下碰撞概率上升（8 hex ≈ 32-bit 熵）；不可用于内容去重；暴露时间信息。
// - 适用：审计/按时间查找、并发不极端或改用 UnixNano/更长后缀/ULID。
// - 改进建议：改用 UnixNano 或 增大随机后缀到 16 hex/全 UUID，或使用 ULID（时间有序且高熵）。
func GeneraTimeBaseFileName(fileExt string) string {
	// 清理扩展名中的前导点
	cleanExt := strings.TrimPrefix(fileExt, ".")

	// 使用 UUID 去掉横线后取前 8 字符作为随机后缀
	suffix := strings.ReplaceAll(uuid.New().String(), "-", "")
	if len(suffix) > 8 {
		suffix = suffix[:8]
	}

	timestamp := time.Now().UnixNano()

	if cleanExt == "" {
		return fmt.Sprintf("%d_%s", timestamp, suffix)
	}

	return fmt.Sprintf("%d_%s.%s", timestamp, suffix, cleanExt)
}

// GenerateFileName 生成文件名
func GenerateFileName(fileContent []byte, fileExt string, typ GenerateFileNameType) string {
	switch typ {
	case GenerateFileNameTypeUUID:
		// 要最大唯一性、低延迟、高并发
		return GeneraUUIDFileName(fileExt)
	case GenerateFileNameTypeContentSHA256:
		// 要去重/验证且不关心能被离线匹配
		return GeneraContentSHA265FileName(fileContent, fileExt)
	case GenerateFileNameTypeHMACContent:
		// 要去重但需防止被猜测
		return GenerateHMACContentFileName(staticHMACSecret, fileContent, fileExt)
	case GenerateFileNameTypeTimeBase:
		// 要可读且按时间排序
		return GeneraTimeBaseFileName(fileExt)
	default:
		return GeneraUUIDFileName(fileExt)
	}
}

// GenerateObjectName 生成对象名
func GenerateObjectName(fileDirectory string, fileContent []byte, fileExt string, typ GenerateFileNameType) string {
	// 生成文件名
	name := GenerateFileName(fileContent, fileExt, typ)

	// 清理首尾斜杠，避免出现双斜杠或以斜杠开头的对象名
	dir := strings.Trim(fileDirectory, "/")

	if dir == "" {
		return name
	}
	return dir + "/" + name
}

func EnsureObjectName(fileDirectory, sourceFileName, contentType string, fileContent []byte, typ GenerateFileNameType) string {
	if typ == "" {
		typ = GenerateFileNameTypeUUID
	}
	fileExt := EnsureFileExtension(sourceFileName, contentType, fileContent)
	return GenerateObjectName(fileDirectory, fileContent, fileExt, typ)
}

// JoinObjectName 拼接对象名
func JoinObjectName(contentType string, filePath, fileName *string) (string, string) {
	fileSuffix := ContentTypeToFileExtension(contentType)

	var _fileName string
	if fileName == nil {
		_fileName = uuid.New().String() + fileSuffix
	} else {
		_fileName = *fileName
	}

	var objectName string
	if filePath != nil {
		objectName = *filePath + "/" + _fileName
	} else {
		objectName = _fileName
	}

	return objectName, _fileName
}

// JoinObjectUrl 拼接对象 URL
func JoinObjectUrl(endpoint, bucketName, objectName string) string {
	trimmedEndpoint := strings.TrimRight(endpoint, "/")
	trimmedBucket := strings.Trim(bucketName, "/")
	trimmedObject := strings.TrimLeft(objectName, "/")

	return fmt.Sprintf("%s/%s/%s", trimmedEndpoint, trimmedBucket, trimmedObject)
}

// ReplaceEndpointHost 替换URL中的Host部分
func ReplaceEndpointHost(rawURL, host, endpoint string) string {
	if rawURL == "" || host == "" {
		return rawURL
	}
	return strings.Replace(rawURL, endpoint, host, -1)
}

func ExtractFileExtension(fileName string) string {
	// 查找最后一个点作为扩展名分隔（点在开头不算）
	idx := strings.LastIndex(fileName, ".")
	if idx <= 0 {
		// 无扩展名或点在首位
		return ""
	}

	ext := strings.ToLower(fileName[idx+1:])
	return ext
}

// EnsureFileExtension 确保文件后缀存在，按顺序从文件名、内容类型、文件内容检测
func EnsureFileExtension(fileName, contentType string, fileContent []byte) string {
	ext := ExtractFileExtension(fileName)
	if ext != "" {
		return ext
	}

	ext = ContentTypeToFileExtension(contentType)
	if ext != "" {
		return ext
	}

	_, ext = DetectFileType(fileContent)
	if ext != "" {
		return strings.TrimPrefix(ext, ".")
	}

	return "bin"
}

// SetDownloadRange 设置下载范围
func SetDownloadRange(opts *minio.GetObjectOptions, start, end *int64) {
	if opts == nil {
		return
	}

	if start != nil && end != nil {
		_ = opts.SetRange(*start, *end)
	} else if start != nil {
		_ = opts.SetRange(*start, 0)
	} else if end != nil {
		_ = opts.SetRange(0, *end)
	}
}
