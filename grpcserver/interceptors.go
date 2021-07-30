package grpcserver

import (
	"context"
	"fmt"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/payfazz/fz-sentry/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ServerInterceptorsOptions struct {
	Logger                       *zap.Logger
	WithPanicRecovery            bool
	WithPrometheus               bool
	AdditionalUnaryInterceptors  []grpc.UnaryServerInterceptor
	AdditionalStreamInterceptors []grpc.StreamServerInterceptor
}

func ServerInterceptors(options ServerInterceptorsOptions) []grpc.ServerOption {
	recoveryOption := grpc_recovery.WithRecoveryHandlerContext(func(ctx context.Context, p interface{}) (err error) {
		logger.GetLogger(ctx).Error(
			fmt.Sprintf("panic: %v", p),
		)
		return status.Error(codes.Internal, "")
	})

	unaryInterceptors := []grpc.UnaryServerInterceptor{
		logger.GrpcUnaryServerInterceptor(options.Logger),
		logger.GrpcEndpointUnaryServerInterceptor(),
	}
	streamInterceptors := []grpc.StreamServerInterceptor{
		logger.GrpcStreamServerInterceptor(options.Logger),
		logger.GrpcEndpointStreamServerInterceptor(),
	}

	if options.WithPanicRecovery {
		unaryInterceptors = append(unaryInterceptors, grpc_recovery.UnaryServerInterceptor(recoveryOption))
		streamInterceptors = append(streamInterceptors, grpc_recovery.StreamServerInterceptor(recoveryOption))
	}

	if options.WithPrometheus {
		unaryInterceptors = append(unaryInterceptors, grpc_prometheus.UnaryServerInterceptor)
		streamInterceptors = append(streamInterceptors, grpc_prometheus.StreamServerInterceptor)
	}

	unaryInterceptors = append(unaryInterceptors, options.AdditionalUnaryInterceptors...)
	streamInterceptors = append(streamInterceptors, options.AdditionalStreamInterceptors...)

	return []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(unaryInterceptors...),
		grpc.ChainStreamInterceptor(streamInterceptors...),
	}
}
