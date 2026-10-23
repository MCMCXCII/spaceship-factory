package payment

import (
	"context"

	"github.com/MCMCXCII/spaceship-factory/payment/internal/api/converter"
	paymentv1 "github.com/MCMCXCII/spaceship-factory/shared/pkg/proto/payment/v1"
)

// PayOrder проводит оплату заказа и возвращает UUID транзакции
func (a *api) PayOrder(ctx context.Context, req *paymentv1.PayOrderRequest) (*paymentv1.PayOrderResponse, error) {
	orderUUID, method, err := converter.PayOrderRequestToModel(req)
	if err != nil {
		return nil, err
	}

	txUUID, err := a.paymentService.Pay(ctx, orderUUID, method)
	if err != nil {
		return nil, err
	}

	return &paymentv1.PayOrderResponse{TransactionUuid: txUUID.String()}, nil
}
