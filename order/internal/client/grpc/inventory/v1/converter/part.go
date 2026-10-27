package converter

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/MCMCXCII/spaceship-factory/order/internal/model"
	inventoryv1 "github.com/MCMCXCII/spaceship-factory/shared/pkg/proto/inventory/v1"
)

// UUIDsToDTO переводит список UUID деталей в строки для proto-запроса.
func UUIDsToDTO(uuids []uuid.UUID) []string {
	result := make([]string, 0, len(uuids))
	for _, u := range uuids {
		result = append(result, u.String())
	}
	return result
}

// PartsToModel переводит детали из ответа InventoryService в доменные модели.
func PartsToModel(parts []*inventoryv1.Part) ([]model.Part, error) {
	result := make([]model.Part, 0, len(parts))
	for _, p := range parts {
		part, err := PartToModel(p)
		if err != nil {
			return nil, err
		}
		result = append(result, part)
	}
	return result, nil
}

// PartToModel переводит деталь из ответа InventoryService в доменную модель.
func PartToModel(p *inventoryv1.Part) (model.Part, error) {
	partUUID, err := uuid.Parse(p.GetUuid())
	if err != nil {
		return model.Part{}, fmt.Errorf("разобрать uuid детали %q: %w", p.GetUuid(), err)
	}

	return model.Part{
		UUID:          partUUID,
		Name:          p.GetName(),
		Price:         p.GetPrice(),
		StockQuantity: p.GetStockQuantity(),
	}, nil
}
