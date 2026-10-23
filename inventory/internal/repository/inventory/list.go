package inventory

import (
	"context"

	"github.com/MCMCXCII/spaceship-factory/inventory/internal/model"
	"github.com/MCMCXCII/spaceship-factory/inventory/internal/repository/converter"
)

func (r *Repository) List(_ context.Context) ([]model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	parts := make([]model.Part, 0, len(r.data))
	for _, rec := range r.data {
		parts = append(parts, converter.PartToModel(rec))
	}

	return parts, nil
}
