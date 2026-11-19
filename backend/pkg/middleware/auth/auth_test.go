package auth

import (
	"context"
	"testing"

	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/transport"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	authenticationV1 "go-wind-cms/api/gen/go/authentication/service/v1"
	identityV1 "go-wind-cms/api/gen/go/identity/service/v1"

	"github.com/tx7do/go-crud/viewer"
	appViewer "go-wind-cms/pkg/entgo/viewer"
	appMetadata "go-wind-cms/pkg/metadata"
)

// ── 测试桩 ──

// headerMap transport.Header 的最小实现。
type headerMap map[string]string

func (h headerMap) Get(key string) string      { return h[key] }
func (h headerMap) Set(key, value string)      { h[key] = value }
func (h headerMap) Add(key, value string)      { h[key] = value }
func (h headerMap) Keys() []string             { return nil }
func (h headerMap) Values(key string) []string { return nil }

// fakeTransporter 测试用 transport.Transporter（kratos http.Transport 字段全未导出）。
type fakeTransporter struct {
	headers   headerMap
	operation string
}

func (f *fakeTransporter) Kind() transport.Kind { return transport.KindHTTP }
func (f *fakeTransporter) Endpoint() string     { return "http://localhost" }
func (f *fakeTransporter) Operation() string    { return f.operation }
func (f *fakeTransporter) RequestHeader() transport.Header {
	return f.headers
}
func (f *fakeTransporter) ReplyHeader() transport.Header { return headerMap{} }

func newServerCtx(headers map[string]string) context.Context {
	if headers == nil {
		headers = map[string]string{}
	}
	tr := &fakeTransporter{headers: headerMap(headers), operation: "/test.Service/Get"}
	return transport.NewServerContext(context.Background(), tr)
}

// stubChecker 恒定返回给定结果的检查器。
func stubChecker(valid bool, payload *authenticationV1.UserTokenPayload) AccessTokenChecker {
	return NewAccessTokenCheckerFromFuncs(
		func(ctx context.Context, accessToken string, skipRedis bool) (bool, *authenticationV1.UserTokenPayload) {
			return valid, payload
		},
		nil,
	)
}

func testPayload() *authenticationV1.UserTokenPayload {
	tid := int64(7)
	ou := int64(3)
	ds := identityV1.DataScope_ALL
	return &authenticationV1.UserTokenPayload{
		UserId:    42,
		TenantId:  &tid,
		OrgUnitId: &ou,
		Username:  stringPtr("alice"),
		Roles:     []string{"tenant:manager"},
		DataScope: &ds,
	}
}

func stringPtr(s string) *string { return &s }

// nextHandler 记录调用并回显 context 内容。
type recorder struct {
	called bool
	ctx    context.Context
	req    any
}

func (r *recorder) handler(ctx context.Context, req any) (any, error) {
	r.called = true
	r.ctx = ctx
	r.req = req
	return "ok", nil
}

// ── 错误哨兵 ──

func TestSentinelErrors(t *testing.T) {
	for _, err := range []error{
		ErrWrongContext,
		ErrMissingBearerToken,
		ErrValidateTokenFailed,
		ErrMissingJwtToken,
		ErrExtractUserInfoFailed,
		ErrExtractSubjectFailed,
		ErrAccessTokenExpired,
		ErrInvalidRequest,
		ErrAccessTokenCheckerNotConfigured,
	} {
		assert.NotNil(t, err)
	}

	// 全部为 kratos Unauthorized，前端按 401 处理。
	for _, err := range []error{
		ErrWrongContext, ErrMissingBearerToken, ErrMissingJwtToken, ErrAccessTokenExpired,
	} {
		assert.True(t, kerrors.IsUnauthorized(err), "%v must be unauthorized", err)
		assert.Equal(t, "UNAUTHORIZED", kerrors.Reason(err))
	}
}

// ── Context 注入/提取 ──

