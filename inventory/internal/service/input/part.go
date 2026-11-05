package input

import (
	"github.com/google/uuid"

	"github.com/MCMCXCII/spaceship-factory/inventory/internal/model"
)

type PartFilter struct {
	UUIDs    []uuid.UUID
	PartType model.PartType
}
