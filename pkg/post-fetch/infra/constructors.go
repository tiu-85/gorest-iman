package infra

import (
	"go.uber.org/fx"

	"tiu-85/gorest-iman/pkg/post-fetch/infra/clients"
	"tiu-85/gorest-iman/pkg/post-fetch/infra/processors"
	"tiu-85/gorest-iman/pkg/post-fetch/infra/services"
)

var (
	Constructors = fx.Provide(
		clients.NewPostApiClient,
		processors.NewPostApiProcessor,
		services.NewPostFetchService,
	)
)
