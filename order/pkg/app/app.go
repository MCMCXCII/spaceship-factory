package app

import (
	"net/http"

	orderAPI "github.com/MCMCXCII/spaceship-factory/order/internal/api/order/v1"
	inventoryClient "github.com/MCMCXCII/spaceship-factory/order/internal/client/grpc/inventory/v1"
	paymentClient "github.com/MCMCXCII/spaceship-factory/order/internal/client/grpc/payment/v1"
	orderRepository "github.com/MCMCXCII/spaceship-factory/order/internal/repository/order"
	orderService "github.com/MCMCXCII/spaceship-factory/order/internal/service/order"
	orderv1 "github.com/MCMCXCII/spaceship-factory/shared/pkg/openapi/order/v1"
	inventoryv1 "github.com/MCMCXCII/spaceship-factory/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/MCMCXCII/spaceship-factory/shared/pkg/proto/payment/v1"
)

// NewHTTPHandler собирает слои OrderService и возвращает HTTP-обработчик ogen
func NewHTTPHandler(
	inventory inventoryv1.InventoryServiceClient,
	payment paymentv1.PaymentServiceClient,
) (http.Handler, error) {
	repo := orderRepository.New()
	service := orderService.New(
		repo,
		inventoryClient.New(inventory),
		paymentClient.New(payment),
	)
	api := orderAPI.NewAPI(service)

	return orderv1.NewServer(api, orderv1.WithErrorHandler(orderAPI.ErrorHandler))
}
