package model

import (
	"time"

	"github.com/google/uuid"
)

// Part - деталь космического корабля
type Part struct {
	UUID          uuid.UUID
	Name          string
	Description   string
	Price         int64 // в копейках
	PartType      PartType
	StockQuantity int64
	CreatedAt     time.Time
}

// PartType - тип детали
type PartType int

const (
	// Не указан
	PartTypeUnspecified PartType = iota
	// Корпус корабля
	PartTypeHull
	// Двигатель
	PartTypeEngine
	// Защитный щит
	PartTypeShield
	// Вооружение
	PartTypeWeapon
)
