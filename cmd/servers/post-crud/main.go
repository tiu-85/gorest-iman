package main

import (
	"context"
	"net"

	"google.golang.org/grpc"

	gi "tiu-85/gorest-iman/pkg/api-gateway/infra"
	"tiu-85/gorest-iman/pkg/common/application"
	"tiu-85/gorest-iman/pkg/common/gen/pbv1"
	ci "tiu-85/gorest-iman/pkg/common/infra"
	"tiu-85/gorest-iman/pkg/common/infra/adapters"
	"tiu-85/gorest-iman/pkg/common/infra/values"
	"tiu-85/gorest-iman/pkg/post-crud/infra"
	"tiu-85/gorest-iman/pkg/post-crud/infra/services"
)

func main() {
	a := application.NewApp(
		ci.Constructors,
		gi.Constructors,
		infra.Constructors,
	)
	a.Demonize(func(
		logger adapters.Logger,
		ctx context.Context,
		cfg *values.Config,

		postCrudService *services.PostCrudService,
	) error {
		lis, err := net.Listen("tcp", cfg.PostCrudServicePort)
		if err != nil {
			logger.Errorf("failed to listening: ", err)
			return err
		}

		logger.Infof("post-crud service is listening on port: %s", cfg.PostCrudServicePort)

		grpcServer := grpc.NewServer()
		pbv1.RegisterPostCrudServiceServer(grpcServer, postCrudService)

		if err := grpcServer.Serve(lis); err != nil {
			logger.Errorf("failed to serve: ", err)
			return err
		}

		return nil
	})
}
