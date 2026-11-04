package input

import "github.com/google/uuid"

// CreateOrderInput — параметры создания заказа
type CreateOrderInput struct {
	HullUUID   uuid.UUID
	EngineUUID uuid.UUID
	ShieldUUID *uuid.UUID
	WeaponUUID *uuid.UUID
}

// PartUUIDs возвращает UUID всех выбранных деталей:
func (in CreateOrderInput) PartUUIDs() []uuid.UUID {
	uuids := make([]uuid.UUID, 0, 4)
	uuids = append(uuids, in.HullUUID, in.EngineUUID)

	if in.ShieldUUID != nil {
		uuids = append(uuids, *in.ShieldUUID)
	}

	if in.WeaponUUID != nil {
		uuids = append(uuids, *in.WeaponUUID)
	}

	return uuids
}
