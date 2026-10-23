package v1

import (
	"context"

	"github.com/MCMCXCII/spaceship-factory/inventory/internal/api/converter"
	inventoryv1 "github.com/MCMCXCII/spaceship-factory/shared/pkg/proto/inventory/v1"
)

func (a *api) GetPart(ctx context.Context, req *inventoryv1.GetPartRequest) (*inventoryv1.GetPartResponse, error) {
	partUUID, err := converter.GetPartToInput(req)
	if err != nil {
		return nil, err
	}

	part, err := a.inventoryService.Get(ctx, partUUID)
	if err != nil {
		return nil, err
	}

	return &inventoryv1.GetPartResponse{Part: converter.PartToDTO(part)}, nil
}