func TestNewContextFromContext(t *testing.T) {
	// 语义钉测：NewContext 存入 typed-nil payload 时 FromContext 的
	// 类型断言成功（ok=true, 值为 nil）——不返回 ErrMissingJwtToken。
	// 消费方（service 层）依赖 proto getter 的 nil 安全（GetUserId()=0）。
	ctx := NewContext(context.Background(), nil)
	got, err := FromContext(ctx)
	assert.NoError(t, err)
	assert.Nil(t, got)

	payload := testPayload()
	ctx = NewContext(context.Background(), payload)
	got, err = FromContext(ctx)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, int64(42), got.GetUserId())
	assert.Equal(t, int64(7), got.GetTenantId())
	assert.Equal(t, []string{"tenant:manager"}, got.GetRoles())
}

func TestFromContextWithoutPayload(t *testing.T) {
	_, err := FromContext(context.Background())
	require.Error(t, err)
	assert.Equal(t, ErrMissingJwtToken, err)
}

// ── composedChecker ──

func TestComposedCheckerDefaults(t *testing.T) {
	// valid 未提供：默认放行；blocker 未提供：默认不阻止。
	c := NewAccessTokenCheckerFromFuncs(nil, nil)
	valid, payload := c.IsValidAccessToken(context.Background(), "tok", false)
	assert.True(t, valid)
	assert.Nil(t, payload)
	assert.False(t, c.IsBlockedAccessToken(context.Background(), "tok"))
}

func TestComposedCheckerFuncs(t *testing.T) {
	var gotToken string
	var gotSkipRedis bool
	valid := func(ctx context.Context, accessToken string, skipRedis bool) (bool, *authenticationV1.UserTokenPayload) {
		gotToken = accessToken
		gotSkipRedis = skipRedis
		return true, testPayload()
	}
	blocked := func(ctx context.Context, accessToken string) bool { return accessToken == "blocked" }

	c := NewAccessTokenCheckerFromFuncs(valid, blocked)

	ok, payload := c.IsValidAccessToken(context.Background(), "tok", true)
	assert.True(t, ok)
	require.NotNil(t, payload)
	assert.Equal(t, int64(42), payload.GetUserId())
	assert.Equal(t, "tok", gotToken)
	assert.True(t, gotSkipRedis)

	assert.True(t, c.IsBlockedAccessToken(context.Background(), "blocked"))
	assert.False(t, c.IsBlockedAccessToken(context.Background(), "tok"))
}

func TestAccessTokenBlockerFunc(t *testing.T) {
	b := AccessTokenBlockerFunc(func(ctx context.Context, token string) bool { return token == "bad" })
	assert.True(t, b.IsBlockedAccessToken(context.Background(), "bad"))
	assert.False(t, b.IsBlockedAccessToken(context.Background(), "good"))
}

// ── Option 装配 ──

func TestOptionsDefaultsAndSetters(t *testing.T) {
	// 走 Server 一层拿不到内部 options，这里通过行为间接验证关键默认值。
	// 默认：injectMetadata/injectEnt/enableAuthz 开，injectOperatorId/injectTenantId 关。
	mw := Server(WithAccessTokenChecker(stubChecker(true, testPayload())))
	require.NotNil(t, mw)

	// withLogger 便于日志重定向（行为无关紧要，只验证不 panic）。
	mw2 := Server(WithAccessTokenChecker(stubChecker(true, testPayload())), WithLogger(nil))
	require.NotNil(t, mw2)

	var opts options
	for _, o := range []Option{
		WithAccessTokenChecker(stubChecker(true, nil)),
		WithAccessTokenCheckerFromFuncs(nil, nil),
		WithEnableCheckRefreshTokenExpiration(true),
		WithEnableCheckScopes(true),
		WithInjectOperatorId(true),
		WithInjectTenantId(true),
		WithInjectEnt(false),
		WithInjectMetadata(false),
		WithEnableAuthority(false),
		WithLogger(nil),
	} {
		o(&opts)
	}
	assert.NotNil(t, opts.accessTokenChecker)
	assert.True(t, opts.enableCheckRefreshTokenExpiration)
	assert.True(t, opts.enableCheckScopes)
	assert.True(t, opts.injectOperatorId)
	assert.True(t, opts.injectTenantId)
	assert.False(t, opts.injectEnt)
	assert.False(t, opts.injectMetadata)
	assert.False(t, opts.enableAuthz)
}

