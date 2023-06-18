package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	AppEnvKey     = "APP_NAME"
	AppVersionKey = "APP_VERSION"
)

func DefaultZapConfig(lvl ...string) zap.Config {
	level := zap.NewAtomicLevel()
	err := level.UnmarshalText([]byte(strings.Join(lvl, "")))
	if err != nil || len(lvl) == 0 {
		level.SetLevel(zap.DebugLevel)
	}
	return zap.Config{
		Level:            level,
		DisableCaller:    false,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		Encoding:         "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "msg",
			LevelKey:       "level",
			NameKey:        "logger",
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			CallerKey:      "file",
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeName:     zapcore.FullNameEncoder,
			TimeKey:        "time",
		},
	}
}

type (
	Logger interface {
		Debug(args ...interface{})
		Info(args ...interface{})
		Warn(args ...interface{})
		Error(args ...interface{})
		Fatal(args ...interface{})
		Debugf(format string, args ...interface{})
		Infof(format string, args ...interface{})
		Printf(format string, args ...interface{})
		Warnf(format string, args ...interface{})
		Errorf(format string, args ...interface{})
		Fatalf(format string, args ...interface{})
		With(args ...interface{}) Logger
		WithCtx(ctx context.Context, args ...interface{}) Logger
		Infow(string, ...interface{})
		Named(string) Logger
		Print(v ...interface{})
		Println(v ...interface{})
		GetZapLogger() *zap.Logger
		SetZapLogger(*zap.Logger) Logger
		ReportOfCall()
		Report(data interface{})
		Reportf(format string, args ...interface{})
	}

	appLogger struct {
		*zap.Config
		*zap.SugaredLogger
	}
)

func (l *appLogger) GetZapLogger() *zap.Logger {
	return l.SugaredLogger.Desugar()
}

func (l *appLogger) SetZapLogger(logger *zap.Logger) Logger {
	l.SugaredLogger = logger.Sugar()
	return l
}

func NewAppLogger(
	cfg zap.Config,
	fields ...zap.Field,
) (Logger, error) {
	baseLogger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	logger := baseLogger.WithOptions(
		zap.Fields(fields...),
	).Sugar()

	defer logger.Sync()

	return &appLogger{
		Config:        &cfg,
		SugaredLogger: logger,
	}, nil
}

func (l *appLogger) With(args ...interface{}) Logger {
	return &appLogger{
		Config:        l.Config,
		SugaredLogger: l.SugaredLogger.With(args...),
	}
}

func (l *appLogger) WithCtx(ctx context.Context, args ...interface{}) Logger {
	logger := &appLogger{
		Config:        l.Config,
		SugaredLogger: l.SugaredLogger.With(args...),
	}

	return logger
}

func (l *appLogger) Named(name string) Logger {
	return &appLogger{
		Config:        l.Config,
		SugaredLogger: l.SugaredLogger.Named(name),
	}
}

func (l *appLogger) Printf(format string, args ...interface{}) {
	l.Debugf(format, args...)
}

func (l *appLogger) Print(v ...interface{}) {
	l.Debug(v...)
}

func (l *appLogger) Println(v ...interface{}) {
	l.Debug(v...)
}

func (l *appLogger) ReportOfCall() {
	pc, _, _, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		methodNameParts := strings.Split(details.Name(), "/")
		if len(methodNameParts) == 0 {
			return
		}

		methodNameParts = strings.Split(methodNameParts[len(methodNameParts)-1], ".")
		if len(methodNameParts) != 3 {
			return
		}

		methodNameParts[1] = strings.Trim(methodNameParts[1], "(*)")
		l.Reportf("call %s", strings.ToTitle(strings.Join(methodNameParts[1:], ".")))
	}
}

func (l *appLogger) Reportf(format string, args ...interface{}) {
	l.Report(fmt.Sprintf(format, args...))
}

func (l *appLogger) Report(data interface{}) {
	logger := l.Named("report")
	switch d := data.(type) {
	case string, []byte, int, uint32, int64:
		logger.Debug(d)
	default:
		buf, err := json.Marshal(d)
		if err != nil {
			l.Debug(err)
			logger.Debug(d)
			return
		}
		logger.Debug(string(buf))
	}
}

func NewDefaultLogger(lvl ...string) Logger {
	l, err := NewAppLogger(DefaultZapConfig(lvl...))
	if err != nil {
		panic(err)
	}
	return l
}
