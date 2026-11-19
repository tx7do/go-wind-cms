package ent

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-kratos/kratos/v2/transport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tx7do/go-crud/viewer"
	broker "github.com/tx7do/kratos-transport/broker"

	authenticationV1 "go-wind-cms/api/gen/go/authentication/service/v1"

	appViewer "go-wind-cms/pkg/entgo/viewer"
	appMetadata "go-wind-cms/pkg/metadata"
)

// ── 测试桩 ──

// headerMap transport.Header 最小实现。
type headerMap map[string]string

func (h headerMap) Get(key string) string      { return h[key] }
func (h headerMap) Set(key, value string)      { h[key] = value }
func (h headerMap) Add(key, value string)      { h[key] = value }
func (h headerMap) Keys() []string             { return nil }
func (h headerMap) Values(key string) []string { return nil }

// fakeTransporter transport.Transporter + khttp.Transporter（HTTP 面完整）。
type fakeTransporter struct {
	headers      headerMap
	request      *http.Request
	operation    string
	pathTemplate string
}

func (f *fakeTransporter) Kind() transport.Kind            { return transport.KindHTTP }
func (f *fakeTransporter) Endpoint() string                { return "http://localhost" }
func (f *fakeTransporter) Operation() string               { return f.operation }
func (f *fakeTransporter) RequestHeader() transport.Header { return f.headers }
func (f *fakeTransporter) ReplyHeader() transport.Header   { return headerMap{} }
func (f *fakeTransporter) Request() *http.Request          { return f.request }
func (f *fakeTransporter) PathTemplate() string            { return f.pathTemplate }

// fakeEvent broker.Event 最小实现（Subscriber 中间件测试用）。
type fakeEvent struct{ msg *broker.Message }

func (e *fakeEvent) Topic() string            { return "test-topic" }
func (e *fakeEvent) Message() *broker.Message { return e.msg }
func (e *fakeEvent) RawMessage() any          { return nil }
func (e *fakeEvent) Ack() error               { return nil }
func (e *fakeEvent) Error() error             { return nil }

// stubResolver TenantResolver 桩。
type stubResolver struct {
	tid int64
	err error
}

func (s *stubResolver) ResolveTenantIDByDomain(ctx context.Context, domain string) (int64, error) {
	return s.tid, s.err
}

func newHTTPCtx(req *http.Request) context.Context {
	tr := &fakeTransporter{
		headers:      headerMap{},
		request:      req,
		operation:    "/test/Get",
		pathTemplate: "/test/{id}",
	}
	return transport.NewServerContext(context.Background(), tr)
}

func viewerOf(ctx context.Context) (viewer.Context, bool) {
	return viewer.FromContext(ctx)
}

// ── metaDataToUserViewerContext ──

func TestMetaDataToUserViewerContext(t *testing.T) {
	assert.Nil(t, metaDataToUserViewerContext(nil, "trace-1"), "nil metadata must map to nil")

	md := &authenticationV1.OperatorMetadata{UserId: 42, TenantId: 7, OrgUnitId: 3}
	vc := metaDataToUserViewerContext(md, "trace-1")
	require.NotNil(t, vc)
	assert.Equal(t, int64(42), vc.UserID())
	assert.Equal(t, int64(7), vc.TenantID())
	assert.Equal(t, int64(3), vc.OrgUnitID())
	assert.Equal(t, "trace-1", vc.TraceID())
	assert.True(t, vc.IsTenantContext())
}

// ── resolveAnonymousViewer ──

func TestResolveAnonymousViewerNoResolver(t *testing.T) {
	// 未注入 resolver（admin/core 场景）：SystemViewer 兜底。
	vc := resolveAnonymousViewer(context.Background(), nil, "trace-1")
	require.NotNil(t, vc)
	assert.True(t, vc.IsSystemContext())
	_, isSystem := vc.(appViewer.SystemViewer)
	assert.True(t, isSystem)
}

func TestResolveAnonymousViewerWithHost(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://tenant-a.example.com/", nil)
	ctx := newHTTPCtx(req)

	resolver := &stubResolver{tid: 9}
	vc := resolveAnonymousViewer(ctx, resolver, "trace-1")
	require.NotNil(t, vc)
	av, ok := vc.(appViewer.AnonymousTenantViewer)
	require.True(t, ok, "resolved anonymous request must yield AnonymousTenantViewer")
	assert.Equal(t, int64(9), av.TenantID())
	assert.Equal(t, "trace-1", av.TraceID())
	assert.True(t, vc.IsTenantContext())
	assert.False(t, vc.IsSystemContext())
	assert.False(t, vc.HasPermission("read", "post"))
}

