package inventory

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	errs "github.com/MCMCXCII/spaceship-factory/inventory/internal/errors"
	"github.com/MCMCXCII/spaceship-factory/inventory/internal/model"
	"github.com/MCMCXCII/spaceship-factory/inventory/internal/repository/converter"
	"github.com/MCMCXCII/spaceship-factory/inventory/internal/repository/record"
)

func (r *repository) Get(ctx context.Context, partUUID uuid.UUID) (model.Part, error) {
	const query = `
		SELECT p.uuid, p.name, p.description, p.part_type, p.price, p.stock_quantity, p.created_at
		FROM parts p
		WHERE p.uuid = $1
	`
	row := r.pool.QueryRow(ctx, query, partUUID)

	var rec record.Part

	err := row.Scan(&rec.UUID, &rec.Name, &rec.Description, &rec.PartType, &rec.Price, &rec.StockQuantity, &rec.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Part{}, errs.ErrPartNotFound
		}

		return model.Part{}, fmt.Errorf("получить деталь: %w", err)
	}

	return converter.PartToModel(rec), nil
}