// ── Server 中间件链路 ──

func TestServerWrongContext(t *testing.T) {
	mw := Server(WithAccessTokenChecker(stubChecker(true, testPayload())))
	h := mw(func(ctx context.Context, req any) (any, error) { return "ok", nil })

	_, err := h(context.Background(), nil)
	require.Error(t, err)
	assert.Equal(t, ErrWrongContext, err)
}

func TestServerMissingBearerToken(t *testing.T) {
	mw := Server(WithAccessTokenChecker(stubChecker(true, testPayload())))
	h := mw(func(ctx context.Context, req any) (any, error) { return "ok", nil })

	_, err := h(newServerCtx(nil), nil)
	require.Error(t, err)
	assert.Equal(t, ErrMissingBearerToken, err)
}

func TestServerCheckerNotConfigured(t *testing.T) {
	mw := Server() // checker 未配置
	h := mw(func(ctx context.Context, req any) (any, error) { return "ok", nil })

	_, err := h(newServerCtx(map[string]string{"Authorization": "Bearer tok"}), nil)
	require.Error(t, err)
	assert.Equal(t, ErrAccessTokenCheckerNotConfigured, err)
}

func TestServerInvalidToken(t *testing.T) {
	mw := Server(WithAccessTokenChecker(stubChecker(false, nil)))
	h := mw(func(ctx context.Context, req any) (any, error) { return "ok", nil })

	_, err := h(newServerCtx(map[string]string{"Authorization": "Bearer tok"}), nil)
	require.Error(t, err)
	assert.Equal(t, ErrAccessTokenExpired, err)
}

func TestServerHappyPath(t *testing.T) {
	payload := testPayload()
	rec := &recorder{}
	mw := Server(
		WithAccessTokenChecker(stubChecker(true, payload)),
		WithInjectMetadata(true),
		WithInjectEnt(true),
	)
	h := mw(rec.handler)

	resp, err := h(newServerCtx(map[string]string{"Authorization": "Bearer tok"}), nil)
	require.NoError(t, err)
	assert.Equal(t, "ok", resp)
	require.True(t, rec.called)

	// claims 注入 context
	got, err := FromContext(rec.ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(42), got.GetUserId())

	// operator metadata 注入 server metadata
	op, err := appMetadata.FromServerContext(rec.ctx)
	require.NoError(t, err)
	require.NotNil(t, op)
	assert.Equal(t, int64(42), int64(op.GetUserId()))
	assert.Equal(t, int64(7), int64(op.GetTenantId()))

	// ent viewer 注入：UserViewer 按租户隔离
	vc, ok := viewer.FromContext(rec.ctx)
	require.True(t, ok)
	require.NotNil(t, vc)
	assert.Equal(t, int64(42), vc.UserID())
	assert.Equal(t, int64(7), vc.TenantID())
	assert.True(t, vc.IsTenantContext())
	assert.False(t, vc.IsSystemContext())
	uv, isUser := vc.(appViewer.UserViewer)
	require.True(t, isUser)
	scopes := uv.DataScope()
	require.Len(t, scopes, 1)
	assert.Equal(t, viewer.ScopeTypeAll, scopes[0].ScopeType)
}

