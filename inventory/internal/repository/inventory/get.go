package inventory

import (
	"context"

	errs "github.com/MCMCXCII/spaceship-factory/inventory/internal/errors"
	"github.com/MCMCXCII/spaceship-factory/inventory/internal/model"
	"github.com/MCMCXCII/spaceship-factory/inventory/internal/repository/converter"
	"github.com/google/uuid"
)

func (r *Repository) Get(_ context.Context, partUUID uuid.UUID) (model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	rec, ok := r.data[partUUID]
	if !ok {
		return model.Part{}, errs.ErrPartNotFound
	}

	return converter.PartToModel(rec), nil
}
