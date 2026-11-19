package jwt

import (
	"testing"

	"github.com/tx7do/go-utils/trans"
	identityV1 "go-wind-cms/api/gen/go/identity/service/v1"
)

func benchPayload() payloadArgs {
	ou := int64(42)
	ds := identityV1.DataScope(1)
	client := "client_123"
	device := "device_1"
	return payloadArgs{
		username: "bench_user",
		uid:      1,
		tid:      2,
		ouid:     &ou,
		roles:    []string{"tenant:manager", "tenant:user"},
		ds:       &ds,
		client:   &client,
		device:   &device,
	}
}

type payloadArgs struct {
	username string
	uid      int64
	tid      int64
	ouid     *int64
	roles    []string
	ds       *identityV1.DataScope
	client   *string
	device   *string
}

// BenchmarkNewUserTokenPayload 签发路径：每次登录构造 payload 的开销。
func BenchmarkNewUserTokenPayload(b *testing.B) {
	a := benchPayload()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = NewUserTokenPayload(a.username, a.uid, a.tid, a.ouid, a.roles, a.ds, a.client, a.device)
	}
}

// BenchmarkNewUserTokenPayloadWithJwtMapClaims 解析路径：中间件从 JWT claims
// 还原 payload 的开销（logging/auth 链路每请求一次）。
func BenchmarkNewUserTokenPayloadWithJwtMapClaims(b *testing.B) {
	src := NewUserTokenPayload(
		"bench_user", 1, 2, trans.Ptr(int64(42)),
		[]string{"tenant:manager"},
		trans.Ptr(identityV1.DataScope(1)),
		trans.Ptr("client_123"), trans.Ptr("device_1"),
	)
	claims := NewUserTokenAuthClaims(src, nil)
	if claims == nil {
		b.Fatal("claims must not be nil")
	}
	// AuthClaims → MapClaims 由上游转换；直接走 MapClaims 入口需要 token 解出。
	// 这里基准 Claims 结构体还原路径。
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := NewUserTokenPayloadWithClaims(claims); err != nil {
			b.Fatal(err)
		}
	}
}
