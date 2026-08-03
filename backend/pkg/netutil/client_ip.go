package netutil

import (
	"context"
	"net/http"

	"github.com/go-kratos/kratos/v2/transport"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
)

// HeaderFromContext 从 kratos context 中提取 HTTP 请求头。
// 用于 service 层读取通过 header 传递的参数（如验证码 id/value）。
// 若上下文非 HTTP 传输或取不到 request，返回 nil。
func HeaderFromContext(ctx context.Context) http.Header {
	tr, ok := transport.FromServerContext(ctx)
	if !ok {
		return nil
	}
	htr, ok := tr.(khttp.Transporter)
	if !ok {
		return nil
	}
	req := htr.Request()
	if req == nil {
		return nil
	}
	return req.Header
}
