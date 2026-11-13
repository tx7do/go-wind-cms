package oss

import (
	"fmt"
	"mime"
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
func GenerateFileName(fileExt string, typ GenerateFileNameType) string {
	switch typ {
	case GenerateFileNameTypeTimeBase:
		// 要可读且按时间排序
		return GeneraTimeBaseFileName(fileExt)
	default:
		// 默认 UUID：最大唯一性、低延迟、高并发，不依赖文件内容
		return GeneraUUIDFileName(fileExt)
	}
}

// GenerateObjectName 生成对象名
func GenerateObjectName(fileDirectory string, fileExt string, typ GenerateFileNameType) string {
	// 生成文件名
	name := GenerateFileName(fileExt, typ)

	// 清理首尾斜杠，避免出现双斜杠或以斜杠开头的对象名
	dir := strings.Trim(fileDirectory, "/")

	if dir == "" {
		return name
	}
	return dir + "/" + name
}

func EnsureObjectName(fileDirectory, sourceFileName, contentType string, typ GenerateFileNameType) string {
	if typ == "" {
		typ = GenerateFileNameTypeUUID
	}
	fileExt := EnsureFileExtension(sourceFileName, contentType)
	return GenerateObjectName(fileDirectory, fileExt, typ)
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

// EnsureFileExtension 确保文件后缀存在，按顺序从文件名、内容类型推断。
// 流式上传下不再读取文件内容做嗅探，取不到时回退到 "bin"。
func EnsureFileExtension(fileName, contentType string) string {
	ext := ExtractFileExtension(fileName)
	if ext != "" {
		return ext
	}

	ext = ContentTypeToFileExtension(contentType)
	if ext != "" {
		return ext
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
