package servergrpc

import (
	"context"
	"time"

	"github.com/juliazadorozhnaya/otus_homework/hw12_13_14_15_calendar/internal/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// LoggingInterceptor - middleware gRPC интерцептор для логирования запросов и ответов.
func LoggingInterceptor(logger server.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		logger.Info("gRPC call started", "method", info.FullMethod, "request", req)

		resp, err := handler(ctx, req)

		logger.Info("gRPC call finished", "method", info.FullMethod, "request", req,
			"response", resp, "error", err, "duration", time.Since(start))

		if err != nil {
			st, _ := status.FromError(err)
			logger.Error("gRPC call error",
				"method", info.FullMethod, "error", st.Message(), "code", st.Code())
		}

		return resp, err
	}
}
