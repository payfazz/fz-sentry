package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc/status"
	"path"
	"runtime"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const (
	MAX_CALLER    = 32 // reserve for max 32 caller stack
	CALLER_OFFSET = 6  // 6 default stack: runtime.goexit, serveStream, handleStream, processUnaryGRPC, grpc.handler
)

func GrpcMiddleware(logger *zap.Logger) endpoint.Middleware {
	return func(f endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, in interface{}) (out interface{}, err error) {
			newCtx := NewRequest(ctx, logger)
			return f(newCtx, in)
		}
	}
}

func GrpcUnaryServerInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		newCtx := NewRequest(ctx, logger)
		return handler(newCtx, req)
	}
}

func GrpcStreamServerInterceptor(logger *zap.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		newCtx := NewRequest(ss.Context(), logger)
		wrappedStream := grpc_middleware.WrapServerStream(ss)
		wrappedStream.WrappedContext = newCtx
		return handler(srv, wrappedStream)
	}
}

func GrpcEndpointUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		log := GetLogger(ctx)
		service := path.Dir(info.FullMethod)[1:]
		method := path.Base(info.FullMethod)

		log.Info(fmt.Sprintf("begin grpc request: %s/%s", service, method),
			zap.String("service", service),
			zap.String("method", method),
		)
		start := time.Now()
		resp, err = handler(ctx, req)
		elapsed := time.Since(start)
		code := status.Code(err).String()

		log.Info(fmt.Sprintf("end grpc request: %s", elapsed),
			zap.Duration("elapsed", elapsed),
			zap.String("service", service),
			zap.String("method", method),
			zap.String("code", code),
		)

		return resp, err
	}
}

func GrpcEndpointStreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		log := GetLogger(ss.Context())
		service := path.Dir(info.FullMethod)[1:]
		method := path.Base(info.FullMethod)

		log.Info(fmt.Sprintf("begin grpc request: %s/%s", service, method),
			zap.String("service", service),
			zap.String("method", method),
		)
		start := time.Now()
		err := handler(srv, ss)
		elapsed := time.Since(start)
		code := status.Code(err).String()

		log.Info(fmt.Sprintf("end grpc request: %s", elapsed),
			zap.Duration("elapsed", elapsed),
			zap.String("service", service),
			zap.String("method", method),
			zap.String("code", code),
		)

		return err
	}
}

// GrpcEndpointMiddleware is middleware for grpc request
//
// Deprecated: funcName will be wrong if using grpc interceptors, please use GrpcEndpointUnaryServerInterceptor and
// GrpcEndpointStreamServerInterceptor instead
func GrpcEndpointMiddleware() endpoint.Middleware {
	return func(f endpoint.Endpoint) endpoint.Endpoint {
		var start time.Time
		return DoGRPC(
			f,
			func(ctx context.Context, log *zap.Logger, in interface{}) error {
				pcs := make([]uintptr, MAX_CALLER)
				n := runtime.Callers(0, pcs)

				funcName := runtime.FuncForPC(pcs[n-CALLER_OFFSET]).Name()
				log.Info(fmt.Sprintf("begin grpc request: %s", funcName))
				start = time.Now()
				return nil
			},
			func(ctx context.Context, log *zap.Logger, out interface{}) error {
				elapsed := time.Since(start)
				log.Info(fmt.Sprintf("end grpc request: %s", elapsed))
				return nil
			},
		)
	}
}

func GrpcRequestMiddleware() endpoint.Middleware {
	return func(f endpoint.Endpoint) endpoint.Endpoint {
		return DoGRPC(
			f,
			func(ctx context.Context, log *zap.Logger, in interface{}) error {
				body, _ := json.Marshal(in)
				log.Debug("grpc request payload",
					zap.String("payload", string(body)),
				)
				return nil
			},
			nil,
		)
	}
}

func GrpcResponseMiddleware() endpoint.Middleware {
	return func(f endpoint.Endpoint) endpoint.Endpoint {
		return DoGRPC(
			f,
			nil,
			func(ctx context.Context, log *zap.Logger, out interface{}) error {
				resp, _ := json.Marshal(out)
				log.Debug("grpc response payload",
					zap.String("payload", string(resp)),
				)
				return nil
			},
		)
	}
}
