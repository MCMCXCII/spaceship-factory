package v1

import (
	"context"

	"github.com/MCMCXCII/spaceship-factory/order/internal/api/converter"
	orderv1 "github.com/MCMCXCII/spaceship-factory/shared/pkg/openapi/order/v1"
)

func (a *api) GetOrder(ctx context.Context, params orderv1.GetOrderParams) (orderv1.GetOrderRes, error) {
	order, err := a.orderService.Get(ctx, params.OrderUUID)
	if err != nil {
		return nil, err
	}

	return converter.OrderToDTO(order), nil
}
