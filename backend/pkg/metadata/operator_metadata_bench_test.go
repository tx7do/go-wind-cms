package metadata

import (
	"context"
	"testing"

	authenticationV1 "go-wind-cms/api/gen/go/authentication/service/v1"
)

func benchOperator() *authenticationV1.OperatorMetadata {
	return &authenticationV1.OperatorMetadata{
		Type:      authenticationV1.OperatorMetadata_USER,
		UserId:    42,
		TenantId:  7,
		OrgUnitId: 3,
	}
}

// BenchmarkEncodeOperatorMetadata 每请求注入路径的编码开销。
func BenchmarkEncodeOperatorMetadata(b *testing.B) {
	op := benchOperator()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = EncodeOperatorMetadata(op)
	}
}

// BenchmarkDecodeOperatorMetadata ent 中间件每请求的解码开销。
func BenchmarkDecodeOperatorMetadata(b *testing.B) {
	enc, err := EncodeOperatorMetadata(benchOperator())
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = DecodeOperatorMetadata(enc)
	}
}

// BenchmarkNewContextFromContext 完整的 NewContext→FromServerContext 链路。
func BenchmarkNewContextFromContext(b *testing.B) {
	op := benchOperator()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx, err := NewContext(context.Background(), op)
		if err != nil {
			b.Fatal(err)
		}
		if _, err := FromServerContext(ctx); err != nil {
			b.Fatal(err)
		}
	}
}
