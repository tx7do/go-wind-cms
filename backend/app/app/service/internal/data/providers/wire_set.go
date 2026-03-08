//go:build wireinject
// +build wireinject

//go:generate go run github.com/google/wire/cmd/wire

// This file defines the dependency injection ProviderSet for the data layer and contains no business logic.
// The build tag `wireinject` excludes this source from normal `go build` and final binaries.
// Run `go generate ./...` or `go run github.com/google/wire/cmd/wire` to regenerate the Wire output (e.g. `wire_gen.go`), which will be included in final builds.
// Keep provider constructors here only; avoid init-time side effects or runtime logic in this file.

package providers

import (
	"github.com/google/wire"

	"go-wind-cms/app/app/service/internal/data"

	"go-wind-cms/pkg/middleware/auth"
)

// ProviderSet is the Wire provider set for data layer.
var ProviderSet = wire.NewSet(
	data.NewRedisClient,
	data.NewMinIoClient,
	data.NewDiscovery,

	data.NewClientType,
	data.NewAuthorizer,

	auth.NewTokenChecker,

	data.NewAuthenticationServiceClient,
	data.NewUserCredentialServiceClient,

	data.NewFileServiceClient,

	data.NewUserServiceClient,
	data.NewTenantServiceClient,
	data.NewRoleServiceClient,
	data.NewOrgUnitServiceClient,
	data.NewPositionServiceClient,

	data.NewPageServiceClient,
	data.NewCategoryServiceClient,
	data.NewPostServiceClient,
	data.NewTagServiceClient,

	data.NewCommentServiceClient,

	data.NewNavigationServiceClient,
	data.NewSiteSettingServiceClient,

	data.NewMediaAssetServiceClient,
)
