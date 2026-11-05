package converter

import (
	"github.com/MCMCXCII/spaceship-factory/inventory/internal/model"
	"github.com/MCMCXCII/spaceship-factory/inventory/internal/repository/record"
)

func PartToRecord(p model.Part) record.Part {
	return record.Part{
		UUID:          p.UUID,
		Name:          p.Name,
		Description:   p.Description,
		Price:         p.Price,
		StockQuantity: p.StockQuantity,
		PartType:      string(p.PartType),
		CreatedAt:     p.CreatedAt,
	}
}

func PartToModel(r record.Part) model.Part {
	return model.Part{
		UUID:          r.UUID,
		Name:          r.Name,
		Description:   r.Description,
		Price:         r.Price,
		StockQuantity: r.StockQuantity,
		PartType:      model.PartType(r.PartType),
		CreatedAt:     r.CreatedAt,
	}
}
