package server

import (
	"errors"
	kratoshttp "github.com/go-kratos/kratos/v2/transport/http"
	nethttp "net/http"
	"strings"
)

const (
	// maxUploadFileBytes 单文件上传上限。
	maxUploadFileBytes int64 = 50 << 20

	// maxUploadBody multipart 上传路由的请求体上限：单文件上限 + multipart 编码余量。
	maxUploadBody = maxUploadFileBytes + 1<<20

	// maxRequestBody 全站请求体上限。需容纳 base64 编码后的上传体（膨胀 4/3）
	// 与 JSON 结构开销；比它更小的请求不受影响。
	maxRequestBody = maxUploadFileBytes*4/3 + 4<<20
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

// requestBodyGuard 挂在整个 REST 服务最外层（mux 之前），给所有请求封顶。
// 服务无认证器（机器对机器 webhook），因此只有总量上限。
// 之所以放在 HTTP Filter 而不是 kratos 中间件：后者的执行点在 handler 内部，
// 晚于 ctx.Bind / FormFile 的 body 读取。
func requestBodyGuard() kratoshttp.FilterFunc {
	return uploadBodyGuardWith(func(*nethttp.Request) int64 { return maxRequestBody })
}

// isMaxBytesError 判断错误链中是否包含请求体超限错误。
func isMaxBytesError(err error) bool {
	var maxBytesErr *nethttp.MaxBytesError
	return errors.As(err, &maxBytesErr)
}
