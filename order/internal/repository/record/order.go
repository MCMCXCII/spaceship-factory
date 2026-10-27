package record

import (
	"time"

	"github.com/google/uuid"
)

// Order — заказ в форме для хранилища
type Order struct {
	UUID            uuid.UUID
	HullUUID        uuid.UUID
	EngineUUID      uuid.UUID
	ShieldUUID      uuid.NullUUID
	WeaponUUID      uuid.NullUUID
	TotalPrice      int64
	TransactionUUID uuid.NullUUID
	PaymentMethod   string
	Status          string
	CreatedAt       time.Time
}

const (
	OrderStatusPendingPayment = "PENDING_PAYMENT"
	OrderStatusPaid           = "PAID"
	OrderStatusCancelled      = "CANCELLED"
)

const (
	PaymentMethodCard          = "CARD"
	PaymentMethodSBP           = "SBP"
	PaymentMethodCreditCard    = "CREDIT_CARD"
	PaymentMethodInvestorMoney = "INVESTOR_MONEY"
)
