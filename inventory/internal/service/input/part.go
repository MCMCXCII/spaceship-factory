package input

import (
	"github.com/google/uuid"

	"github.com/MCMCXCII/spaceship-factory/inventory/internal/model"
)

// PartFilter — параметры выборки деталей.
// Если UUIDs не пуст, возвращаются именно эти детали (PartType игнорируется).
// Иначе фильтрация по PartType; PartTypeUnspecified — все детали.
type PartFilter struct {
	UUIDs    []uuid.UUID
	PartType model.PartType
}
