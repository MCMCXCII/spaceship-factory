package model

import (
	"time"

	"github.com/google/uuid"
)

// Part — деталь космического корабля на складе
type Part struct {
	UUID          uuid.UUID
	Name          string
	Description   string
	Price         int64 // в копейках
	PartType      PartType
	StockQuantity int64
	CreatedAt     time.Time
}

// PartType — тип детали
type PartType string

const (
	// тип не указан; в фильтре означает «все детали»
	PartTypeUnspecified PartType = ""
	// корпус корабля
	PartTypeHull PartType = "HULL"
	// двигатель
	PartTypeEngine PartType = "ENGINE"
	// защитный щит
	PartTypeShield PartType = "SHIELD"
	// вооружение
	PartTypeWeapon PartType = "WEAPON"
)
