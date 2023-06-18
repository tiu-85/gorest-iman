package main

import (
	"context"

	"github.com/gin-gonic/gin"

	gi "tiu-85/gorest-iman/pkg/api-gateway/infra"
	"tiu-85/gorest-iman/pkg/api-gateway/infra/services"
	"tiu-85/gorest-iman/pkg/common/application"
	ci "tiu-85/gorest-iman/pkg/common/infra"
	"tiu-85/gorest-iman/pkg/common/infra/values"
)

func main() {
	a := application.NewApp(
		ci.Constructors,
		gi.Constructors,
	)
	a.Demonize(func(
		ctx context.Context,
		cfg *values.Config,

		postCrudClientService *services.PostCrudClientService,
		postFetchClientService *services.PostFetchClientService,
	) error {
		engine := gin.Default()

		_ = postCrudClientService.RegisterRoutes(engine)
		_ = postFetchClientService.RegisterRoutes(engine)

		return engine.Run(cfg.GatewayPort)
	})
}