func TestResolveAnonymousViewerNoHost(t *testing.T) {
	// 非 HTTP transport（拿不到 Host）：fail-closed NoopContext。
	ctx := transport.NewServerContext(context.Background(), &fakeTransporter{headers: headerMap{}})
	vc := resolveAnonymousViewer(ctx, &stubResolver{tid: 9}, "")
	require.NotNil(t, vc)
	noop, ok := vc.(viewer.Context)
	require.True(t, ok)
	assert.Equal(t, int64(0), noop.TenantID())
	assert.False(t, noop.IsTenantContext())
	assert.False(t, noop.IsSystemContext())
}

func TestResolveAnonymousViewerResolveFailure(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://unknown.example.com/", nil)
	ctx := newHTTPCtx(req)

	// 解析失败：fail-closed NoopContext（不回退 SystemViewer 避免跨租户泄漏）。
	vc := resolveAnonymousViewer(ctx, &stubResolver{err: assert.AnError}, "")
	require.NotNil(t, vc)
	assert.Equal(t, int64(0), vc.TenantID())
	assert.False(t, vc.IsTenantContext())
	assert.False(t, vc.IsSystemContext())

	// 解析到 0：同样 fail-closed。
	vc = resolveAnonymousViewer(ctx, &stubResolver{tid: 0}, "")
	assert.Equal(t, int64(0), vc.TenantID())
	assert.False(t, vc.IsTenantContext())
}

// ── Server 中间件 ──

func TestServerAuthenticatedRequest(t *testing.T) {
	op := appMetadata.NewUserOperator(42, 7, 3, 0)
	mdCtx, err := appMetadata.NewContext(context.Background(), op)
	require.NoError(t, err)

	rec := &recorder{}
	mw := Server()
	resp, err := mw(rec.handler)(mdCtx, nil)
	require.NoError(t, err)
	assert.Equal(t, "ok", resp)

	vc, ok := viewerOf(rec.ctx)
	require.True(t, ok)
	assert.Equal(t, int64(42), vc.UserID())
	assert.Equal(t, int64(7), vc.TenantID())
	assert.Equal(t, int64(3), vc.OrgUnitID())
	assert.True(t, vc.IsTenantContext())
}

func TestServerAnonymousWithResolver(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://tenant-a.example.com/", nil)
	ctx := newHTTPCtx(req)

	rec := &recorder{}
	mw := Server(WithTenantResolver(&stubResolver{tid: 9}))
	_, err := mw(rec.handler)(ctx, nil)
	require.NoError(t, err)

	vc, ok := viewerOf(rec.ctx)
	require.True(t, ok)
	av, isAnon := vc.(appViewer.AnonymousTenantViewer)
	require.True(t, isAnon)
	assert.Equal(t, int64(9), av.TenantID())
}

func TestServerAnonymousFailClosed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://unknown.example.com/", nil)
	ctx := newHTTPCtx(req)

	rec := &recorder{}
	mw := Server(WithTenantResolver(&stubResolver{err: assert.AnError}))
	_, err := mw(rec.handler)(ctx, nil)
	require.NoError(t, err, "middleware itself must not error; fail-closed happens at query time")

	vc, ok := viewerOf(rec.ctx)
	require.True(t, ok)
	// noopContext（go-crud 未导出类型）：租户 0 且非任何上下文。
	assert.Equal(t, int64(0), vc.TenantID())
	assert.False(t, vc.IsTenantContext())
	assert.False(t, vc.IsSystemContext())
}

type recorder struct {
	called bool
	ctx    context.Context
}

func (r *recorder) handler(ctx context.Context, req any) (any, error) {
	r.called = true
	r.ctx = ctx
	return "ok", nil
}

// ── InjectViewerForStream ──

func TestInjectViewerForStreamAuthenticated(t *testing.T) {
	op := appMetadata.NewUserOperator(11, 22, 33, 0)
	ctx, err := appMetadata.NewContext(context.Background(), op)
	require.NoError(t, err)

	out := InjectViewerForStream(ctx)
	vc, ok := viewerOf(out)
	require.True(t, ok)
	assert.Equal(t, int64(11), vc.UserID())
	assert.Equal(t, int64(22), vc.TenantID())
}

