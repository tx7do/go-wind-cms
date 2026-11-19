package serviceid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiscoveryNames(t *testing.T) {
	assert.Equal(t, "gowind-cms/core-service", NewDiscoveryName(CoreService))
	assert.Equal(t, "gowind-cms/admin-service", NewDiscoveryName(AdminService))
	assert.Equal(t, "gowind-cms/app-service", NewDiscoveryName(AppService))
}

func TestMakeDiscoveryAddress(t *testing.T) {
	assert.Equal(t, "discovery:///core-service", MakeDiscoveryAddress(CoreService))
}

func TestDtmServiceAddressPinned(t *testing.T) {
	// 包级初始化变量：DTM 服务发现地址与 DTMService 常量须保持一致。
	assert.Equal(t, MakeDiscoveryAddress(DTMService), DtmServiceAddress)
	assert.Equal(t, "discovery:///dtm-service", DtmServiceAddress)
}

func TestServiceIDConstantsPinned(t *testing.T) {
	assert.Equal(t, "gowind-cms", ProjectName)
	assert.Equal(t, "admin-service", AdminService)
	assert.Equal(t, "app-service", AppService)
	assert.Equal(t, "core-service", CoreService)
	assert.Equal(t, "dtm-service", DTMService)
}
