package social

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
)

// UserAgent 各平台适配器统一的请求标识。
const UserAgent = "go-wind-cms-social/1.0"

// NewMultipartWriter 创建 multipart/form-data 缓冲与写入器（调用方负责 Close）。
func NewMultipartWriter() (*bytes.Buffer, *multipart.Writer) {
	buf := &bytes.Buffer{}
	return buf, multipart.NewWriter(buf)
}

// NewMultipartRequest 组装 multipart 请求（User-Agent + Content-Type）。
func NewMultipartRequest(ctx context.Context, method, url string, buf *bytes.Buffer, fw *multipart.Writer) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(buf.Bytes()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Content-Type", fw.FormDataContentType())
	return req, nil
}

// WriteFilePart 向 multipart 写入一个文件字段。
func WriteFilePart(fw *multipart.Writer, fieldName, filename string, data []byte) error {
	part, err := fw.CreateFormFile(fieldName, filename)
	if err != nil {
		return err
	}
	_, err = part.Write(data)
	return err
}