func TestServerInjectsOperatorId(t *testing.T) {
	payload := testPayload()
	rec := &recorder{}
	mw := Server(WithAccessTokenChecker(stubChecker(true, payload)), WithInjectOperatorId(true))
	h := mw(rec.handler)

	req := struct{ OperatorId *int64 }{}
	_, err := h(newServerCtx(map[string]string{"Authorization": "Bearer tok"}), &req)
	require.NoError(t, err)
	require.NotNil(t, req.OperatorId)
	assert.Equal(t, int64(42), *req.OperatorId)
}

func TestServerInjectOperatorIdNilReq(t *testing.T) {
	payload := testPayload()
	mw := Server(WithAccessTokenChecker(stubChecker(true, payload)), WithInjectOperatorId(true))
	h := mw(func(ctx context.Context, req any) (any, error) { return "ok", nil })

	_, err := h(newServerCtx(map[string]string{"Authorization": "Bearer tok"}), nil)
	require.Error(t, err)
	assert.Equal(t, ErrInvalidRequest, err)
}

func TestServerNilPayloadStaysSafe(t *testing.T) {
	// checker 返回 valid=true 且 payload=nil（异常组合）：中间件不得 panic。
	// 语义钉测（含已知坑）：零值 operator 经 proto.Marshal 得 0 字节 → base64
	// 为空串 → metadata 键丢失，FromServerContext 报 ErrNoOperatorHeader；
	// ent viewer（injectEnt）仍以零值 UserViewer 注入。生产链路 checker 不会
	// 返回 nil payload，此路径属防御面。
	rec := &recorder{}
	mw := Server(
		WithAccessTokenChecker(stubChecker(true, nil)),
		WithInjectMetadata(true),
		WithInjectEnt(true),
	)
	h := mw(rec.handler)

	_, err := h(newServerCtx(map[string]string{"Authorization": "Bearer tok"}), nil)
	require.NoError(t, err)
	require.True(t, rec.called)

	_, mdErr := appMetadata.FromServerContext(rec.ctx)
	require.Error(t, mdErr, "zero operator proto-marshals to empty bytes and loses the metadata key")

	vc, ok := viewer.FromContext(rec.ctx)
	require.True(t, ok)
	require.NotNil(t, vc)
	assert.Equal(t, int64(0), vc.UserID())
}

func TestServerWithoutMetadataStillInjectsClaims(t *testing.T) {
	// 默认（生产 wiring 用 WithInjectMetadata/WithInjectEnt）之外的组合：
	// 只校验 token，claims 必须可用。
	payload := testPayload()
	rec := &recorder{}
	mw := Server(WithAccessTokenChecker(stubChecker(true, payload)), WithInjectMetadata(false), WithInjectEnt(false))
	h := mw(rec.handler)

	_, err := h(newServerCtx(map[string]string{"Authorization": "Bearer tok"}), nil)
	require.NoError(t, err)
	require.True(t, rec.called)
	got, err := FromContext(rec.ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(42), got.GetUserId())
}

func TestServerBearerCaseInsensitive(t *testing.T) {
	// "bearer"（小写）也应被接受（AuthFromMD 用 EqualFold）。
	payload := testPayload()
	rec := &recorder{}
	mw := Server(WithAccessTokenChecker(stubChecker(true, payload)), WithInjectMetadata(false), WithInjectEnt(false))
	h := mw(rec.handler)

	_, err := h(newServerCtx(map[string]string{"Authorization": "bearer tok"}), nil)
	require.NoError(t, err)
	assert.True(t, rec.called)
}

// ── Benchmarks ──

func BenchmarkServerHappyPath(b *testing.B) {
	payload := testPayload()
	mw := Server(
		WithAccessTokenChecker(stubChecker(true, payload)),
		WithInjectMetadata(true),
		WithInjectEnt(true),
	)
	chain := mw(func(ctx context.Context, req any) (any, error) { return struct{}{}, nil })
	ctx := newServerCtx(map[string]string{"Authorization": "Bearer tok"})
	req := &struct{ OperatorId *int64 }{}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = chain(ctx, req)
	}
}
