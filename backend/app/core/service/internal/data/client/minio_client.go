package client

import (
	"go-wind-cms/pkg/oss"

	"github.com/tx7do/kratos-bootstrap/bootstrap"
)

func NewMinIoClient(ctx *bootstrap.Context) *oss.MinIOClient {
	return oss.NewMinIoClient(ctx.GetConfig(), ctx.GetLogger())
}
