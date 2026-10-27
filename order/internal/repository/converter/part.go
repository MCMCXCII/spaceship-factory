package converter

import (
	"github.com/MCMCXCII/spaceship-factory/order/internal/model"
	"github.com/MCMCXCII/spaceship-factory/order/internal/repository/record"
)

func OrderToRecord(o model.Order) record.Order {
	return record.Order{
		UUID:            o.UUID,
		HullUUID:        o.HullUUID,
		EngineUUID:      o.EngineUUID,
		ShieldUUID:      o.ShieldUUID,
		WeaponUUID:      o.WeaponUUID,
		TotalPrice:      o.TotalPrice,
		TransactionUUID: o.TransactionUUID,
		PaymentMethod:   paymentMethodToRecord(o.PaymentMethod),
		Status:          orderStatusToRecord(o.Status),
		CreatedAt:       o.CreatedAt,
	}
}

func OrderToModel(r record.Order) model.Order {
	return model.Order{
		UUID:            r.UUID,
		HullUUID:        r.HullUUID,
		EngineUUID:      r.EngineUUID,
		ShieldUUID:      r.ShieldUUID,
		WeaponUUID:      r.WeaponUUID,
		TotalPrice:      r.TotalPrice,
		TransactionUUID: r.TransactionUUID,
		PaymentMethod:   paymentMethodToModel(r.PaymentMethod),
		Status:          orderStatusToModel(r.Status),
		CreatedAt:       r.CreatedAt,
	}
}

func orderStatusToRecord(s model.OrderStatus) string {
	switch s {
	case model.OrderStatusPendingPayment:
		return record.OrderStatusPendingPayment
	case model.OrderStatusPaid:
		return record.OrderStatusPaid
	case model.OrderStatusCancelled:
		return record.OrderStatusCancelled
	default:
		return ""
	}
}

func orderStatusToModel(s string) model.OrderStatus {
	switch s {
	case record.OrderStatusPendingPayment:
		return model.OrderStatusPendingPayment
	case record.OrderStatusPaid:
		return model.OrderStatusPaid
	case record.OrderStatusCancelled:
		return model.OrderStatusCancelled
	default:
		return model.OrderStatusUnspecified
	}
}

func paymentMethodToRecord(m model.PaymentMethod) string {
	switch m {
	case model.PaymentMethodCard:
		return record.PaymentMethodCard
	case model.PaymentMethodSBP:
		return record.PaymentMethodSBP
	case model.PaymentMethodCreditCard:
		return record.PaymentMethodCreditCard
	case model.PaymentMethodInvestorMoney:
		return record.PaymentMethodInvestorMoney
	default:
		return ""
	}
}

func paymentMethodToModel(s string) model.PaymentMethod {
	switch s {
	case record.PaymentMethodCard:
		return model.PaymentMethodCard
	case record.PaymentMethodSBP:
		return model.PaymentMethodSBP
	case record.PaymentMethodCreditCard:
		return model.PaymentMethodCreditCard
	case record.PaymentMethodInvestorMoney:
		return model.PaymentMethodInvestorMoney
	default:
		return model.PaymentMethodUnspecified
	}
}
