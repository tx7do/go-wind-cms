package server

import (
	"errors"
	kratoshttp "github.com/go-kratos/kratos/v2/transport/http"
	nethttp "net/http"
	"strings"

	auth "go-wind-cms/pkg/middleware/auth"
)

const (
	// maxUploadFileBytes 单文件上传上限。
	maxUploadFileBytes int64 = 50 << 20

	// maxUploadBody multipart 上传路由的请求体上限：单文件上限 + multipart 编码余量。
	maxUploadBody = maxUploadFileBytes + 1<<20

	// maxRequestBody 全站请求体上限。需容纳 base64 编码后的上传体（膨胀 4/3）
	// 与 JSON 结构开销；比它更小的请求不受影响。
	maxRequestBody = maxUploadFileBytes*4/3 + 4<<20

	// maxUnauthenticatedBody 未认证请求体上限。免鉴权端点（登录、验证码、
	// 找回密码、SSO、令牌交换等）的请求体都远小于此值；kratos 生成的 handler
	// 在中间件之前 ctx.Bind（内部 io.ReadAll 整个 body），不设限时未认证请求
	// 即可用超大 body 打爆内存。
	maxUnauthenticatedBody = 1 << 20
)

// uploadBodyGuard 在 handler 之前给上传请求体封顶。
// multipart 的 body 约等于文件本身；JSON 形态走 base64（膨胀 4/3），上限需
// 容纳解码后的文件。Content-Length 已超限时直接 413、不读 body；长度未知
// （chunked）时用 MaxBytesReader 兜住，读取超限会以 *http.MaxBytesError 中止解析。
func uploadBodyGuard() kratoshttp.FilterFunc {
	return uploadBodyGuardWith(func(r *nethttp.Request) int64 {
		if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
			return maxUploadBody
		}
		return maxRequestBody
	})
}

// uploadBodyGuardWith 按请求计算上限，便于测试用小阈值覆盖。
func uploadBodyGuardWith(limitFor func(*nethttp.Request) int64) kratoshttp.FilterFunc {
	return func(next nethttp.Handler) nethttp.Handler {
		return nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
			limit := limitFor(r)
			if r.ContentLength > limit {
				nethttp.Error(w, "request entity too large", nethttp.StatusRequestEntityTooLarge)
				return
			}
			r.Body = nethttp.MaxBytesReader(w, r.Body, limit)
			next.ServeHTTP(w, r)
		})
	}
}

// requestBodyGuard 挂在整个 REST 服务最外层（mux 之前），给所有请求封顶：
//   - 声明长度超过 maxRequestBody：直接 413，不读 body；
//   - 声明长度超过 maxUnauthenticatedBody（或长度未知）：必须携带有效访问令牌，
//     否则 401，不读 body——大 body 的合法场景（上传/头像/知识库文档）都是
//     已登录操作，免鉴权端点没有大 body 的需求；
//   - checker 为 nil（无认证器的机器对机器 webhook 服务）时只有总量上限。
//
// 之所以放在 HTTP Filter 而不是 kratos 中间件：后者的执行点在 handler 内部，
// 晚于 ctx.Bind / FormFile 的 body 读取。
func requestBodyGuard(checker auth.AccessTokenChecker) kratoshttp.FilterFunc {
	return func(next nethttp.Handler) nethttp.Handler {
		return nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
			limit := int64(maxUnauthenticatedBody)

			switch {
			case checker == nil:
				// 无认证器：只有总量上限（webhook 等机器对机器服务）。
				limit = maxRequestBody
			case r.ContentLength > maxRequestBody:
				nethttp.Error(w, "request entity too large", nethttp.StatusRequestEntityTooLarge)
				return
			case r.ContentLength > maxUnauthenticatedBody:
				// 大 body：先认证，再放行读取。
				if !hasValidAccessToken(r, checker) {
					nethttp.Error(w, "missing or invalid bearer token", nethttp.StatusUnauthorized)
					return
				}
				limit = maxRequestBody
			case r.ContentLength < 0:
				// 长度未知（chunked）：有有效令牌才允许超过未认证上限。
				if hasValidAccessToken(r, checker) {
					limit = maxRequestBody
				}
			}

			r.Body = nethttp.MaxBytesReader(w, r.Body, limit)
			next.ServeHTTP(w, r)
		})
	}
}

// hasValidAccessToken 在读取 body 之前校验 Bearer 令牌。
func hasValidAccessToken(r *nethttp.Request, checker auth.AccessTokenChecker) bool {
	if checker == nil {
		return false
	}
	token := bearerToken(r)
	if token == "" {
		return false
	}
	valid, _ := checker.IsValidAccessToken(r.Context(), token, false)
	return valid
}

// bearerToken 从 Authorization 头提取 Bearer 令牌，与 kratos-authn 的
// AuthFromMD 同口径（scheme 大小写不敏感）。
func bearerToken(r *nethttp.Request) string {
	scheme, token, ok := strings.Cut(r.Header.Get("Authorization"), " ")
	if !ok || !strings.EqualFold(scheme, "Bearer") {
		return ""
	}
	return strings.TrimSpace(token)
}

// isMaxBytesError 判断错误链中是否包含请求体超限错误。
func isMaxBytesError(err error) bool {
	var maxBytesErr *nethttp.MaxBytesError
	return errors.As(err, &maxBytesErr)
}
