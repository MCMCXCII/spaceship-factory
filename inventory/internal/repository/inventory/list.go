package inventory

import (
	"context"
	"fmt"

	"github.com/MCMCXCII/spaceship-factory/inventory/internal/model"
	"github.com/MCMCXCII/spaceship-factory/inventory/internal/repository/converter"
	"github.com/MCMCXCII/spaceship-factory/inventory/internal/repository/record"
)

func (r *repository) List(ctx context.Context) ([]model.Part, error) {
	const query = `
		SELECT uuid, name, description, part_type, price, stock_quantity, created_at
		FROM parts`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("получить детали: %w", err)
	}
	defer rows.Close()

	var parts []model.Part
	for rows.Next() {
		var rec record.Part
		if err := rows.Scan(
			&rec.UUID,
			&rec.Name,
			&rec.Description,
			&rec.PartType,
			&rec.Price,
			&rec.StockQuantity,
			&rec.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("прочитать деталь: %w", err)
		}
		parts = append(parts, converter.PartToModel(rec))
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("прочитать детали: %w", err)
	}

	return parts, nil
}