func TestInjectViewerForStreamAnonymous(t *testing.T) {
	// 匿名/异常时回落 SystemViewer（流式 RPC 无 resolver 可用）。
	out := InjectViewerForStream(context.Background())
	vc, ok := viewerOf(out)
	require.True(t, ok)
	assert.True(t, vc.IsSystemContext())
}

// ── broker Publish / Subscriber 头部往返 ──

func TestBrokerPublishSubscriberRoundTrip(t *testing.T) {
	// Publish：viewer 上下文 → 消息头；Subscriber：消息头 → viewer 上下文。
	pubMW := Publish()

	tid := int64(7)
	// UserViewer 不带 datascope 走编码路径（proto codec 对 []viewer.DataScope
	// 非 proto.Message，会 panic——见 Publish 源码注释，属休眠死代码路径）。
	uid := int64(42)
	ouid := int64(3)
	trace := "trace-roundtrip"
	userViewer := appViewer.NewUserViewerWithDataScopes(uid, tid, ouid, trace, nil)

	var sent *broker.Message
	pubHandler := func(ctx context.Context, topic string, msg *broker.Message, opts ...broker.PublishOption) error {
		sent = msg
		return nil
	}
	err := pubMW(pubHandler)(viewer.WithContext(context.Background(), userViewer), "topic", broker.NewMessage("body"))
	require.NoError(t, err)
	require.NotNil(t, sent)

	assert.Equal(t, "42", sent.GetHeader(UserIDHeaderKey))
	assert.Equal(t, "7", sent.GetHeader(TenantIDHeaderKey))
	assert.Equal(t, "3", sent.GetHeader(OrgUnitIDHeaderKey))
	// 语义钉测：Publish 的 trace 头取自 otel span context，而非 viewer.TraceID()
	// ——测试上下文无 span，头为空。
	assert.Empty(t, sent.GetHeader(TraceIDHeaderKey))

	// Publish 不回填 trace 头（取 otel span），手工补 32 位 hex（otel 格式）
	// 验证 Subscriber 侧 TraceIDFromHex 透传。
	sent.SetHeader(TraceIDHeaderKey, "0af7651916cd43dd8448eb211c80319c")

	// Subscriber：从消息头重建 viewer。
	subMW := Subscriber()
	var recv viewer.Context
	subHandler := func(ctx context.Context, evt broker.Event) error {
		v, ok := viewer.FromContext(ctx)
		if ok {
			recv = v
		}
		return nil
	}
	err = subMW(subHandler)(context.Background(), &fakeEvent{msg: sent})
	require.NoError(t, err)
	require.NotNil(t, recv)

	ruv, isUser := recv.(appViewer.UserViewer)
	require.True(t, isUser)
	assert.Equal(t, uid, ruv.UserID())
	assert.Equal(t, tid, ruv.TenantID())
	assert.Equal(t, ouid, ruv.OrgUnitID())
	assert.Equal(t, "0af7651916cd43dd8448eb211c80319c", ruv.TraceID())
}

func TestBrokerPublishWithDatascopeSkipsHeader(t *testing.T) {
	// UserViewer 自带一个 datascope（非 proto.Message）：修复前 codec.Marshal
	// 类型断言必 panic；修复后 warn 跳过 datascope 头，其余头照常。
	pubMW := Publish()
	var sent *broker.Message
	pubHandler := func(ctx context.Context, topic string, msg *broker.Message, opts ...broker.PublishOption) error {
		sent = msg
		return nil
	}

	uv := appViewer.NewUserViewer(42, 7, 3, "trace", 0)
	err := pubMW(pubHandler)(viewer.WithContext(context.Background(), uv), "topic", broker.NewMessage("body"))
	require.NoError(t, err, "publish must not panic on non-proto datascope")
	require.NotNil(t, sent)
	assert.Empty(t, sent.GetHeader(DataScopeHeaderKey))
	assert.Equal(t, "42", sent.GetHeader(UserIDHeaderKey))
}

func TestBrokerSubscriberInvalidTraceID(t *testing.T) {
	// 非法 trace 头（长度非 16/非 hex）：修复前 panic，修复后 warn 跳过。
	subMW := Subscriber()
	var recv viewer.Context
	subHandler := func(ctx context.Context, evt broker.Event) error {
		v, ok := viewer.FromContext(ctx)
		if ok {
			recv = v
		}
		return nil
	}

	msg := broker.NewMessage("body")
	msg.SetHeader(TraceIDHeaderKey, "trace-roundtrip")
	err := subMW(subHandler)(context.Background(), &fakeEvent{msg: msg})
	require.NoError(t, err, "invalid trace id must not panic the consumer")
	require.NotNil(t, recv)
	// span context 不注入，但原值仍进 viewer.TraceID（日志排查用）。
	assert.Equal(t, "trace-roundtrip", recv.TraceID())
	assert.Equal(t, int64(0), recv.UserID())
}

