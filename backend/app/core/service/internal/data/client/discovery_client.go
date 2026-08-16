package client

import (
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	bRegistry "github.com/tx7do/kratos-bootstrap/registry"
)

// NewDiscovery 创建服务发现客户端
func NewDiscovery(ctx *bootstrap.Context) registry.Discovery {
	cfg := ctx.GetConfig()
	if cfg == nil {
		return nil
	}

	ret, err := bRegistry.NewDiscovery(cfg.Registry)
	if err != nil {
		return nil
	}
	return ret
}
