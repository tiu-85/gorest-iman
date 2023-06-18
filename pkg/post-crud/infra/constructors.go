package infra

import (
	"go.uber.org/fx"

	"tiu-85/gorest-iman/pkg/post-crud/infra/services"
)

var (
	Constructors = fx.Provide(
		services.NewPostCrudService,
	)
)
