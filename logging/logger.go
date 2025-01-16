package logging

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/cyber/test-project/config"
)

func defaultLogger() *zap.Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	return logger
}

var Logger = defaultLogger()

func Init(settings config.LoggingConfig) error {
	level := zapcore.InfoLevel

	err := level.Set(settings.Level)
	if err != nil {
		Logger.Error("Failed to set log level", zap.Error(err), zap.String("level", settings.Level))
	}

	atomicLevel := zap.NewAtomicLevelAt(level)
	zapConfig := zap.NewProductionConfig()

	zapConfig.Level = atomicLevel
	zapConfig.OutputPaths = settings.Output
	zapConfig.DisableStacktrace = !settings.LogStackTrace
	zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	Logger, err = zapConfig.Build()
	if err != nil {
		return err
	}

	return nil
}

type key int

const (
	loggerKey key = iota
)

func FromContext(ctx context.Context) *zap.Logger {
	ctxLogger, ok := ctx.Value(loggerKey).(*zap.Logger)
	if !ok {
		return Logger
	}

	return ctxLogger
}

func WithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

type fieldGetter func() zapcore.Field

func DebugField(fieldGetter fieldGetter) zapcore.Field {
	if !Logger.Core().Enabled(zapcore.DebugLevel) {
		return zap.Skip()
	}

	return fieldGetter()
}
