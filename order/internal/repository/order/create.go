package order

import (
	"context"

	"github.com/MCMCXCII/spaceship-factory/order/internal/model"
	"github.com/MCMCXCII/spaceship-factory/order/internal/repository/converter"
)

func (r *Repository) Create(_ context.Context, order model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[order.UUID] = converter.OrderToRecord(order)
	return nil
}
