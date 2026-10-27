package inventory

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/MCMCXCII/spaceship-factory/inventory/internal/model"
	"github.com/MCMCXCII/spaceship-factory/inventory/internal/service/input"
)

func (s *service) List(ctx context.Context, partFilter input.PartFilter) ([]model.Part, error) {
	if len(partFilter.UUIDs) > 0 {
		parts := make([]model.Part, 0, len(partFilter.UUIDs))

		for _, partUUID := range partFilter.UUIDs {
			part, err := s.inventoryRepository.Get(ctx, partUUID)
			if err != nil {
				return nil, fmt.Errorf("inventoryRepository.Get %s: %w", partUUID, err)
			}
			parts = append(parts, part)
		}

		return parts, nil
	}

	all, err := s.inventoryRepository.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("inventoryRepository.List: %w", err)
	}

	parts := make([]model.Part, 0, len(all))
	for _, part := range all {
		if partFilter.PartType == model.PartTypeUnspecified || part.PartType == partFilter.PartType {
			parts = append(parts, part)
		}
	}

	slices.SortFunc(parts, func(a, b model.Part) int {
		return strings.Compare(a.Name, b.Name)
	})

	return parts, nil
}
