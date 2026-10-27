package payment

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	"github.com/MCMCXCII/spaceship-factory/payment/internal/model"
)

type service struct{}

func New() *service {
	return &service{}
}

// Pay проводит оплату заказа и возвращает UUID новой транзакции.
func (s *service) Pay(ctx context.Context, orderUUID uuid.UUID, method model.PaymentMethod) (uuid.UUID, error) {
	txUUID := uuid.New()

	slog.InfoContext(
		ctx, "оплата выполнена",
		"order_uuid", orderUUID,
		"transaction_uuid", txUUID,
		"payment_method", method,
	)

	return txUUID, nil
}
