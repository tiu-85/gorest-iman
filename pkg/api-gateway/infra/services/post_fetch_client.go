package services

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"tiu-85/gorest-iman/pkg/api-gateway/infra/routes"
	"tiu-85/gorest-iman/pkg/common/gen/pbv1"
	"tiu-85/gorest-iman/pkg/common/infra/adapters"
	"tiu-85/gorest-iman/pkg/common/infra/values"
)

type PostFetchClientService struct {
	logger adapters.Logger
	cfg    *values.Config
	client pbv1.PostFetchServiceClient
}

func (c *PostFetchClientService) RegisterRoutes(engine *gin.Engine) error {
	routesGroup := engine.Group("/v1")

	routesGroup.POST("/tasks", func(ctx *gin.Context) {
		routes.CreateTask(ctx, c.client)
	})
	routesGroup.GET("/tasks/:id", func(ctx *gin.Context) {
		routes.GetTask(ctx, c.client)
	})

	return nil
}

func NewPostFetchServiceClient(
	logger adapters.Logger,
	cfg *values.Config,
) *PostFetchClientService {
	connection, err := grpc.Dial(cfg.PostFetchServiceUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		logger.Errorf("could not connect to post-fetch server: ", err)
		return nil
	}

	logger.Infof("successfully connected to post-fetch server: %+v", connection)

	return &PostFetchClientService{
		logger: logger,
		cfg:    cfg,
		client: pbv1.NewPostFetchServiceClient(connection),
	}
}
