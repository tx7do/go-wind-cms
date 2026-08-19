package oss

import (
	"context"
	"io"
)

// uploadReaderKey 用于在 HTTP handler 与 service 之间传递上传文件的流式 reader，
// 避免 proto 的 []byte 字段把整文件载入内存。
type uploadReaderCtxKey struct{}

// WithUploadReader 将上传文件的 io.Reader 与其声明的字节大小注入 context。
// handler 解析 multipart 后调用，service 侧通过 UploadReaderFromContext 取回。
func WithUploadReader(ctx context.Context, reader io.Reader, size int64) context.Context {
	return context.WithValue(ctx, uploadReaderCtxKey{}, uploadReaderPayload{reader: reader, size: size})
}

// UploadReaderFromContext 取出注入的上传 reader 与大小。
// 若未注入返回 (nil, 0)。
func UploadReaderFromContext(ctx context.Context) (io.Reader, int64) {
	v, _ := ctx.Value(uploadReaderCtxKey{}).(uploadReaderPayload)
	return v.reader, v.size
}

type uploadReaderPayload struct {
	reader io.Reader
	size   int64
}
