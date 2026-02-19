package server

import (
	"github.com/dtm-labs/dtm/client/workflow"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/transport/grpc"

	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"github.com/tx7do/kratos-bootstrap/rpc"

	"go-wind-cms/pkg/middleware/ent"
	serviceName "go-wind-cms/pkg/serviceid"
)

type GrpcMiddlewares []middleware.Middleware

func NewGrpcMiddleware(ctx *bootstrap.Context) GrpcMiddlewares {
	var ms GrpcMiddlewares
	ms = append(ms, logging.Server(ctx.GetLogger()))
	ms = append(ms, ent.Server())
	return ms
}

// NewGrpcServer creates a gRPC server.
func NewGrpcServer(
	ctx *bootstrap.Context,
	middlewares GrpcMiddlewares,
) (*grpc.Server, error) {
	cfg := ctx.GetConfig()

	if cfg == nil || cfg.Server == nil || cfg.Server.Grpc == nil {
		return nil, nil
	}

	srv, err := rpc.CreateGrpcServer(cfg, middlewares...)
	if err != nil {
		return nil, err
	}

	en, err := srv.Endpoint()
	if err != nil {
		return nil, err
	}

	log.Infof("grpc server listening on: %s", en.String())

	// 注册操作需要在业务服务启动之后执行，因为当进程crash，dtm会回调业务服务器，继续未完成的任务
	workflow.InitGrpc(serviceName.DtmServiceAddress, en.String(), srv.Server)

	return srv, nil
}
