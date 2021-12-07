package logger

import (
	"context"
	"github.com/gofrs/uuid"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type loggerKeyType struct{}

var loggerKey loggerKeyType

func GetLogger(ctx context.Context) *zap.Logger {
	logger := ctx.Value(loggerKey).(*zap.Logger)
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasSpanID() {
		logger = logger.With(
			zap.String("spanId", spanCtx.SpanID().String()),
		)
	}
	return logger
}

func NewRequest(ctx context.Context, logger *zap.Logger) context.Context {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasTraceID() {
		logger = logger.With(
			zap.String("traceId", spanCtx.TraceID().String()),
		)
	}
	if !spanCtx.HasTraceID() {
		requestId, _ := uuid.NewV4()
		logger = logger.With(
			zap.String("requestId", requestId.String()),
		)
	}
	return context.WithValue(ctx, loggerKey, logger)
}

func New(env string, service string, options ...zap.Option) *zap.Logger {
	cfg := zap.NewDevelopmentConfig()
	if "development" != env {
		cfg = zap.NewProductionConfig()
	}

	cfg.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	logger, _ := cfg.Build(options...)

	logger = logger.With(
		zap.String("serviceName", service),
	)

	return logger
}
