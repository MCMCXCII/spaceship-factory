package payment

import (
	"context"

	"github.com/google/uuid"

	"github.com/MCMCXCII/spaceship-factory/payment/internal/service/input"
)

// PaymentService — бизнес-логика оплаты, нужная API-слою
type PaymentService interface {
	Pay(ctx context.Context, in input.PayOrderInput) (uuid.UUID, error)
}
