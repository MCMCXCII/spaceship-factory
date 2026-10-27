package model

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	UUID            uuid.UUID
	HullUUID        uuid.UUID
	EngineUUID      uuid.UUID
	ShieldUUID      uuid.NullUUID
	WeaponUUID      uuid.NullUUID
	TotalPrice      int64
	TransactionUUID uuid.NullUUID
	PaymentMethod   PaymentMethod
	Status          OrderStatus
	CreatedAt       time.Time
}

// PaymentMethod - способ оплаты товара
type PaymentMethod int

const (
	// PaymentMethodUnspecified — способ оплаты не указан
	PaymentMethodUnspecified PaymentMethod = iota
	// PaymentMethodCard — банковская карта
	PaymentMethodCard
	// PaymentMethodSBP — система быстрых платежей
	PaymentMethodSBP
	// PaymentMethodCreditCard — кредитная карта
	PaymentMethodCreditCard
	// PaymentMethodInvestorMoney — деньги инвестора
	PaymentMethodInvestorMoney
)

// OrderStatus — статус заказа
type OrderStatus int

const (
	// OrderStatusUnspecified — заказ не создан
	OrderStatusUnspecified OrderStatus = iota
	// OrderStatusPendingPayment — заказ ожидает оплаты
	OrderStatusPendingPayment
	// OrderStatusPaid — заказ оплачен
	OrderStatusPaid
	// OrderStatusCancelled — заказ отменён
	OrderStatusCancelled
)
