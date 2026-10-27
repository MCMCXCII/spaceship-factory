package inventory

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/MCMCXCII/spaceship-factory/inventory/internal/model"
)

func (s *service) Get(ctx context.Context, partUUID uuid.UUID) (model.Part, error) {
	part, err := s.inventoryRepository.Get(ctx, partUUID)
	if err != nil {
		return model.Part{}, fmt.Errorf("inventoryRepository.Get: %w", err)
	}

	return part, nil
}
