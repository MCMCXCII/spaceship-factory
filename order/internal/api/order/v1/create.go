package v1

import (
	"context"

	"github.com/MCMCXCII/spaceship-factory/order/internal/api/converter"
	orderv1 "github.com/MCMCXCII/spaceship-factory/shared/pkg/openapi/order/v1"
)

func (a *api) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (orderv1.CreateOrderRes, error) {
	order, err := a.orderService.Create(ctx, converter.CreateOrderRequestToInput(req))
	if err != nil {
		return nil, err
	}

	return converter.OrderToCreateResponse(order), nil
}
