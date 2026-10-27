package app

import (
	"google.golang.org/grpc"

	paymentAPI "github.com/MCMCXCII/spaceship-factory/payment/internal/api/payment/v1"
	"github.com/MCMCXCII/spaceship-factory/payment/internal/interceptor"
	paymentService "github.com/MCMCXCII/spaceship-factory/payment/internal/service/payment"
	paymentv1 "github.com/MCMCXCII/spaceship-factory/shared/pkg/proto/payment/v1"
)

// Interceptors возвращает опции gRPC-сервера с interceptor'ами PaymentService
func Interceptors() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.UnaryInterceptor(interceptor.ErrorInterceptor),
	}
}

// RegisterServices собирает слои PaymentService и регистрирует его в gRPC-сервере
func RegisterServices(s grpc.ServiceRegistrar) {
	service := paymentService.New()
	api := paymentAPI.NewAPI(service)

	paymentv1.RegisterPaymentServiceServer(s, api)
}
