package payment

import (
	"context"

	"github.com/google/uuid"

	"github.com/MCMCXCII/spaceship-factory/payment/internal/model"
)

// PaymentService — бизнес-логика оплаты, нужная API-слою
type PaymentService interface {
	Pay(ctx context.Context, orderUUID uuid.UUID, method model.PaymentMethod) (uuid.UUID, error)
}
