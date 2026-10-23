package converter

import (
	errs "github.com/MCMCXCII/spaceship-factory/inventory/internal/errors"
	"github.com/MCMCXCII/spaceship-factory/inventory/internal/model"
	"github.com/MCMCXCII/spaceship-factory/inventory/internal/service/input"
	inventoryv1 "github.com/MCMCXCII/spaceship-factory/shared/pkg/proto/inventory/v1"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ListPartsRequestToInput переводит запрос ListParts в фильтр сервиса
func ListPartsRequestToInput(req *inventoryv1.ListPartsRequest) (input.PartFilter, error) {
	uuids := make([]uuid.UUID, 0, len(req.Uuids))

	for _, u := range req.GetUuids() {
		partUUID, err := uuid.Parse(u)
		if err != nil {
			return input.PartFilter{}, errs.ErrInvalidUUID
		}
		uuids = append(uuids, partUUID)
	}

	return input.PartFilter{UUIDs: uuids, PartType: partTypeToModel(req.GetPartType())}, nil
}

// GetPartToInput извлекает UUID детали из запроса GetPart
func GetPartToInput(req *inventoryv1.GetPartRequest) (uuid.UUID, error) {
	partUUID, err := uuid.Parse(req.GetUuid())
	if err != nil {
		return uuid.Nil, errs.ErrInvalidUUID
	}

	return partUUID, nil
}

// PartToDTO переводит доменную деталь в proto-сообщение
func PartToDTO(p model.Part) *inventoryv1.Part {
	return &inventoryv1.Part{
		Uuid:          p.UUID.String(),
		Name:          p.Name,
		Description:   p.Description,
		Price:         p.Price,
		PartType:      partTypeToDTO(p.PartType),
		StockQuantity: p.StockQuantity,
		CreatedAt:     timestamppb.New(p.CreatedAt),
	}
}

// PartsToDTO переводит список доменных деталей в proto-сообщения
func PartsToDTO(parts []model.Part) []*inventoryv1.Part {
	result := make([]*inventoryv1.Part, 0, len(parts))
	for _, p := range parts {
		result = append(result, PartToDTO(p))
	}
	return result
}

// partTypeToModel переводит тип детали из proto-enum в доменный
func partTypeToModel(t inventoryv1.PartType) model.PartType {
	switch t {
	case inventoryv1.PartType_PART_TYPE_HULL:
		return model.PartTypeHull
	case inventoryv1.PartType_PART_TYPE_ENGINE:
		return model.PartTypeEngine
	case inventoryv1.PartType_PART_TYPE_SHIELD:
		return model.PartTypeShield
	case inventoryv1.PartType_PART_TYPE_WEAPON:
		return model.PartTypeWeapon
	default:
		return model.PartTypeUnspecified
	}
}

// partTypeToDTO переводит доменный тип детали в proto-enum
func partTypeToDTO(t model.PartType) inventoryv1.PartType {
	switch t {
	case model.PartTypeHull:
		return inventoryv1.PartType_PART_TYPE_HULL
	case model.PartTypeEngine:
		return inventoryv1.PartType_PART_TYPE_ENGINE
	case model.PartTypeShield:
		return inventoryv1.PartType_PART_TYPE_SHIELD
	case model.PartTypeWeapon:
		return inventoryv1.PartType_PART_TYPE_WEAPON
	default:
		return inventoryv1.PartType_PART_TYPE_UNSPECIFIED
	}
}