func TestBrokerPublishWithoutViewer(t *testing.T) {
	// 无 viewer 上下文：消息原样透传（best-effort 语义）。
	pubMW := Publish()
	var sent *broker.Message
	pubHandler := func(ctx context.Context, topic string, msg *broker.Message, opts ...broker.PublishOption) error {
		sent = msg
		return nil
	}
	err := pubMW(pubHandler)(context.Background(), "topic", broker.NewMessage("body"))
	require.NoError(t, err)
	require.NotNil(t, sent)
	assert.Empty(t, sent.GetHeader(UserIDHeaderKey))
	assert.Empty(t, sent.GetHeader(TenantIDHeaderKey))
}

func TestBrokerPublishNilViewerCtx(t *testing.T) {
	// viewer.FromContext 成功但值为 nil：透传不 panic。
	pubMW := Publish()
	var sent *broker.Message
	pubHandler := func(ctx context.Context, topic string, msg *broker.Message, opts ...broker.PublishOption) error {
		sent = msg
		return nil
	}
	err := pubMW(pubHandler)(viewer.WithContext(context.Background(), nil), "topic", broker.NewMessage("body"))
	require.NoError(t, err)
	require.NotNil(t, sent)
	assert.Empty(t, sent.GetHeader(UserIDHeaderKey))
	assert.Empty(t, sent.GetHeader(TenantIDHeaderKey))
}

func TestBrokerSubscriberSystemContext(t *testing.T) {
	// x-system=true 头：注入 SystemViewer 并短路（不落 UserViewer）。
	subMW := Subscriber()
	var recv viewer.Context
	subHandler := func(ctx context.Context, evt broker.Event) error {
		v, ok := viewer.FromContext(ctx)
		if ok {
			recv = v
		}
		return nil
	}

	msg := broker.NewMessage("body")
	msg.SetHeader(SystemContextHeaderKey, "true")
	err := subMW(subHandler)(context.Background(), &fakeEvent{msg: msg})
	require.NoError(t, err)
	require.NotNil(t, recv)
	assert.True(t, recv.IsSystemContext())
}

func TestBrokerSubscriberInvalidHeaders(t *testing.T) {
	// 非法数字/布尔头：不 panic，落零值 UserViewer（warn 日志）。
	subMW := Subscriber()
	var recv viewer.Context
	subHandler := func(ctx context.Context, evt broker.Event) error {
		v, ok := viewer.FromContext(ctx)
		if ok {
			recv = v
		}
		return nil
	}

	msg := broker.NewMessage("body")
	msg.SetHeader(UserIDHeaderKey, "not-a-number")
	msg.SetHeader(TenantIDHeaderKey, "also-bad")
	msg.SetHeader(SystemContextHeaderKey, "not-bool")
	err := subMW(subHandler)(context.Background(), &fakeEvent{msg: msg})
	require.NoError(t, err)
	require.NotNil(t, recv)
	assert.Equal(t, int64(0), recv.UserID())
	assert.Equal(t, int64(0), recv.TenantID())
	assert.False(t, recv.IsSystemContext())
}

func TestBrokerSubscriberEmptyHeaders(t *testing.T) {
	// 无消息头：直接透传，不注入 viewer。
	subMW := Subscriber()
	var handlerCtx context.Context
	subHandler := func(ctx context.Context, evt broker.Event) error {
		handlerCtx = ctx
		return nil
	}

	evt := &fakeEvent{msg: broker.NewMessage("body")}
	err := subMW(subHandler)(context.Background(), evt)
	require.NoError(t, err)
	_, ok := viewer.FromContext(handlerCtx)
	assert.False(t, ok, "headers-less message must not inject a viewer")
}

// ── Benchmark ──

func BenchmarkServerAuthenticated(b *testing.B) {
	op := appMetadata.NewUserOperator(42, 7, 3, 0)
	ctx, _ := appMetadata.NewContext(context.Background(), op)
	mw := Server()
	inner := func(ctx context.Context, req any) (any, error) { return struct{}{}, nil }
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = mw(inner)(ctx, nil)
	}
}
