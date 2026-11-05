package order

import (
	"context"
	"fmt"
	"time"

	errs "github.com/MCMCXCII/spaceship-factory/order/internal/errors"
	"github.com/MCMCXCII/spaceship-factory/order/internal/model"
	"github.com/MCMCXCII/spaceship-factory/order/internal/repository/converter"
)

func (r *repository) Update(ctx context.Context, order model.Order) error {
	return r.txManager.Do(ctx, func(txCtx context.Context) error {
		rec := converter.OrderToRecord(order)

		const query = `
			UPDATE orders
			SET status = $1, transaction_uuid = $2, payment_method = $3, updated_at = $4
			WHERE uuid = $5`

		tag, err := r.getter.DefaultTrOrDB(txCtx, r.pool).Exec(
			txCtx, query,
			rec.Status,
			rec.TransactionUUID,
			rec.PaymentMethod,
			time.Now().UTC(),
			rec.UUID,
		)
		if err != nil {
			return fmt.Errorf("обновить заказ: %w", err)
		}

		if tag.RowsAffected() == 0 {
			return errs.ErrOrderNotFound
		}

		return nil
	})
}
