package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	errs "github.com/MCMCXCII/spaceship-factory/order/internal/errors"
	"github.com/MCMCXCII/spaceship-factory/order/internal/model"
)

func (s *service) Pay(ctx context.Context, orderUUID uuid.UUID, method model.PaymentMethod) (uuid.UUID, error) {
	order, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("получить заказ: %w", err)
	}

	if err := ensurePendingPayment(order); err != nil {
		return uuid.Nil, err
	}

	txUUID, err := s.paymentClient.PayOrder(ctx, order.UUID, method)
	if err != nil {
		return uuid.Nil, fmt.Errorf("оплатить заказ: %w", err)
	}

	order.Status = model.OrderStatusPaid
	order.TransactionUUID = uuid.NullUUID{UUID: txUUID, Valid: true}
	order.PaymentMethod = method

	if err := s.orderRepository.Update(ctx, order); err != nil {
		return uuid.Nil, fmt.Errorf("обновить заказ: %w", err)
	}

	return txUUID, nil
}

func ensurePendingPayment(order model.Order) error {
	switch order.Status {
	case model.OrderStatusPaid:
		return errs.ErrOrderAlreadyPaid
	case model.OrderStatusCancelled:
		return errs.ErrOrderCancelled
	default:
		return nil
	}
}
