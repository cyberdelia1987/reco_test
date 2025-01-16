package appcontext

import (
	"context"

	"go.uber.org/zap"

	"github.com/cyber/test-project/logging"
)

type key uint8

const loggerKey key = iota

func Logger(ctx context.Context) *zap.Logger {
	ctxLogger, ok := ctx.Value(loggerKey).(*zap.Logger)
	if !ok {
		return logging.Logger
	}

	return ctxLogger
}

func WithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}
