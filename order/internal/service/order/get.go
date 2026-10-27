package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/MCMCXCII/spaceship-factory/order/internal/model"
)

func (s *service) Get(ctx context.Context, orderUUID uuid.UUID) (model.Order, error) {
	order, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		return model.Order{}, fmt.Errorf("получить заказ: %w", err)
	}

	return order, nil
}
