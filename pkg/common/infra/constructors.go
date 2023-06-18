package infra

import (
	"os"
	"runtime"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"tiu-85/gorest-iman/pkg/common/infra/adapters"
	"tiu-85/gorest-iman/pkg/common/infra/clients"
	"tiu-85/gorest-iman/pkg/common/infra/providers"
	"tiu-85/gorest-iman/pkg/common/infra/repositories"
	"tiu-85/gorest-iman/pkg/common/infra/values"
)

var (
	Constructors = fx.Provide(
		providers.NewProviderByArgs,
		adapters.NewDefaultConfig,
		NewFxAppLogger,
		NewFxDB,
		clients.NewHTTPClientFactory,

		repositories.NewPostRepoFactory,
		repositories.NewTaskRepoFactory,
	)
)

type FxInDB struct {
	fx.In
	Logger adapters.Logger
	Cfg    *values.Config
}

func NewFxDB(in FxInDB) (adapters.Pool, error) {
	return adapters.NewPool(in.Cfg.DBs, in.Logger)
}

func NewFxAppLogger() (adapters.Logger, error) {
	fields := []zap.Field{
		zap.String("version", os.Getenv(adapters.AppVersionKey)),
		zap.String("app", os.Getenv(adapters.AppEnvKey)),
	}

	zapCfg := adapters.DefaultZapConfig()
	if os.Getenv("ENV") == "local" {
		zapCfg.Encoding = "console"
	}

	logger, err := adapters.NewAppLogger(zapCfg, fields...)
	if err != nil {
		return nil, err
	}
	logger.Debugf("app start cpu: %d", runtime.NumCPU())
	return logger, nil
}
