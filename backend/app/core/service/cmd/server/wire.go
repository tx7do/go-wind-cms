//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/google/wire"

	"github.com/go-kratos/kratos/v2"

	"github.com/tx7do/kratos-bootstrap/bootstrap"

	dataProviders "go-wind-cms/app/core/service/internal/data/providers"
	serverProviders "go-wind-cms/app/core/service/internal/server/providers"
	serviceProviders "go-wind-cms/app/core/service/internal/service/providers"
)

// initApp init kratos application.
func initApp(*bootstrap.Context) (*kratos.App, func(), error) {
	panic(
		wire.Build(
			serverProviders.ProviderSet,
			serviceProviders.ProviderSet,
			dataProviders.ProviderSet,
			newApp,
		),
	)
}
