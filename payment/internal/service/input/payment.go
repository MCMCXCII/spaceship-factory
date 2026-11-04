package input

import (
	"github.com/google/uuid"

	"github.com/MCMCXCII/spaceship-factory/payment/internal/model"
)

type PayOrderInput struct {
	OrderUUID     uuid.UUID
	PaymentMethod model.PaymentMethod
}
