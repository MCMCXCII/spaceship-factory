package v1

import (
	"context"

	"github.com/MCMCXCII/spaceship-factory/inventory/internal/api/converter"
	inventoryv1 "github.com/MCMCXCII/spaceship-factory/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(ctx context.Context, req *inventoryv1.ListPartsRequest) (*inventoryv1.ListPartsResponse, error) {
	partFilter, err := converter.ListPartsRequestToInput(req)
	if err != nil {
		return nil, err
	}

	parts, err := a.inventoryService.List(ctx, partFilter)
	if err != nil {
		return nil, err
	}

	return &inventoryv1.ListPartsResponse{Parts: converter.PartsToDTO(parts)}, nil
}
