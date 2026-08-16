package metadata

import (
	"context"
	"encoding/base64"
	"testing"

	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/stretchr/testify/assert"

	authenticationV1 "go-wind-cms/api/gen/go/authentication/service/v1"
	identityV1 "go-wind-cms/api/gen/go/identity/service/v1"
)

func TestNewContext_FromContext_RoundTrip(t *testing.T) {
	// 使用空的 OperatorMetadata 做回环测试
	info := &authenticationV1.OperatorMetadata{}
	b, err := EncodeOperatorMetadata(info)
	ctx2 := metadata.NewServerContext(context.Background(), metadata.Metadata{
		mdOperatorKey: []string{b},
	})

	_, err = FromServerContext(ctx2)
	assert.Equal(t, err, ErrNoOperatorHeader)
}

func TestFromContext_MissingOrInvalid(t *testing.T) {
	// 丢失 header
	got, err := FromServerContext(context.Background())
	assert.Error(t, err)
	assert.Nil(t, got)

	// 无效 base64
	ctx := metadata.AppendToClientContext(context.Background(), mdOperatorKey, "not-base64")
	got, err = FromServerContext(ctx)
	assert.Error(t, err)
	assert.Nil(t, got)

	// 有效 base64 但不是有效的 proto bytes
	bad := base64.RawStdEncoding.EncodeToString([]byte("not-proto-bytes"))
	ctx = metadata.AppendToClientContext(context.Background(), mdOperatorKey, bad)
	got, err = FromServerContext(ctx)
	assert.Error(t, err)
	assert.Nil(t, got)
}

func TestNewContext_FromContext_InvalidData(t *testing.T) {
	// 使用空的 OperatorMetadata 做回环测试
	info := &authenticationV1.OperatorMetadata{
		UserId:    uint64(123),
		TenantId:  uint64(456),
		OrgUnitId: uint64(789),
		DataScope: identityV1.DataScope_ALL,
		RoleIds:   []uint64{1},
	}
	b, err := EncodeOperatorMetadata(info)
	ctx2 := metadata.NewServerContext(context.Background(), metadata.Metadata{
		mdOperatorKey: []string{b},
	})

	got, err := FromServerContext(ctx2)
	assert.NoError(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, info.UserId, got.UserId)
	assert.Equal(t, info.TenantId, got.TenantId)
	assert.Equal(t, info.OrgUnitId, got.OrgUnitId)
	assert.Equal(t, info.DataScope, got.DataScope)
	assert.Equal(t, info.RoleIds, got.RoleIds)
}

func TestNewContext_FromServerContext_RoundTrip(t *testing.T) {
	// NewContext 写入 server context，FromServerContext 应读回相同的操作员字段。
	info := &authenticationV1.OperatorMetadata{
		UserId:    uint64(123),
		TenantId:  uint64(456),
		OrgUnitId: uint64(789),
		DataScope: identityV1.DataScope_ALL,
		RoleIds:   []uint64{1},
	}
	ctx, err := NewContext(context.Background(), info)
	assert.NoError(t, err)

	got, err := FromServerContext(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, info.UserId, got.UserId)
	assert.Equal(t, info.TenantId, got.TenantId)
	assert.Equal(t, info.OrgUnitId, got.OrgUnitId)
	assert.Equal(t, info.DataScope, got.DataScope)
	assert.Equal(t, info.RoleIds, got.RoleIds)
}

func TestNewContext_PreservesExistingServerMetadata(t *testing.T) {
	// NewContext 采用合并语义：不得覆盖 server context 中已有的其它元数据键。
	pre := metadata.NewServerContext(context.Background(), metadata.Metadata{
		"x-md-global-other": []string{"keep"},
	})
	info := &authenticationV1.OperatorMetadata{UserId: 1}
	ctx, err := NewContext(pre, info)
	assert.NoError(t, err)

	md, ok := metadata.FromServerContext(ctx)
	assert.True(t, ok)
	assert.Equal(t, "keep", md.Get("x-md-global-other"))
	assert.NotEmpty(t, md.Get(mdOperatorKey))
}
