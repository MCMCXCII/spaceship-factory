package order

import (
	"context"

	errs "github.com/MCMCXCII/spaceship-factory/order/internal/errors"
	"github.com/MCMCXCII/spaceship-factory/order/internal/model"
	"github.com/MCMCXCII/spaceship-factory/order/internal/repository/converter"
)

func (r *Repository) Update(_ context.Context, order model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[order.UUID]; !ok {
		return errs.ErrOrderNotFound
	}

	r.data[order.UUID] = converter.OrderToRecord(order)
	return nil
}
