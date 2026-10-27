package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/MCMCXCII/spaceship-factory/inventory/internal/model"
	"github.com/MCMCXCII/spaceship-factory/inventory/internal/service/input"
)

type InventoryService interface {
	Get(context.Context, uuid.UUID) (model.Part, error)
	List(context.Context, input.PartFilter) ([]model.Part, error)
}
