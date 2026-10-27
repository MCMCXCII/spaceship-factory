package v1

import orderv1 "github.com/MCMCXCII/spaceship-factory/shared/pkg/openapi/order/v1"

type api struct {
	orderv1.UnimplementedHandler
	orderService OrderService
}

// NewAPI создаёт HTTP-обработчики OrderService
func NewAPI(orderService OrderService) *api {
	return &api{orderService: orderService}
}
