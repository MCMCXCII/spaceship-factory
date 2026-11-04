package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	errs "github.com/MCMCXCII/spaceship-factory/order/internal/errors"
	"github.com/MCMCXCII/spaceship-factory/order/internal/model"
	"github.com/MCMCXCII/spaceship-factory/order/internal/repository/converter"
	"github.com/MCMCXCII/spaceship-factory/order/internal/repository/record"
)

func (r *repository) Get(ctx context.Context, orderUUID uuid.UUID) (model.Order, error) {
	const query = `
		SELECT o.uuid, o.status, o.transaction_uuid, o.payment_method, o.created_at, o.updated_at,
		       oi.part_uuid, oi.part_type, oi.price
		FROM orders o
		LEFT JOIN order_items oi ON oi.order_uuid = o.uuid
		WHERE o.uuid = $1
		ORDER BY oi.part_type`

	rows, err := r.getter.DefaultTrOrDB(ctx, r.pool).Query(ctx, query, orderUUID)
	if err != nil {
		return model.Order{}, fmt.Errorf("получить заказ: %w", err)
	}
	defer rows.Close()

	var (
		rec   record.Order
		items []record.OrderItem
		found bool
	)

	for rows.Next() {
		var (
			partUUID *uuid.UUID
			partType *string
			price    *int64
		)

		if err := rows.Scan(
			&rec.UUID, &rec.Status, &rec.TransactionUUID, &rec.PaymentMethod, &rec.CreatedAt, &rec.UpdatedAt,
			&partUUID, &partType, &price,
		); err != nil {
			return model.Order{}, fmt.Errorf("прочитать заказ: %w", err)
		}
		found = true

		if partUUID == nil {
			continue
		}

		items = append(items, record.OrderItem{
			OrderUUID: rec.UUID,
			PartUUID:  *partUUID,
			PartType:  *partType,
			Price:     *price,
		})
	}

	if err := rows.Err(); err != nil {
		return model.Order{}, fmt.Errorf("прочитать заказ: %w", err)
	}

	if !found {
		return model.Order{}, errs.ErrOrderNotFound
	}

	return converter.OrderToModel(rec, items), nil
}
