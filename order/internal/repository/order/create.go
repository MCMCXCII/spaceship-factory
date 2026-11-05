package order

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/MCMCXCII/spaceship-factory/order/internal/model"
	"github.com/MCMCXCII/spaceship-factory/order/internal/repository/converter"
)

func (r *repository) Create(ctx context.Context, order model.Order) error {
	return r.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := r.createOrder(txCtx, order); err != nil {
			return err
		}

		return r.createOrderItems(txCtx, order)
	})
}

func (r *repository) createOrder(ctx context.Context, order model.Order) error {
	rec := converter.OrderToRecord(order)

	const query = `
		INSERT INTO orders (uuid, status, created_at)
		VALUES ($1, $2, $3)
	`

	_, err := r.getter.DefaultTrOrDB(ctx, r.pool).Exec(
		ctx, query,
		rec.UUID,
		rec.Status,
		rec.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("создать заказ: %w", err)
	}

	return nil
}

func (r *repository) createOrderItems(ctx context.Context, order model.Order) error {
	if len(order.Items) == 0 {
		return nil
	}

	items := converter.OrderItemsToRecord(order.UUID, order.Items)

	builder := sq.Insert("order_items").
		Columns("order_uuid", "part_uuid", "part_type", "price").
		PlaceholderFormat(sq.Dollar)

	for _, item := range items {
		builder = builder.Values(item.OrderUUID, item.PartUUID, item.PartType, item.Price)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("собрать запррос деталей заказа: %w", err)
	}

	if _, err := r.getter.DefaultTrOrDB(ctx, r.pool).Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("создать детали заказа: %w", err)
	}

	return nil
}
