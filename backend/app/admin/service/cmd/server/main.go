package main

import (
	"context"

	"github.com/go-kratos/kratos/v2"
	kratosTransport "github.com/go-kratos/kratos/v2/transport"


	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	_ "github.com/tx7do/kratos-bootstrap/registry/etcd"
	_ "github.com/tx7do/kratos-bootstrap/tracer"

	"go-wind-cms/app/admin/service/internal/server"
	"go-wind-cms/pkg/serviceid"
)

var version = "1.0.0"

// go build -ldflags "-X main.version=x.y.z"

func newApp(
	ctx *bootstrap.Context,
	servers ...kratosTransport.Server,
) *kratos.App {
	return bootstrap.NewApp(ctx, servers...)
}

func runApp() error {
	ctx := bootstrap.NewContext(
		context.Background(),
		&conf.AppInfo{
			Project: serviceid.ProjectName,
			AppId:   serviceid.AdminService,
			Version: version,
		},
	)
	// server.collab 配置段扩展：须在 bootstrap 读配置前登记（扫描列表届时定格）
	server.RegisterCollabConfigExt(ctx)
	return bootstrap.RunApp(ctx, initApp)
}

func main() {
	if err := runApp(); err != nil {
		panic(err)
	}
}
