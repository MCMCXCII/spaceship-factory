package record

import (
	"time"

	"github.com/google/uuid"
)

const (
	PartTypeHull   = "HULL"
	PartTypeEngine = "ENGINE"
	PartTypeShield = "SHIELD"
	PartTypeWeapon = "WEAPON"
)

type Part struct {
	UUID          uuid.UUID `db:"uuid"`
	Name          string    `db:"name"`
	Description   string    `db:"description"`
	Price         int64     `db:"price"`
	StockQuantity int64     `db:"stock_quantity"`
	PartType      string    `db:"part_type"`
	CreatedAt     time.Time `db:"created_at"`
}
