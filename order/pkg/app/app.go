package app

import (
	"fmt"
	"net/http"

	orderAPI "github.com/MCMCXCII/spaceship-factory/order/internal/api/order/v1"
	inventoryClient "github.com/MCMCXCII/spaceship-factory/order/internal/client/grpc/inventory/v1"
	paymentClient "github.com/MCMCXCII/spaceship-factory/order/internal/client/grpc/payment/v1"
	orderRepository "github.com/MCMCXCII/spaceship-factory/order/internal/repository/order"
	orderService "github.com/MCMCXCII/spaceship-factory/order/internal/service/order"
	orderv1 "github.com/MCMCXCII/spaceship-factory/shared/pkg/openapi/order/v1"
	inventoryv1 "github.com/MCMCXCII/spaceship-factory/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/MCMCXCII/spaceship-factory/shared/pkg/proto/payment/v1"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NewHTTPHandler собирает слои OrderService и возвращает HTTP-обработчик ogen
func NewHTTPHandler(
	pool *pgxpool.Pool,
	inventory inventoryv1.InventoryServiceClient,
	payment paymentv1.PaymentServiceClient,
) (http.Handler, error) {
	txManager, err := manager.New(trmpgx.NewDefaultFactory(pool))
	if err != nil {
		return nil, fmt.Errorf("создать transaction manager: %w", err)
	}

	repo := orderRepository.New(pool, txManager)
	service := orderService.New(
		repo,
		inventoryClient.New(inventory),
		paymentClient.New(payment),
	)
	api := orderAPI.NewAPI(service)

	return orderv1.NewServer(api, orderv1.WithErrorHandler(orderAPI.ErrorHandler))
}
