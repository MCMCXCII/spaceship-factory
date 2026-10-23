package converter

import (
	"github.com/MCMCXCII/spaceship-factory/inventory/internal/model"
	"github.com/MCMCXCII/spaceship-factory/inventory/internal/repository/record"
)

// PartToRecord переводит доменную деталь в запись хранилища.
func PartToRecord(p model.Part) record.Part {
	return record.Part{
		UUID:          p.UUID,
		Name:          p.Name,
		Description:   p.Description,
		Price:         p.Price,
		StockQuantity: p.StockQuantity,
		PartType:      partTypeToRecord(p.PartType),
		CreatedAt:     p.CreatedAt,
	}
}

// PartToModel переводит запись хранилища в доменную деталь.
func PartToModel(r record.Part) model.Part {
	return model.Part{
		UUID:          r.UUID,
		Name:          r.Name,
		Description:   r.Description,
		Price:         r.Price,
		StockQuantity: r.StockQuantity,
		PartType:      partTypeToModel(r.PartType),
		CreatedAt:     r.CreatedAt,
	}
}

// partTypeToRecord переводит доменный тип детали в строку хранилища.
// Неизвестные значения превращаются в пустую строку.
func partTypeToRecord(t model.PartType) string {
	switch t {
	case model.PartTypeHull:
		return record.PartTypeHull
	case model.PartTypeEngine:
		return record.PartTypeEngine
	case model.PartTypeShield:
		return record.PartTypeShield
	case model.PartTypeWeapon:
		return record.PartTypeWeapon
	default:
		return ""
	}
}

// partTypeToModel переводит строку хранилища в доменный тип детали.
// Неизвестные значения превращаются в PartTypeUnspecified.
func partTypeToModel(s string) model.PartType {
	switch s {
	case record.PartTypeHull:
		return model.PartTypeHull
	case record.PartTypeEngine:
		return model.PartTypeEngine
	case record.PartTypeShield:
		return model.PartTypeShield
	case record.PartTypeWeapon:
		return model.PartTypeWeapon
	default:
		return model.PartTypeUnspecified
	}
}
