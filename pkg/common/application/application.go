package application

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/jessevdk/go-flags"

	"go.uber.org/fx"

	"tiu-85/gorest-iman/pkg/common/infra/adapters"
	"tiu-85/gorest-iman/pkg/common/infra/values"
)

type App struct {
	ctx     context.Context
	cancel  context.CancelFunc
	fxApp   *fx.App
	logger  adapters.Logger
	options []fx.Option
}

func NewApp(providers ...interface{}) *App {
	ctx, cancel := context.WithCancel(context.Background())
	zapCfg := adapters.DefaultZapConfig()
	if os.Getenv("ENV") == "local" {
		zapCfg.Encoding = "console"
	}
	preInitLogger, err := adapters.NewAppLogger(zapCfg)
	if err != nil {
		panic(err)
	}
	app := &App{
		options: []fx.Option{fx.NopLogger},
		ctx:     ctx,
		cancel:  cancel,
		logger:  preInitLogger,
	}

	for _, provider := range providers {
		switch p := provider.(type) {
		case fx.Option:
			app.options = append(app.options, p)
		default:
			app.options = append(app.options, fx.Provide(p))
		}
	}
	app.options = append(app.options, fx.Provide(func() context.Context {
		return ctx
	}))

	return app
}

func (a *App) Run(processes ...any) {
	var args values.Args
	_, err := flags.Parse(&args)
	if err != nil {
		a.logger.Fatal(err)
	}

	a.options = append(
		a.options,
		fx.Provide(func() values.Args {
			return args
		}),
		fx.Invoke(processes...),
	)

	if args.RunDry {
		err = fx.ValidateApp(a.options...)
		if err != nil {
			graph, er := fx.VisualizeError(err)
			if er != nil {
				a.logger.Fatalf("%v %v", err, graph)
			} else {
				a.logger.Fatal(err)
			}
		}
		os.Exit(0)
	}

	a.fxApp = fx.New(a.options...)
	go a.listenSignals()
	startCtx, cancel := context.WithTimeout(a.ctx, fx.DefaultTimeout)
	defer cancel()
	err = a.fxApp.Start(startCtx)
	if err != nil {
		a.logger.Fatal(err)
	}
}

func (a *App) Demonize(demons ...any) {
	a.Run(demons...)
	<-a.ctx.Done()
}

func (a *App) Stop() {
	stopCtx, cancel := context.WithTimeout(a.ctx, fx.DefaultTimeout)
	defer cancel()
	err := a.fxApp.Stop(stopCtx)
	if err != nil {
		a.logger.Fatal(err)
	}
	a.fxApp = nil
}

func (a *App) listenSignals() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	for sig := range signals {
		a.logger.Infof("income signal %s", sig)
		a.Stop()
		a.cancel()
		return
	}
}
