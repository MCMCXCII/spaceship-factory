package converter

import (
	"github.com/google/uuid"

	errs "github.com/MCMCXCII/spaceship-factory/order/internal/errors"
	"github.com/MCMCXCII/spaceship-factory/order/internal/model"
	"github.com/MCMCXCII/spaceship-factory/order/internal/service/input"
	orderv1 "github.com/MCMCXCII/spaceship-factory/shared/pkg/openapi/order/v1"
)

// CreateOrderRequestToInput переводит тело запроса на создание заказа во вход сервиса.
// Формат UUID и наличие обязательных полей уже проверил ogen.
func CreateOrderRequestToInput(req *orderv1.CreateOrderRequest) input.CreateOrderInput {
	return input.CreateOrderInput{
		HullUUID:   req.HullUUID,
		EngineUUID: req.EngineUUID,
		ShieldUUID: optNilUUIDToModel(req.ShieldUUID),
		WeaponUUID: optNilUUIDToModel(req.WeaponUUID),
	}
}

// PaymentMethodToModel переводит способ оплаты из HTTP API в доменный.
// Неизвестное значение — errs.ErrInvalidPaymentMethod.
func PaymentMethodToModel(m orderv1.PaymentMethod) (model.PaymentMethod, error) {
	switch m {
	case orderv1.PaymentMethodCARD:
		return model.PaymentMethodCard, nil
	case orderv1.PaymentMethodSBP:
		return model.PaymentMethodSBP, nil
	case orderv1.PaymentMethodCREDITCARD:
		return model.PaymentMethodCreditCard, nil
	case orderv1.PaymentMethodINVESTORMONEY:
		return model.PaymentMethodInvestorMoney, nil
	default:
		return model.PaymentMethodUnspecified, errs.ErrInvalidPaymentMethod
	}
}

// OrderToCreateResponse собирает ответ 201 на создание заказа.
func OrderToCreateResponse(o model.Order) *orderv1.CreateOrderResponse {
	return &orderv1.CreateOrderResponse{
		OrderUUID:  o.UUID,
		TotalPrice: o.TotalPrice,
	}
}

// OrderToDTO переводит доменный заказ в DTO ответа GET.
func OrderToDTO(o model.Order) *orderv1.OrderDto {
	return &orderv1.OrderDto{
		OrderUUID:       o.UUID,
		HullUUID:        o.HullUUID,
		EngineUUID:      o.EngineUUID,
		ShieldUUID:      nullUUIDToDTO(o.ShieldUUID),
		WeaponUUID:      nullUUIDToDTO(o.WeaponUUID),
		TotalPrice:      o.TotalPrice,
		TransactionUUID: nullUUIDToDTO(o.TransactionUUID),
		PaymentMethod:   paymentMethodToDTO(o.PaymentMethod),
		Status:          orderStatusToDTO(o.Status),
		CreatedAt:       o.CreatedAt,
	}
}

// optNilUUIDToModel переводит опциональный UUID из запроса в доменный NullUUID.
// Отсутствующее поле и null дают Valid == false.
func optNilUUIDToModel(v orderv1.OptNilUUID) uuid.NullUUID {
	id, ok := v.Get()
	return uuid.NullUUID{UUID: id, Valid: ok}
}

// nullUUIDToDTO переводит доменный NullUUID в опциональный UUID ответа.
// Valid == false превращается в null.
func nullUUIDToDTO(v uuid.NullUUID) orderv1.OptNilUUID {
	if !v.Valid {
		var out orderv1.OptNilUUID
		out.SetToNull()
		return out
	}
	return orderv1.NewOptNilUUID(v.UUID)
}

// paymentMethodToDTO переводит доменный способ оплаты в поле ответа.
// PaymentMethodUnspecified (заказ не оплачен) превращается в null.
func paymentMethodToDTO(m model.PaymentMethod) orderv1.OptNilPaymentMethod {
	var out orderv1.OptNilPaymentMethod

	switch m {
	case model.PaymentMethodCard:
		out.SetTo(orderv1.PaymentMethodCARD)
	case model.PaymentMethodSBP:
		out.SetTo(orderv1.PaymentMethodSBP)
	case model.PaymentMethodCreditCard:
		out.SetTo(orderv1.PaymentMethodCREDITCARD)
	case model.PaymentMethodInvestorMoney:
		out.SetTo(orderv1.PaymentMethodINVESTORMONEY)
	default:
		out.SetToNull()
	}

	return out
}

// orderStatusToDTO переводит доменный статус заказа в значение API.
func orderStatusToDTO(s model.OrderStatus) orderv1.OrderStatus {
	switch s {
	case model.OrderStatusPaid:
		return orderv1.OrderStatusPAID
	case model.OrderStatusCancelled:
		return orderv1.OrderStatusCANCELLED
	default:
		return orderv1.OrderStatusPENDINGPAYMENT
	}
}
