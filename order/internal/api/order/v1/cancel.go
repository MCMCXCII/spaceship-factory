package v1

import (
	"context"

	orderv1 "github.com/MCMCXCII/spaceship-factory/shared/pkg/openapi/order/v1"
)

func (a *api) CancelOrder(ctx context.Context, params orderv1.CancelOrderParams) (orderv1.CancelOrderRes, error) {
	if err := a.orderService.Cancel(ctx, params.OrderUUID); err != nil {
		return nil, err
	}

	return &orderv1.CancelOrderResponse{}, nil
}
