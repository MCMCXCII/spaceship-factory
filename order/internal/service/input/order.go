package input

import "github.com/google/uuid"

// CreateOrderInput — параметры создания заказа
type CreateOrderInput struct {
	HullUUID   uuid.UUID
	EngineUUID uuid.UUID
	ShieldUUID uuid.NullUUID // Valid == false — щит не заказан
	WeaponUUID uuid.NullUUID // Valid == false — вооружение не заказано
}

// PartUUIDs возвращает UUID всех выбранных деталей:
func (in CreateOrderInput) PartUUIDs() []uuid.UUID {
	uuids := make([]uuid.UUID, 0, 4)
	uuids = append(uuids, in.HullUUID, in.EngineUUID)

	if in.ShieldUUID.Valid {
		uuids = append(uuids, in.ShieldUUID.UUID)
	}

	if in.WeaponUUID.Valid {
		uuids = append(uuids, in.WeaponUUID.UUID)
	}

	return uuids
}
