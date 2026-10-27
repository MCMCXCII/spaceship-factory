package v1

import (
	"context"

	"github.com/MCMCXCII/spaceship-factory/order/internal/api/converter"
	orderv1 "github.com/MCMCXCII/spaceship-factory/shared/pkg/openapi/order/v1"
)

func (a *api) PayOrder(ctx context.Context, req *orderv1.PayOrderRequest, params orderv1.PayOrderParams) (orderv1.PayOrderRes, error) {
	method, err := converter.PaymentMethodToModel(req.PaymentMethod)
	if err != nil {
		return nil, err
	}

	txUUID, err := a.orderService.Pay(ctx, params.OrderUUID, method)
	if err != nil {
		return nil, err
	}

	return &orderv1.PayOrderResponse{TransactionUUID: txUUID}, nil
}
