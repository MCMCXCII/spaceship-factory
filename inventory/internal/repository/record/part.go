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

// Part — деталь в форме для хранилища
type Part struct {
	UUID          uuid.UUID
	Name          string
	Description   string
	Price         int64 // в копейках
	StockQuantity int64
	PartType      string
	CreatedAt     time.Time
}
