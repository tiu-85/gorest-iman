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

type PostCrudClientService struct {
	logger adapters.Logger
	cfg    *values.Config
	client pbv1.PostCrudServiceClient
}

func (c *PostCrudClientService) RegisterRoutes(engine *gin.Engine) error {
	routesGroup := engine.Group("/v1")

	routesGroup.POST("/posts", func(ctx *gin.Context) {
		routes.GetPostList(ctx, c.client)
	})
	routesGroup.GET("/posts/:id", func(ctx *gin.Context) {
		routes.GetPost(ctx, c.client)
	})
	routesGroup.PATCH("/posts/:id", func(ctx *gin.Context) {
		routes.EditPost(ctx, c.client)
	})
	routesGroup.DELETE("/posts/:id", func(ctx *gin.Context) {
		routes.DeletePost(ctx, c.client)
	})

	return nil
}

func NewPostCrudServiceClient(
	logger adapters.Logger,
	cfg *values.Config,
) *PostCrudClientService {
	connection, err := grpc.Dial(cfg.PostCrudServiceUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		logger.Errorf("could not connect to post-crud server: ", err)
		return nil
	}

	logger.Infof("successfully connected to post-crud server: %+v", connection)

	return &PostCrudClientService{
		logger: logger,
		cfg:    cfg,
		client: pbv1.NewPostCrudServiceClient(connection),
	}
}
