package server

import (
	"errors"
	"io"
	nethttp "net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUploadBodyGuardRejectsDeclaredOversizeWithoutReading(t *testing.T) {
	body := &countingBody{}
	nextCalled := false

	h := uploadBodyGuardWith(func(*nethttp.Request) int64 { return 16 })(nethttp.HandlerFunc(func(nethttp.ResponseWriter, *nethttp.Request) {
		nextCalled = true
	}))

	req := httptest.NewRequest(nethttp.MethodPost, "/upload", body)
	req.ContentLength = 17
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != nethttp.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, want %d", rec.Code, nethttp.StatusRequestEntityTooLarge)
	}
	if nextCalled {
		t.Fatal("next handler should not be called for oversize body")
	}
	if body.read {
		t.Fatal("body must not be read when Content-Length already exceeds the limit")
	}
}

func TestRequestBodyGuardCapsEverythingAtGlobalLimit(t *testing.T) {
	body := &countingBody{}
	nextCalled := false

	h := requestBodyGuard()(nethttp.HandlerFunc(func(nethttp.ResponseWriter, *nethttp.Request) {
		nextCalled = true
	}))

	req := httptest.NewRequest(nethttp.MethodPost, "/hook", body)
	req.ContentLength = maxRequestBody + 1
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != nethttp.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, want %d", rec.Code, nethttp.StatusRequestEntityTooLarge)
	}
	if nextCalled {
		t.Fatal("next handler should not be called for an oversize body")
	}
	if body.read {
		t.Fatal("body must not be read when Content-Length already exceeds the global limit")
	}
}

func TestRequestBodyGuardAllowsSmallBody(t *testing.T) {
	nextCalled := false

	h := requestBodyGuard()(nethttp.HandlerFunc(func(nethttp.ResponseWriter, *nethttp.Request) {
		nextCalled = true
	}))

	payload := strings.Repeat("a", 64)
	req := httptest.NewRequest(nethttp.MethodPost, "/hook", strings.NewReader(payload))
	req.ContentLength = int64(len(payload))
	h.ServeHTTP(httptest.NewRecorder(), req)

	if !nextCalled {
		t.Fatal("small request should pass through")
	}
}

func TestRequestBodyGuardCapsChunkedBody(t *testing.T) {
	var readErr error

	h := requestBodyGuard()(nethttp.HandlerFunc(func(_ nethttp.ResponseWriter, r *nethttp.Request) {
		_, readErr = io.ReadAll(r.Body)
	}))

	req := httptest.NewRequest(nethttp.MethodPost, "/hook", strings.NewReader(strings.Repeat("a", int(maxRequestBody)+1024)))
	req.ContentLength = -1

	h.ServeHTTP(httptest.NewRecorder(), req)

	var maxBytesErr *nethttp.MaxBytesError
	if !errors.As(readErr, &maxBytesErr) {
		t.Fatalf("read error = %v, want *http.MaxBytesError", readErr)
	}
}

// countingBody 记录是否被读过；用于断言"拦截发生在读取 body 之前"。
type countingBody struct {
	read bool
}

func (b *countingBody) Read([]byte) (int, error) {
	b.read = true
	return 0, io.EOF
}

func (b *countingBody) Close() error { return nil }
