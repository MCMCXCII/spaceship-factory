package payment

import (
	paymentv1 "github.com/MCMCXCII/spaceship-factory/shared/pkg/proto/payment/v1"
)

// api — gRPC-обработчики PaymentService
type api struct {
	paymentv1.UnimplementedPaymentServiceServer
	paymentService PaymentService
}

// NewAPI создаёт gRPC-обработчики PaymentService
func NewAPI(paymentService PaymentService) *api {
	return &api{paymentService: paymentService}
}
