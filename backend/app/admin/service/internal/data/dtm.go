package data

import (
	"github.com/dtm-labs/dtmdriver"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc/resolver/discovery"

	"google.golang.org/grpc/resolver"

	_ "go-wind-cms/pkg/dtmdriver-kratos"
	dtmdriverKratos "go-wind-cms/pkg/dtmdriver-kratos"
)

func NewDtmDriver(rr registry.Discovery) {
	// 注册Kratos的gRPC解析器的用于动态解析服务地址，用于与Dtm服务通信
	resolver.Register(discovery.NewBuilder(rr, discovery.WithInsecure(true)))

	// 激活 Kratos DTM Driver
	_ = dtmdriver.Use(dtmdriverKratos.Name)

	log.Infof("DTM Driver '%s' is registered and activated.", dtmdriver.GetDriver().GetName())
}
