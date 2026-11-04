package payment

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	errs "github.com/MCMCXCII/spaceship-factory/payment/internal/errors"
	"github.com/MCMCXCII/spaceship-factory/payment/internal/service/input"
)

type service struct{}

func New() *service {
	return &service{}
}

// Pay проводит оплату заказа и возвращает UUID новой транзакции.
func (s *service) Pay(ctx context.Context, in input.PayOrderInput) (uuid.UUID, error) {
	if !in.PaymentMethod.IsValid() {
		return uuid.Nil, errs.ErrInvalidPaymentMethod
	}

	txUUID := uuid.New()

	slog.InfoContext(
		ctx, "оплата выполнена",
		"order_uuid", in.OrderUUID,
		"transaction_uuid", txUUID,
		"payment_method", in.PaymentMethod,
	)

	return txUUID, nil
}
