package infra

import (
	"go.uber.org/fx"

	"tiu-85/gorest-iman/pkg/api-gateway/infra/services"
)

var (
	Constructors = fx.Provide(
		services.NewPostCrudServiceClient,
		services.NewPostFetchServiceClient,
	)
)
