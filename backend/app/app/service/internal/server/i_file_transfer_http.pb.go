package server

import (
	"mime/multipart"
	"bytes"
	"context"
	"io"
	"strconv"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/tx7do/go-utils/trans"

	"github.com/go-kratos/kratos/v2/encoding"
	_ "github.com/go-kratos/kratos/v2/encoding/json"

	"go-wind-cms/app/app/service/internal/service"
	oss "go-wind-cms/pkg/oss"

	storageV1 "go-wind-cms/api/gen/go/storage/service/v1"
)

var codec = encoding.GetCodec("json")

func registerFileTransferServiceHandler(srv *http.Server, svc *service.FileTransferService) {
	r := srv.Route("/")

	// 上传路由在 HTTP Filter 层单独封顶，而不是放在 kratos 中间件链里：
	// ctx.Middleware 由 handler 自己调用，中间件实际在请求体解析之后才执行，
	// 无法阻止 body 先被整包读入。
	uploadGuard := uploadBodyGuard()

	r.POST("app/v1/file/upload", _FileTransferService_PostUploadFile_HTTP_Handler(svc), uploadGuard)
	r.PUT("app/v1/file/upload", _FileTransferService_PutUploadFile_HTTP_Handler(svc), uploadGuard)

	r.GET("app/v1/file/download", _FileTransferService_DownloadFile_HTTP_Handler(svc))

}

const OperationFileTransferServicePostUploadFile = "/app.service.v1.FileTransferService/PostUploadFile"
const OperationFileTransferServicePutUploadFile = "/app.service.v1.FileTransferService/PutUploadFile"

const OperationFileTransferServiceDownloadFile = "/app.service.v1.FileTransferService/DownloadFile"

func _FileTransferService_PostUploadFile_HTTP_Handler(svc *service.FileTransferService) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		http.SetOperation(ctx, OperationFileTransferServicePostUploadFile)
		return serveUploadFile(ctx, svc)
	}
}

func _FileTransferService_PutUploadFile_HTTP_Handler(svc *service.FileTransferService) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		http.SetOperation(ctx, OperationFileTransferServicePutUploadFile)
		return serveUploadFile(ctx, svc)
	}
}

// serveUploadFile POST/PUT 上传的共用实现。
// 关键点：请求体解析放在 ctx.Middleware 返回的 handler 内部。kratos 的中间件链
// 是在 h(ctx, &in) 被调用时才执行的，若在它之前解析 body（FormFile 会把整个
// multipart 读出），认证/鉴权都会被推迟到整包读完之后（未认证即可耗尽资源）。
func serveUploadFile(ctx http.Context, svc *service.FileTransferService) error {
	var in storageV1.UploadFileRequest

	h := ctx.Middleware(func(mctx context.Context, req interface{}) (interface{}, error) {
		aReq := req.(*storageV1.UploadFileRequest)

		file, header, err := ctx.Request().FormFile("file")
		if err != nil {
			if isMaxBytesError(err) {
				return nil, storageV1.ErrorFileTooLarge("request body exceeds the upload body limit")
			}

			// 非 multipart 请求沿用原有行为：交给框架按 Content-Type 解码。
			if bindErr := ctx.Bind(aReq); bindErr != nil {
				log.Errorf("Bind upload file request error: %v", bindErr)
				return nil, bindErr
			}
		} else {
			defer file.Close()

			if merr := applyUploadFileFormMeta(ctx, aReq, header); merr != nil {
				return nil, merr
			}

			// 流式上传：multipart.File 作为 io.Reader 通过 context 注入，
			// service 侧取出后直接喂给 minio-go PutObject（内部按 size 自动
			// 分片），避免整文件载入 []byte 导致 OOM。
			mctx = oss.WithUploadReader(mctx, file, header.Size)
		}

		resp, err := svc.UploadFile(mctx, aReq)
		return resp, err
	})

	// 逻辑处理，取数据
	out, err := h(ctx, &in)
	if err != nil {
		return err
	}

	reply, ok := out.(*storageV1.UploadFileResponse)
	if !ok {
		return storageV1.ErrorUploadFailed("unexpected response type %T", out)
	}

	return ctx.Result(200, reply)
}

