package order

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	errs "github.com/MCMCXCII/spaceship-factory/order/internal/errors"
	"github.com/MCMCXCII/spaceship-factory/order/internal/model"
	"github.com/MCMCXCII/spaceship-factory/order/internal/service/input"
)

func (s *service) Create(ctx context.Context, in input.CreateOrderInput) (model.Order, error) {
	parts, err := s.inventoryClient.ListParts(ctx, in.PartUUIDs())
	if err != nil {
		return model.Order{}, fmt.Errorf("получить детали: %w", err)
	}

	var totalPrice int64
	for _, part := range parts {
		if part.StockQuantity <= 0 {
			return model.Order{}, fmt.Errorf("деталь %s: %w", part.Name, errs.ErrOutOfStock)
		}
		totalPrice += part.Price
	}

	order := model.Order{
		UUID:       uuid.New(),
		HullUUID:   in.HullUUID,
		EngineUUID: in.EngineUUID,
		ShieldUUID: in.ShieldUUID,
		WeaponUUID: in.WeaponUUID,
		TotalPrice: totalPrice,
		Status:     model.OrderStatusPendingPayment,
		CreatedAt:  time.Now(),
	}

	if err := s.orderRepository.Create(ctx, order); err != nil {
		return model.Order{}, fmt.Errorf("сохранить заказ: %w", err)
	}

	return order, nil
}
