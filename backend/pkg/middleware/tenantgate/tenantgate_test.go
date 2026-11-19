package tenantgate

import (
	"context"
	"errors"
	"testing"

	kratosmw "github.com/go-kratos/kratos/v2/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	identityV1 "go-wind-cms/api/gen/go/identity/service/v1"
	"google.golang.org/grpc"

	appMetadata "go-wind-cms/pkg/metadata"
)

// recordingChecker 记录调用并按预设返回。
type recordingChecker struct {
	calls     int
	gotTenant int64
	gotMethod string
	gotPath   string
	err       error
}

func (r *recordingChecker) CheckTenantAccess(ctx context.Context, in *identityV1.CheckTenantAccessRequest, opts ...grpc.CallOption) (*identityV1.CheckTenantAccessResponse, error) {
	r.calls++
	r.gotTenant = int64(in.GetTenantId())
	r.gotMethod = in.GetMethod()
	r.gotPath = in.GetPath()
	if r.err != nil {
		return nil, r.err
	}
	return &identityV1.CheckTenantAccessResponse{}, nil
}

func withOperator(ctx context.Context, tenantID int64) context.Context {
	op := appMetadata.NewUserOperator(42, tenantID, 0, 0)
	c, err := appMetadata.NewContext(ctx, op)
	if err != nil {
		panic(err)
	}
	return c
}

// runGate 以一个"返回 handled"的终结 handler 挂载中间件后执行。
func runGate(t *testing.T, mw kratosmw.Middleware, ctx context.Context, req any) (any, error) {
	t.Helper()
	inner := func(ctx context.Context, req any) (any, error) {
		return "handled", nil
	}
	return mw(inner)(ctx, req)
}

// ── Server 行为 ──

func TestGateNilCheckerPassThrough(t *testing.T) {
	// checker 未装配（生产 fail-open 语义）：请求直通，租户业务不被卡死。
	mw := Server(nil)
	resp, err := runGate(t, mw, context.Background(), nil)
	require.NoError(t, err)
	assert.Equal(t, "handled", resp)
}

func TestGateSkipsAnonymous(t *testing.T) {
	// 无 OperatorMetadata（匿名/白名单请求）：不调用闸门。
	checker := &recordingChecker{}
	mw := Server(checker)
	resp, err := runGate(t, mw, context.Background(), nil)
	require.NoError(t, err)
	assert.Equal(t, "handled", resp)
	assert.Equal(t, 0, checker.calls, "anonymous request must not hit the checker")
}

func TestGateSkipsPlatformIdentity(t *testing.T) {
	// tenant==0（平台身份）跳过闸门。
	checker := &recordingChecker{}
	mw := Server(checker)
	resp, err := runGate(t, mw, withOperator(context.Background(), 0), nil)
	require.NoError(t, err)
	assert.Equal(t, "handled", resp)
	assert.Equal(t, 0, checker.calls)
}

func TestGateChecksTenant(t *testing.T) {
	checker := &recordingChecker{}
	mw := Server(checker)
	resp, err := runGate(t, mw, withOperator(context.Background(), 7), nil)
	require.NoError(t, err)
	assert.Equal(t, "handled", resp)
	assert.Equal(t, 1, checker.calls)
	assert.Equal(t, int64(7), checker.gotTenant)
}

func TestGatePropagatesRejection(t *testing.T) {
	// 配额明确存在且已用满 → core 返回错误 → 闸门拒绝且 handler 不执行。
	checker := &recordingChecker{err: errors.New("api call quota exhausted")}
	mw := Server(checker)
	resp, err := runGate(t, mw, withOperator(context.Background(), 7), nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "quota exhausted")
	assert.Nil(t, resp)
}

func TestGateRequestInfoOptional(t *testing.T) {
	// 无 transport 时 method/path 为空串——判定不应受影响（core 侧日志缺参而已）。
	checker := &recordingChecker{}
	mw := Server(checker)
	_, err := runGate(t, mw, withOperator(context.Background(), 7), nil)
	require.NoError(t, err)
	assert.Equal(t, "", checker.gotMethod)
	assert.Equal(t, "", checker.gotPath)
}

// ── requestInfo ──

func TestRequestInfoWithoutTransport(t *testing.T) {
	method, path := requestInfo(context.Background())
	assert.Equal(t, "", method)
	assert.Equal(t, "", path)
}

// ── Benchmarks ──

func BenchmarkGatePassThrough(b *testing.B) {
	// 平台身份（tenant=0）跳过闸门的开销路径。
	mw := Server(&recordingChecker{})
	ctx := withOperator(context.Background(), 0)
	inner := func(ctx context.Context, req any) (any, error) { return struct{}{}, nil }
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = mw(inner)(ctx, nil)
	}
}