// applyUploadFileFormMeta 读取 multipart 中与文件同行的表单元数据。
func applyUploadFileFormMeta(ctx http.Context, in *storageV1.UploadFileRequest, header *multipart.FileHeader) error {
	sourceFileName := ctx.Request().FormValue("sourceFileName")
	if sourceFileName != "" {
		in.SourceFileName = trans.Ptr(sourceFileName)
	} else {
		in.SourceFileName = trans.Ptr(header.Filename)
	}

	mime := ctx.Request().FormValue("mime")
	if mime != "" {
		in.Mime = trans.Ptr(mime)
	} else {
		in.Mime = trans.Ptr(header.Header.Get("Content-Type"))
	}

	size := ctx.Request().FormValue("size")
	if size != "" {
		if n, parseErr := strconv.ParseInt(size, 10, 64); parseErr == nil {
			in.Size = trans.Ptr(n)
		}
	} else {
		in.Size = trans.Ptr(header.Size)
	}

	storageObject := ctx.Request().FormValue("storageObject")
	if storageObject != "" {
		in.StorageObject = &storageV1.StorageObject{}
		if err := codec.Unmarshal([]byte(storageObject), in.StorageObject); err != nil {
			log.Errorf("Unmarshal upload file storageObject error: %v", err)
			return err
		}
	}

	return nil
}




func _FileTransferService_DownloadFile_HTTP_Handler(svc *service.FileTransferService) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		http.SetOperation(ctx, OperationFileTransferServiceDownloadFile)

		var in storageV1.DownloadFileRequest
		var err error

		if err = ctx.BindQuery(&in); err != nil {
			return err
		}

		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			aReq := req.(*storageV1.DownloadFileRequest)
			var resp *storageV1.DownloadFileResponse
			resp, err = svc.DownloadFile(ctx, aReq)
			return resp, err
		})

		// 逻辑处理，取数据
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}

		reply := out.(*storageV1.DownloadFileResponse)
		rw := ctx.Response()
		if rw == nil {
			return ctx.Result(500, "response writer not available")
		}

		data := reply.GetFile()
		if len(data) == 0 {
			// 若没有文件字节，交由框架默认处理（保持原行为）
			return ctx.Result(200, reply)
		}

		// 基本头部
		mime := reply.GetMime()
		if mime == "" {
			mime = "application/octet-stream"
		}
		rw.Header().Set("Content-Type", mime)

		filename := reply.GetSourceFileName()
		if filename == "" {
			filename = "file"
		}

		var disposition string
		if in.GetDisposition() != "" {
			disposition = in.GetDisposition()
		} else {
			disposition = "attachment; filename=\"" + filename + "\""
		}
		rw.Header().Set("Content-Disposition", disposition)
		rw.Header().Set("Accept-Ranges", "bytes")

		// 使用 bytes.Reader 以便支持高效的部分读取/流式写入
		reader := bytes.NewReader(data)
		total := int64(len(data))

		// 检查请求是否带 Range 头（仅支持单区间，格式 bytes=start-end 或 bytes=start-）
		rangeHeader := ctx.Request().Header.Get("Range")
		if strings.HasPrefix(rangeHeader, "bytes=") {
			r := strings.TrimPrefix(rangeHeader, "bytes=")
			parts := strings.SplitN(r, "-", 2)
			if len(parts) != 2 {
				return ctx.Result(400, "invalid Range header")
			}

			start, err1 := strconv.ParseInt(parts[0], 10, 64)
			if err1 != nil || start < 0 {
				return ctx.Result(400, "invalid Range start")
			}

			var end int64
			if parts[1] == "" {
				end = total - 1
			} else {
				end, err1 = strconv.ParseInt(parts[1], 10, 64)
				if err1 != nil || end < start {
					return ctx.Result(400, "invalid Range end")
				}
			}

			if start >= total {
				return ctx.Result(416, "requested range not satisfiable")
			}
			if end >= total {
				end = total - 1
			}

			length := end - start + 1
			// 设置部分响应头
			rw.Header().Set("Content-Range", "bytes "+strconv.FormatInt(start, 10)+"-"+strconv.FormatInt(end, 10)+"/"+strconv.FormatInt(total, 10))
			rw.Header().Set("Content-Length", strconv.FormatInt(length, 10))

			// 返回 206 并写入指定区间
			rw.WriteHeader(206)
			if _, err = reader.Seek(start, io.SeekStart); err != nil {
				return ctx.Result(500, err.Error())
			}
			if _, err = io.CopyN(rw, reader, length); err != nil {
				return ctx.Result(500, err.Error())
			}
			return nil
		}

		// 无 Range，返回完整内容（200）
		rw.Header().Set("Content-Length", strconv.FormatInt(total, 10))
		rw.WriteHeader(200)
		if _, err = io.Copy(rw, reader); err != nil {
			return ctx.Result(500, err.Error())
		}
		return nil
	}
}
