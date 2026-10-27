package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/MCMCXCII/spaceship-factory/order/internal/model"
)

func (s *service) Cancel(ctx context.Context, orderUUID uuid.UUID) error {
	order, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		return fmt.Errorf("получить заказ: %w", err)
	}

	if err := ensurePendingPayment(order); err != nil {
		return err
	}

	order.Status = model.OrderStatusCancelled

	if err := s.orderRepository.Update(ctx, order); err != nil {
		return fmt.Errorf("обновить заказ: %w", err)
	}

	return nil
}
