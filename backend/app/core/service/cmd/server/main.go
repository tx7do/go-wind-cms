package main

import (
	"context"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/transport/grpc"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	_ "github.com/tx7do/kratos-bootstrap/registry/etcd"
	_ "github.com/tx7do/kratos-bootstrap/tracer"

	"go-wind-cms/pkg/serviceid"
)

var version = "1.0.0"

// go build -ldflags "-X main.version=x.y.z"

func newApp(ctx *bootstrap.Context, gs *grpc.Server) *kratos.App {
	return bootstrap.NewApp(ctx, gs)
}

func runApp() error {
	ctx := bootstrap.NewContext(
		context.Background(),
		&conf.AppInfo{
			Project: serviceid.ProjectName,
			AppId:   serviceid.CoreService,
			Version: version,
		},
	)
	return bootstrap.RunApp(ctx, initApp)
}

func main() {
	if err := runApp(); err != nil {
		panic(err)
	}
}
