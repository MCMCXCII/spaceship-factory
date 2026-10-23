package interceptor

import (
	"context"
	"errors"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	errs "github.com/MCMCXCII/spaceship-factory/payment/internal/errors"
)

// ErrorInterceptor переводит доменные ошибки хендлеров в gRPC-статусы.
func ErrorInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	resp, err := handler(ctx, req)
	if err == nil {
		return resp, nil
	}

	switch {
	case errors.Is(err, errs.ErrInvalidUUID),
		errors.Is(err, errs.ErrInvalidPaymentMethod):
		return nil, status.Error(codes.InvalidArgument, err.Error())
	default:
		slog.ErrorContext(ctx, "необработанная ошибка gRPC-метода",
			"method", info.FullMethod,
			"error", err,
		)
		return nil, status.Error(codes.Internal, "внутренняя ошибка")
	}
}
